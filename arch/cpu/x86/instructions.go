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
	initializeDataMovementInstructions()
	initializeArithmeticInstructions()
	initializeArithmeticMoreInstructions()
	initializeExtendedArithmeticInstructions()
	initializeJumpInstructions()
	initializeStackInstructions()
	initializeLogicalInstructions()
	initializeIncrementDecrementInstructions()
	initializeStringInstructions()
	initializeMiscInstructions()
}

// initializeDataMovementInstructions initializes data movement instructions.
func initializeDataMovementInstructions() {
	initializeMOVRMReg()
	initializeMOVRegRM()
	initializeMOVRegImm()
	initializeMOVMemImm()
}

func initializeMOVRMReg() {
	// MOV r/m8, r8
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

	// MOV r/m16, r16
	MovRMReg16 = &Instruction{
		Name: "mov",
		ParamFunc: func(c *CPU, params ...any) error {
			modrm := params[0].(ModRM)
			displacement := params[1].(int16)

			srcValue := c.getReg16(RegisterParam(modrm.Reg))

			if modrm.Mod == 3 {
				// Register to register
				c.setReg16(RegisterParam(modrm.RM), srcValue)
			} else {
				// Register to memory
				addr := c.GetEffectiveAddress(modrm, displacement, 0)
				c.memory.Write16(addr, srcValue)
			}
			return nil
		},
	}
}

func initializeMOVRegRM() {
	// MOV r8, r/m8
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

	// MOV r16, r/m16
	MovRegRM16 = &Instruction{
		Name: "mov",
		ParamFunc: func(c *CPU, params ...any) error {
			modrm := params[0].(ModRM)
			displacement := params[1].(int16)

			var srcValue uint16
			if modrm.Mod == 3 {
				// Register to register
				srcValue = c.getReg16(RegisterParam(modrm.RM))
			} else {
				// Memory to register
				addr := c.GetEffectiveAddress(modrm, displacement, 0)
				srcValue = c.memory.Read16(addr)
			}

			c.setReg16(RegisterParam(modrm.Reg), srcValue)
			return nil
		},
	}
}

func initializeMOVRegImm() {
	// MOV r8, imm8 (register-specific versions)
	MovRegImm8 = &Instruction{
		Name: "mov",
		ParamFunc: func(c *CPU, params ...any) error {
			_ = params[0].(uint8)
			// Register is determined by opcode
			// This will be called with the appropriate register set
			return nil
		},
	}

	// MOV r16, imm16 (register-specific versions)
	MovRegImm16 = &Instruction{
		Name: "mov",
		ParamFunc: func(c *CPU, params ...any) error {
			_ = params[0].(uint16)
			// Register is determined by opcode
			// This will be called with the appropriate register set
			return nil
		},
	}
}

func initializeMOVMemImm() {
	// MOV r/m8, imm8
	MovMemImm8 = &Instruction{
		Name: "mov",
		ParamFunc: func(c *CPU, params ...any) error {
			modrm := params[0].(ModRM)
			displacement := params[1].(int16)
			immediate := params[2].(uint8)

			if modrm.Mod == 3 {
				// Register
				c.setReg8(RegisterParam(modrm.RM), immediate)
			} else {
				// Memory
				addr := c.GetEffectiveAddress(modrm, displacement, 0)
				c.memory.Write8(addr, immediate)
			}
			return nil
		},
	}

	// MOV r/m16, imm16
	MovMemImm16 = &Instruction{
		Name: "mov",
		ParamFunc: func(c *CPU, params ...any) error {
			modrm := params[0].(ModRM)
			displacement := params[1].(int16)
			immediate := params[2].(uint16)

			if modrm.Mod == 3 {
				// Register
				c.setReg16(RegisterParam(modrm.RM), immediate)
			} else {
				// Memory
				addr := c.GetEffectiveAddress(modrm, displacement, 0)
				c.memory.Write16(addr, immediate)
			}
			return nil
		},
	}
}

// initializeArithmeticInstructions initializes arithmetic instructions (ADD, SUB, ADC, SBB).
func initializeArithmeticInstructions() {
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
}

// initializeArithmeticMoreInstructions initializes remaining arithmetic instructions.
func initializeArithmeticMoreInstructions() {
	initializeSubInstructions()
	initializeCmpInstructions()
}

