package m6502

import (
	"fmt"
	"math"
)

// adc - Add with Carry.
func adc(c *CPU, params ...any) error {
	a := c.A
	value, err := c.memory.ReadAddressModes(true, params...)
	if err != nil {
		return err
	}

	if c.Flags.D != 0 && c.opts.variant != VariantNES6502 {
		if c.opts.variant >= Variant65C02 {
			adcDecimal65C02(c, a, value)
		} else {
			adcDecimalNMOS(c, a, value)
		}
		return nil
	}

	sum := int(c.A) + int(c.Flags.C) + int(value)
	c.A = uint8(sum)
	c.setZN(c.A)

	if sum > math.MaxUint8 {
		c.Flags.C = 1
	} else {
		c.Flags.C = 0
	}
	c.setV((a^value)&0x80 == 0 && (a^c.A)&0x80 != 0)
	return nil
}

// adcDecimalNMOS performs NMOS 6502 BCD decimal mode addition.
// On NMOS 6502, flags behave as follows in decimal mode:
//   - Z: from the 8-bit binary (uncorrected) result
//   - N, V: from the partially-corrected result (binary + lo-nibble adjustment)
//   - A, C: from the fully BCD-corrected result
//
// Reference: http://www.6502.org/tutorials/decimal_mode.html
func adcDecimalNMOS(c *CPU, a, value uint8) {
	carry := int(c.Flags.C)
	binarySum := int(a) + int(value) + carry

	// Z is from the raw binary result.
	c.setZ(uint8(binarySum))

	// BCD correction: each nibble correction always contributes exactly 1 carry to the
	// next nibble, regardless of how far over the threshold the raw nibble sum is.
	lo := int(a&0x0F) + int(value&0x0F) + carry
	loCarry := 0
	if lo >= 10 {
		lo = (lo + 6) & 0x0F
		loCarry = 1
	}
	hi := int(a>>4) + int(value>>4) + loCarry
	hiCarry := 0
	if hi >= 10 {
		hi = (hi + 6) & 0x0F
		hiCarry = 1
	}

	// N and V come from the partially-corrected intermediate (after lo adjustment, before hi).
	// This matches the NMOS 6502 pipeline behavior where N/V are latched from AH at step 3.
	partial := uint8(int(a&0xF0) + int(value&0xF0) + lo + loCarry*16)
	c.setN(partial)
	c.setV((a^value)&0x80 == 0 && (a^partial)&0x80 != 0)

	c.A = uint8((hi << 4) | lo)
	c.Flags.C = uint8(hiCarry)
}

// adcDecimal65C02 performs 65C02 BCD decimal mode addition.
// On 65C02, flags behave as follows in decimal mode:
//   - V: from the lo-nibble-corrected partial result (same source as NMOS)
//   - A, C: from the fully BCD-corrected result
//   - N, Z: from the BCD-corrected result
//
// Reference: http://www.6502.org/tutorials/decimal_mode.html
func adcDecimal65C02(c *CPU, a, value uint8) {
	carry := int(c.Flags.C)

	// BCD correction: each nibble correction contributes exactly 1 carry.
	lo := int(a&0x0F) + int(value&0x0F) + carry
	loCarry := 0
	if lo >= 10 {
		lo = (lo + 6) & 0x0F
		loCarry = 1
	}
	hi := int(a>>4) + int(value>>4) + loCarry

	// V is from the lo-corrected partial result (before hi-nibble correction).
	// This matches hardware captures: same V source as NMOS 6502.
	partial := uint8((hi&0x0F)<<4 | lo)
	c.setV((a^value)&0x80 == 0 && (a^partial)&0x80 != 0)

	hiCarry := 0
	if hi >= 10 {
		hi = (hi + 6) & 0x0F
		hiCarry = 1
	}

	c.A = uint8((hi << 4) | lo)
	c.Flags.C = uint8(hiCarry)

	// N and Z from the BCD-corrected result.
	c.setZN(c.A)
}

// and - AND with accumulator.
func and(c *CPU, params ...any) error {
	value, err := c.memory.ReadAddressModes(true, params...)
	if err != nil {
		return err
	}
	c.A &= value
	c.setZN(c.A)
	return nil
}

