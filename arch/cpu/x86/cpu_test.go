package x86

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/log"
)

func TestNew(t *testing.T) {
	logger := log.NewTestLogger(t)

	tests := []struct {
		name        string
		memory      *Memory
		options     []Option
		expectError bool
	}{
		{
			name:        "valid memory",
			memory:      createTestMemory(t, logger),
			expectError: false,
		},
		{
			name:        "nil memory",
			memory:      nil,
			expectError: true,
		},
		{
			name:        "with DOS defaults",
			memory:      createTestMemory(t, logger),
			options:     []Option{WithDOSDefaults()},
			expectError: false,
		},
		{
			name:        "with BIOS defaults",
			memory:      createTestMemory(t, logger),
			options:     []Option{WithBIOSDefaults()},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu, err := New(tt.memory, tt.options...)

			if tt.expectError {
				assert.ErrorIs(t, err, ErrNilMemory)
				assert.Nil(t, cpu)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cpu)
				assert.Equal(t, uint64(0), cpu.Cycles())
				assert.False(t, cpu.Halted())
			}
		})
	}
}

func TestCPU_State(t *testing.T) {
	logger := log.NewTestLogger(t)
	memory := createTestMemory(t, logger)
	cpu, err := New(memory)
	assert.NoError(t, err)

	// Set some register values
	cpu.AX = 0x1234
	cpu.BX = 0x5678
	cpu.SetCarry(true)
	cpu.SetZero(true)

	state := cpu.State()

	assert.Equal(t, uint16(0x1234), state.AX)
	assert.Equal(t, uint16(0x5678), state.BX)
	assert.True(t, state.Flags.GetCarry())
	assert.True(t, state.Flags.GetZero())
	assert.Equal(t, uint64(0), state.Cycles)
	assert.False(t, state.Halted)
}

func TestCPU_RegisterAccessors(t *testing.T) {
	logger := log.NewTestLogger(t)
	memory := createTestMemory(t, logger)
	cpu, err := New(memory)
	assert.NoError(t, err)

	// Test 8-bit register accessors
	cpu.AX = 0x1234
	assert.Equal(t, uint8(0x34), cpu.AL()) // Low byte
	assert.Equal(t, uint8(0x12), cpu.AH()) // High byte

	cpu.SetAL(0x56)
	assert.Equal(t, uint16(0x1256), cpu.AX)

	cpu.SetAH(0x78)
	assert.Equal(t, uint16(0x7856), cpu.AX)

	// Test other registers
	cpu.BX = 0xABCD
	assert.Equal(t, uint8(0xCD), cpu.BL())
	assert.Equal(t, uint8(0xAB), cpu.BH())

	cpu.CX = 0xEF01
	assert.Equal(t, uint8(0x01), cpu.CL())
	assert.Equal(t, uint8(0xEF), cpu.CH())

	cpu.DX = 0x2345
	assert.Equal(t, uint8(0x45), cpu.DL())
	assert.Equal(t, uint8(0x23), cpu.DH())
}

func TestCPU_SegmentedAddressing(t *testing.T) {
	logger := log.NewTestLogger(t)
	memory := createTestMemory(t, logger)
	cpu, err := New(memory)
	assert.NoError(t, err)

	tests := []struct {
		segment  uint16
		offset   uint16
		expected uint32
	}{
		{0x0000, 0x0000, 0x00000},
		{0x1000, 0x0000, 0x10000},
		{0x0000, 0x1000, 0x01000},
		{0x1234, 0x5678, 0x179B8},          // 0x1234 << 4 + 0x5678 = 0x12340 + 0x5678
		{0xFFFF, 0x000F, 0xFFFF0 + 0x000F}, // Maximum address
	}

	for _, tt := range tests {
		result := cpu.CalculateAddress(tt.segment, tt.offset)
		assert.Equal(t, tt.expected, result)
	}
}

func TestCPU_StackOperations(t *testing.T) {
	logger := log.NewTestLogger(t)
	memory := createTestMemory(t, logger)
	cpu, err := New(memory, WithInitialSS(0x1000), WithInitialSP(0x1000))
	assert.NoError(t, err)

	initialSP := cpu.SP

	// Test 8-bit push/pop
	cpu.push8(0x42)
	assert.Equal(t, initialSP-1, cpu.SP)

	value := cpu.pop8()
	assert.Equal(t, uint8(0x42), value)
	assert.Equal(t, initialSP, cpu.SP)

	// Test 16-bit push/pop
	cpu.push16(0x1234)
	assert.Equal(t, initialSP-2, cpu.SP)

	value16 := cpu.pop16()
	assert.Equal(t, uint16(0x1234), value16)
	assert.Equal(t, initialSP, cpu.SP)
}

func TestCPU_HaltResume(t *testing.T) {
	logger := log.NewTestLogger(t)
	memory := createTestMemory(t, logger)
	cpu, err := New(memory)
	assert.NoError(t, err)

	assert.False(t, cpu.Halted())

	cpu.Halt()
	assert.True(t, cpu.Halted())

	cpu.Resume()
	assert.False(t, cpu.Halted())
}

