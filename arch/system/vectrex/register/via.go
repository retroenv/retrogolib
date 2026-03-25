// Package register contains constants for Vectrex hardware register addresses.
package register

// VIA (MC6522 Versatile Interface Adapter) registers ($D000-$D00F).
// The VIA handles all I/O for the Vectrex including:
// - DAC output for vector display (X/Y position, beam intensity)
// - Sound chip (AY-3-8912) interface via port A
// - Controller/button input via port B
// - Timer for display refresh and game timing
const (
	VIAORB  = 0xD000 // Output Register B (controller buttons, mux select)
	VIAORA  = 0xD001 // Output Register A (DAC data, sound chip data)
	VIADDRB = 0xD002 // Data Direction Register B
	VIADDRA = 0xD003 // Data Direction Register A
	VIAT1CL = 0xD004 // Timer 1 Counter Low
	VIAT1CH = 0xD005 // Timer 1 Counter High
	VIAT1LL = 0xD006 // Timer 1 Latch Low
	VIAT1LH = 0xD007 // Timer 1 Latch High
	VIAT2CL = 0xD008 // Timer 2 Counter Low
	VIAT2CH = 0xD009 // Timer 2 Counter High
	VIASR   = 0xD00A // Shift Register
	VIAACR  = 0xD00B // Auxiliary Control Register
	VIAPCR  = 0xD00C // Peripheral Control Register
	VIAIFR  = 0xD00D // Interrupt Flag Register
	VIAIER  = 0xD00E // Interrupt Enable Register
	VIAORAF = 0xD00F // Output Register A (no handshake)
)

// VIARegisterCount is the number of VIA registers.
const VIARegisterCount = 16

// VIANames maps VIA register addresses to their names.
var VIANames = map[uint16]string{
	VIAORB:  "VIAORB",
	VIAORA:  "VIAORA",
	VIADDRB: "VIADDRB",
	VIADDRA: "VIADDRA",
	VIAT1CL: "VIAT1CL",
	VIAT1CH: "VIAT1CH",
	VIAT1LL: "VIAT1LL",
	VIAT1LH: "VIAT1LH",
	VIAT2CL: "VIAT2CL",
	VIAT2CH: "VIAT2CH",
	VIASR:   "VIASR",
	VIAACR:  "VIAACR",
	VIAPCR:  "VIAPCR",
	VIAIFR:  "VIAIFR",
	VIAIER:  "VIAIER",
	VIAORAF: "VIAORAF",
}

// Port B button bits (active low).
const (
	ButtonRight = 0x01 // Joystick button 1 (right)
	ButtonLeft  = 0x02 // Joystick button 2 (left)
	ButtonDown  = 0x04 // Joystick button 3 (down)
	ButtonUp    = 0x08 // Joystick button 4 (up)
)

// VIA Interrupt Flag/Enable Register bits.
const (
	IRQTimer1 = 0x40 // Timer 1 interrupt
	IRQTimer2 = 0x20 // Timer 2 interrupt
	IRQCB1    = 0x10 // CB1 interrupt
	IRQCB2    = 0x08 // CB2 interrupt
	IRQShift  = 0x04 // Shift register interrupt
	IRQCA1    = 0x02 // CA1 interrupt
	IRQCA2    = 0x01 // CA2 interrupt
	IRQAny    = 0x80 // Any interrupt flag (IFR bit 7) / master enable (IER bit 7)
)
