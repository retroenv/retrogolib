package cpu68000

// DecodedOpcode represents a fully decoded 68000 opcode.
type DecodedOpcode struct {
	Instruction *Instruction
	Size        OperandSize
	SrcMode     uint8  // Source EA mode (3 bits)
	SrcReg      uint8  // Source EA register (3 bits)
	DstMode     uint8  // Destination EA mode (3 bits)
	DstReg      uint8  // Destination EA register (3 bits)
	Timing      uint16 // Base T-states
	Extra       uint16 // Extra data embedded in opcode (condition, quick value, etc.)
}

// eaClass is a set of effective-address categories from the M68000PRM.
// Each instruction accepts only the addressing modes that have all of its
// required categories.
type eaClass uint8

const (
	eaData eaClass = 1 << iota
	eaMemory
	eaControl
	eaAlterable
)

const (
	eaAny              eaClass = 0
	eaDataAlterable            = eaData | eaAlterable
	eaMemoryAlterable          = eaMemory | eaAlterable
	eaControlAlterable         = eaControl | eaAlterable
)

// illegalTiming is the exception processing time for an unassigned opcode.
const illegalTiming = 34

// lineDecoders maps the top 4 bits of an opcode word to a line decoder function.
var lineDecoders = [16]func(opcode uint16) DecodedOpcode{
	decodeLine0, // ORI, ANDI, SUBI, ADDI, EORI, CMPI, BTST/BSET/BCLR/BCHG, MOVEP
	decodeLine1, // MOVE.B
	decodeLine2, // MOVE.L, MOVEA.L
	decodeLine3, // MOVE.W, MOVEA.W
	decodeLine4, // Miscellaneous
	decodeLine5, // ADDQ, SUBQ, Scc, DBcc
	decodeLine6, // Bcc, BRA, BSR
	decodeLine7, // MOVEQ
	decodeLine8, // OR, DIV, SBCD
	decodeLine9, // SUB, SUBA, SUBX
	decodeLineA, // Line A trap
	decodeLineB, // CMP, CMPA, CMPM, EOR
	decodeLineC, // AND, MUL, ABCD, EXG
	decodeLineD, // ADD, ADDA, ADDX
	decodeLineE, // Shift/Rotate
	decodeLineF, // Line F trap
}

// immediateInstructions maps bits 11-9 of a line 0 immediate opcode to the
// instruction. Entry 4 is the immediate bit operation group.
var immediateInstructions = [8]*Instruction{
	insORI, insANDI, insSUBI, insADDI, nil, insEORI, insCMPI, nil,
}

// bitInstructions maps bits 7-6 of a bit operation opcode to the instruction.
var bitInstructions = [4]*Instruction{insBTST, insBCHG, insBCLR, insBSET}

// decodeOpcode decodes a 16-bit opcode word. An unassigned word or an
// addressing mode that the instruction does not accept decodes to ILLEGAL.
func decodeOpcode(opcode uint16) DecodedOpcode {
	line := (opcode >> 12) & 0xF
	return lineDecoders[line](opcode)
}

// eaClasses returns the categories of an addressing mode. It returns zero for
// an encoding that the 68000 does not define.
func eaClasses(mode, reg uint8) eaClass {
	switch mode {
	case 0:
		return eaDataAlterable
	case 1:
		return eaAlterable
	case 2, 5, 6:
		return eaData | eaMemory | eaControl | eaAlterable
	case 3, 4:
		return eaData | eaMemory | eaAlterable

	case 7:
		switch reg {
		case 0, 1:
			return eaData | eaMemory | eaControl | eaAlterable
		case 2, 3:
			return eaData | eaMemory | eaControl
		case 4:
			return eaData | eaMemory
		}
	}
	return 0
}

// validEA reports whether the addressing mode exists and has all required classes.
func validEA(mode, reg uint8, required eaClass) bool {
	classes := eaClasses(mode, reg)
	return classes != 0 && classes&required == required
}

// requireEA returns d when the addressing mode has the required classes.
// Otherwise it returns the ILLEGAL decoding.
func requireEA(d DecodedOpcode, mode, reg uint8, required eaClass) DecodedOpcode {
	if !validEA(mode, reg, required) {
		return illegalOpcode()
	}
	return d
}

// illegalOpcode returns the decoding of an unassigned opcode word.
func illegalOpcode() DecodedOpcode {
	return DecodedOpcode{
		Instruction: insILLEGAL,
		Timing:      illegalTiming,
	}
}

