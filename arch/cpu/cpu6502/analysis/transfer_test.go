package analysis

import (
	"testing"

	"github.com/retroenv/retrogolib/analysis/knownbits"
	"github.com/retroenv/retrogolib/arch/cpu/analysis"
	"github.com/retroenv/retrogolib/arch/cpu/cpu6502"
	"github.com/retroenv/retrogolib/assert"
)

type testMemory struct {
	data [65536]byte
}

func (mm *testMemory) Read(address uint16) byte {
	return mm.data[address]
}

func (mm *testMemory) Write(address uint16, value byte) {
	mm.data[address] = value
}

func TestCanonicalFormsAgainstExecution(t *testing.T) {
	t.Parallel()

	for variant := cpu6502.VariantNMOS6502; variant <= cpu6502.VariantSynertek65C02; variant++ {
		memory := &testMemory{}
		for address := range memory.data {
			memory.data[address] = byte(address*37 + 19)
		}
		bus, err := cpu6502.NewMemory(memory)
		assert.NoError(t, err)
		cpu, err := cpu6502.New(bus, cpu6502.WithVariant(variant))
		assert.NoError(t, err)
		for code := range 256 {
			for _, seed := range []byte{0, 0x47, 0x81, 0xff} {
				checkForm(t, cpu, memory, variant, byte(code), seed)
			}
		}
	}
}

func TestPartialArithmeticAgainstExecution(t *testing.T) {
	t.Parallel()

	memory := &testMemory{}
	bus, err := cpu6502.NewMemory(memory)
	assert.NoError(t, err)
	cpu, err := cpu6502.New(bus, cpu6502.WithVariant(cpu6502.VariantNES6502))
	assert.NoError(t, err)
	for _, code := range []byte{0x69, 0xe9, 0x29, 0x09, 0x49, 0xc9, 0x0a, 0x4a, 0x2a, 0x6a} {
		for _, fact := range []knownbits.Value[uint8]{
			{},
			{
				KnownZero: 0x80,
			},
			{
				KnownOne: 0x0f,
			},
			{
				KnownZero: 0x55,
				KnownOne:  0xaa,
			},
		} {
			checkPartialForm(t, cpu, memory, code, fact)
		}
	}
}

func checkForm(t *testing.T, cpu *cpu6502.CPU, memory *testMemory, variant cpu6502.CPUVariant, code, seed byte) {
	t.Helper()

	table := tableForVariant(variant)
	opcode := table[code]
	if opcode.Instruction == nil {
		return
	}
	input := fixture(code, encodedSize(opcode.Addressing), seed)
	input.Variant = variant
	input.Read = stableRead(memory)
	result := bufferedResult()
	before := input.Before
	install(cpu, memory, input)
	Transfer(input, &result)
	assert.Equal(t, before, input.Before)
	if !result.EffectsComplete {
		assert.Equal(t, analysis.UnsupportedOpcode, result.Reason)
		return
	}
	assert.Equal(t, analysis.Supported, result.Reason)
	assert.True(t, result.Valid(resourceWidth))
	if result.Flow == analysis.FlowStop {
		return
	}
	install(cpu, memory, input)
	assert.NoError(t, cpu.Step())
	after := applyWrites(t, input.Before, result)
	expected := [...]byte{cpu.A, cpu.X, cpu.Y, cpu.SP, cpu.GetFlags()}
	for resource, value := range expected {
		constant, known := after[resource].Constant(8)
		assert.True(t, known)
		assert.Equal(t, value, constant)
	}
	for _, event := range result.Memory {
		if event.Kind == analysis.MemoryWrite {
			address, known := event.Address.Value.Constant(16)
			assert.True(t, known)
			assertContains(t, event.Value, memory.Read(address))
		}
	}
}

