package codedatalog

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/retroenv/retrogolib/arch/system/nes/cartridge"
	"github.com/retroenv/retrogolib/assert"
)

type failedWriter struct {
	calls  int
	failAt int
	short  bool
}

func (w *failedWriter) Write(data []byte) (int, error) {
	w.calls++
	if w.calls != w.failAt {
		return len(data), nil
	}
	if w.short {
		return 0, nil
	}
	return 0, io.ErrClosedPipe
}

func TestWrite(t *testing.T) {
	prg := []PrgFlag{Code | Data, SubEntryPoint}
	chr := []ChrFlag{DrawnOnScreen | ReadProgrammatically}
	var output bytes.Buffer

	assert.NoError(t, Write(&output, prg, chr))
	assert.Equal(t, []byte{0x03, 0x00, 0x03}, output.Bytes())

	cart := cartridge.New()
	cart.PRG = cart.PRG[:len(prg)]
	cart.CHR = cart.CHR[:len(chr)]
	loaded, err := LoadFile(cart, bytes.NewReader(output.Bytes()))
	assert.NoError(t, err)
	assert.Equal(t, []PrgFlag{Code | Data, 0}, loaded)
	assert.Equal(t, []PrgFlag{Code | Data, SubEntryPoint}, prg)
}

func TestWriteEmpty(t *testing.T) {
	var output bytes.Buffer
	assert.NoError(t, Write(&output, nil, nil))
	assert.Empty(t, output.Bytes())
}

func TestWriteCHROnly(t *testing.T) {
	var output bytes.Buffer
	assert.NoError(t, Write(&output, nil, []ChrFlag{ReadProgrammatically}))
	assert.Equal(t, []byte{0x02}, output.Bytes())
}

func TestWriteFailures(t *testing.T) {
	for _, failAt := range []int{1, 2} {
		for _, short := range []bool{false, true} {
			t.Run(fmt.Sprintf("write_%d_short_%t", failAt, short), func(t *testing.T) {
				writer := &failedWriter{
					failAt: failAt,
					short:  short,
				}
				err := Write(writer, []PrgFlag{Code}, []ChrFlag{DrawnOnScreen})
				want := io.ErrClosedPipe
				if short {
					want = io.ErrShortWrite
				}
				assert.ErrorIs(t, err, want)
				assert.Equal(t, failAt, writer.calls)
				space := "PRG"
				if failAt == 2 {
					space = "CHR"
				}
				assert.ErrorContains(t, err, "writing "+space+" flags")
			})
		}
	}
}

func TestPrgFlagNames(t *testing.T) {
	assert.Empty(t, PrgFlag(0).Names())
	assert.Equal(t, []string{"code", "indirect_data"}, (Code | IndirectData).Names())
	assert.Equal(t, []string{
		"code", "data", "rom_bank_mapped_low", "rom_bank_mapped_high",
		"indirect_code", "indirect_data", "pcm_audio", "sub_entry_point",
	}, PrgFlag(0xff).Names())

	names := Code.Names()
	names[0] = "changed"
	assert.Equal(t, []string{"code"}, Code.Names())
}

func TestChrFlagNames(t *testing.T) {
	assert.Empty(t, ChrFlag(0).Names())
	assert.Empty(t, ChrFlag(0x80).Names())
	assert.Equal(t, []string{"read_programmatically"}, ReadProgrammatically.Names())
	assert.Equal(t, []string{"drawn_on_screen", "read_programmatically"}, ChrFlag(0xff).Names())

	names := DrawnOnScreen.Names()
	names[0] = "changed"
	assert.Equal(t, []string{"drawn_on_screen"}, DrawnOnScreen.Names())
}

func TestWriteClearsUndefinedBits(t *testing.T) {
	var output bytes.Buffer
	assert.NoError(t, Write(&output, []PrgFlag{0xff}, []ChrFlag{0xff}))
	assert.Equal(t, []byte{0x7f, 0x03}, output.Bytes())
}

func TestPrgFlagWithMapping(t *testing.T) {
	flags := Code | Data | SubEntryPoint
	for index, address := range []uint16{0x8000, 0xa000, 0xc000, 0xe000} {
		want := flags | PrgFlag(index<<2)
		assert.Equal(t, want, (flags | RomBankMappedLow | RomBankMappedHigh).WithMapping(address))
	}
}

func TestLoadFileSizes(t *testing.T) {
	cart := cartridge.New()
	for _, size := range []int{0, len(cart.PRG) - 1, len(cart.PRG) + 1, len(cart.PRG) + len(cart.CHR) + 1} {
		_, err := LoadFile(cart, bytes.NewReader(make([]byte, size)))
		assert.Error(t, err)
	}
	for _, size := range []int{len(cart.PRG), len(cart.PRG) + len(cart.CHR)} {
		flags, err := LoadFile(cart, bytes.NewReader(make([]byte, size)))
		assert.NoError(t, err)
		assert.Len(t, flags, len(cart.PRG))
	}
}

func TestLoadFileRejectsMesen2(t *testing.T) {
	cart := cartridge.New()
	data := append([]byte("CDLv2"), make([]byte, 4+len(cart.PRG)+len(cart.CHR))...)
	_, err := LoadFile(cart, bytes.NewReader(data))
	assert.ErrorIs(t, err, ErrUnsupportedFormat)
}

func TestLoadFileAllowsRawSignatureBytes(t *testing.T) {
	cart := cartridge.New()
	cart.PRG = cart.PRG[:5]
	flags, err := LoadFile(cart, bytes.NewBufferString("CDLv2"))
	assert.NoError(t, err)
	assert.Len(t, flags, 5)
	assert.Equal(t, PrgFlag('C'), flags[0])
}
