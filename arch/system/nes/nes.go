// Package nes provides support for the Nintendo Entertainment System (NES).
package nes

const (
	// CodeBaseAddress defines the start address of code for the NES.
	CodeBaseAddress = 0x8000

	// IORegisterStartAddress defines the start address of the I/O registers in the NES.
	IORegisterStartAddress = 0x4000

	// IORegisterEndAddress defines the address of the last I/O registers in the NES.
	IORegisterEndAddress = 0x401F

	// NameTableCount defines the number of name tables in the NES.
	NameTableCount = 4

	// NameTableSize defines the size of a name table in bytes.
	NameTableSize = 0x400

	// PaletteSize defines the size of the NES palette in bytes.
	PaletteSize = 32

	// RAMEndAddress defines the end address of RAM in the NES.
	RAMEndAddress = 0x0FFF
)
