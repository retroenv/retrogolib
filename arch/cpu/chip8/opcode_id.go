package chip8

// OpcodeID is a compact numeric identifier for a Chip-8 instruction mnemonic.
// Using OpcodeID instead of string comparisons eliminates hot-path string hashing overhead.
// The zero value InvalidOpcodeID means "not set / unknown".
type OpcodeID uint8

// OpcodeID constants — one per unique mnemonic, in alphabetical order.
// These names are intentionally short (no prefix) to read cleanly in switch statements.
const (
	InvalidOpcodeID OpcodeID = iota // 0 — not set
	Add
	And
	Call
	Cls
	Drw
	Jp
	Ld
	Or
	Ret
	Rnd
	Se
	Shl
	Shr
	Skp
	Sknp
	Sne
	Sub
	Subn
	Xor

	OpcodeIDMax = Xor
)

// NameToOpcodeID maps a lowercase Chip-8 mnemonic to its OpcodeID for O(1) lookup.
var NameToOpcodeID = map[string]OpcodeID{
	AddName:  Add,
	AndName:  And,
	CallName: Call,
	ClsName:  Cls,
	DrwName:  Drw,
	JpName:   Jp,
	LdName:   Ld,
	OrName:   Or,
	RetName:  Ret,
	RndName:  Rnd,
	SeName:   Se,
	ShlName:  Shl,
	ShrName:  Shr,
	SkpName:  Skp,
	SknpName: Sknp,
	SneName:  Sne,
	SubName:  Sub,
	SubnName: Subn,
	XorName:  Xor,
}

// OpcodeIDToName maps an OpcodeID back to its lowercase mnemonic for display/debugging.
var OpcodeIDToName = [OpcodeIDMax + 1]string{
	Add:  AddName,
	And:  AndName,
	Call: CallName,
	Cls:  ClsName,
	Drw:  DrwName,
	Jp:   JpName,
	Ld:   LdName,
	Or:   OrName,
	Ret:  RetName,
	Rnd:  RndName,
	Se:   SeName,
	Shl:  ShlName,
	Shr:  ShrName,
	Skp:  SkpName,
	Sknp: SknpName,
	Sne:  SneName,
	Sub:  SubName,
	Subn: SubnName,
	Xor:  XorName,
}
