package x86

// The values in this file come from the Intel 8086 Family User's Manual,
// October 1979 (order number 9800722-03): table 2-20 for the effective address
// calculation time and table 2-21 for the instruction times. The times are for
// the 8086 and 8088 and assume that the instruction is in the prefetch queue.

// Clock adjustments of tables 2-20 and 2-21.
const (
	// SegmentOverrideClocks is added for a segment override prefix.
	SegmentOverrideClocks = 2
	// WordTransferClocks is added for each word transfer at an odd address on
	// the 8086, and for each word transfer on the 8088.
	WordTransferClocks = 4
)

// eaClocks gives the effective address calculation time for the R/M values of
// ModR/M modes 0, 1, and 2. Mode 0 with R/M 6 is a direct address.
var eaClocks = [3][8]uint8{
	{7, 8, 8, 7, 5, 5, 6, 5},
	{11, 12, 12, 11, 9, 9, 9, 9},
	{11, 12, 12, 11, 9, 9, 9, 9},
}

// OperandForm identifies an operand combination of table 2-21.
type OperandForm uint8

// Operand forms of table 2-21. Accumulator forms use the short encodings
// without a ModR/M byte.
const (
	FormUnknown OperandForm = iota
	FormNone
	FormImmediate
	FormRegister8
	FormRegister16
	FormMemory8
	FormMemory16
	FormSegmentRegister
	FormRegisterRegister
	FormRegisterMemory
	FormMemoryRegister
	FormRegisterImmediate
	FormMemoryImmediate
	FormAccumulatorImmediate
	FormAccumulatorMemory
	FormMemoryAccumulator
	FormAccumulatorRegister
	FormRegisterOne
	FormRegisterCL
	FormMemoryOne
	FormMemoryCL
	FormShortLabel
	FormNearLabel
	FormFarLabel
	FormRegisterPointer
	FormMemoryPointer
	FormMemoryFarPointer
	FormString
	FormRepeatedString
)

// Timing is the 8086 execution time of one operand form.
type Timing struct {
	// Clocks is the minimum time. For a conditional jump, it is the time when
	// the jump is not taken.
	Clocks uint8
	// MaxClocks is the maximum time without the PerCount part. It is equal to
	// Clocks when the time is fixed.
	MaxClocks uint8
	// PerCount is added for each shifted bit or each string repetition.
	PerCount uint8
	// Transfers is the number of memory transfers. For a repeated string form,
	// it is the number for each repetition.
	Transfers uint8
	// EA tells that the effective address time of the memory operand is added.
	EA bool
}

// EAClocks returns the effective address calculation time of a ModR/M memory
// operand (table 2-20). Mode 3 selects a register, which has no effective
// address, so the result is false.
func EAClocks(mod, rm uint8) (uint8, bool) {
	if mod > 2 || rm > 7 {
		return 0, false
	}
	return eaClocks[mod][rm], true
}

// InstructionTiming returns the time of an instruction name and operand form.
// The result is false when table 2-21 has no value for the combination.
func InstructionTiming(name string, form OperandForm) (Timing, bool) {
	forms, ok := instructionTimings[name]
	if !ok {
		return Timing{}, false
	}
	timing, ok := forms[form]

	return timing, ok
}

func fixed(clocks, transfers uint8) Timing {
	return Timing{
		Clocks:    clocks,
		MaxClocks: clocks,
		Transfers: transfers,
	}
}

func withEA(clocks, transfers uint8) Timing {
	return Timing{
		Clocks:    clocks,
		MaxClocks: clocks,
		Transfers: transfers,
		EA:        true,
	}
}

func ranged(minimum, maximum, transfers uint8, ea bool) Timing {
	return Timing{
		Clocks:    minimum,
		MaxClocks: maximum,
		Transfers: transfers,
		EA:        ea,
	}
}

func perCount(clocks, count, transfers uint8, ea bool) Timing {
	return Timing{
		Clocks:    clocks,
		MaxClocks: clocks,
		PerCount:  count,
		Transfers: transfers,
		EA:        ea,
	}
}

