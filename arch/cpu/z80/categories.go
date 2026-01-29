package z80

import "github.com/retroenv/retrogolib/set"

// BranchingInstructions contains all branching and jumping instructions.
var BranchingInstructions = set.NewFromSlice([]string{
	JpAbs.Name, // Jump absolute
	JrRel.Name, // Jump relative
	Call.Name,  // Call subroutine
	Ret.Name,   // Return from subroutine
	Djnz.Name,  // Decrement B and jump if not zero
	Rst.Name,   // Restart (call to fixed address)
})

// NotExecutingFollowingOpcodeInstructions contains all instructions that jump
// to a different address and do not return to execute the following opcode.
var NotExecutingFollowingOpcodeInstructions = set.NewFromSlice([]string{
	JpAbs.Name,  // Jump absolute (unconditional)
	Ret.Name,    // Return from subroutine
	EdReti.Name, // Return from interrupt
	EdRetn.Name, // Return from non-maskable interrupt
	Halt.Name,   // Halt execution
})

// MemoryReadInstructions contains all instructions that can read from memory.
var MemoryReadInstructions = set.NewFromSlice([]string{
	LdImm8.Name,  // Load operations from memory
	AddA.Name,    // Add from memory
	AdcA.Name,    // Add with carry from memory
	SubA.Name,    // Subtract from memory
	SbcA.Name,    // Subtract with carry from memory
	AndA.Name,    // AND from memory
	OrA.Name,     // OR from memory
	XorA.Name,    // XOR from memory
	CpA.Name,     // Compare with memory
	IncReg8.Name, // Increment memory location
	DecReg8.Name, // Decrement memory location
	CBBit.Name,   // Test bit in memory
	CBRes.Name,   // Reset bit in memory
	CBSet.Name,   // Set bit in memory
	CBRl.Name,    // Rotate left memory
	CBRr.Name,    // Rotate right memory
	CBRlc.Name,   // Rotate left circular memory
	CBRrc.Name,   // Rotate right circular memory
	CBSla.Name,   // Shift left arithmetic memory
	CBSra.Name,   // Shift right arithmetic memory
	CBSrl.Name,   // Shift right logical memory
})

// MemoryWriteInstructions contains all instructions that can write to memory.
var MemoryWriteInstructions = set.NewFromSlice([]string{
	LdImm8.Name,    // Load operations to memory
	PushReg16.Name, // Push to stack
	IncReg8.Name,   // Increment memory location
	DecReg8.Name,   // Decrement memory location
	CBRes.Name,     // Reset bit in memory
	CBSet.Name,     // Set bit in memory
	CBRl.Name,      // Rotate left memory
	CBRr.Name,      // Rotate right memory
	CBRlc.Name,     // Rotate left circular memory
	CBRrc.Name,     // Rotate right circular memory
	CBSla.Name,     // Shift left arithmetic memory
	CBSra.Name,     // Shift right arithmetic memory
	CBSrl.Name,     // Shift right logical memory
})

// MemoryReadWriteInstructions contains all instructions that both read and write
// to the same memory location during execution.
var MemoryReadWriteInstructions = set.NewFromSlice([]string{
	IncReg8.Name, // Increment memory (read-modify-write)
	DecReg8.Name, // Decrement memory (read-modify-write)
	CBRes.Name,   // Reset bit in memory (read-modify-write)
	CBSet.Name,   // Set bit in memory (read-modify-write)
	CBRl.Name,    // Rotate left memory (read-modify-write)
	CBRr.Name,    // Rotate right memory (read-modify-write)
	CBRlc.Name,   // Rotate left circular memory (read-modify-write)
	CBRrc.Name,   // Rotate right circular memory (read-modify-write)
	CBSla.Name,   // Shift left arithmetic memory (read-modify-write)
	CBSra.Name,   // Shift right arithmetic memory (read-modify-write)
	CBSrl.Name,   // Shift right logical memory (read-modify-write)
})
