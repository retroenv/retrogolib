package x86

// Core x86 instruction definitions for DOS development.
// This file contains the most commonly used instructions (~585 total).

// Instruction variables for the opcode table.
var (
	// Data Movement Instructions
	MovRMReg8 = &Instruction{
		Name: "mov",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x88, Size: 2, Cycles: 2, HasModRM: true},
		},
	}

	MovRMReg16 = &Instruction{
		Name: "mov",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x89, Size: 2, Cycles: 2, HasModRM: true},
		},
	}

	MovRegRM8 = &Instruction{
		Name: "mov",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x8A, Size: 2, Cycles: 2, HasModRM: true},
		},
	}

	MovRegRM16 = &Instruction{
		Name: "mov",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x8B, Size: 2, Cycles: 2, HasModRM: true},
		},
	}

	MovRegImm8 = &Instruction{
		Name: "mov",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImmediateAddressing: {Opcode: 0xB0, Size: 2, Cycles: 4, HasModRM: false}, // Base for B0-B7
		},
		RegisterOpcodes: map[RegisterParam]OpcodeInfo{
			RegAL: {Opcode: 0xB0, Size: 2, Cycles: 4, HasModRM: false},
			RegCL: {Opcode: 0xB1, Size: 2, Cycles: 4, HasModRM: false},
			RegDL: {Opcode: 0xB2, Size: 2, Cycles: 4, HasModRM: false},
			RegBL: {Opcode: 0xB3, Size: 2, Cycles: 4, HasModRM: false},
			RegAH: {Opcode: 0xB4, Size: 2, Cycles: 4, HasModRM: false},
			RegCH: {Opcode: 0xB5, Size: 2, Cycles: 4, HasModRM: false},
			RegDH: {Opcode: 0xB6, Size: 2, Cycles: 4, HasModRM: false},
			RegBH: {Opcode: 0xB7, Size: 2, Cycles: 4, HasModRM: false},
		},
	}

	MovRegImm16 = &Instruction{
		Name: "mov",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImmediateAddressing: {Opcode: 0xB8, Size: 3, Cycles: 4, HasModRM: false}, // Base for B8-BF
		},
		RegisterOpcodes: map[RegisterParam]OpcodeInfo{
			RegAX: {Opcode: 0xB8, Size: 3, Cycles: 4, HasModRM: false},
			RegCX: {Opcode: 0xB9, Size: 3, Cycles: 4, HasModRM: false},
			RegDX: {Opcode: 0xBA, Size: 3, Cycles: 4, HasModRM: false},
			RegBX: {Opcode: 0xBB, Size: 3, Cycles: 4, HasModRM: false},
			RegSP: {Opcode: 0xBC, Size: 3, Cycles: 4, HasModRM: false},
			RegBP: {Opcode: 0xBD, Size: 3, Cycles: 4, HasModRM: false},
			RegSI: {Opcode: 0xBE, Size: 3, Cycles: 4, HasModRM: false},
			RegDI: {Opcode: 0xBF, Size: 3, Cycles: 4, HasModRM: false},
		},
	}

	MovMemImm8 = &Instruction{
		Name: "mov",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMImmediateAddressing: {Opcode: 0xC6, Size: 3, Cycles: 10, HasModRM: true},
			DirectAddressing:         {Opcode: 0xA2, Size: 3, Cycles: 10, HasModRM: false},
		},
		RegisterOpcodes: map[RegisterParam]OpcodeInfo{
			RegAL: {Opcode: 0xA2, Size: 3, Cycles: 10, HasModRM: false}, // MOV moffs8, AL
		},
	}

	MovMemImm16 = &Instruction{
		Name: "mov",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMImmediateAddressing: {Opcode: 0xC7, Size: 4, Cycles: 10, HasModRM: true},
			DirectAddressing:         {Opcode: 0xA3, Size: 3, Cycles: 10, HasModRM: false},
			ModRMRegisterAddressing:  {Opcode: 0x8C, Size: 2, Cycles: 2, HasModRM: true}, // MOV r/m16, Sreg
		},
		RegisterOpcodes: map[RegisterParam]OpcodeInfo{
			RegAX: {Opcode: 0xA3, Size: 3, Cycles: 10, HasModRM: false}, // MOV moffs16, AX
		},
	}

	// Arithmetic Instructions - ADD
	AddRMReg8 = &Instruction{
		Name: "add",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x00, Size: 2, Cycles: 3, HasModRM: true},
		},
	}
	AddRMReg16 = &Instruction{
		Name: "add",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x01, Size: 2, Cycles: 3, HasModRM: true},
		},
	}
	AddRegRM8 = &Instruction{
		Name: "add",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x02, Size: 2, Cycles: 3, HasModRM: true},
		},
	}
	AddRegRM16 = &Instruction{
		Name: "add",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x03, Size: 2, Cycles: 3, HasModRM: true},
		},
	}
	AddALImm8 = &Instruction{
		Name: "add",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImmediateAddressing: {Opcode: 0x04, Size: 2, Cycles: 4, HasModRM: false},
		},
		RegisterOpcodes: map[RegisterParam]OpcodeInfo{
			RegAL: {Opcode: 0x04, Size: 2, Cycles: 4, HasModRM: false},
		},
	}

	AddAXImm16 = &Instruction{
		Name: "add",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImmediateAddressing: {Opcode: 0x05, Size: 3, Cycles: 4, HasModRM: false},
		},
		RegisterOpcodes: map[RegisterParam]OpcodeInfo{
			RegAX: {Opcode: 0x05, Size: 3, Cycles: 4, HasModRM: false},
		},
	}

	// Arithmetic Instructions - SUB
	SubRMReg8 = &Instruction{
		Name: "sub",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x28, Size: 2, Cycles: 3, HasModRM: true},
		},
	}

	SubRMReg16 = &Instruction{
		Name: "sub",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x29, Size: 2, Cycles: 3, HasModRM: true},
		},
	}

	SubRegRM8 = &Instruction{
		Name: "sub",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x2A, Size: 2, Cycles: 3, HasModRM: true},
		},
	}
	SubRegRM16 = &Instruction{
		Name: "sub",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x2B, Size: 2, Cycles: 3, HasModRM: true},
		},
	}
	SubALImm8 = &Instruction{
		Name: "sub",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImmediateAddressing: {Opcode: 0x2C, Size: 2, Cycles: 4, HasModRM: false},
		},
		RegisterOpcodes: map[RegisterParam]OpcodeInfo{
			RegAL: {Opcode: 0x2C, Size: 2, Cycles: 4, HasModRM: false},
		},
	}
	SubAXImm16 = &Instruction{
		Name: "sub",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImmediateAddressing: {Opcode: 0x2D, Size: 3, Cycles: 4, HasModRM: false},
		},
		RegisterOpcodes: map[RegisterParam]OpcodeInfo{
			RegAX: {Opcode: 0x2D, Size: 3, Cycles: 4, HasModRM: false},
		},
	}

	// Arithmetic Instructions - ADC (Add with Carry)
	AdcRMReg8 = &Instruction{
		Name: "adc",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x10, Size: 2, Cycles: 3, HasModRM: true},
		},
	}
	AdcRMReg16 = &Instruction{
		Name: "adc",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x11, Size: 2, Cycles: 3, HasModRM: true},
		},
	}
	AdcRegRM8 = &Instruction{
		Name: "adc",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x12, Size: 2, Cycles: 3, HasModRM: true},
		},
	}
	AdcRegRM16 = &Instruction{
		Name: "adc",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x13, Size: 2, Cycles: 3, HasModRM: true},
		},
	}
	AdcALImm8 = &Instruction{
		Name: "adc",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImmediateAddressing: {Opcode: 0x14, Size: 2, Cycles: 4, HasModRM: false},
		},
		RegisterOpcodes: map[RegisterParam]OpcodeInfo{
			RegAL: {Opcode: 0x14, Size: 2, Cycles: 4, HasModRM: false},
		},
	}

	AdcAXImm16 = &Instruction{
		Name: "adc",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImmediateAddressing: {Opcode: 0x15, Size: 3, Cycles: 4, HasModRM: false},
		},
		RegisterOpcodes: map[RegisterParam]OpcodeInfo{
			RegAX: {Opcode: 0x15, Size: 3, Cycles: 4, HasModRM: false},
		},
	}

	// Arithmetic Instructions - SBB (Subtract with Borrow)
	SbbRMReg8 = &Instruction{
		Name: "sbb",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x18, Size: 2, Cycles: 3, HasModRM: true},
		},
	}
	SbbRMReg16 = &Instruction{
		Name: "sbb",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x19, Size: 2, Cycles: 3, HasModRM: true},
		},
	}
	SbbRegRM8 = &Instruction{
		Name: "sbb",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x1A, Size: 2, Cycles: 3, HasModRM: true},
		},
	}
	SbbRegRM16 = &Instruction{
		Name: "sbb",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x1B, Size: 2, Cycles: 3, HasModRM: true},
		},
	}
	SbbALImm8 = &Instruction{
		Name: "sbb",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImmediateAddressing: {Opcode: 0x1C, Size: 2, Cycles: 4, HasModRM: false},
		},
		RegisterOpcodes: map[RegisterParam]OpcodeInfo{
			RegAL: {Opcode: 0x1C, Size: 2, Cycles: 4, HasModRM: false},
		},
	}
	SbbAXImm16 = &Instruction{
		Name: "sbb",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImmediateAddressing: {Opcode: 0x1D, Size: 3, Cycles: 4, HasModRM: false},
		},
		RegisterOpcodes: map[RegisterParam]OpcodeInfo{
			RegAX: {Opcode: 0x1D, Size: 3, Cycles: 4, HasModRM: false},
		},
	}

	// Logical Instructions - AND
	AndRMReg8 = &Instruction{
		Name: "and",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x20, Size: 2, Cycles: 3, HasModRM: true},
		},
	}
	AndRMReg16 = &Instruction{
		Name: "and",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x21, Size: 2, Cycles: 3, HasModRM: true},
		},
	}
	AndRegRM8 = &Instruction{
		Name: "and",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x22, Size: 2, Cycles: 3, HasModRM: true},
		},
	}
	AndRegRM16 = &Instruction{
		Name: "and",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x23, Size: 2, Cycles: 3, HasModRM: true},
		},
	}
	AndALImm8 = &Instruction{
		Name: "and",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImmediateAddressing: {Opcode: 0x24, Size: 2, Cycles: 4, HasModRM: false},
		},
		RegisterOpcodes: map[RegisterParam]OpcodeInfo{
			RegAL: {Opcode: 0x24, Size: 2, Cycles: 4, HasModRM: false},
		},
	}
	AndAXImm16 = &Instruction{
		Name: "and",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImmediateAddressing: {Opcode: 0x25, Size: 3, Cycles: 4, HasModRM: false},
		},
		RegisterOpcodes: map[RegisterParam]OpcodeInfo{
			RegAX: {Opcode: 0x25, Size: 3, Cycles: 4, HasModRM: false},
		},
	}

	// Logical Instructions - OR
	OrRMReg8 = &Instruction{
		Name: "or",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x08, Size: 2, Cycles: 3, HasModRM: true},
		},
	}
	OrRMReg16 = &Instruction{
		Name: "or",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x09, Size: 2, Cycles: 3, HasModRM: true},
		},
	}
	OrRegRM8 = &Instruction{
		Name: "or",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x0A, Size: 2, Cycles: 3, HasModRM: true},
		},
	}
	OrRegRM16 = &Instruction{
		Name: "or",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x0B, Size: 2, Cycles: 3, HasModRM: true},
		},
	}
	OrALImm8 = &Instruction{
		Name: "or",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImmediateAddressing: {Opcode: 0x0C, Size: 2, Cycles: 4, HasModRM: false},
		},
		RegisterOpcodes: map[RegisterParam]OpcodeInfo{
			RegAL: {Opcode: 0x0C, Size: 2, Cycles: 4, HasModRM: false},
		},
	}
	OrAXImm16 = &Instruction{
		Name: "or",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImmediateAddressing: {Opcode: 0x0D, Size: 3, Cycles: 4, HasModRM: false},
		},
		RegisterOpcodes: map[RegisterParam]OpcodeInfo{
			RegAX: {Opcode: 0x0D, Size: 3, Cycles: 4, HasModRM: false},
		},
	}

	// Logical Instructions - XOR
	XorRMReg8 = &Instruction{
		Name: "xor",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x30, Size: 2, Cycles: 3, HasModRM: true},
		},
	}
	XorRMReg16 = &Instruction{
		Name: "xor",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x31, Size: 2, Cycles: 3, HasModRM: true},
		},
	}
	XorRegRM8 = &Instruction{
		Name: "xor",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x32, Size: 2, Cycles: 3, HasModRM: true},
		},
	}
	XorRegRM16 = &Instruction{
		Name: "xor",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x33, Size: 2, Cycles: 3, HasModRM: true},
		},
	}
	XorALImm8 = &Instruction{
		Name: "xor",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImmediateAddressing: {Opcode: 0x34, Size: 2, Cycles: 4, HasModRM: false},
		},
		RegisterOpcodes: map[RegisterParam]OpcodeInfo{
			RegAL: {Opcode: 0x34, Size: 2, Cycles: 4, HasModRM: false},
		},
	}
	XorAXImm16 = &Instruction{
		Name: "xor",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImmediateAddressing: {Opcode: 0x35, Size: 3, Cycles: 4, HasModRM: false},
		},
		RegisterOpcodes: map[RegisterParam]OpcodeInfo{
			RegAX: {Opcode: 0x35, Size: 3, Cycles: 4, HasModRM: false},
		},
	}

	// Comparison Instructions
	CmpRMReg8 = &Instruction{
		Name: "cmp",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x38, Size: 2, Cycles: 3, HasModRM: true},
		},
	}

	CmpRMReg16 = &Instruction{
		Name: "cmp",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x39, Size: 2, Cycles: 3, HasModRM: true},
		},
	}

	CmpRegRM8 = &Instruction{
		Name: "cmp",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x3A, Size: 2, Cycles: 3, HasModRM: true},
		},
	}
	CmpRegRM16 = &Instruction{
		Name: "cmp",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x3B, Size: 2, Cycles: 3, HasModRM: true},
		},
	}
	CmpALImm8 = &Instruction{
		Name: "cmp",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImmediateAddressing: {Opcode: 0x3C, Size: 2, Cycles: 4, HasModRM: false},
		},
		RegisterOpcodes: map[RegisterParam]OpcodeInfo{
			RegAL: {Opcode: 0x3C, Size: 2, Cycles: 4, HasModRM: false},
		},
	}
	CmpAXImm16 = &Instruction{
		Name: "cmp",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImmediateAddressing: {Opcode: 0x3D, Size: 3, Cycles: 4, HasModRM: false},
		},
		RegisterOpcodes: map[RegisterParam]OpcodeInfo{
			RegAX: {Opcode: 0x3D, Size: 3, Cycles: 4, HasModRM: false},
		},
	}

	// Increment/Decrement Instructions
	IncReg16 = &Instruction{
		Name: "inc",
		Addressing: map[AddressingMode]OpcodeInfo{
			RegisterAddressing: {Opcode: 0x40, Size: 1, Cycles: 2, HasModRM: false}, // Base for 40-47
		},
		RegisterOpcodes: map[RegisterParam]OpcodeInfo{
			RegAX: {Opcode: 0x40, Size: 1, Cycles: 2, HasModRM: false},
			RegCX: {Opcode: 0x41, Size: 1, Cycles: 2, HasModRM: false},
			RegDX: {Opcode: 0x42, Size: 1, Cycles: 2, HasModRM: false},
			RegBX: {Opcode: 0x43, Size: 1, Cycles: 2, HasModRM: false},
			RegSP: {Opcode: 0x44, Size: 1, Cycles: 2, HasModRM: false},
			RegBP: {Opcode: 0x45, Size: 1, Cycles: 2, HasModRM: false},
			RegSI: {Opcode: 0x46, Size: 1, Cycles: 2, HasModRM: false},
			RegDI: {Opcode: 0x47, Size: 1, Cycles: 2, HasModRM: false},
		},
	}
	IncRM8 = &Instruction{
		Name: "inc",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0xFE, Size: 2, Cycles: 3, HasModRM: true}, // Group 4, /0
		},
	}
	IncRM16 = &Instruction{
		Name: "inc",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0xFF, Size: 2, Cycles: 3, HasModRM: true}, // Group 5, /0
		},
	}
	DecReg16 = &Instruction{
		Name: "dec",
		Addressing: map[AddressingMode]OpcodeInfo{
			RegisterAddressing: {Opcode: 0x48, Size: 1, Cycles: 2, HasModRM: false}, // Base for 48-4F
		},
		RegisterOpcodes: map[RegisterParam]OpcodeInfo{
			RegAX: {Opcode: 0x48, Size: 1, Cycles: 2, HasModRM: false},
			RegCX: {Opcode: 0x49, Size: 1, Cycles: 2, HasModRM: false},
			RegDX: {Opcode: 0x4A, Size: 1, Cycles: 2, HasModRM: false},
			RegBX: {Opcode: 0x4B, Size: 1, Cycles: 2, HasModRM: false},
			RegSP: {Opcode: 0x4C, Size: 1, Cycles: 2, HasModRM: false},
			RegBP: {Opcode: 0x4D, Size: 1, Cycles: 2, HasModRM: false},
			RegSI: {Opcode: 0x4E, Size: 1, Cycles: 2, HasModRM: false},
			RegDI: {Opcode: 0x4F, Size: 1, Cycles: 2, HasModRM: false},
		},
	}

	// Stack Instructions
	PushReg16 = &Instruction{
		Name: "push",
		Addressing: map[AddressingMode]OpcodeInfo{
			RegisterAddressing: {Opcode: 0x50, Size: 1, Cycles: 11, HasModRM: false}, // Base for 50-57
		},
		RegisterOpcodes: map[RegisterParam]OpcodeInfo{
			RegAX: {Opcode: 0x50, Size: 1, Cycles: 11, HasModRM: false},
			RegCX: {Opcode: 0x51, Size: 1, Cycles: 11, HasModRM: false},
			RegDX: {Opcode: 0x52, Size: 1, Cycles: 11, HasModRM: false},
			RegBX: {Opcode: 0x53, Size: 1, Cycles: 11, HasModRM: false},
			RegSP: {Opcode: 0x54, Size: 1, Cycles: 11, HasModRM: false},
			RegBP: {Opcode: 0x55, Size: 1, Cycles: 11, HasModRM: false},
			RegSI: {Opcode: 0x56, Size: 1, Cycles: 11, HasModRM: false},
			RegDI: {Opcode: 0x57, Size: 1, Cycles: 11, HasModRM: false},
		},
	}
	PopReg16 = &Instruction{
		Name: "pop",
		Addressing: map[AddressingMode]OpcodeInfo{
			RegisterAddressing: {Opcode: 0x58, Size: 1, Cycles: 8, HasModRM: false}, // Base for 58-5F
		},
		RegisterOpcodes: map[RegisterParam]OpcodeInfo{
			RegAX: {Opcode: 0x58, Size: 1, Cycles: 8, HasModRM: false},
			RegCX: {Opcode: 0x59, Size: 1, Cycles: 8, HasModRM: false},
			RegDX: {Opcode: 0x5A, Size: 1, Cycles: 8, HasModRM: false},
			RegBX: {Opcode: 0x5B, Size: 1, Cycles: 8, HasModRM: false},
			RegSP: {Opcode: 0x5C, Size: 1, Cycles: 8, HasModRM: false},
			RegBP: {Opcode: 0x5D, Size: 1, Cycles: 8, HasModRM: false},
			RegSI: {Opcode: 0x5E, Size: 1, Cycles: 8, HasModRM: false},
			RegDI: {Opcode: 0x5F, Size: 1, Cycles: 8, HasModRM: false},
		},
	}
	PushCS = &Instruction{
		Name: "push",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0x0E, Size: 1, Cycles: 10, HasModRM: false},
		},
		RegisterOpcodes: map[RegisterParam]OpcodeInfo{
			RegCS: {Opcode: 0x0E, Size: 1, Cycles: 10, HasModRM: false},
		},
	}
	PushDS = &Instruction{
		Name: "push",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0x1E, Size: 1, Cycles: 10, HasModRM: false},
		},
		RegisterOpcodes: map[RegisterParam]OpcodeInfo{
			RegDS: {Opcode: 0x1E, Size: 1, Cycles: 10, HasModRM: false},
		},
	}
	PushES = &Instruction{
		Name: "push",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0x06, Size: 1, Cycles: 10, HasModRM: false},
		},
		RegisterOpcodes: map[RegisterParam]OpcodeInfo{
			RegES: {Opcode: 0x06, Size: 1, Cycles: 10, HasModRM: false},
		},
	}
	PushSS = &Instruction{
		Name: "push",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0x16, Size: 1, Cycles: 10, HasModRM: false},
		},
		RegisterOpcodes: map[RegisterParam]OpcodeInfo{
			RegSS: {Opcode: 0x16, Size: 1, Cycles: 10, HasModRM: false},
		},
	}
	PopDS = &Instruction{
		Name: "pop",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0x1F, Size: 1, Cycles: 8, HasModRM: false},
		},
		RegisterOpcodes: map[RegisterParam]OpcodeInfo{
			RegDS: {Opcode: 0x1F, Size: 1, Cycles: 8, HasModRM: false},
		},
	}
	PopES = &Instruction{
		Name: "pop",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0x07, Size: 1, Cycles: 8, HasModRM: false},
		},
		RegisterOpcodes: map[RegisterParam]OpcodeInfo{
			RegES: {Opcode: 0x07, Size: 1, Cycles: 8, HasModRM: false},
		},
	}
	PopSS = &Instruction{
		Name: "pop",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0x17, Size: 1, Cycles: 8, HasModRM: false},
		},
		RegisterOpcodes: map[RegisterParam]OpcodeInfo{
			RegSS: {Opcode: 0x17, Size: 1, Cycles: 8, HasModRM: false},
		},
	}

	// Jump Instructions - Conditional
	Jo = &Instruction{ // Jump if overflow
		Name: "jo",
		Addressing: map[AddressingMode]OpcodeInfo{
			RelativeAddressing: {Opcode: 0x70, Size: 2, Cycles: 16, HasModRM: false},
		},
	}
	Jno = &Instruction{ // Jump if not overflow
		Name: "jno",
		Addressing: map[AddressingMode]OpcodeInfo{
			RelativeAddressing: {Opcode: 0x71, Size: 2, Cycles: 16, HasModRM: false},
		},
	}
	Jb = &Instruction{ // Jump if below/carry
		Name: "jb",
		Addressing: map[AddressingMode]OpcodeInfo{
			RelativeAddressing: {Opcode: 0x72, Size: 2, Cycles: 16, HasModRM: false},
		},
	}
	Jnb = &Instruction{ // Jump if not below/not carry
		Name: "jnb",
		Addressing: map[AddressingMode]OpcodeInfo{
			RelativeAddressing: {Opcode: 0x73, Size: 2, Cycles: 16, HasModRM: false},
		},
	}
	Jz = &Instruction{ // Jump if zero/equal
		Name: "jz",
		Addressing: map[AddressingMode]OpcodeInfo{
			RelativeAddressing: {Opcode: 0x74, Size: 2, Cycles: 16, HasModRM: false},
		},
	}
	Jnz = &Instruction{ // Jump if not zero/not equal
		Name: "jnz",
		Addressing: map[AddressingMode]OpcodeInfo{
			RelativeAddressing: {Opcode: 0x75, Size: 2, Cycles: 16, HasModRM: false},
		},
	}
	Jbe = &Instruction{ // Jump if below or equal
		Name: "jbe",
		Addressing: map[AddressingMode]OpcodeInfo{
			RelativeAddressing: {Opcode: 0x76, Size: 2, Cycles: 16, HasModRM: false},
		},
	}
	Jnbe = &Instruction{ // Jump if not below or equal
		Name: "jnbe",
		Addressing: map[AddressingMode]OpcodeInfo{
			RelativeAddressing: {Opcode: 0x77, Size: 2, Cycles: 16, HasModRM: false},
		},
	}
	Js = &Instruction{ // Jump if sign
		Name: "js",
		Addressing: map[AddressingMode]OpcodeInfo{
			RelativeAddressing: {Opcode: 0x78, Size: 2, Cycles: 16, HasModRM: false},
		},
	}
	Jns = &Instruction{ // Jump if not sign
		Name: "jns",
		Addressing: map[AddressingMode]OpcodeInfo{
			RelativeAddressing: {Opcode: 0x79, Size: 2, Cycles: 16, HasModRM: false},
		},
	}
	Jp = &Instruction{ // Jump if parity/parity even
		Name: "jp",
		Addressing: map[AddressingMode]OpcodeInfo{
			RelativeAddressing: {Opcode: 0x7A, Size: 2, Cycles: 16, HasModRM: false},
		},
	}
	Jnp = &Instruction{ // Jump if not parity/parity odd
		Name: "jnp",
		Addressing: map[AddressingMode]OpcodeInfo{
			RelativeAddressing: {Opcode: 0x7B, Size: 2, Cycles: 16, HasModRM: false},
		},
	}
	Jl = &Instruction{ // Jump if less
		Name: "jl",
		Addressing: map[AddressingMode]OpcodeInfo{
			RelativeAddressing: {Opcode: 0x7C, Size: 2, Cycles: 16, HasModRM: false},
		},
	}
	Jnl = &Instruction{ // Jump if not less
		Name: "jnl",
		Addressing: map[AddressingMode]OpcodeInfo{
			RelativeAddressing: {Opcode: 0x7D, Size: 2, Cycles: 16, HasModRM: false},
		},
	}
	Jle = &Instruction{ // Jump if less or equal
		Name: "jle",
		Addressing: map[AddressingMode]OpcodeInfo{
			RelativeAddressing: {Opcode: 0x7E, Size: 2, Cycles: 16, HasModRM: false},
		},
	}
	Jnle = &Instruction{ // Jump if not less or equal
		Name: "jnle",
		Addressing: map[AddressingMode]OpcodeInfo{
			RelativeAddressing: {Opcode: 0x7F, Size: 2, Cycles: 16, HasModRM: false},
		},
	}

	// Jump Instructions - Unconditional
	Jmp = &Instruction{ // Unconditional jump
		Name: "jmp",
		Addressing: map[AddressingMode]OpcodeInfo{
			RelativeAddressing:      {Opcode: 0xE9, Size: 3, Cycles: 15, HasModRM: false}, // JMP rel16
			ModRMRegisterAddressing: {Opcode: 0xFF, Size: 2, Cycles: 11, HasModRM: true},  // Group 5, /4
		},
	}
	JmpFar = &Instruction{ // Far jump
		Name: "jmp",
		Addressing: map[AddressingMode]OpcodeInfo{
			SegmentOffsetAddressing: {Opcode: 0xEA, Size: 5, Cycles: 15, HasModRM: false}, // JMP ptr16:16
		},
	}
	Call = &Instruction{ // Call procedure
		Name: "call",
		Addressing: map[AddressingMode]OpcodeInfo{
			RelativeAddressing:      {Opcode: 0xE8, Size: 3, Cycles: 19, HasModRM: false}, // CALL rel16
			ModRMRegisterAddressing: {Opcode: 0xFF, Size: 2, Cycles: 16, HasModRM: true},  // Group 5, /2
		},
	}
	CallFar = &Instruction{ // Far call
		Name: "call",
		Addressing: map[AddressingMode]OpcodeInfo{
			SegmentOffsetAddressing: {Opcode: 0x9A, Size: 5, Cycles: 28, HasModRM: false}, // CALL ptr16:16
		},
	}
	Ret = &Instruction{ // Return
		Name: "ret",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing:   {Opcode: 0xC3, Size: 1, Cycles: 16, HasModRM: false}, // RET
			ImmediateAddressing: {Opcode: 0xC2, Size: 3, Cycles: 20, HasModRM: false}, // RET imm16
		},
	}
	RetFar = &Instruction{ // Far return
		Name: "retf",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing:   {Opcode: 0xCB, Size: 1, Cycles: 34, HasModRM: false}, // RETF
			ImmediateAddressing: {Opcode: 0xCA, Size: 3, Cycles: 25, HasModRM: false}, // RETF imm16
		},
	}

	// Interrupt Instructions
	Int = &Instruction{ // Software interrupt
		Name: "int",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImmediateAddressing: {Opcode: 0xCD, Size: 2, Cycles: 51, HasModRM: false}, // INT imm8
		},
	}
	Into = &Instruction{ // Interrupt on overflow
		Name: "into",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0xCE, Size: 1, Cycles: 53, HasModRM: false},
		},
	}
	Iret = &Instruction{ // Return from interrupt
		Name: "iret",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0xCF, Size: 1, Cycles: 32, HasModRM: false},
		},
	}

	// Flag Instructions
	Clc = &Instruction{ // Resets carry flag to prepare for arithmetic operations
		Name: "clc",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0xF8, Size: 1, Cycles: 2, HasModRM: false},
		},
	}
	Stc = &Instruction{ // Forces carry flag for arithmetic operations
		Name: "stc",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0xF9, Size: 1, Cycles: 2, HasModRM: false},
		},
	}
	Cmc = &Instruction{ // Toggles carry flag state (0→1, 1→0)
		Name: "cmc",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0xF5, Size: 1, Cycles: 2, HasModRM: false},
		},
	}
	Cld = &Instruction{ // Sets string operations to increment (forward direction)
		Name: "cld",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0xFC, Size: 1, Cycles: 2, HasModRM: false},
		},
	}
	Std = &Instruction{ // Sets string operations to decrement (backward direction)
		Name: "std",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0xFD, Size: 1, Cycles: 2, HasModRM: false},
		},
	}
	Cli = &Instruction{ // Disables maskable interrupts (critical sections)
		Name: "cli",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0xFA, Size: 1, Cycles: 2, HasModRM: false},
		},
	}
	Sti = &Instruction{ // Enables maskable interrupts after next instruction
		Name: "sti",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0xFB, Size: 1, Cycles: 2, HasModRM: false},
		},
	}

	// String Instructions
	Movsb = &Instruction{ // Move string byte
		Name: "movsb",
		Addressing: map[AddressingMode]OpcodeInfo{
			StringAddressing: {Opcode: 0xA4, Size: 1, Cycles: 18, HasModRM: false},
		},
	}
	Movsw = &Instruction{ // Move string word
		Name: "movsw",
		Addressing: map[AddressingMode]OpcodeInfo{
			StringAddressing: {Opcode: 0xA5, Size: 1, Cycles: 18, HasModRM: false},
		},
	}
	Cmpsb = &Instruction{ // Compare string byte
		Name: "cmpsb",
		Addressing: map[AddressingMode]OpcodeInfo{
			StringAddressing: {Opcode: 0xA6, Size: 1, Cycles: 22, HasModRM: false},
		},
	}
	Cmpsw = &Instruction{ // Compare string word
		Name: "cmpsw",
		Addressing: map[AddressingMode]OpcodeInfo{
			StringAddressing: {Opcode: 0xA7, Size: 1, Cycles: 22, HasModRM: false},
		},
	}
	Scasb = &Instruction{ // Scan string byte
		Name: "scasb",
		Addressing: map[AddressingMode]OpcodeInfo{
			StringAddressing: {Opcode: 0xAE, Size: 1, Cycles: 15, HasModRM: false},
		},
	}
	Scasw = &Instruction{ // Scan string word
		Name: "scasw",
		Addressing: map[AddressingMode]OpcodeInfo{
			StringAddressing: {Opcode: 0xAF, Size: 1, Cycles: 15, HasModRM: false},
		},
	}
	Lodsb = &Instruction{ // Load string byte
		Name: "lodsb",
		Addressing: map[AddressingMode]OpcodeInfo{
			StringAddressing: {Opcode: 0xAC, Size: 1, Cycles: 12, HasModRM: false},
		},
	}
	Lodsw = &Instruction{ // Load string word
		Name: "lodsw",
		Addressing: map[AddressingMode]OpcodeInfo{
			StringAddressing: {Opcode: 0xAD, Size: 1, Cycles: 12, HasModRM: false},
		},
	}
	Stosb = &Instruction{ // Store string byte
		Name: "stosb",
		Addressing: map[AddressingMode]OpcodeInfo{
			StringAddressing: {Opcode: 0xAA, Size: 1, Cycles: 11, HasModRM: false},
		},
	}
	Stosw = &Instruction{ // Store string word
		Name: "stosw",
		Addressing: map[AddressingMode]OpcodeInfo{
			StringAddressing: {Opcode: 0xAB, Size: 1, Cycles: 11, HasModRM: false},
		},
	}

	// Repeat Prefixes
	Rep = &Instruction{ // Repeat
		Name: "rep",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0xF3, Size: 1, Cycles: 2, HasModRM: false},
		},
	}
	Repz = &Instruction{ // Repeat while zero
		Name: "repz",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0xF3, Size: 1, Cycles: 2, HasModRM: false},
		},
	}
	Repnz = &Instruction{ // Repeat while not zero
		Name: "repnz",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0xF2, Size: 1, Cycles: 2, HasModRM: false},
		},
	}

	// Shift and Rotate Instructions
	Shl = &Instruction{ // Shift left
		Name: "shl",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0xD0, Size: 2, Cycles: 2, HasModRM: true}, // SHL r/m8, 1
		},
	}
	Shr = &Instruction{ // Shift right
		Name: "shr",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0xD0, Size: 2, Cycles: 2, HasModRM: true}, // SHR r/m8, 1
		},
	}
	Sar = &Instruction{ // Shift arithmetic right
		Name: "sar",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0xD0, Size: 2, Cycles: 2, HasModRM: true}, // SAR r/m8, 1
		},
	}
	Rol = &Instruction{ // Rotate left
		Name: "rol",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0xD0, Size: 2, Cycles: 2, HasModRM: true}, // ROL r/m8, 1
		},
	}
	Ror = &Instruction{ // Rotate right
		Name: "ror",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0xD0, Size: 2, Cycles: 2, HasModRM: true}, // ROR r/m8, 1
		},
	}
	Rcl = &Instruction{ // Rotate through carry left
		Name: "rcl",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0xD0, Size: 2, Cycles: 2, HasModRM: true}, // RCL r/m8, 1
		},
	}
	Rcr = &Instruction{ // Rotate through carry right
		Name: "rcr",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0xD0, Size: 2, Cycles: 2, HasModRM: true}, // RCR r/m8, 1
		},
	}

	// Test Instructions
	Test = &Instruction{ // Test (logical AND without storing result)
		Name: "test",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x84, Size: 2, Cycles: 3, HasModRM: true},  // TEST r/m8, r8
			ImmediateAddressing:     {Opcode: 0xA8, Size: 2, Cycles: 4, HasModRM: false}, // TEST AL, imm8
		},
	}

	// Exchange Instructions
	Xchg = &Instruction{ // Exchange
		Name: "xchg",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x86, Size: 2, Cycles: 4, HasModRM: true},  // XCHG r/m8, r8
			RegisterAddressing:      {Opcode: 0x90, Size: 1, Cycles: 3, HasModRM: false}, // XCHG AX, reg16
		},
	}

	// Segment Override Prefixes
	SegES = &Instruction{ // ES segment prefix
		Name: "es:",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0x26, Size: 1, Cycles: 2, HasModRM: false},
		},
	}
	SegCS = &Instruction{ // CS segment prefix
		Name: "cs:",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0x2E, Size: 1, Cycles: 2, HasModRM: false},
		},
	}
	SegSS = &Instruction{ // SS segment prefix
		Name: "ss:",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0x36, Size: 1, Cycles: 2, HasModRM: false},
		},
	}
	SegDS = &Instruction{ // DS segment prefix
		Name: "ds:",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0x3E, Size: 1, Cycles: 2, HasModRM: false},
		},
	}

	// Decimal Arithmetic
	Daa = &Instruction{ // Decimal adjust after addition
		Name: "daa",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0x27, Size: 1, Cycles: 4, HasModRM: false},
		},
	}
	Das = &Instruction{ // Decimal adjust after subtraction
		Name: "das",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0x2F, Size: 1, Cycles: 4, HasModRM: false},
		},
	}
	Aaa = &Instruction{ // ASCII adjust after addition
		Name: "aaa",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0x37, Size: 1, Cycles: 4, HasModRM: false},
		},
	}
	Aas = &Instruction{ // ASCII adjust after subtraction
		Name: "aas",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0x3F, Size: 1, Cycles: 4, HasModRM: false},
		},
	}

	// Multiplication and Division
	Mul = &Instruction{ // Multiply
		Name: "mul",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0xF6, Size: 2, Cycles: 70, HasModRM: true}, // MUL r/m8 (Group 3, /4)
		},
	}
	Imul = &Instruction{ // Signed multiply
		Name: "imul",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0xF6, Size: 2, Cycles: 80, HasModRM: true}, // IMUL r/m8 (Group 3, /5)
		},
	}
	Div = &Instruction{ // Divide
		Name: "div",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0xF6, Size: 2, Cycles: 80, HasModRM: true}, // DIV r/m8 (Group 3, /6)
		},
	}
	Idiv = &Instruction{ // Signed divide
		Name: "idiv",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0xF6, Size: 2, Cycles: 101, HasModRM: true}, // IDIV r/m8 (Group 3, /7)
		},
	}

	// I/O Instructions
	In = &Instruction{ // Input from port
		Name: "in",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImmediateAddressing: {Opcode: 0xE4, Size: 2, Cycles: 10, HasModRM: false}, // IN AL, imm8
			RegisterAddressing:  {Opcode: 0xEC, Size: 1, Cycles: 8, HasModRM: false},  // IN AL, DX
		},
	}
	Out = &Instruction{ // Output to port
		Name: "out",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImmediateAddressing: {Opcode: 0xE6, Size: 2, Cycles: 10, HasModRM: false}, // OUT imm8, AL
			RegisterAddressing:  {Opcode: 0xEE, Size: 1, Cycles: 8, HasModRM: false},  // OUT DX, AL
		},
	}

	// Control Instructions
	Nop = &Instruction{ // No operation
		Name: "nop",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0x90, Size: 1, Cycles: 3, HasModRM: false},
		},
	}
	Hlt = &Instruction{ // Halt
		Name: "hlt",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0xF4, Size: 1, Cycles: 2, HasModRM: false},
		},
	}

	// Other Instructions
	Cbw = &Instruction{ // Convert byte to word
		Name: "cbw",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0x98, Size: 1, Cycles: 2, HasModRM: false},
		},
	}
	Cwd = &Instruction{ // Convert word to double word
		Name: "cwd",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0x99, Size: 1, Cycles: 5, HasModRM: false},
		},
	}
	Xlat = &Instruction{ // Table lookup translation
		Name: "xlat",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0xD7, Size: 1, Cycles: 11, HasModRM: false},
		},
	}
	Lea = &Instruction{ // Load effective address
		Name: "lea",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x8D, Size: 2, Cycles: 2, HasModRM: true},
		},
	}

	// 80186/80188 Instructions (introduced 1982)
	Pusha = &Instruction{ // Push all general registers
		Name: "pusha",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0x60, Size: 1, Cycles: 36, HasModRM: false},
		},
	}
	Popa = &Instruction{ // Pop all general registers
		Name: "popa",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0x61, Size: 1, Cycles: 51, HasModRM: false},
		},
	}
	Bound = &Instruction{ // Check array bounds
		Name: "bound",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x62, Size: 2, Cycles: 33, HasModRM: true},
		},
	}
	Enter = &Instruction{ // Make stack frame
		Name: "enter",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImmediateAddressing: {Opcode: 0xC8, Size: 4, Cycles: 25, HasModRM: false},
		},
	}
	Leave = &Instruction{ // High level procedure leave
		Name: "leave",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0xC9, Size: 1, Cycles: 8, HasModRM: false},
		},
	}
	Insb = &Instruction{ // Input byte from port to ES:[DI]
		Name: "insb",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0x6C, Size: 1, Cycles: 14, HasModRM: false},
		},
	}
	Insw = &Instruction{ // Input word from port to ES:[DI]
		Name: "insw",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0x6D, Size: 1, Cycles: 14, HasModRM: false},
		},
	}
	Outsb = &Instruction{ // Output byte from DS:[SI] to port
		Name: "outsb",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0x6E, Size: 1, Cycles: 14, HasModRM: false},
		},
	}
	Outsw = &Instruction{ // Output word from DS:[SI] to port
		Name: "outsw",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0x6F, Size: 1, Cycles: 14, HasModRM: false},
		},
	}
	PushImm16 = &Instruction{ // Push immediate word
		Name: "push",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImmediateAddressing: {Opcode: 0x68, Size: 3, Cycles: 3, HasModRM: false},
		},
	}
	PushImm8 = &Instruction{ // Push immediate byte (sign-extended)
		Name: "push",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImmediateAddressing: {Opcode: 0x6A, Size: 2, Cycles: 3, HasModRM: false},
		},
	}
	ImulRegRMImm16 = &Instruction{ // Signed multiply register with r/m and immediate word
		Name: "imul",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x69, Size: 4, Cycles: 22, HasModRM: true},
		},
	}
	ImulRegRMImm8 = &Instruction{ // Signed multiply register with r/m and immediate byte
		Name: "imul",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x6B, Size: 3, Cycles: 22, HasModRM: true},
		},
	}

	// 80286 Real Mode Instructions (introduced 1982)
	// Note: Most 80286 additions are for protected mode.
	// The following are usable in real mode:
	Smsw = &Instruction{ // Store Machine Status Word
		Name: "smsw",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x0F01, Size: 3, Cycles: 2, HasModRM: true}, // 0x0F 0x01 /4
		},
	}
	Lmsw = &Instruction{ // Load Machine Status Word
		Name: "lmsw",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x0F01, Size: 3, Cycles: 6, HasModRM: true}, // 0x0F 0x01 /6
		},
	}

	// 80386 Instructions (introduced 1985)
	Bsf = &Instruction{ // Bit scan forward
		Name: "bsf",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x0FBC, Size: 3, Cycles: 10, HasModRM: true},
		},
	}
	Bsr = &Instruction{ // Bit scan reverse
		Name: "bsr",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x0FBD, Size: 3, Cycles: 10, HasModRM: true},
		},
	}
	Bt = &Instruction{ // Bit test
		Name: "bt",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x0FA3, Size: 3, Cycles: 3, HasModRM: true},
		},
	}
	Btc = &Instruction{ // Bit test and complement
		Name: "btc",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x0FBB, Size: 3, Cycles: 6, HasModRM: true},
		},
	}
	Btr = &Instruction{ // Bit test and reset
		Name: "btr",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x0FB3, Size: 3, Cycles: 6, HasModRM: true},
		},
	}
	Bts = &Instruction{ // Bit test and set
		Name: "bts",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x0FAB, Size: 3, Cycles: 6, HasModRM: true},
		},
	}
	Movzx = &Instruction{ // Move with zero extension
		Name: "movzx",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x0FB6, Size: 3, Cycles: 3, HasModRM: true}, // MOVZX r16, r/m8
		},
	}
	Movsx = &Instruction{ // Move with sign extension
		Name: "movsx",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x0FBE, Size: 3, Cycles: 3, HasModRM: true}, // MOVSX r16, r/m8
		},
	}
	Shld = &Instruction{ // Double precision shift left
		Name: "shld",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x0FA4, Size: 4, Cycles: 3, HasModRM: true},
		},
	}
	Shrd = &Instruction{ // Double precision shift right
		Name: "shrd",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x0FAC, Size: 4, Cycles: 3, HasModRM: true},
		},
	}
	Setcc = &Instruction{ // Set byte on condition (generic - specific conditions use 0x0F 0x90-0x9F)
		Name: "setcc",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x0F90, Size: 3, Cycles: 4, HasModRM: true},
		},
	}

	// 80486 Instructions (introduced 1989)
	Bswap = &Instruction{ // Byte swap register (EAX-EDI: 0x0F 0xC8-0xCF)
		Name: "bswap",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0x0FC8, Size: 2, Cycles: 1, HasModRM: false}, // BSWAP EAX
		},
	}
	Cmpxchg = &Instruction{ // Compare and exchange
		Name: "cmpxchg",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x0FB0, Size: 3, Cycles: 6, HasModRM: true}, // CMPXCHG r/m8, r8
		},
	}
	Xadd = &Instruction{ // Exchange and add
		Name: "xadd",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0x0FC0, Size: 3, Cycles: 3, HasModRM: true}, // XADD r/m8, r8
		},
	}
	Invd = &Instruction{ // Invalidate cache
		Name: "invd",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0x0F08, Size: 2, Cycles: 15, HasModRM: false},
		},
	}
	Wbinvd = &Instruction{ // Write back and invalidate cache
		Name: "wbinvd",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0x0F09, Size: 2, Cycles: 255, HasModRM: false}, // Varies greatly (typically 2000+)
		},
	}

	// Undefined/Reserved
	Undefined = &Instruction{ // Placeholder for undefined opcodes
		Name: "undefined",
	}
)

