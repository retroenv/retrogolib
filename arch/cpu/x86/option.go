package x86

// Options contains configuration options for x86 instruction processing.
type Options struct {
	// Memory options
	memorySize uint32
}

// Option represents a CPU configuration option function.
type Option func(*Options)

// NewOptions creates new options with defaults applied.
func NewOptions(options ...Option) Options {
	opts := Options{
		memorySize: 1024 * 1024, // 1MB default
	}

	for _, option := range options {
		option(&opts)
	}

	return opts
}

// WithMemorySize sets the memory size in bytes.
func WithMemorySize(size uint32) Option {
	return func(opts *Options) {
		opts.memorySize = size
	}
}
