// Package parameter provides helpers to output instruction parameters as string compatible with specific assemblers.
package parameter

import (
	"fmt"

	. "github.com/retroenv/retrogolib/addressing"
)

// String returns the parameters as a string that is compatible to the assembler presented by the converter.
// nolint:cyclop
func String(converter Converter, addressing Mode, param any) (string, error) {
	switch addressing {
	case ImpliedAddressing:
		return "", nil
	case ImmediateAddressing:
		return converter.Immediate(param), nil
	case AccumulatorAddressing:
		return converter.Accumulator(), nil
	case AbsoluteAddressing:
		return converter.Absolute(param)
	case AbsoluteXAddressing:
		return converter.AbsoluteX(param)
	case AbsoluteYAddressing:
		return converter.AbsoluteY(param)
	case ZeroPageAddressing:
		return converter.ZeroPage(param)
	case ZeroPageXAddressing:
		return converter.ZeroPageX(param)
	case ZeroPageYAddressing:
		return converter.ZeroPageY(param)
	case RelativeAddressing:
		return converter.Relative(param), nil
	case IndirectAddressing:
		return converter.Indirect(param)
	case IndirectXAddressing:
		return converter.IndirectX(param)
	case IndirectYAddressing:
		return converter.IndirectY(param)
	default:
		return "", fmt.Errorf("unsupported addressing mode %d", addressing)
	}
}
