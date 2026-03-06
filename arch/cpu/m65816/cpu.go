package m65816

import (
	"errors"
	"fmt"
	"sync"
)

// State represents a complete snapshot of the 65816 CPU state.
type State struct {
	C  uint16 // Accumulator (full 16-bit C register)
	X  uint16 // Index X
	Y  uint16 // Index Y
	SP uint16 // Stack pointer
	DP uint16 // Direct Page register
	DB uint8  // Data Bank register
	PB uint8  // Program Bank register
	PC uint16 // Program counter (within bank)
	P  uint8  // Processor status (from Flags.Get())
	E  bool   // Emulation flag

	Cycles uint64
}

// CPU represents a thread-safe WDC 65C816 microprocessor.
type CPU struct {
	mu sync.RWMutex

	// Registers
	C  uint16 // Accumulator (full 16-bit; A = low byte, B = high byte)
	X  uint16 // Index register X
	Y  uint16 // Index register Y
	SP uint16 // Stack pointer
	DP uint16 // Direct Page register
	DB uint8  // Data Bank register
	PB uint8  // Program Bank register
	PC uint16 // Program counter (within current program bank)

	Flags Flags // Processor status register (P)
	E     bool  // Emulation flag (toggled by XCE)

	cycles    uint64
	stopped   bool // STP instruction state
	waiting   bool // WAI instruction state
	pcChanged bool // set by instructions that explicitly set PC (branches, jumps)

	// Interrupt control
	triggerNMI bool
	triggerIRQ bool
	nmiRunning bool
	irqRunning bool

	memory *Memory
	opts   Options

	TraceStep TraceStep // set when tracing is enabled
}

// TraceStep holds information for instruction tracing.
type TraceStep struct {
	PC             uint16
	PB             uint8
	OpcodeOperands []byte
	Opcode         Opcode
	PageCrossed    bool
}

const (
	initialCycles = 7
	// Initial SP is $01FF in emulation mode (stack page 1).
	initialSP = 0x01FF
)

// New creates a new 65816 CPU, reads the reset vector, and initializes registers.
func New(memory *Memory, opts ...Option) (*CPU, error) {
	if memory == nil {
		return nil, errors.New("memory cannot be nil")
	}

	c := &CPU{
		SP:     initialSP,
		cycles: initialCycles,
		memory: memory,
		opts:   NewOptions(opts...),
		E:      true, // Start in emulation mode
	}

	// Force M=1, X=1 in emulation mode
	c.Flags.M = 1
	c.Flags.X = 1
	c.Flags.I = 1

	// Read reset vector (emulation mode vector at $FFFC)
	resetVec := memory.ReadVector(VectorEmuRESET)
	c.PC = resetVec

	return c, nil
}

// A returns the low byte of the accumulator (8-bit mode or low half of 16-bit).
func (c *CPU) A() uint8 { return uint8(c.C) }

// B returns the high byte of the accumulator.
func (c *CPU) B() uint8 { return uint8(c.C >> 8) }

// FullPC returns the 24-bit effective program address (PB:PC).
func (c *CPU) FullPC() uint32 {
	return bank24(c.PB, c.PC)
}

// AccWidth returns the current accumulator width in bytes (1 or 2).
func (c *CPU) AccWidth() int {
	if c.E || c.Flags.M != 0 {
		return 1
	}
	return 2
}

// IdxWidth returns the current index register width in bytes (1 or 2).
func (c *CPU) IdxWidth() int {
	if c.E || c.Flags.X != 0 {
		return 1
	}
	return 2
}

// Cycles returns the total number of cycles executed.
func (c *CPU) Cycles() uint64 { return c.cycles }

// State returns a snapshot of the current CPU state.
func (c *CPU) State() State {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return State{
		C:      c.C,
		X:      c.X,
		Y:      c.Y,
		SP:     c.SP,
		DP:     c.DP,
		DB:     c.DB,
		PB:     c.PB,
		PC:     c.PC,
		P:      c.Flags.Get(),
		E:      c.E,
		Cycles: c.cycles,
	}
}

