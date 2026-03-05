# Z80 Emulator Cross-Language Comparison

A detailed comparison of the retrogolib Z80 implementation against five established
non-Go Z80 emulators across C, C++, and Rust.

## Projects Compared

| Project | Language | LOC | License | Description |
|---------|----------|-----|---------|-------------|
| **retrogolib** | Go | 10,737 | MIT | Multi-arch emulation library with dual-purpose design (emulation + disassembly) |
| **floooh/chips** | C (header-only) | ~4,191 | zlib | Cycle-accurate, pin-centric design; powers ZX Spectrum/CPC/KC85 emulators |
| **redcode/Z80** | ANSI C | ~3,894 | LGPL-3 | Most thoroughly documented undocumented behavior; validated against real hardware |
| **superzazu/z80** | C99 | ~1,814 | MIT | Clean, educational-quality implementation; zexall-passing |
| **kosarev/z80** | C++11 (header-only) | ~4,051 | MIT | Template-based zero-overhead CRTP architecture; dual Z80/i8080 |
| **royaltm/rust-z80emu** | Rust | ~10,247 | LGPL-3 | Three CPU flavour variants (NMOS/CMOS/BM1); powers SPECTRUSTY emulator |

---

## 1. Instruction Dispatch

How each implementation routes an opcode byte to its handler.

| Project | Approach | Details |
|---------|----------|---------|
| **retrogolib** | Function pointer table | 5 tables of 256 entries (`Opcodes`, `CBOpcodes`, `EDOpcodes`, `DDOpcodes`, `FDOpcodes`); each entry holds handler function + metadata for dual emulation/disassembly use |
| **floooh/chips** | Code-generated switch | Python generator (`z80_gen.py`) reads YAML instruction descriptions and emits a ~1,700-case `switch(cpu->step)` where each case is one T-state; hybrid with algorithmic CB decoding |
| **redcode/Z80** | Function pointer table | 5 tables of 256 `Insn` function pointers; each returns T-state count. Explicitly chosen over switch for IM 0 instruction reuse |
| **superzazu/z80** | Hand-written switch | `exec_opcode()` with 256 cases; separate switch functions for CB/ED/DDFD prefixes. Cycle counts via pre-computed lookup arrays |
| **kosarev/z80** | Template + structured decode | CRTP-based dispatch; decodes opcode x/y/z/p/q bit fields (per z80.info decoding scheme) rather than flat 256-case switch |
| **rust-z80emu** | Macro-generated match | Three-layer macro pipeline (`match_instruction!` -> `instruction_dispatch!` -> `run_mnemonic!`) expands at compile time into Rust `match` expressions |

### Analysis

retrogolib's function-pointer tables are closest to redcode/Z80's approach. The key difference is
retrogolib stores instruction metadata (name, operand types, size) alongside each handler, enabling
static analysis and disassembly from the same table -- a feature no other implementation provides.

The floooh/chips approach of code generation from YAML descriptions is the most sophisticated for
maintainability, while kosarev's CRTP and rust-z80emu's macro pipelines are the most architecturally
innovative for their respective languages.

---

## 2. Flag Representation

How the Z80's 8-bit F register is stored and manipulated.

| Project | Storage | Pack/Unpack | Lookup Tables |
|---------|---------|-------------|---------------|
| **retrogolib** | Individual `bool` fields in `Flags` struct | `GetFlags()`/`setFlags()` methods | None; inline computation |
| **floooh/chips** | Packed `uint8_t` byte | Direct bitwise ops | `_z80_szp_flags[256]` for sign+zero+parity |
| **redcode/Z80** | Packed `zuint8` byte | Direct bitwise ops with bitmask macros | `pf_parity_table[256]`; optional `daa_af_table[2048]` |
| **superzazu/z80** | Individual C bitfields (`bool sf:1, zf:1, ...`) | `get_f()`/`set_f()` functions | None; inline computation |
| **kosarev/z80** | Packed `fast_u8` byte | Bitmask constants + optional lazy evaluation | Lazy flag computation defers work until flags are read |
| **rust-z80emu** | `bitflags!` macro struct (`CpuFlags: u8`) | Idiomatic Rust bitflag operations | Precomputed parity via method |

