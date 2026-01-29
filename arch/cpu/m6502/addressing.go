package m6502

// AddressingMode specifies how a 6502 instruction accesses its operands.
// Multiple modes can be combined using bitwise OR for instructions that support variants.
//
// Available modes:
//   - ImpliedAddressing: No operands (NOP, DEX, INX)
//   - AccumulatorAddressing: Operates on accumulator (ASL A, ROL A)
//   - ImmediateAddressing: Constant value (LDA #$10)
//   - AbsoluteAddressing: 16-bit memory address (JMP $1234)
//   - ZeroPageAddressing: 8-bit zero page address (LDA $10)
//   - AbsoluteXAddressing: Absolute address + X register (LDA $1234,X)
//   - ZeroPageXAddressing: Zero page address + X register (LDA $10,X)
//   - AbsoluteYAddressing: Absolute address + Y register (LDA $1234,Y)
//   - ZeroPageYAddressing: Zero page address + Y register (LDX $10,Y)
//   - IndirectAddressing: Indirect through address (JMP ($1234))
//   - IndirectXAddressing: Indexed indirect via X (LDA ($10,X))
//   - IndirectYAddressing: Indirect indexed via Y (LDA ($10),Y)
//   - RelativeAddressing: PC-relative offset for branches (BNE $10)
type AddressingMode int

const (
	NoAddressing      AddressingMode = 0
	ImpliedAddressing AddressingMode = 1 << iota
	AccumulatorAddressing
	ImmediateAddressing
	AbsoluteAddressing
	ZeroPageAddressing
	AbsoluteXAddressing
	ZeroPageXAddressing
	AbsoluteYAddressing
	ZeroPageYAddressing
	IndirectAddressing
	IndirectXAddressing
	IndirectYAddressing
	RelativeAddressing
)

// AccessMode specifies how memory is accessed during instruction execution.
// Used for memory-mapped I/O where read and write operations to the same
// address can have different meanings.
type AccessMode int

const (
	NoAccess        AccessMode = 0 // No memory access
	ReadAccess      AccessMode = 1 // Read-only access
	WriteAccess     AccessMode = 2 // Write-only access
	ReadWriteAccess AccessMode = 3 // Read-modify-write access
)

// AccessModeConstant specifies the access mode for a memory address.
// Used for memory-mapped I/O where the same address has different meanings
// depending on whether it is read from or written to (e.g., NES APU register $4017).
type AccessModeConstant struct {
	Constant string     // Named constant for the address
	Mode     AccessMode // Access mode (read, write, or read-modify-write)
}

// Indexed addressing type definitions.
// These types distinguish between different addressing modes at compile time.
type (
	AbsoluteX uint16 // Absolute + X register (e.g., LDA $1234,X)
	AbsoluteY uint16 // Absolute + Y register (e.g., LDA $1234,Y)
	IndirectX uint16 // Indexed indirect via X (e.g., LDA ($10,X))
	IndirectY uint16 // Indirect indexed via Y (e.g., LDA ($10),Y)
	ZeroPageX uint8  // Zero page + X register (e.g., LDA $10,X)
	ZeroPageY uint8  // Zero page + Y register (e.g., LDX $10,Y)
)

// ZeroPage represents an 8-bit address in the zero page ($00-$FF).
// Zero page addressing is faster than absolute addressing and uses one less byte.
type ZeroPage uint8

// Absolute represents a 16-bit absolute memory address ($0000-$FFFF).
// Used for instructions that access memory at a specific address (e.g., JMP $1234).
type Absolute uint16

// Indirect represents indirect addressing through a memory pointer.
// For JMP ($1234), the CPU reads the target address from memory locations $1234-$1235.
// For indexed indirect, the address is in zero page (e.g., LDA ($10,X) or LDA ($10),Y).
type Indirect uint16

// IndirectResolved represents the final resolved address from indirect addressing.
// Includes handling of the 6502 page-crossing bug where JMP ($12FF) reads from
// $12FF and $1200 instead of $12FF and $1300.
type IndirectResolved uint16

// Accumulator represents the accumulator register as an instruction operand.
// Used for instructions that operate on the accumulator (e.g., ASL A, ROL A).
type Accumulator int
