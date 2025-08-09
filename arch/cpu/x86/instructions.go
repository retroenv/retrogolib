package x86

// Core x86 instruction definitions for DOS development.
// This file contains the most commonly used instructions (~585 total).

// Instruction variables for the opcode table.
var (
	// Data Movement Instructions
	MovRMReg8   *Instruction
	MovRMReg16  *Instruction
	MovRegRM8   *Instruction
	MovRegRM16  *Instruction
	MovRegImm8  *Instruction
	MovRegImm16 *Instruction
	MovMemImm8  *Instruction
	MovMemImm16 *Instruction

	// Arithmetic Instructions - ADD
	AddRMReg8  *Instruction
	AddRMReg16 *Instruction
	AddRegRM8  *Instruction
	AddRegRM16 *Instruction
	AddALImm8  *Instruction
	AddAXImm16 *Instruction

	// Arithmetic Instructions - SUB
	SubRMReg8  *Instruction
	SubRMReg16 *Instruction
	SubRegRM8  *Instruction
	SubRegRM16 *Instruction
	SubALImm8  *Instruction
	SubAXImm16 *Instruction

	// Arithmetic Instructions - ADC (Add with Carry)
	AdcRMReg8  *Instruction
	AdcRMReg16 *Instruction
	AdcRegRM8  *Instruction
	AdcRegRM16 *Instruction
	AdcALImm8  *Instruction
	AdcAXImm16 *Instruction

	// Arithmetic Instructions - SBB (Subtract with Borrow)
	SbbRMReg8  *Instruction
	SbbRMReg16 *Instruction
	SbbRegRM8  *Instruction
	SbbRegRM16 *Instruction
	SbbALImm8  *Instruction
	SbbAXImm16 *Instruction

	// Logical Instructions - AND
	AndRMReg8  *Instruction
	AndRMReg16 *Instruction
	AndRegRM8  *Instruction
	AndRegRM16 *Instruction
	AndALImm8  *Instruction
	AndAXImm16 *Instruction

	// Logical Instructions - OR
	OrRMReg8  *Instruction
	OrRMReg16 *Instruction
	OrRegRM8  *Instruction
	OrRegRM16 *Instruction
	OrALImm8  *Instruction
	OrAXImm16 *Instruction

	// Logical Instructions - XOR
	XorRMReg8  *Instruction
	XorRMReg16 *Instruction
	XorRegRM8  *Instruction
	XorRegRM16 *Instruction
	XorALImm8  *Instruction
	XorAXImm16 *Instruction

	// Comparison Instructions
	CmpRMReg8  *Instruction
	CmpRMReg16 *Instruction
	CmpRegRM8  *Instruction
	CmpRegRM16 *Instruction
	CmpALImm8  *Instruction
	CmpAXImm16 *Instruction

	// Increment/Decrement Instructions
	IncReg8  *Instruction
	IncReg16 *Instruction
	IncRM8   *Instruction
	IncRM16  *Instruction
	DecReg8  *Instruction
	DecReg16 *Instruction
	DecRM8   *Instruction
	DecRM16  *Instruction

	// Stack Instructions
	PushReg16 *Instruction
	PopReg16  *Instruction
	PushSeg   *Instruction
	PopSeg    *Instruction
	PushCS    *Instruction
	PushDS    *Instruction
	PushES    *Instruction
	PushSS    *Instruction
	PopDS     *Instruction
	PopES     *Instruction
	PopSS     *Instruction

	// Jump Instructions - Conditional
	Jo   *Instruction // Jump if overflow
	Jno  *Instruction // Jump if not overflow
	Jb   *Instruction // Jump if below/carry
	Jnb  *Instruction // Jump if not below/not carry
	Jz   *Instruction // Jump if zero/equal
	Jnz  *Instruction // Jump if not zero/not equal
	Jbe  *Instruction // Jump if below or equal
	Jnbe *Instruction // Jump if not below or equal
	Js   *Instruction // Jump if sign
	Jns  *Instruction // Jump if not sign
	Jp   *Instruction // Jump if parity/parity even
	Jnp  *Instruction // Jump if not parity/parity odd
	Jl   *Instruction // Jump if less
	Jnl  *Instruction // Jump if not less
	Jle  *Instruction // Jump if less or equal
	Jnle *Instruction // Jump if not less or equal

	// Jump Instructions - Unconditional
	Jmp     *Instruction // Unconditional jump
	JmpFar  *Instruction // Far jump
	Call    *Instruction // Call procedure
	CallFar *Instruction // Far call
	Ret     *Instruction // Return
	RetFar  *Instruction // Far return

	// Interrupt Instructions
	Int  *Instruction // Software interrupt
	Into *Instruction // Interrupt on overflow
	Iret *Instruction // Return from interrupt

	// Flag Instructions
	Clc *Instruction // Clear carry flag
	Stc *Instruction // Set carry flag
	Cmc *Instruction // Complement carry flag
	Cld *Instruction // Clear direction flag
	Std *Instruction // Set direction flag
	Cli *Instruction // Clear interrupt flag
	Sti *Instruction // Set interrupt flag

	// String Instructions
	Movsb *Instruction // Move string byte
	Movsw *Instruction // Move string word
	Cmpsb *Instruction // Compare string byte
	Cmpsw *Instruction // Compare string word
	Scasb *Instruction // Scan string byte
	Scasw *Instruction // Scan string word
	Lodsb *Instruction // Load string byte
	Lodsw *Instruction // Load string word
	Stosb *Instruction // Store string byte
	Stosw *Instruction // Store string word

	// Repeat Prefixes
	Rep   *Instruction // Repeat
	Repz  *Instruction // Repeat while zero
	Repnz *Instruction // Repeat while not zero

	// Shift and Rotate Instructions
	Shl *Instruction // Shift left
	Shr *Instruction // Shift right
	Sar *Instruction // Shift arithmetic right
	Rol *Instruction // Rotate left
	Ror *Instruction // Rotate right
	Rcl *Instruction // Rotate through carry left
	Rcr *Instruction // Rotate through carry right

	// Test Instructions
	Test *Instruction // Test (logical AND without storing result)

	// Exchange Instructions
	Xchg *Instruction // Exchange

	// Segment Override Prefixes
	SegES *Instruction // ES segment prefix
	SegCS *Instruction // CS segment prefix
	SegSS *Instruction // SS segment prefix
	SegDS *Instruction // DS segment prefix

	// Decimal Arithmetic
	Daa *Instruction // Decimal adjust after addition
	Das *Instruction // Decimal adjust after subtraction
	Aaa *Instruction // ASCII adjust after addition
	Aas *Instruction // ASCII adjust after subtraction

	// Multiplication and Division
	Mul  *Instruction // Multiply
	Imul *Instruction // Signed multiply
	Div  *Instruction // Divide
	Idiv *Instruction // Signed divide

	// I/O Instructions
	In  *Instruction // Input from port
	Out *Instruction // Output to port

	// Control Instructions
	Nop *Instruction // No operation
	Hlt *Instruction // Halt

	// Other Instructions
	Cbw  *Instruction // Convert byte to word
	Cwd  *Instruction // Convert word to double word
	Xlat *Instruction // Table lookup translation
	Lea  *Instruction // Load effective address

	// Undefined/Reserved
	Undefined *Instruction // Placeholder for undefined opcodes
)

