package x86

// Options contains configuration options for x86 instruction processing.
type Options struct {
}

// Option represents a CPU configuration option function.
type Option func(*Options)

// NewOptions creates new options with defaults applied.
func NewOptions(options ...Option) Options {
	opts := Options{}

	for _, option := range options {
		option(&opts)
	}

	return opts
}
