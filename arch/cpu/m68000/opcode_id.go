package m68000

// OpcodeID is a compact numeric identifier for a 68000 instruction mnemonic.
// Using OpcodeID instead of string comparisons eliminates hot-path string hashing overhead.
// The zero value InvalidOpcodeID means "not set / unknown".
type OpcodeID uint8

// OpcodeID constants — one per unique mnemonic, in alphabetical order.
// These names are intentionally short (no prefix) to read cleanly in switch statements.
// Note: the 68000 uses uppercase mnemonics; the const names mirror them in TitleCase.
const (
	InvalidOpcodeID OpcodeID = iota // 0 — not set
	Abcd
	Add
	Adda
	Addi
	Addq
	Addx
	And
	Andi
	Asl
	Asr
	Bcc
	Bchg
	Bclr
	Bra
	Bset
	Bsr
	Btst
	Chk
	Clr
	Cmp
	Cmpa
	Cmpi
	Cmpm
	Dbcc
	Divs
	Divu
	Eor
	Eori
	Exg
	Ext
	Illegal
	Jmp
	Jsr
	Lea
	Link
	Lsl
	Lsr
	Move
	Movea
	Movem
	Movep
	Moveq
	Muls
	Mulu
	Nbcd
	Neg
	Negx
	Nop
	Not
	Or
	Ori
	Pea
	Reset
	Rol
	Ror
	Roxl
	Roxr
	Rte
	Rtr
	Rts
	Sbcd
	Scc
	Stop
	Sub
	Suba
	Subi
	Subq
	Subx
	Swap
	Tas
	Trap
	Trapv
	Tst
	Unlk

	OpcodeIDMax = Unlk
)

// NameToOpcodeID maps an uppercase 68000 mnemonic to its OpcodeID for O(1) lookup.
var NameToOpcodeID = map[string]OpcodeID{
	ABCDName:    Abcd,
	ADDName:     Add,
	ADDAName:    Adda,
	ADDIName:    Addi,
	ADDQName:    Addq,
	ADDXName:    Addx,
	ANDName:     And,
	ANDIName:    Andi,
	ASLName:     Asl,
	ASRName:     Asr,
	BccName:     Bcc,
	BCHGName:    Bchg,
	BCLRName:    Bclr,
	BRAName:     Bra,
	BSETName:    Bset,
	BSRName:     Bsr,
	BTSTName:    Btst,
	CHKName:     Chk,
	CLRName:     Clr,
	CMPName:     Cmp,
	CMPAName:    Cmpa,
	CMPIName:    Cmpi,
	CMPMName:    Cmpm,
	DBccName:    Dbcc,
	DIVSName:    Divs,
	DIVUName:    Divu,
	EORName:     Eor,
	EORIName:    Eori,
	EXGName:     Exg,
	EXTName:     Ext,
	ILLEGALName: Illegal,
	JMPName:     Jmp,
	JSRName:     Jsr,
	LEAName:     Lea,
	LINKName:    Link,
	LSLName:     Lsl,
	LSRName:     Lsr,
	MOVEName:    Move,
	MOVEAName:   Movea,
	MOVEMName:   Movem,
	MOVEPName:   Movep,
	MOVEQName:   Moveq,
	MULSName:    Muls,
	MULUName:    Mulu,
	NBCDName:    Nbcd,
	NEGName:     Neg,
	NEGXName:    Negx,
	NOPName:     Nop,
	NOTName:     Not,
	ORName:      Or,
	ORIName:     Ori,
	PEAName:     Pea,
	RESETName:   Reset,
	ROLName:     Rol,
	RORName:     Ror,
	ROXLName:    Roxl,
	ROXRName:    Roxr,
	RTEName:     Rte,
	RTRName:     Rtr,
	RTSName:     Rts,
	SBCDName:    Sbcd,
	SccName:     Scc,
	STOPName:    Stop,
	SUBName:     Sub,
	SUBAName:    Suba,
	SUBIName:    Subi,
	SUBQName:    Subq,
	SUBXName:    Subx,
	SWAPName:    Swap,
	TASName:     Tas,
	TRAPName:    Trap,
	TRAPVName:   Trapv,
	TSTName:     Tst,
	UNLKName:    Unlk,
}

// OpcodeIDToName maps an OpcodeID back to its uppercase mnemonic for display/debugging.
var OpcodeIDToName = [OpcodeIDMax + 1]string{
	Abcd:    ABCDName,
	Add:     ADDName,
	Adda:    ADDAName,
	Addi:    ADDIName,
	Addq:    ADDQName,
	Addx:    ADDXName,
	And:     ANDName,
	Andi:    ANDIName,
	Asl:     ASLName,
	Asr:     ASRName,
	Bcc:     BccName,
	Bchg:    BCHGName,
	Bclr:    BCLRName,
	Bra:     BRAName,
	Bset:    BSETName,
	Bsr:     BSRName,
	Btst:    BTSTName,
	Chk:     CHKName,
	Clr:     CLRName,
	Cmp:     CMPName,
	Cmpa:    CMPAName,
	Cmpi:    CMPIName,
	Cmpm:    CMPMName,
	Dbcc:    DBccName,
	Divs:    DIVSName,
	Divu:    DIVUName,
	Eor:     EORName,
	Eori:    EORIName,
	Exg:     EXGName,
	Ext:     EXTName,
	Illegal: ILLEGALName,
	Jmp:     JMPName,
	Jsr:     JSRName,
	Lea:     LEAName,
	Link:    LINKName,
	Lsl:     LSLName,
	Lsr:     LSRName,
	Move:    MOVEName,
	Movea:   MOVEAName,
	Movem:   MOVEMName,
	Movep:   MOVEPName,
	Moveq:   MOVEQName,
	Muls:    MULSName,
	Mulu:    MULUName,
	Nbcd:    NBCDName,
	Neg:     NEGName,
	Negx:    NEGXName,
	Nop:     NOPName,
	Not:     NOTName,
	Or:      ORName,
	Ori:     ORIName,
	Pea:     PEAName,
	Reset:   RESETName,
	Rol:     ROLName,
	Ror:     RORName,
	Roxl:    ROXLName,
	Roxr:    ROXRName,
	Rte:     RTEName,
	Rtr:     RTRName,
	Rts:     RTSName,
	Sbcd:    SBCDName,
	Scc:     SccName,
	Stop:    STOPName,
	Sub:     SUBName,
	Suba:    SUBAName,
	Subi:    SUBIName,
	Subq:    SUBQName,
	Subx:    SUBXName,
	Swap:    SWAPName,
	Tas:     TASName,
	Trap:    TRAPName,
	Trapv:   TRAPVName,
	Tst:     TSTName,
	Unlk:    UNLKName,
}
