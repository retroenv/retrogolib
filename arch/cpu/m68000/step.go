package m68000

import "fmt"

// TraceStep contains all info needed to print a trace step.
type TraceStep struct {
	PC     uint32       // Program counter before instruction
	Opcode DecodedOpcode // Decoded opcode
	Words  []uint16     // Instruction words
}

// Step executes the next instruction in the CPU.
func (c *CPU) Step() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.halted {
		c.cycles += 4
		return nil
	}

	if c.stopped {
		c.cycles += 4
		if c.checkInterrupts() {
			c.stopped = false
		}
		return nil
	}

	// Check for pending interrupts.
	c.checkInterrupts()

	pcBefore := c.PC

	// Fetch and decode the opcode word.
	opcodeWord := c.readWord()

	decoded, err := decodeOpcode(opcodeWord)
	if err != nil {
		return fmt.Errorf("decoding opcode at PC=%06X: %w", pcBefore, err)
	}

	if decoded.Instruction == nil {
		return fmt.Errorf("%w: 0x%04X at PC=%06X", ErrUnsupportedOpcode, opcodeWord, pcBefore)
	}

	if c.opts.tracing {
		c.TraceStep = TraceStep{
			PC:     pcBefore,
			Opcode: decoded,
			Words:  []uint16{opcodeWord},
		}
	}

	c.cycles += uint64(decoded.Timing)

	// Execute the instruction.
	if err := c.executeInstruction(decoded); err != nil {
		return fmt.Errorf("executing %s at PC=%06X: %w", decoded.Instruction.Name, pcBefore, err)
	}

	// Check for trace exception.
	if c.sr&MaskTrace != 0 {
		if err := c.processException(VectorTrace); err != nil {
			return fmt.Errorf("processing trace exception: %w", err)
		}
	}

	return nil
}

// executeInstruction dispatches execution to the appropriate handler based on instruction name.
//
//nolint:cyclop
func (c *CPU) executeInstruction(d DecodedOpcode) error {
	ins := d.Instruction

	switch ins {
	case insABCD:
		return c.execABCD(d)
	case insADD:
		return c.execADD(d)
	case insADDA:
		return c.execADDA(d)
	case insADDI:
		return c.execADDI(d)
	case insADDQ:
		return c.execADDQ(d)
	case insADDX:
		return c.execADDX(d)
	case insAND:
		return c.execAND(d)
	case insANDI:
		return c.execANDI(d)
	case insASL:
		return c.execASL(d)
	case insASR:
		return c.execASR(d)
	case insBcc:
		return c.execBcc(d)
	case insBCHG:
		return c.execBCHG(d)
	case insBCLR:
		return c.execBCLR(d)
	case insBRA:
		return c.execBRA(d)
	case insBSET:
		return c.execBSET(d)
	case insBSR:
		return c.execBSR(d)
	case insBTST:
		return c.execBTST(d)
	case insCHK:
		return c.execCHK(d)
	case insCLR:
		return c.execCLR(d)
	case insCMP:
		return c.execCMP(d)
	case insCMPA:
		return c.execCMPA(d)
	case insCMPI:
		return c.execCMPI(d)
	case insCMPM:
		return c.execCMPM(d)
	case insDBcc:
		return c.execDBcc(d)
	case insDIVS:
		return c.execDIVS(d)
	case insDIVU:
		return c.execDIVU(d)
	case insEOR:
		return c.execEOR(d)
	case insEORI:
		return c.execEORI(d)
	case insEXG:
		return c.execEXG(d)
	case insEXT:
		return c.execEXT(d)
	case insILLEGAL:
		return c.execILLEGAL(d)
	case insJMP:
		return c.execJMP(d)
	case insJSR:
		return c.execJSR(d)
	case insLEA:
		return c.execLEA(d)
	case insLINK:
		return c.execLINK(d)
	case insLSL:
		return c.execLSL(d)
	case insLSR:
		return c.execLSR(d)
	case insMOVE:
		return c.execMOVE(d)
	case insMOVEA:
		return c.execMOVEA(d)
	case insMOVEM:
		return c.execMOVEM(d)
	case insMOVEP:
		return c.execMOVEP(d)
	case insMOVEQ:
		return c.execMOVEQ(d)
	case insMULS:
		return c.execMULS(d)
	case insMULU:
		return c.execMULU(d)
	case insNBCD:
		return c.execNBCD(d)
	case insNEG:
		return c.execNEG(d)
	case insNEGX:
		return c.execNEGX(d)
	case insNOP:
		return nil
	case insNOT:
		return c.execNOT(d)
	case insOR:
		return c.execOR(d)
	case insORI:
		return c.execORI(d)
	case insPEA:
		return c.execPEA(d)
	case insRESET:
		return c.execRESET(d)
	case insROL:
		return c.execROL(d)
	case insROR:
		return c.execROR(d)
	case insROXL:
		return c.execROXL(d)
	case insROXR:
		return c.execROXR(d)
	case insRTE:
		return c.execRTE(d)
	case insRTR:
		return c.execRTR(d)
	case insRTS:
		return c.execRTS(d)
	case insSBCD:
		return c.execSBCD(d)
	case insScc:
		return c.execScc(d)
	case insSTOP:
		return c.execSTOP(d)
	case insSUB:
		return c.execSUB(d)
	case insSUBA:
		return c.execSUBA(d)
	case insSUBI:
		return c.execSUBI(d)
	case insSUBQ:
		return c.execSUBQ(d)
	case insSUBX:
		return c.execSUBX(d)
	case insSWAP:
		return c.execSWAP(d)
	case insTAS:
		return c.execTAS(d)
	case insTRAP:
		return c.execTRAP(d)
	case insTRAPV:
		return c.execTRAPV(d)
	case insTST:
		return c.execTST(d)
	case insUNLK:
		return c.execUNLK(d)
	default:
		return fmt.Errorf("%w: %s", ErrUnimplemented, ins.Name)
	}
}
