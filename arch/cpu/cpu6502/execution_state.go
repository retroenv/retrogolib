package cpu6502

import "errors"

// ExecutionState contains CPU state at an instruction boundary.
// Memory, hooks, and the configured CPU variant belong to the caller.
type ExecutionState struct {
	State State

	StallCycles uint16

	JamCycle byte
	Jammed   bool

	IRQLine    bool
	TriggerIRQ bool
	TriggerNMI bool

	IRQPolled bool
	IRQSample bool
	NMIPolled bool
	NMISample bool

	BranchTaken    bool
	ExecutionCycle uint64
	ExecutionPC    uint16
	Interrupt      InterruptSource
	Opcode         byte
}

// SaveExecutionState captures the CPU. Stop instruction execution first.
func (c *CPU) SaveExecutionState() ExecutionState {
	c.mu.RLock()
	defer c.mu.RUnlock()
	state := c.stateLocked()

	return ExecutionState{
		State: state,

		StallCycles: c.stallCycles,

		JamCycle: c.jamCycle,
		Jammed:   c.jammed,

		IRQLine:    c.irqLine,
		TriggerIRQ: c.triggerIrq,
		TriggerNMI: c.triggerNmi,

		IRQPolled: c.irqPolled,
		IRQSample: c.irqSample,
		NMIPolled: c.nmiPolled,
		NMISample: c.nmiSample,

		BranchTaken:    c.branchTaken,
		ExecutionCycle: c.executionCycle,
		ExecutionPC:    c.executionPC,
		Interrupt:      c.interrupt,
		Opcode:         c.opcode,
	}
}

// RestoreExecutionState restores a boundary snapshot without reading memory.
// Stop instruction execution and external interrupt requests first.
func (c *CPU) RestoreExecutionState(s ExecutionState) error {
	f := s.State.Flags
	if c.cycleActive || s.JamCycle > 3 || f.C > 1 || f.Z > 1 || f.I > 1 || f.D > 1 ||
		f.B > 1 || f.U > 1 || f.V > 1 || f.N > 1 {

		return errors.New("invalid CPU execution state")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.A = s.State.A
	c.X = s.State.X
	c.Y = s.State.Y
	c.PC = s.State.PC
	c.SP = s.State.SP
	c.cycles = s.State.Cycles
	c.Flags = s.State.Flags
	c.irqRunning = s.State.Interrupts.IrqRunning
	c.nmiRunning = s.State.Interrupts.NMIRunning

	c.stallCycles = s.StallCycles

	c.jamCycle = s.JamCycle
	c.jammed = s.Jammed

	c.irqLine = s.IRQLine
	c.triggerIrq = s.TriggerIRQ
	c.triggerNmi = s.TriggerNMI

	c.irqPolled = s.IRQPolled
	c.irqSample = s.IRQSample
	c.nmiPolled = s.NMIPolled
	c.nmiSample = s.NMISample

	c.branchTaken = s.BranchTaken
	c.executionCycle = s.ExecutionCycle
	c.executionPC = s.ExecutionPC
	c.interrupt = s.Interrupt
	c.opcode = s.Opcode

	c.TraceStep = TraceStep{}

	return nil
}

// Jammed reports whether KIL has stopped instruction execution.
// Call it from the execution goroutine or while execution is stopped.
func (c *CPU) Jammed() bool { return c.jammed }