func initializeSubInstructions() {
	// SUB Instructions (additional variants)
	SubRMReg8 = &Instruction{
		Name: "sub",
		ParamFunc: func(c *CPU, params ...any) error {
			modrm := params[0].(ModRM)
			displacement := params[1].(int16)

			srcValue := c.getReg8(RegisterParam(modrm.Reg))
			if modrm.Mod == 3 {
				dstValue := c.getReg8(RegisterParam(modrm.RM))
				result := c.sub8(dstValue, srcValue)
				c.setReg8(RegisterParam(modrm.RM), result)
			} else {
				addr := c.GetEffectiveAddress(modrm, displacement, 0)
				dstValue := c.memory.Read8(addr)
				result := c.sub8(dstValue, srcValue)
				c.memory.Write8(addr, result)
			}
			return nil
		},
	}

	SubRMReg16 = &Instruction{
		Name: "sub",
		ParamFunc: func(c *CPU, params ...any) error {
			modrm := params[0].(ModRM)
			displacement := params[1].(int16)

			srcValue := c.getReg16(RegisterParam(modrm.Reg))
			if modrm.Mod == 3 {
				dstValue := c.getReg16(RegisterParam(modrm.RM))
				result := c.sub16(dstValue, srcValue)
				c.setReg16(RegisterParam(modrm.RM), result)
			} else {
				addr := c.GetEffectiveAddress(modrm, displacement, 0)
				dstValue := c.memory.Read16(addr)
				result := c.sub16(dstValue, srcValue)
				c.memory.Write16(addr, result)
			}
			return nil
		},
	}

	SubRegRM8 = &Instruction{
		Name: "sub",
		ParamFunc: func(c *CPU, params ...any) error {
			modrm := params[0].(ModRM)
			displacement := params[1].(int16)

			dstValue := c.getReg8(RegisterParam(modrm.Reg))
			var srcValue uint8
			if modrm.Mod == 3 {
				srcValue = c.getReg8(RegisterParam(modrm.RM))
			} else {
				addr := c.GetEffectiveAddress(modrm, displacement, 0)
				srcValue = c.memory.Read8(addr)
			}

			result := c.sub8(dstValue, srcValue)
			c.setReg8(RegisterParam(modrm.Reg), result)
			return nil
		},
	}

	SubRegRM16 = &Instruction{
		Name: "sub",
		ParamFunc: func(c *CPU, params ...any) error {
			modrm := params[0].(ModRM)
			displacement := params[1].(int16)

			dstValue := c.getReg16(RegisterParam(modrm.Reg))
			var srcValue uint16
			if modrm.Mod == 3 {
				srcValue = c.getReg16(RegisterParam(modrm.RM))
			} else {
				addr := c.GetEffectiveAddress(modrm, displacement, 0)
				srcValue = c.memory.Read16(addr)
			}

			result := c.sub16(dstValue, srcValue)
			c.setReg16(RegisterParam(modrm.Reg), result)
			return nil
		},
	}
}

func initializeCmpInstructions() {
	// CMP Instructions (additional variants)
	CmpRMReg8 = &Instruction{
		Name: "cmp",
		ParamFunc: func(c *CPU, params ...any) error {
			modrm := params[0].(ModRM)
			displacement := params[1].(int16)

			srcValue := c.getReg8(RegisterParam(modrm.Reg))
			if modrm.Mod == 3 {
				dstValue := c.getReg8(RegisterParam(modrm.RM))
				_ = c.sub8(dstValue, srcValue) // Sets flags only
			} else {
				addr := c.GetEffectiveAddress(modrm, displacement, 0)
				dstValue := c.memory.Read8(addr)
				_ = c.sub8(dstValue, srcValue) // Sets flags only
			}
			return nil
		},
	}

	CmpRMReg16 = &Instruction{
		Name: "cmp",
		ParamFunc: func(c *CPU, params ...any) error {
			modrm := params[0].(ModRM)
			displacement := params[1].(int16)

			srcValue := c.getReg16(RegisterParam(modrm.Reg))
			if modrm.Mod == 3 {
				dstValue := c.getReg16(RegisterParam(modrm.RM))
				_ = c.sub16(dstValue, srcValue) // Sets flags only
			} else {
				addr := c.GetEffectiveAddress(modrm, displacement, 0)
				dstValue := c.memory.Read16(addr)
				_ = c.sub16(dstValue, srcValue) // Sets flags only
			}
			return nil
		},
	}

	CmpRegRM8 = &Instruction{
		Name: "cmp",
		ParamFunc: func(c *CPU, params ...any) error {
			modrm := params[0].(ModRM)
			displacement := params[1].(int16)

			dstValue := c.getReg8(RegisterParam(modrm.Reg))
			var srcValue uint8
			if modrm.Mod == 3 {
				srcValue = c.getReg8(RegisterParam(modrm.RM))
			} else {
				addr := c.GetEffectiveAddress(modrm, displacement, 0)
				srcValue = c.memory.Read8(addr)
			}

			_ = c.sub8(dstValue, srcValue) // Sets flags only
			return nil
		},
	}

	CmpRegRM16 = &Instruction{
		Name: "cmp",
		ParamFunc: func(c *CPU, params ...any) error {
			modrm := params[0].(ModRM)
			displacement := params[1].(int16)

			dstValue := c.getReg16(RegisterParam(modrm.Reg))
			var srcValue uint16
			if modrm.Mod == 3 {
				srcValue = c.getReg16(RegisterParam(modrm.RM))
			} else {
				addr := c.GetEffectiveAddress(modrm, displacement, 0)
				srcValue = c.memory.Read16(addr)
			}

			_ = c.sub16(dstValue, srcValue) // Sets flags only
			return nil
		},
	}
}

