package parameter

import (
	"testing"

	"github.com/retroenv/retrogolib/addressing"
	"github.com/retroenv/retrogolib/assert"
)

func TestParameterDefault(t *testing.T) {
	cfg := Config{}
	conv := New(cfg)

	s := conv.IndirectX(addressing.Indirect(0x1000))
	assert.Equal(t, "($1000,X)", s)

	s = conv.IndirectY(addressing.Indirect(0x1000))
	assert.Equal(t, "($1000),Y", s)
}

func TestParameterIndirectNoParentheses(t *testing.T) {
	cfg := Config{
		IndirectNoParentheses: true,
	}
	conv := New(cfg)

	s := conv.IndirectX(addressing.Indirect(0x1000))
	assert.Equal(t, "$1000,X", s)

	s = conv.IndirectY(addressing.Indirect(0x1000))
	assert.Equal(t, "$1000,Y", s)
}
