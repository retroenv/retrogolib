// Package chip8 provides support for the virtual Chip-8 CPU.
package chip8

// Instruction defines a Chip-8 instruction with its opcodes and execution logic.
// Instructions support multiple addressing modes through the Addressing map that
// enables opcode lookup for disassembly and code generation.
type Instruction struct {
	Name string // Instruction mnemonic (lowercase)

	// Opcode lookup map for addressing mode to opcode mapping
	Addressing map[Mode]OpcodeInfo // Maps addressing mode to opcode info

	// Execution handler - receives CPU state and 16-bit opcode parameter
	Emulation func(c *CPU, param uint16) error // Handler for instruction execution
}

// Standard Chip-8 Instructions

// Add adds a value or register to a register (ADD Vx, byte / ADD Vx, Vy / ADD I, Vx).
var Add = &Instruction{
	Name:      "add",
	Emulation: add,
	Addressing: map[Mode]OpcodeInfo{
		RegisterValueAddressing:    Opcode7000,
		RegisterRegisterAddressing: Opcode8004,
		IRegisterAddressing:        OpcodeF01E,
	},
}

// And performs bitwise AND on two registers (AND Vx, Vy).
var And = &Instruction{
	Name:      "and",
	Emulation: and,
	Addressing: map[Mode]OpcodeInfo{
		RegisterRegisterAddressing: Opcode8002,
	},
}

// Call calls a subroutine at address (CALL addr).
var Call = &Instruction{
	Name:      "call",
	Emulation: call,
	Addressing: map[Mode]OpcodeInfo{
		AbsoluteAddressing: Opcode2000,
	},
}

// Cls clears the display screen (CLS).
var Cls = &Instruction{
	Name:      "cls",
	Emulation: cls,
	Addressing: map[Mode]OpcodeInfo{
		ImpliedAddressing: Opcode00E0,
	},
}

// Drw draws n-byte sprite from memory location I at (Vx, Vy), sets VF on collision (DRW Vx, Vy, nibble).
var Drw = &Instruction{
	Name:      "drw",
	Emulation: drw,
	Addressing: map[Mode]OpcodeInfo{
		RegisterRegisterNibbleAddressing: OpcodeD000,
	},
}

// Jp jumps to address, optionally adding V0 to the address (JP addr / JP V0, addr).
var Jp = &Instruction{
	Name:      "jp",
	Emulation: jp,
	Addressing: map[Mode]OpcodeInfo{
		AbsoluteAddressing:   Opcode1000,
		V0AbsoluteAddressing: OpcodeB000,
	},
}

// Ld loads values into registers, timers, or memory (LD Vx, byte / LD I, addr / LD DT, Vx).
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

// Or performs bitwise OR on two registers (OR Vx, Vy).
var Or = &Instruction{
	Name:      "or",
	Emulation: or,
	Addressing: map[Mode]OpcodeInfo{
		RegisterRegisterAddressing: Opcode8001,
	},
}

// Ret returns from a subroutine (RET).
var Ret = &Instruction{
	Name:      "ret",
	Emulation: ret,
	Addressing: map[Mode]OpcodeInfo{
		ImpliedAddressing: Opcode00EE,
	},
}

// Rnd sets Vx to random byte AND immediate value (RND Vx, byte).
var Rnd = &Instruction{
	Name:      "rnd",
	Emulation: rnd,
	Addressing: map[Mode]OpcodeInfo{
		RegisterValueAddressing: OpcodeC000,
	},
}

// Se skips next instruction if register equals value or register (SE Vx, byte / SE Vx, Vy).
var Se = &Instruction{
	Name:      "se",
	Emulation: se,
	Addressing: map[Mode]OpcodeInfo{
		RegisterValueAddressing:    Opcode3000,
		RegisterRegisterAddressing: Opcode5000,
	},
}

// Shl shifts Vx left by 1, stores MSB in VF (SHL Vx).
var Shl = &Instruction{
	Name:      "shl",
	Emulation: shl,
	Addressing: map[Mode]OpcodeInfo{
		RegisterRegisterAddressing: Opcode800E,
	},
}

// Shr shifts Vx right by 1, stores LSB in VF (SHR Vx).
var Shr = &Instruction{
	Name:      "shr",
	Emulation: shr,
	Addressing: map[Mode]OpcodeInfo{
		RegisterRegisterAddressing: Opcode8006,
	},
}

// Skp skips next instruction if key with value of Vx is pressed (SKP Vx).
var Skp = &Instruction{
	Name:      "skp",
	Emulation: skp,
	Addressing: map[Mode]OpcodeInfo{
		RegisterValueAddressing: OpcodeE09E,
	},
}

// Sknp skips next instruction if key with value of Vx is not pressed (SKNP Vx).
var Sknp = &Instruction{
	Name:      "sknp",
	Emulation: sknp,
	Addressing: map[Mode]OpcodeInfo{
		RegisterValueAddressing: OpcodeE0A1,
	},
}

// Sne skips next instruction if register does not equal value or register (SNE Vx, byte / SNE Vx, Vy).
var Sne = &Instruction{
	Name:      "sne",
	Emulation: sne,
	Addressing: map[Mode]OpcodeInfo{
		RegisterValueAddressing:    Opcode4000,
		RegisterRegisterAddressing: Opcode9000,
	},
}

// Sub subtracts Vy from Vx, sets VF = NOT borrow (SUB Vx, Vy).
var Sub = &Instruction{
	Name:      "sub",
	Emulation: sub,
	Addressing: map[Mode]OpcodeInfo{
		RegisterRegisterAddressing: Opcode8005,
	},
}

// Subn subtracts Vx from Vy, stores result in Vx, sets VF = NOT borrow (SUBN Vx, Vy).
var Subn = &Instruction{
	Name:      "subn",
	Emulation: subn,
	Addressing: map[Mode]OpcodeInfo{
		RegisterRegisterAddressing: Opcode8007,
	},
}

// Xor performs bitwise XOR on two registers (XOR Vx, Vy).
var Xor = &Instruction{
	Name:      "xor",
	Emulation: xor,
	Addressing: map[Mode]OpcodeInfo{
		RegisterRegisterAddressing: Opcode8003,
	},
}
