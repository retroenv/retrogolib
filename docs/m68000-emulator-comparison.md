# Motorola 68000 Emulator Comparison: retrogolib vs Best-in-Class

## Overview

This document compares retrogolib's M68000 implementation against the best open-source 68000
emulators across all languages. The 68000 ecosystem is mature, with implementations ranging
from instruction-accurate interpreters to gate-level FPGA cores.

---

## Emulators Compared

| Emulator | Language | License | Accuracy Level | Status |
|----------|----------|---------|----------------|--------|
| **retrogolib** | Go | MIT | Instruction-accurate | Active (2026) |
| **Musashi** | C | MIT-like | Instruction-accurate | Mature, maintained |
| **Moira** | C++ | MIT | Bus-cycle-exact | Active (vAmiga) |
| **MAME 68k** | Generated C++ | GPL-2.0 | Bus-cycle-exact | Active (MAME) |
| **FX68K** | SystemVerilog | GPL-2.0 | Gate-level exact | Stable (MiSTer) |
| **Cyclone68000** | ARM Assembly | Custom | Instruction-accurate | Stable |
| **m68k-rs** | Rust | MIT | Instruction-accurate | Active |
| **r68k** | Rust | MIT | Instruction-accurate | Stable |

---

## Detailed Comparison

### 1. Musashi (C) -- The Gold Standard

**Repository:** https://github.com/kstenerud/Musashi
**Author:** Karl Stenerud
**Used in:** MAME (until 2023), many standalone emulators

Musashi was the de facto reference 68000 emulator for over two decades. It uses a code
generator (`m68kmake`) that reads a table of instruction descriptions and produces optimized
C dispatch code.

| Feature | Musashi | retrogolib |
|---------|---------|------------|
| Instructions | All 68000 + 68010/020/040 | All 68000 |
| Addressing modes | All 14 | All 14 |
| Cycle accuracy | Instruction-level | Instruction-level |
| Prefetch emulation | Basic (2-word) | None |
| Address error | Full (exception + stack frame) | Not checked |
| Bus error | Full | Not implemented |
| Privilege mode | Full | Full |
| Exception stack frames | Full (type 0/2 frames) | Basic (PC + SR only) |
| CPU variants | 68000/010/EC020/020/EC030/040 | 68000 only |
| Code generation | Yes (m68kmake) | No (hand-written) |
| LOC | ~15,000 (generated) | ~5,900 |
| Thread safety | None | Mutex-based |

**Key differences:**
- Musashi supports the entire 680x0 family through a single codebase with compile-time
  variant selection. retrogolib focuses exclusively on the 68000.
- Musashi implements a 2-word prefetch queue that affects timing of certain instructions.
  retrogolib fetches instructions with zero latency.
- Musashi generates full Type 0 and Type 2 exception stack frames (6 or 14 bytes) including
  the faulting address and access type for bus/address errors. retrogolib pushes only PC and SR.
- Musashi has no thread safety; retrogolib protects CPU state with sync.RWMutex.

---

### 2. Moira (C++) -- Modern, Faster Than Musashi

**Repository:** https://github.com/dirkwhoffmann/Moira
**Author:** Dirk W. Hoffmann
**Used in:** vAmiga

Moira achieves bus-cycle-exact accuracy while being *faster* than Musashi through aggressive
C++ template metaprogramming. It fires synchronization callbacks at each individual bus access
within an instruction, enabling cycle-accurate interleaving with other system components.

| Feature | Moira | retrogolib |
|---------|-------|------------|
| Instructions | All 68000 | All 68000 |
| Addressing modes | All 14 | All 14 |
| Cycle accuracy | **Bus-cycle-exact** | Instruction-level |
| Prefetch emulation | **Full (IRC/IRD pipeline)** | None |
| Address error | **Full with stack frame** | Not checked |
| Bus error | **Full** | Not implemented |
| Privilege mode | Full | Full |
| Exception stack frames | **Full (all frame types)** | Basic |
| CPU variants | 68000, 68010 | 68000 |
| Sync callbacks | **Per bus cycle** | None |
| LOC | ~12,000 | ~5,900 |
| Thread safety | None | Mutex-based |

**Key differences:**
- Moira is the accuracy leader among software emulators. It models the 68000's two-word
  prefetch pipeline (IRC and IRD registers) and fires a `sync()` callback on every bus
  access. This is essential for Amiga emulation where copper, blitter, and DMA interleave
  with CPU bus cycles.
- Moira uses C++ templates to specialize instruction handlers at compile time, eliminating
  runtime branching for operand sizes and addressing modes. This makes it faster than Musashi
  despite higher accuracy.
- retrogolib's simpler architecture is easier to understand and modify but cannot model
  mid-instruction timing effects.

