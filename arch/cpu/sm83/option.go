package sm83

// PreExecutionHook runs after operand decoding and before instruction execution.
type PreExecutionHook func(cpu *CPU, opcode uint8, params ...any)

type options struct {
	tracing          bool
	preExecutionHook PreExecutionHook
	initialPC        uint16
	initialSP        uint16
}

// Option defines a CPU parameter.
type Option func(*options)

// WithTracing enables tracing for the program.
func WithTracing() Option {
	return func(options *options) {
		options.tracing = true
	}
}

// WithPreExecutionHook sets a hook that is called before each instruction is executed.
func WithPreExecutionHook(hook PreExecutionHook) Option {
	return func(options *options) {
		options.preExecutionHook = hook
	}
}

// WithInitialPC sets the initial program counter value.
func WithInitialPC(pc uint16) Option {
	return func(options *options) {
		options.initialPC = pc
	}
}

// WithInitialSP sets the initial stack pointer value.
func WithInitialSP(sp uint16) Option {
	return func(options *options) {
		options.initialSP = sp
	}
}

func newOptions(optionList ...Option) options {
	opts := options{
		initialPC: 0x0100,
		initialSP: 0xFFFE,
	}
	for _, option := range optionList {
		if option != nil {
			option(&opts)
		}
	}
	return opts
}
