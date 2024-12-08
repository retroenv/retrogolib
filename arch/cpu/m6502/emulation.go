package m6502

import (
	"fmt"
	"math"

	. "github.com/retroenv/retrogolib/addressing"
)

// adc - Add with Carry.
func adc(c *CPU, params ...any) {
	a := c.A
	value := c.memory.ReadAddressModes(true, params...)
	sum := int(c.A) + int(c.Flags.C) + int(value)
	c.A = uint8(sum)
	c.setZN(c.A)

	if sum > math.MaxUint8 {
		c.Flags.C = 1
	} else {
		c.Flags.C = 0
	}
	c.setV((a^value)&0x80 == 0 && (a^c.A)&0x80 != 0)
}

// and - AND with accumulator.
func and(c *CPU, params ...any) {
	value := c.memory.ReadAddressModes(true, params...)
	c.A &= value
	c.setZN(c.A)
}

// asl - Arithmetic Shift Left.
func asl(c *CPU, params ...any) {
	if hasAccumulatorParam(params...) {
		c.Flags.C = (c.A >> 7) & 1
		c.A <<= 1
		c.setZN(c.A)
		return
	}

	val := c.memory.ReadAddressModes(false, params...)
	c.Flags.C = (val >> 7) & 1
	val <<= 1
	c.setZN(val)
	c.memory.WriteAddressModes(val, params...)
}

// bcc - Branch if Carry Clear.
func bcc(c *CPU, params ...any) {
	c.branch(c.Flags.C == 0, params[0])
}

// bcs - Branch if Carry Set.
func bcs(c *CPU, params ...any) {
	c.branch(c.Flags.C != 0, params[0])
}

// beq - Branch if Equal.
func beq(c *CPU, params ...any) {
	c.branch(c.Flags.Z != 0, params[0])
}

// bit - Bit Test.
func bit(c *CPU, params ...any) {
	value := c.memory.ReadAbsolute(params[0], nil)
	c.setV((value>>6)&1 == 1)
	c.setZ(value & c.A)
	c.setN(value)
}

// bmi - Branch if Minus.
func bmi(c *CPU, params ...any) {
	c.branch(c.Flags.N != 0, params[0])
}

// bne - Branch if Not Equal.
func bne(c *CPU, params ...any) {
	c.branch(c.Flags.Z == 0, params[0])
}

// bpl - Branch if Positive.
func bpl(c *CPU, params ...any) {
	c.branch(c.Flags.N == 0, params[0])
}

// brk - Force Interrupt.
func brk(c *CPU) {
	c.irq()
}

// bvc - Branch if Overflow Clear.
func bvc(c *CPU, params ...any) {
	c.branch(c.Flags.V == 0, params[0])
}

// bvs - Branch if Overflow Set.
func bvs(c *CPU, params ...any) {
	c.branch(c.Flags.V != 0, params[0])
}

// clc - Clear Carry Flag.
func clc(c *CPU) {
	c.Flags.C = 0
}

// cld - Clear Decimal Mode.
func cld(c *CPU) {
	c.Flags.D = 0
}

// cli - Clear Interrupt Disable.
func cli(c *CPU) {
	c.Flags.I = 0
}

// clv - Clear Overflow Flag.
func clv(c *CPU) {
	c.Flags.V = 0
}

// cmp - Compare the contents of A.
func cmp(c *CPU, params ...any) {
	val := c.memory.ReadAddressModes(true, params...)
	c.compare(c.A, val)
}

// cpx - Compare the contents of X.
func cpx(c *CPU, params ...any) {
	val := c.memory.ReadAddressModes(true, params[0])
	c.compare(c.X, val)
}

// cpy - Compare the contents of Y.
func cpy(c *CPU, params ...any) {
	val := c.memory.ReadAddressModes(true, params[0])
	c.compare(c.Y, val)
}

// dec - Decrement memory.
func dec(c *CPU, params ...any) {
	val := c.memory.ReadAddressModes(false, params...)
	val--
	c.memory.WriteAddressModes(val, params...)
	c.setZN(val)
}

