package cpu65816

// PreExecutionHook is a function called before each instruction is executed.
type PreExecutionHook func(cpu *CPU, ins *Instruction, params ...any)

type options struct {
	tracing          bool
	preExecutionHook PreExecutionHook
}

// Option is a functional option for CPU configuration.
type Option func(*options)

// WithTracing enables instruction tracing.
func WithTracing() Option {
	return func(o *options) { o.tracing = true }
}

// WithPreExecutionHook sets a hook called before each instruction executes.
func WithPreExecutionHook(hook PreExecutionHook) Option {
	return func(o *options) { o.preExecutionHook = hook }
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
