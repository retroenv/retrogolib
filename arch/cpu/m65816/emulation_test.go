package m65816

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

// -- NOP --

func TestNOP(t *testing.T) {
	cpu, mem := setupCPU(t)
	writeOp(mem, 0x8000, 0xEA) // NOP
	before := cpu.cycles
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x8001), cpu.PC)
	assert.Equal(t, uint64(2), cpu.cycles-before)
}

// -- LDA / STA 8-bit --

func TestLDA_Immediate8(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 1
	writeOp(mem, 0x8000, 0xA9, 0x42) // LDA #$42
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(0x42), cpu.A())
	assert.Equal(t, uint16(0x8002), cpu.PC)
}

func TestLDA_Immediate16(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 0                        // 16-bit accumulator
	writeOp(mem, 0x8000, 0xA9, 0x34, 0x12) // LDA #$1234
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x1234), cpu.C)
	assert.Equal(t, uint16(0x8003), cpu.PC)
}

func TestLDA_Flags(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 1
	writeOp(mem, 0x8000, 0xA9, 0x00) // LDA #$00
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(1), cpu.Flags.Z)
	assert.Equal(t, uint8(0), cpu.Flags.N)

	cpu.PC = 0x8000
	writeOp(mem, 0x8000, 0xA9, 0x80) // LDA #$80
	err = cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(1), cpu.Flags.N)
}

func TestSTA_DirectPage(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 1
	cpu.C = 0xAB
	cpu.DP = 0x0000
	writeOp(mem, 0x8000, 0x85, 0x10) // STA $10
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(0xAB), mem.data[0x10])
}

// -- ADC 8-bit --

func TestADC_NoCarry(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 1
	cpu.Flags.C = 0
	cpu.C = 0x10
	writeOp(mem, 0x8000, 0x69, 0x20) // ADC #$20
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(0x30), cpu.A())
	assert.Equal(t, uint8(0), cpu.Flags.C)
}

func TestADC_WithCarry(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 1
	cpu.Flags.C = 0
	cpu.C = 0xFF
	writeOp(mem, 0x8000, 0x69, 0x01) // ADC #$01
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(0x00), cpu.A())
	assert.Equal(t, uint8(1), cpu.Flags.C)
	assert.Equal(t, uint8(1), cpu.Flags.Z)
}

// -- SBC 8-bit --

func TestSBC_Basic(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 1
	cpu.Flags.C = 1 // no borrow
	cpu.C = 0x50
	writeOp(mem, 0x8000, 0xE9, 0x20) // SBC #$20
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(0x30), cpu.A())
	assert.Equal(t, uint8(1), cpu.Flags.C)
}

// -- INC/DEC accumulator --

func TestINC_Accumulator(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 1
	cpu.C = 0x41
	writeOp(mem, 0x8000, 0x1A) // INC A
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(0x42), cpu.A())
}

func TestDEC_Accumulator(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 1
	cpu.C = 0x43
	writeOp(mem, 0x8000, 0x3A) // DEC A
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(0x42), cpu.A())
}

// -- ASL --

func TestASL_Accumulator8(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 1
	cpu.C = 0x41
	writeOp(mem, 0x8000, 0x0A) // ASL A
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(0x82), cpu.A())
	assert.Equal(t, uint8(0), cpu.Flags.C)
}

func TestASL_SetsCarry(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 1
	cpu.C = 0x80
	writeOp(mem, 0x8000, 0x0A) // ASL A
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(1), cpu.Flags.C)
	assert.Equal(t, uint8(0x00), cpu.A())
}

// -- AND / ORA / EOR --

func TestAND_Immediate(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 1
	cpu.C = 0xFF
	writeOp(mem, 0x8000, 0x29, 0x0F) // AND #$0F
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(0x0F), cpu.A())
}

func TestORA_Immediate(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 1
	cpu.C = 0x00
	writeOp(mem, 0x8000, 0x09, 0xF0) // ORA #$F0
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(0xF0), cpu.A())
}

func TestEOR_Immediate(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 1
	cpu.C = 0xFF
	writeOp(mem, 0x8000, 0x49, 0x0F) // EOR #$0F
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(0xF0), cpu.A())
}

// -- CMP --

func TestCMP_Equal(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 1
	cpu.C = 0x42
	writeOp(mem, 0x8000, 0xC9, 0x42) // CMP #$42
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(1), cpu.Flags.Z)
	assert.Equal(t, uint8(1), cpu.Flags.C)
	assert.Equal(t, uint8(0), cpu.Flags.N)
}

