package sm83

// OpcodeInfo contains the encoding details for a specific instruction variant.
type OpcodeInfo struct {
	Prefix byte // Prefix byte (0x00 for none, 0xCB for CB-prefixed)
	Opcode byte // Opcode byte (after prefix if applicable)
	Size   byte // Size of opcode in bytes
	Cycles byte // Timing in M-cycles (1 M-cycle = 4 T-states)
}

// Instruction defines an SM83 CPU instruction with its opcodes and execution logic.
// Instructions support multiple addressing modes and register variants through three
// opcode mapping fields that enable bidirectional opcode/instruction lookup.
type Instruction struct {
	Name       string // Instruction mnemonic (lowercase)
	Unofficial bool   // True for undocumented opcodes

	// Opcode lookup maps for bidirectional instruction/opcode mapping
	Addressing          map[AddressingMode]OpcodeInfo   // Maps addressing mode to opcode
	RegisterOpcodes     map[RegisterParam]OpcodeInfo    // Maps single register parameter to opcode
	RegisterPairOpcodes map[[2]RegisterParam]OpcodeInfo // Maps register pairs to opcode (e.g., LD r,r')

	// Execution handlers - exactly one must be set
	NoParamFunc func(c *CPU) error                // Handler for implied addressing
	ParamFunc   func(c *CPU, params ...any) error // Handler for parameterized instructions
}

// HasAddressing returns whether the instruction has any of the passed addressing modes.
func (ins Instruction) HasAddressing(flags ...AddressingMode) bool {
	for _, flag := range flags {
		_, ok := ins.Addressing[flag]
		if ok {
			return ok
		}
	}
	return false
}

// GetOpcodeByRegister returns opcode info for a specific register parameter.
func (ins Instruction) GetOpcodeByRegister(register RegisterParam) (OpcodeInfo, bool) {
	if ins.RegisterOpcodes == nil {
		// Fall back to Addressing map if RegisterOpcodes is not defined
		for _, info := range ins.Addressing {
			return info, true
		}
		return OpcodeInfo{}, false
	}

	info, exists := ins.RegisterOpcodes[register]
	return info, exists
}

// GetOpcodeByRegisterPair returns opcode info for a pair of register parameters.
// Used for instructions like LD r,r' that require both source and destination registers.
func (ins Instruction) GetOpcodeByRegisterPair(dst, src RegisterParam) (OpcodeInfo, bool) {
	if ins.RegisterPairOpcodes == nil {
		return OpcodeInfo{}, false
	}

	info, exists := ins.RegisterPairOpcodes[[2]RegisterParam{dst, src}]
	return info, exists
}

// Instruction name constants for easy access by external packages.
const (
	AdcName  = "adc"
	AddName  = "add"
	AndName  = "and"
	BitName  = "bit"
	CallName = "call"
	CcfName  = "ccf"
	CpName   = "cp"
	CplName  = "cpl"
	DaaName  = "daa"
	DecName  = "dec"
	DiName   = "di"
	EiName   = "ei"
	HaltName = "halt"
	IncName  = "inc"
	JpName   = "jp"
	JrName   = "jr"
	LdName   = "ld"
	LdhName  = "ldh"
	NopName  = "nop"
	OrName   = "or"
	PopName  = "pop"
	PushName = "push"
	ResName  = "res"
	RetName  = "ret"
	RetiName = "reti"
	RlName   = "rl"
	RlaName  = "rla"
	RlcName  = "rlc"
	RlcaName = "rlca"
	RrName   = "rr"
	RraName  = "rra"
	RrcName  = "rrc"
	RrcaName = "rrca"
	RstName  = "rst"
	SbcName  = "sbc"
	ScfName  = "scf"
	SetName  = "set"
	SlaName  = "sla"
	SraName  = "sra"
	SrlName  = "srl"
	StopName = "stop"
	SubName  = "sub"
	SwapName = "swap"
	XorName  = "xor"
)

// ---------------------------------------------------------------------------
// Implied / simple instructions
// ---------------------------------------------------------------------------

// NopInst - No Operation.
var NopInst = &Instruction{
	Name: NopName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x00, Size: 1, Cycles: 1},
	},
	NoParamFunc: nop,
}

// HaltInst - Halt execution until interrupt.
var HaltInst = &Instruction{
	Name: HaltName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x76, Size: 1, Cycles: 1},
	},
	NoParamFunc: halt,
}

// StopInst - Stop CPU and LCD until button press (SM83-unique).
var StopInst = &Instruction{
	Name: StopName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x10, Size: 2, Cycles: 1},
	},
	NoParamFunc: stop,
}

// RlcaInst - Rotate Left Circular Accumulator.
var RlcaInst = &Instruction{
	Name: RlcaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x07, Size: 1, Cycles: 1},
	},
	NoParamFunc: rlcaFunc,
}

// RrcaInst - Rotate Right Circular Accumulator.
var RrcaInst = &Instruction{
	Name: RrcaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x0F, Size: 1, Cycles: 1},
	},
	NoParamFunc: rrcaFunc,
}

// RlaInst - Rotate Left Accumulator through carry.
var RlaInst = &Instruction{
	Name: RlaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x17, Size: 1, Cycles: 1},
	},
	NoParamFunc: rlaFunc,
}

// RraInst - Rotate Right Accumulator through carry.
var RraInst = &Instruction{
	Name: RraName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x1F, Size: 1, Cycles: 1},
	},
	NoParamFunc: rraFunc,
}

// DaaInst - Decimal Adjust Accumulator.
var DaaInst = &Instruction{
	Name: DaaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x27, Size: 1, Cycles: 1},
	},
	NoParamFunc: daa,
}

// CplInst - Complement Accumulator (bitwise NOT).
var CplInst = &Instruction{
	Name: CplName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x2F, Size: 1, Cycles: 1},
	},
	NoParamFunc: cpl,
}

// ScfInst - Set Carry Flag.
var ScfInst = &Instruction{
	Name: ScfName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x37, Size: 1, Cycles: 1},
	},
	NoParamFunc: scf,
}

// CcfInst - Complement Carry Flag.
var CcfInst = &Instruction{
	Name: CcfName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x3F, Size: 1, Cycles: 1},
	},
	NoParamFunc: ccf,
}

