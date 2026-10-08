# Phased Merge Plan for work2

## Scope and verified base

Move the changes from `work2` to `main` as separate, tested candidates. Include
all new CPU and system packages, as well as fixes to existing packages. Extract
changes by behavior. Do not merge the complete source branch in one operation.

This plan follows the file and hunk ownership format used in RetroASM's
`docs/work-branch-changes.md`. Unlike that plan, this plan includes the new
architectures. Paths below are relative to the RetroGoLib repository root.

| Item | Verified value |
| --- | --- |
| Review date | 2026-10-07 |
| Source branch | `work2` |
| Fixed source commit | `4fabd875b5ad29b7ec6e22bed501af4335905e4d` |
| Target ref | `origin/main` |
| Target commit | `9b578e97f6b30c6d3bfeb1bf75a9867927d0ea75` |
| Merge base | `9f81a91c9f86d10f539a9932428e2985a5ea4027` |
| Local `main` | `9f81a91c9f86d10f539a9932428e2985a5ea4027`; behind the target |
| Original branch range | `origin/main...HEAD`, equivalent to merge base through fixed source |
| Original branch delta | 223 files: 161 added, 62 modified; 31,841 insertions, 2,496 deletions |
| Target-to-source endpoint delta | 221 files: 161 added, 60 modified; 31,839 insertions, 2,494 deletions |
| File status detection | `git diff --name-status --find-renames origin/main...HEAD`; no renames or deletions |

A fetch on 2026-10-07 confirmed that the target commit did not move. The
counts cover committed changes at the fixed source commit. They exclude the
uncommitted review corrections recorded in “Review corrections (2026-10-07)”
below. Those corrections change behavior in every CPU package and in the new
system packages; preserve them, and include each correction in the part that
owns its file. Recompute the range statistics after they are committed.

### Work already on the target

The target is not an ancestor of the fixed source. Both tips contain the same
NES RAM correction through different commits: target `9b578e9` and source
`8ca69e4`. Thus `arch/system/nes/nes.go` and `arch/system/nes/nes_test.go` occur
in the three-dot branch inventory but have no endpoint difference. P00 marks
both files as already present. Do not extract them again.

The common base also contains:

- The CHIP-8, 6502-family, x86 metadata, and Z80 packages.
- Architecture and system identifiers in `arch/arch.go` and `arch/system.go`,
  including identifiers for the new packages planned below.
- Z80 `Bus`, `NewWithBus`, full port addresses, `IRQData`, `OnRETI`, and the
  legacy eight-bit I/O adapter. The target's `Step` interrupt path already
  samples device vectors; its separate `CheckInterrupts` path differs.
- 6502 live interrupt lines, stalls, stack events, and opt-in bus-cycle execution.
- Shared utilities, NES 2.0 support, SDL lifecycle support, and PCM audio.
- CPU test-data download targets, except the new CHIP-8 target.

Do not treat existing APIs as additions. In particular, the Z80 parts correct
interrupt acceptance, block-output ordering, metadata, and test coverage.

## Extraction method and common checks

1. Before each part, verify the current remote target and record its SHA.
   Compare it with the fixed source and the progress record below.
2. Prepare the candidate from that target in an isolated checkout. Preserve
   the source worktree. Do not base a candidate on stale local `main`.
3. Extract only the listed files and named hunks. Keep target-only fixes.
   Include a moved declaration and its removal in the same candidate.
4. Adapt imports, private helpers, tests, and docs to the candidate's current
   interfaces. A file listed in several parts must not be copied in full.
5. Run the focused checks for the part. Then format changed Go files and run
   `make build`, `make lint`, `make test`, and `git diff --check` once on the
   complete code candidate. `make test` already enables the race detector.
6. For external tests, provision the data first. Record the corpus revision,
   local corpus patches, command, executed cases, skips, and result. A skipped
   test is not a pass. Run external checks without `-short` and with `-count=1`.
7. Record API changes and test results in the candidate PR. Merge each part
   only after its required checks pass. Use the resulting target for the next
   part. Do not infer candidate correctness from a passing full-branch run.
8. If checks fail, repair that candidate before starting a dependent part.
   Independent parts can proceed when a different part is blocked.

Most source commits combine unrelated changes. Inspect a full commit patch
before considering a cherry-pick. The default method is extraction by hunk.
These commands show the original source patch and the initial endpoint patch:

```sh
git diff 9f81a91c9f86d10f539a9932428e2985a5ea4027 5b533ce3d91f023bbcebe8d6849feed53ff1774b -- arch/cpu/cpu6502
git diff 9b578e97f6b30c6d3bfeb1bf75a9867927d0ea75 5b533ce3d91f023bbcebe8d6849feed53ff1774b -- arch/cpu/cpu6502
```

After extraction, a three-dot diff can still show already-merged patches.
Use endpoint comparisons, target code, and the progress record together.
Do not overwrite a newer target file with its older source version merely
because that file remains in the original inventory.

This task writes the plan only. It does not authorize staging, committing,
pushing, rebasing, or merging. Perform those operations only when requested.
For a documentation-only candidate, review its text and references and run
`git diff --check`; do not run Go checks.

### Dependencies and validation limits

Both tips use Go `1.25.0` and `github.com/ebitengine/purego v0.11.0`.
There is no `go.mod` or `go.sum` branch delta. No dependency upgrade is needed
for the listed source patches. Check workspace overrides before candidate
validation; use `GOWORK=off` if a workspace would substitute other modules.
Keep the target's `build-platforms` target and linter versions.

Use `go test -short -race -count=1 -timeout 60s` with the package paths listed
for each focused check. Before P01, use the target's 10-second timeout unless
a recorded baseline failure requires the P01 change. The complete candidate
still uses `make test`. Run `make build-platforms` for each new CPU package
and for changes that can affect platform builds.

New CPU packages are coherent merge units. Their instruction declarations
refer to execution handlers, and their tests share package helpers. The file
tables give a review order inside each candidate; they are not permission to
merge an incomplete package. Do not extract an early historical version and
leave its known interrupt or memory fixes for a later release.

**Add on target** identifies candidate work not supplied as a complete source
patch. It is a requirement, not an existing test or a passing result.

## Order and dependencies

Dependencies below define the extraction order. P01 is also required before using
the common 60-second test command. The table order is the recommended queue;
independent system packages do not need to wait for unrelated CPU changes.

| Part | Candidate | Required parts | Main review concern |
| --- | --- | --- | --- |
| P00 | Baseline and already-present changes | None | Diverged refs and complete inventory |
| P01 | Configurable short-test deadline | P00 | Preserve target build targets |
| P02 | 6502 option API | P00 | Public `Options` removal |
| P03 | 6502 memory and state contracts | P02 | Nil memory, indexed reads, reset state |
| P04 | 6502 execution and timing fixes | P03 | Interrupt locks, branches, cycle totals |
| P05 | 6502 metadata and variant registries | P04 | Stable IDs and exact encoding sizes |
| P06 | 6502 unofficial-handler extraction | P04 | Preserve the target's bus-cycle write path |
| P07 | Pinned 6502 release gate | P01, P05, P06 | Required data and qualification scope |
| P08 | CHIP-8 metadata organization | P00 | Stable instruction identity |
| P09 | CHIP-8 behavior, state API, and ROM tests | P08 | Default interpreter behavior changes |
| P10 | x86 metadata organization | P00 | Preserve all public entries |
| P11 | Z80 option API and mechanical cleanup | P00 | Public API migration; exclude behavior hunks |
| P12 | Z80 interrupt acceptance | P11 | Shared entry path and instruction boundaries |
| P13 | Z80 ports and unofficial port metadata | P12 | Full port transactions and handler API |
| P14 | Z80 registry, prefix flags, and final evidence | P13 | Complete corpus and exerciser runs |
| P15 | Atari 2600 foundations | P00 | Cartridge mapping limits |
| P16 | CoCo foundations | P00 | SAM definitions |
| P17 | Vectrex foundations | P00 | Register and memory-map definitions |
| P18 | Complete 6809 package | P01 | Indexed modes, interrupts, reverse lookup |
| P19 | Complete SM83 package and integration line | P01 | Interrupt and HALT fetch behavior |
| P20 | Complete 65C816 package and integration line | P01 | Widths, banks, wait states, corpus patch |
| P21 | Complete 68000 package and recorded limits | P01 | Faults, timing, disputed reference vectors |
| P22 | 68000 integration gate | P21 | Blocked until reference discrepancies are resolved |
| P23 | C64 plan, product docs, and final audit | Included code parts | No missing or unassigned hunks |
| P24 | CP/M, Game Boy Color, and Genesis identifiers | P00 | Identifier naming and ordering |
| P25 | Struct-literal formatting in audio, config, and GUI | P00 | No behavior change |

A brace list names every file, for example `arch/cpu/x86/{instruction,instructions}.go`.
Rows without a partial-hunk instruction own the complete source delta for the
listed paths. Added package files are taken together at their fixed-source
versions, subject to target adaptations and the explicit checks below.

