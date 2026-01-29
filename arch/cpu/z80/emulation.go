package z80

import (
	"fmt"
	"math/bits"
)

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

// bit tests bit n of value and sets flags.
func (c *CPU) bit(n uint8, value uint8) {
	bit := (value >> n) & 1
	bitIsZero := bit == 0

	setFlag(&c.Flags.Z, bitIsZero)
	setFlag(&c.Flags.P, bitIsZero) // P/V same as Z for BIT instruction
	c.setH(true)                   // Half carry is always set for BIT
	c.setN(false)                  // Indicates logical (not arithmetic) operation
	c.setXY(value)                 // X/Y from value being tested (not result)

	// S flag is affected differently for BIT instruction
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

// nop performs no operation.
func nop(_ *CPU) error {
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

	// Use the stored current opcode
	opcode := c.currentOpcode
	switch opcode {
	case 0x06: // LD B,n
		c.B = uint8(imm)
	case 0x0E: // LD C,n
		c.C = uint8(imm)
	case 0x16: // LD D,n
		c.D = uint8(imm)
	case 0x1E: // LD E,n
		c.E = uint8(imm)
	case 0x26: // LD H,n
		c.H = uint8(imm)
	case 0x2E: // LD L,n
		c.L = uint8(imm)
	case 0x3E: // LD A,n
		c.A = uint8(imm)
	default:
		// Default to A register for unknown opcodes
		c.A = uint8(imm)
	}
	return nil
}

// ldReg8 loads between 8-bit registers.
func ldReg8(c *CPU, params ...any) error {
	if len(params) < 2 {
		return ErrMissingParameter
	}

	dst, ok1 := params[0].(Register)
	src, ok2 := params[1].(Register)

	if !ok1 || !ok2 {
		return ErrInvalidParameterType
	}

	value := c.GetRegisterValue(uint8(src))
	c.SetRegisterValue(uint8(dst), value)
	return nil
}

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

// addA adds a value to the accumulator.
func addA(c *CPU, params ...any) error {
	if len(params) < 1 {
		return ErrMissingParameter
	}

	var value uint8

	switch param := params[0].(type) {
	case Register:
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
	case Register:
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
	case Register:
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
	case Register:
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
	case Register:
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
	case Register:
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

	// Calculate target address: PC after this 2-byte instruction + offset
	c.PC = uint16(int32(c.PC) + 2 + int32(offset))
	return nil
}

// Additional Z80 emulation functions

// ldReg16 loads a 16-bit immediate value into a register pair.
func ldReg16(c *CPU, params ...any) error {
	if len(params) < 1 {
		return ErrMissingParameter
	}

	imm, ok := params[0].(Immediate16)
	if !ok {
		return ErrInvalidParameterType
	}

	// For now, assume loading into BC (opcode 0x01)
	// In a full implementation, we'd need to determine the target register from the opcode
	c.setBC(uint16(imm))
	return nil
}

// ldIndirect loads between register pairs and memory locations.
func ldIndirect(c *CPU, _ ...any) error {
	// Use the stored current opcode
	opcode := c.currentOpcode

	switch opcode {
	case 0x02: // LD (BC),A - store A at (BC)
		c.memory.Write(c.bc(), c.A)
	case 0x0A: // LD A,(BC) - load A from (BC)
		c.A = c.memory.Read(c.bc())
	case 0x12: // LD (DE),A - store A at (DE)
		c.memory.Write(c.de(), c.A)
	case 0x1A: // LD A,(DE) - load A from (DE)
		c.A = c.memory.Read(c.de())
	default:
		return fmt.Errorf("unsupported indirect load opcode: 0x%02X", opcode)
	}
	return nil
}

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
	c.A = c.AltA
	c.Flags = c.AltFlags
	c.AltA = tempA
	c.AltFlags = tempF
	return nil
}

// addHl adds a 16-bit register pair to HL.
func addHl(c *CPU, _ ...any) error {
	// For opcode 0x09: ADD HL,BC - add BC to HL
	// In a full implementation, we'd determine which register pair from the opcode
	c.setHL(c.add16(c.hl(), c.bc()))
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
		// Calculate target address: PC after this 2-byte instruction + offset
		c.PC = uint16(int32(c.PC) + 2 + int32(offset))
	}
	return nil
}

