package z80

// Instruction defines a Z80 CPU instruction with its opcodes and execution logic.
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
	AdcName       = "adc"
	AddName       = "add"
	AndName       = "and"
	BitName       = "bit"
	CallName      = "call"
	CcfName       = "ccf"
	CpName        = "cp"
	CpdName       = "cpd"
	CpdrName      = "cpdr"
	CpiName       = "cpi"
	CpirName      = "cpir"
	CplName       = "cpl"
	DaaName       = "daa"
	DdcbShiftName = "ddcb-shift"
	DecName       = "dec"
	DiName        = "di"
	DjnzName      = "djnz"
	EiName        = "ei"
	ExName        = "ex"
	ExxName       = "exx"
	FdcbShiftName = "fdcb-shift"
	HaltName      = "halt"
	ImName        = "im"
	InName        = "in"
	IncName       = "inc"
	IndName       = "ind"
	IndrName      = "indr"
	IniName       = "ini"
	InirName      = "inir"
	JpName        = "jp"
	JrName        = "jr"
	LdName        = "ld"
	LddName       = "ldd"
	LddrName      = "lddr"
	LdiName       = "ldi"
	LdirName      = "ldir"
	NegName       = "neg"
	NopName       = "nop"
	OrName        = "or"
	OtdrName      = "otdr"
	OtirName      = "otir"
	OutName       = "out"
	OutdName      = "outd"
	OutiName      = "outi"
	PopName       = "pop"
	PushName      = "push"
	ResName       = "res"
	RetName       = "ret"
	RetiName      = "reti"
	RetnName      = "retn"
	RlName        = "rl"
	RlaName       = "rla"
	RlcName       = "rlc"
	RlcaName      = "rlca"
	RldName       = "rld"
	RrName        = "rr"
	RraName       = "rra"
	RrcName       = "rrc"
	RrcaName      = "rrca"
	RrdName       = "rrd"
	RstName       = "rst"
	SbcName       = "sbc"
	ScfName       = "scf"
	SetName       = "set"
	SlaName       = "sla"
	SraName       = "sra"
	SrlName       = "srl"
	SubName       = "sub"
	XorName       = "xor"

	// Unofficial instruction names
	SllName  = "sll"
	InfName  = "inf"
	OutfName = "outf"
)

// Nop - No Operation.
var Nop = &Instruction{
	Name: NopName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0x00, Opcode: 0x00, Size: 1, Cycles: 4},
	},
	NoParamFunc: nop,
}

// Halt - Halt execution.
var Halt = &Instruction{
	Name: HaltName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0x00, Opcode: 0x76, Size: 1, Cycles: 4},
	},
	NoParamFunc: halt,
}

// LdImm8 - Load 8-bit immediate into register.
var LdImm8 = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Prefix: 0x00, Opcode: 0x3E, Size: 2, Cycles: 7}, // LD A,n (base opcode)
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB: {Prefix: 0x00, Opcode: 0x06, Size: 2, Cycles: 7}, // LD B,n
		RegC: {Prefix: 0x00, Opcode: 0x0E, Size: 2, Cycles: 7}, // LD C,n
		RegD: {Prefix: 0x00, Opcode: 0x16, Size: 2, Cycles: 7}, // LD D,n
		RegE: {Prefix: 0x00, Opcode: 0x1E, Size: 2, Cycles: 7}, // LD E,n
		RegH: {Prefix: 0x00, Opcode: 0x26, Size: 2, Cycles: 7}, // LD H,n
		RegL: {Prefix: 0x00, Opcode: 0x2E, Size: 2, Cycles: 7}, // LD L,n
		RegA: {Prefix: 0x00, Opcode: 0x3E, Size: 2, Cycles: 7}, // LD A,n
	},
	ParamFunc: ldImm8,
}

