# Z80 Emulator Implementation Comparison

A detailed comparison of the retrogolib Z80 implementation against three established
Go-based Z80 emulators.

## Projects Compared

| Project | Repository | LOC | Description |
|---------|-----------|-----|-------------|
| **retrogolib** | retroenv/retrogolib | 10,736 | Multi-arch emulation library (NES, Game Boy, ZX Spectrum) |
| **koron-go** | koron-go/z80 | 11,909 | Standalone Z80 emulator, 100% zexall pass rate |
| **codesqueak** | codesqueak/z80 | 2,731 | Compact emulator focused on undocumented instruction accuracy |
| **voytas** | voytas/z80-go-zx | 6,298 | ZX Spectrum-focused emulator with memory contention |

---

## 1. CPU State Representation

### retrogolib
```go
type CPU struct {
    mu sync.RWMutex
    A, B, C, D, E, H, L        uint8   // Main registers
    AltA, AltB, ..., AltL       uint8   // Shadow registers
    IX, IY                      uint16  // Index registers
    SP, PC                      uint16
    I, R                        uint8
    MEMPTR                      uint16  // Undocumented WZ register
    Flags, AltFlags             Flags   // Struct with individual flag fields
    // ...
}
```
- Individual exported fields for each register
- Separate `Flags` struct with `C`, `N`, `P`, `X`, `H`, `Y`, `Z`, `S` uint8 fields
- Thread-safe via `sync.RWMutex`
- Tracks undocumented MEMPTR/WZ register
- Q register tracked for SCF/CCF X/Y flag behavior

### koron-go
```go
type GPR struct {
    AF, BC, DE, HL Register  // Register = struct { Hi, Lo uint8 }
}
type States struct {
    GPR
    SPR              // IX, IY uint16; SP, PC uint16; IR Register
    Alternate GPR
    IFF1, IFF2 bool
    IM int
}
```
- Register pairs stored as Hi/Lo struct
- Flags packed in `AF.Lo` byte
- No mutex (single-threaded)

### codesqueak
```go
type Registers struct {
    a, b, c, d, e, h, l       byte
    a_, b_, c_, d_, e_, h_, l_ byte
    ix, iy, pc, sp             uint16
    f, f_                      byte    // Packed flag bytes
    i, r                       byte
    iff1, iff2                 bool
    ddMode, fdMode             bool    // Prefix state machine
    tStates                    uint64
}
```
- All fields unexported
- Global `var reg Registers` (no struct method receiver)
- Flags packed in single `f` byte

### voytas
```go
type registers struct {
    A, B, C, D, E, H, L, F         byte
    A_, B_, C_, D_, E_, H_, L_, F_  byte
    IXH, IXL, IYH, IYL             byte   // IX/IY split into halves
    SP, PC                          uint16
    I, R                            byte
    raw       []*byte                      // Pointer array for O(1) lookup
    prefixed  [][]*byte                    // Prefix→register redirect table
    prefix    byte
}
```
- IX/IY split into half-registers (IXH/IXL/IYH/IYL)
- Pointer-based register indirection for DD/FD prefix handling
- Flags packed in `F` byte

### Comparison

| Aspect | retrogolib | koron-go | codesqueak | voytas |
|--------|-----------|----------|------------|--------|
| Register access | Direct exported fields | Hi/Lo struct | Unexported fields | Exported + pointer array |
| Flag storage | Individual uint8 fields | Packed byte (AF.Lo) | Packed byte | Packed byte |
| IX/IY storage | uint16 | uint16 | uint16 | Split IXH/IXL bytes |
| MEMPTR/WZ | ✅ Tracked | ❌ Not tracked | ❌ Not tracked | ❌ Not tracked |
| Q register | ✅ Tracked (SCF/CCF) | ❌ | ❌ | ❌ |
| Thread safety | ✅ sync.RWMutex | ❌ | ❌ | ❌ |