// decodeLine0 decodes line 0: immediate operations and bit operations.
func decodeLine0(opcode uint16) DecodedOpcode {
	if opcode&0xF138 == 0x0108 {
		return decodeLine0Movep(opcode)
	}
	// Bits 11-8 determine the specific operation.
	if opcode&0x0100 != 0 {
		// Bit operations with register (BTST/BCHG/BCLR/BSET Dn,<ea>).
		return decodeLine0BitReg(opcode)
	}

	// Immediate operations.
	return decodeLine0Immediate(opcode)
}

// decodeLine0BitReg decodes bit operations with register source.
// BTST accepts every data address, including an immediate operand.
func decodeLine0BitReg(opcode uint16) DecodedOpcode {
	dn := (opcode >> 9) & 7
	mode := uint8((opcode >> 3) & 7)
	reg := uint8(opcode & 7)
	ins := bitInstructions[(opcode>>6)&3]

	d := DecodedOpcode{
		Instruction: ins,
		Size:        SizeByte,
		SrcMode:     0, // Data register
		SrcReg:      uint8(dn),
		DstMode:     mode,
		DstReg:      reg,
		Timing:      8,
	}
	if ins == insBTST {
		return requireEA(d, mode, reg, eaData)
	}
	return requireEA(d, mode, reg, eaDataAlterable)
}

// decodeLine0Movep decodes MOVEP instruction.
func decodeLine0Movep(opcode uint16) DecodedOpcode {
	dn := (opcode >> 9) & 7
	an := opcode & 7
	dir := (opcode >> 7) & 1
	sz := SizeWord
	if opcode&0x0040 != 0 {
		sz = SizeLong
	}

	d := DecodedOpcode{
		Instruction: insMOVEP,
		Size:        sz,
		Timing:      16,
	}

	if dir != 0 {
		// MOVEP Dn,d16(An)
		d.SrcMode = 0
		d.SrcReg = uint8(dn)
		d.DstMode = 5
		d.DstReg = uint8(an)
	} else {
		// MOVEP d16(An),Dn
		d.SrcMode = 5
		d.SrcReg = uint8(an)
		d.DstMode = 0
		d.DstReg = uint8(dn)
	}

	return d
}

// decodeLine0Immediate decodes immediate operations (ORI, ANDI, SUBI, ADDI, EORI, CMPI)
// and immediate bit operations (BTST/BCHG/BCLR/BSET #imm,<ea>).
func decodeLine0Immediate(opcode uint16) DecodedOpcode {
	op := (opcode >> 9) & 7
	if op == 4 {
		return decodeLine0BitImm(opcode)
	}

	ins := immediateInstructions[op]
	size := sizeFromBits((opcode >> 6) & 3)
	mode := uint8((opcode >> 3) & 7)
	reg := uint8(opcode & 7)
	if ins == nil || size == 0 {
		return illegalOpcode()
	}

	d := DecodedOpcode{
		Instruction: ins,
		Size:        size,
		DstMode:     mode,
		DstReg:      reg,
		Timing:      8,
	}
	if size == SizeLong {
		d.Timing = 16
	}

	if mode == 7 && reg == 4 {
		// Only ORI, ANDI, and EORI have byte and word forms for CCR and SR.
		if size == SizeLong || ins == insSUBI || ins == insADDI || ins == insCMPI {
			return illegalOpcode()
		}
		return d
	}

	return requireEA(d, mode, reg, eaDataAlterable)
}

// decodeLine0BitImm decodes bit operations with immediate source.
// BTST #imm accepts every data address except a second immediate operand.
func decodeLine0BitImm(opcode uint16) DecodedOpcode {
	mode := uint8((opcode >> 3) & 7)
	reg := uint8(opcode & 7)
	ins := bitInstructions[(opcode>>6)&3]

	d := DecodedOpcode{
		Instruction: ins,
		Size:        SizeByte,
		SrcMode:     7,
		SrcReg:      4, // Immediate
		DstMode:     mode,
		DstReg:      reg,
		Timing:      12,
	}
	if ins == insBTST {
		if mode == 7 && reg == 4 {
			return illegalOpcode()
		}
		return requireEA(d, mode, reg, eaData)
	}
	return requireEA(d, mode, reg, eaDataAlterable)
}

// decodeLine1 decodes line 1: MOVE.B.
func decodeLine1(opcode uint16) DecodedOpcode {
	return decodeMOVE(opcode, SizeByte)
}

// decodeLine2 decodes line 2: MOVE.L and MOVEA.L.
func decodeLine2(opcode uint16) DecodedOpcode {
	dstMode := (opcode >> 6) & 7
	if dstMode == 1 {
		return decodeMOVEA(opcode, SizeLong)
	}
	return decodeMOVE(opcode, SizeLong)
}