## Phase A: 6502 and common test controls

### P00 — Establish the target baseline

| File | Action |
| --- | --- |
| `arch/system/nes/nes.go` | Already present: `RAMEndAddress = 0x07FF`. No extraction. |
| `arch/system/nes/nes_test.go` | Already present: matching RAM-end assertion. No extraction. |
| `docs/work-branch-changelog.md` | Keep this plan and its progress record on the source branch. Do not publish fixed branch statistics as product documentation. |
| `Makefile`, `go.mod`, `go.sum` | Inspect the target commands and dependency versions. Dependency files need no source patch. |

Confirm all 215 original paths have an owner below. Record the target baseline
checks once before code extraction. Distinguish a baseline failure from a
candidate regression. Existing identifiers in `arch/arch.go` and `arch/system.go`
need no registration patch for P15–P21.

**Exit:** Known target SHA, preserved worktree, baseline results, and two NES
paths marked already present. No synchronization merge is required.

### P01 — Make short-test deadlines configurable

| File | Hunks |
| --- | --- |
| `Makefile` | Add `TEST_TIMEOUT ?= 60s`; use it in `test` and `test-coverage`. Leave release and integration additions for their CPU parts. |

Keep race detection in `test`. Preserve `build-platforms`, linter versions,
and all existing targets. The source does not add a separate `test-race` target.

**Focused check:** Inspect `make -n test test-coverage` and
`make -n test TEST_TIMEOUT=120s`; run the common gates.
**Exit:** Default and overridden deadlines reach the correct Go commands.

### P02 — Encapsulate 6502 options

| File | Hunks |
| --- | --- |
| `arch/cpu/cpu6502/option.go` | Private `options`, exported `PreExecutionHook`, typed `Option` results, and nil-safe `newOptions`. Preserve all existing variant, cycle-hook, and stack-hook options. |
| `arch/cpu/cpu6502/cpu.go` | Only `CPU.opts` type and `NewOptions` to `newOptions` constructor call. Defer other constructor and state changes. |
| `arch/cpu/cpu6502/option_test.go` | Complete new option tests. |

`Options` and `NewOptions` disappear. Callers use `New(memory, With...())` and
cannot construct custom options through the old exported implementation type.
Do not add compatibility wrappers that the source deliberately removes.
Document this source API change in the PR.

**Focused check:** `./arch/cpu/cpu6502`.
**Exit:** All option constructors compile; nil options work; existing bus and
stack hooks still work. No variant values are renumbered.

### P03 — Correct 6502 memory and state boundaries

| File | Hunks |
| --- | --- |
| `arch/cpu/cpu6502/errors.go` | Add `ErrNilMemory`. |
| `arch/cpu/cpu6502/memory.go` | Typed-nil validation, required address parameters, bounded immediate integers, indexed address calculation, and related helper/error changes. |
| `arch/cpu/cpu6502/memory_test.go` | Complete source delta: missing arguments, indexed reads, nil and typed-nil memory. |
| `arch/cpu/cpu6502/cpu.go` | Nil-memory constructor panic, shared validation error, complete flag snapshot, reset of `branchTaken` and `TraceStep`, initial constant placement, and stack-wrap comment. Leave branch/lock changes for P04. |
| `arch/cpu/cpu6502/cpu_test.go` | `TestNewRejectsNilMemory`, `TestStateIncludesUnusedFlag`, and `TestResetClearsTransientExecutionState`. Defer `TestBranchValidatesTarget`. |
| `arch/cpu/cpu6502/emulation_test.go` | `TestValidateState` error-identity assertion only. |
| `arch/cpu/cpu6502/doc.go` | Correct the `NewMemory(bus)` example, reset-vector setup order, and memory ranges. Threading/accuracy text belongs to P04/P07. |

Keep constructor behavior explicit: `NewMemory` returns an error; `New(nil)`
panics. `AbsoluteX` and `AbsoluteY` must offset the address before the read.
Do not retain the target's addition of the offset to the returned data byte.

**Focused check:** `./arch/cpu/cpu6502`.
**Exit:** Nil inputs, missing operands, flags, reset state, and indexed reads
pass their tests without importing timing or metadata changes.

### P04 — Correct 6502 stepping, interrupt requests, and timing

| File | Hunks |
| --- | --- |
| `arch/cpu/cpu6502/cpu.go` | Error-returning `branch`, its argument checks, `consumeStallCycle`, and execution/threading comments. Preserve existing cycle and stack event fields. |
| `arch/cpu/cpu6502/interrupt.go` | Lock queued `TriggerIrq` and `TriggerNMI` writes. |
| `arch/cpu/cpu6502/emulation.go` | Update all branch callers to return errors; lock RTI handler-state updates. Leave unofficial handler relocation for P06. |
| `arch/cpu/cpu6502/emulation_65c02.go` | Update BRA and BBR/BBS callers to pass target parameters to `branch` and return its errors. |
| `arch/cpu/cpu6502/step.go` | Stall helper, page-cross timing independent of tracing, and branch target/page-cross handling, including self-loops and wraparound. |
| `arch/cpu/cpu6502/opcode.go` | Timing flags for NMOS LSR/INC absolute-X and LAX absolute-Y; timing-field comments. Leave memory-effect methods for P05. |
| `arch/cpu/cpu6502/opcode_65c02.go` | Absolute-X read-modify-write timings and BRA base timing. Leave NOP metadata size for P05. |
| `arch/cpu/cpu6502/cpu_test.go` | Add `TestBranchValidatesTarget`. |
| `arch/cpu/cpu6502/step_test.go` | All five tests and `newProgramCPU`; P02 supplies its `newOptions` call. |
| `arch/cpu/cpu6502/emulation_test.go` | Replace the EOR and ORA placeholder checks with value/flag assertions. |
| `arch/cpu/cpu6502/doc.go` | Execution/threading claims; **adapt on target:** spell the existing IRQ API `TriggerIrq`, not the source comment's `TriggerIRQ`. |

Review both execution paths. The existing `cycle.go`, `cycle_test.go`,
`stack_event.go`, and `stack_event_test.go` are baseline files, not additions.
Do not remove their interrupt polling or cycle-hook behavior while extracting
instruction-path fixes. Public register access still needs caller serialization.

**Focused check:** `./arch/cpu/cpu6502`, including existing cycle and stack-event
tests. **Add on target:** a bounded concurrent-request regression if existing
coverage does not exercise the changed interrupt/stall locking.
**Exit:** Tracing does not change cycles; NMOS/65C02 timing and self-branches
pass; interrupt requests do not race with the tested execution path.

### P05 — Complete 6502 metadata and variant registries

| File | Hunks |
| --- | --- |
| `arch/cpu/cpu6502/instruction.go` | Move names and registry data; correct BRK and indirect JMP sizes; mark KIL unofficial. |
| `arch/cpu/cpu6502/instruction_name.go` | Moved mnemonic constants. |
| `arch/cpu/cpu6502/instruction_registry.go` | `Instructions`, `InstructionsByID`, complete bit-instruction entries, and `InstructionsForVariant`. |
| `arch/cpu/cpu6502/opcode_id.go` | Append BBR/BBS/RMB/SMB IDs and complete forward/reverse name tables. Preserve the original numeric range. |
| `arch/cpu/cpu6502/categories.go` | Branch and memory-effect set corrections. |
| `arch/cpu/cpu6502/opcode.go` | Nil-safe memory-effect methods and exclusion of accumulator/immediate/implied/relative operands. P04 already owns timing changes. |
| `arch/cpu/cpu6502/instruction_65c02.go` | Indirect JMP size. |
| `arch/cpu/cpu6502/opcode_65c02.go` | One-byte implied NOP metadata. |
| `arch/cpu/cpu6502/unofficial.go` | Exact addressing sizes for every changed unofficial instruction definition. |
| `arch/cpu/cpu6502/instruction_registry_test.go` | Variant registry, stable identity lookup, and addressing-size tests. |
| `arch/cpu/cpu6502/opcode_test.go` | Memory-effect cases, including accumulator/immediate exclusions and bit instructions. |

The BRK encoding is one byte; execution still skips its padding byte. KIL's
unofficial flag is required by P07's legal-opcode count. Variant registries
return independent maps with canonical instruction pointers. Do not renumber
existing public `OpcodeID` values.

**Focused check:** `./arch/cpu/cpu6502`, including `TestVerifyOpcodes`.
**Exit:** Forward/reverse mappings in the changed metadata agree; NMOS, WDC, and Synertek
registries have their intended entries and exact sizes.

### P06 — Move unofficial handlers without losing bus-cycle writes

| File | Hunks |
| --- | --- |
| `arch/cpu/cpu6502/emulation.go` | Remove only the moved unofficial handlers after adding their new owner. Keep P04's branch/RTI changes. |
| `arch/cpu/cpu6502/emulation_unofficial.go` | Add the moved handlers, simplified comments, and `axs` flag helper. Use the corrected `shWrite` from the review corrections. |
| `arch/cpu/cpu6502/cycle_test.go` | `TestBusCycleSHFamilyWrites` from the review corrections. |

