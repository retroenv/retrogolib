package m68000

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

// --- ALU Tests ---

func TestADD_DataReg(t *testing.T) {
	// ADD.W D0,D1 = 0xD240
	cpu := newTestCPUWithProgram(t, 0xD240)
	cpu.D[0] = 0x0010
	cpu.D[1] = 0x0020

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x0030), cpu.D[1]&0xFFFF)
	assert.Equal(t, uint8(0), cpu.Flags.Z)
	assert.Equal(t, uint8(0), cpu.Flags.N)
}

func TestADD_Overflow(t *testing.T) {
	cpu := newTestCPUWithProgram(t, 0xD240) // ADD.W D0,D1
	cpu.D[0] = 0x7FFF
	cpu.D[1] = 0x0001

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x8000), cpu.D[1]&0xFFFF)
	assert.Equal(t, uint8(1), cpu.Flags.V)
	assert.Equal(t, uint8(1), cpu.Flags.N)
}

func TestADDI(t *testing.T) {
	// ADDI.W #$0010,D0 = 0x0640, 0x0010
	cpu := newTestCPUWithProgram(t, 0x0640, 0x0010)
	cpu.D[0] = 0x0020

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x0030), cpu.D[0]&0xFFFF)
}

func TestADDQ(t *testing.T) {
	// ADDQ.W #3,D0 = 0x5640
	cpu := newTestCPUWithProgram(t, 0x5640)
	cpu.D[0] = 0x0010

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x0013), cpu.D[0]&0xFFFF)
}

func TestADDA(t *testing.T) {
	// ADDA.W D0,A0 = 0xD0C0
	cpu := newTestCPUWithProgram(t, 0xD0C0)
	cpu.D[0] = 0x0100
	cpu.A[0] = 0x2000

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x2100), cpu.A[0])
}

func TestSUB_DataReg(t *testing.T) {
	// SUB.W D0,D1 = 0x9240
	cpu := newTestCPUWithProgram(t, 0x9240)
	cpu.D[0] = 0x0010
	cpu.D[1] = 0x0030

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x0020), cpu.D[1]&0xFFFF)
}

func TestSUBI(t *testing.T) {
	// SUBI.W #$0010,D0 = 0x0440, 0x0010
	cpu := newTestCPUWithProgram(t, 0x0440, 0x0010)
	cpu.D[0] = 0x0030

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x0020), cpu.D[0]&0xFFFF)
}

func TestSUBQ(t *testing.T) {
	// SUBQ.W #1,D0 = 0x5340
	cpu := newTestCPUWithProgram(t, 0x5340)
	cpu.D[0] = 0x0010

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x000F), cpu.D[0]&0xFFFF)
}

func TestNEG(t *testing.T) {
	// NEG.W D0 = 0x4440
	cpu := newTestCPUWithProgram(t, 0x4440)
	cpu.D[0] = 0x0001

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0xFFFF), cpu.D[0]&0xFFFF)
	assert.Equal(t, uint8(1), cpu.Flags.N)
}

func TestCLR(t *testing.T) {
	// CLR.L D0 = 0x4280
	cpu := newTestCPUWithProgram(t, 0x4280)
	cpu.D[0] = 0x12345678

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), cpu.D[0])
	assert.Equal(t, uint8(1), cpu.Flags.Z)
	assert.Equal(t, uint8(0), cpu.Flags.N)
}

func TestCMP(t *testing.T) {
	// CMP.W D0,D1 = 0xB240
	cpu := newTestCPUWithProgram(t, 0xB240)
	cpu.D[0] = 0x0010
	cpu.D[1] = 0x0010

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(1), cpu.Flags.Z)
	assert.Equal(t, uint8(0), cpu.Flags.N)
	assert.Equal(t, uint8(0), cpu.Flags.C)
}

func TestCMPI(t *testing.T) {
	// CMPI.W #$0010,D0 = 0x0C40, 0x0010
	cpu := newTestCPUWithProgram(t, 0x0C40, 0x0010)
	cpu.D[0] = 0x0020

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(0), cpu.Flags.Z)
	assert.Equal(t, uint8(0), cpu.Flags.N)
}

