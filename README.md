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

* **CPU emulation** - Chip-8, 6502-family, and Z80 emulators with tested instruction implementations
* **Instruction definitions** - x86 definitions from the 8086 through the 80486 for static analysis tools
* **System helpers** - NES cartridge, mapper, register, and parameter support
* **CGO-free GUI support** - Interfaces and SDL integration designed for straightforward cross-compilation
* **Tooling utilities** - Packages for CLI applications, configuration, structured logging, input, assertions, and sets
* **Small dependency footprint** - Go 1.22+ with only `ebitengine/purego` as an external dependency

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
    │  │  ├─ chip8                 Chip-8 virtual machine
    │  │  ├─ cpu6502               MOS 6502-family emulator, including NMOS and 65C02 variants
    │  │  ├─ x86                   Intel x86 instruction definitions from 8086 through 80486
    │  │  └─ z80                   Zilog Z80 emulator, including prefixed and undocumented opcodes
    │  └─ system/nes               Nintendo Entertainment System support
    │     ├─ cartridge             .nes ROM loading and saving
    │     ├─ codedatalog           FCEUX/Mesen-compatible code/data logging
    │     ├─ parameter             assembler-compatible instruction parameter formatting
    │     └─ register              NES memory-register constants
    ├─ assert                      test assertion helpers
    ├─ buildinfo                   embedded build-version metadata formatting
    ├─ cli                         command-line application helpers
    ├─ config                      configuration loading, parsing, and persistence
    ├─ gui                         CGO-free GUI rendering abstractions
    │  ├─ internal/dynlib          dynamic-library helpers for GUI backends
    │  ├─ internal/framebuffer     frame-buffer helpers for GUI backends
    │  └─ sdl2                     SDL2 GUI backend
    ├─ input                       keyboard and controller input helpers
    ├─ log                         nil-safe structured logging built on log/slog
    └─ set                         generic set data structures and operations

For package-level APIs and examples, see the [Go package documentation](https://pkg.go.dev/github.com/retroenv/retrogolib).
