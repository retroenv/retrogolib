// Package m68000 provides a Motorola 68000 CPU emulator with comprehensive instruction set
// support and memory management functionality for retro computing systems.
//
// # Architecture Overview
//
// The 68000 is a 32-bit CISC processor with a 16-bit external data bus and 24-bit
// address bus (16MB addressable space). This implementation provides:
//
//   - Complete 68000 instruction set emulation (~68 mnemonics)
//   - Hierarchical line-based opcode decoder (16 lines from top 4 bits)
//   - 14 addressing modes with effective address resolution
//   - 8 data registers (D0-D7, 32-bit) and 8 address registers (A0-A7)
//   - Dual stack pointers (USP/SSP) for user/supervisor modes
//   - 16-bit status register with CCR and system byte
//   - 256-vector exception model
//   - Big-endian byte order
//   - Cycle-accurate timing for precise emulation
//   - Thread-safe concurrent access through mutex locks
//
// # Usage Example
//
//	mem := m68000.NewBasicMemory()
//	bus := m68000.NewBasicBus(mem)
//	cpu, err := m68000.New(bus)
//	if err != nil {
//	    return fmt.Errorf("creating CPU: %w", err)
//	}
//
//	for !cpu.Halted() {
//	    if err := cpu.Step(); err != nil {
//	        return fmt.Errorf("CPU execution error: %w", err)
//	    }
//	}
package m68000
