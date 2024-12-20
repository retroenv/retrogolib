package chip8

import (
	"fmt"
	"math/rand"
	"time"
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

	rnd rand.Source // Random number generator
}

const (
	displayHeight         = 32
	displayWidth          = 64
	initialProgramCounter = 0x200
)

// New creates a new CPU.
func New() *CPU {
	c := &CPU{
		PC:  initialProgramCounter,
		rnd: rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	// Load fontset into memory
	copy(c.Memory[:], fontSet[:])

	return c
}

// Step executes the next instruction in the CPU.
func (c *CPU) Step() error {
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
