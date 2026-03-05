package z80

// DD prefix instruction implementations - IX register operations.
// These are thin wrappers around shared indexed register operations.

func ddLdIXnn(c *CPU, p ...any) error  { return indexedLdRegNn(c, &c.IX, p...) }
func ddIncIX(c *CPU) error             { c.IX++; return nil }
func ddDecIX(c *CPU) error             { c.IX--; return nil }
func ddAddIXBc(c *CPU, p ...any) error { return indexedAddRegPair(c, &c.IX, c.bc(), p...) }
func ddAddIXDe(c *CPU, p ...any) error { return indexedAddRegPair(c, &c.IX, c.de(), p...) }
func ddAddIXIX(c *CPU, p ...any) error { return indexedAddRegPair(c, &c.IX, c.IX, p...) }
func ddAddIXSp(c *CPU, p ...any) error { return indexedAddRegPair(c, &c.IX, c.SP, p...) }
func ddLdNnIX(c *CPU, p ...any) error  { return indexedLdNnReg(c, c.IX, p...) }
func ddLdIXNn(c *CPU, p ...any) error  { return indexedLdRegFromNn(c, &c.IX, p...) }

// IX indexed load operations - Load register from (IX+d)
func ddLdBIXd(c *CPU, p ...any) error { return indexedLdRegFromMem(c, &c.B, c.IX, p...) }
func ddLdCIXd(c *CPU, p ...any) error { return indexedLdRegFromMem(c, &c.C, c.IX, p...) }
func ddLdDIXd(c *CPU, p ...any) error { return indexedLdRegFromMem(c, &c.D, c.IX, p...) }
func ddLdEIXd(c *CPU, p ...any) error { return indexedLdRegFromMem(c, &c.E, c.IX, p...) }
func ddLdHIXd(c *CPU, p ...any) error { return indexedLdRegFromMem(c, &c.H, c.IX, p...) }
func ddLdLIXd(c *CPU, p ...any) error { return indexedLdRegFromMem(c, &c.L, c.IX, p...) }
func ddLdAIXd(c *CPU, p ...any) error { return indexedLdRegFromMem(c, &c.A, c.IX, p...) }

// IX indexed store operations - Store register to (IX+d)
func ddLdIXdB(c *CPU, p ...any) error { return indexedLdMemFromReg(c, c.B, c.IX, p...) }
func ddLdIXdC(c *CPU, p ...any) error { return indexedLdMemFromReg(c, c.C, c.IX, p...) }
func ddLdIXdD(c *CPU, p ...any) error { return indexedLdMemFromReg(c, c.D, c.IX, p...) }
func ddLdIXdE(c *CPU, p ...any) error { return indexedLdMemFromReg(c, c.E, c.IX, p...) }
func ddLdIXdH(c *CPU, p ...any) error { return indexedLdMemFromReg(c, c.H, c.IX, p...) }
func ddLdIXdL(c *CPU, p ...any) error { return indexedLdMemFromReg(c, c.L, c.IX, p...) }
func ddLdIXdA(c *CPU, p ...any) error { return indexedLdMemFromReg(c, c.A, c.IX, p...) }
func ddLdIXdN(c *CPU, p ...any) error { return indexedLdMemN(c, c.IX, p...) }

// IX indexed INC/DEC
func ddIncIXd(c *CPU, p ...any) error { return indexedIncMem(c, c.IX, p...) }
func ddDecIXd(c *CPU, p ...any) error { return indexedDecMem(c, c.IX, p...) }

// IX arithmetic operations with accumulator
func ddAddAIXd(c *CPU, p ...any) error { return indexedAddA(c, c.IX, p...) }
func ddAdcAIXd(c *CPU, p ...any) error { return indexedAdcA(c, c.IX, p...) }
func ddSubAIXd(c *CPU, p ...any) error { return indexedSubA(c, c.IX, p...) }
func ddSbcAIXd(c *CPU, p ...any) error { return indexedSbcA(c, c.IX, p...) }
func ddAndAIXd(c *CPU, p ...any) error { return indexedAndA(c, c.IX, p...) }
func ddXorAIXd(c *CPU, p ...any) error { return indexedXorA(c, c.IX, p...) }
func ddOrAIXd(c *CPU, p ...any) error  { return indexedOrA(c, c.IX, p...) }
func ddCpAIXd(c *CPU, p ...any) error  { return indexedCpA(c, c.IX, p...) }

// IX stack and jump operations
func ddJpIX(c *CPU) error   { c.PC = c.IX; return nil }
func ddExSpIX(c *CPU) error { return indexedExSp(c, &c.IX) }
func ddPushIX(c *CPU) error { c.push16(c.IX); return nil }
func ddPopIX(c *CPU) error  { c.IX = c.pop16(); return nil }
func ddLdSpIX(c *CPU) error { c.SP = c.IX; return nil }

// DDCB operations - bit operations on (IX+d)
func ddcbShift(c *CPU, p ...any) error { return indexedCBShift(c, c.IX, p...) }
func ddcbBit(c *CPU, p ...any) error   { return indexedCBBit(c, c.IX, p...) }
func ddcbRes(c *CPU, p ...any) error   { return indexedCBRes(c, c.IX, p...) }
func ddcbSet(c *CPU, p ...any) error   { return indexedCBSet(c, c.IX, p...) }
