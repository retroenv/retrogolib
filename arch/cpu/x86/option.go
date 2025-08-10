package x86

// Options contains configuration options for x86 CPU initialization.
type Options struct {
	systemType string

	// Initial register values
	initialIP uint16
	initialSP uint16
	initialCS uint16
	initialDS uint16
	initialES uint16
	initialSS uint16

	// Memory options
	memorySize uint32

	// Interrupt options
	interruptEnabled bool
}

// Option represents a CPU configuration option function.
type Option func(*Options)

// NewOptions creates new options with defaults applied.
func NewOptions(options ...Option) Options {
	opts := Options{
		systemType:       "",
		initialIP:        0x0000,
		initialSP:        0xFFFE,
		initialCS:        0xF000,
		initialDS:        0x0000,
		initialES:        0x0000,
		initialSS:        0x0000,
		memorySize:       1024 * 1024, // 1MB default
		interruptEnabled: false,
	}

	for _, option := range options {
		option(&opts)
	}

	return opts
}

// WithSystemType sets the system type.
func WithSystemType(systemType string) Option {
	return func(opts *Options) {
		opts.systemType = systemType
	}
}

// WithInitialIP sets the initial instruction pointer.
func WithInitialIP(ip uint16) Option {
	return func(opts *Options) {
		opts.initialIP = ip
	}
}

// WithInitialSP sets the initial stack pointer.
func WithInitialSP(sp uint16) Option {
	return func(opts *Options) {
		opts.initialSP = sp
	}
}

// WithInitialCS sets the initial code segment.
func WithInitialCS(cs uint16) Option {
	return func(opts *Options) {
		opts.initialCS = cs
	}
}

// WithInitialDS sets the initial data segment.
func WithInitialDS(ds uint16) Option {
	return func(opts *Options) {
		opts.initialDS = ds
	}
}

// WithInitialES sets the initial extra segment.
func WithInitialES(es uint16) Option {
	return func(opts *Options) {
		opts.initialES = es
	}
}

// WithInitialSS sets the initial stack segment.
func WithInitialSS(ss uint16) Option {
	return func(opts *Options) {
		opts.initialSS = ss
	}
}

// WithMemorySize sets the memory size in bytes.
func WithMemorySize(size uint32) Option {
	return func(opts *Options) {
		opts.memorySize = size
	}
}

// WithInterrupts enables interrupt handling.
func WithInterrupts(enabled bool) Option {
	return func(opts *Options) {
		opts.interruptEnabled = enabled
	}
}

// WithDOSDefaults sets reasonable defaults for DOS development.
func WithDOSDefaults() Option {
	return func(opts *Options) {
		opts.systemType = "dos"
		opts.initialCS = 0x1000 // Typical DOS code segment
		opts.initialDS = 0x1000 // Same as CS for small model
		opts.initialES = 0x1000 // Same as CS/DS
		opts.initialSS = 0x2000 // Stack segment
		opts.initialSP = 0xFFFE // Top of stack
		opts.initialIP = 0x0100 // Standard DOS .COM entry point
		opts.interruptEnabled = true
	}
}

// WithBIOSDefaults sets defaults for BIOS/ROM development.
func WithBIOSDefaults() Option {
	return func(opts *Options) {
		opts.systemType = "bios"
		opts.initialCS = 0xF000 // BIOS ROM segment
		opts.initialDS = 0x0000 // Low memory
		opts.initialES = 0x0000 // Low memory
		opts.initialSS = 0x0000 // Stack in low memory
		opts.initialSP = 0x0400 // After interrupt vector table
		opts.initialIP = 0xFFF0 // BIOS reset vector
		opts.interruptEnabled = false
	}
}
