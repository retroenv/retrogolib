package config

import (
	"io"
	"slices"
	"strings"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestCaseSensitiveNames(t *testing.T) {
	input := "[Video]\nPort=1\nport=2\n[video]\nPort=3"
	_, err := Parse(strings.NewReader(input), Options{})
	assert.ErrorIs(t, err, ErrDuplicateKey)
	cfg, err := Parse(strings.NewReader(input), Options{CaseSensitive: true})
	assert.NoError(t, err)
	entries := slices.Collect(cfg.Entries())
	assert.Equal(t, []SectionInfo{{Name: "Video", Line: 1}, {Name: "video", Line: 4}}, slices.Collect(cfg.Sections()))
	assert.Equal(t, "Port", entries[0].Key)
	assert.Equal(t, "port", entries[1].Key)
	assert.Equal(t, "video", entries[2].Section)
	assert.Equal(t, 5, entries[2].Line)

	var values struct {
		Upper   int `config:"Video.Port,required"`
		Lower   int `config:"Video.port,required"`
		Section int `config:"video.Port,required"`
	}
	assert.NoError(t, cfg.Unmarshal(&values))
	assert.Equal(t, 1, values.Upper)
	assert.Equal(t, 2, values.Lower)
	assert.Equal(t, 3, values.Section)
	values.Upper = 4
	assert.NoError(t, cfg.Marshal(values))
	data, err := cfg.SaveBytes()
	assert.NoError(t, err)
	loaded, err := Parse(strings.NewReader(string(data)), Options{CaseSensitive: true})
	assert.NoError(t, err)
	assert.NoError(t, loaded.Unmarshal(&values))
	assert.Equal(t, 4, values.Upper)
	assert.Equal(t, 2, values.Lower)
}

func TestCaseSensitiveNestedAndAutomaticNames(t *testing.T) {
	cfg, err := Parse(strings.NewReader("Root=7\n[Video]\nPort=8\n[Video.Nested]\nValue=9"), Options{CaseSensitive: true})
	assert.NoError(t, err)
	var values struct {
		Root  int
		Video struct {
			Port   int `config:"Port,required"`
			Nested struct{ Value int }
		}
	}
	assert.NoError(t, cfg.Unmarshal(&values))
	assert.Equal(t, 7, values.Root)
	assert.Equal(t, 8, values.Video.Port)
	assert.Equal(t, 9, values.Video.Nested.Value)
	values.Video.Nested.Value = 10
	assert.NoError(t, cfg.Marshal(values))
	assert.NoError(t, cfg.Unmarshal(&values))
	assert.Equal(t, 10, values.Video.Nested.Value)
}

func TestCaseInsensitiveDefaults(t *testing.T) {
	cfg, err := Parse(strings.NewReader("[Video]\nPort=1"), Options{})
	assert.NoError(t, err)
	assert.Equal(t, "video", slices.Collect(cfg.Sections())[0].Name)
	assert.Equal(t, "port", slices.Collect(cfg.Entries())[0].Key)
	var values struct {
		Port int `config:"VIDEO.PORT,required"`
	}
	assert.NoError(t, cfg.Unmarshal(&values))
	assert.Equal(t, 1, values.Port)
	_, err = Parse(strings.NewReader("[Video]\n[video]"), Options{})
	assert.ErrorIs(t, err, ErrDuplicateSection)
	_, err = Parse(strings.NewReader("[Video]\nPort=1\nPort=2"), Options{CaseSensitive: true})
	assert.ErrorIs(t, err, ErrDuplicateKey)
}

func TestRawINIOptions(t *testing.T) {
	opts := Options{CaseSensitive: true, RawValues: true, CommentPrefixes: ";#", InlineComments: true,
		LiteralSections: []string{"comments"}, AllowRepeatedSections: true}
	input := "\uFEFF; header\r\n[Symbols] ; names\r\nPort = $2000 ; register\n" +
		"Quote = \"literal; # text\" ; trailing\n[comments]\n$8000=Keep; # and = and \\n\n" +
		"[Symbols]\nExpression=0x10+Port\nEscape=\"not an escape\\q\""
	cfg, err := Parse(strings.NewReader(input), opts)
	assert.NoError(t, err)
	entries := slices.Collect(cfg.Entries())
	assert.Len(t, entries, 5)
	assert.Equal(t, "$2000", entries[0].Value.Raw)
	assert.Equal(t, `"literal; # text"`, entries[1].Value.Raw)
	assert.Equal(t, `Keep; # and = and \n`, entries[2].Value.Raw)
	assert.Equal(t, "0x10+Port", entries[3].Value.Raw)
	assert.Equal(t, `"not an escape\q"`, entries[4].Value.Raw)
	assert.Len(t, slices.Collect(cfg.Sections()), 3)
	data, err := cfg.SaveBytes()
	assert.NoError(t, err)
	loaded, err := Parse(strings.NewReader(string(data)), opts)
	assert.True(t, strings.Contains(string(data), "$2000 ; register"))
	assert.NoError(t, err)
	for i, entry := range slices.Collect(loaded.Entries()) {
		assert.Equal(t, entries[i].Value.Raw, entry.Value.Raw)
		assert.Equal(t, entries[i].Key, entry.Key)
	}
}

func TestInlineCommentQuotes(t *testing.T) {
	cfg, err := Parse(strings.NewReader("Text = \"quoted \\\"; # retained\" ; discarded"),
		Options{
			CommentPrefixes: ";#",
			InlineComments:  true,
		})
	assert.NoError(t, err)
	assert.Equal(t, `quoted "; # retained`, slices.Collect(cfg.Entries())[0].Value.Raw)
	for _, input := range []string{"; unsupported by default", "[A] ; unsupported by default"} {
		_, err := Parse(strings.NewReader(input), Options{})
		assert.Error(t, err)
	}
	_, err = Parse(strings.NewReader("[A]\nName=1\n[A]\nname=2"), Options{AllowRepeatedSections: true})
	assert.ErrorIs(t, err, ErrDuplicateKey)
}

func TestOrderedIteratorsStop(t *testing.T) {
	cfg, err := Parse(strings.NewReader("[B]\nZ=1\n[A]\nY=2"), Options{CaseSensitive: true})
	assert.NoError(t, err)
	for entry := range cfg.Entries() {
		assert.Equal(t, "Z", entry.Key)
		break
	}
	for section := range cfg.Sections() {
		assert.Equal(t, "B", section.Name)
		break
	}
}

func TestLiteralSectionsFollowCaseOption(t *testing.T) {
	for _, sensitive := range []bool{false, true} {
		cfg, err := Parse(strings.NewReader("[Comments]\nText=one; two"), Options{
			CaseSensitive: sensitive, InlineComments: true, CommentPrefixes: ";",
			LiteralSections: []string{"comments"},
		})
		assert.NoError(t, err)
		want := "one; two"
		if sensitive {
			want = "one"
		}
		assert.Equal(t, want, slices.Collect(cfg.Entries())[0].Value.Raw)
	}
}

func TestParseReaderLimitsAndErrors(t *testing.T) {
	_, err := Parse(strings.NewReader(strings.Repeat(" ", maxConfigSize+1)), Options{})
	assert.ErrorIs(t, err, ErrConfigTooLarge)
	_, err = Parse(strings.NewReader(strings.Repeat("\n", maxLines)), Options{})
	assert.ErrorIs(t, err, ErrTooManyLines)
	_, err = Parse(failedReader{}, Options{})
	assert.ErrorIs(t, err, io.ErrUnexpectedEOF)
}

type failedReader struct{}

func (failedReader) Read(_ []byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}
