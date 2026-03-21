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

// Instruction name constants for easy access by external packages.
const (
	AddName  = "add"
	AndName  = "and"
	CallName = "call"
	ClsName  = "cls"
	DrwName  = "drw"
	JpName   = "jp"
	LdName   = "ld"
	OrName   = "or"
	RetName  = "ret"
	RndName  = "rnd"
	SeName   = "se"
	ShlName  = "shl"
	ShrName  = "shr"
	SkpName  = "skp"
	SknpName = "sknp"
	SneName  = "sne"
	SubName  = "sub"
	SubnName = "subn"
	XorName  = "xor"
)

// Standard Chip-8 Instructions

// AddInst adds a value or register to a register (ADD Vx, byte / ADD Vx, Vy / ADD I, Vx).
var AddInst = &Instruction{
	Name:      AddName,
	Emulation: add,
	Addressing: map[Mode]OpcodeInfo{
		RegisterValueAddressing:    Opcode7000,
		RegisterRegisterAddressing: Opcode8004,
		IRegisterAddressing:        OpcodeF01E,
	},
}

// AndInst performs bitwise AND on two registers (AND Vx, Vy).
var AndInst = &Instruction{
	Name:      AndName,
	Emulation: and,
	Addressing: map[Mode]OpcodeInfo{
		RegisterRegisterAddressing: Opcode8002,
	},
}

// CallInst calls a subroutine at address (CALL addr).
var CallInst = &Instruction{
	Name:      CallName,
	Emulation: call,
	Addressing: map[Mode]OpcodeInfo{
		AbsoluteAddressing: Opcode2000,
	},
}

// ClsInst clears the display screen (CLS).
var ClsInst = &Instruction{
	Name:      ClsName,
	Emulation: cls,
	Addressing: map[Mode]OpcodeInfo{
		ImpliedAddressing: Opcode00E0,
	},
}

// DrwInst draws n-byte sprite from memory location I at (Vx, Vy), sets VF on collision (DRW Vx, Vy, nibble).
var DrwInst = &Instruction{
	Name:      DrwName,
	Emulation: drw,
	Addressing: map[Mode]OpcodeInfo{
		RegisterRegisterNibbleAddressing: OpcodeD000,
	},
}

// JpInst jumps to address, optionally adding V0 to the address (JP addr / JP V0, addr).
var JpInst = &Instruction{
	Name:      JpName,
	Emulation: jp,
	Addressing: map[Mode]OpcodeInfo{
		AbsoluteAddressing:   Opcode1000,
		V0AbsoluteAddressing: OpcodeB000,
	},
}

// LdInst loads values into registers, timers, or memory (LD Vx, byte / LD I, addr / LD DT, Vx).
var LdInst = &Instruction{
	Name:      LdName,
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

// OrInst performs bitwise OR on two registers (OR Vx, Vy).
var OrInst = &Instruction{
	Name:      OrName,
	Emulation: or,
	Addressing: map[Mode]OpcodeInfo{
		RegisterRegisterAddressing: Opcode8001,
	},
}

// RetInst returns from a subroutine (RET).
var RetInst = &Instruction{
	Name:      RetName,
	Emulation: ret,
	Addressing: map[Mode]OpcodeInfo{
		ImpliedAddressing: Opcode00EE,
	},
}

// RndInst sets Vx to random byte AND immediate value (RND Vx, byte).
var RndInst = &Instruction{
	Name:      RndName,
	Emulation: rnd,
	Addressing: map[Mode]OpcodeInfo{
		RegisterValueAddressing: OpcodeC000,
	},
}

// SeInst skips next instruction if register equals value or register (SE Vx, byte / SE Vx, Vy).
var SeInst = &Instruction{
	Name:      SeName,
	Emulation: se,
	Addressing: map[Mode]OpcodeInfo{
		RegisterValueAddressing:    Opcode3000,
		RegisterRegisterAddressing: Opcode5000,
	},
}

// ShlInst shifts Vx left by 1, stores MSB in VF (SHL Vx).
var ShlInst = &Instruction{
	Name:      ShlName,
	Emulation: shl,
	Addressing: map[Mode]OpcodeInfo{
		RegisterRegisterAddressing: Opcode800E,
	},
}

// ShrInst shifts Vx right by 1, stores LSB in VF (SHR Vx).
var ShrInst = &Instruction{
	Name:      ShrName,
	Emulation: shr,
	Addressing: map[Mode]OpcodeInfo{
		RegisterRegisterAddressing: Opcode8006,
	},
}

// SkpInst skips next instruction if key with value of Vx is pressed (SKP Vx).
var SkpInst = &Instruction{
	Name:      SkpName,
	Emulation: skp,
	Addressing: map[Mode]OpcodeInfo{
		RegisterValueAddressing: OpcodeE09E,
	},
}

// SknpInst skips next instruction if key with value of Vx is not pressed (SKNP Vx).
var SknpInst = &Instruction{
	Name:      SknpName,
	Emulation: sknp,
	Addressing: map[Mode]OpcodeInfo{
		RegisterValueAddressing: OpcodeE0A1,
	},
}

// SneInst skips next instruction if register does not equal value or register (SNE Vx, byte / SNE Vx, Vy).
var SneInst = &Instruction{
	Name:      SneName,
	Emulation: sne,
	Addressing: map[Mode]OpcodeInfo{
		RegisterValueAddressing:    Opcode4000,
		RegisterRegisterAddressing: Opcode9000,
	},
}

// SubInst subtracts Vy from Vx, sets VF = NOT borrow (SUB Vx, Vy).
var SubInst = &Instruction{
	Name:      SubName,
	Emulation: sub,
	Addressing: map[Mode]OpcodeInfo{
		RegisterRegisterAddressing: Opcode8005,
	},
}

// SubnInst subtracts Vx from Vy, stores result in Vx, sets VF = NOT borrow (SUBN Vx, Vy).
var SubnInst = &Instruction{
	Name:      SubnName,
	Emulation: subn,
	Addressing: map[Mode]OpcodeInfo{
		RegisterRegisterAddressing: Opcode8007,
	},
}

// XorInst performs bitwise XOR on two registers (XOR Vx, Vy).
var XorInst = &Instruction{
	Name:      XorName,
	Emulation: xor,
	Addressing: map[Mode]OpcodeInfo{
		RegisterRegisterAddressing: Opcode8003,
	},
}

// Instructions maps instruction names to their information struct.
var Instructions = map[string]*Instruction{
	AddName:  AddInst,
	AndName:  AndInst,
	CallName: CallInst,
	ClsName:  ClsInst,
	DrwName:  DrwInst,
	JpName:   JpInst,
	LdName:   LdInst,
	OrName:   OrInst,
	RetName:  RetInst,
	RndName:  RndInst,
	SeName:   SeInst,
	ShlName:  ShlInst,
	ShrName:  ShrInst,
	SkpName:  SkpInst,
	SknpName: SknpInst,
	SneName:  SneInst,
	SubName:  SubInst,
	SubnName: SubnInst,
	XorName:  XorInst,
}
