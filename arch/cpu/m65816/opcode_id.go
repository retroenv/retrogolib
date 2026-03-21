package m65816

// OpcodeID is a compact numeric identifier for a 65816 instruction mnemonic.
// Using OpcodeID instead of string comparisons eliminates hot-path string hashing overhead.
// The zero value InvalidOpcodeID means "not set / unknown".
type OpcodeID uint8

// OpcodeID constants — one per unique mnemonic, in alphabetical order.
// These names are intentionally short (no prefix) to read cleanly in switch statements.
const (
	InvalidOpcodeID OpcodeID = iota // 0 — not set
	Adc
	And
	Asl
	Bcc
	Bcs
	Beq
	Bit
	Bmi
	Bne
	Bpl
	Bra
	Brk
	Brl
	Bvc
	Bvs
	Clc
	Cld
	Cli
	Clv
	Cmp
	Cop
	Cpx
	Cpy
	Dec
	Dex
	Dey
	Eor
	Inc
	Inx
	Iny
	Jml
	Jmp
	Jsl
	Jsr
	Lda
	Ldx
	Ldy
	Lsr
	Mvn
	Mvp
	Nop
	Ora
	Pea
	Pei
	Per
	Pha
	Phb
	Phd
	Phk
	Php
	Phx
	Phy
	Pla
	Plb
	Pld
	Plp
	Plx
	Ply
	Rep
	Rol
	Ror
	Rti
	Rtl
	Rts
	Sbc
	Sec
	Sed
	Sei
	Sep
	Sta
	Stp
	Stx
	Sty
	Stz
	Tax
	Tay
	Tcd
	Tcs
	Tdc
	Trb
	Tsb
	Tsc
	Tsx
	Txa
	Txs
	Txy
	Tya
	Tyx
	Wai
	Wdm
	Xba
	Xce

	OpcodeIDMax = Xce
)

// NameToOpcodeID maps a lowercase 65816 mnemonic to its OpcodeID for O(1) lookup.
var NameToOpcodeID = map[string]OpcodeID{
	AdcName: Adc,
	AndName: And,
	AslName: Asl,
	BccName: Bcc,
	BcsName: Bcs,
	BeqName: Beq,
	BitName: Bit,
	BmiName: Bmi,
	BneName: Bne,
	BplName: Bpl,
	BraName: Bra,
	BrkName: Brk,
	BrlName: Brl,
	BvcName: Bvc,
	BvsName: Bvs,
	ClcName: Clc,
	CldName: Cld,
	CliName: Cli,
	ClvName: Clv,
	CmpName: Cmp,
	CopName: Cop,
	CpxName: Cpx,
	CpyName: Cpy,
	DecName: Dec,
	DexName: Dex,
	DeyName: Dey,
	EorName: Eor,
	IncName: Inc,
	InxName: Inx,
	InyName: Iny,
	JmlName: Jml,
	JmpName: Jmp,
	JslName: Jsl,
	JsrName: Jsr,
	LdaName: Lda,
	LdxName: Ldx,
	LdyName: Ldy,
	LsrName: Lsr,
	MvnName: Mvn,
	MvpName: Mvp,
	NopName: Nop,
	OraName: Ora,
	PeaName: Pea,
	PeiName: Pei,
	PerName: Per,
	PhaName: Pha,
	PhbName: Phb,
	PhdName: Phd,
	PhkName: Phk,
	PhpName: Php,
	PhxName: Phx,
	PhyName: Phy,
	PlaName: Pla,
	PlbName: Plb,
	PldName: Pld,
	PlpName: Plp,
	PlxName: Plx,
	PlyName: Ply,
	RepName: Rep,
	RolName: Rol,
	RorName: Ror,
	RtiName: Rti,
	RtlName: Rtl,
	RtsName: Rts,
	SbcName: Sbc,
	SecName: Sec,
	SedName: Sed,
	SeiName: Sei,
	SepName: Sep,
	StaName: Sta,
	StpName: Stp,
	StxName: Stx,
	StyName: Sty,
	StzName: Stz,
	TaxName: Tax,
	TayName: Tay,
	TcdName: Tcd,
	TcsName: Tcs,
	TdcName: Tdc,
	TrbName: Trb,
	TsbName: Tsb,
	TscName: Tsc,
	TsxName: Tsx,
	TxaName: Txa,
	TxsName: Txs,
	TxyName: Txy,
	TyaName: Tya,
	TyxName: Tyx,
	WaiName: Wai,
	WdmName: Wdm,
	XbaName: Xba,
	XceName: Xce,
}

// OpcodeIDToName maps an OpcodeID back to its lowercase mnemonic for display/debugging.
var OpcodeIDToName = [OpcodeIDMax + 1]string{
	Adc: AdcName,
	And: AndName,
	Asl: AslName,
	Bcc: BccName,
	Bcs: BcsName,
	Beq: BeqName,
	Bit: BitName,
	Bmi: BmiName,
	Bne: BneName,
	Bpl: BplName,
	Bra: BraName,
	Brk: BrkName,
	Brl: BrlName,
	Bvc: BvcName,
	Bvs: BvsName,
	Clc: ClcName,
	Cld: CldName,
	Cli: CliName,
	Clv: ClvName,
	Cmp: CmpName,
	Cop: CopName,
	Cpx: CpxName,
	Cpy: CpyName,
	Dec: DecName,
	Dex: DexName,
	Dey: DeyName,
	Eor: EorName,
	Inc: IncName,
	Inx: InxName,
	Iny: InyName,
	Jml: JmlName,
	Jmp: JmpName,
	Jsl: JslName,
	Jsr: JsrName,
	Lda: LdaName,
	Ldx: LdxName,
	Ldy: LdyName,
	Lsr: LsrName,
	Mvn: MvnName,
	Mvp: MvpName,
	Nop: NopName,
	Ora: OraName,
	Pea: PeaName,
	Pei: PeiName,
	Per: PerName,
	Pha: PhaName,
	Phb: PhbName,
	Phd: PhdName,
	Phk: PhkName,
	Php: PhpName,
	Phx: PhxName,
	Phy: PhyName,
	Pla: PlaName,
	Plb: PlbName,
	Pld: PldName,
	Plp: PlpName,
	Plx: PlxName,
	Ply: PlyName,
	Rep: RepName,
	Rol: RolName,
	Ror: RorName,
	Rti: RtiName,
	Rtl: RtlName,
	Rts: RtsName,
	Sbc: SbcName,
	Sec: SecName,
	Sed: SedName,
	Sei: SeiName,
	Sep: SepName,
	Sta: StaName,
	Stp: StpName,
	Stx: StxName,
	Sty: StyName,
	Stz: StzName,
	Tax: TaxName,
	Tay: TayName,
	Tcd: TcdName,
	Tcs: TcsName,
	Tdc: TdcName,
	Trb: TrbName,
	Tsb: TsbName,
	Tsc: TscName,
	Tsx: TsxName,
	Txa: TxaName,
	Txs: TxsName,
	Txy: TxyName,
	Tya: TyaName,
	Tyx: TyxName,
	Wai: WaiName,
	Wdm: WdmName,
	Xba: XbaName,
	Xce: XceName,
}
