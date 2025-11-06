package chip8

import "github.com/retroenv/retrogolib/set"

// SkipInstructions contains all instructions that skip the next instruction
// based on conditional evaluation.
var SkipInstructions = set.NewFromSlice([]string{
	Se.Name,
	Sne.Name,
	Skp.Name,
	Sknp.Name,
})
