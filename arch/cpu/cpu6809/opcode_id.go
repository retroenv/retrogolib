package cpu6809

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
	Asla
	Aslb
	Asr
	Asra
	Asrb
	Bcc
	Bcs
	Beq
	Bge
	Bgt
	Bhi
	Bita
	Bitb
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
	Clra
	Clrb
	Cmpa
	Cmpb
	Cmpd
	Cmps
	Cmpu
	Cmpx
	Cmpy
	Com
	Coma
	Comb
	Cwai
	Daa
	Dec
	Deca
	Decb
	Eora
	Eorb
	Exg
	Inc
	Inca
	Incb
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
	Leas
	Leau
	Leax
	Leay
	Lsr
	Lsra
	Lsrb
	Mul
	Neg
	Nega
	Negb
	Nop
	Ora
	Orb
	Orcc
	Pshs
	Pshu
	Puls
	Pulu
	Rol
	Rola
	Rolb
	Ror
	Rora
	Rorb
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
	Tsta
	Tstb

	OpcodeIDMax = Tstb
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
	AslaName:  Asla,
	AslbName:  Aslb,
	AsrName:   Asr,
	AsraName:  Asra,
	AsrbName:  Asrb,
	BccName:   Bcc,
	BcsName:   Bcs,
	BeqName:   Beq,
	BgeName:   Bge,
	BgtName:   Bgt,
	BhiName:   Bhi,
	BitaName:  Bita,
	BitbName:  Bitb,
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
	ClraName:  Clra,
	ClrbName:  Clrb,
	CmpaName:  Cmpa,
	CmpbName:  Cmpb,
	CmpdName:  Cmpd,
	CmpsName:  Cmps,
	CmpuName:  Cmpu,
	CmpxName:  Cmpx,
	CmpyName:  Cmpy,
	ComName:   Com,
	ComaName:  Coma,
	CombName:  Comb,
	CwaiName:  Cwai,
	DaaName:   Daa,
	DecName:   Dec,
	DecaName:  Deca,
	DecbName:  Decb,
	EoraName:  Eora,
	EorbName:  Eorb,
	ExgName:   Exg,
	IncName:   Inc,
	IncaName:  Inca,
	IncbName:  Incb,
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
	LeasName:  Leas,
	LeauName:  Leau,
	LeaxName:  Leax,
	LeayName:  Leay,
	LsrName:   Lsr,
	LsraName:  Lsra,
	LsrbName:  Lsrb,
	MulName:   Mul,
	NegName:   Neg,
	NegaName:  Nega,
	NegbName:  Negb,
	NopName:   Nop,
	OraName:   Ora,
	OrbName:   Orb,
	OrccName:  Orcc,
	PshsName:  Pshs,
	PshuName:  Pshu,
	PulsName:  Puls,
	PuluName:  Pulu,
	RolName:   Rol,
	RolaName:  Rola,
	RolbName:  Rolb,
	RorName:   Ror,
	RoraName:  Rora,
	RorbName:  Rorb,
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
	TstaName:  Tsta,
	TstbName:  Tstb,
}

// OpcodeIDToName maps an OpcodeID back to its lowercase mnemonic for display and debugging.
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
	Asla:  AslaName,
	Aslb:  AslbName,
	Asr:   AsrName,
	Asra:  AsraName,
	Asrb:  AsrbName,
	Bcc:   BccName,
	Bcs:   BcsName,
	Beq:   BeqName,
	Bge:   BgeName,
	Bgt:   BgtName,
	Bhi:   BhiName,
	Bita:  BitaName,
	Bitb:  BitbName,
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
	Clra:  ClraName,
	Clrb:  ClrbName,
	Cmpa:  CmpaName,
	Cmpb:  CmpbName,
	Cmpd:  CmpdName,
	Cmps:  CmpsName,
	Cmpu:  CmpuName,
	Cmpx:  CmpxName,
	Cmpy:  CmpyName,
	Com:   ComName,
	Coma:  ComaName,
	Comb:  CombName,
	Cwai:  CwaiName,
	Daa:   DaaName,
	Dec:   DecName,
	Deca:  DecaName,
	Decb:  DecbName,
	Eora:  EoraName,
	Eorb:  EorbName,
	Exg:   ExgName,
	Inc:   IncName,
	Inca:  IncaName,
	Incb:  IncbName,
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
	Leas:  LeasName,
	Leau:  LeauName,
	Leax:  LeaxName,
	Leay:  LeayName,
	Lsr:   LsrName,
	Lsra:  LsraName,
	Lsrb:  LsrbName,
	Mul:   MulName,
	Neg:   NegName,
	Nega:  NegaName,
	Negb:  NegbName,
	Nop:   NopName,
	Ora:   OraName,
	Orb:   OrbName,
	Orcc:  OrccName,
	Pshs:  PshsName,
	Pshu:  PshuName,
	Puls:  PulsName,
	Pulu:  PuluName,
	Rol:   RolName,
	Rola:  RolaName,
	Rolb:  RolbName,
	Ror:   RorName,
	Rora:  RoraName,
	Rorb:  RorbName,
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
	Tsta:  TstaName,
	Tstb:  TstbName,
}
