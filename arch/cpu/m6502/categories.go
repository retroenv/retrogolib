// Package m6502 provides support for the MOS Technology 6502 CPU.
package m6502

// BranchingInstructions contains all branching instructions.
var BranchingInstructions = map[string]struct{}{
	Bcc.Name: {},
	Bcs.Name: {},
	Beq.Name: {},
	Bmi.Name: {},
	Bne.Name: {},
	Bpl.Name: {},
	Bvc.Name: {},
	Bvs.Name: {},
	Jmp.Name: {},
	Jsr.Name: {},
}

// NotExecutingFollowingOpcodeInstructions contains all instructions that jump
// to a different address and do not return to execute the following opcode.
var NotExecutingFollowingOpcodeInstructions = map[string]struct{}{
	Jmp.Name: {},
	Rti.Name: {},
	Rts.Name: {},
}

// MemoryReadInstructions contains all instructions that can read from an
// absolute memory address.
var MemoryReadInstructions = map[string]struct{}{
	And.Name: {},
	Bit.Name: {},
	Cmp.Name: {},
	Cpx.Name: {},
	Cpy.Name: {},
	Jmp.Name: {},
	Lda.Name: {},
	Ldx.Name: {},
	Ldy.Name: {},
	Lax.Name: {},
}

// MemoryWriteInstructions contains all instructions that can write to an
// absolute memory address.
var MemoryWriteInstructions = map[string]struct{}{
	Sta.Name: {},
	Stx.Name: {},
	Sty.Name: {},
	Sax.Name: {},
}

// MemoryReadWriteInstructions contains all instructions that can read and write
// during instruction execution an absolute memory address.
var MemoryReadWriteInstructions = map[string]struct{}{
	Adc.Name: {},
	Asl.Name: {},
	Dec.Name: {},
	Eor.Name: {},
	Inc.Name: {},
	Lsr.Name: {},
	Ora.Name: {},
	Rol.Name: {},
	Ror.Name: {},
	Sbc.Name: {},
	Dcp.Name: {},
	Isc.Name: {},
	Rla.Name: {},
	Rra.Name: {},
	Slo.Name: {},
	Sre.Name: {},
}
