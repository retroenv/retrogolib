package m68000

import "github.com/retroenv/retrogolib/arch"

// Options contains options for the CPU.
type Options struct {
	tracing    bool
	systemType arch.System
	initialPC  uint32
	initialSP  uint32
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

// WithSystemType sets the target system type for emulation.
func WithSystemType(systemType arch.System) func(*Options) {
	return func(options *Options) {
		options.systemType = systemType
	}
}

// WithInitialPC sets the initial program counter value.
func WithInitialPC(pc uint32) func(*Options) {
	return func(options *Options) {
		options.initialPC = pc
	}
}

// WithInitialSP sets the initial stack pointer value.
func WithInitialSP(sp uint32) func(*Options) {
	return func(options *Options) {
		options.initialSP = sp
	}
}
