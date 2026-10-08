package cpu68000

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/set"
)

func TestMemoryCategoriesMatchHandlers(t *testing.T) {
	// The sets omitted immediate, decimal, shift, stack, and read-before-write
	// instructions that the handlers transfer through memory.
	readOnly := []string{BTSTName, CHKName, CMPName, CMPAName, CMPIName, CMPMName, DIVSName, DIVUName,
		MULSName, MULUName, MOVEAName, RTEName, RTRName, RTSName, TSTName, UNLKName}
	writeOnly := []string{BSRName, JSRName, LINKName, PEAName}
	readWrite := []string{ABCDName, ADDIName, ASLName, CLRName, NBCDName, ROXRName, SBCDName, SccName,
		SUBQName, TASName}
	neither := []string{BRAName, EXGName, JMPName, LEAName, MOVEQName, NOPName, SWAPName, TRAPName}

	for _, name := range readOnly {
		assert.True(t, MemoryReadInstructions.Contains(name), name)
		assert.False(t, MemoryWriteInstructions.Contains(name), name)
	}
	for _, name := range writeOnly {
		assert.False(t, MemoryReadInstructions.Contains(name), name)
		assert.True(t, MemoryWriteInstructions.Contains(name), name)
	}
	for _, name := range readWrite {
		assert.True(t, MemoryReadWriteInstructions.Contains(name), name)
	}
	for _, name := range neither {
		assert.False(t, MemoryReadInstructions.Contains(name), name)
		assert.False(t, MemoryWriteInstructions.Contains(name), name)
	}

	// Every read-modify-write instruction both reads and writes.
	for name := range MemoryReadWriteInstructions {
		assert.True(t, MemoryReadInstructions.Contains(name), name)
		assert.True(t, MemoryWriteInstructions.Contains(name), name)
	}

	// Every name is a registered instruction.
	for _, group := range []set.Set[string]{
		MemoryReadInstructions, MemoryWriteInstructions, MemoryReadWriteInstructions,
	} {
		for name := range group {
			_, ok := Instructions[name]
			assert.True(t, ok, name)
		}
	}
}
