package cpu68000

import "math/bits"

// Core ALU operations: ADD, ADDA, ADDI, ADDQ, ADDX, SUB, SUBA, SUBI, SUBQ, SUBX,
// NEG, NEGX, CLR, CMP, CMPA, CMPI, CMPM, AND, ANDI, OR, ORI, EOR, EORI, NOT, EXT,
// TST, MULU, MULS, DIVU, DIVS, ABCD, SBCD, NBCD.

// setAddFlags sets flags for ADD-type operations.
func (cpu *CPU) setAddFlags(src, dst, result uint32, size OperandSize) {
	msb := msbMask(size)
	sm := src & msb
	dm := dst & msb
	rm := result & msb

	cpu.setFlagN(result, size)
	cpu.setFlagZ(result, size)
	setFlag(&cpu.Flags.V, (sm == dm) && (rm != sm))
	setFlag(&cpu.Flags.C, maskValue(result, size) < maskValue(src, size))
	cpu.Flags.X = cpu.Flags.C
}

// setAddXFlags sets flags for ADDX-type operations (Z flag only cleared, never set).
func (cpu *CPU) setAddXFlags(src, dst, result uint32, size OperandSize) {
	msb := msbMask(size)
	sm := src & msb
	dm := dst & msb
	rm := result & msb

	cpu.setFlagN(result, size)
	if maskValue(result, size) != 0 {
		cpu.Flags.Z = 0
	}
	setFlag(&cpu.Flags.V, (sm == dm) && (rm != sm))
	carry := uint64(src)+uint64(dst)+uint64(cpu.Flags.X) > uint64(maskValue(0xFFFFFFFF, size))
	setFlag(&cpu.Flags.C, carry)
	cpu.Flags.X = cpu.Flags.C
}

// setSubFlags sets flags for SUB-type operations.
func (cpu *CPU) setSubFlags(src, dst, result uint32, size OperandSize) {
	msb := msbMask(size)
	sm := src & msb
	dm := dst & msb
	rm := result & msb

	cpu.setFlagN(result, size)
	cpu.setFlagZ(result, size)
	setFlag(&cpu.Flags.V, (sm != dm) && (rm != dm))
	setFlag(&cpu.Flags.C, maskValue(src, size) > maskValue(dst, size))
	cpu.Flags.X = cpu.Flags.C
}

// setSubXFlags sets flags for SUBX-type operations (Z flag only cleared, never set).
func (cpu *CPU) setSubXFlags(src, dst, result uint32, size OperandSize) {
	msb := msbMask(size)
	sm := src & msb
	dm := dst & msb
	rm := result & msb

	cpu.setFlagN(result, size)
	if maskValue(result, size) != 0 {
		cpu.Flags.Z = 0
	}
	setFlag(&cpu.Flags.V, (sm != dm) && (rm != dm))
	borrow := uint64(maskValue(src, size))+uint64(cpu.Flags.X) > uint64(maskValue(dst, size))
	setFlag(&cpu.Flags.C, borrow)
	cpu.Flags.X = cpu.Flags.C
}

// setCmpFlags sets flags for CMP-type operations (X not affected).
func (cpu *CPU) setCmpFlags(src, dst, result uint32, size OperandSize) {
	msb := msbMask(size)
	sm := src & msb
	dm := dst & msb
	rm := result & msb

	cpu.setFlagN(result, size)
	cpu.setFlagZ(result, size)
	setFlag(&cpu.Flags.V, (sm != dm) && (rm != dm))
	setFlag(&cpu.Flags.C, maskValue(src, size) > maskValue(dst, size))
}

// setLogicFlags sets flags for logic operations (AND, OR, EOR, NOT, etc.).
func (cpu *CPU) setLogicFlags(result uint32, size OperandSize) {
	cpu.setFlagN(result, size)
	cpu.setFlagZ(result, size)
	cpu.Flags.V = 0
	cpu.Flags.C = 0
}

func (cpu *CPU) addDecimal(src, dst uint8) uint8 {
	extend := uint16(cpu.Flags.X)
	binary := uint16(src) + uint16(dst) + extend
	result := binary
	if uint16(src&15)+uint16(dst&15)+extend > 9 {
		result += 6
	}
	carry := result > 0x9F
	if carry {
		result += 0x60
	}
	cpu.setFlagN(uint32(result), SizeByte)
	setFlag(&cpu.Flags.V, binary&0x80 == 0 && result&0x80 != 0)
	setFlag(&cpu.Flags.C, carry)
	cpu.Flags.X = cpu.Flags.C
	if uint8(result) != 0 {
		cpu.Flags.Z = 0
	}
	return uint8(result)
}

