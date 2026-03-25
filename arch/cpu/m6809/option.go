package m6809

// Options contains configuration for the CPU.
type Options struct {
	tracing          bool
	preExecutionHook PreExecutionHook
}

// PreExecutionHook is a function called before each instruction is executed.
type PreExecutionHook func(cpu *CPU, ins *Instruction, params ...any)

// Option is a functional option for CPU configuration.
type Option func(*Options)

// NewOptions creates an Options instance from the provided options.
func NewOptions(opts ...Option) Options {
	o := Options{}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

// WithTracing enables instruction tracing.
func WithTracing() Option {
	return func(o *Options) { o.tracing = true }
}

// WithPreExecutionHook sets a hook called before each instruction executes.
func WithPreExecutionHook(hook PreExecutionHook) Option {
	return func(o *Options) { o.preExecutionHook = hook }
}
