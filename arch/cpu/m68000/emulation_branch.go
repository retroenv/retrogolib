package m68000

// Branch and flow control instructions: Bcc, BRA, BSR, DBcc, Scc, JMP, JSR, NOP.

// Condition code evaluation for Bcc/Scc/DBcc.
// Conditions are encoded in bits 11-8 of the opcode word.
func (c *CPU) evaluateCondition(cond uint16) bool {
	switch {
	case cond == 0: // T (true)
		return true
	case cond == 1: // F (false)
		return false
	case cond <= 5:
		return c.evaluateConditionCarryZero(cond)
	case cond <= 11:
		return c.evaluateConditionSingle(cond)
	default:
		return c.evaluateConditionCompound(cond)
	}
}

// evaluateConditionCarryZero evaluates C/Z flag conditions (2-5).
func (c *CPU) evaluateConditionCarryZero(cond uint16) bool {
	switch cond {
	case 2: // HI: !C && !Z
		return c.Flags.C == 0 && c.Flags.Z == 0
	case 3: // LS: C || Z
		return c.Flags.C != 0 || c.Flags.Z != 0
	case 4: // CC: !C
		return c.Flags.C == 0
	default: // 5: CS: C
		return c.Flags.C != 0
	}
}

// evaluateConditionSingle evaluates single-flag conditions (6-11).
func (c *CPU) evaluateConditionSingle(cond uint16) bool {
	switch cond {
	case 6: // NE: !Z
		return c.Flags.Z == 0
	case 7: // EQ: Z
		return c.Flags.Z != 0
	case 8: // VC: !V
		return c.Flags.V == 0
	case 9: // VS: V
		return c.Flags.V != 0
	case 10: // PL: !N
		return c.Flags.N == 0
	default: // 11: MI: N
		return c.Flags.N != 0
	}
}

// evaluateConditionCompound evaluates compound N/V/Z conditions (12-15).
func (c *CPU) evaluateConditionCompound(cond uint16) bool {
	switch cond {
	case 12: // GE: N == V
		return c.Flags.N == c.Flags.V
	case 13: // LT: N != V
		return c.Flags.N != c.Flags.V
	case 14: // GT: !Z && N == V
		return c.Flags.Z == 0 && c.Flags.N == c.Flags.V
	default: // 15: LE: Z || N != V
		return c.Flags.Z != 0 || c.Flags.N != c.Flags.V
	}
}

func execBcc(c *CPU, d DecodedOpcode) error {
	if !c.evaluateCondition(d.Extra) {
		// Branch not taken. If short branch, PC is already past opcode word.
		// If long branch (disp==0), we need to skip the extension word.
		if d.DstReg == 0 {
			c.PC += 2 // Skip the 16-bit displacement extension word.
		}
		return nil
	}

	return c.takeBranch(d)
}

func execBRA(c *CPU, d DecodedOpcode) error {
	return c.takeBranch(d)
}

func execBSR(c *CPU, d DecodedOpcode) error {
	// Calculate target first, then push return address.
	pcBeforeBranch := c.PC
	disp := int32(int8(d.DstReg))

	if d.DstReg == 0 {
		// 16-bit displacement.
		disp = int32(int16(c.readWord()))
	}

	// Push return address (after all extension words).
	c.push32(c.PC)

	// Branch to target. PC base is the address of the extension word.
	c.PC = uint32(int32(pcBeforeBranch) + disp)
	return nil
}

func execDBcc(c *CPU, d DecodedOpcode) error {
	disp := int16(c.readWord())

	if c.evaluateCondition(d.Extra) {
		// Condition true: no loop, just continue.
		return nil
	}

	// Decrement counter (low word of Dn).
	counter := int16(c.D[d.DstReg]) - 1
	c.D[d.DstReg] = (c.D[d.DstReg] & 0xFFFF0000) | (uint32(counter) & 0xFFFF)

	if counter == -1 {
		// Counter expired: fall through.
		return nil
	}

	// Branch. PC base is the extension word address (PC - 2 since we already read it).
	c.PC = uint32(int32(c.PC-2) + int32(disp))
	return nil
}

func execScc(c *CPU, d DecodedOpcode) error {
	dstEA, err := c.decodeEA(d.DstMode, d.DstReg, SizeByte)
	if err != nil {
		return err
	}

	var result uint32
	if c.evaluateCondition(d.Extra) {
		result = 0xFF
	}
	return c.writeEA(dstEA, result)
}

func execJMP(c *CPU, d DecodedOpcode) error {
	ea, err := c.decodeEA(d.DstMode, d.DstReg, SizeLong)
	if err != nil {
		return err
	}
	c.PC = ea.Address
	return nil
}

func execJSR(c *CPU, d DecodedOpcode) error {
	ea, err := c.decodeEA(d.DstMode, d.DstReg, SizeLong)
	if err != nil {
		return err
	}
	c.push32(c.PC)
	c.PC = ea.Address
	return nil
}

// takeBranch takes a branch with the displacement encoded in the opcode.
func (c *CPU) takeBranch(d DecodedOpcode) error {
	// PC currently points after the opcode word.
	// For short branch (8-bit disp): base is PC after opcode word.
	// For long branch (disp==0): base is PC before extension word, disp is 16-bit.
	pcBase := c.PC
	disp := int32(int8(d.DstReg))

	if d.DstReg == 0 {
		// 16-bit displacement follows the opcode word.
		disp = int32(int16(c.readWord()))
		pcBase = c.PC - 2 // Base is the extension word address.
	}

	c.PC = uint32(int32(pcBase) + disp)
	return nil
}
