package m65816

import "testing"

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

// -- NOP --

func TestNOP(t *testing.T) {
	cpu, mem := setupCPU(t)
	writeOp(mem, 0x8000, 0xEA) // NOP
	before := cpu.cycles
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if cpu.PC != 0x8001 {
		t.Errorf("NOP: PC=%04X, want 8001", cpu.PC)
	}
	if cpu.cycles-before != 2 {
		t.Errorf("NOP: cycles=%d, want 2", cpu.cycles-before)
	}
}

// -- LDA / STA 8-bit --

func TestLDA_Immediate8(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 1
	writeOp(mem, 0x8000, 0xA9, 0x42) // LDA #$42
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if cpu.A() != 0x42 {
		t.Errorf("LDA #$42: A=%02X, want 42", cpu.A())
	}
	if cpu.PC != 0x8002 {
		t.Errorf("LDA #$42: PC=%04X, want 8002", cpu.PC)
	}
}

func TestLDA_Immediate16(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 0                        // 16-bit accumulator
	writeOp(mem, 0x8000, 0xA9, 0x34, 0x12) // LDA #$1234
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if cpu.C != 0x1234 {
		t.Errorf("LDA #$1234: C=%04X, want 1234", cpu.C)
	}
	if cpu.PC != 0x8003 {
		t.Errorf("LDA #$1234: PC=%04X, want 8003", cpu.PC)
	}
}

func TestLDA_Flags(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 1
	writeOp(mem, 0x8000, 0xA9, 0x00) // LDA #$00
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if cpu.Flags.Z != 1 {
		t.Error("LDA #0 should set Z flag")
	}
	if cpu.Flags.N != 0 {
		t.Error("LDA #0 should clear N flag")
	}

	cpu.PC = 0x8000
	writeOp(mem, 0x8000, 0xA9, 0x80) // LDA #$80
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if cpu.Flags.N != 1 {
		t.Error("LDA #$80 should set N flag")
	}
}

func TestSTA_DirectPage(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 1
	cpu.C = 0xAB
	cpu.DP = 0x0000
	writeOp(mem, 0x8000, 0x85, 0x10) // STA $10
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if mem.data[0x10] != 0xAB {
		t.Errorf("STA $10: mem[$10]=%02X, want AB", mem.data[0x10])
	}
}

// -- ADC 8-bit --

func TestADC_NoCarry(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 1
	cpu.Flags.C = 0
	cpu.C = 0x10
	writeOp(mem, 0x8000, 0x69, 0x20) // ADC #$20
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if cpu.A() != 0x30 {
		t.Errorf("ADC: A=%02X, want 30", cpu.A())
	}
	if cpu.Flags.C != 0 {
		t.Error("ADC should not set carry")
	}
}

func TestADC_WithCarry(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 1
	cpu.Flags.C = 0
	cpu.C = 0xFF
	writeOp(mem, 0x8000, 0x69, 0x01) // ADC #$01
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if cpu.A() != 0x00 {
		t.Errorf("ADC overflow: A=%02X, want 00", cpu.A())
	}
	if cpu.Flags.C != 1 {
		t.Error("ADC should set carry on overflow")
	}
	if cpu.Flags.Z != 1 {
		t.Error("ADC should set zero flag")
	}
}

// -- SBC 8-bit --

func TestSBC_Basic(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 1
	cpu.Flags.C = 1 // no borrow
	cpu.C = 0x50
	writeOp(mem, 0x8000, 0xE9, 0x20) // SBC #$20
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if cpu.A() != 0x30 {
		t.Errorf("SBC: A=%02X, want 30", cpu.A())
	}
	if cpu.Flags.C != 1 {
		t.Error("SBC no borrow: C should be 1")
	}
}

// -- INC/DEC accumulator --

func TestINC_Accumulator(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 1
	cpu.C = 0x41
	writeOp(mem, 0x8000, 0x1A) // INC A
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if cpu.A() != 0x42 {
		t.Errorf("INC A: A=%02X, want 42", cpu.A())
	}
}

