package cli

import (
	"errors"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestFlagSet_Parse(t *testing.T) {
	type params struct {
		Input  string `flag:"i" usage:"input file"`
		Output string `flag:"o" usage:"output file"`
	}
	type options struct {
		Debug bool `flag:"debug" usage:"enable debug mode"`
		Count int  `flag:"n" usage:"repeat count" default:"1"`
	}

	var p params
	var o options

	fs := NewFlagSet("test")
	fs.AddSection("Parameters", &p)
	fs.AddSection("Options", &o)

	args, err := fs.Parse([]string{"-i", "input.txt", "-debug", "-n", "5", "extra"})
	assert.NoError(t, err)
	assert.Equal(t, "input.txt", p.Input)
	assert.Equal(t, "", p.Output)
	assert.True(t, o.Debug)
	assert.Equal(t, 5, o.Count)
	assert.Equal(t, []string{"extra"}, args)
}

func TestFlagSet_Defaults(t *testing.T) {
	type options struct {
		Format string `flag:"f" usage:"output format" default:"json"`
		Count  int    `flag:"n" usage:"count" default:"10"`
	}

	var o options
	fs := NewFlagSet("test")
	fs.AddSection("Options", &o)

	_, err := fs.Parse([]string{})
	assert.NoError(t, err)
	assert.Equal(t, "json", o.Format)
	assert.Equal(t, 10, o.Count)
}

func TestFlagSet_SkipFields(t *testing.T) {
	type options struct {
		Internal string // No flag tag - should be skipped
		Visible  string `flag:"v" usage:"visible flag"`
		Skipped  string `flag:"-" usage:"explicitly skipped"`
	}

	var o options
	fs := NewFlagSet("test")
	fs.AddSection("Options", &o)

	assert.Len(t, fs.sections, 1)
	assert.Len(t, fs.sections[0].Flags, 1)
	assert.Equal(t, "v", fs.sections[0].Flags[0].Name)
}

func TestFlagSet_NumericTypes(t *testing.T) {
	type options struct {
		Int     int     `flag:"int" usage:"int value" default:"1"`
		Int64   int64   `flag:"int64" usage:"int64 value" default:"2"`
		Uint    uint    `flag:"uint" usage:"uint value" default:"3"`
		Uint64  uint64  `flag:"uint64" usage:"uint64 value" default:"4"`
		Float64 float64 `flag:"float64" usage:"float64 value" default:"1.5"`
	}

	t.Run("defaults", func(t *testing.T) {
		var o options
		fs := NewFlagSet("test")
		fs.AddSection("Options", &o)

		_, err := fs.Parse([]string{})
		assert.NoError(t, err)
		assert.Equal(t, 1, o.Int)
		assert.Equal(t, int64(2), o.Int64)
		assert.Equal(t, uint(3), o.Uint)
		assert.Equal(t, uint64(4), o.Uint64)
		assert.Equal(t, 1.5, o.Float64)
	})

	t.Run("parsed", func(t *testing.T) {
		var o options
		fs := NewFlagSet("test")
		fs.AddSection("Options", &o)

		_, err := fs.Parse([]string{"-int", "10", "-int64", "20", "-uint", "30", "-uint64", "40", "-float64", "2.5"})
		assert.NoError(t, err)
		assert.Equal(t, 10, o.Int)
		assert.Equal(t, int64(20), o.Int64)
		assert.Equal(t, uint(30), o.Uint)
		assert.Equal(t, uint64(40), o.Uint64)
		assert.Equal(t, 2.5, o.Float64)
	})
}

func TestFlagSet_UnsupportedType(t *testing.T) {
	type options struct {
		Data []string `flag:"data" usage:"data values"`
	}

	var o options
	fs := NewFlagSet("test")
	fs.AddSection("Options", &o)

	assert.Len(t, fs.sections, 1)
	assert.Empty(t, fs.sections[0].Flags)
}

func TestFlagSet_Required(t *testing.T) {
	type options struct {
		Input  string `flag:"i" usage:"input file" required:"true"`
		Output string `flag:"o" usage:"output file"`
	}

	t.Run("missing", func(t *testing.T) {
		var o options
		fs := NewFlagSet("test")
		fs.AddSection("Options", &o)

		_, err := fs.Parse([]string{"-o", "out.txt"})
		assert.Error(t, err)

		var missingErr *MissingFlagsError
		assert.True(t, errors.As(err, &missingErr))
		assert.Equal(t, []string{"i"}, missingErr.Flags)
	})

	t.Run("provided", func(t *testing.T) {
		var o options
		fs := NewFlagSet("test")
		fs.AddSection("Options", &o)

		_, err := fs.Parse([]string{"-i", "in.txt"})
		assert.NoError(t, err)
		assert.Equal(t, "in.txt", o.Input)
	})
}

func TestFlagSet_EnvVars(t *testing.T) {
	type options struct {
		Host string `flag:"host" usage:"server host" env:"TEST_CLI_HOST" default:"localhost"`
		Port int    `flag:"port" usage:"server port" env:"TEST_CLI_PORT" default:"8080"`
	}

	t.Run("from_env", func(t *testing.T) {
		t.Setenv("TEST_CLI_HOST", "example.com")
		t.Setenv("TEST_CLI_PORT", "9000")

		var o options
		fs := NewFlagSet("test")
		fs.AddSection("Options", &o)

		_, err := fs.Parse([]string{})
		assert.NoError(t, err)
		assert.Equal(t, "example.com", o.Host)
		assert.Equal(t, 9000, o.Port)
	})

	t.Run("flag_overrides_env", func(t *testing.T) {
		t.Setenv("TEST_CLI_HOST", "from-env.com")
		t.Setenv("TEST_CLI_PORT", "9000")

		var o options
		fs := NewFlagSet("test")
		fs.AddSection("Options", &o)

		_, err := fs.Parse([]string{"-host", "from-flag.com"})
		assert.NoError(t, err)
		assert.Equal(t, "from-flag.com", o.Host)
	})
}

func TestFlagSet_ShortLongFlags(t *testing.T) {
	type options struct {
		Verbose bool   `flag:"v,verbose" usage:"enable verbose output"`
		Output  string `flag:"o,output" usage:"output file" default:"out.txt"`
		Count   int    `flag:"n,count" usage:"repeat count" default:"1"`
	}

	t.Run("short", func(t *testing.T) {
		var o options
		fs := NewFlagSet("test")
		fs.AddSection("Options", &o)

		_, err := fs.Parse([]string{"-v", "-o", "short.txt", "-n", "5"})
		assert.NoError(t, err)
		assert.True(t, o.Verbose)
		assert.Equal(t, "short.txt", o.Output)
		assert.Equal(t, 5, o.Count)
	})

	t.Run("long", func(t *testing.T) {
		var o options
		fs := NewFlagSet("test")
		fs.AddSection("Options", &o)

		_, err := fs.Parse([]string{"-verbose", "-output", "long.txt"})
		assert.NoError(t, err)
		assert.True(t, o.Verbose)
		assert.Equal(t, "long.txt", o.Output)
	})
}

func TestFlagSet_UsageSections(t *testing.T) {
	type params struct {
		Input string `flag:"i" usage:"input file"`
	}
	type options struct {
		Verbose bool `flag:"v" usage:"verbose output"`
	}

	var p params
	var o options
	fs := NewFlagSet("test")
	fs.AddSection("Parameters", &p)
	fs.AddSection("Options", &o)

	output := captureStdout(func() { fs.ShowUsage() })
	assert.Contains(t, output, "usage: test [options]")
	assert.Contains(t, output, "Parameters:")
	assert.Contains(t, output, "-i string")
	assert.Contains(t, output, "Options:")
	assert.Contains(t, output, "-v")
}

func TestFlagSet_UsageTypesAndDefaults(t *testing.T) {
	type options struct {
		Name    string  `flag:"name" usage:"name value" default:"test"`
		Count   int     `flag:"count" usage:"count value" default:"5"`
		Ratio   float64 `flag:"ratio" usage:"ratio value" default:"1.5"`
		Verbose bool    `flag:"verbose" usage:"enable verbose"`
	}

	var o options
	fs := NewFlagSet("test")
	fs.AddSection("Options", &o)

	output := captureStdout(func() { fs.ShowUsage() })
	assert.Contains(t, output, "-name string")
	assert.Contains(t, output, "-count int")
	assert.Contains(t, output, "-ratio float64")
	assert.Contains(t, output, "-verbose\n")
	assert.Contains(t, output, "(default: test)")
	assert.Contains(t, output, "(default: 5)")
	assert.Contains(t, output, "(default: 1.5)")
}

func TestFlagSet_UsageEnvAndRequired(t *testing.T) {
	type options struct {
		Input string `flag:"i" usage:"input file" required:"true"`
		Host  string `flag:"host" usage:"server host" env:"HOST" default:"localhost"`
	}

	var o options
	fs := NewFlagSet("test")
	fs.AddSection("Options", &o)

	output := captureStdout(func() { fs.ShowUsage() })
	assert.Contains(t, output, "(required)")
	assert.Contains(t, output, "[env: HOST]")
}

func TestFlagSet_UsageShortLong(t *testing.T) {
	type options struct {
		Verbose bool   `flag:"v,verbose" usage:"enable verbose output"`
		Output  string `flag:"o,output" usage:"output file"`
	}

	var o options
	fs := NewFlagSet("test")
	fs.AddSection("Options", &o)

	output := captureStdout(func() { fs.ShowUsage() })
	assert.Contains(t, output, "-v, -verbose")
	assert.Contains(t, output, "-o, -output")
}

func TestFlagSet_PositionalBasic(t *testing.T) {
	type positional struct {
		Input  string `arg:"positional" usage:"input file" required:"true"`
		Output string `arg:"positional" usage:"output file"`
	}

	var p positional
	fs := NewFlagSet("test")
	fs.AddPositional(&p)

	remaining, err := fs.Parse([]string{"in.txt", "out.txt", "extra"})
	assert.NoError(t, err)
	assert.Equal(t, "in.txt", p.Input)
	assert.Equal(t, "out.txt", p.Output)
	assert.Equal(t, []string{"extra"}, remaining)
}

func TestFlagSet_PositionalRequired(t *testing.T) {
	type positional struct {
		Input string `arg:"positional" usage:"input file" required:"true"`
	}

	var p positional
	fs := NewFlagSet("test")
	fs.AddPositional(&p)

	_, err := fs.Parse([]string{})
	assert.Error(t, err)

	var missingErr *MissingArgsError
	assert.True(t, errors.As(err, &missingErr))
	assert.Equal(t, []string{"input"}, missingErr.Args)
}

func TestFlagSet_PositionalVariadic(t *testing.T) {
	type positional struct {
		Files []string `arg:"positional" usage:"files to process"`
	}

	var p positional
	fs := NewFlagSet("test")
	fs.AddPositional(&p)

	remaining, err := fs.Parse([]string{"a.txt", "b.txt", "c.txt"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"a.txt", "b.txt", "c.txt"}, p.Files)
	assert.Empty(t, remaining)
}

func TestFlagSet_PositionalMixed(t *testing.T) {
	type positional struct {
		Input string   `arg:"positional" usage:"input file" required:"true"`
		Files []string `arg:"positional" usage:"additional files"`
	}

	var p positional
	fs := NewFlagSet("test")
	fs.AddPositional(&p)

	remaining, err := fs.Parse([]string{"in.txt", "a.txt", "b.txt"})
	assert.NoError(t, err)
	assert.Equal(t, "in.txt", p.Input)
	assert.Equal(t, []string{"a.txt", "b.txt"}, p.Files)
	assert.Empty(t, remaining)
}

func TestFlagSet_PositionalWithFlags(t *testing.T) {
	type options struct {
		Verbose bool `flag:"v" usage:"verbose output"`
	}
	type positional struct {
		Input string `arg:"positional" usage:"input file" required:"true"`
	}

	var o options
	var p positional
	fs := NewFlagSet("test")
	fs.AddSection("Options", &o)
	fs.AddPositional(&p)

	remaining, err := fs.Parse([]string{"-v", "in.txt"})
	assert.NoError(t, err)
	assert.True(t, o.Verbose)
	assert.Equal(t, "in.txt", p.Input)
	assert.Empty(t, remaining)
}

func TestFlagSet_PositionalUsage(t *testing.T) {
	type positional struct {
		Input string   `arg:"positional" usage:"input file" required:"true"`
		Files []string `arg:"positional" usage:"additional files"`
	}

	var p positional
	fs := NewFlagSet("test")
	fs.AddPositional(&p)

	output := captureStdout(func() { fs.ShowUsage() })
	assert.Contains(t, output, "usage: test input [files...]")
	assert.Contains(t, output, "Positional arguments:")
	assert.Contains(t, output, "input file")
	assert.Contains(t, output, "(required)")
	assert.Contains(t, output, "files...")
}

func TestFlagSet_EmptyFlagName(t *testing.T) {
	type options struct {
		Value string `flag:"" usage:"empty flag name"`
	}

	var o options
	fs := NewFlagSet("test")
	fs.AddSection("Options", &o)

	// Empty flag name should be skipped
	assert.Len(t, fs.sections, 1)
	assert.Empty(t, fs.sections[0].Flags)
}

func TestFlagSet_ParseError(t *testing.T) {
	type options struct {
		Count int `flag:"n" usage:"count"`
	}

	var o options
	fs := NewFlagSet("test")
	fs.AddSection("Options", &o)

	// Invalid integer should cause parse error
	_, err := fs.Parse([]string{"-n", "invalid"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "parsing flags")
}

func TestFlagSet_HelpRequested(t *testing.T) {
	type options struct {
		Verbose bool `flag:"v" usage:"verbose output"`
	}

	t.Run("short_help", func(t *testing.T) {
		var o options
		fs := NewFlagSet("test")
		fs.AddSection("Options", &o)

		_, err := fs.Parse([]string{"-h"})
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrHelpRequested))
	})

	t.Run("long_help", func(t *testing.T) {
		var o options
		fs := NewFlagSet("test")
		fs.AddSection("Options", &o)

		_, err := fs.Parse([]string{"-help"})
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrHelpRequested))
	})
}

