//go:build singlestep

// Package m65816 provides SingleStepTests/65816 JSON-based single-step tests.
//
// To run these tests, download the test data:
//
//	make -C testdata m65816
//
// Then run: go test -tags singlestep ./arch/cpu/m65816/...
package m65816

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

const ssMaxFailures = 10

// TestSingleStep discovers and runs all SingleStepTests/65816 JSON test files.
func TestSingleStep(t *testing.T) {
	dataDir := getSingleStepDir(t)

	entries, err := os.ReadDir(dataDir)
	assert.NoError(t, err)

	found := false
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		filePath := filepath.Join(dataDir, entry.Name())
		name := entry.Name()[:len(entry.Name())-len(".json")]
		t.Run(name, func(t *testing.T) {
			runSS65816File(t, filePath)
		})
		found = true
	}

	if !found {
		t.Skipf("no JSON test files found in %s", dataDir)
	}
}

// ss65816State represents the 65816 CPU state in the SingleStepTests JSON format.
type ss65816State struct {
	PC  uint16      `json:"pc"`
	S   uint16      `json:"s"`
	P   uint8       `json:"p"`
	A   uint16      `json:"a"`
	X   uint16      `json:"x"`
	Y   uint16      `json:"y"`
	DBR uint8       `json:"dbr"`
	D   uint16      `json:"d"`
	PBR uint8       `json:"pbr"`
	E   uint8       `json:"e"`
	RAM [][2]uint32 `json:"ram"`
}

// ss65816TestCase represents a single test case from the 65816 JSON files.
type ss65816TestCase struct {
	Name    string       `json:"name"`
	Initial ss65816State `json:"initial"`
	Final   ss65816State `json:"final"`
	Cycles  [][]any      `json:"cycles"`
}

// ss65816Memory is a sparse 24-bit address space for test isolation.
type ss65816Memory struct {
	data map[uint32]uint8
}

func (m *ss65816Memory) Read(addr uint32) uint8 { return m.data[addr&0xFFFFFF] }
func (m *ss65816Memory) Write(addr uint32, v uint8) {
	m.data[addr&0xFFFFFF] = v
}
func (m *ss65816Memory) ReadWord(addr uint32) uint16 {
	addr &= 0xFFFFFF
	return uint16(m.data[addr]) | uint16(m.data[addr+1])<<8
}
func (m *ss65816Memory) WriteWord(addr uint32, v uint16) {
	addr &= 0xFFFFFF
	m.data[addr] = uint8(v)
	m.data[addr+1] = uint8(v >> 8)
}

var _ BasicMemory = (*ss65816Memory)(nil)

// getSingleStepDir returns the path to the 65816 SingleStepTests v1 directory,
// skipping the test if the data has not been downloaded yet.
func getSingleStepDir(t *testing.T) string {
	t.Helper()

	if dir := os.Getenv("M65816_TESTDATA"); dir != "" {
		if _, err := os.Stat(dir); err != nil {
			t.Skipf("M65816_TESTDATA directory not found: %s", dir)
		}
		return dir
	}

	_, thisFile, _, ok := runtime.Caller(0)
	assert.True(t, ok)

	dir := filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "testdata", "m65816", "65816", "v1")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Skipf("SingleStepTests 65816 data not found at %s (run 'make -C testdata m65816' to download)", dir)
	}

	return dir
}

// runSS65816File executes all test cases from a single JSON file.
func runSS65816File(t *testing.T, path string) {
	t.Helper()

	data, err := os.ReadFile(path)
	assert.NoError(t, err)
	if len(data) == 0 {
		t.Skipf("empty test file %s", path)
	}

	var cases []ss65816TestCase
	err = json.Unmarshal(data, &cases)
	assert.NoError(t, err)

	pass, fail := 0, 0
	for i := range cases {
		if runSS65816Case(t, &cases[i]) {
			pass++
		} else {
			fail++
			if fail >= ssMaxFailures {
				t.Logf("stopping after %d failures", fail)
				break
			}
		}
	}

	t.Logf("%s: %d passed, %d failed of %d", filepath.Base(path), pass, fail, len(cases))
}

