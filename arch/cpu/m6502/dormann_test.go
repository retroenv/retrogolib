// Package m6502 provides Klaus Dormann functional tests for the 6502/65C02 CPU emulator.
//
// To download test data and run:
//
//	make -C testdata m6502
//	make test-integration
package m6502

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

const (
	// dormannStartPC is the entry point for Klaus Dormann tests (code_segment).
	dormannStartPC = 0x0400

	// dormannMaxCycles prevents truly infinite loops on emulator bugs.
	dormannMaxCycles = uint64(200_000_000)

	// dormannProgressInterval prints progress every N cycles.
	dormannProgressInterval = uint64(10_000_000)

	// Success addresses from bin_files/*.lst (last "jmp *" = "test passed, no errors").
	nmos6502SuccessPC = uint16(0x3469) // 6502_functional_test.lst line: jmp * ;test passed, no errors
	c65c02SuccessPC   = uint16(0x24F1) // 65C02_extended_opcodes_test.lst line: jmp * ;test passed, no errors
)

// TestDormann runs the Klaus Dormann 6502/65C02 functional test ROMs.
func TestDormann(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Klaus Dormann tests in short mode")
	}

	dataDir := getDormannDataDir(t)

	tests := []dormannTest{
		{
			name:      "6502 functional test",
			binary:    filepath.Join("bin_files", "6502_functional_test.bin"),
			variant:   VariantNMOS6502,
			startPC:   dormannStartPC,
			successPC: nmos6502SuccessPC,
			maxCycles: dormannMaxCycles,
		},
		{
			name:      "65C02 extended opcodes test",
			binary:    filepath.Join("bin_files", "65C02_extended_opcodes_test.bin"),
			variant:   Variant65C02,
			startPC:   dormannStartPC,
			successPC: c65c02SuccessPC,
			maxCycles: dormannMaxCycles,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runDormannTest(t, dataDir, tc)
		})
	}
}

// dormannTest describes a single Klaus Dormann binary ROM test.
type dormannTest struct {
	name      string
	binary    string // relative path inside the dormann data dir
	variant   CPUVariant
	startPC   uint16
	successPC uint16
	maxCycles uint64
}

// runDormannTest loads a binary ROM and runs it until the CPU halts, then checks the success address.
func runDormannTest(t *testing.T, dataDir string, tc dormannTest) {
	t.Helper()

	path := filepath.Join(dataDir, tc.binary)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Skipf("test binary not found at %s (run 'make -C testdata m6502' to download): %v", path, err)
	}

	assert.LessOrEqual(t, len(data), 0x10000)

	// Load binary flat into 64KB RAM.
	mem := &testMemory{}
	copy(mem.b[:], data)

	// Write start PC into reset vector so New() picks it up.
	mem.b[ResetAddress] = uint8(tc.startPC)
	mem.b[ResetAddress+1] = uint8(tc.startPC >> 8)

	memory, err := NewMemory(mem)
	assert.NoError(t, err)

	cpu := New(memory, WithVariant(tc.variant))
	// Override PC directly after construction (reset vector already set correctly).
	cpu.PC = tc.startPC

	var (
		cycles uint64
		prevPC uint16
		halted bool
	)

	for cycles < tc.maxCycles {
		prevPC = cpu.PC

		stepErr := cpu.Step()
		assert.Nil(t, stepErr)

		cycles++

		if cycles%dormannProgressInterval == 0 {
			fmt.Printf("  %s: %d M cycles, PC=0x%04X\n", tc.name, cycles/1_000_000, cpu.PC)
		}

		if cpu.PC == prevPC {
			halted = true
			break
		}
	}

	assert.True(t, halted)
	assert.Equal(t, tc.successPC, cpu.PC)

	t.Logf("%s: PASSED at 0x%04X after %d cycles", tc.name, cpu.PC, cycles)
}

// getDormannDataDir returns the directory containing Klaus Dormann test data,
// skipping the test if it is not found.
func getDormannDataDir(t *testing.T) string {
	t.Helper()

	if dir := os.Getenv("M6502_DORMANN_TESTDATA"); dir != "" {
		return dir
	}

	_, thisFile, _, ok := runtime.Caller(0)
	assert.True(t, ok)

	dir := filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "testdata", "m6502", "dormann")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Skipf("Klaus Dormann test data not found at %s (run 'make -C testdata m6502' to download)", dir)
	}

	return dir
}
