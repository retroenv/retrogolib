package m68000

// Parameter reading helpers for the instruction handlers.

// readExtWord reads a 16-bit extension word from the instruction stream.
func (c *CPU) readExtWord() uint16 {
	return c.readWord()
}

// readExtLong reads a 32-bit extension long from the instruction stream.
func (c *CPU) readExtLong() uint32 {
	return c.readLong()
}
