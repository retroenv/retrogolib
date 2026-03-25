// Package register contains constants for TRS-80 CoCo hardware register addresses.
package register

// PIA 0 registers ($FF00-$FF03).
// PIA 0 handles keyboard scanning, joystick input, and cassette I/O.
const (
	PIA0DRA = 0xFF00 // Data Register A (keyboard rows / joystick)
	PIA0CRA = 0xFF01 // Control Register A
	PIA0DRB = 0xFF02 // Data Register B (keyboard columns)
	PIA0CRB = 0xFF03 // Control Register B
)

// PIA 1 registers ($FF20-$FF23).
// PIA 1 handles VDG display control, serial port, and sound output.
const (
	PIA1DRA = 0xFF20 // Data Register A (DAC, cassette, serial)
	PIA1CRA = 0xFF21 // Control Register A
	PIA1DRB = 0xFF22 // Data Register B (VDG mode, sound)
	PIA1CRB = 0xFF23 // Control Register B
)

// PIA0Names maps PIA 0 register addresses to their names.
var PIA0Names = map[uint16]string{
	PIA0DRA: "PIA0DRA",
	PIA0CRA: "PIA0CRA",
	PIA0DRB: "PIA0DRB",
	PIA0CRB: "PIA0CRB",
}

// PIA1Names maps PIA 1 register addresses to their names.
var PIA1Names = map[uint16]string{
	PIA1DRA: "PIA1DRA",
	PIA1CRA: "PIA1CRA",
	PIA1DRB: "PIA1DRB",
	PIA1CRB: "PIA1CRB",
}
