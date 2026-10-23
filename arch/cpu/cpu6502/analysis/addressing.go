package analysis

import (
	"github.com/retroenv/retrogolib/analysis/knownbits"
	"github.com/retroenv/retrogolib/arch/cpu/analysis"
	"github.com/retroenv/retrogolib/arch/cpu/cpu6502"
)

func (tt *transfer) operand() knownbits.Value[uint8] {
	switch tt.opcode.Addressing {
	case cpu6502.ImmediateAddressing:
		return knownbits.Value[uint8]{
			KnownZero: uint8(tt.input.Operand.KnownZero),
			KnownOne:  uint8(tt.input.Operand.KnownOne),
		}

	case cpu6502.AccumulatorAddressing:
		return tt.read(A)

	default:
		return tt.readMemory(tt.address())
	}
}

func (tt *transfer) writeOperand(value knownbits.Value[uint8]) {
	if tt.opcode.Addressing == cpu6502.AccumulatorAddressing {
		tt.write(A, value)
		return
	}
	tt.store(tt.address(), value)
}

func (tt *transfer) address() Address {
	if tt.addressReady {
		return tt.operandAddress
	}
	base := tt.input.Operand
	if encodedSize(tt.opcode.Addressing) == 2 {
		base, _ = knownbits.ZeroExtend(8, 16, base)
	}
	switch tt.opcode.Addressing {
	case cpu6502.ZeroPageXAddressing:
		base = tt.index(base, X, true)

	case cpu6502.ZeroPageYAddressing:
		base = tt.index(base, Y, true)

	case cpu6502.AbsoluteXAddressing:
		base = tt.index(base, X, false)

	case cpu6502.AbsoluteYAddressing:
		base = tt.index(base, Y, false)

	case cpu6502.IndirectXAddressing:
		base = tt.pointer(tt.index(base, X, true), true)

	case cpu6502.IndirectYAddressing:
		base = tt.index(tt.pointer(base, true), Y, false)

	case cpu6502.ZeroPageIndirectAddressing:
		base = tt.pointer(base, true)

	case cpu6502.IndirectAddressing:
		base = tt.pointer(base, tt.input.Variant < cpu6502.Variant65C02)

	case cpu6502.AbsoluteXIndirectAddressing:
		base = tt.pointer(tt.index(base, X, false), false)
	}
	tt.operandAddress = Address{
		Value: base,
	}
	tt.addressReady = true
	return tt.operandAddress
}

func (tt *transfer) index(base knownbits.Value[uint16], resource Resource, zeroPage bool) knownbits.Value[uint16] {
	index := tt.read(resource)
	part := knownbits.Value[uint16]{
		KnownZero: uint16(index.KnownZero),
		KnownOne:  uint16(index.KnownOne),
	}
	if zeroPage {
		low, _ := knownbits.Truncate(16, 8, base)
		sum, _ := knownbits.Add(8, low, part, knownbits.BitClear)
		result, _ := knownbits.ZeroExtend(8, 16, sum.Value)
		return result
	}
	wide, _ := knownbits.ZeroExtend(8, 16, part)
	sum, _ := knownbits.Add(16, base, wide, knownbits.BitClear)
	return sum.Value
}

func (tt *transfer) pointer(base knownbits.Value[uint16], wrapPage bool) knownbits.Value[uint16] {
	low := tt.readMemory(Address{
		Value: base,
	})
	next := incrementAddress(base, wrapPage)
	high := tt.readMemory(Address{
		Value: next,
	})
	return knownbits.Value[uint16]{
		KnownZero: uint16(high.KnownZero)<<8 | uint16(low.KnownZero),
		KnownOne:  uint16(high.KnownOne)<<8 | uint16(low.KnownOne),
	}
}

func (tt *transfer) readMemory(address Address) knownbits.Value[uint8] {
	value := knownbits.Value[uint8]{}
	stable := false
	_, known := address.Value.Constant(16)
	if known && tt.input.Read != nil {
		value, stable = tt.input.Read(address)
		if !value.Valid(8) {
			tt.invalid = true
		}
		if !stable || tt.invalid {
			value = knownbits.Value[uint8]{}
		}
	}
	tt.result.Memory = append(tt.result.Memory, analysis.MemoryEvent[uint8, Address]{
		Address: address,
		Bytes:   1,
		Kind:    analysis.MemoryRead,
		Value:   value,
		Stable:  stable && !tt.invalid,
	})
	return value
}

func (tt *transfer) store(address Address, value knownbits.Value[uint8]) {
	tt.result.Memory = append(tt.result.Memory, analysis.MemoryEvent[uint8, Address]{
		Address: address,
		Bytes:   1,
		Kind:    analysis.MemoryWrite,
		Value:   value,
	})
}

func incrementAddress(base knownbits.Value[uint16], wrapPage bool) knownbits.Value[uint16] {
	one, _ := knownbits.NewConstant[uint16](16, 1)
	if !wrapPage {
		sum, _ := knownbits.Add(16, base, one, knownbits.BitClear)
		return sum.Value
	}
	low, _ := knownbits.Extract(16, base, 0, 8)
	one, _ = knownbits.NewConstant[uint16](8, 1)
	sum, _ := knownbits.Add(8, low, one, knownbits.BitClear)
	result, _ := knownbits.Insert(16, base, 0, 8, sum.Value)
	return result
}