// ValidateState checks that CPU state is consistent.
func (c *CPU) ValidateState() error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.memory == nil {
		return errors.New("CPU memory is nil")
	}
	return nil
}

// Reset resets the CPU to its initial post-reset state.
func (c *CPU) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.C = 0
	c.X = 0
	c.Y = 0
	c.SP = initialSP
	c.DP = 0
	c.DB = 0
	c.PB = 0
	c.E = true
	c.Flags = Flags{M: 1, X: 1, I: 1}
	c.cycles = initialCycles
	c.stopped = false
	c.waiting = false
	c.triggerNMI = false
	c.triggerIRQ = false
	c.nmiRunning = false
	c.irqRunning = false

	if c.memory != nil {
		c.PC = c.memory.ReadVector(VectorEmuRESET)
	}
}

// Memory returns the CPU's memory.
func (c *CPU) Memory() *Memory { return c.memory }

// push8 pushes a byte onto the stack and decrements SP.
// In emulation mode, SP wraps within page 1 ($0100-$01FF) after each byte.
func (c *CPU) push8(value uint8) {
	c.memory.WriteByte(bank24(0, c.SP), value)
	c.SP--
	if c.E {
		c.SP = 0x0100 | (c.SP & 0x00FF)
	}
}

// push8raw pushes a byte without page-1 wrap, for 65816-native stack instructions
// that use the full 16-bit SP even in emulation mode.
func (c *CPU) push8raw(value uint8) {
	c.memory.WriteByte(bank24(0, c.SP), value)
	c.SP--
}

// push16 pushes a 16-bit word onto the stack (high byte first).
func (c *CPU) push16(value uint16) {
	c.push8(uint8(value >> 8))
	c.push8(uint8(value))
}

// push16raw pushes a 16-bit word without per-byte page-1 wrap.
func (c *CPU) push16raw(value uint16) {
	c.push8raw(uint8(value >> 8))
	c.push8raw(uint8(value))
}

// pop16raw pops a 16-bit word without per-byte page-1 wrap.
func (c *CPU) pop16raw() uint16 {
	lo := uint16(c.pop8raw())
	hi := uint16(c.pop8raw())
	return hi<<8 | lo
}

// fixEmuSP normalises SP to page 1 after a 65816-native stack instruction.
func (c *CPU) fixEmuSP() {
	if c.E {
		c.SP = 0x0100 | (c.SP & 0x00FF)
	}
}

// push24 pushes a 24-bit value onto the stack (high byte first).
func (c *CPU) push24(value uint32) {
	c.push8(uint8(value >> 16))
	c.push8(uint8(value >> 8))
	c.push8(uint8(value))
}

// pop8 pops a byte from the stack and increments SP.
// In emulation mode, SP wraps within page 1 ($0100-$01FF) after each byte.
func (c *CPU) pop8() uint8 {
	c.SP++
	if c.E {
		c.SP = 0x0100 | (c.SP & 0x00FF)
	}
	return c.memory.ReadByte(bank24(0, c.SP))
}

// pop8raw pops a byte without page-1 wrap, for 65816-native stack instructions
// that use the full 16-bit SP even in emulation mode.
func (c *CPU) pop8raw() uint8 {
	c.SP++
	return c.memory.ReadByte(bank24(0, c.SP))
}

// pop16 pops a 16-bit word from the stack (low byte first).
func (c *CPU) pop16() uint16 {
	lo := uint16(c.pop8())
	hi := uint16(c.pop8())
	return hi<<8 | lo
}

// pop24 pops a 24-bit value from the stack.
func (c *CPU) pop24() uint32 {
	lo := uint32(c.pop8())
	mid := uint32(c.pop8())
	hi := uint32(c.pop8())
	return hi<<16 | mid<<8 | lo
}

// dataAddr forms a 24-bit data address using the Data Bank register.
func (c *CPU) dataAddr(offset uint16) uint32 {
	return bank24(c.DB, offset)
}

// dpAddr forms a 24-bit direct page address.
// In emulation mode with DP=$0000, it wraps within page 0.
func (c *CPU) dpAddr(offset uint8) uint32 {
	return bank24(0, c.DP+uint16(offset))
}