func TestFlagSet_MultipleRequiredFlags(t *testing.T) {
	type options struct {
		Input  string `flag:"i" usage:"input file" required:"true"`
		Output string `flag:"o" usage:"output file" required:"true"`
		Format string `flag:"f" usage:"format" required:"true"`
	}

	var o options
	fs := NewFlagSet("test")
	fs.AddSection("Options", &o)

	// Missing all required flags
	_, err := fs.Parse([]string{})
	assert.Error(t, err)

	var missingErr *MissingFlagsError
	assert.True(t, errors.As(err, &missingErr))
	assert.Equal(t, []string{"i", "o", "f"}, missingErr.Flags)
}

func TestFlagSet_BoolDefaultTrue(t *testing.T) {
	type options struct {
		Enabled bool `flag:"enabled" usage:"enable feature" default:"true"`
	}

	var o options
	fs := NewFlagSet("test")
	fs.AddSection("Options", &o)

	_, err := fs.Parse([]string{})
	assert.NoError(t, err)
	assert.True(t, o.Enabled)
}

func TestFlagSet_ZeroValueDefaults(t *testing.T) {
	type options struct {
		Count   int     `flag:"n" usage:"count" default:"0"`
		Ratio   float64 `flag:"r" usage:"ratio" default:"0.0"`
		Enabled bool    `flag:"e" usage:"enabled" default:"false"`
	}

	var o options
	fs := NewFlagSet("test")
	fs.AddSection("Options", &o)

	_, err := fs.Parse([]string{})
	assert.NoError(t, err)
	assert.Equal(t, 0, o.Count)
	assert.Equal(t, 0.0, o.Ratio)
	assert.False(t, o.Enabled)
}

