package x86

import (
	"fmt"
)

// Step executes a single instruction and returns any error.
func (c *CPU) Step() error {
	if c.halted {
		return nil
	}

	// Check for interrupts
	if c.triggerInt && c.interruptsEnabled {
		c.handleInterrupt()
		c.triggerInt = false
	}

	// Fetch instruction
	addr := c.GetCSIP()
	opcode := c.memory.Read8(addr)
	c.currentOpcode = opcode
	c.IP++

	// Trace step preparation
	if c.opts.tracing {
		c.TraceStep = c.prepareTraceStep(addr, opcode)
	}

	// Decode and execute instruction
	opcodeInfo, exists := GetOpcodeInfo(opcode)
	if !exists {
		return fmt.Errorf("invalid opcode 0x%02X at %04X:%04X", opcode, c.CS, c.IP-1)
	}

	// Execute instruction
	err := c.executeInstruction(opcodeInfo)
	if err != nil {
		return fmt.Errorf("instruction execution failed: %w", err)
	}

	// Update cycles
	c.cycles += uint64(opcodeInfo.Timing)

	// Complete trace step
	if c.opts.tracing {
		c.completeTraceStep()
		if c.opts.tracingCallback != nil {
			c.opts.tracingCallback(c.TraceStep)
		}
	}

	return nil
}

// executeInstruction executes the given instruction.
func (c *CPU) executeInstruction(opcodeInfo Opcode) error {
	instruction := opcodeInfo.Instruction

	// Handle ModR/M instructions
	if opcodeInfo.HasModRM {
		return c.executeModRMInstruction(instruction, opcodeInfo)
	}

	// Handle immediate instructions
	if opcodeInfo.Addressing == ImmediateAddressing {
		return c.executeImmediateInstruction(instruction, opcodeInfo)
	}

	// Handle implied instructions (no operands)
	if opcodeInfo.Addressing == ImpliedAddressing {
		if instruction.NoParamFunc != nil {
			return instruction.NoParamFunc(c)
		}
		return instruction.Execute(c)
	}

	// Handle relative addressing (jumps)
	if opcodeInfo.Addressing == RelativeAddressing {
		return c.executeRelativeInstruction(instruction, opcodeInfo)
	}

	return fmt.Errorf("unsupported addressing mode: %s", opcodeInfo.Addressing)
}

// executeModRMInstruction executes instructions that use ModR/M bytes.
func (c *CPU) executeModRMInstruction(instruction *Instruction, _ Opcode) error {
	// Fetch ModR/M byte
	modrmByte := c.memory.Read8(c.GetCSIP())
	c.IP++

	var modrm ModRM
	modrm.FromByte(modrmByte)

	// Fetch displacement if needed
	var displacement int16
	if modrm.Mod == 1 {
		// 8-bit displacement
		displacement = int16(int8(c.memory.Read8(c.GetCSIP())))
		c.IP++
	} else if modrm.Mod == 2 || (modrm.Mod == 0 && modrm.RM == 6) {
		// 16-bit displacement or direct addressing
		displacement = int16(c.memory.Read16(c.GetCSIP()))
		c.IP += 2
	}

	// Execute instruction with ModR/M parameters
	if instruction.ParamFunc != nil {
		return instruction.ParamFunc(c, modrm, displacement)
	}

	return ErrInvalidInstruction
}

// executeImmediateInstruction executes instructions with immediate operands.
func (c *CPU) executeImmediateInstruction(instruction *Instruction, opcodeInfo Opcode) error {
	var immediate any

	// Determine immediate size
	switch opcodeInfo.Size {
	case 2: // 8-bit immediate
		immediate = c.memory.Read8(c.GetCSIP())
		c.IP++
	case 3: // 16-bit immediate
		immediate = c.memory.Read16(c.GetCSIP())
		c.IP += 2
	default:
		return fmt.Errorf("invalid immediate instruction size: %d", opcodeInfo.Size)
	}

	if instruction.ParamFunc != nil {
		return instruction.ParamFunc(c, immediate)
	}

	return ErrInvalidInstruction
}

