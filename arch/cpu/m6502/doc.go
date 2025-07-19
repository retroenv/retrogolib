// Package m6502 provides a complete MOS Technology 6502 CPU emulation.
//
// The 6502 is an 8-bit microprocessor that was widely used in home computers,
// video game consoles, and other systems in the 1970s and 1980s. This package
// implements a cycle-accurate emulation including:
//
//   - Full instruction set with all addressing modes
//   - Accurate flag handling (N, V, Z, C, I, D, B)
//   - Stack operations and interrupt handling
//   - Memory management with configurable backends
//   - Debugging and tracing capabilities
//
// # Basic Usage
//
//	memory := m6502.NewMemory()
//	cpu := m6502.New(memory)
//
//	// Load program
//	memory.Write(0x8000, 0xA9) // LDA #$42
//	memory.Write(0x8001, 0x42)
//
//	// Set reset vector
//	memory.WriteWord(m6502.ResetAddress, 0x8000)
//
//	// Execute instructions
//	for {
//		if err := cpu.Step(); err != nil {
//			log.Fatal(err)
//		}
//	}
//
// # Memory Layout
//
//   - 0x0000-0x00FF: Zero page (fast access)
//   - 0x0100-0x01FF: Stack
//   - 0x0200-0xFFEF: General memory
//   - 0xFFFA-0xFFFB: NMI vector
//   - 0xFFFC-0xFFFD: Reset vector
//   - 0xFFFE-0xFFFF: IRQ/BRK vector
//
// # Registers
//
//   - A: Accumulator (8-bit)
//   - X, Y: Index registers (8-bit)
//   - PC: Program counter (16-bit)
//   - SP: Stack pointer (8-bit, points into 0x0100-0x01FF)
//   - P: Processor status flags (8-bit)
//
// # Addressing Modes
//
// The 6502 supports various addressing modes:
//   - Immediate: #$42
//   - Zero page: $42
//   - Zero page,X: $42,X
//   - Absolute: $1234
//   - Absolute,X: $1234,X
//   - Absolute,Y: $1234,Y
//   - Indirect: ($1234)
//   - Indexed indirect: ($42,X)
//   - Indirect indexed: ($42),Y
//
// # Thread Safety
//
// CPU operations are protected by a read-write mutex, allowing concurrent
// read access to CPU state while ensuring exclusive access for modifications.
//
// # Accuracy
//
// This implementation includes cycle-accurate timing and historically accurate
// behavior, including the famous JMP ($xxFF) page boundary bug for maximum
// compatibility with original 6502 software.
package m6502
