package parameter

import (
	"fmt"
	"strings"

	"github.com/retroenv/retrogolib/addressing"
)

// Config contains the configuration for the parameter converter.
type Config struct {
	IndirectNoParentheses bool // do not output parentheses for indirect access
}

// Converter converts the opcode parameters to specific assembler compatible output.
type Converter struct {
	cfg Config
}

// New returns a new parameter converter.
func New(cfg Config) Converter {
	return Converter{
		cfg: cfg,
	}
}

// Immediate converts the parameters to the assembler implementation compatible string.
func (c Converter) Immediate(param any) string {
	return fmt.Sprintf("#$%02X", param)
}

// Accumulator converts the parameters to the assembler implementation compatible string.
func (c Converter) Accumulator() string {
	return "a"
}

// Absolute converts the parameters to the assembler implementation compatible string.
func (c Converter) Absolute(param any) string {
	switch val := param.(type) {
	case int, addressing.Absolute:
		return fmt.Sprintf("$%04X", val)
	case string:
		return val
	default:
		panic(fmt.Sprintf("unsupported param type %T", val))
	}
}

// AbsoluteX converts the parameters to the assembler implementation compatible string.
func (c Converter) AbsoluteX(param any) string {
	switch val := param.(type) {
	case int, addressing.Absolute, addressing.AbsoluteX:
		return fmt.Sprintf("$%04X,X", val)
	case string:
		return fmt.Sprintf("%s,X", val)
	default:
		panic(fmt.Sprintf("unsupported param type %T", val))
	}
}

// AbsoluteY converts the parameters to the assembler implementation compatible string.
func (c Converter) AbsoluteY(param any) string {
	switch val := param.(type) {
	case int, addressing.Absolute, addressing.AbsoluteY:
		return fmt.Sprintf("$%04X,Y", val)
	case string:
		return fmt.Sprintf("%s,Y", val)
	default:
		panic(fmt.Sprintf("unsupported param type %T", val))
	}
}

// ZeroPage converts the parameters to the assembler implementation compatible string.
func (c Converter) ZeroPage(param any) string {
	switch val := param.(type) {
	case int, addressing.Absolute, addressing.ZeroPage:
		return fmt.Sprintf("$%02X", val)
	case string:
		return val
	default:
		panic(fmt.Sprintf("unsupported param type %T", val))
	}
}

// ZeroPageX converts the parameters to the assembler implementation compatible string.
func (c Converter) ZeroPageX(param any) string {
	switch val := param.(type) {
	case int, addressing.Absolute, addressing.ZeroPage, addressing.ZeroPageX:
		return fmt.Sprintf("$%02X,X", val)
	case string:
		return fmt.Sprintf("%s,X", val)
	default:
		panic(fmt.Sprintf("unsupported param type %T", val))
	}
}

// ZeroPageY converts the parameters to the assembler implementation compatible string.
func (c Converter) ZeroPageY(param any) string {
	switch val := param.(type) {
	case int, addressing.Absolute, addressing.ZeroPage, addressing.ZeroPageY:
		return fmt.Sprintf("$%02X,Y", val)
	case string:
		return fmt.Sprintf("%s,Y", val)
	default:
		panic(fmt.Sprintf("unsupported param type %T", val))
	}
}

// Relative converts the parameters to the assembler implementation compatible string.
func (c Converter) Relative(param any) string {
	if param == nil {
		return ""
	}
	return fmt.Sprintf("$%04X", param)
}

// Indirect converts the parameters to the assembler implementation compatible string.
func (c Converter) Indirect(param any) string {
	address, ok := param.(addressing.Indirect)
	if ok {
		return fmt.Sprintf("($%04X)", address)
	}
	alias := param.(string)
	return fmt.Sprintf("(%s)", alias)
}

// IndirectX converts the parameters to the assembler implementation compatible string.
func (c Converter) IndirectX(param any) string {
	var builder strings.Builder
	if !c.cfg.IndirectNoParentheses {
		builder.WriteRune('(')
	}

	switch val := param.(type) {
	case addressing.Indirect, addressing.IndirectX:
		builder.WriteString(fmt.Sprintf("$%04X", val))
	case string:
		builder.WriteString(val)
	default:
		panic(fmt.Sprintf("unsupported param type %T", val))
	}

	builder.WriteString(",X")
	if !c.cfg.IndirectNoParentheses {
		builder.WriteRune(')')
	}
	s := builder.String()
	return s
}

// IndirectY converts the parameters to the assembler implementation compatible string.
func (c Converter) IndirectY(param any) string {
	var builder strings.Builder
	if !c.cfg.IndirectNoParentheses {
		builder.WriteRune('(')
	}

	switch val := param.(type) {
	case addressing.Indirect, addressing.IndirectY:
		builder.WriteString(fmt.Sprintf("$%04X", val))
	case string:
		builder.WriteString(val)
	default:
		panic(fmt.Sprintf("unsupported param type %T", val))
	}

	if !c.cfg.IndirectNoParentheses {
		builder.WriteRune(')')
	}
	builder.WriteString(",Y")
	s := builder.String()
	return s
}