func TestDEC_Accumulator(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 1
	cpu.C = 0x43
	writeOp(mem, 0x8000, 0x3A) // DEC A
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if cpu.A() != 0x42 {
		t.Errorf("DEC A: A=%02X, want 42", cpu.A())
	}
}

// -- ASL --

func TestASL_Accumulator8(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 1
	cpu.C = 0x41
	writeOp(mem, 0x8000, 0x0A) // ASL A
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if cpu.A() != 0x82 {
		t.Errorf("ASL A: A=%02X, want 82", cpu.A())
	}
	if cpu.Flags.C != 0 {
		t.Error("ASL should not set carry for $41")
	}
}

func TestASL_SetsCarry(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 1
	cpu.C = 0x80
	writeOp(mem, 0x8000, 0x0A) // ASL A
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if cpu.Flags.C != 1 {
		t.Error("ASL $80 should set carry")
	}
	if cpu.A() != 0x00 {
		t.Errorf("ASL $80: A=%02X, want 00", cpu.A())
	}
}

// -- AND / ORA / EOR --

func TestAND_Immediate(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 1
	cpu.C = 0xFF
	writeOp(mem, 0x8000, 0x29, 0x0F) // AND #$0F
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if cpu.A() != 0x0F {
		t.Errorf("AND: A=%02X, want 0F", cpu.A())
	}
}

func TestORA_Immediate(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 1
	cpu.C = 0x00
	writeOp(mem, 0x8000, 0x09, 0xF0) // ORA #$F0
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if cpu.A() != 0xF0 {
		t.Errorf("ORA: A=%02X, want F0", cpu.A())
	}
}

func TestEOR_Immediate(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 1
	cpu.C = 0xFF
	writeOp(mem, 0x8000, 0x49, 0x0F) // EOR #$0F
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if cpu.A() != 0xF0 {
		t.Errorf("EOR: A=%02X, want F0", cpu.A())
	}
}

// -- CMP --

func TestCMP_Equal(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 1
	cpu.C = 0x42
	writeOp(mem, 0x8000, 0xC9, 0x42) // CMP #$42
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if cpu.Flags.Z != 1 {
		t.Error("CMP equal: Z should be 1")
	}
	if cpu.Flags.C != 1 {
		t.Error("CMP equal: C should be 1")
	}
	if cpu.Flags.N != 0 {
		t.Error("CMP equal: N should be 0")
	}
}

// -- Branch instructions --

func TestBRA_Taken(t *testing.T) {
	cpu, mem := setupCPU(t)
	writeOp(mem, 0x8000, 0x80, 0x05) // BRA +5
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	// target = 0x8000 + 2 + 5 = 0x8007
	if cpu.PC != 0x8007 {
		t.Errorf("BRA: PC=%04X, want 8007", cpu.PC)
	}
}

func TestBNE_NotTaken(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.Z = 1                  // Z set means equal, BNE not taken
	writeOp(mem, 0x8000, 0xD0, 0x05) // BNE +5
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if cpu.PC != 0x8002 {
		t.Errorf("BNE (not taken): PC=%04X, want 8002", cpu.PC)
	}
}

func TestBNE_Taken(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.Z = 0
	writeOp(mem, 0x8000, 0xD0, 0x05) // BNE +5
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if cpu.PC != 0x8007 {
		t.Errorf("BNE (taken): PC=%04X, want 8007", cpu.PC)
	}
}

// -- Stack operations --

func TestPHA_PLA_8bit(t *testing.T) {
	cpu, _ := setupCPU(t)
	cpu.Flags.M = 1
	cpu.SP = 0x01FF
	cpu.C = 0x55
	if err := pha(cpu); err != nil {
		t.Fatal(err)
	}
	cpu.C = 0
	if err := pla(cpu); err != nil {
		t.Fatal(err)
	}
	if cpu.A() != 0x55 {
		t.Errorf("PHA/PLA: A=%02X, want 55", cpu.A())
	}
}

func TestPHX_PLX_8bit(t *testing.T) {
	cpu, _ := setupCPU(t)
	cpu.Flags.X = 1
	cpu.SP = 0x01FF
	cpu.X = 0x77
	if err := phx(cpu); err != nil {
		t.Fatal(err)
	}
	cpu.X = 0
	if err := plx(cpu); err != nil {
		t.Fatal(err)
	}
	if uint8(cpu.X) != 0x77 {
		t.Errorf("PHX/PLX: X=%02X, want 77", cpu.X)
	}
}

