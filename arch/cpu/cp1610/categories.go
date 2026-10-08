package cp1610

import "github.com/retroenv/retrogolib/set"

// BranchingInstructions contains instructions that can change program flow.
var BranchingInstructions = set.NewFromSlice([]string{
	AdcrName,
	AddName,
	AddrName,
	AndName,
	AndrName,
	BName,
	BcName,
	BeqName,
	BescName,
	BextName,
	BgeName,
	BgtName,
	BleName,
	BltName,
	BmiName,
	BncName,
	BneqName,
	BnovName,
	BovName,
	BplName,
	BuscName,
	ComrName,
	DecrName,
	IncrName,
	JumpName,
	MovrName,
	MviName,
	NegrName,
	SubName,
	SubrName,
	XorName,
	XorrName,
})

// NotExecutingFollowingOpcodeInstructions contains instructions that do not
// continue at the next word.
var NotExecutingFollowingOpcodeInstructions = set.NewFromSlice([]string{
	BName,
	HltName,
	JumpName,
})

// MemoryReadInstructions contains instructions that read a data operand.
var MemoryReadInstructions = set.NewFromSlice([]string{
	AddName,
	AndName,
	CmpName,
	MviName,
	SubName,
	XorName,
})

// MemoryWriteInstructions contains instructions that write a data operand.
var MemoryWriteInstructions = set.NewFromSlice([]string{
	MvoName,
})

// MemoryReadWriteInstructions contains instructions that read and write one
// data operand. The CP1610 has no such instruction.
var MemoryReadWriteInstructions = set.New[string]()
