package x86

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

// TestEAClocks checks each row of table 2-20 of the Intel 8086 Family User's
// Manual.
func TestEAClocks(t *testing.T) {
	tests := []struct {
		name   string
		mod    uint8
		rm     uint8
		clocks uint8
	}{
		{name: "displacement only", mod: 0, rm: 6, clocks: 6},
		{name: "SI", mod: 0, rm: 4, clocks: 5},
		{name: "DI", mod: 0, rm: 5, clocks: 5},
		{name: "BX", mod: 0, rm: 7, clocks: 5},
		{name: "BP plus displacement", mod: 1, rm: 6, clocks: 9},
		{name: "SI plus word displacement", mod: 2, rm: 4, clocks: 9},
		{name: "BX plus SI", mod: 0, rm: 0, clocks: 7},
		{name: "BP plus DI", mod: 0, rm: 3, clocks: 7},
		{name: "BX plus DI", mod: 0, rm: 1, clocks: 8},
		{name: "BP plus SI", mod: 0, rm: 2, clocks: 8},
		{name: "BX plus SI plus displacement", mod: 1, rm: 0, clocks: 11},
		{name: "BP plus DI plus displacement", mod: 2, rm: 3, clocks: 11},
		{name: "BX plus DI plus displacement", mod: 1, rm: 1, clocks: 12},
		{name: "BP plus SI plus displacement", mod: 2, rm: 2, clocks: 12},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			clocks, ok := EAClocks(test.mod, test.rm)
			assert.True(t, ok)
			assert.Equal(t, test.clocks, clocks)
		})
	}

	_, ok := EAClocks(3, 0)
	assert.False(t, ok)
	_, ok = EAClocks(0, 8)
	assert.False(t, ok)
}

// TestInstructionTiming checks entries of table 2-21 of the Intel 8086 Family
// User's Manual.
func TestInstructionTiming(t *testing.T) {
	tests := []struct {
		name      string
		form      OperandForm
		clocks    uint8
		maxClocks uint8
		perCount  uint8
		transfers uint8
		ea        bool
	}{
		{name: AddName, form: FormRegisterRegister, clocks: 3, maxClocks: 3, perCount: 0, transfers: 0, ea: false},
		{name: AddName, form: FormMemoryRegister, clocks: 16, maxClocks: 16, perCount: 0, transfers: 2, ea: true},
		{name: SubName, form: FormMemoryImmediate, clocks: 17, maxClocks: 17, perCount: 0, transfers: 2, ea: true},
		{name: CmpName, form: FormMemoryRegister, clocks: 9, maxClocks: 9, perCount: 0, transfers: 1, ea: true},
		{name: CmpName, form: FormMemoryImmediate, clocks: 10, maxClocks: 10, perCount: 0, transfers: 1, ea: true},
		{name: TestName, form: FormRegisterImmediate, clocks: 5, maxClocks: 5, perCount: 0, transfers: 0, ea: false},
		{name: MovName, form: FormRegisterMemory, clocks: 8, maxClocks: 8, perCount: 0, transfers: 1, ea: true},
		{name: MovName, form: FormMemoryRegister, clocks: 9, maxClocks: 9, perCount: 0, transfers: 1, ea: true},
		{name: MovName, form: FormAccumulatorMemory, clocks: 10, maxClocks: 10, perCount: 0, transfers: 1, ea: false},
		{name: XchgName, form: FormAccumulatorRegister, clocks: 3, maxClocks: 3, perCount: 0, transfers: 0, ea: false},
		{name: LeaName, form: FormRegisterMemory, clocks: 2, maxClocks: 2, perCount: 0, transfers: 0, ea: true},
		{name: IncName, form: FormRegister16, clocks: 2, maxClocks: 2, perCount: 0, transfers: 0, ea: false},
		{name: DecName, form: FormRegister8, clocks: 3, maxClocks: 3, perCount: 0, transfers: 0, ea: false},
		{name: NegName, form: FormMemory16, clocks: 16, maxClocks: 16, perCount: 0, transfers: 2, ea: true},
		{name: ShlName, form: FormRegisterCL, clocks: 8, maxClocks: 8, perCount: 4, transfers: 0, ea: false},
		{name: SarName, form: FormMemoryCL, clocks: 20, maxClocks: 20, perCount: 4, transfers: 2, ea: true},
		{name: MulName, form: FormRegister16, clocks: 118, maxClocks: 133, perCount: 0, transfers: 0, ea: false},
		{name: MulName, form: FormMemory8, clocks: 76, maxClocks: 83, perCount: 0, transfers: 1, ea: true},
		{name: DivName, form: FormMemory16, clocks: 150, maxClocks: 168, perCount: 0, transfers: 1, ea: true},
		{name: IdivName, form: FormMemory16, clocks: 171, maxClocks: 190, perCount: 0, transfers: 1, ea: true},
		{name: PushName, form: FormRegister16, clocks: 11, maxClocks: 11, perCount: 0, transfers: 1, ea: false},
		{name: PopName, form: FormRegister16, clocks: 8, maxClocks: 8, perCount: 0, transfers: 1, ea: false},
		{name: CallName, form: FormNearLabel, clocks: 19, maxClocks: 19, perCount: 0, transfers: 1, ea: false},
		{name: RetName, form: FormNone, clocks: 8, maxClocks: 8, perCount: 0, transfers: 1, ea: false},
		{name: JmpName, form: FormShortLabel, clocks: 15, maxClocks: 15, perCount: 0, transfers: 0, ea: false},
		{name: JnzName, form: FormShortLabel, clocks: 4, maxClocks: 16, perCount: 0, transfers: 0, ea: false},
		{name: IntName, form: FormImmediate, clocks: 51, maxClocks: 51, perCount: 0, transfers: 5, ea: false},
		{name: CwdName, form: FormNone, clocks: 5, maxClocks: 5, perCount: 0, transfers: 0, ea: false},
		{name: MovswName, form: FormRepeatedString, clocks: 9, maxClocks: 9, perCount: 17, transfers: 2, ea: false},
		{name: StosbName, form: FormString, clocks: 11, maxClocks: 11, perCount: 0, transfers: 1, ea: false},
		{name: ScasbName, form: FormRepeatedString, clocks: 9, maxClocks: 9, perCount: 15, transfers: 1, ea: false},
	}
	for _, test := range tests {
		timing, ok := InstructionTiming(test.name, test.form)
		assert.True(t, ok, test.name)
		assert.Equal(t, test.clocks, timing.Clocks, test.name)
		assert.Equal(t, test.maxClocks, timing.MaxClocks, test.name)
		assert.Equal(t, test.perCount, timing.PerCount, test.name)
		assert.Equal(t, test.transfers, timing.Transfers, test.name)
		assert.Equal(t, test.ea, timing.EA, test.name)
	}

	_, ok := InstructionTiming(LeaName, FormRegisterRegister)
	assert.False(t, ok)
	_, ok = InstructionTiming("unknown", FormNone)
	assert.False(t, ok)
}

// TestInstructionTimingNames checks that each timing entry uses a known
// instruction name and a known operand form.
func TestInstructionTimingNames(t *testing.T) {
	for name, forms := range instructionTimings {
		_, known := Instructions[name]
		assert.True(t, known, name)
		for form, timing := range forms {
			assert.True(t, form > FormUnknown && form <= FormRepeatedString, name)
			assert.True(t, timing.Clocks > 0 && timing.MaxClocks >= timing.Clocks, name)
		}
	}
}