// DiInst - Disable Interrupts.
var DiInst = &Instruction{
	Name: DiName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xF3, Size: 1, Cycles: 1},
	},
	NoParamFunc: di,
}

// EiInst - Enable Interrupts.
var EiInst = &Instruction{
	Name: EiName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xFB, Size: 1, Cycles: 1},
	},
	NoParamFunc: ei,
}

// RetiInst - Return from Interrupt (SM83: opcode 0xD9, enables interrupts).
var RetiInst = &Instruction{
	Name: RetiName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xD9, Size: 1, Cycles: 4},
	},
	NoParamFunc: reti,
}

// ---------------------------------------------------------------------------
// 8-bit register loads
// ---------------------------------------------------------------------------

// LdImm8 - Load 8-bit immediate into register (LD r,n).
var LdImm8 = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x3E, Size: 2, Cycles: 2},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB: {Opcode: 0x06, Size: 2, Cycles: 2}, // LD B,n
		RegC: {Opcode: 0x0E, Size: 2, Cycles: 2}, // LD C,n
		RegD: {Opcode: 0x16, Size: 2, Cycles: 2}, // LD D,n
		RegE: {Opcode: 0x1E, Size: 2, Cycles: 2}, // LD E,n
		RegH: {Opcode: 0x26, Size: 2, Cycles: 2}, // LD H,n
		RegL: {Opcode: 0x2E, Size: 2, Cycles: 2}, // LD L,n
		RegA: {Opcode: 0x3E, Size: 2, Cycles: 2}, // LD A,n
	},
	ParamFunc: ldImm8,
}

// LdReg8 - Load between 8-bit registers (LD r,r').
var LdReg8 = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Opcode: 0x7F, Size: 1, Cycles: 1},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegLoadHLB: {Opcode: 0x46, Size: 1, Cycles: 2}, // LD B,(HL)
		RegLoadHLC: {Opcode: 0x4E, Size: 1, Cycles: 2}, // LD C,(HL)
		RegLoadHLD: {Opcode: 0x56, Size: 1, Cycles: 2}, // LD D,(HL)
		RegLoadHLE: {Opcode: 0x5E, Size: 1, Cycles: 2}, // LD E,(HL)
		RegLoadHLH: {Opcode: 0x66, Size: 1, Cycles: 2}, // LD H,(HL)
		RegLoadHLL: {Opcode: 0x6E, Size: 1, Cycles: 2}, // LD L,(HL)
		RegLoadHLA: {Opcode: 0x7E, Size: 1, Cycles: 2}, // LD A,(HL)
	},
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegB, RegB}:          {Opcode: 0x40, Size: 1, Cycles: 1}, // LD B,B
		{RegB, RegC}:          {Opcode: 0x41, Size: 1, Cycles: 1}, // LD B,C
		{RegB, RegD}:          {Opcode: 0x42, Size: 1, Cycles: 1}, // LD B,D
		{RegB, RegE}:          {Opcode: 0x43, Size: 1, Cycles: 1}, // LD B,E
		{RegB, RegH}:          {Opcode: 0x44, Size: 1, Cycles: 1}, // LD B,H
		{RegB, RegL}:          {Opcode: 0x45, Size: 1, Cycles: 1}, // LD B,L
		{RegB, RegA}:          {Opcode: 0x47, Size: 1, Cycles: 1}, // LD B,A
		{RegC, RegB}:          {Opcode: 0x48, Size: 1, Cycles: 1}, // LD C,B
		{RegC, RegC}:          {Opcode: 0x49, Size: 1, Cycles: 1}, // LD C,C
		{RegC, RegD}:          {Opcode: 0x4A, Size: 1, Cycles: 1}, // LD C,D
		{RegC, RegE}:          {Opcode: 0x4B, Size: 1, Cycles: 1}, // LD C,E
		{RegC, RegH}:          {Opcode: 0x4C, Size: 1, Cycles: 1}, // LD C,H
		{RegC, RegL}:          {Opcode: 0x4D, Size: 1, Cycles: 1}, // LD C,L
		{RegC, RegA}:          {Opcode: 0x4F, Size: 1, Cycles: 1}, // LD C,A
		{RegD, RegB}:          {Opcode: 0x50, Size: 1, Cycles: 1}, // LD D,B
		{RegD, RegC}:          {Opcode: 0x51, Size: 1, Cycles: 1}, // LD D,C
		{RegD, RegD}:          {Opcode: 0x52, Size: 1, Cycles: 1}, // LD D,D
		{RegD, RegE}:          {Opcode: 0x53, Size: 1, Cycles: 1}, // LD D,E
		{RegD, RegH}:          {Opcode: 0x54, Size: 1, Cycles: 1}, // LD D,H
		{RegD, RegL}:          {Opcode: 0x55, Size: 1, Cycles: 1}, // LD D,L
		{RegD, RegA}:          {Opcode: 0x57, Size: 1, Cycles: 1}, // LD D,A
		{RegE, RegB}:          {Opcode: 0x58, Size: 1, Cycles: 1}, // LD E,B
		{RegE, RegC}:          {Opcode: 0x59, Size: 1, Cycles: 1}, // LD E,C
		{RegE, RegD}:          {Opcode: 0x5A, Size: 1, Cycles: 1}, // LD E,D
		{RegE, RegE}:          {Opcode: 0x5B, Size: 1, Cycles: 1}, // LD E,E
		{RegE, RegH}:          {Opcode: 0x5C, Size: 1, Cycles: 1}, // LD E,H
		{RegE, RegL}:          {Opcode: 0x5D, Size: 1, Cycles: 1}, // LD E,L
		{RegE, RegA}:          {Opcode: 0x5F, Size: 1, Cycles: 1}, // LD E,A
		{RegH, RegB}:          {Opcode: 0x60, Size: 1, Cycles: 1}, // LD H,B
		{RegH, RegC}:          {Opcode: 0x61, Size: 1, Cycles: 1}, // LD H,C
		{RegH, RegD}:          {Opcode: 0x62, Size: 1, Cycles: 1}, // LD H,D
		{RegH, RegE}:          {Opcode: 0x63, Size: 1, Cycles: 1}, // LD H,E
		{RegH, RegH}:          {Opcode: 0x64, Size: 1, Cycles: 1}, // LD H,H
		{RegH, RegL}:          {Opcode: 0x65, Size: 1, Cycles: 1}, // LD H,L
		{RegH, RegA}:          {Opcode: 0x67, Size: 1, Cycles: 1}, // LD H,A
		{RegL, RegB}:          {Opcode: 0x68, Size: 1, Cycles: 1}, // LD L,B
		{RegL, RegC}:          {Opcode: 0x69, Size: 1, Cycles: 1}, // LD L,C
		{RegL, RegD}:          {Opcode: 0x6A, Size: 1, Cycles: 1}, // LD L,D
		{RegL, RegE}:          {Opcode: 0x6B, Size: 1, Cycles: 1}, // LD L,E
		{RegL, RegH}:          {Opcode: 0x6C, Size: 1, Cycles: 1}, // LD L,H
		{RegL, RegL}:          {Opcode: 0x6D, Size: 1, Cycles: 1}, // LD L,L
		{RegL, RegA}:          {Opcode: 0x6F, Size: 1, Cycles: 1}, // LD L,A
		{RegA, RegB}:          {Opcode: 0x78, Size: 1, Cycles: 1}, // LD A,B
		{RegA, RegC}:          {Opcode: 0x79, Size: 1, Cycles: 1}, // LD A,C
		{RegA, RegD}:          {Opcode: 0x7A, Size: 1, Cycles: 1}, // LD A,D
		{RegA, RegE}:          {Opcode: 0x7B, Size: 1, Cycles: 1}, // LD A,E
		{RegA, RegH}:          {Opcode: 0x7C, Size: 1, Cycles: 1}, // LD A,H
		{RegA, RegL}:          {Opcode: 0x7D, Size: 1, Cycles: 1}, // LD A,L
		{RegA, RegA}:          {Opcode: 0x7F, Size: 1, Cycles: 1}, // LD A,A
		{RegHLIndirect, RegB}: {Opcode: 0x70, Size: 1, Cycles: 2}, // LD (HL),B
		{RegHLIndirect, RegC}: {Opcode: 0x71, Size: 1, Cycles: 2}, // LD (HL),C
		{RegHLIndirect, RegD}: {Opcode: 0x72, Size: 1, Cycles: 2}, // LD (HL),D
		{RegHLIndirect, RegE}: {Opcode: 0x73, Size: 1, Cycles: 2}, // LD (HL),E
		{RegHLIndirect, RegH}: {Opcode: 0x74, Size: 1, Cycles: 2}, // LD (HL),H
		{RegHLIndirect, RegL}: {Opcode: 0x75, Size: 1, Cycles: 2}, // LD (HL),L
		{RegHLIndirect, RegA}: {Opcode: 0x77, Size: 1, Cycles: 2}, // LD (HL),A
	},
	ParamFunc: ldReg8,
}

