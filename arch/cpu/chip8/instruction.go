// Package chip8 provides support for the Chip-8 and Chip-48 CPU.
package chip8

// Instruction contains information about a CPU instruction.
type Instruction struct {
	Name string // lowercased instruction name

	Emulation func(c *CPU, param uint16) error // emulation function to execute
}

// Standard Chip-8 Instructions

// Add - adds a value/register to a register.
var Add = &Instruction{
	Name:      "add",
	Emulation: add,
}

// And - performs a bitwise AND operation on two registers.
var And = &Instruction{
	Name:      "and",
	Emulation: and,
}

// Call - Call subroutine.
var Call = &Instruction{
	Name:      "call",
	Emulation: call,
}

// Cls - Clear screen.
var Cls = &Instruction{
	Name:      "cls",
	Emulation: cls,
}

// Drw - Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.
var Drw = &Instruction{
	Name:      "drw",
	Emulation: drw,
}

// Jp - jumps to an address and optionally adds V0 to the address.
var Jp = &Instruction{
	Name:      "jp",
	Emulation: jp,
}

// Ld - Set Vx = kk.
var Ld = &Instruction{
	Name:      "ld",
	Emulation: ld,
}

// Or - performs a bitwise OR operation on two registers.
var Or = &Instruction{
	Name:      "or",
	Emulation: or,
}

// Ret - Return from a subroutine.
var Ret = &Instruction{
	Name:      "ret",
	Emulation: ret,
}

// Rnd - Set Vx = random byte AND kk.
var Rnd = &Instruction{
	Name:      "rnd",
	Emulation: rnd,
}

// Se - Skip next instruction if the register equals a value/register.
var Se = &Instruction{
	Name:      "se",
	Emulation: se,
}

// Shl - Set Vx = Vx SHL 1.
var Shl = &Instruction{
	Name:      "shl",
	Emulation: shl,
}

// Shr - Set Vx = Vx SHR 1.
var Shr = &Instruction{
	Name:      "shr",
	Emulation: shr,
}

// Skp - Skip next instruction if key with the value of Vx is pressed.
var Skp = &Instruction{
	Name:      "skp",
	Emulation: skp,
}

// Sknp - Skip next instruction if key with the value of Vx is not pressed.
var Sknp = &Instruction{
	Name:      "sknp",
	Emulation: sknp,
}

// Sne - Skip next instruction if the register does not equal a value/register.
var Sne = &Instruction{
	Name:      "sne",
	Emulation: sne,
}

// Sub - Set Vx = Vx - Vy, set VF = NOT borrow.
var Sub = &Instruction{
	Name:      "sub",
	Emulation: sub,
}

// Subn - Set Vx = Vy - Vx, set VF = NOT borrow.
var Subn = &Instruction{
	Name:      "subn",
	Emulation: subn,
}

// Xor - performs a bitwise XOR operation on two registers.
var Xor = &Instruction{
	Name:      "xor",
	Emulation: xor,
}