// arithmeticTimings is the row of ADC, ADD, AND, OR, SBB, SUB, and XOR.
func arithmeticTimings() map[OperandForm]Timing {
	return map[OperandForm]Timing{
		FormRegisterRegister:     fixed(3, 0),
		FormRegisterMemory:       withEA(9, 1),
		FormMemoryRegister:       withEA(16, 2),
		FormRegisterImmediate:    fixed(4, 0),
		FormMemoryImmediate:      withEA(17, 2),
		FormAccumulatorImmediate: fixed(4, 0),
	}
}

// unaryTimings is the row of INC and DEC.
func unaryTimings() map[OperandForm]Timing {
	return map[OperandForm]Timing{
		FormRegister16: fixed(2, 0),
		FormRegister8:  fixed(3, 0),
		FormMemory8:    withEA(15, 2),
		FormMemory16:   withEA(15, 2),
	}
}

// complementTimings is the row of NEG and NOT.
func complementTimings() map[OperandForm]Timing {
	return map[OperandForm]Timing{
		FormRegister8:  fixed(3, 0),
		FormRegister16: fixed(3, 0),
		FormMemory8:    withEA(16, 2),
		FormMemory16:   withEA(16, 2),
	}
}

// shiftTimings is the row of RCL, RCR, ROL, ROR, SAR, SHL, and SHR.
func shiftTimings() map[OperandForm]Timing {
	return map[OperandForm]Timing{
		FormRegisterOne: fixed(2, 0),
		FormRegisterCL:  perCount(8, 4, 0, false),
		FormMemoryOne:   withEA(15, 2),
		FormMemoryCL:    perCount(20, 4, 2, true),
	}
}

// multiplyTimings gives the rows of MUL, IMUL, DIV, and IDIV. The arguments
// are the minimum and maximum times for byte and word register operands. A
// memory operand adds six clocks to both limits.
func multiplyTimings(byteMin, byteMax, wordMin, wordMax uint8) map[OperandForm]Timing {
	return map[OperandForm]Timing{
		FormRegister8:  ranged(byteMin, byteMax, 0, false),
		FormRegister16: ranged(wordMin, wordMax, 0, false),
		FormMemory8:    ranged(byteMin+6, byteMax+6, 1, true),
		FormMemory16:   ranged(wordMin+6, wordMax+6, 1, true),
	}
}

// conditionalJumpTimings is the row of a conditional jump: 16 clocks when the
// jump is taken and 4 clocks when it is not taken.
func conditionalJumpTimings() map[OperandForm]Timing {
	return map[OperandForm]Timing{
		FormShortLabel: ranged(4, 16, 0, false),
	}
}

func impliedTiming(clocks uint8) map[OperandForm]Timing {
	return map[OperandForm]Timing{
		FormNone: fixed(clocks, 0),
	}
}

// stringTimings gives the row of a string instruction: the time of one
// operation and the time of a repeated operation as 9 plus repeat clocks for
// each repetition.
func stringTimings(single, repeat, transfers uint8) map[OperandForm]Timing {
	return map[OperandForm]Timing{
		FormString:         fixed(single, transfers),
		FormRepeatedString: perCount(9, repeat, transfers, false),
	}
}

