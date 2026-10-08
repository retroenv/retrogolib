package cpu65816

func (c *CPU) dataAddr(offset uint16) uint32 {
	return bank24(c.DB, offset)
}

func (c *CPU) dpAddr(offset uint8) uint32 {
	return bank24(0, c.DP+uint16(offset))
}

func (c *CPU) readMem8(addr uint32) uint8 {
	return c.memory.Read(addr & 0xFFFFFF)
}

func (c *CPU) writeMem8(addr uint32, value uint8) {
	c.memory.Write(addr&0xFFFFFF, value)
}

// readMem16 reads a word and keeps its high byte in the bank of the address.
func (c *CPU) readMem16(addr uint32) uint16 {
	addr &= 0xFFFFFF
	lo := uint16(c.memory.Read(addr))
	bank := addr & 0xFF0000
	hi := uint16(c.memory.Read(bank | uint32(uint16(addr)+1)))
	return hi<<8 | lo
}

// readData16 reads data across a bank boundary in the full 24-bit address space.
func (c *CPU) readData16(addr uint32) uint16 {
	addr &= 0xFFFFFF
	lo := uint16(c.memory.Read(addr))
	hi := uint16(c.memory.Read((addr + 1) & 0xFFFFFF))
	return hi<<8 | lo
}

// readDPWord keeps both pointer bytes in the direct-page block when emulation
// mode uses a page-aligned direct page.
func (c *CPU) readDPWord(dpOffset uint8) uint16 {
	if c.E && c.DP&0xFF == 0 {
		dpPage := uint32(c.DP)
		lo := uint16(c.memory.Read(dpPage | uint32(dpOffset)))
		hi := uint16(c.memory.Read(dpPage | uint32(dpOffset+1)))
		return hi<<8 | lo
	}
	return c.readMem16(bank24(0, c.DP+uint16(dpOffset)))
}

// writeMem16 writes a word and lets its high byte cross into the next bank.
func (c *CPU) writeMem16(addr uint32, value uint16) {
	addr &= 0xFFFFFF
	c.memory.Write(addr, uint8(value))
	c.memory.Write((addr+1)&0xFFFFFF, uint8(value>>8))
}

// writeMem16Bank writes a word and keeps its high byte in the bank of the address.
func (c *CPU) writeMem16Bank(addr uint32, value uint16) {
	addr &= 0xFFFFFF
	c.memory.Write(addr, uint8(value))
	bank := addr & 0xFF0000
	c.memory.Write(bank|uint32(uint16(addr)+1), uint8(value>>8))
}

// readOperandWord reads the 16-bit operand at addr. Direct page and stack
// relative operands stay in bank 0; all other operands can cross a bank.
func (c *CPU) readOperandWord(param any, addr uint32) uint16 {
	if wrapsBank(param) {
		return c.readMem16(addr)
	}
	return c.readData16(addr)
}

// writeOperandWord writes the 16-bit operand at addr. Direct page and stack
// relative operands stay in bank 0; all other operands can cross a bank.
func (c *CPU) writeOperandWord(param any, addr uint32, value uint16) {
	if wrapsBank(param) {
		c.writeMem16Bank(addr, value)
		return
	}
	c.writeMem16(addr, value)
}

// readMem24 keeps all pointer bytes in the bank containing the pointer.
func (c *CPU) readMem24(addr uint32) uint32 {
	addr &= 0xFFFFFF
	bank := addr & 0xFF0000
	lo := uint32(c.memory.Read(addr))
	mid := uint32(c.memory.Read(bank | uint32(uint16(addr)+1)))
	hi := uint32(c.memory.Read(bank | uint32(uint16(addr)+2)))
	return hi<<16 | mid<<8 | lo
}

// wrapsBank reports if 16-bit accesses through the operand must stay in bank 0.
func wrapsBank(param any) bool {
	switch param.(type) {
	case DirectPage, DirectPageX, DirectPageY, StackRel:
		return true
	default:
		return false
	}
}
