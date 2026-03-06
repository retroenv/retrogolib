//go:build singlestep

// Package m6502 provides SingleStepTests/65x02 JSON-based single-step tests.
//
// To run these tests, download the test data:
//
//	make -C testdata m6502
//
// Then run: go test -tags singlestep ./arch/cpu/m6502/...
package m6502

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

const (
	// ssMaxFailures limits the number of failures reported per file.
	ssMaxFailures = 10
)

// ss6502State represents the CPU state in the SingleStepTests JSON format.
type ss6502State struct {
	PC  uint16      `json:"pc"`
	S   uint8       `json:"s"`
	A   uint8       `json:"a"`
	X   uint8       `json:"x"`
	Y   uint8       `json:"y"`
	P   uint8       `json:"p"`
	RAM [][2]uint32 `json:"ram"`
}

// ss6502TestCase represents a single test case from the JSON files.
type ss6502TestCase struct {
	Name    string      `json:"name"`
	Initial ss6502State `json:"initial"`
	Final   ss6502State `json:"final"`
	Cycles  [][]any     `json:"cycles"`
}

// ssSparseMemory implements BasicMemory using a sparse map for test isolation.
type ssSparseMemory struct {
	data map[uint16]uint8
}

// Read returns the byte at the given address, or 0 if not set.
func (m *ssSparseMemory) Read(address uint16) uint8 {
	return m.data[address]
}

// Write stores a byte at the given address.
func (m *ssSparseMemory) Write(address uint16, value uint8) {
	m.data[address] = value
}

// getSingleStepDir returns the path to the m6502 SingleStepTests data directory,
// skipping the test if it is not found.
func getSingleStepDir(t *testing.T) string {
	t.Helper()

	if dir := os.Getenv("M6502_TESTDATA"); dir != "" {
		return dir
	}

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to determine source file location")
	}

	dir := filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "testdata", "m6502", "65x02")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Skipf("SingleStepTests m6502 data not found at %s (run 'make -C testdata m6502' to download)", dir)
	}

	return dir
}

// TestSingleStep discovers and runs all SingleStepTests/65x02 JSON test files.
func TestSingleStep(t *testing.T) {
	dataDir := getSingleStepDir(t)

	entries, err := os.ReadDir(dataDir)
	if err != nil {
		t.Fatalf("failed to read test data directory %s: %v", dataDir, err)
	}

	found := false
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		subDir := filepath.Join(dataDir, entry.Name())
		runSingleStepSubdir(t, subDir, entry.Name())
		found = true
	}

	if !found {
		t.Skipf("no test subdirectories found in %s", dataDir)
	}
}

// runSingleStepSubdir runs all JSON test files in a subdirectory.
func runSingleStepSubdir(t *testing.T, dir, variant string) {
	t.Helper()

	versionDir := filepath.Join(dir, "v1")
	if _, err := os.Stat(versionDir); err == nil {
		dir = versionDir
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Logf("skipping %s: cannot read directory: %v", variant, err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(dir, entry.Name())
		opName := entry.Name()[:len(entry.Name())-len(".json")]

		cpuVariant := ssVariantForDir(variant)
		t.Run(fmt.Sprintf("%s/%s", variant, opName), func(t *testing.T) {
			runSingleStepFile(t, filePath, cpuVariant)
		})
	}
}

// ssVariantForDir maps a test data directory name to a CPU variant.
func ssVariantForDir(dir string) CPUVariant {
	switch dir {
	case "nes6502":
		return VariantNES6502
	case "rockwell65c02", "synertek65c02", "wdc65c02":
		return Variant65C02
	default:
		return VariantNMOS6502
	}
}

// runSingleStepFile executes all test cases from a single JSON file.
func runSingleStepFile(t *testing.T, path string, variant CPUVariant) {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read test file %s: %v", path, err)
	}
	if len(data) == 0 {
		t.Skipf("empty test file %s", path)
	}

	var testCases []ss6502TestCase
	if err := json.Unmarshal(data, &testCases); err != nil {
		t.Fatalf("failed to parse test file %s: %v", path, err)
	}

	failures := 0
	for _, tc := range testCases {
		if !runSingleStepCase(t, tc, variant) {
			failures++
			if failures >= ssMaxFailures {
				t.Logf("stopping after %d failures", failures)
				return
			}
		}
	}
}