// initializeExtendedArithmeticInstructions initializes ADC and SBB instructions.
func initializeExtendedArithmeticInstructions() {
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
}

// makeConditionalJump creates a conditional jump instruction with the given condition function.
func makeConditionalJump(name string, condition func(*CPU) bool) *Instruction {
	return &Instruction{
		Name: name,
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			if condition(c) {
				c.IP = uint16(int32(c.IP) + int32(offset))
			}
			return nil
		},
	}
}

// initializeJumpInstructions initializes jump and control flow instructions.
func initializeJumpInstructions() {
	initializeConditionalJumps()
	initializeUnconditionalJumps()
}

func initializeConditionalJumps() {
	// Conditional Jump Instructions
	Jo = makeConditionalJump("jo", func(c *CPU) bool { return c.Flags.GetOverflow() })
	Jno = makeConditionalJump("jno", func(c *CPU) bool { return !c.Flags.GetOverflow() })
	Jb = makeConditionalJump("jb", func(c *CPU) bool { return c.Flags.GetCarry() })
	Jnb = makeConditionalJump("jnb", func(c *CPU) bool { return !c.Flags.GetCarry() })
	Jz = makeConditionalJump("jz", func(c *CPU) bool { return c.Flags.GetZero() })
	Jnz = makeConditionalJump("jnz", func(c *CPU) bool { return !c.Flags.GetZero() })
	Jbe = makeConditionalJump("jbe", func(c *CPU) bool { return c.Flags.GetCarry() || c.Flags.GetZero() })
	Jnbe = makeConditionalJump("jnbe", func(c *CPU) bool { return !c.Flags.GetCarry() && !c.Flags.GetZero() })
	Js = makeConditionalJump("js", func(c *CPU) bool { return c.Flags.GetSign() })
	Jns = makeConditionalJump("jns", func(c *CPU) bool { return !c.Flags.GetSign() })
	Jp = makeConditionalJump("jp", func(c *CPU) bool { return c.Flags.GetParity() })
	Jnp = makeConditionalJump("jnp", func(c *CPU) bool { return !c.Flags.GetParity() })
	Jl = makeConditionalJump("jl", func(c *CPU) bool { return c.Flags.GetSign() != c.Flags.GetOverflow() })
	Jnl = makeConditionalJump("jnl", func(c *CPU) bool { return c.Flags.GetSign() == c.Flags.GetOverflow() })
	Jle = makeConditionalJump("jle", func(c *CPU) bool { return c.Flags.GetZero() || (c.Flags.GetSign() != c.Flags.GetOverflow()) })
	Jnle = makeConditionalJump("jnle", func(c *CPU) bool { return !c.Flags.GetZero() && (c.Flags.GetSign() == c.Flags.GetOverflow()) })
}

func initializeUnconditionalJumps() {
	// Unconditional Jump Instructions
	Jmp = &Instruction{
		Name: "jmp",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			c.IP = uint16(int32(c.IP) + int32(offset))
			return nil
		},
	}

	JmpFar = &Instruction{
		Name: "jmp",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(uint16)
			segment := params[1].(uint16)
			c.IP = offset
			c.CS = segment
			return nil
		},
	}

	// Call Instructions
	Call = &Instruction{
		Name: "call",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			c.push16(c.IP)
			c.IP = uint16(int32(c.IP) + int32(offset))
			return nil
		},
	}

	CallFar = &Instruction{
		Name: "call",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(uint16)
			segment := params[1].(uint16)
			c.push16(c.CS)
			c.push16(c.IP)
			c.IP = offset
			c.CS = segment
			return nil
		},
	}

	// Return Instructions
	Ret = &Instruction{
		Name: "ret",
		ParamFunc: func(c *CPU, params ...any) error {
			if len(params) > 0 {
				// RET imm16 - pop return address and adjust stack
				imm := params[0].(uint16)
				c.IP = c.pop16()
				c.SP += imm
			} else {
				// RET - pop return address
				c.IP = c.pop16()
			}
			return nil
		},
	}

	RetFar = &Instruction{
		Name: "retf",
		ParamFunc: func(c *CPU, params ...any) error {
			if len(params) > 0 {
				// RETF imm16 - pop return address and segment, adjust stack
				imm := params[0].(uint16)
				c.IP = c.pop16()
				c.CS = c.pop16()
				c.SP += imm
			} else {
				// RETF - pop return address and segment
				c.IP = c.pop16()
				c.CS = c.pop16()
			}
			return nil
		},
	}
}

