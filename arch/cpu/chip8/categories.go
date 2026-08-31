package chip8

import "github.com/retroenv/retrogolib/set"

// SkipInstructions contains all instructions that skip the next instruction
// based on conditional evaluation.
var SkipInstructions = set.NewFromSlice([]string{
	SeName,   // SE Vx, Vy / SE Vx, byte - skip if equal
	SneName,  // SNE Vx, Vy / SNE Vx, byte - skip if not equal
	SkpName,  // SKP Vx - skip if key pressed
	SknpName, // SKNP Vx - skip if key not pressed
})

// MemoryReadInstructions contains all instructions that read from memory.
// These instructions access the main memory array to load data.
var MemoryReadInstructions = set.NewFromSlice([]string{
	DrwName, // DRW Vx, Vy, n - reads sprite data from memory at I
	LdName,  // LD Vx, [I] - reads registers from memory at I
})

// MemoryWriteInstructions contains all instructions that write to memory.
// These instructions modify the main memory array.
var MemoryWriteInstructions = set.NewFromSlice([]string{
	LdName, // LD [I], Vx and LD B, Vx write to memory at I.
})
