package parameter

import (
	"testing"

	"github.com/retroenv/retrogolib/arch/cpu/m6502"
	"github.com/retroenv/retrogolib/assert"
)

func TestParameterAbsolute(t *testing.T) {
	cfg := Config{
		AbsolutePrefix: "a:",
	}
	conv := New(cfg)
	var s string
	var err error

	s, err = conv.Absolute(m6502.Absolute(0x1000))
	assert.NoError(t, err)
	assert.Equal(t, "a:$1000", s)

	s, err = conv.AbsoluteX(m6502.Absolute(0x1000))
	assert.NoError(t, err)
	assert.Equal(t, "a:$1000,X", s)

	s, err = conv.AbsoluteY(m6502.Absolute(0x1000))
	assert.NoError(t, err)
	assert.Equal(t, "a:$1000,Y", s)
}

func TestParameterZeroPage(t *testing.T) {
	cfg := Config{
		ZeroPagePrefix: "<",
	}
	conv := New(cfg)
	var s string
	var err error

	s, err = conv.ZeroPage(m6502.ZeroPage(0x10))
	assert.NoError(t, err)
	assert.Equal(t, "<$10", s)

	s, err = conv.ZeroPageX(m6502.ZeroPage(0x10))
	assert.NoError(t, err)
	assert.Equal(t, "<$10,X", s)

	s, err = conv.ZeroPageY(m6502.ZeroPage(0x10))
	assert.NoError(t, err)
	assert.Equal(t, "<$10,Y", s)
}

func TestParameterIndirect(t *testing.T) {
	cfg := Config{
		IndirectPrefix: "[",
		IndirectSuffix: "]",
	}
	conv := New(cfg)
	var s string
	var err error

	s, err = conv.Indirect(m6502.Indirect(0x1000))
	assert.NoError(t, err)
	assert.Equal(t, "[$1000]", s)

	s, err = conv.IndirectX(m6502.Indirect(0x1000))
	assert.NoError(t, err)
	assert.Equal(t, "[$1000,X]", s)

	s, err = conv.IndirectY(m6502.Indirect(0x1000))
	assert.NoError(t, err)
	assert.Equal(t, "[$1000],Y", s)
}

func TestParameterImmediate(t *testing.T) {
	t.Parallel()
	conv := New(Config{})

	s := conv.Immediate(0x42)
	assert.Equal(t, "#$42", s)

	s = conv.Immediate(0xFF)
	assert.Equal(t, "#$FF", s)

	s = conv.Immediate(0x00)
	assert.Equal(t, "#$00", s)
}

func TestParameterAccumulator(t *testing.T) {
	t.Parallel()
	conv := New(Config{})

	s := conv.Accumulator()
	assert.Equal(t, "a", s)
}

func TestParameterRelative(t *testing.T) {
	t.Parallel()
	conv := New(Config{})

	s := conv.Relative(0x1000)
	assert.Equal(t, "$1000", s)

	s = conv.Relative(nil)
	assert.Equal(t, "", s)
}

func TestParameterStringValues(t *testing.T) {
	t.Parallel()
	cfg := Config{
		AbsolutePrefix: "abs:",
		ZeroPagePrefix: "zp:",
		IndirectPrefix: "(",
		IndirectSuffix: ")",
	}
	conv := New(cfg)

	// Test string parameters
	s, err := conv.Absolute("LABEL")
	assert.NoError(t, err)
	assert.Equal(t, "abs:LABEL", s)

	s, err = conv.ZeroPage("ZP_LABEL")
	assert.NoError(t, err)
	assert.Equal(t, "zp:ZP_LABEL", s)

	s, err = conv.Indirect("IND_LABEL")
	assert.NoError(t, err)
	assert.Equal(t, "(IND_LABEL)", s)
}

func TestParameterErrorCases(t *testing.T) {
	t.Parallel()
	conv := New(Config{})

	// Test unsupported types
	_, err := conv.Absolute([]int{1, 2, 3})
	assert.ErrorContains(t, err, "unsupported param type")

	_, err = conv.ZeroPage(map[string]int{"key": 1})
	assert.ErrorContains(t, err, "unsupported param type")

	_, err = conv.Indirect(42.5)
	assert.ErrorContains(t, err, "unsupported param type")
}

func TestConverterConfiguration(t *testing.T) {
	t.Parallel()

	// Test different configurations
	cfg1 := Config{
		AbsolutePrefix: "ABS:",
		ZeroPagePrefix: "ZP:",
		IndirectPrefix: "<",
		IndirectSuffix: ">",
	}
	conv1 := New(cfg1)

	s, err := conv1.Absolute(0x2000)
	assert.NoError(t, err)
	assert.Equal(t, "ABS:$2000", s)

	s, err = conv1.ZeroPage(0x80)
	assert.NoError(t, err)
	assert.Equal(t, "ZP:$80", s)

	s, err = conv1.Indirect(m6502.Indirect(0x3000))
	assert.NoError(t, err)
	assert.Equal(t, "<$3000>", s)

	// Test empty configuration
	cfg2 := Config{}
	conv2 := New(cfg2)

	s, err = conv2.Absolute(0x4000)
	assert.NoError(t, err)
	assert.Equal(t, "$4000", s)
}

func TestParameterIntTypes(t *testing.T) {
	t.Parallel()
	conv := New(Config{})

	// Test different addressing mode types
	s, err := conv.AbsoluteX(m6502.AbsoluteX(0x5000))
	assert.NoError(t, err)
	assert.Equal(t, "$5000,X", s)

	s, err = conv.AbsoluteY(m6502.AbsoluteY(0x6000))
	assert.NoError(t, err)
	assert.Equal(t, "$6000,Y", s)

	s, err = conv.ZeroPageX(m6502.ZeroPageX(0x70))
	assert.NoError(t, err)
	assert.Equal(t, "$70,X", s)

	s, err = conv.ZeroPageY(m6502.ZeroPageY(0x80))
	assert.NoError(t, err)
	assert.Equal(t, "$80,Y", s)

	s, err = conv.IndirectX(m6502.IndirectX(0x7000))
	assert.NoError(t, err)
	assert.Equal(t, "$7000,X", s)

	s, err = conv.IndirectY(m6502.IndirectY(0x8000))
	assert.NoError(t, err)
	assert.Equal(t, "$8000,Y", s)
}