---

### 3. MAME Microcode Core (Generated C++) -- Die-Accurate

**Repository:** https://github.com/mamedev/mame (src/devices/cpu/m68000/)
**Author:** Olivier Galibert
**Used in:** MAME (since February 2023)

MAME replaced Musashi with a new core derived from actual die photography of a real 68000
chip. The microcode was extracted from silicon analysis and used to generate C++ emulation
code. This core was also used to generate the SingleStepTests JSON test vectors.

| Feature | MAME 68k | retrogolib |
|---------|----------|------------|
| Instructions | All 68000 + variants | All 68000 |
| Cycle accuracy | **Bus-cycle-exact** | Instruction-level |
| Derivation | **Die photography** | Reference manuals |
| Prefetch emulation | **Full** | None |
| Address/Bus error | **Full** | Not implemented |
| Exception frames | **Full** | Basic |
| Test vectors | **Generated SingleStepTests** | Unit tests |
| LOC | ~20,000+ (generated) | ~5,900 |

**Key differences:**
- The MAME core's microcode is derived from actual silicon, making it definitionally
  accurate for the original Motorola 68000. retrogolib was written from documentation.
- This core generated the SingleStepTests JSON vectors (73,000+ test cases) that serve
  as the ground truth for validating other emulators.
- Currently slower than both Musashi and Moira due to the microcode interpretation overhead.

---

### 4. FX68K (SystemVerilog) -- Gate-Level FPGA Core

**Repository:** https://github.com/ijor/fx68k
**Author:** Jorge Cwik
**Used in:** MiSTer FPGA cores (Genesis, Amiga, Atari ST, Neo Geo)

FX68K is not a software emulator but an RTL (Register Transfer Level) description of the
68000 that runs on FPGA hardware. It replicates the actual silicon behavior at the gate
level and is functionally indistinguishable from a real 68000 chip.

| Feature | FX68K | retrogolib |
|---------|-------|------------|
| Accuracy | **Gate-level exact** | Instruction-level |
| Platform | FPGA (Cyclone V) | Software (any CPU) |
| Cycle accuracy | **Exact (real hardware timing)** | Basic |
| Address/Bus error | **Exact** | Not implemented |
| Power consumption | ~5W (FPGA) | N/A |

FX68K represents the ultimate accuracy target but is in a fundamentally different category
(hardware vs software emulation).

---

### 5. Cyclone68000 (ARM Assembly) -- Speed-Optimized

**Repository:** https://github.com/notaz/cyclone68000
**Author:** Dave (FinalDave), maintained by notaz
**Used in:** PicoDrive (Sega Genesis emulator for handhelds)

Cyclone is a highly optimized 68000 emulator written in ARM assembly, designed for
resource-constrained devices like the GP2X and early smartphones. It leverages ARM's
condition codes and barrel shifter to map 68000 operations efficiently.

| Feature | Cyclone | retrogolib |
|---------|---------|------------|
| Instructions | All 68000 | All 68000 |
| Cycle accuracy | Instruction-level | Instruction-level |
| Platform | ARM only | Cross-platform (Go) |
| Performance | **Extremely fast** | Moderate |
| Address error | Partial | Not checked |
| Code generation | Table-driven | Hand-written |
| LOC | ~12,000 (ARM asm) | ~5,900 |

**Key differences:**
- Cyclone sacrifices portability for raw speed on ARM. retrogolib runs anywhere Go compiles.
- Cyclone uses a run-slice model where you specify a cycle budget and it runs until exhausted.
  retrogolib uses step-at-a-time execution.

---

### 6. m68k-rs (Rust) -- Modern Rust Implementation

**Repository:** https://github.com/benletchford/m68k-rs (and related crates)
**Author:** Ben Letchford and others

The Rust 68000 ecosystem has several implementations. m68k-rs stands out for supporting
the broadest CPU family (68000 through 68040 including FPU and MMU).

| Feature | m68k-rs | retrogolib |
|---------|---------|------------|
| Instructions | All 68000-68040 | All 68000 |
| Addressing modes | All 14 | All 14 |
| Cycle accuracy | Instruction-level | Instruction-level |
| FPU/MMU | Yes (68040) | No |
| Unsafe code | **Zero** | N/A (Go) |
| Test suites | SingleStepTests + Musashi + custom | Unit tests |
| CPU variants | 68000/010/020/030/040 | 68000 only |

**Key differences:**
- m68k-rs validates against three independent test suites including SingleStepTests.
  retrogolib uses hand-written unit tests.
- m68k-rs supports the full 680x0 family. retrogolib targets only the original 68000.
- Both achieve memory safety through their respective languages (Rust ownership, Go GC).

