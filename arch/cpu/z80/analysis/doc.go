// Package analysis computes local Z80 load, logical, compare, and carry effects.
//
// Views bind byte registers and individual flags to pair storage. Main and shadow
// banks are distinct. Writes preserve sibling bits. Copies describe input values;
// equal partial masks do not prove register identity. Exchanges and other forms
// outside the pilot fail closed.
//
// Before precedes opcode fetch. Each supported transfer increments the low seven
// bits of R once, or twice for a prefix, and preserves R bit seven. LD A,R uses
// that incremented value. Q is the opaque byte latch used by SCF/CCF for X/Y.
// A caller can supply its known bits; each transfer conservatively clobbers Q.
// The package does not infer instruction history. IFF2 is a separate one-bit
// input for LD A,I and LD A,R parity/overflow.
//
// The declared data scope contains the pair, index, special, flag, and mode
// resources below. PC, cycles, interrupt acceptance, and internal WZ are outside
// this scope. Flow represents PC control; callers own interrupt and call policy.
// Logical addresses require caller-owned physical mapping and memory stability.
// Events exclude opcode fetches and bus timing. Result must be nonnil and its
// buffers belong to the caller.
package analysis
