package cpu6502

// CPUVariant represents a CPU variant within the 6502 family.
type CPUVariant int

const (
	VariantNMOS6502      CPUVariant = iota // Original NMOS 6502
	VariantNES6502                         // NES 2A03/2A07: NMOS 6502 with decimal mode disabled
	Variant6507                            // MOS 6507: 6502 with 13-bit address bus, no IRQ/NMI pins (Atari 2600)
	Variant6510                            // MOS 6510: 6502 with built-in 6-bit I/O port at $0000-$0001 (Commodore 64)
	Variant65C02                           // WDC 65C02 (base), includes Rockwell extensions (RMB/SMB/BBR/BBS)
	VariantSynertek65C02                   // Synertek 65C02: 65C02 without Rockwell bit-manipulation extensions
)

// PreExecutionHook is called after instruction decoding and before execution.
type PreExecutionHook func(cpu *CPU, ins *Instruction, params ...any)

// Options contains options for the CPU.
type Options struct {
	variant          CPUVariant
	tracing          bool
	preExecutionHook PreExecutionHook
}

// Option configures a CPU constructor option.
type Option func(*Options)

// NewOptions creates a new options instance from the passed options.
func NewOptions(optionList ...Option) Options {
	opts := Options{}
	for _, option := range optionList {
		if option == nil {
			continue
		}
		option(&opts)
	}
	return opts
}

// WithTracing enables tracing for the program.
func WithTracing() Option {
	return func(options *Options) {
		options.tracing = true
	}
}

// WithPreExecutionHook sets a hook that is called before each instruction is executed.
// It can be used to read a memory value before the instruction overwrites it.
func WithPreExecutionHook(hook PreExecutionHook) Option {
	return func(options *Options) {
		options.preExecutionHook = hook
	}
}

// WithVariant sets the CPU variant.
func WithVariant(v CPUVariant) Option {
	return func(options *Options) {
		options.variant = v
	}
}
