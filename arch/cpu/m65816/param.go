package m65816

import "fmt"

type paramReaderFunc func(c *CPU) ([]any, []byte, bool)

var paramReader = map[AddressingMode]paramReaderFunc{
	ImpliedAddressing:                        paramReaderImplied,
	AccumulatorAddressing:                    paramReaderAccumulator,
	ImmediateAddressing:                      paramReaderImmediate,
	DirectPageAddressing:                     paramReaderDP,
	DirectPageIndexedXAddressing:             paramReaderDPX,
	DirectPageIndexedYAddressing:             paramReaderDPY,
	DirectPageIndirectAddressing:             paramReaderDPIndirect,
	DirectPageIndexedXIndirectAddressing:     paramReaderDPXIndirect,
	DirectPageIndirectIndexedYAddressing:     paramReaderDPIndirectY,
	DirectPageIndirectLongAddressing:         paramReaderDPIndirectLong,
	DirectPageIndirectLongIndexedYAddressing: paramReaderDPIndirectLongY,
	AbsoluteAddressing:                       paramReaderAbsolute,
	AbsoluteIndexedXAddressing:               paramReaderAbsoluteX,
	AbsoluteIndexedYAddressing:               paramReaderAbsoluteY,
	AbsoluteIndirectAddressing:               paramReaderAbsoluteIndirect,
	AbsoluteIndexedXIndirectAddressing:       paramReaderAbsoluteXIndirect,
	AbsoluteLongAddressing:                   paramReaderAbsoluteLong,
	AbsoluteLongIndexedXAddressing:           paramReaderAbsoluteLongX,
	AbsoluteIndirectLongAddressing:           paramReaderAbsoluteIndirectLong,
	StackRelativeAddressing:                  paramReaderSR,
	StackRelativeIndirectIndexedYAddressing:  paramReaderSRIndirectY,
	RelativeAddressing:                       paramReaderRelative,
	RelativeLongAddressing:                   paramReaderRelativeLong,
	BlockMoveAddressing:                      paramReaderBlockMove,
}

// readOpParams reads the instruction operand bytes for the given addressing mode.
// Returns the decoded params, raw operand bytes, page-crossed flag, and error.
func readOpParams(c *CPU, mode AddressingMode, op Opcode) ([]any, []byte, bool, error) {
	fn, ok := paramReader[mode]
	if !ok {
		return nil, nil, false, fmt.Errorf("%w: mode 0x%x", ErrUnsupportedAddressingMode, mode)
	}
	// For immediate mode, pass the width flag via a size-aware reader
	if mode == ImmediateAddressing {
		return paramReaderImmediateWidth(c, op)
	}
	params, opcodes, pageCrossed := fn(c)
	return params, opcodes, pageCrossed, nil
}

func paramReaderImplied(_ *CPU) ([]any, []byte, bool) {
	return nil, nil, false
}

func paramReaderAccumulator(_ *CPU) ([]any, []byte, bool) {
	return []any{Accumulator{}}, nil, false
}

// paramReaderImmediateWidth reads an immediate value whose size depends on M or X flags.
func paramReaderImmediateWidth(c *CPU, op Opcode) ([]any, []byte, bool, error) {
	switch op.WidthFlag {
	case WidthM:
		if c.AccWidth() == 2 {
			w := c.fetchWord(1)
			return []any{Immediate16(w)}, []byte{uint8(w), uint8(w >> 8)}, false, nil
		}
	case WidthX:
		if c.IdxWidth() == 2 {
			w := c.fetchWord(1)
			return []any{Immediate16(w)}, []byte{uint8(w), uint8(w >> 8)}, false, nil
		}
	}
	b := c.fetchByte(1)
	return []any{Immediate8(b)}, []byte{b}, false, nil
}

func paramReaderImmediate(c *CPU) ([]any, []byte, bool) {
	b := c.fetchByte(1)
	return []any{Immediate8(b)}, []byte{b}, false
}

func paramReaderDP(c *CPU) ([]any, []byte, bool) {
	b := c.fetchByte(1)
	return []any{DirectPage(b)}, []byte{b}, false
}

func paramReaderDPX(c *CPU) ([]any, []byte, bool) {
	b := c.fetchByte(1)
	return []any{DirectPageX(b)}, []byte{b}, false
}

func paramReaderDPY(c *CPU) ([]any, []byte, bool) {
	b := c.fetchByte(1)
	return []any{DirectPageY(b)}, []byte{b}, false
}

