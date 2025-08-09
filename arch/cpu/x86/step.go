package x86

import "fmt"

// TraceStep contains information about an executed instruction for debugging and tracing.
type TraceStep struct {
	// Instruction details
	IP          uint16 // instruction pointer before execution
	CS          uint16 // code segment before execution
	Opcode      uint8  // instruction opcode
	Instruction string // assembly instruction text

	// Register state before execution
	PreAX uint16
	PreBX uint16
	PreCX uint16
	PreDX uint16
	PreSI uint16
	PreDI uint16
	PreBP uint16
	PreSP uint16

	// Segment registers before execution
	PreCS uint16
	PreDS uint16
	PreES uint16
	PreSS uint16

	// Flags before execution
	PreFlags Flags

	// Register state after execution
	PostAX uint16
	PostBX uint16
	PostCX uint16
	PostDX uint16
	PostSI uint16
	PostDI uint16
	PostBP uint16
	PostSP uint16

	// Segment registers after execution
	PostCS uint16
	PostDS uint16
	PostES uint16
	PostSS uint16

	// Flags after execution
	PostFlags Flags

	// Execution details
	Cycles uint64 // total cycles before this instruction
	Timing uint8  // cycles consumed by this instruction
	Size   uint8  // instruction size in bytes

	// Memory access (if any)
	MemoryRead    bool
	MemoryWrite   bool
	MemoryAddress uint32
	MemoryValue   uint16
}

// String returns a formatted string representation of the trace step.
func (ts TraceStep) String() string {
	// Format: CS:IP OPCODE INSTRUCTION AX=XXXX BX=XXXX ... FLAGS=XXXX CY=XXXXXXXX
	return fmt.Sprintf("%04X:%04X %02X %-12s AX=%04X BX=%04X CX=%04X DX=%04X SI=%04X DI=%04X BP=%04X SP=%04X FL=%04X CY=%08X",
		ts.CS, ts.IP, ts.Opcode, ts.Instruction,
		ts.PostAX, ts.PostBX, ts.PostCX, ts.PostDX,
		ts.PostSI, ts.PostDI, ts.PostBP, ts.PostSP,
		uint16(ts.PostFlags), ts.Cycles)
}

// DetailedString returns a detailed multi-line representation showing before/after states.
func (ts TraceStep) DetailedString() string {
	var result string

	// Instruction line
	result += fmt.Sprintf("%04X:%04X %02X %-12s (size=%d, cycles=%d)\n",
		ts.CS, ts.IP, ts.Opcode, ts.Instruction, ts.Size, ts.Timing)

	// Register changes
	result += "Registers:\n"
	if ts.PreAX != ts.PostAX {
		result += fmt.Sprintf("  AX: %04X -> %04X\n", ts.PreAX, ts.PostAX)
	}
	if ts.PreBX != ts.PostBX {
		result += fmt.Sprintf("  BX: %04X -> %04X\n", ts.PreBX, ts.PostBX)
	}
	if ts.PreCX != ts.PostCX {
		result += fmt.Sprintf("  CX: %04X -> %04X\n", ts.PreCX, ts.PostCX)
	}
	if ts.PreDX != ts.PostDX {
		result += fmt.Sprintf("  DX: %04X -> %04X\n", ts.PreDX, ts.PostDX)
	}
	if ts.PreSI != ts.PostSI {
		result += fmt.Sprintf("  SI: %04X -> %04X\n", ts.PreSI, ts.PostSI)
	}
	if ts.PreDI != ts.PostDI {
		result += fmt.Sprintf("  DI: %04X -> %04X\n", ts.PreDI, ts.PostDI)
	}
	if ts.PreBP != ts.PostBP {
		result += fmt.Sprintf("  BP: %04X -> %04X\n", ts.PreBP, ts.PostBP)
	}
	if ts.PreSP != ts.PostSP {
		result += fmt.Sprintf("  SP: %04X -> %04X\n", ts.PreSP, ts.PostSP)
	}

	// Segment register changes
	if ts.PreCS != ts.PostCS || ts.PreDS != ts.PostDS || ts.PreES != ts.PostES || ts.PreSS != ts.PostSS {
		result += "Segments:\n"
		if ts.PreCS != ts.PostCS {
			result += fmt.Sprintf("  CS: %04X -> %04X\n", ts.PreCS, ts.PostCS)
		}
		if ts.PreDS != ts.PostDS {
			result += fmt.Sprintf("  DS: %04X -> %04X\n", ts.PreDS, ts.PostDS)
		}
		if ts.PreES != ts.PostES {
			result += fmt.Sprintf("  ES: %04X -> %04X\n", ts.PreES, ts.PostES)
		}
		if ts.PreSS != ts.PostSS {
			result += fmt.Sprintf("  SS: %04X -> %04X\n", ts.PreSS, ts.PostSS)
		}
	}

	// Flag changes
	if ts.PreFlags != ts.PostFlags {
		result += fmt.Sprintf("Flags: %04X -> %04X\n", uint16(ts.PreFlags), uint16(ts.PostFlags))
		result += ts.formatFlagChanges()
	}

	// Memory access
	if ts.MemoryRead || ts.MemoryWrite {
		switch {
		case ts.MemoryRead && ts.MemoryWrite:
			result += fmt.Sprintf("Memory: R/W %06X = %04X\n", ts.MemoryAddress, ts.MemoryValue)
		case ts.MemoryRead:
			result += fmt.Sprintf("Memory: R %06X = %04X\n", ts.MemoryAddress, ts.MemoryValue)
		default: // MemoryWrite only
			result += fmt.Sprintf("Memory: W %06X = %04X\n", ts.MemoryAddress, ts.MemoryValue)
		}
	}

	return result
}

