package z80

// OpcodeMap provides a better mapping for Z80 opcodes that solves the duplicate key problem.
// Instead of using Instruction.Addressing map which causes conflicts, we create a direct
// opcode-to-info mapping with parameter differentiation.
type OpcodeMap struct {
	opcodeToInfo map[byte]OpcodeDetail
}

// OpcodeDetail extends OpcodeInfo with parameter information for assembler tooling.
type OpcodeDetail struct {
	OpcodeInfo
	Instruction *Instruction      // Reference to the instruction
	Addressing  AddressingMode    // Addressing mode used
	Params      []string          // Specific parameters (e.g., ["bc"] for DEC BC)
}

// NewOpcodeMap creates a complete opcode mapping for Z80 with parameter differentiation.
func NewOpcodeMap() *OpcodeMap {
	om := &OpcodeMap{
		opcodeToInfo: make(map[byte]OpcodeDetail),
	}
	om.buildOpcodeMap()
	return om
}

// GetOpcodeByBytes returns opcode information for a given opcode byte.
func (om *OpcodeMap) GetOpcodeByBytes(opcode byte) *OpcodeDetail {
	if info, exists := om.opcodeToInfo[opcode]; exists {
		return &info
	}
	return nil
}

// GetOpcodeByInstructionAndParams finds an opcode by instruction name and parameters.
// This is useful for assembler tooling.
func (om *OpcodeMap) GetOpcodeByInstructionAndParams(instruction string, addressing AddressingMode, params []string) *OpcodeDetail {
	for _, info := range om.opcodeToInfo {
		if info.Instruction.Name == instruction && info.Addressing == addressing {
			if paramsMatch(info.Params, params) {
				return &info
			}
		}
	}
	return nil
}

// GetInstructionVariants returns all variants of a given instruction.
func (om *OpcodeMap) GetInstructionVariants(instruction string) []OpcodeDetail {
	var variants []OpcodeDetail
	for _, info := range om.opcodeToInfo {
		if info.Instruction.Name == instruction {
			variants = append(variants, info)
		}
	}
	return variants
}

// buildOpcodeMap populates the opcode map with detailed parameter information.
func (om *OpcodeMap) buildOpcodeMap() {
	// Build from existing Opcodes array but add parameter differentiation
	for opcode, opcodeInfo := range Opcodes {
		if opcodeInfo.Instruction == nil {
			continue // Skip empty slots
		}
		
		params := om.getParamsForOpcode(byte(opcode), opcodeInfo)
		
		om.opcodeToInfo[byte(opcode)] = OpcodeDetail{
			OpcodeInfo: OpcodeInfo{
				Opcode: byte(opcode),
				Size:   opcodeInfo.Size,
				Cycles: opcodeInfo.Timing,
			},
			Instruction: opcodeInfo.Instruction,
			Addressing:  opcodeInfo.Addressing,
			Params:      params,
		}
	}
}