// asl - Arithmetic Shift Left.
func asl(c *CPU, params ...any) error {
	if hasAccumulatorParam(params...) {
		c.Flags.C = (c.A >> 7) & 1
		c.A <<= 1
		c.setZN(c.A)
		return nil
	}

	val, err := c.memory.ReadAddressModes(false, params...)
	if err != nil {
		return err
	}
	c.Flags.C = (val >> 7) & 1
	val <<= 1
	c.setZN(val)
	return c.memory.WriteAddressModes(val, params...)
}

// bcc - Branch if Carry Clear.
func bcc(c *CPU, params ...any) error {
	c.branch(c.Flags.C == 0, params[0])
	return nil
}

// bcs - Branch if Carry Set.
func bcs(c *CPU, params ...any) error {
	c.branch(c.Flags.C != 0, params[0])
	return nil
}

// beq - Branch if Equal.
func beq(c *CPU, params ...any) error {
	c.branch(c.Flags.Z != 0, params[0])
	return nil
}

// bit - Bit Test.
func bit(c *CPU, params ...any) error {
	value, err := c.memory.ReadAbsolute(params[0], nil)
	if err != nil {
		return err
	}
	c.setV((value>>6)&1 == 1)
	c.setZ(value & c.A)
	c.setN(value)
	return nil
}

// bmi - Branch if Minus.
func bmi(c *CPU, params ...any) error {
	c.branch(c.Flags.N != 0, params[0])
	return nil
}

// bne - Branch if Not Equal.
func bne(c *CPU, params ...any) error {
	c.branch(c.Flags.Z == 0, params[0])
	return nil
}

// bpl - Branch if Positive.
func bpl(c *CPU, params ...any) error {
	c.branch(c.Flags.N == 0, params[0])
	return nil
}

// brk - Force Interrupt.
func brk(c *CPU) error {
	// BRK is a 2-byte instruction, the second byte is a signature/padding byte
	c.push16(c.PC + 2) // Push PC+2 to skip the signature byte

	// The B flag should be set when pushing the status to distinguish BRK from IRQ
	c.Flags.B = 1
	f := c.GetFlags()
	f |= 0b0010_0000 // Ensure unused flag is set
	c.push(f)
	c.Flags.I = 1 // Disable interrupts

	// 65C02: Clear D flag after pushing status on BRK
	if c.opts.variant >= Variant65C02 {
		c.Flags.D = 0
	}

	c.PC = c.irqAddress

	c.mu.Lock()
	c.triggerIrq = false
	c.irqRunning = true
	c.mu.Unlock()
	return nil
}

// bvc - Branch if Overflow Clear.
func bvc(c *CPU, params ...any) error {
	c.branch(c.Flags.V == 0, params[0])
	return nil
}

// bvs - Branch if Overflow Set.
func bvs(c *CPU, params ...any) error {
	c.branch(c.Flags.V != 0, params[0])
	return nil
}

// clc - Clear Carry Flag.
func clc(c *CPU) error {
	c.Flags.C = 0
	return nil
}

// cld - Clear Decimal Mode.
func cld(c *CPU) error {
	c.Flags.D = 0
	return nil
}

// cli - Clear Interrupt Disable.
func cli(c *CPU) error {
	c.Flags.I = 0
	return nil
}

// clv - Clear Overflow Flag.
func clv(c *CPU) error {
	c.Flags.V = 0
	return nil
}

// cmp - Compare the contents of A.
func cmp(c *CPU, params ...any) error {
	val, err := c.memory.ReadAddressModes(true, params...)
	if err != nil {
		return err
	}
	c.compare(c.A, val)
	return nil
}

// cpx - Compare the contents of X.
func cpx(c *CPU, params ...any) error {
	val, err := c.memory.ReadAddressModes(true, params[0])
	if err != nil {
		return err
	}
	c.compare(c.X, val)
	return nil
}

// cpy - Compare the contents of Y.
func cpy(c *CPU, params ...any) error {
	val, err := c.memory.ReadAddressModes(true, params[0])
	if err != nil {
		return err
	}
	c.compare(c.Y, val)
	return nil
}

