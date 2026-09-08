# Work Branch Changelog

Tracks the committed changes introduced by `work2` since its common ancestor
with the current remote-tracking `main` branch.

**Last Updated:** 2026-09-07

---

## Current Branch State

- Comparison range: `origin/main...HEAD`, currently
  `1cac5af...dfc8a64` on `work2`, using the available remote-tracking ref.
- Merge base: `2901134` (`cartridge: fix iNES header mirroring flags`, 2026-09-01).
- Committed branch delta: 193 files, with 151 added and 42 modified; 27,799
  insertions and 1,531 deletions.
- `git diff --name-status --find-renames origin/main...HEAD` reports no deleted
  or renamed files in the current range.
- These statistics cover committed changes only. The review fixes below and
  this documentation refresh are uncommitted and excluded from the totals.

## Changes Already Absorbed From `main`

The merge base already contains earlier branch work that is therefore absent
from the current branch-side diff:

- The base CHIP-8, 6502/65C02, x86, and Z80 implementations and their initial
  architecture registration.
- Earlier shared test infrastructure and test data support; this branch only
  contributes the CPU targets and CHIP-8 additions described below.
- Earlier CLI, configuration, and logging cleanups.
- The NES cartridge iNES mirroring correction at the merge base.

The `set.Sorted` and `set.SortedFunc` changes still appear in the branch-side
diff because they follow the merge base. Their implementation and tests are
also present at the current `origin/main` tip; they are not unique to `work2`.

## Branch-Specific Changes

### New CPU Emulators

- **WDC 65C816:** Adds `arch/cpu/cpu65816/` with emulation and native modes,
  24-bit addressing, interrupts, tracing, structured instruction metadata, and
  unit, opcode, effective-address, and SingleStepTests coverage.
- **Motorola 68000:** Adds `arch/cpu/cpu68000/` with bus and memory abstractions,
  effective-address handling, interrupts, tracing, structured instruction
  metadata, and unit, opcode, memory, and SingleStepTests coverage.
- **Motorola 6809:** Adds `arch/cpu/cpu6809/` with indexed addressing, page
  0/10/11 opcode tables, interrupts, tracing, structured instruction metadata,
  and unit and opcode coverage.
- **Sharp SM83:** Adds `arch/cpu/sm83/` with Game Boy CPU behavior, CB-prefixed
  instructions, interrupts, tracing, structured instruction metadata, and
  SingleStepTests coverage.

### Existing CPU Refactors and Correctness

- **Cross-architecture structure:** Aligns instruction names, instruction
  registries, option handling, constructor setup, and focused tests across the
  CHIP-8, 6502, 65C816, 68000, 6809, SM83, x86, and Z80 packages. Small internal
  declarations remain with their cohesive implementation instead of occupying
  single-purpose files.
- **Public APIs:** Replaces exposed option implementation details with private
  option state and typed functional options. Obsolete compatibility APIs are
  removed rather than retained as deprecated wrappers.
- **CHIP-8:** Makes COSMAC VIP behavior the default and adds explicit quirks for
  later interpreter variants. Corrects display-wait, key-release, shift, logic,
  jump, load/store, sprite, arithmetic-flag, and bounds behavior; expands unit
  tests and adds Timendus ROM conformance tests with screenshot comparison.
- **6502/65C02:** Reorganizes instruction metadata and unofficial operations,
  tightens constructor, memory, interrupt, stepping, and option behavior, and
  adds focused CPU, interrupt, option, opcode, and step tests. Adds
  `InstructionsForVariant`, variant-specific memory-effect classification,
  and corrected branch and indexed read-modify-write timing. `BRK` metadata
  describes its one-byte encoding; interrupt return still skips the padding byte.
- **x86 and Z80:** Separates substantial instruction-name and registry data from
  core instruction definitions. Z80 also adopts private option state, a typed
  pre-execution hook, nil-option handling, and registry and option tests.
