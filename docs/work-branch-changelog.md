# Work Branch Changelog

Tracks the remaining `work2` changes relative to `main`.

**Last Updated:** 2026-04-27

---

## Current Branch State

- `main` was merged into `work2` in commit `6042731`.
- Post-merge duplicate declarations were removed in commit `787ffc6`.
- The source of truth for this document is the current `main...work2` diff, not older
  merge-plan assumptions.

## What Is Already In `main`

These areas no longer appear in `git diff --name-only main...work2`, so `work2` and
`main` currently agree on them:

- Test infrastructure changes such as `Makefile` and `testdata`
- Architecture registration updates in `arch/`
- 6502 / 65C02 integration and support work in `arch/cpu/m6502/`
- Most Z80 bus, opcode, and integration-test work
- x86 support in `arch/cpu/x86/`

## What Is Still Branch-Specific

These areas still differ from `main` and remain on `work2`:

### Group 1: m65816 CPU
**Status:** Branch-specific
**Files:** Entire `arch/cpu/m65816/` package

- 65C816 CPU emulation package
- Includes unit tests and `singlestep_test.go`

### Group 2: m68000 CPU
**Status:** Branch-specific
**Files:** Entire `arch/cpu/m68000/` package

- Motorola 68000 CPU emulation package
- Includes unit tests and `singlestep_test.go`

### Group 3: m6809 CPU
**Status:** Branch-specific
**Files:** Entire `arch/cpu/m6809/` package

- Motorola 6809 CPU emulation package
- Includes unit tests

### Group 4: SM83 CPU
**Status:** Branch-specific
**Files:** Entire `arch/cpu/sm83/` package

- Game Boy / Game Boy Color CPU package
- Includes `singlestep_test.go`

### Group 5: Z80 Follow-Up Changes
**Status:** Small branch-specific delta
**Files:**

- `arch/cpu/z80/categories.go`
- `arch/cpu/z80/instruction.go`
- `arch/cpu/z80/opcode_test.go`

### Group 6: Atari 2600 System
**Status:** Branch-specific
**Files:** Entire `arch/system/atari2600/` package

- System package
- Cartridge support
- Register definitions and tests

### Group 7: CoCo System
**Status:** Branch-specific
**Files:** Entire `arch/system/coco/` package

- System package
- Register definitions and tests

### Group 8: Vectrex System
**Status:** Branch-specific
**Files:** Entire `arch/system/vectrex/` package

- System package
- VIA/register definitions and tests

### Group 9: Documentation
**Status:** Branch-specific
**Files:**

- `docs/cpu-implementation-plan-65816.md`
- `docs/cpu-implementation-plan-65c02-68000.md`
- `docs/cpu-implementation-plan-sm83.md`
- `docs/m68000-emulator-comparison.md`
- `docs/m68000-gap-closure-plan.md`
- `docs/supported-systems.md`
- `docs/system-implementation-plan-atari2600.md`
- `docs/system-implementation-plan-c64.md`
- `docs/system-implementation-plan-gameboy.md`
- `docs/z80-emulator-comparison-cross-language.md`
- `docs/z80-emulator-comparison.md`
- `docs/z80-gap-closure-plan.md`
- `docs/work-branch-changelog.md`

## Merge Summary

Current `main...work2` diff summary:

- 127 files changed
- 29,084 insertions
- 5 deletions

Most remaining branch-specific work is additive new-package work. The only small
follow-up delta inside an existing package is the three-file Z80 change listed above.