// decodeLine3 decodes line 3: MOVE.W and MOVEA.W.
func decodeLine3(opcode uint16) DecodedOpcode {
	dstMode := (opcode >> 6) & 7
	if dstMode == 1 {
		return decodeMOVEA(opcode, SizeWord)
	}
	return decodeMOVE(opcode, SizeWord)
}

// decodeMOVE decodes a MOVE instruction. Note: MOVE uses a special encoding
// where destination is encoded as register:mode (reversed from source).
// A byte source cannot be an address register.
func decodeMOVE(opcode uint16, size OperandSize) DecodedOpcode {
	srcMode := uint8((opcode >> 3) & 7)
	srcReg := uint8(opcode & 7)
	dstReg := uint8((opcode >> 9) & 7)
	dstMode := uint8((opcode >> 6) & 7)

	source := eaAny
	if size == SizeByte {
		source = eaData
	}
	if !validEA(srcMode, srcReg, source) {
		return illegalOpcode()
	}

	d := DecodedOpcode{
		Instruction: insMOVE,
		Size:        size,
		SrcMode:     srcMode,
		SrcReg:      srcReg,
		DstMode:     dstMode,
		DstReg:      dstReg,
		Timing:      4,
	}
	return requireEA(d, dstMode, dstReg, eaDataAlterable)
}

// decodeMOVEA decodes a MOVEA instruction.
func decodeMOVEA(opcode uint16, size OperandSize) DecodedOpcode {
	srcMode := uint8((opcode >> 3) & 7)
	srcReg := uint8(opcode & 7)
	dstReg := (opcode >> 9) & 7

	d := DecodedOpcode{
		Instruction: insMOVEA,
		Size:        size,
		SrcMode:     srcMode,
		SrcReg:      srcReg,
		DstMode:     1, // Address register direct
		DstReg:      uint8(dstReg),
		Timing:      4,
	}
	return requireEA(d, srcMode, srcReg, eaAny)
}

// decodeLine4 decodes line 4: Miscellaneous instructions.
func decodeLine4(opcode uint16) DecodedOpcode {
	mode := uint8((opcode >> 3) & 7)
	reg := uint8(opcode & 7)

	// Check for specific encodings.
	switch {
	case opcode&0xFFF8 == 0x4E70:
		return decodeLine4Special(opcode)
	case opcode == 0x4AFC:
		return illegalOpcode()

	case opcode&0xFFF0 == 0x4E40:
		return DecodedOpcode{
			Instruction: insTRAP,
			Extra:       opcode & 0x0F,
			Timing:      34,
		}

	case opcode&0xFFF8 == 0x4E50:
		return DecodedOpcode{
			Instruction: insLINK,
			DstReg:      reg,
			Timing:      16,
		}

	case opcode&0xFFF8 == 0x4E58:
		return DecodedOpcode{
			Instruction: insUNLK,
			DstReg:      reg,
			Timing:      12,
		}

	case opcode&0xFFF8 == 0x4E60:
		return decodeLine4MoveUSP(reg, true)
	case opcode&0xFFF8 == 0x4E68:
		return decodeLine4MoveUSP(reg, false)

	case opcode&0xFFC0 == 0x4E80:
		return requireEA(DecodedOpcode{
			Instruction: insJSR,
			DstMode:     mode,
			DstReg:      reg,
			Timing:      16,
		}, mode, reg, eaControl)

	case opcode&0xFFC0 == 0x4EC0:
		return requireEA(DecodedOpcode{
			Instruction: insJMP,
			DstMode:     mode,
			DstReg:      reg,
			Timing:      8,
		}, mode, reg, eaControl)
	}

	return decodeLine4Group(opcode, mode, reg)
}

// decodeLine4MoveUSP decodes MOVE An,USP and MOVE USP,An.
func decodeLine4MoveUSP(reg uint8, toUSP bool) DecodedOpcode {
	if toUSP {
		return DecodedOpcode{
			Instruction: insMOVE,
			Size:        SizeLong,
			SrcMode:     1,
			SrcReg:      reg,
			DstMode:     7,
			DstReg:      5,
			Extra:       1,
			Timing:      4,
		}
	}
	return DecodedOpcode{
		Instruction: insMOVE,
		Size:        SizeLong,
		SrcMode:     7,
		SrcReg:      5,
		DstMode:     1,
		DstReg:      reg,
		Extra:       2,
		Timing:      4,
	}
}

