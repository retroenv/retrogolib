package chip8

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestCls(t *testing.T) {
	c := New()
	c.Display[0] = true
	c.Display[displayWidth+1] = true
	assert.NoError(t, cls(c, 0))

	for i := 0; i < 64; i++ {
		assert.False(t, c.Display[i], "Display[%d] is not false", i)
	}
}
