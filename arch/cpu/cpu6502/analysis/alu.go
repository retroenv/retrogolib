package analysis

import (
	"github.com/retroenv/retrogolib/analysis/knownbits"
	"github.com/retroenv/retrogolib/arch/cpu/analysis"
	"github.com/retroenv/retrogolib/arch/cpu/cpu6502"
)

func (tt *transfer) alu(name string) bool {
	switch name {
	case cpu6502.AndName, cpu6502.OraName, cpu6502.EorName:
		left := tt.read(A)
		right := tt.operand()
		value := logical(name, left, right)
		tt.write(A, value)
		tt.zn(value)
		tt.logicalCopy(name, right)

	case cpu6502.AdcName, cpu6502.SbcName:
		value := tt.arithmetic(tt.read(A), tt.operand(), tt.flag(0), name == cpu6502.SbcName)
		tt.write(A, value)

	case cpu6502.CmpName, cpu6502.CpxName, cpu6502.CpyName:
		resource := A
		switch name {
		case cpu6502.CpxName:
			resource = X

		case cpu6502.CpyName:
			resource = Y
		}
		tt.compare(tt.read(resource), tt.operand())

	case cpu6502.AslName, cpu6502.LsrName, cpu6502.RolName, cpu6502.RorName:
		tt.writeOperand(tt.shift(name, tt.operand()))

	case cpu6502.IncName, cpu6502.DecName:
		value := changeByte(tt.operand(), name == cpu6502.DecName)
		tt.writeOperand(value)
		tt.zn(value)

	case cpu6502.BitName, cpu6502.TrbName, cpu6502.TsbName:
		tt.bitTest(name)

	default:
		return tt.unofficialALU(name)
	}
	return true
}

func (tt *transfer) arithmetic(left, right knownbits.Value[uint8], carry knownbits.Bit,
	subtract bool) knownbits.Value[uint8] {

	if tt.input.Variant != cpu6502.VariantNES6502 {
		decimal := tt.flag(3)
		if decimal != knownbits.BitClear {
			tt.result.Reason = analysis.UnsupportedMode
			if decimal == knownbits.BitUnknown {
				tt.result.Reason = analysis.UnknownMode
			}
			for _, position := range []knownbits.Width{0, 1, 6, 7} {
				tt.writeFlag(position, knownbits.BitUnknown)
			}
			return knownbits.Value[uint8]{}
		}
	}
	var value knownbits.Value[uint8]
	var outgoing, overflow knownbits.Bit
	if subtract {
		difference, _ := knownbits.Sub(8, left, right, invert(carry))
		value = difference.Value
		outgoing = invert(difference.BorrowOut)
		overflow = difference.SignedOverflow
	} else {
		addition, _ := knownbits.Add(8, left, right, carry)
		value = addition.Value
		outgoing = addition.CarryOut
		overflow = addition.SignedOverflow
	}
	tt.writeFlag(0, outgoing)
	tt.writeFlag(6, overflow)
	tt.zn(value)
	return value
}

func (tt *transfer) compare(left, right knownbits.Value[uint8]) {
	difference, _ := knownbits.Sub(8, left, right, knownbits.BitClear)
	tt.writeFlag(0, invert(difference.BorrowOut))
	tt.zn(difference.Value)
}

func (tt *transfer) shift(name string, value knownbits.Value[uint8]) knownbits.Value[uint8] {
	var result knownbits.Value[uint8]
	var carry knownbits.Bit
	switch name {
	case cpu6502.AslName:
		result, _ = knownbits.ShiftLeft(8, value, 1)
		carry, _ = knownbits.TestBit(8, value, 7)

	case cpu6502.LsrName:
		result, _ = knownbits.ShiftRight(8, value, 1)
		carry, _ = knownbits.TestBit(8, value, 0)

	case cpu6502.RolName:
		result, carry, _ = knownbits.RotateLeftCarry(8, value, tt.flag(0))

	case cpu6502.RorName:
		result, carry, _ = knownbits.RotateRightCarry(8, value, tt.flag(0))
	}
	tt.writeFlag(0, carry)
	tt.zn(result)
	return result
}

func (tt *transfer) bitTest(name string) {
	left := tt.read(A)
	right := tt.operand()
	intersection, _ := knownbits.And(8, left, right)
	zero, _ := knownbits.Zero(8, intersection)
	tt.writeFlag(1, zero)
	switch name {
	case cpu6502.TrbName:
		complement, _ := knownbits.Not(8, left)
		value, _ := knownbits.And(8, right, complement)
		tt.writeOperand(value)

	case cpu6502.TsbName:
		value, _ := knownbits.Or(8, left, right)
		tt.writeOperand(value)

	case cpu6502.BitName:
		if tt.opcode.Addressing != cpu6502.ImmediateAddressing {
			overflow, _ := knownbits.TestBit(8, right, 6)
			negative, _ := knownbits.Sign(8, right)
			tt.writeFlag(6, overflow)
			tt.writeFlag(7, negative)
		}
	}
}

func (tt *transfer) logicalCopy(name string, right knownbits.Value[uint8]) {
	constant, known := right.Constant(8)
	identity := name == cpu6502.AndName && constant == 0xff || name != cpu6502.AndName && constant == 0
	if !known || !identity {
		return
	}
	tt.result.Copies = append(tt.result.Copies, analysis.Copy[Resource]{
		Source:      byteView(A),
		Destination: byteView(A),
	})
}

func logical(name string, left, right knownbits.Value[uint8]) knownbits.Value[uint8] {
	switch name {
	case cpu6502.AndName:
		result, _ := knownbits.And(8, left, right)
		return result

	case cpu6502.OraName:
		result, _ := knownbits.Or(8, left, right)
		return result

	default:
		result, _ := knownbits.Xor(8, left, right)
		return result
	}
}

func changeByte(value knownbits.Value[uint8], subtract bool) knownbits.Value[uint8] {
	one, _ := knownbits.NewConstant[uint8](8, 1)
	if subtract {
		result, _ := knownbits.Sub(8, value, one, knownbits.BitClear)
		return result.Value
	}
	result, _ := knownbits.Add(8, value, one, knownbits.BitClear)
	return result.Value
}