// decodeLine4Special decodes the special instructions at 0x4E7x.
func decodeLine4Special(opcode uint16) DecodedOpcode {
	switch opcode {
	case 0x4E70:
		return DecodedOpcode{
			Instruction: insRESET,
			Timing:      132,
		}

	case 0x4E71:
		return DecodedOpcode{
			Instruction: insNOP,
			Timing:      4,
		}

	case 0x4E72:
		return DecodedOpcode{
			Instruction: insSTOP,
			Timing:      4,
		}

	case 0x4E73:
		return DecodedOpcode{
			Instruction: insRTE,
			Timing:      20,
		}

	case 0x4E75:
		return DecodedOpcode{
			Instruction: insRTS,
			Timing:      16,
		}

	case 0x4E76:
		return DecodedOpcode{
			Instruction: insTRAPV,
			Timing:      4,
		}

	case 0x4E77:
		return DecodedOpcode{
			Instruction: insRTR,
			Timing:      20,
		}

	default:
		return illegalOpcode()
	}
}

// decodeLine4Group decodes the remaining line 4 instructions.
func decodeLine4Group(opcode uint16, mode, reg uint8) DecodedOpcode {
	if d, ok := decodeLine4Unary(opcode, mode, reg); ok {
		return d
	}

	return decodeLine4Extended(opcode, mode, reg)
}

// decodeLine4Unary decodes unary ALU operations and MOVE to/from SR/CCR in line 4.
// MOVE to CCR and MOVE to SR read any data address; the others alter a data address.
func decodeLine4Unary(opcode uint16, mode, reg uint8) (DecodedOpcode, bool) {
	op := (opcode >> 6) & 0x3F
	var d DecodedOpcode
	var ok bool
	if op <= 0x0A {
		d, ok = decodeLine4UnaryLow(op, mode, reg)
	} else {
		d, ok = decodeLine4UnaryHigh(op, mode, reg)
	}
	if !ok {
		return d, false
	}

	required := eaDataAlterable
	if d.Extra == 4 || d.Extra == 5 {
		required = eaData
	}
	return requireEA(d, mode, reg, required), true
}

// line4UnaryLow maps the low opcode bits of line 4 to NEGX, MOVE from SR, and CLR.
var line4UnaryLow = map[uint16]DecodedOpcode{
	0x00: {
		Instruction: insNEGX,
		Size:        SizeByte,
		Timing:      4,
	},
	0x01: {
		Instruction: insNEGX,
		Size:        SizeWord,
		Timing:      4,
	},
	0x02: {
		Instruction: insNEGX,
		Size:        SizeLong,
		Timing:      6,
	},
	0x03: {
		Instruction: insMOVE,
		Size:        SizeWord,
		Extra:       3,
		Timing:      6,
	},
	0x08: {
		Instruction: insCLR,
		Size:        SizeByte,
		Timing:      4,
	},
	0x09: {
		Instruction: insCLR,
		Size:        SizeWord,
		Timing:      4,
	},
	0x0A: {
		Instruction: insCLR,
		Size:        SizeLong,
		Timing:      6,
	},
}

func decodeLine4UnaryLow(op uint16, mode, reg uint8) (DecodedOpcode, bool) {
	d, ok := line4UnaryLow[op]
	if !ok {
		return DecodedOpcode{}, false
	}
	d.DstMode, d.DstReg = mode, reg
	return d, true
}

func decodeLine4UnaryHigh(op uint16, mode, reg uint8) (DecodedOpcode, bool) {
	if op >= 0x10 && op <= 0x13 {
		return decodeLine4NegAndMoveCCR(op, mode, reg), true
	}

	switch op {
	case 0x18: // NOT.B
		return DecodedOpcode{
			Instruction: insNOT,
			Size:        SizeByte,
			DstMode:     mode,
			DstReg:      reg,
			Timing:      4,
		}, true

	case 0x19: // NOT.W
		return DecodedOpcode{
			Instruction: insNOT,
			Size:        SizeWord,
			DstMode:     mode,
			DstReg:      reg,
			Timing:      4,
		}, true

	case 0x1A: // NOT.L
		return DecodedOpcode{
			Instruction: insNOT,
			Size:        SizeLong,
			DstMode:     mode,
			DstReg:      reg,
			Timing:      6,
		}, true

	case 0x1B: // MOVE to SR
		return DecodedOpcode{
			Instruction: insMOVE,
			Size:        SizeWord,
			SrcMode:     mode,
			SrcReg:      reg,
			Extra:       5,
			Timing:      12,
		}, true

	case 0x20: // NBCD
		return DecodedOpcode{
			Instruction: insNBCD,
			Size:        SizeByte,
			DstMode:     mode,
			DstReg:      reg,
			Timing:      8,
		}, true

	default:
		return DecodedOpcode{}, false
	}
}

