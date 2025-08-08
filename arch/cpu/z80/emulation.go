package z80

// inc8 increments an 8-bit value and sets flags appropriately.
func (c *CPU) inc8(value uint8) uint8 {
	result := value + 1

	// Set flags efficiently
	c.setSZ(result)
	c.setH((value & 0x0F) == 0x0F) // Half carry if lower nibble was 0xF
	c.setPOverflow(value == 0x7F)  // Overflow if incrementing 0x7F
	c.Flags.N = 0                  // Clear N flag for increment (direct assignment)

	return result
}

// dec8 decrements an 8-bit value and sets flags appropriately.
func (c *CPU) dec8(value uint8) uint8 {
	result := value - 1

	// Set flags efficiently
	c.setSZ(result)
	c.setH((value & 0x0F) == 0x00) // Half carry if lower nibble was 0x0
	c.setPOverflow(value == 0x80)  // Overflow if decrementing 0x80
	c.Flags.N = 1                  // Set N flag for decrement (direct assignment)

	return result
}

// add8 adds two 8-bit values and sets flags appropriately.
func (c *CPU) add8(a, b uint8) uint8 {
	result16 := uint16(a) + uint16(b)
	result := uint8(result16)

	// Set flags
	c.setSZ(result)
	c.setC(result16 > 0xFF)                                     // Carry if result > 255
	c.setH((a&0x0F)+(b&0x0F) > 0x0F)                            // Half carry
	c.setPOverflow(((a ^ b ^ 0x80) & (result ^ a) & 0x80) != 0) // Overflow
	c.setN(false)                                               // Clear N flag for addition

	return result
}

// sub8 subtracts two 8-bit values and sets flags appropriately.
func (c *CPU) sub8(a, b uint8) uint8 {
	result16 := uint16(a) - uint16(b)
	result := uint8(result16)

	// Set flags
	c.setSZ(result)
	c.setC(a < b)                                        // Carry if a < b
	c.setH((a & 0x0F) < (b & 0x0F))                      // Half carry
	c.setPOverflow(((a ^ b) & (a ^ result) & 0x80) != 0) // Overflow
	c.setN(true)                                         // Set N flag for subtraction

	return result
}

// add16 adds two 16-bit values and sets carry/half-carry flags.
func (c *CPU) add16(a, b uint16) uint16 {
	result32 := uint32(a) + uint32(b)
	result := uint16(result32)

	// Set flags (only C, H, and N for 16-bit operations)
	c.setC(result32 > 0xFFFF)              // Carry if result > 65535
	c.setH((a&0x0FFF)+(b&0x0FFF) > 0x0FFF) // Half carry on bit 11
	c.setN(false)                          // Clear N flag for addition

	return result
}

// and8 performs bitwise AND on two 8-bit values and sets flags.
func (c *CPU) and8(a, b uint8) uint8 {
	result := a & b

	// Set flags
	c.setSZP(result)
	c.setH(true)  // Half carry is always set for AND
	c.setN(false) // Clear N flag
	c.setC(false) // Clear carry flag

	return result
}

// or8 performs bitwise OR on two 8-bit values and sets flags.
func (c *CPU) or8(a, b uint8) uint8 {
	result := a | b

	// Set flags
	c.setSZP(result)
	c.setH(false) // Clear half carry
	c.setN(false) // Clear N flag
	c.setC(false) // Clear carry flag

	return result
}

// xor8 performs bitwise XOR on two 8-bit values and sets flags.
func (c *CPU) xor8(a, b uint8) uint8 {
	result := a ^ b

	// Set flags
	c.setSZP(result)
	c.setH(false) // Clear half carry
	c.setN(false) // Clear N flag
	c.setC(false) // Clear carry flag

	return result
}

// cp compares two 8-bit values (like SUB but doesn't store result).
func (c *CPU) cp(a, b uint8) {
	result16 := uint16(a) - uint16(b)
	result := uint8(result16)

	// Set flags
	c.setSZ(result)
	c.setC(a < b)                                        // Carry if a < b
	c.setH((a & 0x0F) < (b & 0x0F))                      // Half carry
	c.setPOverflow(((a ^ b) & (a ^ result) & 0x80) != 0) // Overflow
	c.setN(true)                                         // Set N flag for subtraction
}

