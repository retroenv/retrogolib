package x86

import (
	"fmt"

	"github.com/retroenv/retrogolib/log"
)

// Memory represents the x86 memory system with segmented addressing.
// The x86 uses a 20-bit address space (1MB) accessed through 16-bit segments.
type Memory struct {
	data   []uint8
	size   uint32
	logger *log.Logger
}

// Memory size constants
const (
	MaxMemorySize = 1024 * 1024 // 1MB maximum for 8086/8088
	MinMemorySize = 64 * 1024   // 64KB minimum
	SegmentSize   = 64 * 1024   // 64KB segment size
	AddressMask   = 0x000FFFFF  // 20-bit address mask
)

// NewMemory creates a new memory instance with the specified size.
// Pre-allocates memory buffer with exact capacity for optimal performance.
func NewMemory(size uint32, logger *log.Logger) (*Memory, error) {
	if size < MinMemorySize {
		return nil, fmt.Errorf("memory size %d is below minimum %d", size, MinMemorySize)
	}
	if size > MaxMemorySize {
		return nil, fmt.Errorf("memory size %d exceeds maximum %d", size, MaxMemorySize)
	}

	return &Memory{
		data:   make([]uint8, size),
		size:   size,
		logger: logger,
	}, nil
}

// Size returns the total memory size in bytes.
func (m *Memory) Size() uint32 {
	return m.size
}

// Data returns a copy of the memory data.
// Pre-allocates buffer with exact capacity for optimal performance.
func (m *Memory) Data() []uint8 {
	data := make([]uint8, len(m.data))
	copy(data, m.data)
	return data
}

// Clear fills memory with the specified value.
func (m *Memory) Clear(value uint8) {
	for i := range m.data {
		m.data[i] = value
	}
}

// Read8 reads a byte from the specified linear address.
func (m *Memory) Read8(addr uint32) uint8 {
	addr &= AddressMask // Mask to 20-bit address space
	if addr >= m.size {
		if m.logger != nil {
			m.logger.Debug("memory read beyond bounds",
				log.String("address", fmt.Sprintf("0x%06X", addr)),
				log.String("size", fmt.Sprintf("0x%06X", m.size)))
		}
		return 0xFF // Return default value for out-of-bounds reads
	}
	return m.data[addr]
}

// Read16 reads a word (16-bit) from the specified linear address.
// Uses little-endian byte order (low byte first).
func (m *Memory) Read16(addr uint32) uint16 {
	low := uint16(m.Read8(addr))
	high := uint16(m.Read8(addr + 1))
	return high<<8 | low
}

// Write8 writes a byte to the specified linear address.
func (m *Memory) Write8(addr uint32, value uint8) {
	addr &= AddressMask // Mask to 20-bit address space
	if addr >= m.size {
		if m.logger != nil {
			m.logger.Debug("memory write beyond bounds",
				log.String("address", fmt.Sprintf("0x%06X", addr)),
				log.String("value", fmt.Sprintf("0x%02X", value)),
				log.String("size", fmt.Sprintf("0x%06X", m.size)))
		}
		return // Ignore out-of-bounds writes
	}
	m.data[addr] = value
}

// Write16 writes a word (16-bit) to the specified linear address.
// Uses little-endian byte order (low byte first).
func (m *Memory) Write16(addr uint32, value uint16) {
	m.Write8(addr, uint8(value))      // Low byte
	m.Write8(addr+1, uint8(value>>8)) // High byte
}

// ReadSegmented reads a byte using segment:offset addressing.
func (m *Memory) ReadSegmented(segment, offset uint16) uint8 {
	addr := uint32(segment)<<4 + uint32(offset)
	return m.Read8(addr)
}

// ReadSegmented16 reads a word using segment:offset addressing.
func (m *Memory) ReadSegmented16(segment, offset uint16) uint16 {
	addr := uint32(segment)<<4 + uint32(offset)
	return m.Read16(addr)
}

// WriteSegmented writes a byte using segment:offset addressing.
func (m *Memory) WriteSegmented(segment, offset uint16, value uint8) {
	addr := uint32(segment)<<4 + uint32(offset)
	m.Write8(addr, value)
}

// WriteSegmented16 writes a word using segment:offset addressing.
func (m *Memory) WriteSegmented16(segment, offset uint16, value uint16) {
	addr := uint32(segment)<<4 + uint32(offset)
	m.Write16(addr, value)
}

// LoadData loads data into memory at the specified linear address.
func (m *Memory) LoadData(addr uint32, data []uint8) error {
	if addr >= m.size {
		return fmt.Errorf("load address 0x%06X is beyond memory size 0x%06X", addr, m.size)
	}
	if addr+uint32(len(data)) > m.size {
		return fmt.Errorf("load data exceeds memory bounds: addr=0x%06X, len=%d, size=0x%06X",
			addr, len(data), m.size)
	}

	copy(m.data[addr:], data)

	if m.logger != nil {
		m.logger.Debug("loaded data into memory",
			log.String("address", fmt.Sprintf("0x%06X", addr)),
			log.Int("size", len(data)))
	}

	return nil
}

// LoadSegmentedData loads data into memory using segment:offset addressing.
func (m *Memory) LoadSegmentedData(segment, offset uint16, data []uint8) error {
	addr := uint32(segment)<<4 + uint32(offset)
	return m.LoadData(addr, data)
}

// Dump returns a formatted hex dump of memory from start to end addresses.
func (m *Memory) Dump(start, end uint32) []string {
	if start >= m.size {
		return nil
	}
	if end > m.size {
		end = m.size
	}

	const bytesPerLine = 16
	lines := make([]string, 0, (end-start+bytesPerLine-1)/bytesPerLine)

	for addr := start; addr < end; addr += bytesPerLine {
		var line string
		line = fmt.Sprintf("%06X: ", addr)

		// Hex bytes
		for i := range bytesPerLine {
			if addr+uint32(i) < end {
				line += fmt.Sprintf("%02X ", m.data[addr+uint32(i)])
			} else {
				line += "   "
			}
		}

		// ASCII representation
		line += " |"
		for i := range bytesPerLine {
			if addr+uint32(i) >= end {
				break
			}
			b := m.data[addr+uint32(i)]
			if b >= 32 && b <= 126 {
				line += string(rune(b))
			} else {
				line += "."
			}
		}
		line += "|"

		lines = append(lines, line)
	}

	return lines
}

// ValidateAddress checks if an address is valid within memory bounds.
// For x86 compatibility, this checks both the original address and the masked 20-bit address.
func (m *Memory) ValidateAddress(addr uint32) bool {
	// Check if the original address is within the actual memory size
	if addr >= m.size {
		return false
	}
	// Also check the masked address (20-bit address space)
	maskedAddr := addr & AddressMask
	return maskedAddr < m.size
}

// ValidateSegmentedAddress checks if a segment:offset address is valid.
func (m *Memory) ValidateSegmentedAddress(segment, offset uint16) bool {
	addr := uint32(segment)<<4 + uint32(offset)
	return m.ValidateAddress(addr)
}

// GetLinearAddress converts segment:offset to linear address.
func (m *Memory) GetLinearAddress(segment, offset uint16) uint32 {
	return (uint32(segment)<<4 + uint32(offset)) & AddressMask
}
