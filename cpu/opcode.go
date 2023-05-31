package cpu

import (
	. "github.com/retroenv/retrogolib/addressing"
)

// Opcode is a NES CPU opcode that contains the instruction info and used addressing mode.
type Opcode struct {
	Instruction    *Instruction
	Addressing     Mode
	Timing         byte
	PageCrossCycle bool
}

// ReadsMemory returns whether the instruction accesses memory reading.
func (opcode Opcode) ReadsMemory(memoryReadInstructions map[string]struct{}) bool {
	switch opcode.Addressing {
	case ImmediateAddressing, ImpliedAddressing, RelativeAddressing:
		return false
	}

	_, ok := memoryReadInstructions[opcode.Instruction.Name]
	return ok
}

// WritesMemory returns whether the instruction accesses memory writing.
func (opcode Opcode) WritesMemory(memoryWriteInstructions map[string]struct{}) bool {
	switch opcode.Addressing {
	case ImmediateAddressing, ImpliedAddressing, RelativeAddressing:
		return false
	}

	_, ok := memoryWriteInstructions[opcode.Instruction.Name]
	return ok
}

// ReadWritesMemory returns whether the instruction accesses memory reading and writing.
func (opcode Opcode) ReadWritesMemory(memoryReadWriteInstructions map[string]struct{}) bool {
	switch opcode.Addressing {
	case ImmediateAddressing, ImpliedAddressing, RelativeAddressing:
		return false
	}

	_, ok := memoryReadWriteInstructions[opcode.Instruction.Name]
	return ok
}