**Assessment:** retrogolib is the only implementation tracking the undocumented MEMPTR/WZ
and Q registers, and the only one providing thread safety. The individual flag fields
trade performance for clarity and correctness.

MEMPTR and Q are internal Z80 registers that affect only the undocumented X/Y flag
bits (bits 3/5). MEMPTR's high byte is used as the source for X/Y flags in
`BIT n,(HL)` and indexed BIT instructions. Q tracks the previous flag state to
produce correct X/Y bits for SCF/CCF. Both are required solely to pass zexall's
undocumented flag bit verification — no real software depends on these values.

---

## 2. Flag Handling

### retrogolib — Individual Fields
```go
type Flags struct {
    C, N, P, X, H, Y, Z, S uint8  // Each 0 or 1
}

func (c *CPU) GetFlags() uint8 {
    return c.Flags.C | c.Flags.N<<1 | c.Flags.P<<2 | c.Flags.X<<3 |
        c.Flags.H<<4 | c.Flags.Y<<5 | c.Flags.Z<<6 | c.Flags.S<<7
}

func (c *CPU) setSZ(value uint8) {
    c.setS(value)
    c.setZ(value)
    c.setXY(value)
}
```
- Each flag is a separate uint8 (0 or 1)
- Batch helpers: `setSZ()`, `setSZP()`, `setXY()`
- Clear intent per operation

### koron-go — Packed Byte with Mask Operations
```go
func (cpu *CPU) updateFlagArith8(r, a, b uint16, subtract bool) {
    c := r ^ a ^ b
    var nand uint8 = maskS53 | maskZ | maskH | maskPV | maskN | maskC
    var or uint8
    or |= uint8(r) & maskS53
    if uint8(r) == 0 { or |= maskZ }
    or |= uint8(c) & maskH
    or |= uint8((c>>6)^(c>>5)) & maskPV
    if subtract { or |= maskN }
    or |= uint8(r>>8) & maskC
    cpu.AF.Lo = cpu.AF.Lo&^nand | or
}
```
- Single bulk operation: `cpu.AF.Lo = cpu.AF.Lo&^nand | or`
- All flags computed in one pass
- Most performant approach

### codesqueak — Packed Byte with Boolean Setters
```go
func setSBool(b bool) {
    if b { reg.f = reg.f | flagS } else { reg.f = reg.f & 0x7F }
}
```
- Individual set/reset per flag
- Pre-computed parity table (256 entries)

### voytas — Inline Bitwise Operations
```go
z80.Reg.F = (FS | FY | FX) & z80.Reg.A
if z80.Reg.A == 0 { z80.Reg.F |= FZ }
z80.Reg.F |= (a ^ n ^ z80.Reg.A) & FH
```
- Flags computed inline at each instruction site
- No helper functions — maximum performance but code duplication

### Comparison

| Aspect | retrogolib | koron-go | codesqueak | voytas |
|--------|-----------|----------|------------|--------|
| Storage | Individual uint8 fields | Packed byte | Packed byte | Packed byte |
| Parity calc | `bits.OnesCount8()` | `bits.OnesCount8()` | Pre-computed table | Pre-computed table |
| Update style | Per-flag setter methods | Bulk mask & OR | Per-flag bool setters | Inline bitwise |
| X/Y flags | ✅ Full support | ✅ Full support | ✅ Full support | ✅ Full support |
| Performance | Moderate (method calls) | Best (single write) | Moderate | Good (inline) |
| Readability | Best (explicit) | Low (mask math) | Good | Low (scattered) |

**Assessment:** retrogolib trades raw performance for readability with its individual flag
fields. For an emulation library that prioritizes correctness and maintainability, this
is a reasonable trade-off. koron-go's bulk mask approach is faster but harder to debug.

**Potential improvement:** Consider a hybrid approach — keep individual fields for clarity
but add a `packFlags()` method that uses the koron-go mask pattern for hot paths, or
benchmark to verify if the current approach is actually a bottleneck.

---

