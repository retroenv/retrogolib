package m68000

import "fmt"

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

// lineDecoders maps the top 4 bits of an opcode word to a line decoder function.
var lineDecoders = [16]func(opcode uint16) (DecodedOpcode, error){
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

// decodeOpcode decodes a 16-bit opcode word into a DecodedOpcode.
func decodeOpcode(opcode uint16) (DecodedOpcode, error) {
	line := (opcode >> 12) & 0xF
	return lineDecoders[line](opcode)
}

// decodeLine0 decodes line 0: immediate operations and bit operations.
func decodeLine0(opcode uint16) (DecodedOpcode, error) {
	// Bits 11-8 determine the specific operation.
	if opcode&0x0100 != 0 {
		// Bit operations with register (BTST/BCHG/BCLR/BSET Dn,<ea>).
		return decodeLine0BitReg(opcode)
	}

	// Check for MOVEP.
	if opcode&0x0138 == 0x0108 {
		return decodeLine0Movep(opcode)
	}

	// Immediate operations.
	return decodeLine0Immediate(opcode)
}

// decodeLine0BitReg decodes bit operations with register source.
func decodeLine0BitReg(opcode uint16) (DecodedOpcode, error) {
	dn := (opcode >> 9) & 7
	mode := (opcode >> 3) & 7
	reg := opcode & 7

	bitOp := (opcode >> 6) & 3

	var ins *Instruction

	switch bitOp {
	case 0:
		ins = insBTST
	case 1:
		ins = insBCHG
	case 2:
		ins = insBCLR
	case 3:
		ins = insBSET
	}

	return DecodedOpcode{
		Instruction: ins,
		Size:        SizeByte,
		SrcMode:     0, // Data register
		SrcReg:      uint8(dn),
		DstMode:     uint8(mode),
		DstReg:      uint8(reg),
		Timing:      8,
	}, nil
}

// decodeLine0Movep decodes MOVEP instruction.
func decodeLine0Movep(opcode uint16) (DecodedOpcode, error) {
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

	if dir == 0 {
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

	return d, nil
}

// decodeLine0Immediate decodes immediate operations (ORI, ANDI, SUBI, ADDI, EORI, CMPI)
// and immediate bit operations (BTST/BCHG/BCLR/BSET #imm,<ea>).
func decodeLine0Immediate(opcode uint16) (DecodedOpcode, error) {
	op := (opcode >> 9) & 7
	size := sizeFromBits((opcode >> 6) & 3)
	mode := (opcode >> 3) & 7
	reg := opcode & 7

	var ins *Instruction
	var timing uint16 = 8

	switch op {
	case 0: // ORI
		ins = insORI
		if mode == 7 && reg == 4 {
			ins = insORI // ORI to CCR/SR
		}
	case 1: // ANDI
		ins = insANDI
	case 2: // SUBI
		ins = insSUBI
	case 3: // ADDI
		ins = insADDI
	case 4: // Bit operations with immediate
		return decodeLine0BitImm(opcode)
	case 5: // EORI
		ins = insEORI
	case 6: // CMPI
		ins = insCMPI
	default:
		return DecodedOpcode{}, fmt.Errorf("%w: 0x%04X", ErrUnsupportedOpcode, opcode)
	}

	if size == SizeLong {
		timing = 16
	}

	return DecodedOpcode{
		Instruction: ins,
		Size:        size,
		DstMode:     uint8(mode),
		DstReg:      uint8(reg),
		Timing:      timing,
	}, nil
}

// decodeLine0BitImm decodes bit operations with immediate source.
func decodeLine0BitImm(opcode uint16) (DecodedOpcode, error) {
	bitOp := (opcode >> 6) & 3
	mode := (opcode >> 3) & 7
	reg := opcode & 7

	var ins *Instruction

	switch bitOp {
	case 0:
		ins = insBTST
	case 1:
		ins = insBCHG
	case 2:
		ins = insBCLR
	case 3:
		ins = insBSET
	}

	return DecodedOpcode{
		Instruction: ins,
		Size:        SizeByte,
		SrcMode:     7,
		SrcReg:      4, // Immediate
		DstMode:     uint8(mode),
		DstReg:      uint8(reg),
		Timing:      12,
	}, nil
}

// decodeLine1 decodes line 1: MOVE.B.
func decodeLine1(opcode uint16) (DecodedOpcode, error) {
	return decodeMOVE(opcode, SizeByte)
}

// decodeLine2 decodes line 2: MOVE.L and MOVEA.L.
func decodeLine2(opcode uint16) (DecodedOpcode, error) {
	dstMode := (opcode >> 6) & 7
	if dstMode == 1 {
		return decodeMOVEA(opcode, SizeLong)
	}
	return decodeMOVE(opcode, SizeLong)
}

// decodeLine3 decodes line 3: MOVE.W and MOVEA.W.
func decodeLine3(opcode uint16) (DecodedOpcode, error) {
	dstMode := (opcode >> 6) & 7
	if dstMode == 1 {
		return decodeMOVEA(opcode, SizeWord)
	}
	return decodeMOVE(opcode, SizeWord)
}

// decodeMOVE decodes a MOVE instruction. Note: MOVE uses a special encoding
// where destination is encoded as register:mode (reversed from source).
func decodeMOVE(opcode uint16, size OperandSize) (DecodedOpcode, error) {
	srcMode := (opcode >> 3) & 7
	srcReg := opcode & 7
	dstReg := (opcode >> 9) & 7
	dstMode := (opcode >> 6) & 7

	return DecodedOpcode{
		Instruction: insMOVE,
		Size:        size,
		SrcMode:     uint8(srcMode),
		SrcReg:      uint8(srcReg),
		DstMode:     uint8(dstMode),
		DstReg:      uint8(dstReg),
		Timing:      4,
	}, nil
}

// decodeMOVEA decodes a MOVEA instruction.
func decodeMOVEA(opcode uint16, size OperandSize) (DecodedOpcode, error) {
	srcMode := (opcode >> 3) & 7
	srcReg := opcode & 7
	dstReg := (opcode >> 9) & 7

	return DecodedOpcode{
		Instruction: insMOVEA,
		Size:        size,
		SrcMode:     uint8(srcMode),
		SrcReg:      uint8(srcReg),
		DstMode:     1, // Address register direct
		DstReg:      uint8(dstReg),
		Timing:      4,
	}, nil
}

// decodeLine4 decodes line 4: Miscellaneous instructions.
func decodeLine4(opcode uint16) (DecodedOpcode, error) {
	mode := (opcode >> 3) & 7
	reg := opcode & 7

	// Check for specific encodings.
	switch {
	case opcode&0xFFF8 == 0x4E70:
		return decodeLine4Special(opcode)
	case opcode == 0x4AFC:
		return DecodedOpcode{Instruction: insILLEGAL, Timing: 34}, nil
	case opcode&0xFFF0 == 0x4E40:
		return DecodedOpcode{Instruction: insTRAP, Extra: opcode & 0x0F, Timing: 34}, nil
	case opcode&0xFFF8 == 0x4E50:
		return DecodedOpcode{Instruction: insLINK, DstReg: uint8(reg), Timing: 16}, nil
	case opcode&0xFFF8 == 0x4E58:
		return DecodedOpcode{Instruction: insUNLK, DstReg: uint8(reg), Timing: 12}, nil
	case opcode&0xFFF8 == 0x4E60:
		return decodeLine4MoveUSP(reg, true)
	case opcode&0xFFF8 == 0x4E68:
		return decodeLine4MoveUSP(reg, false)
	case opcode&0xFFC0 == 0x4E80:
		return DecodedOpcode{Instruction: insJSR, DstMode: uint8(mode), DstReg: uint8(reg), Timing: 16}, nil
	case opcode&0xFFC0 == 0x4EC0:
		return DecodedOpcode{Instruction: insJMP, DstMode: uint8(mode), DstReg: uint8(reg), Timing: 8}, nil
	}

	return decodeLine4Group(opcode)
}

// decodeLine4MoveUSP decodes MOVE An,USP and MOVE USP,An.
func decodeLine4MoveUSP(reg uint16, toUSP bool) (DecodedOpcode, error) {
	if toUSP {
		return DecodedOpcode{
			Instruction: insMOVE, Size: SizeLong,
			SrcMode: 1, SrcReg: uint8(reg), DstMode: 7, DstReg: 5,
			Extra: 1, Timing: 4,
		}, nil
	}
	return DecodedOpcode{
		Instruction: insMOVE, Size: SizeLong,
		SrcMode: 7, SrcReg: 5, DstMode: 1, DstReg: uint8(reg),
		Extra: 2, Timing: 4,
	}, nil
}

// decodeLine4Special decodes the special instructions at 0x4E7x.
func decodeLine4Special(opcode uint16) (DecodedOpcode, error) {
	switch opcode {
	case 0x4E70:
		return DecodedOpcode{Instruction: insRESET, Timing: 132}, nil
	case 0x4E71:
		return DecodedOpcode{Instruction: insNOP, Timing: 4}, nil
	case 0x4E72:
		return DecodedOpcode{Instruction: insSTOP, Timing: 4}, nil
	case 0x4E73:
		return DecodedOpcode{Instruction: insRTE, Timing: 20}, nil
	case 0x4E75:
		return DecodedOpcode{Instruction: insRTS, Timing: 16}, nil
	case 0x4E76:
		return DecodedOpcode{Instruction: insTRAPV, Timing: 4}, nil
	case 0x4E77:
		return DecodedOpcode{Instruction: insRTR, Timing: 20}, nil
	default:
		return DecodedOpcode{}, fmt.Errorf("%w: 0x%04X", ErrUnsupportedOpcode, opcode)
	}
}

// decodeLine4Group decodes the remaining line 4 instructions.
func decodeLine4Group(opcode uint16) (DecodedOpcode, error) {
	mode := (opcode >> 3) & 7
	reg := opcode & 7

	if d, ok := decodeLine4Unary(opcode, mode, reg); ok {
		return d, nil
	}

	return decodeLine4Extended(opcode, mode, reg)
}

// decodeLine4Unary decodes unary ALU operations and MOVE to/from SR/CCR in line 4.
func decodeLine4Unary(opcode, mode, reg uint16) (DecodedOpcode, bool) {
	op := (opcode >> 6) & 0x3F
	if op <= 0x0A {
		return decodeLine4UnaryLow(op, mode, reg)
	}
	return decodeLine4UnaryHigh(op, mode, reg)
}

func decodeLine4UnaryLow(op, mode, reg uint16) (DecodedOpcode, bool) {
	switch op {
	case 0x00: // NEGX.B
		return DecodedOpcode{Instruction: insNEGX, Size: SizeByte, DstMode: uint8(mode), DstReg: uint8(reg), Timing: 4}, true
	case 0x01: // NEGX.W
		return DecodedOpcode{Instruction: insNEGX, Size: SizeWord, DstMode: uint8(mode), DstReg: uint8(reg), Timing: 4}, true
	case 0x02: // NEGX.L
		return DecodedOpcode{Instruction: insNEGX, Size: SizeLong, DstMode: uint8(mode), DstReg: uint8(reg), Timing: 6}, true
	case 0x03: // MOVE from SR
		return DecodedOpcode{Instruction: insMOVE, Size: SizeWord, DstMode: uint8(mode), DstReg: uint8(reg), Extra: 3, Timing: 6}, true
	case 0x08: // CLR.B
		return DecodedOpcode{Instruction: insCLR, Size: SizeByte, DstMode: uint8(mode), DstReg: uint8(reg), Timing: 4}, true
	case 0x09: // CLR.W
		return DecodedOpcode{Instruction: insCLR, Size: SizeWord, DstMode: uint8(mode), DstReg: uint8(reg), Timing: 4}, true
	case 0x0A: // CLR.L
		return DecodedOpcode{Instruction: insCLR, Size: SizeLong, DstMode: uint8(mode), DstReg: uint8(reg), Timing: 6}, true
	default:
		return DecodedOpcode{}, false
	}
}

func decodeLine4UnaryHigh(op, mode, reg uint16) (DecodedOpcode, bool) {
	switch op {
	case 0x10: // NEG.B
		return DecodedOpcode{Instruction: insNEG, Size: SizeByte, DstMode: uint8(mode), DstReg: uint8(reg), Timing: 4}, true
	case 0x11: // NEG.W
		return DecodedOpcode{Instruction: insNEG, Size: SizeWord, DstMode: uint8(mode), DstReg: uint8(reg), Timing: 4}, true
	case 0x12: // NEG.L
		return DecodedOpcode{Instruction: insNEG, Size: SizeLong, DstMode: uint8(mode), DstReg: uint8(reg), Timing: 6}, true
	case 0x13: // MOVE to CCR
		return DecodedOpcode{Instruction: insMOVE, Size: SizeWord, SrcMode: uint8(mode), SrcReg: uint8(reg), Extra: 4, Timing: 12}, true
	case 0x18: // NOT.B
		return DecodedOpcode{Instruction: insNOT, Size: SizeByte, DstMode: uint8(mode), DstReg: uint8(reg), Timing: 4}, true
	case 0x19: // NOT.W
		return DecodedOpcode{Instruction: insNOT, Size: SizeWord, DstMode: uint8(mode), DstReg: uint8(reg), Timing: 4}, true
	case 0x1A: // NOT.L
		return DecodedOpcode{Instruction: insNOT, Size: SizeLong, DstMode: uint8(mode), DstReg: uint8(reg), Timing: 6}, true
	case 0x1B: // MOVE to SR
		return DecodedOpcode{Instruction: insMOVE, Size: SizeWord, SrcMode: uint8(mode), SrcReg: uint8(reg), Extra: 5, Timing: 12}, true
	case 0x20: // NBCD
		return DecodedOpcode{Instruction: insNBCD, Size: SizeByte, DstMode: uint8(mode), DstReg: uint8(reg), Timing: 8}, true
	default:
		return DecodedOpcode{}, false
	}
}

// decodeLine4Extended decodes SWAP, PEA, EXT, MOVEM, TST, TAS, Scc, LEA, CHK.
func decodeLine4Extended(opcode, mode, reg uint16) (DecodedOpcode, error) {
	switch {
	case opcode&0xFFF8 == 0x4840:
		return DecodedOpcode{Instruction: insSWAP, DstReg: uint8(reg), Timing: 4}, nil
	case opcode&0xFFC0 == 0x4840:
		return DecodedOpcode{Instruction: insPEA, Size: SizeLong, DstMode: uint8(mode), DstReg: uint8(reg), Timing: 12}, nil
	case opcode&0xFFF8 == 0x4880:
		return DecodedOpcode{Instruction: insEXT, Size: SizeWord, DstReg: uint8(reg), Timing: 4}, nil
	case opcode&0xFFF8 == 0x48C0:
		return DecodedOpcode{Instruction: insEXT, Size: SizeLong, DstReg: uint8(reg), Timing: 4}, nil
	case opcode&0xFB80 == 0x4880:
		sz := movemSize(opcode)
		return DecodedOpcode{Instruction: insMOVEM, Size: sz, DstMode: uint8(mode), DstReg: uint8(reg), Extra: 0, Timing: 8}, nil
	case opcode&0xFB80 == 0x4C80:
		sz := movemSize(opcode)
		return DecodedOpcode{Instruction: insMOVEM, Size: sz, SrcMode: uint8(mode), SrcReg: uint8(reg), Extra: 1, Timing: 12}, nil
	case opcode&0xFF00 == 0x4A00:
		return decodeLine4TstTas(opcode, mode, reg)
	case opcode&0xF0C0 == 0x50C0 && mode != 1:
		cond := (opcode >> 8) & 0xF
		return DecodedOpcode{Instruction: insScc, Size: SizeByte, DstMode: uint8(mode), DstReg: uint8(reg), Extra: cond, Timing: 4}, nil
	case opcode&0xF1C0 == 0x41C0:
		an := (opcode >> 9) & 7
		return DecodedOpcode{Instruction: insLEA, Size: SizeLong, SrcMode: uint8(mode), SrcReg: uint8(reg), DstReg: uint8(an), Timing: 4}, nil
	case opcode&0xF1C0 == 0x4180:
		dn := (opcode >> 9) & 7
		return DecodedOpcode{Instruction: insCHK, Size: SizeWord, SrcMode: uint8(mode), SrcReg: uint8(reg), DstReg: uint8(dn), Timing: 10}, nil
	default:
		return DecodedOpcode{}, fmt.Errorf("%w: 0x%04X", ErrUnsupportedOpcode, opcode)
	}
}

// decodeLine4TstTas decodes TST and TAS instructions.
func decodeLine4TstTas(opcode, mode, reg uint16) (DecodedOpcode, error) {
	size := sizeFromBits((opcode >> 6) & 3)
	if size == 0 {
		return DecodedOpcode{Instruction: insTAS, Size: SizeByte, DstMode: uint8(mode), DstReg: uint8(reg), Timing: 4}, nil
	}
	return DecodedOpcode{Instruction: insTST, Size: size, DstMode: uint8(mode), DstReg: uint8(reg), Timing: 4}, nil
}

// movemSize returns the operand size for MOVEM from the opcode bit.
func movemSize(opcode uint16) OperandSize {
	if opcode&0x0040 != 0 {
		return SizeLong
	}
	return SizeWord
}

// decodeLine5 decodes line 5: ADDQ, SUBQ, Scc, DBcc.
func decodeLine5(opcode uint16) (DecodedOpcode, error) {
	mode := (opcode >> 3) & 7
	reg := opcode & 7
	sizeBits := (opcode >> 6) & 3

	if sizeBits == 3 {
		// Scc or DBcc
		if mode == 1 {
			// DBcc Dn,displacement
			cond := (opcode >> 8) & 0xF
			return DecodedOpcode{Instruction: insDBcc, DstReg: uint8(reg), Extra: cond, Timing: 10}, nil
		}
		// Scc
		cond := (opcode >> 8) & 0xF
		return DecodedOpcode{Instruction: insScc, Size: SizeByte, DstMode: uint8(mode), DstReg: uint8(reg), Extra: cond, Timing: 4}, nil
	}

	// ADDQ or SUBQ
	data := (opcode >> 9) & 7
	if data == 0 {
		data = 8
	}
	size := sizeFromBits(sizeBits)

	if opcode&0x0100 == 0 {
		return DecodedOpcode{Instruction: insADDQ, Size: size, DstMode: uint8(mode), DstReg: uint8(reg), Extra: data, Timing: 4}, nil
	}
	return DecodedOpcode{Instruction: insSUBQ, Size: size, DstMode: uint8(mode), DstReg: uint8(reg), Extra: data, Timing: 4}, nil
}

// decodeLine6 decodes line 6: Bcc, BRA, BSR.
func decodeLine6(opcode uint16) (DecodedOpcode, error) {
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
	}, nil
}

// decodeLine7 decodes line 7: MOVEQ.
func decodeLine7(opcode uint16) (DecodedOpcode, error) {
	if opcode&0x0100 != 0 {
		return DecodedOpcode{}, fmt.Errorf("%w: 0x%04X", ErrUnsupportedOpcode, opcode)
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
	}, nil
}

// decodeLine8 decodes line 8: OR, DIV, SBCD.
func decodeLine8(opcode uint16) (DecodedOpcode, error) {
	dn := (opcode >> 9) & 7
	mode := (opcode >> 3) & 7
	reg := opcode & 7
	opMode := (opcode >> 6) & 7

	// SBCD
	if opMode == 4 {
		return DecodedOpcode{
			Instruction: insSBCD,
			Size:        SizeByte,
			SrcMode:     uint8(mode),
			SrcReg:      uint8(reg),
			DstReg:      uint8(dn),
			Extra:       opcode & 0x8, // RM bit
			Timing:      6,
		}, nil
	}

	// DIVU
	if opMode == 3 {
		return DecodedOpcode{
			Instruction: insDIVU,
			Size:        SizeWord,
			SrcMode:     uint8(mode),
			SrcReg:      uint8(reg),
			DstReg:      uint8(dn),
			Timing:      140,
		}, nil
	}

	// DIVS
	if opMode == 7 {
		return DecodedOpcode{
			Instruction: insDIVS,
			Size:        SizeWord,
			SrcMode:     uint8(mode),
			SrcReg:      uint8(reg),
			DstReg:      uint8(dn),
			Timing:      158,
		}, nil
	}

	// OR
	size := sizeFromBits(opMode & 3)
	d := DecodedOpcode{
		Instruction: insOR,
		Size:        size,
		Timing:      4,
	}

	if opMode < 3 {
		// OR <ea>,Dn
		d.SrcMode = uint8(mode)
		d.SrcReg = uint8(reg)
		d.DstMode = 0
		d.DstReg = uint8(dn)
	} else {
		// OR Dn,<ea>
		d.SrcMode = 0
		d.SrcReg = uint8(dn)
		d.DstMode = uint8(mode)
		d.DstReg = uint8(reg)
	}

	return d, nil
}

// decodeLine9 decodes line 9: SUB, SUBA, SUBX.
func decodeLine9(opcode uint16) (DecodedOpcode, error) {
	return decodeAddSub(opcode, insSUB, insSUBA, insSUBX)
}

// decodeLineA decodes line A: Line A emulator trap.
func decodeLineA(opcode uint16) (DecodedOpcode, error) {
	return DecodedOpcode{
		Instruction: insILLEGAL,
		Extra:       opcode & 0x0FFF,
		Timing:      34,
	}, nil
}

// decodeLineB decodes line B: CMP, CMPA, CMPM, EOR.
func decodeLineB(opcode uint16) (DecodedOpcode, error) {
	dn := (opcode >> 9) & 7
	mode := (opcode >> 3) & 7
	reg := opcode & 7
	opMode := (opcode >> 6) & 7

	// CMPA
	if opMode == 3 || opMode == 7 {
		sz := SizeWord
		if opMode == 7 {
			sz = SizeLong
		}
		return DecodedOpcode{
			Instruction: insCMPA,
			Size:        sz,
			SrcMode:     uint8(mode),
			SrcReg:      uint8(reg),
			DstReg:      uint8(dn),
			Timing:      6,
		}, nil
	}

	// CMPM
	if opMode >= 4 && mode == 1 {
		size := sizeFromBits(opMode & 3)
		return DecodedOpcode{
			Instruction: insCMPM,
			Size:        size,
			SrcMode:     3,
			SrcReg:      uint8(reg),
			DstMode:     3,
			DstReg:      uint8(dn),
			Timing:      12,
		}, nil
	}

	// EOR Dn,<ea>
	if opMode >= 4 {
		size := sizeFromBits(opMode & 3)
		return DecodedOpcode{
			Instruction: insEOR,
			Size:        size,
			SrcMode:     0,
			SrcReg:      uint8(dn),
			DstMode:     uint8(mode),
			DstReg:      uint8(reg),
			Timing:      4,
		}, nil
	}

	// CMP <ea>,Dn
	size := sizeFromBits(opMode & 3)
	return DecodedOpcode{
		Instruction: insCMP,
		Size:        size,
		SrcMode:     uint8(mode),
		SrcReg:      uint8(reg),
		DstMode:     0,
		DstReg:      uint8(dn),
		Timing:      4,
	}, nil
}

// decodeLineC decodes line C: AND, MUL, ABCD, EXG.
func decodeLineC(opcode uint16) (DecodedOpcode, error) {
	dn := (opcode >> 9) & 7
	mode := (opcode >> 3) & 7
	reg := opcode & 7
	opMode := (opcode >> 6) & 7

	switch opMode {
	case 4: // ABCD
		return DecodedOpcode{Instruction: insABCD, Size: SizeByte, SrcReg: uint8(reg), DstReg: uint8(dn), Extra: opcode & 0x8, Timing: 6}, nil
	case 3: // MULU
		return DecodedOpcode{Instruction: insMULU, Size: SizeWord, SrcMode: uint8(mode), SrcReg: uint8(reg), DstReg: uint8(dn), Timing: 70}, nil
	case 7: // MULS
		return DecodedOpcode{Instruction: insMULS, Size: SizeWord, SrcMode: uint8(mode), SrcReg: uint8(reg), DstReg: uint8(dn), Timing: 70}, nil
	}

	// EXG variants.
	if d, ok := decodeLineCExg(opMode, mode, dn, reg); ok {
		return d, nil
	}

	// AND
	size := sizeFromBits(opMode & 3)
	d := DecodedOpcode{Instruction: insAND, Size: size, Timing: 4}
	if opMode < 3 {
		d.SrcMode = uint8(mode)
		d.SrcReg = uint8(reg)
		d.DstMode = 0
		d.DstReg = uint8(dn)
	} else {
		d.SrcMode = 0
		d.SrcReg = uint8(dn)
		d.DstMode = uint8(mode)
		d.DstReg = uint8(reg)
	}
	return d, nil
}

// decodeLineCExg decodes EXG instruction variants within line C.
func decodeLineCExg(opMode, mode, dn, reg uint16) (DecodedOpcode, bool) {
	switch {
	case opMode == 5 && mode == 0: // EXG Dn,Dn
		return DecodedOpcode{Instruction: insEXG, SrcReg: uint8(dn), DstReg: uint8(reg), Extra: 0, Timing: 6}, true
	case opMode == 5 && mode == 1: // EXG An,An
		return DecodedOpcode{Instruction: insEXG, SrcReg: uint8(dn), DstReg: uint8(reg), Extra: 1, Timing: 6}, true
	case opMode == 6 && mode == 1: // EXG Dn,An
		return DecodedOpcode{Instruction: insEXG, SrcReg: uint8(dn), DstReg: uint8(reg), Extra: 2, Timing: 6}, true
	default:
		return DecodedOpcode{}, false
	}
}

// decodeLineD decodes line D: ADD, ADDA, ADDX.
func decodeLineD(opcode uint16) (DecodedOpcode, error) {
	return decodeAddSub(opcode, insADD, insADDA, insADDX)
}

// decodeAddSub decodes ADD/SUB family instructions (lines 9, D).
func decodeAddSub(opcode uint16, insBase, insAddr, insExtended *Instruction) (DecodedOpcode, error) {
	dn := (opcode >> 9) & 7
	mode := (opcode >> 3) & 7
	reg := opcode & 7
	opMode := (opcode >> 6) & 7

	// ADDA/SUBA
	if opMode == 3 || opMode == 7 {
		sz := SizeWord
		if opMode == 7 {
			sz = SizeLong
		}
		return DecodedOpcode{
			Instruction: insAddr,
			Size:        sz,
			SrcMode:     uint8(mode),
			SrcReg:      uint8(reg),
			DstReg:      uint8(dn),
			Timing:      8,
		}, nil
	}

	// ADDX/SUBX
	if (opMode == 4 || opMode == 5 || opMode == 6) && (mode == 0 || mode == 1) {
		size := sizeFromBits(opMode & 3)
		rmBit := uint16(0)
		if mode == 1 {
			rmBit = 1
		}
		return DecodedOpcode{
			Instruction: insExtended,
			Size:        size,
			SrcReg:      uint8(reg),
			DstReg:      uint8(dn),
			Extra:       rmBit,
			Timing:      4,
		}, nil
	}

	// ADD/SUB
	size := sizeFromBits(opMode & 3)
	d := DecodedOpcode{
		Instruction: insBase,
		Size:        size,
		Timing:      4,
	}

	if opMode < 3 {
		d.SrcMode = uint8(mode)
		d.SrcReg = uint8(reg)
		d.DstMode = 0
		d.DstReg = uint8(dn)
	} else {
		d.SrcMode = 0
		d.SrcReg = uint8(dn)
		d.DstMode = uint8(mode)
		d.DstReg = uint8(reg)
	}

	return d, nil
}

// shiftRotateInstructions maps (type << 1 | direction) to instruction.
// Type: 0=AS, 1=LS, 2=ROX, 3=RO. Direction: 0=right, 1=left.
var shiftRotateInstructions = [8]*Instruction{
	insASR, insASL, insLSR, insLSL, insROXR, insROXL, insROR, insROL,
}

// decodeLineE decodes line E: Shift/Rotate instructions.
func decodeLineE(opcode uint16) (DecodedOpcode, error) {
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
	}, nil
}

// decodeLineEMemory decodes memory shift/rotate instructions.
func decodeLineEMemory(opcode uint16) (DecodedOpcode, error) {
	mode := (opcode >> 3) & 7
	reg := opcode & 7
	typ := (opcode >> 9) & 3
	dr := (opcode >> 8) & 1

	ins := shiftRotateInstructions[typ<<1|dr]

	return DecodedOpcode{
		Instruction: ins,
		Size:        SizeWord,
		DstMode:     uint8(mode),
		DstReg:      uint8(reg),
		Extra:       0x40, // Flag to indicate memory operation
		Timing:      8,
	}, nil
}

// decodeLineF decodes line F: Line F emulator trap.
func decodeLineF(opcode uint16) (DecodedOpcode, error) {
	return DecodedOpcode{
		Instruction: insILLEGAL,
		Extra:       opcode & 0x0FFF,
		Timing:      34,
	}, nil
}
