package analysis

import (
	"github.com/retroenv/retrogolib/analysis/knownbits"
	"github.com/retroenv/retrogolib/arch/cpu/analysis"
	"github.com/retroenv/retrogolib/arch/cpu/z80"
)

// Resource identifies canonical storage. AF has ID zero.
type Resource uint8

const (
	AF Resource = iota
	BC
	DE
	HL
	SP
	IX
	IY
	ShadowAF
	ShadowBC
	ShadowDE
	ShadowHL
	I
	R
	Q
	IFF2
	ResourceCount
)

// Address is a logical CPU address before caller-owned physical mapping.
type Address struct {
	Value knownbits.Value[uint16]
}

// Read returns a byte scalar and its stability proof in the caller's context.
// It must not mutate the input snapshot or result buffers.
type Read func(Address) (knownbits.Value[uint16], bool)

// Input holds a canonical encoded form and a snapshot before opcode fetch.
// Operand contains zero, one, or two encoded bytes in little-endian order.
// For LD (IX+d),n and LD (IY+d),n, its low byte is d and its high byte is n.
type Input struct {
	Opcode  z80.OpcodeInfo
	Operand knownbits.Value[uint16]
	Before  [ResourceCount]knownbits.Value[uint16]
	Read    Read
}

// StorageWidth returns the width of a declared resource, or zero for an invalid ID.
func StorageWidth(resource Resource) knownbits.Width {
	switch {
	case resource < I:
		return 16

	case resource == I || resource == R || resource == Q:
		return 8

	case resource == IFF2:
		return 1

	default:
		return 0
	}
}

// RegisterView binds canonical CPU register names to their logical storage.
func RegisterView(register z80.RegisterParam) (analysis.View[Resource], bool) {
	switch register {
	case z80.RegA:
		return view(AF, 8, 8), true

	case z80.RegB, z80.RegC:
		return view(BC, byteOffset(register == z80.RegB), 8), true

	case z80.RegD, z80.RegE:
		return view(DE, byteOffset(register == z80.RegD), 8), true

	case z80.RegH, z80.RegL:
		return view(HL, byteOffset(register == z80.RegH), 8), true

	case z80.RegBC:
		return view(BC, 0, 16), true

	case z80.RegDE:
		return view(DE, 0, 16), true

	case z80.RegHL:
		return view(HL, 0, 16), true

	case z80.RegAF:
		return view(AF, 0, 16), true

	case z80.RegSP:
		return view(SP, 0, 16), true

	case z80.RegIX:
		return view(IX, 0, 16), true

	case z80.RegIY:
		return view(IY, 0, 16), true

	case z80.RegI:
		return view(I, 0, 8), true

	case z80.RegR:
		return view(R, 0, 8), true

	default:
		return analysis.View[Resource]{}, false
	}
}

// Transfer fills result from an immutable input snapshot and reuses its buffers.
func Transfer(input Input, result *analysis.Result[Resource, uint16, Address]) {
	result.Reset()
	opcode, indexed, reason := decode(input)
	if reason != analysis.Supported {
		result.Reason = reason
		return
	}
	tt := transfer{
		input:   input,
		opcode:  opcode,
		indexed: indexed,
		result:  result,
		after:   input.Before,
	}
	result.Reason = analysis.Supported
	result.Flow = analysis.FlowNext
	tt.refresh()
	result.EffectsComplete = tt.execute()
	tt.finish()
}

type transfer struct {
	input   Input
	opcode  z80.Opcode
	indexed bool
	result  *analysis.Result[Resource, uint16, Address]
	after   [ResourceCount]knownbits.Value[uint16]
	changed [ResourceCount]uint16
	invalid bool
}

func (tt *transfer) read(part analysis.View[Resource]) knownbits.Value[uint16] {
	tt.result.Reads = append(tt.result.Reads, part)
	before := tt.input.Before[part.Storage]
	if part.Storage == R {
		before = tt.after[R]
	}
	value, _ := knownbits.Extract(StorageWidth(part.Storage), before, part.Offset, part.Bits)
	return value
}

func (tt *transfer) write(part analysis.View[Resource], value knownbits.Value[uint16]) {
	storage := part.Storage
	tt.after[storage], _ = knownbits.Insert(StorageWidth(storage), tt.after[storage], part.Offset, part.Bits, value)
	mask := (^uint16(0) >> (16 - part.Bits)) << part.Offset
	tt.changed[storage] |= mask
}

func (tt *transfer) copy(destination, source analysis.View[Resource]) {
	value := tt.read(source)
	tt.write(destination, value)
	if source.Storage != R {
		tt.result.Copies = append(tt.result.Copies, analysis.Copy[Resource]{
			Source:      source,
			Destination: destination,
		})
	}
}

func (tt *transfer) refresh() {
	part := view(R, 0, 8)
	tt.result.Reads = append(tt.result.Reads, part)
	low, _ := knownbits.Truncate(8, 7, tt.input.Before[R])
	count := uint16(1)
	if tt.input.Opcode.Prefix != 0 {
		count = 2
	}
	constant, _ := knownbits.NewConstant[uint16](7, count)
	incremented, _ := knownbits.Add(7, low, constant, knownbits.BitClear)
	value, _ := knownbits.Insert(8, tt.input.Before[R], 0, 7, incremented.Value)
	tt.write(part, value)
}

func (tt *transfer) finish() {
	if !tt.result.EffectsComplete || tt.invalid {
		tt.result.Reason = analysis.UnsupportedOpcode
		if tt.invalid {
			tt.result.Reason = analysis.InvalidInput
		}
		tt.result.EffectsComplete = false
		tt.result.Flow = analysis.FlowUnknown
		tt.result.Copies = tt.result.Copies[:0]
		return
	}
	tt.write(view(Q, 0, 8), knownbits.Value[uint16]{})
	for resource, mask := range tt.changed {
		tt.finishStorage(Resource(resource), mask)
	}
}

func (tt *transfer) finishStorage(resource Resource, mask uint16) {
	width := StorageWidth(resource)
	for position := knownbits.Width(0); position < width; {
		if mask&(1<<position) == 0 {
			position++
			continue
		}
		end := position + 1
		for end < width && mask&(1<<end) != 0 {
			end++
		}
		part := view(resource, position, end-position)
		value, _ := knownbits.Extract(width, tt.after[resource], part.Offset, part.Bits)
		rule := analysis.Preserve
		if part.Bits == width {
			rule = analysis.Replace
		}
		tt.result.Writes = append(tt.result.Writes, analysis.Write[Resource, uint16]{
			View:  part,
			Value: value,
			Rule:  rule,
		})
		position = end
	}
}

func view(storage Resource, offset, bits knownbits.Width) analysis.View[Resource] {
	return analysis.View[Resource]{
		Storage: storage,
		Offset:  offset,
		Bits:    bits,
	}
}

func byteOffset(high bool) knownbits.Width {
	if high {
		return 8
	}
	return 0
}
