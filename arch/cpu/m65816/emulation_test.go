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
	cpu.Flags.M = 0 // 16-bit accumulator
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
	cpu.Flags.Z = 1 // Z set means equal, BNE not taken
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