// jrCond performs conditional relative jump.
func jrCond(c *CPU, params ...any) error {
	if len(params) < 1 {
		return ErrMissingParameter
	}
	offset, ok := params[0].(Relative)
	if !ok {
		return ErrInvalidParameterType
	}

	var shouldJump bool
	opcode := c.currentOpcode
	switch opcode {
	case 0x20: // JR NZ,e - jump if zero flag is clear
		shouldJump = c.Flags.Z == 0
	case 0x28: // JR Z,e - jump if zero flag is set
		shouldJump = c.Flags.Z != 0
	case 0x30: // JR NC,e - jump if carry flag is clear
		shouldJump = c.Flags.C == 0
	case 0x38: // JR C,e - jump if carry flag is set
		shouldJump = c.Flags.C != 0
	default:
		return fmt.Errorf("unsupported jrCond opcode: 0x%02X", opcode)
	}

	if shouldJump {
		// Calculate target address: PC after this 2-byte instruction + offset
		c.PC = uint16(int32(c.PC) + 2 + int32(offset))
		// Add extra cycles for taken branch (JR taken = 12 cycles, not taken = 7 cycles)
		c.cycles += 5
	}
	return nil
}

// ldExtended loads using extended addressing.
func ldExtended(c *CPU, params ...any) error {
	if len(params) < 1 {
		return ErrMissingParameter
	}
	address, ok := params[0].(Extended)
	if !ok {
		return ErrInvalidParameterType
	}

	opcode := c.currentOpcode
	switch opcode {
	case 0x22: // LD (nn),HL - store HL to memory address nn
		c.memory.WriteWord(uint16(address), c.hl())
	case 0x2A: // LD HL,(nn) - load HL from memory address nn
		value := c.memory.ReadWord(uint16(address))
		c.setHL(value)
	case 0x32: // LD (nn),A - store A to memory address nn
		c.memory.Write(uint16(address), c.A)
	case 0x3A: // LD A,(nn) - load A from memory address nn
		value := c.memory.Read(uint16(address))
		c.A = value
	default:
		return fmt.Errorf("unsupported ldExtended opcode: 0x%02X", opcode)
	}
	return nil
}

// daaAdditionMode calculates correction for DAA in addition mode.
func (c *CPU) daaAdditionMode() (uint8, bool) {
	correction := uint8(0)
	carrySet := false

	if c.Flags.H != 0 || (c.A&0x0F) > 9 {
		correction |= 0x06
	}
	if c.Flags.C != 0 || c.A > 0x99 {
		correction |= 0x60
		carrySet = true
	}
	return correction, carrySet
}

// daaSubtractionMode calculates correction for DAA in subtraction mode.
func (c *CPU) daaSubtractionMode() (uint8, bool) {
	correction := uint8(0)
	carrySet := false

	if c.Flags.H != 0 || (c.A&0x0F) > 9 {
		correction |= 0x06
	}
	if c.Flags.C != 0 {
		correction |= 0x60
		carrySet = true
	}
	return correction, carrySet
}

