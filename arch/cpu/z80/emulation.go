package z80

import "fmt"

// inc8 increments an 8-bit value and sets flags appropriately.
func (c *CPU) inc8(value uint8) uint8 {
	result := value + 1

	c.setSZ(result)
	c.setH((value & 0x0F) == 0x0F) // Half carry if lower nibble was 0xF
	c.setPOverflow(value == 0x7F)  // Overflow if incrementing 0x7F
	c.setN(false)                  // Indicates arithmetic (not logical) operation

	return result
}

// dec8 decrements an 8-bit value and sets flags appropriately.
func (c *CPU) dec8(value uint8) uint8 {
	result := value - 1

	c.setSZ(result)
	c.setH((value & 0x0F) == 0x00) // Half carry if lower nibble was 0x0
	c.setPOverflow(value == 0x80)  // Overflow if decrementing 0x80
	c.setN(true)                   // Indicates subtraction operation for BCD correction

	return result
}

// add8 adds two 8-bit values and sets flags appropriately.
func (c *CPU) add8(a, b uint8) uint8 {
	result16 := uint16(a) + uint16(b)
	result := uint8(result16)

	c.setSZ(result)
	c.setC(result16 > 0xFF)                                     // Carry if result > 255
	c.setH((a&0x0F)+(b&0x0F) > 0x0F)                            // Half carry from bit 3 to 4
	c.setPOverflow(((a ^ b ^ 0x80) & (result ^ a) & 0x80) != 0) // Two's complement overflow detection
	c.setN(false)                                               // Indicates addition for BCD correction

	return result
}

// sub8 subtracts two 8-bit values and sets flags appropriately.
func (c *CPU) sub8(a, b uint8) uint8 {
	result16 := uint16(a) - uint16(b)
	result := uint8(result16)

	c.setSZ(result)
	c.setC(a < b)                                        // Borrow if minuend < subtrahend
	c.setH((a & 0x0F) < (b & 0x0F))                      // Half borrow from bit 3
	c.setPOverflow(((a ^ b) & (a ^ result) & 0x80) != 0) // Two's complement overflow detection
	c.setN(true)                                         // Indicates subtraction for BCD correction

	return result
}

// add16 adds two 16-bit values and sets carry/half-carry flags.
func (c *CPU) add16(a, b uint16) uint16 {
	result32 := uint32(a) + uint32(b)
	result := uint16(result32)

	// Update limited flags for 16-bit arithmetic (Z80 16-bit ops don't affect S,Z,P)
	c.setC(result32 > 0xFFFF)              // Carry if result > 65535
	c.setH((a&0x0FFF)+(b&0x0FFF) > 0x0FFF) // Half carry from bit 11 to bit 12
	c.setN(false)                          // Indicates addition for BCD correction

	return result
}

// and8 performs bitwise AND on two 8-bit values and sets flags.
func (c *CPU) and8(a, b uint8) uint8 {
	result := a & b

	c.setSZP(result)
	c.setH(true)  // H always set for Z80 AND instruction
	c.setN(false) // Indicates logical (not arithmetic) operation
	c.setC(false) // Logical operations clear carry

	return result
}

// or8 performs bitwise OR on two 8-bit values and sets flags.
func (c *CPU) or8(a, b uint8) uint8 {
	result := a | b

	c.setSZP(result)
	c.setH(false) // Logical operations clear half carry
	c.setN(false) // Indicates logical (not arithmetic) operation
	c.setC(false) // Logical operations clear carry

	return result
}

// xor8 performs bitwise XOR on two 8-bit values and sets flags.
func (c *CPU) xor8(a, b uint8) uint8 {
	result := a ^ b

	c.setSZP(result)
	c.setH(false) // Logical operations clear half carry
	c.setN(false) // Indicates logical (not arithmetic) operation
	c.setC(false) // Logical operations clear carry

	return result
}

// cp compares two 8-bit values (like SUB but doesn't store result).
func (c *CPU) cp(a, b uint8) {
	result16 := uint16(a) - uint16(b)
	result := uint8(result16)

	// Set flags
	c.setS(result)                                       // S from result
	c.setZ(result)                                       // Z from result
	c.setXY(b)                                           // X/Y from operand (not result) - Z80 quirk for CP
	c.setC(a < b)                                        // Carry if a < b
	c.setH((a & 0x0F) < (b & 0x0F))                      // Half carry
	c.setPOverflow(((a ^ b) & (a ^ result) & 0x80) != 0) // Overflow
	c.setN(true)                                         // Indicates subtraction for BCD correction
}

