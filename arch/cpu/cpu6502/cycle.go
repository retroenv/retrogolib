package cpu6502

import "fmt"

// BusCycle describes the bus access that the CPU will attempt next.
type BusCycle struct {
	// Address is the address on the CPU bus.
	Address uint16
	// Value is the byte to write. It is zero for a read.
	Value byte
	// Write selects a write cycle. A false value selects a read cycle.
	Write bool
}

// CycleHook clocks external devices before a CPU bus access. Return true to
// take a read cycle for DMA. The CPU retries that read on the next cycle. A
// write cannot be stopped. The hook must perform any DMA access itself.
// Cycles returns the number of completed cycles while the hook runs. The hook
// must not call Step, Reset, or CheckInterrupts on this CPU.
type CycleHook func(BusCycle) bool

func (c *CPU) stepCycles() error {
	if c.opts.variant >= Variant65C02 {
		return fmt.Errorf("%w: bus-cycle execution requires an NMOS CPU", ErrUnsupportedAddressingMode)
	}
	if c.jammed {
		address := uint16(0xffff)
		if c.jamCycle == 1 || c.jamCycle == 2 {
			address = 0xfffe
		}
		c.jamCycle = min(c.jamCycle+1, 3)
		c.readCycle(address)
		return nil
	}
	if c.stallCycles != 0 {
		c.stallCycles--
		c.readCycle(c.PC)
		return nil
	}

	c.cycleActive = true
	defer func() { c.cycleActive = false }()
	c.executionPC = c.PC
	c.executionCycle = c.cycles
	c.interrupt = InterruptNone
	c.branchTaken = false
	c.opcode = c.readCycle(c.PC)
	opcode := Opcodes[c.opcode]
	if opcode.Instruction == nil {
		return fmt.Errorf("%w: 0x%02x at PC=0x%04x", ErrUnknownOpcode, c.opcode, c.PC)
	}
	if c.opts.tracing {
		c.TraceStep = TraceStep{
			PC:             c.PC,
			Opcode:         opcode,
			OpcodeOperands: []byte{c.opcode},
		}
	}
	return c.executeCycles(opcode)
}

func (c *CPU) executeCycles(opcode Opcode) error {
	ins := opcode.Instruction
	if ins == JsrInst || ins == BrkInst || ins == RtsInst || ins == RtiInst {
		return c.controlCycles(ins)
	}
	switch opcode.Addressing {
	case ImpliedAddressing, AccumulatorAddressing:
		c.readCycle(c.PC + 1)
		return c.impliedCycles(opcode)
	case ImmediateAddressing:
		value := c.readOperand(1)
		c.observeCycles(ins, int(value))
		if err := ins.ParamFunc(c, value); err != nil {
			return err
		}
		c.PC += 2
		return nil
	case RelativeAddressing:
		return c.branchCycles(ins)
	default:
		return c.memoryCycles(opcode)
	}
}

func (c *CPU) impliedCycles(opcode Opcode) error {
	ins := opcode.Instruction
	if ins == KilInst {
		c.jammed = true
		c.jamCycle = 0
	}
	if ins == PlaInst || ins == PlpInst {
		c.readCycle(StackBase + uint16(c.SP))
	}
	c.observeCycles(ins)
	var err error
	switch {
	case ins.NoParamFunc != nil:
		err = ins.NoParamFunc(c)
	case opcode.Addressing == AccumulatorAddressing:
		err = ins.ParamFunc(c, Accumulator(0))
	default:
		err = ins.ParamFunc(c)
	}
	c.PC++
	return err
}

func (c *CPU) branchCycles(ins *Instruction) error {
	offset := int8(c.readOperand(1))
	irq, nmi := c.irqPolled, c.nmiPolled
	next := c.PC + 2
	target := uint16(int32(next) + int32(offset))
	c.observeCycles(ins, Absolute(target))
	if err := ins.ParamFunc(c, Absolute(target)); err != nil {
		return err
	}
	if !c.branchTaken {
		c.PC = next
		return nil
	}
	c.readCycle(next)
	if next&0xff00 != target&0xff00 {
		c.readCycle(next&0xff00 | target&0xff)
		irq = irq || c.irqPolled
		nmi = nmi || c.nmiPolled
	}
	// A taken branch polls after its opcode fetch and before a page fixup.
	// The extra cycle within the same page does not add a polling point.
	c.irqPolled, c.nmiPolled = irq, nmi
	return nil
}

