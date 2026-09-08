package cpu68000

// Exception vector numbers.
const (
	VectorResetSSP    = 0  // Reset: initial SSP
	VectorResetPC     = 1  // Reset: initial PC
	VectorBusError    = 2  // Bus error
	VectorAddressErr  = 3  // Address error
	VectorIllegal     = 4  // Illegal instruction
	VectorDivZero     = 5  // Divide by zero
	VectorCHK         = 6  // CHK instruction
	VectorTRAPV       = 7  // TRAPV instruction
	VectorPrivilege   = 8  // Privilege violation
	VectorTrace       = 9  // Trace
	VectorLineA       = 10 // Line 1010 emulator
	VectorLineF       = 11 // Line 1111 emulator
	VectorSpurious    = 24 // Spurious interrupt
	VectorAutoVector1 = 25 // Autovector level 1
	VectorAutoVector2 = 26 // Autovector level 2
	VectorAutoVector3 = 27 // Autovector level 3
	VectorAutoVector4 = 28 // Autovector level 4
	VectorAutoVector5 = 29 // Autovector level 5
	VectorAutoVector6 = 30 // Autovector level 6
	VectorAutoVector7 = 31 // Autovector level 7
	VectorTrap0       = 32 // TRAP #0
	VectorTrap15      = 47 // TRAP #15
)

// TriggerIRQ queues an interrupt at level 1-7; level 7 is non-maskable.
// The highest requested level remains pending until accepted. Other values are ignored.
func (cpu *CPU) TriggerIRQ(level uint8) {
	if level == 0 || level > 7 {
		return
	}
	cpu.mu.Lock()
	defer cpu.mu.Unlock()
	cpu.pendingIRQ = max(cpu.pendingIRQ, level)
}

// processException processes an exception with the given vector number.
// It saves the current state and loads the new PC from the vector table.
func (cpu *CPU) processException(vector int) error {
	// Save current SR.
	oldSR := cpu.GetSR()
	switch vector {
	case VectorTrace:
		cpu.cycles += 34
	case VectorCHK:
		cpu.cycles += 30
	case VectorDivZero:
		cpu.cycles += 38
	default:
		cpu.cycles = cpu.stepCycles + 34
	}
	cpu.exceptionRaised = true
	cpu.exceptionAccess = vector == VectorIllegal || vector == VectorPrivilege || vector == VectorTrace ||
		vector == VectorLineA || vector == VectorLineF

	// Enter supervisor mode and clear trace.
	cpu.sr |= MaskSupervisor
	cpu.sr &^= MaskTrace

	// If switching from user to supervisor, swap stack pointers.
	if oldSR&MaskSupervisor == 0 {
		cpu.USP = cpu.sp
		cpu.sp = cpu.SSP
	}

	// Push PC and SR onto the supervisor stack.
	pc := cpu.PC
	if vector == VectorIllegal || vector == VectorPrivilege || vector == VectorLineA || vector == VectorLineF {
		pc = cpu.instructionPC
	}
	cpu.push32(pc)
	cpu.push16(oldSR)

	// Load new PC from vector table.
	vectorAddr := uint32(vector) * 4
	cpu.PC = vectorAddr
	cpu.operandPCOffset = 0
	cpu.PC = cpu.readBusLong(vectorAddr)

	cpu.stopped = false

	return nil
}

// processInterruptException processes an interrupt exception for the given level.
func (cpu *CPU) processInterruptException(level uint8) {
	// Save current SR.
	oldSR := cpu.GetSR()
	cpu.exceptionAccess, cpu.exceptionRaised = true, true

	// Enter supervisor mode, clear trace, set interrupt mask.
	cpu.sr |= MaskSupervisor
	cpu.sr &^= MaskTrace
	cpu.sr = (cpu.sr & ^uint16(MaskIPM)) | (uint16(level) << FlagIPM0)

	// If switching from user to supervisor, swap stack pointers.
	if oldSR&MaskSupervisor == 0 {
		cpu.USP = cpu.sp
		cpu.sp = cpu.SSP
	}

	// Push PC and SR.
	cpu.push32(cpu.PC)
	cpu.push16(oldSR)

	// Get vector from bus.
	vector := cpu.bus.IRQAcknowledge(level)

	// Load new PC from vector table.
	vectorAddr := vector * 4
	cpu.PC = vectorAddr
	cpu.PC = cpu.readBusLong(vectorAddr)

	cpu.stopped = false
	cpu.cycles += 44
}

// checkInterrupts checks for pending interrupts and processes them.
// Returns true if an interrupt was processed.
func (cpu *CPU) checkInterrupts() bool {
	level := max(cpu.bus.IRQLevel(), cpu.pendingIRQ)
	if level == 0 {
		return false
	}

	mask := cpu.InterruptMask()

	// Level 7 is non-maskable. Other levels must be higher than mask.
	if level < 7 && level <= mask {
		return false
	}

	if level == cpu.pendingIRQ {
		cpu.pendingIRQ = 0
	}
	cpu.processInterruptException(level)
	return true
}

func (cpu *CPU) acceptAccessFault(fault *accessError) error {
	err := catchAccessFault(func() error {
		cpu.SetSR((cpu.GetSR() | MaskSupervisor) &^ MaskTrace)
		cpu.exceptionAccess, cpu.exceptionRaised = true, true
		cpu.stopped = false
		cpu.cycles = cpu.stepCycles + fault.cycles + 50
		// The original 68000 has no format word. Its seven-word fault frame
		// starts with SSW, address, IR, followed by the ordinary SR/PC frame.
		cpu.push32(fault.pc)
		cpu.push16(fault.sr)
		cpu.push16(fault.word)
		cpu.push32(fault.address)
		cpu.push16(fault.status)
		cpu.PC = cpu.readBusLong(uint32(fault.vector) * 4)
		if cpu.PC&1 != 0 {
			cpu.raiseAccessFault(cpu.PC, false, programSpace, VectorAddressErr, ErrAddressError)
		}
		return nil
	})
	if err != nil {
		cpu.halted, cpu.faultHalted = true, true
	}
	return nil
}
