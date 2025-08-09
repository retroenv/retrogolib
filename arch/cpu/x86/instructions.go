package x86

// Core x86 instruction definitions for DOS development.
// This file contains the most commonly used instructions (~585 total).

// Instruction variables for the opcode table.
var (
	// Data Movement Instructions
	MovRMReg8 = &Instruction{
		Name:      "mov",
		ParamFunc: movRMReg8,
	}

	MovRMReg16 = &Instruction{
		Name:      "mov",
		ParamFunc: movRMReg16,
	}

	MovRegRM8 = &Instruction{
		Name:      "mov",
		ParamFunc: movRegRM8,
	}

	MovRegRM16 = &Instruction{
		Name:      "mov",
		ParamFunc: movRegRM16,
	}

	MovRegImm8 = &Instruction{
		Name:      "mov",
		ParamFunc: movRegImm8,
	}

	MovRegImm16 = &Instruction{
		Name:      "mov",
		ParamFunc: movRegImm16,
	}

	MovMemImm8 = &Instruction{
		Name:      "mov",
		ParamFunc: movMemImm8,
	}

	MovMemImm16 = &Instruction{
		Name:      "mov",
		ParamFunc: movMemImm16,
	}

	// Arithmetic Instructions - ADD
	AddRMReg8 = &Instruction{
		Name: "add",
		ParamFunc: func(c *CPU, params ...any) error {
			// ADD r/m8, r8 implementation
			return nil
		},
	}
	AddRMReg16 = &Instruction{
		Name: "add",
		ParamFunc: func(c *CPU, params ...any) error {
			// ADD r/m16, r16 implementation
			return nil
		},
	}
	AddRegRM8 = &Instruction{
		Name: "add",
		ParamFunc: func(c *CPU, params ...any) error {
			// ADD r8, r/m8 implementation
			return nil
		},
	}
	AddRegRM16 = &Instruction{
		Name: "add",
		ParamFunc: func(c *CPU, params ...any) error {
			// ADD r16, r/m16 implementation
			return nil
		},
	}
	AddALImm8 = &Instruction{
		Name:      "add",
		ParamFunc: addALImm8,
	}

	AddAXImm16 = &Instruction{
		Name:      "add",
		ParamFunc: addAXImm16,
	}

	// Arithmetic Instructions - SUB
	SubRMReg8 = &Instruction{
		Name:      "sub",
		ParamFunc: subRMReg8,
	}

	SubRMReg16 = &Instruction{
		Name:      "sub",
		ParamFunc: subRMReg16,
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

	// Arithmetic Instructions - ADC (Add with Carry)
	AdcRMReg8 = &Instruction{
		Name: "adc",
		ParamFunc: func(c *CPU, params ...any) error {
			// ADC r/m8, r8 implementation
			return nil
		},
	}
	AdcRMReg16 = &Instruction{
		Name: "adc",
		ParamFunc: func(c *CPU, params ...any) error {
			// ADC r/m16, r16 implementation
			return nil
		},
	}
	AdcRegRM8 = &Instruction{
		Name: "adc",
		ParamFunc: func(c *CPU, params ...any) error {
			// ADC r8, r/m8 implementation
			return nil
		},
	}
	AdcRegRM16 = &Instruction{
		Name: "adc",
		ParamFunc: func(c *CPU, params ...any) error {
			// ADC r16, r/m16 implementation
			return nil
		},
	}
	AdcALImm8 = &Instruction{
		Name:      "adc",
		ParamFunc: adcALImm8,
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

	// Arithmetic Instructions - SBB (Subtract with Borrow)
	SbbRMReg8 = &Instruction{
		Name: "sbb",
		ParamFunc: func(c *CPU, params ...any) error {
			// SBB r/m8, r8 implementation
			return nil
		},
	}
	SbbRMReg16 = &Instruction{
		Name: "sbb",
		ParamFunc: func(c *CPU, params ...any) error {
			// SBB r/m16, r16 implementation
			return nil
		},
	}
	SbbRegRM8 = &Instruction{
		Name: "sbb",
		ParamFunc: func(c *CPU, params ...any) error {
			// SBB r8, r/m8 implementation
			return nil
		},
	}
	SbbRegRM16 = &Instruction{
		Name: "sbb",
		ParamFunc: func(c *CPU, params ...any) error {
			// SBB r16, r/m16 implementation
			return nil
		},
	}
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

	// Logical Instructions - AND
	AndRMReg8 = &Instruction{
		Name: "and",
		ParamFunc: func(c *CPU, params ...any) error {
			// AND r/m8, r8 implementation
			return nil
		},
	}
	AndRMReg16 = &Instruction{
		Name: "and",
		ParamFunc: func(c *CPU, params ...any) error {
			// AND r/m16, r16 implementation
			return nil
		},
	}
	AndRegRM8 = &Instruction{
		Name: "and",
		ParamFunc: func(c *CPU, params ...any) error {
			// AND r8, r/m8 implementation
			return nil
		},
	}
	AndRegRM16 = &Instruction{
		Name: "and",
		ParamFunc: func(c *CPU, params ...any) error {
			// AND r16, r/m16 implementation
			return nil
		},
	}
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

	// Logical Instructions - OR
	OrRMReg8 = &Instruction{
		Name: "or",
		ParamFunc: func(c *CPU, params ...any) error {
			// OR r/m8, r8 implementation
			return nil
		},
	}
	OrRMReg16 = &Instruction{
		Name: "or",
		ParamFunc: func(c *CPU, params ...any) error {
			// OR r/m16, r16 implementation
			return nil
		},
	}
	OrRegRM8 = &Instruction{
		Name: "or",
		ParamFunc: func(c *CPU, params ...any) error {
			// OR r8, r/m8 implementation
			return nil
		},
	}
	OrRegRM16 = &Instruction{
		Name: "or",
		ParamFunc: func(c *CPU, params ...any) error {
			// OR r16, r/m16 implementation
			return nil
		},
	}
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

	// Logical Instructions - XOR
	XorRMReg8 = &Instruction{
		Name: "xor",
		ParamFunc: func(c *CPU, params ...any) error {
			// XOR r/m8, r8 implementation
			return nil
		},
	}
	XorRMReg16 = &Instruction{
		Name: "xor",
		ParamFunc: func(c *CPU, params ...any) error {
			// XOR r/m16, r16 implementation
			return nil
		},
	}
	XorRegRM8 = &Instruction{
		Name: "xor",
		ParamFunc: func(c *CPU, params ...any) error {
			// XOR r8, r/m8 implementation
			return nil
		},
	}
	XorRegRM16 = &Instruction{
		Name: "xor",
		ParamFunc: func(c *CPU, params ...any) error {
			// XOR r16, r/m16 implementation
			return nil
		},
	}
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

	// Comparison Instructions
	CmpRMReg8 = &Instruction{
		Name:      "cmp",
		ParamFunc: cmpRMReg8,
	}

	CmpRMReg16 = &Instruction{
		Name:      "cmp",
		ParamFunc: cmpRMReg16,
	}

	CmpRegRM8 = &Instruction{
		Name: "cmp",
		ParamFunc: func(c *CPU, params ...any) error {
			// CMP r8, r/m8 implementation
			return nil
		},
	}
	CmpRegRM16 = &Instruction{
		Name: "cmp",
		ParamFunc: func(c *CPU, params ...any) error {
			// CMP r16, r/m16 implementation
			return nil
		},
	}
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

	// Increment/Decrement Instructions
	IncReg8 = &Instruction{
		Name: "inc",
		ParamFunc: func(c *CPU, params ...any) error {
			// Register determined by opcode
			return nil
		},
	}
	IncReg16 = &Instruction{
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
	DecReg8 = &Instruction{
		Name: "dec",
		ParamFunc: func(c *CPU, params ...any) error {
			// Register determined by opcode
			return nil
		},
	}
	DecReg16 = &Instruction{
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

	// Stack Instructions
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
	PushSeg = &Instruction{
		Name: "push",
		ParamFunc: func(c *CPU, params ...any) error {
			// Segment determined by opcode
			return nil
		},
	}
	PopSeg = &Instruction{
		Name: "pop",
		ParamFunc: func(c *CPU, params ...any) error {
			// Segment determined by opcode
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
	PushES = &Instruction{
		Name: "push",
		NoParamFunc: func(c *CPU) error {
			c.push16(c.ES)
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
	PopDS = &Instruction{
		Name: "pop",
		NoParamFunc: func(c *CPU) error {
			c.DS = c.pop16()
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
	PopSS = &Instruction{
		Name: "pop",
		NoParamFunc: func(c *CPU) error {
			c.SS = c.pop16()
			return nil
		},
	}

	// Jump Instructions - Conditional
	Jo = &Instruction{ // Jump if overflow
		Name: "jo",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			if c.Flags.GetOverflow() {
				c.IP = uint16(int32(c.IP) + int32(offset))
			}
			return nil
		},
	}
	Jno = &Instruction{ // Jump if not overflow
		Name: "jno",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			if !c.Flags.GetOverflow() {
				c.IP = uint16(int32(c.IP) + int32(offset))
			}
			return nil
		},
	}
	Jb = &Instruction{ // Jump if below/carry
		Name: "jb",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			if c.Flags.GetCarry() {
				c.IP = uint16(int32(c.IP) + int32(offset))
			}
			return nil
		},
	}
	Jnb = &Instruction{ // Jump if not below/not carry
		Name: "jnb",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			if !c.Flags.GetCarry() {
				c.IP = uint16(int32(c.IP) + int32(offset))
			}
			return nil
		},
	}
	Jz = &Instruction{ // Jump if zero/equal
		Name: "jz",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			if c.Flags.GetZero() {
				c.IP = uint16(int32(c.IP) + int32(offset))
			}
			return nil
		},
	}
	Jnz = &Instruction{ // Jump if not zero/not equal
		Name: "jnz",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			if !c.Flags.GetZero() {
				c.IP = uint16(int32(c.IP) + int32(offset))
			}
			return nil
		},
	}
	Jbe = &Instruction{ // Jump if below or equal
		Name: "jbe",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			if c.Flags.GetCarry() || c.Flags.GetZero() {
				c.IP = uint16(int32(c.IP) + int32(offset))
			}
			return nil
		},
	}
	Jnbe = &Instruction{ // Jump if not below or equal
		Name: "jnbe",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			if !c.Flags.GetCarry() && !c.Flags.GetZero() {
				c.IP = uint16(int32(c.IP) + int32(offset))
			}
			return nil
		},
	}
	Js = &Instruction{ // Jump if sign
		Name: "js",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			if c.Flags.GetSign() {
				c.IP = uint16(int32(c.IP) + int32(offset))
			}
			return nil
		},
	}
	Jns = &Instruction{ // Jump if not sign
		Name: "jns",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			if !c.Flags.GetSign() {
				c.IP = uint16(int32(c.IP) + int32(offset))
			}
			return nil
		},
	}
	Jp = &Instruction{ // Jump if parity/parity even
		Name: "jp",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			if c.Flags.GetParity() {
				c.IP = uint16(int32(c.IP) + int32(offset))
			}
			return nil
		},
	}
	Jnp = &Instruction{ // Jump if not parity/parity odd
		Name: "jnp",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			if !c.Flags.GetParity() {
				c.IP = uint16(int32(c.IP) + int32(offset))
			}
			return nil
		},
	}
	Jl = &Instruction{ // Jump if less
		Name: "jl",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			if c.Flags.GetSign() != c.Flags.GetOverflow() {
				c.IP = uint16(int32(c.IP) + int32(offset))
			}
			return nil
		},
	}
	Jnl = &Instruction{ // Jump if not less
		Name: "jnl",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			if c.Flags.GetSign() == c.Flags.GetOverflow() {
				c.IP = uint16(int32(c.IP) + int32(offset))
			}
			return nil
		},
	}
	Jle = &Instruction{ // Jump if less or equal
		Name: "jle",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			if c.Flags.GetZero() || (c.Flags.GetSign() != c.Flags.GetOverflow()) {
				c.IP = uint16(int32(c.IP) + int32(offset))
			}
			return nil
		},
	}
	Jnle = &Instruction{ // Jump if not less or equal
		Name: "jnle",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			if !c.Flags.GetZero() && (c.Flags.GetSign() == c.Flags.GetOverflow()) {
				c.IP = uint16(int32(c.IP) + int32(offset))
			}
			return nil
		},
	}

	// Jump Instructions - Unconditional
	Jmp = &Instruction{ // Unconditional jump
		Name: "jmp",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			c.IP = uint16(int32(c.IP) + int32(offset))
			return nil
		},
	}
	JmpFar = &Instruction{ // Far jump
		Name: "jmp",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(uint16)
			segment := params[1].(uint16)
			c.IP = offset
			c.CS = segment
			return nil
		},
	}
	Call = &Instruction{ // Call procedure
		Name: "call",
		ParamFunc: func(c *CPU, params ...any) error {
			offset := params[0].(int16)
			c.push16(c.IP)
			c.IP = uint16(int32(c.IP) + int32(offset))
			return nil
		},
	}
	CallFar = &Instruction{ // Far call
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
	Ret = &Instruction{ // Return
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
	RetFar = &Instruction{ // Far return
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

	// Interrupt Instructions
	Int = &Instruction{ // Software interrupt
		Name: "int",
		ParamFunc: func(c *CPU, params ...any) error {
			vector := params[0].(uint8)
			c.TriggerInterrupt(vector)
			return nil
		},
	}
	Into = &Instruction{ // Interrupt on overflow
		Name: "into",
		NoParamFunc: func(c *CPU) error {
			if c.Flags.GetOverflow() {
				c.TriggerInterrupt(4) // Interrupt 4 for overflow
			}
			return nil
		},
	}
	Iret = &Instruction{ // Return from interrupt
		Name: "iret",
		NoParamFunc: func(c *CPU) error {
			c.IP = c.pop16()
			c.CS = c.pop16()
			flags := c.pop16()
			c.Flags = Flags(flags)
			return nil
		},
	}

	// Flag Instructions
	Clc = &Instruction{ // Clear carry flag
		Name: "clc",
		NoParamFunc: func(c *CPU) error {
			c.SetCarry(false)
			return nil
		},
	}
	Stc = &Instruction{ // Set carry flag
		Name: "stc",
		NoParamFunc: func(c *CPU) error {
			c.SetCarry(true)
			return nil
		},
	}
	Cmc = &Instruction{ // Complement carry flag
		Name: "cmc",
		NoParamFunc: func(c *CPU) error {
			c.SetCarry(!c.Flags.GetCarry())
			return nil
		},
	}
	Cld = &Instruction{ // Clear direction flag
		Name: "cld",
		NoParamFunc: func(c *CPU) error {
			c.SetDirection(false)
			return nil
		},
	}
	Std = &Instruction{ // Set direction flag
		Name: "std",
		NoParamFunc: func(c *CPU) error {
			c.SetDirection(true)
			return nil
		},
	}
	Cli = &Instruction{ // Clear interrupt flag
		Name: "cli",
		NoParamFunc: func(c *CPU) error {
			c.SetInterrupt(false)
			return nil
		},
	}
	Sti = &Instruction{ // Set interrupt flag
		Name: "sti",
		NoParamFunc: func(c *CPU) error {
			c.SetInterrupt(true)
			return nil
		},
	}

	// String Instructions
	Movsb = &Instruction{ // Move string byte
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
	Movsw = &Instruction{ // Move string word
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
	Cmpsb = &Instruction{ // Compare string byte
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
	Cmpsw = &Instruction{ // Compare string word
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
	Scasb = &Instruction{ // Scan string byte
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
	Scasw = &Instruction{ // Scan string word
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
	Lodsb = &Instruction{ // Load string byte
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
	Lodsw = &Instruction{ // Load string word
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
	Stosb = &Instruction{ // Store string byte
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
	Stosw = &Instruction{ // Store string word
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

	// Repeat Prefixes
	Rep = &Instruction{ // Repeat
		Name: "rep",
		NoParamFunc: func(c *CPU) error {
			// Repeat prefix - implementation depends on following instruction
			return nil
		},
	}
	Repz = &Instruction{ // Repeat while zero
		Name: "repz",
		NoParamFunc: func(c *CPU) error {
			// Repeat while zero prefix - implementation depends on following instruction
			return nil
		},
	}
	Repnz = &Instruction{ // Repeat while not zero
		Name: "repnz",
		NoParamFunc: func(c *CPU) error {
			// Repeat while not zero prefix - implementation depends on following instruction
			return nil
		},
	}

	// Shift and Rotate Instructions
	Shl = &Instruction{ // Shift left
		Name: "shl",
		ParamFunc: func(c *CPU, params ...any) error {
			// Shift left implementation
			return nil
		},
	}
	Shr = &Instruction{ // Shift right
		Name: "shr",
		ParamFunc: func(c *CPU, params ...any) error {
			// Shift right implementation
			return nil
		},
	}
	Sar = &Instruction{ // Shift arithmetic right
		Name: "sar",
		ParamFunc: func(c *CPU, params ...any) error {
			// Shift arithmetic right implementation
			return nil
		},
	}
	Rol = &Instruction{ // Rotate left
		Name: "rol",
		ParamFunc: func(c *CPU, params ...any) error {
			// Rotate left implementation
			return nil
		},
	}
	Ror = &Instruction{ // Rotate right
		Name: "ror",
		ParamFunc: func(c *CPU, params ...any) error {
			// Rotate right implementation
			return nil
		},
	}
	Rcl = &Instruction{ // Rotate through carry left
		Name: "rcl",
		ParamFunc: func(c *CPU, params ...any) error {
			// Rotate through carry left implementation
			return nil
		},
	}
	Rcr = &Instruction{ // Rotate through carry right
		Name: "rcr",
		ParamFunc: func(c *CPU, params ...any) error {
			// Rotate through carry right implementation
			return nil
		},
	}

	// Test Instructions
	Test = &Instruction{ // Test (logical AND without storing result)
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

	// Exchange Instructions
	Xchg = &Instruction{ // Exchange
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

	// Segment Override Prefixes
	SegES = &Instruction{ // ES segment prefix
		Name: "es:",
		NoParamFunc: func(c *CPU) error {
			// Segment override prefix - actual implementation would set a flag
			// for the next instruction to use ES segment
			return nil
		},
	}
	SegCS = &Instruction{ // CS segment prefix
		Name: "cs:",
		NoParamFunc: func(c *CPU) error {
			// CS segment override prefix
			return nil
		},
	}
	SegSS = &Instruction{ // SS segment prefix
		Name: "ss:",
		NoParamFunc: func(c *CPU) error {
			// SS segment override prefix
			return nil
		},
	}
	SegDS = &Instruction{ // DS segment prefix
		Name: "ds:",
		NoParamFunc: func(c *CPU) error {
			// DS segment override prefix
			return nil
		},
	}

	// Decimal Arithmetic
	Daa = &Instruction{ // Decimal adjust after addition
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
	Das = &Instruction{ // Decimal adjust after subtraction
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
	Aaa = &Instruction{ // ASCII adjust after addition
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
	Aas = &Instruction{ // ASCII adjust after subtraction
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

	// Multiplication and Division
	Mul = &Instruction{ // Multiply
		Name: "mul",
		ParamFunc: func(c *CPU, params ...any) error {
			// Multiply implementation
			return nil
		},
	}
	Imul = &Instruction{ // Signed multiply
		Name: "imul",
		ParamFunc: func(c *CPU, params ...any) error {
			// Signed multiply implementation
			return nil
		},
	}
	Div = &Instruction{ // Divide
		Name: "div",
		ParamFunc: func(c *CPU, params ...any) error {
			// Divide implementation
			return nil
		},
	}
	Idiv = &Instruction{ // Signed divide
		Name: "idiv",
		ParamFunc: func(c *CPU, params ...any) error {
			// Signed divide implementation
			return nil
		},
	}

	// I/O Instructions
	In = &Instruction{ // Input from port
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
	Out = &Instruction{ // Output to port
		Name: "out",
		ParamFunc: func(c *CPU, params ...any) error {
			// I/O port output - implementation would depend on system
			// For now, do nothing
			return nil
		},
	}

	// Control Instructions
	Nop = &Instruction{ // No operation
		Name: "nop",
		NoParamFunc: func(c *CPU) error {
			// Do nothing
			return nil
		},
	}
	Hlt = &Instruction{ // Halt
		Name: "hlt",
		NoParamFunc: func(c *CPU) error {
			c.Halt()
			return nil
		},
	}

	// Other Instructions
	Cbw = &Instruction{ // Convert byte to word
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
	Cwd = &Instruction{ // Convert word to double word
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
	Xlat = &Instruction{ // Table lookup translation
		Name: "xlat",
		NoParamFunc: func(c *CPU) error {
			addr := c.CalculateAddress(c.DS, c.BX+uint16(c.AL()))
			c.SetAL(c.memory.Read8(addr))
			return nil
		},
	}
	Lea = &Instruction{ // Load effective address
		Name: "lea",
		ParamFunc: func(c *CPU, params ...any) error {
			modrm := params[0].(ModRM)
			displacement := params[1].(int16)
			offset, _ := c.calculateOffset(modrm, displacement, 0)
			c.setReg16(RegisterParam(modrm.Reg), offset)
			return nil
		},
	}

	// Undefined/Reserved
	Undefined = &Instruction{ // Placeholder for undefined opcodes
		Name: "undefined",
		NoParamFunc: func(c *CPU) error {
			return ErrInvalidOpcode
		},
	}
)

// init initializes all instruction definitions.
func init() {
	InitializeOpcodeMaps()
}
