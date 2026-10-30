package x86

// CPUGeneration identifies the minimum processor for an opcode encoding.
// Unknown encodings require an explicit classification before use.
type CPUGeneration uint8

const (
	CPUUnknown CPUGeneration = iota
	CPU8086
	CPU80186
	CPU80286
	CPU80386
	CPU80486
)

// String returns the processor name used in diagnostics.
func (cpu CPUGeneration) String() string {
	names := [...]string{"unknown", "8086", "80186", "80286", "80386", "80486"}
	if int(cpu) >= len(names) {
		return names[0]
	}
	return names[cpu]
}

// Supports reports whether the processor supports this opcode generation.
// Operand sizes, execution mode, privileges, and coprocessors need separate
// checks. This method does not establish complete instruction legality.
func (cpu CPUGeneration) Supports(info OpcodeInfo) bool {
	minimum := info.MinimumCPU()
	return minimum != CPUUnknown && minimum <= cpu && cpu <= CPU80486
}

// MinimumCPU classifies the opcode and prefix values used by this package.
// Different forms of one mnemonic can require different processors. In
// particular, immediate PUSH and immediate shift forms require the 80186.
func (o OpcodeInfo) MinimumCPU() CPUGeneration {
	if o.Size == 0 {
		return CPUUnknown
	}
	if o.IsTwoByte() {
		return minimumExtendedCPU(o.Opcode)
	}
	switch o.Opcode {
	case 0x63:
		return CPU80286
	case 0x64, 0x65, 0x66, 0x67:
		return CPU80386
	case 0x60, 0x61, 0x62, 0x68, 0x69, 0x6a, 0x6b, 0x6c, 0x6d, 0x6e, 0x6f,
		0xc0, 0xc1, 0xc8, 0xc9:
		return CPU80186
	}

	if o.Opcode == 0x0f || o.Opcode == 0xd6 || o.Opcode == 0xf1 {
		return CPUUnknown
	}
	return CPU8086
}

func minimumExtendedCPU(opcode uint16) CPUGeneration {
	switch {
	case opcode == 0x0f01:
		return CPU80286
	case opcode >= 0x0fc8 && opcode <= 0x0fcf,
		opcode == 0x0fb0, opcode == 0x0fb1, opcode == 0x0fc0, opcode == 0x0fc1,
		opcode == 0x0f08, opcode == 0x0f09:
		return CPU80486
	case opcode >= 0x0f80 && opcode <= 0x0f9f,
		opcode == 0x0fbc, opcode == 0x0fbd, opcode == 0x0fa3, opcode == 0x0fbb,
		opcode == 0x0fb3, opcode == 0x0fab, opcode == 0x0fb6, opcode == 0x0fb7,
		opcode == 0x0fbe, opcode == 0x0fbf, opcode == 0x0fa4, opcode == 0x0fa5,
		opcode == 0x0fac, opcode == 0x0fad:
		return CPU80386
	default:
		return CPUUnknown
	}
}
