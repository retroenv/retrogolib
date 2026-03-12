package sm83

// Options contains options for the CPU.
type Options struct {
	tracing bool

	preExecutionHook preExecutionHook

	initialPC uint16
	initialSP uint16
}

// Option defines a CPU parameter.
type Option func(*Options)

// NewOptions creates a new options instance from the passed options.
func NewOptions(optionList ...Option) Options {
	opts := Options{
		initialPC: 0x0100, // Game Boy default entry point
		initialSP: 0xFFFE, // Game Boy default stack pointer
	}
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

// WithPreExecutionHook sets a hook that is called before each instruction is executed.
func WithPreExecutionHook(hook preExecutionHook) func(*Options) {
	return func(options *Options) {
		options.preExecutionHook = hook
	}
}

// WithInitialPC sets the initial program counter value.
func WithInitialPC(pc uint16) func(*Options) {
	return func(options *Options) {
		options.initialPC = pc
	}
}

// WithInitialSP sets the initial stack pointer value.
func WithInitialSP(sp uint16) func(*Options) {
	return func(options *Options) {
		options.initialSP = sp
	}
}

type preExecutionHook func(cpu *CPU, opcode uint8, params ...any)
