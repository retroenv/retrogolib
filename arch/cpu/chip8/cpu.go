package chip8

import (
	"fmt"
	"sync"
)

const (
	displayHeight         = 32
	displayWidth          = 64
	initialProgramCounter = 0x200
)

type keyWait struct {
	register uint16
	key      int8
	active   bool
}

// State represents complete Chip-8 VM state for save/load and debugging.
type State struct {
	Memory       [4096]byte                         // Full 4KB memory
	V            [16]byte                           // General-purpose registers V0-VF
	Stack        [16]uint16                         // Call stack
	Display      [displayWidth * displayHeight]byte // Display buffer
	Key          [16]bool                           // Keypad state
	I            uint16                             // Index register
	PC           uint16                             // Program counter
	SP           uint8                              // Stack pointer
	DelayTimer   byte                               // Delay timer value
	SoundTimer   byte                               // Sound timer value
	RedrawScreen bool                               // Screen redraw flag

	keyWait       keyWait
	drewThisFrame bool
}

// CPU represents a CHIP-8 virtual machine with serialized execution and snapshots.
type CPU struct {
	// Memory and registers
	Memory [4096]byte // 4KB memory ($000-$FFF)
	V      [16]byte   // 16 general-purpose registers (V0-VF, VF used as flag)
	I      uint16     // Index register
	PC     uint16     // Program counter

	// Stack for subroutine calls
	Stack [16]uint16 // Call stack (16 levels deep)
	SP    uint8      // Stack pointer

	// Timers (count down at 60Hz when non-zero)
	DelayTimer byte // Delay timer for timing events
	SoundTimer byte // Sound timer (beep when non-zero)

	// Input
	Key [16]bool // Hexadecimal keypad state (0-F)

	// Display
	Display      [displayWidth * displayHeight]byte // 64x32 monochrome display
	RedrawScreen bool                               // Set when screen needs redraw

	quirks        Quirks
	keyWait       keyWait
	drewThisFrame bool
	mu            sync.RWMutex // Thread-safe access protection
}

// New creates a new CPU.
func New(optionList ...Option) *CPU {
	opts := newOptions(optionList...)
	c := &CPU{
		PC:      initialProgramCounter,
		quirks:  opts.quirks,
		keyWait: newKeyWait(),
	}

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

// UpdateTimers advances the 60 Hz timer and display cadence by one tick.
func (c *CPU) UpdateTimers() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.drewThisFrame = false
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
	c.keyWait = newKeyWait()
	c.drewThisFrame = false

	c.Memory = [len(c.Memory)]byte{}
	copy(c.Memory[:], fontSet[:])
	c.V = [len(c.V)]byte{}
	c.Stack = [len(c.Stack)]uint16{}
	c.Display = [len(c.Display)]byte{}
	c.Key = [len(c.Key)]bool{}
}

// State returns a copy of the CPU state for safe access.
func (c *CPU) State() State {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var state State
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
	state.keyWait = c.keyWait
	state.drewThisFrame = c.drewThisFrame

	return state
}

// SetState sets the CPU state from a snapshot.
func (c *CPU) SetState(state State) {
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
	c.keyWait = state.keyWait
	c.drewThisFrame = state.drewThisFrame
}

// updatePC increments the program counter to the next instruction and optionally skips the following instruction.
func (c *CPU) updatePC(skipInstruction bool) {
	if skipInstruction {
		c.PC += 4
	} else {
		c.PC += 2
	}
}

func newKeyWait() keyWait {
	return keyWait{key: -1}
}