func TestAND(t *testing.T) {
	// AND.W D0,D1 = 0xC240
	cpu := newTestCPUWithProgram(t, 0xC240)
	cpu.D[0] = 0xFF00
	cpu.D[1] = 0x1234

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x1200), cpu.D[1]&0xFFFF)
}

func TestANDI(t *testing.T) {
	// ANDI.W #$00FF,D0 = 0x0240, 0x00FF
	cpu := newTestCPUWithProgram(t, 0x0240, 0x00FF)
	cpu.D[0] = 0x1234

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x0034), cpu.D[0]&0xFFFF)
}

func TestOR(t *testing.T) {
	// OR.W D0,D1 = 0x8240
	cpu := newTestCPUWithProgram(t, 0x8240)
	cpu.D[0] = 0xFF00
	cpu.D[1] = 0x00FF

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0xFFFF), cpu.D[1]&0xFFFF)
}

func TestORI(t *testing.T) {
	// ORI.W #$FF00,D0 = 0x0040, 0xFF00
	cpu := newTestCPUWithProgram(t, 0x0040, 0xFF00)
	cpu.D[0] = 0x00FF

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0xFFFF), cpu.D[0]&0xFFFF)
}

func TestEOR(t *testing.T) {
	// EOR.W D1,D0 = 0xB340
	cpu := newTestCPUWithProgram(t, 0xB340)
	cpu.D[0] = 0xFF00
	cpu.D[1] = 0xFFFF

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x00FF), cpu.D[0]&0xFFFF)
}

func TestEORI(t *testing.T) {
	// EORI.W #$FFFF,D0 = 0x0A40, 0xFFFF
	cpu := newTestCPUWithProgram(t, 0x0A40, 0xFFFF)
	cpu.D[0] = 0xFF00

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x00FF), cpu.D[0]&0xFFFF)
}

func TestNOT(t *testing.T) {
	// NOT.W D0 = 0x4640
	cpu := newTestCPUWithProgram(t, 0x4640)
	cpu.D[0] = 0xFF00

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x00FF), cpu.D[0]&0xFFFF)
}

func TestEXT_Word(t *testing.T) {
	// EXT.W D0 = 0x4880
	cpu := newTestCPUWithProgram(t, 0x4880)
	cpu.D[0] = 0x000000FF

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x0000FFFF), cpu.D[0]&0xFFFF)
	assert.Equal(t, uint8(1), cpu.Flags.N)
}

func TestEXT_Long(t *testing.T) {
	// EXT.L D0 = 0x48C0
	cpu := newTestCPUWithProgram(t, 0x48C0)
	cpu.D[0] = 0x0000FF00

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0xFFFFFF00), cpu.D[0])
}

func TestTST(t *testing.T) {
	// TST.W D0 = 0x4A40
	cpu := newTestCPUWithProgram(t, 0x4A40)
	cpu.D[0] = 0x0000

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(1), cpu.Flags.Z)
	assert.Equal(t, uint8(0), cpu.Flags.N)
}

func TestMULU(t *testing.T) {
	// MULU D0,D1 = 0xC2C0
	cpu := newTestCPUWithProgram(t, 0xC2C0)
	cpu.D[0] = 0x0010
	cpu.D[1] = 0x0010

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x0100), cpu.D[1])
}

func TestMULS(t *testing.T) {
	// MULS D0,D1 = 0xC3C0
	cpu := newTestCPUWithProgram(t, 0xC3C0)
	cpu.D[0] = 0xFFFF // -1
	cpu.D[1] = 0x0010 // 16

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0xFFFFFFF0), cpu.D[1]) // -16
}