func decodeLine4NegAndMoveCCR(op uint16, mode, reg uint8) DecodedOpcode {
	switch op {
	case 0x10: // NEG.B
		return DecodedOpcode{
			Instruction: insNEG,
			Size:        SizeByte,
			DstMode:     mode,
			DstReg:      reg,
			Timing:      4,
		}

	case 0x11: // NEG.W
		return DecodedOpcode{
			Instruction: insNEG,
			Size:        SizeWord,
			DstMode:     mode,
			DstReg:      reg,
			Timing:      4,
		}

	case 0x12: // NEG.L
		return DecodedOpcode{
			Instruction: insNEG,
			Size:        SizeLong,
			DstMode:     mode,
			DstReg:      reg,
			Timing:      6,
		}

	default: // MOVE to CCR
		return DecodedOpcode{
			Instruction: insMOVE,
			Size:        SizeWord,
			SrcMode:     mode,
			SrcReg:      reg,
			Extra:       4,
			Timing:      12,
		}
	}
}

// decodeLine4Extended decodes SWAP, PEA, EXT, MOVEM, TST, TAS, LEA, CHK.
func decodeLine4Extended(opcode uint16, mode, reg uint8) DecodedOpcode {
	if decoded, ok := decodeLine4SwapExtMovem(opcode, mode, reg); ok {
		return decoded
	}

	switch {
	case opcode&0xFF00 == 0x4A00:
		return decodeLine4TstTas(opcode, mode, reg)

	case opcode&0xF1C0 == 0x41C0:
		an := (opcode >> 9) & 7
		return requireEA(DecodedOpcode{
			Instruction: insLEA,
			Size:        SizeLong,
			SrcMode:     mode,
			SrcReg:      reg,
			DstReg:      uint8(an),
			Timing:      4,
		}, mode, reg, eaControl)

	case opcode&0xF1C0 == 0x4180:
		dn := (opcode >> 9) & 7
		return requireEA(DecodedOpcode{
			Instruction: insCHK,
			Size:        SizeWord,
			SrcMode:     mode,
			SrcReg:      reg,
			DstReg:      uint8(dn),
			Timing:      10,
		}, mode, reg, eaData)

	default:
		return illegalOpcode()
	}
}

// decodeLine4SwapExtMovem decodes SWAP, PEA, EXT, and MOVEM. MOVEM to memory
// accepts control alterable and predecrement addresses; MOVEM to registers
// accepts control and postincrement addresses.
func decodeLine4SwapExtMovem(opcode uint16, mode, reg uint8) (DecodedOpcode, bool) {
	switch {
	case opcode&0xFFF8 == 0x4840:
		return DecodedOpcode{
			Instruction: insSWAP,
			DstReg:      reg,
			Timing:      4,
		}, true

	case opcode&0xFFC0 == 0x4840:
		return requireEA(DecodedOpcode{
			Instruction: insPEA,
			Size:        SizeLong,
			DstMode:     mode,
			DstReg:      reg,
			Timing:      12,
		}, mode, reg, eaControl), true

	case opcode&0xFFF8 == 0x4880:
		return DecodedOpcode{
			Instruction: insEXT,
			Size:        SizeWord,
			DstReg:      reg,
			Timing:      4,
		}, true

	case opcode&0xFFF8 == 0x48C0:
		return DecodedOpcode{
			Instruction: insEXT,
			Size:        SizeLong,
			DstReg:      reg,
			Timing:      4,
		}, true

	case opcode&0xFF80 == 0x4880, opcode&0xFF80 == 0x4C80:
		return decodeLine4Movem(opcode, mode, reg), true

	default:
		return DecodedOpcode{}, false
	}
}

// decodeLine4Movem decodes MOVEM in both directions. MOVEM to memory accepts
// control alterable and predecrement addresses; MOVEM to registers accepts
// control and postincrement addresses.
func decodeLine4Movem(opcode uint16, mode, reg uint8) DecodedOpcode {
	if opcode&0x0400 == 0 {
		d := DecodedOpcode{
			Instruction: insMOVEM,
			Size:        movemSize(opcode),
			DstMode:     mode,
			DstReg:      reg,
			Timing:      8,
		}
		if mode == 4 {
			return d
		}
		return requireEA(d, mode, reg, eaControlAlterable)
	}

	d := DecodedOpcode{
		Instruction: insMOVEM,
		Size:        movemSize(opcode),
		SrcMode:     mode,
		SrcReg:      reg,
		Extra:       1,
		Timing:      12,
	}
	if mode == 3 {
		return d
	}
	return requireEA(d, mode, reg, eaControl)
}

