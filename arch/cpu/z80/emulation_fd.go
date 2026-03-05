package z80

// FD prefix instruction implementations - IY register operations.
// These are thin wrappers around shared indexed register operations.

func fdLdIYnn(c *CPU, p ...any) error  { return indexedLdRegNn(c, &c.IY, p...) }
func fdIncIY(c *CPU) error             { c.IY++; return nil }
func fdDecIY(c *CPU) error             { c.IY--; return nil }
func fdAddIYBc(c *CPU, p ...any) error { return indexedAddRegPair(c, &c.IY, c.bc(), p...) }
func fdAddIYDe(c *CPU, p ...any) error { return indexedAddRegPair(c, &c.IY, c.de(), p...) }
func fdAddIYIY(c *CPU, p ...any) error { return indexedAddRegPair(c, &c.IY, c.IY, p...) }
func fdAddIYSp(c *CPU, p ...any) error { return indexedAddRegPair(c, &c.IY, c.SP, p...) }
func fdLdNnIY(c *CPU, p ...any) error  { return indexedLdNnReg(c, c.IY, p...) }
func fdLdIYNn(c *CPU, p ...any) error  { return indexedLdRegFromNn(c, &c.IY, p...) }

// IY indexed load operations - Load register from (IY+d)
func fdLdBIYd(c *CPU, p ...any) error { return indexedLdRegFromMem(c, &c.B, c.IY, p...) }
func fdLdCIYd(c *CPU, p ...any) error { return indexedLdRegFromMem(c, &c.C, c.IY, p...) }
func fdLdDIYd(c *CPU, p ...any) error { return indexedLdRegFromMem(c, &c.D, c.IY, p...) }
func fdLdEIYd(c *CPU, p ...any) error { return indexedLdRegFromMem(c, &c.E, c.IY, p...) }
func fdLdHIYd(c *CPU, p ...any) error { return indexedLdRegFromMem(c, &c.H, c.IY, p...) }
func fdLdLIYd(c *CPU, p ...any) error { return indexedLdRegFromMem(c, &c.L, c.IY, p...) }
func fdLdAIYd(c *CPU, p ...any) error { return indexedLdRegFromMem(c, &c.A, c.IY, p...) }

// IY indexed store operations - Store register to (IY+d)
func fdLdIYdB(c *CPU, p ...any) error { return indexedLdMemFromReg(c, c.B, c.IY, p...) }
func fdLdIYdC(c *CPU, p ...any) error { return indexedLdMemFromReg(c, c.C, c.IY, p...) }
func fdLdIYdD(c *CPU, p ...any) error { return indexedLdMemFromReg(c, c.D, c.IY, p...) }
func fdLdIYdE(c *CPU, p ...any) error { return indexedLdMemFromReg(c, c.E, c.IY, p...) }
func fdLdIYdH(c *CPU, p ...any) error { return indexedLdMemFromReg(c, c.H, c.IY, p...) }
func fdLdIYdL(c *CPU, p ...any) error { return indexedLdMemFromReg(c, c.L, c.IY, p...) }
func fdLdIYdA(c *CPU, p ...any) error { return indexedLdMemFromReg(c, c.A, c.IY, p...) }
func fdLdIYdN(c *CPU, p ...any) error { return indexedLdMemN(c, c.IY, p...) }

// IY indexed INC/DEC
func fdIncIYd(c *CPU, p ...any) error { return indexedIncMem(c, c.IY, p...) }
func fdDecIYd(c *CPU, p ...any) error { return indexedDecMem(c, c.IY, p...) }

// IY arithmetic operations with accumulator
func fdAddAIYd(c *CPU, p ...any) error { return indexedAddA(c, c.IY, p...) }
func fdAdcAIYd(c *CPU, p ...any) error { return indexedAdcA(c, c.IY, p...) }
func fdSubAIYd(c *CPU, p ...any) error { return indexedSubA(c, c.IY, p...) }
func fdSbcAIYd(c *CPU, p ...any) error { return indexedSbcA(c, c.IY, p...) }
func fdAndAIYd(c *CPU, p ...any) error { return indexedAndA(c, c.IY, p...) }
func fdXorAIYd(c *CPU, p ...any) error { return indexedXorA(c, c.IY, p...) }
func fdOrAIYd(c *CPU, p ...any) error  { return indexedOrA(c, c.IY, p...) }
func fdCpAIYd(c *CPU, p ...any) error  { return indexedCpA(c, c.IY, p...) }

// IY stack and jump operations
func fdJpIY(c *CPU) error   { c.PC = c.IY; return nil }
func fdExSpIY(c *CPU) error { return indexedExSp(c, &c.IY) }
func fdPushIY(c *CPU) error { c.push16(c.IY); return nil }
func fdPopIY(c *CPU) error  { c.IY = c.pop16(); return nil }
func fdLdSpIY(c *CPU) error { c.SP = c.IY; return nil }

// FDCB operations - bit operations on (IY+d)
func fdcbShift(c *CPU, p ...any) error { return indexedCBShift(c, c.IY, p...) }
func fdcbBit(c *CPU, p ...any) error   { return indexedCBBit(c, c.IY, p...) }
func fdcbRes(c *CPU, p ...any) error   { return indexedCBRes(c, c.IY, p...) }
func fdcbSet(c *CPU, p ...any) error   { return indexedCBSet(c, c.IY, p...) }
