package z80

// OpcodeID is a compact numeric identifier for a Z80 instruction mnemonic.
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
	Cpd
	Cpdr
	Cpi
	Cpir
	Cpl
	Daa
	DdcbShiftOp
	Dec
	Di
	Djnz
	Ei
	Ex
	Exx
	FdcbShiftOp
	Halt
	Im
	In
	Inc
	Ind
	Indr
	Inf
	Ini
	Inir
	Jp
	Jr
	Ld
	Ldd
	Lddr
	Ldi
	Ldir
	Neg
	Nop
	Or
	Otdr
	Otir
	Out
	Outd
	Outf
	Outi
	Pop
	Push
	Res
	Ret
	Reti
	Retn
	Rl
	Rla
	Rlc
	Rlca
	Rld
	Rr
	Rra
	Rrc
	Rrca
	Rrd
	Rst
	Sbc
	Scf
	Set
	Sla
	Sll
	Sra
	Srl
	Sub
	Xor

	OpcodeIDMax = Xor
)

// NameToOpcodeID maps a lowercase Z80 mnemonic to its OpcodeID for O(1) lookup.
var NameToOpcodeID = map[string]OpcodeID{
	AdcName:       Adc,
	AddName:       Add,
	AndName:       And,
	BitName:       BitOp,
	CallName:      Call,
	CcfName:       Ccf,
	CpName:        Cp,
	CpdName:       Cpd,
	CpdrName:      Cpdr,
	CpiName:       Cpi,
	CpirName:      Cpir,
	CplName:       Cpl,
	DaaName:       Daa,
	DdcbShiftName: DdcbShiftOp,
	DecName:       Dec,
	DiName:        Di,
	DjnzName:      Djnz,
	EiName:        Ei,
	ExName:        Ex,
	ExxName:       Exx,
	FdcbShiftName: FdcbShiftOp,
	HaltName:      Halt,
	ImName:        Im,
	InName:        In,
	IncName:       Inc,
	IndName:       Ind,
	IndrName:      Indr,
	InfName:       Inf,
	IniName:       Ini,
	InirName:      Inir,
	JpName:        Jp,
	JrName:        Jr,
	LdName:        Ld,
	LddName:       Ldd,
	LddrName:      Lddr,
	LdiName:       Ldi,
	LdirName:      Ldir,
	NegName:       Neg,
	NopName:       Nop,
	OrName:        Or,
	OtdrName:      Otdr,
	OtirName:      Otir,
	OutName:       Out,
	OutdName:      Outd,
	OutfName:      Outf,
	OutiName:      Outi,
	PopName:       Pop,
	PushName:      Push,
	ResName:       Res,
	RetName:       Ret,
	RetiName:      Reti,
	RetnName:      Retn,
	RlName:        Rl,
	RlaName:       Rla,
	RlcName:       Rlc,
	RlcaName:      Rlca,
	RldName:       Rld,
	RrName:        Rr,
	RraName:       Rra,
	RrcName:       Rrc,
	RrcaName:      Rrca,
	RrdName:       Rrd,
	RstName:       Rst,
	SbcName:       Sbc,
	ScfName:       Scf,
	SetName:       Set,
	SlaName:       Sla,
	SllName:       Sll,
	SraName:       Sra,
	SrlName:       Srl,
	SubName:       Sub,
	XorName:       Xor,
}

// OpcodeIDToName maps an OpcodeID back to its lowercase mnemonic for display/debugging.
var OpcodeIDToName = [OpcodeIDMax + 1]string{
	Adc:         AdcName,
	Add:         AddName,
	And:         AndName,
	BitOp:       BitName,
	Call:        CallName,
	Ccf:         CcfName,
	Cp:          CpName,
	Cpd:         CpdName,
	Cpdr:        CpdrName,
	Cpi:         CpiName,
	Cpir:        CpirName,
	Cpl:         CplName,
	Daa:         DaaName,
	DdcbShiftOp: DdcbShiftName,
	Dec:         DecName,
	Di:          DiName,
	Djnz:        DjnzName,
	Ei:          EiName,
	Ex:          ExName,
	Exx:         ExxName,
	FdcbShiftOp: FdcbShiftName,
	Halt:        HaltName,
	Im:          ImName,
	In:          InName,
	Inc:         IncName,
	Ind:         IndName,
	Indr:        IndrName,
	Inf:         InfName,
	Ini:         IniName,
	Inir:        InirName,
	Jp:          JpName,
	Jr:          JrName,
	Ld:          LdName,
	Ldd:         LddName,
	Lddr:        LddrName,
	Ldi:         LdiName,
	Ldir:        LdirName,
	Neg:         NegName,
	Nop:         NopName,
	Or:          OrName,
	Otdr:        OtdrName,
	Otir:        OtirName,
	Out:         OutName,
	Outd:        OutdName,
	Outf:        OutfName,
	Outi:        OutiName,
	Pop:         PopName,
	Push:        PushName,
	Res:         ResName,
	Ret:         RetName,
	Reti:        RetiName,
	Retn:        RetnName,
	Rl:          RlName,
	Rla:         RlaName,
	Rlc:         RlcName,
	Rlca:        RlcaName,
	Rld:         RldName,
	Rr:          RrName,
	Rra:         RraName,
	Rrc:         RrcName,
	Rrca:        RrcaName,
	Rrd:         RrdName,
	Rst:         RstName,
	Sbc:         SbcName,
	Scf:         ScfName,
	Set:         SetName,
	Sla:         SlaName,
	Sll:         SllName,
	Sra:         SraName,
	Srl:         SrlName,
	Sub:         SubName,
	Xor:         XorName,
}