func paramReaderDPIndirect(c *CPU) ([]any, []byte, bool) {
	dp := c.fetchByte(1)
	// Pointer read uses DP page wrap in emulation mode (DP_lo=0)
	addr := uint32(c.readDPWord(dp))
	eff := c.dataAddr(uint16(addr))
	return []any{DPIndirect(eff)}, []byte{dp}, false
}

func paramReaderDPXIndirect(c *CPU) ([]any, []byte, bool) {
	dp := c.fetchByte(1)
	idx := uint16(dp) + (c.X & 0xFF)
	if c.IdxWidth() == 2 {
		idx = uint16(dp) + c.X
	}
	var addr uint32
	if c.E && c.DP&0xFF == 0 {
		// Emulation mode, DP page-aligned: (dp+X) and pointer bytes wrap within DP 256-byte page
		dpOffset := uint8(idx)
		dpPage := uint32(c.DP)
		lo := uint32(c.memory.ReadByte(dpPage | uint32(dpOffset)))
		hi := uint32(c.memory.ReadByte(dpPage | uint32(dpOffset+1))) // +1 wraps at 8 bits
		addr = hi<<8 | lo
	} else {
		ptr := bank24(0, c.DP+idx)
		addr = uint32(c.readMem16(ptr))
	}
	eff := bank24(c.DB, uint16(addr))
	return []any{DPIndirectX(eff)}, []byte{dp}, false
}

func paramReaderDPIndirectY(c *CPU) ([]any, []byte, bool) {
	dp := c.fetchByte(1)
	// Pointer read uses DP page wrap in emulation mode (DP_lo=0)
	base := uint32(c.readDPWord(dp))
	var yVal uint32
	if c.IdxWidth() == 2 {
		yVal = uint32(c.Y)
	} else {
		yVal = uint32(c.Y & 0xFF)
	}
	eff := c.dataAddr(uint16(base)) + yVal
	basePage := c.dataAddr(uint16(base)) & 0xFF00
	pageCrossed := eff&0xFF00 != basePage
	return []any{DPIndirectY(eff)}, []byte{dp}, pageCrossed
}

func paramReaderDPIndirectLong(c *CPU) ([]any, []byte, bool) {
	dp := c.fetchByte(1)
	ptr := c.dpAddr(dp)
	addr := c.readMem24(ptr) & 0xFFFFFF
	return []any{DPIndirectLong(addr)}, []byte{dp}, false
}

func paramReaderDPIndirectLongY(c *CPU) ([]any, []byte, bool) {
	dp := c.fetchByte(1)
	ptr := c.dpAddr(dp)
	base := c.readMem24(ptr) & 0xFFFFFF
	var yVal uint32
	if c.IdxWidth() == 2 {
		yVal = uint32(c.Y)
	} else {
		yVal = uint32(c.Y & 0xFF)
	}
	eff := base + yVal
	return []any{DPIndLongY(eff)}, []byte{dp}, false
}

func paramReaderAbsolute(c *CPU) ([]any, []byte, bool) {
	b1 := c.fetchByte(1)
	b2 := c.fetchByte(2)
	addr := uint16(b2)<<8 | uint16(b1)
	return []any{Absolute16(addr)}, []byte{b1, b2}, false
}

func paramReaderAbsoluteX(c *CPU) ([]any, []byte, bool) {
	b1 := c.fetchByte(1)
	b2 := c.fetchByte(2)
	base := uint16(b2)<<8 | uint16(b1)
	var xVal uint16
	if c.IdxWidth() == 2 {
		xVal = c.X
	} else {
		xVal = c.X & 0xFF
	}
	baseAddr := c.dataAddr(base)
	eff := baseAddr + uint32(xVal)
	pageCrossed := eff&0xFF00 != baseAddr&0xFF00
	return []any{AbsoluteX16(eff)}, []byte{b1, b2}, pageCrossed
}

func paramReaderAbsoluteY(c *CPU) ([]any, []byte, bool) {
	b1 := c.fetchByte(1)
	b2 := c.fetchByte(2)
	base := uint16(b2)<<8 | uint16(b1)
	var yVal uint16
	if c.IdxWidth() == 2 {
		yVal = c.Y
	} else {
		yVal = c.Y & 0xFF
	}
	baseAddr := c.dataAddr(base)
	eff := baseAddr + uint32(yVal)
	pageCrossed := eff&0xFF00 != baseAddr&0xFF00
	return []any{AbsoluteY16(eff)}, []byte{b1, b2}, pageCrossed
}