// executeRelativeInstruction executes relative jump instructions.
func (c *CPU) executeRelativeInstruction(instruction *Instruction, opcodeInfo Opcode) error {
	// Read relative offset
	var offset int16
	switch opcodeInfo.Size {
	case 2:
		// 8-bit relative offset (sign-extended)
		offset = int16(int8(c.memory.Read8(c.GetCSIP())))
		c.IP++
	case 3:
		// 16-bit relative offset
		offset = int16(c.memory.Read16(c.GetCSIP()))
		c.IP += 2
	default:
		return fmt.Errorf("invalid relative instruction size: %d", opcodeInfo.Size)
	}

	if instruction.ParamFunc != nil {
		return instruction.ParamFunc(c, offset)
	}

	return ErrInvalidInstruction
}

// handleInterrupt processes interrupt requests.
func (c *CPU) handleInterrupt() {
	if !c.interruptsEnabled {
		return
	}

	// Save flags and return address
	c.push16(c.GetFlags())
	c.push16(c.CS)
	c.push16(c.IP)

	// Disable interrupts
	c.SetInterrupt(false)

	// Load interrupt vector
	vectorAddr := uint32(c.intVector) * 4
	c.IP = c.memory.Read16(vectorAddr)
	c.CS = c.memory.Read16(vectorAddr + 2)
}

// Arithmetic operations

// add8 adds two 8-bit values and sets flags.
func (c *CPU) add8(a, b uint8) uint8 {
	result16 := uint16(a) + uint16(b)
	result := uint8(result16)

	c.SetCarry(result16 > 0xFF)
	c.SetSZP8(result)
	c.SetOverflow(((a ^ b ^ 0x80) & (result ^ a) & 0x80) != 0)
	c.SetAuxCarry((a&0x0F)+(b&0x0F) > 0x0F)

	return result
}

// add16 adds two 16-bit values and sets flags.
func (c *CPU) add16(a, b uint16) uint16 {
	result32 := uint32(a) + uint32(b)
	result := uint16(result32)

	c.SetCarry(result32 > 0xFFFF)
	c.SetSZP16(result)
	c.SetOverflow(((a ^ b ^ 0x8000) & (result ^ a) & 0x8000) != 0)
	c.SetAuxCarry((a&0x0F)+(b&0x0F) > 0x0F)

	return result
}

// sub8 subtracts two 8-bit values and sets flags.
func (c *CPU) sub8(a, b uint8) uint8 {
	result16 := uint16(a) - uint16(b)
	result := uint8(result16)

	c.SetCarry(a < b)
	c.SetSZP8(result)
	c.SetOverflow(((a ^ b) & (result ^ a) & 0x80) != 0)
	c.SetAuxCarry((a & 0x0F) < (b & 0x0F))

	return result
}

// sub16 subtracts two 16-bit values and sets flags.
func (c *CPU) sub16(a, b uint16) uint16 {
	result32 := uint32(a) - uint32(b)
	result := uint16(result32)

	c.SetCarry(a < b)
	c.SetSZP16(result)
	c.SetOverflow(((a ^ b) & (result ^ a) & 0x8000) != 0)
	c.SetAuxCarry((a & 0x0F) < (b & 0x0F))

	return result
}

// and8 performs bitwise AND on 8-bit values and sets flags.
func (c *CPU) and8(a, b uint8) uint8 {
	result := a & b
	c.SetCarry(false)
	c.SetOverflow(false)
	c.SetSZP8(result)
	return result
}

// and16 performs bitwise AND on 16-bit values and sets flags.
func (c *CPU) and16(a, b uint16) uint16 {
	result := a & b
	c.SetCarry(false)
	c.SetOverflow(false)
	c.SetSZP16(result)
	return result
}

// or8 performs bitwise OR on 8-bit values and sets flags.
func (c *CPU) or8(a, b uint8) uint8 {
	result := a | b
	c.SetCarry(false)
	c.SetOverflow(false)
	c.SetSZP8(result)
	return result
}

// or16 performs bitwise OR on 16-bit values and sets flags.
func (c *CPU) or16(a, b uint16) uint16 {
	result := a | b
	c.SetCarry(false)
	c.SetOverflow(false)
	c.SetSZP16(result)
	return result
}

// Register access helpers