// rlca rotates accumulator left circular and sets carry.
func (c *CPU) rlca(value uint8) uint8 {
	carry := (value & 0x80) >> 7
	result := (value << 1) | carry

	c.setC(carry != 0)
	c.setH(false)
	c.setN(false)

	return result
}

// rrca rotates accumulator right circular and sets carry.
func (c *CPU) rrca(value uint8) uint8 {
	carry := value & 0x01
	result := (value >> 1) | (carry << 7)

	c.setC(carry != 0)
	c.setH(false)
	c.setN(false)

	return result
}

// rlc rotates value left circular and sets all flags.
func (c *CPU) rlc(value uint8) uint8 {
	carry := (value & 0x80) >> 7
	result := (value << 1) | carry

	c.setSZP(result)
	c.setC(carry != 0)
	c.setH(false)
	c.setN(false)

	return result
}

// rl rotates value left through carry and sets all flags.
func (c *CPU) rl(value uint8) uint8 {
	newCarry := (value & 0x80) >> 7
	result := (value << 1) | c.Flags.C

	c.setSZP(result)
	c.setC(newCarry != 0)
	c.setH(false)
	c.setN(false)

	return result
}

// sla shifts value left arithmetic and sets all flags.
func (c *CPU) sla(value uint8) uint8 {
	carry := (value & 0x80) >> 7
	result := value << 1

	c.setSZP(result)
	c.setC(carry != 0)
	c.setH(false)
	c.setN(false)

	return result
}

// sra shifts value right arithmetic and sets all flags.
func (c *CPU) sra(value uint8) uint8 {
	carry := value & 0x01
	result := (value >> 1) | (value & 0x80) // Keep sign bit

	c.setSZP(result)
	c.setC(carry != 0)
	c.setH(false)
	c.setN(false)

	return result
}

// srl shifts value right logical and sets all flags.
func (c *CPU) srl(value uint8) uint8 {
	carry := value & 0x01
	result := value >> 1

	c.setSZP(result)
	c.setC(carry != 0)
	c.setH(false)
	c.setN(false)

	return result
}

// bit tests bit n of value and sets flags.
func (c *CPU) bit(n uint8, value uint8) {
	bit := (value >> n) & 1

	setFlag(&c.Flags.Z, bit == 0)
	c.setH(true)  // Half carry is always set for BIT
	c.setN(false) // Clear N flag
	// S and P flags are affected differently for BIT instruction
	if n == 7 {
		setFlag(&c.Flags.S, bit != 0)
	}
}

// set sets bit n of value.
func (c *CPU) setBit(n uint8, value uint8) uint8 {
	return value | (1 << n)
}

// res resets bit n of value.
func (c *CPU) res(n uint8, value uint8) uint8 {
	return value & ^(1 << n)
}

// neg negates the accumulator (two's complement).
func (c *CPU) neg(value uint8) uint8 {
	result := uint8(-int8(value))

	c.setSZP(result)
	c.setC(value != 0)            // Carry set unless original value was 0
	c.setH((value & 0x0F) != 0)   // Half carry
	c.setPOverflow(value == 0x80) // Overflow if negating 0x80
	c.setN(true)                  // Set N flag for negation

	return result
}

// adc adds with carry.
func (c *CPU) adc(a, b uint8) uint8 {
	carry := c.Flags.C
	result16 := uint16(a) + uint16(b) + uint16(carry)
	result := uint8(result16)

	// Set flags
	c.setSZ(result)
	c.setC(result16 > 0xFF)
	c.setH((a&0x0F)+(b&0x0F)+carry > 0x0F)
	c.setPOverflow(((a ^ b ^ 0x80) & (result ^ a) & 0x80) != 0)
	c.setN(false)

	return result
}

// sbc subtracts with carry.
func (c *CPU) sbc(a, b uint8) uint8 {
	carry := c.Flags.C
	result16 := uint16(a) - uint16(b) - uint16(carry)
	result := uint8(result16)

	// Set flags
	c.setSZ(result)
	c.setC(result16 > 0xFF) // Borrow occurred
	c.setH((a & 0x0F) < (b&0x0F)+carry)
	c.setPOverflow(((a ^ b) & (a ^ result) & 0x80) != 0)
	c.setN(true)

	return result
}

// nop performs no operation.
func nop(c *CPU) error {
	return nil
}

