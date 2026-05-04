package chip8

import "github.com/retroenv/retrogolib/set"

// SkipInstructions contains all instructions that skip the next instruction
// based on conditional evaluation.
var SkipInstructions = set.NewFromSlice([]string{
	SeInst.Name,   // SE Vx, Vy / SE Vx, byte - skip if equal
	SneInst.Name,  // SNE Vx, Vy / SNE Vx, byte - skip if not equal
	SkpInst.Name,  // SKP Vx - skip if key pressed
	SknpInst.Name, // SKNP Vx - skip if key not pressed
})

// MemoryReadInstructions contains all instructions that read from memory.
// These instructions access the main memory array to load data.
var MemoryReadInstructions = set.NewFromSlice([]string{
	DrwInst.Name, // DRW Vx, Vy, n - reads sprite data from memory at I
	LdInst.Name,  // LD Vx, [I] - reads registers from memory at I
})

// MemoryWriteInstructions contains all instructions that write to memory.
// These instructions modify the main memory array.
var MemoryWriteInstructions = set.NewFromSlice([]string{
	LdInst.Name, // LD [I], Vx - writes registers to memory at I
	LdInst.Name, // LD B, Vx - writes BCD representation to memory at I
})
