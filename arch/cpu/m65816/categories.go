package m65816

import "github.com/retroenv/retrogolib/set"

// BranchingInstructions contains all instructions that change the program counter.
var BranchingInstructions = set.NewFromSlice([]string{
	BccName,
	BcsName,
	BeqName,
	BmiName,
	BneName,
	BplName,
	BraName,
	BrlName,
	BvcName,
	BvsName,
	JmlName,
	JmpName,
	JslName,
	JsrName,
})

// NotExecutingFollowingOpcodeInstructions contains all instructions that do not
// continue to execute the following opcode after execution.
var NotExecutingFollowingOpcodeInstructions = set.NewFromSlice([]string{
	BrkName,
	CopName,
	JmlName,
	JmpName,
	RtiName,
	RtlName,
	RtsName,
	StpName,
})

// MemoryReadInstructions contains instructions that read from an absolute memory address.
var MemoryReadInstructions = set.NewFromSlice([]string{
	AdcName,
	AndName,
	BitName,
	CmpName,
	CpxName,
	CpyName,
	EorName,
	JmlName,
	JmpName,
	LdaName,
	LdxName,
	LdyName,
	OraName,
	SbcName,
	TrbName,
	TsbName,
})

// MemoryWriteInstructions contains instructions that write to an absolute memory address.
var MemoryWriteInstructions = set.NewFromSlice([]string{
	StaName,
	StxName,
	StyName,
	StzName,
})

// MemoryReadWriteInstructions contains instructions that both read and write
// to an absolute memory address during execution.
var MemoryReadWriteInstructions = set.NewFromSlice([]string{
	AslName,
	DecName,
	IncName,
	LsrName,
	RolName,
	RorName,
	TrbName,
	TsbName,
})