// LdReg8 - Load between 8-bit registers.
var LdReg8 = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0x00, Opcode: 0x7F, Size: 1, Cycles: 4}, // LD A,A (base opcode, others calculated)
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegLoadHLB: {Prefix: 0x00, Opcode: 0x46, Size: 1, Cycles: 7}, // LD B,(HL)
		RegLoadHLC: {Prefix: 0x00, Opcode: 0x4E, Size: 1, Cycles: 7}, // LD C,(HL)
		RegLoadHLD: {Prefix: 0x00, Opcode: 0x56, Size: 1, Cycles: 7}, // LD D,(HL)
		RegLoadHLE: {Prefix: 0x00, Opcode: 0x5E, Size: 1, Cycles: 7}, // LD E,(HL)
		RegLoadHLH: {Prefix: 0x00, Opcode: 0x66, Size: 1, Cycles: 7}, // LD H,(HL)
		RegLoadHLL: {Prefix: 0x00, Opcode: 0x6E, Size: 1, Cycles: 7}, // LD L,(HL)
		RegLoadHLA: {Prefix: 0x00, Opcode: 0x7E, Size: 1, Cycles: 7}, // LD A,(HL)
	},
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegB, RegB}: {Prefix: 0x00, Opcode: 0x40, Size: 1, Cycles: 4}, // LD B,B
		{RegB, RegC}: {Prefix: 0x00, Opcode: 0x41, Size: 1, Cycles: 4}, // LD B,C
		{RegB, RegD}: {Prefix: 0x00, Opcode: 0x42, Size: 1, Cycles: 4}, // LD B,D
		{RegB, RegE}: {Prefix: 0x00, Opcode: 0x43, Size: 1, Cycles: 4}, // LD B,E
		{RegB, RegH}: {Prefix: 0x00, Opcode: 0x44, Size: 1, Cycles: 4}, // LD B,H
		{RegB, RegL}: {Prefix: 0x00, Opcode: 0x45, Size: 1, Cycles: 4}, // LD B,L
		{RegB, RegA}: {Prefix: 0x00, Opcode: 0x47, Size: 1, Cycles: 4}, // LD B,A
		{RegC, RegB}: {Prefix: 0x00, Opcode: 0x48, Size: 1, Cycles: 4}, // LD C,B
		{RegC, RegC}: {Prefix: 0x00, Opcode: 0x49, Size: 1, Cycles: 4}, // LD C,C
		{RegC, RegD}: {Prefix: 0x00, Opcode: 0x4A, Size: 1, Cycles: 4}, // LD C,D
		{RegC, RegE}: {Prefix: 0x00, Opcode: 0x4B, Size: 1, Cycles: 4}, // LD C,E
		{RegC, RegH}: {Prefix: 0x00, Opcode: 0x4C, Size: 1, Cycles: 4}, // LD C,H
		{RegC, RegL}: {Prefix: 0x00, Opcode: 0x4D, Size: 1, Cycles: 4}, // LD C,L
		{RegC, RegA}: {Prefix: 0x00, Opcode: 0x4F, Size: 1, Cycles: 4}, // LD C,A
		{RegD, RegB}: {Prefix: 0x00, Opcode: 0x50, Size: 1, Cycles: 4}, // LD D,B
		{RegD, RegC}: {Prefix: 0x00, Opcode: 0x51, Size: 1, Cycles: 4}, // LD D,C
		{RegD, RegD}: {Prefix: 0x00, Opcode: 0x52, Size: 1, Cycles: 4}, // LD D,D
		{RegD, RegE}: {Prefix: 0x00, Opcode: 0x53, Size: 1, Cycles: 4}, // LD D,E
		{RegD, RegH}: {Prefix: 0x00, Opcode: 0x54, Size: 1, Cycles: 4}, // LD D,H
		{RegD, RegL}: {Prefix: 0x00, Opcode: 0x55, Size: 1, Cycles: 4}, // LD D,L
		{RegD, RegA}: {Prefix: 0x00, Opcode: 0x57, Size: 1, Cycles: 4}, // LD D,A
		{RegE, RegB}: {Prefix: 0x00, Opcode: 0x58, Size: 1, Cycles: 4}, // LD E,B
		{RegE, RegC}: {Prefix: 0x00, Opcode: 0x59, Size: 1, Cycles: 4}, // LD E,C
		{RegE, RegD}: {Prefix: 0x00, Opcode: 0x5A, Size: 1, Cycles: 4}, // LD E,D
		{RegE, RegE}: {Prefix: 0x00, Opcode: 0x5B, Size: 1, Cycles: 4}, // LD E,E
		{RegE, RegH}: {Prefix: 0x00, Opcode: 0x5C, Size: 1, Cycles: 4}, // LD E,H
		{RegE, RegL}: {Prefix: 0x00, Opcode: 0x5D, Size: 1, Cycles: 4}, // LD E,L
		{RegE, RegA}: {Prefix: 0x00, Opcode: 0x5F, Size: 1, Cycles: 4}, // LD E,A
		{RegH, RegB}: {Prefix: 0x00, Opcode: 0x60, Size: 1, Cycles: 4}, // LD H,B
		{RegH, RegC}: {Prefix: 0x00, Opcode: 0x61, Size: 1, Cycles: 4}, // LD H,C
		{RegH, RegD}: {Prefix: 0x00, Opcode: 0x62, Size: 1, Cycles: 4}, // LD H,D
		{RegH, RegE}: {Prefix: 0x00, Opcode: 0x63, Size: 1, Cycles: 4}, // LD H,E
		{RegH, RegH}: {Prefix: 0x00, Opcode: 0x64, Size: 1, Cycles: 4}, // LD H,H
		{RegH, RegL}: {Prefix: 0x00, Opcode: 0x65, Size: 1, Cycles: 4}, // LD H,L
		{RegH, RegA}: {Prefix: 0x00, Opcode: 0x67, Size: 1, Cycles: 4}, // LD H,A
		{RegL, RegB}: {Prefix: 0x00, Opcode: 0x68, Size: 1, Cycles: 4}, // LD L,B
		{RegL, RegC}: {Prefix: 0x00, Opcode: 0x69, Size: 1, Cycles: 4}, // LD L,C
		{RegL, RegD}: {Prefix: 0x00, Opcode: 0x6A, Size: 1, Cycles: 4}, // LD L,D
		{RegL, RegE}: {Prefix: 0x00, Opcode: 0x6B, Size: 1, Cycles: 4}, // LD L,E
		{RegL, RegH}: {Prefix: 0x00, Opcode: 0x6C, Size: 1, Cycles: 4}, // LD L,H
		{RegL, RegL}: {Prefix: 0x00, Opcode: 0x6D, Size: 1, Cycles: 4}, // LD L,L
		{RegL, RegA}: {Prefix: 0x00, Opcode: 0x6F, Size: 1, Cycles: 4}, // LD L,A
		{RegA, RegB}: {Prefix: 0x00, Opcode: 0x78, Size: 1, Cycles: 4}, // LD A,B
		{RegA, RegC}: {Prefix: 0x00, Opcode: 0x79, Size: 1, Cycles: 4}, // LD A,C
		{RegA, RegD}: {Prefix: 0x00, Opcode: 0x7A, Size: 1, Cycles: 4}, // LD A,D
		{RegA, RegE}: {Prefix: 0x00, Opcode: 0x7B, Size: 1, Cycles: 4}, // LD A,E
		{RegA, RegH}: {Prefix: 0x00, Opcode: 0x7C, Size: 1, Cycles: 4}, // LD A,H
		{RegA, RegL}: {Prefix: 0x00, Opcode: 0x7D, Size: 1, Cycles: 4}, // LD A,L
		{RegA, RegA}: {Prefix: 0x00, Opcode: 0x7F, Size: 1, Cycles: 4}, // LD A,A
	},
	ParamFunc: ldReg8,
}

