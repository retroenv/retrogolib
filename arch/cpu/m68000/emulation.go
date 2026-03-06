package m68000

// Core ALU operations: ADD, ADDA, ADDI, ADDQ, ADDX, SUB, SUBA, SUBI, SUBQ, SUBX,
// NEG, NEGX, CLR, CMP, CMPA, CMPI, CMPM, AND, ANDI, OR, ORI, EOR, EORI, NOT, EXT,
// TST, MULU, MULS, DIVU, DIVS, ABCD, SBCD, NBCD.

func execADD(c *CPU, d DecodedOpcode) error {
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
	dst, err := c.readEA(dstEA)
	if err != nil {
		return err
	}

	result := src + dst
	c.setAddFlags(src, dst, result, d.Size)
	return c.writeEA(dstEA, result)
}

func execADDA(c *CPU, d DecodedOpcode) error {
	srcEA, err := c.decodeEA(d.SrcMode, d.SrcReg, d.Size)
	if err != nil {
		return err
	}
	src, err := c.readEA(srcEA)
	if err != nil {
		return err
	}
	src = signExtend(src, d.Size)
	c.setRegA(d.DstReg, c.getRegA(d.DstReg)+src)
	return nil
}

func execADDI(c *CPU, d DecodedOpcode) error {
	imm := c.readImmediate(d.Size)

	dstEA, err := c.decodeEA(d.DstMode, d.DstReg, d.Size)
	if err != nil {
		return err
	}
	dst, err := c.readEA(dstEA)
	if err != nil {
		return err
	}

	result := imm + dst
	c.setAddFlags(imm, dst, result, d.Size)
	return c.writeEA(dstEA, result)
}

func execADDQ(c *CPU, d DecodedOpcode) error {
	imm := uint32(d.Extra)
	if d.DstMode == 1 {
		// ADDQ to address register: no flags affected, full 32-bit.
		c.setRegA(d.DstReg, c.getRegA(d.DstReg)+imm)
		return nil
	}

	dstEA, err := c.decodeEA(d.DstMode, d.DstReg, d.Size)
	if err != nil {
		return err
	}
	dst, err := c.readEA(dstEA)
	if err != nil {
		return err
	}

	result := dst + imm
	c.setAddFlags(imm, dst, result, d.Size)
	return c.writeEA(dstEA, result)
}

func execADDX(c *CPU, d DecodedOpcode) error {
	x := uint32(c.Flags.X)

	var src, dst uint32
	if d.Extra == 0 {
		// Register to register.
		src = c.getRegD(d.SrcReg, d.Size)
		dst = c.getRegD(d.DstReg, d.Size)
	} else {
		// Memory to memory (predecrement).
		c.setRegA(d.SrcReg, c.getRegA(d.SrcReg)-uint32(d.Size))
		sv, err := c.readMemory(c.getRegA(d.SrcReg), d.Size)
		if err != nil {
			return err
		}
		src = sv

		c.setRegA(d.DstReg, c.getRegA(d.DstReg)-uint32(d.Size))
		dv, err := c.readMemory(c.getRegA(d.DstReg), d.Size)
		if err != nil {
			return err
		}
		dst = dv
	}

	result := src + dst + x
	c.setAddXFlags(src, dst, result, d.Size)

	if d.Extra == 0 {
		c.setRegD(d.DstReg, result, d.Size)
	} else {
		return c.writeMemory(c.getRegA(d.DstReg), result, d.Size)
	}

	return nil
}

func execSUB(c *CPU, d DecodedOpcode) error {
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
	dst, err := c.readEA(dstEA)
	if err != nil {
		return err
	}

	result := dst - src
	c.setSubFlags(src, dst, result, d.Size)
	return c.writeEA(dstEA, result)
}

func execSUBA(c *CPU, d DecodedOpcode) error {
	srcEA, err := c.decodeEA(d.SrcMode, d.SrcReg, d.Size)
	if err != nil {
		return err
	}
	src, err := c.readEA(srcEA)
	if err != nil {
		return err
	}
	src = signExtend(src, d.Size)
	c.setRegA(d.DstReg, c.getRegA(d.DstReg)-src)
	return nil
}

