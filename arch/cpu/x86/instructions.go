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
	}

	SubRMReg16 = &Instruction{
		Name: "sub",
	}

	SubRegRM8 = &Instruction{
		Name: "sub",
	}
	SubRegRM16 = &Instruction{
		Name: "sub",
	}
	SubALImm8 = &Instruction{
		Name: "sub",
	}
	SubAXImm16 = &Instruction{
		Name: "sub",
	}

	// Arithmetic Instructions - ADC (Add with Carry)
	AdcRMReg8 = &Instruction{
		Name: "adc",
	}
	AdcRMReg16 = &Instruction{
		Name: "adc",
	}
	AdcRegRM8 = &Instruction{
		Name: "adc",
	}
	AdcRegRM16 = &Instruction{
		Name: "adc",
	}
	AdcALImm8 = &Instruction{
		Name: "adc",
	}

	AdcAXImm16 = &Instruction{
		Name: "adc",
	}

	// Arithmetic Instructions - SBB (Subtract with Borrow)
	SbbRMReg8 = &Instruction{
		Name: "sbb",
	}
	SbbRMReg16 = &Instruction{
		Name: "sbb",
	}
	SbbRegRM8 = &Instruction{
		Name: "sbb",
	}
	SbbRegRM16 = &Instruction{
		Name: "sbb",
	}
	SbbALImm8 = &Instruction{
		Name: "sbb",
	}
	SbbAXImm16 = &Instruction{
		Name: "sbb",
	}

	// Logical Instructions - AND
	AndRMReg8 = &Instruction{
		Name: "and",
	}
	AndRMReg16 = &Instruction{
		Name: "and",
	}
	AndRegRM8 = &Instruction{
		Name: "and",
	}
	AndRegRM16 = &Instruction{
		Name: "and",
	}
	AndALImm8 = &Instruction{
		Name: "and",
	}
	AndAXImm16 = &Instruction{
		Name: "and",
	}

	// Logical Instructions - OR
	OrRMReg8 = &Instruction{
		Name: "or",
	}
	OrRMReg16 = &Instruction{
		Name: "or",
	}
	OrRegRM8 = &Instruction{
		Name: "or",
	}
	OrRegRM16 = &Instruction{
		Name: "or",
	}
	OrALImm8 = &Instruction{
		Name: "or",
	}
	OrAXImm16 = &Instruction{
		Name: "or",
	}

	// Logical Instructions - XOR
	XorRMReg8 = &Instruction{
		Name: "xor",
	}
	XorRMReg16 = &Instruction{
		Name: "xor",
	}
	XorRegRM8 = &Instruction{
		Name: "xor",
	}
	XorRegRM16 = &Instruction{
		Name: "xor",
	}
	XorALImm8 = &Instruction{
		Name: "xor",
	}
	XorAXImm16 = &Instruction{
		Name: "xor",
	}

	// Comparison Instructions
	CmpRMReg8 = &Instruction{
		Name: "cmp",
	}

	CmpRMReg16 = &Instruction{
		Name: "cmp",
	}

	CmpRegRM8 = &Instruction{
		Name: "cmp",
	}
	CmpRegRM16 = &Instruction{
		Name: "cmp",
	}
	CmpALImm8 = &Instruction{
		Name: "cmp",
	}
	CmpAXImm16 = &Instruction{
		Name: "cmp",
	}

	// Increment/Decrement Instructions
	IncReg8 = &Instruction{
		Name: "inc",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0xFE, Size: 2, Cycles: 2, HasModRM: true}, // Group 4
		},
	}
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
	DecReg8 = &Instruction{
		Name: "dec",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0xFE, Size: 2, Cycles: 3, HasModRM: true}, // Group 4, /1
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
	DecRM8 = &Instruction{
		Name: "dec",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0xFE, Size: 2, Cycles: 3, HasModRM: true}, // Group 4, /1
		},
	}
	DecRM16 = &Instruction{
		Name: "dec",
		Addressing: map[AddressingMode]OpcodeInfo{
			ModRMRegisterAddressing: {Opcode: 0xFF, Size: 2, Cycles: 3, HasModRM: true}, // Group 5, /1
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
	PushSeg = &Instruction{
		Name: "push",
	}
	PopSeg = &Instruction{
		Name: "pop",
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
	}
	PopES = &Instruction{
		Name: "pop",
	}
	PopSS = &Instruction{
		Name: "pop",
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
	}
	Jnb = &Instruction{ // Jump if not below/not carry
		Name: "jnb",
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
	}
	Jnbe = &Instruction{ // Jump if not below or equal
		Name: "jnbe",
	}
	Js = &Instruction{ // Jump if sign
		Name: "js",
	}
	Jns = &Instruction{ // Jump if not sign
		Name: "jns",
	}
	Jp = &Instruction{ // Jump if parity/parity even
		Name: "jp",
	}
	Jnp = &Instruction{ // Jump if not parity/parity odd
		Name: "jnp",
	}
	Jl = &Instruction{ // Jump if less
		Name: "jl",
	}
	Jnl = &Instruction{ // Jump if not less
		Name: "jnl",
	}
	Jle = &Instruction{ // Jump if less or equal
		Name: "jle",
	}
	Jnle = &Instruction{ // Jump if not less or equal
		Name: "jnle",
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
	}

	// Interrupt Instructions
	Int = &Instruction{ // Software interrupt
		Name: "int",
	}
	Into = &Instruction{ // Interrupt on overflow
		Name: "into",
	}
	Iret = &Instruction{ // Return from interrupt
		Name: "iret",
	}

	// Flag Instructions
	Clc = &Instruction{ // Clear carry flag
		Name: "clc",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0xF8, Size: 1, Cycles: 2, HasModRM: false},
		},
	}
	Stc = &Instruction{ // Set carry flag
		Name: "stc",
		Addressing: map[AddressingMode]OpcodeInfo{
			ImpliedAddressing: {Opcode: 0xF9, Size: 1, Cycles: 2, HasModRM: false},
		},
	}
	Cmc = &Instruction{ // Complement carry flag
		Name: "cmc",
	}
	Cld = &Instruction{ // Clear direction flag
		Name: "cld",
	}
	Std = &Instruction{ // Set direction flag
		Name: "std",
	}
	Cli = &Instruction{ // Clear interrupt flag
		Name: "cli",
	}
	Sti = &Instruction{ // Set interrupt flag
		Name: "sti",
	}

	// String Instructions
	Movsb = &Instruction{ // Move string byte
		Name: "movsb",
	}
	Movsw = &Instruction{ // Move string word
		Name: "movsw",
	}
	Cmpsb = &Instruction{ // Compare string byte
		Name: "cmpsb",
	}
	Cmpsw = &Instruction{ // Compare string word
		Name: "cmpsw",
	}
	Scasb = &Instruction{ // Scan string byte
		Name: "scasb",
	}
	Scasw = &Instruction{ // Scan string word
		Name: "scasw",
	}
	Lodsb = &Instruction{ // Load string byte
		Name: "lodsb",
	}
	Lodsw = &Instruction{ // Load string word
		Name: "lodsw",
	}
	Stosb = &Instruction{ // Store string byte
		Name: "stosb",
	}
	Stosw = &Instruction{ // Store string word
		Name: "stosw",
	}

	// Repeat Prefixes
	Rep = &Instruction{ // Repeat
		Name: "rep",
	}
	Repz = &Instruction{ // Repeat while zero
		Name: "repz",
	}
	Repnz = &Instruction{ // Repeat while not zero
		Name: "repnz",
	}

	// Shift and Rotate Instructions
	Shl = &Instruction{ // Shift left
		Name: "shl",
	}
	Shr = &Instruction{ // Shift right
		Name: "shr",
	}
	Sar = &Instruction{ // Shift arithmetic right
		Name: "sar",
	}
	Rol = &Instruction{ // Rotate left
		Name: "rol",
	}
	Ror = &Instruction{ // Rotate right
		Name: "ror",
	}
	Rcl = &Instruction{ // Rotate through carry left
		Name: "rcl",
	}
	Rcr = &Instruction{ // Rotate through carry right
		Name: "rcr",
	}

	// Test Instructions
	Test = &Instruction{ // Test (logical AND without storing result)
		Name: "test",
	}

	// Exchange Instructions
	Xchg = &Instruction{ // Exchange
		Name: "xchg",
	}

	// Segment Override Prefixes
	SegES = &Instruction{ // ES segment prefix
		Name: "es:",
	}
	SegCS = &Instruction{ // CS segment prefix
		Name: "cs:",
	}
	SegSS = &Instruction{ // SS segment prefix
		Name: "ss:",
	}
	SegDS = &Instruction{ // DS segment prefix
		Name: "ds:",
	}

	// Decimal Arithmetic
	Daa = &Instruction{ // Decimal adjust after addition
		Name: "daa",
	}
	Das = &Instruction{ // Decimal adjust after subtraction
		Name: "das",
	}
	Aaa = &Instruction{ // ASCII adjust after addition
		Name: "aaa",
	}
	Aas = &Instruction{ // ASCII adjust after subtraction
		Name: "aas",
	}

	// Multiplication and Division
	Mul = &Instruction{ // Multiply
		Name: "mul",
	}
	Imul = &Instruction{ // Signed multiply
		Name: "imul",
	}
	Div = &Instruction{ // Divide
		Name: "div",
	}
	Idiv = &Instruction{ // Signed divide
		Name: "idiv",
	}

	// I/O Instructions
	In = &Instruction{ // Input from port
		Name: "in",
	}
	Out = &Instruction{ // Output to port
		Name: "out",
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
	}
	Cwd = &Instruction{ // Convert word to double word
		Name: "cwd",
	}
	Xlat = &Instruction{ // Table lookup translation
		Name: "xlat",
	}
	Lea = &Instruction{ // Load effective address
		Name: "lea",
	}

	// Undefined/Reserved
	Undefined = &Instruction{ // Placeholder for undefined opcodes
		Name: "undefined",
	}
)

// init initializes all instruction definitions.
func init() {
	InitializeOpcodeMaps()
}
