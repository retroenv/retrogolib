package z80

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

const (
	blockTransferCodeAddress        = 0x2800
	blockTransferSourceAddress      = 0x4000
	blockTransferDestinationAddress = 0x6000
	blockTransferLdirOpcode         = 0xb0
	blockTransferLddrOpcode         = 0xb8
	blockTransferFinalCycles        = 16
	blockTransferRepeatCycles       = 21
	blockTransferRepeatExtraCycles  = 5
	blockTransferDocumentedFlags    = 0xd7
	blockTransferParityFlag         = 0x04
	blockTransferFill               = 0x69
)

type blockTransferFixture struct {
	cpu    *CPU
	memory *BasicMemory

	initial  State
	expected [0x10000]uint8

	source      uint16
	destination uint16
	steps       uint16
	decrement   bool
}

func TestRepeatingBlockTransferTiming(t *testing.T) {
	t.Parallel()

	// Zilog UM008011-0816, pages 132 and 136, specifies 21 T-states for a
	// repeating step and 16 T-states for the final step.
	for _, opcode := range []byte{blockTransferLdirOpcode, blockTransferLddrOpcode} {
		for _, count := range []uint16{1, 2, 17, 255} {
			for _, flags := range []byte{0x12, 0x13, 0x52, 0x53, 0x92, 0x93, 0xd2, 0xd3,
				0x16, 0x17, 0x56, 0x57, 0x96, 0x97, 0xd6, 0xd7} {
				name := LdirName
				if opcode == blockTransferLddrOpcode {
					name = LddrName
				}
				t.Run(fmt.Sprintf("%s/count=%d/flags=%02x", name, count, flags), func(t *testing.T) {
					t.Parallel()
					fixture := newBlockTransferFixture(t, opcode, count, flags)
					for step := uint16(1); step <= count; step++ {
						fixture.step(t, count-step)
					}

					wantCycles := uint64(blockTransferRepeatCycles)*uint64(count) - blockTransferRepeatExtraCycles
					assert.Equal(t, wantCycles, fixture.cpu.Cycles())
					assert.True(t, bytes.Equal(fixture.expected[:], fixture.memory.data[:]))
				})
			}
		}
	}
}

func TestRepeatingBlockTransferZeroCountWrapsOnFirstStep(t *testing.T) {
	t.Parallel()

	for _, opcode := range []byte{blockTransferLdirOpcode, blockTransferLddrOpcode} {
		t.Run(fmt.Sprintf("%02x", opcode), func(t *testing.T) {
			t.Parallel()
			fixture := newBlockTransferFixture(t, opcode, 0, 0xd3)
			fixture.step(t, 0xffff)
			assert.Equal(t, uint64(blockTransferRepeatCycles), fixture.cpu.Cycles())
			assert.True(t, bytes.Equal(fixture.expected[:], fixture.memory.data[:]))
		})
	}
}

func TestRepeatingBlockTransferTimingMetadata(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		opcode      byte
		instruction *Instruction
	}{
		{
			opcode:      blockTransferLdirOpcode,
			instruction: EdLdir,
		},
		{
			opcode:      blockTransferLddrOpcode,
			instruction: EdLddr,
		},
	} {
		t.Run(test.instruction.Name, func(t *testing.T) {
			t.Parallel()
			fixture := newBlockTransferFixture(t, test.opcode, 1, 0)
			decoded, prefix, err := fixture.cpu.decodeEDInstruction()
			assert.NoError(t, err)
			assert.Equal(t, byte(PrefixED), prefix)
			assert.True(t, decoded.Instruction == test.instruction)
			assert.Equal(t, ImpliedAddressing, decoded.Addressing)
			assert.Equal(t, byte(2), decoded.Size)
			assert.Equal(t, byte(blockTransferFinalCycles), decoded.Timing)

			reverse, known := test.instruction.Addressing[ImpliedAddressing]
			assert.True(t, known)
			assert.Equal(t, byte(PrefixED), reverse.Prefix)
			assert.Equal(t, test.opcode, reverse.Opcode)
			assert.Equal(t, byte(2), reverse.Size)
			assert.Equal(t, byte(blockTransferFinalCycles), reverse.Cycles)
		})
	}
}

