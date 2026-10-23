package analysis

import (
	"github.com/retroenv/retrogolib/analysis/knownbits"
	"github.com/retroenv/retrogolib/arch/cpu/analysis"
	"github.com/retroenv/retrogolib/arch/cpu/cpu6502"
)

// Resource identifies byte storage. A has ID zero.
type Resource uint8

const (
	A Resource = iota
	X
	Y
	SP
	P
	ResourceCount
)

// Address is a logical CPU address, before caller-owned physical mapping.
type Address struct {
	Value knownbits.Value[uint16]
}

// Read returns a scalar and a proof that this read is stable in the caller's context.
// It must not mutate the input state or result buffers.
type Read func(Address) (knownbits.Value[uint8], bool)

// Input contains one encoded instruction and an immutable input snapshot.
// Opcode.Size must match the selected variant's canonical addressing form.
// Operand contains the encoded operand field, with width 8 or 16 as required.
// PC is the abstract instruction address. Its zero value is unknown.
type Input struct {
	Opcode  cpu6502.OpcodeInfo
	Variant cpu6502.CPUVariant
	Operand knownbits.Value[uint16]
	PC      knownbits.Value[uint16]

	Before [ResourceCount]knownbits.Value[uint8]
	Read   Read
}

// Transfer fills result and reuses its caller-owned buffers. Result must be nonnil.
// All reads use Before. Copies refer to that snapshot, including self copies.
func Transfer(input Input, result *analysis.Result[Resource, uint8, Address]) {
	result.Reset()
	opcode, reason := decode(input)
	if reason != analysis.Supported {
		result.Reason = reason
		return
	}
	tt := transfer{
		input:  input,
		opcode: opcode,
		result: result,
		after:  input.Before,
	}
	result.Reason = analysis.Supported
	result.Flow = analysis.FlowNext
	result.EffectsComplete = tt.execute()
	tt.finish()
}

type transfer struct {
	input          Input
	opcode         cpu6502.Opcode
	result         *analysis.Result[Resource, uint8, Address]
	after          [ResourceCount]knownbits.Value[uint8]
	changed        [ResourceCount]uint8
	invalid        bool
	operandAddress Address
	addressReady   bool
}

func (tt *transfer) read(resource Resource) knownbits.Value[uint8] {
	tt.result.Reads = append(tt.result.Reads, byteView(resource))
	return tt.input.Before[resource]
}

func (tt *transfer) flag(position knownbits.Width) knownbits.Bit {
	view := analysis.View[Resource]{
		Storage: P,
		Offset:  position,
		Bits:    1,
	}
	tt.result.Reads = append(tt.result.Reads, view)
	value, _ := knownbits.TestBit(8, tt.input.Before[P], position)
	return value
}

func (tt *transfer) write(resource Resource, value knownbits.Value[uint8]) {
	tt.after[resource] = value
	tt.changed[resource] = 0xff
}

func (tt *transfer) writeFlag(position knownbits.Width, value knownbits.Bit) {
	part := knownbits.Value[uint8]{}
	switch value {
	case knownbits.BitClear:
		part.KnownZero = 1

	case knownbits.BitSet:
		part.KnownOne = 1
	}
	tt.after[P], _ = knownbits.Insert(8, tt.after[P], position, 1, part)
	tt.changed[P] |= 1 << position
}

func (tt *transfer) zn(value knownbits.Value[uint8]) {
	zero, _ := knownbits.Zero(8, value)
	negative, _ := knownbits.Sign(8, value)
	tt.writeFlag(1, zero)
	tt.writeFlag(7, negative)
}

func (tt *transfer) copy(destination, source Resource, flags bool) {
	value := tt.read(source)
	tt.write(destination, value)
	tt.result.Copies = append(tt.result.Copies, analysis.Copy[Resource]{
		Source:      byteView(source),
		Destination: byteView(destination),
	})
	if flags {
		tt.zn(value)
	}
}

func (tt *transfer) finish() {
	if !tt.result.EffectsComplete {
		tt.result.Reason = analysis.UnsupportedOpcode
		tt.result.Flow = analysis.FlowUnknown
	}
	if tt.invalid {
		tt.result.Reason = analysis.InvalidInput
		tt.result.EffectsComplete = false
		tt.result.Flow = analysis.FlowUnknown
		tt.result.Copies = tt.result.Copies[:0]
		return
	}
	for resource, mask := range tt.changed {
		if mask == 0 {
			continue
		}
		if Resource(resource) != P || mask == 0xff {
			tt.result.Writes = append(tt.result.Writes, analysis.Write[Resource, uint8]{
				View:  byteView(Resource(resource)),
				Value: tt.after[resource],
				Rule:  analysis.Replace,
			})
			continue
		}
		tt.finishFlags(mask)
	}
}

func (tt *transfer) finishFlags(mask uint8) {
	for position := range knownbits.Width(8) {
		if mask&(1<<position) == 0 {
			continue
		}
		part, _ := knownbits.Extract(8, tt.after[P], position, 1)
		tt.result.Writes = append(tt.result.Writes, analysis.Write[Resource, uint8]{
			View: analysis.View[Resource]{
				Storage: P,
				Offset:  position,
				Bits:    1,
			},
			Value: part,
			Rule:  analysis.Preserve,
		})
	}
}

func byteView(resource Resource) analysis.View[Resource] {
	return analysis.View[Resource]{
		Storage: resource,
		Bits:    8,
	}
}

func decode(input Input) (cpu6502.Opcode, analysis.Reason) {
	if input.Variant < cpu6502.VariantNMOS6502 || input.Variant > cpu6502.VariantSynertek65C02 {
		return cpu6502.Opcode{}, analysis.UnsupportedVariant
	}
	opcode := cpu6502.Opcodes[input.Opcode.Opcode]
	switch input.Variant {
	case cpu6502.Variant65C02:
		opcode = cpu6502.Opcodes65C02[input.Opcode.Opcode]

	case cpu6502.VariantSynertek65C02:
		opcode = cpu6502.OpcodesSynertek65C02[input.Opcode.Opcode]
	}
	if opcode.Instruction == nil {
		return opcode, analysis.UnsupportedOpcode
	}
	size := encodedSize(opcode.Addressing)
	if size != input.Opcode.Size {
		return opcode, analysis.IllegalOperand
	}
	width := knownbits.Width(8)
	if size == 3 {
		width = 16
	}
	if !input.Operand.Valid(width) || !input.PC.Valid(16) {
		return opcode, analysis.InvalidInput
	}
	for _, value := range input.Before {
		if !value.Valid(8) {
			return opcode, analysis.InvalidInput
		}
	}
	return opcode, analysis.Supported
}

func encodedSize(mode cpu6502.AddressingMode) byte {
	switch mode {
	case cpu6502.ImpliedAddressing, cpu6502.AccumulatorAddressing:
		return 1

	case cpu6502.AbsoluteAddressing, cpu6502.AbsoluteXAddressing, cpu6502.AbsoluteYAddressing,
		cpu6502.IndirectAddressing, cpu6502.AbsoluteXIndirectAddressing, cpu6502.ZeroPageRelativeAddressing:

		return 3

	default:
		return 2
	}
}

func invert(value knownbits.Bit) knownbits.Bit {
	switch value {
	case knownbits.BitClear:
		return knownbits.BitSet

	case knownbits.BitSet:
		return knownbits.BitClear

	default:
		return knownbits.BitUnknown
	}
}