// runSingleStepCase executes a single test case and returns true if it passed.
func runSingleStepCase(t *testing.T, tc ss6502TestCase, variant CPUVariant) bool {
	t.Helper()

	mem := &ssSparseMemory{data: make(map[uint16]uint8, len(tc.Initial.RAM)+2)}
	for _, entry := range tc.Initial.RAM {
		mem.data[uint16(entry[0])] = uint8(entry[1])
	}
	// Ensure reset vector is set so New() can read the initial PC.
	// Only write each byte if the test case didn't specify it, to preserve the test state.
	if _, ok := mem.data[ResetAddress]; !ok {
		mem.data[ResetAddress] = uint8(tc.Initial.PC)
	}
	if _, ok := mem.data[ResetAddress+1]; !ok {
		mem.data[ResetAddress+1] = uint8(tc.Initial.PC >> 8)
	}

	memory, err := NewMemory(mem)
	if err != nil {
		t.Errorf("[%s] failed to create memory: %v", tc.Name, err)
		return false
	}

	cpu := New(memory, WithVariant(variant))

	// Load initial CPU state.
	cpu.PC = tc.Initial.PC
	cpu.SP = tc.Initial.S
	cpu.A = tc.Initial.A
	cpu.X = tc.Initial.X
	cpu.Y = tc.Initial.Y
	cpu.setFlags(tc.Initial.P)

	if err := cpu.Step(); err != nil {
		t.Errorf("[%s] Step() error: %v", tc.Name, err)
		return false
	}

	return verifySingleStepCase(t, tc, cpu, mem)
}

// verifySingleStepCase compares the CPU state against the expected final state.
func verifySingleStepCase(t *testing.T, tc ss6502TestCase, cpu *CPU, mem *ssSparseMemory) bool {
	t.Helper()

	passed := true

	if cpu.PC != tc.Final.PC {
		t.Errorf("[%s] PC: got 0x%04X, want 0x%04X", tc.Name, cpu.PC, tc.Final.PC)
		passed = false
	}

	if cpu.SP != tc.Final.S {
		t.Errorf("[%s] SP: got 0x%02X, want 0x%02X", tc.Name, cpu.SP, tc.Final.S)
		passed = false
	}

	if cpu.A != tc.Final.A {
		t.Errorf("[%s] A: got 0x%02X, want 0x%02X", tc.Name, cpu.A, tc.Final.A)
		passed = false
	}

	if cpu.X != tc.Final.X {
		t.Errorf("[%s] X: got 0x%02X, want 0x%02X", tc.Name, cpu.X, tc.Final.X)
		passed = false
	}

	if cpu.Y != tc.Final.Y {
		t.Errorf("[%s] Y: got 0x%02X, want 0x%02X", tc.Name, cpu.Y, tc.Final.Y)
		passed = false
	}

	// Compare P flags: mask out bit 5 (unused, always 1) and bit 4 (B flag, varies).
	const pMask = uint8(0b1100_1111)
	gotP := cpu.GetFlags() & pMask
	wantP := tc.Final.P & pMask
	if gotP != wantP {
		t.Errorf("[%s] P: got 0x%02X, want 0x%02X (masked: got 0x%02X, want 0x%02X)",
			tc.Name, cpu.GetFlags(), tc.Final.P, gotP, wantP)
		passed = false
	}

	for _, entry := range tc.Final.RAM {
		addr := uint16(entry[0])
		want := uint8(entry[1])
		got := mem.Read(addr)
		if got != want {
			t.Errorf("[%s] RAM[0x%04X]: got 0x%02X, want 0x%02X", tc.Name, addr, got, want)
			passed = false
		}
	}

	return passed
}

// This assertion verifies that ssSparseMemory implements BasicMemory at compile time.
var _ BasicMemory = (*ssSparseMemory)(nil)
