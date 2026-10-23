package analysis

import (
	"github.com/retroenv/retrogolib/analysis/knownbits"
	"github.com/retroenv/retrogolib/arch/cpu/cpu6502"
)

func (tt *transfer) unofficialALU(name string) bool {
	switch name {
	case cpu6502.DcpName, cpu6502.IscName:
		value := changeByte(tt.operand(), name == cpu6502.DcpName)
		tt.writeOperand(value)
		if name == cpu6502.DcpName {
			tt.compare(tt.read(A), value)
		} else {
			tt.write(A, tt.arithmetic(tt.read(A), value, tt.flag(0), true))
		}

	case cpu6502.SloName, cpu6502.SreName, cpu6502.RlaName, cpu6502.RraName:
		tt.combinedShift(name)

	case cpu6502.AncName, cpu6502.AlrName:
		value, _ := knownbits.And(8, tt.read(A), tt.operand())
		if name == cpu6502.AlrName {
			value = tt.shift(cpu6502.LsrName, value)
		} else {
			carry, _ := knownbits.Sign(8, value)
			tt.writeFlag(0, carry)
			tt.zn(value)
		}
		tt.write(A, value)

	case cpu6502.AxsName:
		left, _ := knownbits.And(8, tt.read(A), tt.read(X))
		difference, _ := knownbits.Sub(8, left, tt.operand(), knownbits.BitClear)
		tt.write(X, difference.Value)
		tt.writeFlag(0, invert(difference.BorrowOut))
		tt.zn(difference.Value)

	case cpu6502.LasName:
		value, _ := knownbits.And(8, tt.operand(), tt.read(SP))
		tt.write(A, value)
		tt.write(X, value)
		tt.write(SP, value)
		tt.zn(value)

	default:
		return false
	}
	return true
}

func (tt *transfer) combinedShift(name string) {
	shift := cpu6502.AslName
	operation := cpu6502.OraName
	switch name {
	case cpu6502.SreName:
		shift = cpu6502.LsrName
		operation = cpu6502.EorName

	case cpu6502.RlaName:
		shift = cpu6502.RolName
		operation = cpu6502.AndName

	case cpu6502.RraName:
		shift = cpu6502.RorName
	}
	value := tt.shift(shift, tt.operand())
	tt.writeOperand(value)
	left := tt.read(A)
	if name == cpu6502.RraName {
		carry, _ := knownbits.TestBit(8, tt.after[P], 0)
		tt.write(A, tt.arithmetic(left, value, carry, false))
		return
	}
	result := logical(operation, left, value)
	tt.write(A, result)
	tt.zn(result)
}