func TestDIVU(t *testing.T) {
	// DIVU D0,D1 = 0x82C0
	cpu := newTestCPUWithProgram(t, 0x82C0)
	cpu.D[0] = 0x0003
	cpu.D[1] = 0x000A // 10 / 3 = 3 remainder 1

	err := cpu.Step()
	assert.NoError(t, err)
	quotient := cpu.D[1] & 0xFFFF
	remainder := cpu.D[1] >> 16
	assert.Equal(t, uint32(3), quotient)
	assert.Equal(t, uint32(1), remainder)
}

func TestDIVU_ByZero(t *testing.T) {
	// DIVU D0,D1 = 0x82C0
	cpu := newTestCPUWithProgram(t, 0x82C0)
	cpu.D[0] = 0x0000
	cpu.D[1] = 0x000A
	// Set up divide-by-zero vector.
	cpu.bus.WriteLong(VectorDivZero*4, 0x00002000)

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x00002000), cpu.PC)
}

// --- MOVE Tests ---

func TestMOVE_DataRegToDataReg(t *testing.T) {
	// MOVE.L D0,D1 = 0x2200
	cpu := newTestCPUWithProgram(t, 0x2200)
	cpu.D[0] = 0x12345678

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x12345678), cpu.D[1])
}

func TestMOVEA(t *testing.T) {
	// MOVEA.L D0,A1 = 0x2240
	cpu := newTestCPUWithProgram(t, 0x2240)
	cpu.D[0] = 0x12345678

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x12345678), cpu.A[1])
}

func TestMOVEQ(t *testing.T) {
	// MOVEQ #-1,D0 = 0x70FF
	cpu := newTestCPUWithProgram(t, 0x70FF)

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0xFFFFFFFF), cpu.D[0])
	assert.Equal(t, uint8(1), cpu.Flags.N)
}

func TestMOVEQ_Positive(t *testing.T) {
	// MOVEQ #42,D3 = 0x762A
	cpu := newTestCPUWithProgram(t, 0x762A)

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(42), cpu.D[3])
	assert.Equal(t, uint8(0), cpu.Flags.N)
	assert.Equal(t, uint8(0), cpu.Flags.Z)
}

func TestLEA(t *testing.T) {
	// LEA (A0),A1 = 0x43D0
	cpu := newTestCPUWithProgram(t, 0x43D0)
	cpu.A[0] = 0x00005000

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x00005000), cpu.A[1])
}

func TestSWAP(t *testing.T) {
	// SWAP D0 = 0x4840
	cpu := newTestCPUWithProgram(t, 0x4840)
	cpu.D[0] = 0x12345678

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x56781234), cpu.D[0])
}

func TestEXG_DataRegs(t *testing.T) {
	// EXG D0,D1 = 0xC141
	cpu := newTestCPUWithProgram(t, 0xC141)
	cpu.D[0] = 0x11111111
	cpu.D[1] = 0x22222222

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x22222222), cpu.D[0])
	assert.Equal(t, uint32(0x11111111), cpu.D[1])
}

// --- Branch Tests ---

func TestBRA_Short(t *testing.T) {
	// BRA.S +4 = 0x6004
	cpu := newTestCPUWithProgram(t, 0x6004)

	err := cpu.Step()
	assert.NoError(t, err)
	// PC after reading opcode is 0x1002. Short branch: 0x1002 + 4 = 0x1006.
	assert.Equal(t, uint32(0x1006), cpu.PC)
}

func TestBRA_Long(t *testing.T) {
	// BRA.W +$0100 = 0x6000, 0x0100
	cpu := newTestCPUWithProgram(t, 0x6000, 0x0100)

	err := cpu.Step()
	assert.NoError(t, err)
	// PC base for 16-bit displacement is extension word address: 0x1002.
	assert.Equal(t, uint32(0x1102), cpu.PC)
}

func TestBcc_Taken(t *testing.T) {
	// BEQ.S +4 = 0x6704
	cpu := newTestCPUWithProgram(t, 0x6704)
	cpu.Flags.Z = 1

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x1006), cpu.PC)
}

