package cpu6809

// PreExecutionHook runs after operand decoding and before instruction execution.
// Params is empty for implied instructions and otherwise contains one typed operand.
type PreExecutionHook func(cpu *CPU, ins *Instruction, params ...any)

type options struct {
	tracing          bool
	preExecutionHook PreExecutionHook
}

// Option configures a CPU.
type Option func(*options)

// WithTracing enables instruction tracing.
func WithTracing() Option {
	return func(o *options) { o.tracing = true }
}

// WithPreExecutionHook sets a hook called before each instruction executes.
func WithPreExecutionHook(hook PreExecutionHook) Option {
	return func(o *options) { o.preExecutionHook = hook }
}

func newOptions(opts ...Option) options {
	o := options{}
	for _, opt := range opts {
		if opt != nil {
			opt(&o)
		}
	}
	return o
}
