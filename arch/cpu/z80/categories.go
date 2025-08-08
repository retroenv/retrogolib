// Package z80 provides support for the Zilog Z80 CPU.
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

// ConditionalBranchingInstructions contains conditional branching instructions.
var ConditionalBranchingInstructions = set.NewFromSlice([]string{
	JrRel.Name, // JR cc,e (conditional relative jumps)
	JpAbs.Name, // JP cc,nn (conditional absolute jumps)
	Call.Name,  // CALL cc,nn (conditional calls)
	Ret.Name,   // RET cc (conditional returns)
	Djnz.Name,  // DJNZ (conditional on B != 0)
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

// IOInstructions contains all input/output instructions.
var IOInstructions = set.NewFromSlice([]string{
	InPort.Name,  // Input from port
	OutPort.Name, // Output to port
	EdIni.Name,   // Input and increment
	EdInd.Name,   // Input and decrement
	EdOuti.Name,  // Output and increment
	EdOutd.Name,  // Output and decrement
	EdInir.Name,  // Input, increment, and repeat
	EdIndr.Name,  // Input, decrement, and repeat
	EdOtir.Name,  // Output, increment, and repeat
	EdOtdr.Name,  // Output, decrement, and repeat
})

// StackInstructions contains all stack manipulation instructions.
var StackInstructions = set.NewFromSlice([]string{
	PushReg16.Name, // Push register pair to stack
	PopReg16.Name,  // Pop register pair from stack
	Call.Name,      // Call subroutine (uses stack)
	Ret.Name,       // Return from subroutine (uses stack)
	EdReti.Name,    // Return from interrupt (uses stack)
	EdRetn.Name,    // Return from NMI (uses stack)
	Rst.Name,       // Restart (uses stack)
})

// ArithmeticInstructions contains all arithmetic instructions.
var ArithmeticInstructions = set.NewFromSlice([]string{
	AddA.Name,    // Add
	AdcA.Name,    // Add with carry
	SubA.Name,    // Subtract
	SbcA.Name,    // Subtract with carry
	IncReg8.Name, // Increment
	DecReg8.Name, // Decrement
	EdNeg.Name,   // Negate (2's complement)
	Daa.Name,     // Decimal adjust accumulator
	Cpl.Name,     // Complement accumulator
})

// LogicalInstructions contains all logical instructions.
var LogicalInstructions = set.NewFromSlice([]string{
	AndA.Name,  // Logical AND
	OrA.Name,   // Logical OR
	XorA.Name,  // Logical XOR
	CpA.Name,   // Compare (logical subtraction)
	CBBit.Name, // Test bit
	CBRes.Name, // Reset bit
	CBSet.Name, // Set bit
})

// RotateShiftInstructions contains all rotate and shift instructions.
var RotateShiftInstructions = set.NewFromSlice([]string{
	Rlca.Name,  // Rotate left circular accumulator
	Rrca.Name,  // Rotate right circular accumulator
	Rla.Name,   // Rotate left accumulator through carry
	Rra.Name,   // Rotate right accumulator through carry
	CBRlc.Name, // Rotate left circular
	CBRrc.Name, // Rotate right circular
	CBRl.Name,  // Rotate left through carry
	CBRr.Name,  // Rotate right through carry
	CBSla.Name, // Shift left arithmetic
	CBSra.Name, // Shift right arithmetic
	CBSrl.Name, // Shift right logical
})

// FlagInstructions contains instructions that primarily affect flags.
var FlagInstructions = set.NewFromSlice([]string{
	Scf.Name,   // Set carry flag
	Ccf.Name,   // Complement carry flag
	Di.Name,    // Disable interrupts
	Ei.Name,    // Enable interrupts
	EdIm0.Name, // Set interrupt mode
})