// IncReg8 - Increment 8-bit register.
var IncReg8 = &Instruction{
	Name: IncName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0x00, Opcode: 0x3C, Size: 1, Cycles: 4}, // INC A (base opcode)
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB: {Prefix: 0x00, Opcode: 0x04, Size: 1, Cycles: 4}, // INC B
		RegC: {Prefix: 0x00, Opcode: 0x0C, Size: 1, Cycles: 4}, // INC C
		RegD: {Prefix: 0x00, Opcode: 0x14, Size: 1, Cycles: 4}, // INC D
		RegE: {Prefix: 0x00, Opcode: 0x1C, Size: 1, Cycles: 4}, // INC E
		RegH: {Prefix: 0x00, Opcode: 0x24, Size: 1, Cycles: 4}, // INC H
		RegL: {Prefix: 0x00, Opcode: 0x2C, Size: 1, Cycles: 4}, // INC L
		RegA: {Prefix: 0x00, Opcode: 0x3C, Size: 1, Cycles: 4}, // INC A
	},
	ParamFunc: incReg8,
}

// DecReg8 - Decrement 8-bit register.
var DecReg8 = &Instruction{
	Name: DecName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0x00, Opcode: 0x3D, Size: 1, Cycles: 4}, // DEC A (base opcode)
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB: {Prefix: 0x00, Opcode: 0x05, Size: 1, Cycles: 4}, // DEC B
		RegC: {Prefix: 0x00, Opcode: 0x0D, Size: 1, Cycles: 4}, // DEC C
		RegD: {Prefix: 0x00, Opcode: 0x15, Size: 1, Cycles: 4}, // DEC D
		RegE: {Prefix: 0x00, Opcode: 0x1D, Size: 1, Cycles: 4}, // DEC E
		RegH: {Prefix: 0x00, Opcode: 0x25, Size: 1, Cycles: 4}, // DEC H
		RegL: {Prefix: 0x00, Opcode: 0x2D, Size: 1, Cycles: 4}, // DEC L
		RegA: {Prefix: 0x00, Opcode: 0x3D, Size: 1, Cycles: 4}, // DEC A
	},
	ParamFunc: decReg8,
}

// AddA - Add to accumulator.
var AddA = &Instruction{
	Name: AddName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing:         {Prefix: 0x00, Opcode: 0x87, Size: 1, Cycles: 4}, // ADD A,A (base opcode)
		RegisterIndirectAddressing: {Prefix: 0x00, Opcode: 0x86, Size: 1, Cycles: 7}, // ADD A,(HL)
		ImmediateAddressing:        {Prefix: 0x00, Opcode: 0xC6, Size: 2, Cycles: 7}, // ADD A,n
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB: {Prefix: 0x00, Opcode: 0x80, Size: 1, Cycles: 4}, // ADD A,B
		RegC: {Prefix: 0x00, Opcode: 0x81, Size: 1, Cycles: 4}, // ADD A,C
		RegD: {Prefix: 0x00, Opcode: 0x82, Size: 1, Cycles: 4}, // ADD A,D
		RegE: {Prefix: 0x00, Opcode: 0x83, Size: 1, Cycles: 4}, // ADD A,E
		RegH: {Prefix: 0x00, Opcode: 0x84, Size: 1, Cycles: 4}, // ADD A,H
		RegL: {Prefix: 0x00, Opcode: 0x85, Size: 1, Cycles: 4}, // ADD A,L
		RegA: {Prefix: 0x00, Opcode: 0x87, Size: 1, Cycles: 4}, // ADD A,A
	},
	ParamFunc: addA,
}

// SubA - Subtract from accumulator.
var SubA = &Instruction{
	Name: SubName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing:         {Prefix: 0x00, Opcode: 0x97, Size: 1, Cycles: 4}, // SUB A (base opcode)
		RegisterIndirectAddressing: {Prefix: 0x00, Opcode: 0x96, Size: 1, Cycles: 7}, // SUB (HL)
		ImmediateAddressing:        {Prefix: 0x00, Opcode: 0xD6, Size: 2, Cycles: 7}, // SUB n
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB: {Prefix: 0x00, Opcode: 0x90, Size: 1, Cycles: 4}, // SUB B
		RegC: {Prefix: 0x00, Opcode: 0x91, Size: 1, Cycles: 4}, // SUB C
		RegD: {Prefix: 0x00, Opcode: 0x92, Size: 1, Cycles: 4}, // SUB D
		RegE: {Prefix: 0x00, Opcode: 0x93, Size: 1, Cycles: 4}, // SUB E
		RegH: {Prefix: 0x00, Opcode: 0x94, Size: 1, Cycles: 4}, // SUB H
		RegL: {Prefix: 0x00, Opcode: 0x95, Size: 1, Cycles: 4}, // SUB L
		RegA: {Prefix: 0x00, Opcode: 0x97, Size: 1, Cycles: 4}, // SUB A
	},
	ParamFunc: subA,
}