// runSS65816Case sets up the CPU from the test's initial state, executes one
// Step, and compares the result against the expected final state.
func runSS65816Case(t *testing.T, tc *ss65816TestCase) bool {
	t.Helper()

	mem := &ss65816Memory{data: make(map[uint32]uint8, len(tc.Initial.RAM)*2)}
	for _, entry := range tc.Initial.RAM {
		mem.data[entry[0]&0xFFFFFF] = uint8(entry[1])
	}

	wrapped, err := NewMemory(mem)
	assert.NoError(t, err)

	cpu, err := New(wrapped)
	assert.NoError(t, err)

	// Load initial CPU state directly (bypass SetP side-effects for clean load).
	cpu.E = tc.Initial.E != 0
	cpu.Flags.Set(tc.Initial.P)
	cpu.PC = tc.Initial.PC
	cpu.SP = tc.Initial.S
	// In emulation mode the 65816 hardware forces SP high byte to $01.
	// The test data includes arbitrary SP values; normalize to page 1.
	if cpu.E {
		cpu.SP = 0x0100 | (cpu.SP & 0x00FF)
	}
	cpu.C = tc.Initial.A
	cpu.X = tc.Initial.X
	cpu.Y = tc.Initial.Y
	cpu.DB = tc.Initial.DBR
	cpu.DP = tc.Initial.D
	cpu.PB = tc.Initial.PBR

	err = cpu.Step()
	assert.NoError(t, err)

	return verifySS65816Case(t, tc, cpu, mem)
}

// verifySS65816Case compares all CPU registers and RAM writes against the
// expected final state, returning true if everything matches.
func verifySS65816Case(t *testing.T, tc *ss65816TestCase, cpu *CPU, mem *ss65816Memory) bool {
	t.Helper()

	var diffs []string

	if cpu.PC != tc.Final.PC {
		diffs = append(diffs, fmt.Sprintf("PC: got %04X, want %04X", cpu.PC, tc.Final.PC))
	}
	if cpu.SP != tc.Final.S {
		diffs = append(diffs, fmt.Sprintf("S: got %04X, want %04X", cpu.SP, tc.Final.S))
	}
	if cpu.C != tc.Final.A {
		diffs = append(diffs, fmt.Sprintf("C(A): got %04X, want %04X", cpu.C, tc.Final.A))
	}
	if cpu.X != tc.Final.X {
		diffs = append(diffs, fmt.Sprintf("X: got %04X, want %04X", cpu.X, tc.Final.X))
	}
	if cpu.Y != tc.Final.Y {
		diffs = append(diffs, fmt.Sprintf("Y: got %04X, want %04X", cpu.Y, tc.Final.Y))
	}
	if cpu.DB != tc.Final.DBR {
		diffs = append(diffs, fmt.Sprintf("DB: got %02X, want %02X", cpu.DB, tc.Final.DBR))
	}
	if cpu.DP != tc.Final.D {
		diffs = append(diffs, fmt.Sprintf("D: got %04X, want %04X", cpu.DP, tc.Final.D))
	}
	if cpu.PB != tc.Final.PBR {
		diffs = append(diffs, fmt.Sprintf("PB: got %02X, want %02X", cpu.PB, tc.Final.PBR))
	}
	wantE := tc.Final.E != 0
	if cpu.E != wantE {
		diffs = append(diffs, fmt.Sprintf("E: got %v, want %v", cpu.E, wantE))
	}
	gotP, wantP := cpu.Flags.Get(), tc.Final.P
	if gotP != wantP {
		diffs = append(diffs, fmt.Sprintf("P: got %08b (%02X), want %08b (%02X)", gotP, gotP, wantP, wantP))
	}
	for _, entry := range tc.Final.RAM {
		addr := entry[0] & 0xFFFFFF
		want := uint8(entry[1])
		if got := mem.data[addr]; got != want {
			diffs = append(diffs, fmt.Sprintf("RAM[%06X]: got %02X, want %02X", addr, got, want))
		}
	}

	if len(diffs) == 0 {
		return true
	}
	for _, d := range diffs {
		t.Errorf("[%s] %s", tc.Name, d)
	}
	return false
}