func execSUBI(c *CPU, d DecodedOpcode) error {
	imm := c.readImmediate(d.Size)

	dstEA, err := c.decodeEA(d.DstMode, d.DstReg, d.Size)
	if err != nil {
		return err
	}
	dst, err := c.readEA(dstEA)
	if err != nil {
		return err
	}

	result := dst - imm
	c.setSubFlags(imm, dst, result, d.Size)
	return c.writeEA(dstEA, result)
}

func execSUBQ(c *CPU, d DecodedOpcode) error {
	imm := uint32(d.Extra)
	if d.DstMode == 1 {
		c.setRegA(d.DstReg, c.getRegA(d.DstReg)-imm)
		return nil
	}

	dstEA, err := c.decodeEA(d.DstMode, d.DstReg, d.Size)
	if err != nil {
		return err
	}
	dst, err := c.readEA(dstEA)
	if err != nil {
		return err
	}

	result := dst - imm
	c.setSubFlags(imm, dst, result, d.Size)
	return c.writeEA(dstEA, result)
}

func execSUBX(c *CPU, d DecodedOpcode) error {
	x := uint32(c.Flags.X)

	var src, dst uint32
	if d.Extra == 0 {
		src = c.getRegD(d.SrcReg, d.Size)
		dst = c.getRegD(d.DstReg, d.Size)
	} else {
		c.setRegA(d.SrcReg, c.getRegA(d.SrcReg)-uint32(d.Size))
		sv, err := c.readMemory(c.getRegA(d.SrcReg), d.Size)
		if err != nil {
			return err
		}
		src = sv

		c.setRegA(d.DstReg, c.getRegA(d.DstReg)-uint32(d.Size))
		dv, err := c.readMemory(c.getRegA(d.DstReg), d.Size)
		if err != nil {
			return err
		}
		dst = dv
	}

	result := dst - src - x
	c.setSubXFlags(src, dst, result, d.Size)

	if d.Extra == 0 {
		c.setRegD(d.DstReg, result, d.Size)
	} else {
		return c.writeMemory(c.getRegA(d.DstReg), result, d.Size)
	}

	return nil
}

func execNEG(c *CPU, d DecodedOpcode) error {
	dstEA, err := c.decodeEA(d.DstMode, d.DstReg, d.Size)
	if err != nil {
		return err
	}
	dst, err := c.readEA(dstEA)
	if err != nil {
		return err
	}

	result := uint32(0) - dst
	c.setSubFlags(dst, 0, result, d.Size)
	return c.writeEA(dstEA, result)
}

func execNEGX(c *CPU, d DecodedOpcode) error {
	dstEA, err := c.decodeEA(d.DstMode, d.DstReg, d.Size)
	if err != nil {
		return err
	}
	dst, err := c.readEA(dstEA)
	if err != nil {
		return err
	}

	x := uint32(c.Flags.X)
	result := uint32(0) - dst - x
	c.setSubXFlags(dst, 0, result, d.Size)
	return c.writeEA(dstEA, result)
}

func execCLR(c *CPU, d DecodedOpcode) error {
	dstEA, err := c.decodeEA(d.DstMode, d.DstReg, d.Size)
	if err != nil {
		return err
	}

	c.Flags.N = 0
	c.Flags.Z = 1
	c.Flags.V = 0
	c.Flags.C = 0
	return c.writeEA(dstEA, 0)
}

func execCMP(c *CPU, d DecodedOpcode) error {
	srcEA, err := c.decodeEA(d.SrcMode, d.SrcReg, d.Size)
	if err != nil {
		return err
	}
	src, err := c.readEA(srcEA)
	if err != nil {
		return err
	}
	dst := c.getRegD(d.DstReg, d.Size)

	result := dst - src
	c.setCmpFlags(src, dst, result, d.Size)
	return nil
}

