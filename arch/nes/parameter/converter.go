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
func (c Converter) Absolute(param any) (string, error) {
	var builder strings.Builder
	builder.WriteString(c.cfg.AbsolutePrefix)

	switch val := param.(type) {
	case int, addressing.Absolute:
		builder.WriteString(fmt.Sprintf("$%04X", val))
	case string:
		builder.WriteString(val)
	default:
		return "", fmt.Errorf("unsupported param type %T", val)
	}

	return builder.String(), nil
}

// AbsoluteX converts the parameters to the assembler implementation compatible string.
func (c Converter) AbsoluteX(param any) (string, error) {
	var builder strings.Builder
	builder.WriteString(c.cfg.AbsolutePrefix)

	switch val := param.(type) {
	case int, addressing.Absolute, addressing.AbsoluteX:
		builder.WriteString(fmt.Sprintf("$%04X", val))
	case string:
		builder.WriteString(val)
	default:
		return "", fmt.Errorf("unsupported param type %T", val)
	}

	builder.WriteString(",X")
	return builder.String(), nil
}

// AbsoluteY converts the parameters to the assembler implementation compatible string.
func (c Converter) AbsoluteY(param any) (string, error) {
	var builder strings.Builder
	builder.WriteString(c.cfg.AbsolutePrefix)

	switch val := param.(type) {
	case int, addressing.Absolute, addressing.AbsoluteY:
		builder.WriteString(fmt.Sprintf("$%04X", val))
	case string:
		builder.WriteString(val)
	default:
		return "", fmt.Errorf("unsupported param type %T", val)
	}

	builder.WriteString(",Y")
	return builder.String(), nil
}

// ZeroPage converts the parameters to the assembler implementation compatible string.
func (c Converter) ZeroPage(param any) (string, error) {
	var builder strings.Builder
	builder.WriteString(c.cfg.ZeroPagePrefix)

	switch val := param.(type) {
	case int, addressing.Absolute, addressing.ZeroPage:
		builder.WriteString(fmt.Sprintf("$%02X", val))
	case string:
		builder.WriteString(val)
	default:
		return "", fmt.Errorf("unsupported param type %T", val)
	}

	return builder.String(), nil
}

// ZeroPageX converts the parameters to the assembler implementation compatible string.
func (c Converter) ZeroPageX(param any) (string, error) {
	var builder strings.Builder
	builder.WriteString(c.cfg.ZeroPagePrefix)

	switch val := param.(type) {
	case int, addressing.Absolute, addressing.ZeroPage, addressing.ZeroPageX:
		builder.WriteString(fmt.Sprintf("$%02X", val))
	case string:
		builder.WriteString(val)
	default:
		return "", fmt.Errorf("unsupported param type %T", val)
	}

	builder.WriteString(",X")
	return builder.String(), nil
}

// ZeroPageY converts the parameters to the assembler implementation compatible string.
func (c Converter) ZeroPageY(param any) (string, error) {
	var builder strings.Builder
	builder.WriteString(c.cfg.ZeroPagePrefix)

	switch val := param.(type) {
	case int, addressing.Absolute, addressing.ZeroPage, addressing.ZeroPageY:
		builder.WriteString(fmt.Sprintf("$%02X", val))
	case string:
		builder.WriteString(val)
	default:
		return "", fmt.Errorf("unsupported param type %T", val)
	}

	builder.WriteString(",Y")
	return builder.String(), nil
}

// Relative converts the parameters to the assembler implementation compatible string.
func (c Converter) Relative(param any) string {
	if param == nil {
		return ""
	}
	return fmt.Sprintf("$%04X", param)
}

// Indirect converts the parameters to the assembler implementation compatible string.
func (c Converter) Indirect(param any) (string, error) {
	var builder strings.Builder
	builder.WriteString(c.cfg.IndirectPrefix)

	address, ok := param.(addressing.Indirect)
	if ok {
		builder.WriteString(fmt.Sprintf("$%04X", address))
	} else {
		alias, ok := param.(string)
		if !ok {
			return "", fmt.Errorf("unsupported param type %T", param)
		}
		builder.WriteString(alias)
	}

	builder.WriteString(c.cfg.IndirectSuffix)
	return builder.String(), nil
}

// IndirectX converts the parameters to the assembler implementation compatible string.
func (c Converter) IndirectX(param any) (string, error) {
	var builder strings.Builder
	builder.WriteString(c.cfg.IndirectPrefix)

	switch val := param.(type) {
	case addressing.Indirect, addressing.IndirectX:
		builder.WriteString(fmt.Sprintf("$%04X", val))
	case string:
		builder.WriteString(val)
	default:
		return "", fmt.Errorf("unsupported param type %T", val)
	}

	builder.WriteString(",X")
	builder.WriteString(c.cfg.IndirectSuffix)
	return builder.String(), nil
}

// IndirectY converts the parameters to the assembler implementation compatible string.
func (c Converter) IndirectY(param any) (string, error) {
	var builder strings.Builder
	builder.WriteString(c.cfg.IndirectPrefix)

	switch val := param.(type) {
	case addressing.Indirect, addressing.IndirectY:
		builder.WriteString(fmt.Sprintf("$%04X", val))
	case string:
		builder.WriteString(val)
	default:
		return "", fmt.Errorf("unsupported param type %T", val)
	}

	builder.WriteString(c.cfg.IndirectSuffix)
	builder.WriteString(",Y")
	return builder.String(), nil
}
