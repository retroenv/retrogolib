package m6502

import (
	"fmt"

	. "github.com/retroenv/retrogolib/addressing"
)

type paramReaderFunc func(c *CPU) ([]any, []byte, bool)

var paramReader = map[Mode]paramReaderFunc{
	ImpliedAddressing:     paramReaderImplied,
	ImmediateAddressing:   paramReaderImmediate,
	AccumulatorAddressing: paramReaderAccumulator,
	AbsoluteAddressing:    paramReaderAbsolute,
	AbsoluteXAddressing:   paramReaderAbsoluteX,
	AbsoluteYAddressing:   paramReaderAbsoluteY,
	ZeroPageAddressing:    paramReaderZeroPage,
	ZeroPageXAddressing:   paramReaderZeroPageX,
	ZeroPageYAddressing:   paramReaderZeroPageY,
	RelativeAddressing:    paramReaderRelative,
	IndirectAddressing:    paramReaderIndirect,
	IndirectXAddressing:   paramReaderIndirectX,
	IndirectYAddressing:   paramReaderIndirectY,
}

// readOpParams reads the opcode parameters after the first opcode byte
// and translates it into emulator specific types.
func readOpParams(c *CPU, addressing Mode) ([]any, []byte, bool, error) {
	fun, ok := paramReader[addressing]
	if !ok {
		return nil, nil, false, fmt.Errorf("unsupported addressing mode %00x", addressing)
	}

	params, opcodes, pageCrossed := fun(c)
	return params, opcodes, pageCrossed, nil
}

func paramReaderImplied(_ *CPU) ([]any, []byte, bool) {
	return nil, nil, false
}

func paramReaderImmediate(c *CPU) ([]any, []byte, bool) {
	b := c.memory.Read(c.PC + 1)
	params := []any{int(b)}
	opcodes := []byte{b}
	return params, opcodes, false
}

func paramReaderAccumulator(_ *CPU) ([]any, []byte, bool) {
	params := []any{Accumulator(0)}
	return params, nil, false
}

func paramReaderAbsolute(c *CPU) ([]any, []byte, bool) {
	b1 := uint16(c.memory.Read(c.PC + 1))
	b2 := uint16(c.memory.Read(c.PC + 2))

	params := []any{Absolute(b2<<8 | b1)}
	opcodes := []byte{byte(b1), byte(b2)}
	return params, opcodes, false
}

func paramReaderAbsoluteX(c *CPU) ([]any, []byte, bool) {
	b1 := uint16(c.memory.Read(c.PC + 1))
	b2 := uint16(c.memory.Read(c.PC + 2))
	w := b2<<8 | b1
	_, pageCrossed := offsetAddress(w, c.X)

	params := []any{Absolute(w), &c.X}
	opcodes := []byte{byte(b1), byte(b2)}
	return params, opcodes, pageCrossed
}

func paramReaderAbsoluteY(c *CPU) ([]any, []byte, bool) {
	b1 := uint16(c.memory.Read(c.PC + 1))
	b2 := uint16(c.memory.Read(c.PC + 2))
	w := b2<<8 | b1
	_, pageCrossed := offsetAddress(w, c.Y)

	params := []any{Absolute(w), &c.Y}
	opcodes := []byte{byte(b1), byte(b2)}
	return params, opcodes, pageCrossed
}

func paramReaderZeroPage(c *CPU) ([]any, []byte, bool) {
	b := c.memory.Read(c.PC + 1)

	params := []any{Absolute(b)}
	opcodes := []byte{b}
	return params, opcodes, false
}

func paramReaderZeroPageX(c *CPU) ([]any, []byte, bool) {
	b := c.memory.Read(c.PC + 1)

	params := []any{ZeroPage(b), &c.X}
	opcodes := []byte{b}
	return params, opcodes, false
}

func paramReaderZeroPageY(c *CPU) ([]any, []byte, bool) {
	b := c.memory.Read(c.PC + 1)

	params := []any{ZeroPage(b), &c.Y}
	opcodes := []byte{b}
	return params, opcodes, false
}

func paramReaderRelative(c *CPU) ([]any, []byte, bool) {
	offset := uint16(c.memory.Read(c.PC + 1))

	var address uint16
	if offset < 0x80 {
		address = c.PC + 2 + offset
	} else {
		address = c.PC + 2 + offset - 0x100
	}

	params := []any{Absolute(address)}
	opcodes := []byte{byte(offset)}
	return params, opcodes, false
}

func paramReaderIndirect(c *CPU) ([]any, []byte, bool) {
	address := c.memory.ReadWordBug(c.PC + 1)
	b1 := uint16(c.memory.Read(c.PC + 1))
	b2 := uint16(c.memory.Read(c.PC + 2))

	params := []any{Indirect(address)}
	opcodes := []byte{byte(b1), byte(b2)}
	return params, opcodes, false
}

func paramReaderIndirectX(c *CPU) ([]any, []byte, bool) {
	b := c.memory.Read(c.PC + 1)
	offset := uint16(b + c.X)

	address := c.memory.ReadWordBug(offset)
	params := []any{IndirectResolved(address), &c.X}

	opcodes := []byte{b}
	return params, opcodes, false
}

func paramReaderIndirectY(c *CPU) ([]any, []byte, bool) {
	b := c.memory.Read(c.PC + 1)

	var pageCrossed bool

	address := c.memory.ReadWordBug(uint16(b))
	address, pageCrossed = offsetAddress(address, c.Y)
	params := []any{IndirectResolved(address), &c.Y}

	opcodes := []byte{b}
	return params, opcodes, pageCrossed
}

// offsetAddress returns the offset address and whether it crosses a page boundary.
func offsetAddress(address uint16, offset byte) (uint16, bool) {
	newAddress := address + uint16(offset)
	pageCrossed := newAddress&0xff00 != address&0xff00
	return newAddress, pageCrossed
}

// hasAccumulatorParam returns whether the passed or missing parameter
// indicates usage of the accumulator register.
func hasAccumulatorParam(params ...any) bool {
	if params == nil {
		return true
	}
	param := params[0]
	_, ok := param.(Accumulator)
	return ok
}