// neg negates the accumulator (two's complement).
func (c *CPU) neg(value uint8) uint8 {
	result := uint8(-int8(value))

	c.setSZP(result)
	c.setC(value != 0)            // Carry set unless original value was 0
	c.setH((value & 0x0F) != 0)   // Half carry
	c.setPOverflow(value == 0x80) // Overflow if negating 0x80
	c.setN(true)                  // Indicates subtraction-based operation (two's complement)

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

// Rotate operations

// rlca rotates accumulator left circular and sets carry.
func (c *CPU) rlca(value uint8) uint8 {
	carry := (value & 0x80) >> 7
	result := (value << 1) | carry

	c.setC(carry != 0)
	c.setH(false)
	c.setN(false)
	c.setXY(result)

	return result
}

// rrca rotates accumulator right circular and sets carry.
func (c *CPU) rrca(value uint8) uint8 {
	carry := value & 0x01
	result := (value >> 1) | (carry << 7)

	c.setC(carry != 0)
	c.setH(false)
	c.setN(false)
	c.setXY(result)

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

// rrc rotates value right circular and sets all flags.
func (c *CPU) rrc(value uint8) uint8 {
	carry := value & 0x01
	result := (value >> 1) | (carry << 7)

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

// rr rotates value right through carry and sets all flags.
func (c *CPU) rr(value uint8) uint8 {
	newCarry := value & 0x01
	result := (value >> 1) | (c.Flags.C << 7)

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

// sll shifts value left logical (undocumented) and sets all flags.
func (c *CPU) sll(value uint8) uint8 {
	carry := (value & 0x80) >> 7
	result := (value << 1) | 0x01 // Set bit 0

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

// Bit operations

// bit tests bit n of value and sets flags.
// For BIT on register operands, X/Y come from the register value.
// For BIT on (HL), X/Y come from MEMPTR high byte - caller must handle this.
func (c *CPU) bit(n uint8, value uint8) {
	bit := (value >> n) & 1
	bitIsZero := bit == 0

	setFlag(&c.Flags.Z, bitIsZero)
	setFlag(&c.Flags.P, bitIsZero) // P/V same as Z for BIT instruction
	setFlag(&c.Flags.S, n == 7 && bit != 0)
	c.setH(true)
	c.setN(false)
	c.setXY(value) // X/Y from value for register BIT ops
}

// bitMemptr tests bit n of value, setting X/Y from MEMPTR high byte.
// Used for BIT n,(HL) and BIT n,(IX+d)/(IY+d).
func (c *CPU) bitMemptr(n uint8, value uint8, memptrHigh uint8) {
	bit := (value >> n) & 1
	bitIsZero := bit == 0

	setFlag(&c.Flags.Z, bitIsZero)
	setFlag(&c.Flags.P, bitIsZero)
	setFlag(&c.Flags.S, n == 7 && bit != 0)
	c.setH(true)
	c.setN(false)
	c.setXY(memptrHigh)
}

// setBit sets bit n of value.
func (c *CPU) setBit(n uint8, value uint8) uint8 {
	return value | (1 << n)
}

// res resets bit n of value.
func (c *CPU) res(n uint8, value uint8) uint8 {
	return value & ^(1 << n)
}

// Control and flag operations

// nop performs no operation.
func nop(_ *CPU) error {
	return nil
}

// halt halts the CPU execution.
func halt(c *CPU) error {
	c.halted = true
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

// scf sets the carry flag.
func scf(c *CPU) error {
	c.setXY(c.A | (c.GetFlags() & ^c.q))
	c.setC(true)
	c.setH(false)
	c.setN(false)
	return nil
}

// ccf complements the carry flag.
func ccf(c *CPU) error {
	c.setXY(c.A | (c.GetFlags() & ^c.q))
	oldCarry := c.Flags.C != 0
	c.setC(c.Flags.C == 0)
	c.setH(oldCarry) // H = old carry
	c.setN(false)
	return nil
}

// daa performs decimal adjust accumulator (MAME-verified algorithm).
func daa(c *CPU) error {
	originalA := c.A
	correction := uint8(0)
	carry := c.Flags.C != 0

	if c.Flags.H != 0 || (c.A&0x0F) > 9 {
		correction |= 0x06
	}
	if carry || c.A > 0x99 {
		correction |= 0x60
		carry = true
	}

	if c.Flags.N != 0 {
		c.A -= correction
	} else {
		c.A += correction
	}

	// Set flags
	c.setS(c.A)
	c.setZ(c.A)
	c.setP(c.A)
	c.setXY(c.A)
	c.setC(carry)
	c.setH((originalA^c.A)&0x10 != 0)
	// N flag remains unchanged

	return nil
}

// cpl complements the accumulator.
func cpl(c *CPU) error {
	c.A = ^c.A
	c.setH(true)
	c.setN(true)
	c.setXY(c.A)
	return nil
}

// resolveOperand extracts an 8-bit value from a Register or Immediate8 parameter.
func resolveOperand(c *CPU, params ...any) (uint8, error) {
	if len(params) < 1 {
		return 0, ErrMissingParameter
	}
	switch param := params[0].(type) {
	case Register:
		return c.GetRegisterValue(uint8(param)), nil
	case Immediate8:
		return uint8(param), nil
	default:
		return 0, ErrInvalidParameterType
	}
}

// Accumulator wrapper functions

// addA adds a value to the accumulator.
func addA(c *CPU, params ...any) error {
	value, err := resolveOperand(c, params...)
	if err != nil {
		return err
	}
	c.A = c.add8(c.A, value)
	return nil
}

// subA subtracts a value from the accumulator.
func subA(c *CPU, params ...any) error {
	value, err := resolveOperand(c, params...)
	if err != nil {
		return err
	}
	c.A = c.sub8(c.A, value)
	return nil
}

// andA performs logical AND with the accumulator.
func andA(c *CPU, params ...any) error {
	value, err := resolveOperand(c, params...)
	if err != nil {
		return err
	}
	c.A = c.and8(c.A, value)
	return nil
}

// orA performs logical OR with the accumulator.
func orA(c *CPU, params ...any) error {
	value, err := resolveOperand(c, params...)
	if err != nil {
		return err
	}
	c.A = c.or8(c.A, value)
	return nil
}

// xorA performs logical XOR with the accumulator.
func xorA(c *CPU, params ...any) error {
	value, err := resolveOperand(c, params...)
	if err != nil {
		return err
	}
	c.A = c.xor8(c.A, value)
	return nil
}

// cpA compares a value with the accumulator.
func cpA(c *CPU, params ...any) error {
	value, err := resolveOperand(c, params...)
	if err != nil {
		return err
	}
	c.cp(c.A, value)
	return nil
}

// adcA adds with carry to accumulator.
func adcA(c *CPU, params ...any) error {
	value, err := resolveOperand(c, params...)
	if err != nil {
		return err
	}
	c.A = c.adc(c.A, value)
	return nil
}

// sbcA subtracts with carry from accumulator.
func sbcA(c *CPU, params ...any) error {
	value, err := resolveOperand(c, params...)
	if err != nil {
		return err
	}
	c.A = c.sbc(c.A, value)
	return nil
}

// Register 8-bit operations

// incReg8 increments an 8-bit register.
func incReg8(c *CPU, params ...any) error {
	if len(params) < 2 {
		return ErrMissingParameter
	}

	// For INC operations, the destination register (params[1]) is what gets incremented
	dstReg, ok := params[1].(Register)
	if !ok {
		return ErrInvalidParameterType
	}

	value := c.GetRegisterValue(uint8(dstReg))
	result := c.inc8(value)
	c.SetRegisterValue(uint8(dstReg), result)
	return nil
}

// decReg8 decrements an 8-bit register.
func decReg8(c *CPU, params ...any) error {
	if len(params) < 2 {
		return ErrMissingParameter
	}

	// For DEC operations, the destination register (params[1]) is what gets decremented
	dstReg, ok := params[1].(Register)
	if !ok {
		return ErrInvalidParameterType
	}

	value := c.GetRegisterValue(uint8(dstReg))
	result := c.dec8(value)
	c.SetRegisterValue(uint8(dstReg), result)
	return nil
}

// incIndirect increments memory location pointed to by register pair.
func incIndirect(c *CPU, _ ...any) error {
	address := c.hl()
	value := c.bus.Read(address)
	result := c.inc8(value)
	c.bus.Write(address, result)
	return nil
}

// decIndirect decrements memory location pointed to by register pair.
func decIndirect(c *CPU, _ ...any) error {
	address := c.hl()
	value := c.bus.Read(address)
	result := c.dec8(value)
	c.bus.Write(address, result)
	return nil
}

// Register 16-bit operations

// incReg16 increments a 16-bit register pair.
func incReg16(c *CPU, _ ...any) error {
	opcode := c.currentOpcode
	switch opcode {
	case 0x03: // INC BC
		c.setBC(c.bc() + 1)
	case 0x13: // INC DE
		c.setDE(c.de() + 1)
	case 0x23: // INC HL
		c.setHL(c.hl() + 1)
	case 0x33: // INC SP
		c.SP++
	default:
		return fmt.Errorf("unsupported incReg16 opcode: 0x%02X", opcode)
	}
	return nil
}

// decReg16 decrements a 16-bit register pair.
func decReg16(c *CPU, _ ...any) error {
	opcode := c.currentOpcode
	switch opcode {
	case 0x0B: // DEC BC
		c.setBC(c.bc() - 1)
	case 0x1B: // DEC DE
		c.setDE(c.de() - 1)
	case 0x2B: // DEC HL
		c.setHL(c.hl() - 1)
	case 0x3B: // DEC SP
		c.SP--
	default:
		return fmt.Errorf("unsupported decReg16 opcode: 0x%02X", opcode)
	}
	return nil
}

// addHl adds a 16-bit register pair to HL.
func addHl(c *CPU, _ ...any) error {
	hl := c.hl()
	c.MEMPTR = hl + 1

	opcode := c.currentOpcode
	var value uint16
	switch opcode {
	case 0x09: // ADD HL,BC
		value = c.bc()
	case 0x19: // ADD HL,DE
		value = c.de()
	case 0x29: // ADD HL,HL
		value = hl
	case 0x39: // ADD HL,SP
		value = c.SP
	default:
		value = c.bc()
	}

	result := c.add16(hl, value)
	c.setXY(uint8(result >> 8))
	c.setHL(result)
	return nil
}

// Rotate wrapper functions

// rlcaFunc performs rotate left circular accumulator.
func rlcaFunc(c *CPU) error {
	c.A = c.rlca(c.A)
	return nil
}

// rrcaFunc performs rotate right circular accumulator.
func rrcaFunc(c *CPU) error {
	c.A = c.rrca(c.A)
	return nil
}

// rla performs rotate left accumulator through carry.
func rla(c *CPU) error {
	newCarry := c.A >> 7
	c.A = (c.A << 1) | c.Flags.C
	c.setC(newCarry != 0)
	c.setH(false)
	c.setN(false)
	c.setXY(c.A)
	return nil
}

// rra performs rotate right accumulator through carry.
func rra(c *CPU) error {
	newCarry := c.A & 0x01
	c.A = (c.A >> 1) | (c.Flags.C << 7)
	c.setC(newCarry != 0)
	c.setH(false)
	c.setN(false)
	c.setXY(c.A)
	return nil
}

// I/O operations

// outPort outputs to port.
func outPort(c *CPU, params ...any) error {
	if len(params) < 1 {
		return ErrMissingParameter
	}

	var portAddr uint8
	switch param := params[0].(type) {
	case Port:
		portAddr = uint8(param)
	default:
		return ErrInvalidParameterType
	}

	// OUT (n),A - Output accumulator to port, address = A<<8 | n
	address := uint16(c.A)<<8 | uint16(portAddr)
	c.writePort(address, c.A)
	c.MEMPTR = uint16(portAddr+1) | uint16(c.A)<<8

	return nil
}

// inPort inputs from port.
func inPort(c *CPU, params ...any) error {
	if len(params) < 1 {
		return ErrMissingParameter
	}

	var portAddr uint8
	switch param := params[0].(type) {
	case Port:
		portAddr = uint8(param)
	default:
		return ErrInvalidParameterType
	}

	// IN A,(n) - Input from port to accumulator, address = A<<8 | n
	address := uint16(c.A)<<8 | uint16(portAddr)
	c.MEMPTR = address + 1
	c.A = c.readPort(address)

	return nil
}
