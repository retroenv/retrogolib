// Package cpu68000 emulates the original Motorola 68000, with 32-bit registers,
// big-endian memory access, and a 24-bit external address bus.
//
// A host supplies a Bus, optionally implementing BusErrorHandler to reject memory
// transfers. Step executes an instruction or services an interrupt. Misaligned word
// accesses and rejected transfers raise guest exceptions; faults during bus/address
// error handling halt execution until Reset. CPU long accesses use two word transfers.
//
// Cycles accounts for addressing modes and operand-dependent instruction timing.
// Instruction prefetch, wait states, and complete bus-cycle timing are not emulated.
// The conformance results and remaining discrepancies are documented in
// docs/cpu68000-gap-closure-plan.md.
//
// Step, Reset, and interrupt triggers are synchronized. Callers must serialize direct
// register access and status updates with execution. Bus callbacks run under the CPU
// lock and must not call locking CPU methods. State is a debugging snapshot, not a
// complete save/restore format.
package cpu68000
