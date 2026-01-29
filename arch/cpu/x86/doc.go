// Package x86 provides Intel x86 instruction definitions for static analysis.
//
// This package provides instruction set definitions from 8086 through 80486,
// with opcode tables, addressing modes, and instruction metadata for
// disassemblers and assemblers.
//
// # Supported Generations
//
//   - 8086/8088 (1978): Base 16-bit instruction set
//   - 80186 (1982): PUSHA/POPA, ENTER/LEAVE, BOUND, INS/OUTS
//   - 80286 (1982): SMSW/LMSW
//   - 80386 (1985): BSF/BSR, BT/BTC/BTR/BTS, MOVZX/MOVSX, SHLD/SHRD
//   - 80486 (1989): CMPXCHG, XADD, BSWAP, INVD/WBINVD
//
// # Usage
//
//	// Look up by opcode
//	op := x86.Opcodes[0x90]
//	fmt.Println(op.Instruction.Name) // "nop"
//
//	// Look up by name
//	inst := x86.Instructions["mov"]
//
//	// Parse ModR/M byte
//	var modrm x86.ModRM
//	modrm.FromByte(0xC0)
package x86