func TestFlagSet_PositionalEmptyVariadic(t *testing.T) {
	type positional struct {
		Files []string `arg:"positional" usage:"files to process"`
	}

	var p positional
	fs := NewFlagSet("test")
	fs.AddPositional(&p)

	remaining, err := fs.Parse([]string{})
	assert.NoError(t, err)
	// Empty args result in empty slice, not nil
	assert.Empty(t, p.Files)
	assert.Empty(t, remaining)
}

func TestFlagSet_MixedPositionalAndFlags(t *testing.T) {
	type options struct {
		Verbose bool `flag:"v" usage:"verbose"`
		Debug   bool `flag:"d" usage:"debug"`
	}
	type positional struct {
		Command string   `arg:"positional" usage:"command" required:"true"`
		Args    []string `arg:"positional" usage:"command arguments"`
	}

	var o options
	var p positional
	fs := NewFlagSet("test")
	fs.AddSection("Options", &o)
	fs.AddPositional(&p)

	remaining, err := fs.Parse([]string{"-v", "-d", "build", "--flag1", "val1"})
	assert.NoError(t, err)
	assert.True(t, o.Verbose)
	assert.True(t, o.Debug)
	assert.Equal(t, "build", p.Command)
	assert.Equal(t, []string{"--flag1", "val1"}, p.Args)
	assert.Empty(t, remaining)
}

func TestFlagSet_StructuredErrors(t *testing.T) {
	t.Run("MissingFlagsError", func(t *testing.T) {
		type options struct {
			Input string `flag:"i" usage:"input" required:"true"`
		}

		var o options
		fs := NewFlagSet("test")
		fs.AddSection("Options", &o)

		_, err := fs.Parse([]string{})
		assert.Error(t, err)

		var missingErr *MissingFlagsError
		assert.True(t, errors.As(err, &missingErr))
		assert.Equal(t, []string{"i"}, missingErr.Flags)
	})

	t.Run("MissingArgsError", func(t *testing.T) {
		type positional struct {
			File string `arg:"positional" usage:"file" required:"true"`
		}

		var p positional
		fs := NewFlagSet("test")
		fs.AddPositional(&p)

		_, err := fs.Parse([]string{})
		assert.Error(t, err)

		var missingErr *MissingArgsError
		assert.True(t, errors.As(err, &missingErr))
		assert.Equal(t, []string{"file"}, missingErr.Args)
	})
}