// halt halts the CPU execution.
func halt(c *CPU) error {
	c.halted = true
	return nil
}

// ldImm8 loads an 8-bit immediate value into a register.
func ldImm8(c *CPU, params ...any) error {
	if len(params) < 1 {
		return ErrMissingParameter
	}

	imm, ok := params[0].(Immediate8)
	if !ok {
		return ErrInvalidParameterType
	}

	// For now, load into A register (this would need opcode analysis for other registers)
	c.A = uint8(imm)
	return nil
}

// ldReg8 loads between 8-bit registers.
func ldReg8(c *CPU, params ...any) error {
	if len(params) < 2 {
		return ErrMissingParameter
	}

	dst, ok1 := params[0].(Register8)
	src, ok2 := params[1].(Register8)

	if !ok1 || !ok2 {
		return ErrInvalidParameterType
	}

	value := c.GetRegisterValue(uint8(src))
	c.SetRegisterValue(uint8(dst), value)
	return nil
}

// incReg8 increments an 8-bit register.
func incReg8(c *CPU, params ...any) error {
	if len(params) < 1 {
		return ErrMissingParameter
	}

	reg, ok := params[0].(Register8)
	if !ok {
		return ErrInvalidParameterType
	}

	value := c.GetRegisterValue(uint8(reg))
	result := c.inc8(value)
	c.SetRegisterValue(uint8(reg), result)
	return nil
}

// decReg8 decrements an 8-bit register.
func decReg8(c *CPU, params ...any) error {
	if len(params) < 1 {
		return ErrMissingParameter
	}

	reg, ok := params[0].(Register8)
	if !ok {
		return ErrInvalidParameterType
	}

	value := c.GetRegisterValue(uint8(reg))
	result := c.dec8(value)
	c.SetRegisterValue(uint8(reg), result)
	return nil
}

// addA adds a value to the accumulator.
func addA(c *CPU, params ...any) error {
	if len(params) < 1 {
		return ErrMissingParameter
	}

	var value uint8

	switch param := params[0].(type) {
	case Register8:
		value = c.GetRegisterValue(uint8(param))
	case Immediate8:
		value = uint8(param)
	default:
		return ErrInvalidParameterType
	}

	c.A = c.add8(c.A, value)
	return nil
}

// subA subtracts a value from the accumulator.
func subA(c *CPU, params ...any) error {
	if len(params) < 1 {
		return ErrMissingParameter
	}

	var value uint8

	switch param := params[0].(type) {
	case Register8:
		value = c.GetRegisterValue(uint8(param))
	case Immediate8:
		value = uint8(param)
	default:
		return ErrInvalidParameterType
	}

	c.A = c.sub8(c.A, value)
	return nil
}

// andA performs logical AND with the accumulator.
func andA(c *CPU, params ...any) error {
	if len(params) < 1 {
		return ErrMissingParameter
	}

	var value uint8

	switch param := params[0].(type) {
	case Register8:
		value = c.GetRegisterValue(uint8(param))
	case Immediate8:
		value = uint8(param)
	default:
		return ErrInvalidParameterType
	}

	c.A = c.and8(c.A, value)
	return nil
}

// orA performs logical OR with the accumulator.
func orA(c *CPU, params ...any) error {
	if len(params) < 1 {
		return ErrMissingParameter
	}

	var value uint8

	switch param := params[0].(type) {
	case Register8:
		value = c.GetRegisterValue(uint8(param))
	case Immediate8:
		value = uint8(param)
	default:
		return ErrInvalidParameterType
	}

	c.A = c.or8(c.A, value)
	return nil
}

// xorA performs logical XOR with the accumulator.
func xorA(c *CPU, params ...any) error {
	if len(params) < 1 {
		return ErrMissingParameter
	}

	var value uint8

	switch param := params[0].(type) {
	case Register8:
		value = c.GetRegisterValue(uint8(param))
	case Immediate8:
		value = uint8(param)
	default:
		return ErrInvalidParameterType
	}

	c.A = c.xor8(c.A, value)
	return nil
}

// cpA compares a value with the accumulator.
func cpA(c *CPU, params ...any) error {
	if len(params) < 1 {
		return ErrMissingParameter
	}

	var value uint8

	switch param := params[0].(type) {
	case Register8:
		value = c.GetRegisterValue(uint8(param))
	case Immediate8:
		value = uint8(param)
	default:
		return ErrInvalidParameterType
	}

	c.cp(c.A, value)
	return nil
}