// readMem8 reads a byte from a 24-bit address.
func (c *CPU) readMem8(addr uint32) uint8 {
	return c.memory.ReadByte(addr & 0xFFFFFF)
}

// writeMem8 writes a byte to a 24-bit address.
func (c *CPU) writeMem8(addr uint32, value uint8) {
	c.memory.WriteByte(addr&0xFFFFFF, value)
}

// readMem16 reads a 16-bit word (little-endian) from a 24-bit address.
// The hi byte wraps within the same 64KB bank — used for pointer fetches
// (DP indirect, absolute indirect, stack-relative indirect) where the 65816
// keeps pointer bytes inside the bank containing the pointer itself.
func (c *CPU) readMem16(addr uint32) uint16 {
	addr &= 0xFFFFFF
	lo := uint16(c.memory.ReadByte(addr))
	bank := addr & 0xFF0000
	hi := uint16(c.memory.ReadByte(bank | uint32(uint16(addr)+1)))
	return hi<<8 | lo
}

// readData16 reads a 16-bit word (little-endian) from a 24-bit address.
// The hi byte is at addr+1 in full 24-bit address space — used for data reads
// (absolute, indexed-absolute, direct-page memory operands) where the 65816
// allows the address to cross a bank boundary.
func (c *CPU) readData16(addr uint32) uint16 {
	addr &= 0xFFFFFF
	lo := uint16(c.memory.ReadByte(addr))
	hi := uint16(c.memory.ReadByte((addr + 1) & 0xFFFFFF))
	return hi<<8 | lo
}

// readDPWord reads a 16-bit indirect pointer from a direct-page offset.
// In emulation mode with DP_lo=0, both pointer bytes wrap within the DP 256-byte page.
func (c *CPU) readDPWord(dpOffset uint8) uint16 {
	if c.E && c.DP&0xFF == 0 {
		dpPage := uint32(c.DP)
		lo := uint16(c.memory.ReadByte(dpPage | uint32(dpOffset)))
		hi := uint16(c.memory.ReadByte(dpPage | uint32(dpOffset+1))) // +1 wraps at 8 bits
		return hi<<8 | lo
	}
	return c.readMem16(bank24(0, c.DP+uint16(dpOffset)))
}

// writeMem16 writes a 16-bit word (little-endian) to a 24-bit address.
func (c *CPU) writeMem16(addr uint32, value uint16) {
	c.memory.WriteWord(addr&0xFFFFFF, value)
}

// readMem24 reads a 24-bit (3-byte) long pointer, wrapping within the same bank.
// All three bytes stay in the same 64KB bank as addr — used for [dp] and [abs] pointer fetches.
func (c *CPU) readMem24(addr uint32) uint32 {
	addr &= 0xFFFFFF
	bank := addr & 0xFF0000
	lo := uint32(c.memory.ReadByte(addr))
	mid := uint32(c.memory.ReadByte(bank | uint32(uint16(addr)+1)))
	hi := uint32(c.memory.ReadByte(bank | uint32(uint16(addr)+2)))
	return hi<<16 | mid<<8 | lo
}

// branch performs a relative branch if the condition is true.
// addr is the pre-computed absolute branch target address.
func (c *CPU) branch(taken bool, addr uint16) {
	if !taken {
		return
	}
	c.PC = addr
	c.pcChanged = true
	c.cycles++ // extra cycle when branch taken
}

// GetP returns the current processor status byte.
func (c *CPU) GetP() uint8 {
	return c.Flags.Get()
}

// SetP sets the processor status from a byte, handling E-mode constraints.
func (c *CPU) SetP(p uint8) {
	c.Flags.Set(p)
	if c.E {
		// Emulation mode forces M=1, X=1
		c.Flags.M = 1
		c.Flags.X = 1
	} else if c.Flags.X != 0 {
		// X flag transition to 8-bit: zero high bytes of X and Y
		c.X &= 0x00FF
		c.Y &= 0x00FF
	}
}

