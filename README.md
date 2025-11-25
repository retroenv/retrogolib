## RetroGoLib - a Golang library for retro console tooling development

[![Build status](https://github.com/retroenv/retrogolib/actions/workflows/go.yaml/badge.svg?branch=main)](https://github.com/retroenv/retrogolib/actions)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/retroenv/retrogolib)
[![Go Report Card](https://goreportcard.com/badge/github.com/retroenv/retrogolib)](https://goreportcard.com/report/github.com/retroenv/retrogolib)
[![codecov](https://codecov.io/gh/retroenv/retrogolib/branch/main/graph/badge.svg?token=jiBBxNmmVB)](https://app.codecov.io/gh/retroenv/retrogolib)

## Installation

```bash
go get github.com/retroenv/retrogolib
```

**Requirements:**
- Go 1.22 or later
- No CGO dependencies

## Overview

RetroGoLib is a Golang library designed to simplify the development of tools for retro consoles.
It provides a comprehensive set of functionalities for creating emulators, debugging tools, and other
retro console utilities, all while maintaining minimal dependencies and focusing on clean, maintainable code.

### Key Design Principles
- **Minimal Dependencies**: Only one external dependency (ebitengine/purego)
- **CGO-Free**: SDL support without CGO for easier cross-compilation
- **Type Safety**: Extensive use of Go generics for type-safe APIs
- **Thread Safety**: CPU implementations with proper synchronization patterns
- **Testing**: Comprehensive test coverage with consistent assertion patterns

## Supported Systems

### CPUs
- **6502**: Full instruction set with accurate timing
- **Chip-8**: Complete virtual machine implementation
- **Z80**: Complete Z80 CPU emulation with array-based opcode tables

### Consoles
- **NES (Nintendo Entertainment System)**: Cartridge formats, memory mapping

## Features

### CPU Emulation
- **6502 CPU**: Full instruction set with memory management, stack operations, and interrupt support
- **Chip-8 Virtual CPU**: Complete virtual machine with display, timers, and input handling
- **Z80 CPU**: Complete Z80 instruction set with 16-bit registers, prefix instructions (ED/DD/FD), and interrupts

### System Support
- **NES (Nintendo Entertainment System)**: Cartridge handling, memory mapping, and parameter conversion

## Package Overview

    ├─ app              common application/service helpers
    ├─ arch/cpu/chip8   Chip-8 virtual CPU support
    ├─ arch/cpu/m6502   6502 CPU support
    ├─ arch/cpu/z80     Z80 CPU support
    ├─ arch/system/nes  NES common types and helpers
    ├─ assert           test assertion helpers
    ├─ buildinfo        show version info that is embedded in the binary
    ├─ config           configuration management
    ├─ gui              GUI support - SDL without need for CGO
    ├─ input            hardware controller/keyboard helpers
    ├─ log              fast and structured logging based on slog
    ├─ set              generic set data structure with comprehensive operations

## API Documentation

For detailed API documentation, visit [pkg.go.dev](https://pkg.go.dev/github.com/retroenv/retrogolib).

## License

This project is licensed under the Apache License Version 2.0 - see the LICENSE file for details.
