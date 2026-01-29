package z80

// GameBoyMemory implements Game Boy specific memory banking (MBC1).
// This is an example of hardware-specific memory controller.
type GameBoyMemory struct {
	data [0x10000]uint8

	// Memory banking support (for Game Boy MBC1)
	romBank  uint8
	ramBank  uint8
	mbc1Mode bool // MBC1 mode flag for banking
}

// NewGameBoyMemory creates a new Game Boy memory controller with MBC1 support.
func NewGameBoyMemory() *GameBoyMemory {
	return &GameBoyMemory{
		romBank: 1, // ROM bank 1 is default for Game Boy
		ramBank: 0,
	}
}

// Read reads a byte from memory at the given address.
func (mem *GameBoyMemory) Read(address uint16) uint8 {
	return mem.data[address]
}

// Write writes a byte to memory at the given address.
// Handles Game Boy MBC1 banking register writes.
func (mem *GameBoyMemory) Write(address uint16, value uint8) {
	// Allow writes to all memory areas
	mem.data[address] = value

	// Update banking state for Game Boy MBC1 compatibility
	switch {
	case address >= 0x2000 && address < 0x4000:
		// ROM bank number (0x2000-0x3FFF)
		bank := value & 0x1F
		if bank == 0 {
			bank = 1 // Bank 0 is not directly accessible (hardware constraint)
		}
		mem.romBank = bank

	case address >= 0x4000 && address < 0x6000:
		// RAM bank number or upper ROM bank bits (0x4000-0x5FFF)
		if mem.mbc1Mode {
			mem.ramBank = value & 0x03
		} else {
			// Upper 2 bits of ROM bank for larger ROMs
			mem.romBank = (mem.romBank & 0x1F) | ((value & 0x03) << 5)
		}

	case address >= 0x6000 && address < 0x8000:
		// Banking mode select (0x6000-0x7FFF)
		mem.mbc1Mode = (value & 0x01) != 0
	}
}

// ReadWord reads a 16-bit word from memory at the given address (little-endian).
func (mem *GameBoyMemory) ReadWord(address uint16) uint16 {
	low := uint16(mem.Read(address))
	high := uint16(mem.Read(address + 1))
	return high<<8 | low
}

// WriteWord writes a 16-bit word to memory at the given address (little-endian).
func (mem *GameBoyMemory) WriteWord(address uint16, value uint16) {
	mem.Write(address, uint8(value))
	mem.Write(address+1, uint8(value>>8))
}

// LoadROM loads ROM data into memory starting at address 0.
func (mem *GameBoyMemory) LoadROM(data []byte) {
	if data == nil {
		return
	}

	n := min(len(data), len(mem.data))
	if n > 0 {
		copy(mem.data[:n], data[:n])
	}
}

// LoadProgram loads program data into memory starting at address 0.
// This is an alias for LoadROM for backward compatibility.
func (mem *GameBoyMemory) LoadProgram(data []byte) {
	mem.LoadROM(data)
}

// GetROMBank returns the current ROM bank number.
func (mem *GameBoyMemory) GetROMBank() uint8 {
	return mem.romBank
}

// GetRAMBank returns the current RAM bank number.
func (mem *GameBoyMemory) GetRAMBank() uint8 {
	return mem.ramBank
}

// SetBankingMode sets the MBC1 banking mode.
func (mem *GameBoyMemory) SetBankingMode(mode bool) {
	mem.mbc1Mode = mode
}