func (c *CPU) memoryCycles(opcode Opcode) error {
	access := instructionAccess(opcode.Instruction)
	base, address, err := c.operandAddressCycles(opcode.Addressing, access)
	if err != nil {
		return err
	}
	ins := opcode.Instruction
	c.observeAddressCycles(opcode, base, address)
	if ins == JmpInst {
		if opcode.Addressing == IndirectAddressing {
			address = c.readWordCycles(address, true)
		}
		c.PC = address
		return nil
	}
	if ins == ShaInst || ins == ShxInst || ins == ShyInst || ins == TasInst {
		register := &c.Y
		if opcode.Addressing == AbsoluteXAddressing {
			register = &c.X
		}
		err = ins.ParamFunc(c, Absolute(base), register)
	} else {
		err = c.dataCycles(ins, access, address)
	}
	c.PC += uint16(addressingModeSize(opcode.Addressing))
	return err
}

// dataCycles separates the bus transfers from the arithmetic. A compound
// unofficial instruction uses the result latch for its second operation.
func (c *CPU) dataCycles(ins *Instruction, access memoryAccess, address uint16) error {
	var value byte
	if access != writeAccess {
		value = c.readCycle(address)
	}
	if access == modifyAccess {
		c.writeCycle(address, value)
	}
	if err := ins.ParamFunc(c, &value); err != nil {
		return err
	}
	if access != readAccess {
		c.writeCycle(address, value)
	}
	return nil
}

func (c *CPU) operandAddressCycles(mode AddressingMode, access memoryAccess) (base, address uint16, err error) {
	low := c.readOperand(1)
	switch mode {
	case ZeroPageAddressing:
		base = uint16(low)
		return base, base, nil
	case ZeroPageXAddressing, ZeroPageYAddressing:
		c.readCycle(uint16(low))
		index := c.X
		if mode == ZeroPageYAddressing {
			index = c.Y
		}
		return uint16(low), uint16(low + index), nil
	case IndirectXAddressing:
		c.readCycle(uint16(low))
		base = c.readWordCycles(uint16(low+c.X), true)
		return base, base, nil
	case IndirectYAddressing:
		base = c.readWordCycles(uint16(low), true)
		return base, c.indexedAddressCycles(base, c.Y, access), nil
	case AbsoluteAddressing, IndirectAddressing, AbsoluteXAddressing, AbsoluteYAddressing:
		base = uint16(low) | uint16(c.readOperand(2))<<8
		if mode == AbsoluteXAddressing {
			return base, c.indexedAddressCycles(base, c.X, access), nil
		}
		if mode == AbsoluteYAddressing {
			return base, c.indexedAddressCycles(base, c.Y, access), nil
		}
		return base, base, nil
	default:
		return 0, 0, fmt.Errorf("%w: bus-cycle mode 0x%02x", ErrUnsupportedAddressingMode, mode)
	}
}

func (c *CPU) indexedAddressCycles(base uint16, index byte, access memoryAccess) uint16 {
	address := base + uint16(index)
	crossed := base&0xff00 != address&0xff00
	if c.opts.tracing {
		c.TraceStep.PageCrossed = crossed
	}
	if crossed || access != readAccess {
		c.readCycle(base&0xff00 | address&0xff)
	}
	return address
}

func (c *CPU) controlCycles(ins *Instruction) error {
	c.observeCycles(ins)
	if ins == JsrInst {
		low := c.readOperand(1)
		c.readCycle(StackBase + uint16(c.SP))
		c.push16(c.PC + 2)
		high := c.readOperand(2)
		c.PC = uint16(low) | uint16(high)<<8
		return nil
	}
	c.readCycle(c.PC + 1)
	if ins == BrkInst {
		c.interrupt = InterruptBRK
		c.irqRunning = true
		c.push16(c.PC + 2)
		vector := c.selectInterruptVector(IrqAddress)
		c.push(c.GetFlags() | 0x30)
		c.Flags.I = 1
		c.PC = c.readWordCycles(vector, false)
		c.irqPolled, c.nmiPolled = false, false
		return nil
	}
	c.readCycle(StackBase + uint16(c.SP))
	if err := ins.NoParamFunc(c); err != nil {
		return err
	}
	if ins == RtsInst {
		c.readCycle(c.PC - 1)
	}
	return nil
}

func (c *CPU) interruptCycles(vector uint16, source InterruptSource) {
	c.cycleActive = true
	defer func() { c.cycleActive = false }()
	c.executionPC = c.PC
	c.executionCycle = c.cycles
	c.opcode = 0
	c.interrupt = source
	c.readCycle(c.PC)
	c.readCycle(c.PC)
	c.push16(c.PC)
	vector = c.selectInterruptVector(vector)
	c.push(c.GetFlags()&^0x10 | 0x20)
	c.Flags.I = 1
	c.PC = c.readWordCycles(vector, false)
	c.irqPolled, c.nmiPolled = false, false
}

