package analysis

import (
	"github.com/retroenv/retrogolib/analysis/knownbits"
	"github.com/retroenv/retrogolib/arch/cpu/analysis"
	"github.com/retroenv/retrogolib/arch/cpu/z80"
)

func (tt *transfer) execute() bool {
	if tt.input.Opcode.Prefix == z80.PrefixED {
		return tt.extendedLoad()
	}
	code := tt.input.Opcode.Opcode
	switch {
	case code >= 0x40 && code <= 0x7f && code != 0x76:
		tt.byteLoad(code)

	case code&0xc7 == 0x06:
		tt.immediateByteLoad(code)

	case code&0xcf == 0x01:
		tt.write(tt.pair((code>>4)&3), tt.input.Operand)

	case code >= 0xa0 && code <= 0xbf || code&0xe7 == 0xe6:
		tt.byteALU(code)

	case code == 0x37 || code == 0x3f:
		tt.carryControl(code == 0x3f)

	case code == 0:
		return true

	default:
		return tt.otherLoad(code)
	}

	return true
}

func (tt *transfer) indexResource() Resource {
	if tt.input.Opcode.Prefix == z80.PrefixFD {
		return IY
	}
	return IX
}

func (tt *transfer) pair(code byte) analysis.View[Resource] {
	resource := [...]Resource{BC, DE, HL, SP}[code]
	if tt.indexed && resource == HL {
		resource = tt.indexResource()
	}

	return view(resource, 0, 16)
}

func (tt *transfer) byteRegister(code byte, memoryForm bool) analysis.View[Resource] {
	registers := [...]z80.RegisterParam{z80.RegB, z80.RegC, z80.RegD, z80.RegE, z80.RegH, z80.RegL, z80.RegNone, z80.RegA}
	part, _ := RegisterView(registers[code])
	if tt.indexed && !memoryForm && (code == 4 || code == 5) {
		part.Storage = tt.indexResource()
	}

	return part
}

func decode(input Input) (z80.Opcode, bool, analysis.Reason) {
	opcode, indexed, reason := decodeOpcode(input.Opcode)
	if reason != analysis.Supported {
		return opcode, indexed, reason
	}
	return opcode, indexed, validateInput(input, opcode)
}

func decodeOpcode(info z80.OpcodeInfo) (z80.Opcode, bool, analysis.Reason) {
	opcode := z80.Opcodes[info.Opcode]
	indexed := false
	prefix := info.Prefix
	switch prefix {
	case 0:

	case z80.PrefixED:
		opcode = z80.EDOpcodes[info.Opcode]

	case z80.PrefixDD, z80.PrefixFD:
		table := &z80.DDOpcodes
		if prefix == z80.PrefixFD {
			table = &z80.FDOpcodes
		}

		candidate := table[info.Opcode]
		if candidate.Instruction != nil {
			opcode = candidate
			indexed = true
		} else {
			opcode.Size++
		}

	default:
		return z80.Opcode{}, false, analysis.UnsupportedOpcode
	}

	if opcode.Instruction == nil || opcode.Size == 0 || info.Opcode == z80.PrefixDD ||
		info.Opcode == z80.PrefixFD || info.Opcode == z80.PrefixCB || info.Opcode == z80.PrefixED {

		return opcode, indexed, analysis.UnsupportedOpcode
	}
	return opcode, indexed, analysis.Supported
}

func validateInput(input Input, opcode z80.Opcode) analysis.Reason {
	if input.Opcode.Size != opcode.Size {
		return analysis.IllegalOperand
	}
	operandBytes := opcode.Size - 1
	if input.Opcode.Prefix != 0 {
		operandBytes--
	}

	width := knownbits.Width(8)
	if operandBytes == 2 {
		width = 16
	}

	if operandBytes > 2 || !input.Operand.Valid(width) {
		return analysis.InvalidInput
	}
	for resource, value := range input.Before {
		if !value.Valid(StorageWidth(Resource(resource))) {
			return analysis.InvalidInput
		}
	}

	return analysis.Supported
}
