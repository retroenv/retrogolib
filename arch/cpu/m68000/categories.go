package m68000

import "github.com/retroenv/retrogolib/set"

// BranchingInstructions contains all branching and jumping instructions.
var BranchingInstructions = set.NewFromSlice([]string{
	BccName,
	BRAName,
	BSRName,
	DBccName,
	JMPName,
	JSRName,
	RTEName,
	RTRName,
	RTSName,
})

// NotExecutingFollowingOpcodeInstructions contains all instructions that jump
// to a different address and do not return to execute the following opcode.
var NotExecutingFollowingOpcodeInstructions = set.NewFromSlice([]string{
	BRAName,
	JMPName,
	RTEName,
	RTSName,
	STOPName,
})

// MemoryReadInstructions contains all instructions that can read from memory.
var MemoryReadInstructions = set.NewFromSlice([]string{
	ADDName,
	ANDName,
	CMPName,
	CMPAName,
	CMPIName,
	CMPMName,
	DIVSName,
	DIVUName,
	EORName,
	MOVEName,
	MOVEAName,
	MOVEMName,
	MOVEPName,
	MULSName,
	MULUName,
	ORName,
	SUBName,
	TSTName,
})

// MemoryWriteInstructions contains all instructions that can write to memory.
var MemoryWriteInstructions = set.NewFromSlice([]string{
	ADDName,
	ANDName,
	BCHGName,
	BCLRName,
	BSETName,
	CLRName,
	EORName,
	MOVEName,
	MOVEMName,
	MOVEPName,
	NEGName,
	NEGXName,
	NOTName,
	ORName,
	SUBName,
})

// MemoryReadWriteInstructions contains all instructions that both read and write
// to the same memory location during execution.
var MemoryReadWriteInstructions = set.NewFromSlice([]string{
	ADDName,
	ANDName,
	BCHGName,
	BCLRName,
	BSETName,
	EORName,
	NEGName,
	NEGXName,
	NOTName,
	ORName,
	SUBName,
})