func execCMPA(c *CPU, d DecodedOpcode) error {
	srcEA, err := c.decodeEA(d.SrcMode, d.SrcReg, d.Size)
	if err != nil {
		return err
	}
	src, err := c.readEA(srcEA)
	if err != nil {
		return err
	}
	src = signExtend(src, d.Size)
	dst := c.getRegA(d.DstReg)

	result := dst - src
	c.setCmpFlags(src, dst, result, SizeLong)
	return nil
}

func execCMPI(c *CPU, d DecodedOpcode) error {
	imm := c.readImmediate(d.Size)

	dstEA, err := c.decodeEA(d.DstMode, d.DstReg, d.Size)
	if err != nil {
		return err
	}
	dst, err := c.readEA(dstEA)
	if err != nil {
		return err
	}

	result := dst - imm
	c.setCmpFlags(imm, dst, result, d.Size)
	return nil
}

func execCMPM(c *CPU, d DecodedOpcode) error {
	srcAddr := c.getRegA(d.SrcReg)
	src, err := c.readMemory(srcAddr, d.Size)
	if err != nil {
		return err
	}
	c.setRegA(d.SrcReg, srcAddr+uint32(d.Size))

	dstAddr := c.getRegA(d.DstReg)
	dst, err := c.readMemory(dstAddr, d.Size)
	if err != nil {
		return err
	}
	c.setRegA(d.DstReg, dstAddr+uint32(d.Size))

	result := dst - src
	c.setCmpFlags(src, dst, result, d.Size)
	return nil
}

func execAND(c *CPU, d DecodedOpcode) error {
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
	dst, err := c.readEA(dstEA)
	if err != nil {
		return err
	}

	result := src & dst
	c.setLogicFlags(result, d.Size)
	return c.writeEA(dstEA, result)
}

func execANDI(c *CPU, d DecodedOpcode) error {
	imm := c.readImmediate(d.Size)

	// ANDI to CCR/SR special cases.
	if d.DstMode == 7 && d.DstReg == 4 {
		if d.Size == SizeByte {
			c.SetCCR(c.GetCCR() & uint8(imm))
			return nil
		}
		if !c.IsSupervisor() {
			return c.processException(VectorPrivilege)
		}
		c.SetSR(c.GetSR() & uint16(imm))
		return nil
	}

	dstEA, err := c.decodeEA(d.DstMode, d.DstReg, d.Size)
	if err != nil {
		return err
	}
	dst, err := c.readEA(dstEA)
	if err != nil {
		return err
	}

	result := imm & dst
	c.setLogicFlags(result, d.Size)
	return c.writeEA(dstEA, result)
}

func execOR(c *CPU, d DecodedOpcode) error {
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
	dst, err := c.readEA(dstEA)
	if err != nil {
		return err
	}

	result := src | dst
	c.setLogicFlags(result, d.Size)
	return c.writeEA(dstEA, result)
}

func execORI(c *CPU, d DecodedOpcode) error {
	imm := c.readImmediate(d.Size)

	// ORI to CCR/SR special cases.
	if d.DstMode == 7 && d.DstReg == 4 {
		if d.Size == SizeByte {
			c.SetCCR(c.GetCCR() | uint8(imm))
			return nil
		}
		if !c.IsSupervisor() {
			return c.processException(VectorPrivilege)
		}
		c.SetSR(c.GetSR() | uint16(imm))
		return nil
	}

	dstEA, err := c.decodeEA(d.DstMode, d.DstReg, d.Size)
	if err != nil {
		return err
	}
	dst, err := c.readEA(dstEA)
	if err != nil {
		return err
	}

	result := imm | dst
	c.setLogicFlags(result, d.Size)
	return c.writeEA(dstEA, result)
}

func execEOR(c *CPU, d DecodedOpcode) error {
	src := c.getRegD(d.SrcReg, d.Size)

	dstEA, err := c.decodeEA(d.DstMode, d.DstReg, d.Size)
	if err != nil {
		return err
	}
	dst, err := c.readEA(dstEA)
	if err != nil {
		return err
	}

	result := src ^ dst
	c.setLogicFlags(result, d.Size)
	return c.writeEA(dstEA, result)
}