// dec - Decrement memory.
func dec(c *CPU, params ...any) error {
	val, err := c.memory.ReadAddressModes(false, params...)
	if err != nil {
		return err
	}
	val--
	if err = c.memory.WriteAddressModes(val, params...); err != nil {
		return err
	}
	c.setZN(val)
	return nil
}

// dex - Decrement X Register.
func dex(c *CPU) error {
	c.X--
	c.setZN(c.X)
	return nil
}

// dey - Decrement Y Register.
func dey(c *CPU) error {
	c.Y--
	c.setZN(c.Y)
	return nil
}

// eor - Exclusive OR - XOR.
func eor(c *CPU, params ...any) error {
	value, err := c.memory.ReadAddressModes(true, params...)
	if err != nil {
		return err
	}
	c.A ^= value
	c.setZN(c.A)
	return nil
}

// inc - Increments memory.
func inc(c *CPU, params ...any) error {
	val, err := c.memory.ReadAddressModes(false, params...)
	if err != nil {
		return err
	}
	val++
	if err = c.memory.WriteAddressModes(val, params...); err != nil {
		return err
	}
	c.setZN(val)
	return nil
}

// inx - Increment X Register.
func inx(c *CPU) error {
	c.X++
	c.setZN(c.X)
	return nil
}

// iny - Increment Y Register.
func iny(c *CPU) error {
	c.Y++
	c.setZN(c.Y)
	return nil
}

// jmp - jump to address.
func jmp(c *CPU, params ...any) error {
	param := params[0]
	switch address := param.(type) {
	case Absolute:
		c.PC = uint16(address)
	case Indirect:
		c.PC = c.memory.ReadWordBug(uint16(address))
	default:
		return fmt.Errorf("%w: jmp mode type %T", ErrUnsupportedAddressingMode, param)
	}
	return nil
}

// jsr - jump to subroutine.
func jsr(c *CPU, params ...any) error {
	if len(params) == 0 {
		return fmt.Errorf("%w: jsr missing address parameter", ErrMissingParameter)
	}

	addr, ok := params[0].(Absolute)
	if !ok {
		return fmt.Errorf("%w: jsr invalid address parameter type", ErrInvalidParameterType)
	}

	returnAddr := c.PC + 2
	c.push(uint8(returnAddr >> 8))
	c.push(uint8(returnAddr & 0xFF))
	c.PC = uint16(addr)
	return nil
}

// lda - Load Accumulator - load a byte into A.
func lda(c *CPU, params ...any) error {
	val, err := c.memory.ReadAddressModes(true, params...)
	if err != nil {
		return err
	}
	c.A = val
	c.setZN(c.A)
	return nil
}

// ldx - Load X Register - load a byte into X.
func ldx(c *CPU, params ...any) error {
	val, err := c.memory.ReadAddressModes(true, params...)
	if err != nil {
		return err
	}
	c.X = val
	c.setZN(c.X)
	return nil
}

// ldy - Load Y Register - load a byte into Y.
func ldy(c *CPU, params ...any) error {
	val, err := c.memory.ReadAddressModes(true, params...)
	if err != nil {
		return err
	}
	c.Y = val
	c.setZN(c.Y)
	return nil
}

// lsr - Logical Shift Right.
func lsr(c *CPU, params ...any) error {
	if hasAccumulatorParam(params...) {
		c.Flags.C = c.A & 1
		c.A >>= 1
		c.setZN(c.A)
		return nil
	}

	val, err := c.memory.ReadAddressModes(false, params...)
	if err != nil {
		return err
	}
	c.Flags.C = val & 1
	val >>= 1
	c.setZN(val)
	return c.memory.WriteAddressModes(val, params...)
}

// nop - No Operation.
func nop(_ *CPU) error {
	return nil
}

// kil - KIL/JAM: halt the CPU. On real hardware this freezes the processor;
// the test-visible effect is PC advancing by 1 (past the opcode byte only).
func kil(_ *CPU) error {
	return nil
}

