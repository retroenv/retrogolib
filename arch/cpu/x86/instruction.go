package x86

// Instruction contains information about an x86 CPU instruction.
type Instruction struct {
	Name       string // lowercased instruction name
	Unofficial bool   // unofficial instructions (not part of original 8086/8088)

	Addressing      map[AddressingMode]OpcodeInfo // addressing mode mapping to opcode info
	RegisterOpcodes map[RegisterParam]OpcodeInfo  // register-specific opcode mapping
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
