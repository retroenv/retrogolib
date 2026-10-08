package cpu68000

type options struct {
	tracing   bool
	initialPC uint32
	initialSP uint32
}

// Option defines a CPU parameter.
type Option func(*options)

// WithTracing enables tracing for the program.
func WithTracing() Option {
	return func(options *options) {
		options.tracing = true
	}
}

// WithInitialPC sets the initial program counter value.
func WithInitialPC(pc uint32) Option {
	return func(options *options) {
		options.initialPC = pc
	}
}

// WithInitialSP sets the initial stack pointer value.
func WithInitialSP(sp uint32) Option {
	return func(options *options) {
		options.initialSP = sp
	}
}

func newOptions(optionList ...Option) options {
	opts := options{}
	for _, option := range optionList {
		if option != nil {
			option(&opts)
		}
	}
	return opts
}