// ---------------------------------------------------------------------------
// 16-bit register loads
// ---------------------------------------------------------------------------

// LdReg16 - Load 16-bit immediate into register pair (LD rr,nn).
var LdReg16 = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x01, Size: 3, Cycles: 3},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegBC: {Opcode: 0x01, Size: 3, Cycles: 3}, // LD BC,nn
		RegDE: {Opcode: 0x11, Size: 3, Cycles: 3}, // LD DE,nn
		RegHL: {Opcode: 0x21, Size: 3, Cycles: 3}, // LD HL,nn
		RegSP: {Opcode: 0x31, Size: 3, Cycles: 3}, // LD SP,nn
	},
	ParamFunc: ldReg16,
}

// ---------------------------------------------------------------------------
// Indirect loads
// ---------------------------------------------------------------------------

// LdIndirect - Load indirect through BC/DE (LD (BC),A / LD A,(BC) / LD (DE),A / LD A,(DE)).
var LdIndirect = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Opcode: 0x02, Size: 1, Cycles: 2},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegBCIndirect: {Opcode: 0x02, Size: 1, Cycles: 2}, // LD (BC),A
		RegLoadBC:     {Opcode: 0x0A, Size: 1, Cycles: 2}, // LD A,(BC)
		RegDEIndirect: {Opcode: 0x12, Size: 1, Cycles: 2}, // LD (DE),A
		RegLoadDE:     {Opcode: 0x1A, Size: 1, Cycles: 2}, // LD A,(DE)
	},
	ParamFunc: ldIndirect,
}

// ---------------------------------------------------------------------------
// SM83-specific load instructions
// ---------------------------------------------------------------------------

// LdHLPlusA - LD (HL+),A — Store A to (HL), then increment HL.
var LdHLPlusA = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x22, Size: 1, Cycles: 2},
	},
	NoParamFunc: ldHLPlusA,
}

// LdAHLPlus - LD A,(HL+) — Load (HL) into A, then increment HL.
var LdAHLPlus = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x2A, Size: 1, Cycles: 2},
	},
	NoParamFunc: ldAHLPlus,
}

// LdHLMinusA - LD (HL-),A — Store A to (HL), then decrement HL.
var LdHLMinusA = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x32, Size: 1, Cycles: 2},
	},
	NoParamFunc: ldHLMinusA,
}

// LdAHLMinus - LD A,(HL-) — Load (HL) into A, then decrement HL.
var LdAHLMinus = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x3A, Size: 1, Cycles: 2},
	},
	NoParamFunc: ldAHLMinus,
}

// LdAddrSP - LD (nn),SP — Store SP to address nn (SM83-unique at 0x08).
var LdAddrSP = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ExtendedAddressing: {Opcode: 0x08, Size: 3, Cycles: 5},
	},
	ParamFunc: ldAddrSP,
}

// LdSPHL - LD SP,HL — Load HL into SP.
var LdSPHL = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Opcode: 0xF9, Size: 1, Cycles: 2},
	},
	NoParamFunc: ldSPHL,
}

// LdHLSPOffset - LD HL,SP+e — Load SP plus signed 8-bit offset into HL (SM83-unique).
var LdHLSPOffset = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xF8, Size: 2, Cycles: 3},
	},
	ParamFunc: ldHLSPOffset,
}

// LdAddrA - LD (nn),A — Store A to absolute address (SM83-unique at 0xEA).
var LdAddrA = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ExtendedAddressing: {Opcode: 0xEA, Size: 3, Cycles: 4},
	},
	ParamFunc: ldAddrA,
}

// LdAAddr - LD A,(nn) — Load from absolute address into A (SM83-unique at 0xFA).
var LdAAddr = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ExtendedAddressing: {Opcode: 0xFA, Size: 3, Cycles: 4},
	},
	ParamFunc: ldAAddr,
}