func TestBcc_NotTaken(t *testing.T) {
	// BEQ.S +4 = 0x6704
	cpu := newTestCPUWithProgram(t, 0x6704)
	cpu.Flags.Z = 0

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x1002), cpu.PC)
}

func TestBSR(t *testing.T) {
	// BSR.S +4 = 0x6104
	cpu := newTestCPUWithProgram(t, 0x6104)
	oldSP := cpu.sp

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x1006), cpu.PC)
	assert.Equal(t, oldSP-4, cpu.sp)
	// Return address should be 0x1002 (after opcode word).
	assert.Equal(t, uint32(0x1002), cpu.bus.ReadLong(cpu.sp))
}

func TestJMP(t *testing.T) {
	// JMP (A0) = 0x4ED0
	cpu := newTestCPUWithProgram(t, 0x4ED0)
	cpu.A[0] = 0x00002000

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x00002000), cpu.PC)
}

func TestJSR(t *testing.T) {
	// JSR (A0) = 0x4E90
	cpu := newTestCPUWithProgram(t, 0x4E90)
	cpu.A[0] = 0x00002000
	oldSP := cpu.sp

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x00002000), cpu.PC)
	assert.Equal(t, oldSP-4, cpu.sp)
}

func TestRTS(t *testing.T) {
	// RTS = 0x4E75
	cpu := newTestCPUWithProgram(t, 0x4E75)
	cpu.push32(0x00003000) // Push return address.

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x00003000), cpu.PC)
}

func TestNOP(t *testing.T) {
	// NOP = 0x4E71
	cpu := newTestCPUWithProgram(t, 0x4E71)
	pcBefore := cpu.PC

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, pcBefore+2, cpu.PC)
}

func TestDBcc_Loop(t *testing.T) {
	// DBF D0,displacement = 0x51C8, displacement
	// DBRA (false condition, always loops) D0 with displacement -4
	cpu := newTestCPUWithProgram(t, 0x51C8, 0xFFFC) // displacement = -4

	cpu.D[0] = 2 // Counter

	// First iteration.
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), cpu.D[0]&0xFFFF)
	// PC should go back: 0x1002 + (-4) = 0x0FFE.
	assert.Equal(t, uint32(0x0FFE), cpu.PC)
}

func TestScc_True(t *testing.T) {
	// ST D0 (Set True) = 0x50C0
	cpu := newTestCPUWithProgram(t, 0x50C0)

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0xFF), cpu.D[0]&0xFF)
}

func TestScc_False(t *testing.T) {
	// SF D0 (Set False) = 0x51C0
	cpu := newTestCPUWithProgram(t, 0x51C0)
	cpu.D[0] = 0xFF

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x00), cpu.D[0]&0xFF)
}

// --- Shift/Rotate Tests ---

func TestASL_Register(t *testing.T) {
	// ASL.W #1,D0 = 0xE340
	cpu := newTestCPUWithProgram(t, 0xE340)
	cpu.D[0] = 0x4000

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x8000), cpu.D[0]&0xFFFF)
	assert.Equal(t, uint8(1), cpu.Flags.N)
}

func TestLSR_Register(t *testing.T) {
	// LSR.W #1,D0 = 0xE248
	cpu := newTestCPUWithProgram(t, 0xE248)
	cpu.D[0] = 0x0002

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x0001), cpu.D[0]&0xFFFF)
}

func TestROL_Register(t *testing.T) {
	// ROL.W #1,D0 = 0xE358
	cpu := newTestCPUWithProgram(t, 0xE358)
	cpu.D[0] = 0x8000

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x0001), cpu.D[0]&0xFFFF)
	assert.Equal(t, uint8(1), cpu.Flags.C)
}

func TestROR_Register(t *testing.T) {
	// ROR.W #1,D0 = 0xE258
	cpu := newTestCPUWithProgram(t, 0xE258)
	cpu.D[0] = 0x0001

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x8000), cpu.D[0]&0xFFFF)
	assert.Equal(t, uint8(1), cpu.Flags.C)
}