- **SM83 metadata:** Includes high-memory loads through C in the `LDH` registry.
- **Sets:** Adds `Sorted` for ordered values and `SortedFunc` for custom
  comparators; both return a sorted slice without changing the set.

### System Foundations

- **Atari 2600:** Adds system and memory-map definitions, TIA and RIOT register
  constants, cartridge loading, supported banking schemes, and tests under
  `arch/system/atari2600/`.
- **TRS-80 Color Computer:** Adds CoCo memory-map, hardware, interrupt, PIA, and
  SAM definitions with tests under `arch/system/coco/`.
- **Vectrex:** Adds Vectrex memory-map, hardware, interrupt, and VIA definitions
  with tests under `arch/system/vectrex/`.

### Test and Documentation Integration

- Extends the root integration-test target for the 65C816, 68000, and SM83
  suites.
- Adds a `testdata/Makefile` target for the Timendus CHIP-8 test suite and keeps
  the existing CPU test-data targets in the aggregate workflow.
- Updates `README.md` with the new CPU packages.
- Adds gap-closure plans for the Motorola 68000 and Z80, a Commodore 64 system
  implementation plan, and this branch changelog.

## Package Review and Uncommitted Fixes

The review checked changes in every affected package group, focusing on public
contracts, instruction metadata, interrupt boundaries, memory access, and system
constants. The table records the additional work from this review; it does not
claim exhaustive hardware conformance.

| Package group | Review result |
| --- | --- |
| `arch/cpu/chip8` | Reviewed VIP defaults, quirks, waits, flag aliases, bounds, and API changes; six Timendus ROM checks pass. |
| `arch/cpu/cpu6502` | Reviewed variant registries, memory effects, constructor checks, interrupt dispatch, and cycle corrections; existing short tests pass. |
| `arch/cpu/cpu65816` | Fixed interrupt dispatch from `Step`, masked-IRQ wake from WAI, and emulation-mode program-bank clearing. |
| `arch/cpu/cpu68000` | Implemented the no-op `TriggerIRQ`, separated interrupt entry from handler execution, and fixed multibyte memory wrapping. Conformance gaps remain. |
| `arch/cpu/cpu6809` | Fixed zero-cycle hardware interrupt entry: NMI/IRQ take 19 cycles and FIRQ takes 10. |
| `arch/cpu/sm83` | Fixed interrupt step boundaries, delayed EI cancellation by DI, and HALT-bug operand fetching. Reused the decoded opcode for immediate operands. |
| `arch/cpu/x86` | Reviewed instruction-name and registry extraction; no additional behavioral fix identified. |
| `arch/cpu/z80` | Corrected INF/OUTF encoding and execution through the existing ED70/ED71 implementations, added reverse metadata, and fixed a declaration-order lint issue. |
| `arch/system/atari2600`, `register`, `cartridge` | Reviewed system/register definitions; fixed 3F bank size/count, write hotspots, and hotspot address mirroring. |
| `arch/system/coco`, `register` | Corrected SAM rate, RAM-size, and memory-map descriptions and replaced misnamed memory-map constants. |
| `arch/system/vectrex`, `register` | Reviewed memory-map and register foundations; existing short tests pass. |
| `set` | Reviewed non-mutating sorted projections; implementation and tests already match the remote main tip. |

### Interrupt and Fetch Corrections

