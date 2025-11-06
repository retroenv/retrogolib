package m6502

import "github.com/retroenv/retrogolib/set"

// BranchingInstructions contains all branching instructions.
var BranchingInstructions = set.NewFromSlice([]string{
	Bcc.Name,
	Bcs.Name,
	Beq.Name,
	Bmi.Name,
	Bne.Name,
	Bpl.Name,
	Bvc.Name,
	Bvs.Name,
	Jmp.Name,
	Jsr.Name,
})

// NotExecutingFollowingOpcodeInstructions contains all instructions that jump
// to a different address and do not return to execute the following opcode.
var NotExecutingFollowingOpcodeInstructions = set.NewFromSlice([]string{
	Jmp.Name,
	Rti.Name,
	Rts.Name,
})

// MemoryReadInstructions contains all instructions that can read from an
// absolute memory address.
var MemoryReadInstructions = set.NewFromSlice([]string{
	And.Name,
	Bit.Name,
	Cmp.Name,
	Cpx.Name,
	Cpy.Name,
	Jmp.Name,
	Lda.Name,
	Ldx.Name,
	Ldy.Name,
	Lax.Name,
})

// MemoryWriteInstructions contains all instructions that can write to an
// absolute memory address.
var MemoryWriteInstructions = set.NewFromSlice([]string{
	Sta.Name,
	Stx.Name,
	Sty.Name,
	Sax.Name,
})

// MemoryReadWriteInstructions contains all instructions that can read and write
// during instruction execution an absolute memory address.
var MemoryReadWriteInstructions = set.NewFromSlice([]string{
	Adc.Name,
	Asl.Name,
	Dec.Name,
	Eor.Name,
	Inc.Name,
	Lsr.Name,
	Ora.Name,
	Rol.Name,
	Ror.Name,
	Sbc.Name,
	Dcp.Name,
	Isc.Name,
	Rla.Name,
	Rra.Name,
	Slo.Name,
	Sre.Name,
})
