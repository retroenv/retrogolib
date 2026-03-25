// Package m6809 provides a complete Motorola 6809 CPU emulation.
//
// The 6809 is an 8-bit microprocessor designed by Motorola, featuring an
// advanced instruction set with position-independent code support, hardware
// multiply, and a comprehensive indexed addressing mode system.
//
// # Target Systems
//
//   - TRS-80 Color Computer (CoCo)
//   - Vectrex
//   - Dragon 32/64
//   - Williams arcade hardware (Defender, Robotron, Joust)
//
// # Basic Usage
//
//	type myMem struct { data [65536]byte }
//	func (m *myMem) Read(addr uint16) uint8     { return m.data[addr] }
//	func (m *myMem) Write(addr uint16, v uint8) { m.data[addr] = v }
//	func (m *myMem) ReadWord(addr uint16) uint16 {
//	    return uint16(m.data[addr])<<8 | uint16(m.data[addr+1])
//	}
//	func (m *myMem) WriteWord(addr uint16, v uint16) {
//	    m.data[addr] = uint8(v >> 8)
//	    m.data[addr+1] = uint8(v)
//	}
//
//	mem := &myMem{}
//	wrapped, _ := m6809.NewMemory(mem)
//	cpu, _ := m6809.New(wrapped)
//	for { if err := cpu.Step(); err != nil { break } }
//
// # Memory Layout
//
//   - $0000-$7FFF: RAM (system dependent)
//   - $8000-$FEFF: ROM (system dependent)
//   - $FFF0-$FFF1: Reserved
//   - $FFF2-$FFF3: SWI3 vector
//   - $FFF4-$FFF5: SWI2 vector
//   - $FFF6-$FFF7: FIRQ vector
//   - $FFF8-$FFF9: IRQ vector
//   - $FFFA-$FFFB: SWI vector
//   - $FFFC-$FFFD: NMI vector
//   - $FFFE-$FFFF: RESET vector
//
// # Registers
//
//   - A, B: 8-bit accumulators (combine as 16-bit D register, A=high, B=low)
//   - X, Y: 16-bit index registers
//   - U: 16-bit user stack pointer
//   - S: 16-bit system stack pointer
//   - DP: 8-bit direct page register
//   - PC: 16-bit program counter
//   - CC: Condition code register (E F H I N Z V C)
//
// # Addressing Modes
//
// The 6809 supports implied, immediate (8/16-bit), direct page, extended (16-bit),
// indexed (with 5-bit/8-bit/16-bit offsets, accumulator offsets, auto-increment/
// decrement, PC-relative, and indirect variants), relative (8/16-bit), and
// register-to-register (TFR/EXG) addressing modes.
//
// # Prefix Opcodes
//
// The 6809 uses two prefix bytes to extend the instruction set:
//   - $10: Page 2 instructions (long branches, CMPD, CMPY, LDY, STY, LDS, STS, SWI2)
//   - $11: Page 3 instructions (CMPU, CMPS, SWI3)
package m6809