// daa performs decimal adjust accumulator.
func daa(c *CPU) error {
	originalA := c.A // Save original value for H flag calculation
	var correction uint8
	var carrySet bool

	if c.Flags.N == 0 {
		correction, carrySet = c.daaAdditionMode()
		c.A += correction
	} else {
		correction, carrySet = c.daaSubtractionMode()
		c.A -= correction
	}

	// Set flags
	c.setS(c.A)
	c.setZ(c.A)
	c.setP(c.A)  // Parity flag
	c.setXY(c.A) // Set undocumented X and Y flags
	c.setC(carrySet)
	c.setH((originalA^c.A)&0x10 != 0) // H = XOR of original and result (bit 4)
	// N flag remains unchanged

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
func incIndirect(c *CPU, _ ...any) error {
	// For opcode 0x34: INC (HL) - increment memory at HL
	address := c.hl()
	value := c.memory.Read(address)
	newValue := value + 1
	c.memory.Write(address, newValue)

	// Set flags (S, Z, H, PV, N)
	c.setS(newValue)
	c.setZ(newValue)
	c.setH((value & 0x0F) == 0x0F) // Half carry when low nibble overflows
	c.setPOverflow(value == 0x7F)  // Overflow when 0x7F -> 0x80
	c.setN(false)
	return nil
}

// decIndirect decrements memory location pointed to by register pair.
func decIndirect(c *CPU, _ ...any) error {
	// For opcode 0x35: DEC (HL) - decrement memory at HL
	address := c.hl()
	value := c.memory.Read(address)
	newValue := value - 1
	c.memory.Write(address, newValue)

	// Set flags (S, Z, H, PV, N)
	c.setS(newValue)
	c.setZ(newValue)
	c.setH((value & 0x0F) == 0x00) // Half carry when low nibble underflows
	c.setPOverflow(value == 0x80)  // Overflow when 0x80 -> 0x7F
	c.setN(true)
	return nil
}

// ldIndirectImm loads immediate value to indirect memory location.
func ldIndirectImm(c *CPU, params ...any) error {
	// For opcode 0x36: LD (HL),n - load immediate to memory at HL
	if len(params) < 1 {
		return ErrMissingParameter
	}

	var immediate uint8
	switch v := params[0].(type) {
	case Immediate8:
		immediate = uint8(v)
	case uint8:
		immediate = v
	default:
		return ErrInvalidParameterType
	}

	address := c.hl()
	c.memory.Write(address, immediate)
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
	case Register:
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
	case Register:
		value = c.GetRegisterValue(uint8(param))
	case Immediate8:
		value = uint8(param)
	default:
		return ErrInvalidParameterType
	}

	c.A = c.sbc(c.A, value)
	return nil
}

// checkCondition returns true if the condition is met based on the opcode.
func (c *CPU) checkCondition(opcode uint8) bool {
	switch opcode {
	// NZ (Not Zero) - Z flag clear
	case 0xC0, 0xC2, 0xC4:
		return c.Flags.Z == 0
	// Z (Zero) - Z flag set
	case 0xC8, 0xCA, 0xCC:
		return c.Flags.Z != 0
	// NC (Not Carry) - C flag clear
	case 0xD0, 0xD2, 0xD4:
		return c.Flags.C == 0
	// C (Carry) - C flag set
	case 0xD8, 0xDA, 0xDC:
		return c.Flags.C != 0
	// PO (Parity Odd) - P flag clear
	case 0xE0, 0xE2, 0xE4:
		return c.Flags.P == 0
	// PE (Parity Even) - P flag set
	case 0xE8, 0xEA, 0xEC:
		return c.Flags.P != 0
	// P (Plus/Positive) - S flag clear
	case 0xF0, 0xF2, 0xF4:
		return c.Flags.S == 0
	// M (Minus/Negative) - S flag set
	case 0xF8, 0xFA, 0xFC:
		return c.Flags.S != 0
	default:
		return false
	}
}

// retCond performs conditional return.
func retCond(c *CPU) error {
	if c.checkCondition(c.currentOpcode) {
		c.PC = c.pop16()
		// Add extra cycles for taken return (RET taken = 11 cycles, not taken = 5 cycles)
		c.cycles += 6
	}
	return nil
}

// popReg16 pops 16-bit register from stack.
func popReg16(c *CPU, _ ...any) error {
	value := c.pop16()
	opcode := c.currentOpcode
	switch opcode {
	case 0xC1: // POP BC
		c.setBC(value)
	case 0xD1: // POP DE
		c.setDE(value)
	case 0xE1: // POP HL
		c.setHL(value)
	case 0xF1: // POP AF
		c.setAF(value)
	default:
		return fmt.Errorf("unsupported popReg16 opcode: 0x%02X", opcode)
	}
	return nil
}