// initializeStackInstructions initializes stack operation instructions.
func initializeStackInstructions() {
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
}

// initializeLogicalInstructions initializes logical operation instructions.
func initializeLogicalInstructions() {
	initializeORInstructions()
	initializeANDInstructions()
	initializeXORInstructions()
	initializeLogicalSUBCMPInstructions()
}

func initializeORInstructions() {
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
}

func initializeANDInstructions() {
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
}

func initializeXORInstructions() {
	// XOR Instructions
	XorALImm8 = &Instruction{
		Name: "xor",
		ParamFunc: func(c *CPU, params ...any) error {
			immediate := params[0].(uint8)
			result := c.AL() ^ immediate
			c.SetAL(result)
			c.SetCarry(false)
			c.SetOverflow(false)
			c.SetSZP8(result)
			return nil
		},
	}

	XorAXImm16 = &Instruction{
		Name: "xor",
		ParamFunc: func(c *CPU, params ...any) error {
			immediate := params[0].(uint16)
			result := c.AX ^ immediate
			c.AX = result
			c.SetCarry(false)
			c.SetOverflow(false)
			c.SetSZP16(result)
			return nil
		},
	}
}

func initializeLogicalSUBCMPInstructions() {
	// SUB Instructions
	SubALImm8 = &Instruction{
		Name: "sub",
		ParamFunc: func(c *CPU, params ...any) error {
			immediate := params[0].(uint8)
			result := c.sub8(c.AL(), immediate)
			c.SetAL(result)
			return nil
		},
	}

	SubAXImm16 = &Instruction{
		Name: "sub",
		ParamFunc: func(c *CPU, params ...any) error {
			immediate := params[0].(uint16)
			result := c.sub16(c.AX, immediate)
			c.AX = result
			return nil
		},
	}

	// CMP Instructions
	CmpALImm8 = &Instruction{
		Name: "cmp",
		ParamFunc: func(c *CPU, params ...any) error {
			immediate := params[0].(uint8)
			_ = c.sub8(c.AL(), immediate) // Sets flags but doesn't store result
			return nil
		},
	}

	CmpAXImm16 = &Instruction{
		Name: "cmp",
		ParamFunc: func(c *CPU, params ...any) error {
			immediate := params[0].(uint16)
			_ = c.sub16(c.AX, immediate) // Sets flags but doesn't store result
			return nil
		},
	}
}

// initializeIncrementDecrementInstructions initializes INC/DEC instructions.
func initializeIncrementDecrementInstructions() {
	initializeINCInstructions()
	initializeDECInstructions()
	initializePUSHPOPInstructions()
}

func initializeINCInstructions() {
	// INC register instructions
	IncReg16 = &Instruction{
		Name: "inc",
		ParamFunc: func(c *CPU, params ...any) error {
			// Register determined by opcode
			return nil
		},
	}

	IncReg8 = &Instruction{
		Name: "inc",
		ParamFunc: func(c *CPU, params ...any) error {
			// Register determined by opcode
			return nil
		},
	}

	IncRM8 = &Instruction{
		Name: "inc",
		ParamFunc: func(c *CPU, params ...any) error {
			modrm := params[0].(ModRM)
			displacement := params[1].(int16)

			if modrm.Mod == 3 {
				value := c.getReg8(RegisterParam(modrm.RM))
				result := value + 1
				c.setReg8(RegisterParam(modrm.RM), result)
				c.SetOverflow(value == 0x7F)
				c.SetSZP8(result)
				c.SetAuxCarry((value & 0x0F) == 0x0F)
			} else {
				addr := c.GetEffectiveAddress(modrm, displacement, 0)
				value := c.memory.Read8(addr)
				result := value + 1
				c.memory.Write8(addr, result)
				c.SetOverflow(value == 0x7F)
				c.SetSZP8(result)
				c.SetAuxCarry((value & 0x0F) == 0x0F)
			}
			return nil
		},
	}

	IncRM16 = &Instruction{
		Name: "inc",
		ParamFunc: func(c *CPU, params ...any) error {
			modrm := params[0].(ModRM)
			displacement := params[1].(int16)

			if modrm.Mod == 3 {
				value := c.getReg16(RegisterParam(modrm.RM))
				result := value + 1
				c.setReg16(RegisterParam(modrm.RM), result)
				c.SetOverflow(value == 0x7FFF)
				c.SetSZP16(result)
				c.SetAuxCarry((value & 0x0F) == 0x0F)
			} else {
				addr := c.GetEffectiveAddress(modrm, displacement, 0)
				value := c.memory.Read16(addr)
				result := value + 1
				c.memory.Write16(addr, result)
				c.SetOverflow(value == 0x7FFF)
				c.SetSZP16(result)
				c.SetAuxCarry((value & 0x0F) == 0x0F)
			}
			return nil
		},
	}
}

