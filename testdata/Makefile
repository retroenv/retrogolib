# Makefile for managing CPU test data.
#
# Downloads and updates external test suites used by CPU emulator integration tests.
# All cloned repositories are excluded from git via .gitignore.
#
# Usage:
#   make all       - download/update all test data
#   make cpu6502   - 65x02 SingleStepTests + Klaus Dormann functional tests
#   make cpu65816  - 65816 SingleStepTests (with open PR fixes applied)
#   make cpu68000  - 680x0 SingleStepTests
#   make sm83      - SM83 (Game Boy) SingleStepTests
#   make z80       - Z80 SingleStepTests + ZEXALL exerciser
#   make clean     - remove all downloaded test data

REPO_ROOT := $(shell git -C $(dir $(abspath $(lastword $(MAKEFILE_LIST)))) rev-parse --show-toplevel)

TESTDATA_ROOT := $(REPO_ROOT)/testdata

CPU6502_TESTDATA  := $(TESTDATA_ROOT)/cpu6502
CPU65816_TESTDATA := $(TESTDATA_ROOT)/cpu65816
CPU68000_TESTDATA := $(TESTDATA_ROOT)/cpu68000
SM83_TESTDATA     := $(TESTDATA_ROOT)/sm83
Z80_TESTDATA      := $(TESTDATA_ROOT)/z80
ZEXALL_TESTDATA   := $(TESTDATA_ROOT)/zexall

DIRTYHAIRY_65816 := https://github.com/DirtyHairy/65816.git
DORMANN_REPO     := https://github.com/Klaus2m5/6502_65C02_functional_tests
SINGLESTEP_BASE  := https://github.com/SingleStepTests
ZEXALL_REPO      := https://github.com/agn453/ZEXALL

.PHONY: all cpu6502 cpu6502-dormann cpu6502-singlestep cpu65816 cpu68000 sm83 z80 z80-singlestep zexall clean

all: cpu6502 cpu65816 cpu68000 sm83 z80

# 6502/65C02: SingleStepTests/65x02 + Klaus Dormann functional tests
cpu6502: cpu6502-singlestep cpu6502-dormann

cpu6502-singlestep:
	@echo "==> CPU 6502: SingleStepTests/65x02"
	@if [ -d "$(CPU6502_TESTDATA)/65x02/.git" ]; then \
		echo "  Updating..."; \
		git -C "$(CPU6502_TESTDATA)/65x02" pull --ff-only; \
	else \
		echo "  Cloning..."; \
		mkdir -p "$(CPU6502_TESTDATA)" && git clone --depth=1 $(SINGLESTEP_BASE)/65x02.git "$(CPU6502_TESTDATA)/65x02"; \
	fi

cpu6502-dormann:
	@echo "==> CPU 6502: Klaus Dormann functional tests"
	@mkdir -p "$(CPU6502_TESTDATA)/dormann"
	@if [ -d "$(CPU6502_TESTDATA)/dormann/.git" ]; then \
		echo "  Updating..."; \
		git -C "$(CPU6502_TESTDATA)/dormann" pull --ff-only; \
	else \
		echo "  Cloning (sparse, bin_files only)..."; \
		git clone --depth=1 --filter=blob:none --sparse $(DORMANN_REPO).git "$(CPU6502_TESTDATA)/dormann" && \
		git -C "$(CPU6502_TESTDATA)/dormann" sparse-checkout set bin_files && \
		git -C "$(CPU6502_TESTDATA)/dormann" checkout; \
	fi

# 65816: SingleStepTests/65816 with open PR fix applied
# PR #5 (DirtyHairy/fix_sbc_direct_x): fix SBC (dp,X) page wrap in emulation mode (v1/e1.e.json)
# PR #7 (DirtyHairy/fix_sbc_indirect_absolute_x): NOT applied — incorrect fix (final.ram vs cycles inconsistency)
cpu65816:
	@echo "==> CPU 65816: SingleStepTests/65816"
	@if [ -d "$(CPU65816_TESTDATA)/65816/.git" ]; then \
		echo "  Updating..."; \
		git -C "$(CPU65816_TESTDATA)/65816" pull --ff-only; \
	else \
		echo "  Cloning..."; \
		mkdir -p "$(CPU65816_TESTDATA)" && git clone --depth=1 $(SINGLESTEP_BASE)/65816.git "$(CPU65816_TESTDATA)/65816"; \
	fi
	@echo "  Applying PR #5 fix (SBC dp,X page wrap in emulation mode)..."
	@if ! git -C "$(CPU65816_TESTDATA)/65816" remote | grep -q dirtyhairy; then \
		git -C "$(CPU65816_TESTDATA)/65816" remote add dirtyhairy $(DIRTYHAIRY_65816); \
	fi
	@git -C "$(CPU65816_TESTDATA)/65816" fetch --depth=1 dirtyhairy fix_sbc_direct_x
	@git -C "$(CPU65816_TESTDATA)/65816" checkout dirtyhairy/fix_sbc_direct_x -- v1/e1.e.json
	@echo "  Done."

# 68000: SingleStepTests/680x0
cpu68000:
	@echo "==> CPU 68000: SingleStepTests/680x0"
	@if [ -d "$(CPU68000_TESTDATA)/680x0/.git" ]; then \
		echo "  Updating..."; \
		git -C "$(CPU68000_TESTDATA)/680x0" pull --ff-only; \
	else \
		echo "  Cloning..."; \
		mkdir -p "$(CPU68000_TESTDATA)" && git clone --depth=1 $(SINGLESTEP_BASE)/680x0.git "$(CPU68000_TESTDATA)/680x0"; \
	fi

# SM83: SingleStepTests/sm83
sm83:
	@echo "==> SM83: SingleStepTests/sm83"
	@if [ -d "$(SM83_TESTDATA)/.git" ]; then \
		echo "  Updating..."; \
		git -C "$(SM83_TESTDATA)" pull --ff-only; \
	else \
		echo "  Cloning..."; \
		git clone --depth=1 $(SINGLESTEP_BASE)/sm83.git "$(SM83_TESTDATA)"; \
	fi

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

clean:
	@echo "==> Removing all downloaded test data"
	rm -rf "$(CPU6502_TESTDATA)/65x02"
	rm -rf "$(CPU6502_TESTDATA)/dormann"
	rm -rf "$(CPU65816_TESTDATA)/65816"
	rm -rf "$(CPU68000_TESTDATA)/680x0"
	rm -rf "$(SM83_TESTDATA)"
	rm -rf "$(Z80_TESTDATA)"
	rm -rf "$(ZEXALL_TESTDATA)"