## 3. Instruction Dispatch

### retrogolib — Opcode Table + Function Pointers
```go
var Opcodes = [256]Opcode{
    {Instruction: Nop, Addressing: ImpliedAddressing, Timing: 4, Size: 1},
    {Instruction: LdReg16, Addressing: ImmediateAddressing, Timing: 10, Size: 3},
    // ...
}

// In Step():
ins := opcode.Instruction
if ins.NoParamFunc != nil {
    err := ins.NoParamFunc(c)
} else {
    params, operands, _ := readOpParams(c, opcode.Addressing)
    err := ins.ParamFunc(c, params...)
}
```
- Static [256]Opcode arrays (one per prefix: base, CB, DD, ED, FD)
- Each opcode carries Instruction pointer, addressing mode, timing, size
- Two-path dispatch: NoParamFunc vs ParamFunc with dynamic params
- Bidirectional lookup (opcode↔instruction) for disassembly

### koron-go — Giant Switch (936 Cases)
```go
func (cpu *CPU) executeOne() {
    switch c0 := cpu.fetchM1(); c0 {
    case 0x00: oopNOP(cpu)
    case 0x01: xopLDbcnn(cpu)
    // ... 936 cases, nested switches for prefixes
    }
}
```
- Manually written switch with 936+ cases
- Nested switches for CB/DD/ED/FD prefixes
- Each instruction is a dedicated function

### codesqueak — Algorithmic Decode (x/y/z Extraction)
```go
x, y, z := basicDecode(inst)  // x=(inst>>6)&3, y=(inst>>3)&7, z=inst&7
switch x {
case 0: decodeX0(y, z)
case 1: store8r(load8r(z), y)  // LD r,r' inline
case 2: decodeX2(y, z)         // ALU ops
default: decodeX3(y, z)
}
```
- Bit-field extraction from opcode
- 4 main blocks, further decoded by y/z
- Most compact dispatch

### voytas — Switch with Named Constants
```go
switch opcode {
case nop:          // 0x00
case halt:         // 0x76
case ld_bc_nn:     // 0x01
case prefix_cb:    prefixCB()
case useIX:        z80.Reg.prefix = useIX; continue
// ... 100+ cases
}
```
- Named opcode constants for readability
- DD/FD handled by setting prefix and re-looping

### Comparison

| Aspect | retrogolib | koron-go | codesqueak | voytas |
|--------|-----------|----------|------------|--------|
| Method | Table + func ptrs | Giant switch | Algorithmic decode | Switch + constants |
| Cases | 256 per table | 936 (nested) | ~50 (decoded) | ~100 |
| Prefix handling | Separate [256] tables | Nested switch | State machine | Prefix + re-loop |
| Disassembly support | ✅ Bidirectional | ❌ Execute only | ❌ Execute only | ❌ Execute only |
| Static analysis | ✅ Categories, addressing | ❌ | ❌ | ❌ |
| Extensibility | Best (data-driven) | Low (manual) | Good (algorithmic) | Moderate |

**Assessment:** retrogolib's table-driven approach is unique — it supports both emulation
AND static analysis/disassembly through the same data structures. This is a major
architectural advantage over the other implementations which only support execution.
The dynamic parameter passing (`params ...any`) has a performance cost but enables
clean separation of param reading from instruction execution.

**Potential improvement:** The `params ...any` interface causes heap allocations on every
parameterized instruction. Consider typed parameter structs or pre-allocated param
slots for hot-path optimization.

---

## 4. Memory Interface

| Aspect | retrogolib | koron-go | codesqueak | voytas |
|--------|-----------|----------|------------|--------|
| Interface | `Read`/`Write`/`ReadWord`/`WriteWord` | `Get`/`Set` | `Get`/`Put`/`Load` | `Read`/`Write` |
| Word ops | Interface methods | Internal helpers | Manual | Manual |
| Bundled impl | `BasicMemory` (64KB flat) + `GameBoyMemory` | `DumbMemory` + `MapMemory` + `im0data` | Simple test RAM | `BasicMemory` |
| Bank switching | Via interface impl | Via interface impl | Via interface impl | Via interface impl |