// AndA - AND with accumulator.
var AndA = &Instruction{
	Name: AndName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing:         {Prefix: 0x00, Opcode: 0xA7, Size: 1, Cycles: 4}, // AND A (base opcode)
		RegisterIndirectAddressing: {Prefix: 0x00, Opcode: 0xA6, Size: 1, Cycles: 7}, // AND (HL)
		ImmediateAddressing:        {Prefix: 0x00, Opcode: 0xE6, Size: 2, Cycles: 7}, // AND n
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB: {Prefix: 0x00, Opcode: 0xA0, Size: 1, Cycles: 4}, // AND B
		RegC: {Prefix: 0x00, Opcode: 0xA1, Size: 1, Cycles: 4}, // AND C
		RegD: {Prefix: 0x00, Opcode: 0xA2, Size: 1, Cycles: 4}, // AND D
		RegE: {Prefix: 0x00, Opcode: 0xA3, Size: 1, Cycles: 4}, // AND E
		RegH: {Prefix: 0x00, Opcode: 0xA4, Size: 1, Cycles: 4}, // AND H
		RegL: {Prefix: 0x00, Opcode: 0xA5, Size: 1, Cycles: 4}, // AND L
		RegA: {Prefix: 0x00, Opcode: 0xA7, Size: 1, Cycles: 4}, // AND A
	},
	ParamFunc: andA,
}

// OrA - OR with accumulator.
var OrA = &Instruction{
	Name: OrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing:         {Prefix: 0x00, Opcode: 0xB7, Size: 1, Cycles: 4}, // OR A (base opcode)
		RegisterIndirectAddressing: {Prefix: 0x00, Opcode: 0xB6, Size: 1, Cycles: 7}, // OR (HL)
		ImmediateAddressing:        {Prefix: 0x00, Opcode: 0xF6, Size: 2, Cycles: 7}, // OR n
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB: {Prefix: 0x00, Opcode: 0xB0, Size: 1, Cycles: 4}, // OR B
		RegC: {Prefix: 0x00, Opcode: 0xB1, Size: 1, Cycles: 4}, // OR C
		RegD: {Prefix: 0x00, Opcode: 0xB2, Size: 1, Cycles: 4}, // OR D
		RegE: {Prefix: 0x00, Opcode: 0xB3, Size: 1, Cycles: 4}, // OR E
		RegH: {Prefix: 0x00, Opcode: 0xB4, Size: 1, Cycles: 4}, // OR H
		RegL: {Prefix: 0x00, Opcode: 0xB5, Size: 1, Cycles: 4}, // OR L
		RegA: {Prefix: 0x00, Opcode: 0xB7, Size: 1, Cycles: 4}, // OR A
	},
	ParamFunc: orA,
}

// XorA - XOR with accumulator.
var XorA = &Instruction{
	Name: XorName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing:         {Prefix: 0x00, Opcode: 0xAF, Size: 1, Cycles: 4}, // XOR A (base opcode)
		RegisterIndirectAddressing: {Prefix: 0x00, Opcode: 0xAE, Size: 1, Cycles: 7}, // XOR (HL)
		ImmediateAddressing:        {Prefix: 0x00, Opcode: 0xEE, Size: 2, Cycles: 7}, // XOR n
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB: {Prefix: 0x00, Opcode: 0xA8, Size: 1, Cycles: 4}, // XOR B
		RegC: {Prefix: 0x00, Opcode: 0xA9, Size: 1, Cycles: 4}, // XOR C
		RegD: {Prefix: 0x00, Opcode: 0xAA, Size: 1, Cycles: 4}, // XOR D
		RegE: {Prefix: 0x00, Opcode: 0xAB, Size: 1, Cycles: 4}, // XOR E
		RegH: {Prefix: 0x00, Opcode: 0xAC, Size: 1, Cycles: 4}, // XOR H
		RegL: {Prefix: 0x00, Opcode: 0xAD, Size: 1, Cycles: 4}, // XOR L
		RegA: {Prefix: 0x00, Opcode: 0xAF, Size: 1, Cycles: 4}, // XOR A
	},
	ParamFunc: xorA,
}

// CpA - Compare with accumulator.
var CpA = &Instruction{
	Name: CpName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing:         {Prefix: 0x00, Opcode: 0xBF, Size: 1, Cycles: 4}, // CP A (base opcode)
		RegisterIndirectAddressing: {Prefix: 0x00, Opcode: 0xBE, Size: 1, Cycles: 7}, // CP (HL)
		ImmediateAddressing:        {Prefix: 0x00, Opcode: 0xFE, Size: 2, Cycles: 7}, // CP n
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB: {Prefix: 0x00, Opcode: 0xB8, Size: 1, Cycles: 4}, // CP B
		RegC: {Prefix: 0x00, Opcode: 0xB9, Size: 1, Cycles: 4}, // CP C
		RegD: {Prefix: 0x00, Opcode: 0xBA, Size: 1, Cycles: 4}, // CP D
		RegE: {Prefix: 0x00, Opcode: 0xBB, Size: 1, Cycles: 4}, // CP E
		RegH: {Prefix: 0x00, Opcode: 0xBC, Size: 1, Cycles: 4}, // CP H
		RegL: {Prefix: 0x00, Opcode: 0xBD, Size: 1, Cycles: 4}, // CP L
		RegA: {Prefix: 0x00, Opcode: 0xBF, Size: 1, Cycles: 4}, // CP A
	},
	ParamFunc: cpA,
}

