package parameter

import (
	"fmt"
	"strings"

	"github.com/retroenv/retrogolib/addressing"
)

// Config contains the configuration for the parameter converter.
type Config struct {
	ZeroPagePrefix string
	AbsolutePrefix string
	IndirectPrefix string
	IndirectSuffix string
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
	var builder strings.Builder
	builder.WriteString(c.cfg.AbsolutePrefix)

	switch val := param.(type) {
	case int, addressing.Absolute:
		builder.WriteString(fmt.Sprintf("$%04X", val))
	case string:
		builder.WriteString(val)
	default:
		panic(fmt.Sprintf("unsupported param type %T", val))
	}

	s := builder.String()
	return s
}

// AbsoluteX converts the parameters to the assembler implementation compatible string.
func (c Converter) AbsoluteX(param any) string {
	var builder strings.Builder
	builder.WriteString(c.cfg.AbsolutePrefix)

	switch val := param.(type) {
	case int, addressing.Absolute, addressing.AbsoluteX:
		builder.WriteString(fmt.Sprintf("$%04X", val))
	case string:
		builder.WriteString(val)
	default:
		panic(fmt.Sprintf("unsupported param type %T", val))
	}

	builder.WriteString(",X")
	s := builder.String()
	return s
}

// AbsoluteY converts the parameters to the assembler implementation compatible string.
func (c Converter) AbsoluteY(param any) string {
	var builder strings.Builder
	builder.WriteString(c.cfg.AbsolutePrefix)

	switch val := param.(type) {
	case int, addressing.Absolute, addressing.AbsoluteY:
		builder.WriteString(fmt.Sprintf("$%04X", val))
	case string:
		builder.WriteString(val)
	default:
		panic(fmt.Sprintf("unsupported param type %T", val))
	}

	builder.WriteString(",Y")
	s := builder.String()
	return s
}

// ZeroPage converts the parameters to the assembler implementation compatible string.
func (c Converter) ZeroPage(param any) string {
	var builder strings.Builder
	builder.WriteString(c.cfg.ZeroPagePrefix)

	switch val := param.(type) {
	case int, addressing.Absolute, addressing.ZeroPage:
		builder.WriteString(fmt.Sprintf("$%02X", val))
	case string:
		builder.WriteString(val)
	default:
		panic(fmt.Sprintf("unsupported param type %T", val))
	}

	s := builder.String()
	return s
}

// ZeroPageX converts the parameters to the assembler implementation compatible string.
func (c Converter) ZeroPageX(param any) string {
	var builder strings.Builder
	builder.WriteString(c.cfg.ZeroPagePrefix)

	switch val := param.(type) {
	case int, addressing.Absolute, addressing.ZeroPage, addressing.ZeroPageX:
		builder.WriteString(fmt.Sprintf("$%02X", val))
	case string:
		builder.WriteString(val)
	default:
		panic(fmt.Sprintf("unsupported param type %T", val))
	}

	builder.WriteString(",X")
	s := builder.String()
	return s
}

// ZeroPageY converts the parameters to the assembler implementation compatible string.
func (c Converter) ZeroPageY(param any) string {
	var builder strings.Builder
	builder.WriteString(c.cfg.ZeroPagePrefix)

	switch val := param.(type) {
	case int, addressing.Absolute, addressing.ZeroPage, addressing.ZeroPageY:
		builder.WriteString(fmt.Sprintf("$%02X", val))
	case string:
		builder.WriteString(val)
	default:
		panic(fmt.Sprintf("unsupported param type %T", val))
	}

	builder.WriteString(",Y")
	s := builder.String()
	return s
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
	var builder strings.Builder
	builder.WriteString(c.cfg.IndirectPrefix)

	address, ok := param.(addressing.Indirect)
	if ok {
		builder.WriteString(fmt.Sprintf("$%04X", address))
	} else {
		alias := param.(string)
		builder.WriteString(alias)
	}

	builder.WriteString(c.cfg.IndirectSuffix)
	s := builder.String()
	return s
}

// IndirectX converts the parameters to the assembler implementation compatible string.
func (c Converter) IndirectX(param any) string {
	var builder strings.Builder
	builder.WriteString(c.cfg.IndirectPrefix)

	switch val := param.(type) {
	case addressing.Indirect, addressing.IndirectX:
		builder.WriteString(fmt.Sprintf("$%04X", val))
	case string:
		builder.WriteString(val)
	default:
		panic(fmt.Sprintf("unsupported param type %T", val))
	}

	builder.WriteString(",X")
	builder.WriteString(c.cfg.IndirectSuffix)
	s := builder.String()
	return s
}

// IndirectY converts the parameters to the assembler implementation compatible string.
func (c Converter) IndirectY(param any) string {
	var builder strings.Builder
	builder.WriteString(c.cfg.IndirectPrefix)

	switch val := param.(type) {
	case addressing.Indirect, addressing.IndirectY:
		builder.WriteString(fmt.Sprintf("$%04X", val))
	case string:
		builder.WriteString(val)
	default:
		panic(fmt.Sprintf("unsupported param type %T", val))
	}

	builder.WriteString(c.cfg.IndirectSuffix)
	builder.WriteString(",Y")
	s := builder.String()
	return s
}
