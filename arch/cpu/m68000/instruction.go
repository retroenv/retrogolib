package m68000

// execFunc is the type of an instruction execution handler.
type execFunc func(c *CPU, d DecodedOpcode) error

// Instruction defines a 68000 CPU instruction with its execution logic.
// The 68000 uses a hierarchical opcode decoder rather than flat tables,
// so instructions are simpler than the Z80 equivalent.
type Instruction struct {
	Name string   // instruction mnemonic (uppercase)
	exec execFunc // execution handler (nil for NOP)
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
	insABCD    = &Instruction{Name: ABCDName, exec: (*CPU).execABCD}
	insADD     = &Instruction{Name: ADDName, exec: (*CPU).execADD}
	insADDA    = &Instruction{Name: ADDAName, exec: (*CPU).execADDA}
	insADDI    = &Instruction{Name: ADDIName, exec: (*CPU).execADDI}
	insADDQ    = &Instruction{Name: ADDQName, exec: (*CPU).execADDQ}
	insADDX    = &Instruction{Name: ADDXName, exec: (*CPU).execADDX}
	insAND     = &Instruction{Name: ANDName, exec: (*CPU).execAND}
	insANDI    = &Instruction{Name: ANDIName, exec: (*CPU).execANDI}
	insASL     = &Instruction{Name: ASLName, exec: (*CPU).execASL}
	insASR     = &Instruction{Name: ASRName, exec: (*CPU).execASR}
	insBcc     = &Instruction{Name: BccName, exec: (*CPU).execBcc}
	insBCHG    = &Instruction{Name: BCHGName, exec: (*CPU).execBCHG}
	insBCLR    = &Instruction{Name: BCLRName, exec: (*CPU).execBCLR}
	insBRA     = &Instruction{Name: BRAName, exec: (*CPU).execBRA}
	insBSET    = &Instruction{Name: BSETName, exec: (*CPU).execBSET}
	insBSR     = &Instruction{Name: BSRName, exec: (*CPU).execBSR}
	insBTST    = &Instruction{Name: BTSTName, exec: (*CPU).execBTST}
	insCHK     = &Instruction{Name: CHKName, exec: (*CPU).execCHK}
	insCLR     = &Instruction{Name: CLRName, exec: (*CPU).execCLR}
	insCMP     = &Instruction{Name: CMPName, exec: (*CPU).execCMP}
	insCMPA    = &Instruction{Name: CMPAName, exec: (*CPU).execCMPA}
	insCMPI    = &Instruction{Name: CMPIName, exec: (*CPU).execCMPI}
	insCMPM    = &Instruction{Name: CMPMName, exec: (*CPU).execCMPM}
	insDBcc    = &Instruction{Name: DBccName, exec: (*CPU).execDBcc}
	insDIVS    = &Instruction{Name: DIVSName, exec: (*CPU).execDIVS}
	insDIVU    = &Instruction{Name: DIVUName, exec: (*CPU).execDIVU}
	insEOR     = &Instruction{Name: EORName, exec: (*CPU).execEOR}
	insEORI    = &Instruction{Name: EORIName, exec: (*CPU).execEORI}
	insEXG     = &Instruction{Name: EXGName, exec: (*CPU).execEXG}
	insEXT     = &Instruction{Name: EXTName, exec: (*CPU).execEXT}
	insILLEGAL = &Instruction{Name: ILLEGALName, exec: (*CPU).execILLEGAL}
	insJMP     = &Instruction{Name: JMPName, exec: (*CPU).execJMP}
	insJSR     = &Instruction{Name: JSRName, exec: (*CPU).execJSR}
	insLEA     = &Instruction{Name: LEAName, exec: (*CPU).execLEA}
	insLINK    = &Instruction{Name: LINKName, exec: (*CPU).execLINK}
	insLSL     = &Instruction{Name: LSLName, exec: (*CPU).execLSL}
	insLSR     = &Instruction{Name: LSRName, exec: (*CPU).execLSR}
	insMOVE    = &Instruction{Name: MOVEName, exec: (*CPU).execMOVE}
	insMOVEA   = &Instruction{Name: MOVEAName, exec: (*CPU).execMOVEA}
	insMOVEM   = &Instruction{Name: MOVEMName, exec: (*CPU).execMOVEM}
	insMOVEP   = &Instruction{Name: MOVEPName, exec: (*CPU).execMOVEP}
	insMOVEQ   = &Instruction{Name: MOVEQName, exec: (*CPU).execMOVEQ}
	insMULS    = &Instruction{Name: MULSName, exec: (*CPU).execMULS}
	insMULU    = &Instruction{Name: MULUName, exec: (*CPU).execMULU}
	insNBCD    = &Instruction{Name: NBCDName, exec: (*CPU).execNBCD}
	insNEG     = &Instruction{Name: NEGName, exec: (*CPU).execNEG}
	insNEGX    = &Instruction{Name: NEGXName, exec: (*CPU).execNEGX}
	insNOP     = &Instruction{Name: NOPName}
	insNOT     = &Instruction{Name: NOTName, exec: (*CPU).execNOT}
	insOR      = &Instruction{Name: ORName, exec: (*CPU).execOR}
	insORI     = &Instruction{Name: ORIName, exec: (*CPU).execORI}
	insPEA     = &Instruction{Name: PEAName, exec: (*CPU).execPEA}
	insRESET   = &Instruction{Name: RESETName, exec: (*CPU).execRESET}
	insROL     = &Instruction{Name: ROLName, exec: (*CPU).execROL}
	insROR     = &Instruction{Name: RORName, exec: (*CPU).execROR}
	insROXL    = &Instruction{Name: ROXLName, exec: (*CPU).execROXL}
	insROXR    = &Instruction{Name: ROXRName, exec: (*CPU).execROXR}
	insRTE     = &Instruction{Name: RTEName, exec: (*CPU).execRTE}
	insRTR     = &Instruction{Name: RTRName, exec: (*CPU).execRTR}
	insRTS     = &Instruction{Name: RTSName, exec: (*CPU).execRTS}
	insSBCD    = &Instruction{Name: SBCDName, exec: (*CPU).execSBCD}
	insScc     = &Instruction{Name: SccName, exec: (*CPU).execScc}
	insSTOP    = &Instruction{Name: STOPName, exec: (*CPU).execSTOP}
	insSUB     = &Instruction{Name: SUBName, exec: (*CPU).execSUB}
	insSUBA    = &Instruction{Name: SUBAName, exec: (*CPU).execSUBA}
	insSUBI    = &Instruction{Name: SUBIName, exec: (*CPU).execSUBI}
	insSUBQ    = &Instruction{Name: SUBQName, exec: (*CPU).execSUBQ}
	insSUBX    = &Instruction{Name: SUBXName, exec: (*CPU).execSUBX}
	insSWAP    = &Instruction{Name: SWAPName, exec: (*CPU).execSWAP}
	insTAS     = &Instruction{Name: TASName, exec: (*CPU).execTAS}
	insTRAP    = &Instruction{Name: TRAPName, exec: (*CPU).execTRAP}
	insTRAPV   = &Instruction{Name: TRAPVName, exec: (*CPU).execTRAPV}
	insTST     = &Instruction{Name: TSTName, exec: (*CPU).execTST}
	insUNLK    = &Instruction{Name: UNLKName, exec: (*CPU).execUNLK}
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
