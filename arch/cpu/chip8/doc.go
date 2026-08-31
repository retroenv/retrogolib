// Package chip8 provides a complete Chip-8 virtual machine implementation.
//
// The Chip-8 is a simple, interpreted programming language that was first used on
// some microcomputers in the mid-1970s. This package implements a full Chip-8
// virtual machine including:
//
//   - CPU with 16 general-purpose registers (V0-VF)
//   - 4KB of memory
//   - 64x32 monochrome display
//   - 16-key hexadecimal keypad
//   - Sound and delay timers
//   - Built-in font set for hexadecimal digits
//
// # Basic Usage
//
//	cpu := chip8.New()
//
//	// Load program into memory starting at 0x200
//	copy(cpu.Memory[0x200:], program)
//
//	// Execute instructions
//	for {
//		if err := cpu.Step(); err != nil {
//			log.Fatal(err)
//		}
//
//		// Handle display updates and input...
//	}
//
// Call UpdateTimers once per 60 Hz tick. Besides decrementing the delay and
// sound timers, that tick releases the original interpreter's display wait.
//
// # Memory Layout
//
//   - 0x000-0x1FF: Reserved (interpreter and font data)
//   - 0x200-0xFFF: Program memory (3584 bytes)
//
// # Registers
//
//   - V0-VE: General-purpose registers
//   - VF: Flag register (used for carry, borrow, collision detection)
//   - I: 16-bit index register, normally used with 12-bit memory addresses
//   - PC: Program counter
//   - SP: Stack pointer
//
// # Display
//
// The display is 64x32 pixels, monochrome. Drawing is performed using XOR
// operations, and collision detection sets the VF flag when sprites overlap
// existing pixels.
//
// # Compatibility
//
// New defaults to COSMAC VIP CHIP-8 behavior. WithQuirks enables individual
// behaviors used by later interpreters without silently mixing dialects.
//
// # Safety
//
// This implementation validates instruction fetches, memory-backed operands,
// stack access, key indices, and font indices before use.
package chip8
