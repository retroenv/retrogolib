package cpu6502

// StackOperation identifies a CPU stack push or pull.
type StackOperation uint8

const (
	// StackPush decrements the stack pointer after a write.
	StackPush StackOperation = iota
	// StackPull increments the stack pointer before a read.
	StackPull
)

// InterruptSource identifies an interrupt that caused a stack operation.
type InterruptSource uint8

const (
	// InterruptNone means that an instruction caused the stack operation.
	InterruptNone InterruptSource = iota
	// InterruptNMI identifies non-maskable interrupt entry.
	InterruptNMI
	// InterruptIRQ identifies maskable interrupt entry.
	InterruptIRQ
	// InterruptBRK identifies BRK instruction entry.
	InterruptBRK
)

// StackEvent reports one stack pointer change at the operation boundary.
// TXS, Reset, and direct assignments to SP do not produce events.
type StackEvent struct {
	// After is the stack pointer after the push or pull.
	After byte
	// Before is the stack pointer before the push or pull.
	Before byte
	// Operation identifies whether the event is a push or pull.
	Operation StackOperation

	// Cycle is the CPU cycle count at the start of the instruction or interrupt.
	// All events from that instruction or interrupt share this value; it is not
	// the cycle of the individual stack bus access.
	Cycle uint64
	// Interrupt identifies the interrupt entry that caused the stack operation.
	// It is InterruptNone for ordinary instructions.
	Interrupt InterruptSource
	// Opcode is the executed instruction byte. It is zero for IRQ/NMI entry,
	// where no instruction is decoded; Interrupt distinguishes these from BRK.
	Opcode byte
	// PC is the instruction address, or the interrupted PC for IRQ/NMI entry.
	PC uint16
}

// StackEventHook receives events synchronously in stack operation order.
// It runs after SP changes: after the memory write for a push, and before the
// memory read for a pull. The instruction or interrupt is still in progress.
// Hooks must not reenter execution or mutate the CPU or its memory.
type StackEventHook func(StackEvent)

// Wrapped reports whether SP wrapped from 0x00 to 0xff on a push or from
// 0xff to 0x00 on a pull. Stack accesses remain within page one. A wrap is valid
// hardware behavior and does not by itself establish stack overflow or underflow.
func (e StackEvent) Wrapped() bool {
	return e.Operation == StackPush && e.Before == 0x00 ||
		e.Operation == StackPull && e.Before == 0xff
}