// -- Register transfers --

func TestTAX_8bit(t *testing.T) {
	cpu, _ := setupCPU(t)
	cpu.Flags.X = 1
	cpu.C = 0x42
	if err := tax(cpu); err != nil {
		t.Fatal(err)
	}
	if uint8(cpu.X) != 0x42 {
		t.Errorf("TAX: X=%02X, want 42", cpu.X)
	}
}

func TestTCD_TDC(t *testing.T) {
	cpu, _ := setupCPU(t)
	cpu.C = 0x1234
	if err := tcd(cpu); err != nil {
		t.Fatal(err)
	}
	if cpu.DP != 0x1234 {
		t.Errorf("TCD: DP=%04X, want 1234", cpu.DP)
	}
	cpu.C = 0
	if err := tdc(cpu); err != nil {
		t.Fatal(err)
	}
	if cpu.C != 0x1234 {
		t.Errorf("TDC: C=%04X, want 1234", cpu.C)
	}
}

func TestTCS_TSC(t *testing.T) {
	cpu, _ := setupCPU(t)
	cpu.E = false
	cpu.C = 0x0150
	if err := tcs(cpu); err != nil {
		t.Fatal(err)
	}
	if cpu.SP != 0x0150 {
		t.Errorf("TCS: SP=%04X, want 0150", cpu.SP)
	}
	cpu.C = 0
	if err := tsc(cpu); err != nil {
		t.Fatal(err)
	}
	if cpu.C != 0x0150 {
		t.Errorf("TSC: C=%04X, want 0150", cpu.C)
	}
}

// -- Processor status --

func TestSEP_REP(t *testing.T) {
	cpu, _ := setupCPU(t)
	cpu.E = false
	cpu.Flags.M = 0
	cpu.Flags.X = 0
	// SEP #$30 sets M and X flags
	if err := sep(cpu, Immediate8(0x30)); err != nil {
		t.Fatal(err)
	}
	if cpu.Flags.M != 1 {
		t.Error("SEP #$30 should set M flag")
	}
	if cpu.Flags.X != 1 {
		t.Error("SEP #$30 should set X flag")
	}
	// REP #$30 clears M and X flags
	if err := rep(cpu, Immediate8(0x30)); err != nil {
		t.Fatal(err)
	}
	if cpu.Flags.M != 0 {
		t.Error("REP #$30 should clear M flag")
	}
	if cpu.Flags.X != 0 {
		t.Error("REP #$30 should clear X flag")
	}
}

func TestXCE_ToNative(t *testing.T) {
	cpu, _ := newTestCPU(t)
	// Start in emulation mode (E=true)
	if !cpu.E {
		t.Fatal("expected emulation mode")
	}
	cpu.Flags.C = 0 // carry = 0
	if err := xce(cpu); err != nil {
		t.Fatal(err)
	}
	if cpu.E {
		t.Error("XCE with C=0 should switch to native mode")
	}
	if cpu.Flags.C != 1 {
		t.Error("XCE should set carry to old emulation flag (was 1)")
	}
}

func TestXBA(t *testing.T) {
	cpu, _ := setupCPU(t)
	cpu.C = 0x1234
	if err := xba(cpu); err != nil {
		t.Fatal(err)
	}
	if cpu.C != 0x3412 {
		t.Errorf("XBA: C=%04X, want 3412", cpu.C)
	}
}

// -- STP / WAI --

func TestSTP(t *testing.T) {
	cpu, _ := setupCPU(t)
	if err := stp(cpu); err != nil {
		t.Fatal(err)
	}
	if !cpu.stopped {
		t.Error("STP should set stopped flag")
	}
	// Step should be no-op when stopped
	before := cpu.PC
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if cpu.PC != before {
		t.Error("Step when stopped should not advance PC")
	}
}

// -- 16-bit accumulator --

