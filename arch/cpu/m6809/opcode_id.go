package m6809

// OpcodeID is a compact numeric identifier for a 6809 instruction mnemonic.
// Using OpcodeID instead of string comparisons eliminates hot-path string hashing overhead.
// The zero value InvalidOpcodeID means "not set / unknown".
type OpcodeID uint8

// OpcodeID constants - one per unique mnemonic, in alphabetical order.
// These names are intentionally short (no prefix) to read cleanly in switch statements.
const (
	InvalidOpcodeID OpcodeID = iota // 0 - not set
	Abx
	Adca
	Adcb
	Adda
	Addb
	Addd
	Anda
	Andb
	Andcc
	Asl
	Asr
	Bita
	Bitb
	Bcc
	Bcs
	Beq
	Bge
	Bgt
	Bhi
	Ble
	Bls
	Blt
	Bmi
	Bne
	Bpl
	Bra
	Brn
	Bsr
	Bvc
	Bvs
	Clr
	Cmpa
	Cmpb
	Cmpd
	Cmps
	Cmpu
	Cmpx
	Cmpy
	Com
	Cwai
	Daa
	Dec
	Eora
	Eorb
	Exg
	Inc
	Jmp
	Jsr
	Lbcc
	Lbcs
	Lbeq
	Lbge
	Lbgt
	Lbhi
	Lble
	Lbls
	Lblt
	Lbmi
	Lbne
	Lbpl
	Lbra
	Lbrn
	Lbsr
	Lbvc
	Lbvs
	Lda
	Ldb
	Ldd
	Lds
	Ldu
	Ldx
	Ldy
	Leax
	Leay
	Leas
	Leau
	Lsr
	Mul
	Neg
	Nop
	Ora
	Orb
	Orcc
	Pshs
	Pshu
	Puls
	Pulu
	Rol
	Ror
	Rti
	Rts
	Sbca
	Sbcb
	Sex
	Sta
	Stb
	Std
	Sts
	Stu
	Stx
	Sty
	Suba
	Subb
	Subd
	Swi
	Swi2
	Swi3
	Sync
	Tfr
	Tst

	OpcodeIDMax = Tst
)

// NameToOpcodeID maps a lowercase 6809 mnemonic to its OpcodeID for O(1) lookup.
var NameToOpcodeID = map[string]OpcodeID{
	AbxName:   Abx,
	AdcaName:  Adca,
	AdcbName:  Adcb,
	AddaName:  Adda,
	AddbName:  Addb,
	AdddName:  Addd,
	AndaName:  Anda,
	AndbName:  Andb,
	AndccName: Andcc,
	AslName:   Asl,
	AsrName:   Asr,
	BitaName:  Bita,
	BitbName:  Bitb,
	BccName:   Bcc,
	BcsName:   Bcs,
	BeqName:   Beq,
	BgeName:   Bge,
	BgtName:   Bgt,
	BhiName:   Bhi,
	BleName:   Ble,
	BlsName:   Bls,
	BltName:   Blt,
	BmiName:   Bmi,
	BneName:   Bne,
	BplName:   Bpl,
	BraName:   Bra,
	BrnName:   Brn,
	BsrName:   Bsr,
	BvcName:   Bvc,
	BvsName:   Bvs,
	ClrName:   Clr,
	CmpaName:  Cmpa,
	CmpbName:  Cmpb,
	CmpdName:  Cmpd,
	CmpsName:  Cmps,
	CmpuName:  Cmpu,
	CmpxName:  Cmpx,
	CmpyName:  Cmpy,
	ComName:   Com,
	CwaiName:  Cwai,
	DaaName:   Daa,
	DecName:   Dec,
	EoraName:  Eora,
	EorbName:  Eorb,
	ExgName:   Exg,
	IncName:   Inc,
	JmpName:   Jmp,
	JsrName:   Jsr,
	LbccName:  Lbcc,
	LbcsName:  Lbcs,
	LbeqName:  Lbeq,
	LbgeName:  Lbge,
	LbgtName:  Lbgt,
	LbhiName:  Lbhi,
	LbleName:  Lble,
	LblsName:  Lbls,
	LbltName:  Lblt,
	LbmiName:  Lbmi,
	LbneName:  Lbne,
	LbplName:  Lbpl,
	LbraName:  Lbra,
	LbrnName:  Lbrn,
	LbsrName:  Lbsr,
	LbvcName:  Lbvc,
	LbvsName:  Lbvs,
	LdaName:   Lda,
	LdbName:   Ldb,
	LddName:   Ldd,
	LdsName:   Lds,
	LduName:   Ldu,
	LdxName:   Ldx,
	LdyName:   Ldy,
	LeaxName:  Leax,
	LeayName:  Leay,
	LeasName:  Leas,
	LeauName:  Leau,
	LsrName:   Lsr,
	MulName:   Mul,
	NegName:   Neg,
	NopName:   Nop,
	OraName:   Ora,
	OrbName:   Orb,
	OrccName:  Orcc,
	PshsName:  Pshs,
	PshuName:  Pshu,
	PulsName:  Puls,
	PuluName:  Pulu,
	RolName:   Rol,
	RorName:   Ror,
	RtiName:   Rti,
	RtsName:   Rts,
	SbcaName:  Sbca,
	SbcbName:  Sbcb,
	SexName:   Sex,
	StaName:   Sta,
	StbName:   Stb,
	StdName:   Std,
	StsName:   Sts,
	StuName:   Stu,
	StxName:   Stx,
	StyName:   Sty,
	SubaName:  Suba,
	SubbName:  Subb,
	SubdName:  Subd,
	SwiName:   Swi,
	Swi2Name:  Swi2,
	Swi3Name:  Swi3,
	SyncName:  Sync,
	TfrName:   Tfr,
	TstName:   Tst,
}

