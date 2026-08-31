package z80

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestInstructionRegistry(t *testing.T) {
	t.Parallel()

	for name := range NameToOpcodeID {
		assert.NotNil(t, Instructions[name], "instruction %s", name)
	}
	for name, ins := range Instructions {
		assert.Equal(t, name, ins.Name)
	}
}
