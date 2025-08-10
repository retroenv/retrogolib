package x86

// Core x86 instruction definitions for DOS development.
// This file contains the most commonly used instructions (~585 total).

// Instruction variables for the opcode table.
var (
	// Data Movement Instructions
	MovRMReg8 = &Instruction{
		Name: "mov",
	}

	MovRMReg16 = &Instruction{
		Name: "mov",
	}

	MovRegRM8 = &Instruction{
		Name: "mov",
	}

	MovRegRM16 = &Instruction{
		Name: "mov",
	}

	MovRegImm8 = &Instruction{
		Name: "mov",
	}

	MovRegImm16 = &Instruction{
		Name: "mov",
	}

	MovMemImm8 = &Instruction{
		Name: "mov",
	}

	MovMemImm16 = &Instruction{
		Name: "mov",
	}

	// Arithmetic Instructions - ADD
	AddRMReg8 = &Instruction{
		Name: "add",
	}
	AddRMReg16 = &Instruction{
		Name: "add",
	}
	AddRegRM8 = &Instruction{
		Name: "add",
	}
	AddRegRM16 = &Instruction{
		Name: "add",
	}
	AddALImm8 = &Instruction{
		Name: "add",
	}

	AddAXImm16 = &Instruction{
		Name: "add",
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
	}
	IncReg16 = &Instruction{
		Name: "inc",
	}
	IncRM8 = &Instruction{
		Name: "inc",
	}
	IncRM16 = &Instruction{
		Name: "inc",
	}
	DecReg8 = &Instruction{
		Name: "dec",
	}
	DecReg16 = &Instruction{
		Name: "dec",
	}
	DecRM8 = &Instruction{
		Name: "dec",
	}
	DecRM16 = &Instruction{
		Name: "dec",
	}

	// Stack Instructions
	PushReg16 = &Instruction{
		Name: "push",
	}
	PopReg16 = &Instruction{
		Name: "pop",
	}
	PushSeg = &Instruction{
		Name: "push",
	}
	PopSeg = &Instruction{
		Name: "pop",
	}
	PushCS = &Instruction{
		Name: "push",
	}
	PushDS = &Instruction{
		Name: "push",
	}
	PushES = &Instruction{
		Name: "push",
	}
	PushSS = &Instruction{
		Name: "push",
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
	}
	Jno = &Instruction{ // Jump if not overflow
		Name: "jno",
	}
	Jb = &Instruction{ // Jump if below/carry
		Name: "jb",
	}
	Jnb = &Instruction{ // Jump if not below/not carry
		Name: "jnb",
	}
	Jz = &Instruction{ // Jump if zero/equal
		Name: "jz",
	}
	Jnz = &Instruction{ // Jump if not zero/not equal
		Name: "jnz",
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
	}
	JmpFar = &Instruction{ // Far jump
		Name: "jmp",
	}
	Call = &Instruction{ // Call procedure
		Name: "call",
	}
	CallFar = &Instruction{ // Far call
		Name: "call",
	}
	Ret = &Instruction{ // Return
		Name: "ret",
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
	}
	Stc = &Instruction{ // Set carry flag
		Name: "stc",
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
	}
	Hlt = &Instruction{ // Halt
		Name: "hlt",
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