// -- Branch instructions --

func TestBRA_Taken(t *testing.T) {
	cpu, mem := setupCPU(t)
	writeOp(mem, 0x8000, 0x80, 0x05) // BRA +5
	err := cpu.Step()
	assert.NoError(t, err)
	// target = 0x8000 + 2 + 5 = 0x8007
	assert.Equal(t, uint16(0x8007), cpu.PC)
}

func TestBNE_NotTaken(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.Z = 1                  // Z set means equal, BNE not taken
	writeOp(mem, 0x8000, 0xD0, 0x05) // BNE +5
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x8002), cpu.PC)
}

func TestBNE_Taken(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.Z = 0
	writeOp(mem, 0x8000, 0xD0, 0x05) // BNE +5
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x8007), cpu.PC)
}

// -- Cycle accuracy --

// LDA abs,X: 4 cycles with no page crossing, 5 cycles when X crosses a page boundary.
func TestLDA_AbsX_PageCross_Cycles(t *testing.T) {
	// No page cross: base=$8100, X=$01 → eff=$8101 (same page $81)
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 1
	cpu.Flags.X = 1
	cpu.DB = 0x00
	cpu.X = 0x01
	writeOp(mem, 0x8000, 0xBD, 0x00, 0x81) // LDA $8100,X
	before := cpu.cycles
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint64(4), cpu.cycles-before)

	// Page cross: base=$81FF, X=$01 → eff=$8200 (crosses $81→$82)
	cpu.PC = 0x8000
	cpu.X = 0x01
	writeOp(mem, 0x8000, 0xBD, 0xFF, 0x81) // LDA $81FF,X
	before = cpu.cycles
	err = cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint64(5), cpu.cycles-before)
}

// BNE cycle counts: 2 not taken, 3 taken same page, 3 taken cross page (native),
// 4 taken cross page (emulation mode).
func TestBranch_Cycles(t *testing.T) {
	// Not taken: 2 cycles
	cpu, mem := setupCPU(t)
	cpu.Flags.Z = 1                  // Z set → BNE not taken
	writeOp(mem, 0x8000, 0xD0, 0x05) // BNE +5
	before := cpu.cycles
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), cpu.cycles-before)

	// Taken, same page: 3 cycles
	cpu.PC = 0x8000
	cpu.Flags.Z = 0                  // BNE taken
	writeOp(mem, 0x8000, 0xD0, 0x05) // BNE +5 → target=$8007 (page $80, same as nextPC $8002)
	before = cpu.cycles
	err = cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint64(3), cpu.cycles-before)

	// Taken, cross page, native mode: 3 cycles (no page-cross penalty in native mode)
	// BNE at $82FE, offset=-128 ($80) → nextPC=$8300, target=$8280 (page $82 ≠ $83)
	cpu.PC = 0x82FE
	cpu.E = false // native mode (already set by setupCPU)
	cpu.Flags.Z = 0
	writeOp(mem, 0x82FE, 0xD0, 0x80) // BNE -128
	before = cpu.cycles
	err = cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint64(3), cpu.cycles-before)

	// Taken, cross page, emulation mode: 4 cycles
	cpu.PC = 0x82FE
	cpu.E = true
	cpu.Flags.M = 1
	cpu.Flags.X = 1
	cpu.Flags.Z = 0
	writeOp(mem, 0x82FE, 0xD0, 0x80) // BNE -128
	before = cpu.cycles
	err = cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint64(4), cpu.cycles-before)
}

// BRA (always taken): 3 cycles same page, 4 cycles cross page in emulation mode.
func TestBRA_Cycles(t *testing.T) {
	// Same page: 3 cycles (native)
	cpu, mem := setupCPU(t)
	writeOp(mem, 0x8000, 0x80, 0x05) // BRA +5 → target=$8007 (same page as nextPC=$8002)
	before := cpu.cycles
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint64(3), cpu.cycles-before)

	// Cross page, emulation mode: 4 cycles
	// BRA at $82FE, offset=-128 → nextPC=$8300, target=$8280 (page $82 ≠ $83)
	cpu.PC = 0x82FE
	cpu.E = true
	cpu.Flags.M = 1
	cpu.Flags.X = 1
	writeOp(mem, 0x82FE, 0x80, 0x80) // BRA -128
	before = cpu.cycles
	err = cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint64(4), cpu.cycles-before)
}

// -- Stack operations --

