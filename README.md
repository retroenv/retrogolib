## RetroGoLib - a Golang library for retro console tooling development

[![Build status](https://github.com/retroenv/retrogolib/actions/workflows/go.yaml/badge.svg?branch=main)](https://github.com/retroenv/retrogolib/actions)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/retroenv/retrogolib)
[![Go Report Card](https://goreportcard.com/badge/github.com/retroenv/retrogolib)](https://goreportcard.com/report/github.com/retroenv/retrogolib)
[![codecov](https://codecov.io/gh/retroenv/retrogolib/branch/main/graph/badge.svg?token=jiBBxNmmVB)](https://app.codecov.io/gh/retroenv/retrogolib)

## Overview

RetroGoLib is a Golang library designed to simplify the development of tools for retro consoles.
It provides a set of functionalities that can be used to create tools for various retro consoles,
like emulators and debugging tools, all while keeping dependencies to a minimum.

The GUI renderer supports SDL without the use of CGO/SDL2, which makes it easier to build and distribute.

## Current system support

CPUs:
- 6502
- Chip-8

Systems:
- NES (Nintendo Entertainment System)

## Project layout

    ├─ app              common application/service helpers
    ├─ arch/cpu/chip8   Chip-8 virtual CPU support
    ├─ arch/cpu/m6502   6502 CPU support
    ├─ arch/nes         NES common types and helpers
    ├─ assert           test assertion helpers
    ├─ buildinfo        show version info that is embedded in the binary
    ├─ gui              GUI support - SDL without need for CGO
    ├─ input            hardware controller/keyboard helpers
    ├─ log              fast and structured logging based on slog
