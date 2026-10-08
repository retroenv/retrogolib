// Package cp1610 provides instruction metadata and execution for the General
// Instrument CP1610 processor used in the Intellivision.
//
// The CPU reads 16-bit words from a 65,536-word address space. Instructions use
// the low 10 bits of each opcode word. R7 is the program counter, and R6 is the
// stack pointer. The caller supplies memory and controls video and audio devices.
// The CPU does not model bus wait states or undocumented SDBD combinations.
package cp1610
