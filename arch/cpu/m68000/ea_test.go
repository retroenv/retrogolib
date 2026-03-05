package m68000

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestDecodeEA_DataRegDirect(t *testing.T) {
	cpu := newTestCPU(t)
	cpu.D[3] = 0x12345678

	ea, err := cpu.decodeEA(0, 3, SizeLong)
	assert.NoError(t, err)
	assert.Equal(t, uint8(0), ea.Mode)
	assert.Equal(t, uint8(3), ea.Reg)

	val, err := cpu.readEA(ea)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x12345678), val)
}

func TestDecodeEA_DataRegDirect_Byte(t *testing.T) {
	cpu := newTestCPU(t)
	cpu.D[0] = 0x12345678

	ea, err := cpu.decodeEA(0, 0, SizeByte)
	assert.NoError(t, err)
	val, err := cpu.readEA(ea)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x78), val)
}

func TestDecodeEA_AddrRegDirect(t *testing.T) {
	cpu := newTestCPU(t)
	cpu.A[2] = 0x00ABCDEF

	ea, err := cpu.decodeEA(1, 2, SizeLong)
	assert.NoError(t, err)
	val, err := cpu.readEA(ea)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x00ABCDEF), val)
}

func TestDecodeEA_AddrRegIndirect(t *testing.T) {
	cpu := newTestCPU(t)
	cpu.A[0] = 0x2000
	cpu.bus.WriteWord(0x2000, 0x1234)

	ea, err := cpu.decodeEA(2, 0, SizeWord)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x2000), ea.Address)

	val, err := cpu.readEA(ea)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x1234), val)
}

func TestDecodeEA_PostIncrement(t *testing.T) {
	cpu := newTestCPU(t)
	cpu.A[1] = 0x3000
	cpu.bus.WriteWord(0x3000, 0x5678)

	ea, err := cpu.decodeEA(3, 1, SizeWord)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x3000), ea.Address)
	assert.Equal(t, uint32(0x3002), cpu.A[1]) // Incremented by word size

	val, err := cpu.readEA(ea)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x5678), val)
}

func TestDecodeEA_PostIncrement_A7Byte(t *testing.T) {
	cpu := newTestCPU(t)
	origSP := cpu.sp
	cpu.bus.Write(origSP, 0x42)

	ea, err := cpu.decodeEA(3, 7, SizeByte)
	assert.NoError(t, err)
	assert.Equal(t, origSP, ea.Address)
	// A7 byte operations increment by 2 to maintain alignment.
	assert.Equal(t, origSP+2, cpu.sp)
}

func TestDecodeEA_PreDecrement(t *testing.T) {
	cpu := newTestCPU(t)
	cpu.A[2] = 0x4002
	cpu.bus.WriteWord(0x4000, 0xABCD)

	ea, err := cpu.decodeEA(4, 2, SizeWord)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x4000), ea.Address)
	assert.Equal(t, uint32(0x4000), cpu.A[2])

	val, err := cpu.readEA(ea)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0xABCD), val)
}

func TestDecodeEA_Displacement(t *testing.T) {
	cpu := newTestCPU(t)
	cpu.A[0] = 0x5000
	// Write displacement word in instruction stream.
	cpu.bus.WriteWord(cpu.PC, 0x0010) // d16 = +16
	cpu.bus.WriteWord(0x5010, 0x9999)

	ea, err := cpu.decodeEA(5, 0, SizeWord)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x5010), ea.Address)

	val, err := cpu.readEA(ea)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x9999), val)
}

func TestDecodeEA_Indexed(t *testing.T) {
	cpu := newTestCPU(t)
	cpu.A[0] = 0x6000
	cpu.D[1] = 0x0010
	// Extension word: D1.W, displacement = 4.
	// Format: D/A=0 | Reg=001 | W/L=0 | 000 | disp=00000100
	cpu.bus.WriteWord(cpu.PC, 0x1004)
	cpu.bus.WriteWord(0x6014, 0x7777)

	ea, err := cpu.decodeEA(6, 0, SizeWord)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x6014), ea.Address)

	val, err := cpu.readEA(ea)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x7777), val)
}

func TestDecodeEA_AbsShort(t *testing.T) {
	cpu := newTestCPU(t)
	cpu.bus.WriteWord(cpu.PC, 0x2000) // Absolute short address
	cpu.bus.WriteWord(0x2000, 0x1111)

	ea, err := cpu.decodeEA(7, 0, SizeWord)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x2000), ea.Address)

	val, err := cpu.readEA(ea)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x1111), val)
}

func TestDecodeEA_AbsLong(t *testing.T) {
	cpu := newTestCPU(t)
	cpu.bus.WriteLong(cpu.PC, 0x00123456) // Absolute long address
	cpu.bus.WriteWord(0x123456, 0x2222)

	ea, err := cpu.decodeEA(7, 1, SizeWord)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x123456), ea.Address)

	val, err := cpu.readEA(ea)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x2222), val)
}

func TestDecodeEA_PCDisplacement(t *testing.T) {
	cpu := newTestCPU(t)
	pcBefore := cpu.PC
	cpu.bus.WriteWord(cpu.PC, 0x0020) // d16 = +32
	cpu.bus.WriteWord(pcBefore+0x0020, 0x3333)

	ea, err := cpu.decodeEA(7, 2, SizeWord)
	assert.NoError(t, err)
	assert.Equal(t, pcBefore+0x0020, ea.Address)

	val, err := cpu.readEA(ea)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x3333), val)
}

func TestDecodeEA_Immediate(t *testing.T) {
	cpu := newTestCPU(t)
	cpu.bus.WriteWord(cpu.PC, 0x00FF) // Immediate byte value

	ea, err := cpu.decodeEA(7, 4, SizeByte)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0xFF), ea.Value)
}

func TestDecodeEA_Immediate_Word(t *testing.T) {
	cpu := newTestCPU(t)
	cpu.bus.WriteWord(cpu.PC, 0x1234)

	ea, err := cpu.decodeEA(7, 4, SizeWord)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x1234), ea.Value)
}

func TestDecodeEA_Immediate_Long(t *testing.T) {
	cpu := newTestCPU(t)
	cpu.bus.WriteLong(cpu.PC, 0x12345678)

	ea, err := cpu.decodeEA(7, 4, SizeLong)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x12345678), ea.Value)
}

func TestWriteEA_DataReg(t *testing.T) {
	cpu := newTestCPU(t)
	cpu.D[0] = 0xFFFFFFFF

	ea := EffectiveAddress{Mode: 0, Reg: 0, Size: SizeByte}
	err := cpu.writeEA(ea, 0x42)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0xFFFFFF42), cpu.D[0]) // Preserves upper bits
}

func TestWriteEA_Memory(t *testing.T) {
	cpu := newTestCPU(t)

	ea := EffectiveAddress{Mode: 2, Reg: 0, Size: SizeWord, Address: 0x2000}
	err := cpu.writeEA(ea, 0x1234)
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x1234), cpu.bus.ReadWord(0x2000))
}