func TestLDA16_STA16(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 0 // 16-bit accumulator
	cpu.DP = 0x0000

	// LDA #$ABCD
	writeOp(mem, 0x8000, 0xA9, 0xCD, 0xAB)
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if cpu.C != 0xABCD {
		t.Errorf("LDA16: C=%04X, want ABCD", cpu.C)
	}

	// STA $20 (direct page)
	cpu.PC = 0x8003
	writeOp(mem, 0x8003, 0x85, 0x20)
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
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
	if err := inx(cpu); err != nil {
		t.Fatal(err)
	}
	if cpu.X != 0xFFFF {
		t.Errorf("INX16: X=%04X, want FFFF", cpu.X)
	}
}

func TestDEX16(t *testing.T) {
	cpu, _ := setupCPU(t)
	cpu.Flags.X = 0
	cpu.X = 0x0001
	if err := dex(cpu); err != nil {
		t.Fatal(err)
	}
	if cpu.X != 0x0000 {
		t.Errorf("DEX16: X=%04X, want 0000", cpu.X)
	}
	if cpu.Flags.Z != 1 {
		t.Error("DEX to zero should set Z flag")
	}
}

// -- ADC/SBC 16-bit --

func TestADC_16bit(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 0 // 16-bit accumulator
	cpu.Flags.C = 0
	cpu.C = 0x1000
	writeOp(mem, 0x8000, 0x69, 0x34, 0x12) // ADC #$1234
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if cpu.C != 0x2234 {
		t.Errorf("ADC16: C=%04X, want 2234", cpu.C)
	}
	if cpu.Flags.C != 0 {
		t.Error("ADC16: unexpected carry")
	}
}

func TestADC_16bit_Carry(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 0
	cpu.Flags.C = 0
	cpu.C = 0xFF00
	writeOp(mem, 0x8000, 0x69, 0x00, 0x01) // ADC #$0100
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if cpu.C != 0x0000 {
		t.Errorf("ADC16 carry: C=%04X, want 0000", cpu.C)
	}
	if cpu.Flags.C != 1 {
		t.Error("ADC16 carry: expected carry set")
	}
	if cpu.Flags.Z != 1 {
		t.Error("ADC16 carry: expected zero flag")
	}
}

func TestSBC_16bit(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 0
	cpu.Flags.C = 1 // no borrow
	cpu.C = 0x1234
	writeOp(mem, 0x8000, 0xE9, 0x34, 0x02) // SBC #$0234
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if cpu.C != 0x1000 {
		t.Errorf("SBC16: C=%04X, want 1000", cpu.C)
	}
	if cpu.Flags.C != 1 {
		t.Error("SBC16: no-borrow C should be 1")
	}
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
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if mem.data[0x0030] != 0xBB {
		t.Errorf("STA dp,X 8bit: mem[0030]=%02X, want BB", mem.data[0x0030])
	}
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
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if mem.data[0x01F0] != 0xAA {
		t.Errorf("STA dp,X 16bit: mem[01F0]=%02X, want AA", mem.data[0x01F0])
	}
	// Verify wrong address ($00F0) was NOT written (old bug would have written there)
	if mem.data[0x00F0] != 0x00 {
		t.Errorf("STA dp,X 16bit: wrote to wrong address 00F0 (old bug)")
	}
}

func TestSTADP_IndexedX_EmulationWrap(t *testing.T) {
	// In emulation mode with DP=$0000, (dp+X) wraps within page 0
	cpu, mem := newTestCPU(t)
	// Stays in emulation mode (E=true)
	cpu.DP = 0x0000
	cpu.C = 0x77
	cpu.X = 0x20
	writeOp(mem, 0x8000, 0x95, 0xF0) // STA $F0,X -> (0xF0+0x20)&0xFF = 0x10
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if mem.data[0x0010] != 0x77 {
		t.Errorf("STA dp,X emu wrap: mem[0010]=%02X, want 77", mem.data[0x0010])
	}
}

// -- JSR / RTS --

