// Package chip8 provides support for the virtual Chip-8 CPU.
package chip8

// Instruction contains information about a CPU instruction.
type Instruction struct {
	Name string // lowercased instruction name

	Addressing map[Mode]OpcodeInfo

	Emulation func(c *CPU, param uint16) error // emulation function to execute
}

// Standard Chip-8 Instructions

// Add - adds a value/register to a register.
var Add = &Instruction{
	Name:      "add",
	Emulation: add,
	Addressing: map[Mode]OpcodeInfo{
		RegisterValueAddressing:    Opcode7000,
		RegisterRegisterAddressing: Opcode8004,
		IRegisterAddressing:        OpcodeF01E,
	},
}

// And - performs a bitwise AND operation on two registers.
var And = &Instruction{
	Name:      "and",
	Emulation: and,
	Addressing: map[Mode]OpcodeInfo{
		RegisterRegisterAddressing: Opcode8002,
	},
}

// Call - Call subroutine.
var Call = &Instruction{
	Name:      "call",
	Emulation: call,
	Addressing: map[Mode]OpcodeInfo{
		AbsoluteAddressing: Opcode2000,
	},
}

// Cls - Clear screen.
var Cls = &Instruction{
	Name:      "cls",
	Emulation: cls,
	Addressing: map[Mode]OpcodeInfo{
		ImpliedAddressing: Opcode00E0,
	},
}

// Drw - Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.
var Drw = &Instruction{
	Name:      "drw",
	Emulation: drw,
	Addressing: map[Mode]OpcodeInfo{
		RegisterRegisterNibbleAddressing: OpcodeD000,
	},
}

// Jp - jumps to an address and optionally adds V0 to the address.
var Jp = &Instruction{
	Name:      "jp",
	Emulation: jp,
	Addressing: map[Mode]OpcodeInfo{
		AbsoluteAddressing:   Opcode1000,
		V0AbsoluteAddressing: OpcodeB000,
	},
}

// Ld - Set Vx = kk.
var Ld = &Instruction{
	Name:      "ld",
	Emulation: ld,
	Addressing: map[Mode]OpcodeInfo{
		RegisterValueAddressing:     Opcode6000,
		RegisterRegisterAddressing:  Opcode8000,
		IAbsoluteAddressing:         OpcodeA000,
		RegisterDTAddressing:        OpcodeF007,
		RegisterKAddressing:         OpcodeF00A,
		DTRegisterAddressing:        OpcodeF015,
		STRegisterAddressing:        OpcodeF018,
		FRegisterAddressing:         OpcodeF029,
		BRegisterAddressing:         OpcodeF033,
		IIndirectRegisterAddressing: OpcodeF055,
		RegisterIndirectIAddressing: OpcodeF065,
	},
}

// Or - performs a bitwise OR operation on two registers.
var Or = &Instruction{
	Name:      "or",
	Emulation: or,
	Addressing: map[Mode]OpcodeInfo{
		RegisterRegisterAddressing: Opcode8001,
	},
}

// Ret - Return from a subroutine.
var Ret = &Instruction{
	Name:      "ret",
	Emulation: ret,
	Addressing: map[Mode]OpcodeInfo{
		ImpliedAddressing: Opcode00EE,
	},
}

// Rnd - Set Vx = random byte AND kk.
var Rnd = &Instruction{
	Name:      "rnd",
	Emulation: rnd,
	Addressing: map[Mode]OpcodeInfo{
		RegisterValueAddressing: OpcodeC000,
	},
}

// Se - Skip next instruction if the register equals a value/register.
var Se = &Instruction{
	Name:      "se",
	Emulation: se,
	Addressing: map[Mode]OpcodeInfo{
		RegisterValueAddressing:    Opcode3000,
		RegisterRegisterAddressing: Opcode5000,
	},
}

// Shl - Set Vx = Vx SHL 1.
var Shl = &Instruction{
	Name:      "shl",
	Emulation: shl,
	Addressing: map[Mode]OpcodeInfo{
		RegisterRegisterAddressing: Opcode800E,
	},
}

// Shr - Set Vx = Vx SHR 1.
var Shr = &Instruction{
	Name:      "shr",
	Emulation: shr,
	Addressing: map[Mode]OpcodeInfo{
		RegisterRegisterAddressing: Opcode8006,
	},
}

// Skp - Skip next instruction if key with the value of Vx is pressed.
var Skp = &Instruction{
	Name:      "skp",
	Emulation: skp,
	Addressing: map[Mode]OpcodeInfo{
		RegisterValueAddressing: OpcodeE09E,
	},
}

// Sknp - Skip next instruction if key with the value of Vx is not pressed.
var Sknp = &Instruction{
	Name:      "sknp",
	Emulation: sknp,
	Addressing: map[Mode]OpcodeInfo{
		RegisterValueAddressing: OpcodeE0A1,
	},
}

// Sne - Skip next instruction if the register does not equal a value/register.
var Sne = &Instruction{
	Name:      "sne",
	Emulation: sne,
	Addressing: map[Mode]OpcodeInfo{
		RegisterValueAddressing:    Opcode4000,
		RegisterRegisterAddressing: Opcode9000,
	},
}

// Sub - Set Vx = Vx - Vy, set VF = NOT borrow.
var Sub = &Instruction{
	Name:      "sub",
	Emulation: sub,
	Addressing: map[Mode]OpcodeInfo{
		RegisterRegisterAddressing: Opcode8005,
	},
}

// Subn - Set Vx = Vy - Vx, set VF = NOT borrow.
var Subn = &Instruction{
	Name:      "subn",
	Emulation: subn,
	Addressing: map[Mode]OpcodeInfo{
		RegisterRegisterAddressing: Opcode8007,
	},
}

// Xor - performs a bitwise XOR operation on two registers.
var Xor = &Instruction{
	Name:      "xor",
	Emulation: xor,
	Addressing: map[Mode]OpcodeInfo{
		RegisterRegisterAddressing: Opcode8003,
	},
}