// Instructions maps instruction names to their information struct.
// This enables name-based lookup for disassemblers and tooling.
var Instructions = map[string]*Instruction{
	AaaName:     Aaa,
	AasName:     Aas,
	AdcName:     AdcALImm8, // Representative ADC instruction
	AddName:     AddALImm8, // Representative ADD instruction
	AndName:     AndALImm8, // Representative AND instruction
	BoundName:   Bound,
	BsfName:     Bsf,
	BsrName:     Bsr,
	BswapName:   Bswap,
	BtName:      Bt,
	BtcName:     Btc,
	BtrName:     Btr,
	BtsName:     Bts,
	CallName:    Call,
	CbwName:     Cbw,
	ClcName:     Clc,
	CldName:     Cld,
	CliName:     Cli,
	CmcName:     Cmc,
	CmpName:     CmpALImm8, // Representative CMP instruction
	CmpsbName:   Cmpsb,
	CmpswName:   Cmpsw,
	CmpxchgName: Cmpxchg,
	CwdName:     Cwd,
	DaaName:     Daa,
	DasName:     Das,
	DecName:     DecReg16, // Representative DEC instruction
	DivName:     Div,
	EnterName:   Enter,
	HltName:     Hlt,
	IdivName:    Idiv,
	ImulName:    Imul,
	InName:      In,
	IncName:     IncReg16, // Representative INC instruction
	InsbName:    Insb,
	InswName:    Insw,
	IntName:     Int,
	IntoName:    Into,
	InvdName:    Invd,
	IretName:    Iret,
	JbName:      Jb,
	JbeName:     Jbe,
	JlName:      Jl,
	JleName:     Jle,
	JmpName:     Jmp,
	JnbName:     Jnb,
	JnbeName:    Jnbe,
	JnlName:     Jnl,
	JnleName:    Jnle,
	JnoName:     Jno,
	JnpName:     Jnp,
	JnsName:     Jns,
	JnzName:     Jnz,
	JoName:      Jo,
	JpName:      Jp,
	JsName:      Js,
	JzName:      Jz,
	LeaName:     Lea,
	LeaveName:   Leave,
	LmswName:    Lmsw,
	LodsbName:   Lodsb,
	LodswName:   Lodsw,
	MovName:     MovRegImm16, // Representative MOV instruction
	MovsbName:   Movsb,
	MovswName:   Movsw,
	MovsxName:   Movsx,
	MovzxName:   Movzx,
	MulName:     Mul,
	NopName:     Nop,
	OrName:      OrALImm8, // Representative OR instruction
	OutName:     Out,
	OutsbName:   Outsb,
	OutswName:   Outsw,
	PopName:     PopReg16, // Representative POP instruction
	PopaName:    Popa,
	PushName:    PushReg16, // Representative PUSH instruction
	PushaName:   Pusha,
	RclName:     Rcl,
	RcrName:     Rcr,
	RepName:     Rep,
	RepnzName:   Repnz,
	RepzName:    Repz,
	RetName:     Ret,
	RetfName:    RetFar,
	RolName:     Rol,
	RorName:     Ror,
	SarName:     Sar,
	SbbName:     SbbALImm8, // Representative SBB instruction
	ScasbName:   Scasb,
	ScaswName:   Scasw,
	SetccName:   Setcc,
	ShldName:    Shld,
	ShlName:     Shl,
	ShrdName:    Shrd,
	ShrName:     Shr,
	SmswName:    Smsw,
	StcName:     Stc,
	StdName:     Std,
	StiName:     Sti,
	StosbName:   Stosb,
	StoswName:   Stosw,
	SubName:     SubALImm8, // Representative SUB instruction
	TestName:    Test,
	WbinvdName:  Wbinvd,
	XaddName:    Xadd,
	XchgName:    Xchg,
	XlatName:    Xlat,
	XorName:     XorALImm8, // Representative XOR instruction

	// Segment override prefixes
	SegCSName: SegCS,
	SegDSName: SegDS,
	SegESName: SegES,
	SegSSName: SegSS,
}
