package cpu6809

import "github.com/retroenv/retrogolib/set"

// BranchingInstructions contains branch, jump, and subroutine-call instructions.
var BranchingInstructions = set.NewFromSlice([]string{
	BccInst.Name,
	BcsInst.Name,
	BeqInst.Name,
	BgeInst.Name,
	BgtInst.Name,
	BhiInst.Name,
	BleInst.Name,
	BlsInst.Name,
	BltInst.Name,
	BmiInst.Name,
	BneInst.Name,
	BplInst.Name,
	BraInst.Name,
	BsrInst.Name,
	BvcInst.Name,
	BvsInst.Name,
	JmpInst.Name,
	JsrInst.Name,
	LbccInst.Name,
	LbcsInst.Name,
	LbeqInst.Name,
	LbgeInst.Name,
	LbgtInst.Name,
	LbhiInst.Name,
	LbleInst.Name,
	LblsInst.Name,
	LbltInst.Name,
	LbmiInst.Name,
	LbneInst.Name,
	LbplInst.Name,
	LbraInst.Name,
	LbsrInst.Name,
	LbvcInst.Name,
	LbvsInst.Name,
})

// NotExecutingFollowingOpcodeInstructions contains all instructions that do not
// continue to execute the following opcode after execution.
var NotExecutingFollowingOpcodeInstructions = set.NewFromSlice([]string{
	BraInst.Name,
	JmpInst.Name,
	LbraInst.Name,
	RtiInst.Name,
	RtsInst.Name,
	SwiInst.Name,
	Swi2Inst.Name,
	Swi3Inst.Name,
})

// MemoryReadInstructions contains instructions that read from an absolute memory address.
var MemoryReadInstructions = set.NewFromSlice([]string{
	AdcaInst.Name,
	AdcbInst.Name,
	AddaInst.Name,
	AddbInst.Name,
	AdddInst.Name,
	AndaInst.Name,
	AndbInst.Name,
	BitaInst.Name,
	BitbInst.Name,
	CmpaInst.Name,
	CmpbInst.Name,
	CmpdInst.Name,
	CmpsInst.Name,
	CmpuInst.Name,
	CmpxInst.Name,
	CmpyInst.Name,
	EoraInst.Name,
	EorbInst.Name,
	JmpInst.Name,
	LdaInst.Name,
	LdbInst.Name,
	LddInst.Name,
	LdsInst.Name,
	LduInst.Name,
	LdxInst.Name,
	LdyInst.Name,
	OraInst.Name,
	OrbInst.Name,
	SbcaInst.Name,
	SbcbInst.Name,
	SubaInst.Name,
	SubbInst.Name,
	SubdInst.Name,
})

// MemoryWriteInstructions contains instructions that write to an absolute memory address.
var MemoryWriteInstructions = set.NewFromSlice([]string{
	StaInst.Name,
	StbInst.Name,
	StdInst.Name,
	StsInst.Name,
	StuInst.Name,
	StxInst.Name,
	StyInst.Name,
})

// MemoryReadWriteInstructions contains instructions that both read and write
// to an absolute memory address during execution.
var MemoryReadWriteInstructions = set.NewFromSlice([]string{
	AslInst.Name,
	AsrInst.Name,
	ClrInst.Name,
	ComInst.Name,
	DecInst.Name,
	IncInst.Name,
	LsrInst.Name,
	NegInst.Name,
	RolInst.Name,
	RorInst.Name,
})