// --- Bit Manipulation Tests ---

func TestBTST_DataReg(t *testing.T) {
	// BTST D0,D1 = 0x0101
	cpu := newTestCPUWithProgram(t, 0x0101)
	cpu.D[0] = 4 // Test bit 4.
	cpu.D[1] = 0x0010

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(0), cpu.Flags.Z) // Bit 4 is set.
}

func TestBTST_DataReg_NotSet(t *testing.T) {
	cpu := newTestCPUWithProgram(t, 0x0101) // BTST D0,D1
	cpu.D[0] = 5                            // Test bit 5.
	cpu.D[1] = 0x0010

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(1), cpu.Flags.Z) // Bit 5 is not set.
}

func TestBSET_DataReg(t *testing.T) {
	// BSET D0,D1 = 0x01C1
	cpu := newTestCPUWithProgram(t, 0x01C1)
	cpu.D[0] = 3
	cpu.D[1] = 0x0000

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x0008), cpu.D[1])
	assert.Equal(t, uint8(1), cpu.Flags.Z) // Was clear before set.
}

func TestBCLR_DataReg(t *testing.T) {
	// BCLR D0,D1 = 0x0181
	cpu := newTestCPUWithProgram(t, 0x0181)
	cpu.D[0] = 3
	cpu.D[1] = 0x000F

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x0007), cpu.D[1])
	assert.Equal(t, uint8(0), cpu.Flags.Z) // Was set before clear.
}

func TestBCHG_DataReg(t *testing.T) {
	// BCHG D0,D1 = 0x0141
	cpu := newTestCPUWithProgram(t, 0x0141)
	cpu.D[0] = 0
	cpu.D[1] = 0x0001

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x0000), cpu.D[1])
	assert.Equal(t, uint8(0), cpu.Flags.Z)
}

// --- System Instruction Tests ---

func TestTRAP(t *testing.T) {
	// TRAP #0 = 0x4E40
	cpu := newTestCPUWithProgram(t, 0x4E40)
	cpu.bus.WriteLong(uint32(VectorTrap0)*4, 0x00002000)

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x00002000), cpu.PC)
}

func TestTRAPV_NoOverflow(t *testing.T) {
	// TRAPV = 0x4E76
	cpu := newTestCPUWithProgram(t, 0x4E76)
	cpu.Flags.V = 0

	err := cpu.Step()
	assert.NoError(t, err)
	// Should not trap.
	assert.Equal(t, uint32(0x1002), cpu.PC)
}

func TestTRAPV_Overflow(t *testing.T) {
	cpu := newTestCPUWithProgram(t, 0x4E76)
	cpu.Flags.V = 1
	cpu.bus.WriteLong(uint32(VectorTRAPV)*4, 0x00003000)

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x00003000), cpu.PC)
}

func TestRESET_Supervisor(t *testing.T) {
	// RESET = 0x4E70
	cpu := newTestCPUWithProgram(t, 0x4E70)
	// Already in supervisor mode.

	err := cpu.Step()
	assert.NoError(t, err)
}

func TestRESET_UserMode(t *testing.T) {
	cpu := newTestCPUWithProgram(t, 0x4E70)
	cpu.SetSR(cpu.GetSR() & ^uint16(MaskSupervisor)) // Switch to user mode.
	cpu.bus.WriteLong(uint32(VectorPrivilege)*4, 0x00004000)

	err := cpu.Step()
	assert.NoError(t, err)
	// Should trigger privilege violation.
	assert.Equal(t, uint32(0x00004000), cpu.PC)
}

func TestSTOP_Supervisor(t *testing.T) {
	// STOP #$2000 = 0x4E72, 0x2000
	cpu := newTestCPUWithProgram(t, 0x4E72, 0x2000)

	err := cpu.Step()
	assert.NoError(t, err)
	assert.True(t, cpu.stopped)
}

