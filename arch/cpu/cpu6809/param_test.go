package cpu6809

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestInvalidIndirectAutoUpdateModes(t *testing.T) {
	tests := []struct {
		name     string
		postbyte uint8
	}{
		{name: "postincrement by one", postbyte: 0x90},
		{name: "predecrement by one", postbyte: 0x92},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu, mem := newTestCPU(t)
			cpu.X = 0x1000
			mem.data[0x8000] = 0xA6 // LDA indexed
			mem.data[0x8001] = tt.postbyte

			err := cpu.Step()
			assert.ErrorIs(t, err, ErrInvalidIndexPostbyte)
			assert.Equal(t, uint16(0x1000), cpu.X)
		})
	}
}

func TestParameterizedHandlerRejectsWrongOperandType(t *testing.T) {
	cpu, _ := newTestCPU(t)

	err := BccInst.paramFunc(cpu, Immediate8(0))
	assert.ErrorIs(t, err, ErrInvalidParameterType)
}

func TestIndexedModeCycles(t *testing.T) {
	tests := []struct {
		name     string
		postbyte uint8
		want     uint64
	}{
		{name: "five-bit offset", postbyte: 0x00, want: 1},
		{name: "no offset", postbyte: 0x84},
		{name: "postincrement", postbyte: 0x80, want: 2},
		{name: "16-bit offset", postbyte: 0x89, want: 4},
		{name: "indirect no offset", postbyte: 0x94, want: 3},
		{name: "indirect 16-bit offset", postbyte: 0x99, want: 7},
		{name: "extended indirect", postbyte: 0x9F, want: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, indexedModeCycles(tt.postbyte))
		})
	}
}
