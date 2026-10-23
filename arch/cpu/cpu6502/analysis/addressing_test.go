package analysis

import (
	"testing"

	"github.com/retroenv/retrogolib/analysis/knownbits"
	"github.com/retroenv/retrogolib/arch/cpu/analysis"
	"github.com/retroenv/retrogolib/arch/cpu/cpu6502"
	"github.com/retroenv/retrogolib/assert"
)

func TestPointerWrapAndSingleReadModifyWrite(t *testing.T) {
	t.Parallel()

	input := fixture(0x03, 2, 1)
	input.Operand, _ = knownbits.NewConstant[uint16](8, 0xff)
	input.Before[X], _ = knownbits.NewConstant[uint8](8, 0)
	memory := &testMemory{}
	memory.Write(0xff, 0xff)
	memory.Write(0, 0xff)
	memory.Write(0xffff, 0x81)
	input.Read = stableRead(memory)
	result := bufferedResult()
	Transfer(input, &result)
	assert.True(t, result.EffectsComplete)
	assert.Len(t, result.Memory, 4)
	for index, address := range []uint16{0xff, 0, 0xffff, 0xffff} {
		actual, known := result.Memory[index].Address.Value.Constant(16)
		assert.True(t, known)
		assert.Equal(t, address, actual)
	}
	assert.Equal(t, analysis.MemoryRead, result.Memory[2].Kind)
	assert.Equal(t, analysis.MemoryWrite, result.Memory[3].Kind)
	assert.Equal(t, byte(2), result.Memory[3].Value.KnownOne)
	after := applyWrites(t, input.Before, result)
	assert.Equal(t, byte(3), after[A].KnownOne)
	input = fixture(0xb5, 2, 0)
	input.Operand, _ = knownbits.NewConstant[uint16](8, 0xff)
	input.Before[X], _ = knownbits.NewConstant[uint8](8, 1)
	Transfer(input, &result)
	address, known := result.Memory[0].Address.Value.Constant(16)
	assert.True(t, known)
	assert.Equal(t, uint16(0), address)
}

func TestIndirectJumpVariantRule(t *testing.T) {
	t.Parallel()

	for _, variant := range []cpu6502.CPUVariant{cpu6502.VariantNES6502, cpu6502.Variant65C02} {
		input := fixture(0x6c, 3, 1)
		input.Variant = variant
		input.Operand, _ = knownbits.NewConstant[uint16](16, 0x12ff)
		result := bufferedResult()
		Transfer(input, &result)
		assert.Len(t, result.Memory, 2)
		assert.Equal(t, analysis.FlowJump, result.Flow)
		address, known := result.Memory[1].Address.Value.Constant(16)
		assert.True(t, known)
		expected := uint16(0x1200)
		if variant == cpu6502.Variant65C02 {
			expected = 0x1300
		}
		assert.Equal(t, expected, address)
	}
}

func TestCallAndStackWrap(t *testing.T) {
	t.Parallel()

	input := fixture(0x20, 3, 1)
	input.PC, _ = knownbits.NewConstant[uint16](16, 0xffff)
	input.Before[SP], _ = knownbits.NewConstant[uint8](8, 0)
	result := bufferedResult()
	Transfer(input, &result)
	assert.Equal(t, analysis.FlowCall, result.Flow)
	assert.True(t, result.EffectsComplete)
	assert.Len(t, result.Memory, 2)
	for index, address := range []uint16{0x100, 0x1ff} {
		actual, known := result.Memory[index].Address.Value.Constant(16)
		assert.True(t, known)
		assert.Equal(t, address, actual)
		assert.Equal(t, analysis.MemoryWrite, result.Memory[index].Kind)
	}
	after := applyWrites(t, input.Before, result)
	assert.Equal(t, byte(0xfe), after[SP].KnownOne)
	assert.Equal(t, byte(0), result.Memory[0].Value.KnownOne)
	assert.Equal(t, byte(1), result.Memory[1].Value.KnownOne)
	input = fixture(0x68, 1, 1)
	input.Before[SP], _ = knownbits.NewConstant[uint8](8, 0xff)
	Transfer(input, &result)
	address, known := result.Memory[0].Address.Value.Constant(16)
	assert.True(t, known)
	assert.Equal(t, uint16(0x100), address)
}