func (fixture *blockTransferFixture) step(t *testing.T, remaining uint16) {
	t.Helper()

	value := fixture.expected[fixture.source]
	fixture.expected[fixture.destination] = value
	before := fixture.cpu.Cycles()
	assert.NoError(t, fixture.cpu.Step())
	wantCycles := uint64(blockTransferFinalCycles)
	wantPC := uint16(blockTransferCodeAddress + 2)
	wantFlags := fixture.initial.Flags
	if remaining != 0 {
		wantCycles = blockTransferRepeatCycles
		wantPC = blockTransferCodeAddress
	}

	assert.Equal(t, wantCycles, fixture.cpu.Cycles()-before)
	assert.Equal(t, remaining, fixture.cpu.BC())
	assert.Equal(t, wantPC, fixture.cpu.PC)
	assert.Equal(t, value, fixture.memory.Read(fixture.destination))
	if fixture.decrement {
		fixture.source--
		fixture.destination--
	} else {
		fixture.source++
		fixture.destination++
	}

	fixture.steps++
	assert.Equal(t, fixture.source, fixture.cpu.HL())
	assert.Equal(t, fixture.destination, fixture.cpu.DE())
	retained := wantFlags.S<<7 | wantFlags.Z<<6 | wantFlags.C
	if remaining != 0 {
		retained |= blockTransferParityFlag
	}
	actualFlags := fixture.cpu.GetFlags() & blockTransferDocumentedFlags
	assert.Equal(t, retained, actualFlags)
	refresh := fixture.initial.R&0x80 | (fixture.initial.R+byte(fixture.steps*2))&0x7f
	assert.Equal(t, refresh, fixture.cpu.R)
	fixture.checkUnrelatedState(t)
}

func (fixture *blockTransferFixture) checkUnrelatedState(t *testing.T) {
	t.Helper()

	state := fixture.cpu.State()
	state.B = fixture.initial.B
	state.C = fixture.initial.C
	state.D = fixture.initial.D
	state.E = fixture.initial.E
	state.H = fixture.initial.H
	state.L = fixture.initial.L
	state.PC = fixture.initial.PC
	state.R = fixture.initial.R
	state.Cycles = fixture.initial.Cycles
	// The step checks documented flags. Internal WZ and undocumented flags
	// are outside this timing test.
	state.Flags = fixture.initial.Flags
	state.MEMPTR = fixture.initial.MEMPTR
	assert.Equal(t, fixture.initial, state)
}

func newBlockTransferFixture(t *testing.T, opcode byte, count uint16, flags byte) *blockTransferFixture {
	t.Helper()

	memory := NewBasicMemory()
	for index := range memory.data {
		memory.data[index] = blockTransferFill
	}
	memory.Write(blockTransferCodeAddress, PrefixED)
	memory.Write(blockTransferCodeAddress+1, opcode)
	for offset := range max(count, 1) {
		address := uint16(blockTransferSourceAddress) + offset
		memory.Write(address, byte(offset*37+0x42))
	}

	cpu, err := New(memory, WithInitialPC(blockTransferCodeAddress), WithInitialSP(0xd000))
	assert.NoError(t, err)
	initializeBlockTransferRegisters(cpu, flags)
	fixture := &blockTransferFixture{
		cpu:    cpu,
		memory: memory,

		expected: memory.data,

		source:      blockTransferSourceAddress,
		destination: blockTransferDestinationAddress,
		decrement:   opcode == blockTransferLddrOpcode,
	}
	if fixture.decrement {
		fixture.source += max(count, 1) - 1
		fixture.destination += max(count, 1) - 1
	}

	cpu.setBC(count)
	cpu.setHL(fixture.source)
	cpu.setDE(fixture.destination)
	fixture.initial = cpu.State()
	return fixture
}

func initializeBlockTransferRegisters(cpu *CPU, flags byte) {
	cpu.A = 0xa5
	cpu.AltA = 0x13
	cpu.AltB = 0x24
	cpu.AltC = 0x35
	cpu.AltD = 0x46
	cpu.AltE = 0x57
	cpu.AltH = 0x68
	cpu.AltL = 0x79
	cpu.IX = 0xabc1
	cpu.IY = 0xdef2
	cpu.I = 0x87
	cpu.R = 0xf8
	cpu.AltFlags = Flags{
		S: 1,
		Z: 1,
		C: 1,
	}
	cpu.setFlags(flags)
}