func TestPHA_PLA_8bit(t *testing.T) {
	cpu, _ := setupCPU(t)
	cpu.Flags.M = 1
	cpu.SP = 0x01FF
	cpu.C = 0x55
	err := pha(cpu)
	assert.NoError(t, err)
	cpu.C = 0
	err = pla(cpu)
	assert.NoError(t, err)
	assert.Equal(t, uint8(0x55), cpu.A())
}

func TestPHX_PLX_8bit(t *testing.T) {
	cpu, _ := setupCPU(t)
	cpu.Flags.X = 1
	cpu.SP = 0x01FF
	cpu.X = 0x77
	err := phx(cpu)
	assert.NoError(t, err)
	cpu.X = 0
	err = plx(cpu)
	assert.NoError(t, err)
	assert.Equal(t, uint8(0x77), uint8(cpu.X))
}

// -- Register transfers --

func TestTAX_8bit(t *testing.T) {
	cpu, _ := setupCPU(t)
	cpu.Flags.X = 1
	cpu.C = 0x42
	err := tax(cpu)
	assert.NoError(t, err)
	assert.Equal(t, uint8(0x42), uint8(cpu.X))
}

func TestTCD_TDC(t *testing.T) {
	cpu, _ := setupCPU(t)
	cpu.C = 0x1234
	err := tcd(cpu)
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x1234), cpu.DP)
	cpu.C = 0
	err = tdc(cpu)
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x1234), cpu.C)
}

func TestTCS_TSC(t *testing.T) {
	cpu, _ := setupCPU(t)
	cpu.E = false
	cpu.C = 0x0150
	err := tcs(cpu)
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x0150), cpu.SP)
	cpu.C = 0
	err = tsc(cpu)
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x0150), cpu.C)
}

// -- Processor status --

func TestSEP_REP(t *testing.T) {
	cpu, _ := setupCPU(t)
	cpu.E = false
	cpu.Flags.M = 0
	cpu.Flags.X = 0
	// SEP #$30 sets M and X flags
	err := sep(cpu, Immediate8(0x30))
	assert.NoError(t, err)
	assert.Equal(t, uint8(1), cpu.Flags.M)
	assert.Equal(t, uint8(1), cpu.Flags.X)
	// REP #$30 clears M and X flags
	err = rep(cpu, Immediate8(0x30))
	assert.NoError(t, err)
	assert.Equal(t, uint8(0), cpu.Flags.M)
	assert.Equal(t, uint8(0), cpu.Flags.X)
}

func TestXCE_ToNative(t *testing.T) {
	cpu, _ := newTestCPU(t)
	// Start in emulation mode (E=true)
	assert.True(t, cpu.E)
	cpu.Flags.C = 0 // carry = 0
	err := xce(cpu)
	assert.NoError(t, err)
	assert.False(t, cpu.E)
	assert.Equal(t, uint8(1), cpu.Flags.C)
}

func TestXBA(t *testing.T) {
	cpu, _ := setupCPU(t)
	cpu.C = 0x1234
	err := xba(cpu)
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x3412), cpu.C)
}

// -- STP / WAI --

func TestSTP(t *testing.T) {
	cpu, _ := setupCPU(t)
	err := stp(cpu)
	assert.NoError(t, err)
	assert.True(t, cpu.stopped)
	// Step should be no-op when stopped
	before := cpu.PC
	err = cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, before, cpu.PC)
}

// -- 16-bit accumulator --

func TestLDA16_STA16(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 0 // 16-bit accumulator
	cpu.DP = 0x0000

	// LDA #$ABCD
	writeOp(mem, 0x8000, 0xA9, 0xCD, 0xAB)
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0xABCD), cpu.C)

	// STA $20 (direct page)
	cpu.PC = 0x8003
	writeOp(mem, 0x8003, 0x85, 0x20)
	err = cpu.Step()
	assert.NoError(t, err)
	lo := mem.data[0x20]
	hi := mem.data[0x21]
	if lo != 0xCD || hi != 0xAB {
		t.Errorf("STA16: mem=$%02X%02X, want ABCD", hi, lo)
	}
}

// -- INX/INY/DEX/DEY 16-bit --

func TestINX16(t *testing.T) {
	cpu, _ := setupCPU(t)
	cpu.Flags.X = 0
	cpu.X = 0xFFFE
	err := inx(cpu)
	assert.NoError(t, err)
	assert.Equal(t, uint16(0xFFFF), cpu.X)
}

