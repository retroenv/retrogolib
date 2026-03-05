//go:build dormann

// Package m6502 provides Klaus Dormann functional tests for the 6502 CPU emulator.
//
// To run these tests, download the Klaus Dormann test binaries:
//
//	mkdir -p arch/cpu/m6502/testdata/dormann
//	cd arch/cpu/m6502/testdata/dormann
//	wget https://github.com/Klaus2m5/6502_65C02_functional_tests/raw/master/bin_files/6502_functional_test.bin
//	wget https://github.com/Klaus2m5/6502_65C02_functional_tests/raw/master/bin_files/65C02_extended_opcodes_test.bin
//	wget https://github.com/Klaus2m5/6502_65C02_functional_tests/raw/master/bin_files/6502_decimal_test.bin
//
// Then run: go test -tags dormann -timeout 30m ./arch/cpu/m6502/...
package m6502

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

const (
	// dormannStartPC is the standard entry point for Klaus Dormann tests.
	dormannStartPC = 0x0400

	// dormannMaxCycles limits execution to prevent infinite loops on bugs.
	dormannMaxCycles = uint64(200_000_000)

	// dormannProgressInterval controls how often progress is reported.
	dormannProgressInterval = uint64(10_000_000)

	// nmos6502SuccessPC is the known success address for 6502_functional_test.bin.
	nmos6502SuccessPC = uint16(0x3B1C)

	// decimalErrorAddr is the address checked for errors in the decimal test.
	decimalErrorAddr = uint16(0x000B)
)

// dormannTest describes a single Klaus Dormann binary ROM test.
type dormannTest struct {
	name        string
	binaryPath  string
	startPC     uint16
	successPC   uint16 // 0 means use decimal test detection.
	maxCycles   uint64
}

// TestDormann runs the Klaus Dormann 6502 functional tests.
func TestDormann(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Klaus Dormann tests in short mode")
	}

	dataDir := os.Getenv("M6502_TESTDATA")
	if dataDir == "" {
		dataDir = filepath.Join("testdata", "dormann")
	}

	tests := []dormannTest{
		{
			name:       "6502 functional test",
			binaryPath: filepath.Join(dataDir, "6502_functional_test.bin"),
			startPC:    dormannStartPC,
			successPC:  nmos6502SuccessPC,
			maxCycles:  dormannMaxCycles,
		},
		{
			name:       "6502 decimal test",
			binaryPath: filepath.Join(dataDir, "6502_decimal_test.bin"),
			startPC:    dormannStartPC,
			successPC:  0, // decimal test uses memory-based error detection
			maxCycles:  dormannMaxCycles,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runDormannTest(t, tc)
		})
	}
}

// runDormannTest executes a single Dormann binary ROM test.
func runDormannTest(t *testing.T, tc dormannTest) {
	t.Helper()

	data, err := os.ReadFile(tc.binaryPath)
	if err != nil {
		t.Skipf("test binary not found at %s - download it to run this test: %v", tc.binaryPath, err)
	}

	if len(data) > 0x10000 {
		t.Fatalf("binary too large: %d bytes (max 65536)", len(data))
	}

	mem := &testMemory{}
	for i, b := range data {
		mem.b[i] = b
	}

	memory, err := NewMemory(mem)
	assert.NoError(t, err)

	// Write the start PC into the reset vector so New() loads it correctly.
	memory.WriteWord(ResetAddress, tc.startPC)

	cpu := New(memory)
	cpu.PC = tc.startPC

	var (
		cycles   uint64
		prevPC   uint16
		halted   bool
	)

	for cycles < tc.maxCycles {
		prevPC = cpu.PC

		if err := cpu.Step(); err != nil {
			t.Fatalf("CPU error at PC=0x%04X after %d cycles: %v", cpu.PC, cycles, err)
		}

		cycles++

		if cycles%dormannProgressInterval == 0 {
			fmt.Printf("  %s: %d M cycles, PC=0x%04X\n", tc.name, cycles/1_000_000, cpu.PC)
		}

		if cpu.PC == prevPC {
			halted = true
			break
		}
	}

	if !halted {
		t.Fatalf("%s: did not halt after %d cycles (last PC=0x%04X)", tc.name, cycles, cpu.PC)
	}

	if tc.successPC != 0 {
		if cpu.PC != tc.successPC {
			t.Fatalf("%s: halted at 0x%04X, expected success address 0x%04X", tc.name, cpu.PC, tc.successPC)
		}
		t.Logf("%s: PASSED at 0x%04X after %d cycles", tc.name, cpu.PC, cycles)
		return
	}

	// Decimal test: check error byte at known address.
	errVal := memory.Read(decimalErrorAddr)
	if errVal != 0 {
		t.Fatalf("%s: decimal test failed with error code 0x%02X at address 0x%04X", tc.name, errVal, decimalErrorAddr)
	}
	t.Logf("%s: PASSED after %d cycles", tc.name, cycles)
}
