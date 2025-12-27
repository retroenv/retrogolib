package x86

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/log"
)

func TestNewMemory(t *testing.T) {
	logger := log.NewTestLogger(t)

	tests := []struct {
		name        string
		size        uint32
		expectError bool
	}{
		{
			name:        "valid size 1MB",
			size:        1024 * 1024,
			expectError: false,
		},
		{
			name:        "minimum size 64KB",
			size:        64 * 1024,
			expectError: false,
		},
		{
			name:        "below minimum size",
			size:        32 * 1024,
			expectError: true,
		},
		{
			name:        "above maximum size",
			size:        2 * 1024 * 1024,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memory, err := NewMemory(tt.size, logger)

			if tt.expectError {
				assert.NotNil(t, err)
				assert.Nil(t, memory)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, memory)
				assert.Equal(t, tt.size, memory.Size())
			}
		})
	}
}

func TestMemory_ReadWrite8(t *testing.T) {
	logger := log.NewTestLogger(t)
	memory, err := NewMemory(65536, logger)
	assert.NoError(t, err)

	tests := []struct {
		addr  uint32
		value uint8
	}{
		{0x0000, 0x00},
		{0x0001, 0xFF},
		{0x0100, 0x42},
		{0xFFFF, 0xAB}, // Last valid address
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			memory.Write8(tt.addr, tt.value)
			result := memory.Read8(tt.addr)
			assert.Equal(t, tt.value, result)
		})
	}
}

func TestMemory_ReadWrite16(t *testing.T) {
	logger := log.NewTestLogger(t)
	memory, err := NewMemory(65536, logger)
	assert.NoError(t, err)

	tests := []struct {
		addr  uint32
		value uint16
	}{
		{0x0000, 0x0000},
		{0x0002, 0xFFFF},
		{0x0100, 0x1234},
		{0x0200, 0xABCD},
		{0xFFFE, 0x5678}, // Last valid 16-bit address
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			memory.Write16(tt.addr, tt.value)
			result := memory.Read16(tt.addr)
			assert.Equal(t, tt.value, result)
		})
	}
}

func TestMemory_LittleEndian(t *testing.T) {
	logger := log.NewTestLogger(t)
	memory, err := NewMemory(65536, logger)
	assert.NoError(t, err)

	// Write 16-bit value 0x1234
	memory.Write16(0x0100, 0x1234)

	// Check little-endian byte order
	lowByte := memory.Read8(0x0100)  // Should be 0x34 (low byte)
	highByte := memory.Read8(0x0101) // Should be 0x12 (high byte)

	assert.Equal(t, uint8(0x34), lowByte)
	assert.Equal(t, uint8(0x12), highByte)
}

func TestMemory_SegmentedAddressing(t *testing.T) {
	logger := log.NewTestLogger(t)
	memory, err := NewMemory(1024*1024, logger)
	assert.NoError(t, err)

	tests := []struct {
		segment uint16
		offset  uint16
		value   uint8
	}{
		{0x0000, 0x0000, 0x11},
		{0x1000, 0x0000, 0x22},
		{0x0000, 0x1000, 0x33},
		{0x1234, 0x5678, 0x44},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			memory.WriteSegmented(tt.segment, tt.offset, tt.value)
			result := memory.ReadSegmented(tt.segment, tt.offset)
			assert.Equal(t, tt.value, result)
		})
	}
}

func TestMemory_SegmentedAddressing16(t *testing.T) {
	logger := log.NewTestLogger(t)
	memory, err := NewMemory(1024*1024, logger)
	assert.NoError(t, err)

	tests := []struct {
		segment uint16
		offset  uint16
		value   uint16
	}{
		{0x0000, 0x0000, 0x1122},
		{0x1000, 0x0000, 0x3344},
		{0x0000, 0x1000, 0x5566},
		{0x1234, 0x5678, 0x7788},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			memory.WriteSegmented16(tt.segment, tt.offset, tt.value)
			result := memory.ReadSegmented16(tt.segment, tt.offset)
			assert.Equal(t, tt.value, result)
		})
	}
}

func TestMemory_OutOfBounds(t *testing.T) {
	logger := log.NewTestLogger(t)
	memory, err := NewMemory(65536, logger)
	assert.NoError(t, err)

	// Test read beyond bounds returns default value
	result := memory.Read8(70000)
	assert.Equal(t, uint8(0xFF), result)

	// Test write beyond bounds is ignored (no crash)
	memory.Write8(70000, 0x42)
	// Should not crash or affect memory
}

func TestMemory_LoadData(t *testing.T) {
	logger := log.NewTestLogger(t)
	memory, err := NewMemory(65536, logger)
	assert.NoError(t, err)

	data := []uint8{0x01, 0x02, 0x03, 0x04, 0x05}

	err = memory.LoadData(0x100, data)
	assert.NoError(t, err)

	// Verify data was loaded correctly
	for i, expectedValue := range data {
		actualValue := memory.Read8(uint32(0x100 + i))
		assert.Equal(t, expectedValue, actualValue)
	}
}