// JpAbs - Jump absolute.
var JpAbs = &Instruction{
	Name: JpName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ExtendedAddressing: {Prefix: 0x00, Opcode: 0xC3, Size: 3, Cycles: 10}, // JP nn
	},
	ParamFunc: jpAbs,
}

// JrRel - Jump relative.
var JrRel = &Instruction{
	Name: JrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RelativeAddressing: {Prefix: 0x00, Opcode: 0x18, Size: 2, Cycles: 12}, // JR e
	},
	ParamFunc: jrRel,
}

// Additional Z80 instructions

// LdReg16 - Load 16-bit register with immediate value.
var LdReg16 = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Prefix: 0x00, Opcode: 0x01, Size: 3, Cycles: 10}, // LD BC,nn (base opcode)
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegBC: {Prefix: 0x00, Opcode: 0x01, Size: 3, Cycles: 10}, // LD BC,nn
		RegDE: {Prefix: 0x00, Opcode: 0x11, Size: 3, Cycles: 10}, // LD DE,nn
		RegHL: {Prefix: 0x00, Opcode: 0x21, Size: 3, Cycles: 10}, // LD HL,nn
		RegSP: {Prefix: 0x00, Opcode: 0x31, Size: 3, Cycles: 10}, // LD SP,nn
	},
	ParamFunc: ldReg16,
}

// LdIndirect - Load indirect (register pair to memory or memory to register).
var LdIndirect = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0x00, Opcode: 0x02, Size: 1, Cycles: 7}, // LD (BC),A (base opcode)
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegLoadBC:     {Prefix: 0x00, Opcode: 0x0A, Size: 1, Cycles: 7}, // LD A,(BC)
		RegBCIndirect: {Prefix: 0x00, Opcode: 0x02, Size: 1, Cycles: 7}, // LD (BC),A
		RegLoadDE:     {Prefix: 0x00, Opcode: 0x1A, Size: 1, Cycles: 7}, // LD A,(DE)
		RegDEIndirect: {Prefix: 0x00, Opcode: 0x12, Size: 1, Cycles: 7}, // LD (DE),A
		RegB:          {Prefix: 0x00, Opcode: 0x70, Size: 1, Cycles: 7}, // LD (HL),B
		RegC:          {Prefix: 0x00, Opcode: 0x71, Size: 1, Cycles: 7}, // LD (HL),C
		RegD:          {Prefix: 0x00, Opcode: 0x72, Size: 1, Cycles: 7}, // LD (HL),D
		RegE:          {Prefix: 0x00, Opcode: 0x73, Size: 1, Cycles: 7}, // LD (HL),E
		RegH:          {Prefix: 0x00, Opcode: 0x74, Size: 1, Cycles: 7}, // LD (HL),H
		RegL:          {Prefix: 0x00, Opcode: 0x75, Size: 1, Cycles: 7}, // LD (HL),L
		RegA:          {Prefix: 0x00, Opcode: 0x77, Size: 1, Cycles: 7}, // LD (HL),A
	},
	ParamFunc: ldIndirect,
}

// IncReg16 - Increment 16-bit register.
var IncReg16 = &Instruction{
	Name: IncName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0x00, Opcode: 0x03, Size: 1, Cycles: 6}, // INC BC (base opcode)
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegBC: {Prefix: 0x00, Opcode: 0x03, Size: 1, Cycles: 6}, // INC BC
		RegDE: {Prefix: 0x00, Opcode: 0x13, Size: 1, Cycles: 6}, // INC DE
		RegHL: {Prefix: 0x00, Opcode: 0x23, Size: 1, Cycles: 6}, // INC HL
		RegSP: {Prefix: 0x00, Opcode: 0x33, Size: 1, Cycles: 6}, // INC SP
	},
	ParamFunc: incReg16,
}

// DecReg16 - Decrement 16-bit register.
var DecReg16 = &Instruction{
	Name: DecName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0x00, Opcode: 0x0B, Size: 1, Cycles: 6}, // DEC BC (base opcode)
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegBC: {Prefix: 0x00, Opcode: 0x0B, Size: 1, Cycles: 6}, // DEC BC
		RegDE: {Prefix: 0x00, Opcode: 0x1B, Size: 1, Cycles: 6}, // DEC DE
		RegHL: {Prefix: 0x00, Opcode: 0x2B, Size: 1, Cycles: 6}, // DEC HL
		RegSP: {Prefix: 0x00, Opcode: 0x3B, Size: 1, Cycles: 6}, // DEC SP
	},
	ParamFunc: decReg16,
}

// Rlca - Rotate Left Circular Accumulator.
var Rlca = &Instruction{
	Name: RlcaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0x00, Opcode: 0x07, Size: 1, Cycles: 4},
	},
	NoParamFunc: rlca,
}

// Rrca - Rotate Right Circular Accumulator.
var Rrca = &Instruction{
	Name: RrcaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0x00, Opcode: 0x0F, Size: 1, Cycles: 4},
	},
	NoParamFunc: rrca,
}

// Rla - Rotate Left Accumulator through carry.
var Rla = &Instruction{
	Name: RlaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0x00, Opcode: 0x17, Size: 1, Cycles: 4},
	},
	NoParamFunc: rla,
}

// Rra - Rotate Right Accumulator through carry.
var Rra = &Instruction{
	Name: RraName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0x00, Opcode: 0x1F, Size: 1, Cycles: 4},
	},
	NoParamFunc: rra,
}

// ExAf - Exchange AF with AF'.
var ExAf = &Instruction{
	Name: ExName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0x00, Opcode: 0x08, Size: 1, Cycles: 4}, // EX AF,AF'
	},
	NoParamFunc: exAf,
}