// getReg8 gets an 8-bit register value.
func (c *CPU) getReg8(reg RegisterParam) uint8 {
	switch reg {
	case RegAL:
		return c.AL()
	case RegCL:
		return c.CL()
	case RegDL:
		return c.DL()
	case RegBL:
		return c.BL()
	case RegAH:
		return c.AH()
	case RegCH:
		return c.CH()
	case RegDH:
		return c.DH()
	case RegBH:
		return c.BH()
	default:
		return 0
	}
}

// setReg8 sets an 8-bit register value.
func (c *CPU) setReg8(reg RegisterParam, value uint8) {
	switch reg {
	case RegAL:
		c.SetAL(value)
	case RegCL:
		c.SetCL(value)
	case RegDL:
		c.SetDL(value)
	case RegBL:
		c.SetBL(value)
	case RegAH:
		c.SetAH(value)
	case RegCH:
		c.SetCH(value)
	case RegDH:
		c.SetDH(value)
	case RegBH:
		c.SetBH(value)
	}
}

// getReg16 gets a 16-bit register value.
func (c *CPU) getReg16(reg RegisterParam) uint16 {
	switch reg {
	case RegAX:
		return c.AX
	case RegCX:
		return c.CX
	case RegDX:
		return c.DX
	case RegBX:
		return c.BX
	case RegSP:
		return c.SP
	case RegBP:
		return c.BP
	case RegSI:
		return c.SI
	case RegDI:
		return c.DI
	case RegES:
		return c.ES
	case RegCS:
		return c.CS
	case RegSS:
		return c.SS
	case RegDS:
		return c.DS
	default:
		return 0
	}
}

// setReg16 sets a 16-bit register value.
func (c *CPU) setReg16(reg RegisterParam, value uint16) {
	switch reg {
	case RegAX:
		c.AX = value
	case RegCX:
		c.CX = value
	case RegDX:
		c.DX = value
	case RegBX:
		c.BX = value
	case RegSP:
		c.SP = value
	case RegBP:
		c.BP = value
	case RegSI:
		c.SI = value
	case RegDI:
		c.DI = value
	case RegES:
		c.ES = value
	case RegCS:
		c.CS = value
	case RegSS:
		c.SS = value
	case RegDS:
		c.DS = value
	}
}

// Tracing support

// prepareTraceStep prepares a trace step before instruction execution.
func (c *CPU) prepareTraceStep(_ uint32, opcode uint8) TraceStep {
	return TraceStep{
		IP:     c.IP - 1,
		CS:     c.CS,
		Opcode: opcode,

		// Capture pre-execution state
		PreAX:    c.AX,
		PreBX:    c.BX,
		PreCX:    c.CX,
		PreDX:    c.DX,
		PreSI:    c.SI,
		PreDI:    c.DI,
		PreBP:    c.BP,
		PreSP:    c.SP,
		PreCS:    c.CS,
		PreDS:    c.DS,
		PreES:    c.ES,
		PreSS:    c.SS,
		PreFlags: c.Flags,

		Cycles: c.cycles,
	}
}

// completeTraceStep completes the trace step after instruction execution.
func (c *CPU) completeTraceStep() {
	c.TraceStep.PostAX = c.AX
	c.TraceStep.PostBX = c.BX
	c.TraceStep.PostCX = c.CX
	c.TraceStep.PostDX = c.DX
	c.TraceStep.PostSI = c.SI
	c.TraceStep.PostDI = c.DI
	c.TraceStep.PostBP = c.BP
	c.TraceStep.PostSP = c.SP
	c.TraceStep.PostCS = c.CS
	c.TraceStep.PostDS = c.DS
	c.TraceStep.PostES = c.ES
	c.TraceStep.PostSS = c.SS
	c.TraceStep.PostFlags = c.Flags
}

// Instruction emulation functions

// Data Movement Instructions

// movRMReg8 implements MOV r/m8, r8.
func movRMReg8(c *CPU, params ...any) error {
	modrm := params[0].(ModRM)
	displacement := params[1].(int16)
	srcValue := c.getReg8(RegisterParam(modrm.Reg))
	if modrm.Mod == 3 {
		c.setReg8(RegisterParam(modrm.RM), srcValue)
	} else {
		addr := c.GetEffectiveAddress(modrm, displacement, 0)
		c.memory.Write8(addr, srcValue)
	}
	return nil
}

