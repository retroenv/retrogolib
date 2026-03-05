package m68000

// Instruction defines a 68000 CPU instruction with its execution logic.
// The 68000 uses a hierarchical opcode decoder rather than flat tables,
// so instructions are simpler than the Z80 equivalent.
type Instruction struct {
	Name string // instruction mnemonic (uppercase)
}

// Instruction name constants sorted alphabetically.
const (
	ABCDName    = "ABCD"
	ADDName     = "ADD"
	ADDAName    = "ADDA"
	ADDIName    = "ADDI"
	ADDQName    = "ADDQ"
	ADDXName    = "ADDX"
	ANDName     = "AND"
	ANDIName    = "ANDI"
	ASLName     = "ASL"
	ASRName     = "ASR"
	BccName     = "Bcc"
	BCHGName    = "BCHG"
	BCLRName    = "BCLR"
	BRAName     = "BRA"
	BSETName    = "BSET"
	BSRName     = "BSR"
	BTSTName    = "BTST"
	CHKName     = "CHK"
	CLRName     = "CLR"
	CMPName     = "CMP"
	CMPAName    = "CMPA"
	CMPIName    = "CMPI"
	CMPMName    = "CMPM"
	DBccName    = "DBcc"
	DIVSName    = "DIVS"
	DIVUName    = "DIVU"
	EORName     = "EOR"
	EORIName    = "EORI"
	EXGName     = "EXG"
	EXTName     = "EXT"
	ILLEGALName = "ILLEGAL"
	JMPName     = "JMP"
	JSRName     = "JSR"
	LEAName     = "LEA"
	LINKName    = "LINK"
	LSLName     = "LSL"
	LSRName     = "LSR"
	MOVEName    = "MOVE"
	MOVEAName   = "MOVEA"
	MOVEMName   = "MOVEM"
	MOVEPName   = "MOVEP"
	MOVEQName   = "MOVEQ"
	MULSName    = "MULS"
	MULUName    = "MULU"
	NBCDName    = "NBCD"
	NEGName     = "NEG"
	NEGXName    = "NEGX"
	NOPName     = "NOP"
	NOTName     = "NOT"
	ORName      = "OR"
	ORIName     = "ORI"
	PEAName     = "PEA"
	RESETName   = "RESET"
	ROLName     = "ROL"
	RORName     = "ROR"
	ROXLName    = "ROXL"
	ROXRName    = "ROXR"
	RTEName     = "RTE"
	RTRName     = "RTR"
	RTSName     = "RTS"
	SBCDName    = "SBCD"
	SccName     = "Scc"
	STOPName    = "STOP"
	SUBName     = "SUB"
	SUBAName    = "SUBA"
	SUBIName    = "SUBI"
	SUBQName    = "SUBQ"
	SUBXName    = "SUBX"
	SWAPName    = "SWAP"
	TASName     = "TAS"
	TRAPName    = "TRAP"
	TRAPVName   = "TRAPV"
	TSTName     = "TST"
	UNLKName    = "UNLK"
)