// AddHl - Add register pair to HL.
var AddHl = &Instruction{
	Name: AddName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0x00, Opcode: 0x09, Size: 1, Cycles: 11}, // ADD HL,BC (base opcode)
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegBC: {Prefix: 0x00, Opcode: 0x09, Size: 1, Cycles: 11}, // ADD HL,BC
		RegDE: {Prefix: 0x00, Opcode: 0x19, Size: 1, Cycles: 11}, // ADD HL,DE
		RegHL: {Prefix: 0x00, Opcode: 0x29, Size: 1, Cycles: 11}, // ADD HL,HL
		RegSP: {Prefix: 0x00, Opcode: 0x39, Size: 1, Cycles: 11}, // ADD HL,SP
	},
	ParamFunc: addHl,
}

// Djnz - Decrement B and Jump if Not Zero.
var Djnz = &Instruction{
	Name: DjnzName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RelativeAddressing: {Prefix: 0x00, Opcode: 0x10, Size: 2, Cycles: 8},
	},
	ParamFunc: djnz,
}

// JrCond - Conditional Jump Relative.
var JrCond = &Instruction{
	Name: JrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RelativeAddressing: {Prefix: 0x00, Opcode: 0x20, Size: 2, Cycles: 7}, // JR NZ,e (base opcode)
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegCondNZ: {Prefix: 0x00, Opcode: 0x20, Size: 2, Cycles: 7}, // JR NZ,e
		RegCondZ:  {Prefix: 0x00, Opcode: 0x28, Size: 2, Cycles: 7}, // JR Z,e
		RegCondNC: {Prefix: 0x00, Opcode: 0x30, Size: 2, Cycles: 7}, // JR NC,e
		RegCondC:  {Prefix: 0x00, Opcode: 0x38, Size: 2, Cycles: 7}, // JR C,e
	},
	ParamFunc: jrCond,
}

// LdExtended - Load using extended addressing (nn).
var LdExtended = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ExtendedAddressing: {Prefix: 0x00, Opcode: 0x22, Size: 3, Cycles: 16}, // LD (nn),HL (base opcode)
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegLoadExtHL: {Prefix: 0x00, Opcode: 0x2A, Size: 3, Cycles: 16}, // LD HL,(nn)
		RegStoreExtA: {Prefix: 0x00, Opcode: 0x32, Size: 3, Cycles: 13}, // LD (nn),A
		RegLoadExtA:  {Prefix: 0x00, Opcode: 0x3A, Size: 3, Cycles: 13}, // LD A,(nn)
	},
	ParamFunc: ldExtended,
}

// Daa - Decimal Adjust Accumulator.
var Daa = &Instruction{
	Name: DaaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0x00, Opcode: 0x27, Size: 1, Cycles: 4},
	},
	NoParamFunc: daa,
}

// Cpl - Complement Accumulator.
var Cpl = &Instruction{
	Name: CplName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0x00, Opcode: 0x2F, Size: 1, Cycles: 4},
	},
	NoParamFunc: cpl,
}

// IncIndirect - Increment indirect memory location.
var IncIndirect = &Instruction{
	Name: IncName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0x00, Opcode: 0x34, Size: 1, Cycles: 11}, // INC (HL)
	},
	ParamFunc: incIndirect,
}

// DecIndirect - Decrement indirect memory location.
var DecIndirect = &Instruction{
	Name: DecName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0x00, Opcode: 0x35, Size: 1, Cycles: 11}, // DEC (HL)
	},
	ParamFunc: decIndirect,
}

// LdIndirectImm - Load immediate to indirect memory location.
var LdIndirectImm = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0x00, Opcode: 0x36, Size: 2, Cycles: 10}, // LD (HL),n
	},
	ParamFunc: ldIndirectImm,
}

// Scf - Set Carry Flag.
var Scf = &Instruction{
	Name: ScfName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0x00, Opcode: 0x37, Size: 1, Cycles: 4},
	},
	NoParamFunc: scf,
}

// Ccf - Complement Carry Flag.
var Ccf = &Instruction{
	Name: CcfName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0x00, Opcode: 0x3F, Size: 1, Cycles: 4},
	},
	NoParamFunc: ccf,
}

// AdcA - Add with Carry to Accumulator.
var AdcA = &Instruction{
	Name: AdcName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing:         {Prefix: 0x00, Opcode: 0x8F, Size: 1, Cycles: 4}, // ADC A,A
		RegisterIndirectAddressing: {Prefix: 0x00, Opcode: 0x8E, Size: 1, Cycles: 7}, // ADC A,(HL)
		ImmediateAddressing:        {Prefix: 0x00, Opcode: 0xCE, Size: 2, Cycles: 7}, // ADC A,n
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB: {Prefix: 0x00, Opcode: 0x88, Size: 1, Cycles: 4}, // ADC A,B
		RegC: {Prefix: 0x00, Opcode: 0x89, Size: 1, Cycles: 4}, // ADC A,C
		RegD: {Prefix: 0x00, Opcode: 0x8A, Size: 1, Cycles: 4}, // ADC A,D
		RegE: {Prefix: 0x00, Opcode: 0x8B, Size: 1, Cycles: 4}, // ADC A,E
		RegH: {Prefix: 0x00, Opcode: 0x8C, Size: 1, Cycles: 4}, // ADC A,H
		RegL: {Prefix: 0x00, Opcode: 0x8D, Size: 1, Cycles: 4}, // ADC A,L
		RegA: {Prefix: 0x00, Opcode: 0x8F, Size: 1, Cycles: 4}, // ADC A,A
	},
	ParamFunc: adcA,
}

