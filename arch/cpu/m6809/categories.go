package m6809

import "github.com/retroenv/retrogolib/set"

// BranchingInstructions contains all instructions that change the program counter.
var BranchingInstructions = set.NewFromSlice([]string{
	BccName,
	BcsName,
	BeqName,
	BgeName,
	BgtName,
	BhiName,
	BleName,
	BlsName,
	BltName,
	BmiName,
	BneName,
	BplName,
	BraName,
	BsrName,
	BvcName,
	BvsName,
	JmpName,
	JsrName,
	LbccName,
	LbcsName,
	LbeqName,
	LbgeName,
	LbgtName,
	LbhiName,
	LbleName,
	LblsName,
	LbltName,
	LbmiName,
	LbneName,
	LbplName,
	LbraName,
	LbsrName,
	LbvcName,
	LbvsName,
})

// NotExecutingFollowingOpcodeInstructions contains all instructions that do not
// continue to execute the following opcode after execution.
var NotExecutingFollowingOpcodeInstructions = set.NewFromSlice([]string{
	BraName,
	JmpName,
	LbraName,
	RtiName,
	RtsName,
	SwiName,
	Swi2Name,
	Swi3Name,
})

// MemoryReadInstructions contains instructions that read from an absolute memory address.
var MemoryReadInstructions = set.NewFromSlice([]string{
	AdcaName,
	AdcbName,
	AddaName,
	AddbName,
	AdddName,
	AndaName,
	AndbName,
	BitaName,
	BitbName,
	CmpaName,
	CmpbName,
	CmpdName,
	CmpsName,
	CmpuName,
	CmpxName,
	CmpyName,
	EoraName,
	EorbName,
	JmpName,
	LdaName,
	LdbName,
	LddName,
	LdsName,
	LduName,
	LdxName,
	LdyName,
	OraName,
	OrbName,
	SbcaName,
	SbcbName,
	SubaName,
	SubbName,
	SubdName,
})

// MemoryWriteInstructions contains instructions that write to an absolute memory address.
var MemoryWriteInstructions = set.NewFromSlice([]string{
	StaName,
	StbName,
	StdName,
	StsName,
	StuName,
	StxName,
	StyName,
})

// MemoryReadWriteInstructions contains instructions that both read and write
// to an absolute memory address during execution.
var MemoryReadWriteInstructions = set.NewFromSlice([]string{
	AslName,
	AsrName,
	ClrName,
	ComName,
	DecName,
	IncName,
	LsrName,
	NegName,
	RolName,
	RorName,
})
