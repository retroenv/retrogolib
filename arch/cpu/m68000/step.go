package m68000

import "fmt"

// TraceStep contains all info needed to print a trace step.
type TraceStep struct {
	PC     uint32        // Program counter before instruction
	Opcode DecodedOpcode // Decoded opcode
	Words  []uint16      // Instruction words
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

	// Execute the instruction via its handler.
	ins := decoded.Instruction
	if ins.exec != nil {
		if err := ins.exec(c, decoded); err != nil {
			return fmt.Errorf("executing %s at PC=%06X: %w", ins.Name, pcBefore, err)
		}
	}

	// Check for trace exception.
	if c.sr&MaskTrace != 0 {
		if err := c.processException(VectorTrace); err != nil {
			return fmt.Errorf("processing trace exception: %w", err)
		}
	}

	return nil
}
