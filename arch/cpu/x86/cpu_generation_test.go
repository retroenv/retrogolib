package x86

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestMinimumCPUByEncoding(t *testing.T) {
	for _, test := range []struct {
		opcode  uint16
		minimum CPUGeneration
	}{
		{
			opcode:  0x50,
			minimum: CPU8086,
		},
		{
			opcode:  0x68,
			minimum: CPU80186,
		},
		{
			opcode:  0xd1,
			minimum: CPU8086,
		},
		{
			opcode:  0xc1,
			minimum: CPU80186,
		},
		{
			opcode:  0xc8,
			minimum: CPU80186,
		},
		{
			opcode:  0x0f01,
			minimum: CPU80286,
		},
		{
			opcode:  0x66,
			minimum: CPU80386,
		},
		{
			opcode:  0x0fb6,
			minimum: CPU80386,
		},
		{
			opcode:  0x0fc8,
			minimum: CPU80486,
		},
		{
			opcode:  0xffff,
			minimum: CPUUnknown,
		},
	} {
		info := OpcodeInfo{
			Opcode: test.opcode,
			Size:   1,
		}
		assert.Equal(t, test.minimum, info.MinimumCPU())
		for cpu := CPUUnknown; cpu <= CPU80486; cpu++ {
			assert.Equal(t, test.minimum != CPUUnknown && cpu >= test.minimum, cpu.Supports(info))
		}
		assert.False(t, CPUGeneration(255).Supports(info))
	}
	assert.Equal(t, CPUUnknown, (OpcodeInfo{}).MinimumCPU())
	assert.Equal(t, "unknown", CPUGeneration(255).String())
	for cpu := CPU8086; cpu <= CPU80486; cpu++ {
		assert.NotEqual(t, "unknown", cpu.String())
	}
}

func TestRegisteredInstructionEncodingsHaveCPUClassification(t *testing.T) {
	for name, instruction := range Instructions {
		for mode, info := range instruction.Addressing {
			assert.NotEqual(t, CPUUnknown, info.MinimumCPU(), "%s %s", name, mode)
		}
		for register, info := range instruction.RegisterOpcodes {
			assert.NotEqual(t, CPUUnknown, info.MinimumCPU(), "%s %s", name, register)
		}
	}
	assert.True(t, CPU8086.Supports(PushReg16.Addressing[RegisterAddressing]))
	assert.False(t, CPU8086.Supports(PushImm16.Addressing[ImmediateAddressing]))
	assert.True(t, CPU80186.Supports(PushImm16.Addressing[ImmediateAddressing]))
}
