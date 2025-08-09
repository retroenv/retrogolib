# Z80 OpcodeMap Solution

## Problem

The original Z80 `Instruction.Addressing` map approach had a fundamental issue: multiple opcodes shared the same `Instruction + AddressingMode` combination, causing key collisions.

**Example of the problem:**
- `DecReg16` with `RegisterAddressing` was used by:
  - 0x0B (DEC BC)  
  - 0x1B (DEC DE)
  - 0x2B (DEC HL) 
  - 0x3B (DEC SP)

This made it impossible for assembler tooling to uniquely identify which opcode to use for a given instruction + addressing combination.

## Solution: OpcodeMap

The `OpcodeMap` provides a better mapping system that solves the duplicate key problem:

### Key Features

1. **Direct opcode-to-info mapping**: `map[byte]OpcodeDetail` - Each opcode byte maps to complete information
2. **Parameter differentiation**: `OpcodeDetail.Params []string` - Specific parameters distinguish between variants
3. **Bi-directional lookup**:
   - By opcode byte: `GetOpcodeByBytes(0x0B) -> DecReg16 with params ["bc"]`
   - By instruction+params: `GetOpcodeByInstructionAndParams("dec", RegisterAddressing, ["bc"]) -> 0x0B`

### OpcodeDetail Structure

```go
type OpcodeDetail struct {
    OpcodeInfo          // Opcode, Size, Cycles
    Instruction *Instruction      // Reference to instruction
    Addressing  AddressingMode    // Addressing mode used  
    Params      []string          // Specific parameters (e.g., ["bc"] for DEC BC)
}
```

### Usage Examples

```go
opcodeMap := NewOpcodeMap()

// Find opcode by instruction and parameters (assembler use case)
detail := opcodeMap.GetOpcodeByInstructionAndParams("dec", RegisterAddressing, []string{"bc"})
// Returns: opcode 0x0B, params ["bc"]

// Find instruction by opcode (disassembler use case)  
detail := opcodeMap.GetOpcodeByBytes(0x0B)
// Returns: "dec" instruction with params ["bc"]

// Get all variants of an instruction
variants := opcodeMap.GetInstructionVariants("dec")
// Returns: All DEC variants (8-bit regs, 16-bit regs, indirect)
```

### Parameter Differentiation Examples

| Instruction | Addressing | Params | Opcode | Assembly |
|-------------|------------|---------|---------|----------|
| "dec" | RegisterAddressing | ["bc"] | 0x0B | DEC BC |
| "dec" | RegisterAddressing | ["de"] | 0x1B | DEC DE |
| "dec" | RegisterAddressing | ["hl"] | 0x2B | DEC HL |
| "dec" | RegisterAddressing | ["sp"] | 0x3B | DEC SP |
| "dec" | RegisterAddressing | ["a"] | 0x3D | DEC A |
| "ld" | ImmediateAddressing | ["a", "n"] | 0x3E | LD A,n |
| "ld" | RegisterAddressing | ["b", "c"] | 0x41 | LD B,C |

## Benefits for Tooling

1. **Assemblers**: Can uniquely map assembly syntax to opcodes
2. **Disassemblers**: Can generate proper assembly with specific register names
3. **Debuggers**: Can display instructions with actual parameters
4. **Code analysis**: Can distinguish between different instruction variants

## Compatibility

The `OpcodeMap` is built from the existing `Opcodes` array, so it maintains full compatibility with the current Z80 implementation while providing enhanced functionality for tooling.