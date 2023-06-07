// Package cpu provides general CPU related type support.
package cpu

import (
	. "github.com/retroenv/retrogolib/addressing"
)

// AddressingInfo contains the opcode and timing info for an instruction addressing mode.
type AddressingInfo struct {
	Opcode byte
	Size   int
}

// Instruction contains information about a NES CPU instruction.
type Instruction struct {
	Name       string
	Unofficial bool

	// instruction has no parameters
	NoParamFunc func()
	// instruction has parameters
	ParamFunc func(params ...any)

	// maps addressing mode to cpu cycles
	Addressing map[Mode]AddressingInfo
}

// HasAddressing returns whether the instruction has any of the passed addressing modes.
func (ins Instruction) HasAddressing(flags ...Mode) bool {
	for _, flag := range flags {
		_, ok := ins.Addressing[flag]
		if ok {
			return ok
		}
	}
	return false
}
