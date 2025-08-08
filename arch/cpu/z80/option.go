package z80

// Options contains configuration for the CPU emulator.
type Options struct {
	// DisableUnofficialOpcodes disables support for undocumented Z80 instructions.
	DisableUnofficialOpcodes bool

	// TraceExecution enables instruction tracing for debugging.
	TraceExecution bool
}

// Option is a function that modifies CPU options.
type Option func(*Options)

// NewOptions creates a new Options struct with the given options applied.
func NewOptions(options ...Option) Options {
	opts := Options{}
	for _, option := range options {
		option(&opts)
	}
	return opts
}

// WithUnofficialOpcodesDisabled disables support for undocumented Z80 instructions.
func WithUnofficialOpcodesDisabled() Option {
	return func(opts *Options) {
		opts.DisableUnofficialOpcodes = true
	}
}

// WithTraceExecution enables instruction tracing for debugging.
func WithTraceExecution() Option {
	return func(opts *Options) {
		opts.TraceExecution = true
	}
}