func TestMemory_LoadDataErrors(t *testing.T) {
	logger := log.NewTestLogger(t)
	memory, err := NewMemory(65536, logger)
	assert.NoError(t, err)

	tests := []struct {
		name string
		addr uint32
		data []uint8
	}{
		{
			name: "address beyond bounds",
			addr: 70000,
			data: []uint8{0x01, 0x02},
		},
		{
			name: "data exceeds bounds",
			addr: 65530,
			data: []uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := memory.LoadData(tt.addr, tt.data)
			assert.NotNil(t, err)
		})
	}
}

func TestMemory_LoadSegmentedData(t *testing.T) {
	logger := log.NewTestLogger(t)
	memory, err := NewMemory(1024*1024, logger)
	assert.NoError(t, err)

	data := []uint8{0xDE, 0xAD, 0xBE, 0xEF}

	err = memory.LoadSegmentedData(0x1000, 0x0100, data)
	assert.NoError(t, err)

	// Verify data was loaded at correct linear address
	linearAddr := uint32(0x1000)<<4 + uint32(0x0100)
	for i, expectedValue := range data {
		actualValue := memory.Read8(linearAddr + uint32(i))
		assert.Equal(t, expectedValue, actualValue)
	}
}

func TestMemory_Clear(t *testing.T) {
	logger := log.NewTestLogger(t)
	memory, err := NewMemory(65536, logger)
	assert.NoError(t, err)

	// Write some data
	memory.Write8(0x100, 0x42)
	memory.Write8(0x200, 0x43)
	memory.Write8(0x300, 0x44)

	// Clear memory with specific value
	memory.Clear(0xAA)

	// Verify memory is cleared
	assert.Equal(t, uint8(0xAA), memory.Read8(0x100))
	assert.Equal(t, uint8(0xAA), memory.Read8(0x200))
	assert.Equal(t, uint8(0xAA), memory.Read8(0x300))
	assert.Equal(t, uint8(0xAA), memory.Read8(0x000))
}

func TestMemory_ValidateAddress(t *testing.T) {
	logger := log.NewTestLogger(t)
	memory, err := NewMemory(65536, logger)
	assert.NoError(t, err)

	tests := []struct {
		addr  uint32
		valid bool
	}{
		{0x000, true},
		{0xFFFF, true},    // Last valid address
		{0x10000, false},  // Beyond memory size
		{0x100000, false}, // Far beyond
		{0xFFFFFF, false}, // Maximum possible
	}

	for _, tt := range tests {
		result := memory.ValidateAddress(tt.addr)
		assert.Equal(t, tt.valid, result)
	}
}

func TestMemory_ValidateSegmentedAddress(t *testing.T) {
	logger := log.NewTestLogger(t)
	memory, err := NewMemory(1024*1024, logger)
	assert.NoError(t, err)

	tests := []struct {
		segment uint16
		offset  uint16
		valid   bool
	}{
		{0x0000, 0x0000, true},
		{0xF000, 0xFFFF, true},  // 0xFFFFF = 1MB-1, valid for 1MB memory
		{0x0000, 0xFFFF, true},  // Valid within 1MB
		{0xFFFF, 0x000F, true},  // 0xFFFFF = 1MB-1, valid for 1MB memory
		{0xFFFF, 0x0010, false}, // 0x100000 = 1MB, exceeds 1MB memory bounds
	}

	for _, tt := range tests {
		result := memory.ValidateSegmentedAddress(tt.segment, tt.offset)
		assert.Equal(t, tt.valid, result)
	}
}

func TestMemory_GetLinearAddress(t *testing.T) {
	logger := log.NewTestLogger(t)
	memory, err := NewMemory(1024*1024, logger)
	assert.NoError(t, err)

	tests := []struct {
		segment  uint16
		offset   uint16
		expected uint32
	}{
		{0x0000, 0x0000, 0x00000},
		{0x1000, 0x0000, 0x10000},
		{0x0000, 0x1000, 0x01000},
		{0x1234, 0x5678, 0x179B8},
		{0xFFFF, 0x000F, 0xFFFF0 + 0x000F}, // Maximum valid
	}

	for _, tt := range tests {
		result := memory.GetLinearAddress(tt.segment, tt.offset)
		// Mask the result to 20-bit like the implementation does
		expected := tt.expected & AddressMask
		assert.Equal(t, expected, result)
	}
}

func TestMemory_Data(t *testing.T) {
	logger := log.NewTestLogger(t)
	memory, err := NewMemory(65536, logger)
	assert.NoError(t, err)

	// Write some test data
	memory.Write8(0x000, 0x11)
	memory.Write8(0x001, 0x22)
	memory.Write8(0x002, 0x33)

	// Get data copy
	data := memory.Data()

	// Verify data is correct
	assert.Equal(t, uint8(0x11), data[0x000])
	assert.Equal(t, uint8(0x22), data[0x001])
	assert.Equal(t, uint8(0x33), data[0x002])

	// Modify copy should not affect original
	data[0x000] = 0xFF
	assert.Equal(t, uint8(0x11), memory.Read8(0x000))
}

func TestMemory_Dump(t *testing.T) {
	logger := log.NewTestLogger(t)
	memory, err := NewMemory(65536, logger)
	assert.NoError(t, err)

	// Write test pattern
	for i := range uint32(32) {
		memory.Write8(i, uint8(i))
	}

	dump := memory.Dump(0, 32)
	assert.NotEmpty(t, dump)
	assert.True(t, len(dump) >= 2) // Should have at least 2 lines for 32 bytes
}