The committed fixed source is not a behavior-preserving copy for `shWrite`:
its new file always calls `memory.Write`, while the target uses `writeCycle`
when `cycleActive` is true. The unchanged `cycle.go:memoryCycles` calls the
SHA, SHX, SHY, and TAS handlers directly, so the committed hunk drops the
write callback and one cycle in bus-cycle mode; five NMOS single-step files
(`93`, `9B`, `9C`, `9E`, `9F`) fail with it. The uncommitted review correction
restores the `cycleActive` dispatch and adds the regression test. Extract the
corrected file, not the committed hunk.

**Focused check:** `./arch/cpu/cpu6502` with existing and added bus-cycle tests.
**Exit:** No duplicate handlers; all moved behavior is retained; the SH-family
bus-cycle test passes without downloaded data.

### P07 — Add pinned release qualification

| File | Hunks |
| --- | --- |
| `arch/cpu/cpu6502/release_test.go` | Required pinned checkout validation, legal-opcode vector gate, Dormann gate, and malformed-data unit tests. |
| `arch/cpu/cpu6502/dormann_test.go` | Required-data mode, instruction-count limit/names, real cycle reporting, and updated callers. |
| `Makefile` | Add only `test-6502-release`; preserve the P01 timeout and existing integration lines. |
| `docs/cpu6502-release-qualification.md` | Use the reviewed working-tree procedure, data revisions, and default-path coverage limits. |
| `arch/cpu/cpu6502/doc.go` | Retain instruction-level accuracy limits; do not imply that the release runner qualifies `WithCycleHook`. |

The release tests use helpers from the already-present `singlestep_test.go`.
Do not create duplicate sparse-memory or vector types. `testdata/Makefile`
needs no 6502 change: its moving developer datasets are separate from the
pinned gate.

**Focused checks:** `./arch/cpu/cpu6502`, then the documented
`CPU6502_TESTDATA=... CPU6502_DORMANN_TESTDATA=... make test-6502-release`.
**Exit:** Execute all 1,510,000 legal NMOS vectors and both functional binaries.
Missing, dirty, wrong-revision, or empty required data must fail. Record the
candidate SHA and logs; this plan does not claim that the gate has passed.

## Phase B: Existing CHIP-8, x86, and Z80 packages

### P08 — Separate CHIP-8 metadata owners

| File | Hunks |
| --- | --- |
| `arch/cpu/chip8/instruction.go` | Remove moved names and the `Instructions` map. Preserve instruction definitions. |
| `arch/cpu/chip8/instruction_name.go` | Add the moved names. |
| `arch/cpu/chip8/instruction_registry.go` | Add the moved canonical registry. |
| `arch/cpu/chip8/categories.go` | Use name constants and remove the duplicate LD write-set entry. |

**Focused check:** `./arch/cpu/chip8`.
**Exit:** Existing names, instruction pointers, and opcode lookup remain valid.
Do not include state or interpreter behavior changes in this candidate.

### P09 — Correct CHIP-8 interpreter behavior and add ROM evidence

| File | Hunks |
| --- | --- |
| `arch/cpu/chip8/option.go` | `Quirks`, functional options, and nil-option handling. |
| `arch/cpu/chip8/cpu.go` | `State`/`State()` rename, option-aware constructor, key-wait state, display cadence, reset, and snapshot handling. |
| `arch/cpu/chip8/emulation.go` | Key release, VF operand ordering, shifts, logic flags, load/store I updates, jump quirk, sprite wait/wrap, and memory bounds. |
| `arch/cpu/chip8/emulation_test.go` | All adjusted and added behavior, state-round-trip, quirk, flag-alias, bounds, and display-wait cases. |
| `arch/cpu/chip8/rom_test.go` | Complete six-ROM runner, screenshot comparison, keypad input, and shared helpers. |
| `arch/cpu/chip8/doc.go` | VIP default, explicit quirks, 60 Hz `UpdateTimers`, and supported bounds behavior. |
| `testdata/Makefile` | All source hunks: CHIP-8 variables, target, help, `.PHONY`, `all`, and `clean` entries. |

Keep the quirk options and default behavior in one candidate. Publish the
`CPUState`/`GetState` to `State`/`State()` migration. A plain `New()` remains a
valid call but now uses VIP behavior. Hosts must drive the 60 Hz timer tick.

**Focused checks:** `./arch/cpu/chip8`; provision `make -C testdata chip8`,
record its revision, then run
`go test -v -count=1 -timeout 120s ./arch/cpu/chip8 -run '^TestROMConformance$'`.
**Exit:** All six ROM subtests execute and pass. Snapshot and key/display waits
pass the unit tests. Do not add this ROM runner to the root integration target
as though such a source hunk existed; the source has no such addition.

### P10 — Separate x86 metadata owners

| File | Hunks |
| --- | --- |
| `arch/cpu/x86/instruction.go` | Remove moved instruction-name constants only. |
| `arch/cpu/x86/instruction_name.go` | Add those constants, including segment prefixes. |
| `arch/cpu/x86/instructions.go` | Remove the moved `Instructions` map only. |
| `arch/cpu/x86/instruction_registry.go` | Add the unchanged registry and representative instruction entries. |

**Focused check:** `./arch/cpu/x86`.
**Exit:** Public symbols and mappings remain unchanged. This is metadata
organization; it does not add an x86 emulator.

### P11 — Encapsulate Z80 options and apply mechanical cleanup

| File | Hunks |
| --- | --- |
| `arch/cpu/z80/option.go` | Private options, exported hook type, typed option results, nil-safe `newOptions`. |
| `arch/cpu/z80/option_test.go` | All option tests and their fixture. |
| `arch/cpu/z80/cpu.go` | Option storage, constructor/newCPU signatures, receiver renames, and snapshot/callback contract comments. Defer typed IM, `eiPending`, and port-helper replacement. |
| `arch/cpu/z80/cpu_test.go` | `TestNew` uses `newOptions`; defer bus-fixture and legacy-port additions. |
| `arch/cpu/z80/{emulation_helpers,emulation_jump,flag,param}.go` | Entire source deltas: receiver renames only. |
| `arch/cpu/z80/{emulation_dd_undoc,emulation_fd_undoc}.go` | Receiver renames and unused self-load parameters only. |
| `arch/cpu/z80/{emulation,emulation_ed,step,interrupt}.go` | Receiver renames only. Leave semantic changes for P12–P14. |
| `arch/cpu/z80/memory.go` | Method ordering and adapter receiver renames only. Defer new behavior-specific comments until P12/P13. |

As with P02, removal of exported `Options`/`NewOptions` needs a migration note.
Do not import `eiPending` behavior, typed interrupt mode assignments, or direct
bus-call changes merely because receiver changes share their diff hunks.

**Focused check:** `./arch/cpu/z80`.
**Exit:** Options and existing constructors work; behavior is unchanged.

### P12 — Unify Z80 interrupt acceptance

| File | Hunks |
| --- | --- |
| `arch/cpu/z80/cpu.go` | Typed `im`, exported-state conversion back to `uint8`, and `eiPending` field. |
| `arch/cpu/z80/interrupt.go` | One locked acceptance path for `CheckInterrupts` and `Step`; vector sampling before stack writes; NMI/IFF2, refresh, HALT, and parity handling. |
| `arch/cpu/z80/step.go` | Early return for accepted interrupts, HALT refresh, EI delay lifecycle, and removal of the duplicate interrupt implementation. Leave Q-prefix changes for P14. |
| `arch/cpu/z80/emulation.go` | DI cancels and EI sets `eiPending`. |
| `arch/cpu/z80/emulation_ed.go` | Use typed interrupt-mode constants in IM handlers only. |
| `arch/cpu/z80/singlestep_test.go` | Convert IM when setting/comparing state. Leave port runner changes for P13. |
| `arch/cpu/z80/interrupt_test.go` | Complete new interrupt regression suite and `stackMappedIRQBus`. |
| `arch/cpu/z80/cpu_test.go` | Extend `testBus` with IRQ, RETI, and port recording fields/methods used by the new tests. The existing `singleStepPort` type already lives in `singlestep_test.go`. |
| `arch/cpu/z80/ed_mirror_test.go` | `TestEDInterruptModeMirrors` and `TestEDReturnMirrorsAndNotification`; omit P13/P14 tests until their changes arrive. |
| `arch/cpu/z80/memory.go` | IRQData/OnRETI callback contracts and CPU-lock warning. Leave block-port ordering text for P13. |

Existing RETI notification and device-vector support are baseline behavior.
The new work makes both interrupt entry points agree and tests that behavior.
Retain the non-RST IM 0 fallback; do not claim arbitrary IM 0 execution.

**Focused check:** `./arch/cpu/z80`.
**Exit:** Accepted interrupts consume their own step; EI/DI/HALT, nested NMI,
IM 2 wraparound, refresh bits, and LD A,I/R parity tests pass for both entry paths.