// decodeLine4TstTas decodes TST and TAS instructions.
func decodeLine4TstTas(opcode uint16, mode, reg uint8) DecodedOpcode {
	size := sizeFromBits((opcode >> 6) & 3)
	if size == 0 {
		return requireEA(DecodedOpcode{
			Instruction: insTAS,
			Size:        SizeByte,
			DstMode:     mode,
			DstReg:      reg,
			Timing:      4,
		}, mode, reg, eaDataAlterable)
	}
	return requireEA(DecodedOpcode{
		Instruction: insTST,
		Size:        size,
		DstMode:     mode,
		DstReg:      reg,
		Timing:      4,
	}, mode, reg, eaDataAlterable)
}

// movemSize returns the operand size for MOVEM from the opcode bit.
func movemSize(opcode uint16) OperandSize {
	if opcode&0x0040 != 0 {
		return SizeLong
	}
	return SizeWord
}

// decodeLine5 decodes line 5: ADDQ, SUBQ, Scc, DBcc.
func decodeLine5(opcode uint16) DecodedOpcode {
	mode := uint8((opcode >> 3) & 7)
	reg := uint8(opcode & 7)
	sizeBits := (opcode >> 6) & 3

	if sizeBits == 3 {
		// Scc or DBcc
		cond := (opcode >> 8) & 0xF
		if mode == 1 {
			// DBcc Dn,displacement
			return DecodedOpcode{
				Instruction: insDBcc,
				DstReg:      reg,
				Extra:       cond,
				Timing:      10,
			}
		}
		// Scc
		return requireEA(DecodedOpcode{
			Instruction: insScc,
			Size:        SizeByte,
			DstMode:     mode,
			DstReg:      reg,
			Extra:       cond,
			Timing:      4,
		}, mode, reg, eaDataAlterable)
	}

	// ADDQ or SUBQ. A byte operation cannot alter an address register.
	data := (opcode >> 9) & 7
	if data == 0 {
		data = 8
	}
	size := sizeFromBits(sizeBits)
	if size == SizeByte && mode == 1 {
		return illegalOpcode()
	}

	ins := insADDQ
	if opcode&0x0100 != 0 {
		ins = insSUBQ
	}
	return requireEA(DecodedOpcode{
		Instruction: ins,
		Size:        size,
		DstMode:     mode,
		DstReg:      reg,
		Extra:       data,
		Timing:      4,
	}, mode, reg, eaAlterable)
}

// decodeLine6 decodes line 6: Bcc, BRA, BSR.
func decodeLine6(opcode uint16) DecodedOpcode {
	cond := (opcode >> 8) & 0xF
	disp := opcode & 0xFF

	var ins *Instruction

	switch cond {
	case 0:
		ins = insBRA
	case 1:
		ins = insBSR
	default:
		ins = insBcc
	}

	return DecodedOpcode{
		Instruction: ins,
		Extra:       cond,
		DstReg:      uint8(disp), // 8-bit displacement stored in DstReg for short branch
		Timing:      10,
	}
}

// decodeLine7 decodes line 7: MOVEQ.
func decodeLine7(opcode uint16) DecodedOpcode {
	if opcode&0x0100 != 0 {
		return illegalOpcode()
	}

	dn := (opcode >> 9) & 7
	data := opcode & 0xFF

	return DecodedOpcode{
		Instruction: insMOVEQ,
		Size:        SizeLong,
		DstMode:     0,
		DstReg:      uint8(dn),
		Extra:       data,
		Timing:      4,
	}
}

// decodeLine8 decodes line 8: OR, DIV, SBCD.
func decodeLine8(opcode uint16) DecodedOpcode {
	if opcode&0xF1F0 == 0x8100 {
		return decodeBCD(opcode, insSBCD)
	}

	opMode := (opcode >> 6) & 7
	switch opMode {
	case 3:
		return decodeMulDiv(opcode, insDIVU, 140)
	case 7:
		return decodeMulDiv(opcode, insDIVS, 158)
	default:
		return decodeDyadic(opcode, insOR, opMode, eaData)
	}
}

// decodeLine9 decodes line 9: SUB, SUBA, SUBX.
func decodeLine9(opcode uint16) DecodedOpcode {
	return decodeAddSub(opcode, insSUB, insSUBA, insSUBX)
}

// decodeLineA decodes line A: Line A emulator trap.
func decodeLineA(opcode uint16) DecodedOpcode {
	return DecodedOpcode{
		Instruction: insILLEGAL,
		Extra:       opcode,
		Timing:      illegalTiming,
	}
}