var instructionTimings = map[string]map[OperandForm]Timing{
	AdcName: arithmeticTimings(),
	AddName: arithmeticTimings(),
	AndName: arithmeticTimings(),
	OrName:  arithmeticTimings(),
	SbbName: arithmeticTimings(),
	SubName: arithmeticTimings(),
	XorName: arithmeticTimings(),
	CmpName: {
		FormRegisterRegister:     fixed(3, 0),
		FormRegisterMemory:       withEA(9, 1),
		FormMemoryRegister:       withEA(9, 1),
		FormRegisterImmediate:    fixed(4, 0),
		FormMemoryImmediate:      withEA(10, 1),
		FormAccumulatorImmediate: fixed(4, 0),
	},
	TestName: {
		FormRegisterRegister:     fixed(3, 0),
		FormRegisterMemory:       withEA(9, 1),
		FormMemoryRegister:       withEA(9, 1),
		FormAccumulatorImmediate: fixed(4, 0),
		FormRegisterImmediate:    fixed(5, 0),
		FormMemoryImmediate:      withEA(11, 1),
	},
	MovName: {
		FormMemoryAccumulator: fixed(10, 1),
		FormAccumulatorMemory: fixed(10, 1),
		FormRegisterRegister:  fixed(2, 0),
		FormRegisterMemory:    withEA(8, 1),
		FormMemoryRegister:    withEA(9, 1),
		FormRegisterImmediate: fixed(4, 0),
		FormMemoryImmediate:   withEA(10, 1),
	},
	XchgName: {
		FormAccumulatorRegister: fixed(3, 0),
		FormMemoryRegister:      withEA(17, 2),
		FormRegisterMemory:      withEA(17, 2),
		FormRegisterRegister:    fixed(4, 0),
	},
	LeaName: {
		FormRegisterMemory: withEA(2, 0),
	},
	IncName: unaryTimings(),
	DecName: unaryTimings(),
	NegName: complementTimings(),
	NotName: complementTimings(),
	RclName: shiftTimings(),
	RcrName: shiftTimings(),
	RolName: shiftTimings(),
	RorName: shiftTimings(),
	SarName: shiftTimings(),
	ShlName: shiftTimings(),
	ShrName: shiftTimings(),

	MulName:  multiplyTimings(70, 77, 118, 133),
	ImulName: multiplyTimings(80, 98, 128, 154),
	DivName:  multiplyTimings(80, 90, 144, 162),
	IdivName: multiplyTimings(101, 112, 165, 184),

	PushName: {
		FormRegister16:      fixed(11, 1),
		FormSegmentRegister: fixed(10, 1),
		FormMemory16:        withEA(16, 2),
	},
	PopName: {
		FormRegister16:      fixed(8, 1),
		FormSegmentRegister: fixed(8, 1),
		FormMemory16:        withEA(17, 2),
	},
	CallName: {
		FormNearLabel:        fixed(19, 1),
		FormFarLabel:         fixed(28, 2),
		FormMemoryPointer:    withEA(21, 2),
		FormRegisterPointer:  fixed(16, 1),
		FormMemoryFarPointer: withEA(37, 4),
	},
	JmpName: {
		FormShortLabel:       fixed(15, 0),
		FormNearLabel:        fixed(15, 0),
		FormFarLabel:         fixed(15, 0),
		FormMemoryPointer:    withEA(18, 1),
		FormRegisterPointer:  fixed(11, 0),
		FormMemoryFarPointer: withEA(24, 2),
	},
	RetName: {
		FormNone:      fixed(8, 1),
		FormImmediate: fixed(12, 1),
	},
	RetfName: {
		FormNone:      fixed(18, 2),
		FormImmediate: fixed(17, 2),
	},
	IntName: {
		FormImmediate: fixed(51, 5),
	},

	JbName:   conditionalJumpTimings(),
	JbeName:  conditionalJumpTimings(),
	JlName:   conditionalJumpTimings(),
	JleName:  conditionalJumpTimings(),
	JnbName:  conditionalJumpTimings(),
	JnbeName: conditionalJumpTimings(),
	JnlName:  conditionalJumpTimings(),
	JnleName: conditionalJumpTimings(),
	JnoName:  conditionalJumpTimings(),
	JnpName:  conditionalJumpTimings(),
	JnsName:  conditionalJumpTimings(),
	JnzName:  conditionalJumpTimings(),
	JoName:   conditionalJumpTimings(),
	JpName:   conditionalJumpTimings(),
	JsName:   conditionalJumpTimings(),
	JzName:   conditionalJumpTimings(),

	CbwName:   impliedTiming(2),
	CwdName:   impliedTiming(5),
	ClcName:   impliedTiming(2),
	StcName:   impliedTiming(2),
	CldName:   impliedTiming(2),
	StdName:   impliedTiming(2),
	CliName:   impliedTiming(2),
	StiName:   impliedTiming(2),
	NopName:   impliedTiming(3),
	RepName:   impliedTiming(2),
	RepzName:  impliedTiming(2),
	RepnzName: impliedTiming(2),

	MovsbName: stringTimings(18, 17, 2),
	MovswName: stringTimings(18, 17, 2),
	CmpsbName: stringTimings(22, 22, 2),
	CmpswName: stringTimings(22, 22, 2),
	ScasbName: stringTimings(15, 15, 1),
	ScaswName: stringTimings(15, 15, 1),
	LodsbName: stringTimings(12, 13, 1),
	LodswName: stringTimings(12, 13, 1),
	StosbName: stringTimings(11, 10, 1),
	StoswName: stringTimings(11, 10, 1),
}
