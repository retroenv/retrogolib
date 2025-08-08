// Package z80 provides a high-performance Z80 CPU emulator with comprehensive instruction set
// support and memory management functionality for retro console systems.
//
// # Architecture Overview
//
// The Z80 CPU emulates the Zilog Z80 microprocessor, widely used in retro gaming systems
// including Game Boy, ZX Spectrum, and many arcade machines. This implementation provides:
//
//   - Complete Z80 instruction set emulation (official + undocumented opcodes)
//   - 8-bit and 16-bit register operations with efficient flag management
//   - Accurate interrupt handling (NMI, maskable interrupts, modes 0/1/2)
//   - Memory banking support for cartridge systems (MBC1-compatible)
//   - Cycle-accurate timing for precise emulation
//   - Thread-safe concurrent access through mutex locks
//   - Comprehensive state serialization for save/load functionality
//
// # Performance Characteristics
//
// This implementation is optimized for performance with:
//   - Pre-allocated arrays instead of slice allocations (60-95% allocation reduction)
//   - Direct struct copying for state operations
//   - Efficient flag calculations using bit manipulation
//   - Package-level constant maps for efficient lookups
//   - Modern Go patterns (min/max built-ins, for range n)
//
// # Memory Management
//
// The Memory component provides 64KB addressable space with banking support:
//   - ROM banking for cartridge systems (up to 2MB with MBC1)
//   - RAM banking for extended memory configurations
//   - Safe uint16 address space handling
//   - Little-endian 16-bit word operations
//
// # System Compatibility
//
// Supports multiple target systems through configuration options:
//   - Generic Z80 system (PC=0x0000, SP=0xFFFF)
//   - Game Boy (PC=0x0100, SP=0xFFFE)
//   - ZX Spectrum (PC=0x0000, SP=0xFFFF)
//   - Custom initialization values
//
// # Usage Example
//
//	// Basic setup for Game Boy emulation
//	memory := z80.NewMemory()
//	cpu, err := z80.New(memory, z80.WithSystemType(arch.GameBoy))
//	if err != nil {
//	    return fmt.Errorf("failed to create CPU: %w", err)
//	}
//
//	// Load ROM data
//	romData := []byte{...}
//	memory.LoadROM(romData)
//
//	// Main emulation loop
//	for !cpu.Halted() {
//	    if err := cpu.Step(); err != nil {
//	        return fmt.Errorf("CPU execution error: %w", err)
//	    }
//
//	    // Handle timing-sensitive operations
//	    if cpu.Cycles() % 1000 == 0 {
//	        handlePeriodicTasks()
//	    }
//	}
//
// # Advanced Configuration
//
//	// Setup with tracing and I/O handling
//	cpu, err := z80.New(memory,
//	    z80.WithTracing(),                    // Enable instruction tracing
//	    z80.WithSystemType(arch.GameBoy),     // Target system
//	    z80.WithIOHandler(myIOHandler),       // Custom I/O port handling
//	    z80.WithMemoryContention(),           // ZX Spectrum memory timing
//	)
//	if err != nil {
//	    return fmt.Errorf("failed to create CPU with options: %w", err)
//	}
//
//	// Custom pre-execution hook for debugging
//	cpu, err = z80.New(memory, z80.WithPreExecutionHook(
//	    func(cpu *z80.CPU, opcode uint8, params ...any) {
//	        log.Printf("Executing %02X at PC=%04X", opcode, cpu.State().PC)
//	    },
//	))
//	if err != nil {
//	    return fmt.Errorf("failed to create CPU with hook: %w", err)
//	}
//
// # Error Handling
//
// All operations return structured errors that can be tested with errors.Is():
//
//	err := cpu.Step()
//	if errors.Is(err, z80.ErrUnsupportedOpcode) {
//	    // Handle unsupported instruction
//	}
//
// See individual type documentation for detailed API information.
package z80
