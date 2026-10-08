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
		{name: "extended indirect with Y bits", postbyte: 0xBF, want: 5},
		{name: "extended indirect with U bits", postbyte: 0xDF, want: 5},
		{name: "extended indirect with S bits", postbyte: 0xFF, want: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, indexedModeCycles(tt.postbyte))
		})
	}
}

func TestExtendedIndirectRegisterBitsAreIgnored(t *testing.T) {
	// The register bits of an extended indirect postbyte do not change the decode or the cycles.
	for _, postbyte := range []uint8{0x9F, 0xBF, 0xDF, 0xFF} {
		cpu, mem := newTestCPU(t)
		mem.data[0x8000] = 0xA6 // LDA indexed
		mem.data[0x8001] = postbyte
		mem.WriteWord(0x8002, 0x1000)
		mem.WriteWord(0x1000, 0x2000)
		mem.data[0x2000] = 0x42

		assert.NoError(t, cpu.Step())
		assert.Equal(t, uint8(0x42), cpu.A, "postbyte 0x%02X", postbyte)
		assert.Equal(t, uint16(0x8004), cpu.PC, "postbyte 0x%02X", postbyte)
		assert.Equal(t, uint64(9), cpu.Cycles(), "postbyte 0x%02X", postbyte)
	}
}