// jpAbs performs absolute jump.
func jpAbs(c *CPU, params ...any) error {
	if len(params) < 1 {
		return ErrMissingParameter
	}

	addr, ok := params[0].(Extended)
	if !ok {
		return ErrInvalidParameterType
	}

	c.PC = uint16(addr)
	return nil
}

// jrRel performs relative jump.
func jrRel(c *CPU, params ...any) error {
	if len(params) < 1 {
		return ErrMissingParameter
	}

	offset, ok := params[0].(Relative)
	if !ok {
		return ErrInvalidParameterType
	}

	c.PC = uint16(int32(c.PC) + int32(offset))
	return nil
}

// Additional Z80 emulation functions

// ldReg16 loads a 16-bit immediate value into a register pair.
func ldReg16(c *CPU, params ...any) error {
	// Implementation placeholder - would need opcode analysis for specific register pair
	return nil
}

// ldIndirect loads between register pairs and memory locations.
func ldIndirect(c *CPU, params ...any) error {
	// Implementation placeholder - would need specific addressing mode handling
	return nil
}

// incReg16 increments a 16-bit register pair.
func incReg16(c *CPU, params ...any) error {
	// Implementation placeholder - would need opcode analysis for specific register pair
	return nil
}

// decReg16 decrements a 16-bit register pair.
func decReg16(c *CPU, params ...any) error {
	// Implementation placeholder - would need opcode analysis for specific register pair
	return nil
}

// rlca performs rotate left circular accumulator.
func rlca(c *CPU) error {
	c.A = c.rlca(c.A)
	return nil
}

// rrca performs rotate right circular accumulator.
func rrca(c *CPU) error {
	c.A = c.rrca(c.A)
	return nil
}

// rla performs rotate left accumulator through carry.
func rla(c *CPU) error {
	c.A = c.rl(c.A)
	return nil
}

// rra performs rotate right accumulator through carry.
func rra(c *CPU) error {
	newCarry := c.A & 0x01
	c.A = (c.A >> 1) | (c.Flags.C << 7)
	c.setC(newCarry != 0)
	c.setH(false)
	c.setN(false)
	return nil
}

// exAf exchanges AF with AF'.
func exAf(c *CPU) error {
	// Exchange AF with shadow AF'
	tempA := c.A
	tempF := c.Flags
	c.A = c.A_
	c.Flags = c.Flags_
	c.A_ = tempA
	c.Flags_ = tempF
	return nil
}

// addHl adds a 16-bit register pair to HL.
func addHl(c *CPU, params ...any) error {
	// Implementation placeholder - would need specific register pair identification
	return nil
}

// djnz decrements B and jumps if not zero.
func djnz(c *CPU, params ...any) error {
	c.B--
	if c.B != 0 {
		if len(params) < 1 {
			return ErrMissingParameter
		}
		offset, ok := params[0].(Relative)
		if !ok {
			return ErrInvalidParameterType
		}
		c.PC = uint16(int32(c.PC) + int32(offset))
	}
	return nil
}

// jrCond performs conditional relative jump.
func jrCond(c *CPU, params ...any) error {
	// Implementation placeholder - would need condition checking
	return nil
}

// ldExtended loads using extended addressing.
func ldExtended(c *CPU, params ...any) error {
	// Implementation placeholder - would need specific addressing mode handling
	return nil
}

// daa performs decimal adjust accumulator.
func daa(c *CPU) error {
	// Implementation placeholder - complex BCD adjustment logic
	return nil
}

// cpl complements the accumulator.
func cpl(c *CPU) error {
	c.A = ^c.A
	c.setH(true)
	c.setN(true)
	return nil
}

// incIndirect increments memory location pointed to by register pair.
func incIndirect(c *CPU, params ...any) error {
	// Implementation placeholder - would need specific addressing
	return nil
}

// decIndirect decrements memory location pointed to by register pair.
func decIndirect(c *CPU, params ...any) error {
	// Implementation placeholder - would need specific addressing
	return nil
}

// ldIndirectImm loads immediate value to indirect memory location.
func ldIndirectImm(c *CPU, params ...any) error {
	// Implementation placeholder - would need specific addressing
	return nil
}