### P13 — Correct Z80 port transfers and INF/OUTF metadata

| File | Hunks |
| --- | --- |
| `arch/cpu/z80/emulation_ed.go` | Decrement B before OUTI/OUTD port writes; replace remaining port wrappers with bus calls. Include repeating instructions through their shared handlers. |
| `arch/cpu/z80/{cpu,emulation}.go` | Replace port wrapper calls with direct `Bus` calls. |
| `arch/cpu/z80/unofficial.go` | Remove port wrappers and obsolete `inf`/`outf` block handlers; make INF/OUTF use ED 70/71 metadata and `ParamFunc`. |
| `arch/cpu/z80/instruction_ed.go` | Add addressing metadata for `EdInFC` and `EdOut0C`. |
| `arch/cpu/z80/unofficial_test.go` | `TestUndocumentedPortInstructionDefinitions` only. |
| `arch/cpu/z80/ed_mirror_test.go` | Add `TestFullBusPortAddresses`. |
| `arch/cpu/z80/cpu_test.go` | Add `TestLegacyBusCompatibility` and `legacyTestIO`; retain the P12 bus fixture. |
| `arch/cpu/z80/singlestep_test.go` | `singleStepBus`, ordered full-address port comparison, and `NewWithBus` construction. Preserve P12's IM conversions. |
| `arch/cpu/z80/memory.go` | Complete input/output address and block-ordering comments. |

Publish the INF/OUTF direct-handler migration from `NoParamFunc` to `ParamFunc`.
Remove private `readPort`/`writePort` only after all callers are converted.
Keep the eight-bit legacy adapter behavior unchanged.

**Focused check:** `./arch/cpu/z80`.
**Exit:** Address, data, direction, and ordering checks pass for every port
form. Full-corpus success is checked after P14's prefix correction; do not
weaken the runner to hide that outstanding failure.

### P14 — Finish Z80 registry, prefix flags, and validation record

| File | Hunks |
| --- | --- |
| `arch/cpu/z80/instruction.go` | Move name constants and fix the changed instruction comments. |
| `arch/cpu/z80/instruction_name.go` | Add the moved constants. |
| `arch/cpu/z80/instruction_registry.go` | Build representative lookup entries from base/prefix tables; include DDCB/FDCB and INF/OUTF. |
| `arch/cpu/z80/instruction_registry_test.go` | Complete registry test. |
| `arch/cpu/z80/categories.go` | Correct Call/Halt comments. |
| `arch/cpu/z80/opcode_test.go` | Correct the RST test label. |
| `arch/cpu/z80/step.go` | Clear Q on ignored DD/FD prefixes. |
| `arch/cpu/z80/ed_mirror_test.go` | Add `TestEDNegMirrors` and `TestIgnoredPrefixResetsQ`; all tests in this file are now included. |
| `arch/cpu/z80/doc.go` | Complete rewritten package contracts and accuracy limits. |
| `docs/z80-gap-closure-plan.md` | Use the reviewed procedure. Retain historical results separately from candidate results. |

**Focused checks:** `./arch/cpu/z80`, then the documented pinned single-step
command and `TestZexdoc`/`TestZexall` commands from the Z80 plan. The root
integration target already invokes these tests; no Z80 Makefile line is added.
**Exit:** Execute the 1,604,000 recorded vectors and all 67 groups in each
exerciser. Record the ZEX repository revision, which the old record omitted.
State checks do not establish T-state or memory-bus-cycle accuracy.

## Phase C: System foundations

These packages define hardware constants and helpers. They do not implement
complete machines. Their tests import shared utilities, not the new CPU
packages, so they do not depend on P18–P21.

### P15 — Add Atari 2600 definitions and cartridge helpers

| File | Complete contents to include |
| --- | --- |
| `arch/system/atari2600/atari2600.go` | Address ranges, reset vector, RAM/mirror constants, cartridge sizes, and 13-bit address mask. |
| `arch/system/atari2600/atari2600_test.go` | Memory ranges, mirrors, vectors, sizes, and address masking. |
| `arch/system/atari2600/doc.go` | System foundation scope and CPU relationship. |
| `arch/system/atari2600/register/tia.go` | TIA read/write register addresses and name maps. |
| `arch/system/atari2600/register/riot.go` | RIOT registers and console/joystick bit definitions. |
| `arch/system/atari2600/register/register_test.go` | Register addresses, completeness, and input bits. |
| `arch/system/atari2600/cartridge/cartridge.go` | BankingScheme, image loading, size detection, bank offsets/counts, and hotspot lookup. |
| `arch/system/atari2600/cartridge/cartridge_test.go` | Valid/invalid images, all supported bank offsets and hotspots, 3F 2 KB banks, and address mirrors. |

Use the final source implementation, including 3F fixes. Size-based scheme
detection is a heuristic. `TriggerBank` identifies hotspots; it does not run a
complete bank-switching bus. The host must apply write values and fixed-bank
behavior. Do not present these helpers as a complete cartridge emulator.

**Focused packages:** `./arch/system/atari2600`,
`./arch/system/atari2600/register`, `./arch/system/atari2600/cartridge`.
**Exit:** All range and cartridge tests pass; package docs state the limits.

### P16 — Add CoCo definitions

| File | Complete contents to include |
| --- | --- |
| `arch/system/coco/coco.go` | Memory, I/O, interrupt vectors, and cartridge constants. |
| `arch/system/coco/coco_test.go` | Address ranges, vectors, sizes, and overlap checks. |
| `arch/system/coco/doc.go` | Package scope. |
| `arch/system/coco/register/pia.go` | PIA register addresses. |
| `arch/system/coco/register/sam.go` | SAM set/clear addresses, CPU rate, RAM size, and memory-map definitions. |
| `arch/system/coco/register/register_test.go` | PIA completeness and corrected SAM control addresses. |

The target has no earlier CoCo package. Add the corrected `SAMR0*`, `SAMR1*`,
and `SAMTY*` names directly. There is no main-branch migration from the old
branch-only `SAMRate*` names. CPU rate, RAM size, and memory mapping are distinct.

**Focused packages:** `./arch/system/coco`, `./arch/system/coco/register`.
**Exit:** Definitions and tests pass without introducing a CPU dependency.

### P17 — Add Vectrex definitions

| File | Complete contents to include |
| --- | --- |
| `arch/system/vectrex/vectrex.go` | RAM, ROM, cartridge, VIA, and interrupt-vector constants. |
| `arch/system/vectrex/vectrex_test.go` | Ranges, vectors, cartridge sizes, mirrors, and overlap checks. |
| `arch/system/vectrex/doc.go` | Package scope. |
| `arch/system/vectrex/register/via.go` | VIA registers, port B signal bits, PSG button bits, and interrupt bits. |
| `arch/system/vectrex/register/register_test.go` | Register addresses, completeness, buttons, and IRQ bits. |

**Focused packages:** `./arch/system/vectrex`, `./arch/system/vectrex/register`.
**Exit:** Definition tests pass. No claim of vector display, audio, or VIA execution.

## Phase D: New CPU packages

All paths in each package table are new relative to the verified target.
Include every listed file in that package's candidate. Shared dependencies
are already on the target: `assert`, `set`, and, for 68000 options, `arch`.
The new CPU packages do not import one another.

### P18 — Add the complete Motorola 6809 package

| File | Complete contents to include |
| --- | --- |
| `arch/cpu/cpu6809/addressing.go` | Addressing-mode contract. |
| `arch/cpu/cpu6809/categories.go` | Branch and memory-effect sets. |
| `arch/cpu/cpu6809/cpu.go` | Registers, State, constructor/reset, D register, stacks, and trace state. |
| `arch/cpu/cpu6809/cpu_test.go` | CPU lifecycle, register/stack tests, nil options, tracing, and shared `testMem`/`newTestCPU`. |
| `arch/cpu/cpu6809/doc.go` | Architecture and usage documentation. |
| `arch/cpu/cpu6809/emulation.go` | Arithmetic, logic, shifts, loads, stores, and operand helpers. |
| `arch/cpu/cpu6809/emulation_branch.go` | Conditional branches, control transfers, and register transfers. |
| `arch/cpu/cpu6809/emulation_stack.go` | PSH/PUL masks, transfer order, and byte counts. |
| `arch/cpu/cpu6809/emulation_system.go` | RTS/RTI, software interrupts, CWAI, and SYNC. |
| `arch/cpu/cpu6809/emulation_test.go` | Instruction execution and flag tests. |
| `arch/cpu/cpu6809/errors.go` | Package error contracts. |
| `arch/cpu/cpu6809/flag.go` | Condition-code packing and arithmetic helpers. |
| `arch/cpu/cpu6809/instruction.go` | Instruction definitions and addressing metadata types. |
| `arch/cpu/cpu6809/instruction_name.go` | Mnemonic constants. |
| `arch/cpu/cpu6809/instruction_registry.go` | Canonical mnemonic lookup registry. |
| `arch/cpu/cpu6809/interrupt.go` | Queued NMI/IRQ/FIRQ, masks, wait modes, and entry cycles. |
| `arch/cpu/cpu6809/interrupt_test.go` | Pending/masked requests, SYNC wake, CWAI single stacking, and 19/19/10-cycle entry. |
| `arch/cpu/cpu6809/memory.go` | BasicMemory, wrapper, and endian helpers. |
| `arch/cpu/cpu6809/opcode.go` | Base opcode table, metadata methods, and lookup APIs. |
| `arch/cpu/cpu6809/opcode_10.go` | Page `$10` table. |
| `arch/cpu/cpu6809/opcode_11.go` | Page `$11` table. |
| `arch/cpu/cpu6809/opcode_id.go` | Stable identities and name mappings. |
| `arch/cpu/cpu6809/opcode_test.go` | Base and prefixed consistency, reverse mappings, IDs, timings, and categories. |
| `arch/cpu/cpu6809/option.go` | Tracing and typed execution hook options. |
| `arch/cpu/cpu6809/param.go` | Indexed postbytes, effective addresses, typed operands, and indexed cycle costs. |
| `arch/cpu/cpu6809/param_test.go` | Invalid indirect updates, invalid operand types, and indexed timing. |
| `arch/cpu/cpu6809/step.go` | Base/prefix decode, interrupt dispatch, operand decode, and execution. |