// getParamsForOpcode determines the specific parameters for an opcode.
// This differentiates between opcodes that have the same instruction+addressing.
func (om *OpcodeMap) getParamsForOpcode(opcode byte, opcodeInfo Opcode) []string {
	instruction := opcodeInfo.Instruction.Name
	addressing := opcodeInfo.Addressing
	
	switch instruction {
	case "inc":
		if addressing == RegisterAddressing {
			// 16-bit register increment
			switch opcode {
			case 0x03: return []string{"bc"}
			case 0x13: return []string{"de"}
			case 0x23: return []string{"hl"}
			case 0x33: return []string{"sp"}
			default:
				// 8-bit register increment - handle specific opcodes
				switch opcode {
				case 0x04: return []string{"b"}
				case 0x0C: return []string{"c"}
				case 0x14: return []string{"d"}
				case 0x1C: return []string{"e"}
				case 0x24: return []string{"h"}
				case 0x2C: return []string{"l"}
				case 0x3C: return []string{"a"}
				case 0x34: return []string{"(hl)"} // Actually handled by RegisterIndirectAddressing case above
				}
			}
		} else if addressing == RegisterIndirectAddressing {
			return []string{"(hl)"}
		}
		
	case "dec":
		if addressing == RegisterAddressing {
			// 16-bit register decrement  
			switch opcode {
			case 0x0B: return []string{"bc"}
			case 0x1B: return []string{"de"}
			case 0x2B: return []string{"hl"}
			case 0x3B: return []string{"sp"}
			default:
				// 8-bit register decrement - handle specific opcodes
				switch opcode {
				case 0x05: return []string{"b"}
				case 0x0D: return []string{"c"}
				case 0x15: return []string{"d"}
				case 0x1D: return []string{"e"}
				case 0x25: return []string{"h"}
				case 0x2D: return []string{"l"}
				case 0x3D: return []string{"a"}
				case 0x35: return []string{"(hl)"} // Actually handled by RegisterIndirectAddressing case above
				}
			}
		} else if addressing == RegisterIndirectAddressing {
			return []string{"(hl)"}
		}
		
	case "ld":
		if addressing == ImmediateAddressing {
			if opcodeInfo.Size == 3 {
				// 16-bit immediate load
				switch opcode {
				case 0x01: return []string{"bc", "nn"}
				case 0x11: return []string{"de", "nn"}
				case 0x21: return []string{"hl", "nn"}
				case 0x31: return []string{"sp", "nn"}
				}
			} else if opcodeInfo.Size == 2 {
				// 8-bit immediate load - need to handle specific opcodes
				switch opcode {
				case 0x06: return []string{"b", "n"}
				case 0x0E: return []string{"c", "n"} 
				case 0x16: return []string{"d", "n"}
				case 0x1E: return []string{"e", "n"}
				case 0x26: return []string{"h", "n"}
				case 0x2E: return []string{"l", "n"}
				case 0x3E: return []string{"a", "n"}
				}
			}
		} else if addressing == RegisterAddressing {
			if opcode >= 0x40 && opcode <= 0x7F && opcode != 0x76 {
				// Register-to-register loads
				dst := get8BitRegFromOpcode((opcode-0x40)/8, 0)
				src := get8BitRegFromOpcode((opcode-0x40)%8, 0)
				return []string{dst, src}
			} else if opcode == 0xF9 {
				return []string{"sp", "hl"}
			}
		} else if addressing == RegisterIndirectAddressing {
			switch opcode {
			case 0x02: return []string{"(bc)", "a"}
			case 0x0A: return []string{"a", "(bc)"}
			case 0x12: return []string{"(de)", "a"}
			case 0x1A: return []string{"a", "(de)"}
			case 0x36: return []string{"(hl)", "n"}
			default:
				if opcode >= 0x70 && opcode <= 0x77 {
					src := get8BitRegFromOpcode(opcode, 0x70)
					return []string{"(hl)", src}
				} else if opcode >= 0x46 && opcode <= 0x7E && (opcode-0x46)%8 == 6 {
					dst := get8BitRegFromOpcode((opcode-0x46)/8, 0)
					return []string{dst, "(hl)"}
				}
			}
		} else if addressing == ExtendedAddressing {
			switch opcode {
			case 0x22: return []string{"(nn)", "hl"}
			case 0x2A: return []string{"hl", "(nn)"}
			case 0x32: return []string{"(nn)", "a"}
			case 0x3A: return []string{"a", "(nn)"}
			}
		}
		
	case "add":
		if addressing == RegisterAddressing {
			if opcode >= 0x80 && opcode <= 0x87 {
				// ADD A,r
				src := get8BitRegFromOpcode(opcode, 0x80)
				return []string{"a", src}
			} else {
				// ADD HL,rr
				switch opcode {
				case 0x09: return []string{"hl", "bc"}
				case 0x19: return []string{"hl", "de"}
				case 0x29: return []string{"hl", "hl"}
				case 0x39: return []string{"hl", "sp"}
				}
			}
		} else if addressing == ImmediateAddressing {
			return []string{"a", "n"}
		} else if addressing == RegisterIndirectAddressing {
			return []string{"a", "(hl)"}
		}
		
	// Add more cases for other instructions as needed
	case "push", "pop":
		regPair := getRegPairFromStackOpcode(opcode)
		return []string{regPair}
		
	case "rst":
		address := getRstAddress(opcode)
		return []string{address}
	}
	
	// Default: no specific parameters
	return nil
}

// Helper functions
func get8BitRegFromOpcode(offset byte, base byte) string {
	registers := []string{"b", "c", "d", "e", "h", "l", "(hl)", "a"}
	index := (offset - base) % 8
	if int(index) < len(registers) {
		return registers[index]
	}
	return "?"
}

func getRegPairFromStackOpcode(opcode byte) string {
	regPairs := map[byte]string{
		0xC1: "bc", 0xC5: "bc",
		0xD1: "de", 0xD5: "de",
		0xE1: "hl", 0xE5: "hl",
		0xF1: "af", 0xF5: "af",
	}
	if reg, exists := regPairs[opcode]; exists {
		return reg
	}
	return "?"
}

func getRstAddress(opcode byte) string {
	addresses := map[byte]string{
		0xC7: "00h", 0xCF: "08h", 0xD7: "10h", 0xDF: "18h",
		0xE7: "20h", 0xEF: "28h", 0xF7: "30h", 0xFF: "38h",
	}
	if addr, exists := addresses[opcode]; exists {
		return addr
	}
	return "?"
}

func paramsMatch(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}