// scf sets the carry flag.
func scf(c *CPU) error {
	c.setC(true)
	c.setH(false)
	c.setN(false)
	return nil
}

// ccf complements the carry flag.
func ccf(c *CPU) error {
	c.setC(c.Flags.C == 0) // Complement carry flag
	c.setH(false)
	c.setN(false)
	return nil
}

// adcA adds with carry to accumulator.
func adcA(c *CPU, params ...any) error {
	if len(params) < 1 {
		return ErrMissingParameter
	}

	var value uint8
	switch param := params[0].(type) {
	case Register8:
		value = c.GetRegisterValue(uint8(param))
	case Immediate8:
		value = uint8(param)
	default:
		return ErrInvalidParameterType
	}

	c.A = c.adc(c.A, value)
	return nil
}

// sbcA subtracts with carry from accumulator.
func sbcA(c *CPU, params ...any) error {
	if len(params) < 1 {
		return ErrMissingParameter
	}

	var value uint8
	switch param := params[0].(type) {
	case Register8:
		value = c.GetRegisterValue(uint8(param))
	case Immediate8:
		value = uint8(param)
	default:
		return ErrInvalidParameterType
	}

	c.A = c.sbc(c.A, value)
	return nil
}

// retCond performs conditional return.
func retCond(c *CPU) error {
	// Implementation placeholder - would need condition checking
	return nil
}

// popReg16 pops 16-bit register from stack.
func popReg16(c *CPU, params ...any) error {
	// Implementation placeholder - would need stack operations
	return nil
}

// jpCond performs conditional jump.
func jpCond(c *CPU, params ...any) error {
	// Implementation placeholder - would need condition checking
	return nil
}

// callCond performs conditional call.
func callCond(c *CPU, params ...any) error {
	// Implementation placeholder - would need condition checking and stack operations
	return nil
}

// pushReg16 pushes 16-bit register to stack.
func pushReg16(c *CPU, params ...any) error {
	// Implementation placeholder - would need stack operations
	return nil
}

// rst performs restart (call to fixed address).
func rst(c *CPU, params ...any) error {
	// Implementation placeholder - would need stack operations and address calculation
	return nil
}

// ret returns from subroutine.
func ret(c *CPU) error {
	// Implementation placeholder - would need stack operations
	return nil
}

// call calls a subroutine.
func call(c *CPU, params ...any) error {
	// Implementation placeholder - would need stack operations
	return nil
}

// outPort outputs to port.
func outPort(c *CPU, params ...any) error {
	// Implementation placeholder - would need port I/O
	return nil
}

// inPort inputs from port.
func inPort(c *CPU, params ...any) error {
	// Implementation placeholder - would need port I/O
	return nil
}

// exx exchanges register pairs.
func exx(c *CPU) error {
	// Exchange BC, DE, HL with BC', DE', HL'
	tempB := c.B
	tempC := c.C
	tempD := c.D
	tempE := c.E
	tempH := c.H
	tempL := c.L

	c.B = c.B_
	c.C = c.C_
	c.D = c.D_
	c.E = c.E_
	c.H = c.H_
	c.L = c.L_

	c.B_ = tempB
	c.C_ = tempC
	c.D_ = tempD
	c.E_ = tempE
	c.H_ = tempH
	c.L_ = tempL
	return nil
}

// exSp exchanges top of stack with register pair.
func exSp(c *CPU, params ...any) error {
	// Implementation placeholder - would need stack operations
	return nil
}

// jpIndirect performs indirect jump.
func jpIndirect(c *CPU, params ...any) error {
	// Implementation placeholder - would need indirect addressing
	return nil
}

// exDeHl exchanges DE with HL.
func exDeHl(c *CPU) error {
	tempD := c.D
	tempE := c.E
	c.D = c.H
	c.E = c.L
	c.H = tempD
	c.L = tempE
	return nil
}

// di disables interrupts.
func di(c *CPU) error {
	c.iff1 = false
	c.iff2 = false
	return nil
}

// ei enables interrupts.
func ei(c *CPU) error {
	c.iff1 = true
	c.iff2 = true
	return nil
}

// ldSp loads SP from HL.
func ldSp(c *CPU, params ...any) error {
	c.SP = uint16(c.H)<<8 | uint16(c.L)
	return nil
}
