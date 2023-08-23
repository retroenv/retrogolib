package parameter

import (
	"testing"

	"github.com/retroenv/retrogolib/addressing"
	"github.com/retroenv/retrogolib/assert"
)

func TestParameterAbsolute(t *testing.T) {
	cfg := Config{
		AbsolutePrefix: "a:",
	}
	conv := New(cfg)
	var s string

	s = conv.Absolute(addressing.Absolute(0x1000))
	assert.Equal(t, "a:$1000", s)

	s = conv.AbsoluteX(addressing.Absolute(0x1000))
	assert.Equal(t, "a:$1000,X", s)

	s = conv.AbsoluteY(addressing.Absolute(0x1000))
	assert.Equal(t, "a:$1000,Y", s)
}

func TestParameterZeroPage(t *testing.T) {
	cfg := Config{
		ZeroPagePrefix: "<",
	}
	conv := New(cfg)
	var s string

	s = conv.ZeroPage(addressing.ZeroPage(0x10))
	assert.Equal(t, "<$10", s)

	s = conv.ZeroPageX(addressing.ZeroPage(0x10))
	assert.Equal(t, "<$10,X", s)

	s = conv.ZeroPageY(addressing.ZeroPage(0x10))
	assert.Equal(t, "<$10,Y", s)
}

func TestParameterIndirect(t *testing.T) {
	cfg := Config{
		IndirectPrefix: "[",
		IndirectSuffix: "]",
	}
	conv := New(cfg)
	var s string

	s = conv.Indirect(addressing.Indirect(0x1000))
	assert.Equal(t, "[$1000]", s)

	s = conv.IndirectX(addressing.Indirect(0x1000))
	assert.Equal(t, "[$1000,X]", s)

	s = conv.IndirectY(addressing.Indirect(0x1000))
	assert.Equal(t, "[$1000],Y", s)
}
