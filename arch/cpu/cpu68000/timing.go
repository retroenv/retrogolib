package cpu68000

// instructionCycles includes operand transfers and extension words, using the
// MC68000 user's manual section 8. MUL/DIV and MOVEM add operand-dependent work
// in their handlers after fetching the source or register mask exactly once.
func (cpu *CPU) instructionCycles(op DecodedOpcode) uint64 {
	src := eaCycles(op.SrcMode, op.SrcReg, op.Size)

	switch op.Instruction {
	case insMOVE:
		return moveCycles(op)
	case insMOVEA:
		return 4 + src
	case insADD, insSUB, insAND, insOR, insEOR, insADDA, insSUBA, insCMP, insCMPA:
		return aluCycles(op)
	case insADDI, insSUBI, insANDI, insORI, insEORI, insCMPI:
		return immediateCycles(op)
	case insADDQ, insSUBQ, insCLR, insNEG, insNEGX, insNOT, insTST, insNBCD, insTAS:
		return unaryCycles(op)
	case insASL, insASR, insLSL, insLSR, insROL, insROR, insROXL, insROXR:
		if op.Extra&0x40 != 0 {
			return 8 + eaCycles(op.DstMode, op.DstReg, SizeWord)
		}
		return sizeCycles(op.Size, 6, 8) + 2*uint64(cpu.shiftCount(op))
	case insMULU, insMULS, insDIVU, insDIVS:
		return src
	case insCHK:
		return 10 + src
	case insJMP, insJSR, insLEA, insPEA, insMOVEM, insMOVEP:
		return controlCycles(op)
	case insBTST, insBSET, insBCLR, insBCHG:
		return bitCycles(op)
	case insABCD, insSBCD, insADDX, insSUBX, insCMPM:
		return multiprecisionCycles(op)
	case insBcc, insDBcc, insScc, insBSR:
		return cpu.conditionalCycles(op)
	default:
		return uint64(op.Timing)
	}
}

func (cpu *CPU) conditionalCycles(op DecodedOpcode) uint64 {
	if op.Instruction == insBSR {
		return 18
	}
	condition := cpu.evaluateCondition(op.Extra)
	switch op.Instruction {
	case insBcc:
		if condition {
			return 10
		}
		if op.DstReg == 0 {
			return 12
		}
		return 8
	case insDBcc:
		if condition {
			return 12
		}
		if uint16(cpu.D[op.DstReg]) == 0 {
			return 14
		}
		return 10
	default:
		if op.DstMode != 0 {
			return 8 + eaCycles(op.DstMode, op.DstReg, SizeByte)
		}
		if condition {
			return 6
		}
		return 4
	}
}

func unaryCycles(op DecodedOpcode) uint64 {
	dst := eaCycles(op.DstMode, op.DstReg, op.Size)
	switch op.Instruction {
	case insADDQ, insSUBQ:
		if op.DstMode == 1 {
			return 8
		}
		if op.DstMode == 0 {
			return sizeCycles(op.Size, 4, 8)
		}
		return sizeCycles(op.Size, 8, 12) + dst
	case insTST:
		return 4 + dst
	case insNBCD:
		if op.DstMode == 0 {
			return 6
		}
		return 8 + dst
	case insTAS:
		if op.DstMode == 0 {
			return 4
		}
		return 14 + dst
	default:
		if op.DstMode == 0 {
			return sizeCycles(op.Size, 4, 6)
		}
		return sizeCycles(op.Size, 8, 12) + dst
	}
}

func multiprecisionCycles(op DecodedOpcode) uint64 {
	switch op.Instruction {
	case insABCD, insSBCD:
		if op.Extra != 0 {
			return 18
		}
		return 6
	case insCMPM:
		return sizeCycles(op.Size, 12, 20)
	default:
		if op.Extra != 0 {
			return sizeCycles(op.Size, 18, 30)
		}
		return sizeCycles(op.Size, 4, 8)
	}
}

func eaCycles(mode, reg uint8, size OperandSize) uint64 {
	if mode < 2 {
		return 0
	}
	cycles := sizeCycles(size, 4, 8)
	switch mode {
	case 4:
		return cycles + 2
	case 5:
		return cycles + 4
	case 6:
		return cycles + 6
	case 7:
		switch reg {
		case 0, 2:
			return cycles + 4
		case 1:
			return cycles + 8
		case 3:
			return cycles + 6
		}
	}
	return cycles
}

func sizeCycles(size OperandSize, shortCycles, longCycles uint64) uint64 {
	if size == SizeLong {
		return longCycles
	}
	return shortCycles
}