// OpcodeIDToName maps an OpcodeID back to its lowercase mnemonic for display/debugging.
var OpcodeIDToName = [OpcodeIDMax + 1]string{
	Abx:   AbxName,
	Adca:  AdcaName,
	Adcb:  AdcbName,
	Adda:  AddaName,
	Addb:  AddbName,
	Addd:  AdddName,
	Anda:  AndaName,
	Andb:  AndbName,
	Andcc: AndccName,
	Asl:   AslName,
	Asr:   AsrName,
	Bita:  BitaName,
	Bitb:  BitbName,
	Bcc:   BccName,
	Bcs:   BcsName,
	Beq:   BeqName,
	Bge:   BgeName,
	Bgt:   BgtName,
	Bhi:   BhiName,
	Ble:   BleName,
	Bls:   BlsName,
	Blt:   BltName,
	Bmi:   BmiName,
	Bne:   BneName,
	Bpl:   BplName,
	Bra:   BraName,
	Brn:   BrnName,
	Bsr:   BsrName,
	Bvc:   BvcName,
	Bvs:   BvsName,
	Clr:   ClrName,
	Cmpa:  CmpaName,
	Cmpb:  CmpbName,
	Cmpd:  CmpdName,
	Cmps:  CmpsName,
	Cmpu:  CmpuName,
	Cmpx:  CmpxName,
	Cmpy:  CmpyName,
	Com:   ComName,
	Cwai:  CwaiName,
	Daa:   DaaName,
	Dec:   DecName,
	Eora:  EoraName,
	Eorb:  EorbName,
	Exg:   ExgName,
	Inc:   IncName,
	Jmp:   JmpName,
	Jsr:   JsrName,
	Lbcc:  LbccName,
	Lbcs:  LbcsName,
	Lbeq:  LbeqName,
	Lbge:  LbgeName,
	Lbgt:  LbgtName,
	Lbhi:  LbhiName,
	Lble:  LbleName,
	Lbls:  LblsName,
	Lblt:  LbltName,
	Lbmi:  LbmiName,
	Lbne:  LbneName,
	Lbpl:  LbplName,
	Lbra:  LbraName,
	Lbrn:  LbrnName,
	Lbsr:  LbsrName,
	Lbvc:  LbvcName,
	Lbvs:  LbvsName,
	Lda:   LdaName,
	Ldb:   LdbName,
	Ldd:   LddName,
	Lds:   LdsName,
	Ldu:   LduName,
	Ldx:   LdxName,
	Ldy:   LdyName,
	Leax:  LeaxName,
	Leay:  LeayName,
	Leas:  LeasName,
	Leau:  LeauName,
	Lsr:   LsrName,
	Mul:   MulName,
	Neg:   NegName,
	Nop:   NopName,
	Ora:   OraName,
	Orb:   OrbName,
	Orcc:  OrccName,
	Pshs:  PshsName,
	Pshu:  PshuName,
	Puls:  PulsName,
	Pulu:  PuluName,
	Rol:   RolName,
	Ror:   RorName,
	Rti:   RtiName,
	Rts:   RtsName,
	Sbca:  SbcaName,
	Sbcb:  SbcbName,
	Sex:   SexName,
	Sta:   StaName,
	Stb:   StbName,
	Std:   StdName,
	Sts:   StsName,
	Stu:   StuName,
	Stx:   StxName,
	Sty:   StyName,
	Suba:  SubaName,
	Subb:  SubbName,
	Subd:  SubdName,
	Swi:   SwiName,
	Swi2:  Swi2Name,
	Swi3:  Swi3Name,
	Sync:  SyncName,
	Tfr:   TfrName,
	Tst:   TstName,
}