// An NMI detected before the status push can replace the IRQ or BRK vector.
// The stacked PC and status still come from the original interrupt sequence.
func (c *CPU) selectInterruptVector(vector uint16) uint16 {
	if c.triggerNmi {
		c.triggerNmi = false
		c.nmiRunning = true
		return NMIAddress
	}
	return vector
}

func (c *CPU) resetCycles() {
	c.jammed = false
	c.jamCycle = 0
	c.triggerIrq = false
	c.triggerNmi = false
	c.irqRunning = false
	c.nmiRunning = false
	c.stallCycles = 0
	c.Flags.I = 1
	c.readCycle(c.PC)
	c.readCycle(c.PC)
	for range 3 {
		c.readCycle(StackBase + uint16(c.SP))
		c.SP--
	}
	c.PC = c.readWordCycles(ResetAddress, false)
	c.irqPolled, c.nmiPolled = false, false
}

func (c *CPU) readOperand(offset uint16) byte {
	value := c.readCycle(c.PC + offset)
	if c.opts.tracing {
		c.TraceStep.OpcodeOperands = append(c.TraceStep.OpcodeOperands, value)
	}
	return value
}

func (c *CPU) readCycle(address uint16) byte {
	for {
		c.irqPolled, c.nmiPolled = c.irqSample, c.nmiSample
		held := c.opts.cycleHook(BusCycle{Address: address})
		c.cycles++
		if !held {
			value := c.memory.Read(address)
			c.sampleInterrupts()
			return value
		}
		c.sampleInterrupts()
	}
}

func (c *CPU) writeCycle(address uint16, value byte) {
	c.irqPolled, c.nmiPolled = c.irqSample, c.nmiSample
	c.opts.cycleHook(BusCycle{
		Address: address,
		Value:   value,
		Write:   true,
	})
	c.cycles++
	c.memory.Write(address, value)
	c.sampleInterrupts()
}

// sampleInterrupts records the inputs for the next cycle. Most instructions
// use the sample from their second-to-last cycle. CLI, SEI, and PLP change I
// after this sample; RTI changes I before it.
// https://www.nesdev.org/wiki/CPU_interrupts
func (c *CPU) sampleInterrupts() {
	c.irqSample = (c.triggerIrq || c.irqLine) && c.Flags.I == 0
	c.nmiSample = c.triggerNmi
}

func (c *CPU) readWordCycles(address uint16, wrap bool) uint16 {
	low := c.readCycle(address)
	next := address + 1
	if wrap {
		next = address&0xff00 | uint16(byte(next))
	}
	return uint16(low) | uint16(c.readCycle(next))<<8
}

func (c *CPU) observeCycles(ins *Instruction, params ...any) {
	if c.opts.preExecutionHook != nil {
		c.opts.preExecutionHook(c, ins, params...)
	}
}

func (c *CPU) observeAddressCycles(opcode Opcode, base, address uint16) {
	if c.opts.preExecutionHook == nil {
		return
	}
	switch opcode.Addressing {
	case AbsoluteXAddressing:
		c.observeCycles(opcode.Instruction, Absolute(base), &c.X)
	case AbsoluteYAddressing:
		c.observeCycles(opcode.Instruction, Absolute(base), &c.Y)
	case ZeroPageXAddressing:
		c.observeCycles(opcode.Instruction, ZeroPage(base), &c.X)
	case ZeroPageYAddressing:
		c.observeCycles(opcode.Instruction, ZeroPage(base), &c.Y)
	case IndirectXAddressing:
		c.observeCycles(opcode.Instruction, IndirectResolved(address), &c.X)
	case IndirectYAddressing:
		c.observeCycles(opcode.Instruction, IndirectResolved(address), &c.Y)
	case IndirectAddressing:
		c.observeCycles(opcode.Instruction, Indirect(base))
	default:
		c.observeCycles(opcode.Instruction, Absolute(base))
	}
}

type memoryAccess byte

const (
	readAccess memoryAccess = iota
	writeAccess
	modifyAccess
)

func instructionAccess(ins *Instruction) memoryAccess {
	switch ins {
	case AslInst, LsrInst, RolInst, RorInst, IncInst, DecInst, SloInst, SreInst, RlaInst, RraInst, DcpInst, IscInst:
		return modifyAccess
	case StaInst, StxInst, StyInst, SaxInst, ShaInst, ShxInst, ShyInst, TasInst:
		return writeAccess
	default:
		return readAccess
	}
}
