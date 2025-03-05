// Package parameter provides helpers to output instruction parameters as string compatible with specific assemblers.
package parameter

import (
	"fmt"

	"github.com/retroenv/retrogolib/arch/cpu/m6502"
)

// String returns the parameters as a string that is compatible to the assembler presented by the converter.
// nolint:cyclop
func String(converter Converter, addressing m6502.AddressingMode, param any) (string, error) {
	switch addressing {
	case m6502.ImpliedAddressing:
		return "", nil
	case m6502.ImmediateAddressing:
		return converter.Immediate(param), nil
	case m6502.AccumulatorAddressing:
		return converter.Accumulator(), nil
	case m6502.AbsoluteAddressing:
		return converter.Absolute(param)
	case m6502.AbsoluteXAddressing:
		return converter.AbsoluteX(param)
	case m6502.AbsoluteYAddressing:
		return converter.AbsoluteY(param)
	case m6502.ZeroPageAddressing:
		return converter.ZeroPage(param)
	case m6502.ZeroPageXAddressing:
		return converter.ZeroPageX(param)
	case m6502.ZeroPageYAddressing:
		return converter.ZeroPageY(param)
	case m6502.RelativeAddressing:
		return converter.Relative(param), nil
	case m6502.IndirectAddressing:
		return converter.Indirect(param)
	case m6502.IndirectXAddressing:
		return converter.IndirectX(param)
	case m6502.IndirectYAddressing:
		return converter.IndirectY(param)
	default:
		return "", fmt.Errorf("unsupported addressing mode %d", addressing)
	}
}