**Assessment:** retrogolib provides the richest default implementation with both a flat
memory model and a Game Boy-specific memory mapper. The `ReadWord`/`WriteWord` interface
methods reduce error-prone manual byte assembly. All projects correctly use interfaces
for hardware abstraction.

---

## 5. I/O Port Handling

| Aspect | retrogolib | koron-go | codesqueak | voytas |
|--------|-----------|----------|------------|--------|
| Interface | `IOHandler` (optional) | `IO` (required field) | `IO` interface | `IOBus` interface |
| Port address | uint8 | uint8 | byte | (hi, lo) byte pair |
| Null safety | Returns 0xFF if nil | Returns 0 if nil | Assumed non-nil | Returns 0xFF if nil |
| 16-bit port addr | ❌ | ❌ | ❌ | ✅ (hi=B, lo=port) |

**Assessment:** voytas is the most hardware-accurate with its (hi, lo) port address pair,
matching real Z80 behavior where register B is placed on the high address bus during
IN/OUT (C) instructions. retrogolib and others simplify to 8-bit port addresses.

**Potential improvement:** Consider expanding the `IOHandler` interface to accept a
16-bit address for hardware-accurate port decoding, or add a `ReadPort16`/`WritePort16`
variant.

---

## 6. Prefix Handling (CB/DD/ED/FD)

### retrogolib
- Separate `[256]Opcode` tables: `CBOpcodes`, `DDOpcodes`, `EDOpcodes`, `FDOpcodes`
- DD/FD fallthrough: advances PC, executes unprefixed with +4 T-states
- DDCB/FDCB: Dynamic instruction selection based on opcode range

### koron-go
- Nested switch inside main dispatch (inline prefix handling)
- DDCB/FDCB: Three-level nested switch (DD → CB → operation)
- Offset fetched inline: `d := cpu.fetch()`

### codesqueak
- State machine: `ddMode`/`fdMode` booleans set on prefix
- Prefix-specific T-state lookup tables (5 total)
- Register redirection handled by mode flags

### voytas
- Prefix sets `z80.Reg.prefix` and `continue`s the main loop
- Register pointer array automatically redirects HL→IX or HL→IY
- Most elegant DD/FD handling through indirection

### Comparison

| Aspect | retrogolib | koron-go | codesqueak | voytas |
|--------|-----------|----------|------------|--------|
| DD/FD tables | Separate [256] arrays | Inline nested switch | State machine | Prefix + pointer redirect |
| DD/FD fallthrough | ✅ +4 T-states | ❌ | ❌ | ❌ |
| DDCB/FDCB | Dynamic instruction | 3-level switch | Separate handler | Separate handler |
| Undocumented DD/FD | ✅ Passthrough | ✅ IXH/IXL opcodes | ✅ IXH/IXL | ✅ IXH/IXL |

**Assessment:** retrogolib correctly handles the undocumented DD/FD passthrough behavior
(executing unprefixed instruction with +4 T-states), which most other implementations
skip. voytas has the most elegant DD/FD handling through pointer indirection.

---

## 7. Testing Strategy

| Aspect | retrogolib | koron-go | codesqueak | voytas |
|--------|-----------|----------|------------|--------|
| Test LOC | ~1,700 | ~2,500 | ~230 | ~3,000 |
| Test functions | 52 | ~80 | 1 | 112 |
| zexdoc | ✅ (67/67) | ✅ (67/67) | ❌ | ✅ |
| zexall | ✅ (67/67) | ✅ (67/67) | ❌ | ✅ |
| Single-step tests | ✅ (JSON-based) | ❌ | ❌ | ❌ |
| Unit tests | ✅ Per-operation | ✅ Per-operation | ❌ (integration only) | ✅ Per-instruction |
| CRC validation | ❌ | ✅ | ❌ | ❌ |