func TestDEX16(t *testing.T) {
	cpu, _ := setupCPU(t)
	cpu.Flags.X = 0
	cpu.X = 0x0001
	err := dex(cpu)
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x0000), cpu.X)
	assert.Equal(t, uint8(1), cpu.Flags.Z)
}

// -- ADC/SBC 16-bit --

func TestADC_16bit(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 0 // 16-bit accumulator
	cpu.Flags.C = 0
	cpu.C = 0x1000
	writeOp(mem, 0x8000, 0x69, 0x34, 0x12) // ADC #$1234
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x2234), cpu.C)
	assert.Equal(t, uint8(0), cpu.Flags.C)
}

func TestADC_16bit_Carry(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 0
	cpu.Flags.C = 0
	cpu.C = 0xFF00
	writeOp(mem, 0x8000, 0x69, 0x00, 0x01) // ADC #$0100
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x0000), cpu.C)
	assert.Equal(t, uint8(1), cpu.Flags.C)
	assert.Equal(t, uint8(1), cpu.Flags.Z)
}

func TestSBC_16bit(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 0
	cpu.Flags.C = 1 // no borrow
	cpu.C = 0x1234
	writeOp(mem, 0x8000, 0xE9, 0x34, 0x02) // SBC #$0234
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x1000), cpu.C)
	assert.Equal(t, uint8(1), cpu.Flags.C)
}

// -- DirectPage indexed addressing --

func TestSTADP_IndexedX_8bit(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.E = false
	cpu.Flags.M = 1 // 8-bit acc
	cpu.Flags.X = 1 // 8-bit X
	cpu.DP = 0x0000
	cpu.C = 0xBB
	cpu.X = 0x10
	writeOp(mem, 0x8000, 0x95, 0x20) // STA $20,X  ->  EA = 0x0000 + 0x20 + 0x10 = 0x0030
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(0xBB), mem.data[0x0030])
}

func TestSTADP_IndexedX_16bit(t *testing.T) {
	// 16-bit X: EA = DP + dp + X (no 8-bit truncation)
	cpu, mem := setupCPU(t)
	cpu.E = false
	cpu.Flags.M = 1 // 8-bit acc
	cpu.Flags.X = 0 // 16-bit X
	cpu.DP = 0x0000
	cpu.C = 0xAA
	cpu.X = 0x0100                   // high byte of X is non-zero: tests the fix
	writeOp(mem, 0x8000, 0x95, 0xF0) // STA $F0,X  ->  EA = 0x0000 + 0xF0 + 0x0100 = 0x01F0
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(0xAA), mem.data[0x01F0])
	// Verify wrong address ($00F0) was NOT written (old bug would have written there)
	assert.Equal(t, uint8(0x00), mem.data[0x00F0])
}

func TestSTADP_IndexedX_EmulationWrap(t *testing.T) {
	// In emulation mode with DP=$0000, (dp+X) wraps within page 0
	cpu, mem := newTestCPU(t)
	// Stays in emulation mode (E=true)
	cpu.DP = 0x0000
	cpu.C = 0x77
	cpu.X = 0x20
	writeOp(mem, 0x8000, 0x95, 0xF0) // STA $F0,X -> (0xF0+0x20)&0xFF = 0x10
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(0x77), mem.data[0x0010])
}

// -- JSR / RTS --

func TestJSR_RTS(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.SP = 0x01FF
	// JSR $9000 at $8000; subroutine at $9000 contains RTS
	writeOp(mem, 0x8000, 0x20, 0x00, 0x90) // JSR $9000
	writeOp(mem, 0x9000, 0x60)             // RTS
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x9000), cpu.PC)
	assert.Equal(t, uint16(0x01FD), cpu.SP)
	// Return address $8002 pushed: high byte ($80) at $01FF, low byte ($02) at $01FE
	if mem.data[0x01FE] != 0x02 || mem.data[0x01FF] != 0x80 {
		t.Errorf("JSR: stack hi=%02X lo=%02X, want hi=80 lo=02", mem.data[0x01FF], mem.data[0x01FE])
	}

	err = cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x8003), cpu.PC)
	assert.Equal(t, uint16(0x01FF), cpu.SP)
}

// -- JSL / RTL --