func initializeDECInstructions() {
	// DEC register instructions
	DecReg16 = &Instruction{
		Name: "dec",
		ParamFunc: func(c *CPU, params ...any) error {
			// Register determined by opcode
			return nil
		},
	}

	DecReg8 = &Instruction{
		Name: "dec",
		ParamFunc: func(c *CPU, params ...any) error {
			// Register determined by opcode
			return nil
		},
	}

	DecRM8 = &Instruction{
		Name: "dec",
		ParamFunc: func(c *CPU, params ...any) error {
			modrm := params[0].(ModRM)
			displacement := params[1].(int16)

			if modrm.Mod == 3 {
				value := c.getReg8(RegisterParam(modrm.RM))
				result := value - 1
				c.setReg8(RegisterParam(modrm.RM), result)
				c.SetOverflow(value == 0x80)
				c.SetSZP8(result)
				c.SetAuxCarry((value & 0x0F) == 0x00)
			} else {
				addr := c.GetEffectiveAddress(modrm, displacement, 0)
				value := c.memory.Read8(addr)
				result := value - 1
				c.memory.Write8(addr, result)
				c.SetOverflow(value == 0x80)
				c.SetSZP8(result)
				c.SetAuxCarry((value & 0x0F) == 0x00)
			}
			return nil
		},
	}

	DecRM16 = &Instruction{
		Name: "dec",
		ParamFunc: func(c *CPU, params ...any) error {
			modrm := params[0].(ModRM)
			displacement := params[1].(int16)

			if modrm.Mod == 3 {
				value := c.getReg16(RegisterParam(modrm.RM))
				result := value - 1
				c.setReg16(RegisterParam(modrm.RM), result)
				c.SetOverflow(value == 0x8000)
				c.SetSZP16(result)
				c.SetAuxCarry((value & 0x0F) == 0x00)
			} else {
				addr := c.GetEffectiveAddress(modrm, displacement, 0)
				value := c.memory.Read16(addr)
				result := value - 1
				c.memory.Write16(addr, result)
				c.SetOverflow(value == 0x8000)
				c.SetSZP16(result)
				c.SetAuxCarry((value & 0x0F) == 0x00)
			}
			return nil
		},
	}
}

func initializePUSHPOPInstructions() {
	// PUSH/POP register instructions
	PushReg16 = &Instruction{
		Name: "push",
		ParamFunc: func(c *CPU, params ...any) error {
			// Register determined by opcode
			return nil
		},
	}

	PopReg16 = &Instruction{
		Name: "pop",
		ParamFunc: func(c *CPU, params ...any) error {
			// Register determined by opcode
			return nil
		},
	}
}

// initializeStringInstructions initializes string operation instructions.
func initializeStringInstructions() {
	initializeCompareStringInstructions()
	initializeScanStringInstructions()
	initializeLoadStringInstructions()
	initializeStoreStringInstructions()
	initializeRepeatPrefixes()
}

func initializeCompareStringInstructions() {
	// Compare String instructions
	Cmpsb = &Instruction{
		Name: "cmpsb",
		NoParamFunc: func(c *CPU) error {
			srcAddr := c.CalculateAddress(c.DS, c.SI)
			dstAddr := c.CalculateAddress(c.ES, c.DI)
			srcValue := c.memory.Read8(srcAddr)
			dstValue := c.memory.Read8(dstAddr)
			_ = c.sub8(dstValue, srcValue) // Sets flags
			if c.Flags.GetDirection() {
				c.SI--
				c.DI--
			} else {
				c.SI++
				c.DI++
			}
			return nil
		},
	}

	Cmpsw = &Instruction{
		Name: "cmpsw",
		NoParamFunc: func(c *CPU) error {
			srcAddr := c.CalculateAddress(c.DS, c.SI)
			dstAddr := c.CalculateAddress(c.ES, c.DI)
			srcValue := c.memory.Read16(srcAddr)
			dstValue := c.memory.Read16(dstAddr)
			_ = c.sub16(dstValue, srcValue) // Sets flags
			if c.Flags.GetDirection() {
				c.SI -= 2
				c.DI -= 2
			} else {
				c.SI += 2
				c.DI += 2
			}
			return nil
		},
	}
}

