package analysis

import (
	"testing"

	"github.com/retroenv/retrogolib/analysis/knownbits"
	"github.com/retroenv/retrogolib/arch/cpu/analysis"
	"github.com/retroenv/retrogolib/arch/cpu/z80"
	"github.com/retroenv/retrogolib/assert"
)

func TestPairAndFlagViews(t *testing.T) {
	t.Parallel()

	upper, valid := RegisterView(z80.RegB)
	assert.True(t, valid)
	lower, valid := RegisterView(z80.RegC)
	assert.True(t, valid)
	whole, valid := RegisterView(z80.RegBC)
	assert.True(t, valid)
	assert.False(t, upper.Overlaps(lower))
	assert.True(t, upper.Overlaps(whole))
	assert.True(t, lower.Overlaps(whole))
	accumulator, valid := RegisterView(z80.RegA)
	assert.True(t, valid)
	assert.False(t, accumulator.Overlaps(view(AF, 0, 8)))
	assert.True(t, view(AF, 0, 8).Overlaps(view(AF, 0, 1)))
	assert.False(t, view(ShadowAF, 0, 16).Overlaps(view(AF, 0, 16)))
	_, valid = RegisterView(z80.RegNone)
	assert.False(t, valid)
	assert.Equal(t, knownbits.Width(0), StorageWidth(ResourceCount))
}

func TestByteWritesPreserveSiblingAndCopySnapshot(t *testing.T) {
	t.Parallel()

	input := fixture(0, 0x06, 0)
	input.Before[BC], _ = knownbits.NewConstant[uint16](16, 0x9934)
	input.Operand, _ = knownbits.NewConstant[uint16](8, 0x12)
	result := bufferedResult()
	Transfer(input, &result)
	after := applyWrites(t, input.Before, result)
	assert.Equal(t, uint16(0x1234), after[BC].KnownOne)
	assert.Equal(t, input.Before[AF], after[AF])
	input.Opcode.Opcode = 0x41
	input.Opcode.Size = 1
	Transfer(input, &result)
	after = applyWrites(t, input.Before, result)
	assert.Equal(t, uint16(0x3434), after[BC].KnownOne)
	assert.Len(t, result.Copies, 1)
	assert.Equal(t, view(BC, 0, 8), result.Copies[0].Source)
	assert.Equal(t, view(BC, 8, 8), result.Copies[0].Destination)
	assert.True(t, result.Valid(StorageWidth))
}

func TestIndexedWriteAndSignedWrap(t *testing.T) {
	t.Parallel()

	input := fixture(z80.PrefixDD, 0x66, 0)
	input.Before[HL], _ = knownbits.NewConstant[uint16](16, 0x1234)
	input.Before[IX], _ = knownbits.NewConstant[uint16](16, 0)
	input.Operand, _ = knownbits.NewConstant[uint16](8, 0xff)
	input.Read = func(Address) (knownbits.Value[uint16], bool) {
		return knownbits.NewConstant[uint16](8, 0xaa)
	}
	result := bufferedResult()
	Transfer(input, &result)
	after := applyWrites(t, input.Before, result)
	assert.Equal(t, uint16(0xaa34), after[HL].KnownOne)
	assert.Equal(t, input.Before[IX], after[IX])
	assert.Len(t, result.Memory, 1)
	address, known := result.Memory[0].Address.Value.Constant(16)
	assert.True(t, known)
	assert.Equal(t, uint16(0xffff), address)
	input = fixture(z80.PrefixDD, 0x26, 0)
	input.Before[IX], _ = knownbits.NewConstant[uint16](16, 0x1234)
	input.Operand, _ = knownbits.NewConstant[uint16](8, 0xaa)
	Transfer(input, &result)
	after = applyWrites(t, input.Before, result)
	assert.Equal(t, uint16(0xaa34), after[IX].KnownOne)
	input = fixture(z80.PrefixDD, 0x36, 0)
	input.Before[IX], _ = knownbits.NewConstant[uint16](16, 0)
	input.Operand, _ = knownbits.NewConstant[uint16](16, 0xaaff)
	Transfer(input, &result)
	assert.Len(t, result.Memory, 1)
	assert.Equal(t, analysis.MemoryWrite, result.Memory[0].Kind)
	assert.Equal(t, uint16(0xaa), result.Memory[0].Value.KnownOne)
}

func TestLittleEndianWordStoreWrap(t *testing.T) {
	t.Parallel()

	input := fixture(0, 0x22, 0)
	input.Operand, _ = knownbits.NewConstant[uint16](16, 0xffff)
	input.Before[HL], _ = knownbits.NewConstant[uint16](16, 0x1234)
	result := bufferedResult()
	Transfer(input, &result)
	assert.Len(t, result.Memory, 2)
	for index, number := range []uint16{0xffff, 0} {
		address, known := result.Memory[index].Address.Value.Constant(16)
		assert.True(t, known)
		assert.Equal(t, number, address)
	}
	assert.Equal(t, uint16(0x34), result.Memory[0].Value.KnownOne)
	assert.Equal(t, uint16(0x12), result.Memory[1].Value.KnownOne)
}
