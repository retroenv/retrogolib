package z80

type preExecutionHook func(cpu *CPU, opcode uint8, params ...any)

// IOHandler defines the interface for handling I/O port operations.
type IOHandler interface {
	ReadPort(port uint8) uint8
	WritePort(port uint8, value uint8)
}

// Options contains options for the CPU.
type Options struct {
	tracing                  bool
	disableUnofficialOpcodes bool
	preExecutionHook         preExecutionHook
	ioHandler                IOHandler
}

// Option defines a CPU parameter.
type Option func(*Options)

// NewOptions creates a new options instance from the passed options.
func NewOptions(optionList ...Option) Options {
	opts := Options{}
	for _, option := range optionList {
		option(&opts)
	}
	return opts
}

// WithTracing enables tracing for the program.
func WithTracing() func(*Options) {
	return func(options *Options) {
		options.tracing = true
	}
}

// WithUnofficialOpcodesDisabled disables support for undocumented Z80 instructions.
func WithUnofficialOpcodesDisabled() func(*Options) {
	return func(options *Options) {
		options.disableUnofficialOpcodes = true
	}
}

// WithPreExecutionHook sets a hook that is called before each instruction is executed.
// It can be used to read a memory value before the instruction overwrites it.
func WithPreExecutionHook(hook preExecutionHook) func(*Options) {
	return func(options *Options) {
		options.preExecutionHook = hook
	}
}

// WithIOHandler sets an I/O handler for port operations.
func WithIOHandler(handler IOHandler) func(*Options) {
	return func(options *Options) {
		options.ioHandler = handler
	}
}
