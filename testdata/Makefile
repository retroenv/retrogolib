# Makefile for managing CPU test data.
#
# Downloads and updates external test suites used by CPU emulator integration tests.
# All cloned repositories are excluded from git via .gitignore.
#
# Usage:
#   make all        - download/update all test data
#   make z80        - Z80 SingleStepTests
#   make m6502      - 65x02 SingleStepTests + Klaus Dormann functional tests
#   make m68000     - 680x0 SingleStepTests
#   make clean      - remove all downloaded test data

REPO_ROOT := $(shell git -C $(dir $(abspath $(lastword $(MAKEFILE_LIST)))) rev-parse --show-toplevel)

Z80_TESTDATA    := $(REPO_ROOT)/arch/cpu/z80/testdata
M6502_TESTDATA  := $(REPO_ROOT)/arch/cpu/m6502/testdata
M68000_TESTDATA := $(REPO_ROOT)/arch/cpu/m68000/testdata

SINGLESTEP_BASE := https://github.com/SingleStepTests
DORMANN_REPO    := https://github.com/Klaus2m5/6502_65C02_functional_tests

.PHONY: all z80 m6502 m68000 clean

all: z80 m6502 m68000

# Z80: SingleStepTests/z80
z80:
	@echo "==> Z80: SingleStepTests/z80"
	@if [ -d "$(Z80_TESTDATA)/singlestep/.git" ]; then \
		echo "  Updating..."; \
		git -C "$(Z80_TESTDATA)/singlestep" pull --ff-only; \
	else \
		echo "  Cloning..."; \
		git clone --depth=1 $(SINGLESTEP_BASE)/z80.git "$(Z80_TESTDATA)/singlestep"; \
	fi

# m6502/65C02: SingleStepTests/65x02 + Klaus Dormann functional tests
m6502: m6502-singlestep m6502-dormann

m6502-singlestep:
	@echo "==> m6502: SingleStepTests/65x02"
	@if [ -d "$(M6502_TESTDATA)/65x02/.git" ]; then \
		echo "  Updating..."; \
		git -C "$(M6502_TESTDATA)/65x02" pull --ff-only; \
	else \
		echo "  Cloning..."; \
		git clone --depth=1 $(SINGLESTEP_BASE)/65x02.git "$(M6502_TESTDATA)/65x02"; \
	fi

m6502-dormann:
	@echo "==> m6502: Klaus Dormann functional tests"
	@mkdir -p "$(M6502_TESTDATA)/dormann"
	@if [ -d "$(M6502_TESTDATA)/dormann/.git" ]; then \
		echo "  Updating..."; \
		git -C "$(M6502_TESTDATA)/dormann" pull --ff-only; \
	else \
		echo "  Cloning (sparse, bin_files only)..."; \
		git clone --depth=1 --filter=blob:none --sparse $(DORMANN_REPO).git "$(M6502_TESTDATA)/dormann"; \
		git -C "$(M6502_TESTDATA)/dormann" sparse-checkout set bin_files; \
		git -C "$(M6502_TESTDATA)/dormann" checkout; \
	fi

# m68000: SingleStepTests/680x0
m68000:
	@echo "==> m68000: SingleStepTests/680x0"
	@if [ -d "$(M68000_TESTDATA)/680x0/.git" ]; then \
		echo "  Updating..."; \
		git -C "$(M68000_TESTDATA)/680x0" pull --ff-only; \
	else \
		echo "  Cloning..."; \
		git clone --depth=1 $(SINGLESTEP_BASE)/680x0.git "$(M68000_TESTDATA)/680x0"; \
	fi

clean:
	@echo "==> Removing all downloaded test data"
	rm -rf "$(Z80_TESTDATA)/singlestep"
	rm -rf "$(M6502_TESTDATA)/65x02"
	rm -rf "$(M6502_TESTDATA)/dormann"
	rm -rf "$(M68000_TESTDATA)/680x0"