// Instruction variable definitions.
var (
	insABCD    = &Instruction{Name: ABCDName}
	insADD     = &Instruction{Name: ADDName}
	insADDA    = &Instruction{Name: ADDAName}
	insADDI    = &Instruction{Name: ADDIName}
	insADDQ    = &Instruction{Name: ADDQName}
	insADDX    = &Instruction{Name: ADDXName}
	insAND     = &Instruction{Name: ANDName}
	insANDI    = &Instruction{Name: ANDIName}
	insASL     = &Instruction{Name: ASLName}
	insASR     = &Instruction{Name: ASRName}
	insBcc     = &Instruction{Name: BccName}
	insBCHG    = &Instruction{Name: BCHGName}
	insBCLR    = &Instruction{Name: BCLRName}
	insBRA     = &Instruction{Name: BRAName}
	insBSET    = &Instruction{Name: BSETName}
	insBSR     = &Instruction{Name: BSRName}
	insBTST    = &Instruction{Name: BTSTName}
	insCHK     = &Instruction{Name: CHKName}
	insCLR     = &Instruction{Name: CLRName}
	insCMP     = &Instruction{Name: CMPName}
	insCMPA    = &Instruction{Name: CMPAName}
	insCMPI    = &Instruction{Name: CMPIName}
	insCMPM    = &Instruction{Name: CMPMName}
	insDBcc    = &Instruction{Name: DBccName}
	insDIVS    = &Instruction{Name: DIVSName}
	insDIVU    = &Instruction{Name: DIVUName}
	insEOR     = &Instruction{Name: EORName}
	insEORI    = &Instruction{Name: EORIName}
	insEXG     = &Instruction{Name: EXGName}
	insEXT     = &Instruction{Name: EXTName}
	insILLEGAL = &Instruction{Name: ILLEGALName}
	insJMP     = &Instruction{Name: JMPName}
	insJSR     = &Instruction{Name: JSRName}
	insLEA     = &Instruction{Name: LEAName}
	insLINK    = &Instruction{Name: LINKName}
	insLSL     = &Instruction{Name: LSLName}
	insLSR     = &Instruction{Name: LSRName}
	insMOVE    = &Instruction{Name: MOVEName}
	insMOVEA   = &Instruction{Name: MOVEAName}
	insMOVEM   = &Instruction{Name: MOVEMName}
	insMOVEP   = &Instruction{Name: MOVEPName}
	insMOVEQ   = &Instruction{Name: MOVEQName}
	insMULS    = &Instruction{Name: MULSName}
	insMULU    = &Instruction{Name: MULUName}
	insNBCD    = &Instruction{Name: NBCDName}
	insNEG     = &Instruction{Name: NEGName}
	insNEGX    = &Instruction{Name: NEGXName}
	insNOP     = &Instruction{Name: NOPName}
	insNOT     = &Instruction{Name: NOTName}
	insOR      = &Instruction{Name: ORName}
	insORI     = &Instruction{Name: ORIName}
	insPEA     = &Instruction{Name: PEAName}
	insRESET   = &Instruction{Name: RESETName}
	insROL     = &Instruction{Name: ROLName}
	insROR     = &Instruction{Name: RORName}
	insROXL    = &Instruction{Name: ROXLName}
	insROXR    = &Instruction{Name: ROXRName}
	insRTE     = &Instruction{Name: RTEName}
	insRTR     = &Instruction{Name: RTRName}
	insRTS     = &Instruction{Name: RTSName}
	insSBCD    = &Instruction{Name: SBCDName}
	insScc     = &Instruction{Name: SccName}
	insSTOP    = &Instruction{Name: STOPName}
	insSUB     = &Instruction{Name: SUBName}
	insSUBA    = &Instruction{Name: SUBAName}
	insSUBI    = &Instruction{Name: SUBIName}
	insSUBQ    = &Instruction{Name: SUBQName}
	insSUBX    = &Instruction{Name: SUBXName}
	insSWAP    = &Instruction{Name: SWAPName}
	insTAS     = &Instruction{Name: TASName}
	insTRAP    = &Instruction{Name: TRAPName}
	insTRAPV   = &Instruction{Name: TRAPVName}
	insTST     = &Instruction{Name: TSTName}
	insUNLK    = &Instruction{Name: UNLKName}
)

// Instructions maps instruction names to their definitions.
var Instructions = map[string]*Instruction{
	ABCDName:    insABCD,
	ADDName:     insADD,
	ADDAName:    insADDA,
	ADDIName:    insADDI,
	ADDQName:    insADDQ,
	ADDXName:    insADDX,
	ANDName:     insAND,
	ANDIName:    insANDI,
	ASLName:     insASL,
	ASRName:     insASR,
	BccName:     insBcc,
	BCHGName:    insBCHG,
	BCLRName:    insBCLR,
	BRAName:     insBRA,
	BSETName:    insBSET,
	BSRName:     insBSR,
	BTSTName:    insBTST,
	CHKName:     insCHK,
	CLRName:     insCLR,
	CMPName:     insCMP,
	CMPAName:    insCMPA,
	CMPIName:    insCMPI,
	CMPMName:    insCMPM,
	DBccName:    insDBcc,
	DIVSName:    insDIVS,
	DIVUName:    insDIVU,
	EORName:     insEOR,
	EORIName:    insEORI,
	EXGName:     insEXG,
	EXTName:     insEXT,
	ILLEGALName: insILLEGAL,
	JMPName:     insJMP,
	JSRName:     insJSR,
	LEAName:     insLEA,
	LINKName:    insLINK,
	LSLName:     insLSL,
	LSRName:     insLSR,
	MOVEName:    insMOVE,
	MOVEAName:   insMOVEA,
	MOVEMName:   insMOVEM,
	MOVEPName:   insMOVEP,
	MOVEQName:   insMOVEQ,
	MULSName:    insMULS,
	MULUName:    insMULU,
	NBCDName:    insNBCD,
	NEGName:     insNEG,
	NEGXName:    insNEGX,
	NOPName:     insNOP,
	NOTName:     insNOT,
	ORName:      insOR,
	ORIName:     insORI,
	PEAName:     insPEA,
	RESETName:   insRESET,
	ROLName:     insROL,
	RORName:     insROR,
	ROXLName:    insROXL,
	ROXRName:    insROXR,
	RTEName:     insRTE,
	RTRName:     insRTR,
	RTSName:     insRTS,
	SBCDName:    insSBCD,
	SccName:     insScc,
	STOPName:    insSTOP,
	SUBName:     insSUB,
	SUBAName:    insSUBA,
	SUBIName:    insSUBI,
	SUBQName:    insSUBQ,
	SUBXName:    insSUBX,
	SWAPName:    insSWAP,
	TASName:     insTAS,
	TRAPName:    insTRAP,
	TRAPVName:   insTRAPV,
	TSTName:     insTST,
	UNLKName:    insUNLK,
}