func execEORI(c *CPU, d DecodedOpcode) error {
	imm := c.readImmediate(d.Size)

	// EORI to CCR/SR special cases.
	if d.DstMode == 7 && d.DstReg == 4 {
		if d.Size == SizeByte {
			c.SetCCR(c.GetCCR() ^ uint8(imm))
			return nil
		}
		if !c.IsSupervisor() {
			return c.processException(VectorPrivilege)
		}
		c.SetSR(c.GetSR() ^ uint16(imm))
		return nil
	}

	dstEA, err := c.decodeEA(d.DstMode, d.DstReg, d.Size)
	if err != nil {
		return err
	}
	dst, err := c.readEA(dstEA)
	if err != nil {
		return err
	}

	result := imm ^ dst
	c.setLogicFlags(result, d.Size)
	return c.writeEA(dstEA, result)
}

func execNOT(c *CPU, d DecodedOpcode) error {
	dstEA, err := c.decodeEA(d.DstMode, d.DstReg, d.Size)
	if err != nil {
		return err
	}
	dst, err := c.readEA(dstEA)
	if err != nil {
		return err
	}

	result := ^dst
	c.setLogicFlags(result, d.Size)
	return c.writeEA(dstEA, result)
}

func execEXT(c *CPU, d DecodedOpcode) error {
	var result uint32
	if d.Size == SizeWord {
		result = uint32(int16(int8(c.D[d.DstReg])))
		c.setRegD(d.DstReg, result, SizeWord)
	} else {
		result = uint32(int32(int16(c.D[d.DstReg])))
		c.setRegD(d.DstReg, result, SizeLong)
	}
	c.setLogicFlags(result, d.Size)
	return nil
}

func execTST(c *CPU, d DecodedOpcode) error {
	dstEA, err := c.decodeEA(d.DstMode, d.DstReg, d.Size)
	if err != nil {
		return err
	}
	val, err := c.readEA(dstEA)
	if err != nil {
		return err
	}
	c.setLogicFlags(val, d.Size)
	return nil
}

func execMULU(c *CPU, d DecodedOpcode) error {
	srcEA, err := c.decodeEA(d.SrcMode, d.SrcReg, SizeWord)
	if err != nil {
		return err
	}
	src, err := c.readEA(srcEA)
	if err != nil {
		return err
	}

	dst := c.D[d.DstReg] & 0xFFFF
	result := dst * src
	c.D[d.DstReg] = result

	c.setFlagN(result, SizeLong)
	c.setFlagZ(result, SizeLong)
	c.Flags.V = 0
	c.Flags.C = 0
	return nil
}

func execMULS(c *CPU, d DecodedOpcode) error {
	srcEA, err := c.decodeEA(d.SrcMode, d.SrcReg, SizeWord)
	if err != nil {
		return err
	}
	src, err := c.readEA(srcEA)
	if err != nil {
		return err
	}

	dst := int32(int16(c.D[d.DstReg]))
	result := dst * int32(int16(src))
	c.D[d.DstReg] = uint32(result)

	c.setFlagN(uint32(result), SizeLong)
	c.setFlagZ(uint32(result), SizeLong)
	c.Flags.V = 0
	c.Flags.C = 0
	return nil
}

func execDIVU(c *CPU, d DecodedOpcode) error {
	srcEA, err := c.decodeEA(d.SrcMode, d.SrcReg, SizeWord)
	if err != nil {
		return err
	}
	src, err := c.readEA(srcEA)
	if err != nil {
		return err
	}

	if src == 0 {
		return c.processException(VectorDivZero)
	}

	dividend := c.D[d.DstReg]
	quotient := dividend / src
	remainder := dividend % src

	if quotient > 0xFFFF {
		c.Flags.V = 1
		c.Flags.C = 0
		return nil
	}

	c.D[d.DstReg] = (remainder << 16) | (quotient & 0xFFFF)
	c.setFlagN(quotient, SizeWord)
	c.setFlagZ(quotient, SizeWord)
	c.Flags.V = 0
	c.Flags.C = 0
	return nil
}