// instrSize returns the actual size of an instruction given the opcode's WidthFlag
// and current M/X flag state.
func (c *CPU) instrSize(op Opcode) int {
	size := int(op.Instruction.Addressing[op.Addressing].BaseSize)
	switch op.WidthFlag {
	case WidthM:
		if c.AccWidth() == 2 {
			size++
		}
	case WidthX:
		if c.IdxWidth() == 2 {
			size++
		}
	}
	return size
}

// fetchByte reads the next byte from PC (in PB bank) without advancing PC.
func (c *CPU) fetchByte(offset uint16) uint8 {
	return c.memory.ReadByte(bank24(c.PB, c.PC+offset))
}

// fetchWord reads the next word starting at PC+offset.
func (c *CPU) fetchWord(offset uint16) uint16 {
	lo := uint16(c.fetchByte(offset))
	hi := uint16(c.fetchByte(offset + 1))
	return hi<<8 | lo
}

// fetchLong reads 3 bytes from PC+offset.
func (c *CPU) fetchLong(offset uint16) uint32 {
	b0 := uint32(c.fetchByte(offset))
	b1 := uint32(c.fetchByte(offset + 1))
	b2 := uint32(c.fetchByte(offset + 2))
	return b2<<16 | b1<<8 | b0
}

// resolveDP resolves a direct page address to a 24-bit address.
// Handles emulation mode page-0 wrap when DP low byte is $00.
func (c *CPU) resolveDP(dp uint8) uint32 {
	if c.E && c.DP&0xFF == 0 {
		// Emulation mode with DP page-aligned: result stays within the DP 256-byte block.
		return uint32(c.DP) | uint32(dp)
	}
	// Direct page is always in bank 0; wrap at 16 bits if DP+dp overflows.
	return bank24(0, c.DP+uint16(dp))
}

// resolveEA resolves an effective address parameter to a readable 24-bit address.
// Returns the address and an error if the type is unexpected.
func (c *CPU) resolveEA(param any) (uint32, error) {
	switch p := param.(type) {
	case Immediate8:
		return 0, fmt.Errorf("resolveEA called on immediate8")
	case Immediate16:
		return 0, fmt.Errorf("resolveEA called on immediate16")
	case DirectPage:
		return c.resolveDP(uint8(p)), nil
	case DirectPageX:
		dp := uint8(p)
		if c.E && c.DP&0xFF == 0 {
			// Emulation mode, DP page-aligned: (dp+X) wraps within page 0
			return uint32(c.DP) | uint32(dp+uint8(c.X&0xFF)), nil
		}
		if c.IdxWidth() == 1 {
			return (uint32(c.DP) + uint32(dp) + uint32(c.X&0xFF)) & 0xFFFF, nil
		}
		return (uint32(c.DP) + uint32(dp) + uint32(c.X)) & 0xFFFF, nil
	case DirectPageY:
		dp := uint8(p)
		if c.E && c.DP&0xFF == 0 {
			// Emulation mode, DP page-aligned: (dp+Y) wraps within page 0
			return uint32(c.DP) | uint32(dp+uint8(c.Y&0xFF)), nil
		}
		if c.IdxWidth() == 1 {
			return (uint32(c.DP) + uint32(dp) + uint32(c.Y&0xFF)) & 0xFFFF, nil
		}
		return (uint32(c.DP) + uint32(dp) + uint32(c.Y)) & 0xFFFF, nil
	case DPIndirect:
		return uint32(p), nil
	case DPIndirectX:
		return uint32(p), nil
	case DPIndirectY:
		return uint32(p), nil
	case DPIndirectLong:
		return uint32(p), nil
	case DPIndLongY:
		return uint32(p), nil
	case Absolute16:
		return c.dataAddr(uint16(p)), nil
	case AbsoluteX16:
		return uint32(p), nil
	case AbsoluteY16:
		return uint32(p), nil
	case AbsLong:
		return uint32(p), nil
	case AbsLongX:
		return uint32(p), nil
	case StackRel:
		return bank24(0, c.SP+uint16(p)), nil
	case SRIndY:
		return uint32(p), nil
	default:
		return 0, fmt.Errorf("%w: type %T", ErrUnsupportedAddressingMode, param)
	}
}
