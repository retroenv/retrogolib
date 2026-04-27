package m6502

import "github.com/retroenv/retrogolib/set"

// BranchingInstructions contains all branching instructions.
var BranchingInstructions = set.NewFromSlice([]string{
	BccInst.Name,
	BcsInst.Name,
	BeqInst.Name,
	BmiInst.Name,
	BneInst.Name,
	BplInst.Name,
	BraInst.Name,
	BvcInst.Name,
	BvsInst.Name,
	JmpInst.Name,
	JsrInst.Name,
})

// NotExecutingFollowingOpcodeInstructions contains all instructions that jump
// to a different address and do not return to execute the following opcode.
var NotExecutingFollowingOpcodeInstructions = set.NewFromSlice([]string{
	BrkInst.Name, // BRK jumps to IRQ handler, doesn't continue to next instruction
	JmpInst.Name,
	RtiInst.Name,
	RtsInst.Name,
})

// MemoryReadInstructions contains all instructions that can read from an
// absolute memory address.
var MemoryReadInstructions = set.NewFromSlice([]string{
	AndInst.Name,
	BitInst.Name,
	CmpInst.Name,
	CpxInst.Name,
	CpyInst.Name,
	JmpInst.Name,
	LdaInst.Name,
	LdxInst.Name,
	LdyInst.Name,
	LaxInst.Name,
})

// MemoryWriteInstructions contains all instructions that can write to an
// absolute memory address.
var MemoryWriteInstructions = set.NewFromSlice([]string{
	SaxInst.Name,
	StaInst.Name,
	StxInst.Name,
	StyInst.Name,
	StzInst.Name,
})

// MemoryReadWriteInstructions contains all instructions that can read and write
// during instruction execution an absolute memory address.
var MemoryReadWriteInstructions = set.NewFromSlice([]string{
	AdcInst.Name,
	AslInst.Name,
	DcpInst.Name,
	DecInst.Name,
	EorInst.Name,
	IncInst.Name,
	IscInst.Name,
	LsrInst.Name,
	OraInst.Name,
	RlaInst.Name,
	RolInst.Name,
	RorInst.Name,
	RraInst.Name,
	SbcInst.Name,
	SloInst.Name,
	SreInst.Name,
	TrbInst.Name,
	TsbInst.Name,
})