**Assessment:** Both retrogolib and koron-go achieve 100% zexall pass rate (67/67),
validating full undocumented flag behavior. retrogolib's JSON-based single-step tests
are unique and excellent for regression testing. voytas has the most granular
per-instruction unit tests (112 functions). koron-go's CRC-based validation approach
is a complementary technique worth studying.

---

## 8. Undocumented Instruction Support

| Feature | retrogolib | koron-go | codesqueak | voytas |
|---------|-----------|----------|------------|--------|
| IXH/IXL/IYH/IYL arithmetic | ✅ (DD/FD undoc files) | ✅ | ✅ | ✅ |
| SLL (undocumented shift) | ✅ | ✅ (SL1) | ✅ | ✅ |
| DD/FD passthrough | ✅ | ❓ | ❌ | ❌ |
| Undocumented X/Y flags | ✅ | ✅ | ✅ | ✅ |
| MEMPTR/WZ register | ✅ | ❌ | ❌ | ❌ |
| Q register (SCF/CCF) | ✅ | ❌ | ❌ | ❌ |
| Undocumented ED NEG mirrors | ❌ | ✅ | ✅ | ✅ |
| Undocumented ED IM mirrors | ❌ | ✅ | ✅ | ✅ |

**Assessment:** retrogolib leads in MEMPTR/WZ and Q register tracking — these are the
hardest undocumented behaviors to get right and are validated by the 100% zexall pass
rate. It may be missing some undocumented ED prefix mirrors (multiple opcodes mapping
to NEG, IM, etc.) that other implementations include but which are not tested by zexall.

**Potential improvement:** Add undocumented ED opcode mirrors (e.g., ED 4C/54/5C/64/6C/7C
all executing NEG; ED 4E/66/6E mapping to IM 0/0/0) for completeness.

---

## 9. Thread Safety

| Feature | retrogolib | koron-go | codesqueak | voytas |
|---------|-----------|----------|------------|--------|
| Concurrency model | sync.RWMutex | None | None | None |
| Safe State() access | ✅ | ❌ | ❌ | ❌ |
| Safe register reads | ✅ (exported methods) | ❌ | ❌ | ❌ |
| Multiple instances | ✅ | ✅ | ❌ (global state) | ✅ |

**Assessment:** retrogolib is the only implementation designed for concurrent access.
This is important for emulators with separate render/audio/CPU threads. codesqueak's
use of global state is a significant limitation — it cannot run multiple Z80 instances.

---

## 10. Architecture & Extensibility

| Feature | retrogolib | koron-go | codesqueak | voytas |
|---------|-----------|----------|------------|--------|
| Disassembly support | ✅ (RegisterOpcodes, RegisterPairOpcodes) | ❌ | ❌ | ❌ |
| Static analysis | ✅ (categories, addressing modes) | ❌ | ❌ | ❌ |
| Instruction categories | ✅ (branching, memory R/W sets) | ❌ | ❌ | ❌ |
| System variants | ✅ (Game Boy, ZX Spectrum, Generic) | ❌ | ❌ | ✅ (ZX Spectrum) |
| Option pattern | ✅ (functional options) | ❌ | ❌ | ❌ |
| Pre-execution hooks | ✅ | ❌ | ❌ | ✅ (Trap function) |
| Tracing | ✅ (TraceStep) | ❌ | ❌ | ❌ |

**Assessment:** retrogolib is architecturally the most versatile — it's designed as a
library that supports emulation, disassembly, and static analysis. The other
implementations are pure emulators. The instruction category sets
(`BranchingInstructions`, `MemoryReadInstructions`, etc.) enable higher-level tools
to reason about code without executing it.

---

## 11. Performance Characteristics