func (cpu *CPU) subtractDecimal(src, dst uint8) uint8 {
	extend := int16(cpu.Flags.X)
	binary := int16(dst) - int16(src) - extend
	result := binary
	if int16(dst&15)-int16(src&15)-extend < 0 {
		result -= 6
	}
	borrow := result < 0
	if binary < 0 {
		result -= 0x60
	}
	cpu.setFlagN(uint32(result), SizeByte)
	setFlag(&cpu.Flags.V, binary&0x80 != 0 && result&0x80 == 0)
	setFlag(&cpu.Flags.C, borrow)
	cpu.Flags.X = cpu.Flags.C
	if uint8(result) != 0 {
		cpu.Flags.Z = 0
	}
	return uint8(result)
}

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
		sv, err := c.readPredecrement(d.SrcReg, d.Size)
		if err != nil {
			return err
		}
		src = sv

		dv, err := c.readPredecrement(d.DstReg, d.Size)
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
		sv, err := c.readPredecrement(d.SrcReg, d.Size)
		if err != nil {
			return err
		}
		src = sv

		dv, err := c.readPredecrement(d.DstReg, d.Size)
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

	// The original 68000 reads the destination before clearing it.
	if _, err := c.readEA(dstEA); err != nil {
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
	c.setRegA(d.SrcReg, srcAddr+incrementSize(d.SrcReg, d.Size))
	src, err := c.readMemory(srcAddr, d.Size)
	if err != nil {
		return err
	}

	dstAddr := c.getRegA(d.DstReg)
	c.setRegA(d.DstReg, dstAddr+incrementSize(d.DstReg, d.Size))
	dst, err := c.readMemory(dstAddr, d.Size)
	if err != nil {
		return err
	}

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
	if d.DstMode == 7 && d.DstReg == 4 && d.Size != SizeByte && !c.IsSupervisor() {
		return c.processException(VectorPrivilege)
	}

	imm := c.readImmediate(d.Size)

	// ANDI to CCR/SR special cases.
	if d.DstMode == 7 && d.DstReg == 4 {
		if d.Size == SizeByte {
			c.SetCCR(c.GetCCR() & uint8(imm))
			return nil
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
	if d.DstMode == 7 && d.DstReg == 4 && d.Size != SizeByte && !c.IsSupervisor() {
		return c.processException(VectorPrivilege)
	}

	imm := c.readImmediate(d.Size)

	// ORI to CCR/SR special cases.
	if d.DstMode == 7 && d.DstReg == 4 {
		if d.Size == SizeByte {
			c.SetCCR(c.GetCCR() | uint8(imm))
			return nil
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
	if d.DstMode == 7 && d.DstReg == 4 && d.Size != SizeByte && !c.IsSupervisor() {
		return c.processException(VectorPrivilege)
	}

	imm := c.readImmediate(d.Size)

	// EORI to CCR/SR special cases.
	if d.DstMode == 7 && d.DstReg == 4 {
		if d.Size == SizeByte {
			c.SetCCR(c.GetCCR() ^ uint8(imm))
			return nil
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

	c.cycles += 38 + 2*uint64(bits.OnesCount16(uint16(src)))
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

	transitions := uint16(src) ^ uint16(src<<1)
	c.cycles += 38 + 2*uint64(bits.OnesCount16(transitions))
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
		c.setLogicFlags(c.D[d.DstReg]>>16, SizeWord)
		return c.processException(VectorDivZero)
	}

	dividend := c.D[d.DstReg]
	c.cycles += divideUnsignedCycles(dividend, uint16(src))
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
	c.cycles += divideSignedCycles(dividend, int16(src))
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
		c.setRegA(d.SrcReg, c.getRegA(d.SrcReg)-incrementSize(d.SrcReg, SizeByte))
		src = c.readByte(c.getRegA(d.SrcReg))
		c.setRegA(d.DstReg, c.getRegA(d.DstReg)-incrementSize(d.DstReg, SizeByte))
		dst = c.readByte(c.getRegA(d.DstReg))
	}

	result := c.addDecimal(src, dst)

	if d.Extra&0x8 == 0 {
		c.D[d.DstReg] = (c.D[d.DstReg] & 0xFFFFFF00) | uint32(result)
	} else {
		c.writeByte(c.getRegA(d.DstReg), result)
	}

	return nil
}

func execSBCD(c *CPU, d DecodedOpcode) error {
	var src, dst uint8

	if d.Extra&0x8 == 0 {
		src = uint8(c.D[d.SrcReg])
		dst = uint8(c.D[d.DstReg])
	} else {
		c.setRegA(d.SrcReg, c.getRegA(d.SrcReg)-incrementSize(d.SrcReg, SizeByte))
		src = c.readByte(c.getRegA(d.SrcReg))
		c.setRegA(d.DstReg, c.getRegA(d.DstReg)-incrementSize(d.DstReg, SizeByte))
		dst = c.readByte(c.getRegA(d.DstReg))
	}

	result := c.subtractDecimal(src, dst)

	if d.Extra&0x8 == 0 {
		c.D[d.DstReg] = (c.D[d.DstReg] & 0xFFFFFF00) | uint32(result)
	} else {
		c.writeByte(c.getRegA(d.DstReg), result)
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

	result := c.subtractDecimal(uint8(dst), 0)
	return c.writeEA(dstEA, uint32(result))
}
