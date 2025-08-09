package x86

// Instruction contains information about an x86 CPU instruction.
type Instruction struct {
	Name       string // lowercased instruction name
	Unofficial bool   // unofficial instructions (not part of original 8086/8088)

	Addressing      map[AddressingMode]OpcodeInfo // addressing mode mapping to opcode info
	RegisterOpcodes map[RegisterParam]OpcodeInfo  // register-specific opcode mapping

	NoParamFunc func(c *CPU) error                // emulation function with no parameters
	ParamFunc   func(c *CPU, params ...any) error // emulation function with parameters
}

// HasAddressing returns whether the instruction has any of the passed addressing modes.
func (ins Instruction) HasAddressing(modes ...AddressingMode) bool {
	for _, mode := range modes {
		if _, exists := ins.Addressing[mode]; exists {
			return true
		}
	}
	return false
}

// GetOpcodeByRegister returns opcode info for a specific register parameter.
func (ins Instruction) GetOpcodeByRegister(register RegisterParam) (OpcodeInfo, bool) {
	if ins.RegisterOpcodes == nil {
		// Fall back to Addressing map if RegisterOpcodes is not defined
		for _, info := range ins.Addressing {
			return info, true
		}
		return OpcodeInfo{}, false
	}

	info, exists := ins.RegisterOpcodes[register]
	return info, exists
}

// GetAllRegisterVariants returns all register variants for this instruction.
func (ins Instruction) GetAllRegisterVariants() map[RegisterParam]OpcodeInfo {
	if ins.RegisterOpcodes == nil {
		return nil
	}

	// Return a copy to prevent external modification
	// Pre-allocate with exact capacity for optimal performance
	variants := make(map[RegisterParam]OpcodeInfo, len(ins.RegisterOpcodes))
	for reg, info := range ins.RegisterOpcodes {
		variants[reg] = info
	}
	return variants
}

// GetOpcodeInfo returns opcode info for the specified addressing mode.
func (ins Instruction) GetOpcodeInfo(mode AddressingMode) (OpcodeInfo, bool) {
	info, exists := ins.Addressing[mode]
	return info, exists
}

// SupportsRegister returns whether the instruction supports the specified register.
func (ins Instruction) SupportsRegister(register RegisterParam) bool {
	if ins.RegisterOpcodes == nil {
		return false
	}
	_, exists := ins.RegisterOpcodes[register]
	return exists
}

// GetSupportedAddressingModes returns all supported addressing modes.
func (ins Instruction) GetSupportedAddressingModes() []AddressingMode {
	modes := make([]AddressingMode, 0, len(ins.Addressing))
	for mode := range ins.Addressing {
		modes = append(modes, mode)
	}
	return modes
}

// GetSupportedRegisters returns all supported register parameters.
func (ins Instruction) GetSupportedRegisters() []RegisterParam {
	registers := make([]RegisterParam, 0, len(ins.RegisterOpcodes))
	for register := range ins.RegisterOpcodes {
		registers = append(registers, register)
	}
	return registers
}

// IsValid returns whether the instruction has valid opcode mappings.
func (ins Instruction) IsValid() bool {
	return len(ins.Addressing) > 0 || len(ins.RegisterOpcodes) > 0
}

// RequiresParameters returns whether the instruction requires parameters.
func (ins Instruction) RequiresParameters() bool {
	return ins.ParamFunc != nil && ins.NoParamFunc == nil
}

// SupportsNoParameters returns whether the instruction can execute without parameters.
func (ins Instruction) SupportsNoParameters() bool {
	return ins.NoParamFunc != nil
}

// Execute executes the instruction with the given parameters.
func (ins Instruction) Execute(c *CPU, params ...any) error {
	if len(params) == 0 && ins.NoParamFunc != nil {
		return ins.NoParamFunc(c)
	}
	if ins.ParamFunc != nil {
		return ins.ParamFunc(c, params...)
	}
	return ErrInvalidInstruction
}

// InstructionSet contains all x86 instructions indexed by name.
type InstructionSet map[string]*Instruction

// NewInstructionSet creates a new instruction set with all x86 instructions.
func NewInstructionSet() InstructionSet {
	return make(InstructionSet)
}

// AddInstruction adds an instruction to the instruction set.
func (is InstructionSet) AddInstruction(name string, instruction *Instruction) {
	is[name] = instruction
}

// GetInstruction retrieves an instruction by name.
func (is InstructionSet) GetInstruction(name string) (*Instruction, bool) {
	instruction, exists := is[name]
	return instruction, exists
}

// GetInstructionNames returns all instruction names in the set.
func (is InstructionSet) GetInstructionNames() []string {
	names := make([]string, 0, len(is))
	for name := range is {
		names = append(names, name)
	}
	return names
}

// Count returns the number of instructions in the set.
func (is InstructionSet) Count() int {
	return len(is)
}

// HasInstruction returns whether the instruction set contains the named instruction.
func (is InstructionSet) HasInstruction(name string) bool {
	_, exists := is[name]
	return exists
}
