package cpu6502

// OpcodeID is a compact numeric identifier for a 6502 instruction mnemonic.
// Using OpcodeID instead of string comparisons eliminates hot-path string hashing overhead.
// The zero value InvalidOpcodeID means "not set / unknown".
type OpcodeID uint8

// OpcodeID constants — one per unique mnemonic. The original range is
// alphabetical and stable; later ISA additions are appended without renumbering it.
// These names are intentionally short where they do not collide with instruction values.
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
	Bra // 65C02
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
	Kil // unofficial — halts the CPU
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
	Phx // 65C02
	Phy // 65C02
	Pla
	Plp
	Plx // 65C02
	Ply // 65C02
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
	Stz // 65C02
	Tas // unofficial
	Tax
	Tay
	Trb // 65C02
	Tsb // 65C02
	Tsx
	Txa
	Txs
	Tya

	// Rockwell/WDC bit instructions were added after the original stable ID
	// range and intentionally remain appended to avoid renumbering it.
	Bbr0ID
	Bbr1ID
	Bbr2ID
	Bbr3ID
	Bbr4ID
	Bbr5ID
	Bbr6ID
	Bbr7ID
	Bbs0ID
	Bbs1ID
	Bbs2ID
	Bbs3ID
	Bbs4ID
	Bbs5ID
	Bbs6ID
	Bbs7ID
	Rmb0ID
	Rmb1ID
	Rmb2ID
	Rmb3ID
	Rmb4ID
	Rmb5ID
	Rmb6ID
	Rmb7ID
	Smb0ID
	Smb1ID
	Smb2ID
	Smb3ID
	Smb4ID
	Smb5ID
	Smb6ID
	Smb7ID

	OpcodeIDMax = Smb7ID
)

// NameToOpcodeID maps a lowercase 6502 mnemonic to its OpcodeID for O(1) lookup.
var NameToOpcodeID = map[string]OpcodeID{
	AdcName:   Adc,
	AlrName:   Alr,
	AncName:   Anc,
	AndName:   And,
	AneName:   Ane,
	ArrName:   Arr,
	AslName:   Asl,
	AxsName:   Axs,
	Bbr0.Name: Bbr0ID,
	Bbr1.Name: Bbr1ID,
	Bbr2.Name: Bbr2ID,
	Bbr3.Name: Bbr3ID,
	Bbr4.Name: Bbr4ID,
	Bbr5.Name: Bbr5ID,
	Bbr6.Name: Bbr6ID,
	Bbr7.Name: Bbr7ID,
	BccName:   Bcc,
	BcsName:   Bcs,
	BeqName:   Beq,
	BitName:   Bit,
	BmiName:   Bmi,
	BneName:   Bne,
	BplName:   Bpl,
	BraName:   Bra,
	BrkName:   Brk,
	Bbs0.Name: Bbs0ID,
	Bbs1.Name: Bbs1ID,
	Bbs2.Name: Bbs2ID,
	Bbs3.Name: Bbs3ID,
	Bbs4.Name: Bbs4ID,
	Bbs5.Name: Bbs5ID,
	Bbs6.Name: Bbs6ID,
	Bbs7.Name: Bbs7ID,
	BvcName:   Bvc,
	BvsName:   Bvs,
	ClcName:   Clc,
	CldName:   Cld,
	CliName:   Cli,
	ClvName:   Clv,
	CmpName:   Cmp,
	CpxName:   Cpx,
	CpyName:   Cpy,
	DcpName:   Dcp,
	DecName:   Dec,
	DexName:   Dex,
	DeyName:   Dey,
	EorName:   Eor,
	IncName:   Inc,
	InxName:   Inx,
	InyName:   Iny,
	IscName:   Isc,
	JmpName:   Jmp,
	JsrName:   Jsr,
	KilName:   Kil,
	LasName:   Las,
	LaxName:   Lax,
	LdaName:   Lda,
	LdxName:   Ldx,
	LdyName:   Ldy,
	LsrName:   Lsr,
	LxaName:   Lxa,
	NopName:   Nop,
	OraName:   Ora,
	PhaName:   Pha,
	PhpName:   Php,
	PhxName:   Phx,
	PhyName:   Phy,
	PlaName:   Pla,
	PlpName:   Plp,
	PlxName:   Plx,
	PlyName:   Ply,
	RlaName:   Rla,
	Rmb0.Name: Rmb0ID,
	Rmb1.Name: Rmb1ID,
	Rmb2.Name: Rmb2ID,
	Rmb3.Name: Rmb3ID,
	Rmb4.Name: Rmb4ID,
	Rmb5.Name: Rmb5ID,
	Rmb6.Name: Rmb6ID,
	Rmb7.Name: Rmb7ID,
	RolName:   Rol,
	RorName:   Ror,
	RraName:   Rra,
	RtiName:   Rti,
	RtsName:   Rts,
	SaxName:   Sax,
	SbcName:   Sbc,
	SecName:   Sec,
	SedName:   Sed,
	SeiName:   Sei,
	ShaName:   Sha,
	ShxName:   Shx,
	ShyName:   Shy,
	SloName:   Slo,
	Smb0.Name: Smb0ID,
	Smb1.Name: Smb1ID,
	Smb2.Name: Smb2ID,
	Smb3.Name: Smb3ID,
	Smb4.Name: Smb4ID,
	Smb5.Name: Smb5ID,
	Smb6.Name: Smb6ID,
	Smb7.Name: Smb7ID,
	SreName:   Sre,
	StaName:   Sta,
	StxName:   Stx,
	StyName:   Sty,
	StzName:   Stz,
	TasName:   Tas,
	TaxName:   Tax,
	TayName:   Tay,
	TrbName:   Trb,
	TsbName:   Tsb,
	TsxName:   Tsx,
	TxaName:   Txa,
	TxsName:   Txs,
	TyaName:   Tya,
}

