## RetroGoLib - a Golang library for retro console tooling development

[![Build status](https://github.com/retroenv/retrogolib/actions/workflows/go.yaml/badge.svg?branch=main)](https://github.com/retroenv/retrogolib/actions)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/retroenv/retrogolib)
[![Go Report Card](https://goreportcard.com/badge/github.com/retroenv/retrogolib)](https://goreportcard.com/report/github.com/retroenv/retrogolib)
[![codecov](https://codecov.io/gh/retroenv/retrogolib/branch/main/graph/badge.svg?token=jiBBxNmmVB)](https://app.codecov.io/gh/retroenv/retrogolib)

## Installation

```bash
go get github.com/retroenv/retrogolib
```

## Overview

RetroGoLib is a Go library for building retro computing tools such as emulators, debuggers,
disassemblers, and system-specific utilities. It focuses on small dependencies, reusable CPU and
system components, and clean APIs that compose well in standalone tools.

## Highlights

- Go 1.22 or later
- Single external dependency: [`ebitengine/purego`](https://github.com/ebitengine/purego)
- No CGO dependencies, including SDL support for easier cross-compilation
- Reusable CPU emulation packages with tests and system helpers
- Supporting utility packages for CLI apps, configuration, logging, input, and assertions

## CPU Packages

These packages are currently implemented in the repository:

- `arch/cpu/chip8`: Chip-8 virtual machine
- `arch/cpu/m6502`: MOS 6502, including 65C02-related support
- `arch/cpu/x86`: x86 instruction definitions for 8086 through 80486
- `arch/cpu/z80`: Zilog Z80, including prefixed and undocumented opcode support

## System Support

Concrete system helper packages currently include:

- `arch/system/nes`: NES cartridge, mapper, register, and parameter helpers

## Package Overview

    ├─ app               common application and service helpers
    ├─ arch              shared architecture and system identifiers
    ├─ arch/cpu/*        CPU emulation and virtual-machine packages
    ├─ arch/system/*     system-specific constants and helpers
    ├─ assert            test assertion helpers
    ├─ buildinfo         version metadata helpers
    ├─ cli               command-line parsing and related utilities
    ├─ config            configuration loading, parsing, and persistence
    ├─ gui               graphical interface helpers with CGO-free SDL integration
    ├─ input             keyboard and controller input helpers
    ├─ log               structured logging helpers built on slog
    └─ set               generic set data structures and operations

## API Documentation

For detailed package documentation, visit [pkg.go.dev](https://pkg.go.dev/github.com/retroenv/retrogolib).

## License

This project is licensed under the Apache License Version 2.0 - see the LICENSE file for details.
