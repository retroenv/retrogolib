package z80

import (
	"fmt"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestIndexHalfOpcodeMetadata(t *testing.T) {
	for _, group := range []struct {
		prefix byte
		high   RegisterParam
		low    RegisterParam
		table  [256]Opcode
	}{
		{
			prefix: PrefixDD,
			high:   RegIXH,
			low:    RegIXL,
			table:  DDOpcodes,
		},
		{
			prefix: PrefixFD,
			high:   RegIYH,
			low:    RegIYL,
			table:  FDOpcodes,
		},
	} {
		t.Run(fmt.Sprintf("%02x", group.prefix), func(t *testing.T) {
			registers := []RegisterParam{RegB, RegC, RegD, RegE, group.high, group.low, RegNone, RegA}
			for _, lane := range []byte{4, 5} {
				for _, base := range []byte{0x04, 0x05, 0x06} {
					opcode := base + lane*8
					testIndexHalfMetadata(t, group.prefix, opcode, group.table[opcode], registers[lane], RegNone)
				}
				for operation := range byte(8) {
					opcode := 0x80 + operation*8 + lane
					testIndexHalfMetadata(t, group.prefix, opcode, group.table[opcode], registers[lane], RegNone)
				}
			}
			for dst, destination := range registers {
				for src, source := range registers {
					if destination == RegNone || source == RegNone || (dst != 4 && dst != 5 && src != 4 && src != 5) {
						continue
					}
					opcode := byte(0x40 + dst*8 + src)
					testIndexHalfMetadata(t, group.prefix, opcode, group.table[opcode], destination, source)
				}
			}
		})
	}
}

func TestIndexHalfRegisterNames(t *testing.T) {
	for index, register := range []RegisterParam{RegIXH, RegIXL, RegIYH, RegIYL} {
		assert.Equal(t, []string{"ixh", "ixl", "iyh", "iyl"}[index], register.String())
	}
}

func testIndexHalfMetadata(t *testing.T, prefix, opcode byte, forward Opcode, dst, src RegisterParam) {
	t.Helper()
	assert.NotNil(t, forward.Instruction)
	assert.True(t, forward.Instruction.Unofficial)
	addressing := RegisterAddressing
	if forward.Size == 3 {
		addressing = ImmediateAddressing
	}
	assert.True(t, forward.Instruction.HasAddressing(addressing))
	info, ok := forward.Instruction.RegisterOpcodes[dst]
	if src != RegNone {
		info, ok = forward.Instruction.RegisterPairOpcodes[[2]RegisterParam{dst, src}]
	}
	assert.True(t, ok, "prefix %02x opcode %02x", prefix, opcode)
	assert.Equal(t, OpcodeInfo{
		Prefix: prefix,
		Opcode: opcode,
		Size:   forward.Size,
		Cycles: forward.Timing,
	}, info)
}