func TestCPU_InterruptHandling(t *testing.T) {
	logger := log.NewTestLogger(t)
	memory := createTestMemory(t, logger)
	cpu, err := New(memory)
	assert.NoError(t, err)

	// Test interrupt enable/disable
	assert.False(t, cpu.interruptsEnabled)

	cpu.EnableInterrupts()
	assert.True(t, cpu.interruptsEnabled)

	cpu.DisableInterrupts()
	assert.False(t, cpu.interruptsEnabled)

	// Test interrupt triggering
	cpu.TriggerInterrupt(0x21) // DOS interrupt
	assert.True(t, cpu.triggerInt)
	assert.Equal(t, uint8(0x21), cpu.intVector)
}

// Test with DOS defaults
func TestCPU_DOSDefaults(t *testing.T) {
	logger := log.NewTestLogger(t)
	memory := createTestMemory(t, logger)
	cpu, err := New(memory, WithDOSDefaults())
	assert.NoError(t, err)

	assert.Equal(t, uint16(0x1000), cpu.CS)
	assert.Equal(t, uint16(0x1000), cpu.DS)
	assert.Equal(t, uint16(0x1000), cpu.ES)
	assert.Equal(t, uint16(0x2000), cpu.SS)
	assert.Equal(t, uint16(0xFFFE), cpu.SP)
	assert.Equal(t, uint16(0x0100), cpu.IP)
	assert.True(t, cpu.opts.interruptEnabled)
}

// Test with BIOS defaults
func TestCPU_BIOSDefaults(t *testing.T) {
	logger := log.NewTestLogger(t)
	memory := createTestMemory(t, logger)
	cpu, err := New(memory, WithBIOSDefaults())
	assert.NoError(t, err)

	assert.Equal(t, uint16(0xF000), cpu.CS)
	assert.Equal(t, uint16(0x0000), cpu.DS)
	assert.Equal(t, uint16(0x0000), cpu.ES)
	assert.Equal(t, uint16(0x0000), cpu.SS)
	assert.Equal(t, uint16(0x0400), cpu.SP)
	assert.Equal(t, uint16(0xFFF0), cpu.IP)
	assert.False(t, cpu.opts.interruptEnabled)
}

func TestCPU_RegisterHelpers(t *testing.T) {
	logger := log.NewTestLogger(t)
	memory := createTestMemory(t, logger)
	cpu, err := New(memory)
	assert.NoError(t, err)

	// Test 8-bit register access
	cpu.AX = 0x1234
	assert.Equal(t, uint8(0x34), cpu.getReg8(RegAL))
	assert.Equal(t, uint8(0x12), cpu.getReg8(RegAH))

	cpu.setReg8(RegAL, 0x56)
	assert.Equal(t, uint16(0x1256), cpu.AX)

	// Test 16-bit register access
	cpu.setReg16(RegBX, 0xABCD)
	assert.Equal(t, uint16(0xABCD), cpu.getReg16(RegBX))
	assert.Equal(t, uint16(0xABCD), cpu.BX)

	// Test segment register access
	cpu.setReg16(RegDS, 0x2000)
	assert.Equal(t, uint16(0x2000), cpu.getReg16(RegDS))
	assert.Equal(t, uint16(0x2000), cpu.DS)
}

// createTestMemory creates a test memory instance.
func createTestMemory(t *testing.T, logger *log.Logger) *Memory {
	t.Helper()
	memory, err := NewMemory(1024*1024, logger) // 1MB
	assert.NoError(t, err)
	return memory
}

// Benchmark tests
func BenchmarkCPU_New(b *testing.B) {
	logger := log.NewTestLogger(b)
	memory := createTestMemoryForBench(b, logger)

	b.ResetTimer()
	for range b.N {
		_, _ = New(memory)
	}
}

func BenchmarkCPU_RegisterAccess(b *testing.B) {
	logger := log.NewTestLogger(b)
	memory := createTestMemoryForBench(b, logger)
	cpu, _ := New(memory)

	b.ResetTimer()
	for i := range b.N {
		cpu.SetAL(uint8(i))
		_ = cpu.AL()
	}
}

func BenchmarkCPU_StackOperations(b *testing.B) {
	logger := log.NewTestLogger(b)
	memory := createTestMemoryForBench(b, logger)
	cpu, _ := New(memory, WithInitialSS(0x1000), WithInitialSP(0x8000))

	b.ResetTimer()
	for i := range b.N {
		cpu.push16(uint16(i))
		_ = cpu.pop16()
	}
}

func createTestMemoryForBench(b *testing.B, logger *log.Logger) *Memory {
	b.Helper()
	memory, err := NewMemory(1024*1024, logger)
	if err != nil {
		b.Fatal(err)
	}
	return memory
}