// OpcodeIDToName maps an OpcodeID back to its lowercase mnemonic for display/debugging.
var OpcodeIDToName = [OpcodeIDMax + 1]string{
	Adc:    AdcName,
	Alr:    AlrName,
	Anc:    AncName,
	And:    AndName,
	Ane:    AneName,
	Arr:    ArrName,
	Asl:    AslName,
	Axs:    AxsName,
	Bbr0ID: Bbr0.Name,
	Bbr1ID: Bbr1.Name,
	Bbr2ID: Bbr2.Name,
	Bbr3ID: Bbr3.Name,
	Bbr4ID: Bbr4.Name,
	Bbr5ID: Bbr5.Name,
	Bbr6ID: Bbr6.Name,
	Bbr7ID: Bbr7.Name,
	Bcc:    BccName,
	Bcs:    BcsName,
	Beq:    BeqName,
	Bit:    BitName,
	Bmi:    BmiName,
	Bne:    BneName,
	Bpl:    BplName,
	Bra:    BraName,
	Brk:    BrkName,
	Bbs0ID: Bbs0.Name,
	Bbs1ID: Bbs1.Name,
	Bbs2ID: Bbs2.Name,
	Bbs3ID: Bbs3.Name,
	Bbs4ID: Bbs4.Name,
	Bbs5ID: Bbs5.Name,
	Bbs6ID: Bbs6.Name,
	Bbs7ID: Bbs7.Name,
	Bvc:    BvcName,
	Bvs:    BvsName,
	Clc:    ClcName,
	Cld:    CldName,
	Cli:    CliName,
	Clv:    ClvName,
	Cmp:    CmpName,
	Cpx:    CpxName,
	Cpy:    CpyName,
	Dcp:    DcpName,
	Dec:    DecName,
	Dex:    DexName,
	Dey:    DeyName,
	Eor:    EorName,
	Inc:    IncName,
	Inx:    InxName,
	Iny:    InyName,
	Isc:    IscName,
	Jmp:    JmpName,
	Jsr:    JsrName,
	Kil:    KilName,
	Las:    LasName,
	Lax:    LaxName,
	Lda:    LdaName,
	Ldx:    LdxName,
	Ldy:    LdyName,
	Lsr:    LsrName,
	Lxa:    LxaName,
	Nop:    NopName,
	Ora:    OraName,
	Pha:    PhaName,
	Php:    PhpName,
	Phx:    PhxName,
	Phy:    PhyName,
	Pla:    PlaName,
	Plp:    PlpName,
	Plx:    PlxName,
	Ply:    PlyName,
	Rla:    RlaName,
	Rmb0ID: Rmb0.Name,
	Rmb1ID: Rmb1.Name,
	Rmb2ID: Rmb2.Name,
	Rmb3ID: Rmb3.Name,
	Rmb4ID: Rmb4.Name,
	Rmb5ID: Rmb5.Name,
	Rmb6ID: Rmb6.Name,
	Rmb7ID: Rmb7.Name,
	Rol:    RolName,
	Ror:    RorName,
	Rra:    RraName,
	Rti:    RtiName,
	Rts:    RtsName,
	Sax:    SaxName,
	Sbc:    SbcName,
	Sec:    SecName,
	Sed:    SedName,
	Sei:    SeiName,
	Sha:    ShaName,
	Shx:    ShxName,
	Shy:    ShyName,
	Slo:    SloName,
	Smb0ID: Smb0.Name,
	Smb1ID: Smb1.Name,
	Smb2ID: Smb2.Name,
	Smb3ID: Smb3.Name,
	Smb4ID: Smb4.Name,
	Smb5ID: Smb5.Name,
	Smb6ID: Smb6.Name,
	Smb7ID: Smb7.Name,
	Sre:    SreName,
	Sta:    StaName,
	Stx:    StxName,
	Sty:    StyName,
	Stz:    StzName,
	Tas:    TasName,
	Tax:    TaxName,
	Tay:    TayName,
	Trb:    TrbName,
	Tsb:    TsbName,
	Tsx:    TsxName,
	Txa:    TxaName,
	Txs:    TxsName,
	Tya:    TyaName,
}