// SbcA - Subtract with Carry from Accumulator.
var SbcA = &Instruction{
	Name: SbcName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing:         {Prefix: 0x00, Opcode: 0x9F, Size: 1, Cycles: 4}, // SBC A,A
		RegisterIndirectAddressing: {Prefix: 0x00, Opcode: 0x9E, Size: 1, Cycles: 7}, // SBC A,(HL)
		ImmediateAddressing:        {Prefix: 0x00, Opcode: 0xDE, Size: 2, Cycles: 7}, // SBC A,n
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB: {Prefix: 0x00, Opcode: 0x98, Size: 1, Cycles: 4}, // SBC A,B
		RegC: {Prefix: 0x00, Opcode: 0x99, Size: 1, Cycles: 4}, // SBC A,C
		RegD: {Prefix: 0x00, Opcode: 0x9A, Size: 1, Cycles: 4}, // SBC A,D
		RegE: {Prefix: 0x00, Opcode: 0x9B, Size: 1, Cycles: 4}, // SBC A,E
		RegH: {Prefix: 0x00, Opcode: 0x9C, Size: 1, Cycles: 4}, // SBC A,H
		RegL: {Prefix: 0x00, Opcode: 0x9D, Size: 1, Cycles: 4}, // SBC A,L
		RegA: {Prefix: 0x00, Opcode: 0x9F, Size: 1, Cycles: 4}, // SBC A,A
	},
	ParamFunc: sbcA,
}

// RetCond - Conditional Return.
var RetCond = &Instruction{
	Name: RetName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0x00, Opcode: 0xC0, Size: 1, Cycles: 5}, // RET NZ (base opcode)
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegCondNZ: {Prefix: 0x00, Opcode: 0xC0, Size: 1, Cycles: 5}, // RET NZ
		RegCondZ:  {Prefix: 0x00, Opcode: 0xC8, Size: 1, Cycles: 5}, // RET Z
		RegCondNC: {Prefix: 0x00, Opcode: 0xD0, Size: 1, Cycles: 5}, // RET NC
		RegCondC:  {Prefix: 0x00, Opcode: 0xD8, Size: 1, Cycles: 5}, // RET C
		RegCondPO: {Prefix: 0x00, Opcode: 0xE0, Size: 1, Cycles: 5}, // RET PO
		RegCondPE: {Prefix: 0x00, Opcode: 0xE8, Size: 1, Cycles: 5}, // RET PE
		RegCondP:  {Prefix: 0x00, Opcode: 0xF0, Size: 1, Cycles: 5}, // RET P
		RegCondM:  {Prefix: 0x00, Opcode: 0xF8, Size: 1, Cycles: 5}, // RET M
	},
	NoParamFunc: retCond,
}

// PopReg16 - Pop 16-bit register from stack.
var PopReg16 = &Instruction{
	Name: PopName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0x00, Opcode: 0xC1, Size: 1, Cycles: 10}, // POP BC (base opcode)
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegBC: {Prefix: 0x00, Opcode: 0xC1, Size: 1, Cycles: 10}, // POP BC
		RegDE: {Prefix: 0x00, Opcode: 0xD1, Size: 1, Cycles: 10}, // POP DE
		RegHL: {Prefix: 0x00, Opcode: 0xE1, Size: 1, Cycles: 10}, // POP HL
		RegAF: {Prefix: 0x00, Opcode: 0xF1, Size: 1, Cycles: 10}, // POP AF
	},
	ParamFunc: popReg16,
}

// JpCond - Conditional Jump.
var JpCond = &Instruction{
	Name: JpName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ExtendedAddressing: {Prefix: 0x00, Opcode: 0xC2, Size: 3, Cycles: 10}, // JP NZ,nn (base opcode)
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegCondNZ: {Prefix: 0x00, Opcode: 0xC2, Size: 3, Cycles: 10}, // JP NZ,nn
		RegCondZ:  {Prefix: 0x00, Opcode: 0xCA, Size: 3, Cycles: 10}, // JP Z,nn
		RegCondNC: {Prefix: 0x00, Opcode: 0xD2, Size: 3, Cycles: 10}, // JP NC,nn
		RegCondC:  {Prefix: 0x00, Opcode: 0xDA, Size: 3, Cycles: 10}, // JP C,nn
		RegCondPO: {Prefix: 0x00, Opcode: 0xE2, Size: 3, Cycles: 10}, // JP PO,nn
		RegCondPE: {Prefix: 0x00, Opcode: 0xEA, Size: 3, Cycles: 10}, // JP PE,nn
		RegCondP:  {Prefix: 0x00, Opcode: 0xF2, Size: 3, Cycles: 10}, // JP P,nn
		RegCondM:  {Prefix: 0x00, Opcode: 0xFA, Size: 3, Cycles: 10}, // JP M,nn
	},
	ParamFunc: jpCond,
}

