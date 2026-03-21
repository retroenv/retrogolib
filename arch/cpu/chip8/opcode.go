package chip8

type Opcode struct {
	Info        OpcodeInfo
	Instruction *Instruction // CPU instruction
}

// OpcodeInfo contains the opcode and timing info for an instruction addressing mode.
type OpcodeInfo struct {
	Value uint16 // Opcode value
	Mask  uint16 // Mask to apply to opcode before comparing
}

// Opcodes maps the first nibble of the opcode to the list of opcodes.
var Opcodes = [16][]Opcode{
	0x0: {
		{Info: Opcode00E0, Instruction: ClsInst},
		{Info: Opcode00EE, Instruction: RetInst},
	},
	0x1: {
		{Info: Opcode1000, Instruction: JpInst},
	},
	0x2: {
		{Info: Opcode2000, Instruction: CallInst},
	},
	0x3: {
		{Info: Opcode3000, Instruction: SeInst},
	},
	0x4: {
		{Info: Opcode4000, Instruction: SneInst},
	},
	0x5: {
		{Info: Opcode5000, Instruction: SeInst},
	},
	0x6: {
		{Info: Opcode6000, Instruction: LdInst},
	},
	0x7: {
		{Info: Opcode7000, Instruction: AddInst},
	},
	0x8: {
		{Info: Opcode8000, Instruction: LdInst},
		{Info: Opcode8001, Instruction: OrInst},
		{Info: Opcode8002, Instruction: AndInst},
		{Info: Opcode8003, Instruction: XorInst},
		{Info: Opcode8004, Instruction: AddInst},
		{Info: Opcode8005, Instruction: SubInst},
		{Info: Opcode8006, Instruction: ShrInst},
		{Info: Opcode8007, Instruction: SubnInst},
		{Info: Opcode800E, Instruction: ShlInst},
	},
	0x9: {
		{Info: Opcode9000, Instruction: SneInst},
	},
	0xA: {
		{Info: OpcodeA000, Instruction: LdInst},
	},
	0xB: {
		{Info: OpcodeB000, Instruction: JpInst},
	},
	0xC: {
		{Info: OpcodeC000, Instruction: RndInst},
	},
	0xD: {
		{Info: OpcodeD000, Instruction: DrwInst},
	},
	0xE: {
		{Info: OpcodeE09E, Instruction: SkpInst},
		{Info: OpcodeE0A1, Instruction: SknpInst},
	},
	0xF: {
		{Info: OpcodeF007, Instruction: LdInst},
		{Info: OpcodeF00A, Instruction: LdInst},
		{Info: OpcodeF015, Instruction: LdInst},
		{Info: OpcodeF018, Instruction: LdInst},
		{Info: OpcodeF01E, Instruction: AddInst},
		{Info: OpcodeF029, Instruction: LdInst},
		{Info: OpcodeF033, Instruction: LdInst},
		{Info: OpcodeF055, Instruction: LdInst},
		{Info: OpcodeF065, Instruction: LdInst},
	},
}

var (
	Opcode00E0 = OpcodeInfo{Value: 0x00E0, Mask: 0xFFFF}
	Opcode00EE = OpcodeInfo{Value: 0x00EE, Mask: 0xFFFF}
	Opcode1000 = OpcodeInfo{Value: 0x1000, Mask: 0xF000}
	Opcode2000 = OpcodeInfo{Value: 0x2000, Mask: 0xF000}
	Opcode3000 = OpcodeInfo{Value: 0x3000, Mask: 0xF000}
	Opcode4000 = OpcodeInfo{Value: 0x4000, Mask: 0xF000}
	Opcode5000 = OpcodeInfo{Value: 0x5000, Mask: 0xF00F}
	Opcode6000 = OpcodeInfo{Value: 0x6000, Mask: 0xF000}
	Opcode7000 = OpcodeInfo{Value: 0x7000, Mask: 0xF000}
	Opcode8000 = OpcodeInfo{Value: 0x8000, Mask: 0xF00F}
	Opcode8001 = OpcodeInfo{Value: 0x8001, Mask: 0xF00F}
	Opcode8002 = OpcodeInfo{Value: 0x8002, Mask: 0xF00F}
	Opcode8003 = OpcodeInfo{Value: 0x8003, Mask: 0xF00F}
	Opcode8004 = OpcodeInfo{Value: 0x8004, Mask: 0xF00F}
	Opcode8005 = OpcodeInfo{Value: 0x8005, Mask: 0xF00F}
	Opcode8006 = OpcodeInfo{Value: 0x8006, Mask: 0xF00F}
	Opcode8007 = OpcodeInfo{Value: 0x8007, Mask: 0xF00F}
	Opcode800E = OpcodeInfo{Value: 0x800E, Mask: 0xF00F}
	Opcode9000 = OpcodeInfo{Value: 0x9000, Mask: 0xF00F}
	OpcodeA000 = OpcodeInfo{Value: 0xA000, Mask: 0xF000}
	OpcodeB000 = OpcodeInfo{Value: 0xB000, Mask: 0xF000}
	OpcodeC000 = OpcodeInfo{Value: 0xC000, Mask: 0xF000}
	OpcodeD000 = OpcodeInfo{Value: 0xD000, Mask: 0xF000}
	OpcodeE09E = OpcodeInfo{Value: 0xE09E, Mask: 0xF0FF}
	OpcodeE0A1 = OpcodeInfo{Value: 0xE0A1, Mask: 0xF0FF}
	OpcodeF007 = OpcodeInfo{Value: 0xF007, Mask: 0xF0FF}
	OpcodeF00A = OpcodeInfo{Value: 0xF00A, Mask: 0xF0FF}
	OpcodeF015 = OpcodeInfo{Value: 0xF015, Mask: 0xF0FF}
	OpcodeF018 = OpcodeInfo{Value: 0xF018, Mask: 0xF0FF}
	OpcodeF01E = OpcodeInfo{Value: 0xF01E, Mask: 0xF0FF}
	OpcodeF029 = OpcodeInfo{Value: 0xF029, Mask: 0xF0FF}
	OpcodeF033 = OpcodeInfo{Value: 0xF033, Mask: 0xF0FF}
	OpcodeF055 = OpcodeInfo{Value: 0xF055, Mask: 0xF0FF}
	OpcodeF065 = OpcodeInfo{Value: 0xF065, Mask: 0xF0FF}
)
