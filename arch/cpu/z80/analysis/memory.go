package analysis

import (
	"github.com/retroenv/retrogolib/analysis/knownbits"
	"github.com/retroenv/retrogolib/arch/cpu/analysis"
)

func (tt *transfer) indirectAddress() Address {
	resource := HL
	if tt.indexed {
		resource = tt.indexResource()
	}

	base := tt.read(view(resource, 0, 16))
	if tt.indexed {
		width := knownbits.Width(8)
		if tt.input.Opcode.Size == 4 {
			width = 16
		}

		displacement, _ := knownbits.Extract(width, tt.input.Operand, 0, 8)
		signed, _ := knownbits.SignExtend(8, 16, displacement)
		sum, _ := knownbits.Add(16, base, signed, knownbits.BitClear)
		base = sum.Value
	}

	return Address{
		Value: base,
	}
}

func (tt *transfer) readMemory(address Address) knownbits.Value[uint16] {
	value := knownbits.Value[uint16]{}
	stable := false
	_, known := address.Value.Constant(16)
	if known && tt.input.Read != nil {
		value, stable = tt.input.Read(address)
		if !value.Valid(8) {
			tt.invalid = true
		}

		if !stable || tt.invalid {
			value = knownbits.Value[uint16]{}
		}
	}

	tt.result.Memory = append(tt.result.Memory, analysis.MemoryEvent[uint16, Address]{
		Address: address,
		Bytes:   1,
		Kind:    analysis.MemoryRead,
		Value:   value,
		Stable:  stable && !tt.invalid,
	})

	return value
}

func (tt *transfer) store(address Address, value knownbits.Value[uint16]) {
	tt.result.Memory = append(tt.result.Memory, analysis.MemoryEvent[uint16, Address]{
		Address: address,
		Bytes:   1,
		Kind:    analysis.MemoryWrite,
		Value:   value,
	})
}

func (tt *transfer) wordMemory(part analysis.View[Resource], load bool) {
	address := Address{
		Value: tt.input.Operand,
	}
	one, _ := knownbits.NewConstant[uint16](16, 1)
	sum, _ := knownbits.Add(16, address.Value, one, knownbits.BitClear)
	next := Address{
		Value: sum.Value,
	}
	if load {
		low := tt.readMemory(address)
		high := tt.readMemory(next)
		value, _ := knownbits.Concat(8, high, 8, low)
		tt.write(part, value)

		return
	}
	value := tt.read(part)
	low, _ := knownbits.Extract(16, value, 0, 8)
	high, _ := knownbits.Extract(16, value, 8, 8)
	tt.store(address, low)
	tt.store(next, high)
}
