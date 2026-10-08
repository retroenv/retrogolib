# retrogolib

[![CI](https://github.com/retroenv/retrogolib/actions/workflows/go.yaml/badge.svg?branch=main)](https://github.com/retroenv/retrogolib/actions/workflows/go.yaml)
[![Codecov](https://codecov.io/gh/retroenv/retrogolib/graph/badge.svg)](https://codecov.io/gh/retroenv/retrogolib)
[![Release](https://img.shields.io/github/v/release/retroenv/retrogolib)](https://github.com/retroenv/retrogolib/releases/latest)
[![Go Reference](https://pkg.go.dev/badge/github.com/retroenv/retrogolib.svg)](https://pkg.go.dev/github.com/retroenv/retrogolib)
[![License](https://img.shields.io/github/license/retroenv/retrogolib)](LICENSE)
![LLM assisted: human reviewed](https://img.shields.io/badge/LLM%20assisted-human%20reviewed-6f42c1)

A Go library of reusable components for retro-computing tools: emulators,
debuggers, disassemblers, and system-specific utilities.

## Features

* **CPU emulation** - Chip-8, 6502-family, 65C816, 6809, 68000, SM83, and Z80 emulators with tested instruction implementations
* **Instruction definitions** - x86 definitions from the 8086 through the 80486 for static analysis tools
* **System helpers** - NES cartridge, mapper, register, and parameter support; Atari 2600, CoCo, and Vectrex memory-map and register definitions
* **CGO-free GUI support** - Cross-platform rendering through SDL2 on Linux, Windows, macOS, and FreeBSD
* **CGO-free audio support** - Backend-neutral PCM playback with SDL2 output, concurrency-safe controls, and asynchronous error reporting
* **Tooling utilities** - Packages for CLI applications, configuration, structured logging, input, assertions, and sets
* **Small dependency footprint** - Go 1.25+ with only `ebitengine/purego` as an external dependency

## Installation

Add the module to an existing Go project:

```bash
go get github.com/retroenv/retrogolib
```

Then import the package needed by your tool:

```go
import "github.com/retroenv/retrogolib/arch/cpu/cpu6502"
```

## Packages

    ├─ app                         common application and service helpers
    ├─ arch                        shared architecture constants and types
    │  ├─ cpu
    │  │  ├─ chip8                 Chip-8 virtual machine with configurable quirks
    │  │  ├─ cpu6502               MOS 6502-family emulator, including NMOS and 65C02 variants
    │  │  ├─ cpu65816              WDC 65C816 emulator with emulation and native modes
    │  │  ├─ cpu6809               Motorola 6809 emulator
    │  │  ├─ cpu68000              Motorola 68000 emulator with bus-error and timing support
    │  │  ├─ sm83                  Sharp SM83 (Game Boy) emulator
    │  │  ├─ x86                   Intel x86 instruction definitions from 8086 through 80486
    │  │  └─ z80                   Zilog Z80 emulator, including prefixed and undocumented opcodes
    │  └─ system
    │     ├─ atari2600             Atari 2600 memory map and address constants
    │     │  ├─ cartridge          ROM loading and bank-switching scheme helpers
    │     │  └─ register           TIA and RIOT register constants
    │     ├─ coco                  TRS-80 Color Computer memory map and vectors
    │     │  └─ register           SAM and PIA register constants
    │     ├─ nes                   Nintendo Entertainment System support
    │     │  ├─ cartridge          .nes ROM loading and saving
    │     │  ├─ codedatalog        FCEUX/Mesen-compatible code/data logging
    │     │  ├─ parameter          assembler-compatible instruction parameter formatting
    │     │  └─ register           NES memory-register constants
    │     └─ vectrex               Vectrex memory map and vectors
    │        └─ register           VIA register and PSG button constants
    ├─ assert                      test assertion helpers
    ├─ audio                       PCM formats and playback controls
    │  └─ sdl2                     SDL2 audio backend
    ├─ buildinfo                   embedded build-version metadata formatting
    ├─ cli                         command-line application helpers
    ├─ config                      configuration loading, parsing, and persistence
    ├─ gui                         CGO-free GUI rendering abstractions
    │  ├─ internal/framebuffer     frame-buffer helpers for GUI backends
    │  └─ sdl2                     SDL2 GUI backend
    ├─ input                       keyboard and controller input helpers
    ├─ internal                    shared dynamic-library and SDL2 subsystem helpers
    ├─ log                         nil-safe structured logging built on log/slog
    └─ set                         generic set data structures and operations

For package-level APIs and examples, see the [Go package documentation](https://pkg.go.dev/github.com/retroenv/retrogolib).
