package chip8

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

	Display      [64 * 32]bool // Monochrome display (64x32)
	RedrawScreen bool          // Indicates if the screen needs to be redrawn
}

const (
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

func (c *CPU) fetchOpcode() uint16 {
	return uint16(c.Memory[c.PC])<<8 | uint16(c.Memory[c.PC+1])
}

// updatePC increments the program counter to the next instruction and optionally skips the following instruction.
func (c *CPU) updatePC(skipInstruction bool) {
	if skipInstruction {
		c.PC += 4
	} else {
		c.PC += 2
	}
}
