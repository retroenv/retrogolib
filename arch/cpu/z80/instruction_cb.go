package z80

// CB prefix instructions - Rotate, shift, and bit operations

// CBRlc rotates register left circular (RLC r, CB prefix).
var CBRlc = &Instruction{
	Name: RlcName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xCB, Opcode: 0x00, Size: 2, Cycles: 8},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB:          {Prefix: 0xCB, Opcode: 0x00, Size: 2, Cycles: 8},
		RegC:          {Prefix: 0xCB, Opcode: 0x01, Size: 2, Cycles: 8},
		RegD:          {Prefix: 0xCB, Opcode: 0x02, Size: 2, Cycles: 8},
		RegE:          {Prefix: 0xCB, Opcode: 0x03, Size: 2, Cycles: 8},
		RegH:          {Prefix: 0xCB, Opcode: 0x04, Size: 2, Cycles: 8},
		RegL:          {Prefix: 0xCB, Opcode: 0x05, Size: 2, Cycles: 8},
		RegHLIndirect: {Prefix: 0xCB, Opcode: 0x06, Size: 2, Cycles: 15},
		RegA:          {Prefix: 0xCB, Opcode: 0x07, Size: 2, Cycles: 8},
	},
	ParamFunc: cbRlc,
}

// CBRrc rotates register right circular (RRC r, CB prefix).
var CBRrc = &Instruction{
	Name: RrcName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xCB, Opcode: 0x08, Size: 2, Cycles: 8},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB:          {Prefix: 0xCB, Opcode: 0x08, Size: 2, Cycles: 8},
		RegC:          {Prefix: 0xCB, Opcode: 0x09, Size: 2, Cycles: 8},
		RegD:          {Prefix: 0xCB, Opcode: 0x0A, Size: 2, Cycles: 8},
		RegE:          {Prefix: 0xCB, Opcode: 0x0B, Size: 2, Cycles: 8},
		RegH:          {Prefix: 0xCB, Opcode: 0x0C, Size: 2, Cycles: 8},
		RegL:          {Prefix: 0xCB, Opcode: 0x0D, Size: 2, Cycles: 8},
		RegHLIndirect: {Prefix: 0xCB, Opcode: 0x0E, Size: 2, Cycles: 15},
		RegA:          {Prefix: 0xCB, Opcode: 0x0F, Size: 2, Cycles: 8},
	},
	ParamFunc: cbRrc,
}

// CBRl rotates register left through carry (RL r, CB prefix).
var CBRl = &Instruction{
	Name: RlName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xCB, Opcode: 0x10, Size: 2, Cycles: 8},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB:          {Prefix: 0xCB, Opcode: 0x10, Size: 2, Cycles: 8},
		RegC:          {Prefix: 0xCB, Opcode: 0x11, Size: 2, Cycles: 8},
		RegD:          {Prefix: 0xCB, Opcode: 0x12, Size: 2, Cycles: 8},
		RegE:          {Prefix: 0xCB, Opcode: 0x13, Size: 2, Cycles: 8},
		RegH:          {Prefix: 0xCB, Opcode: 0x14, Size: 2, Cycles: 8},
		RegL:          {Prefix: 0xCB, Opcode: 0x15, Size: 2, Cycles: 8},
		RegHLIndirect: {Prefix: 0xCB, Opcode: 0x16, Size: 2, Cycles: 15},
		RegA:          {Prefix: 0xCB, Opcode: 0x17, Size: 2, Cycles: 8},
	},
	ParamFunc: cbRl,
}

// CBRr rotates register right through carry (RR r, CB prefix).
var CBRr = &Instruction{
	Name: RrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xCB, Opcode: 0x18, Size: 2, Cycles: 8},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB:          {Prefix: 0xCB, Opcode: 0x18, Size: 2, Cycles: 8},
		RegC:          {Prefix: 0xCB, Opcode: 0x19, Size: 2, Cycles: 8},
		RegD:          {Prefix: 0xCB, Opcode: 0x1A, Size: 2, Cycles: 8},
		RegE:          {Prefix: 0xCB, Opcode: 0x1B, Size: 2, Cycles: 8},
		RegH:          {Prefix: 0xCB, Opcode: 0x1C, Size: 2, Cycles: 8},
		RegL:          {Prefix: 0xCB, Opcode: 0x1D, Size: 2, Cycles: 8},
		RegHLIndirect: {Prefix: 0xCB, Opcode: 0x1E, Size: 2, Cycles: 15},
		RegA:          {Prefix: 0xCB, Opcode: 0x1F, Size: 2, Cycles: 8},
	},
	ParamFunc: cbRr,
}