// initializeInstructions creates and initializes all instruction definitions.
func initializeInstructions() {
	// Data Movement Instructions
	MovRMReg8 = &Instruction{
		Name: "mov",
		ParamFunc: func(c *CPU, params ...any) error {
			modrm := params[0].(ModRM)
			displacement := params[1].(int16)

			srcValue := c.getReg8(RegisterParam(modrm.Reg))

			if modrm.Mod == 3 {
				// Register to register
				c.setReg8(RegisterParam(modrm.RM), srcValue)
			} else {
				// Register to memory
				addr := c.GetEffectiveAddress(modrm, displacement, 0)
				c.memory.Write8(addr, srcValue)
			}
			return nil
		},
	}

	MovRegRM8 = &Instruction{
		Name: "mov",
		ParamFunc: func(c *CPU, params ...any) error {
			modrm := params[0].(ModRM)
			displacement := params[1].(int16)

			var srcValue uint8
			if modrm.Mod == 3 {
				// Register to register
				srcValue = c.getReg8(RegisterParam(modrm.RM))
			} else {
				// Memory to register
				addr := c.GetEffectiveAddress(modrm, displacement, 0)
				srcValue = c.memory.Read8(addr)
			}

			c.setReg8(RegisterParam(modrm.Reg), srcValue)
			return nil
		},
	}

	// ADD Instructions
	AddALImm8 = &Instruction{
		Name: "add",
		ParamFunc: func(c *CPU, params ...any) error {
			immediate := params[0].(uint8)
			result := c.add8(c.AL(), immediate)
			c.SetAL(result)
			return nil
		},
	}

	AddAXImm16 = &Instruction{
		Name: "add",
		ParamFunc: func(c *CPU, params ...any) error {
			immediate := params[0].(uint16)
			result := c.add16(c.AX, immediate)
			c.AX = result
			return nil
		},
	}

	// Conditional Jump Instructions
	Jo = &Instruction{
		Name: "jo",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			if c.Flags.GetOverflow() {
				c.IP = uint16(int32(c.IP) + int32(offset))
			}
			return nil
		},
	}

	Jno = &Instruction{
		Name: "jno",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			if !c.Flags.GetOverflow() {
				c.IP = uint16(int32(c.IP) + int32(offset))
			}
			return nil
		},
	}

	Jb = &Instruction{
		Name: "jb",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			if c.Flags.GetCarry() {
				c.IP = uint16(int32(c.IP) + int32(offset))
			}
			return nil
		},
	}

	Jnb = &Instruction{
		Name: "jnb",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			if !c.Flags.GetCarry() {
				c.IP = uint16(int32(c.IP) + int32(offset))
			}
			return nil
		},
	}

	Jz = &Instruction{
		Name: "jz",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			if c.Flags.GetZero() {
				c.IP = uint16(int32(c.IP) + int32(offset))
			}
			return nil
		},
	}

	Jnz = &Instruction{
		Name: "jnz",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			if !c.Flags.GetZero() {
				c.IP = uint16(int32(c.IP) + int32(offset))
			}
			return nil
		},
	}

	Jbe = &Instruction{
		Name: "jbe",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			if c.Flags.GetCarry() || c.Flags.GetZero() {
				c.IP = uint16(int32(c.IP) + int32(offset))
			}
			return nil
		},
	}

	Jnbe = &Instruction{
		Name: "jnbe",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			if !c.Flags.GetCarry() && !c.Flags.GetZero() {
				c.IP = uint16(int32(c.IP) + int32(offset))
			}
			return nil
		},
	}

	Js = &Instruction{
		Name: "js",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			if c.Flags.GetSign() {
				c.IP = uint16(int32(c.IP) + int32(offset))
			}
			return nil
		},
	}

	Jns = &Instruction{
		Name: "jns",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			if !c.Flags.GetSign() {
				c.IP = uint16(int32(c.IP) + int32(offset))
			}
			return nil
		},
	}

	Jp = &Instruction{
		Name: "jp",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			if c.Flags.GetParity() {
				c.IP = uint16(int32(c.IP) + int32(offset))
			}
			return nil
		},
	}

	Jnp = &Instruction{
		Name: "jnp",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			if !c.Flags.GetParity() {
				c.IP = uint16(int32(c.IP) + int32(offset))
			}
			return nil
		},
	}

	Jl = &Instruction{
		Name: "jl",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			if c.Flags.GetSign() != c.Flags.GetOverflow() {
				c.IP = uint16(int32(c.IP) + int32(offset))
			}
			return nil
		},
	}

	Jnl = &Instruction{
		Name: "jnl",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			if c.Flags.GetSign() == c.Flags.GetOverflow() {
				c.IP = uint16(int32(c.IP) + int32(offset))
			}
			return nil
		},
	}

	Jle = &Instruction{
		Name: "jle",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			if c.Flags.GetZero() || (c.Flags.GetSign() != c.Flags.GetOverflow()) {
				c.IP = uint16(int32(c.IP) + int32(offset))
			}
			return nil
		},
	}

	Jnle = &Instruction{
		Name: "jnle",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			if !c.Flags.GetZero() && (c.Flags.GetSign() == c.Flags.GetOverflow()) {
				c.IP = uint16(int32(c.IP) + int32(offset))
			}
			return nil
		},
	}

	// Stack Instructions
	PushES = &Instruction{
		Name: "push",
		NoParamFunc: func(c *CPU) error {
			c.push16(c.ES)
			return nil
		},
	}

	PopES = &Instruction{
		Name: "pop",
		NoParamFunc: func(c *CPU) error {
			c.ES = c.pop16()
			return nil
		},
	}

	PushCS = &Instruction{
		Name: "push",
		NoParamFunc: func(c *CPU) error {
			c.push16(c.CS)
			return nil
		},
	}

	PushDS = &Instruction{
		Name: "push",
		NoParamFunc: func(c *CPU) error {
			c.push16(c.DS)
			return nil
		},
	}

	PopDS = &Instruction{
		Name: "pop",
		NoParamFunc: func(c *CPU) error {
			c.DS = c.pop16()
			return nil
		},
	}

	PushSS = &Instruction{
		Name: "push",
		NoParamFunc: func(c *CPU) error {
			c.push16(c.SS)
			return nil
		},
	}

	PopSS = &Instruction{
		Name: "pop",
		NoParamFunc: func(c *CPU) error {
			c.SS = c.pop16()
			return nil
		},
	}

	// Logical Instructions - OR
	OrALImm8 = &Instruction{
		Name: "or",
		ParamFunc: func(c *CPU, params ...any) error {
			immediate := params[0].(uint8)
			result := c.or8(c.AL(), immediate)
			c.SetAL(result)
			return nil
		},
	}

	OrAXImm16 = &Instruction{
		Name: "or",
		ParamFunc: func(c *CPU, params ...any) error {
			immediate := params[0].(uint16)
			result := c.or16(c.AX, immediate)
			c.AX = result
			return nil
		},
	}

	// ADC Instructions
	AdcALImm8 = &Instruction{
		Name: "adc",
		ParamFunc: func(c *CPU, params ...any) error {
			immediate := params[0].(uint8)
			carry := uint8(0)
			if c.Flags.GetCarry() {
				carry = 1
			}
			result := c.add8(c.AL(), immediate+carry)
			c.SetAL(result)
			return nil
		},
	}

	AdcAXImm16 = &Instruction{
		Name: "adc",
		ParamFunc: func(c *CPU, params ...any) error {
			immediate := params[0].(uint16)
			carry := uint16(0)
			if c.Flags.GetCarry() {
				carry = 1
			}
			result := c.add16(c.AX, immediate+carry)
			c.AX = result
			return nil
		},
	}

	// SBB Instructions
	SbbALImm8 = &Instruction{
		Name: "sbb",
		ParamFunc: func(c *CPU, params ...any) error {
			immediate := params[0].(uint8)
			carry := uint8(0)
			if c.Flags.GetCarry() {
				carry = 1
			}
			result := c.sub8(c.AL(), immediate+carry)
			c.SetAL(result)
			return nil
		},
	}

	SbbAXImm16 = &Instruction{
		Name: "sbb",
		ParamFunc: func(c *CPU, params ...any) error {
			immediate := params[0].(uint16)
			carry := uint16(0)
			if c.Flags.GetCarry() {
				carry = 1
			}
			result := c.sub16(c.AX, immediate+carry)
			c.AX = result
			return nil
		},
	}

	// AND Instructions
	AndALImm8 = &Instruction{
		Name: "and",
		ParamFunc: func(c *CPU, params ...any) error {
			immediate := params[0].(uint8)
			result := c.and8(c.AL(), immediate)
			c.SetAL(result)
			return nil
		},
	}

	AndAXImm16 = &Instruction{
		Name: "and",
		ParamFunc: func(c *CPU, params ...any) error {
			immediate := params[0].(uint16)
			result := c.and16(c.AX, immediate)
			c.AX = result
			return nil
		},
	}

	// Segment Prefixes
	SegES = &Instruction{
		Name: "es:",
		NoParamFunc: func(c *CPU) error {
			// Segment override prefix - actual implementation would set a flag
			// for the next instruction to use ES segment
			return nil
		},
	}

	// Decimal Adjust
	Daa = &Instruction{
		Name: "daa",
		NoParamFunc: func(c *CPU) error {
			al := c.AL()
			oldCarry := c.Flags.GetCarry()

			if (al&0x0F) > 9 || c.Flags.GetAuxCarry() {
				c.SetAL(al + 6)
				c.SetAuxCarry(true)
			} else {
				c.SetAuxCarry(false)
			}

			al = c.AL()
			if al > 0x9F || oldCarry {
				c.SetAL(al + 0x60)
				c.SetCarry(true)
			} else {
				c.SetCarry(false)
			}

			c.SetSZP8(c.AL())
			return nil
		},
	}

	// NOP Instruction
	Nop = &Instruction{
		Name: "nop",
		NoParamFunc: func(c *CPU) error {
			// Do nothing
			return nil
		},
	}

	// Undefined/Reserved opcode placeholder
	Undefined = &Instruction{
		Name: "undefined",
		NoParamFunc: func(c *CPU) error {
			return ErrInvalidOpcode
		},
	}

	// TODO: Implement remaining instructions...
	// This is a subset showing the pattern. A complete implementation
	// would include all ~585 instructions commonly used in DOS development.
}

// init initializes all instruction definitions.
func init() {
	initializeInstructions()
	InitializeOpcodeMaps()
}