// ---------------------------------------------------------------------------
// LDH instructions (SM83-unique high memory access)
// ---------------------------------------------------------------------------

// LdhNA - LDH (n),A — Store A to high memory address $FF00+n.
var LdhNA = &Instruction{
	Name: LdhName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xE0, Size: 2, Cycles: 3},
	},
	ParamFunc: ldhNA,
}

// LdhAN - LDH A,(n) — Load from high memory address $FF00+n into A.
var LdhAN = &Instruction{
	Name: LdhName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xF0, Size: 2, Cycles: 3},
	},
	ParamFunc: ldhAN,
}

// LdCA - LD (C),A — Store A to $FF00+C.
var LdCA = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xE2, Size: 1, Cycles: 2},
	},
	NoParamFunc: ldCA,
}

// LdAC - LD A,(C) — Load from $FF00+C into A.
var LdAC = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xF2, Size: 1, Cycles: 2},
	},
	NoParamFunc: ldAC,
}

// ---------------------------------------------------------------------------
// Indirect memory operations
// ---------------------------------------------------------------------------

// LdIndirectImm - LD (HL),n — Load immediate to memory at (HL).
var LdIndirectImm = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Opcode: 0x36, Size: 2, Cycles: 3},
	},
	ParamFunc: ldIndirectImm,
}

// IncIndirect - INC (HL) — Increment memory at (HL).
var IncIndirect = &Instruction{
	Name: IncName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Opcode: 0x34, Size: 1, Cycles: 3},
	},
	ParamFunc: incIndirect,
}

// DecIndirect - DEC (HL) — Decrement memory at (HL).
var DecIndirect = &Instruction{
	Name: DecName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Opcode: 0x35, Size: 1, Cycles: 3},
	},
	ParamFunc: decIndirect,
}

// ---------------------------------------------------------------------------
// 8-bit ALU instructions
// ---------------------------------------------------------------------------

// AddA - ADD A,r — Add register or immediate to accumulator.
var AddA = &Instruction{
	Name: AddName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing:  {Opcode: 0x87, Size: 1, Cycles: 1},
		ImmediateAddressing: {Opcode: 0xC6, Size: 2, Cycles: 2},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB:          {Opcode: 0x80, Size: 1, Cycles: 1}, // ADD A,B
		RegC:          {Opcode: 0x81, Size: 1, Cycles: 1}, // ADD A,C
		RegD:          {Opcode: 0x82, Size: 1, Cycles: 1}, // ADD A,D
		RegE:          {Opcode: 0x83, Size: 1, Cycles: 1}, // ADD A,E
		RegH:          {Opcode: 0x84, Size: 1, Cycles: 1}, // ADD A,H
		RegL:          {Opcode: 0x85, Size: 1, Cycles: 1}, // ADD A,L
		RegHLIndirect: {Opcode: 0x86, Size: 1, Cycles: 2}, // ADD A,(HL)
		RegA:          {Opcode: 0x87, Size: 1, Cycles: 1}, // ADD A,A
	},
	ParamFunc: addA,
}

// AdcA - ADC A,r — Add with carry to accumulator.
var AdcA = &Instruction{
	Name: AdcName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing:  {Opcode: 0x8F, Size: 1, Cycles: 1},
		ImmediateAddressing: {Opcode: 0xCE, Size: 2, Cycles: 2},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB:          {Opcode: 0x88, Size: 1, Cycles: 1}, // ADC A,B
		RegC:          {Opcode: 0x89, Size: 1, Cycles: 1}, // ADC A,C
		RegD:          {Opcode: 0x8A, Size: 1, Cycles: 1}, // ADC A,D
		RegE:          {Opcode: 0x8B, Size: 1, Cycles: 1}, // ADC A,E
		RegH:          {Opcode: 0x8C, Size: 1, Cycles: 1}, // ADC A,H
		RegL:          {Opcode: 0x8D, Size: 1, Cycles: 1}, // ADC A,L
		RegHLIndirect: {Opcode: 0x8E, Size: 1, Cycles: 2}, // ADC A,(HL)
		RegA:          {Opcode: 0x8F, Size: 1, Cycles: 1}, // ADC A,A
	},
	ParamFunc: adcA,
}

// SubA - SUB r — Subtract register or immediate from accumulator.
var SubA = &Instruction{
	Name: SubName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing:  {Opcode: 0x97, Size: 1, Cycles: 1},
		ImmediateAddressing: {Opcode: 0xD6, Size: 2, Cycles: 2},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB:          {Opcode: 0x90, Size: 1, Cycles: 1}, // SUB B
		RegC:          {Opcode: 0x91, Size: 1, Cycles: 1}, // SUB C
		RegD:          {Opcode: 0x92, Size: 1, Cycles: 1}, // SUB D
		RegE:          {Opcode: 0x93, Size: 1, Cycles: 1}, // SUB E
		RegH:          {Opcode: 0x94, Size: 1, Cycles: 1}, // SUB H
		RegL:          {Opcode: 0x95, Size: 1, Cycles: 1}, // SUB L
		RegHLIndirect: {Opcode: 0x96, Size: 1, Cycles: 2}, // SUB (HL)
		RegA:          {Opcode: 0x97, Size: 1, Cycles: 1}, // SUB A
	},
	ParamFunc: subA,
}

// SbcA - SBC A,r — Subtract with carry from accumulator.
var SbcA = &Instruction{
	Name: SbcName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing:  {Opcode: 0x9F, Size: 1, Cycles: 1},
		ImmediateAddressing: {Opcode: 0xDE, Size: 2, Cycles: 2},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB:          {Opcode: 0x98, Size: 1, Cycles: 1}, // SBC A,B
		RegC:          {Opcode: 0x99, Size: 1, Cycles: 1}, // SBC A,C
		RegD:          {Opcode: 0x9A, Size: 1, Cycles: 1}, // SBC A,D
		RegE:          {Opcode: 0x9B, Size: 1, Cycles: 1}, // SBC A,E
		RegH:          {Opcode: 0x9C, Size: 1, Cycles: 1}, // SBC A,H
		RegL:          {Opcode: 0x9D, Size: 1, Cycles: 1}, // SBC A,L
		RegHLIndirect: {Opcode: 0x9E, Size: 1, Cycles: 2}, // SBC A,(HL)
		RegA:          {Opcode: 0x9F, Size: 1, Cycles: 1}, // SBC A,A
	},
	ParamFunc: sbcA,
}

