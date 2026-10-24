package analysis

import (
	"github.com/retroenv/retrogolib/analysis/knownbits"
	"github.com/retroenv/retrogolib/arch/cpu/analysis"
	"github.com/retroenv/retrogolib/arch/cpu/cpu6502"
)

func (tt *transfer) execute() bool {
	name := tt.opcode.Instruction.Name
	return tt.loadStore(name) || tt.transfers(name) || tt.alu(name) || tt.system(name) || tt.stack(name) || tt.bitMemory()
}

func (tt *transfer) loadStore(name string) bool {
	switch name {
	case cpu6502.LdaName, cpu6502.LdxName, cpu6502.LdyName, cpu6502.LaxName:
		resource := A
		switch name {
		case cpu6502.LdxName:
			resource = X

		case cpu6502.LdyName:
			resource = Y
		}

		value := tt.operand()
		tt.write(resource, value)
		tt.zn(value)
		if name == cpu6502.LaxName {
			tt.write(X, value)
		}

	case cpu6502.StaName, cpu6502.StxName, cpu6502.StyName, cpu6502.SaxName, cpu6502.StzName:
		value, _ := knownbits.NewConstant[uint8](8, 0)
		switch name {
		case cpu6502.StaName:
			value = tt.read(A)

		case cpu6502.StxName:
			value = tt.read(X)

		case cpu6502.StyName:
			value = tt.read(Y)

		case cpu6502.SaxName:
			value, _ = knownbits.And(8, tt.read(A), tt.read(X))
		}

		tt.writeOperand(value)

	default:
		return false
	}

	return true
}

func (tt *transfer) transfers(name string) bool {
	switch name {
	case cpu6502.TaxName:
		tt.copy(X, A, true)

	case cpu6502.TayName:
		tt.copy(Y, A, true)

	case cpu6502.TxaName:
		tt.copy(A, X, true)

	case cpu6502.TyaName:
		tt.copy(A, Y, true)

	case cpu6502.TsxName:
		tt.copy(X, SP, true)

	case cpu6502.TxsName:
		tt.copy(SP, X, false)

	case cpu6502.InxName, cpu6502.DexName:
		tt.delta(X, name == cpu6502.DexName)

	case cpu6502.InyName, cpu6502.DeyName:
		tt.delta(Y, name == cpu6502.DeyName)

	default:
		return false
	}

	return true
}

func (tt *transfer) delta(resource Resource, subtract bool) {
	value := changeByte(tt.read(resource), subtract)
	tt.write(resource, value)
	tt.zn(value)
}

func (tt *transfer) system(name string) bool {
	switch name {
	case cpu6502.ClcName, cpu6502.SecName:
		value := knownbits.BitClear
		if name == cpu6502.SecName {
			value = knownbits.BitSet
		}

		tt.writeFlag(0, value)

	case cpu6502.CldName, cpu6502.SedName:
		value := knownbits.BitClear
		if name == cpu6502.SedName {
			value = knownbits.BitSet
		}

		tt.writeFlag(3, value)

	case cpu6502.CliName, cpu6502.SeiName:
		value := knownbits.BitClear
		if name == cpu6502.SeiName {
			value = knownbits.BitSet
		}

		tt.writeFlag(2, value)

	case cpu6502.ClvName:
		tt.writeFlag(6, knownbits.BitClear)

	case cpu6502.NopName:
		if tt.opcode.ReadsMemory(cpu6502.MemoryReadInstructions) {
			tt.operand()
		}

	case cpu6502.JmpName:
		tt.result.Flow = analysis.FlowJump
		if tt.opcode.Addressing != cpu6502.AbsoluteAddressing {
			tt.address()
		}

	case cpu6502.KilName:
		tt.result.Flow = analysis.FlowStop

	default:
		return tt.branch(name)
	}

	return true
}

func (tt *transfer) branch(name string) bool {
	position := knownbits.Width(0)
	var set bool
	switch name {
	case cpu6502.BccName, cpu6502.BcsName:
		set = name == cpu6502.BcsName

	case cpu6502.BeqName, cpu6502.BneName:
		position = 1
		set = name == cpu6502.BeqName

	case cpu6502.BmiName, cpu6502.BplName:
		position = 7
		set = name == cpu6502.BmiName

	case cpu6502.BvsName, cpu6502.BvcName:
		position = 6
		set = name == cpu6502.BvsName

	case cpu6502.BraName:
		tt.result.Flow = analysis.FlowBranch
		tt.result.Condition = knownbits.BitSet
		return true

	default:
		return false
	}

	tt.result.Flow = analysis.FlowBranch
	condition := tt.flag(position)
	if !set {
		condition = invert(condition)
	}

	tt.result.Condition = condition

	return true
}
