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

func TestErrorConditions(t *testing.T) {
	c := New()

	// Test stack underflow
	c.SP = 0
	err := ret(c, 0)
	assert.ErrorContains(t, err, "stack underflow")
	assert.ErrorIs(t, err, ErrStackUnderflow, "error should be ErrStackUnderflow")

	// Test stack overflow
	c.SP = 16
	err = call(c, 0x200)
	assert.ErrorContains(t, err, "stack overflow")
	assert.ErrorIs(t, err, ErrStackOverflow, "error should be ErrStackOverflow")

	// Test key index out of bounds
	c.V[0] = 16
	err = skp(c, 0)
	assert.ErrorContains(t, err, "key index out of bounds")
	assert.ErrorIs(t, err, ErrKeyIndexOutOfBounds, "error should be ErrKeyIndexOutOfBounds")

	// Test font index out of bounds
	c.V[0] = 16
	err = c.ldFVx(0)
	assert.ErrorContains(t, err, "font index out of bounds")
	assert.ErrorIs(t, err, ErrFontIndexOutOfBounds, "error should be ErrFontIndexOutOfBounds")
}

func TestMemoryBounds(t *testing.T) {
	c := New()

	// Test memory out of bounds in Step
	c.PC = 4095
	err := c.Step()
	assert.ErrorContains(t, err, "memory")
	assert.ErrorIs(t, err, ErrMemoryOutOfBounds, "error should be ErrMemoryOutOfBounds")

	// Test memory bounds in BCD operation
	c.I = 4094
	err = c.ldBVx(0)
	assert.ErrorContains(t, err, "memory")
	assert.ErrorIs(t, err, ErrMemoryOutOfBounds, "error should be ErrMemoryOutOfBounds")
}

func TestCPUState(t *testing.T) {
	t.Parallel()
	c := New()

	// Test initial state
	assert.Equal(t, uint16(0x200), c.PC)
	assert.Equal(t, uint16(0), c.I)
	assert.Equal(t, uint8(0), c.SP)
	assert.Equal(t, uint8(0), c.DelayTimer)
	assert.Equal(t, uint8(0), c.SoundTimer)

	// Test register initialization
	for i := range 16 {
		assert.Equal(t, uint8(0), c.V[i], "V[%d] should be 0", i)
	}

	// Test stack initialization
	for i := range 16 {
		assert.Equal(t, uint16(0), c.Stack[i], "Stack[%d] should be 0", i)
	}
}

func TestAddOverflow(t *testing.T) {
	t.Parallel()
	c := New()

	// Test addition with carry
	c.V[0] = 0xFF
	c.V[1] = 0x02
	assert.NoError(t, add(c, 0x8010))
	assert.Equal(t, uint8(0x01), c.V[0])
	assert.Equal(t, uint8(0x01), c.V[0xF], "Carry flag should be set")

	// Test addition without carry
	c.V[0] = 0x10
	c.V[1] = 0x20
	assert.NoError(t, add(c, 0x8010))
	assert.Equal(t, uint8(0x30), c.V[0])
	assert.Equal(t, uint8(0x00), c.V[0xF], "Carry flag should be cleared")
}

func TestSubBorrow(t *testing.T) {
	t.Parallel()
	c := New()

	// Test subtraction without borrow
	c.V[0] = 0x30
	c.V[1] = 0x10
	assert.NoError(t, sub(c, 0x8015))
	assert.Equal(t, uint8(0x20), c.V[0])
	assert.Equal(t, uint8(0x01), c.V[0xF], "No borrow flag should be set")

	// Test subtraction with borrow
	c.V[0] = 0x10
	c.V[1] = 0x30
	assert.NoError(t, sub(c, 0x8015))
	assert.Equal(t, uint8(0xE0), c.V[0])
	assert.Equal(t, uint8(0x00), c.V[0xF], "Borrow flag should be cleared")
}

func TestRandomization(t *testing.T) {
	t.Parallel()
	c := New()

	// Run random test multiple times to ensure it's actually randomizing
	var different bool
	firstResult := uint8(0)
	for i := range 100 {
		assert.NoError(t, rnd(c, 0x00FF))
		if i == 0 {
			firstResult = c.V[0]
		} else if c.V[0] != firstResult {
			different = true
			break
		}
	}
	assert.True(t, different, "RND should produce different results")
}