func TestTAS(t *testing.T) {
	// TAS D0 = 0x4AC0
	cpu := newTestCPUWithProgram(t, 0x4AC0)
	cpu.D[0] = 0x00

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x80), cpu.D[0]&0xFF)
	assert.Equal(t, uint8(1), cpu.Flags.Z) // Was zero before.
}

func TestILLEGAL(t *testing.T) {
	// ILLEGAL = 0x4AFC
	cpu := newTestCPUWithProgram(t, 0x4AFC)
	cpu.bus.WriteLong(uint32(VectorIllegal)*4, 0x00005000)

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x00005000), cpu.PC)
}

// --- Condition Code Tests ---

func TestConditions(t *testing.T) {
	cpu := newTestCPU(t)

	tests := []struct {
		name  string
		cond  uint16
		flags Flags
		want  bool
	}{
		{"T", 0, Flags{}, true},
		{"F", 1, Flags{}, false},
		{"HI: !C&&!Z", 2, Flags{C: 0, Z: 0}, true},
		{"HI: C=1", 2, Flags{C: 1, Z: 0}, false},
		{"LS: C=1", 3, Flags{C: 1, Z: 0}, true},
		{"CC: !C", 4, Flags{C: 0}, true},
		{"CS: C", 5, Flags{C: 1}, true},
		{"NE: !Z", 6, Flags{Z: 0}, true},
		{"EQ: Z", 7, Flags{Z: 1}, true},
		{"VC: !V", 8, Flags{V: 0}, true},
		{"VS: V", 9, Flags{V: 1}, true},
		{"PL: !N", 10, Flags{N: 0}, true},
		{"MI: N", 11, Flags{N: 1}, true},
		{"GE: N=V=0", 12, Flags{N: 0, V: 0}, true},
		{"GE: N=V=1", 12, Flags{N: 1, V: 1}, true},
		{"LT: N!=V", 13, Flags{N: 1, V: 0}, true},
		{"GT: Z=0,N=V", 14, Flags{Z: 0, N: 0, V: 0}, true},
		{"LE: Z=1", 15, Flags{Z: 1, N: 0, V: 0}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu.Flags = tt.flags
			got := cpu.evaluateCondition(tt.cond)
			assert.Equal(t, tt.want, got)
		})
	}
}

// --- Exception Tests ---

func TestProcessException(t *testing.T) {
	cpu := newTestCPU(t)
	cpu.PC = 0x1000
	cpu.bus.WriteLong(uint32(VectorIllegal)*4, 0x2000)
	oldSP := cpu.sp

	err := cpu.processException(VectorIllegal)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x2000), cpu.PC)
	assert.True(t, cpu.IsSupervisor())
	assert.Equal(t, oldSP-6, cpu.sp) // Pushed PC (4) + SR (2)
}

// --- LINK/UNLK Tests ---

func TestLINK_UNLK(t *testing.T) {
	// LINK A6,#-8 = 0x4E56, 0xFFF8
	cpu := newTestCPUWithProgram(t, 0x4E56, 0xFFF8)
	cpu.A[6] = 0x12345678
	oldSP := cpu.sp

	err := cpu.Step()
	assert.NoError(t, err)
	// A6 should be old SP - 4 (after pushing old A6).
	assert.Equal(t, oldSP-4, cpu.A[6])
	// SP should be A6 + displacement (-8).
	assert.Equal(t, oldSP-12, cpu.sp)

	// Test UNLK A6 = 0x4E5E.
	cpu.bus.WriteWord(cpu.PC, 0x4E5E)
	err = cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x12345678), cpu.A[6])
}

// newTestCPUWithProgram creates a CPU with program data at PC.
func newTestCPUWithProgram(t *testing.T, program ...uint16) *CPU {
	t.Helper()
	mem := NewBasicMemory()
	bus := NewBasicBus(mem)
	cpu, err := New(bus, WithInitialPC(0x1000), WithInitialSP(0x10000))
	assert.NoError(t, err)

	for i, word := range program {
		mem.WriteWord(0x1000+uint32(i)*2, word)
	}
	return cpu
}
