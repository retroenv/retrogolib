// Package m65816 provides a complete WDC 65C816 (65816) CPU emulation.
//
// The 65816 is a 16-bit successor to the 65C02, designed by Western Design Center.
// It is backward-compatible with the 65C02 in emulation mode and provides an
// expanded architecture in native mode:
//
//   - 24-bit address space (16 MB)
//   - 16-bit accumulator (C), index registers (X, Y), and stack pointer
//   - Mode-switchable register widths (M flag for accumulator, X flag for index)
//   - Direct Page register (replaces fixed zero page)
//   - Data Bank and Program Bank registers for 24-bit addressing
//   - Dual-mode operation: Emulation (65C02 compatible) and Native
//
// # Target Systems
//
//   - Super Nintendo Entertainment System (SNES/Super Famicom)
//   - Apple IIGS
//
// # Basic Usage
//
//	type myMem struct { data [16 * 1024 * 1024]byte }
//	func (m *myMem) ReadByte(addr uint32) uint8      { return m.data[addr&0xFFFFFF] }
//	func (m *myMem) WriteByte(addr uint32, v uint8)  { m.data[addr&0xFFFFFF] = v }
//	func (m *myMem) ReadWord(addr uint32) uint16 {
//	    lo := uint16(m.data[addr&0xFFFFFF])
//	    hi := uint16(m.data[(addr+1)&0xFFFFFF])
//	    return hi<<8 | lo
//	}
//	func (m *myMem) WriteWord(addr uint32, v uint16) {
//	    m.data[addr&0xFFFFFF] = uint8(v)
//	    m.data[(addr+1)&0xFFFFFF] = uint8(v >> 8)
//	}
//
//	mem := &myMem{}
//	cpu := m65816.New(mem)
//	for { if err := cpu.Step(); err != nil { break } }
//
// # Memory Layout
//
//   - Bank $00: Direct page (DP register), stack, zero page compatibility
//   - $00:FFE4-$00:FFE5: COP vector (native)
//   - $00:FFE6-$00:FFE7: BRK vector (native)
//   - $00:FFE8-$00:FFE9: ABORT vector (native)
//   - $00:FFEA-$00:FFEB: NMI vector (native)
//   - $00:FFEC-$00:FFED: Reserved (native)
//   - $00:FFEE-$00:FFEF: IRQ vector (native)
//   - $00:FFF4-$00:FFF5: COP vector (emulation)
//   - $00:FFF8-$00:FFF9: ABORT vector (emulation)
//   - $00:FFFA-$00:FFFB: NMI vector (emulation)
//   - $00:FFFC-$00:FFFD: RESET vector (emulation)
//   - $00:FFFE-$00:FFFF: IRQ/BRK vector (emulation)
//
// # Registers
//
//   - C: 16-bit accumulator (A = low byte, B = high byte)
//   - X, Y: 16-bit (or 8-bit when X flag set) index registers
//   - SP: 16-bit stack pointer
//   - DP: 16-bit Direct Page register
//   - DB: 8-bit Data Bank register
//   - PB: 8-bit Program Bank register
//   - PC: 16-bit program counter (within bank)
//   - P: Processor status (N V M X D I Z C)
//   - E: Emulation flag (toggled by XCE instruction)
//
// # Processor Modes
//
// Emulation Mode (E=1): Backward-compatible 65C02 behavior.
// 8-bit accumulator and index registers, stack fixed to page 1.
//
// Native Mode (E=0): Full 65816 capabilities.
// M flag controls accumulator width (M=0: 16-bit, M=1: 8-bit).
// X flag controls index register width (X=0: 16-bit, X=1: 8-bit).
package m65816