Review metadata and reverse mappings first, then operands and execution, then
interrupt and wait-state behavior. Keep the final interrupt cycle fix with the
first package addition. The source has no external 6809 corpus runner and no
root integration line for it. Do not invent a passing external qualification.

**Focused package:** `./arch/cpu/cpu6809`; run common gates and cross-builds.
**Exit:** All 27 files build together; all three opcode tables have tested
reverse lookup; masked requests and CWAI/SYNC behavior pass.

### P19 — Add the complete Sharp SM83 package

| File | Complete contents to include |
| --- | --- |
| `arch/cpu/sm83/addressing.go` | Addressing modes and typed operand/register selectors. |
| `arch/cpu/sm83/categories.go` | Instruction categories. |
| `arch/cpu/sm83/cpu.go` | Registers, constructor, state, memory access, and register/stack helpers. |
| `arch/cpu/sm83/doc.go` | Package behavior and usage. |
| `arch/cpu/sm83/emulation.go` | Arithmetic, flags, HALT/STOP, and EI/DI. |
| `arch/cpu/sm83/emulation_branch.go` | Jumps, calls, returns, conditions, and RST. |
| `arch/cpu/sm83/emulation_cb.go` | CB shifts, rotations, swaps, and bit operations. |
| `arch/cpu/sm83/emulation_load.go` | Register/memory loads, high-memory loads, SP operations, and stack transfers. |
| `arch/cpu/sm83/errors.go` | Package errors. |
| `arch/cpu/sm83/flag.go` | SM83 flag packing and setters. |
| `arch/cpu/sm83/instruction.go` | Instruction definitions and opcode metadata contracts. |
| `arch/cpu/sm83/instruction_name.go` | Mnemonic constants. |
| `arch/cpu/sm83/instruction_registry.go` | Registry from base/CB tables, including LDH through C. |
| `arch/cpu/sm83/instruction_registry_test.go` | Registry completeness. |
| `arch/cpu/sm83/interrupt.go` | Enabled/pending interrupt selection and service. |
| `arch/cpu/sm83/memory.go` | Memory contract and flat memory implementation. |
| `arch/cpu/sm83/opcode.go` | Base table. |
| `arch/cpu/sm83/opcode_cb.go` | CB table. |
| `arch/cpu/sm83/opcode_id.go` | IDs and name mappings. |
| `arch/cpu/sm83/option.go` | Typed options, tracing, hook, initial PC/SP. |
| `arch/cpu/sm83/option_test.go` | Option checks. |
| `arch/cpu/sm83/param.go` | Operand decoding with HALT-bug fetch behavior and decoded immediate reuse. |
| `arch/cpu/sm83/singlestep_test.go` | External state-vector runner and comparison helpers. |
| `arch/cpu/sm83/step.go` | Interrupt boundary, delayed IME, HALT wake, and base/CB decode. |
| `arch/cpu/sm83/step_test.go` | Interrupt dispatch, EI/DI delay, HALT-bug byte/word/CB/branch cases, and EI/HALT return address. |
| `Makefile` | Add only the SM83 command in `test-integration`. |

Keep the five-machine-cycle interrupt step separate from the first handler
instruction. Include the final HALT-bug and immediate-read fixes in the first
candidate. The source has no standalone `opcode_test.go` or general
`emulation_test.go` for SM83. **Add on target:**
`arch/cpu/sm83/opcode_test.go` with bidirectional mapping checks. The existing
registry test checks names, not each opcode encoding.

**Focused package:** `./arch/cpu/sm83`; run common gates and cross-builds.
Provision `make -C testdata sm83`, record the revision, then run
`go test -v -race -count=1 -timeout 10m ./arch/cpu/sm83 -run '^TestSingleStep$'`.
The data path is `testdata/sm83/v1`; no environment override exists.
**Exit:** Unit and external cases execute and pass. Do not add a complete
Game Boy system claim.

### P20 — Add the complete WDC 65C816 package

| File | Complete contents to include |
| --- | --- |
| `arch/cpu/cpu65816/addressing.go` | Addressing modes and accumulator operand type. |
| `arch/cpu/cpu65816/categories.go` | Instruction categories. |
| `arch/cpu/cpu65816/cpu.go` | Registers, native/emulation mode, widths, constructor/reset, state, and tracing. |
| `arch/cpu/cpu65816/cpu_test.go` | Lifecycle, A/B register widths, stacks, PC banks, and shared test-memory helpers. |
| `arch/cpu/cpu65816/doc.go` | Architecture and usage documentation. |
| `arch/cpu/cpu65816/effective_address.go` | Banked effective-address resolution. |
| `arch/cpu/cpu65816/effective_address_test.go` | Address modes and invalid operand checks. |
| `arch/cpu/cpu65816/emulation.go` | Width-aware arithmetic, logic, shifts, and operand helpers. |
| `arch/cpu/cpu65816/emulation_branch.go` | Short/long branches, jumps, and calls. |
| `arch/cpu/cpu65816/emulation_move.go` | Loads, stores, transfers, and block moves. |
| `arch/cpu/cpu65816/emulation_stack.go` | Stack instructions and effective-address pushes. |
| `arch/cpu/cpu65816/emulation_system.go` | Status/mode changes, wait/stop, software interrupts, and returns. |
| `arch/cpu/cpu65816/emulation_test.go` | Instruction, mode, width, and flag tests. |
| `arch/cpu/cpu65816/errors.go` | Package errors. |
| `arch/cpu/cpu65816/flag.go` | Status packing and width-specific flag helpers. |
| `arch/cpu/cpu65816/instruction.go` | Instruction definitions and opcode metadata contracts. |
| `arch/cpu/cpu65816/instruction_name.go` | Mnemonics. |
| `arch/cpu/cpu65816/instruction_registry.go` | Canonical mnemonic lookup registry. |
| `arch/cpu/cpu65816/interrupt.go` | Queued interrupt acceptance and native/emulation stack entry. |
| `arch/cpu/cpu65816/interrupt_test.go` | Step dispatch, masked-IRQ WAI wake, and STP reset requirement. |
| `arch/cpu/cpu65816/memory.go` | Memory interface, wrapper, and bank helpers. |
| `arch/cpu/cpu65816/memory_access.go` | Program/data-bank access and wrapping rules. |
| `arch/cpu/cpu65816/opcode.go` | Complete table, width flags, timing, and metadata methods. |
| `arch/cpu/cpu65816/opcode_id.go` | Identities and name mappings. |
| `arch/cpu/cpu65816/opcode_test.go` | Table completeness, lookup, timing, consistency, and width flags. |
| `arch/cpu/cpu65816/option.go` | Typed tracing/hook options. |
| `arch/cpu/cpu65816/option_test.go` | Option checks. |
| `arch/cpu/cpu65816/param.go` | Width-aware operand readers and address decoding. |
| `arch/cpu/cpu65816/singlestep_test.go` | Tagged external runner, sparse memory, and state comparisons. |
| `arch/cpu/cpu65816/stack.go` | Native/emulation stack width and wrap behavior. |
| `arch/cpu/cpu65816/step.go` | Interrupt dispatch, decode, execution, and timing. |
| `Makefile` | Add only the tagged 65C816 command in `test-integration`. |

Do not separate the final WAI/STP and program-bank interrupt fixes from the
new package. The existing `testdata/Makefile` already downloads this corpus
and applies DirtyHairy's `fix_sbc_direct_x` patch to `v1/e1.e.json`. Record both
the corpus HEAD and that patch's fetched commit; a dirty corpus is intentional
here and must be explained. This differs from P07's clean pinned checkouts.

