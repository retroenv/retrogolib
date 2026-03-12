# Makefile for managing CPU test data.
#
# Downloads and updates external test suites used by CPU emulator integration tests.
# All cloned repositories are excluded from git via .gitignore.
#
# Usage:
#   make all        - download/update all test data
#   make z80        - Z80 SingleStepTests + ZEXALL exerciser
#   make m6502      - 65x02 SingleStepTests + Klaus Dormann functional tests
#   make m68000     - 680x0 SingleStepTests
#   make m65816     - 65816 SingleStepTests (with open PR fixes applied)
#   make clean      - remove all downloaded test data

REPO_ROOT := $(shell git -C $(dir $(abspath $(lastword $(MAKEFILE_LIST)))) rev-parse --show-toplevel)

TESTDATA_ROOT    := $(REPO_ROOT)/testdata
Z80_TESTDATA     := $(TESTDATA_ROOT)/z80
ZEXALL_TESTDATA  := $(TESTDATA_ROOT)/zexall
M6502_TESTDATA   := $(TESTDATA_ROOT)/m6502
M68000_TESTDATA  := $(TESTDATA_ROOT)/m68000
M65816_TESTDATA  := $(TESTDATA_ROOT)/m65816

SINGLESTEP_BASE  := https://github.com/SingleStepTests
DORMANN_REPO     := https://github.com/Klaus2m5/6502_65C02_functional_tests
ZEXALL_REPO      := https://github.com/agn453/ZEXALL
DIRTYHAIRY_65816 := https://github.com/DirtyHairy/65816.git

.PHONY: all z80 z80-singlestep zexall m6502 m6502-singlestep m6502-dormann m68000 m65816 clean

all: z80 m6502 m68000 m65816

# Z80: SingleStepTests/z80 + ZEXALL exerciser
z80: z80-singlestep zexall

# Z80: SingleStepTests/z80
z80-singlestep:
	@echo "==> Z80: SingleStepTests/z80"
	@if [ -d "$(Z80_TESTDATA)/.git" ]; then \
		echo "  Updating..."; \
		git -C "$(Z80_TESTDATA)" pull --ff-only; \
	else \
		echo "  Cloning..."; \
		git clone --depth=1 $(SINGLESTEP_BASE)/z80.git "$(Z80_TESTDATA)"; \
	fi

# Z80: ZEXALL/ZEXDOC exerciser binaries
zexall:
	@echo "==> Z80: ZEXALL exerciser"
	@if [ -d "$(ZEXALL_TESTDATA)/.git" ]; then \
		echo "  Updating..."; \
		git -C "$(ZEXALL_TESTDATA)" pull --ff-only; \
	else \
		echo "  Cloning..."; \
		git clone --depth=1 $(ZEXALL_REPO).git "$(ZEXALL_TESTDATA)"; \
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
		mkdir -p "$(M6502_TESTDATA)" && git clone --depth=1 $(SINGLESTEP_BASE)/65x02.git "$(M6502_TESTDATA)/65x02"; \
	fi

m6502-dormann:
	@echo "==> m6502: Klaus Dormann functional tests"
	@mkdir -p "$(M6502_TESTDATA)/dormann"
	@if [ -d "$(M6502_TESTDATA)/dormann/.git" ]; then \
		echo "  Updating..."; \
		git -C "$(M6502_TESTDATA)/dormann" pull --ff-only; \
	else \
		echo "  Cloning (sparse, bin_files only)..."; \
		git clone --depth=1 --filter=blob:none --sparse $(DORMANN_REPO).git "$(M6502_TESTDATA)/dormann" && \
		git -C "$(M6502_TESTDATA)/dormann" sparse-checkout set bin_files && \
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
		mkdir -p "$(M68000_TESTDATA)" && git clone --depth=1 $(SINGLESTEP_BASE)/680x0.git "$(M68000_TESTDATA)/680x0"; \
	fi

# m65816: SingleStepTests/65816 with open PR fix applied
# PR #5 (DirtyHairy/fix_sbc_direct_x): fix SBC (dp,X) page wrap in emulation mode (v1/e1.e.json)
# PR #7 (DirtyHairy/fix_sbc_indirect_absolute_x): NOT applied — incorrect fix (final.ram vs cycles inconsistency)
m65816:
	@echo "==> m65816: SingleStepTests/65816"
	@if [ -d "$(M65816_TESTDATA)/65816/.git" ]; then \
		echo "  Updating..."; \
		git -C "$(M65816_TESTDATA)/65816" pull --ff-only; \
	else \
		echo "  Cloning..."; \
		mkdir -p "$(M65816_TESTDATA)" && git clone --depth=1 $(SINGLESTEP_BASE)/65816.git "$(M65816_TESTDATA)/65816"; \
	fi
	@echo "  Applying PR #5 fix (SBC dp,X page wrap in emulation mode)..."
	@if ! git -C "$(M65816_TESTDATA)/65816" remote | grep -q dirtyhairy; then \
		git -C "$(M65816_TESTDATA)/65816" remote add dirtyhairy $(DIRTYHAIRY_65816); \
	fi
	@git -C "$(M65816_TESTDATA)/65816" fetch --depth=1 dirtyhairy fix_sbc_direct_x
	@git -C "$(M65816_TESTDATA)/65816" checkout dirtyhairy/fix_sbc_direct_x -- v1/e1.e.json
	@echo "  Done."

clean:
	@echo "==> Removing all downloaded test data"
	rm -rf "$(Z80_TESTDATA)"
	rm -rf "$(ZEXALL_TESTDATA)"
	rm -rf "$(M6502_TESTDATA)/65x02"
	rm -rf "$(M6502_TESTDATA)/dormann"
	rm -rf "$(M68000_TESTDATA)/680x0"
	rm -rf "$(M65816_TESTDATA)/65816"
