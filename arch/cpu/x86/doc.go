// Package x86 provides Intel x86 (8086/8088) CPU emulation for DOS development.
//
// This package implements the Intel 8086/8088 instruction set commonly used
// in DOS development, including the most frequently used instructions for
// retro computing applications.
//
// Features:
//   - Complete 8086/8088 instruction set emulation
//   - Real mode memory addressing (segmented memory)
//   - Interrupt handling (hardware and software)
//   - Flag register management
//   - Thread-safe CPU state management
//   - Cycle-accurate timing
//   - State serialization for save/restore
//
// The implementation focuses on DOS-era compatibility and includes
// approximately 585 core instructions commonly used in DOS development.
//
// Example usage:
//
//	memory := x86.NewMemory(1024 * 1024) // 1MB memory
//	cpu, err := x86.New(memory)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Load program and execute
//	cpu.SetCS(0x1000)
//	cpu.SetIP(0x0000)
//	err = cpu.Step()
package x86