// decodeLineB decodes line B: CMP, CMPA, CMPM, EOR.
func decodeLineB(opcode uint16) DecodedOpcode {
	dn := uint8((opcode >> 9) & 7)
	mode := uint8((opcode >> 3) & 7)
	reg := uint8(opcode & 7)
	opMode := (opcode >> 6) & 7

	// CMPA
	if opMode == 3 || opMode == 7 {
		sz := SizeWord
		if opMode == 7 {
			sz = SizeLong
		}
		return requireEA(DecodedOpcode{
			Instruction: insCMPA,
			Size:        sz,
			SrcMode:     mode,
			SrcReg:      reg,
			DstReg:      dn,
			Timing:      6,
		}, mode, reg, eaAny)
	}

	size := sizeFromBits(opMode & 3)

	// CMPM
	if opMode >= 4 && mode == 1 {
		return DecodedOpcode{
			Instruction: insCMPM,
			Size:        size,
			SrcMode:     3,
			SrcReg:      reg,
			DstMode:     3,
			DstReg:      dn,
			Timing:      12,
		}
	}

	// EOR Dn,<ea>
	if opMode >= 4 {
		return requireEA(DecodedOpcode{
			Instruction: insEOR,
			Size:        size,
			SrcMode:     0,
			SrcReg:      dn,
			DstMode:     mode,
			DstReg:      reg,
			Timing:      4,
		}, mode, reg, eaDataAlterable)
	}

	// CMP <ea>,Dn. A byte source cannot be an address register.
	return requireEA(DecodedOpcode{
		Instruction: insCMP,
		Size:        size,
		SrcMode:     mode,
		SrcReg:      reg,
		DstMode:     0,
		DstReg:      dn,
		Timing:      4,
	}, mode, reg, byteSourceClass(size))
}

// decodeLineC decodes line C: AND, MUL, ABCD, EXG.
func decodeLineC(opcode uint16) DecodedOpcode {
	if opcode&0xF1F0 == 0xC100 {
		return decodeBCD(opcode, insABCD)
	}

	opMode := (opcode >> 6) & 7
	switch opMode {
	case 3:
		return decodeMulDiv(opcode, insMULU, 70)
	case 7:
		return decodeMulDiv(opcode, insMULS, 70)
	}

	// EXG variants.
	if d, ok := decodeLineCExg(opcode, opMode); ok {
		return d
	}

	return decodeDyadic(opcode, insAND, opMode, eaData)
}

// decodeLineCExg decodes EXG instruction variants within line C.
func decodeLineCExg(opcode, opMode uint16) (DecodedOpcode, bool) {
	dn := uint8((opcode >> 9) & 7)
	mode := (opcode >> 3) & 7
	reg := uint8(opcode & 7)

	d := DecodedOpcode{
		Instruction: insEXG,
		SrcReg:      dn,
		DstReg:      reg,
		Timing:      6,
	}
	switch {
	case opMode == 5 && mode == 0: // EXG Dn,Dn
		d.Extra = 0
	case opMode == 5 && mode == 1: // EXG An,An
		d.Extra = 1
	case opMode == 6 && mode == 1: // EXG Dn,An
		d.Extra = 2
	default:
		return DecodedOpcode{}, false
	}
	return d, true
}

// decodeLineD decodes line D: ADD, ADDA, ADDX.
func decodeLineD(opcode uint16) DecodedOpcode {
	return decodeAddSub(opcode, insADD, insADDA, insADDX)
}

// decodeAddSub decodes ADD/SUB family instructions (lines 9, D).
func decodeAddSub(opcode uint16, insBase, insAddr, insExtended *Instruction) DecodedOpcode {
	dn := uint8((opcode >> 9) & 7)
	mode := uint8((opcode >> 3) & 7)
	reg := uint8(opcode & 7)
	opMode := (opcode >> 6) & 7

	// ADDA/SUBA
	if opMode == 3 || opMode == 7 {
		sz := SizeWord
		if opMode == 7 {
			sz = SizeLong
		}
		return requireEA(DecodedOpcode{
			Instruction: insAddr,
			Size:        sz,
			SrcMode:     mode,
			SrcReg:      reg,
			DstReg:      dn,
			Timing:      8,
		}, mode, reg, eaAny)
	}

	// ADDX/SUBX
	if opMode >= 4 && mode <= 1 {
		return DecodedOpcode{
			Instruction: insExtended,
			Size:        sizeFromBits(opMode & 3),
			SrcReg:      reg,
			DstReg:      dn,
			Extra:       uint16(mode),
			Timing:      4,
		}
	}

	size := sizeFromBits(opMode & 3)
	return decodeDyadic(opcode, insBase, opMode, byteSourceClass(size))
}