// ora - OR with Accumulator.
func ora(c *CPU, params ...any) error {
	value, err := c.memory.ReadAddressModes(true, params...)
	if err != nil {
		return err
	}
	c.A |= value
	c.setZN(c.A)
	return nil
}

// pha - Push Accumulator - push A content to stack.
func pha(c *CPU) error {
	c.push(c.A)
	return nil
}

// php - Push Processor Status - push status flags to stack.
func php(c *CPU) error {
	f := c.GetFlags()
	f |= 0b0001_0000 // break is set to 1
	c.push(f)
	return nil
}

// pla - Pull Accumulator - pull A content from stack.
func pla(c *CPU) error {
	c.A = c.pop()
	c.setZN(c.A)
	return nil
}

// plp - Pull Processor Status - pull status flags from stack.
func plp(c *CPU) error {
	f := c.pop()
	f &= 0b1110_1111 // break flag is ignored
	f |= 0b0010_0000 // unused flag is set
	c.setFlags(f)
	return nil
}

// rol - Rotate Left.
func rol(c *CPU, params ...any) error {
	cFlag := c.Flags.C
	if hasAccumulatorParam(params...) {
		c.Flags.C = (c.A >> 7) & 1
		c.A = (c.A << 1) | cFlag
		c.setZN(c.A)
		return nil
	}

	val, err := c.memory.ReadAddressModes(false, params...)
	if err != nil {
		return err
	}
	c.Flags.C = (val >> 7) & 1
	val = (val << 1) | cFlag
	c.setZN(val)
	return c.memory.WriteAddressModes(val, params...)
}

// ror - Rotate Right.
func ror(c *CPU, params ...any) error {
	cFlag := c.Flags.C
	if hasAccumulatorParam(params...) {
		c.Flags.C = c.A & 1
		c.A = (c.A >> 1) | (cFlag << 7)
		c.setZN(c.A)
		return nil
	}

	val, err := c.memory.ReadAddressModes(false, params...)
	if err != nil {
		return err
	}
	c.Flags.C = val & 1
	val = (val >> 1) | (cFlag << 7)
	c.setZN(val)
	return c.memory.WriteAddressModes(val, params...)
}

// rti - Return from Interrupt.
func rti(c *CPU) error {
	b := c.pop()
	b &= 0b1110_1111 // break flag is ignored
	b |= 0b0010_0000 // unused flag is set
	c.setFlags(b)
	c.PC = c.pop16()

	// lock is already taken
	c.irqRunning = false
	c.nmiRunning = false
	return nil
}

// rts - return from subroutine.
func rts(c *CPU) error {
	c.PC = c.pop16() + 1
	return nil
}

// sbc - subtract with Carry.
func sbc(c *CPU, params ...any) error {
	a := c.A
	value, err := c.memory.ReadAddressModes(true, params...)
	if err != nil {
		return err
	}

	if c.Flags.D != 0 && c.opts.variant != VariantNES6502 {
		if c.opts.variant >= Variant65C02 {
			sbcDecimal65C02(c, a, value)
		} else {
			sbcDecimalNMOS(c, a, value)
		}
		return nil
	}

	sub := int(c.A) - int(value) - (1 - int(c.Flags.C))
	c.A = uint8(sub)
	c.setZN(c.A)

	if sub >= 0 {
		c.Flags.C = 1
	} else {
		c.Flags.C = 0
	}
	c.setV((a^value)&0x80 != 0 && (a^c.A)&0x80 != 0)
	return nil
}

// sbcDecimalNMOS performs NMOS 6502 BCD decimal mode subtraction.
// On NMOS, N/Z/V are set from the binary result, not the BCD-corrected result.
// Reference: http://www.6502.org/tutorials/decimal_mode.html
func sbcDecimalNMOS(c *CPU, a, value uint8) {
	borrow := 1 - int(c.Flags.C)
	binaryDiff := int(a) - int(value) - borrow

	// N, Z, V flags are set from the BINARY result on NMOS 6502.
	binaryResult := uint8(binaryDiff)
	c.setZN(binaryResult)
	c.setV((a^value)&0x80 != 0 && (a^binaryResult)&0x80 != 0)

	// BCD correction for lower nibble.
	lo := int(a&0x0F) - int(value&0x0F) - borrow
	halfBorrow := 0
	if lo < 0 {
		lo -= 6
		halfBorrow = 1
	}

	// BCD correction for upper nibble.
	hi := int(a>>4) - int(value>>4) - halfBorrow
	if hi < 0 {
		hi -= 6
	}

	c.A = uint8((lo & 0x0F) | (hi << 4))
	if binaryDiff >= 0 {
		c.Flags.C = 1
	} else {
		c.Flags.C = 0
	}
}