func checkPartialForm(t *testing.T, cpu *cpu6502.CPU, memory *testMemory, code byte, fact knownbits.Value[uint8]) {
	t.Helper()

	opcode := cpu6502.Opcodes[code]
	input := fixture(code, encodedSize(opcode.Addressing), 0x47)
	input.Before[A] = fact
	input.Before[P], _ = knownbits.New[uint8](8, 0xda, 0x24)
	result := bufferedResult()
	Transfer(input, &result)
	assert.True(t, result.EffectsComplete)
	assert.Equal(t, analysis.Supported, result.Reason)
	after := applyWrites(t, input.Before, result)
	for number := range 256 {
		if byte(number)&fact.KnownZero != 0 || byte(number)&fact.KnownOne != fact.KnownOne {
			continue
		}
		for carry := range byte(2) {
			concrete := input
			concrete.Before[A], _ = knownbits.NewConstant[uint8](8, byte(number))
			concrete.Before[P], _ = knownbits.NewConstant[uint8](8, 0x24|carry)
			install(cpu, memory, concrete)
			assert.NoError(t, cpu.Step())
			assertContains(t, after[A], cpu.A)
			assertContains(t, after[P], cpu.GetFlags())
		}
	}
}

func fixture(code, size, seed byte) Input {
	input := Input{
		Opcode: cpu6502.OpcodeInfo{
			Opcode: code,
			Size:   size,
		},
		Variant: cpu6502.VariantNES6502,
	}
	input.PC, _ = knownbits.NewConstant[uint16](16, 0x200)
	for resource, value := range [...]byte{seed, seed ^ 0x43, seed ^ 0x18, 0x80, 0x20 | seed&0xc7} {
		input.Before[resource], _ = knownbits.NewConstant[uint8](8, value)
	}
	width := knownbits.Width(8)
	operand := uint16(0x7f)
	if size == 3 {
		width = 16
		operand = 0x12ff
	}
	input.Operand, _ = knownbits.NewConstant(width, operand)
	return input
}

func install(cpu *cpu6502.CPU, memory *testMemory, input Input) {
	pc, _ := input.PC.Constant(16)
	cpu.PC = pc
	cpu.A = input.Before[A].KnownOne
	cpu.X = input.Before[X].KnownOne
	cpu.Y = input.Before[Y].KnownOne
	cpu.SP = input.Before[SP].KnownOne
	status := input.Before[P].KnownOne
	cpu.Flags = cpu6502.Flags{
		C: status & 1,
		Z: status >> 1 & 1,
		I: status >> 2 & 1,
		D: status >> 3 & 1,
		B: status >> 4 & 1,
		U: status >> 5 & 1,
		V: status >> 6 & 1,
		N: status >> 7 & 1,
	}
	memory.Write(pc, input.Opcode.Opcode)
	memory.Write(pc+1, byte(input.Operand.KnownOne))
	memory.Write(pc+2, byte(input.Operand.KnownOne>>8))
}

func stableRead(memory *testMemory) Read {
	return func(address Address) (knownbits.Value[uint8], bool) {
		value, _ := address.Value.Constant(16)
		return knownbits.NewConstant[uint8](8, memory.Read(value))
	}
}

func applyWrites(t *testing.T, before [ResourceCount]knownbits.Value[uint8],
	result analysis.Result[Resource, uint8, Address]) [ResourceCount]knownbits.Value[uint8] {

	t.Helper()

	for _, write := range result.Writes {
		value, valid := analysis.Apply(8, before[write.View.Storage], write)
		assert.True(t, valid)
		before[write.View.Storage] = value
	}
	return before
}

func assertContains(t *testing.T, fact knownbits.Value[uint8], value byte) {
	t.Helper()

	assert.True(t, fact.Valid(8))
	assert.Equal(t, byte(0), value&fact.KnownZero)
	assert.Equal(t, fact.KnownOne, value&fact.KnownOne)
}

func tableForVariant(variant cpu6502.CPUVariant) *[256]cpu6502.Opcode {
	switch variant {
	case cpu6502.Variant65C02:
		return &cpu6502.Opcodes65C02

	case cpu6502.VariantSynertek65C02:
		return &cpu6502.OpcodesSynertek65C02

	default:
		return &cpu6502.Opcodes
	}
}

func resourceWidth(resource Resource) knownbits.Width {
	if resource >= ResourceCount {
		return 0
	}
	return 8
}

func bufferedResult() analysis.Result[Resource, uint8, Address] {
	return analysis.Result[Resource, uint8, Address]{
		Reads:  make([]analysis.View[Resource], 0, 16),
		Writes: make([]analysis.Write[Resource, uint8], 0, 12),
		Copies: make([]analysis.Copy[Resource], 0, 4),
		Memory: make([]analysis.MemoryEvent[uint8, Address], 0, 8),
	}
}