func initializeScanStringInstructions() {
	// Scan String instructions
	Scasb = &Instruction{
		Name: "scasb",
		NoParamFunc: func(c *CPU) error {
			dstAddr := c.CalculateAddress(c.ES, c.DI)
			dstValue := c.memory.Read8(dstAddr)
			_ = c.sub8(c.AL(), dstValue) // Sets flags
			if c.Flags.GetDirection() {
				c.DI--
			} else {
				c.DI++
			}
			return nil
		},
	}

	Scasw = &Instruction{
		Name: "scasw",
		NoParamFunc: func(c *CPU) error {
			dstAddr := c.CalculateAddress(c.ES, c.DI)
			dstValue := c.memory.Read16(dstAddr)
			_ = c.sub16(c.AX, dstValue) // Sets flags
			if c.Flags.GetDirection() {
				c.DI -= 2
			} else {
				c.DI += 2
			}
			return nil
		},
	}
}

func initializeLoadStringInstructions() {
	// Load String instructions
	Lodsb = &Instruction{
		Name: "lodsb",
		NoParamFunc: func(c *CPU) error {
			srcAddr := c.CalculateAddress(c.DS, c.SI)
			value := c.memory.Read8(srcAddr)
			c.SetAL(value)
			if c.Flags.GetDirection() {
				c.SI--
			} else {
				c.SI++
			}
			return nil
		},
	}

	Lodsw = &Instruction{
		Name: "lodsw",
		NoParamFunc: func(c *CPU) error {
			srcAddr := c.CalculateAddress(c.DS, c.SI)
			value := c.memory.Read16(srcAddr)
			c.AX = value
			if c.Flags.GetDirection() {
				c.SI -= 2
			} else {
				c.SI += 2
			}
			return nil
		},
	}
}

func initializeStoreStringInstructions() {
	// Store String instructions
	Stosb = &Instruction{
		Name: "stosb",
		NoParamFunc: func(c *CPU) error {
			dstAddr := c.CalculateAddress(c.ES, c.DI)
			c.memory.Write8(dstAddr, c.AL())
			if c.Flags.GetDirection() {
				c.DI--
			} else {
				c.DI++
			}
			return nil
		},
	}

	Stosw = &Instruction{
		Name: "stosw",
		NoParamFunc: func(c *CPU) error {
			dstAddr := c.CalculateAddress(c.ES, c.DI)
			c.memory.Write16(dstAddr, c.AX)
			if c.Flags.GetDirection() {
				c.DI -= 2
			} else {
				c.DI += 2
			}
			return nil
		},
	}
}

func initializeRepeatPrefixes() {
	// Repeat prefixes
	Rep = &Instruction{
		Name: "rep",
		NoParamFunc: func(c *CPU) error {
			// Repeat prefix - implementation depends on following instruction
			return nil
		},
	}

	Repz = &Instruction{
		Name: "repz",
		NoParamFunc: func(c *CPU) error {
			// Repeat while zero prefix - implementation depends on following instruction
			return nil
		},
	}

	Repnz = &Instruction{
		Name: "repnz",
		NoParamFunc: func(c *CPU) error {
			// Repeat while not zero prefix - implementation depends on following instruction
			return nil
		},
	}
}

// initializeMiscInstructions initializes miscellaneous instructions.
func initializeMiscInstructions() {
	initializeSegmentPrefixes()
	initializeASCIIArithmetic()
	initializeIOInstructions()
	initializeNOPAndReserved()
	initializeFlagInstructions()
	initializeInterruptInstructions()
}

func initializeSegmentPrefixes() {
	// Segment Prefixes
	SegES = &Instruction{
		Name: "es:",
		NoParamFunc: func(c *CPU) error {
			// Segment override prefix - actual implementation would set a flag
			// for the next instruction to use ES segment
			return nil
		},
	}

	// Additional segment prefixes
	SegCS = &Instruction{
		Name: "cs:",
		NoParamFunc: func(c *CPU) error {
			// CS segment override prefix
			return nil
		},
	}

	SegSS = &Instruction{
		Name: "ss:",
		NoParamFunc: func(c *CPU) error {
			// SS segment override prefix
			return nil
		},
	}

	SegDS = &Instruction{
		Name: "ds:",
		NoParamFunc: func(c *CPU) error {
			// DS segment override prefix
			return nil
		},
	}
}