---

### 7. r68k (Rust) -- Research-Focused

**Repository:** https://github.com/marhel/r68k
**Author:** Martin Hellspong

r68k is a pure Rust 68000 emulator focused on correctness and clean architecture. It uses
an enum-based instruction representation similar to retrogolib's struct-based approach.

| Feature | r68k | retrogolib |
|---------|------|------------|
| Instructions | All 68000 | All 68000 |
| Addressing modes | All 14 | All 14 |
| Decoder | Enum-based | Line-based hierarchical |
| Thread safety | Rust ownership | Mutex-based |

---

## Feature Matrix

| Feature | retrogolib | Musashi | Moira | MAME | Cyclone | m68k-rs |
|---------|-----------|---------|-------|------|---------|---------|
| All 68000 instructions | Yes | Yes | Yes | Yes | Yes | Yes |
| All 14 addressing modes | Yes | Yes | Yes | Yes | Yes | Yes |
| Byte/Word/Long sizes | Yes | Yes | Yes | Yes | Yes | Yes |
| Supervisor/User mode | Yes | Yes | Yes | Yes | Yes | Yes |
| All 16 conditions (Bcc) | Yes | Yes | Yes | Yes | Yes | Yes |
| BCD arithmetic | Yes | Yes | Yes | Yes | Yes | Yes |
| MOVEM register masks | Yes | Yes | Yes | Yes | Yes | Yes |
| A7 byte alignment | Yes | Yes | Yes | Yes | Yes | Yes |
| Divide by zero exception | Yes | Yes | Yes | Yes | Yes | Yes |
| CHK exception | Yes | Yes | Yes | Yes | Yes | Yes |
| TRAP vectors (0-15) | Yes | Yes | Yes | Yes | Yes | Yes |
| Line A/F traps | Yes | Yes | Yes | Yes | Yes | Yes |
| **Address error check** | **No** | Yes | Yes | Yes | Partial | Yes |
| **Bus error** | **No** | Yes | Yes | Yes | Partial | Partial |
| **Prefetch queue** | **No** | Basic | Full | Full | No | No |
| **Bus-cycle callbacks** | **No** | No | Yes | Yes | No | No |
| **Full exception frames** | **No** | Yes | Yes | Yes | Partial | Partial |
| 680x0 family support | No | Yes | Partial | Yes | No | Yes |
| Thread safety | **Yes** | No | No | No | No | N/A |
| Tracing support | **Yes** | Yes | Yes | Yes | No | No |
| Hierarchical decoder | **Yes** | Table | Template | uCode | Table | Match |

---

## Accuracy Tiers

### Tier 1: Gate-Level Exact
- **FX68K** -- FPGA core from die analysis, indistinguishable from real silicon.

### Tier 2: Bus-Cycle-Exact
- **Moira** -- Fires sync callbacks per bus access. Full prefetch pipeline.
- **MAME 68k** -- Microcode derived from die photography. Silicon-accurate timing.

### Tier 3: Instruction-Accurate with Full Exception Model
- **Musashi** -- Complete exception stack frames, basic prefetch, all 680x0 variants.
- **m68k-rs** -- Validated against SingleStepTests, full 680x0 family.

### Tier 4: Instruction-Accurate
- **retrogolib** -- All instructions, all addressing modes, basic exception handling.
- **Cyclone68000** -- Speed-optimized, ARM-only.
- **r68k** -- Clean architecture, research-focused.

---

## retrogolib Strengths

1. **Clean architecture**: Line-based hierarchical opcode decoder is elegant and maintainable
   compared to flat lookup tables or code generators. Each of the 16 line decoders is a
   focused function handling a logical group of instructions.

2. **Thread safety**: Only 68000 emulator with built-in mutex protection for concurrent
   access. All other implementations are single-threaded.

3. **Dual-purpose design**: Instruction definitions support both emulation and static
   analysis (disassembly, code generation) through category sets and instruction metadata.

4. **Tracing infrastructure**: Built-in trace step recording for debugging and analysis.

5. **Go ecosystem**: Cross-platform, memory-safe, easy to integrate into Go applications.
   No C FFI, no unsafe code, no code generation step.

6. **Compact codebase**: ~5,900 LOC delivers all 75 instructions, all 14 addressing modes,
   privilege modes, and exception handling. Musashi needs ~15,000 LOC (generated) for
   comparable instruction coverage.

7. **Complete instruction set**: Every standard 68000 instruction is implemented with
   all three operand sizes. No stubs or TODOs.

---

## retrogolib Gaps

### Gap 1: Address Error Checking (High Priority)
Word and long accesses at odd addresses should trigger an Address Error exception
(vector 3). The real 68000 generates a special stack frame with the faulting address,
access type (read/write), and function code. retrogolib silently allows misaligned access.

