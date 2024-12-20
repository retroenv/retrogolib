package chip8

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestCls(t *testing.T) {
	c := New()
	c.Display[0] = 1
	c.Display[displayWidth+1] = 1
	assert.NoError(t, cls(c, 0))

	for i := range 64 {
		assert.Equal(t, 0, c.Display[i], "Display[%d] is not 0", i)
	}
}

func TestRet(t *testing.T) {
	c := New()
	c.Stack[0] = 0x200
	c.SP = 1
	assert.NoError(t, ret(c, 0))
	assert.Equal(t, uint16(0x200), c.PC)
	assert.Equal(t, uint8(0), c.SP)
}

func TestJp(t *testing.T) {
	c := New()
	assert.NoError(t, jp(c, 0x1123))
	assert.Equal(t, uint16(0x123), c.PC)
}

func TestCall(t *testing.T) {
	c := New()
	c.PC = 0x200
	assert.NoError(t, call(c, 0x123))
	assert.Equal(t, uint16(0x123), c.PC)
	assert.Equal(t, uint16(0x200), c.Stack[0])
	assert.Equal(t, uint8(1), c.SP)
}

func TestSe(t *testing.T) {
	c := New()
	c.V[0] = 0x12
	assert.NoError(t, se(c, 0x3000))
	assert.Equal(t, uint16(0x202), c.PC)
}

func TestSne(t *testing.T) {
	c := New()
	c.V[0] = 0x12
	assert.NoError(t, sne(c, 0x4012))
	assert.Equal(t, uint16(0x202), c.PC)
}

func TestOr(t *testing.T) {
	c := New()
	c.V[0] = 0x12
	c.V[1] = 0x34
	assert.NoError(t, or(c, 0x0010))
	assert.Equal(t, uint8(0x36), c.V[0])
}

func TestXor(t *testing.T) {
	c := New()
	c.V[0] = 0x12
	c.V[1] = 0x34
	assert.NoError(t, xor(c, 0x0010))
	assert.Equal(t, uint8(0x26), c.V[0])
}

func TestAdd(t *testing.T) {
	c := New()

	c.V[0] = 0x12
	assert.NoError(t, add(c, 0x7034))
	assert.Equal(t, uint8(0x46), c.V[0])

	c.V[0] = 0x12
	c.V[1] = 0x34
	assert.NoError(t, add(c, 0x8010))
	assert.Equal(t, uint8(0x46), c.V[0])
}

func TestSub(t *testing.T) {
	c := New()
	c.V[0] = 0x34
	c.V[1] = 0x12
	assert.NoError(t, sub(c, 0x0010))
	assert.Equal(t, uint8(0x22), c.V[0])
}

func TestLd(t *testing.T) {
	c := New()
	assert.NoError(t, ld(c, 0x6012))
	assert.Equal(t, uint8(0x12), c.V[0])
}

func TestAnd(t *testing.T) {
	c := New()
	c.V[0] = 0x12
	c.V[1] = 0x34
	assert.NoError(t, and(c, 0x0010))
	assert.Equal(t, uint8(0x10), c.V[0])
}

func TestDrw(t *testing.T) {
	c := New()
	c.Memory[0] = 0b11110000
	c.Memory[1] = 0b00001111
	c.Memory[2] = 0b11110000
	c.Memory[3] = 0b00001111
	c.Display[0] = 1
	c.Display[1] = 1
	c.Display[displayWidth] = 1
	c.Display[displayWidth+1] = 1
	assert.NoError(t, drw(c, 0x0003))

	assert.Equal(t, uint8(1), c.V[0xF])
	assert.Equal(t, 0, c.Display[0])
	assert.Equal(t, 0, c.Display[1])
	assert.Equal(t, 1, c.Display[displayWidth])
	assert.Equal(t, 1, c.Display[displayWidth+1])
}

func TestRnd(t *testing.T) {
	c := New()
	assert.NoError(t, rnd(c, 0x00ff))
	assert.NotEqual(t, uint8(0), c.V[0])
}

func TestShl(t *testing.T) {
	c := New()
	c.V[0] = 0b10000000

	assert.NoError(t, shl(c, 0))
	assert.Equal(t, uint8(1), c.V[0xF])
	assert.Equal(t, uint8(0), c.V[0])

	assert.NoError(t, shl(c, 0))
	assert.Equal(t, uint8(0), c.V[0xF])
	assert.Equal(t, uint8(0), c.V[0])
}

func TestShr(t *testing.T) {
	c := New()
	c.V[0] = 0b00000001

	assert.NoError(t, shr(c, 0))
	assert.Equal(t, uint8(1), c.V[0xF])
	assert.Equal(t, uint8(0), c.V[0])

	assert.NoError(t, shr(c, 0))
	assert.Equal(t, uint8(0), c.V[0xF])
	assert.Equal(t, uint8(0), c.V[0])
}

func TestSkp(t *testing.T) {
	c := New()

	c.Key[0] = true
	assert.NoError(t, skp(c, 0))
	assert.Equal(t, uint16(0x204), c.PC)

	c.Key[0] = false
	assert.NoError(t, skp(c, 0))
	assert.Equal(t, uint16(0x206), c.PC)
}

func TestSknp(t *testing.T) {
	c := New()

	c.Key[0] = false
	assert.NoError(t, sknp(c, 0))
	assert.Equal(t, uint16(0x204), c.PC)

	c.Key[0] = true
	assert.NoError(t, sknp(c, 0))
	assert.Equal(t, uint16(0x206), c.PC)
}

func TestSubn(t *testing.T) {
	c := New()
	c.V[0] = 0x12
	c.V[1] = 0x34
	assert.NoError(t, subn(c, 0x0010))
	assert.Equal(t, uint8(0x22), c.V[0])
}
