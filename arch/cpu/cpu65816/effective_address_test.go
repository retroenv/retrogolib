package cpu65816

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestResolveEA(t *testing.T) {
	t.Parallel()

	cpu := &CPU{DB: 0x12, DP: 0x1000, SP: 0x2000}
	tests := []struct {
		name  string
		param any
		want  uint32
	}{
		{name: "direct page", param: DirectPage(0x34), want: 0x001034},
		{name: "direct page X", param: DirectPageX(0x34), want: 0x001034},
		{name: "absolute", param: Absolute16(0x3456), want: 0x123456},
		{name: "resolved indirect", param: DPIndirect(0x234567), want: 0x234567},
		{name: "resolved long", param: AbsLong(0x345678), want: 0x345678},
		{name: "stack relative", param: StackRel(0x34), want: 0x002034},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := cpu.resolveEA(test.param)
			assert.NoError(t, err)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestResolveEARejectsNonAddressOperands(t *testing.T) {
	t.Parallel()

	cpu := &CPU{}
	for _, param := range []any{Immediate8(1), Immediate16(1), "invalid"} {
		_, err := cpu.resolveEA(param)
		assert.ErrorIs(t, err, ErrUnsupportedAddressingMode)
	}
}
