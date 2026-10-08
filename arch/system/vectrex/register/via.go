// Package register contains constants for Vectrex hardware register addresses.
package register

// VIA (6522 Versatile Interface Adapter) registers ($D000-$D00F).
// The VIA handles all I/O for the Vectrex including:
// - DAC output for vector display (X/Y position, beam intensity) through port A
// - Sound chip (AY-3-8912) data through port A and control lines on port B
// - Analog joystick comparison and multiplexer selection through port B
// - Timer for display refresh and game timing
//
// The joystick buttons are not VIA inputs. The PSG I/O port supplies them;
// see PSGIOPortA.
const (
	VIAORB  = 0xD000 // Output Register B (multiplexer select, PSG control, ramp)
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

// VIA port B signal bits (VIAORB and VIADDRB).
const (
	PortBSwitch  = 0x01 // PB0: analog multiplexer sample/hold switch (output, active low)
	PortBSel0    = 0x02 // PB1: multiplexer select bit 0 (output)
	PortBSel1    = 0x04 // PB2: multiplexer select bit 1 (output)
	PortBBC1     = 0x08 // PB3: PSG BC1 control (output)
	PortBBDIR    = 0x10 // PB4: PSG BDIR control (output)
	PortBCompare = 0x20 // PB5: comparator result for joystick position (input)
	PortBCart    = 0x40 // PB6: cartridge detect line (input)
	PortBRamp    = 0x80 // PB7: integrator ramp control (output, active low)
)

// PSGIOPortA is the AY-3-8912 register that reads the joystick buttons.
// The BIOS reads it through the VIA port A with BDIR/BC1 on port B.
const PSGIOPortA = 0x0E

// Joystick button bits in PSGIOPortA (active low).
const (
	Player1Button1 = 0x01
	Player1Button2 = 0x02
	Player1Button3 = 0x04
	Player1Button4 = 0x08
	Player2Button1 = 0x10
	Player2Button2 = 0x20
	Player2Button3 = 0x40
	Player2Button4 = 0x80
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