### Analysis

retrogolib and superzazu/z80 are the only implementations using individual flag fields rather than
a packed byte. This trades a small performance cost on pack/unpack for clearer, more readable flag
manipulation code. Both require conversion functions when the full F byte is needed (e.g., PUSH AF,
POP AF).

kosarev/z80's optional lazy flag evaluation is unique -- flags are stored as an operand+result pair
and only computed into actual flag bits when read. This can avoid unnecessary flag computation for
instructions whose flag results are never consumed, at the cost of complexity.

---

## 3. MEMPTR/WZ Register

The undocumented internal register whose high byte provides X/Y flag bits for BIT n,(HL) and
indexed BIT instructions.

| Project | Tracked | Storage | Usage |
|---------|---------|---------|-------|
| **retrogolib** | Yes | `MEMPTR uint16` | Updated in memory/IO/jump/block instructions; high byte used for BIT (HL) X/Y flags |
| **floooh/chips** | Yes | `uint16_t wz` (union with `wzh`/`wzl`) | Full tracking; `(cpu->wz >> 8) & (Z80_YF\|Z80_XF)` for BIT |
| **redcode/Z80** | Yes | `ZInt16 memptr` (union with `MEMPTRH`/`MEMPTRL`) | Most extensively documented; project wiki is a reference resource |
| **superzazu/z80** | Yes | `uint16_t mem_ptr` | Updated in jumps, calls, memory access instructions |
| **kosarev/z80** | Yes | `reg16_value wz` | Accessed via `on_get_wz()`/`on_set_wz()` |
| **rust-z80emu** | Yes | `RegisterPair memptr` | Full tracking with flavour-dependent `memptr_mix()` for LD (rr),A |

### Analysis

All six implementations track MEMPTR. The redcode/Z80 project stands out with the most
comprehensive documentation (dedicated wiki pages). The rust-z80emu implementation is notable for
encoding flavour-specific MEMPTR behavior -- the BM1 (Soviet clone) variant zeroes the high byte
in `memptr_mix()`, matching real BM1 silicon behavior.

retrogolib's implementation is complete and correct, covering all standard MEMPTR update points.

---

## 4. Q Register

The internal flag-modification tracker that enables correct undocumented X/Y flag behavior in
SCF/CCF instructions. Q tracks whether the previous instruction modified flags; SCF/CCF use this
to compute X/Y as `(A | (F & ~Q)) & 0x28`.

