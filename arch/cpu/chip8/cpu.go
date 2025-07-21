package chip8

import (
	"fmt"
	"sync"
)

type CPU struct {
	Memory [4096]byte // 4KB of memory

	V  [16]byte // 16 general-purpose registers V0-VF
	I  uint16   // Index register
	PC uint16   // Program counter

	Stack [16]uint16 // Call stack
	SP    uint8      // Stack pointer

	DelayTimer byte // Delay timer
	SoundTimer byte // Sound timer

	Key [16]bool // Hexadecimal keypad state

	Display      [displayWidth * displayHeight]byte // Monochrome display (64x32)
	RedrawScreen bool                               // Indicates if the screen needs to be redrawn

	mu sync.RWMutex // Mutex for thread-safe access
}

const (
	displayHeight         = 32
	displayWidth          = 64
	initialProgramCounter = 0x200
)

// New creates a new CPU.
func New() *CPU {
	c := &CPU{
		PC: initialProgramCounter,
	}

	// Load fontset into memory
	copy(c.Memory[:], fontSet[:])

	return c
}

// Step executes the next instruction in the CPU.
func (c *CPU) Step() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.PC >= uint16(len(c.Memory)-1) {
		return fmt.Errorf("%w: PC=0x%03X", ErrMemoryOutOfBounds, c.PC)
	}

	w := uint16(c.Memory[c.PC])<<8 | uint16(c.Memory[c.PC+1])
	idx := byte(w >> 12)
	opcodes := Opcodes[idx]

	for _, opcode := range opcodes {
		if opcode.Info.Mask&w == opcode.Info.Value {
			return opcode.Instruction.Emulation(c, w)
		}
	}

	return fmt.Errorf("unknown opcode: %04X", w)
}

// updatePC increments the program counter to the next instruction and optionally skips the following instruction.
func (c *CPU) updatePC(skipInstruction bool) {
	if skipInstruction {
		c.PC += 4
	} else {
		c.PC += 2
	}
}

// UpdateTimers decrements the delay and sound timers.
func (c *CPU) UpdateTimers() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.DelayTimer > 0 {
		c.DelayTimer--
	}
	if c.SoundTimer > 0 {
		c.SoundTimer--
	}
}

// Reset resets the CPU to its initial state.
func (c *CPU) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.PC = initialProgramCounter
	c.SP = 0
	c.I = 0
	c.DelayTimer = 0
	c.SoundTimer = 0
	c.RedrawScreen = false

	// Clear registers
	for i := range c.V {
		c.V[i] = 0
	}

	// Clear memory (except font data)
	for i := len(fontSet); i < len(c.Memory); i++ {
		c.Memory[i] = 0
	}

	// Clear display
	for i := range c.Display {
		c.Display[i] = 0
	}

	// Clear stack
	for i := range c.Stack {
		c.Stack[i] = 0
	}

	// Clear keys
	for i := range c.Key {
		c.Key[i] = false
	}
}

// GetState returns a copy of the CPU state for safe access.
func (c *CPU) GetState() CPUState {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var state CPUState
	copy(state.Memory[:], c.Memory[:])
	copy(state.V[:], c.V[:])
	copy(state.Stack[:], c.Stack[:])
	copy(state.Display[:], c.Display[:])
	copy(state.Key[:], c.Key[:])

	state.I = c.I
	state.PC = c.PC
	state.SP = c.SP
	state.DelayTimer = c.DelayTimer
	state.SoundTimer = c.SoundTimer
	state.RedrawScreen = c.RedrawScreen

	return state
}

// SetState sets the CPU state from a snapshot.
func (c *CPU) SetState(state CPUState) {
	c.mu.Lock()
	defer c.mu.Unlock()

	copy(c.Memory[:], state.Memory[:])
	copy(c.V[:], state.V[:])
	copy(c.Stack[:], state.Stack[:])
	copy(c.Display[:], state.Display[:])
	copy(c.Key[:], state.Key[:])

	c.I = state.I
	c.PC = state.PC
	c.SP = state.SP
	c.DelayTimer = state.DelayTimer
	c.SoundTimer = state.SoundTimer
	c.RedrawScreen = state.RedrawScreen
}

// CPUState represents a snapshot of the CPU state.
type CPUState struct {
	Memory       [4096]byte
	V            [16]byte
	Stack        [16]uint16
	Display      [displayWidth * displayHeight]byte
	Key          [16]bool
	I            uint16
	PC           uint16
	SP           uint8
	DelayTimer   byte
	SoundTimer   byte
	RedrawScreen bool
}