// sbcDecimal65C02 performs 65C02 BCD decimal mode subtraction.
// On 65C02:
//   - V: from the binary result
//   - C: from the binary result (set if no borrow)
//   - N, Z: from the BCD-corrected result
//   - A: BCD-corrected (lo digit borrows when (a&0xF) < (val&0xF)+borrow;
//     hi digit borrows when binary result is negative)
func sbcDecimal65C02(c *CPU, a, value uint8) {
	borrow := 1 - int(c.Flags.C)
	binaryDiff := int(a) - int(value) - borrow

	// V and C from the binary result.
	binaryResult := uint8(binaryDiff)
	c.setV((a^value)&0x80 != 0 && (a^binaryResult)&0x80 != 0)
	if binaryDiff >= 0 {
		c.Flags.C = 1
	} else {
		c.Flags.C = 0
	}

	// BCD correction: lo nibble borrows when (a&0xF) < (value&0xF) + borrow.
	temp := binaryDiff
	if int(a&0x0F) < int(value&0x0F)+borrow {
		temp -= 6
	}
	// Hi nibble borrows when the binary result was negative.
	if binaryDiff < 0 {
		temp -= 0x60
	}

	c.A = uint8(temp)
	// N and Z are set from the BCD-corrected result on 65C02.
	c.setZN(c.A)
}

// sec - Set Carry Flag.
func sec(c *CPU) error {
	c.Flags.C = 1
	return nil
}

// sed - Set Decimal Flag.
func sed(c *CPU) error {
	c.Flags.D = 1
	return nil
}

// sei - Set Interrupt Disable.
func sei(c *CPU) error {
	c.Flags.I = 1
	return nil
}

// sta - Store Accumulator.
func sta(c *CPU, params ...any) error {
	return c.memory.WriteAddressModes(c.A, params...)
}

// stx - Store X Register.
func stx(c *CPU, params ...any) error {
	return c.memory.WriteAddressModes(c.X, params...)
}

// sty - Store Y Register.
func sty(c *CPU, params ...any) error {
	return c.memory.WriteAddressModes(c.Y, params...)
}

// tax - Transfer Accumulator to X.
func tax(c *CPU) error {
	c.X = c.A
	c.setZN(c.X)
	return nil
}

// tay - Transfer Accumulator to Y.
func tay(c *CPU) error {
	c.Y = c.A
	c.setZN(c.Y)
	return nil
}

// tsx - Transfer Stack Pointer to X.
func tsx(c *CPU) error {
	c.X = c.SP
	c.setZN(c.X)
	return nil
}

// txa - Transfer X to Accumulator.
func txa(c *CPU) error {
	c.A = c.X
	c.setZN(c.A)
	return nil
}

// txs - Transfer X to Stack Pointer.
func txs(c *CPU) error {
	c.SP = c.X
	return nil
}

// tya - Transfer Y to Accumulator.
func tya(c *CPU) error {
	c.A = c.Y
	c.setZN(c.A)
	return nil
}

// unofficial instructions

func dcp(c *CPU, params ...any) error {
	if err := dec(c, params...); err != nil {
		return err
	}
	val, err := c.memory.ReadAddressModes(false, params...)
	if err != nil {
		return err
	}
	c.compare(c.A, val)
	return nil
}

func isc(c *CPU, params ...any) error {
	if err := inc(c, params...); err != nil {
		return err
	}
	return sbc(c, params...)
}

