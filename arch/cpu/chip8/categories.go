package chip8

import "github.com/retroenv/retrogolib/set"

// SkipInstructions contains all instructions that skip the next instruction
// based on conditional evaluation.
var SkipInstructions = set.NewFromSlice([]string{
	Se.Name,   // SE Vx, Vy / SE Vx, byte - skip if equal
	Sne.Name,  // SNE Vx, Vy / SNE Vx, byte - skip if not equal
	Skp.Name,  // SKP Vx - skip if key pressed
	Sknp.Name, // SKNP Vx - skip if key not pressed
})

// MemoryReadInstructions contains all instructions that read from memory.
// These instructions access the main memory array to load data.
var MemoryReadInstructions = set.NewFromSlice([]string{
	Drw.Name, // DRW Vx, Vy, n - reads sprite data from memory at I
	Ld.Name,  // LD Vx, [I] - reads registers from memory at I
})

// MemoryWriteInstructions contains all instructions that write to memory.
// These instructions modify the main memory array.
var MemoryWriteInstructions = set.NewFromSlice([]string{
	Ld.Name, // LD [I], Vx - writes registers to memory at I
	Ld.Name, // LD B, Vx - writes BCD representation to memory at I
})
