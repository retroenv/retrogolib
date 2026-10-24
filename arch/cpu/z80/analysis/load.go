package analysis

import (
	"github.com/retroenv/retrogolib/analysis/knownbits"
	"github.com/retroenv/retrogolib/arch/cpu/analysis"
)

func (tt *transfer) byteLoad(code byte) {
	destination := (code >> 3) & 7
	source := code & 7
	memory := destination == 6 || source == 6
	switch {
	case source == 6:
		value := tt.readMemory(tt.indirectAddress())
		tt.write(tt.byteRegister(destination, true), value)

	case destination == 6:
		value := tt.read(tt.byteRegister(source, true))
		tt.store(tt.indirectAddress(), value)

	default:
		tt.copy(tt.byteRegister(destination, memory), tt.byteRegister(source, memory))
	}
}

func (tt *transfer) immediateByteLoad(code byte) {
	destination := (code >> 3) & 7
	value := tt.input.Operand
	if tt.indexed && destination == 6 {
		value, _ = knownbits.Extract(16, value, 8, 8)
	}

	if destination == 6 {
		tt.store(tt.indirectAddress(), value)
	} else {
		tt.write(tt.byteRegister(destination, false), value)
	}
}

func (tt *transfer) otherLoad(code byte) bool {
	accumulator := view(AF, 8, 8)
	switch code {
	case 0x02, 0x12, 0x0a, 0x1a:
		address := Address{
			Value: tt.read(tt.pair((code >> 4) & 1)),
		}
		if code&8 == 0 {
			tt.store(address, tt.read(accumulator))
		} else {
			tt.write(accumulator, tt.readMemory(address))
		}

	case 0x32, 0x3a:
		address := Address{
			Value: tt.input.Operand,
		}
		if code == 0x32 {
			tt.store(address, tt.read(accumulator))
		} else {
			tt.write(accumulator, tt.readMemory(address))
		}

	case 0x22, 0x2a:
		tt.wordMemory(tt.pair(2), code == 0x2a)

	case 0xf9:
		tt.copy(view(SP, 0, 16), tt.pair(2))

	default:
		return false
	}

	return true
}

func (tt *transfer) extendedLoad() bool {
	code := tt.input.Opcode.Opcode
	if code&0xc7 == 0x43 {
		tt.wordMemory(tt.pair((code>>4)&3), code&8 != 0)
		return true
	}
	if code != 0x47 && code != 0x4f && code != 0x57 && code != 0x5f {
		return false
	}
	resource := I
	if code&8 != 0 {
		resource = R
	}

	accumulator := view(AF, 8, 8)
	if code&0x10 == 0 {
		tt.copy(view(resource, 0, 8), accumulator)
		return true
	}
	tt.copy(accumulator, view(resource, 0, 8))
	value := tt.after[AF]
	value, _ = knownbits.Extract(16, value, 8, 8)
	tt.szxy(value)
	tt.writeFlag(4, knownbits.BitClear)
	tt.writeFlag(1, knownbits.BitClear)
	iff := tt.read(view(IFF2, 0, 1))
	flag, _ := knownbits.TestBit(1, iff, 0)
	tt.writeFlag(2, flag)

	return true
}

func (tt *transfer) logicalCopy(operation byte, right knownbits.Value[uint16], self bool) {
	constant, known := right.Constant(8)
	identity := self && operation != 1
	identity = identity || known && (operation == 0 && constant == 0xff || operation != 0 && constant == 0)
	if !identity {
		return
	}
	tt.result.Copies = append(tt.result.Copies, analysis.Copy[Resource]{
		Source:      view(AF, 8, 8),
		Destination: view(AF, 8, 8),
	})
}