// dex - Decrement X Register.
func dex(c *CPU) {
	c.X--
	c.setZN(c.X)
}

// dey - Decrement Y Register.
func dey(c *CPU) {
	c.Y--
	c.setZN(c.Y)
}

// eor - Exclusive OR - XOR.
func eor(c *CPU, params ...any) {
	value := c.memory.ReadAddressModes(true, params...)
	c.A ^= value
	c.setZN(c.A)
}

// inc - Increments memory.
func inc(c *CPU, params ...any) {
	val := c.memory.ReadAddressModes(false, params...)
	val++
	c.memory.WriteAddressModes(val, params...)
	c.setZN(val)
}

// inx - Increment X Register.
func inx(c *CPU) {
	c.X++
	c.setZN(c.X)
}

// iny - Increment Y Register.
func iny(c *CPU) {
	c.Y++
	c.setZN(c.Y)
}

// jmp - jump to address.
func jmp(c *CPU, params ...any) {
	param := params[0]
	switch address := param.(type) {
	case Absolute:
		c.PC = uint16(address)
	case Indirect:
		c.PC = c.memory.ReadWordBug(uint16(address))

	default:
		panic(fmt.Sprintf("unsupported jmp mode type %T", param))
	}
}

// jsr - jump to subroutine.
func jsr(c *CPU, params ...any) {
	c.push16(c.PC + 2)

	addr := params[0].(Absolute)
	c.PC = uint16(addr)
}

// lda - Load Accumulator - load a byte into A.
func lda(c *CPU, params ...any) {
	c.A = c.memory.ReadAddressModes(true, params...)
	c.setZN(c.A)
}

// ldx - Load X Register - load a byte into X.
func ldx(c *CPU, params ...any) {
	c.X = c.memory.ReadAddressModes(true, params...)
	c.setZN(c.X)
}

// ldy - Load Y Register - load a byte into Y.
func ldy(c *CPU, params ...any) {
	c.Y = c.memory.ReadAddressModes(true, params...)
	c.setZN(c.Y)
}

// lsr - Logical Shift Right.
func lsr(c *CPU, params ...any) {
	if hasAccumulatorParam(params...) {
		c.Flags.C = c.A & 1
		c.A >>= 1
		c.setZN(c.A)
		return
	}

	val := c.memory.ReadAddressModes(false, params...)
	c.Flags.C = val & 1
	val >>= 1
	c.setZN(val)
	c.memory.WriteAddressModes(val, params...)
}

// nop - No Operation.
func nop(_ *CPU) {
}

// ora - OR with Accumulator.
func ora(c *CPU, params ...any) {
	value := c.memory.ReadAddressModes(true, params...)
	c.A |= value
	c.setZN(c.A)
}

// pha - Push Accumulator - push A content to stack.
func pha(c *CPU) {
	c.push(c.A)
}

// php - Push Processor Status - push status flags to stack.
func php(c *CPU) {
	f := c.GetFlags()
	f |= 0b0001_0000 // break is set to 1
	c.push(f)
}

// pla - Pull Accumulator - pull A content from stack.
func pla(c *CPU) {
	c.A = c.pop()
	c.setZN(c.A)
}

// plp - Pull Processor Status - pull status flags from stack.
func plp(c *CPU) {
	f := c.pop()
	f &= 0b1110_1111 // break flag is ignored
	f |= 0b0010_0000 // unused flag is set
	c.setFlags(f)
}

// rol - Rotate Left.
func rol(c *CPU, params ...any) {
	cFlag := c.Flags.C
	if hasAccumulatorParam(params...) {
		c.Flags.C = (c.A >> 7) & 1
		c.A = (c.A << 1) | cFlag
		c.setZN(c.A)
		return
	}

	val := c.memory.ReadAddressModes(false, params...)
	c.Flags.C = (val >> 7) & 1
	val = (val << 1) | cFlag
	c.setZN(val)
	c.memory.WriteAddressModes(val, params...)
}