// jpCond performs conditional jump.
func jpCond(c *CPU, params ...any) error {
	if len(params) < 1 {
		return ErrMissingParameter
	}
	address, ok := params[0].(Extended)
	if !ok {
		return ErrInvalidParameterType
	}

	if c.checkCondition(c.currentOpcode) {
		c.PC = uint16(address)
	}
	return nil
}

// callCond performs conditional call.
func callCond(c *CPU, params ...any) error {
	if len(params) < 1 {
		return ErrMissingParameter
	}
	address, ok := params[0].(Extended)
	if !ok {
		return ErrInvalidParameterType
	}

	if c.checkCondition(c.currentOpcode) {
		c.push16(c.PC)
		c.PC = uint16(address)
		// Add extra cycles for taken call (CALL taken = 17 cycles, not taken = 10 cycles)
		c.cycles += 7
	}
	return nil
}

// pushReg16 pushes 16-bit register to stack.
func pushReg16(c *CPU, _ ...any) error {
	opcode := c.currentOpcode
	var value uint16
	switch opcode {
	case 0xC5: // PUSH BC
		value = c.bc()
	case 0xD5: // PUSH DE
		value = c.de()
	case 0xE5: // PUSH HL
		value = c.hl()
	case 0xF5: // PUSH AF
		value = c.af()
	default:
		return fmt.Errorf("unsupported pushReg16 opcode: 0x%02X", opcode)
	}
	c.push16(value)
	return nil
}

// rst performs restart (call to fixed address).
func rst(c *CPU, _ ...any) error {
	// RST pushes current PC to stack and jumps to fixed address
	c.push16(c.PC)

	// Calculate restart vector from opcode
	opcode := c.currentOpcode
	var vector uint16
	switch opcode {
	case 0xC7: // RST 00H
		vector = 0x0000
	case 0xCF: // RST 08H
		vector = 0x0008
	case 0xD7: // RST 10H
		vector = 0x0010
	case 0xDF: // RST 18H
		vector = 0x0018
	case 0xE7: // RST 20H
		vector = 0x0020
	case 0xEF: // RST 28H
		vector = 0x0028
	case 0xF7: // RST 30H
		vector = 0x0030
	case 0xFF: // RST 38H
		vector = 0x0038
	default:
		return fmt.Errorf("unsupported rst opcode: 0x%02X", opcode)
	}

	c.PC = vector
	return nil
}

// ret returns from subroutine.
func ret(c *CPU) error {
	// Pop return address from stack and jump to it
	c.PC = c.pop16()
	return nil
}

// call calls a subroutine.
func call(c *CPU, params ...any) error {
	if len(params) < 1 {
		return ErrMissingParameter
	}

	addr, ok := params[0].(Extended)
	if !ok {
		return ErrInvalidParameterType
	}

	// Push return address (PC + 3 for the 3-byte CALL instruction) to stack
	// The PC still points to the CALL instruction when this function runs
	returnAddr := c.PC + 3
	c.push16(returnAddr)

	// Jump to the called address
	c.PC = uint16(addr)
	return nil
}

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

	// OUT (n),A - Output accumulator to port
	if c.opts.ioHandler != nil {
		c.opts.ioHandler.WritePort(portAddr, c.A)
	}

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

	// IN A,(n) - Input from port to accumulator
	if c.opts.ioHandler != nil {
		c.A = c.opts.ioHandler.ReadPort(portAddr)
	} else {
		c.A = 0xFF
	}

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

	c.B = c.AltB
	c.C = c.AltC
	c.D = c.AltD
	c.E = c.AltE
	c.H = c.AltH
	c.L = c.AltL

	c.AltB = tempB
	c.AltC = tempC
	c.AltD = tempD
	c.AltE = tempE
	c.AltH = tempH
	c.AltL = tempL
	return nil
}

