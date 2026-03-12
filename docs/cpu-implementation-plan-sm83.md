# CPU Implementation Plan: SM83 (Sharp LR35902)

## Context

The SM83 (Sharp LR35902) is a custom 8-bit CPU used in the Nintendo Game Boy (1989) and
Game Boy Color (1998). Often described as "a modified Z80," the SM83 shares the Z80's basic
register set and many ALU instructions but removes entire subsystems (shadow registers, index
registers, I/O ports, three prefix groups) and adds unique instructions for memory-mapped I/O.
The differences are substantial enough to require a dedicated package rather than a Z80 variant.

**Package:** `arch/cpu/sm83/`
**Status:** COMPLETE (Phases 1-3)
**Last Updated:** 2026-03-12

---

## 1. SM83 Architecture Overview

### 1.1 Registers

| Register | Width  | Description                                              |
|----------|--------|----------------------------------------------------------|
| A        | 8-bit  | Accumulator                                              |
| F        | 8-bit  | Flags (upper nibble only: Z, N, H, C; bits 3-0 always 0)|
| B        | 8-bit  | General purpose                                          |
| C        | 8-bit  | General purpose (also used for $FF00+C addressing)       |
| D        | 8-bit  | General purpose                                          |
| E        | 8-bit  | General purpose                                          |
| H        | 8-bit  | General purpose (HL pair used as memory pointer)         |
| L        | 8-bit  | General purpose                                          |
| SP       | 16-bit | Stack pointer                                            |
| PC       | 16-bit | Program counter                                          |
| IME      | 1-bit  | Interrupt Master Enable (not part of any register)       |

**Register pairs** for 16-bit operations: AF, BC, DE, HL, SP.

No shadow registers (AF', BC', DE', HL'). No index registers (IX, IY). No interrupt
vector register (I) or refresh register (R).

### 1.2 Flags

Four flags stored in the upper nibble of the F register:

| Bit | Flag | Name       | Description                              |
|-----|------|------------|------------------------------------------|
| 7   | Z    | Zero       | Set when result is zero                  |
| 6   | N    | Subtract   | Set when last operation was subtraction  |
| 5   | H    | Half-carry | Set on carry from bit 3 to bit 4         |
| 4   | C    | Carry      | Set on carry from bit 7                  |
| 3-0 | -    | Unused     | Always 0                                 |

### 1.3 Addressing Modes

| Mode              | Syntax        | Example          | Description                          |
|-------------------|---------------|------------------|--------------------------------------|
| Implied           | -             | NOP, HALT        | No operand                           |
| Register          | r             | INC B            | 8-bit register operand               |
| Register pair     | rr            | INC BC           | 16-bit register pair operand         |
| Immediate 8-bit   | n             | LD A,$42         | 8-bit value follows opcode           |
| Immediate 16-bit  | nn            | LD BC,$1234      | 16-bit value follows opcode (LE)     |
| Register indirect | (rr)          | LD A,(HL)        | Memory at address in register pair   |
| Direct            | (nn)          | LD A,($C000)     | Memory at 16-bit immediate address   |

Additionally, the SM83 uses specialized high-RAM addressing for I/O: `LDH (n),A` accesses
`$FF00+n` and `LD (C),A` accesses `$FF00+C`.

### 1.4 Condition Codes

| Code | Flag Test  | Description           |
|------|------------|-----------------------|
| NZ   | Z = 0      | Not zero              |
| Z    | Z = 1      | Zero                  |
| NC   | C = 0      | No carry              |
| C    | C = 1      | Carry                 |

The Z80's additional conditions (PO, PE, P, M) based on Sign and Parity/Overflow flags
are not available.

---

## 2. Key Differences from Z80

| Feature              | Z80                                  | SM83                                          |
|----------------------|--------------------------------------|-----------------------------------------------|
| Shadow registers     | AF', BC', DE', HL'                   | **None**                                      |
| Index registers      | IX, IY (16-bit)                      | **None**                                      |
| I/R registers        | Interrupt vector (I), Refresh (R)    | **None**                                      |
| I/O instructions     | IN/OUT (256 ports)                   | **None** (memory-mapped I/O only)             |
| Prefix groups        | CB, DD, ED, FD                       | **CB only**                                   |
| Flags                | S, Z, H, P/V, N, C + X, Y (8 bits)  | Z, N, H, C only (4 bits)                     |
| Interrupt modes      | IM 0, IM 1, IM 2                    | **Single mode** (5 fixed vectors)             |
| HALT behavior        | Waits for interrupt                  | **HALT bug** (PC not incremented when IME=0)  |
| Condition codes      | NZ, Z, NC, C, PO, PE, P, M          | NZ, Z, NC, C only                            |
| Clock speed          | Variable                             | 4.19 MHz (fixed, 8.39 MHz CGB double speed)  |