// ror - Rotate Right.
func ror(c *CPU, params ...any) {
	cFlag := c.Flags.C
	if hasAccumulatorParam(params...) {
		c.Flags.C = c.A & 1
		c.A = (c.A >> 1) | (cFlag << 7)
		c.setZN(c.A)
		return
	}

	val := c.memory.ReadAddressModes(false, params...)
	c.Flags.C = val & 1
	val = (val >> 1) | (cFlag << 7)
	c.setZN(val)
	c.memory.WriteAddressModes(val, params...)
}

// rti - Return from Interrupt.
func rti(c *CPU) {
	b := c.pop()
	b &= 0b1110_1111 // break flag is ignored
	b |= 0b0010_0000 // unused flag is set
	c.setFlags(b)
	c.PC = c.pop16()

	// lock is already taken
	c.irqRunning = false
	c.nmiRunning = false
}

// rts - return from subroutine.
func rts(c *CPU) {
	c.PC = c.pop16() + 1
}

// sbc - subtract with Carry.
func sbc(c *CPU, params ...any) {
	a := c.A
	value := c.memory.ReadAddressModes(true, params...)
	sub := int(c.A) - int(value) - (1 - int(c.Flags.C))
	c.A = uint8(sub)
	c.setZN(c.A)

	if sub >= 0 {
		c.Flags.C = 1
	} else {
		c.Flags.C = 0
	}
	c.setV((a^value)&0x80 != 0 && (a^c.A)&0x80 != 0)
}

// sec - Set Carry Flag.
func sec(c *CPU) {
	c.Flags.C = 1
}

// sed - Set Decimal Flag.
func sed(c *CPU) {
	c.Flags.D = 1
}

// sei - Set Interrupt Disable.
func sei(c *CPU) {
	c.Flags.I = 1
}

// sta - Store Accumulator.
func sta(c *CPU, params ...any) {
	c.memory.WriteAddressModes(c.A, params...)
}

// stx - Store X Register.
func stx(c *CPU, params ...any) {
	c.memory.WriteAddressModes(c.X, params...)
}

// sty - Store Y Register.
func sty(c *CPU, params ...any) {
	c.memory.WriteAddressModes(c.Y, params...)
}

// tax - Transfer Accumulator to X.
func tax(c *CPU) {
	c.X = c.A
	c.setZN(c.X)
}

// tay - Transfer Accumulator to Y.
func tay(c *CPU) {
	c.Y = c.A
	c.setZN(c.Y)
}

// tsx - Transfer Stack Pointer to X.
func tsx(c *CPU) {
	c.X = c.SP
	c.setZN(c.X)
}

// txa - Transfer X to Accumulator.
func txa(c *CPU) {
	c.A = c.X
	c.setZN(c.A)
}

// txs - Transfer X to Stack Pointer.
func txs(c *CPU) {
	c.SP = c.X
}

// tya - Transfer Y to Accumulator.
func tya(c *CPU) {
	c.A = c.Y
	c.setZN(c.A)
}

// unofficial instructions

func dcp(c *CPU, params ...any) {
	dec(c, params...)
	val := c.memory.ReadAddressModes(false, params...)
	c.compare(c.A, val)
}

func isc(c *CPU, params ...any) {
	inc(c, params...)
	sbc(c, params...)
}

func lax(c *CPU, params ...any) {
	c.A = c.memory.ReadAddressModes(false, params...)
	c.X = c.A
	c.setZN(c.A)
}

func nopUnofficial(c *CPU, params ...any) {
	if len(params) > 0 {
		c.memory.ReadAddressModes(false, params...)
	}
}

func rla(c *CPU, params ...any) {
	rol(c, params...)
	and(c, params...)
}

func rra(c *CPU, params ...any) {
	ror(c, params...)
	adc(c, params...)
}

func sax(c *CPU, params ...any) {
	val := c.A & c.X
	c.memory.WriteAddressModes(val, params...)
}

func slo(c *CPU, params ...any) {
	asl(c, params...)
	ora(c, params...)
}

func sre(c *CPU, params ...any) {
	lsr(c, params...)
	eor(c, params...)
}