// CallCond - Conditional Call.
var CallCond = &Instruction{
	Name: CallName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ExtendedAddressing: {Prefix: 0x00, Opcode: 0xC4, Size: 3, Cycles: 10}, // CALL NZ,nn (base opcode)
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegCondNZ: {Prefix: 0x00, Opcode: 0xC4, Size: 3, Cycles: 10}, // CALL NZ,nn
		RegCondZ:  {Prefix: 0x00, Opcode: 0xCC, Size: 3, Cycles: 10}, // CALL Z,nn
		RegCondNC: {Prefix: 0x00, Opcode: 0xD4, Size: 3, Cycles: 10}, // CALL NC,nn
		RegCondC:  {Prefix: 0x00, Opcode: 0xDC, Size: 3, Cycles: 10}, // CALL C,nn
		RegCondPO: {Prefix: 0x00, Opcode: 0xE4, Size: 3, Cycles: 10}, // CALL PO,nn
		RegCondPE: {Prefix: 0x00, Opcode: 0xEC, Size: 3, Cycles: 10}, // CALL PE,nn
		RegCondP:  {Prefix: 0x00, Opcode: 0xF4, Size: 3, Cycles: 10}, // CALL P,nn
		RegCondM:  {Prefix: 0x00, Opcode: 0xFC, Size: 3, Cycles: 10}, // CALL M,nn
	},
	ParamFunc: callCond,
}

// PushReg16 - Push 16-bit register to stack.
var PushReg16 = &Instruction{
	Name: PushName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0x00, Opcode: 0xC5, Size: 1, Cycles: 11}, // PUSH BC (base opcode)
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegBC: {Prefix: 0x00, Opcode: 0xC5, Size: 1, Cycles: 11}, // PUSH BC
		RegDE: {Prefix: 0x00, Opcode: 0xD5, Size: 1, Cycles: 11}, // PUSH DE
		RegHL: {Prefix: 0x00, Opcode: 0xE5, Size: 1, Cycles: 11}, // PUSH HL
		RegAF: {Prefix: 0x00, Opcode: 0xF5, Size: 1, Cycles: 11}, // PUSH AF
	},
	ParamFunc: pushReg16,
}

// Rst - Restart (call to fixed address).
var Rst = &Instruction{
	Name: RstName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0x00, Opcode: 0xC7, Size: 1, Cycles: 11}, // RST 00H (base opcode)
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegRst00: {Prefix: 0x00, Opcode: 0xC7, Size: 1, Cycles: 11}, // RST 00H
		RegRst08: {Prefix: 0x00, Opcode: 0xCF, Size: 1, Cycles: 11}, // RST 08H
		RegRst10: {Prefix: 0x00, Opcode: 0xD7, Size: 1, Cycles: 11}, // RST 10H
		RegRst18: {Prefix: 0x00, Opcode: 0xDF, Size: 1, Cycles: 11}, // RST 18H
		RegRst20: {Prefix: 0x00, Opcode: 0xE7, Size: 1, Cycles: 11}, // RST 20H
		RegRst28: {Prefix: 0x00, Opcode: 0xEF, Size: 1, Cycles: 11}, // RST 28H
		RegRst30: {Prefix: 0x00, Opcode: 0xF7, Size: 1, Cycles: 11}, // RST 30H
		RegRst38: {Prefix: 0x00, Opcode: 0xFF, Size: 1, Cycles: 11}, // RST 38H
	},
	ParamFunc: rst,
}

// Ret - Return from subroutine.
var Ret = &Instruction{
	Name: RetName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0x00, Opcode: 0xC9, Size: 1, Cycles: 10},
	},
	NoParamFunc: ret,
}

// Call - Call subroutine.
var Call = &Instruction{
	Name: CallName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ExtendedAddressing: {Prefix: 0x00, Opcode: 0xCD, Size: 3, Cycles: 17},
	},
	ParamFunc: call,
}

// OutPort - Output to port.
var OutPort = &Instruction{
	Name: OutName,
	Addressing: map[AddressingMode]OpcodeInfo{
		PortAddressing: {Prefix: 0x00, Opcode: 0xD3, Size: 2, Cycles: 11}, // OUT (n),A
	},
	ParamFunc: outPort,
}

// InPort - Input from port.
var InPort = &Instruction{
	Name: InName,
	Addressing: map[AddressingMode]OpcodeInfo{
		PortAddressing: {Prefix: 0x00, Opcode: 0xDB, Size: 2, Cycles: 11}, // IN A,(n)
	},
	ParamFunc: inPort,
}

// Exx - Exchange register pairs.
var Exx = &Instruction{
	Name: ExxName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0x00, Opcode: 0xD9, Size: 1, Cycles: 4},
	},
	NoParamFunc: exx,
}

// ExSp - Exchange top of stack with register pair.
var ExSp = &Instruction{
	Name: ExName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0x00, Opcode: 0xE3, Size: 1, Cycles: 19}, // EX (SP),HL
	},
	ParamFunc: exSp,
}

// JpIndirect - Jump indirect.
var JpIndirect = &Instruction{
	Name: JpName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0x00, Opcode: 0xE9, Size: 1, Cycles: 4}, // JP (HL)
	},
	ParamFunc: jpIndirect,
}

// ExDeHl - Exchange DE with HL.
var ExDeHl = &Instruction{
	Name: ExName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0x00, Opcode: 0xEB, Size: 1, Cycles: 4}, // EX DE,HL
	},
	NoParamFunc: exDeHl,
}

// Di - Disable Interrupts.
var Di = &Instruction{
	Name: DiName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0x00, Opcode: 0xF3, Size: 1, Cycles: 4},
	},
	NoParamFunc: di,
}

// Ei - Enable Interrupts.
var Ei = &Instruction{
	Name: EiName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0x00, Opcode: 0xFB, Size: 1, Cycles: 4},
	},
	NoParamFunc: ei,
}

// LdSp - Load SP from HL.
var LdSp = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0x00, Opcode: 0xF9, Size: 1, Cycles: 6}, // LD SP,HL
	},
	ParamFunc: ldSp,
}
