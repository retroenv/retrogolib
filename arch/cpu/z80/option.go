package z80

type preExecutionHook func(cpu *CPU, opcode uint8, params ...any)

// Options contains options for the CPU.
type Options struct {
	tracing                  bool
	disableUnofficialOpcodes bool
	preExecutionHook         preExecutionHook
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
