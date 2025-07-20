package gameboy

import (
	"errors"
	"fmt"
)

// CartridgeHeader represents the cartridge header structure.
type CartridgeHeader struct {
	EntryPoint       [4]byte  // 0x0100-0x0103: Entry point
	NintendoLogo     [48]byte // 0x0104-0x0133: Nintendo logo
	Title            [15]byte // 0x0134-0x0142: Title
	ManufacturerCode [4]byte  // 0x013F-0x0142: Manufacturer code
	CGBFlag          uint8    // 0x0143: CGB flag
	NewLicenseeCode  [2]byte  // 0x0144-0x0145: New licensee code
	SGBFlag          uint8    // 0x0146: SGB flag
	CartridgeType    uint8    // 0x0147: Cartridge type
	ROMSize          uint8    // 0x0148: ROM size
	RAMSize          uint8    // 0x0149: RAM size
	DestinationCode  uint8    // 0x014A: Destination code
	OldLicenseeCode  uint8    // 0x014B: Old licensee code
	MaskROMVersion   uint8    // 0x014C: Mask ROM version
	HeaderChecksum   uint8    // 0x014D: Header checksum
	GlobalChecksum   [2]byte  // 0x014E-0x014F: Global checksum
}

// ParseCartridgeHeader parses the cartridge header from ROM data.
func ParseCartridgeHeader(rom []byte) (*CartridgeHeader, error) {
	if len(rom) < 0x150 {
		return nil, errors.New("ROM too small to contain header")
	}

	header := &CartridgeHeader{}

	// Copy each field from the ROM
	copy(header.EntryPoint[:], rom[0x0100:0x0104])
	copy(header.NintendoLogo[:], rom[0x0104:0x0134])
	copy(header.Title[:], rom[0x0134:0x0143])
	copy(header.ManufacturerCode[:], rom[0x013F:0x0143])
	header.CGBFlag = rom[0x0143]
	copy(header.NewLicenseeCode[:], rom[0x0144:0x0146])
	header.SGBFlag = rom[0x0146]
	header.CartridgeType = rom[0x0147]
	header.ROMSize = rom[0x0148]
	header.RAMSize = rom[0x0149]
	header.DestinationCode = rom[0x014A]
	header.OldLicenseeCode = rom[0x014B]
	header.MaskROMVersion = rom[0x014C]
	header.HeaderChecksum = rom[0x014D]
	copy(header.GlobalChecksum[:], rom[0x014E:0x0150])

	return header, nil
}

// GetTitle returns the title as a string, trimmed of null bytes.
func (h *CartridgeHeader) GetTitle() string {
	title := make([]byte, 0, len(h.Title))
	for _, b := range h.Title {
		if b == 0 {
			break
		}
		title = append(title, b)
	}
	return string(title)
}

// GetCartridgeTypeName returns the human-readable cartridge type name.
func (h *CartridgeHeader) GetCartridgeTypeName() string {
	if name := h.getBasicCartridgeTypeName(); name != "" {
		return name
	}
	if name := h.getMBCCartridgeTypeName(); name != "" {
		return name
	}
	if name := h.getSpecialCartridgeTypeName(); name != "" {
		return name
	}
	return fmt.Sprintf("UNKNOWN (0x%02X)", h.CartridgeType)
}

// getBasicCartridgeTypeName returns names for basic cartridge types.
func (h *CartridgeHeader) getBasicCartridgeTypeName() string {
	switch h.CartridgeType {
	case CartridgeTypeROMOnly:
		return "ROM ONLY"
	case CartridgeTypeROMRAM:
		return "ROM+RAM"
	case CartridgeTypeROMRAMBattery:
		return "ROM+RAM+BATTERY"
	case CartridgeTypeMMM01:
		return "MMM01"
	case CartridgeTypeMMM01RAM:
		return "MMM01+RAM"
	case CartridgeTypeMMM01RAMBattery:
		return "MMM01+RAM+BATTERY"
	}
	return ""
}

