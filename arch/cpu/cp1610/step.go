package cp1610

import "fmt"

// Step executes one instruction or accepts a pending interrupt.
func (c *CPU) Step() error {
	if c.halted {
		return nil
	}
	if c.serviceIRQ() {
		return nil
	}

	pc := c.R[ProgramCounter]
	word := c.memory.Read(pc) & 0x03FF
	op, ok := GetOpcodeInfo(word)
	if !ok {
		return fmt.Errorf("%w: %03x at %04x", ErrInvalidOpcode, word, pc)
	}
	c.R[ProgramCounter]++
	doubleByte := c.Flags.D
	c.Flags.D = false
	c.lastInterruptible = op.Interruptible
	c.cycles += uint64(op.Cycles)

	switch {
	case word < 8:
		return c.executeControl(word)
	case word < 0x30:
		c.executeUnary(op)
	case word < 0x40:
		c.executeStatus(op)
	case word < 0x80:
		c.executeShift(op)
	case word < 0x200:
		c.executeRegisterPair(op)
	case word < 0x240:
		c.executeBranch(op)
	default:
		c.executeMemory(op, doubleByte)
	}
	return nil
}

func (c *CPU) executeControl(word uint16) error {
	switch word {
	case 0:
		c.halted = true
	case 1:
		c.Flags.D = true
	case 2:
		c.Flags.I = true
	case 3:
		c.Flags.I = false
	case 4:
		return c.executeJump()
	case 6:
		c.Flags.C = false
	case 7:
		c.Flags.C = true
	}
	return nil
}

func (c *CPU) executeJump() error {
	second := c.memory.Read(c.R[ProgramCounter]) & 0x03FF
	third := c.memory.Read(c.R[ProgramCounter]+1) & 0x03FF
	mode := second & 3
	if mode == 3 {
		return fmt.Errorf("%w: jump interrupt mode %d", ErrInvalidOpcode, mode)
	}
	returnAddress := c.R[ProgramCounter] + 2
	register := uint8((second>>8)&3) + 4
	if register != ProgramCounter {
		c.R[register] = returnAddress
	}
	if mode != 0 {
		c.Flags.I = mode == 1
	}
	c.R[ProgramCounter] = (second&0xFC)<<8 | third
	return nil
}

func (c *CPU) executeUnary(op Opcode) {
	reg := op.Destination
	value := c.R[reg]
	switch op.Instruction.Name {
	case IncrName:
		value++
		c.setSZ(value)

	case DecrName:
		value--
		c.setSZ(value)

	case ComrName:
		value = ^value
		c.setSZ(value)

	case NegrName:
		value = c.subtract(value, 0)

	case AdcrName:
		carry := uint16(0)
		if c.Flags.C {
			carry = 1
		}
		value = c.add(value, carry)
	}
	c.R[reg] = value
}

func (c *CPU) executeStatus(op Opcode) {
	switch op.Instruction.Name {
	case GswdName:
		c.R[op.Destination] = c.statusWord()
	case RswdName:
		c.setStatusWord(c.R[op.Source])
	}
}

func (c *CPU) executeRegisterPair(op Opcode) {
	left := c.R[op.Source]
	right := c.R[op.Destination]
	result := right
	switch op.Instruction.Name {
	case MovrName:
		result = left
		c.setSZ(result)

	case AddrName:
		result = c.add(left, right)

	case SubrName:
		result = c.subtract(left, right)

	case CmprName:
		c.subtract(left, right)
		return

	case AndrName:
		result = left & right
		c.setSZ(result)

	case XorrName:
		result = left ^ right
		c.setSZ(result)
	}
	c.R[op.Destination] = result
}

func (c *CPU) executeBranch(op Opcode) {
	displacement := c.memory.Read(c.R[ProgramCounter])
	c.R[ProgramCounter]++
	condition := op.Source
	taken := c.branchCondition(condition)
	if !taken {
		return
	}
	if op.Variant == 0 {
		c.R[ProgramCounter] += displacement
	} else {
		c.R[ProgramCounter] -= displacement + 1
	}
	c.cycles += 2
}

func (c *CPU) branchCondition(condition uint8) bool {
	if condition&0x10 != 0 {
		return false // EBCI is not connected on the Intellivision.
	}
	var taken bool
	switch condition & 7 {
	case 0:
		taken = true
	case 1:
		taken = c.Flags.C
	case 2:
		taken = c.Flags.O
	case 3:
		taken = !c.Flags.S
	case 4:
		taken = c.Flags.Z
	case 5:
		taken = c.Flags.S != c.Flags.O
	case 6:
		taken = c.Flags.Z || c.Flags.S != c.Flags.O
	case 7:
		taken = c.Flags.S != c.Flags.C
	}
	if condition&8 != 0 {
		return !taken
	}
	return taken
}

func (c *CPU) executeMemory(op Opcode, doubleByte bool) {
	if op.Instruction.Name == MvoName {
		c.storeOperand(op)
		return
	}
	value := c.readOperand(op, doubleByte)
	dest := op.Destination
	right := c.R[dest]
	var result uint16
	switch op.Instruction.Name {
	case MviName:
		c.R[dest] = value
		return

	case AddName:
		result = c.add(value, right)

	case SubName:
		result = c.subtract(value, right)

	case CmpName:
		c.subtract(value, right)
		return

	case AndName:
		result = value & right
		c.setSZ(result)

	case XorName:
		result = value ^ right
		c.setSZ(result)
	}
	c.R[dest] = result
}

func (c *CPU) readOperand(op Opcode, doubleByte bool) uint16 {
	if op.Addressing == DirectAddressing {
		address := c.memory.Read(c.R[ProgramCounter])
		c.R[ProgramCounter]++
		return c.memory.Read(address)
	}
	addressRegister := op.Source
	if addressRegister == StackPointer {
		c.R[StackPointer]--
	}
	address := c.R[addressRegister]
	value := c.memory.Read(address)
	if addressRegister >= 4 && addressRegister != StackPointer {
		c.R[addressRegister]++
	}
	if doubleByte {
		second := c.memory.Read(c.R[addressRegister])
		if addressRegister >= 4 && addressRegister != StackPointer {
			c.R[addressRegister]++
		}
		value = value&0xFF | second<<8
		c.cycles += 2
	}
	return value
}

func (c *CPU) storeOperand(op Opcode) {
	if op.Addressing == DirectAddressing {
		address := c.memory.Read(c.R[ProgramCounter])
		c.R[ProgramCounter]++
		c.memory.Write(address, c.R[op.Destination])
		return
	}
	addressRegister := op.Source
	address := c.R[addressRegister]
	if addressRegister == ProgramCounter {
		c.R[ProgramCounter]++
	}
	value := c.R[op.Destination]
	c.memory.Write(address, value)
	if addressRegister >= 4 && addressRegister != ProgramCounter {
		c.R[addressRegister]++
	}
}