func moveCycles(op DecodedOpcode) uint64 {
	switch op.Extra {
	case 1, 2:
		return 4
	case 3:
		if op.DstMode == 0 {
			return 6
		}
		return 8 + eaCycles(op.DstMode, op.DstReg, SizeWord)
	case 4, 5:
		return 12 + eaCycles(op.SrcMode, op.SrcReg, SizeWord)
	}
	dst := eaCycles(op.DstMode, op.DstReg, op.Size)
	if op.DstMode == 4 {
		dst -= 2
	}
	return 4 + eaCycles(op.SrcMode, op.SrcReg, op.Size) + dst
}

func aluCycles(op DecodedOpcode) uint64 {
	src := eaCycles(op.SrcMode, op.SrcReg, op.Size)
	if op.Instruction == insCMPA {
		return 6 + src
	}
	if op.Instruction == insCMP {
		return sizeCycles(op.Size, 4, 6) + src
	}
	if op.DstMode >= 2 {
		return sizeCycles(op.Size, 8, 12) + eaCycles(op.DstMode, op.DstReg, op.Size)
	}
	base := uint64(4)
	if op.DstMode == 1 {
		base = 8
	}
	if op.Size == SizeLong {
		base = 6
		if op.SrcMode < 2 || op.SrcMode == 7 && op.SrcReg == 4 {
			base = 8
		}
	}
	return base + src
}

func immediateCycles(op DecodedOpcode) uint64 {
	if op.DstMode == 7 && op.DstReg == 4 {
		return 20
	}
	if op.DstMode == 0 {
		if op.Instruction == insCMPI || op.Instruction == insANDI {
			return sizeCycles(op.Size, 8, 14)
		}
		return sizeCycles(op.Size, 8, 16)
	}
	base := sizeCycles(op.Size, 12, 20)
	if op.Instruction == insCMPI {
		base = sizeCycles(op.Size, 8, 12)
	}
	return base + eaCycles(op.DstMode, op.DstReg, op.Size)
}

func controlCycles(op DecodedOpcode) uint64 {
	mode, reg := op.DstMode, op.DstReg
	if op.Instruction == insLEA || op.Instruction == insMOVEM && op.Extra != 0 {
		mode, reg = op.SrcMode, op.SrcReg
	}
	extra := uint64(0)
	if mode >= 5 {
		extra = eaCycles(mode, reg, SizeWord) - 4
	}
	switch op.Instruction {
	case insMOVEP:
		return sizeCycles(op.Size, 16, 24)
	case insMOVEM:
		return uint64(op.Timing) + extra
	case insLEA, insPEA:
		if mode == 6 || mode == 7 && reg == 3 {
			extra += 2
		}
	default:
		switch extra {
		case 4:
			extra = 2
		case 8:
			extra = 4
		}
	}
	return uint64(op.Timing) + extra
}

func bitCycles(op DecodedOpcode) uint64 {
	base := uint64(8)
	if op.Instruction == insBTST {
		base = 4
	}
	if op.DstMode == 0 {
		if op.Instruction == insBTST || op.Instruction == insBCLR {
			base += 2
		}
	} else {
		base += eaCycles(op.DstMode, op.DstReg, SizeByte)
	}
	if op.SrcMode == 7 {
		base += 4
	}
	return base
}

// divideUnsignedCycles counts the restoring divider's 15 quotient iterations.
func divideUnsignedCycles(dividend uint32, divisor uint16) uint64 {
	if dividend>>16 >= uint32(divisor) {
		return 10
	}
	cycles := uint64(76)
	trial := uint32(divisor) << 16
	for range 15 {
		carry := dividend&0x80000000 != 0
		dividend <<= 1
		if carry {
			dividend -= trial
		} else {
			cycles += 4
			if dividend >= trial {
				dividend -= trial
				cycles -= 2
			}
		}
	}
	return cycles
}

func divideSignedCycles(dividend int32, divisor int16) uint64 {
	cycles := uint64(12)
	numerator, denominator := int64(dividend), int64(divisor)
	if numerator < 0 {
		numerator = -numerator
		cycles += 2
	}
	if denominator < 0 {
		denominator = -denominator
	}
	if numerator>>16 >= denominator {
		return cycles + 4
	}
	quotient := uint32(numerator / denominator)
	cycles += 110
	if divisor >= 0 {
		if dividend >= 0 {
			cycles -= 2
		} else {
			cycles += 2
		}
	}
	for range 15 {
		if quotient&0x8000 == 0 {
			cycles += 2
		}
		quotient <<= 1
	}
	return cycles
}
