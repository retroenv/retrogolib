// Package z80 provides support for the Zilog Z80 CPU.
package z80

import "github.com/retroenv/retrogolib/set"

// BranchingInstructions contains all branching and jumping instructions.
var BranchingInstructions = set.NewFromSlice([]string{
	"jp",   // Jump absolute
	"jr",   // Jump relative
	"call", // Call subroutine
	"ret",  // Return from subroutine
	"djnz", // Decrement B and jump if not zero
	"rst",  // Restart (call to fixed address)
})

// ConditionalBranchingInstructions contains conditional branching instructions.
var ConditionalBranchingInstructions = set.NewFromSlice([]string{
	"jr",   // JR cc,e (conditional relative jumps)
	"jp",   // JP cc,nn (conditional absolute jumps)
	"call", // CALL cc,nn (conditional calls)
	"ret",  // RET cc (conditional returns)
	"djnz", // DJNZ (conditional on B != 0)
})

// NotExecutingFollowingOpcodeInstructions contains all instructions that jump
// to a different address and do not return to execute the following opcode.
var NotExecutingFollowingOpcodeInstructions = set.NewFromSlice([]string{
	"jp",   // Jump absolute (unconditional)
	"ret",  // Return from subroutine
	"reti", // Return from interrupt
	"retn", // Return from non-maskable interrupt
	"halt", // Halt execution
})

// MemoryReadInstructions contains all instructions that can read from memory.
var MemoryReadInstructions = set.NewFromSlice([]string{
	"ld",  // Load operations from memory
	"add", // Add from memory
	"adc", // Add with carry from memory
	"sub", // Subtract from memory
	"sbc", // Subtract with carry from memory
	"and", // AND from memory
	"or",  // OR from memory
	"xor", // XOR from memory
	"cp",  // Compare with memory
	"inc", // Increment memory location
	"dec", // Decrement memory location
	"bit", // Test bit in memory
	"res", // Reset bit in memory
	"set", // Set bit in memory
	"rl",  // Rotate left memory
	"rr",  // Rotate right memory
	"rlc", // Rotate left circular memory
	"rrc", // Rotate right circular memory
	"sla", // Shift left arithmetic memory
	"sra", // Shift right arithmetic memory
	"srl", // Shift right logical memory
})

// MemoryWriteInstructions contains all instructions that can write to memory.
var MemoryWriteInstructions = set.NewFromSlice([]string{
	"ld",   // Load operations to memory
	"push", // Push to stack
	"inc",  // Increment memory location
	"dec",  // Decrement memory location
	"res",  // Reset bit in memory
	"set",  // Set bit in memory
	"rl",   // Rotate left memory
	"rr",   // Rotate right memory
	"rlc",  // Rotate left circular memory
	"rrc",  // Rotate right circular memory
	"sla",  // Shift left arithmetic memory
	"sra",  // Shift right arithmetic memory
	"srl",  // Shift right logical memory
})

// MemoryReadWriteInstructions contains all instructions that both read and write
// to the same memory location during execution.
var MemoryReadWriteInstructions = set.NewFromSlice([]string{
	"inc", // Increment memory (read-modify-write)
	"dec", // Decrement memory (read-modify-write)
	"res", // Reset bit in memory (read-modify-write)
	"set", // Set bit in memory (read-modify-write)
	"rl",  // Rotate left memory (read-modify-write)
	"rr",  // Rotate right memory (read-modify-write)
	"rlc", // Rotate left circular memory (read-modify-write)
	"rrc", // Rotate right circular memory (read-modify-write)
	"sla", // Shift left arithmetic memory (read-modify-write)
	"sra", // Shift right arithmetic memory (read-modify-write)
	"srl", // Shift right logical memory (read-modify-write)
})

// IOInstructions contains all input/output instructions.
var IOInstructions = set.NewFromSlice([]string{
	"in",   // Input from port
	"out",  // Output to port
	"ini",  // Input and increment
	"ind",  // Input and decrement
	"outi", // Output and increment
	"outd", // Output and decrement
	"inir", // Input, increment, and repeat
	"indr", // Input, decrement, and repeat
	"otir", // Output, increment, and repeat
	"otdr", // Output, decrement, and repeat
})

// StackInstructions contains all stack manipulation instructions.
var StackInstructions = set.NewFromSlice([]string{
	"push", // Push register pair to stack
	"pop",  // Pop register pair from stack
	"call", // Call subroutine (uses stack)
	"ret",  // Return from subroutine (uses stack)
	"reti", // Return from interrupt (uses stack)
	"retn", // Return from NMI (uses stack)
	"rst",  // Restart (uses stack)
})

// ArithmeticInstructions contains all arithmetic instructions.
var ArithmeticInstructions = set.NewFromSlice([]string{
	"add", // Add
	"adc", // Add with carry
	"sub", // Subtract
	"sbc", // Subtract with carry
	"inc", // Increment
	"dec", // Decrement
	"neg", // Negate (2's complement)
	"daa", // Decimal adjust accumulator
	"cpl", // Complement accumulator
})

// LogicalInstructions contains all logical instructions.
var LogicalInstructions = set.NewFromSlice([]string{
	"and", // Logical AND
	"or",  // Logical OR
	"xor", // Logical XOR
	"cp",  // Compare (logical subtraction)
	"bit", // Test bit
	"res", // Reset bit
	"set", // Set bit
})

// RotateShiftInstructions contains all rotate and shift instructions.
var RotateShiftInstructions = set.NewFromSlice([]string{
	"rlca", // Rotate left circular accumulator
	"rrca", // Rotate right circular accumulator
	"rla",  // Rotate left accumulator through carry
	"rra",  // Rotate right accumulator through carry
	"rlc",  // Rotate left circular
	"rrc",  // Rotate right circular
	"rl",   // Rotate left through carry
	"rr",   // Rotate right through carry
	"sla",  // Shift left arithmetic
	"sra",  // Shift right arithmetic
	"srl",  // Shift right logical
})

// FlagInstructions contains instructions that primarily affect flags.
var FlagInstructions = set.NewFromSlice([]string{
	"scf", // Set carry flag
	"ccf", // Complement carry flag
	"di",  // Disable interrupts
	"ei",  // Enable interrupts
	"im",  // Set interrupt mode
})