// AndA - AND r — Logical AND with accumulator.
var AndA = &Instruction{
	Name: AndName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing:  {Opcode: 0xA7, Size: 1, Cycles: 1},
		ImmediateAddressing: {Opcode: 0xE6, Size: 2, Cycles: 2},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB:          {Opcode: 0xA0, Size: 1, Cycles: 1}, // AND B
		RegC:          {Opcode: 0xA1, Size: 1, Cycles: 1}, // AND C
		RegD:          {Opcode: 0xA2, Size: 1, Cycles: 1}, // AND D
		RegE:          {Opcode: 0xA3, Size: 1, Cycles: 1}, // AND E
		RegH:          {Opcode: 0xA4, Size: 1, Cycles: 1}, // AND H
		RegL:          {Opcode: 0xA5, Size: 1, Cycles: 1}, // AND L
		RegHLIndirect: {Opcode: 0xA6, Size: 1, Cycles: 2}, // AND (HL)
		RegA:          {Opcode: 0xA7, Size: 1, Cycles: 1}, // AND A
	},
	ParamFunc: andA,
}

// XorA - XOR r — Logical XOR with accumulator.
var XorA = &Instruction{
	Name: XorName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing:  {Opcode: 0xAF, Size: 1, Cycles: 1},
		ImmediateAddressing: {Opcode: 0xEE, Size: 2, Cycles: 2},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB:          {Opcode: 0xA8, Size: 1, Cycles: 1}, // XOR B
		RegC:          {Opcode: 0xA9, Size: 1, Cycles: 1}, // XOR C
		RegD:          {Opcode: 0xAA, Size: 1, Cycles: 1}, // XOR D
		RegE:          {Opcode: 0xAB, Size: 1, Cycles: 1}, // XOR E
		RegH:          {Opcode: 0xAC, Size: 1, Cycles: 1}, // XOR H
		RegL:          {Opcode: 0xAD, Size: 1, Cycles: 1}, // XOR L
		RegHLIndirect: {Opcode: 0xAE, Size: 1, Cycles: 2}, // XOR (HL)
		RegA:          {Opcode: 0xAF, Size: 1, Cycles: 1}, // XOR A
	},
	ParamFunc: xorA,
}

// OrA - OR r — Logical OR with accumulator.
var OrA = &Instruction{
	Name: OrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing:  {Opcode: 0xB7, Size: 1, Cycles: 1},
		ImmediateAddressing: {Opcode: 0xF6, Size: 2, Cycles: 2},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB:          {Opcode: 0xB0, Size: 1, Cycles: 1}, // OR B
		RegC:          {Opcode: 0xB1, Size: 1, Cycles: 1}, // OR C
		RegD:          {Opcode: 0xB2, Size: 1, Cycles: 1}, // OR D
		RegE:          {Opcode: 0xB3, Size: 1, Cycles: 1}, // OR E
		RegH:          {Opcode: 0xB4, Size: 1, Cycles: 1}, // OR H
		RegL:          {Opcode: 0xB5, Size: 1, Cycles: 1}, // OR L
		RegHLIndirect: {Opcode: 0xB6, Size: 1, Cycles: 2}, // OR (HL)
		RegA:          {Opcode: 0xB7, Size: 1, Cycles: 1}, // OR A
	},
	ParamFunc: orA,
}

// CpA - CP r — Compare register or immediate with accumulator.
var CpA = &Instruction{
	Name: CpName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing:  {Opcode: 0xBF, Size: 1, Cycles: 1},
		ImmediateAddressing: {Opcode: 0xFE, Size: 2, Cycles: 2},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB:          {Opcode: 0xB8, Size: 1, Cycles: 1}, // CP B
		RegC:          {Opcode: 0xB9, Size: 1, Cycles: 1}, // CP C
		RegD:          {Opcode: 0xBA, Size: 1, Cycles: 1}, // CP D
		RegE:          {Opcode: 0xBB, Size: 1, Cycles: 1}, // CP E
		RegH:          {Opcode: 0xBC, Size: 1, Cycles: 1}, // CP H
		RegL:          {Opcode: 0xBD, Size: 1, Cycles: 1}, // CP L
		RegHLIndirect: {Opcode: 0xBE, Size: 1, Cycles: 2}, // CP (HL)
		RegA:          {Opcode: 0xBF, Size: 1, Cycles: 1}, // CP A
	},
	ParamFunc: cpA,
}

// ---------------------------------------------------------------------------
// 8-bit INC/DEC
// ---------------------------------------------------------------------------

// IncReg8 - INC r — Increment 8-bit register.
var IncReg8 = &Instruction{
	Name: IncName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Opcode: 0x3C, Size: 1, Cycles: 1},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB: {Opcode: 0x04, Size: 1, Cycles: 1}, // INC B
		RegC: {Opcode: 0x0C, Size: 1, Cycles: 1}, // INC C
		RegD: {Opcode: 0x14, Size: 1, Cycles: 1}, // INC D
		RegE: {Opcode: 0x1C, Size: 1, Cycles: 1}, // INC E
		RegH: {Opcode: 0x24, Size: 1, Cycles: 1}, // INC H
		RegL: {Opcode: 0x2C, Size: 1, Cycles: 1}, // INC L
		RegA: {Opcode: 0x3C, Size: 1, Cycles: 1}, // INC A
	},
	ParamFunc: incReg8,
}

// DecReg8 - DEC r — Decrement 8-bit register.
var DecReg8 = &Instruction{
	Name: DecName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Opcode: 0x3D, Size: 1, Cycles: 1},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB: {Opcode: 0x05, Size: 1, Cycles: 1}, // DEC B
		RegC: {Opcode: 0x0D, Size: 1, Cycles: 1}, // DEC C
		RegD: {Opcode: 0x15, Size: 1, Cycles: 1}, // DEC D
		RegE: {Opcode: 0x1D, Size: 1, Cycles: 1}, // DEC E
		RegH: {Opcode: 0x25, Size: 1, Cycles: 1}, // DEC H
		RegL: {Opcode: 0x2D, Size: 1, Cycles: 1}, // DEC L
		RegA: {Opcode: 0x3D, Size: 1, Cycles: 1}, // DEC A
	},
	ParamFunc: decReg8,
}

