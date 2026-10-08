package chip8

import (
	"encoding/json"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestStateSerializesKeyWait(t *testing.T) {
	t.Parallel()

	c := New()
	c.keyWait = KeyWait{
		Active:   true,
		Pressed:  true,
		Key:      5,
		Register: 2,
	}
	c.drewThisFrame = true

	// External save states serialize the exported fields only.
	data, err := json.Marshal(c.State())
	assert.NoError(t, err)
	var state State
	assert.NoError(t, json.Unmarshal(data, &state))

	restored := New()
	restored.SetState(state)
	assert.Equal(t, c.keyWait, restored.keyWait)
	assert.True(t, restored.drewThisFrame)

	// The restored wait completes on the key release.
	restored.Key[5] = false
	assert.NoError(t, restored.ldVxK(0))
	assert.Equal(t, uint8(5), restored.V[2])
	assert.Equal(t, uint16(0x202), restored.PC)
	assert.Equal(t, KeyWait{}, restored.keyWait)
}

func TestLdVxKRejectsInvalidRestoredIndices(t *testing.T) {
	t.Parallel()

	c := New()
	c.keyWait = KeyWait{
		Active:   true,
		Register: 16,
	}
	assert.ErrorIs(t, c.ldVxK(0), ErrRegisterOutOfBounds)

	c = New()
	c.keyWait = KeyWait{
		Active:  true,
		Pressed: true,
		Key:     16,
	}
	assert.ErrorIs(t, c.ldVxK(0), ErrRegisterOutOfBounds)
}