func execDIVS(c *CPU, d DecodedOpcode) error {
	srcEA, err := c.decodeEA(d.SrcMode, d.SrcReg, SizeWord)
	if err != nil {
		return err
	}
	src, err := c.readEA(srcEA)
	if err != nil {
		return err
	}

	if src == 0 {
		return c.processException(VectorDivZero)
	}

	dividend := int32(c.D[d.DstReg])
	divisor := int32(int16(src))
	quotient := dividend / divisor
	remainder := dividend % divisor

	if quotient > 32767 || quotient < -32768 {
		c.Flags.V = 1
		c.Flags.C = 0
		return nil
	}

	c.D[d.DstReg] = (uint32(int16(remainder)) << 16) | (uint32(int16(quotient)) & 0xFFFF)
	c.setFlagN(uint32(int16(quotient)), SizeWord)
	c.setFlagZ(uint32(int16(quotient)), SizeWord)
	c.Flags.V = 0
	c.Flags.C = 0
	return nil
}

func execABCD(c *CPU, d DecodedOpcode) error {
	var src, dst uint8

	if d.Extra&0x8 == 0 {
		src = uint8(c.D[d.SrcReg])
		dst = uint8(c.D[d.DstReg])
	} else {
		c.setRegA(d.SrcReg, c.getRegA(d.SrcReg)-1)
		src = c.bus.Read(c.getRegA(d.SrcReg))
		c.setRegA(d.DstReg, c.getRegA(d.DstReg)-1)
		dst = c.bus.Read(c.getRegA(d.DstReg))
	}

	x := c.Flags.X
	low := (dst & 0x0F) + (src & 0x0F) + x
	if low > 9 {
		low += 6
	}
	high := uint16(dst>>4) + uint16(src>>4)
	if low > 0x0F {
		high++
	}
	if high > 9 {
		high += 6
	}

	result := uint8((high << 4) | uint16(low&0x0F))
	setFlag(&c.Flags.X, high > 0x0F)
	setFlag(&c.Flags.C, high > 0x0F)
	if result != 0 {
		c.Flags.Z = 0
	}

	if d.Extra&0x8 == 0 {
		c.D[d.DstReg] = (c.D[d.DstReg] & 0xFFFFFF00) | uint32(result)
	} else {
		c.bus.Write(c.getRegA(d.DstReg), result)
	}

	return nil
}

func execSBCD(c *CPU, d DecodedOpcode) error {
	var src, dst uint8

	if d.Extra&0x8 == 0 {
		src = uint8(c.D[d.SrcReg])
		dst = uint8(c.D[d.DstReg])
	} else {
		c.setRegA(d.SrcReg, c.getRegA(d.SrcReg)-1)
		src = c.bus.Read(c.getRegA(d.SrcReg))
		c.setRegA(d.DstReg, c.getRegA(d.DstReg)-1)
		dst = c.bus.Read(c.getRegA(d.DstReg))
	}

	x := c.Flags.X
	low := int16(dst&0x0F) - int16(src&0x0F) - int16(x)
	borrow := uint16(0)
	if low < 0 {
		low += 10
		borrow = 1
	}
	high := int16(dst>>4) - int16(src>>4) - int16(borrow)
	if high < 0 {
		high += 10
		setFlag(&c.Flags.X, true)
		setFlag(&c.Flags.C, true)
	} else {
		setFlag(&c.Flags.X, false)
		setFlag(&c.Flags.C, false)
	}

	result := uint8((uint16(high) << 4) | uint16(low&0x0F))
	if result != 0 {
		c.Flags.Z = 0
	}

	if d.Extra&0x8 == 0 {
		c.D[d.DstReg] = (c.D[d.DstReg] & 0xFFFFFF00) | uint32(result)
	} else {
		c.bus.Write(c.getRegA(d.DstReg), result)
	}

	return nil
}