// getMBCCartridgeTypeName returns names for MBC cartridge types.
func (h *CartridgeHeader) getMBCCartridgeTypeName() string {
	if name := h.getMBC1And2CartridgeTypeName(); name != "" {
		return name
	}
	if name := h.getMBC3CartridgeTypeName(); name != "" {
		return name
	}
	if name := h.getMBC5AndHigherCartridgeTypeName(); name != "" {
		return name
	}
	return ""
}

// getMBC1And2CartridgeTypeName returns names for MBC1 and MBC2 cartridge types.
func (h *CartridgeHeader) getMBC1And2CartridgeTypeName() string {
	switch h.CartridgeType {
	case CartridgeTypeMBC1:
		return "MBC1"
	case CartridgeTypeMBC1RAM:
		return "MBC1+RAM"
	case CartridgeTypeMBC1RAMBattery:
		return "MBC1+RAM+BATTERY"
	case CartridgeTypeMBC2:
		return "MBC2"
	case CartridgeTypeMBC2Battery:
		return "MBC2+BATTERY"
	}
	return ""
}

// getMBC3CartridgeTypeName returns names for MBC3 cartridge types.
func (h *CartridgeHeader) getMBC3CartridgeTypeName() string {
	switch h.CartridgeType {
	case CartridgeTypeMBC3TimerBattery:
		return "MBC3+TIMER+BATTERY"
	case CartridgeTypeMBC3TimerRAMBattery:
		return "MBC3+TIMER+RAM+BATTERY"
	case CartridgeTypeMBC3:
		return "MBC3"
	case CartridgeTypeMBC3RAM:
		return "MBC3+RAM"
	case CartridgeTypeMBC3RAMBattery:
		return "MBC3+RAM+BATTERY"
	}
	return ""
}

// getMBC5AndHigherCartridgeTypeName returns names for MBC5+ cartridge types.
func (h *CartridgeHeader) getMBC5AndHigherCartridgeTypeName() string {
	switch h.CartridgeType {
	case CartridgeTypeMBC5:
		return "MBC5"
	case CartridgeTypeMBC5RAM:
		return "MBC5+RAM"
	case CartridgeTypeMBC5RAMBattery:
		return "MBC5+RAM+BATTERY"
	case CartridgeTypeMBC5Rumble:
		return "MBC5+RUMBLE"
	case CartridgeTypeMBC5RumbleRAM:
		return "MBC5+RUMBLE+RAM"
	case CartridgeTypeMBC5RumbleRAMBattery:
		return "MBC5+RUMBLE+RAM+BATTERY"
	case CartridgeTypeMBC6:
		return "MBC6"
	case CartridgeTypeMBC7SensorRumbleRAMBattery:
		return "MBC7+SENSOR+RUMBLE+RAM+BATTERY"
	}
	return ""
}

// getSpecialCartridgeTypeName returns names for special cartridge types.
func (h *CartridgeHeader) getSpecialCartridgeTypeName() string {
	switch h.CartridgeType {
	case CartridgeTypePocketCamera:
		return "POCKET CAMERA"
	case CartridgeTypeBandaiTAMA5:
		return "BANDAI TAMA5"
	case CartridgeTypeHuC3:
		return "HuC3"
	case CartridgeTypeHuC1RAMBattery:
		return "HuC1+RAM+BATTERY"
	}
	return ""
}

// GetROMSizeKB returns the ROM size in kilobytes.
func (h *CartridgeHeader) GetROMSizeKB() int {
	switch h.ROMSize {
	case ROMSize32KB:
		return 32
	case ROMSize64KB:
		return 64
	case ROMSize128KB:
		return 128
	case ROMSize256KB:
		return 256
	case ROMSize512KB:
		return 512
	case ROMSize1MB:
		return 1024
	case ROMSize2MB:
		return 2048
	case ROMSize4MB:
		return 4096
	case ROMSize8MB:
		return 8192
	default:
		return 0
	}
}

// GetRAMSizeKB returns the RAM size in kilobytes.
func (h *CartridgeHeader) GetRAMSizeKB() int {
	switch h.RAMSize {
	case RAMSizeNone:
		return 0
	case RAMSize2KB:
		return 2
	case RAMSize8KB:
		return 8
	case RAMSize32KB:
		return 32
	case RAMSize128KB:
		return 128
	case RAMSize64KB:
		return 64
	default:
		return 0
	}
}

