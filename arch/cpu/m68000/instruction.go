package m68000

// Instruction defines a 68000 CPU instruction with its execution logic.
// The 68000 uses a hierarchical opcode decoder rather than flat tables,
// so instructions are simpler than the Z80 equivalent.
type Instruction struct {
	Name string   // instruction mnemonic (uppercase)
	exec execFunc // execution handler (nil for NOP)
}

// execFunc is the type of an instruction execution handler.
type execFunc func(c *CPU, d DecodedOpcode) error

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
	insABCD    = &Instruction{Name: ABCDName, exec: execABCD}
	insADD     = &Instruction{Name: ADDName, exec: execADD}
	insADDA    = &Instruction{Name: ADDAName, exec: execADDA}
	insADDI    = &Instruction{Name: ADDIName, exec: execADDI}
	insADDQ    = &Instruction{Name: ADDQName, exec: execADDQ}
	insADDX    = &Instruction{Name: ADDXName, exec: execADDX}
	insAND     = &Instruction{Name: ANDName, exec: execAND}
	insANDI    = &Instruction{Name: ANDIName, exec: execANDI}
	insASL     = &Instruction{Name: ASLName, exec: execASL}
	insASR     = &Instruction{Name: ASRName, exec: execASR}
	insBcc     = &Instruction{Name: BccName, exec: execBcc}
	insBCHG    = &Instruction{Name: BCHGName, exec: execBCHG}
	insBCLR    = &Instruction{Name: BCLRName, exec: execBCLR}
	insBRA     = &Instruction{Name: BRAName, exec: execBRA}
	insBSET    = &Instruction{Name: BSETName, exec: execBSET}
	insBSR     = &Instruction{Name: BSRName, exec: execBSR}
	insBTST    = &Instruction{Name: BTSTName, exec: execBTST}
	insCHK     = &Instruction{Name: CHKName, exec: execCHK}
	insCLR     = &Instruction{Name: CLRName, exec: execCLR}
	insCMP     = &Instruction{Name: CMPName, exec: execCMP}
	insCMPA    = &Instruction{Name: CMPAName, exec: execCMPA}
	insCMPI    = &Instruction{Name: CMPIName, exec: execCMPI}
	insCMPM    = &Instruction{Name: CMPMName, exec: execCMPM}
	insDBcc    = &Instruction{Name: DBccName, exec: execDBcc}
	insDIVS    = &Instruction{Name: DIVSName, exec: execDIVS}
	insDIVU    = &Instruction{Name: DIVUName, exec: execDIVU}
	insEOR     = &Instruction{Name: EORName, exec: execEOR}
	insEORI    = &Instruction{Name: EORIName, exec: execEORI}
	insEXG     = &Instruction{Name: EXGName, exec: execEXG}
	insEXT     = &Instruction{Name: EXTName, exec: execEXT}
	insILLEGAL = &Instruction{Name: ILLEGALName, exec: execILLEGAL}
	insJMP     = &Instruction{Name: JMPName, exec: execJMP}
	insJSR     = &Instruction{Name: JSRName, exec: execJSR}
	insLEA     = &Instruction{Name: LEAName, exec: execLEA}
	insLINK    = &Instruction{Name: LINKName, exec: execLINK}
	insLSL     = &Instruction{Name: LSLName, exec: execLSL}
	insLSR     = &Instruction{Name: LSRName, exec: execLSR}
	insMOVE    = &Instruction{Name: MOVEName, exec: execMOVE}
	insMOVEA   = &Instruction{Name: MOVEAName, exec: execMOVEA}
	insMOVEM   = &Instruction{Name: MOVEMName, exec: execMOVEM}
	insMOVEP   = &Instruction{Name: MOVEPName, exec: execMOVEP}
	insMOVEQ   = &Instruction{Name: MOVEQName, exec: execMOVEQ}
	insMULS    = &Instruction{Name: MULSName, exec: execMULS}
	insMULU    = &Instruction{Name: MULUName, exec: execMULU}
	insNBCD    = &Instruction{Name: NBCDName, exec: execNBCD}
	insNEG     = &Instruction{Name: NEGName, exec: execNEG}
	insNEGX    = &Instruction{Name: NEGXName, exec: execNEGX}
	insNOP     = &Instruction{Name: NOPName}
	insNOT     = &Instruction{Name: NOTName, exec: execNOT}
	insOR      = &Instruction{Name: ORName, exec: execOR}
	insORI     = &Instruction{Name: ORIName, exec: execORI}
	insPEA     = &Instruction{Name: PEAName, exec: execPEA}
	insRESET   = &Instruction{Name: RESETName, exec: execRESET}
	insROL     = &Instruction{Name: ROLName, exec: execROL}
	insROR     = &Instruction{Name: RORName, exec: execROR}
	insROXL    = &Instruction{Name: ROXLName, exec: execROXL}
	insROXR    = &Instruction{Name: ROXRName, exec: execROXR}
	insRTE     = &Instruction{Name: RTEName, exec: execRTE}
	insRTR     = &Instruction{Name: RTRName, exec: execRTR}
	insRTS     = &Instruction{Name: RTSName, exec: execRTS}
	insSBCD    = &Instruction{Name: SBCDName, exec: execSBCD}
	insScc     = &Instruction{Name: SccName, exec: execScc}
	insSTOP    = &Instruction{Name: STOPName, exec: execSTOP}
	insSUB     = &Instruction{Name: SUBName, exec: execSUB}
	insSUBA    = &Instruction{Name: SUBAName, exec: execSUBA}
	insSUBI    = &Instruction{Name: SUBIName, exec: execSUBI}
	insSUBQ    = &Instruction{Name: SUBQName, exec: execSUBQ}
	insSUBX    = &Instruction{Name: SUBXName, exec: execSUBX}
	insSWAP    = &Instruction{Name: SWAPName, exec: execSWAP}
	insTAS     = &Instruction{Name: TASName, exec: execTAS}
	insTRAP    = &Instruction{Name: TRAPName, exec: execTRAP}
	insTRAPV   = &Instruction{Name: TRAPVName, exec: execTRAPV}
	insTST     = &Instruction{Name: TSTName, exec: execTST}
	insUNLK    = &Instruction{Name: UNLKName, exec: execUNLK}
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
