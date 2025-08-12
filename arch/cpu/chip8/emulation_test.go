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
	c.Stack[0] = 0x202 // Stack contains return address (after CALL instruction)
	c.SP = 1
	assert.NoError(t, ret(c, 0))
	assert.Equal(t, uint16(0x202), c.PC) // PC restored to return address
	assert.Equal(t, uint8(0), c.SP)
}

func TestCallReturnInteraction(t *testing.T) {
	c := New()

	// Simulate CALL from address 0x200
	c.PC = 0x200
	assert.NoError(t, call(c, 0x300)) // CALL 0x300
	assert.Equal(t, uint16(0x300), c.PC)
	assert.Equal(t, uint16(0x202), c.Stack[0]) // Return address saved
	assert.Equal(t, uint8(1), c.SP)

	// Simulate RET
	assert.NoError(t, ret(c, 0))
	assert.Equal(t, uint16(0x202), c.PC) // Back to instruction after CALL
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
	assert.Equal(t, uint16(0x202), c.Stack[0]) // Should save return address (PC+2)
	assert.Equal(t, uint8(1), c.SP)
}

func TestSe(t *testing.T) {
	c := New()
	c.V[0] = 0x12
	assert.NoError(t, se(c, 0x3012))     // SE V0, 0x12 - should skip because V0 == 0x12
	assert.Equal(t, uint16(0x204), c.PC) // Should skip next instruction
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

func TestLdVxK(t *testing.T) {
	c := New()

	// Test waiting for key press (no key pressed)
	err := c.ldVxK(0)
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x200), c.PC) // PC should not advance

	// Test key press detected
	c.Key[5] = true
	err = c.ldVxK(0)
	assert.NoError(t, err)
	assert.Equal(t, uint8(5), c.V[0])
	assert.Equal(t, uint16(0x202), c.PC) // PC should advance
}

func TestJpV0(t *testing.T) {
	c := New()
	c.V[0] = 0x10
	assert.NoError(t, jp(c, 0xB100))     // JP V0, 0x100
	assert.Equal(t, uint16(0x110), c.PC) // PC should be 0x100 + V0 (0x10)
}

func TestSeRegisterRegister(t *testing.T) {
	c := New()
	c.V[0] = 0x12
	c.V[1] = 0x12
	assert.NoError(t, se(c, 0x5010))     // SE V0, V1
	assert.Equal(t, uint16(0x204), c.PC) // Should skip next instruction

	c.PC = 0x200
	c.V[1] = 0x13
	assert.NoError(t, se(c, 0x5010))     // SE V0, V1
	assert.Equal(t, uint16(0x202), c.PC) // Should not skip
}

func TestSneRegisterRegister(t *testing.T) {
	c := New()
	c.V[0] = 0x12
	c.V[1] = 0x13
	assert.NoError(t, sne(c, 0x9010))    // SNE V0, V1
	assert.Equal(t, uint16(0x204), c.PC) // Should skip next instruction

	c.PC = 0x200
	c.V[1] = 0x12
	assert.NoError(t, sne(c, 0x9010))    // SNE V0, V1
	assert.Equal(t, uint16(0x202), c.PC) // Should not skip
}

func TestLdDT(t *testing.T) {
	c := New()
	c.DelayTimer = 0x42
	assert.NoError(t, ldF(c, 0xF007)) // LD V0, DT
	assert.Equal(t, uint8(0x42), c.V[0])
	assert.Equal(t, uint16(0x202), c.PC)
}

func TestLdST(t *testing.T) {
	c := New()
	c.V[0] = 0x42
	assert.NoError(t, ldF(c, 0xF018)) // LD ST, V0
	assert.Equal(t, uint8(0x42), c.SoundTimer)
	assert.Equal(t, uint16(0x202), c.PC)
}

func TestLdBVxCorrectness(t *testing.T) {
	c := New()
	c.V[0] = 123
	c.I = 0x300
	assert.NoError(t, c.ldBVx(0))
	assert.Equal(t, uint8(1), c.Memory[0x300]) // Hundreds
	assert.Equal(t, uint8(2), c.Memory[0x301]) // Tens
	assert.Equal(t, uint8(3), c.Memory[0x302]) // Ones
}

