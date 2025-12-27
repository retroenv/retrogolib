// Package x86 provides Intel x86 real mode CPU support for retro computing tooling.
//
// This package implements the Intel x86 CPU architecture in real mode,
// covering instruction sets from 8086/8088 through 80486, including register
// management and comprehensive instruction definitions for static analysis,
// disassembly, and assembler development.
//
// # Supported CPU Generations
//
//   - 8086/8088 (1978): Base instruction set, segmented memory, 16-bit operations
//   - 80186/80188 (1982): Enhanced instructions (PUSHA/POPA, ENTER/LEAVE, BOUND, string I/O)
//   - 80286 (1982): Real mode enhancements (SMSW/LMSW for machine status)
//   - 80386 (1985): Bit manipulation (BSF/BSR/BT/BTC/BTR/BTS), move extensions (MOVZX/MOVSX), double-precision shifts
//   - 80486 (1989): Atomic operations (CMPXCHG/XADD), byte swap (BSWAP), cache control
//
// # Architecture Features
//
//   - Complete real mode instruction set (256 single-byte + two-byte opcodes with 0x0F prefix)
//   - Real mode memory addressing (segmented memory, 1MB address space)
//   - Comprehensive opcode tables with timing and size information
//   - ModR/M byte support for complex addressing modes
//   - Interrupt and flag register management
//   - Thread-safe CPU state access
//   - State serialization for save/restore
//
// The implementation focuses on static analysis and tooling rather than runtime
// emulation, making it ideal for assemblers, disassemblers, and code analysis tools.
//
// Example usage:
//
//	memory := x86.NewMemory(1024 * 1024) // 1MB memory
//	cpu, err := x86.New(memory)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Configure CPU state for analysis
//	cpu.SetCS(0x1000)
//	cpu.SetIP(0x0000)
//	state := cpu.State()
package x86
