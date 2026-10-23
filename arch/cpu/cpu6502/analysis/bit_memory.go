package analysis

import (
	"github.com/retroenv/retrogolib/analysis/knownbits"
	"github.com/retroenv/retrogolib/arch/cpu/analysis"
	"github.com/retroenv/retrogolib/arch/cpu/cpu6502"
)

func (tt *transfer) bitMemory() bool {
	if tt.input.Variant < cpu6502.Variant65C02 {
		return false
	}
	code := tt.input.Opcode.Opcode
	position := knownbits.Width((code >> 4) & 7)
	set := code&0x80 != 0
	switch {
	case tt.opcode.Addressing == cpu6502.ZeroPageRelativeAddressing && code&0x0f == 0x0f:
		low, _ := knownbits.Truncate(16, 8, tt.input.Operand)
		address, _ := knownbits.ZeroExtend(8, 16, low)
		value := tt.readMemory(Address{
			Value: address,
		})
		condition, _ := knownbits.TestBit(8, value, position)
		if !set {
			condition = invert(condition)
		}
		tt.result.Flow = analysis.FlowBranch
		tt.result.Condition = condition

	case tt.opcode.Addressing == cpu6502.ZeroPageAddressing && code&0x0f == 7:
		value := tt.operand()
		mask, _ := knownbits.NewConstant[uint8](8, 1<<position)
		if set {
			value, _ = knownbits.Or(8, value, mask)
		} else {
			complement, _ := knownbits.Not(8, mask)
			value, _ = knownbits.And(8, value, complement)
		}
		tt.writeOperand(value)

	default:
		return false
	}
	return true
}