| Project | Tracked | Implementation |
|---------|---------|---------------|
| **retrogolib** | Yes | `q uint8` field; captures F after flag-modifying instructions |
| **floooh/chips** | No | SCF/CCF use only `cpu->a` for X/Y bits |
| **redcode/Z80** | Yes | Conditional via `Z80_WITH_Q`; `FLAGS` macro copies F to Q on flag-affecting ops; `Q_0` zeroes Q on non-flag ops. Supports per-model XQ/YQ option flags |
| **superzazu/z80** | No | No Q tracking |
| **kosarev/z80** | No | No Q tracking |
| **rust-z80emu** | Yes | Encoded in `Flavour` trait; `flags_modified`/`last_flags_modified` booleans; `get_q()` returns flavour-specific result (NMOS/BM1 track, CMOS doesn't) |

### Analysis

Only three of six implementations track Q: retrogolib, redcode/Z80, and rust-z80emu. This is the
most obscure undocumented Z80 behavior -- it was only fully characterized in 2018. Even floooh/chips,
widely regarded as one of the best Z80 emulators, does not implement it.

retrogolib's Q implementation is straightforward and correct. The redcode/Z80 approach is the most
configurable (per-model XQ/YQ factors), while rust-z80emu elegantly encodes it as part of the
CPU flavour system.

---

## 5. DD/FD Prefix Handling

How each implementation handles the IX/IY prefix opcodes, especially undocumented passthrough
behavior where invalid DD/FD opcodes execute as their unprefixed equivalents.

| Project | Approach | Undocumented IXH/IXL/IYH/IYL | Passthrough |
|---------|----------|-------------------------------|-------------|
| **retrogolib** | Separate DD/FD opcode tables (256 entries each) + undoc files | Yes, via `emulation_dd_undoc.go` / `emulation_fd_undoc.go` | Yes, invalid opcodes execute unprefixed (+4 T-states) |
| **floooh/chips** | `hlx_idx` register remapping; `hlx[3]` union overlays HL/IX/IY | Yes, automatic via register substitution | Yes, non-indirect opcodes fall through to unprefixed handler with IX/IY substituted |
| **redcode/Z80** | Shared `xy_insn_table` with temporary `XY` register copy | Yes, via `o_p_table`/`w_table` referencing `xy.uint8_values` | Yes, `xy_illegal()` falls through to `insn_table` |
| **superzazu/z80** | Single `exec_opcode_ddfd()` with `uint16_t* iz` pointer | Yes, via `IZH`/`IZL` macros | Yes, unrecognized opcodes call `exec_opcode()` |
| **kosarev/z80** | `iregp_kind` enum set on prefix, resets after decode | Yes, through iregp register substitution | Yes, decoder dispatches to unprefixed path |
| **rust-z80emu** | `maybe_prefix` option in `'repeat` loop | Yes, through prefix-aware instruction macros | Yes, invalid DD/FD opcodes set `maybe_prefix = None` and continue loop |

### Analysis

The floooh/chips approach is the most elegant: a union overlay of HL/IX/IY as a 3-element array,
with `hlx_idx` selecting which one to use. This means the same generated code handles HL, IX, and
IY operations with zero additional dispatch logic.

retrogolib uses separate DD/FD opcode tables, which is the most explicit approach. This adds some
code size but makes the mapping completely transparent and easy to debug. The dedicated undocumented
instruction files (`emulation_dd_undoc.go`, `emulation_fd_undoc.go`) are a clean separation of
concerns that no other implementation provides.

---

## 6. Cycle Accuracy

| Project | Granularity | Approach |
|---------|-------------|----------|
| **retrogolib** | Instruction-level | Cycle count stored per opcode in instruction table; added after execution |
| **floooh/chips** | T-state (cycle-exact) | Each `z80_tick()` call = 1 T-state; `cpu->step` counter tracks position within microcode. Supports WAIT pin stalling |
| **redcode/Z80** | Instruction-level | Handler functions return T-state count; total accumulated in `self->cycles` |
| **superzazu/z80** | Instruction-level | Pre-computed cycle lookup arrays (`cyc_00[256]`, etc.) |
| **kosarev/z80** | T-state level | Explicit `on_tick()` calls at each bus cycle within instruction execution |
| **rust-z80emu** | T-state level | `Clock` trait with `add_m1()`, `add_mreq()`, `add_io()` called at each bus cycle |

### Analysis

floooh/chips achieves true cycle-exact emulation with WAIT pin support, which is essential for
accurate ZX Spectrum memory contention emulation. kosarev/z80 and rust-z80emu also provide T-state
granularity through their callback mechanisms.

retrogolib, redcode/Z80, and superzazu/z80 use instruction-level cycle counting. This is sufficient
for most emulation purposes but cannot model mid-instruction memory contention effects. For systems
like the ZX Spectrum 48K where video memory access causes cycle stretching during instruction
execution, T-state-level timing would be needed.

---

## 7. Interrupt Handling

| Project | IM 0 | IM 1 | IM 2 | NMI | IFF1/IFF2 |
|---------|------|------|------|-----|-----------|
| **retrogolib** | Basic (RST 38h equivalent) | Yes | Yes | Yes | Yes |
| **floooh/chips** | Full (executes data bus instruction) | Yes | Yes | Yes (edge-triggered) | Yes |
| **redcode/Z80** | Full (with `Z80_WITH_FULL_IM0` trampoline mechanism) | Yes | Yes | Yes (with rejection logic) | Yes, including LD A,{I\|R} bug |
| **superzazu/z80** | Yes | Yes | Yes | Yes | Yes |
| **kosarev/z80** | Yes | Yes | Yes | Yes | Yes |
| **rust-z80emu** | Yes | Yes | Yes | Yes | Yes, with flavour-specific `ACCEPTING_INT_RESETS_IFF2_EARLY` |

### Analysis

IM 0 is the most complex interrupt mode -- the CPU reads an instruction from the data bus during
the interrupt acknowledge cycle and executes it. Most software only puts RST instructions on the
bus, but the Z80 technically supports any instruction (including multi-byte ones).

floooh/chips and redcode/Z80 implement full IM 0 with actual data bus instruction execution.
redcode/Z80's approach is particularly sophisticated, using a trampoline mechanism that temporarily
redirects memory callbacks to the interrupting device.

retrogolib currently treats IM 0 as equivalent to RST 38h, which is correct for the vast majority
of Z80 software. Full IM 0 with arbitrary instruction execution is a potential enhancement.

The redcode/Z80 implementation also uniquely emulates the Zilog NMOS bug where `LD A,I` and
`LD A,R` reset the P/V flag when an INT is accepted during those instructions.

---

## 8. Hardware Interface

How each implementation communicates with the host system (memory, I/O ports).

| Project | Approach | Details |
|---------|----------|---------|
| **retrogolib** | Go interface (`Memory`) | `Read(addr)`, `Write(addr, val)`, `ReadWord()`, `WriteWord()`, `ReadPort()`, `WritePort()` |
| **floooh/chips** | Pin-based `uint64_t` | All 40 CPU pins packed into single value; caller inspects MREQ/IORQ/RD/WR signals and feeds data back next tick |
| **redcode/Z80** | C callback function pointers | `read`, `write`, `in`, `out`, `halt`, `nop`, `nmia`, `inta`, `int_fetch` callbacks on Z80 struct |
| **superzazu/z80** | C callback function pointers | `read_byte`, `write_byte`, `port_in`, `port_out` callbacks |
| **kosarev/z80** | C++ CRTP method overrides | User class inherits from `z80_cpu<MyEmulator>` and overrides `on_read()`, `on_write()`, `on_input()`, `on_output()` |
| **rust-z80emu** | Rust traits (`Memory`, `Io`, `Clock`) | `read_mem()`, `write_mem()`, `read_opcode()` (with IR refresh value), `read_io()`, `write_io()` |

### Analysis

The floooh/chips pin-based interface is the most hardware-faithful: it models the actual electrical
signals on the Z80's pins, enabling accurate system-level emulation. The downside is that the
caller must understand bus protocols.

retrogolib's Go interface approach is clean and idiomatic. It's most similar to rust-z80emu's trait
approach, with both providing clean abstraction boundaries that allow different hardware
implementations.

kosarev/z80's CRTP approach eliminates virtual dispatch overhead entirely -- all method calls are
resolved at compile time. This gives C++-like performance with the flexibility of callbacks.

---

## 9. Thread Safety

| Project | Thread-Safe | Mechanism |
|---------|-------------|-----------|
| **retrogolib** | Yes | `sync.RWMutex` on all public CPU accessors |
| **floooh/chips** | No | Single-threaded; caller manages synchronization |
| **redcode/Z80** | No | Single-threaded |
| **superzazu/z80** | No | Single-threaded |
| **kosarev/z80** | No | Single-threaded |
| **rust-z80emu** | No | Single-threaded (Rust's ownership model prevents data races at compile time) |

### Analysis

retrogolib is the only implementation with built-in thread safety. This is valuable for debuggers,
profilers, and tools that need to inspect CPU state from a different thread while emulation runs.
All other implementations leave synchronization to the caller, which is the traditional approach
for emulators where the CPU step function is called from a single emulation thread.

---

## 10. Dual-Purpose Architecture

| Project | Emulation | Disassembly | Static Analysis |
|---------|-----------|-------------|-----------------|
| **retrogolib** | Yes | Yes (instruction metadata in opcode tables) | Yes (operand types, sizes, addressing modes) |
| **floooh/chips** | Yes | No | No |
| **redcode/Z80** | Yes | No | No |
| **superzazu/z80** | Yes | No | No |
| **kosarev/z80** | Yes | Partial (Python bindings can disassemble) | No |
| **rust-z80emu** | Yes | No | No |

### Analysis

retrogolib is unique in storing instruction metadata (mnemonic name, operand types, byte size,
cycle count) alongside each handler in the same opcode tables. This enables disassembly and static
analysis without maintaining a separate instruction database -- a significant advantage for tool
development.

kosarev/z80 provides some disassembly capability through its Python bindings, but this is a
separate code path rather than an integrated feature of the instruction tables.

---

## 11. Test Infrastructure

| Project | ZEXDOC | ZEXALL | FUSE | SingleStep | Other |
|---------|--------|--------|------|------------|-------|
| **retrogolib** | Yes (67/67) | Yes (67/67) | No | Yes (1,609 JSON test vectors) | Unit tests |
| **floooh/chips** | Yes | Yes | Partial (minor XF/YF exceptions) | No | Generated test framework |
| **redcode/Z80** | Yes | Yes | No | No | Hardware-validated against ZX Spectrum 48K |
| **superzazu/z80** | Yes | Yes | No | No | Basic unit tests |
| **kosarev/z80** | Yes | Yes | No | No | cputest, 8080pre, 8080exer, 8080exm |
| **rust-z80emu** | Implied (via SPECTRUSTY) | Implied | Yes (via SPECTRUSTY) | No | Cycle tests, flavour tests, shuffle integration test |

### Analysis

retrogolib has the most comprehensive test infrastructure among these implementations, combining
three distinct test approaches: ZEXDOC/ZEXALL (exerciser-based validation), SingleStep (bus-level
cycle accuracy from Ares Z80 core translation), and unit tests.

The redcode/Z80 project is unique in validating against physical Z80 hardware, which provides
confidence in undocumented behavior accuracy that software-only tests cannot match.

---

## 12. Code Size and Complexity

| Project | Total LOC | Files | Instruction Handlers | Generated Code |
|---------|-----------|-------|---------------------|----------------|
| **retrogolib** | 10,737 | 39 | ~394 | No |
| **floooh/chips** | ~4,191 | 1 (+codegen) | ~1,700 switch cases | Yes (Python -> C) |
| **redcode/Z80** | ~3,894 | 2 | ~256 function pointers | No (but has table generators) |
| **superzazu/z80** | ~1,814 | 2 | ~256 switch cases | No |
| **kosarev/z80** | ~4,051 | 1 | Structured decode | No |
| **rust-z80emu** | ~10,247 | 20 | Macro-expanded | No (but macro-generated) |

### Analysis

retrogolib and rust-z80emu are the largest implementations. retrogolib's size comes from its
dual-purpose architecture (separate instruction table definitions + emulation handlers + metadata)
and modular file organization. rust-z80emu's size comes from three CPU flavour variants and
comprehensive trait abstractions.

superzazu/z80 is the most compact at ~1,814 LOC -- a testament to how concise a correct Z80
emulator can be in C, at the cost of no Q register and limited undocumented behavior support.

floooh/chips achieves remarkable density through code generation, packing cycle-exact emulation
of the entire Z80 into a single header file.

---

## 13. Unique Strengths by Project

### retrogolib
- **Dual-purpose opcode tables** -- emulation + disassembly from single source of truth
- **Thread safety** -- built-in RWMutex for concurrent debugger access
- **Q register tracking** -- one of only 3 implementations
- **Three test suites** -- ZEXDOC + ZEXALL + SingleStep vectors
- **Go ecosystem** -- clean interface-based design, modular file organization

### floooh/chips
- **Cycle-exact T-state emulation** -- true pin-level accuracy with WAIT support
- **Pin-centric interface** -- models actual Z80 electrical signals
- **Code generation** -- YAML-driven, maintainable instruction definitions
- **Performance** -- ~556 MHz emulated clock speed
- **Tick merging** -- batches idle cycles for efficiency

### redcode/Z80
- **Most documented undocumented behavior** -- MEMPTR wiki is a community reference
- **Hardware validation** -- tested against real Sinclair ZX Spectrum 48K
- **Full IM 0** -- trampoline mechanism for arbitrary instruction execution
- **LD A,{I|R} interrupt bug** -- Zilog NMOS silicon bug emulated
- **Per-model Q factors** -- XQ/YQ option flags for different Z80 manufacturers

### superzazu/z80
- **Minimal and clean** -- ~1,814 LOC, easy to read and understand
- **Individual flag bitfields** -- similar approach to retrogolib
- **Educational quality** -- ideal reference implementation

### kosarev/z80
- **Zero-overhead CRTP** -- all dispatch resolved at compile time
- **Lazy flag evaluation** -- defers flag computation until read
- **Structured opcode decoding** -- x/y/z/p/q bit field decomposition
- **Dual Z80/i8080** -- single codebase supports both CPUs
- **Module replacement** -- any standard component can be swapped

### rust-z80emu
- **CPU flavour variants** -- NMOS/CMOS/BM1 encode real silicon differences
- **Flavour-specific MEMPTR** -- BM1 clone zeroes high byte in `memptr_mix()`
- **Trait-based composition** -- idiomatic Rust `Memory`/`Io`/`Clock` traits
- **Zero-cost debugging** -- debugger hooks compile away in release builds
- **Battle-tested** -- powers the SPECTRUSTY ZX Spectrum emulator

---

## 14. Potential Improvements for retrogolib

Based on this cross-language comparison, areas where retrogolib could learn from other
implementations:

### From floooh/chips
- **T-state granularity**: Consider adding optional cycle-level callbacks for systems
  requiring mid-instruction timing accuracy (ZX Spectrum memory contention)
- **Code generation**: YAML-driven instruction definitions could reduce maintenance burden

### From redcode/Z80
- **Full IM 0**: Support executing arbitrary instructions from the data bus during
  interrupt acknowledge, not just RST 38h
- **LD A,{I|R} interrupt bug**: The Zilog NMOS bug where P/V is reset when INT fires
  during these instructions
- **Per-model Q behavior**: Different Z80 manufacturers (Zilog, NEC, ST) have subtly
  different Q register behavior

### From rust-z80emu
- **CPU flavour variants**: NMOS vs CMOS vs clone differences in undocumented behavior
  could be encoded as configuration options
- **Flavour-specific MEMPTR**: The BM1 clone's different `memptr_mix` behavior

### From kosarev/z80
- **Lazy flag evaluation**: Could improve performance for instruction sequences where
  flags are set but never read before being overwritten

---

## 15. Summary

retrogolib's Z80 implementation holds its own against the best non-Go implementations:

| Dimension | retrogolib Standing |
|-----------|-------------------|
| Instruction dispatch | On par (function pointer tables, like redcode/Z80) |
| Flag handling | Unique approach (individual fields, shared only with superzazu) |
| MEMPTR/WZ | Complete (all implementations track this) |
| Q register | Top tier (only 3 of 6 implement this) |
| DD/FD passthrough | Complete with dedicated undocumented instruction files |
| Cycle accuracy | Instruction-level (3 of 6 do T-state level) |
| Interrupt handling | Good; IM 0 could be enhanced to full bus instruction execution |
| Thread safety | Unique among all compared implementations |
| Dual-purpose design | Unique -- no other implementation integrates disassembly metadata |
| Test coverage | Strongest (3 complementary test suites) |

The implementation is accurate, well-tested, and architecturally unique in its dual-purpose
design and thread safety. The main gaps relative to the best-in-class implementations are
T-state-level cycle accuracy and full IM 0 interrupt mode support.