func execNBCD(c *CPU, d DecodedOpcode) error {
	dstEA, err := c.decodeEA(d.DstMode, d.DstReg, SizeByte)
	if err != nil {
		return err
	}
	dst, err := c.readEA(dstEA)
	if err != nil {
		return err
	}

	src := uint8(dst)
	x := c.Flags.X

	low := int16(0) - int16(src&0x0F) - int16(x)
	borrow := uint16(0)
	if low < 0 {
		low += 10
		borrow = 1
	}
	high := int16(0) - int16(src>>4) - int16(borrow)
	if high < 0 {
		high += 10
		setFlag(&c.Flags.X, true)
		setFlag(&c.Flags.C, true)
	} else {
		setFlag(&c.Flags.X, false)
		setFlag(&c.Flags.C, false)
	}

	result := uint8((uint16(high) << 4) | uint16(low&0x0F))
	if result != 0 {
		c.Flags.Z = 0
	}

	return c.writeEA(dstEA, uint32(result))
}

// setAddFlags sets flags for ADD-type operations.
func (c *CPU) setAddFlags(src, dst, result uint32, size OperandSize) {
	msb := msbMask(size)
	sm := src & msb
	dm := dst & msb
	rm := result & msb

	c.setFlagN(result, size)
	c.setFlagZ(result, size)
	setFlag(&c.Flags.V, (sm == dm) && (rm != sm))
	setFlag(&c.Flags.C, maskValue(result, size) < maskValue(src, size))
	c.Flags.X = c.Flags.C
}

// setAddXFlags sets flags for ADDX-type operations (Z flag only cleared, never set).
func (c *CPU) setAddXFlags(src, dst, result uint32, size OperandSize) {
	msb := msbMask(size)
	sm := src & msb
	dm := dst & msb
	rm := result & msb

	c.setFlagN(result, size)
	if maskValue(result, size) != 0 {
		c.Flags.Z = 0
	}
	setFlag(&c.Flags.V, (sm == dm) && (rm != sm))
	setFlag(&c.Flags.C, maskValue(result, size) < maskValue(src, size))
	c.Flags.X = c.Flags.C
}

// setSubFlags sets flags for SUB-type operations.
func (c *CPU) setSubFlags(src, dst, result uint32, size OperandSize) {
	msb := msbMask(size)
	sm := src & msb
	dm := dst & msb
	rm := result & msb

	c.setFlagN(result, size)
	c.setFlagZ(result, size)
	setFlag(&c.Flags.V, (sm != dm) && (rm != dm))
	setFlag(&c.Flags.C, maskValue(src, size) > maskValue(dst, size))
	c.Flags.X = c.Flags.C
}

// setSubXFlags sets flags for SUBX-type operations (Z flag only cleared, never set).
func (c *CPU) setSubXFlags(src, dst, result uint32, size OperandSize) {
	msb := msbMask(size)
	sm := src & msb
	dm := dst & msb
	rm := result & msb

	c.setFlagN(result, size)
	if maskValue(result, size) != 0 {
		c.Flags.Z = 0
	}
	setFlag(&c.Flags.V, (sm != dm) && (rm != dm))
	setFlag(&c.Flags.C, maskValue(src, size) > maskValue(dst, size))
	c.Flags.X = c.Flags.C
}

// setCmpFlags sets flags for CMP-type operations (X not affected).
func (c *CPU) setCmpFlags(src, dst, result uint32, size OperandSize) {
	msb := msbMask(size)
	sm := src & msb
	dm := dst & msb
	rm := result & msb

	c.setFlagN(result, size)
	c.setFlagZ(result, size)
	setFlag(&c.Flags.V, (sm != dm) && (rm != dm))
	setFlag(&c.Flags.C, maskValue(src, size) > maskValue(dst, size))
}

// setLogicFlags sets flags for logic operations (AND, OR, EOR, NOT, etc.).
func (c *CPU) setLogicFlags(result uint32, size OperandSize) {
	c.setFlagN(result, size)
	c.setFlagZ(result, size)
	c.Flags.V = 0
	c.Flags.C = 0
}
