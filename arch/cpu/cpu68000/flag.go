package cpu68000

// Flags contains the condition code register (CCR) flags of the CPU.
// The CCR is the low byte of the status register.
//
// Bit layout of CCR:
//
//	Bit  4   3   2   1   0
//	Flag X   N   Z   V   C
//
// X (Extend): Like carry but not affected by all operations.
// N (Negative): Set if the result is negative (MSB set).
// Z (Zero): Set if the result is zero.
// V (Overflow): Set on arithmetic overflow.
// C (Carry): Set on carry/borrow.
type Flags struct {
	C uint8 // carry flag
	V uint8 // overflow flag
	Z uint8 // zero flag
	N uint8 // negative flag
	X uint8 // extend flag
}

// Status register system byte bit positions.
const (
	FlagTrace      = 15 // Trace mode enable
	FlagSupervisor = 13 // Supervisor/user mode
	FlagIPM2       = 10 // Interrupt priority mask bit 2
	FlagIPM1       = 9  // Interrupt priority mask bit 1
	FlagIPM0       = 8  // Interrupt priority mask bit 0
)

// Status register masks.
const (
	MaskTrace      = 1 << FlagTrace
	MaskSupervisor = 1 << FlagSupervisor
	MaskIPM        = 7 << FlagIPM0                        // Interrupt priority mask (3 bits)
	MaskSystem     = MaskTrace | MaskSupervisor | MaskIPM // Implemented system bits.
	MaskCCR        = 0x001F                               // CCR byte mask (X, N, Z, V, C)
)

// GetCCR returns the condition code register as a byte.
func (cpu *CPU) GetCCR() uint8 {
	return cpu.Flags.C |
		cpu.Flags.V<<1 |
		cpu.Flags.Z<<2 |
		cpu.Flags.N<<3 |
		cpu.Flags.X<<4
}

// GetSR returns the full 16-bit status register.
func (cpu *CPU) GetSR() uint16 {
	return cpu.sr&MaskSystem | uint16(cpu.GetCCR())
}

// SetCCR sets the condition code register from a byte value.
func (cpu *CPU) SetCCR(ccr uint8) {
	cpu.Flags.C = ccr & 1
	cpu.Flags.V = (ccr >> 1) & 1
	cpu.Flags.Z = (ccr >> 2) & 1
	cpu.Flags.N = (ccr >> 3) & 1
	cpu.Flags.X = (ccr >> 4) & 1
}

// SetSR sets the full 16-bit status register.
// This may cause a privilege mode switch.
func (cpu *CPU) SetSR(sr uint16) {
	oldSupervisor := cpu.sr & MaskSupervisor

	cpu.sr = sr & MaskSystem
	cpu.SetCCR(uint8(sr & MaskCCR))

	newSupervisor := cpu.sr & MaskSupervisor

	// Handle privilege mode switch.
	if oldSupervisor != 0 && newSupervisor == 0 {
		// Switching from supervisor to user mode: save SSP, load USP.
		cpu.SSP = cpu.sp
		cpu.sp = cpu.USP
	} else if oldSupervisor == 0 && newSupervisor != 0 {
		// Switching from user to supervisor mode: save USP, load SSP.
		cpu.USP = cpu.sp
		cpu.sp = cpu.SSP
	}
}

// IsSupervisor returns whether the CPU is in supervisor mode.
func (cpu *CPU) IsSupervisor() bool {
	return cpu.sr&MaskSupervisor != 0
}

// InterruptMask returns the current interrupt priority mask (0-7).
func (cpu *CPU) InterruptMask() uint8 {
	return uint8((cpu.sr & MaskIPM) >> FlagIPM0)
}

// setFlagN sets the negative flag based on the MSB of a value for the given size.
func (cpu *CPU) setFlagN(value uint32, size OperandSize) {
	switch size {
	case SizeByte:
		setFlag(&cpu.Flags.N, value&0x80 != 0)
	case SizeWord:
		setFlag(&cpu.Flags.N, value&0x8000 != 0)
	case SizeLong:
		setFlag(&cpu.Flags.N, value&0x80000000 != 0)
	}
}

// setFlagZ sets the zero flag based on a masked value for the given size.
func (cpu *CPU) setFlagZ(value uint32, size OperandSize) {
	setFlag(&cpu.Flags.Z, maskValue(value, size) == 0)
}

// setFlag sets a flag to 1 if condition is true, 0 otherwise.
func setFlag(flag *uint8, condition bool) {
	if condition {
		*flag = 1
	} else {
		*flag = 0
	}
}

// maskValue masks a value to the appropriate size.
func maskValue(value uint32, size OperandSize) uint32 {
	switch size {
	case SizeByte:
		return value & 0xFF
	case SizeWord:
		return value & 0xFFFF
	case SizeLong:
		return value
	default:
		return value
	}
}

// signExtend sign-extends a value from the given size to 32 bits.
func signExtend(value uint32, size OperandSize) uint32 {
	switch size {
	case SizeByte:
		return uint32(int32(int8(value)))
	case SizeWord:
		return uint32(int32(int16(value)))
	default:
		return value
	}
}

// msbMask returns the mask for the most significant bit of the given size.
func msbMask(size OperandSize) uint32 {
	switch size {
	case SizeByte:
		return 0x80
	case SizeWord:
		return 0x8000
	case SizeLong:
		return 0x80000000
	default:
		return 0
	}
}
