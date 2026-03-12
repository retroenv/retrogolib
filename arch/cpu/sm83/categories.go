package sm83

import "github.com/retroenv/retrogolib/set"

// BranchingInstructions contains all branching instructions.
var BranchingInstructions = set.NewFromSlice([]string{
	JpAbs.Name,
	JrRel.Name,
	Call.Name,
	Ret.Name,
	Rst.Name,
})

// NotExecutingFollowingOpcodeInstructions contains instructions that don't return to the next opcode.
var NotExecutingFollowingOpcodeInstructions = set.NewFromSlice([]string{
	JpAbs.Name,
	Ret.Name,
	Reti.Name,
	Halt.Name,
})

// MemoryReadInstructions contains instructions that can read from memory.
var MemoryReadInstructions = set.NewFromSlice([]string{
	LdImm8.Name,
	AddA.Name,
	AdcA.Name,
	SubA.Name,
	SbcA.Name,
	AndA.Name,
	OrA.Name,
	XorA.Name,
	CpA.Name,
	IncReg8.Name,
	DecReg8.Name,
	CBBit.Name,
	CBRes.Name,
	CBSet.Name,
	CBRl.Name,
	CBRr.Name,
	CBRlc.Name,
	CBRrc.Name,
	CBSla.Name,
	CBSra.Name,
	CBSrl.Name,
	CBSwap.Name,
})

// MemoryWriteInstructions contains instructions that can write to memory.
var MemoryWriteInstructions = set.NewFromSlice([]string{
	LdImm8.Name,
	PushReg16.Name,
	IncReg8.Name,
	DecReg8.Name,
	CBRes.Name,
	CBSet.Name,
	CBRl.Name,
	CBRr.Name,
	CBRlc.Name,
	CBRrc.Name,
	CBSla.Name,
	CBSra.Name,
	CBSrl.Name,
	CBSwap.Name,
})

// MemoryReadWriteInstructions contains instructions that both read and write the same memory location.
var MemoryReadWriteInstructions = set.NewFromSlice([]string{
	IncReg8.Name,
	DecReg8.Name,
	CBRes.Name,
	CBSet.Name,
	CBRl.Name,
	CBRr.Name,
	CBRlc.Name,
	CBRrc.Name,
	CBSla.Name,
	CBSra.Name,
	CBSrl.Name,
	CBSwap.Name,
})