| Aspect | retrogolib | koron-go | codesqueak | voytas |
|--------|-----------|----------|------------|--------|
| Dispatch cost | Table lookup + func ptr | Switch (compiler-optimized) | Bit decode + switch | Switch |
| Flag update cost | Multiple field writes | Single byte write | Single byte write | Inline byte write |
| Param passing | `...any` (heap alloc) | Direct register access | Direct register access | Direct register access |
| Mutex overhead | RWMutex per Step() | None | None | None |
| Memory interface | Virtual dispatch | Virtual dispatch | Pointer deref | Virtual dispatch |
| Parity calculation | Runtime `bits.OnesCount8` | Runtime `bits.OnesCount8` | Pre-computed table | Pre-computed table |

**Assessment:** retrogolib has higher per-instruction overhead due to: (1) mutex
lock/unlock on every Step(), (2) `...any` parameter passing causing allocations,
(3) individual flag field writes. For most use cases this is negligible, but for
cycle-accurate emulation at full speed, it may become relevant.

**Potential improvements:**
1. Move mutex to a higher-level `Run(cycles)` method instead of per-Step
2. Consider pre-computed parity table (256 bytes) instead of runtime calculation
3. Profile `...any` param allocation — may be optimizable with sync.Pool or typed params

---

## 12. Interrupt Handling

| Aspect | retrogolib | koron-go | codesqueak | voytas |
|--------|-----------|----------|------------|--------|
| IM 0 | ✅ (simplified: RST 38H) | ✅ (full: executes data bus instruction) | ✅ (state only) | ✅ (simplified: RST 38H) |
| IM 1 | ✅ | ✅ | ✅ (state only) | ✅ |
| IM 2 | ✅ (simplified) | ✅ (full) | ✅ (state only) | ✅ |
| NMI | ✅ | ✅ | ❌ (no trigger mechanism) | ✅ |
| RETI handler | ❌ | ✅ (callback) | ❌ | ❌ |
| RETN handler | ❌ | ✅ (callback) | ❌ | ❌ |

**Assessment:** koron-go has the most complete interrupt handling, including IM 0 data
bus emulation via pseudo-memory injection and RETI/RETN callbacks for peripheral
notification. retrogolib's implementation is functional but simplified.

**Potential improvement:** Add `RETIHandler`/`RETNHandler` callbacks to notify peripherals
when interrupt service routines complete. Consider full IM 0 data bus instruction
execution for hardware accuracy.

---

## Summary: Strengths and Gaps

### retrogolib Unique Strengths
1. **Dual-purpose architecture** — supports both emulation and static analysis/disassembly
2. **MEMPTR/WZ tracking** — only implementation with this undocumented register
3. **Q register tracking** — correct SCF/CCF undocumented flag behavior
4. **Thread safety** — safe for multi-threaded emulator architectures
5. **System variants** — configurable for Game Boy, ZX Spectrum, or generic
6. **Tracing infrastructure** — built-in instruction tracing with pre-execution hooks
7. **DD/FD passthrough** — correct undocumented behavior with +4 T-states

### Areas for Improvement
1. **Undocumented ED mirrors** — add NEG/IM/RETN opcode mirrors
2. **IM 0 accuracy** — execute data bus instruction instead of assuming RST 38H
3. **RETI/RETN callbacks** — notify peripherals on interrupt return
4. **16-bit I/O port address** — expose B register on high address bus
5. **Performance** — profile and optimize mutex, `...any` params, flag writes
6. **Pre-computed parity table** — replace runtime `bits.OnesCount8()` with lookup

### Overall Assessment

retrogolib's Z80 implementation is architecturally the most sophisticated of the four
projects compared. Its dual-purpose design (emulation + static analysis), thread safety,
and undocumented register tracking set it apart. It matches koron-go's 100% zexall pass
rate, validating that the MEMPTR and Q register tracking produces correct undocumented
flag behavior.

For pure emulation speed, koron-go's packed flags and switch dispatch would be faster.
For a reusable library that tools can build on, retrogolib's table-driven architecture
with bidirectional opcode lookup is the superior design.
