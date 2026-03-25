package m6809

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestLDA_Immediate(t *testing.T) {
	cpu, mem := newTestCPU(t)
	mem.data[0x8000] = 0x86 // LDA #$42
	mem.data[0x8001] = 0x42
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(0x42), cpu.A)
	assert.Equal(t, uint16(0x8002), cpu.PC)
}

func TestLDB_Immediate(t *testing.T) {
	cpu, mem := newTestCPU(t)
	mem.data[0x8000] = 0xC6 // LDB #$FF
	mem.data[0x8001] = 0xFF
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(0xFF), cpu.B)
	assert.Equal(t, uint8(1), cpu.Flags.N) // negative
}

func TestLDD_Immediate(t *testing.T) {
	cpu, mem := newTestCPU(t)
	mem.data[0x8000] = 0xCC // LDD #$1234
	mem.data[0x8001] = 0x12
	mem.data[0x8002] = 0x34
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x1234), cpu.D())
}

func TestADDA_Immediate(t *testing.T) {
	cpu, mem := newTestCPU(t)
	cpu.A = 0x10
	mem.data[0x8000] = 0x8B // ADDA #$20
	mem.data[0x8001] = 0x20
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(0x30), cpu.A)
}

func TestSUBA_Immediate(t *testing.T) {
	cpu, mem := newTestCPU(t)
	cpu.A = 0x30
	mem.data[0x8000] = 0x80 // SUBA #$10
	mem.data[0x8001] = 0x10
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(0x20), cpu.A)
}

func TestANDA_Immediate(t *testing.T) {
	cpu, mem := newTestCPU(t)
	cpu.A = 0xFF
	mem.data[0x8000] = 0x84 // ANDA #$0F
	mem.data[0x8001] = 0x0F
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(0x0F), cpu.A)
}

func TestNOP(t *testing.T) {
	cpu, mem := newTestCPU(t)
	mem.data[0x8000] = 0x12 // NOP
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x8001), cpu.PC)
}

func TestCLRA(t *testing.T) {
	cpu, mem := newTestCPU(t)
	cpu.A = 0xFF
	mem.data[0x8000] = 0x4F // CLRA
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(0), cpu.A)
	assert.Equal(t, uint8(1), cpu.Flags.Z)
	assert.Equal(t, uint8(0), cpu.Flags.N)
}

func TestINCA(t *testing.T) {
	cpu, mem := newTestCPU(t)
	cpu.A = 0x7F
	mem.data[0x8000] = 0x4C // INCA
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(0x80), cpu.A)
	assert.Equal(t, uint8(1), cpu.Flags.V) // overflow: 0x7F -> 0x80
	assert.Equal(t, uint8(1), cpu.Flags.N) // negative
}

func TestDECA(t *testing.T) {
	cpu, mem := newTestCPU(t)
	cpu.A = 0x01
	mem.data[0x8000] = 0x4A // DECA
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(0x00), cpu.A)
	assert.Equal(t, uint8(1), cpu.Flags.Z)
}

func TestBRA(t *testing.T) {
	cpu, mem := newTestCPU(t)
	mem.data[0x8000] = 0x20 // BRA +5
	mem.data[0x8001] = 0x05
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x8007), cpu.PC) // 0x8000 + 2 + 5
}

func TestBEQ_Taken(t *testing.T) {
	cpu, mem := newTestCPU(t)
	cpu.Flags.Z = 1
	mem.data[0x8000] = 0x27 // BEQ +3
	mem.data[0x8001] = 0x03
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x8005), cpu.PC) // 0x8000 + 2 + 3
}

func TestBEQ_NotTaken(t *testing.T) {
	cpu, mem := newTestCPU(t)
	cpu.Flags.Z = 0
	mem.data[0x8000] = 0x27 // BEQ +3
	mem.data[0x8001] = 0x03
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x8002), cpu.PC) // not taken, advance normally
}

func TestSTA_Extended(t *testing.T) {
	cpu, mem := newTestCPU(t)
	cpu.A = 0x42
	mem.data[0x8000] = 0xB7 // STA $1000
	mem.data[0x8001] = 0x10
	mem.data[0x8002] = 0x00
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(0x42), mem.data[0x1000])
}

func TestLDA_Direct(t *testing.T) {
	cpu, mem := newTestCPU(t)
	cpu.DP = 0x10
	mem.data[0x1042] = 0xAB
	mem.data[0x8000] = 0x96 // LDA <$42
	mem.data[0x8001] = 0x42
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(0xAB), cpu.A)
}