func TestJSR_RTS(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.SP = 0x01FF
	// JSR $9000 at $8000; subroutine at $9000 contains RTS
	writeOp(mem, 0x8000, 0x20, 0x00, 0x90) // JSR $9000
	writeOp(mem, 0x9000, 0x60)             // RTS
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if cpu.PC != 0x9000 {
		t.Errorf("JSR: PC=%04X, want 9000", cpu.PC)
	}
	if cpu.SP != 0x01FD {
		t.Errorf("JSR: SP=%04X, want 01FD", cpu.SP)
	}
	// Return address $8002 pushed: high byte ($80) at $01FF, low byte ($02) at $01FE
	if mem.data[0x01FE] != 0x02 || mem.data[0x01FF] != 0x80 {
		t.Errorf("JSR: stack hi=%02X lo=%02X, want hi=80 lo=02", mem.data[0x01FF], mem.data[0x01FE])
	}

	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if cpu.PC != 0x8003 {
		t.Errorf("RTS: PC=%04X, want 8003", cpu.PC)
	}
	if cpu.SP != 0x01FF {
		t.Errorf("RTS: SP=%04X, want 01FF", cpu.SP)
	}
}

// -- JSL / RTL --

func TestJSL_RTL(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.SP = 0x01FF
	cpu.PB = 0x00
	// JSL $019000 at $8000; subroutine at bank $01:$9000 contains RTL
	writeOp(mem, 0x8000, 0x22, 0x00, 0x90, 0x01) // JSL $019000
	mem.data[0x019000] = 0x6B                    // RTL (in bank $01)
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if cpu.PC != 0x9000 {
		t.Errorf("JSL: PC=%04X, want 9000", cpu.PC)
	}
	if cpu.PB != 0x01 {
		t.Errorf("JSL: PB=%02X, want 01", cpu.PB)
	}
	if cpu.SP != 0x01FC {
		t.Errorf("JSL: SP=%04X, want 01FC", cpu.SP)
	}
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if cpu.PC != 0x8004 {
		t.Errorf("RTL: PC=%04X, want 8004", cpu.PC)
	}
	if cpu.PB != 0x00 {
		t.Errorf("RTL: PB=%02X, want 00", cpu.PB)
	}
	if cpu.SP != 0x01FF {
		t.Errorf("RTL: SP=%04X, want 01FF", cpu.SP)
	}
}

// -- JMP / JML --

func TestJMP_Absolute(t *testing.T) {
	cpu, mem := setupCPU(t)
	writeOp(mem, 0x8000, 0x4C, 0x34, 0x12) // JMP $1234
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if cpu.PC != 0x1234 {
		t.Errorf("JMP: PC=%04X, want 1234", cpu.PC)
	}
}

func TestJML_Long(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.PB = 0x00
	writeOp(mem, 0x8000, 0x5C, 0x00, 0x90, 0x01) // JML $019000
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if cpu.PC != 0x9000 {
		t.Errorf("JML: PC=%04X, want 9000", cpu.PC)
	}
	if cpu.PB != 0x01 {
		t.Errorf("JML: PB=%02X, want 01", cpu.PB)
	}
}

// -- PEA / PEI / PER --

func TestPEA(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.SP = 0x01FF
	writeOp(mem, 0x8000, 0xF4, 0x34, 0x12) // PEA $1234
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if cpu.SP != 0x01FD {
		t.Errorf("PEA: SP=%04X, want 01FD", cpu.SP)
	}
	// $1234 pushed: lo=$34 at $01FE, hi=$12 at $01FF
	if mem.data[0x01FE] != 0x34 || mem.data[0x01FF] != 0x12 {
		t.Errorf("PEA: stack=%02X%02X, want 1234", mem.data[0x01FF], mem.data[0x01FE])
	}
	if cpu.PC != 0x8003 {
		t.Errorf("PEA: PC=%04X, want 8003", cpu.PC)
	}
}

func TestPEI(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.SP = 0x01FF
	cpu.DP = 0x0000
	// Store $5678 at direct page offset $10
	mem.data[0x0010] = 0x78
	mem.data[0x0011] = 0x56
	writeOp(mem, 0x8000, 0xD4, 0x10) // PEI ($10)
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if cpu.SP != 0x01FD {
		t.Errorf("PEI: SP=%04X, want 01FD", cpu.SP)
	}
	got := uint16(mem.data[0x01FF])<<8 | uint16(mem.data[0x01FE])
	if got != 0x5678 {
		t.Errorf("PEI: pushed=%04X, want 5678", got)
	}
	if cpu.PC != 0x8002 {
		t.Errorf("PEI: PC=%04X, want 8002", cpu.PC)
	}
}

