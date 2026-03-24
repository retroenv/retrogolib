# Work Branch Changelog

Tracks the files changed on `work2` compared to `main`. The document is grouped by subsystem so
individual chunks can be extracted without copying the entire file.

**Last Updated:** 2026-03-24

## Section Index

- [Architecture Registration](#architecture-registration)
- [MOS 6502 / 65C02](#mos-6502--65c02)
- [WDC 65C816](#wdc-65c816)
- [Motorola 68000](#motorola-68000)
- [Sharp SM83](#sharp-sm83)
- [Zilog Z80](#zilog-z80)
- [Atari 2600 System](#atari-2600-system)
- [Documentation](#documentation)

---

## Architecture Registration

- `arch/arch.go`, `arch/arch_test.go`: added CPU architecture constants for `M65C02`, `M65816`,
  `M68000`, and `SM83`, with test updates for the expanded registry.
- `arch/system.go`, `arch/system_test.go`: added system constants for `AppleIIGS`, `Atari2600`,
  and `SNES`, with matching completeness test updates.

---

## MOS 6502 / 65C02

- `arch/cpu/m6502/addressing.go`, `arch/cpu/m6502/param.go`, `arch/cpu/m6502/memory.go`:
  added 65C02-only addressing support for zero-page indirect, absolute indexed indirect, and
  zero-page relative operands, plus the memory helpers needed to resolve them correctly.
- `arch/cpu/m6502/cpu.go`, `arch/cpu/m6502/option.go`, `arch/cpu/m6502/step.go`:
  added CPU variant selection, branch tracking, 6507 interrupt gating, and opcode-table
  selection based on the active variant.
- `arch/cpu/m6502/emulation.go`, `arch/cpu/m6502/interrupt.go`:
  added 65C02 decimal-mode behavior fixes, D-flag clearing on interrupts, and shared helpers
  needed by the variant-specific handlers.
- `arch/cpu/m6502/categories.go`, `arch/cpu/m6502/instruction.go`, `arch/cpu/m6502/opcode.go`,
  `arch/cpu/m6502/opcode_id.go`, `arch/cpu/m6502/opcode_test.go`, `arch/cpu/m6502/errors_test.go`:
  updated the instruction metadata, opcode tables, mnemonic IDs, and validation coverage to
  include the 65C02 and variant-aware 6502 changes.
- `arch/cpu/m6502/emulation_6507_test.go`, `arch/cpu/m6502/emulation_65c02.go`,
  `arch/cpu/m6502/emulation_65c02_test.go`, `arch/cpu/m6502/instruction_65c02.go`,
  `arch/cpu/m6502/opcode_65c02.go`, `arch/cpu/m6502/dormann_test.go`,
  `arch/cpu/m6502/singlestep_test.go`: added 6507 behavior tests, 65C02 instruction handlers,
  65C02 opcode coverage, and ROM/SingleStep validation for the expanded CPU model.

---

## WDC 65C816

- `arch/cpu/m65816/addressing.go`, `arch/cpu/m65816/categories.go`, `arch/cpu/m65816/cpu.go`,
  `arch/cpu/m65816/flag.go`, `arch/cpu/m65816/instruction.go`, `arch/cpu/m65816/opcode.go`,
  `arch/cpu/m65816/opcode_id.go`, `arch/cpu/m65816/option.go`, `arch/cpu/m65816/memory.go`,
  `arch/cpu/m65816/param.go`, `arch/cpu/m65816/step.go`, `arch/cpu/m65816/errors.go`:
  introduced the complete 65816 core model, including 24-bit addressing, register width control,
  opcode metadata, and the supporting helpers for decoding and execution.
- `arch/cpu/m65816/emulation.go`, `arch/cpu/m65816/emulation_branch.go`,
  `arch/cpu/m65816/emulation_move.go`, `arch/cpu/m65816/emulation_stack.go`,
  `arch/cpu/m65816/emulation_system.go`, `arch/cpu/m65816/interrupt.go`:
  implemented the main execution paths for arithmetic, branching, data movement, stack/system
  operations, and interrupts.
- `arch/cpu/m65816/doc.go`, `arch/cpu/m65816/cpu_test.go`,
  `arch/cpu/m65816/emulation_test.go`, `arch/cpu/m65816/opcode_test.go`,
  `arch/cpu/m65816/singlestep_test.go`: added package documentation and coverage for the new CPU.

---

## Motorola 68000

- `arch/cpu/m68000/addressing.go`, `arch/cpu/m68000/categories.go`, `arch/cpu/m68000/cpu.go`,
  `arch/cpu/m68000/flag.go`, `arch/cpu/m68000/instruction.go`, `arch/cpu/m68000/opcode.go`,
  `arch/cpu/m68000/opcode_id.go`, `arch/cpu/m68000/option.go`, `arch/cpu/m68000/memory.go`,
  `arch/cpu/m68000/param.go`, `arch/cpu/m68000/step.go`, `arch/cpu/m68000/interrupt.go`,
  `arch/cpu/m68000/errors.go`, `arch/cpu/m68000/doc.go`: introduced the full 68000 CPU model,
  including registers, addressing, opcode decoding, memory, interrupts, and package docs.
- `arch/cpu/m68000/emulation.go`, `arch/cpu/m68000/emulation_bit.go`,
  `arch/cpu/m68000/emulation_branch.go`, `arch/cpu/m68000/emulation_move.go`,
  `arch/cpu/m68000/emulation_shift.go`, `arch/cpu/m68000/emulation_system.go`, `arch/cpu/m68000/ea.go`:
  added the instruction execution layer and effective-address helpers for the 68000 core.
- `arch/cpu/m68000/cpu_test.go`, `arch/cpu/m68000/ea_test.go`,
  `arch/cpu/m68000/emulation_test.go`, `arch/cpu/m68000/memory_test.go`,
  `arch/cpu/m68000/opcode_test.go`, `arch/cpu/m68000/singlestep_test.go`:
  added coverage for CPU state, effective-address resolution, instruction execution, memory
  behavior, opcode decoding, and integration test vectors.

---

## Sharp SM83

- `arch/cpu/sm83/addressing.go`, `arch/cpu/sm83/categories.go`, `arch/cpu/sm83/cpu.go`,
  `arch/cpu/sm83/flag.go`, `arch/cpu/sm83/instruction.go`, `arch/cpu/sm83/opcode.go`,
  `arch/cpu/sm83/opcode_cb.go`, `arch/cpu/sm83/opcode_id.go`, `arch/cpu/sm83/option.go`,
  `arch/cpu/sm83/memory.go`, `arch/cpu/sm83/param.go`, `arch/cpu/sm83/step.go`,
  `arch/cpu/sm83/errors.go`, `arch/cpu/sm83/interrupt.go`, `arch/cpu/sm83/doc.go`:
  introduced the full Game Boy CPU model, including its reduced register set, CB-prefix support,
  interrupt model, and package documentation.
- `arch/cpu/sm83/emulation.go`, `arch/cpu/sm83/emulation_branch.go`,
  `arch/cpu/sm83/emulation_cb.go`, `arch/cpu/sm83/emulation_load.go`:
  implemented the SM83 execution paths for arithmetic, branching, prefixed bit ops, and loads.
- `arch/cpu/sm83/singlestep_test.go`: added integration test coverage for the SM83 single-step
  JSON vectors.

---

## Zilog Z80

- `arch/cpu/z80/categories.go`, `arch/cpu/z80/instruction.go`, `arch/cpu/z80/opcode_test.go`:
  adjusted category names, instruction comments, and opcode-test labels so the generated Z80
  tables and tests stay aligned with the current naming style.

---

## Atari 2600 System

- `arch/system/atari2600/atari2600.go`, `arch/system/atari2600/doc.go`,
  `arch/system/atari2600/atari2600_test.go`: added the Atari 2600 system package and its
  top-level validation.
- `arch/system/atari2600/cartridge/cartridge.go`, `arch/system/atari2600/cartridge/cartridge_test.go`:
  added cartridge parsing and bank-switch detection for the 2600 ROM formats supported by the
  system package.
- `arch/system/atari2600/register/riot.go`, `arch/system/atari2600/register/tia.go`,
  `arch/system/atari2600/register/register_test.go`: added TIA and RIOT register definitions and
  tests for the console-specific hardware map.

---

## Documentation

- `docs/cpu-implementation-plan-65816.md`, `docs/cpu-implementation-plan-65c02-68000.md`,
  `docs/cpu-implementation-plan-sm83.md`: implementation plans for the completed 65816, 65C02,
  68000, and SM83 work.
- `docs/m68000-emulator-comparison.md`, `docs/m68000-gap-closure-plan.md`:
  background and follow-up notes for the 68000 emulator.
- `docs/supported-systems.md`: consolidated overview of implemented CPUs, system packages, and
  remaining gaps.
- `docs/system-implementation-plan-atari2600.md`, `docs/system-implementation-plan-c64.md`,
  `docs/system-implementation-plan-gameboy.md`: system-level plans for the 2600, C64, and Game Boy.
- `docs/z80-emulator-comparison.md`, `docs/z80-emulator-comparison-cross-language.md`,
  `docs/z80-gap-closure-plan.md`: Z80 implementation comparisons and remaining gap analysis.
- `docs/work-branch-changelog.md`: this file, refreshed to match the current `main...HEAD` diff.