func TestJSL_RTL(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.SP = 0x01FF
	cpu.PB = 0x00
	// JSL $019000 at $8000; subroutine at bank $01:$9000 contains RTL
	writeOp(mem, 0x8000, 0x22, 0x00, 0x90, 0x01) // JSL $019000
	mem.data[0x019000] = 0x6B                    // RTL (in bank $01)
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x9000), cpu.PC)
	assert.Equal(t, uint8(0x01), cpu.PB)
	assert.Equal(t, uint16(0x01FC), cpu.SP)
	err = cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x8004), cpu.PC)
	assert.Equal(t, uint8(0x00), cpu.PB)
	assert.Equal(t, uint16(0x01FF), cpu.SP)
}

// -- JMP / JML --

func TestJMP_Absolute(t *testing.T) {
	cpu, mem := setupCPU(t)
	writeOp(mem, 0x8000, 0x4C, 0x34, 0x12) // JMP $1234
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x1234), cpu.PC)
}

func TestJML_Long(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.PB = 0x00
	writeOp(mem, 0x8000, 0x5C, 0x00, 0x90, 0x01) // JML $019000
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x9000), cpu.PC)
	assert.Equal(t, uint8(0x01), cpu.PB)
}

// -- PEA / PEI / PER --

func TestPEA(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.SP = 0x01FF
	writeOp(mem, 0x8000, 0xF4, 0x34, 0x12) // PEA $1234
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x01FD), cpu.SP)
	// $1234 pushed: lo=$34 at $01FE, hi=$12 at $01FF
	if mem.data[0x01FE] != 0x34 || mem.data[0x01FF] != 0x12 {
		t.Errorf("PEA: stack=%02X%02X, want 1234", mem.data[0x01FF], mem.data[0x01FE])
	}
	assert.Equal(t, uint16(0x8003), cpu.PC)
}

func TestPEI(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.SP = 0x01FF
	cpu.DP = 0x0000
	// Store $5678 at direct page offset $10
	mem.data[0x0010] = 0x78
	mem.data[0x0011] = 0x56
	writeOp(mem, 0x8000, 0xD4, 0x10) // PEI ($10)
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x01FD), cpu.SP)
	got := uint16(mem.data[0x01FF])<<8 | uint16(mem.data[0x01FE])
	assert.Equal(t, uint16(0x5678), got)
	assert.Equal(t, uint16(0x8002), cpu.PC)
}

func TestPER(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.SP = 0x01FF
	// PER +$0100: effective = 0x8000 + 3 + 0x0100 = 0x8103
	writeOp(mem, 0x8000, 0x62, 0x00, 0x01) // PER +$0100
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x01FD), cpu.SP)
	got := uint16(mem.data[0x01FF])<<8 | uint16(mem.data[0x01FE])
	assert.Equal(t, uint16(0x8103), got)
	assert.Equal(t, uint16(0x8003), cpu.PC)
}

// -- MVN block move --

func TestMVN_SingleByte(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.X = 0 // 16-bit index registers
	cpu.C = 0x0000  // 1 byte to copy (C+1)
	cpu.X = 0x1000
	cpu.Y = 0x2000
	mem.data[0x1000] = 0xAA
	err := mvn(cpu, BlockMove{Src: 0x00, Dst: 0x00})
	assert.NoError(t, err)
	assert.Equal(t, uint8(0xAA), mem.data[0x2000])
	assert.Equal(t, uint16(0xFFFF), cpu.C)
	assert.Equal(t, uint16(0x1001), cpu.X)
	assert.Equal(t, uint16(0x2001), cpu.Y)
}

func TestMVN_ThreeBytes(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.X = 0 // 16-bit index registers
	cpu.C = 0x0002  // 3 bytes to copy
	cpu.X = 0x1000
	cpu.Y = 0x2000
	mem.data[0x1000] = 0xAA
	mem.data[0x1001] = 0xBB
	mem.data[0x1002] = 0xCC
	err := mvn(cpu, BlockMove{Src: 0x00, Dst: 0x00})
	assert.NoError(t, err)
	if mem.data[0x2000] != 0xAA || mem.data[0x2001] != 0xBB || mem.data[0x2002] != 0xCC {
		t.Errorf("MVN 3-byte: dst=%02X%02X%02X, want AABBCC",
			mem.data[0x2000], mem.data[0x2001], mem.data[0x2002])
	}
	assert.Equal(t, uint16(0xFFFF), cpu.C)
	assert.Equal(t, uint16(0x1003), cpu.X)
}

// -- BRK (native mode) --

