// Package analysis computes local 6502-family data and status effects.
//
// Transfer uses the canonical variant opcode tables. It models A, X, Y, SP,
// and the status byte. PC and instruction fetches are represented by Flow,
// not register or memory writes. Calls describe only their local stack effects;
// callers must apply a callee summary or a conservative barrier.
//
// Addresses are logical 16-bit CPU addresses. A 6507 board's address masking,
// 6510 port mapping, NES mirrors, banks, and volatile memory remain caller policy.
// Read must prove a stable value in the current context. Unknown addresses do
// not invoke Read. Ordered events retain all semantic reads and writes, including
// pointer and stack accesses, but exclude dummy bus accesses and instruction fetches.
//
// Binary ADC and SBC are supported. NES ignores decimal mode. Other variants
// require a proved clear D flag; otherwise affected facts become unknown and
// Reason reports the refused mode. Unsupported instructions fail closed.
package analysis
