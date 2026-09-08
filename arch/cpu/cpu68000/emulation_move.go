package cpu68000

import "math/bits"

// Data movement instructions: MOVE, MOVEA, MOVEQ, MOVEM, MOVEP, EXG, LEA, PEA,
// LINK, UNLK, SWAP.

// getMovemReg returns the value of register i (0-7=D0-D7, 8-15=A0-A7).
func (cpu *CPU) getMovemReg(i uint8) uint32 {
	if i < 8 {
		return cpu.D[i]
	}
	return cpu.getRegA(i - 8)
}

// setMovemReg sets the value of register i (0-7=D0-D7, 8-15=A0-A7).
func (cpu *CPU) setMovemReg(i uint8, value uint32) {
	if i < 8 {
		cpu.D[i] = value
	} else {
		cpu.setRegA(i-8, value)
	}
}

func (cpu *CPU) writeMovePredecrement(reg uint8, value uint32, size OperandSize) error {
	cpu.setLogicFlags(value, size)
	cpu.operandPCOffset = 0
	// MOVE's destination decrement overlaps the next instruction prefetch.
	cpu.accessCycles += 4
	if size != SizeLong {
		cpu.setRegA(reg, cpu.getRegA(reg)-incrementSize(reg, size))
		return cpu.writeMemory(cpu.getRegA(reg), value, size)
	}
	cpu.setRegA(reg, cpu.getRegA(reg)-2)
	cpu.writeBusWord(cpu.getRegA(reg), uint16(value))
	cpu.setRegA(reg, cpu.getRegA(reg)-2)
	cpu.writeBusWord(cpu.getRegA(reg), uint16(value>>16))
	return nil
}

func execMOVE(c *CPU, d DecodedOpcode) error {
	if d.Extra != 0 {
		return execMOVESpecial(c, d)
	}

	// Regular MOVE.
	srcEA, err := c.decodeEA(d.SrcMode, d.SrcReg, d.Size)
	if err != nil {
		return err
	}
	src, err := c.readEA(srcEA)
	if err != nil {
		return err
	}

	if d.DstMode == 4 {
		return c.writeMovePredecrement(d.DstReg, src, d.Size)
	}
	dstMode := d.DstMode
	if dstMode == 3 {
		dstMode = 2
	}
	dstEA, err := c.decodeEA(dstMode, d.DstReg, d.Size)
	if err != nil {
		return err
	}

	c.setLogicFlags(src, d.Size)
	if d.DstMode == 7 && d.DstReg == 1 && d.SrcMode >= 2 && (d.SrcMode != 7 || d.SrcReg != 4) {
		c.operandPCOffset = -4
	}
	if err := c.writeEA(dstEA, src); err != nil {
		return err
	}
	if d.DstMode == 3 {
		c.setRegA(d.DstReg, dstEA.Address+incrementSize(d.DstReg, d.Size))
	}
	return nil
}

// execMOVESpecial handles MOVE to/from SR/CCR/USP.
func execMOVESpecial(c *CPU, d DecodedOpcode) error {
	switch d.Extra {
	case 1: // MOVE An,USP
		if !c.IsSupervisor() {
			return c.processException(VectorPrivilege)
		}
		c.USP = c.getRegA(d.SrcReg)
		return nil

	case 2: // MOVE USP,An
		if !c.IsSupervisor() {
			return c.processException(VectorPrivilege)
		}
		c.setRegA(d.DstReg, c.USP)
		return nil

	case 3: // MOVE from SR
		sr := c.GetSR()
		dstEA, err := c.decodeEA(d.DstMode, d.DstReg, SizeWord)
		if err != nil {
			return err
		}
		if _, err := c.readEA(dstEA); err != nil {
			return err
		}
		return c.writeEA(dstEA, uint32(sr))

	case 4, 5: // MOVE to CCR/SR.
		if d.Extra == 5 && !c.IsSupervisor() {
			return c.processException(VectorPrivilege)
		}
		srcEA, err := c.decodeEA(d.SrcMode, d.SrcReg, SizeWord)
		if err != nil {
			return err
		}
		src, err := c.readEA(srcEA)
		if err != nil {
			return err
		}
		if d.Extra == 4 {
			c.SetCCR(uint8(src))
		} else {
			c.SetSR(uint16(src))
		}
		return nil

	default:
		return nil
	}
}

func execMOVEA(c *CPU, d DecodedOpcode) error {
	srcEA, err := c.decodeEA(d.SrcMode, d.SrcReg, d.Size)
	if err != nil {
		return err
	}
	src, err := c.readEA(srcEA)
	if err != nil {
		return err
	}

	// MOVEA sign-extends word to long. No flags affected.
	src = signExtend(src, d.Size)
	c.setRegA(d.DstReg, src)
	return nil
}

func execMOVEQ(c *CPU, d DecodedOpcode) error {
	// Sign-extend 8-bit immediate to 32-bit.
	value := uint32(int32(int8(d.Extra)))
	c.D[d.DstReg] = value
	c.setLogicFlags(value, SizeLong)
	return nil
}

func execMOVEM(c *CPU, d DecodedOpcode) error {
	mask := c.readWord()
	c.cycles += uint64(bits.OnesCount16(mask)) * sizeCycles(d.Size, 4, 8)

	if d.Extra == 0 {
		// Register to memory.
		return execMOVEMToMem(c, d, mask)
	}
	// Memory to register.
	return execMOVEMToReg(c, d, mask)
}