// decodeDyadic decodes the <ea>,Dn and Dn,<ea> forms of ADD, SUB, AND, and OR.
// The source must have the given class; the memory destination must be alterable.
func decodeDyadic(opcode uint16, ins *Instruction, opMode uint16, source eaClass) DecodedOpcode {
	dn := uint8((opcode >> 9) & 7)
	mode := uint8((opcode >> 3) & 7)
	reg := uint8(opcode & 7)

	d := DecodedOpcode{
		Instruction: ins,
		Size:        sizeFromBits(opMode & 3),
		Timing:      4,
	}
	if opMode < 3 {
		d.SrcMode = mode
		d.SrcReg = reg
		d.DstMode = 0
		d.DstReg = dn
		return requireEA(d, mode, reg, source)
	}

	d.SrcMode = 0
	d.SrcReg = dn
	d.DstMode = mode
	d.DstReg = reg
	return requireEA(d, mode, reg, eaMemoryAlterable)
}

// byteSourceClass returns the source class of ADD, SUB, and CMP. Only their
// byte operations exclude address registers.
func byteSourceClass(size OperandSize) eaClass {
	if size == SizeByte {
		return eaData
	}
	return eaAny
}

// decodeBCD decodes the register and predecrement forms of ABCD and SBCD.
func decodeBCD(opcode uint16, ins *Instruction) DecodedOpcode {
	return DecodedOpcode{
		Instruction: ins,
		Size:        SizeByte,
		SrcReg:      uint8(opcode & 7),
		DstReg:      uint8((opcode >> 9) & 7),
		Extra:       opcode & 0x8, // RM bit
		Timing:      6,
	}
}

// decodeMulDiv decodes MULU, MULS, DIVU, and DIVS, which read a data address.
func decodeMulDiv(opcode uint16, ins *Instruction, timing uint16) DecodedOpcode {
	mode := uint8((opcode >> 3) & 7)
	reg := uint8(opcode & 7)
	return requireEA(DecodedOpcode{
		Instruction: ins,
		Size:        SizeWord,
		SrcMode:     mode,
		SrcReg:      reg,
		DstReg:      uint8((opcode >> 9) & 7),
		Timing:      timing,
	}, mode, reg, eaData)
}

// shiftRotateInstructions maps (type << 1 | direction) to instruction.
// Type: 0=AS, 1=LS, 2=ROX, 3=RO. Direction: 0=right, 1=left.
var shiftRotateInstructions = [8]*Instruction{
	insASR, insASL, insLSR, insLSL, insROXR, insROXL, insROR, insROL,
}

// decodeLineE decodes line E: Shift/Rotate instructions.
func decodeLineE(opcode uint16) DecodedOpcode {
	reg := opcode & 7

	// Memory shift/rotate (size = word, count = 1).
	if (opcode>>6)&3 == 3 {
		return decodeLineEMemory(opcode)
	}

	// Register shift/rotate.
	size := sizeFromBits((opcode >> 6) & 3)
	count := (opcode >> 9) & 7
	dr := (opcode >> 8) & 1
	ir := (opcode >> 5) & 1
	typ := (opcode >> 3) & 3

	ins := shiftRotateInstructions[typ<<1|dr]

	extra := count
	if ir != 0 {
		extra |= 0x20 // Flag to indicate count is in register
	}

	return DecodedOpcode{
		Instruction: ins,
		Size:        size,
		DstReg:      uint8(reg),
		Extra:       extra,
		Timing:      6,
	}
}

// decodeLineEMemory decodes memory shift/rotate instructions. Bit 11 must be
// zero on the 68000; later processors use it for bit field operations.
func decodeLineEMemory(opcode uint16) DecodedOpcode {
	if opcode&0x0800 != 0 {
		return illegalOpcode()
	}

	mode := uint8((opcode >> 3) & 7)
	reg := uint8(opcode & 7)
	typ := (opcode >> 9) & 3
	dr := (opcode >> 8) & 1

	ins := shiftRotateInstructions[typ<<1|dr]

	return requireEA(DecodedOpcode{
		Instruction: ins,
		Size:        SizeWord,
		DstMode:     mode,
		DstReg:      reg,
		Extra:       0x40, // Flag to indicate memory operation
		Timing:      8,
	}, mode, reg, eaMemoryAlterable)
}

// decodeLineF decodes line F: Line F emulator trap.
func decodeLineF(opcode uint16) DecodedOpcode {
	return DecodedOpcode{
		Instruction: insILLEGAL,
		Extra:       opcode,
		Timing:      illegalTiming,
	}
}
