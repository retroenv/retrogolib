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
	var err error

	s, err = conv.Absolute(addressing.Absolute(0x1000))
	assert.NoError(t, err)
	assert.Equal(t, "a:$1000", s)

	s, err = conv.AbsoluteX(addressing.Absolute(0x1000))
	assert.NoError(t, err)
	assert.Equal(t, "a:$1000,X", s)

	s, err = conv.AbsoluteY(addressing.Absolute(0x1000))
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

	s, err = conv.ZeroPage(addressing.ZeroPage(0x10))
	assert.NoError(t, err)
	assert.Equal(t, "<$10", s)

	s, err = conv.ZeroPageX(addressing.ZeroPage(0x10))
	assert.NoError(t, err)
	assert.Equal(t, "<$10,X", s)

	s, err = conv.ZeroPageY(addressing.ZeroPage(0x10))
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

	s, err = conv.Indirect(addressing.Indirect(0x1000))
	assert.NoError(t, err)
	assert.Equal(t, "[$1000]", s)

	s, err = conv.IndirectX(addressing.Indirect(0x1000))
	assert.NoError(t, err)
	assert.Equal(t, "[$1000,X]", s)

	s, err = conv.IndirectY(addressing.Indirect(0x1000))
	assert.NoError(t, err)
	assert.Equal(t, "[$1000],Y", s)
}
