package z80

import "github.com/retroenv/retrogolib/arch"

// IOHandler defines the interface for handling I/O port operations.
type IOHandler interface {
	ReadPort(port uint8) uint8
	WritePort(port uint8, value uint8)
}

// PreExecutionHook runs after operand decoding and before instruction execution.
// Params is empty for implied instructions and otherwise contains decoded operands.
type PreExecutionHook func(cpu *CPU, opcode uint8, params ...any)

type options struct {
	tracing          bool
	preExecutionHook PreExecutionHook
	ioHandler        IOHandler
	systemType       arch.System
	initialPC        uint16
	initialSP        uint16
}

// Option defines a CPU parameter.
type Option func(*options)

// WithTracing enables tracing for the program.
func WithTracing() Option {
	return func(options *options) {
		options.tracing = true
	}
}

// WithPreExecutionHook sets a hook that is called before each instruction is executed.
// It can be used to read a memory value before the instruction overwrites it.
func WithPreExecutionHook(hook PreExecutionHook) Option {
	return func(options *options) {
		options.preExecutionHook = hook
	}
}

// WithIOHandler sets an I/O handler for port operations.
func WithIOHandler(handler IOHandler) Option {
	return func(options *options) {
		options.ioHandler = handler
	}
}

// WithSystemType sets the target system type for emulation.
func WithSystemType(systemType arch.System) Option {
	return func(options *options) {
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
func WithInitialPC(pc uint16) Option {
	return func(options *options) {
		options.initialPC = pc
	}
}

// WithInitialSP sets the initial stack pointer value.
func WithInitialSP(sp uint16) Option {
	return func(options *options) {
		options.initialSP = sp
	}
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