func TestBRK_Native(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.E = false
	cpu.SP = 0x01FF
	cpu.PB = 0x00
	cpu.Flags.D = 1 // should be cleared by BRK
	// Set native BRK vector at $FFE6 -> $9000
	mem.data[0xFFE6] = 0x00
	mem.data[0xFFE7] = 0x90
	writeOp(mem, 0x8000, 0x00, 0x00) // BRK + signature byte
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x9000), cpu.PC)
	assert.Equal(t, uint8(0x00), cpu.PB)
	assert.Equal(t, uint8(1), cpu.Flags.I)
	assert.Equal(t, uint8(0), cpu.Flags.D)
	assert.Equal(t, uint16(0x01FB), cpu.SP)
}

// -- ADC decimal mode (BCD) --

func TestADC_Decimal8_Basic(t *testing.T) {
	cpu, _ := setupCPU(t)
	cpu.Flags.D = 1
	cpu.Flags.M = 1
	cpu.Flags.C = 0
	cpu.C = 0x0025 // BCD 25
	err := adc(cpu, Immediate8(0x13))
	assert.NoError(t, err)
	assert.Equal(t, uint8(0x38), cpu.A())
	assert.Equal(t, uint8(0), cpu.Flags.C)
}

func TestADC_Decimal8_CarryOut(t *testing.T) {
	cpu, _ := setupCPU(t)
	cpu.Flags.D = 1
	cpu.Flags.M = 1
	cpu.Flags.C = 0
	cpu.C = 0x99 // BCD 99
	err := adc(cpu, Immediate8(0x01))
	assert.NoError(t, err)
	assert.Equal(t, uint8(0x00), cpu.A())
	assert.Equal(t, uint8(1), cpu.Flags.C)
	assert.Equal(t, uint8(1), cpu.Flags.Z)
}

func TestADC_Decimal8_LowNibbleAdjust(t *testing.T) {
	// Low nibble overflow: 5+5=10 → adjust to 0 with carry to high nibble
	cpu, _ := setupCPU(t)
	cpu.Flags.D = 1
	cpu.Flags.M = 1
	cpu.Flags.C = 0
	cpu.C = 0x15 // BCD 15
	err := adc(cpu, Immediate8(0x15))
	assert.NoError(t, err)
	assert.Equal(t, uint8(0x30), cpu.A())
	assert.Equal(t, uint8(0), cpu.Flags.C)
}

func TestADC_Decimal16(t *testing.T) {
	cpu, _ := setupCPU(t)
	cpu.Flags.D = 1
	cpu.Flags.M = 0 // 16-bit accumulator
	cpu.Flags.C = 0
	cpu.C = 0x1234 // BCD 1234
	err := adc(cpu, Immediate16(0x4321))
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x5555), cpu.C)
	assert.Equal(t, uint8(0), cpu.Flags.C)
}

func TestADC_Decimal16_CarryOut(t *testing.T) {
	cpu, _ := setupCPU(t)
	cpu.Flags.D = 1
	cpu.Flags.M = 0
	cpu.Flags.C = 0
	cpu.C = 0x9999 // BCD 9999
	err := adc(cpu, Immediate16(0x0001))
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x0000), cpu.C)
	assert.Equal(t, uint8(1), cpu.Flags.C)
}

// -- SBC decimal mode (BCD) --

func TestSBC_Decimal8_Basic(t *testing.T) {
	cpu, _ := setupCPU(t)
	cpu.Flags.D = 1
	cpu.Flags.M = 1
	cpu.Flags.C = 1 // no borrow
	cpu.C = 0x50    // BCD 50
	err := sbc(cpu, Immediate8(0x20))
	assert.NoError(t, err)
	assert.Equal(t, uint8(0x30), cpu.A())
	assert.Equal(t, uint8(1), cpu.Flags.C)
}

func TestSBC_Decimal8_LowNibbleBorrow(t *testing.T) {
	// Low nibble borrow: 20-05 → low nibble 0-5 < 0, borrow
	cpu, _ := setupCPU(t)
	cpu.Flags.D = 1
	cpu.Flags.M = 1
	cpu.Flags.C = 1 // no borrow
	cpu.C = 0x20    // BCD 20
	err := sbc(cpu, Immediate8(0x05))
	assert.NoError(t, err)
	assert.Equal(t, uint8(0x15), cpu.A())
	assert.Equal(t, uint8(1), cpu.Flags.C)
}

