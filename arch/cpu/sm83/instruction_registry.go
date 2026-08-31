package sm83

// Instructions maps instruction names to representative definitions for tooling.
var Instructions = buildInstructionRegistry(&Opcodes, &CBOpcodes)

func buildInstructionRegistry(tables ...*[256]Opcode) map[string]*Instruction {
	instructions := make(map[string]*Instruction, len(NameToOpcodeID))
	for _, table := range tables {
		for _, opcode := range table {
			ins := opcode.Instruction
			if ins == nil || ins.Name == "" {
				continue
			}
			if _, exists := instructions[ins.Name]; !exists {
				instructions[ins.Name] = ins
			}
		}
	}
	return instructions
}
