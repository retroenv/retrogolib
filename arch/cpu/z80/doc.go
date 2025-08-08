// Package z80 provides a Z80 CPU emulator with support for the full instruction set
// and memory management functionality needed for retro console systems like the Game Boy.
//
// The Z80 CPU emulates the Zilog Z80 microprocessor, which is widely used in
// retro gaming systems and computers. This implementation provides:
//
//   - Full Z80 instruction set emulation
//   - 16-bit and 8-bit register management
//   - Flag register operations
//   - Stack operations
//   - Interrupt handling
//   - Memory management with banking support
//   - Cycle-accurate timing
//
// The CPU state can be saved and restored for debugging and save state functionality.
// Thread-safe access is provided through mutex locks.
//
// Example usage:
//
//	memory := z80.NewMemory()
//	cpu := z80.New(memory)
//
//	// Execute instructions
//	for !cpu.Halted() {
//	    err := cpu.Step()
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	}
package z80
