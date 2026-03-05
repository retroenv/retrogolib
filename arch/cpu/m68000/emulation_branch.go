package m68000

// Branch and flow control instructions: Bcc, BRA, BSR, DBcc, Scc, JMP, JSR, NOP.

// Condition code evaluation for Bcc/Scc/DBcc.
// Conditions are encoded in bits 11-8 of the opcode word.
func (c *CPU) evaluateCondition(cond uint16) bool {
	switch cond {
	case 0: // T (true)
		return true
	case 1: // F (false)
		return false
	case 2: // HI (higher): !C && !Z
		return c.Flags.C == 0 && c.Flags.Z == 0
	case 3: // LS (lower or same): C || Z
		return c.Flags.C != 0 || c.Flags.Z != 0
	case 4: // CC (carry clear): !C
		return c.Flags.C == 0
	case 5: // CS (carry set): C
		return c.Flags.C != 0
	case 6: // NE (not equal): !Z
		return c.Flags.Z == 0
	case 7: // EQ (equal): Z
		return c.Flags.Z != 0
	case 8: // VC (overflow clear): !V
		return c.Flags.V == 0
	case 9: // VS (overflow set): V
		return c.Flags.V != 0
	case 10: // PL (plus): !N
		return c.Flags.N == 0
	case 11: // MI (minus): N
		return c.Flags.N != 0
	case 12: // GE (greater or equal): (N && V) || (!N && !V)
		return c.Flags.N == c.Flags.V
	case 13: // LT (less than): (N && !V) || (!N && V)
		return c.Flags.N != c.Flags.V
	case 14: // GT (greater than): (N && V && !Z) || (!N && !V && !Z)
		return c.Flags.Z == 0 && c.Flags.N == c.Flags.V
	case 15: // LE (less or equal): Z || (N && !V) || (!N && V)
		return c.Flags.Z != 0 || c.Flags.N != c.Flags.V
	default:
		return false
	}
}

func (c *CPU) execBcc(d DecodedOpcode) error {
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

func (c *CPU) execBRA(d DecodedOpcode) error {
	return c.takeBranch(d)
}

func (c *CPU) execBSR(d DecodedOpcode) error {
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

func (c *CPU) execDBcc(d DecodedOpcode) error {
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

func (c *CPU) execScc(d DecodedOpcode) error {
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

func (c *CPU) execJMP(d DecodedOpcode) error {
	ea, err := c.decodeEA(d.DstMode, d.DstReg, SizeLong)
	if err != nil {
		return err
	}
	c.PC = ea.Address
	return nil
}

func (c *CPU) execJSR(d DecodedOpcode) error {
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
