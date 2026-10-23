package analysis

import (
	"testing"

	"github.com/retroenv/retrogolib/analysis/knownbits"
	"github.com/retroenv/retrogolib/arch/cpu/analysis"
	"github.com/retroenv/retrogolib/arch/cpu/z80"
	"github.com/retroenv/retrogolib/assert"
)

func TestCanonicalFormsAgainstExecution(t *testing.T) {
	t.Parallel()

	memory := z80.NewBasicMemory()
	for address := range 65536 {
		memory.Write(uint16(address), byte(address*37+19))
	}
	cpu, err := z80.New(memory)
	assert.NoError(t, err)
	for _, prefix := range []byte{0, z80.PrefixDD, z80.PrefixFD, z80.PrefixED} {
		for code := range 256 {
			for _, seed := range []byte{0, 0x47, 0xff} {
				checkCanonicalForm(t, cpu, memory, prefix, byte(code), seed)
			}
		}
	}
}

func TestPartialALUAgainstExecution(t *testing.T) {
	t.Parallel()

	memory := z80.NewBasicMemory()
	cpu, err := z80.New(memory)
	assert.NoError(t, err)
	for _, code := range []byte{0xe6, 0xee, 0xf6, 0xfe, 0xaf, 0xbf, 0xa7, 0xb7} {
		for _, zero := range []uint16{0, 0x80, 0xf0, 0xff} {
			input := fixture(0, code, 0x47)
			input.Before[AF], _ = knownbits.New[uint16](16, zero<<8, 0)
			result := bufferedResult()
			Transfer(input, &result)
			assert.True(t, result.EffectsComplete)
			after := applyWrites(t, input.Before, result)
			for number := range 256 {
				if uint16(number)&zero != 0 {
					continue
				}
				concrete := input
				concrete.Before[AF], _ = knownbits.NewConstant[uint16](16, uint16(number)<<8)
				install(t, cpu, memory, concrete)
				assert.NoError(t, cpu.Step())
				assertContains(t, 16, after[AF], cpu.AF())
			}
		}
	}
}

func checkCanonicalForm(t *testing.T, cpu *z80.CPU, memory *z80.BasicMemory, prefix, code, seed byte) {
	t.Helper()

	input := fixture(prefix, code, seed)
	input.Read = stableRead(memory)
	install(t, cpu, memory, input)
	result := bufferedResult()
	before := input.Before
	Transfer(input, &result)
	assert.Equal(t, before, input.Before)
	if !result.EffectsComplete {
		assert.True(t, result.Reason != analysis.Supported)
		return
	}
	assert.Equal(t, analysis.Supported, result.Reason)
	assert.True(t, result.Valid(StorageWidth))
	assert.NoError(t, cpu.Step())
	after := applyWrites(t, before, result)
	expected := concreteResources(cpu)
	for resource, value := range expected {
		if Resource(resource) != Q {
			constant, known := after[resource].Constant(StorageWidth(Resource(resource)))
			assert.True(t, known)
			assert.Equal(t, value, constant)
		}
	}
	for _, event := range result.Memory {
		if event.Kind == analysis.MemoryWrite {
			address, known := event.Address.Value.Constant(16)
			assert.True(t, known)
			assertContains(t, 8, event.Value, uint16(memory.Read(address)))
		}
	}
}

func fixture(prefix, code, seed byte) Input {
	info := z80.OpcodeInfo{
		Prefix: prefix,
		Opcode: code,
		Size:   1,
	}
	opcode := z80.Opcodes[code]
	switch prefix {
	case z80.PrefixED:
		opcode = z80.EDOpcodes[code]

	case z80.PrefixDD, z80.PrefixFD:
		table := &z80.DDOpcodes
		if prefix == z80.PrefixFD {
			table = &z80.FDOpcodes
		}
		if table[code].Instruction != nil {
			opcode = table[code]
		} else {
			opcode.Size++
		}
	}
	info.Size = opcode.Size
	input := Input{
		Opcode: info,
	}
	for resource := range ResourceCount {
		width := StorageWidth(resource)
		input.Before[resource], _ = knownbits.NewConstant(width, uint16(seed)*257+uint16(resource))
	}
	input.Before[Q], _ = knownbits.NewConstant[uint16](8, 0)
	width := knownbits.Width(8)
	operand := uint16(0x7f)
	operandBytes := int(info.Size) - 1
	if prefix != 0 {
		operandBytes--
	}
	if operandBytes == 2 {
		width = 16
		operand = 0x12ff
	}
	input.Operand, _ = knownbits.NewConstant(width, operand)
	return input
}

