package cpu6502

import "github.com/retroenv/retrogolib/set"

// BranchingInstructions contains all branching instructions.
var BranchingInstructions = set.NewFromSlice([]string{
	Bbr0.Name,
	Bbr1.Name,
	Bbr2.Name,
	Bbr3.Name,
	Bbr4.Name,
	Bbr5.Name,
	Bbr6.Name,
	Bbr7.Name,
	Bbs0.Name,
	Bbs1.Name,
	Bbs2.Name,
	Bbs3.Name,
	Bbs4.Name,
	Bbs5.Name,
	Bbs6.Name,
	Bbs7.Name,
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

// MemoryReadInstructions contains instructions that can read a memory operand.
var MemoryReadInstructions = set.NewFromSlice([]string{
	AdcInst.Name,
	AndInst.Name,
	Bbr0.Name,
	Bbr1.Name,
	Bbr2.Name,
	Bbr3.Name,
	Bbr4.Name,
	Bbr5.Name,
	Bbr6.Name,
	Bbr7.Name,
	Bbs0.Name,
	Bbs1.Name,
	Bbs2.Name,
	Bbs3.Name,
	Bbs4.Name,
	Bbs5.Name,
	Bbs6.Name,
	Bbs7.Name,
	BitInst.Name,
	CmpInst.Name,
	CpxInst.Name,
	CpyInst.Name,
	EorInst.Name,
	JmpInst.Name,
	LasInst.Name,
	LaxInst.Name,
	LdaInst.Name,
	LdxInst.Name,
	LdyInst.Name,
	NopInst.Name,
	OraInst.Name,
	SbcInst.Name,
	TrbInst.Name,
	TsbInst.Name,
})

// MemoryWriteInstructions contains instructions that can write a memory operand.
var MemoryWriteInstructions = set.NewFromSlice([]string{
	SaxInst.Name,
	ShaInst.Name,
	ShxInst.Name,
	ShyInst.Name,
	StaInst.Name,
	StxInst.Name,
	StyInst.Name,
	StzInst.Name,
	TasInst.Name,
})

// MemoryReadWriteInstructions contains instructions that can read and write the
// same memory operand during execution.
var MemoryReadWriteInstructions = set.NewFromSlice([]string{
	AslInst.Name,
	DcpInst.Name,
	DecInst.Name,
	IncInst.Name,
	IscInst.Name,
	LsrInst.Name,
	RlaInst.Name,
	Rmb0.Name,
	Rmb1.Name,
	Rmb2.Name,
	Rmb3.Name,
	Rmb4.Name,
	Rmb5.Name,
	Rmb6.Name,
	Rmb7.Name,
	RolInst.Name,
	RorInst.Name,
	RraInst.Name,
	SloInst.Name,
	Smb0.Name,
	Smb1.Name,
	Smb2.Name,
	Smb3.Name,
	Smb4.Name,
	Smb5.Name,
	Smb6.Name,
	Smb7.Name,
	SreInst.Name,
	TrbInst.Name,
	TsbInst.Name,
})
