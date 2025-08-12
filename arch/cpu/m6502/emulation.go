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

	c.push16(c.PC + 2)
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
