package m6502

import (
	"fmt"
	"math"

	. "github.com/retroenv/retrogolib/addressing"
)

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

func and(c *CPU, params ...any) {
	value := c.memory.ReadAddressModes(true, params...)
	c.A &= value
	c.setZN(c.A)
}

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

func bcc(c *CPU, params ...any) {
	c.branch(c.Flags.C == 0, params[0])
}

func bcs(c *CPU, params ...any) {
	c.branch(c.Flags.C != 0, params[0])
}

func beq(c *CPU, params ...any) {
	c.branch(c.Flags.Z != 0, params[0])
}

func bit(c *CPU, params ...any) {
	value := c.memory.ReadAbsolute(params[0], nil)
	c.setV((value>>6)&1 == 1)
	c.setZ(value & c.A)
	c.setN(value)
}

func bmi(c *CPU, params ...any) {
	c.branch(c.Flags.N != 0, params[0])
}

func bne(c *CPU, params ...any) {
	c.branch(c.Flags.Z == 0, params[0])
}

func bpl(c *CPU, params ...any) {
	c.branch(c.Flags.N == 0, params[0])
}

// Brk - Force Interrupt.
func brk(c *CPU) {
	c.irq()
}

func bvc(c *CPU, params ...any) {
	c.branch(c.Flags.V == 0, params[0])
}

func bvs(c *CPU, params ...any) {
	c.branch(c.Flags.V != 0, params[0])
}

func clc(c *CPU) {
	c.Flags.C = 0
}

func cld(c *CPU) {
	c.Flags.D = 0
}

func cli(c *CPU) {
	c.Flags.I = 0
}

func clv(c *CPU) {
	c.Flags.V = 0
}

func cmp(c *CPU, params ...any) {
	val := c.memory.ReadAddressModes(true, params...)
	c.compare(c.A, val)
}

func cpx(c *CPU, params ...any) {
	val := c.memory.ReadAddressModes(true, params[0])
	c.compare(c.X, val)
}

func cpy(c *CPU, params ...any) {
	val := c.memory.ReadAddressModes(true, params[0])
	c.compare(c.Y, val)
}

func dec(c *CPU, params ...any) {
	val := c.memory.ReadAddressModes(false, params...)
	val--
	c.memory.WriteAddressModes(val, params...)
	c.setZN(val)
}

func dex(c *CPU) {
	c.X--
	c.setZN(c.X)
}

func dey(c *CPU) {
	c.Y--
	c.setZN(c.Y)
}

func eor(c *CPU, params ...any) {
	value := c.memory.ReadAddressModes(true, params...)
	c.A ^= value
	c.setZN(c.A)
}

func inc(c *CPU, params ...any) {
	val := c.memory.ReadAddressModes(false, params...)
	val++
	c.memory.WriteAddressModes(val, params...)
	c.setZN(val)
}

func inx(c *CPU) {
	c.X++
	c.setZN(c.X)
}

func iny(c *CPU) {
	c.Y++
	c.setZN(c.Y)
}

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

func jsr(c *CPU, params ...any) {
	c.push16(c.PC + 2)

	addr := params[0].(Absolute)
	c.PC = uint16(addr)
}

func lda(c *CPU, params ...any) {
	c.A = c.memory.ReadAddressModes(true, params...)
	c.setZN(c.A)
}

func ldx(c *CPU, params ...any) {
	c.X = c.memory.ReadAddressModes(true, params...)
	c.setZN(c.X)
}

func ldy(c *CPU, params ...any) {
	c.Y = c.memory.ReadAddressModes(true, params...)
	c.setZN(c.Y)
}

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

func nop(_ *CPU) {
}

func ora(c *CPU, params ...any) {
	value := c.memory.ReadAddressModes(true, params...)
	c.A |= value
	c.setZN(c.A)
}

func pha(c *CPU) {
	c.push(c.A)
}

func php(c *CPU) {
	f := c.getFlags()
	f |= 0b0001_0000 // break is set to 1
	c.push(f)
}

func pla(c *CPU) {
	c.A = c.pop()
	c.setZN(c.A)
}

func plp(c *CPU) {
	f := c.pop()
	f &= 0b1110_1111 // break flag is ignored
	f |= 0b0010_0000 // unused flag is set
	c.setFlags(f)
}

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

func rts(c *CPU) {
	c.PC = c.pop16() + 1
}

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

func sec(c *CPU) {
	c.Flags.C = 1
}

func sed(c *CPU) {
	c.Flags.D = 1
}

func sei(c *CPU) {
	c.Flags.I = 1
}

func sta(c *CPU, params ...any) {
	c.memory.WriteAddressModes(c.A, params...)
}

func stx(c *CPU, params ...any) {
	c.memory.WriteAddressModes(c.X, params...)
}

func sty(c *CPU, params ...any) {
	c.memory.WriteAddressModes(c.Y, params...)
}

func tax(c *CPU) {
	c.X = c.A
	c.setZN(c.X)
}

func tay(c *CPU) {
	c.Y = c.A
	c.setZN(c.Y)
}

func tsx(c *CPU) {
	c.X = c.SP
	c.setZN(c.X)
}

func txa(c *CPU) {
	c.A = c.X
	c.setZN(c.A)
}

func txs(c *CPU) {
	c.SP = c.X
}

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
