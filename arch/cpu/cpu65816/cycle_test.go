package cpu65816

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

// TestCycleRules checks the datasheet cycle rules for register width, direct
// page alignment, indexing, and native mode. The rules were missing before,
// so 16-bit operations and DL!=0 accesses were undercounted.
func TestCycleRules(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(cpu *CPU)
		code   []uint8
		cycles uint64
	}{
		{name: "LDA abs M=1", setup: func(*CPU) {}, code: []uint8{0xAD, 0x00, 0x20}, cycles: 4},
		{name: "LDA abs M=0", setup: func(c *CPU) { c.Flags.M = 0 }, code: []uint8{0xAD, 0x00, 0x20}, cycles: 5},
		{name: "LDA # M=0", setup: func(c *CPU) { c.Flags.M = 0 }, code: []uint8{0xA9, 0x00, 0x20}, cycles: 3},
		{name: "LDX abs X=0", setup: func(c *CPU) { c.Flags.X = 0 }, code: []uint8{0xAE, 0x00, 0x20}, cycles: 5},
		{name: "LDA dp DL=0", setup: func(*CPU) {}, code: []uint8{0xA5, 0x10}, cycles: 3},
		{name: "LDA dp DL!=0", setup: func(c *CPU) { c.DP = 0x0001 }, code: []uint8{0xA5, 0x10}, cycles: 4},
		{name: "LDA dp DL!=0 M=0", setup: func(c *CPU) { c.DP = 0x0001; c.Flags.M = 0 }, code: []uint8{0xA5, 0x10}, cycles: 5},
		{name: "LDA abs,X X=1 no cross", setup: func(c *CPU) { c.X = 1 }, code: []uint8{0xBD, 0x00, 0x20}, cycles: 4},
		{name: "LDA abs,X X=1 cross", setup: func(c *CPU) { c.X = 1 }, code: []uint8{0xBD, 0xFF, 0x20}, cycles: 5},
		{name: "LDA abs,X X=0 no cross", setup: func(c *CPU) { c.Flags.X = 0; c.X = 1 }, code: []uint8{0xBD, 0x00, 0x20}, cycles: 5},
		{name: "STA abs,X X=1 no cross", setup: func(c *CPU) { c.X = 1 }, code: []uint8{0x9D, 0x00, 0x20}, cycles: 5},
		{name: "ASL abs M=0", setup: func(c *CPU) { c.Flags.M = 0 }, code: []uint8{0x0E, 0x00, 0x20}, cycles: 8},
		{name: "ASL A M=0", setup: func(c *CPU) { c.Flags.M = 0 }, code: []uint8{0x0A}, cycles: 2},
		{name: "TSB dp DL!=0 M=0", setup: func(c *CPU) { c.DP = 0x0001; c.Flags.M = 0 }, code: []uint8{0x04, 0x10}, cycles: 8},
		{name: "PHA M=0", setup: func(c *CPU) { c.Flags.M = 0 }, code: []uint8{0x48}, cycles: 4},
		{name: "PLY X=0", setup: func(c *CPU) { c.Flags.X = 0 }, code: []uint8{0x7A}, cycles: 5},
		{name: "PEI DL!=0", setup: func(c *CPU) { c.DP = 0x0001 }, code: []uint8{0xD4, 0x10}, cycles: 7},
		{name: "BRK native", setup: func(*CPU) {}, code: []uint8{0x00, 0x00}, cycles: 8},
		{name: "BRK emulation", setup: func(c *CPU) { c.E = true }, code: []uint8{0x00, 0x00}, cycles: 7},
		{name: "RTI native", setup: func(*CPU) {}, code: []uint8{0x40}, cycles: 7},
		{name: "COP emulation", setup: func(c *CPU) { c.E = true }, code: []uint8{0x02, 0x00}, cycles: 7},
		{name: "MVN one byte", setup: func(c *CPU) { c.C = 0 }, code: []uint8{0x54, 0x00, 0x00}, cycles: 7},
		{name: "XCE", setup: func(*CPU) {}, code: []uint8{0xFB}, cycles: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu, mem := setupCPU(t)
			tt.setup(cpu)
			writeOp(mem, 0x8000, tt.code...)
			before := cpu.cycles

			assert.NoError(t, cpu.Step())
			assert.Equal(t, tt.cycles, cpu.cycles-before)
		})
	}
}

// TestExtraCyclesUsesWidthBeforeExecution checks that a width change made by
// the instruction does not change its own cycle count.
func TestExtraCyclesUsesWidthBeforeExecution(t *testing.T) {
	cpu, mem := setupCPU(t)
	cpu.Flags.M = 0
	mem.data[0x01FF] = 0x20 // P with M=1
	cpu.SP = 0x01FE
	writeOp(mem, 0x8000, 0x28) // PLP
	before := cpu.cycles

	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint64(4), cpu.cycles-before)
	assert.Equal(t, uint8(1), cpu.Flags.M)
}
