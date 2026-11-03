package cpu6502

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestMemoryInstructionCategories(t *testing.T) {
	tests := []struct {
		name      string
		read      bool
		readWrite bool
		write     bool
	}{
		// These instructions change a register or the flags, not the memory.
		{
			name: AdcName,
			read: true,
		},
		{
			name: EorName,
			read: true,
		},
		{
			name: OraName,
			read: true,
		},
		{
			name: SbcName,
			read: true,
		},
		{
			name: LdaName,
			read: true,
		},
		{
			name:  StaName,
			write: true,
		},
		{
			name:      AslName,
			readWrite: true,
		},
		{
			name:      IncName,
			readWrite: true,
		},
		{
			name:      RlaName,
			readWrite: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.read, MemoryReadInstructions.Contains(tt.name))
			assert.Equal(t, tt.readWrite, MemoryReadWriteInstructions.Contains(tt.name))
			assert.Equal(t, tt.write, MemoryWriteInstructions.Contains(tt.name))
		})
	}
}

func TestMemoryInstructionCategoriesDoNotOverlap(t *testing.T) {
	for name := range MemoryReadWriteInstructions {
		assert.False(t, MemoryReadInstructions.Contains(name), name)
		assert.False(t, MemoryWriteInstructions.Contains(name), name)
	}
	for name := range MemoryReadInstructions {
		assert.False(t, MemoryWriteInstructions.Contains(name), name)
	}
}