**Focused package:** `./arch/cpu/cpu65816`; run common gates and cross-builds.
Provision `make -C testdata cpu65816`, record the revisions and file patch,
then run `go test -v -race -tags singlestep -count=1 -timeout 10m ./arch/cpu/cpu65816 -run '^TestSingleStep$'`.
`CPU65816_TESTDATA`, when set, names the `v1` directory itself.
**Exit:** All files execute with the recorded patch. Inspect missing/empty-file
skips. **Add on target:** extend `arch/cpu/cpu65816/opcode_test.go` to compare
reverse opcode values. Its current consistency test checks addressing-mode
presence only. The runner stops each file after ten failures; a failed run
does not establish a complete executed-vector count. Do not claim complete bus-cycle accuracy.

### P21 — Add the complete Motorola 68000 package with explicit limits

| File | Complete contents to include |
| --- | --- |
| `arch/cpu/cpu68000/addressing.go` | Addressing modes, operand sizes, and size decode. |
| `arch/cpu/cpu68000/categories.go` | Instruction categories. |
| `arch/cpu/cpu68000/cpu.go` | Registers, supervisor/user stacks, reset vectors, state, halt/recovery, and options. |
| `arch/cpu/cpu68000/cpu_test.go` | Constructor/vector/state/stack-bank/halt/cycle tests and shared CPU helper. |
| `arch/cpu/cpu68000/doc.go` | Package architecture and usage. |
| `arch/cpu/cpu68000/ea.go` | Effective-address decode, extension words, reads/writes, and update order. |
| `arch/cpu/cpu68000/ea_test.go` | Addressing, operand sizes, A7 byte updates, and extension cases. |
| `arch/cpu/cpu68000/emulation.go` | Arithmetic, decimal operations, logic, comparisons, and multiply/divide. |
| `arch/cpu/cpu68000/emulation_bit.go` | BTST/BSET/BCLR/BCHG behavior. |
| `arch/cpu/cpu68000/emulation_branch.go` | Conditions, branches, DBcc/Scc, jumps, and calls. |
| `arch/cpu/cpu68000/emulation_move.go` | MOVE variants, MOVEM/MOVEP, register/address transfers, LINK/UNLK. |
| `arch/cpu/cpu68000/emulation_shift.go` | Register/memory shifts and rotations. |
| `arch/cpu/cpu68000/emulation_system.go` | Traps, CHK, exception returns, STOP/RESET, illegal operations, and TAS. |
| `arch/cpu/cpu68000/emulation_test.go` | Execution regressions, overlapping operands, shift references, and divide-by-zero PC behavior. |
| `arch/cpu/cpu68000/errors.go` | Package and memory-fault errors. |
| `arch/cpu/cpu68000/fault_test.go` | Fault frame, partial transfers, reset recovery, double faults, privilege/fetch ordering, and host-panic boundaries. |
| `arch/cpu/cpu68000/flag.go` | CCR/SR, supervisor state, masks, and sign helpers. |
| `arch/cpu/cpu68000/instruction.go` | Execution contract and instruction definitions. |
| `arch/cpu/cpu68000/instruction_name.go` | Mnemonics. |
| `arch/cpu/cpu68000/instruction_registry.go` | Canonical mnemonic lookup registry. |
| `arch/cpu/cpu68000/interrupt.go` | Queued/bus IRQs, exception entry, vectors, and interrupt cycles. |
| `arch/cpu/cpu68000/interrupt_test.go` | One-step IRQ service, priority/masking, and STOP behavior. |
| `arch/cpu/cpu68000/memory.go` | Memory/Bus interfaces, BasicMemory/BasicBus, big-endian transfers, and 24-bit masking. |
| `arch/cpu/cpu68000/memory_access.go` | Checked CPU transfers, optional `BusErrorHandler`, alignment checks, and fault boundaries. |
| `arch/cpu/cpu68000/memory_test.go` | Byte/word/long access, ROM loading, masks, and address wrapping. |
| `arch/cpu/cpu68000/opcode.go` | Complete instruction-family decoder, including overlapping opcode families. |
| `arch/cpu/cpu68000/opcode_id.go` | IDs and name mappings. |
| `arch/cpu/cpu68000/opcode_test.go` | Decoder-family and overlapping-encoding tests. |
| `arch/cpu/cpu68000/option.go` | Typed system, tracing, and initial-state options. |
| `arch/cpu/cpu68000/option_test.go` | Option checks. |
| `arch/cpu/cpu68000/singlestep_test.go` | Tagged full-corpus runner; cap diagnostics, not executed vectors. |
| `arch/cpu/cpu68000/step.go` | Decode/execute boundary, checked fetches, trace, IRQ entry, and stack-pointer synchronization. |
| `arch/cpu/cpu68000/timing.go` | Effective-address and operand-dependent instruction timing. |
| `arch/cpu/cpu68000/timing_test.go` | Addressing, instruction, and divide timing checks. |
| `docs/cpu68000-gap-closure-plan.md` | Reviewed fault/timing contracts, pinned corpus procedure, and disputed-vector record. |

Use three review passes on one buildable candidate: metadata/decode; operands
and instruction execution; checked transfers, exceptions, and timing. Keep the
14-byte original-68000 fault frame and all later corpus-driven fixes. Do not
split checked access from the handlers that must use it.

**Focused package:** `./arch/cpu/cpu68000`; run common gates and cross-builds.
Run the pinned full-corpus command in the 68000 plan and record every count.
The prior result is 996,321 passes and 3,739 failures, not a passing gate.

**Exit:** All normal required checks pass. Publish the full-corpus result and
explicit scope limits with the candidate. Review the decoder's reverse-mapping
contract; document any architectural limitation instead of inventing a passing
`TestVerifyOpcodes` command. This package can be reviewed as an instruction-level
implementation with unresolved reference discrepancies. If project acceptance
requires zero external mismatches, hold P21 together with P22.

Do not add the 68000 line to the required aggregate integration target in this
part. Keep the tagged runner and every failing comparison intact and directly
runnable. Omitting the aggregate line is not permission to report conformance.

### P22 — Resolve 68000 reference discrepancies before aggregate integration

| File | Action |
| --- | --- |
| `Makefile` | Add the source's tagged 68000 `test-integration` line only after this part's exit condition. |
| `arch/cpu/cpu68000/singlestep_test.go` | Retain all comparisons and vectors. **Add on target only if needed:** justified runner fixes established by independent evidence. |
| `arch/cpu/cpu68000/{emulation_shift,emulation,emulation_test}.go` | **Add on target only if needed:** independently verified ASR, ASL, or DIVU corrections and focused regressions. Do not change correct behavior merely to match disputed vectors. |
| `docs/cpu68000-gap-closure-plan.md` | Record each resolution, reference evidence, corpus revision/patch, and new totals. |

This part has required work beyond the fixed source. The old full run reports
3,736 ASR C/X mismatches, two ASL byte cases that change upper register bits,
and one DIVU-by-zero saved-PC/flag mismatch. Resolve each class through an
independent reference or a traceable corpus correction. Keep the original
inputs and the correction record. Do not skip failures, weaken assertions,
or convert failure into success through a shell wrapper.

**Focused checks:** Relevant 68000 unit regressions and the complete pinned
corpus. Once those pass, run the resulting aggregate integration target with
all its datasets present. The source's aggregate commands use the race
detector and no finite timeout; apply a recorded CI job deadline.
**Exit:** All required external comparisons pass against identified data.
Until then, record P22 as planned with an unresolved acceptance condition.
P23 can record that deferral, but must not call the complete plan finished.

## Phase E: Product documentation and final reconciliation

### P23 — Publish the remaining docs and audit every source hunk

| File | Action |
| --- | --- |
| `docs/system-implementation-plan-c64.md` | Publish the reviewed plan as planned system work. `Variant6510` already exists; C64 system registration and devices do not. |
| `docs/cpu6502-release-qualification.md` | P07 owns the procedure; verify its commands against the merged target. |
| `docs/z80-gap-closure-plan.md` | P14 owns the record; preserve historical results and add candidate evidence separately. |
| `docs/cpu68000-gap-closure-plan.md` | P21/P22 own the record; preserve unresolved limits until evidence closes them. |
| `docs/work-branch-changelog.md` | Update source-branch progress and remaining hunks. Keep branch tracking out of product documentation. |
| `README.md` | Use the uncommitted source delta: feature list and package tree for the new CPU and system packages. List only packages actually merged. |
| `Makefile` | Change the integration-target description to the source's CPU wording; verify all included package lines and P22's status. |

API notes must accompany their code parts, not wait for this final audit.
For the new system packages, distinguish constants/helpers from device
execution. For CPUs, distinguish instruction-state/cycle totals from complete
bus emulation. State snapshots do not imply deterministic save/restore of
pending interrupts, internal latches, and wait states.