func initializeASCIIArithmetic() {
	initializeAAInstructions()
	initializeDASInstruction()
}

func initializeAAInstructions() {
	// ASCII Arithmetic Instructions - AAA/AAS
	Aaa = &Instruction{
		Name: "aaa",
		NoParamFunc: func(c *CPU) error {
			al := c.AL()
			if (al&0x0F) > 9 || c.Flags.GetAuxCarry() {
				c.SetAL((al + 6) & 0x0F)
				c.SetAH(c.AH() + 1)
				c.SetCarry(true)
				c.SetAuxCarry(true)
			} else {
				c.SetAL(al & 0x0F)
				c.SetCarry(false)
				c.SetAuxCarry(false)
			}
			return nil
		},
	}

	Aas = &Instruction{
		Name: "aas",
		NoParamFunc: func(c *CPU) error {
			al := c.AL()
			if (al&0x0F) > 9 || c.Flags.GetAuxCarry() {
				c.SetAL((al - 6) & 0x0F)
				c.SetAH(c.AH() - 1)
				c.SetCarry(true)
				c.SetAuxCarry(true)
			} else {
				c.SetAL(al & 0x0F)
				c.SetCarry(false)
				c.SetAuxCarry(false)
			}
			return nil
		},
	}
}

func initializeDASInstruction() {
	Das = &Instruction{
		Name: "das",
		NoParamFunc: func(c *CPU) error {
			al := c.AL()
			oldCarry := c.Flags.GetCarry()

			if (al&0x0F) > 9 || c.Flags.GetAuxCarry() {
				c.SetAL(al - 6)
				c.SetAuxCarry(true)
			} else {
				c.SetAuxCarry(false)
			}

			al = c.AL()
			if al > 0x9F || oldCarry {
				c.SetAL(al - 0x60)
				c.SetCarry(true)
			} else {
				c.SetCarry(false)
			}

			c.SetSZP8(c.AL())
			return nil
		},
	}
}

