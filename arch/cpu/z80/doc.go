// Package z80 provides Z80 instruction metadata and an instruction-level CPU emulator.
// It includes base, CB, ED, DD, and FD opcode tables and undocumented instructions.
//
// New accepts Memory and an optional eight-bit IOHandler. NewWithBus accepts a Bus
// with full 16-bit port addresses, interrupt data, and RETI notification. BasicMemory
// provides a flat 64 KiB address space; custom Memory implementations can map devices.
//
// Step accepts an interrupt, idles in HALT, or executes one instruction. Interrupts
// can release HALT, and EI delays IRQ acceptance until the next instruction completes.
// IM 0 supports device-supplied RST opcodes, with a RST 38h fallback for other values.
// IM 2 uses the device's vector byte. LD A,I/R implements the NMOS interrupt flag quirk.
//
// Timing is tracked per instruction and interrupt. T-state callbacks, contention,
// and arbitrary IM 0 instructions are not implemented. Validation results and scope
// are recorded in docs/z80-gap-closure-plan.md.
//
// Step, CheckInterrupts, and interrupt triggers are synchronized. Callers must
// serialize direct register access and interrupt configuration with execution.
// Bus callbacks and execution hooks run under the CPU lock and must not call locking
// CPU methods. State is a debugging snapshot, not a complete save/restore format.
package z80