- **65C816:** `Step` now accepts queued interrupts before fetching an opcode.
  NMI takes priority; a masked IRQ wakes WAI and remains pending until enabled.
  STP still requires reset. Interrupt entry clears PB in both modes, and trigger
  methods only queue requests under the mutex. Regression tests cover native
  and emulation stack frames, RTI, WAI, and STP. These behaviors follow the
  [WDC W65C816S datasheet](https://www.westerndesigncenter.com/wdc/documentation/w65c816s.pdf).
- **68000:** `TriggerIRQ(level)` queues the highest requested level, retaining
  masked requests. Bus and queued interrupts now finish a step at the vector
  entry with 44 cycles, including when waking STOP. Tests cover both request
  paths, stack contents, masking, and STOP wakeup.
- **6809:** Normal hardware entry now charges its stacking and vector cycles.
  CWAI retains its existing prepaid entry timing and avoids stacking twice.
  New tests check all three interrupt vectors, cycle counts, and stack sizes.
- **SM83:** An accepted interrupt consumes its own five-machine-cycle step.
  `EI; DI` cancels the delayed enable. HALT with an already pending interrupt
  suppresses one fetch increment; an interrupt arriving after HALT began wakes
  it normally. Tests cover byte/word operands, CB prefixes, relative branches,
  and the return address after `EI; HALT`, following
  [Pan Docs on interrupts](https://gbdev.io/pandocs/Interrupts.html) and
  [HALT](https://gbdev.io/pandocs/halt.html).

### Memory, Metadata, and System Corrections

- **68000 memory:** Mask every byte of word/long accesses to the 24-bit bus.
  Previously, even an aligned long access at `$FFFFFE` could panic. Regression
  tests verify big-endian reads and writes across the address-space boundary.
- **Z80:** INF/OUTF previously encoded IND/OUTD (`ED AA`/`ED AB`) and modified
  B, HL, and memory. They now use `ED 70`/`ED 71`, 12-cycle metadata, and the
  canonical port handlers. Tests check metadata, preserved registers/flags,
  and the corresponding SingleStepTests vectors.
- **Atari cartridges:** 3F uses 2 KB banks, so a 64 KB image has 32 banks.
  Writes through `$003F` are hotspots, and CPU addresses are masked to the
  6507's 13-bit bus before lookup. Tests cover offsets, bounds, and mirrors.
  ROM-size detection is documented as a heuristic; hosts still implement
  switching and keep the last 3F segment fixed. This matches
  [Stella's 3F cartridge implementation](https://github.com/stella-emu/stella/blob/master/src/emucore/Cart3F.hxx).
- **CoCo SAM:** R0/R1 select CPU rate, M0/M1 select RAM size, and TY at
  `$FFDE/$FFDF` selects the memory map. Corrected names and address tests follow
  [MAME's MC6883 SAM implementation](https://github.com/mamedev/mame/blob/master/src/devices/machine/6883sam.cpp).
- **Documentation:** Refreshed committed counts and package coverage here and
  corrected the 68000 plan's obsolete claim that SingleStepTests integration
  had not been implemented.

## API and Behavior Migration

- CHIP-8 callers use `State` and `State()` instead of `CPUState` and
  `GetState()`. COSMAC VIP is the default; select explicit quirks for other
  interpreter behavior.
- 6502 and Z80 callers configure constructors through typed functional options
  instead of the removed exported `Options`/`NewOptions` implementation API.
  The 6502 interrupt trigger is now `TriggerIRQ`, replacing `TriggerIrq`.
- With the review fixes, 65C816 callers can rely on `Step` for interrupt
  dispatch. Interrupt entry in 65C816, 68000, and SM83 is a separate step from
  the first handler instruction; instruction-count loops must account for it.
- CoCo callers using `SAMRateClear`/`SAMRateSet` must choose the intended
  operation: `SAMR0*`/`SAMR1*` for CPU rate, or the corrected
  `SAMTYClear`/`SAMTYSet` for memory mapping. The misleading names are removed.
- 3F consumers must use the corrected 2 KB offsets and bank count. A
  `TriggerBank` result of zero identifies a write hotspot for this scheme;
  the written value supplies the bank number.
- Direct users of Z80 INF/OUTF handlers must use `ParamFunc`; their obsolete
  block-I/O `NoParamFunc` handlers have been removed.

## Files

| Status | Count | Files | Purpose |
| --- | ---: | --- | --- |
| Modified | 1 | `Makefile` | Extends CPU integration-test coverage. |
| Modified | 1 | `README.md` | Lists the new CPU packages. |
| Added / Modified | 4 / 6 | `arch/cpu/chip8/` | Adds compatibility options, registry separation, correctness fixes, and Timendus ROM tests. |
| Added / Modified | 8 / 20 | `arch/cpu/cpu6502/` | Refactors structure, variant metadata, timing, and focused tests. |
| Added | 30 | `arch/cpu/cpu65816/` | Adds the WDC 65C816 emulator and tests. |
| Added | 29 | `arch/cpu/cpu68000/` | Adds the Motorola 68000 emulator and tests. |
| Added | 27 | `arch/cpu/cpu6809/` | Adds the Motorola 6809 emulator and tests. |
| Added | 24 | `arch/cpu/sm83/` | Adds the Sharp SM83 emulator and tests. |
| Added / Modified | 2 / 2 | `arch/cpu/x86/` | Moves instruction names and registry data into cohesive files. |
| Added / Modified | 4 / 9 | `arch/cpu/z80/` | Aligns registry, option, and undocumented-instruction metadata and tests. |
| Added | 8 | `arch/system/atari2600/` | Adds system/register definitions, cartridge metadata, and tests. |
| Added | 6 | `arch/system/coco/` | Adds CoCo system/register definitions and tests. |
| Added | 5 | `arch/system/vectrex/` | Adds Vectrex system/register definitions and tests. |
| Added | 4 | `docs/` | Adds implementation plans and branch tracking. |
| Modified | 2 | `set/` | Adds sorted projections and tests, also present at the remote main tip. |
| Modified | 1 | `testdata/Makefile` | Integrates the Timendus CHIP-8 ROM suite. |

The grouped counts above total the exact 151 added and 42 modified files
reported for `origin/main...HEAD`; no row represents a rename.

## Verification

Review validation on 2026-09-07:

| Check | Result |
| --- | --- |
| `go fmt ./...` | Pass. |
| `make lint` | Pass: zero golangci-lint issues and no retrogolint violations. |
| `make test` | Pass: repository short tests with the race detector. |
| CHIP-8 `TestROMConformance` | Pass: six Timendus ROM checks. |
| SM83 `TestSingleStep` and interrupt/fetch regressions | Pass. |
| 65C816 `TestSingleStep` with `-tags singlestep` | Pass. |
| Z80 `TestSingleStep/ed_70.json` and `ed_71.json` | Pass. |
| 68000 `TestSingleStep` with `-tags singlestep` | Fails: 206,547 passed and 992 failed among 207,539 executed vectors. |
| `git diff --check` and committed range/count reconciliation | Pass. |

The 68000 runner stops each file after ten failures, so these counts are not a
full-corpus pass rate. Running the same suite with a Go overlay containing the
unchanged `dfc8a64` implementations produced the same counts and failure
patterns. These failures predate the review fixes. The aggregate
`make test-integration` target is therefore not reported as passing.

## Remaining Gaps

- **68000 conformance:** Prioritize address errors and exception frames,
  active/saved stack-pointer reporting, and BCD/flag discrepancies exposed by
  the vector suite. Bus errors, timing, and prefetch fidelity also remain in
  [the gap-closure plan](cpu68000-gap-closure-plan.md). Unit-test success does
  not establish full instruction or bus-cycle accuracy.
- **System foundations:** Atari 2600, CoCo, and Vectrex additions provide
  definitions and helpers, not complete machines with video, audio, timers,
  and peripheral execution. Atari scheme detection needs a richer interface
  before ambiguous ROM formats can be selected reliably.
- **CPU state contracts:** Public register fields require callers to serialize
  direct access. State snapshots in the new CPU packages are not complete
  save/restore APIs; pending interrupts and wait/stop state need explicit
  design before promising deterministic restoration.
