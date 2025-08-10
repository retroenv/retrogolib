// Package x86 provides Intel x86 (8086/8088) CPU support for DOS development tooling.
//
// This package implements the Intel 8086/8088 CPU architecture commonly used
// in DOS development, including register management and instruction definitions
// for retro computing applications.
//
// Features:
//   - Complete 8086/8088 CPU state management
//   - Real mode memory addressing (segmented memory)
//   - Interrupt handling support
//   - Flag register management
//   - Thread-safe CPU state access
//   - State serialization for save/restore
//
// The implementation focuses on DOS-era compatibility for static analysis,
// disassembly, and development tooling applications.
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