func TestSBC_Decimal8_WithBorrow(t *testing.T) {
	// 10-10 with borrow (C=0): result = 10-10-1 = -1 = 99 in BCD, carry=0
	cpu, _ := setupCPU(t)
	cpu.Flags.D = 1
	cpu.Flags.M = 1
	cpu.Flags.C = 0 // borrow in
	cpu.C = 0x10    // BCD 10
	err := sbc(cpu, Immediate8(0x10))
	assert.NoError(t, err)
	assert.Equal(t, uint8(0x99), cpu.A())
	assert.Equal(t, uint8(0), cpu.Flags.C)
}

func TestSBC_Decimal16(t *testing.T) {
	cpu, _ := setupCPU(t)
	cpu.Flags.D = 1
	cpu.Flags.M = 0 // 16-bit accumulator
	cpu.Flags.C = 1 // no borrow
	cpu.C = 0x5678  // BCD 5678
	err := sbc(cpu, Immediate16(0x1234))
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x4444), cpu.C)
	assert.Equal(t, uint8(1), cpu.Flags.C)
}

// -- MVP block move --

func TestMVP_ThreeBytes(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.X = 0 // 16-bit index registers
	cpu.C = 0x0002  // 3 bytes to copy
	cpu.X = 0x1002  // start from high address (moving backwards)
	cpu.Y = 0x2002
	mem.data[0x1000] = 0xAA
	mem.data[0x1001] = 0xBB
	mem.data[0x1002] = 0xCC
	err := mvp(cpu, BlockMove{Src: 0x00, Dst: 0x00})
	assert.NoError(t, err)
	if mem.data[0x2000] != 0xAA || mem.data[0x2001] != 0xBB || mem.data[0x2002] != 0xCC {
		t.Errorf("MVP 3-byte: dst=%02X%02X%02X, want AABBCC",
			mem.data[0x2000], mem.data[0x2001], mem.data[0x2002])
	}
	assert.Equal(t, uint16(0xFFFF), cpu.C)
	assert.Equal(t, uint16(0x0FFF), cpu.X)
}

// -- RTI --

func TestRTI_Native(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.E = false
	// Manually arrange stack (lowest to highest address): P, PC_lo, PC_hi, PB
	// Stack grows down; SP points to last written byte.
	// Push order was: PB ($00), PC ($8005), P ($A5)
	cpu.SP = 0x01FB
	mem.data[0x01FC] = 0xA5 // P
	mem.data[0x01FD] = 0x05 // PC low
	mem.data[0x01FE] = 0x80 // PC high
	mem.data[0x01FF] = 0x00 // PB
	err := rti(cpu)
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x8005), cpu.PC)
	assert.Equal(t, uint8(0x00), cpu.PB)
	assert.Equal(t, uint8(0xA5), cpu.Flags.Get())
	assert.Equal(t, uint16(0x01FF), cpu.SP)
}

// -- PHB / PLB --

func TestPHB_PLB(t *testing.T) {
	cpu, _ := setupCPU(t)
	cpu.SP = 0x01FF
	cpu.DB = 0x42
	err := phb(cpu)
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x01FE), cpu.SP)
	cpu.DB = 0x00
	err = plb(cpu)
	assert.NoError(t, err)
	assert.Equal(t, uint8(0x42), cpu.DB)
	assert.Equal(t, uint16(0x01FF), cpu.SP)
}

// -- PHD / PLD --

func TestPHD_PLD(t *testing.T) {
	cpu, _ := setupCPU(t)
	cpu.SP = 0x01FF
	cpu.DP = 0x1234
	err := phd(cpu)
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x01FD), cpu.SP)
	cpu.DP = 0x0000
	err = pld(cpu)
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x1234), cpu.DP)
	assert.Equal(t, uint16(0x01FF), cpu.SP)
}

// -- WAI --

func TestWAI_HaltsAndResumes(t *testing.T) {
	cpu, mem := setupCPU(t)
	// Setup native NMI vector -> $9000
	mem.data[0xFFEA] = 0x00
	mem.data[0xFFEB] = 0x90
	cpu.PC = 0x8000
	writeOp(mem, 0x8000, 0xCB) // WAI

	// Step executes WAI and sets waiting=true
	err := cpu.Step()
	assert.NoError(t, err)
	assert.True(t, cpu.waiting)
	beforePC := cpu.PC

	// Further steps do nothing while waiting
	err = cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, beforePC, cpu.PC)

	// TriggerNMI clears waiting and queues interrupt
	cpu.TriggerNMI()
	assert.False(t, cpu.waiting)

	// CheckInterrupts dispatches the NMI (separate from Step in this emulator)
	handled := cpu.CheckInterrupts()
	assert.True(t, handled)
	assert.Equal(t, uint16(0x9000), cpu.PC)
}

