package cp1610

import "github.com/retroenv/retrogolib/set"

// Opcode describes one 10-bit instruction word.
type Opcode struct {
	Instruction *Instruction
	OperandKey
	Size          uint8
	Cycles        uint8
	Interruptible bool
}

// Opcodes maps each valid 10-bit word to its instruction.
// Instructions maps each mnemonic to its encodings.
var Opcodes, Instructions = buildOpcodes()

// ReadsMemory reports whether the opcode reads a data operand from memory.
func (opcode Opcode) ReadsMemory(memoryReadInstructions set.Set[string]) bool {
	if opcode.Instruction == nil || opcode.Addressing != DirectAddressing && opcode.Addressing != IndirectAddressing {
		return false
	}
	return memoryReadInstructions.Contains(opcode.Instruction.Name)
}

// WritesMemory reports whether the opcode writes a data operand to memory.
func (opcode Opcode) WritesMemory(memoryWriteInstructions set.Set[string]) bool {
	if opcode.Instruction == nil || opcode.Addressing != DirectAddressing && opcode.Addressing != IndirectAddressing {
		return false
	}
	return memoryWriteInstructions.Contains(opcode.Instruction.Name)
}

// ReadWritesMemory reports whether the opcode reads and writes one data operand.
func (opcode Opcode) ReadWritesMemory(memoryReadWriteInstructions set.Set[string]) bool {
	if opcode.Instruction == nil || opcode.Addressing != DirectAddressing && opcode.Addressing != IndirectAddressing {
		return false
	}
	return memoryReadWriteInstructions.Contains(opcode.Instruction.Name)
}

// GetOpcodeInfo finds an opcode by its 10-bit word.
func GetOpcodeInfo(word uint16) (Opcode, bool) {
	if word >= uint16(len(Opcodes)) {
		return Opcode{}, false
	}
	op := Opcodes[word]
	return op, op.Instruction != nil
}

// GetEncoding finds the instruction word for a mnemonic and operand form.
func GetEncoding(name string, key OperandKey) (OpcodeInfo, bool) {
	ins, ok := Instructions[name]
	if !ok {
		return OpcodeInfo{}, false
	}
	info, ok := ins.Encodings[key]
	return info, ok
}

func buildOpcodes() ([1024]Opcode, map[string]*Instruction) {
	var table [1024]Opcode
	instructions := make(map[string]*Instruction)
	for word := range table {
		name, op, ok := decodeOpcode(uint16(word))
		if !ok {
			continue
		}
		ins := instructions[name]
		if ins == nil {
			ins = &Instruction{
				Name:      name,
				Encodings: make(map[OperandKey]OpcodeInfo),
			}
			instructions[name] = ins
		}
		op.Instruction = ins
		table[word] = op
		ins.Encodings[op.OperandKey] = OpcodeInfo{
			Word:   uint16(word),
			Size:   op.Size,
			Cycles: op.Cycles,
		}
	}
	return table, instructions
}

func decodeOpcode(word uint16) (string, Opcode, bool) {
	switch {
	case word <= 7:
		return decodeControl(word)

	case word < 0x30:
		return decodeUnary(word), Opcode{
			OperandKey: OperandKey{
				Addressing:  RegisterAddressing,
				Destination: uint8(word & 7),
			},
			Size:          1,
			Cycles:        6,
			Interruptible: true,
		}, true

	case word < 0x40:
		return decodeStatus(word)

	case word < 0x80:
		return decodeShift(word)

	case word < 0x200:
		return decodeRegisterPair(word)

	case word < 0x240:
		return decodeBranch(word), Opcode{
			OperandKey: OperandKey{
				Addressing: RelativeAddressing,
				Source:     uint8(word & 0x1F),
				Variant:    uint8((word >> 5) & 1),
			},
			Size:          2,
			Cycles:        7,
			Interruptible: true,
		}, true

	default:
		return decodeMemory(word)
	}
}

func decodeBranch(word uint16) string {
	if word&0x10 != 0 {
		return BextName
	}
	names := [...]string{
		BName, BcName, BovName, BplName, BeqName, BltName, BleName, BuscName,
		NoppName, BncName, BnovName, BmiName, BneqName, BgeName, BgtName, BescName,
	}
	return names[word&0x0F]
}

func decodeControl(word uint16) (string, Opcode, bool) {
	names := [...]string{HltName, SdbdName, EisName, DisName, JumpName, TciName, ClrcName, SetcName}
	op := Opcode{
		OperandKey: OperandKey{Addressing: ImpliedAddressing},
		Size:       1,
		Cycles:     4,
	}
	if word == 4 {
		op.Addressing = SpecialAddressing
		op.Size = 3
		op.Cycles = 12
		op.Interruptible = true
	}
	return names[word], op, true
}

func decodeUnary(word uint16) string {
	names := [...]string{IncrName, DecrName, ComrName, NegrName, AdcrName}
	return names[(word-8)>>3]
}

func decodeStatus(word uint16) (string, Opcode, bool) {
	op := Opcode{
		OperandKey:    OperandKey{Addressing: RegisterAddressing},
		Size:          1,
		Cycles:        6,
		Interruptible: true,
	}
	switch {
	case word < 0x34:
		op.Destination = uint8(word & 3)
		return GswdName, op, true

	case word < 0x36:
		op.Addressing = ImpliedAddressing
		op.Variant = uint8(word & 1)
		return NopName, op, true

	case word < 0x38:
		op.Addressing = ImpliedAddressing
		op.Variant = uint8(word & 1)
		return SinName, op, true

	default:
		op.Source = uint8(word & 7)
		return RswdName, op, true
	}
}

func decodeShift(word uint16) (string, Opcode, bool) {
	names := [...]string{SwapName, SllName, RlcName, SllcName, SlrName, SarName, RrcName, SarcName}
	op := Opcode{
		OperandKey: OperandKey{
			Addressing:  RegisterAddressing,
			Destination: uint8(word & 3),
			Variant:     uint8((word >> 2) & 1),
		},
		Size:   1,
		Cycles: 6,
	}
	if op.Variant != 0 {
		op.Cycles = 8
	}
	return names[(word-0x40)>>3], op, true
}

func decodeRegisterPair(word uint16) (string, Opcode, bool) {
	names := [...]string{MovrName, AddrName, SubrName, CmprName, AndrName, XorrName}
	op := Opcode{
		OperandKey: OperandKey{
			Addressing:  RegisterAddressing,
			Source:      uint8((word >> 3) & 7),
			Destination: uint8(word & 7),
		},
		Size:          1,
		Cycles:        6,
		Interruptible: true,
	}
	if op.Destination >= 6 {
		op.Cycles++
	}
	return names[(word-0x80)>>6], op, true
}

func decodeMemory(word uint16) (string, Opcode, bool) {
	names := [...]string{MvoName, MviName, AddName, SubName, CmpName, AndName, XorName}
	family := (word - 0x240) >> 6
	modeRegister := uint8((word >> 3) & 7)
	op := Opcode{
		OperandKey: OperandKey{
			Addressing:  IndirectAddressing,
			Source:      modeRegister,
			Destination: uint8(word & 7),
		},
		Size:          1,
		Cycles:        8,
		Interruptible: family != 0,
	}
	switch modeRegister {
	case 0:
		op.Addressing = DirectAddressing
		op.Size = 2
		op.Cycles = 10

	case 6:
		op.Cycles = 11

	case 7:
		op.Addressing = ImmediateAddressing
		op.Size = 2
	}
	if family == 0 {
		op.Cycles = 9
		if modeRegister == 0 {
			op.Cycles = 11
		}
	}
	return names[family], op, true
}