func install(t *testing.T, cpu *z80.CPU, memory *z80.BasicMemory, input Input) {
	t.Helper()

	cpu.PC = 0x100
	memory.Write(0x100, 0)
	assert.NoError(t, cpu.Step())
	cpu.A = byte(input.Before[AF].KnownOne >> 8)
	cpu.B = byte(input.Before[BC].KnownOne >> 8)
	cpu.C = byte(input.Before[BC].KnownOne)
	cpu.D = byte(input.Before[DE].KnownOne >> 8)
	cpu.E = byte(input.Before[DE].KnownOne)
	cpu.H = byte(input.Before[HL].KnownOne >> 8)
	cpu.L = byte(input.Before[HL].KnownOne)
	cpu.AltA = byte(input.Before[ShadowAF].KnownOne >> 8)
	cpu.AltB = byte(input.Before[ShadowBC].KnownOne >> 8)
	cpu.AltC = byte(input.Before[ShadowBC].KnownOne)
	cpu.AltD = byte(input.Before[ShadowDE].KnownOne >> 8)
	cpu.AltE = byte(input.Before[ShadowDE].KnownOne)
	cpu.AltH = byte(input.Before[ShadowHL].KnownOne >> 8)
	cpu.AltL = byte(input.Before[ShadowHL].KnownOne)
	cpu.IX = input.Before[IX].KnownOne
	cpu.IY = input.Before[IY].KnownOne
	cpu.SP = input.Before[SP].KnownOne
	cpu.I = byte(input.Before[I].KnownOne)
	cpu.R = byte(input.Before[R].KnownOne)
	cpu.Flags = flags(byte(input.Before[AF].KnownOne))
	cpu.AltFlags = flags(byte(input.Before[ShadowAF].KnownOne))
	if input.Before[IFF2].KnownOne != 0 {
		cpu.EnableInterrupts()
	} else {
		cpu.DisableInterrupts()
	}
	cpu.PC = 0x200
	pc := cpu.PC
	if input.Opcode.Prefix != 0 {
		memory.Write(pc, input.Opcode.Prefix)
		pc++
	}
	memory.Write(pc, input.Opcode.Opcode)
	memory.Write(pc+1, byte(input.Operand.KnownOne))
	memory.Write(pc+2, byte(input.Operand.KnownOne>>8))
}

func flags(value byte) z80.Flags {
	return z80.Flags{
		C: value & 1,
		N: value >> 1 & 1,
		P: value >> 2 & 1,
		X: value >> 3 & 1,
		H: value >> 4 & 1,
		Y: value >> 5 & 1,
		Z: value >> 6 & 1,
		S: value >> 7 & 1,
	}
}

func concreteResources(cpu *z80.CPU) [ResourceCount]uint16 {
	iff := uint16(0)
	if cpu.State().Interrupts.IFF2 {
		iff = 1
	}
	alt := cpu.AltFlags
	altF := alt.C | alt.N<<1 | alt.P<<2 | alt.X<<3 | alt.H<<4 | alt.Y<<5 | alt.Z<<6 | alt.S<<7
	return [ResourceCount]uint16{
		AF:       cpu.AF(),
		BC:       cpu.BC(),
		DE:       cpu.DE(),
		HL:       cpu.HL(),
		SP:       cpu.SP,
		IX:       cpu.IX,
		IY:       cpu.IY,
		ShadowAF: uint16(cpu.AltA)<<8 | uint16(altF),
		ShadowBC: uint16(cpu.AltB)<<8 | uint16(cpu.AltC),
		ShadowDE: uint16(cpu.AltD)<<8 | uint16(cpu.AltE),
		ShadowHL: uint16(cpu.AltH)<<8 | uint16(cpu.AltL),
		I:        uint16(cpu.I),
		R:        uint16(cpu.R),
		IFF2:     iff,
	}
}

func stableRead(memory *z80.BasicMemory) Read {
	return func(address Address) (knownbits.Value[uint16], bool) {
		value, _ := address.Value.Constant(16)
		return knownbits.NewConstant[uint16](8, uint16(memory.Read(value)))
	}
}

func applyWrites(t *testing.T, before [ResourceCount]knownbits.Value[uint16],
	result analysis.Result[Resource, uint16, Address]) [ResourceCount]knownbits.Value[uint16] {

	t.Helper()

	for _, write := range result.Writes {
		resource := write.View.Storage
		value, valid := analysis.Apply(StorageWidth(resource), before[resource], write)
		assert.True(t, valid)
		before[resource] = value
	}
	return before
}

func assertContains(t *testing.T, width knownbits.Width, fact knownbits.Value[uint16], value uint16) {
	t.Helper()

	assert.True(t, fact.Valid(width))
	assert.Equal(t, uint16(0), value&fact.KnownZero)
	assert.Equal(t, fact.KnownOne, value&fact.KnownOne)
}

func bufferedResult() analysis.Result[Resource, uint16, Address] {
	return analysis.Result[Resource, uint16, Address]{
		Reads:  make([]analysis.View[Resource], 0, 16),
		Writes: make([]analysis.Write[Resource, uint16], 0, 16),
		Copies: make([]analysis.Copy[Resource], 0, 4),
		Memory: make([]analysis.MemoryEvent[uint16, Address], 0, 8),
	}
}