func TestTFR(t *testing.T) {
	cpu, mem := newTestCPU(t)
	cpu.SetD(0x1234)
	mem.data[0x8000] = 0x1F // TFR D,X
	mem.data[0x8001] = 0x01 // src=D(0), dst=X(1)
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x1234), cpu.X)
}

func TestEXG(t *testing.T) {
	cpu, mem := newTestCPU(t)
	cpu.SetD(0x1234)
	cpu.X = 0x5678
	mem.data[0x8000] = 0x1E // EXG D,X
	mem.data[0x8001] = 0x01 // src=D(0), dst=X(1)
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x5678), cpu.D())
	assert.Equal(t, uint16(0x1234), cpu.X)
}

func TestBSR_RTS(t *testing.T) {
	cpu, mem := newTestCPU(t)
	cpu.S = 0x0200
	mem.data[0x8000] = 0x8D // BSR +2
	mem.data[0x8001] = 0x02
	// Subroutine at $8004
	mem.data[0x8004] = 0x39 // RTS

	err := cpu.Step() // BSR
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x8004), cpu.PC)

	err = cpu.Step() // RTS
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x8002), cpu.PC) // return address
}

func TestMUL(t *testing.T) {
	cpu, mem := newTestCPU(t)
	cpu.A = 0x0A            // 10
	cpu.B = 0x14            // 20
	mem.data[0x8000] = 0x3D // MUL
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(200), cpu.D()) // 10 * 20 = 200
}

func TestSEX(t *testing.T) {
	cpu, mem := newTestCPU(t)
	cpu.B = 0x80            // negative
	mem.data[0x8000] = 0x1D // SEX
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(0xFF), cpu.A) // sign extended
}

func TestABX(t *testing.T) {
	cpu, mem := newTestCPU(t)
	cpu.X = 0x1000
	cpu.B = 0x42
	mem.data[0x8000] = 0x3A // ABX
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x1042), cpu.X)
}

func TestPSHS_PULS(t *testing.T) {
	cpu, mem := newTestCPU(t)
	cpu.S = 0x0200
	cpu.A = 0x42
	cpu.B = 0x43

	// PSHS A,B (mask = 0x06)
	mem.data[0x8000] = 0x34
	mem.data[0x8001] = 0x06
	err := cpu.Step()
	assert.NoError(t, err)

	// Clear registers
	cpu.A = 0
	cpu.B = 0

	// PULS A,B (mask = 0x06)
	mem.data[0x8002] = 0x35
	mem.data[0x8003] = 0x06
	err = cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(0x42), cpu.A)
	assert.Equal(t, uint8(0x43), cpu.B)
}

func TestPage2_CMPD(t *testing.T) {
	cpu, mem := newTestCPU(t)
	cpu.SetD(0x1234)
	mem.data[0x8000] = 0x10 // page 2 prefix
	mem.data[0x8001] = 0x83 // CMPD #$1234
	mem.data[0x8002] = 0x12
	mem.data[0x8003] = 0x34
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(1), cpu.Flags.Z) // equal
}

func TestPage3_SWI3(t *testing.T) {
	cpu, mem := newTestCPU(t)
	cpu.S = 0x0200
	// Set SWI3 vector
	mem.WriteWord(VectorSWI3, 0x9000)
	mem.data[0x8000] = 0x11 // page 3 prefix
	mem.data[0x8001] = 0x3F // SWI3
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x9000), cpu.PC)
}

func TestIndexed_ConstantOffset(t *testing.T) {
	cpu, mem := newTestCPU(t)
	cpu.X = 0x1000
	mem.data[0x1005] = 0x42
	mem.data[0x8000] = 0xA6 // LDA indexed
	mem.data[0x8001] = 0x05 // 5-bit offset: +5 from X (bit 7=0, reg=X=00, offset=00101)
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(0x42), cpu.A)
}

func TestIndexed_AutoIncrement(t *testing.T) {
	cpu, mem := newTestCPU(t)
	cpu.X = 0x1000
	mem.data[0x1000] = 0x42
	mem.data[0x8000] = 0xA6 // LDA indexed
	mem.data[0x8001] = 0x80 // ,X+ (postincrement by 1)
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(0x42), cpu.A)
	assert.Equal(t, uint16(0x1001), cpu.X) // X incremented
}

func TestLEAX(t *testing.T) {
	cpu, mem := newTestCPU(t)
	cpu.X = 0x1000
	mem.data[0x8000] = 0x30 // LEAX indexed
	mem.data[0x8001] = 0x01 // 5-bit offset: +1 from X
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x1001), cpu.X)
}

func TestInvalidOpcode(t *testing.T) {
	cpu, mem := newTestCPU(t)
	mem.data[0x8000] = 0x01 // illegal opcode
	err := cpu.Step()
	assert.Error(t, err)
}