// CBSla shifts register left arithmetic (SLA r, CB prefix).
var CBSla = &Instruction{
	Name: SlaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xCB, Opcode: 0x20, Size: 2, Cycles: 8},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB:          {Prefix: 0xCB, Opcode: 0x20, Size: 2, Cycles: 8},
		RegC:          {Prefix: 0xCB, Opcode: 0x21, Size: 2, Cycles: 8},
		RegD:          {Prefix: 0xCB, Opcode: 0x22, Size: 2, Cycles: 8},
		RegE:          {Prefix: 0xCB, Opcode: 0x23, Size: 2, Cycles: 8},
		RegH:          {Prefix: 0xCB, Opcode: 0x24, Size: 2, Cycles: 8},
		RegL:          {Prefix: 0xCB, Opcode: 0x25, Size: 2, Cycles: 8},
		RegHLIndirect: {Prefix: 0xCB, Opcode: 0x26, Size: 2, Cycles: 15},
		RegA:          {Prefix: 0xCB, Opcode: 0x27, Size: 2, Cycles: 8},
	},
	ParamFunc: cbSla,
}

// CBSra shifts register right arithmetic (SRA r, CB prefix).
var CBSra = &Instruction{
	Name: SraName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xCB, Opcode: 0x28, Size: 2, Cycles: 8},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB:          {Prefix: 0xCB, Opcode: 0x28, Size: 2, Cycles: 8},
		RegC:          {Prefix: 0xCB, Opcode: 0x29, Size: 2, Cycles: 8},
		RegD:          {Prefix: 0xCB, Opcode: 0x2A, Size: 2, Cycles: 8},
		RegE:          {Prefix: 0xCB, Opcode: 0x2B, Size: 2, Cycles: 8},
		RegH:          {Prefix: 0xCB, Opcode: 0x2C, Size: 2, Cycles: 8},
		RegL:          {Prefix: 0xCB, Opcode: 0x2D, Size: 2, Cycles: 8},
		RegHLIndirect: {Prefix: 0xCB, Opcode: 0x2E, Size: 2, Cycles: 15},
		RegA:          {Prefix: 0xCB, Opcode: 0x2F, Size: 2, Cycles: 8},
	},
	ParamFunc: cbSra,
}

// CBSll shifts register left logical (SLL r, CB prefix, undocumented).
var CBSll = &Instruction{
	Name: SllName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xCB, Opcode: 0x30, Size: 2, Cycles: 8},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB:          {Prefix: 0xCB, Opcode: 0x30, Size: 2, Cycles: 8},
		RegC:          {Prefix: 0xCB, Opcode: 0x31, Size: 2, Cycles: 8},
		RegD:          {Prefix: 0xCB, Opcode: 0x32, Size: 2, Cycles: 8},
		RegE:          {Prefix: 0xCB, Opcode: 0x33, Size: 2, Cycles: 8},
		RegH:          {Prefix: 0xCB, Opcode: 0x34, Size: 2, Cycles: 8},
		RegL:          {Prefix: 0xCB, Opcode: 0x35, Size: 2, Cycles: 8},
		RegHLIndirect: {Prefix: 0xCB, Opcode: 0x36, Size: 2, Cycles: 15},
		RegA:          {Prefix: 0xCB, Opcode: 0x37, Size: 2, Cycles: 8},
	},
	ParamFunc: cbSll, // undocumented
}

// CBSrl shifts register right logical (SRL r, CB prefix).
var CBSrl = &Instruction{
	Name: SrlName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xCB, Opcode: 0x38, Size: 2, Cycles: 8},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB:          {Prefix: 0xCB, Opcode: 0x38, Size: 2, Cycles: 8},
		RegC:          {Prefix: 0xCB, Opcode: 0x39, Size: 2, Cycles: 8},
		RegD:          {Prefix: 0xCB, Opcode: 0x3A, Size: 2, Cycles: 8},
		RegE:          {Prefix: 0xCB, Opcode: 0x3B, Size: 2, Cycles: 8},
		RegH:          {Prefix: 0xCB, Opcode: 0x3C, Size: 2, Cycles: 8},
		RegL:          {Prefix: 0xCB, Opcode: 0x3D, Size: 2, Cycles: 8},
		RegHLIndirect: {Prefix: 0xCB, Opcode: 0x3E, Size: 2, Cycles: 15},
		RegA:          {Prefix: 0xCB, Opcode: 0x3F, Size: 2, Cycles: 8},
	},
	ParamFunc: cbSrl,
}

// CBBit tests bit in register (BIT b,r, CB prefix).
var CBBit = &Instruction{
	Name: BitName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xCB, Opcode: 0x40, Size: 2, Cycles: 8},
	},
	ParamFunc: cbBit,
}

// CBRes resets bit in register (RES b,r, CB prefix).
var CBRes = &Instruction{
	Name: ResName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xCB, Opcode: 0x80, Size: 2, Cycles: 8},
	},
	ParamFunc: cbRes,
}

// CBSet sets bit in register (SET b,r, CB prefix).
var CBSet = &Instruction{
	Name: SetName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xCB, Opcode: 0xC0, Size: 2, Cycles: 8},
	},
	ParamFunc: cbSet,
}
