package analysis

import (
	"github.com/retroenv/retrogolib/analysis/knownbits"
	"github.com/retroenv/retrogolib/arch/cpu/analysis"
	"github.com/retroenv/retrogolib/arch/cpu/cpu6502"
)

func (tt *transfer) stack(name string) bool {
	switch name {
	case cpu6502.PhaName:
		tt.push(tt.read(A))

	case cpu6502.PhxName:
		tt.push(tt.read(X))

	case cpu6502.PhyName:
		tt.push(tt.read(Y))

	case cpu6502.PhpName:
		tt.push(stackStatus(tt.read(P)))

	case cpu6502.PlaName, cpu6502.PlxName, cpu6502.PlyName:
		resource := A
		switch name {
		case cpu6502.PlxName:
			resource = X

		case cpu6502.PlyName:
			resource = Y
		}

		value := tt.pull()
		tt.write(resource, value)
		tt.zn(value)

	case cpu6502.PlpName:
		tt.restoreStatus()

	case cpu6502.JsrName:
		tt.result.Flow = analysis.FlowCall
		tt.pushPC(tt.returnPC())

	case cpu6502.RtsName:
		tt.result.Flow = analysis.FlowReturn
		tt.pull()
		tt.pull()

	case cpu6502.RtiName:
		tt.result.Flow = analysis.FlowReturn
		tt.restoreStatus()
		tt.pull()
		tt.pull()

	case cpu6502.BrkName:
		tt.trap()

	default:
		return false
	}

	return true
}

func (tt *transfer) push(value knownbits.Value[uint8]) {
	tt.read(SP)
	address := stackAddress(tt.after[SP])
	tt.store(address, value)
	tt.write(SP, changeByte(tt.after[SP], true))
}

func (tt *transfer) pull() knownbits.Value[uint8] {
	tt.read(SP)
	tt.write(SP, changeByte(tt.after[SP], false))
	address := stackAddress(tt.after[SP])

	return tt.readMemory(address)
}

func (tt *transfer) pushPC(pc knownbits.Value[uint16]) {
	high := knownbits.Value[uint8]{
		KnownZero: uint8(pc.KnownZero >> 8),
		KnownOne:  uint8(pc.KnownOne >> 8),
	}
	low := knownbits.Value[uint8]{
		KnownZero: uint8(pc.KnownZero),
		KnownOne:  uint8(pc.KnownOne),
	}
	tt.push(high)
	tt.push(low)
}

func (tt *transfer) returnPC() knownbits.Value[uint16] {
	two, _ := knownbits.NewConstant[uint16](16, 2)
	address, _ := knownbits.Add(16, tt.input.PC, two, knownbits.BitClear)
	return address.Value
}

func (tt *transfer) restoreStatus() {
	value := tt.pull()
	clearBreak, _ := knownbits.NewConstant[uint8](8, 0xef)
	value, _ = knownbits.And(8, value, clearBreak)
	unused, _ := knownbits.NewConstant[uint8](8, 0x20)
	value, _ = knownbits.Or(8, value, unused)
	tt.write(P, value)
}

func (tt *transfer) trap() {
	tt.result.Flow = analysis.FlowTrap
	tt.pushPC(tt.returnPC())
	tt.push(stackStatus(tt.read(P)))
	tt.writeFlag(2, knownbits.BitSet)
	if tt.input.Variant >= cpu6502.Variant65C02 {
		tt.writeFlag(3, knownbits.BitClear)
	}

	vector, _ := knownbits.NewConstant[uint16](16, cpu6502.IrqAddress)
	tt.pointer(vector, false)
}

func stackStatus(before knownbits.Value[uint8]) knownbits.Value[uint8] {
	forced, _ := knownbits.NewConstant[uint8](8, 0x30)
	result, _ := knownbits.Or(8, before, forced)
	return result
}

func stackAddress(sp knownbits.Value[uint8]) Address {
	return Address{
		Value: knownbits.Value[uint16]{
			KnownZero: 0xfe00 | uint16(sp.KnownZero),
			KnownOne:  0x0100 | uint16(sp.KnownOne),
		},
	}
}