// execMOVEMToMem moves registers to memory.
func execMOVEMToMem(c *CPU, d DecodedOpcode, mask uint16) error {
	if d.DstMode == 4 {
		// Predecrement mode: register order is reversed (A7 first, D0 last).
		addr := c.getRegA(d.DstReg)
		for i := range 16 {
			if mask&(1<<uint(i)) == 0 {
				continue
			}

			addr -= uint32(d.Size)
			val := c.getMovemReg(15 - uint8(i))

			if d.Size == SizeLong {
				c.writeBusWord(addr+2, uint16(val))
				c.writeBusWord(addr, uint16(val>>16))
			} else if err := c.writeMemory(addr, val, d.Size); err != nil {
				return err
			}
		}
		c.setRegA(d.DstReg, addr)
		return nil
	}

	ea, err := c.decodeEA(d.DstMode, d.DstReg, d.Size)
	if err != nil {
		return err
	}

	addr := ea.Address
	for i := range 16 {
		if mask&(1<<uint(i)) == 0 {
			continue
		}

		val := c.getMovemReg(uint8(i))

		if err := c.writeMemory(addr, val, d.Size); err != nil {
			return err
		}
		addr += uint32(d.Size)
	}

	return nil
}

// execMOVEMToReg moves memory to registers.
func execMOVEMToReg(c *CPU, d DecodedOpcode, mask uint16) error {
	ea, err := c.decodeEA(d.SrcMode, d.SrcReg, d.Size)
	if err != nil {
		return err
	}
	if d.SrcMode == 3 && d.Size == SizeLong {
		c.setRegA(d.SrcReg, ea.Address+2)
	}

	addr := ea.Address
	for i := range 16 {
		if mask&(1<<uint(i)) == 0 {
			continue
		}

		val, err := c.readMemory(addr, d.Size)
		if err != nil {
			return err
		}

		// Sign-extend word to long.
		if d.Size == SizeWord {
			val = signExtend(val, SizeWord)
		}

		c.setMovemReg(uint8(i), val)
		addr += uint32(d.Size)
	}

	// Update address register for postincrement mode.
	if d.SrcMode == 3 {
		c.setRegA(d.SrcReg, addr)
	}

	return nil
}

func execMOVEP(c *CPU, d DecodedOpcode) error {
	if d.SrcMode == 0 {
		// MOVEP Dn,d16(An): register to memory.
		disp := int16(c.readWord())
		addr := uint32(int32(c.getRegA(d.DstReg)) + int32(disp))
		val := c.D[d.SrcReg]

		if d.Size == SizeLong {
			c.writeByte(addr, uint8(val>>24))
			c.writeByte(addr+2, uint8(val>>16))
			c.writeByte(addr+4, uint8(val>>8))
			c.writeByte(addr+6, uint8(val))
		} else {
			c.writeByte(addr, uint8(val>>8))
			c.writeByte(addr+2, uint8(val))
		}
		return nil
	}

	// MOVEP d16(An),Dn: memory to register.
	disp := int16(c.readWord())
	addr := uint32(int32(c.getRegA(d.SrcReg)) + int32(disp))

	if d.Size == SizeLong {
		b0 := uint32(c.readByte(addr))
		b1 := uint32(c.readByte(addr + 2))
		b2 := uint32(c.readByte(addr + 4))
		b3 := uint32(c.readByte(addr + 6))
		c.D[d.DstReg] = (b0 << 24) | (b1 << 16) | (b2 << 8) | b3
	} else {
		b0 := uint32(c.readByte(addr))
		b1 := uint32(c.readByte(addr + 2))
		c.D[d.DstReg] = (c.D[d.DstReg] & 0xFFFF0000) | (b0 << 8) | b1
	}
	return nil
}

func execEXG(c *CPU, d DecodedOpcode) error {
	switch d.Extra {
	case 0: // EXG Dn,Dn
		c.D[d.SrcReg], c.D[d.DstReg] = c.D[d.DstReg], c.D[d.SrcReg]
	case 1: // EXG An,An
		srcA := c.getRegA(d.SrcReg)
		dstA := c.getRegA(d.DstReg)
		c.setRegA(d.SrcReg, dstA)
		c.setRegA(d.DstReg, srcA)
	case 2: // EXG Dn,An
		dn := c.D[d.SrcReg]
		an := c.getRegA(d.DstReg)
		c.D[d.SrcReg] = an
		c.setRegA(d.DstReg, dn)
	}
	return nil
}

func execLEA(c *CPU, d DecodedOpcode) error {
	ea, err := c.decodeEA(d.SrcMode, d.SrcReg, SizeLong)
	if err != nil {
		return err
	}
	c.setRegA(d.DstReg, ea.Address)
	return nil
}

func execPEA(c *CPU, d DecodedOpcode) error {
	ea, err := c.decodeEA(d.DstMode, d.DstReg, SizeLong)
	if err != nil {
		return err
	}
	c.push32(ea.Address)
	return nil
}

func execLINK(c *CPU, d DecodedOpcode) error {
	// Push current An.
	value := c.getRegA(d.DstReg)
	if d.DstReg == 7 {
		value -= 4
	}
	c.push32(value)
	// An = SP.
	c.setRegA(d.DstReg, c.sp)
	// SP += displacement.
	disp := int16(c.readWord())
	c.sp = uint32(int32(c.sp) + int32(disp))
	return nil
}

func execUNLK(c *CPU, d DecodedOpcode) error {
	c.sp = c.getRegA(d.DstReg)
	c.setRegA(d.DstReg, c.pop32())
	return nil
}

func execSWAP(c *CPU, d DecodedOpcode) error {
	val := c.D[d.DstReg]
	result := (val>>16)&0xFFFF | (val&0xFFFF)<<16
	c.D[d.DstReg] = result

	c.setFlagN(result, SizeLong)
	c.setFlagZ(result, SizeLong)
	c.Flags.V = 0
	c.Flags.C = 0
	return nil
}
