package z80

import "github.com/retroenv/retrogolib/arch"

type preExecutionHook func(cpu *CPU, opcode uint8, params ...any)

// IOHandler defines the interface for handling I/O port operations.
type IOHandler interface {
	ReadPort(port uint8) uint8
	WritePort(port uint8, value uint8)
}

// Options contains options for the CPU.
type Options struct {
	tracing bool

	preExecutionHook preExecutionHook
	ioHandler        IOHandler
	systemType       arch.System

	initialPC uint16
	initialSP uint16
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

// WithPreExecutionHook sets a hook that is called before each instruction is executed.
// It can be used to read a memory value before the instruction overwrites it.
func WithPreExecutionHook(hook preExecutionHook) func(*Options) {
	return func(options *Options) {
		options.preExecutionHook = hook
	}
}

// WithIOHandler sets an I/O handler for port operations.
func WithIOHandler(handler IOHandler) func(*Options) {
	return func(options *Options) {
		options.ioHandler = handler
	}
}

// WithSystemType sets the target system type for emulation.
func WithSystemType(systemType arch.System) func(*Options) {
	return func(options *Options) {
		options.systemType = systemType
		// Set system-specific defaults
		switch systemType {
		case arch.GameBoy:
			options.initialPC = 0x0100
			options.initialSP = 0xFFFE
		case arch.ZXSpectrum:
			options.initialPC = 0x0000
			options.initialSP = 0xFFFF
		default: // Generic or other systems
			options.initialPC = 0x0000
			options.initialSP = 0xFFFF
		}
	}
}

// WithInitialPC sets the initial program counter value.
func WithInitialPC(pc uint16) func(*Options) {
	return func(options *Options) {
		options.initialPC = pc
	}
}

// WithInitialSP sets the initial stack pointer value.
func WithInitialSP(sp uint16) func(*Options) {
	return func(options *Options) {
		options.initialSP = sp
	}
}