// movRMReg16 implements MOV r/m16, r16.
func movRMReg16(c *CPU, params ...any) error {
	modrm := params[0].(ModRM)
	displacement := params[1].(int16)
	srcValue := c.getReg16(RegisterParam(modrm.Reg))
	if modrm.Mod == 3 {
		c.setReg16(RegisterParam(modrm.RM), srcValue)
	} else {
		addr := c.GetEffectiveAddress(modrm, displacement, 0)
		c.memory.Write16(addr, srcValue)
	}
	return nil
}

// movRegRM8 implements MOV r8, r/m8.
func movRegRM8(c *CPU, params ...any) error {
	modrm := params[0].(ModRM)
	displacement := params[1].(int16)
	var srcValue uint8
	if modrm.Mod == 3 {
		srcValue = c.getReg8(RegisterParam(modrm.RM))
	} else {
		addr := c.GetEffectiveAddress(modrm, displacement, 0)
		srcValue = c.memory.Read8(addr)
	}
	c.setReg8(RegisterParam(modrm.Reg), srcValue)
	return nil
}

// movRegRM16 implements MOV r16, r/m16.
func movRegRM16(c *CPU, params ...any) error {
	modrm := params[0].(ModRM)
	displacement := params[1].(int16)
	var srcValue uint16
	if modrm.Mod == 3 {
		srcValue = c.getReg16(RegisterParam(modrm.RM))
	} else {
		addr := c.GetEffectiveAddress(modrm, displacement, 0)
		srcValue = c.memory.Read16(addr)
	}
	c.setReg16(RegisterParam(modrm.Reg), srcValue)
	return nil
}

// movRegImm8 implements MOV r8, imm8.
func movRegImm8(c *CPU, params ...any) error {
	_ = params[0].(uint8)
	return nil
}

// movRegImm16 implements MOV r16, imm16.
func movRegImm16(c *CPU, params ...any) error {
	_ = params[0].(uint16)
	return nil
}

// movMemImm8 implements MOV r/m8, imm8.
func movMemImm8(c *CPU, params ...any) error {
	modrm := params[0].(ModRM)
	displacement := params[1].(int16)
	immediate := params[2].(uint8)
	if modrm.Mod == 3 {
		c.setReg8(RegisterParam(modrm.RM), immediate)
	} else {
		addr := c.GetEffectiveAddress(modrm, displacement, 0)
		c.memory.Write8(addr, immediate)
	}
	return nil
}

// movMemImm16 implements MOV r/m16, imm16.
func movMemImm16(c *CPU, params ...any) error {
	modrm := params[0].(ModRM)
	displacement := params[1].(int16)
	immediate := params[2].(uint16)
	if modrm.Mod == 3 {
		c.setReg16(RegisterParam(modrm.RM), immediate)
	} else {
		addr := c.GetEffectiveAddress(modrm, displacement, 0)
		c.memory.Write16(addr, immediate)
	}
	return nil
}

// Arithmetic Instructions

// addALImm8 implements ADD AL, imm8.
func addALImm8(c *CPU, params ...any) error {
	immediate := params[0].(uint8)
	result := c.add8(c.AL(), immediate)
	c.SetAL(result)
	return nil
}

// addAXImm16 implements ADD AX, imm16.
func addAXImm16(c *CPU, params ...any) error {
	immediate := params[0].(uint16)
	result := c.add16(c.AX, immediate)
	c.AX = result
	return nil
}

// subRMReg8 implements SUB r/m8, r8.
func subRMReg8(c *CPU, params ...any) error {
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
}

// subRMReg16 implements SUB r/m16, r16.
func subRMReg16(c *CPU, params ...any) error {
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
}

// cmpRMReg8 implements CMP r/m8, r8.
func cmpRMReg8(c *CPU, params ...any) error {
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
}

// cmpRMReg16 implements CMP r/m16, r16.
func cmpRMReg16(c *CPU, params ...any) error {
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
}

// adcALImm8 implements ADC AL, imm8.
func adcALImm8(c *CPU, params ...any) error {
	immediate := params[0].(uint8)
	carry := uint8(0)
	if c.Flags.GetCarry() {
		carry = 1
	}
	result := c.add8(c.AL(), immediate+carry)
	c.SetAL(result)
	return nil
}