func paramReaderAbsoluteIndirect(c *CPU) ([]any, []byte, bool) {
	b1 := c.fetchByte(1)
	b2 := c.fetchByte(2)
	ptr := uint32(uint16(b2)<<8 | uint16(b1))
	addr := uint32(c.readMem16(ptr))
	// JMP (abs) stays in current bank
	eff := bank24(c.PB, uint16(addr))
	return []any{DPIndirect(eff)}, []byte{b1, b2}, false
}

func paramReaderAbsoluteXIndirect(c *CPU) ([]any, []byte, bool) {
	b1 := c.fetchByte(1)
	b2 := c.fetchByte(2)
	base := uint16(b2)<<8 | uint16(b1)
	ptr := bank24(c.PB, base+c.X)
	addr := uint32(c.readMem16(ptr))
	eff := bank24(c.PB, uint16(addr))
	return []any{DPIndirectX(eff)}, []byte{b1, b2}, false
}

func paramReaderAbsoluteLong(c *CPU) ([]any, []byte, bool) {
	b1 := c.fetchByte(1)
	b2 := c.fetchByte(2)
	b3 := c.fetchByte(3)
	addr := uint32(b3)<<16 | uint32(b2)<<8 | uint32(b1)
	return []any{AbsLong(addr)}, []byte{b1, b2, b3}, false
}

func paramReaderAbsoluteLongX(c *CPU) ([]any, []byte, bool) {
	b1 := c.fetchByte(1)
	b2 := c.fetchByte(2)
	b3 := c.fetchByte(3)
	base := uint32(b3)<<16 | uint32(b2)<<8 | uint32(b1)
	var xVal uint32
	if c.IdxWidth() == 2 {
		xVal = uint32(c.X)
	} else {
		xVal = uint32(c.X & 0xFF)
	}
	eff := (base + xVal) & 0xFFFFFF
	return []any{AbsLongX(eff)}, []byte{b1, b2, b3}, false
}

func paramReaderAbsoluteIndirectLong(c *CPU) ([]any, []byte, bool) {
	b1 := c.fetchByte(1)
	b2 := c.fetchByte(2)
	ptr := uint32(uint16(b2)<<8 | uint16(b1))
	addr := c.readMem24(ptr) & 0xFFFFFF
	return []any{AbsLong(addr)}, []byte{b1, b2}, false
}

func paramReaderSR(c *CPU) ([]any, []byte, bool) {
	b := c.fetchByte(1)
	return []any{StackRel(b)}, []byte{b}, false
}

func paramReaderSRIndirectY(c *CPU) ([]any, []byte, bool) {
	sr := c.fetchByte(1)
	ptr := bank24(0, c.SP+uint16(sr))
	base := uint32(c.readMem16(ptr))
	var yVal uint32
	if c.IdxWidth() == 2 {
		yVal = uint32(c.Y)
	} else {
		yVal = uint32(c.Y & 0xFF)
	}
	eff := c.dataAddr(uint16(base)) + yVal
	return []any{SRIndY(eff)}, []byte{sr}, false
}

func paramReaderRelative(c *CPU) ([]any, []byte, bool) {
	offset := int8(c.fetchByte(1))
	// Branch target: PC+2 (after the 2-byte instruction) + signed offset
	nextPC := c.PC + 2
	target := uint16(int32(nextPC) + int32(offset))
	// Page crossing: target lands in a different 256-byte page than the next instruction.
	// On the 65816, this penalty applies only in emulation mode (handled in step.go).
	pageCrossed := (target & 0xFF00) != (nextPC & 0xFF00)
	return []any{target}, []byte{uint8(offset)}, pageCrossed
}

func paramReaderRelativeLong(c *CPU) ([]any, []byte, bool) {
	b1 := c.fetchByte(1)
	b2 := c.fetchByte(2)
	offset := int16(uint16(b2)<<8 | uint16(b1))
	target := uint16(int32(c.PC) + 3 + int32(offset))
	return []any{target}, []byte{b1, b2}, false
}

func paramReaderBlockMove(c *CPU) ([]any, []byte, bool) {
	// Encoding: dst bank, src bank (note: reversed in machine code)
	dst := c.fetchByte(1)
	src := c.fetchByte(2)
	return []any{BlockMove{Src: src, Dst: dst}}, []byte{dst, src}, false
}
