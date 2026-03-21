package sm83

// OpcodeID is a compact numeric identifier for an SM83 instruction mnemonic.
// Using OpcodeID instead of string comparisons eliminates hot-path string hashing overhead.
// The zero value InvalidOpcodeID means "not set / unknown".
type OpcodeID uint8

// OpcodeID constants — one per unique mnemonic, in alphabetical order.
// These names are intentionally short (no prefix) to read cleanly in switch statements.
const (
	InvalidOpcodeID OpcodeID = iota // 0 — not set
	Adc
	Add
	And
	BitOp // "bit" mnemonic — renamed to avoid conflict with the Bit type in addressing.go
	Call
	Ccf
	Cp
	Cpl
	Daa
	Dec
	Di
	Ei
	Halt
	Inc
	Jp
	Jr
	Ld
	Ldh
	Nop
	Or
	Pop
	Push
	Res
	Ret
	Reti
	Rl
	Rla
	Rlc
	Rlca
	Rr
	Rra
	Rrc
	Rrca
	Rst
	Sbc
	Scf
	Set
	Sla
	Sra
	Srl
	Stop
	Sub
	Swap
	Xor

	OpcodeIDMax = Xor
)

// NameToOpcodeID maps a lowercase SM83 mnemonic to its OpcodeID for O(1) lookup.
var NameToOpcodeID = map[string]OpcodeID{
	AdcName:  Adc,
	AddName:  Add,
	AndName:  And,
	BitName:  BitOp,
	CallName: Call,
	CcfName:  Ccf,
	CpName:   Cp,
	CplName:  Cpl,
	DaaName:  Daa,
	DecName:  Dec,
	DiName:   Di,
	EiName:   Ei,
	HaltName: Halt,
	IncName:  Inc,
	JpName:   Jp,
	JrName:   Jr,
	LdName:   Ld,
	LdhName:  Ldh,
	NopName:  Nop,
	OrName:   Or,
	PopName:  Pop,
	PushName: Push,
	ResName:  Res,
	RetName:  Ret,
	RetiName: Reti,
	RlName:   Rl,
	RlaName:  Rla,
	RlcName:  Rlc,
	RlcaName: Rlca,
	RrName:   Rr,
	RraName:  Rra,
	RrcName:  Rrc,
	RrcaName: Rrca,
	RstName:  Rst,
	SbcName:  Sbc,
	ScfName:  Scf,
	SetName:  Set,
	SlaName:  Sla,
	SraName:  Sra,
	SrlName:  Srl,
	StopName: Stop,
	SubName:  Sub,
	SwapName: Swap,
	XorName:  Xor,
}

// OpcodeIDToName maps an OpcodeID back to its lowercase mnemonic for display/debugging.
var OpcodeIDToName = [OpcodeIDMax + 1]string{
	Adc:   AdcName,
	Add:   AddName,
	And:   AndName,
	BitOp: BitName,
	Call:  CallName,
	Ccf:   CcfName,
	Cp:    CpName,
	Cpl:   CplName,
	Daa:   DaaName,
	Dec:   DecName,
	Di:    DiName,
	Ei:    EiName,
	Halt:  HaltName,
	Inc:   IncName,
	Jp:    JpName,
	Jr:    JrName,
	Ld:    LdName,
	Ldh:   LdhName,
	Nop:   NopName,
	Or:    OrName,
	Pop:   PopName,
	Push:  PushName,
	Res:   ResName,
	Ret:   RetName,
	Reti:  RetiName,
	Rl:    RlName,
	Rla:   RlaName,
	Rlc:   RlcName,
	Rlca:  RlcaName,
	Rr:    RrName,
	Rra:   RraName,
	Rrc:   RrcName,
	Rrca:  RrcaName,
	Rst:   RstName,
	Sbc:   SbcName,
	Scf:   ScfName,
	Set:   SetName,
	Sla:   SlaName,
	Sra:   SraName,
	Srl:   SrlName,
	Stop:  StopName,
	Sub:   SubName,
	Swap:  SwapName,
	Xor:   XorName,
}
