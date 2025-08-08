package z80

// Memory represents the Z80 CPU memory with support for banking and memory mapping.
type Memory struct {
	data [0x10000]uint8 // 64KB address space

	// Memory banking support (for Game Boy)
	romBank uint8
	ramBank uint8

	// Memory control registers
	mbc1Mode bool // MBC1 mode flag for banking
}

// NewMemory creates a new Z80 memory instance.
func NewMemory() *Memory {
	return &Memory{
		romBank: 1, // ROM bank 1 is default for Game Boy
		ramBank: 0,
	}
}

// Read reads a byte from memory at the given address.
func (m *Memory) Read(address uint16) uint8 {
	return m.data[address]
}

// Write writes a byte to memory at the given address.
func (m *Memory) Write(address uint16, value uint8) {
	// For basic Z80 emulation, allow writes to all memory areas
	// In a full Game Boy emulator, banking logic would be more complex
	m.data[address] = value

	// Optional: Update banking state for Game Boy compatibility
	// This allows the banking registers to be tracked without preventing writes
	switch {
	case address >= 0x2000 && address < 0x4000:
		// ROM bank number (0x2000-0x3FFF)
		bank := value & 0x1F
		if bank == 0 {
			bank = 1 // Bank 0 is not directly accessible (hardware constraint)
		}
		m.romBank = bank

	case address >= 0x4000 && address < 0x6000:
		// RAM bank number or upper ROM bank bits (0x4000-0x5FFF)
		if m.mbc1Mode {
			m.ramBank = value & 0x03
		} else {
			// Upper 2 bits of ROM bank for larger ROMs
			m.romBank = (m.romBank & 0x1F) | ((value & 0x03) << 5)
		}

	case address >= 0x6000 && address < 0x8000:
		// Banking mode select (0x6000-0x7FFF)
		m.mbc1Mode = (value & 0x01) != 0
	}
}

// ReadWord reads a 16-bit word from memory at the given address (little-endian).
func (m *Memory) ReadWord(address uint16) uint16 {
	low := uint16(m.Read(address))
	high := uint16(m.Read(address + 1))
	return high<<8 | low
}

// WriteWord writes a 16-bit word to memory at the given address (little-endian).
func (m *Memory) WriteWord(address uint16, value uint16) {
	m.Write(address, uint8(value))
	m.Write(address+1, uint8(value>>8))
}

// LoadROM loads ROM data into memory starting at address 0.
func (m *Memory) LoadROM(data []byte) {
	if data == nil {
		return
	}

	n := min(len(data), len(m.data))
	if n > 0 {
		copy(m.data[:n], data[:n])
	}
}

// GetROMBank returns the current ROM bank number.
func (m *Memory) GetROMBank() uint8 {
	return m.romBank
}

// GetRAMBank returns the current RAM bank number.
func (m *Memory) GetRAMBank() uint8 {
	return m.ramBank
}

// SetBankingMode sets the MBC1 banking mode.
func (m *Memory) SetBankingMode(mode bool) {
	m.mbc1Mode = mode
}