// las implements LAS/LAR (0xBB): result = mem & SP; A = X = SP = result.
func las(c *CPU, params ...any) error {
	val, err := c.memory.ReadAddressModes(false, params...)
	if err != nil {
		return err
	}
	result := val & c.SP
	c.A = result
	c.X = result
	c.SP = result
	c.setZN(result)
	return nil
}

func lax(c *CPU, params ...any) error {
	val, err := c.memory.ReadAddressModes(false, params...)
	if err != nil {
		return err
	}
	c.A = val
	c.X = c.A
	c.setZN(c.A)
	return nil
}

// lxa implements LXA (0xAB): A = X = (A | magic) & imm.
// This is a highly unstable instruction; the magic constant varies by chip
// (observed values: 0x00, 0xEE, 0xFF). Using 0xFF makes the result deterministic.
func lxa(c *CPU, params ...any) error {
	val, err := c.memory.ReadAddressModes(true, params...)
	if err != nil {
		return err
	}
	// Magic constant 0xEE matches SingleStepTests/65x02 hardware captures.
	const magicConstant = 0xEE
	c.A = (c.A | magicConstant) & val
	c.X = c.A
	c.setZN(c.A)
	return nil
}

func nopUnofficial(c *CPU, params ...any) error {
	if len(params) > 0 {
		_, err := c.memory.ReadAddressModes(false, params...)
		return err
	}
	return nil
}

func rla(c *CPU, params ...any) error {
	if err := rol(c, params...); err != nil {
		return err
	}
	return and(c, params...)
}

func rra(c *CPU, params ...any) error {
	if err := ror(c, params...); err != nil {
		return err
	}
	return adc(c, params...)
}

func sax(c *CPU, params ...any) error {
	val := c.A & c.X
	return c.memory.WriteAddressModes(val, params...)
}

// sha implements SHA/AHX: stores A & X & (base_addr_hi + 1).
// Handles both AbsoluteY (0x9F) and IndirectY (0x93) addressing.
// On page cross the write address high byte is replaced by the stored value.
func sha(c *CPU, params ...any) error {
	baseAddr, indexReg := shBaseAddr(c, params)
	shWrite(c, c.A&c.X, baseAddr, indexReg)
	return nil
}

// shBaseAddr extracts the base address (before indexing) and index register value
// for SHA which supports both AbsoluteY and IndirectY addressing modes.
// For IndirectY the PC is re-read because params already carry the resolved address.
func shBaseAddr(c *CPU, params []any) (uint16, uint8) {
	if _, ok := params[0].(Absolute); ok {
		return uint16(params[0].(Absolute)), *params[1].(*uint8)
	}
	// IndirectY: params[0] is IndirectResolved (already has Y added), re-read base.
	zp := c.memory.Read(c.PC + 1)
	baseAddr := c.memory.ReadWordBug(uint16(zp))
	return baseAddr, c.Y
}

// shWrite performs the "SH" store with page-crossing address corruption.
// value is ANDed with (base_addr_hi + 1). On page cross the write address
// high byte is replaced by the stored value (hardware bus conflict behavior).
func shWrite(c *CPU, value uint8, baseAddr uint16, indexReg uint8) {
	andValue := value & (byte(baseAddr>>8) + 1)
	effectiveAddr := baseAddr + uint16(indexReg)
	pageCrossed := effectiveAddr&0xFF00 != baseAddr&0xFF00

	var writeAddr uint16
	if pageCrossed {
		writeAddr = (uint16(andValue) << 8) | (effectiveAddr & 0xFF)
	} else {
		writeAddr = effectiveAddr
	}
	c.memory.Write(writeAddr, andValue)
}

// shx implements SHX/SXA (0x9E): stores X & (base_addr_hi + 1).
func shx(c *CPU, params ...any) error {
	baseAddr := uint16(params[0].(Absolute))
	indexReg := *params[1].(*uint8) // Y
	shWrite(c, c.X, baseAddr, indexReg)
	return nil
}

// shy implements SHY/SYA (0x9C): stores Y & (base_addr_hi + 1).
func shy(c *CPU, params ...any) error {
	baseAddr := uint16(params[0].(Absolute))
	indexReg := *params[1].(*uint8) // X
	shWrite(c, c.Y, baseAddr, indexReg)
	return nil
}

