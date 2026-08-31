package z80

// Instructions maps instruction names to representative definitions for tooling.
var Instructions = buildInstructionRegistry(&Opcodes, &CBOpcodes, &EDOpcodes, &DDOpcodes, &FDOpcodes)

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
	// DD/FD CB instructions are decoded as four-byte prefix sequences rather
	// than entries in the ordinary prefix tables.
	instructions[DdcbShiftName] = DdcbShift
	instructions[FdcbShiftName] = FdcbShift
	instructions[InfName] = INF
	instructions[OutfName] = OUTF
	return instructions
}
