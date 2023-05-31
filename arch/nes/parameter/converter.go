package parameter

// Converter is an interface for the conversion of the instruction parameters to
// specific assembler implementation outputs.
type Converter interface {
	Absolute(param any) string
	AbsoluteX(param any) string
	AbsoluteY(param any) string
	Accumulator() string
	Immediate(param any) string
	Indirect(param any) string
	IndirectX(param any) string
	IndirectY(param any) string
	Relative(param any) string
	ZeroPage(param any) string
	ZeroPageX(param any) string
	ZeroPageY(param any) string
}