func initializeIOInstructions() {
	// I/O Instructions
	In = &Instruction{
		Name: "in",
		ParamFunc: func(c *CPU, params ...any) error {
			// I/O port input - implementation would depend on system
			// For now, just return 0
			if len(params) > 0 {
				// IN AL/AX, imm8 or IN AL/AX, DX
				// Set AL or AX to 0 for now
				c.SetAL(0)
			}
			return nil
		},
	}

	Out = &Instruction{
		Name: "out",
		ParamFunc: func(c *CPU, params ...any) error {
			// I/O port output - implementation would depend on system
			// For now, do nothing
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
}

func initializeNOPAndReserved() {
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
}

func initializeFlagInstructions() {
	// Flag Instructions
	Clc = &Instruction{
		Name: "clc",
		NoParamFunc: func(c *CPU) error {
			c.SetCarry(false)
			return nil
		},
	}

	Stc = &Instruction{
		Name: "stc",
		NoParamFunc: func(c *CPU) error {
			c.SetCarry(true)
			return nil
		},
	}

	Cmc = &Instruction{
		Name: "cmc",
		NoParamFunc: func(c *CPU) error {
			c.SetCarry(!c.Flags.GetCarry())
			return nil
		},
	}

	Cld = &Instruction{
		Name: "cld",
		NoParamFunc: func(c *CPU) error {
			c.SetDirection(false)
			return nil
		},
	}

	Std = &Instruction{
		Name: "std",
		NoParamFunc: func(c *CPU) error {
			c.SetDirection(true)
			return nil
		},
	}

	Cli = &Instruction{
		Name: "cli",
		NoParamFunc: func(c *CPU) error {
			c.SetInterrupt(false)
			return nil
		},
	}

	Sti = &Instruction{
		Name: "sti",
		NoParamFunc: func(c *CPU) error {
			c.SetInterrupt(true)
			return nil
		},
	}
}

func initializeInterruptInstructions() {
	initializeBasicInterrupts()
	initializeConversionInstructions()
	initializeRemainingStringInstructions()
	initializeMiscellaneousRemaining()
}

func initializeBasicInterrupts() {
	// Interrupt Instructions
	Int = &Instruction{
		Name: "int",
		ParamFunc: func(c *CPU, params ...any) error {
			vector := params[0].(uint8)
			c.TriggerInterrupt(vector)
			return nil
		},
	}

	Into = &Instruction{
		Name: "into",
		NoParamFunc: func(c *CPU) error {
			if c.Flags.GetOverflow() {
				c.TriggerInterrupt(4) // Interrupt 4 for overflow
			}
			return nil
		},
	}

	Iret = &Instruction{
		Name: "iret",
		NoParamFunc: func(c *CPU) error {
			c.IP = c.pop16()
			c.CS = c.pop16()
			flags := c.pop16()
			c.Flags = Flags(flags)
			return nil
		},
	}
}

func initializeConversionInstructions() {
	// Conversion Instructions
	Cbw = &Instruction{
		Name: "cbw",
		NoParamFunc: func(c *CPU) error {
			// Convert byte in AL to word in AX (sign extend)
			if c.AL()&0x80 != 0 {
				c.SetAH(0xFF)
			} else {
				c.SetAH(0x00)
			}
			return nil
		},
	}

	Cwd = &Instruction{
		Name: "cwd",
		NoParamFunc: func(c *CPU) error {
			// Convert word in AX to double word in DX:AX (sign extend)
			if c.AX&0x8000 != 0 {
				c.DX = 0xFFFF
			} else {
				c.DX = 0x0000
			}
			return nil
		},
	}
}

func initializeRemainingStringInstructions() {
	// String Instructions
	Movsb = &Instruction{
		Name: "movsb",
		NoParamFunc: func(c *CPU) error {
			srcAddr := c.CalculateAddress(c.DS, c.SI)
			dstAddr := c.CalculateAddress(c.ES, c.DI)
			value := c.memory.Read8(srcAddr)
			c.memory.Write8(dstAddr, value)
			if c.Flags.GetDirection() {
				c.SI--
				c.DI--
			} else {
				c.SI++
				c.DI++
			}
			return nil
		},
	}

	Movsw = &Instruction{
		Name: "movsw",
		NoParamFunc: func(c *CPU) error {
			srcAddr := c.CalculateAddress(c.DS, c.SI)
			dstAddr := c.CalculateAddress(c.ES, c.DI)
			value := c.memory.Read16(srcAddr)
			c.memory.Write16(dstAddr, value)
			if c.Flags.GetDirection() {
				c.SI -= 2
				c.DI -= 2
			} else {
				c.SI += 2
				c.DI += 2
			}
			return nil
		},
	}
}

func initializeMiscellaneousRemaining() {
	// Miscellaneous Instructions
	Xlat = &Instruction{
		Name: "xlat",
		NoParamFunc: func(c *CPU) error {
			addr := c.CalculateAddress(c.DS, c.BX+uint16(c.AL()))
			c.SetAL(c.memory.Read8(addr))
			return nil
		},
	}

	Lea = &Instruction{
		Name: "lea",
		ParamFunc: func(c *CPU, params ...any) error {
			modrm := params[0].(ModRM)
			displacement := params[1].(int16)
			offset, _ := c.calculateOffset(modrm, displacement, 0)
			c.setReg16(RegisterParam(modrm.Reg), offset)
			return nil
		},
	}

	// Halt Instruction
	Hlt = &Instruction{
		Name: "hlt",
		NoParamFunc: func(c *CPU) error {
			c.Halt()
			return nil
		},
	}

	// TEST Instruction
	Test = &Instruction{
		Name: "test",
		ParamFunc: func(c *CPU, params ...any) error {
			// TEST performs AND operation but doesn't store result
			if len(params) == 1 {
				// TEST AL/AX, imm
				imm := params[0]
				switch v := imm.(type) {
				case uint8:
					_ = c.and8(c.AL(), v)
				case uint16:
					_ = c.and16(c.AX, v)
				}
			} else {
				// TEST r/m, r with ModR/M
				modrm := params[0].(ModRM)
				displacement := params[1].(int16)
				regValue := c.getReg16(RegisterParam(modrm.Reg))
				var memValue uint16
				if modrm.Mod == 3 {
					memValue = c.getReg16(RegisterParam(modrm.RM))
				} else {
					addr := c.GetEffectiveAddress(modrm, displacement, 0)
					memValue = c.memory.Read16(addr)
				}
				_ = c.and16(regValue, memValue)
			}
			return nil
		},
	}

	// Exchange Instruction
	Xchg = &Instruction{
		Name: "xchg",
		ParamFunc: func(c *CPU, params ...any) error {
			// XCHG AX, r16 (single-byte opcodes)
			if len(params) == 0 {
				// NOP (XCHG AX, AX)
				return nil
			}
			// Implementation depends on parameters
			return nil
		},
	}
}

// init initializes all instruction definitions.
func init() {
	initializeInstructions()
	InitializeOpcodeMaps()
}