### Unique SM83 Instructions

| Opcode | Mnemonic      | Description                                          |
|--------|---------------|------------------------------------------------------|
| $10    | STOP          | Enter low-power mode (Z80: DJNZ)                    |
| $08    | LD (nn),SP    | Store SP at 16-bit address (Z80: EX AF,AF')          |
| $E0    | LDH (n),A     | Store A at $FF00+n (Z80: RET PO)                    |
| $F0    | LDH A,(n)     | Load A from $FF00+n (Z80: RET P)                    |
| $E2    | LD (C),A      | Store A at $FF00+C (Z80: JP PO,nn)                  |
| $F2    | LD A,(C)      | Load A from $FF00+C (Z80: JP P,nn)                  |
| $E8    | ADD SP,e      | Add signed 8-bit to SP (Z80: RET PE)                |
| $F8    | LD HL,SP+e    | Load SP + signed 8-bit into HL (Z80: RET M)         |
| $22    | LD (HL+),A    | Store A at (HL), increment HL (Z80: LD (nn),HL)     |
| $2A    | LD A,(HL+)    | Load A from (HL), increment HL (Z80: LD HL,(nn))    |
| $32    | LD (HL-),A    | Store A at (HL), decrement HL (Z80: LD (nn),A)      |
| $3A    | LD A,(HL-)    | Load A from (HL), decrement HL (Z80: LD A,(nn))     |
| $EA    | LD (nn),A     | Store A at 16-bit address (Z80: JP PE,nn)            |
| $FA    | LD A,(nn)     | Load A from 16-bit address (Z80: JP P,nn)            |
| $D9    | RETI          | Return and enable interrupts (Z80: EXX)              |
| $CB 3x | SWAP r        | Swap upper and lower nibbles (Z80: SLL r)            |

---

## 3. Complete Instruction Set

The SM83 has 245 valid base opcodes (11 unused) and 256 CB-prefixed opcodes, totaling
501 valid instructions.

### Base Opcodes ($00-$FF)

| Range     | Group                  | Instructions                                   |
|-----------|------------------------|------------------------------------------------|
| $00       | Control                | NOP                                            |
| $01-$31   | 16-bit loads           | LD rr,nn (BC/DE/HL/SP)                         |
| $02-$3A   | Indirect loads         | LD (BC/DE/HL+/HL-),A and reverse               |
| $03-$3B   | 16-bit arithmetic      | INC rr, DEC rr                                 |
| $04-$3D   | 8-bit INC/DEC          | INC r, DEC r (all 8 registers)                 |
| $06-$3E   | Immediate loads        | LD r,n (all 8 registers)                       |
| $07-$0F   | Rotate accumulator     | RLCA, RRCA, RLA, RRA                           |
| $08       | Special load           | LD (nn),SP                                     |
| $09-$39   | 16-bit ADD             | ADD HL,rr                                      |
| $10       | Power control          | STOP                                           |
| $18-$38   | Relative jumps         | JR e, JR cc,e                                  |
| $20-$30   | Conditional JR         | JR NZ/Z/NC/C,e                                 |
| $27       | BCD adjust             | DAA                                            |
| $2F       | Complement             | CPL                                            |
| $37       | Set carry              | SCF                                            |
| $3F       | Complement carry       | CCF                                            |
| $40-$7F   | Register loads         | LD r,r' (64 entries, $76 = HALT)               |
| $76       | Control                | HALT                                           |
| $80-$BF   | ALU register ops       | ADD/ADC/SUB/SBC/AND/XOR/OR/CP A,r              |
| $C0-$D8   | Conditional returns    | RET cc                                         |
| $C1-$F1   | Stack ops              | POP rr, PUSH rr (AF/BC/DE/HL)                  |
| $C2-$DA   | Conditional jumps      | JP cc,nn                                       |
| $C3       | Unconditional jump     | JP nn                                          |
| $C4-$DC   | Conditional calls      | CALL cc,nn                                     |
| $C6-$FE   | ALU immediate ops      | ADD/ADC/SUB/SBC/AND/XOR/OR/CP A,n              |
| $C7-$FF   | Restart                | RST $00/$08/$10/$18/$20/$28/$30/$38             |
| $C9       | Return                 | RET                                            |
| $CB       | Prefix                 | CB prefix (see section 4)                      |
| $CD       | Call                   | CALL nn                                        |
| $D9       | Return interrupt       | RETI                                           |
| $E0/$F0   | High-RAM I/O           | LDH (n),A / LDH A,(n)                         |
| $E2/$F2   | Register-indirect I/O  | LD (C),A / LD A,(C)                            |
| $E8       | SP arithmetic          | ADD SP,e                                       |
| $E9       | Jump indirect          | JP HL                                          |
| $EA/$FA   | Direct loads           | LD (nn),A / LD A,(nn)                          |
| $F3/$FB   | Interrupt control      | DI / EI                                        |
| $F8       | SP-relative load       | LD HL,SP+e                                     |
| $F9       | SP load                | LD SP,HL                                       |

---

## 4. CB-Prefix Instructions

All 256 CB-prefixed opcodes ($CB $00-$CB $FF) are valid. They follow a regular encoding:

| Range       | Instruction | Description                                      |
|-------------|-------------|--------------------------------------------------|
| $CB $00-$07 | RLC r       | Rotate left circular                             |
| $CB $08-$0F | RRC r       | Rotate right circular                            |
| $CB $10-$17 | RL r        | Rotate left through carry                        |
| $CB $18-$1F | RR r        | Rotate right through carry                       |
| $CB $20-$27 | SLA r       | Shift left arithmetic                            |
| $CB $28-$2F | SRA r       | Shift right arithmetic (preserves sign bit)      |
| $CB $30-$37 | SWAP r      | **Swap upper and lower nibbles** (replaces Z80 SLL) |
| $CB $38-$3F | SRL r       | Shift right logical                              |
| $CB $40-$7F | BIT b,r     | Test bit b (8 bits x 8 registers = 64 opcodes)   |
| $CB $80-$BF | RES b,r     | Reset (clear) bit b                              |
| $CB $C0-$FF | SET b,r     | Set bit b                                        |

**Register encoding** (bits 2-0 of second byte): B=0, C=1, D=2, E=3, H=4, L=5, (HL)=6, A=7.

Operations on (HL) access memory at the address in the HL register pair and take additional
cycles compared to register-only operations.

---

## 5. Illegal Opcodes

11 base opcodes are undefined and should trigger an illegal opcode error:

| Opcode | Z80 Equivalent     |
|--------|--------------------|
| $D3   | OUT (n),A           |
| $DB   | IN A,(n)            |
| $DD   | IX prefix           |
| $E3   | EX (SP),HL          |
| $E4   | CALL PO,nn          |
| $EB   | EX DE,HL            |
| $EC   | CALL PE,nn          |
| $ED   | ED prefix           |
| $F4   | CALL P,nn           |
| $FC   | CALL M,nn           |
| $FD   | IY prefix           |

These opcodes correspond to Z80 instructions that relied on removed features (I/O ports,
index registers, parity/sign conditions, register exchange). On real hardware, executing
these opcodes causes undefined behavior; the implementation treats them as errors.

---

## 6. Interrupt System

### 6.1 Interrupt Vectors

The SM83 uses a single interrupt mode with 5 fixed vectors:

| Priority | Vector | Bit | Source    |
|----------|--------|-----|----------|
| Highest  | $0040  | 0   | V-Blank  |
|          | $0048  | 1   | LCD STAT |
|          | $0050  | 2   | Timer    |
|          | $0058  | 3   | Serial   |
| Lowest   | $0060  | 4   | Joypad   |

### 6.2 Control Registers

| Register | Address | Description                                |
|----------|---------|--------------------------------------------|
| IE       | $FFFF   | Interrupt Enable -- which interrupts are enabled |
| IF       | $FF0F   | Interrupt Flag -- which interrupts are pending   |
| IME      | (internal) | Interrupt Master Enable -- global enable    |

- **EI** ($FB) sets IME=1 (takes effect after the next instruction).
- **DI** ($F3) sets IME=0 immediately.
- **RETI** ($D9) returns from interrupt and sets IME=1.

### 6.3 Interrupt Dispatch

When `IME=1` and `(IE & IF) != 0`:
1. IME is cleared (IME=0)
2. The highest-priority pending bit in `(IE & IF)` is identified
3. The corresponding IF bit is cleared
4. PC is pushed onto the stack
5. PC is set to the corresponding vector address
6. Total dispatch cost: 20 cycles (5 machine cycles)

### 6.4 HALT Bug

When HALT is executed with `IME=0` and an interrupt is pending (`IE & IF != 0`):
- The CPU resumes execution but **fails to increment PC** for the next instruction
- The byte immediately after HALT is read and executed twice
- This is a well-documented silicon bug present on all SM83 revisions

---

## 7. File Structure

```
arch/cpu/sm83/
    doc.go                -- Package documentation
    addressing.go         -- 7 addressing modes
    instruction.go        -- Instruction type and definitions
    opcode.go             -- 256-entry base opcode table (245 valid + 11 illegal)
    opcode_cb.go          -- 256-entry CB-prefix opcode table (all valid)
    categories.go         -- Instruction category sets for static analysis
    errors.go             -- Package-specific errors (illegal opcode)
    flag.go               -- Flag definitions (Z, N, H, C)
    cpu.go                -- CPU state struct, registers, initialization
    option.go             -- Functional options (WithPC, WithSP, etc.)
    memory.go             -- Memory interface (Read/Write/ReadWord/WriteWord)
    step.go               -- Fetch/decode/execute cycle
    param.go              -- Operand reading and addressing mode resolution
    interrupt.go          -- Interrupt dispatch (5 vectors, IME, IE/IF, HALT bug)
    emulation.go          -- ALU handlers (ADD, ADC, SUB, SBC, AND, OR, XOR, CP,
                             INC, DEC, DAA, CPL, CCF, SCF, RLCA, RRCA, RLA, RRA)
    emulation_load.go     -- Load/store handlers (LD variants, LDH, LD (HL+/-),
                             PUSH, POP, LD (nn),SP, ADD SP,e, LD HL,SP+e)
    emulation_branch.go   -- Branch handlers (JP, JR, CALL, RET, RETI, RST,
                             conditional variants with NZ/Z/NC/C)
    emulation_cb.go       -- CB-prefix handlers (RLC, RRC, RL, RR, SLA, SRA,
                             SWAP, SRL, BIT, RES, SET)
    singlestep_test.go    -- SingleStepTests integration test
```

---

## 8. Implementation Status

| Phase | Description                  | Status   |
|-------|------------------------------|----------|
| 1     | Static Analysis Foundation   | COMPLETE |
| 2     | CPU Emulation Core           | COMPLETE |
| 3     | Instruction Handlers         | COMPLETE |

### Phase 1: Static Analysis Foundation

- All addressing modes defined in `addressing.go`
- Complete instruction definitions in `instruction.go`
- Full 256-entry base opcode table in `opcode.go` (245 valid, 11 illegal)
- Full 256-entry CB-prefix opcode table in `opcode_cb.go` (all valid, including SWAP)
- Instruction categories for static analysis in `categories.go`
- Flag constants (Z, N, H, C) in `flag.go`
- Illegal opcode error handling in `errors.go`

### Phase 2: CPU Emulation Core

- CPU state struct with all registers (A, F, B, C, D, E, H, L, SP, PC, IME)
- Memory interface for bus integration
- Fetch/decode/execute step loop in `step.go`
- Operand reading for all addressing modes in `param.go`
- Functional options for CPU configuration
- Interrupt dispatch with 5 fixed vectors, IME management, and HALT bug emulation

### Phase 3: Instruction Handlers

- ALU operations with correct 4-flag behavior (no S, P/V flags)
- All load/store variants including SM83-unique LDH, LD (HL+/-), ADD SP,e, LD HL,SP+e
- Branch/call/return with 4 condition codes (NZ, Z, NC, C)
- All 256 CB-prefix operations including SWAP (replacing Z80 SLL at $30-$37)

---

## 9. Verification

### Build and Lint

All standard checks pass:

- `go build ./arch/cpu/sm83/...`
- `go vet ./arch/cpu/sm83/...`
- `golangci-lint run ./arch/cpu/sm83/...`

### Testing

- **SingleStepTests** (`singlestep_test.go`): Integration test using the
  [SingleStepTests](https://github.com/SingleStepTests) JSON test suite for the SM83.
  Each test provides initial CPU/memory state, executes one instruction, and verifies
  the final CPU state, memory contents, and cycle count against expected values. This
  validates all 501 valid opcodes (245 base + 256 CB-prefix) including flag behavior,
  memory access patterns, and cycle accuracy.

---

## References

- **Pan Docs** (gbdev.io/pandocs): Comprehensive Game Boy technical reference
- **SM83 Instruction Set** (gbdev.io/gb-opcodes): Complete opcode table with flags and cycles
- **SingleStepTests**: Per-instruction JSON test vectors for SM83 validation
- **Game Boy CPU Manual**: Community instruction set reference with flag effects
- **TCAGBD** (The Cycle-Accurate Game Boy Docs): Detailed timing documentation
- **System plan**: `docs/system-implementation-plan-gameboy.md` covers the full Game Boy
  system including memory map, I/O registers, and cartridge format