Compare the final target with the fixed source and inspect each remaining
hunk. Mark it merged, already present, superseded by a target repair, explicitly
deferred, or branch tracking only. A small diff is not proof of completion;
the P06 repair intentionally differs from the source. The P22 integration
line remains deferred until its acceptance condition is met.

**Checks:** Documentation links and command text, `git diff --check`, and the
file/hunk inventory. Reuse code results when no code or test inputs changed.
**Exit:** No unexplained path or hunk. If P22 remains deferred, identify it as
remaining work instead of reporting all source changes merged.

### P24 — Add CP/M, Game Boy Color, and Genesis system identifiers

| File | Hunks |
| --- | --- |
| `arch/system.go` | `CPM`, `GameBoyColor` (`"gbc"`), and `Genesis` constants; alphabetical placement; corrected `GameBoy` comment. |
| `arch/system_test.go` | Parse, constant, and `SupportedSystems` cases for the three identifiers. |

`"gbc"` is the only abbreviated identifier. Dependent tools already use it as
the profile name, so keep it and its comment. `Generic` sorts before `Genesis`.

**Focused check:** `./arch`.
**Exit:** All three identifiers parse and appear once in `SupportedSystems`.

### P25 — Apply struct-literal formatting without behavior change

| File | Hunks |
| --- | --- |
| `audio/sdl2/sdl.go` | One field per line in multi-line struct literals. |
| `audio/format_test.go`, `audio/sdl2/sdl_integration_test.go`, `audio/sdl2/sdl_test.go`, `config/options_test.go`, `gui/sdl2/sdl_test.go` | Same formatting in test literals. |

These hunks satisfy a newer `retrogolint` rule than the pinned `v1.0.6`. The
pinned linter accepts both layouts. Verify with `gofmt -l` and the package tests.

**Focused check:** `./audio/...`, `./config`, `./gui/sdl2`.
**Exit:** No semantic diff; `git diff -w` shows only layout changes.

## Shared-file ownership

This table is the extraction checklist for files whose complete source diff
must be divided across candidates. Phase tables above name tests and symbols.

| File | Part ownership |
| --- | --- |
| `Makefile` | P01 timeout; P07 release target; P19 SM83 integration; P20 65C816 integration; P22 68000 integration; P23 help text. Preserve baseline platform targets. |
| `README.md` | P23 package listings; no other part changes it. |
| `arch/cpu/cpu6502/cpu.go` | P02 options; P03 construction/state/reset/constants; P04 branch/stall/interrupt contract. |
| `arch/cpu/cpu6502/cpu_test.go` | P03 constructor/state/reset tests; P04 branch validation. |
| `arch/cpu/cpu6502/emulation.go` | P04 branch/RTI fixes; P06 unofficial-handler removal. |
| `arch/cpu/cpu6502/emulation_test.go` | P03 nil-memory error; P04 EOR/ORA assertions. |
| `arch/cpu/cpu6502/opcode.go` | P04 runtime timing; P05 memory-effect methods. |
| `arch/cpu/cpu6502/opcode_65c02.go` | P04 timing; P05 implied NOP size. |
| `arch/cpu/cpu6502/doc.go` | P03 example/map; P04 execution and corrected `TriggerIrq` name; P07 qualification limits. |
| `arch/cpu/z80/cpu.go` | P11 options/receivers/contracts; P12 IM type/state cast/EI latch; P13 direct port call. |
| `arch/cpu/z80/cpu_test.go` | P11 option helper; P12 recording bus fixture; P13 legacy-port test and fixture. |
| `arch/cpu/z80/emulation.go` | P11 receivers; P12 EI/DI; P13 direct bus calls. |
| `arch/cpu/z80/emulation_ed.go` | P11 receivers; P12 IM constants; P13 output ordering and direct bus calls. |
| `arch/cpu/z80/step.go` | P11 receivers; P12 interrupt/HALT/EI lifecycle and old handler removal; P14 ignored-prefix Q. |
| `arch/cpu/z80/interrupt.go` | P11 receivers; P12 typed modes and common acceptance path. |
| `arch/cpu/z80/memory.go` | P11 ordering/receivers; P12 callback contracts; P13 port-order contracts. |
| `arch/cpu/z80/singlestep_test.go` | P12 IM casts; P13 ordered full-address port runner. |
| `arch/cpu/z80/ed_mirror_test.go` | P12 IM/return mirrors; P13 port addresses; P14 NEG mirrors and Q-prefix regression. |
| `docs/cpu68000-gap-closure-plan.md` | P21 existing evidence; P22 new resolution evidence; P23 final status check. |

Keep helper ownership explicit. `newProgramCPU` arrives with P04.
Z80's `testBus` changes arrive with P12 and use the pre-existing
`singleStepPort`. P07 uses the existing 6502 single-step helpers. New CPU
packages carry all their shared test helpers in their first candidate.

## Complete path inventory

The phase tables assign all 223 paths in the original three-dot range.
Counts below describe unique paths; mixed files count once. These statistics
use the exact range and file statuses recorded at the start of this plan.

| Status | Count | Paths | Owner |
| --- | ---: | --- | --- |
| Modified | 1 | `Makefile` | P01/P07/P19/P20/P22/P23 |
| Added / Modified | 4 / 6 | `arch/cpu/chip8/` | P08/P09 |
| Added / Modified | 8 / 20 | `arch/cpu/cpu6502/` | P02–P07 |
| Added | 31 | `arch/cpu/cpu65816/` | P20 |
| Added | 34 | `arch/cpu/cpu68000/` | P21; later P22 repairs if required |
| Added | 27 | `arch/cpu/cpu6809/` | P18 |
| Added | 25 | `arch/cpu/sm83/` | P19 |
| Added / Modified | 2 / 2 | `arch/cpu/x86/` | P10 |
| Added / Modified | 6 / 22 | `arch/cpu/z80/` | P11–P14 |
| Added | 8 | `arch/system/atari2600/` | P15 |
| Added | 6 | `arch/system/coco/` | P16 |
| Added | 5 | `arch/system/vectrex/` | P17 |
| Modified | 2 | `arch/system/nes/{nes,nes_test}.go` | P00: already present on target |
| Modified | 2 | `arch/system.go`, `arch/system_test.go` | P24 |
| Modified | 6 | `audio/`, `config/options_test.go`, `gui/sdl2/sdl_test.go` | P25 |
| Added | 5 | `docs/` | P00/P07/P14/P21–P23 |
| Modified | 1 | `testdata/Makefile` | P09 |

Total: 161 added and 62 modified files. No row is a rename. Excluding the two
identical NES files leaves 221 endpoint-different files. One of those is this
source-only tracking document. Thus 220 original product paths remain for
extraction, including the conditional P22 Makefile hunk and the reviewed C64
plan. A path count does not measure completion of mixed-file hunks. The
uncommitted review corrections add `README.md` and six new test files to this
inventory when they are committed.

Candidate-only work is separate from the 223-path inventory: P22's
evidence-driven corrections. The review corrections below supply the P04
concurrency regression, the P06 cycle tests, the reverse-mapping tests for the
new packages, and the README update on the source branch.

## Historical evidence and unresolved limits

These results come from the earlier branch record. They are not results of
this planning task or proof that an extracted candidate passes.

| Earlier check | Recorded result |
| --- | --- |
| `go fmt ./...`, `make lint`, `make test` | Passed through the recorded 2026-09-16 validation; short tests used the race detector. |
| CHIP-8 ROM tests | Six checks passed. |
| SM83 single-step and interrupt/fetch tests | Passed. |
| Tagged 65C816 single-step tests | Passed with the recorded corpus patch. |
| Z80 single-step | 1,604,000 vectors passed. |
| Z80 ZEXDOC / ZEXALL | All 67 groups passed in each; corpus revision was not recorded. |
| Tagged 68000 single-step | 996,321 passed; 3,739 failed; 1,000,060 total. |
| Pinned 6502 release gate | Procedure existed; no executed run was recorded before 2026-10-07. |

Preserve the 68000 discrepancies and the Z80 arbitrary-IM-0/T-state limits
in their package documents. The new system packages remain foundations.
The C64 plan remains future system work. The source's SH-family cycle-write
regression is an explicit target adaptation in P06.

## Review corrections (2026-10-07)

A code review of the complete branch on 2026-10-07 found defects that the
external corpora could not detect, and added the tests that the plan listed
as candidate-only work. The corrections are uncommitted working-tree changes
on the source branch. Each correction belongs to the part that owns its file
and must travel with that part. The behavior changes below are also public
contract changes; record them in the candidate PRs.

### Z80 (P12–P14)

- Q follows the hardware rule: Q is F after an instruction that writes F
  through the ALU flag logic and zero after any other instruction. POP AF and
  EX AF,AF' leave Q zero, as the corpus records. The runner now compares the
  final Q latch. The former model latched F after every instruction.
- An idle HALT cycle is an instruction boundary: it ends the EI delay and the
  LD A,I/R quirk window, and it resets Q. `Halt()` followed by `TriggerIRQ()`
  no longer leaves the CPU halted with the request pending.
