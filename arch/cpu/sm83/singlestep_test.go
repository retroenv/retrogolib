package sm83

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

// TestSingleStep runs the SingleStepTests SM83 test suite.
// Each JSON file contains 1000 test cases that verify single-instruction execution
// against known-correct hardware traces.
func TestSingleStep(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping SingleStepTests in short mode")
	}

	dir := getSingleStepDir(t)

	files, err := filepath.Glob(filepath.Join(dir, "v1", "*.json"))
	assert.NoError(t, err, "globbing test files")
	assert.NotEqual(t, 0, len(files), "no SingleStepTests JSON files found")

	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			t.Parallel()
			runSingleStepFile(t, file)
		})
	}
}

// singleStepState represents the CPU state in a SingleStepTests SM83 JSON test case.
type singleStepState struct {
	PC  uint16   `json:"pc"`
	SP  uint16   `json:"sp"`
	A   uint8    `json:"a"`
	B   uint8    `json:"b"`
	C   uint8    `json:"c"`
	D   uint8    `json:"d"`
	E   uint8    `json:"e"`
	F   uint8    `json:"f"`
	H   uint8    `json:"h"`
	L   uint8    `json:"l"`
	IME uint8    `json:"ime"`
	EI  uint8    `json:"ei"`
	RAM [][2]int `json:"ram"`
}

// singleStepTest represents a single test case from SingleStepTests.
type singleStepTest struct {
	Name    string          `json:"name"`
	Initial singleStepState `json:"initial"`
	Final   singleStepState `json:"final"`
}

// getSingleStepDir returns the path to the sm83 SingleStepTests data directory,
// skipping the test if it is not found.
func getSingleStepDir(t *testing.T) string {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	assert.True(t, ok)

	dir := filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "testdata", "sm83")
	if _, err := os.Stat(filepath.Join(dir, "v1")); err != nil {
		t.Skipf("SingleStepTests sm83 data not found at %s (run 'make -C testdata sm83' to download)", dir)
	}

	return dir
}

// runSingleStepFile runs all test cases from a single JSON file.
func runSingleStepFile(t *testing.T, path string) {
	t.Helper()

	data, err := os.ReadFile(path)
	assert.NoError(t, err, "reading %s", path)

	var tests []singleStepTest
	err = json.Unmarshal(data, &tests)
	assert.NoError(t, err, "parsing %s", path)

	for i := range tests {
		tc := &tests[i]
		err = runSingleStepCase(tc)
		assert.NoError(t, err, "%s", tc.Name)
	}
}

// runSingleStepCase executes a single test case and returns an error describing the first mismatch.
func runSingleStepCase(tc *singleStepTest) error {
	mem := NewBasicMemory()

	// Set initial memory state.
	for _, entry := range tc.Initial.RAM {
		mem.Write(uint16(entry[0]), uint8(entry[1]))
	}

	cpu, err := New(mem, WithInitialPC(tc.Initial.PC), WithInitialSP(tc.Initial.SP))
	if err != nil {
		return fmt.Errorf("creating CPU: %w", err)
	}

	// Set initial CPU state.
	setSingleStepState(cpu, &tc.Initial)

	// Execute one instruction.
	if err := cpu.Step(); err != nil {
		return fmt.Errorf("Step: %w", err)
	}

	// Compare final state.
	return compareSingleStepState(cpu, mem, &tc.Final)
}

// setSingleStepState sets all CPU registers and flags from a test state.
func setSingleStepState(cpu *CPU, s *singleStepState) {
	cpu.PC = s.PC
	cpu.SP = s.SP
	cpu.A = s.A
	cpu.B = s.B
	cpu.C = s.C
	cpu.D = s.D
	cpu.E = s.E
	cpu.setFlags(s.F)
	cpu.H = s.H
	cpu.L = s.L
	cpu.ime = s.IME != 0
	cpu.imeDelay = s.EI != 0
}

// compareSingleStepState compares the CPU state against expected final state.
func compareSingleStepState(cpu *CPU, mem *BasicMemory, expected *singleStepState) error {
	checks := []struct {
		name string
		got  uint8
		want uint8
	}{
		{"A", cpu.A, expected.A}, {"B", cpu.B, expected.B},
		{"C", cpu.C, expected.C}, {"D", cpu.D, expected.D},
		{"E", cpu.E, expected.E}, {"H", cpu.H, expected.H},
		{"L", cpu.L, expected.L}, {"F", cpu.GetFlags(), expected.F},
	}
	for _, c := range checks {
		if c.got != c.want {
			return fmt.Errorf("%s: got 0x%02X, want 0x%02X", c.name, c.got, c.want)
		}
	}

	checks16 := []struct {
		name string
		got  uint16
		want uint16
	}{
		{"PC", cpu.PC, expected.PC}, {"SP", cpu.SP, expected.SP},
	}
	for _, c := range checks16 {
		if c.got != c.want {
			return fmt.Errorf("%s: got 0x%04X, want 0x%04X", c.name, c.got, c.want)
		}
	}

	// Compare IME state.
	if boolToUint8(cpu.ime) != expected.IME {
		return fmt.Errorf("IME: got %d, want %d", boolToUint8(cpu.ime), expected.IME)
	}

	// Compare memory.
	for _, entry := range expected.RAM {
		addr := uint16(entry[0])
		want := uint8(entry[1])
		got := mem.Read(addr)
		if got != want {
			return fmt.Errorf("RAM[0x%04X]: got 0x%02X, want 0x%02X", addr, got, want)
		}
	}

	return nil
}

// boolToUint8 converts a boolean to 0 or 1.
func boolToUint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}