// ---------------------------------------------------------------------------
// 16-bit INC/DEC
// ---------------------------------------------------------------------------

// IncReg16 - INC rr — Increment 16-bit register pair.
var IncReg16 = &Instruction{
	Name: IncName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Opcode: 0x03, Size: 1, Cycles: 2},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegBC: {Opcode: 0x03, Size: 1, Cycles: 2}, // INC BC
		RegDE: {Opcode: 0x13, Size: 1, Cycles: 2}, // INC DE
		RegHL: {Opcode: 0x23, Size: 1, Cycles: 2}, // INC HL
		RegSP: {Opcode: 0x33, Size: 1, Cycles: 2}, // INC SP
	},
	ParamFunc: incReg16,
}

// DecReg16 - DEC rr — Decrement 16-bit register pair.
var DecReg16 = &Instruction{
	Name: DecName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Opcode: 0x0B, Size: 1, Cycles: 2},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegBC: {Opcode: 0x0B, Size: 1, Cycles: 2}, // DEC BC
		RegDE: {Opcode: 0x1B, Size: 1, Cycles: 2}, // DEC DE
		RegHL: {Opcode: 0x2B, Size: 1, Cycles: 2}, // DEC HL
		RegSP: {Opcode: 0x3B, Size: 1, Cycles: 2}, // DEC SP
	},
	ParamFunc: decReg16,
}

// ---------------------------------------------------------------------------
// 16-bit ADD
// ---------------------------------------------------------------------------

// AddHL - ADD HL,rr — Add register pair to HL.
var AddHL = &Instruction{
	Name: AddName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Opcode: 0x09, Size: 1, Cycles: 2},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegBC: {Opcode: 0x09, Size: 1, Cycles: 2}, // ADD HL,BC
		RegDE: {Opcode: 0x19, Size: 1, Cycles: 2}, // ADD HL,DE
		RegHL: {Opcode: 0x29, Size: 1, Cycles: 2}, // ADD HL,HL
		RegSP: {Opcode: 0x39, Size: 1, Cycles: 2}, // ADD HL,SP
	},
	ParamFunc: addHL,
}

// AddSPE - ADD SP,e — Add signed 8-bit immediate to SP (SM83-unique).
var AddSPE = &Instruction{
	Name: AddName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xE8, Size: 2, Cycles: 4},
	},
	ParamFunc: addSPE,
}

// ---------------------------------------------------------------------------
// Stack operations
// ---------------------------------------------------------------------------

// PushReg16 - PUSH rr — Push register pair to stack.
var PushReg16 = &Instruction{
	Name: PushName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Opcode: 0xC5, Size: 1, Cycles: 4},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegBC: {Opcode: 0xC5, Size: 1, Cycles: 4}, // PUSH BC
		RegDE: {Opcode: 0xD5, Size: 1, Cycles: 4}, // PUSH DE
		RegHL: {Opcode: 0xE5, Size: 1, Cycles: 4}, // PUSH HL
		RegAF: {Opcode: 0xF5, Size: 1, Cycles: 4}, // PUSH AF
	},
	ParamFunc: pushReg16,
}

// PopReg16 - POP rr — Pop register pair from stack.
var PopReg16 = &Instruction{
	Name: PopName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Opcode: 0xC1, Size: 1, Cycles: 3},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegBC: {Opcode: 0xC1, Size: 1, Cycles: 3}, // POP BC
		RegDE: {Opcode: 0xD1, Size: 1, Cycles: 3}, // POP DE
		RegHL: {Opcode: 0xE1, Size: 1, Cycles: 3}, // POP HL
		RegAF: {Opcode: 0xF1, Size: 1, Cycles: 3}, // POP AF
	},
	ParamFunc: popReg16,
}

// ---------------------------------------------------------------------------
// Jump instructions
// ---------------------------------------------------------------------------

// JpAbs - JP nn — Jump to absolute address.
var JpAbs = &Instruction{
	Name: JpName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ExtendedAddressing: {Opcode: 0xC3, Size: 3, Cycles: 4},
	},
	ParamFunc: jpAbs,
}

// JpCond - JP cc,nn — Conditional jump to absolute address.
var JpCond = &Instruction{
	Name: JpName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ExtendedAddressing: {Opcode: 0xC2, Size: 3, Cycles: 3},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegCondNZ: {Opcode: 0xC2, Size: 3, Cycles: 3}, // JP NZ,nn
		RegCondZ:  {Opcode: 0xCA, Size: 3, Cycles: 3}, // JP Z,nn
		RegCondNC: {Opcode: 0xD2, Size: 3, Cycles: 3}, // JP NC,nn
		RegCondC:  {Opcode: 0xDA, Size: 3, Cycles: 3}, // JP C,nn
	},
	ParamFunc: jpCond,
}

// JpHL - JP (HL) — Jump to address in HL.
var JpHL = &Instruction{
	Name: JpName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Opcode: 0xE9, Size: 1, Cycles: 1},
	},
	ParamFunc: jpHL,
}

// JrRel - JR e — Jump relative.
var JrRel = &Instruction{
	Name: JrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RelativeAddressing: {Opcode: 0x18, Size: 2, Cycles: 3},
	},
	ParamFunc: jrRel,
}

// JrCond - JR cc,e — Conditional jump relative.
var JrCond = &Instruction{
	Name: JrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RelativeAddressing: {Opcode: 0x20, Size: 2, Cycles: 2},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegCondNZ: {Opcode: 0x20, Size: 2, Cycles: 2}, // JR NZ,e
		RegCondZ:  {Opcode: 0x28, Size: 2, Cycles: 2}, // JR Z,e
		RegCondNC: {Opcode: 0x30, Size: 2, Cycles: 2}, // JR NC,e
		RegCondC:  {Opcode: 0x38, Size: 2, Cycles: 2}, // JR C,e
	},
	ParamFunc: jrCond,
}

// ---------------------------------------------------------------------------
// Call instructions
// ---------------------------------------------------------------------------

// CallInst - CALL nn — Call subroutine.
var CallInst = &Instruction{
	Name: CallName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ExtendedAddressing: {Opcode: 0xCD, Size: 3, Cycles: 6},
	},
	ParamFunc: call,
}

