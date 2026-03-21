package m6502

// OpcodeID is a compact numeric identifier for a 6502 instruction mnemonic.
// Using OpcodeID instead of string comparisons eliminates hot-path string hashing overhead.
// The zero value InvalidOpcodeID means "not set / unknown".
type OpcodeID uint8

// OpcodeID constants — one per unique mnemonic, in alphabetical order.
// These names are intentionally short (no prefix) to read cleanly in switch statements.
const (
	InvalidOpcodeID OpcodeID = iota // 0 — not set
	Adc
	Alr // unofficial
	Anc // unofficial
	And
	Ane // unofficial
	Arr // unofficial
	Asl
	Axs // unofficial
	Bcc
	Bcs
	Beq
	Bit
	Bmi
	Bne
	Bpl
	Brk
	Bvc
	Bvs
	Clc
	Cld
	Cli
	Clv
	Cmp
	Cpx
	Cpy
	Dcp // unofficial
	Dec
	Dex
	Dey
	Eor
	Inc
	Inx
	Iny
	Isc // unofficial
	Jmp
	Jsr
	Las // unofficial
	Lax // unofficial
	Lda
	Ldx
	Ldy
	Lsr
	Lxa // unofficial
	Nop
	Ora
	Pha
	Php
	Pla
	Plp
	Rla // unofficial
	Rol
	Ror
	Rra // unofficial
	Rti
	Rts
	Sax // unofficial
	Sbc
	Sec
	Sed
	Sei
	Sha // unofficial
	Shx // unofficial
	Shy // unofficial
	Slo // unofficial
	Sre // unofficial
	Sta
	Stx
	Sty
	Tas // unofficial
	Tax
	Tay
	Tsx
	Txa
	Txs
	Tya

	OpcodeIDMax = Tya
)

// NameToOpcodeID maps a lowercase 6502 mnemonic to its OpcodeID for O(1) lookup.
var NameToOpcodeID = map[string]OpcodeID{
	AdcName: Adc,
	AlrName: Alr,
	AncName: Anc,
	AndName: And,
	AneName: Ane,
	ArrName: Arr,
	AslName: Asl,
	AxsName: Axs,
	BccName: Bcc,
	BcsName: Bcs,
	BeqName: Beq,
	BitName: Bit,
	BmiName: Bmi,
	BneName: Bne,
	BplName: Bpl,
	BrkName: Brk,
	BvcName: Bvc,
	BvsName: Bvs,
	ClcName: Clc,
	CldName: Cld,
	CliName: Cli,
	ClvName: Clv,
	CmpName: Cmp,
	CpxName: Cpx,
	CpyName: Cpy,
	DcpName: Dcp,
	DecName: Dec,
	DexName: Dex,
	DeyName: Dey,
	EorName: Eor,
	IncName: Inc,
	InxName: Inx,
	InyName: Iny,
	IscName: Isc,
	JmpName: Jmp,
	JsrName: Jsr,
	LasName: Las,
	LaxName: Lax,
	LdaName: Lda,
	LdxName: Ldx,
	LdyName: Ldy,
	LsrName: Lsr,
	LxaName: Lxa,
	NopName: Nop,
	OraName: Ora,
	PhaName: Pha,
	PhpName: Php,
	PlaName: Pla,
	PlpName: Plp,
	RlaName: Rla,
	RolName: Rol,
	RorName: Ror,
	RraName: Rra,
	RtiName: Rti,
	RtsName: Rts,
	SaxName: Sax,
	SbcName: Sbc,
	SecName: Sec,
	SedName: Sed,
	SeiName: Sei,
	ShaName: Sha,
	ShxName: Shx,
	ShyName: Shy,
	SloName: Slo,
	SreName: Sre,
	StaName: Sta,
	StxName: Stx,
	StyName: Sty,
	TasName: Tas,
	TaxName: Tax,
	TayName: Tay,
	TsxName: Tsx,
	TxaName: Txa,
	TxsName: Txs,
	TyaName: Tya,
}

// OpcodeIDToName maps an OpcodeID back to its lowercase mnemonic for display/debugging.
var OpcodeIDToName = [OpcodeIDMax + 1]string{
	Adc: AdcName,
	Alr: AlrName,
	Anc: AncName,
	And: AndName,
	Ane: AneName,
	Arr: ArrName,
	Asl: AslName,
	Axs: AxsName,
	Bcc: BccName,
	Bcs: BcsName,
	Beq: BeqName,
	Bit: BitName,
	Bmi: BmiName,
	Bne: BneName,
	Bpl: BplName,
	Brk: BrkName,
	Bvc: BvcName,
	Bvs: BvsName,
	Clc: ClcName,
	Cld: CldName,
	Cli: CliName,
	Clv: ClvName,
	Cmp: CmpName,
	Cpx: CpxName,
	Cpy: CpyName,
	Dcp: DcpName,
	Dec: DecName,
	Dex: DexName,
	Dey: DeyName,
	Eor: EorName,
	Inc: IncName,
	Inx: InxName,
	Iny: InyName,
	Isc: IscName,
	Jmp: JmpName,
	Jsr: JsrName,
	Las: LasName,
	Lax: LaxName,
	Lda: LdaName,
	Ldx: LdxName,
	Ldy: LdyName,
	Lsr: LsrName,
	Lxa: LxaName,
	Nop: NopName,
	Ora: OraName,
	Pha: PhaName,
	Php: PhpName,
	Pla: PlaName,
	Plp: PlpName,
	Rla: RlaName,
	Rol: RolName,
	Ror: RorName,
	Rra: RraName,
	Rti: RtiName,
	Rts: RtsName,
	Sax: SaxName,
	Sbc: SbcName,
	Sec: SecName,
	Sed: SedName,
	Sei: SeiName,
	Sha: ShaName,
	Shx: ShxName,
	Shy: ShyName,
	Slo: SloName,
	Sre: SreName,
	Sta: StaName,
	Stx: StxName,
	Sty: StyName,
	Tas: TasName,
	Tax: TaxName,
	Tay: TayName,
	Tsx: TsxName,
	Txa: TxaName,
	Txs: TxsName,
	Tya: TyaName,
}