func TestLdIVxCorrectness(t *testing.T) {
	c := New()
	c.I = 0x300
	for i := range uint16(4) {
		c.V[i] = uint8(i + 1)
	}
	assert.NoError(t, c.ldIVx(3)) // Store V0-V3
	for i := range uint16(4) {
		assert.Equal(t, uint8(i+1), c.Memory[0x300+i])
	}
}

func TestLdVxICorrectness(t *testing.T) {
	c := New()
	c.I = 0x300
	for i := range uint16(4) {
		c.Memory[0x300+i] = uint8(i + 10)
	}
	assert.NoError(t, c.ldVxI(3)) // Load V0-V3
	for i := range uint16(4) {
		assert.Equal(t, uint8(i+10), c.V[i])
	}
}

func TestAddI(t *testing.T) {
	c := New()
	c.I = 0x100
	c.V[0] = 0x50
	assert.NoError(t, add(c, 0xF01E)) // ADD I, V0
	assert.Equal(t, uint16(0x150), c.I)
	assert.Equal(t, uint16(0x202), c.PC)
}

func TestUpdateTimers(t *testing.T) {
	c := New()
	c.DelayTimer = 10
	c.SoundTimer = 5

	c.UpdateTimers()
	assert.Equal(t, uint8(9), c.DelayTimer)
	assert.Equal(t, uint8(4), c.SoundTimer)

	// Test that timers don't go below 0
	for range 10 {
		c.UpdateTimers()
	}
	assert.Equal(t, uint8(0), c.DelayTimer)
	assert.Equal(t, uint8(0), c.SoundTimer)
}

func TestReset(t *testing.T) {
	c := New()

	// Modify state
	c.PC = 0x300
	c.SP = 5
	c.I = 0x123
	c.V[0] = 0x42
	c.Stack[0] = 0x200
	c.Display[0] = 1
	c.Key[0] = true
	c.DelayTimer = 10
	c.SoundTimer = 5
	c.RedrawScreen = true

	c.Reset()

	// Verify reset state
	assert.Equal(t, uint16(0x200), c.PC)
	assert.Equal(t, uint8(0), c.SP)
	assert.Equal(t, uint16(0), c.I)
	assert.Equal(t, uint8(0), c.V[0])
	assert.Equal(t, uint16(0), c.Stack[0])
	assert.Equal(t, uint8(0), c.Display[0])
	assert.False(t, c.Key[0])
	assert.Equal(t, uint8(0), c.DelayTimer)
	assert.Equal(t, uint8(0), c.SoundTimer)
	assert.False(t, c.RedrawScreen)

	// Verify font data is preserved
	assert.Equal(t, fontSet[0], c.Memory[0])
}

func TestGetSetState(t *testing.T) {
	c := New()

	// Modify state
	c.PC = 0x300
	c.I = 0x123
	c.V[0] = 0x42
	c.DelayTimer = 10

	// Get state
	state := c.GetState()
	assert.Equal(t, uint16(0x300), state.PC)
	assert.Equal(t, uint16(0x123), state.I)
	assert.Equal(t, uint8(0x42), state.V[0])
	assert.Equal(t, uint8(10), state.DelayTimer)

	// Create new CPU and set state
	c2 := New()
	c2.SetState(state)
	assert.Equal(t, uint16(0x300), c2.PC)
	assert.Equal(t, uint16(0x123), c2.I)
	assert.Equal(t, uint8(0x42), c2.V[0])
	assert.Equal(t, uint8(10), c2.DelayTimer)
}

func TestDrwWrapping(t *testing.T) {
	c := New()

	// Test X wrapping
	c.V[0] = 62 // Near right edge
	c.V[1] = 0
	c.I = 0x200
	c.Memory[0x200] = 0xFF // 8 pixels wide sprite

	assert.NoError(t, drw(c, 0xD011)) // DRW V0, V1, 1

	// Should draw 2 pixels on right edge, rest should wrap but are cut off
	assert.Equal(t, uint8(1), c.Display[62])
	assert.Equal(t, uint8(1), c.Display[63])

	// Test Y wrapping
	c = New()
	c.V[0] = 0
	c.V[1] = 31 // Bottom row
	c.I = 0x200
	c.Memory[0x200] = 0xFF
	c.Memory[0x201] = 0xFF // This would wrap but should be cut off

	assert.NoError(t, drw(c, 0xD012)) // DRW V0, V1, 2

	// Only first row should be drawn
	assert.Equal(t, uint8(1), c.Display[31*64])
}