// CallCond - CALL cc,nn — Conditional call.
var CallCond = &Instruction{
	Name: CallName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ExtendedAddressing: {Opcode: 0xC4, Size: 3, Cycles: 3},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegCondNZ: {Opcode: 0xC4, Size: 3, Cycles: 3}, // CALL NZ,nn
		RegCondZ:  {Opcode: 0xCC, Size: 3, Cycles: 3}, // CALL Z,nn
		RegCondNC: {Opcode: 0xD4, Size: 3, Cycles: 3}, // CALL NC,nn
		RegCondC:  {Opcode: 0xDC, Size: 3, Cycles: 3}, // CALL C,nn
	},
	ParamFunc: callCond,
}

// ---------------------------------------------------------------------------
// Return instructions
// ---------------------------------------------------------------------------

// RetInst - RET — Return from subroutine.
var RetInst = &Instruction{
	Name: RetName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xC9, Size: 1, Cycles: 4},
	},
	NoParamFunc: ret,
}

// RetCond - RET cc — Conditional return.
var RetCond = &Instruction{
	Name: RetName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xC0, Size: 1, Cycles: 2},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegCondNZ: {Opcode: 0xC0, Size: 1, Cycles: 2}, // RET NZ
		RegCondZ:  {Opcode: 0xC8, Size: 1, Cycles: 2}, // RET Z
		RegCondNC: {Opcode: 0xD0, Size: 1, Cycles: 2}, // RET NC
		RegCondC:  {Opcode: 0xD8, Size: 1, Cycles: 2}, // RET C
	},
	NoParamFunc: retCond,
}

// ---------------------------------------------------------------------------
// RST instructions
// ---------------------------------------------------------------------------

// RstInst - RST n — Restart (call to fixed address).
var RstInst = &Instruction{
	Name: RstName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xC7, Size: 1, Cycles: 4},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegRst00: {Opcode: 0xC7, Size: 1, Cycles: 4}, // RST 00H
		RegRst08: {Opcode: 0xCF, Size: 1, Cycles: 4}, // RST 08H
		RegRst10: {Opcode: 0xD7, Size: 1, Cycles: 4}, // RST 10H
		RegRst18: {Opcode: 0xDF, Size: 1, Cycles: 4}, // RST 18H
		RegRst20: {Opcode: 0xE7, Size: 1, Cycles: 4}, // RST 20H
		RegRst28: {Opcode: 0xEF, Size: 1, Cycles: 4}, // RST 28H
		RegRst30: {Opcode: 0xF7, Size: 1, Cycles: 4}, // RST 30H
		RegRst38: {Opcode: 0xFF, Size: 1, Cycles: 4}, // RST 38H
	},
	ParamFunc: rst,
}

// ---------------------------------------------------------------------------
// CB-prefix instructions — Rotate, shift, and bit operations
// ---------------------------------------------------------------------------

// CBRlc - RLC r — Rotate register left circular (CB prefix).
var CBRlc = &Instruction{
	Name: RlcName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xCB, Opcode: 0x00, Size: 2, Cycles: 2},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB:          {Prefix: 0xCB, Opcode: 0x00, Size: 2, Cycles: 2},
		RegC:          {Prefix: 0xCB, Opcode: 0x01, Size: 2, Cycles: 2},
		RegD:          {Prefix: 0xCB, Opcode: 0x02, Size: 2, Cycles: 2},
		RegE:          {Prefix: 0xCB, Opcode: 0x03, Size: 2, Cycles: 2},
		RegH:          {Prefix: 0xCB, Opcode: 0x04, Size: 2, Cycles: 2},
		RegL:          {Prefix: 0xCB, Opcode: 0x05, Size: 2, Cycles: 2},
		RegHLIndirect: {Prefix: 0xCB, Opcode: 0x06, Size: 2, Cycles: 4},
		RegA:          {Prefix: 0xCB, Opcode: 0x07, Size: 2, Cycles: 2},
	},
	ParamFunc: cbRlc,
}

// CBRrc - RRC r — Rotate register right circular (CB prefix).
var CBRrc = &Instruction{
	Name: RrcName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xCB, Opcode: 0x08, Size: 2, Cycles: 2},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB:          {Prefix: 0xCB, Opcode: 0x08, Size: 2, Cycles: 2},
		RegC:          {Prefix: 0xCB, Opcode: 0x09, Size: 2, Cycles: 2},
		RegD:          {Prefix: 0xCB, Opcode: 0x0A, Size: 2, Cycles: 2},
		RegE:          {Prefix: 0xCB, Opcode: 0x0B, Size: 2, Cycles: 2},
		RegH:          {Prefix: 0xCB, Opcode: 0x0C, Size: 2, Cycles: 2},
		RegL:          {Prefix: 0xCB, Opcode: 0x0D, Size: 2, Cycles: 2},
		RegHLIndirect: {Prefix: 0xCB, Opcode: 0x0E, Size: 2, Cycles: 4},
		RegA:          {Prefix: 0xCB, Opcode: 0x0F, Size: 2, Cycles: 2},
	},
	ParamFunc: cbRrc,
}

// CBRl - RL r — Rotate register left through carry (CB prefix).
var CBRl = &Instruction{
	Name: RlName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xCB, Opcode: 0x10, Size: 2, Cycles: 2},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB:          {Prefix: 0xCB, Opcode: 0x10, Size: 2, Cycles: 2},
		RegC:          {Prefix: 0xCB, Opcode: 0x11, Size: 2, Cycles: 2},
		RegD:          {Prefix: 0xCB, Opcode: 0x12, Size: 2, Cycles: 2},
		RegE:          {Prefix: 0xCB, Opcode: 0x13, Size: 2, Cycles: 2},
		RegH:          {Prefix: 0xCB, Opcode: 0x14, Size: 2, Cycles: 2},
		RegL:          {Prefix: 0xCB, Opcode: 0x15, Size: 2, Cycles: 2},
		RegHLIndirect: {Prefix: 0xCB, Opcode: 0x16, Size: 2, Cycles: 4},
		RegA:          {Prefix: 0xCB, Opcode: 0x17, Size: 2, Cycles: 2},
	},
	ParamFunc: cbRl,
}

