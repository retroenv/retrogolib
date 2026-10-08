package sm83

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/set"
)

// opcodeTables lists the opcode tables with their prefix byte.
var opcodeTables = []struct {
	name   string
	prefix byte
	table  *[256]Opcode
}{
	{name: "base", prefix: 0x00, table: &Opcodes},
	{name: "CB", prefix: PrefixCB, table: &CBOpcodes},
}

// illegalOpcodes lists the 11 undefined base opcodes plus the CB prefix byte itself.
var illegalOpcodes = set.NewFromSlice([]byte{
	0xCB, 0xD3, 0xDB, 0xDD, 0xE3, 0xE4, 0xEB, 0xEC, 0xED, 0xF4, 0xFC, 0xFD,
})

// opcodeInfos returns every OpcodeInfo of an instruction.
func opcodeInfos(ins *Instruction) []OpcodeInfo {
	var infos []OpcodeInfo
	for _, info := range ins.Addressing {
		infos = append(infos, info)
	}
	for _, info := range ins.RegisterOpcodes {
		infos = append(infos, info)
	}
	for _, info := range ins.RegisterPairOpcodes {
		infos = append(infos, info)
	}
	for _, registers := range ins.BitOpcodes {
		for _, info := range registers {
			infos = append(infos, info)
		}
	}
	return infos
}

// TestVerifyOpcodes ensures bidirectional opcode mapping consistency for the
// base and CB tables without skipping any entry. Every table entry must have
// an OpcodeInfo with the same prefix, opcode byte, size and cycle count, and
// every OpcodeInfo must map forward to the same instruction pointer.
func TestVerifyOpcodes(t *testing.T) {
	t.Parallel()

	reachable := set.Set[*Instruction]{}

	for _, page := range opcodeTables {
		for b, op := range page.table {
			ins := op.Instruction
			if ins == nil {
				assert.True(t, page.prefix == 0 && illegalOpcodes.Contains(byte(b)),
					"%s opcode 0x%02X is undefined but not an illegal opcode", page.name, b)
				continue
			}
			assert.False(t, page.prefix == 0 && illegalOpcodes.Contains(byte(b)),
				"%s opcode 0x%02X is illegal but defined", page.name, b)
			reachable[ins] = struct{}{}

			found := false
			for _, info := range opcodeInfos(ins) {
				if info.Prefix != page.prefix || info.Opcode != byte(b) {
					continue
				}
				found = true
				assert.Equal(t, op.Size, info.Size, "%s opcode 0x%02X: %s size", page.name, b, ins.Name)
				assert.Equal(t, op.Timing, info.Cycles, "%s opcode 0x%02X: %s cycles", page.name, b, ins.Name)
			}
			assert.True(t, found, "%s opcode 0x%02X: %s has no reverse mapping", page.name, b, ins.Name)
		}
	}

	for ins := range reachable {
		assert.NotNil(t, Instructions[ins.Name], "instruction %s is not registered", ins.Name)
		for _, info := range opcodeInfos(ins) {
			var table *[256]Opcode
			for _, page := range opcodeTables {
				if page.prefix == info.Prefix {
					table = page.table
				}
			}
			assert.NotNil(t, table, "instruction %s has unknown prefix 0x%02X", ins.Name, info.Prefix)
			if table == nil {
				continue
			}
			assert.True(t, table[info.Opcode].Instruction == ins,
				"instruction %s opcode 0x%02X 0x%02X maps to a different instruction", ins.Name, info.Prefix, info.Opcode)
		}
	}

	for name, ins := range Instructions {
		assert.True(t, reachable.Contains(ins), "registered instruction %s is not in any opcode table", name)
	}
}

// TestBitOpcodesComplete verifies that BIT, RES and SET map every bit and register.
func TestBitOpcodesComplete(t *testing.T) {
	t.Parallel()

	registers := []RegisterParam{RegB, RegC, RegD, RegE, RegH, RegL, RegHLIndirect, RegA}
	tests := []struct {
		ins      *Instruction
		base     byte
		hlCycles byte
	}{
		{ins: CBBit, base: 0x40, hlCycles: 3},
		{ins: CBRes, base: 0x80, hlCycles: 4},
		{ins: CBSet, base: 0xC0, hlCycles: 4},
	}

	for _, tt := range tests {
		t.Run(tt.ins.Name, func(t *testing.T) {
			t.Parallel()

			for bit := range Bit(8) {
				for i, register := range registers {
					info, ok := tt.ins.GetOpcodeByBit(bit, register)
					assert.True(t, ok, "bit %d register %s", bit, register)
					assert.Equal(t, PrefixCB, info.Prefix)
					assert.Equal(t, tt.base+byte(bit)*8+byte(i), info.Opcode, "bit %d register %s", bit, register)
					wantCycles := byte(2)
					if register == RegHLIndirect {
						wantCycles = tt.hlCycles
					}
					assert.Equal(t, wantCycles, info.Cycles, "bit %d register %s", bit, register)
				}
			}
		})
	}

	_, ok := CBBit.GetOpcodeByBit(8, RegA)
	assert.False(t, ok)
	_, ok = NopInst.GetOpcodeByBit(0, RegA)
	assert.False(t, ok)
}

// TestOpcodeProperties validates timing and size constraints for all opcodes.
func TestOpcodeProperties(t *testing.T) {
	t.Parallel()

	for _, page := range opcodeTables {
		for i, op := range page.table {
			if op.Instruction == nil {
				continue
			}
			assert.True(t, op.Timing > 0 && op.Timing <= 6,
				"%s opcode 0x%02X (%s) has invalid timing: %d", page.name, i, op.Instruction.Name, op.Timing)
			assert.True(t, op.Size > 0 && op.Size <= MaxOpcodeSize,
				"%s opcode 0x%02X (%s) has invalid size: %d", page.name, i, op.Instruction.Name, op.Size)
		}
	}
}

// TestLdhCategories verifies that LDH counts as a memory read and write instruction.
func TestLdhCategories(t *testing.T) {
	t.Parallel()

	assert.True(t, MemoryReadInstructions.Contains(LdhInst.Name))
	assert.True(t, MemoryWriteInstructions.Contains(LdhInst.Name))
	assert.False(t, MemoryReadWriteInstructions.Contains(LdhInst.Name))

	for _, opcode := range []byte{0xE0, 0xE2, 0xF0, 0xF2} {
		assert.True(t, Opcodes[opcode].Instruction == LdhInst, "opcode 0x%02X", opcode)
	}
	for _, register := range []RegisterParam{RegHighMem, RegLoadHighMem, RegCIndirect, RegLoadCIndirect} {
		_, ok := LdhInst.GetOpcodeByRegister(register)
		assert.True(t, ok, "register %s", register)
	}
}