// exSp exchanges top of stack with register pair.
func exSp(c *CPU, _ ...any) error {
	// EX (SP),HL - Exchange HL with word at top of stack
	// Read word from stack (SP and SP+1)
	low := c.memory.Read(c.SP)
	high := c.memory.Read(c.SP + 1)
	stackValue := uint16(high)<<8 | uint16(low)

	// Get current HL value
	hlValue := c.hl()

	// Write HL to stack
	c.memory.Write(c.SP, uint8(hlValue))
	c.memory.Write(c.SP+1, uint8(hlValue>>8))

	// Set HL to old stack value
	c.setHL(stackValue)

	return nil
}

// jpIndirect performs indirect jump.
func jpIndirect(c *CPU, _ ...any) error {
	// JP (HL) - Jump to address in HL register
	c.PC = c.hl()
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
func ldSp(c *CPU, _ ...any) error {
	c.SP = uint16(c.H)<<8 | uint16(c.L)
	return nil
}

// boolToUint8 converts a boolean to 1 or 0 as uint8.
func boolToUint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

// calculateParity calculates the parity of a byte (true if even parity).
func calculateParity(value uint8) bool {
	count := bits.OnesCount8(value)
	return count%2 == 0
}

// performShiftRotateOperation performs shift/rotate operations based on opcode.
func performShiftRotateOperation(value, opcode, oldCarry uint8) (uint8, bool) {
	switch {
	case opcode <= 0x07: // RLC
		carry := (value & 0x80) != 0
		return (value << 1) | boolToUint8(carry), carry
	case opcode <= 0x0F: // RRC
		carry := (value & 0x01) != 0
		return (value >> 1) | (boolToUint8(carry) << 7), carry
	case opcode <= 0x17: // RL
		carry := (value & 0x80) != 0
		return (value << 1) | oldCarry, carry
	case opcode <= 0x1F: // RR
		carry := (value & 0x01) != 0
		return (value >> 1) | (oldCarry << 7), carry
	case opcode <= 0x27: // SLA
		carry := (value & 0x80) != 0
		return value << 1, carry
	case opcode <= 0x2F: // SRA
		carry := (value & 0x01) != 0
		return (value >> 1) | (value & 0x80), carry
	case opcode <= 0x37: // SLL
		carry := (value & 0x80) != 0
		return (value << 1) | 0x01, carry
	default: // SRL
		carry := (value & 0x01) != 0
		return value >> 1, carry
	}
}

// setShiftRotateFlags sets flags for shift/rotate operations.
func setShiftRotateFlags(c *CPU, result uint8, carry bool) {
	c.setSZ(result)
	c.setPOverflow(calculateParity(result))
	c.setH(false)
	c.setN(false)
	c.setC(carry)
}

// Helper methods for emulation tests

// inc16 increments a 16-bit value.
func (c *CPU) inc16(value uint16) uint16 {
	return value + 1
}

// dec16 decrements a 16-bit value.
func (c *CPU) dec16(value uint16) uint16 {
	return value - 1
}

// addHL adds a 16-bit value to HL register pair.
func (c *CPU) addHL(hl, value uint16) uint16 {
	result32 := uint32(hl) + uint32(value)
	result := uint16(result32)

	// Set flags for 16-bit addition
	c.setC(result32 > 0xFFFF)                   // Carry if result > 65535
	c.setH((hl&0x0FFF)+(value&0x0FFF) > 0x0FFF) // Half carry on bit 11
	c.setN(false)                               // Clear N flag for addition
	// Note: Z and S flags are not affected by ADD HL

	return result
}

// Helper methods for load operations

// ldReg8 loads a value to a register (test helper).
func (c *CPU) ldReg8(dst *uint8, src uint8) {
	*dst = src
}

// ldImm8 loads immediate value to register (test helper).
func (c *CPU) ldImm8(dst *uint8, value uint8) {
	*dst = value
}

// ldMemToReg8 loads memory value to register.
func (c *CPU) ldMemToReg8(dst *uint8, addr uint16) {
	*dst = c.memory.Read(addr)
}

// ldRegToMem8 loads register value to memory.
func (c *CPU) ldRegToMem8(addr uint16, src uint8) {
	c.memory.Write(addr, src)
}

// Helper methods for memory block operations

// ldi executes Load and Increment operation.
func (c *CPU) ldi() {
	hl := c.hl()
	de := c.de()
	bc := c.bc()

	// Copy byte from (HL) to (DE)
	value := c.memory.Read(hl)
	c.memory.Write(de, value)

	// Increment HL and DE, decrement BC
	c.setHL(hl + 1)
	c.setDE(de + 1)
	c.setBC(bc - 1)

	// Set P/V flag based on BC
	c.setPOverflow(bc != 1) // P/V set if BC-1 != 0
	c.setH(false)
	c.setN(false)
}

// ldd executes Load and Decrement operation.
func (c *CPU) ldd() {
	hl := c.hl()
	de := c.de()
	bc := c.bc()

	// Copy byte from (HL) to (DE)
	value := c.memory.Read(hl)
	c.memory.Write(de, value)

	// Decrement HL, DE, and BC
	c.setHL(hl - 1)
	c.setDE(de - 1)
	c.setBC(bc - 1)

	// Set P/V flag based on BC
	c.setPOverflow(bc != 1) // P/V set if BC-1 != 0
	c.setH(false)
	c.setN(false)
}

// Helper methods for jump and call operations

// jp executes unconditional jump.
func (c *CPU) jp(addr uint16) {
	c.PC = addr
}

// jr executes relative jump.
func (c *CPU) jr(offset int8) {
	c.PC = uint16(int32(c.PC) + int32(offset))
}

// jpZ executes conditional jump if Z flag is set.
func (c *CPU) jpZ(addr uint16) {
	if c.Flags.Z == 1 {
		c.PC = addr
	}
}

// djnz executes Decrement and Jump if Not Zero.
func (c *CPU) djnz(offset int8) {
	c.B--
	if c.B != 0 {
		c.PC = uint16(int32(c.PC) + int32(offset))
	}
}

// Helper methods for exchange operations

// exDEHL exchanges DE and HL register pairs.
func (c *CPU) exDEHL() {
	d, e := c.D, c.E
	c.D, c.E = c.H, c.L
	c.H, c.L = d, e
}

// exx exchanges BC, DE, HL with shadow registers.
func (c *CPU) exx() {
	c.B, c.AltB = c.AltB, c.B
	c.C, c.AltC = c.AltC, c.C
	c.D, c.AltD = c.AltD, c.D
	c.E, c.AltE = c.AltE, c.E
	c.H, c.AltH = c.AltH, c.H
	c.L, c.AltL = c.AltL, c.L
}

// exAF exchanges AF with shadow AF.
func (c *CPU) exAF() {
	c.A, c.AltA = c.AltA, c.A
	c.Flags, c.AltFlags = c.AltFlags, c.Flags
}

// getFlags returns current flags as uint8.
func (c *CPU) getFlags() uint8 {
	return c.GetFlags()
}

// setFlagsFromUint8 sets flags struct from uint8 value.
func (c *CPU) setFlagsFromUint8(flags *Flags, value uint8) {
	flags.C = value & 0x01
	flags.N = (value >> 1) & 0x01
	flags.P = (value >> 2) & 0x01
	flags.X = (value >> 3) & 0x01
	flags.H = (value >> 4) & 0x01
	flags.Y = (value >> 5) & 0x01
	flags.Z = (value >> 6) & 0x01
	flags.S = (value >> 7) & 0x01
}

// getFlagsAsUint8 gets flags struct as uint8 value.
func (c *CPU) getFlagsAsUint8(flags Flags) uint8 {
	return flags.C | (flags.N << 1) | (flags.P << 2) | (flags.X << 3) |
		(flags.H << 4) | (flags.Y << 5) | (flags.Z << 6) | (flags.S << 7)
}
