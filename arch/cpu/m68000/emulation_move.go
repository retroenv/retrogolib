package m68000

// Data movement instructions: MOVE, MOVEA, MOVEQ, MOVEM, MOVEP, EXG, LEA, PEA,
// LINK, UNLK, SWAP.

func (c *CPU) execMOVE(d DecodedOpcode) error {
	if d.Extra != 0 {
		return c.execMOVESpecial(d)
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

	dstEA, err := c.decodeEA(d.DstMode, d.DstReg, d.Size)
	if err != nil {
		return err
	}

	c.setLogicFlags(src, d.Size)
	return c.writeEA(dstEA, src)
}

// execMOVESpecial handles MOVE to/from SR/CCR/USP.
func (c *CPU) execMOVESpecial(d DecodedOpcode) error {
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
		return c.writeEA(dstEA, uint32(sr))

	case 4: // MOVE to CCR
		srcEA, err := c.decodeEA(d.SrcMode, d.SrcReg, SizeWord)
		if err != nil {
			return err
		}
		src, err := c.readEA(srcEA)
		if err != nil {
			return err
		}
		c.SetCCR(uint8(src))
		return nil

	case 5: // MOVE to SR
		if !c.IsSupervisor() {
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
		c.SetSR(uint16(src))
		return nil

	default:
		return nil
	}
}

func (c *CPU) execMOVEA(d DecodedOpcode) error {
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

func (c *CPU) execMOVEQ(d DecodedOpcode) error {
	// Sign-extend 8-bit immediate to 32-bit.
	value := uint32(int32(int8(d.Extra)))
	c.D[d.DstReg] = value
	c.setLogicFlags(value, SizeLong)
	return nil
}

func (c *CPU) execMOVEM(d DecodedOpcode) error {
	mask := c.readWord()

	if d.Extra == 0 {
		// Register to memory.
		return c.execMOVEMToMem(d, mask)
	}
	// Memory to register.
	return c.execMOVEMToReg(d, mask)
}

// execMOVEMToMem moves registers to memory.
func (c *CPU) execMOVEMToMem(d DecodedOpcode, mask uint16) error {
	if d.DstMode == 4 {
		// Predecrement mode: register order is reversed (A7 first, D0 last).
		addr := c.getRegA(d.DstReg)
		for i := 15; i >= 0; i-- {
			if mask&(1<<uint(i)) == 0 {
				continue
			}

			addr -= uint32(d.Size)
			val := c.getMovemReg(15 - uint8(i))

			if err := c.writeMemory(addr, val, d.Size); err != nil {
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
func (c *CPU) execMOVEMToReg(d DecodedOpcode, mask uint16) error {
	ea, err := c.decodeEA(d.SrcMode, d.SrcReg, d.Size)
	if err != nil {
		return err
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

// getMovemReg returns the value of register i (0-7=D0-D7, 8-15=A0-A7).
func (c *CPU) getMovemReg(i uint8) uint32 {
	if i < 8 {
		return c.D[i]
	}
	return c.getRegA(i - 8)
}

// setMovemReg sets the value of register i (0-7=D0-D7, 8-15=A0-A7).
func (c *CPU) setMovemReg(i uint8, value uint32) {
	if i < 8 {
		c.D[i] = value
	} else {
		c.setRegA(i-8, value)
	}
}

func (c *CPU) execMOVEP(d DecodedOpcode) error {
	if d.SrcMode == 0 {
		// MOVEP Dn,d16(An): register to memory.
		disp := int16(c.readWord())
		addr := uint32(int32(c.getRegA(d.DstReg)) + int32(disp))
		val := c.D[d.SrcReg]

		if d.Size == SizeLong {
			c.bus.Write(addr, uint8(val>>24))
			c.bus.Write(addr+2, uint8(val>>16))
			c.bus.Write(addr+4, uint8(val>>8))
			c.bus.Write(addr+6, uint8(val))
		} else {
			c.bus.Write(addr, uint8(val>>8))
			c.bus.Write(addr+2, uint8(val))
		}
		return nil
	}

	// MOVEP d16(An),Dn: memory to register.
	disp := int16(c.readWord())
	addr := uint32(int32(c.getRegA(d.SrcReg)) + int32(disp))

	if d.Size == SizeLong {
		b0 := uint32(c.bus.Read(addr))
		b1 := uint32(c.bus.Read(addr + 2))
		b2 := uint32(c.bus.Read(addr + 4))
		b3 := uint32(c.bus.Read(addr + 6))
		c.D[d.DstReg] = (b0 << 24) | (b1 << 16) | (b2 << 8) | b3
	} else {
		b0 := uint32(c.bus.Read(addr))
		b1 := uint32(c.bus.Read(addr + 2))
		c.D[d.DstReg] = (c.D[d.DstReg] & 0xFFFF0000) | (b0 << 8) | b1
	}
	return nil
}

func (c *CPU) execEXG(d DecodedOpcode) error {
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

func (c *CPU) execLEA(d DecodedOpcode) error {
	ea, err := c.decodeEA(d.SrcMode, d.SrcReg, SizeLong)
	if err != nil {
		return err
	}
	c.setRegA(d.DstReg, ea.Address)
	return nil
}

func (c *CPU) execPEA(d DecodedOpcode) error {
	ea, err := c.decodeEA(d.DstMode, d.DstReg, SizeLong)
	if err != nil {
		return err
	}
	c.push32(ea.Address)
	return nil
}

func (c *CPU) execLINK(d DecodedOpcode) error {
	// Push current An.
	c.push32(c.getRegA(d.DstReg))
	// An = SP.
	c.setRegA(d.DstReg, c.sp)
	// SP += displacement.
	disp := int16(c.readWord())
	c.sp = uint32(int32(c.sp) + int32(disp))
	return nil
}

func (c *CPU) execUNLK(d DecodedOpcode) error {
	c.sp = c.getRegA(d.DstReg)
	c.setRegA(d.DstReg, c.pop32())
	return nil
}

func (c *CPU) execSWAP(d DecodedOpcode) error {
	val := c.D[d.DstReg]
	result := (val>>16)&0xFFFF | (val&0xFFFF)<<16
	c.D[d.DstReg] = result

	c.setFlagN(result, SizeLong)
	c.setFlagZ(result, SizeLong)
	c.Flags.V = 0
	c.Flags.C = 0
	return nil
}