// formatFlagChanges returns a string showing which flags changed.
func (ts TraceStep) formatFlagChanges() string {
	if ts.PreFlags == ts.PostFlags {
		return ""
	}

	var changes []string

	if ts.PreFlags.GetCarry() != ts.PostFlags.GetCarry() {
		if ts.PostFlags.GetCarry() {
			changes = append(changes, "+CF")
		} else {
			changes = append(changes, "-CF")
		}
	}

	if ts.PreFlags.GetZero() != ts.PostFlags.GetZero() {
		if ts.PostFlags.GetZero() {
			changes = append(changes, "+ZF")
		} else {
			changes = append(changes, "-ZF")
		}
	}

	if ts.PreFlags.GetSign() != ts.PostFlags.GetSign() {
		if ts.PostFlags.GetSign() {
			changes = append(changes, "+SF")
		} else {
			changes = append(changes, "-SF")
		}
	}

	if ts.PreFlags.GetOverflow() != ts.PostFlags.GetOverflow() {
		if ts.PostFlags.GetOverflow() {
			changes = append(changes, "+OF")
		} else {
			changes = append(changes, "-OF")
		}
	}

	if ts.PreFlags.GetParity() != ts.PostFlags.GetParity() {
		if ts.PostFlags.GetParity() {
			changes = append(changes, "+PF")
		} else {
			changes = append(changes, "-PF")
		}
	}

	if ts.PreFlags.GetAuxCarry() != ts.PostFlags.GetAuxCarry() {
		if ts.PostFlags.GetAuxCarry() {
			changes = append(changes, "+AF")
		} else {
			changes = append(changes, "-AF")
		}
	}

	if ts.PreFlags.GetInterrupt() != ts.PostFlags.GetInterrupt() {
		if ts.PostFlags.GetInterrupt() {
			changes = append(changes, "+IF")
		} else {
			changes = append(changes, "-IF")
		}
	}

	if ts.PreFlags.GetDirection() != ts.PostFlags.GetDirection() {
		if ts.PostFlags.GetDirection() {
			changes = append(changes, "+DF")
		} else {
			changes = append(changes, "-DF")
		}
	}

	if ts.PreFlags.GetTrap() != ts.PostFlags.GetTrap() {
		if ts.PostFlags.GetTrap() {
			changes = append(changes, "+TF")
		} else {
			changes = append(changes, "-TF")
		}
	}

	result := "  Changed: "
	for i, change := range changes {
		if i > 0 {
			result += ", "
		}
		result += change
	}
	result += "\n"

	return result
}

// GetMemoryAccess returns memory access information as a formatted string.
func (ts TraceStep) GetMemoryAccess() string {
	if !ts.MemoryRead && !ts.MemoryWrite {
		return ""
	}

	var accessType string
	switch {
	case ts.MemoryRead && ts.MemoryWrite:
		accessType = "RW"
	case ts.MemoryRead:
		accessType = "R"
	default: // MemoryWrite only
		accessType = "W"
	}

	return fmt.Sprintf("%s:%06X=%04X", accessType, ts.MemoryAddress, ts.MemoryValue)
}