**Impact:** Programs with alignment bugs will run incorrectly instead of trapping.
Software that deliberately tests for address errors (OS trap handlers) won't work.

**Reference:** Musashi `m68k_read_memory_16()` checks alignment and calls
`m68k_exception_address_error()`.

### Gap 2: Bus Error Support (Medium Priority)
Bus errors (vector 2) occur when external hardware signals an invalid memory access
(e.g., accessing non-existent memory). The Bus interface has no mechanism for memory
implementations to signal bus errors.

**Impact:** Systems with memory-mapped I/O that need to report access violations
(Amiga, Atari ST) cannot properly emulate bus error behavior.

### Gap 3: Exception Stack Frames (Medium Priority)
The 68000 has two stack frame formats:
- **Type 0 (6 bytes):** SR + PC. Used for most exceptions.
- **Type 2 (14 bytes):** SR + PC + instruction word + fault address + access info.
  Used for bus error and address error.

retrogolib only implements Type 0 frames. Type 2 frames are needed for proper
bus/address error recovery.

**Impact:** OS-level code that examines exception stack frames (page fault handlers,
memory protection) won't work correctly.

### Gap 4: Prefetch Queue Emulation (Low Priority)
The real 68000 has a 2-word prefetch pipeline. Certain self-modifying code patterns
and hardware register accesses depend on the timing of when instructions are fetched
versus when they are executed.

**Impact:** Affects accuracy for:
- Amiga copper/blitter timing (critical for demos)
- Self-modifying code (rare in normal programs)
- Certain copy protection schemes

**Reference:** Moira models the full IRC/IRD prefetch pipeline. Musashi has a basic
2-word prefetch.

### Gap 5: Cycle Timing Accuracy (Low Priority)
retrogolib assigns fixed cycle counts per instruction. The real 68000 has variable
timing based on:
- Addressing mode (each mode adds different cycle overhead)
- Operand size (long operations take more cycles than byte/word)
- Multiplication/division (timing depends on operand values)
- Memory wait states (system-dependent)

**Impact:** Cycle-counted timing loops and raster effects won't be accurate.
Acceptable for functional emulation but not for demo-accurate or hardware-accurate
emulation.

### Gap 6: TAS Atomic Bus Cycle (Low Priority)
The TAS (Test and Set) instruction performs an indivisible read-modify-write bus
cycle. The real 68000 holds the bus for the entire operation. retrogolib implements
TAS as separate read and write operations.

**Impact:** Affects multiprocessor systems (rare for 68000) and certain Sega Genesis
cartridge detection schemes.

---

## Gap Closure Recommendations

### Phase 1: Address Error Checking
- Add alignment validation to `readEA` and `writeEA` for word/long accesses.
- Trigger `processException(VectorAddressError)` on misaligned access.
- Extend `processException` to save Type 2 stack frames for address errors.

### Phase 2: Exception Stack Frames
- Implement Type 2 stack frame (14 bytes) with faulting address and access info.
- Track instruction word and access type during instruction execution.
- Apply Type 2 frames to address error and bus error exceptions.

### Phase 3: Bus Error Support
- Extend the Bus interface or Memory interface to return errors on read/write.
- Translate bus-level errors into 68000 Bus Error exceptions with Type 2 frames.

### Phase 4: Variable Cycle Timing
- Add EA cycle costs based on addressing mode (Motorola M68000UM Table 8-4).
- Add size-dependent timing for long word operations.
- Add data-dependent timing for MULU/MULS/DIVU/DIVS.

### Phase 5: Prefetch Queue (Optional)
- Model the 2-word prefetch pipeline (IRC and IRD registers).
- Fetch next instruction word during execute phase of current instruction.
- Only needed for cycle-exact Amiga/demo emulation.

---

## Conclusion

retrogolib's 68000 implementation is a solid, complete, instruction-accurate emulator
that covers the full 68000 instruction set with proper privilege modes and basic exception
handling. At ~5,900 LOC it is remarkably compact compared to established implementations.

The main gap versus Musashi and Moira is in exception handling fidelity (address errors,
stack frames) and timing accuracy (prefetch, variable cycles). These gaps can be closed
incrementally -- phases 1-2 would bring retrogolib to parity with Musashi for functional
emulation, while phases 4-5 would approach Moira's accuracy level.

For use cases like Sega Genesis, Atari ST, or early Macintosh emulation where instruction
correctness matters more than cycle-exact timing, retrogolib is already suitable. For
Amiga emulation where copper/blitter interleaving requires bus-cycle accuracy, Moira
remains the reference implementation.