- `EnableInterrupts`, `DisableInterrupts`, `SetInterruptMode`,
  `GetInterruptMode`, and `InterruptsEnabled` take the CPU lock.
- `State.Interrupts.IM` has the `InterruptMode` type.
- With tracing enabled, an accepted interrupt fills `TraceStep` with the
  handler address and the acceptance T-states.
- New tests in `interrupt_test.go` and `q_test.go`; `WithIOHandler` and the
  legacy adapter are documented. `docs/z80-gap-closure-plan.md` records the
  contract changes.

### 6809 (P18)

- The 22 inherent accumulator forms (NEGA/NEGB through CLRA/CLRB) have their
  own mnemonics, `OpcodeID` values, and registry entries. The former metadata
  reused the memory form's name, so a disassembler printed `neg` for `$40`.
- NMI is inhibited after reset until the first program load of S (LDS, LEAS,
  TFR/EXG into S, PULU of S).
- Extended indirect postbytes `$BF`, `$DF`, and `$FF` cost 5 cycles like `$9F`.
- TFR/EXG mixed sizes follow observed silicon: A and B read as `$FF00|r`, CC
  and DP as the duplicated byte, invalid codes as `$FFFF`.
- `OpcodeInfo` carries `Cycles`; `OpcodeID` values are renumbered
  alphabetically (22 new IDs); one `TestVerifyOpcodes` covers all three pages
  with pointer identity in both directions.

### SM83 (P19)

- BIT/RES/SET carry `BitOpcodes` reverse mappings for all 64 (bit, register)
  encodings; new `opcode_test.go` verifies both tables with zero skips.
- STOP has its own state that only a pending joypad interrupt ends; HALT keeps
  its wake rule. Idle Steps add one machine cycle.
- The four LDH forms are one `LdhInst` with `RegisterOpcodes`; `ldh` is in the
  memory read and write sets.
- `State` gains `IMEDelay`, `HaltBug`, and `Stopped`, and `SetState` restores
  them. Dead error, operand, and parameter declarations are removed.
- The corpus runner compares cycle counts; opcodes `$10` and `$76` are the only
  exclusions, with the reason recorded in the test.

### 65C816 (P20)

- MVN/MVP move one byte per Step and keep PC on the opcode until C wraps; the
  14-iteration harness rule left production code. One byte costs 7 cycles.
- Branch page-cross penalties apply only to taken branches.
- Direct-page and stack-relative 16-bit reads, writes, and read-modify-write
  operations stay in bank 0.
- Cycle rules for M and X width, `DL != 0`, indexed page crossing, 16-bit
  read-modify-write, and native-mode BRK/COP are implemented through
  per-opcode `CycleRule` flags. The corpus runner compares cycle counts; WAI,
  STP, and unfinished block moves are the only exclusions.
- `ErrNilMemory`, `ErrInvalidParameterType`, and `ErrMissingParameter` are
  used; `State.Interrupts` mirrors cpu6502; `BranchTarget` replaces untyped
  branch operands; idle WAI/STP Steps add one cycle; `ReadLong`-style helpers
  and `ReadVector` are gone from the `Memory` contract; `ReadWritesMemory` is
  added; `TestVerifyOpcodes` compares reverse opcode values.

### 68000 (P21–P22)

- Every decoder validates addressing-mode classes and size bits; 7,935
  formerly executed illegal encodings now raise vector 4. The new
  `opcode_map_test.go` checks all 65,536 words against the official opcode map.
- Unassigned opcode words raise the illegal-instruction exception instead of
  returning a host error; the unused error values are removed.
- Level-7 interrupts are edge-sensitive; a held level 7 no longer re-enters
  the handler on every Step.
- DIVS by zero sets N=V=C=0 and Z=1 (WinUAE reference).
- Trace stays pending across TRAP, TRAPV, CHK, and divide-by-zero exceptions.
- Scc to memory reads the destination before it writes.
- `ReadLong`/`WriteLong` leave the `Memory` interface; `WithSystemType` is
  removed; the acknowledged vector is masked to eight bits; memory category
  sets are regenerated from the handlers.
- Corpus result after the corrections: 996,321 passed, 3,739 failed, same
  classes. `docs/cpu68000-gap-closure-plan.md` records the independent
  references for the disputed vectors.

### 6502 (P03–P07) and CHIP-8 (P09)

- `shWrite` keeps the `cycleActive` dispatch (see P06); `TestBusCycleSHFamilyWrites`
  covers opcodes `93`, `9B`, `9C`, `9E`, and `9F` with and without page crossing.
- 65C02 INC/DEC abs,X return to a fixed 7 cycles; only ASL/LSR/ROL/ROR abs,X
  are 6+1 on the 65C02. The committed 6+1 metadata for `DE`/`FE` mismatched
  half of the wdc65c02 corpus vectors.
- 65C02 ADC/SBC add one cycle in decimal mode. `$5C` stays an 8-cycle NOP as
  the W65C02S datasheet states; the emulator-generated wdc65c02 corpus records
  4 cycles for it and is not followed on this point.
- The documented concurrency contract is now true: `StallCycles`, `SetIRQ`,
  `CheckInterrupts`, and the bus-cycle interrupt sampling take the lock. A
  `-race` test drives triggers from another goroutine in both execution modes.
  The bus-cycle path pays one uncontended read lock per cycle.
- `Step` documentation no longer claims interrupt service. TRB/TSB leave the
  memory read set. Name tables and `Instructions` are sorted; numeric
  `OpcodeID` values are unchanged. The required-data Dormann test returns after
  its assertion.
- CHIP-8 `State` exports `KeyWait` and `DrewThisFrame`, so an external
  serializer can round-trip a pending FX0A wait; restored indices are validated.
- New tests: independent `InstructionsForVariant` maps, negative-int address
  reads, 65C02 timing cases, and the CHIP-8 snapshot cases.

### System packages (P15–P17) and identifiers (P24)

- Vectrex: the VIA is a 6522, mirrored every 16 bytes through `$D000-$D7FF`;
  `$D800-$DFFF` selects RAM and VIA together and is marked as unusable; the
  joystick buttons come from PSG register 14, not VIA port B; port B signal
  bits and SWI2/SWI3 vectors are added. The former `$D800` VIA mirror and
  `ButtonRight`-style port B constants are removed.
- Atari 2600: the SWCHB difficulty comments read `0=B/amateur, 1=A/pro`; the
  package doc lists TIA read registers at `$0000-$000D`.
- CoCo: the package doc lists PIA mirrors and places GIME at `$FF90-$FFBF`.
- `arch/system.go`: `Generic` sorts before `Genesis`; the Game Boy comments no
  longer overlap.
- `README.md`: feature list and package tree include all new packages.

### Checks run on 2026-10-07

| Check | Result |
| --- | --- |
| `make build`, `make build-platforms`, `go vet ./...` | Passed before and after the corrections. |
| `golangci-lint run ./...` | 0 issues. |
| `retrogolint` `v1.0.6` (pinned) | 0 issues with `-exclude-dirs testdata`; the local `testdata/` clones contain `README.MD` files that the tool tries to parse. |
| `make test` | Passed. |
| Z80 single-step | 1,604 files passed with Q comparison. |
| SM83 single-step | 500 files, 500,000 cases passed with cycle comparison. |
| Tagged 65C816 single-step | 512 files, 5,120,000 cases passed with cycle comparison. |
| Tagged 68000 single-step | 996,321 passed; 3,739 failed; unchanged classes. |
| 6502 single-step (local corpus) | nes6502, 6502, synertek65c02, rockwell65c02: 256 files each passed; wdc65c02: 254 passed, 2 empty files skipped. |
| Pinned 6502 release gate | `TestReleaseNMOS` 151 files, 1,510,000 vectors passed; `TestReleaseDormann` both binaries passed. |

## Progress record

P00–P25 are planned. No candidate was extracted or merged. The two NES paths
are already on the verified target. The 2026-10-02 planning task checked local
refs, source diffs, target interfaces, test/helper ownership, and documentation
inventory without running Go checks. The 2026-10-07 review ran the checks in
the table above and left its corrections uncommitted on the source branch.

For each future part, record:

```text
Part:
Target base:
Fixed source SHA and any separately committed documentation input:
Candidate commit:
PR:
Included files and hunks:
Excluded or deferred mixed hunks:
Target-only behavior preserved:
Public API or behavior migration:
Candidate-only tests or repairs:
Dependency and workspace configuration:
Focused checks and results:
External data revisions, patches, executed cases, skips, and results:
Common gates and results:
Target merge commit and CI result:
Remaining source hunks and acceptance conditions:
```

For each source hunk not copied, record the path, symbol or test, disposition,
reason, and owning part. In particular, retain entries for the already-present
NES fix, P06's superseded `shWrite` hunk, this source-only plan, and any pending
P22 integration hunk. Recompute live range statistics after new commits; keep
the fixed-source values above as the original extraction inventory.
