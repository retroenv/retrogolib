package cpu68000

import (
	"errors"
	"fmt"
)

// TraceStep contains all info needed to print a trace step.
type TraceStep struct {
	PC     uint32        // Program counter before instruction
	Opcode DecodedOpcode // Decoded opcode
	Words  []uint16      // Instruction words
}

// Step services a pending interrupt, idles while stopped, or executes one instruction.
func (cpu *CPU) Step() error {
	cpu.mu.Lock()
	defer cpu.mu.Unlock()
	defer cpu.synchronizeStackPointer()

	cpu.stepCycles = cpu.cycles
	cpu.instructionPC = cpu.PC
	cpu.exceptionAccess, cpu.exceptionRaised = false, false
	cpu.operandPCOffset = -2
	cpu.accessCycles = 0
	err := catchAccessFault(cpu.executeStep)
	var fault *accessError
	if errors.As(err, &fault) {
		return cpu.acceptAccessFault(fault)
	}
	return err
}

func (cpu *CPU) synchronizeStackPointer() {
	if cpu.IsSupervisor() {
		cpu.SSP = cpu.sp
	} else {
		cpu.USP = cpu.sp
	}
}

func (cpu *CPU) executeStep() error {
	if cpu.halted {
		cpu.cycles += 4
		return nil
	}

	if cpu.checkInterrupts() {
		return nil
	}
	if cpu.stopped {
		cpu.cycles += 4
		return nil
	}

	pcBefore := cpu.PC

	// Fetch and decode the opcode word.
	opcodeWord := cpu.readWord()
	cpu.instructionWord = opcodeWord
	cpu.accessCycles = 0
	trace := cpu.sr&MaskTrace != 0

	decoded, err := decodeOpcode(opcodeWord)
	if err != nil {
		return fmt.Errorf("decoding opcode at PC=%06X: %w", pcBefore, err)
	}

	if decoded.Instruction == nil {
		return fmt.Errorf("%w: 0x%04X at PC=%06X", ErrUnsupportedOpcode, opcodeWord, pcBefore)
	}

	if cpu.opts.tracing {
		cpu.TraceStep = TraceStep{
			PC:     pcBefore,
			Opcode: decoded,
			Words:  []uint16{opcodeWord},
		}
	}

	cpu.cycles += cpu.instructionCycles(decoded)

	// Execute the instruction via its handler.
	ins := decoded.Instruction
	if ins.exec != nil {
		if err := ins.exec(cpu, decoded); err != nil {
			return fmt.Errorf("executing %s at PC=%06X: %w", ins.Name, pcBefore, err)
		}
	}
	cpu.checkInstructionAddress()

	// Check for trace exception.
	if trace && !cpu.exceptionRaised {
		if err := cpu.processException(VectorTrace); err != nil {
			return fmt.Errorf("processing trace exception: %w", err)
		}
	}

	return nil
}

func (cpu *CPU) checkInstructionAddress() {
	if cpu.PC&1 != 0 {
		// The next instruction prefetch faults before this step retires.
		target := cpu.PC
		cpu.PC -= 4
		cpu.exceptionAccess = true
		cpu.raiseAccessFault(target, false, programSpace, VectorAddressErr, ErrAddressError)
	}
}
