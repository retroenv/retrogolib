package x86

// ModRM represents the ModR/M byte used in x86 instruction encoding.
// Bit pattern: [Mod:2][Reg:3][R/M:3]
type ModRM struct {
	Mod uint8 // Mode field (bits 7-6)
	Reg uint8 // Register field (bits 5-3)
	RM  uint8 // R/M field (bits 2-0)
}

// NewModRM creates a ModR/M byte from its components.
func NewModRM(mod, reg, rm uint8) ModRM {
	return ModRM{
		Mod: mod & 0x03,
		Reg: reg & 0x07,
		RM:  rm & 0x07,
	}
}

// FromByte creates a ModR/M from a raw byte value.
func (m *ModRM) FromByte(value uint8) {
	m.Mod = (value >> 6) & 0x03
	m.Reg = (value >> 3) & 0x07
	m.RM = value & 0x07
}

// ToByte converts the ModR/M to a raw byte value.
func (m ModRM) ToByte() uint8 {
	return (m.Mod << 6) | (m.Reg << 3) | m.RM
}