// -- BRK/COP emulation mode --

func TestBRK_EmulationMode(t *testing.T) {
	cpu, mem := newTestCPU(t) // starts in emulation mode (E=true), SP=$01FF
	mem.WriteWord(VectorEmuIRQ, 0x9000)
	writeOp(mem, 0x8000, 0x00, 0x00) // BRK + signature byte

	err := cpu.Step()
	assert.NoError(t, err)

	assert.Equal(t, uint16(0x9000), cpu.PC)
	// 3-byte push: push16(PC+2) + push8(P|B) → SP goes $01FF→$01FE→$01FD→$01FC
	assert.Equal(t, uint16(0x01FC), cpu.SP)
	// Stack layout: $01FF=hi(PC+2=$8002), $01FE=lo, $01FD=P|Break
	assert.Equal(t, uint8(0x80), mem.data[0x01FF])
	assert.Equal(t, uint8(0x02), mem.data[0x01FE])
	assert.NotEqual(t, uint8(0), mem.data[0x01FD]&MaskBreak)
	assert.Equal(t, uint8(1), cpu.Flags.I)
	assert.Equal(t, uint8(0), cpu.Flags.D)
}

func TestCOP_EmulationMode(t *testing.T) {
	cpu, mem := newTestCPU(t) // starts in emulation mode (E=true), SP=$01FF
	mem.WriteWord(VectorEmuCOP, 0xA000)
	writeOp(mem, 0x8000, 0x02, 0x00) // COP + signature byte

	err := cpu.Step()
	assert.NoError(t, err)

	assert.Equal(t, uint16(0xA000), cpu.PC)
	assert.Equal(t, uint16(0x01FC), cpu.SP)
	assert.Equal(t, uint8(0x80), mem.data[0x01FF])
	assert.Equal(t, uint8(0x02), mem.data[0x01FE])
	assert.Equal(t, uint8(1), cpu.Flags.I)
	assert.Equal(t, uint8(0), cpu.Flags.D)
}

// -- Mode switch sequence --

func TestModeSwitchSequence(t *testing.T) {
	cpu, mem := newTestCPU(t) // starts in emulation mode (E=true)
	// CLC ($18) → XCE ($FB) → REP #$30 ($C2 $30) → LDA #$1234 ($A9 $34 $12)
	writeOp(mem, 0x8000, 0x18, 0xFB, 0xC2, 0x30, 0xA9, 0x34, 0x12)

	err := cpu.Step() // CLC: clear carry
	assert.NoError(t, err)
	err = cpu.Step() // XCE: E=0 (native), C=1
	assert.NoError(t, err)
	assert.False(t, cpu.E)
	err = cpu.Step() // REP #$30: clear M and X
	assert.NoError(t, err)
	assert.Equal(t, uint8(0), cpu.Flags.M)
	assert.Equal(t, uint8(0), cpu.Flags.X)
	err = cpu.Step() // LDA #$1234 in 16-bit mode
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x1234), cpu.C)
}

// -- Bank boundary crossing --

func TestAbsoluteX_CrossesBankBoundary(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.X = 0 // 16-bit index registers
	cpu.Flags.M = 1 // 8-bit accumulator
	cpu.DB = 0x01   // data bank $01
	cpu.X = 0x0100  // index $0100

	// LDA $FF00,X → eff = bank24(DB=$01, $FF00) + $0100 = $01FF00 + $0100 = $020000
	writeOp(mem, 0x8000, 0xBD, 0x00, 0xFF)
	mem.data[0x020000] = 0x55

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint8(0x55), cpu.A())
}

// -- Flags --

func TestFlagsRoundtrip(t *testing.T) {
	cpu, _ := setupCPU(t)
	p := uint8(0b10110101)
	cpu.Flags.Set(p)
	got := cpu.Flags.Get()
	assert.Equal(t, p, got)
}

// -- helpers --

func setupCPU(t *testing.T) (*CPU, *testMem) {
	t.Helper()
	cpu, mem := newTestCPU(t)
	// Switch to native mode for most tests
	cpu.E = false
	cpu.Flags.M = 1
	cpu.Flags.X = 1
	cpu.PC = 0x8000
	return cpu, mem
}

func writeOp(mem *testMem, pc uint16, bytes ...uint8) {
	for i, b := range bytes {
		mem.data[uint32(pc)+uint32(i)] = b
	}
}