func TestPER(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.SP = 0x01FF
	// PER +$0100: effective = 0x8000 + 3 + 0x0100 = 0x8103
	writeOp(mem, 0x8000, 0x62, 0x00, 0x01) // PER +$0100
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if cpu.SP != 0x01FD {
		t.Errorf("PER: SP=%04X, want 01FD", cpu.SP)
	}
	got := uint16(mem.data[0x01FF])<<8 | uint16(mem.data[0x01FE])
	if got != 0x8103 {
		t.Errorf("PER: pushed=%04X, want 8103", got)
	}
	if cpu.PC != 0x8003 {
		t.Errorf("PER: PC=%04X, want 8003", cpu.PC)
	}
}

// -- MVN block move --

func TestMVN_SingleByte(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.C = 0x0000 // 1 byte to copy (C+1)
	cpu.X = 0x1000
	cpu.Y = 0x2000
	mem.data[0x1000] = 0xAA
	if err := mvn(cpu, BlockMove{Src: 0x00, Dst: 0x00}); err != nil {
		t.Fatal(err)
	}
	if mem.data[0x2000] != 0xAA {
		t.Errorf("MVN 1-byte: dst=%02X, want AA", mem.data[0x2000])
	}
	if cpu.C != 0xFFFF {
		t.Errorf("MVN 1-byte: C=%04X, want FFFF", cpu.C)
	}
	if cpu.X != 0x1001 {
		t.Errorf("MVN 1-byte: X=%04X, want 1001", cpu.X)
	}
	if cpu.Y != 0x2001 {
		t.Errorf("MVN 1-byte: Y=%04X, want 2001", cpu.Y)
	}
}

func TestMVN_ThreeBytes(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.C = 0x0002 // 3 bytes to copy
	cpu.X = 0x1000
	cpu.Y = 0x2000
	mem.data[0x1000] = 0xAA
	mem.data[0x1001] = 0xBB
	mem.data[0x1002] = 0xCC
	if err := mvn(cpu, BlockMove{Src: 0x00, Dst: 0x00}); err != nil {
		t.Fatal(err)
	}
	if mem.data[0x2000] != 0xAA || mem.data[0x2001] != 0xBB || mem.data[0x2002] != 0xCC {
		t.Errorf("MVN 3-byte: dst=%02X%02X%02X, want AABBCC",
			mem.data[0x2000], mem.data[0x2001], mem.data[0x2002])
	}
	if cpu.C != 0xFFFF {
		t.Errorf("MVN 3-byte: C=%04X, want FFFF", cpu.C)
	}
	if cpu.X != 0x1003 {
		t.Errorf("MVN 3-byte: X=%04X, want 1003", cpu.X)
	}
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
	if err := cpu.Step(); err != nil {
		t.Fatal(err)
	}
	if cpu.PC != 0x9000 {
		t.Errorf("BRK: PC=%04X, want 9000", cpu.PC)
	}
	if cpu.PB != 0x00 {
		t.Errorf("BRK: PB=%02X, want 00", cpu.PB)
	}
	if cpu.Flags.I != 1 {
		t.Error("BRK: I flag should be set")
	}
	if cpu.Flags.D != 0 {
		t.Error("BRK: D flag should be cleared")
	}
	if cpu.SP != 0x01FB {
		t.Errorf("BRK: SP=%04X, want 01FB (pushed PB+retAddr+P)", cpu.SP)
	}
}

// -- Flags --

func TestFlagsRoundtrip(t *testing.T) {
	cpu, _ := setupCPU(t)
	p := uint8(0b10110101)
	cpu.Flags.Set(p)
	got := cpu.Flags.Get()
	if got != p {
		t.Errorf("Flags roundtrip: got %08b, want %08b", got, p)
	}
}