// HasBattery returns true if the cartridge has battery backup.
func (h *CartridgeHeader) HasBattery() bool {
	switch h.CartridgeType {
	case CartridgeTypeMBC1RAMBattery,
		CartridgeTypeMBC2Battery,
		CartridgeTypeROMRAMBattery,
		CartridgeTypeMMM01RAMBattery,
		CartridgeTypeMBC3TimerBattery,
		CartridgeTypeMBC3TimerRAMBattery,
		CartridgeTypeMBC3RAMBattery,
		CartridgeTypeMBC5RAMBattery,
		CartridgeTypeMBC5RumbleRAMBattery,
		CartridgeTypeMBC7SensorRumbleRAMBattery,
		CartridgeTypeHuC1RAMBattery:
		return true
	default:
		return false
	}
}

// HasRAM returns true if the cartridge has external RAM.
func (h *CartridgeHeader) HasRAM() bool {
	switch h.CartridgeType {
	case CartridgeTypeMBC1RAM,
		CartridgeTypeMBC1RAMBattery,
		CartridgeTypeROMRAM,
		CartridgeTypeROMRAMBattery,
		CartridgeTypeMMM01RAM,
		CartridgeTypeMMM01RAMBattery,
		CartridgeTypeMBC3TimerRAMBattery,
		CartridgeTypeMBC3RAM,
		CartridgeTypeMBC3RAMBattery,
		CartridgeTypeMBC5RAM,
		CartridgeTypeMBC5RAMBattery,
		CartridgeTypeMBC5RumbleRAM,
		CartridgeTypeMBC5RumbleRAMBattery,
		CartridgeTypeMBC7SensorRumbleRAMBattery,
		CartridgeTypeHuC1RAMBattery:
		return true
	case CartridgeTypeMBC2,
		CartridgeTypeMBC2Battery:
		return true // MBC2 has built-in RAM
	default:
		return false
	}
}

// HasTimer returns true if the cartridge has a real-time clock.
func (h *CartridgeHeader) HasTimer() bool {
	switch h.CartridgeType {
	case CartridgeTypeMBC3TimerBattery,
		CartridgeTypeMBC3TimerRAMBattery:
		return true
	default:
		return false
	}
}

// GetMBCType returns the MBC type number.
func (h *CartridgeHeader) GetMBCType() int {
	switch h.CartridgeType {
	case CartridgeTypeROMOnly,
		CartridgeTypeROMRAM,
		CartridgeTypeROMRAMBattery:
		return 0
	case CartridgeTypeMBC1,
		CartridgeTypeMBC1RAM,
		CartridgeTypeMBC1RAMBattery:
		return 1
	case CartridgeTypeMBC2,
		CartridgeTypeMBC2Battery:
		return 2
	case CartridgeTypeMBC3TimerBattery,
		CartridgeTypeMBC3TimerRAMBattery,
		CartridgeTypeMBC3,
		CartridgeTypeMBC3RAM,
		CartridgeTypeMBC3RAMBattery:
		return 3
	case CartridgeTypeMBC5,
		CartridgeTypeMBC5RAM,
		CartridgeTypeMBC5RAMBattery,
		CartridgeTypeMBC5Rumble,
		CartridgeTypeMBC5RumbleRAM,
		CartridgeTypeMBC5RumbleRAMBattery:
		return 5
	case CartridgeTypeMBC6:
		return 6
	case CartridgeTypeMBC7SensorRumbleRAMBattery:
		return 7
	default:
		return -1 // Unknown or special type
	}
}

// ValidateHeaderChecksum validates the header checksum.
func (h *CartridgeHeader) ValidateHeaderChecksum(rom []byte) bool {
	if len(rom) < 0x150 {
		return false
	}

	var checksum uint8
	for i := 0x0134; i <= 0x014C; i++ {
		checksum = checksum - rom[i] - 1
	}

	return checksum == h.HeaderChecksum
}
