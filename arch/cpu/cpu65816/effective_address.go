package cpu65816

import "fmt"

func (addr DPIndirect) resolvedAddress() uint32     { return uint32(addr) }
func (addr DPIndirectX) resolvedAddress() uint32    { return uint32(addr) }
func (addr DPIndirectY) resolvedAddress() uint32    { return uint32(addr) }
func (addr DPIndirectLong) resolvedAddress() uint32 { return uint32(addr) }
func (addr DPIndLongY) resolvedAddress() uint32     { return uint32(addr) }
func (addr AbsoluteX16) resolvedAddress() uint32    { return uint32(addr) }
func (addr AbsoluteY16) resolvedAddress() uint32    { return uint32(addr) }
func (addr AbsLong) resolvedAddress() uint32        { return uint32(addr) }
func (addr AbsLongX) resolvedAddress() uint32       { return uint32(addr) }
func (addr SRIndY) resolvedAddress() uint32         { return uint32(addr) }

// resolveDP handles emulation-mode wrapping for a direct-page operand.
func (c *CPU) resolveDP(dp uint8) uint32 {
	if c.E && c.DP&0xFF == 0 {
		return uint32(c.DP) | uint32(dp)
	}
	return bank24(0, c.DP+uint16(dp))
}

// resolveDPIndexed masks the index at its active width and applies direct-page wrapping.
func (c *CPU) resolveDPIndexed(dp uint8, index uint16) uint32 {
	if c.E && c.DP&0xFF == 0 {
		return uint32(c.DP) | uint32(dp+uint8(index&0xFF))
	}
	if c.IdxWidth() == 1 {
		index &= 0xFF
	}
	return (uint32(c.DP) + uint32(dp) + uint32(index)) & 0xFFFF
}

// resolveEA resolves a decoded operand to its readable 24-bit address.
func (c *CPU) resolveEA(param any) (uint32, error) {
	if addr, ok := param.(resolvedAddress); ok {
		return addr.resolvedAddress(), nil
	}

	switch value := param.(type) {
	case Immediate8, Immediate16:
		return 0, fmt.Errorf("%w: immediate operand type %T", ErrUnsupportedAddressingMode, param)
	case DirectPage:
		return c.resolveDP(uint8(value)), nil
	case DirectPageX:
		return c.resolveDPIndexed(uint8(value), c.X), nil
	case DirectPageY:
		return c.resolveDPIndexed(uint8(value), c.Y), nil
	case Absolute16:
		return c.dataAddr(uint16(value)), nil
	case StackRel:
		return bank24(0, c.SP+uint16(value)), nil
	default:
		return 0, fmt.Errorf("%w: type %T", ErrUnsupportedAddressingMode, param)
	}
}

type resolvedAddress interface {
	resolvedAddress() uint32
}
