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
	result += ts.formatRegisterChanges()

	// Segment register changes
	result += ts.formatSegmentChanges()

	// Flag changes
	result += ts.formatFlags()

	// Memory access
	result += ts.formatMemoryAccess()

	return result
}

// formatFlagChanges returns a string showing which flags changed.
func (ts TraceStep) formatFlagChanges() string {
	if ts.PreFlags == ts.PostFlags {
		return ""
	}

	flagChecks := []struct {
		name    string
		preVal  bool
		postVal bool
	}{
		{"CF", ts.PreFlags.GetCarry(), ts.PostFlags.GetCarry()},
		{"ZF", ts.PreFlags.GetZero(), ts.PostFlags.GetZero()},
		{"SF", ts.PreFlags.GetSign(), ts.PostFlags.GetSign()},
		{"OF", ts.PreFlags.GetOverflow(), ts.PostFlags.GetOverflow()},
		{"PF", ts.PreFlags.GetParity(), ts.PostFlags.GetParity()},
		{"AF", ts.PreFlags.GetAuxCarry(), ts.PostFlags.GetAuxCarry()},
		{"IF", ts.PreFlags.GetInterrupt(), ts.PostFlags.GetInterrupt()},
		{"DF", ts.PreFlags.GetDirection(), ts.PostFlags.GetDirection()},
		{"TF", ts.PreFlags.GetTrap(), ts.PostFlags.GetTrap()},
	}

	var changes []string
	for _, flag := range flagChecks {
		if flag.preVal != flag.postVal {
			prefix := "+"
			if !flag.postVal {
				prefix = "-"
			}
			changes = append(changes, prefix+flag.name)
		}
	}

	if len(changes) == 0 {
		return ""
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

// formatRegisterChanges formats general purpose register changes.
func (ts TraceStep) formatRegisterChanges() string {
	regChanges := []struct {
		name      string
		preValue  uint16
		postValue uint16
	}{
		{"AX", ts.PreAX, ts.PostAX},
		{"BX", ts.PreBX, ts.PostBX},
		{"CX", ts.PreCX, ts.PostCX},
		{"DX", ts.PreDX, ts.PostDX},
		{"SI", ts.PreSI, ts.PostSI},
		{"DI", ts.PreDI, ts.PostDI},
		{"BP", ts.PreBP, ts.PostBP},
		{"SP", ts.PreSP, ts.PostSP},
	}

	var result string
	hasChanges := false
	for _, reg := range regChanges {
		if reg.preValue != reg.postValue {
			if !hasChanges {
				result += "Registers:\n"
				hasChanges = true
			}
			result += fmt.Sprintf("  %s: %04X -> %04X\n", reg.name, reg.preValue, reg.postValue)
		}
	}
	return result
}

// formatSegmentChanges formats segment register changes.
func (ts TraceStep) formatSegmentChanges() string {
	segChanges := []struct {
		name      string
		preValue  uint16
		postValue uint16
	}{
		{"CS", ts.PreCS, ts.PostCS},
		{"DS", ts.PreDS, ts.PostDS},
		{"ES", ts.PreES, ts.PostES},
		{"SS", ts.PreSS, ts.PostSS},
	}

	var result string
	hasChanges := false
	for _, seg := range segChanges {
		if seg.preValue != seg.postValue {
			if !hasChanges {
				result += "Segments:\n"
				hasChanges = true
			}
			result += fmt.Sprintf("  %s: %04X -> %04X\n", seg.name, seg.preValue, seg.postValue)
		}
	}
	return result
}

// formatFlags formats flag changes.
func (ts TraceStep) formatFlags() string {
	if ts.PreFlags == ts.PostFlags {
		return ""
	}
	result := fmt.Sprintf("Flags: %04X -> %04X\n", uint16(ts.PreFlags), uint16(ts.PostFlags))
	result += ts.formatFlagChanges()
	return result
}

// formatMemoryAccess formats memory access information.
func (ts TraceStep) formatMemoryAccess() string {
	if !ts.MemoryRead && !ts.MemoryWrite {
		return ""
	}

	var accessType string
	switch {
	case ts.MemoryRead && ts.MemoryWrite:
		accessType = "R/W"
	case ts.MemoryRead:
		accessType = "R"
	default:
		accessType = "W"
	}

	return fmt.Sprintf("Memory: %s %06X = %04X\n", accessType, ts.MemoryAddress, ts.MemoryValue)
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