func slo(c *CPU, params ...any) error {
	if err := asl(c, params...); err != nil {
		return err
	}
	return ora(c, params...)
}

func sre(c *CPU, params ...any) error {
	if err := lsr(c, params...); err != nil {
		return err
	}
	return eor(c, params...)
}

// tas implements TAS/XAS (0x9B): SP = A & X, then stores SP & (base_addr_hi + 1).
func tas(c *CPU, params ...any) error {
	c.SP = c.A & c.X
	baseAddr := uint16(params[0].(Absolute))
	indexReg := *params[1].(*uint8) // Y
	shWrite(c, c.SP, baseAddr, indexReg)
	return nil
}

// alr - AND with accumulator, then LSR.
func alr(c *CPU, params ...any) error {
	if err := and(c, params...); err != nil {
		return err
	}
	// LSR on accumulator
	c.Flags.C = c.A & 1
	c.A >>= 1
	c.setZN(c.A)
	return nil
}

// anc - AND with accumulator, copy N flag to C flag.
func anc(c *CPU, params ...any) error {
	if err := and(c, params...); err != nil {
		return err
	}
	// Copy N flag to C flag
	c.Flags.C = c.Flags.N
	return nil
}

// ane implements ANE/XAA (0x8B): A = (A | magic) & X & imm.
// This is a highly unstable instruction; the magic constant varies by chip
// (observed values: 0x00, 0xEE, 0xFF). Using 0xFF makes the result deterministic.
func ane(c *CPU, params ...any) error {
	val, err := c.memory.ReadAddressModes(true, params...)
	if err != nil {
		return err
	}
	// Magic constant 0xEE matches SingleStepTests/65x02 hardware captures.
	const magicConstant = 0xEE
	c.A = (c.A | magicConstant) & c.X & val
	c.setZN(c.A)
	return nil
}

// arr - AND with accumulator, then ROR with special flag behavior.
// In binary mode (D=0): C = bit6(result), V = bit6 XOR bit5.
// In decimal mode (D=1): BCD corrections are applied but conditioned on the AND
// result (pre-ROR), matching observed NMOS 6502 hardware behavior.
func arr(c *CPU, params ...any) error {
	if err := and(c, params...); err != nil {
		return err
	}
	// After and(), c.A holds the AND result. Save it for decimal correction conditions.
	andResult := c.A

	// ROR on accumulator using pre-AND carry as bit 7.
	oldCarry := c.Flags.C
	c.A = (andResult >> 1) | (oldCarry << 7)
	c.setZN(c.A)

	// ARR-specific V flag: XOR of bits 6 and 5 of the ROR result.
	c.Flags.V = (c.A>>6)&1 ^ (c.A>>5)&1

	if c.Flags.D != 0 && c.opts.variant != VariantNES6502 {
		// Decimal mode: nibble corrections conditioned on the AND result,
		// applied to the ROR result (N/V/Z already set from ROR result above).
		r := c.A
		if (andResult & 0x0F) >= 5 {
			r = (r & 0xF0) | ((r&0x0F + 6) & 0x0F)
		}
		if (andResult & 0xF0) >= 0x50 {
			r += 0x60
			c.Flags.C = 1
		} else {
			c.Flags.C = 0
		}
		c.A = r
	} else {
		// Binary mode (or NES which disables decimal): C = bit6 of ROR result.
		c.Flags.C = (c.A >> 6) & 1
	}
	return nil
}

// axs - (A AND X) minus immediate, store in X.
func axs(c *CPU, params ...any) error {
	value, err := c.memory.ReadAddressModes(true, params...)
	if err != nil {
		return err
	}

	// Calculate (A AND X) - immediate
	val := c.A & c.X
	result := int(val) - int(value)

	// Set carry if no borrow (result >= 0)
	if result >= 0 {
		c.Flags.C = 1
	} else {
		c.Flags.C = 0
	}

	// Store result in X
	c.X = uint8(result)
	c.setZN(c.X)
	return nil
}