// CBRr - RR r — Rotate register right through carry (CB prefix).
var CBRr = &Instruction{
	Name: RrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xCB, Opcode: 0x18, Size: 2, Cycles: 2},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB:          {Prefix: 0xCB, Opcode: 0x18, Size: 2, Cycles: 2},
		RegC:          {Prefix: 0xCB, Opcode: 0x19, Size: 2, Cycles: 2},
		RegD:          {Prefix: 0xCB, Opcode: 0x1A, Size: 2, Cycles: 2},
		RegE:          {Prefix: 0xCB, Opcode: 0x1B, Size: 2, Cycles: 2},
		RegH:          {Prefix: 0xCB, Opcode: 0x1C, Size: 2, Cycles: 2},
		RegL:          {Prefix: 0xCB, Opcode: 0x1D, Size: 2, Cycles: 2},
		RegHLIndirect: {Prefix: 0xCB, Opcode: 0x1E, Size: 2, Cycles: 4},
		RegA:          {Prefix: 0xCB, Opcode: 0x1F, Size: 2, Cycles: 2},
	},
	ParamFunc: cbRr,
}

// CBSla - SLA r — Shift register left arithmetic (CB prefix).
var CBSla = &Instruction{
	Name: SlaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xCB, Opcode: 0x20, Size: 2, Cycles: 2},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB:          {Prefix: 0xCB, Opcode: 0x20, Size: 2, Cycles: 2},
		RegC:          {Prefix: 0xCB, Opcode: 0x21, Size: 2, Cycles: 2},
		RegD:          {Prefix: 0xCB, Opcode: 0x22, Size: 2, Cycles: 2},
		RegE:          {Prefix: 0xCB, Opcode: 0x23, Size: 2, Cycles: 2},
		RegH:          {Prefix: 0xCB, Opcode: 0x24, Size: 2, Cycles: 2},
		RegL:          {Prefix: 0xCB, Opcode: 0x25, Size: 2, Cycles: 2},
		RegHLIndirect: {Prefix: 0xCB, Opcode: 0x26, Size: 2, Cycles: 4},
		RegA:          {Prefix: 0xCB, Opcode: 0x27, Size: 2, Cycles: 2},
	},
	ParamFunc: cbSla,
}

// CBSra - SRA r — Shift register right arithmetic (CB prefix).
var CBSra = &Instruction{
	Name: SraName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xCB, Opcode: 0x28, Size: 2, Cycles: 2},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB:          {Prefix: 0xCB, Opcode: 0x28, Size: 2, Cycles: 2},
		RegC:          {Prefix: 0xCB, Opcode: 0x29, Size: 2, Cycles: 2},
		RegD:          {Prefix: 0xCB, Opcode: 0x2A, Size: 2, Cycles: 2},
		RegE:          {Prefix: 0xCB, Opcode: 0x2B, Size: 2, Cycles: 2},
		RegH:          {Prefix: 0xCB, Opcode: 0x2C, Size: 2, Cycles: 2},
		RegL:          {Prefix: 0xCB, Opcode: 0x2D, Size: 2, Cycles: 2},
		RegHLIndirect: {Prefix: 0xCB, Opcode: 0x2E, Size: 2, Cycles: 4},
		RegA:          {Prefix: 0xCB, Opcode: 0x2F, Size: 2, Cycles: 2},
	},
	ParamFunc: cbSra,
}

// CBSwap - SWAP r — Swap upper and lower nibbles (SM83-unique, replaces Z80 SLL).
var CBSwap = &Instruction{
	Name: SwapName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xCB, Opcode: 0x30, Size: 2, Cycles: 2},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB:          {Prefix: 0xCB, Opcode: 0x30, Size: 2, Cycles: 2},
		RegC:          {Prefix: 0xCB, Opcode: 0x31, Size: 2, Cycles: 2},
		RegD:          {Prefix: 0xCB, Opcode: 0x32, Size: 2, Cycles: 2},
		RegE:          {Prefix: 0xCB, Opcode: 0x33, Size: 2, Cycles: 2},
		RegH:          {Prefix: 0xCB, Opcode: 0x34, Size: 2, Cycles: 2},
		RegL:          {Prefix: 0xCB, Opcode: 0x35, Size: 2, Cycles: 2},
		RegHLIndirect: {Prefix: 0xCB, Opcode: 0x36, Size: 2, Cycles: 4},
		RegA:          {Prefix: 0xCB, Opcode: 0x37, Size: 2, Cycles: 2},
	},
	ParamFunc: cbSwap,
}

// CBSrl - SRL r — Shift register right logical (CB prefix).
var CBSrl = &Instruction{
	Name: SrlName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xCB, Opcode: 0x38, Size: 2, Cycles: 2},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB:          {Prefix: 0xCB, Opcode: 0x38, Size: 2, Cycles: 2},
		RegC:          {Prefix: 0xCB, Opcode: 0x39, Size: 2, Cycles: 2},
		RegD:          {Prefix: 0xCB, Opcode: 0x3A, Size: 2, Cycles: 2},
		RegE:          {Prefix: 0xCB, Opcode: 0x3B, Size: 2, Cycles: 2},
		RegH:          {Prefix: 0xCB, Opcode: 0x3C, Size: 2, Cycles: 2},
		RegL:          {Prefix: 0xCB, Opcode: 0x3D, Size: 2, Cycles: 2},
		RegHLIndirect: {Prefix: 0xCB, Opcode: 0x3E, Size: 2, Cycles: 4},
		RegA:          {Prefix: 0xCB, Opcode: 0x3F, Size: 2, Cycles: 2},
	},
	ParamFunc: cbSrl,
}

// CBBit - BIT b,r — Test bit in register (CB prefix, 0x40-0x7F).
var CBBit = &Instruction{
	Name: BitName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xCB, Opcode: 0x40, Size: 2, Cycles: 2},
	},
	ParamFunc: cbBit,
}

// CBRes - RES b,r — Reset bit in register (CB prefix, 0x80-0xBF).
var CBRes = &Instruction{
	Name: ResName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xCB, Opcode: 0x80, Size: 2, Cycles: 2},
	},
	ParamFunc: cbRes,
}

// CBSet - SET b,r — Set bit in register (CB prefix, 0xC0-0xFF).
var CBSet = &Instruction{
	Name: SetName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xCB, Opcode: 0xC0, Size: 2, Cycles: 2},
	},
	ParamFunc: cbSet,
}
