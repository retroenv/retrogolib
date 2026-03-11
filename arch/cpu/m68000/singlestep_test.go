//go:build singlestep

package m68000

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func newTestMemory() *testMemory {
	return &testMemory{data: make(map[uint32]uint8)}
}

// TestSingleStep runs the SingleStepTests/680x0 JSON test suite.
// Download test data: git clone https://github.com/SingleStepTests/680x0.git testdata/680x0
func TestSingleStep(t *testing.T) {
	testDir := getTestDataDir(t)

	files, err := filepath.Glob(filepath.Join(testDir, "*.json.gz"))
	assert.NoError(t, err)

	if len(files) == 0 {
		t.Skip("no .json.gz test files found in", testDir)
	}

	sort.Strings(files)

	var totalPass, totalFail int

	for _, file := range files {
		name := strings.TrimSuffix(filepath.Base(file), ".json.gz")
		t.Run(name, func(t *testing.T) {
			pass, fail := runTestFile(t, file)
			totalPass += pass
			totalFail += fail
		})
	}

	t.Logf("overall: %d passed, %d failed out of %d total",
		totalPass, totalFail, totalPass+totalFail)
}

// testState represents the CPU state from a JSON test case.
type testState struct {
	D0       uint32      `json:"d0"`
	D1       uint32      `json:"d1"`
	D2       uint32      `json:"d2"`
	D3       uint32      `json:"d3"`
	D4       uint32      `json:"d4"`
	D5       uint32      `json:"d5"`
	D6       uint32      `json:"d6"`
	D7       uint32      `json:"d7"`
	A0       uint32      `json:"a0"`
	A1       uint32      `json:"a1"`
	A2       uint32      `json:"a2"`
	A3       uint32      `json:"a3"`
	A4       uint32      `json:"a4"`
	A5       uint32      `json:"a5"`
	A6       uint32      `json:"a6"`
	USP      uint32      `json:"usp"`
	SSP      uint32      `json:"ssp"`
	SR       uint16      `json:"sr"`
	PC       uint32      `json:"pc"`
	Prefetch [2]uint16   `json:"prefetch"`
	RAM      [][2]uint32 `json:"ram"`
}

// testCase represents a single test from the JSON file.
type testCase struct {
	Name    string    `json:"name"`
	Initial testState `json:"initial"`
	Final   testState `json:"final"`
	Length  int       `json:"length"`
}

// testMemory implements Memory using a sparse map for test cases.
type testMemory struct {
	data map[uint32]uint8
}

// testBusForSingleStep wraps testMemory into a Bus with no IRQ activity.
type testBusForSingleStep struct {
	Memory
}

// Read reads a byte from memory at the given address.
func (m *testMemory) Read(address uint32) uint8 {
	return m.data[address&addressMask]
}

// ReadWord reads a 16-bit word from memory at the given address (big-endian).
func (m *testMemory) ReadWord(address uint32) uint16 {
	addr := address & addressMask
	return uint16(m.data[addr])<<8 | uint16(m.data[addr+1])
}

// ReadLong reads a 32-bit long word from memory at the given address (big-endian).
func (m *testMemory) ReadLong(address uint32) uint32 {
	addr := address & addressMask
	return uint32(m.data[addr])<<24 |
		uint32(m.data[addr+1])<<16 |
		uint32(m.data[addr+2])<<8 |
		uint32(m.data[addr+3])
}

// Write writes a byte to memory at the given address.
func (m *testMemory) Write(address uint32, value uint8) {
	m.data[address&addressMask] = value
}

// WriteWord writes a 16-bit word to memory at the given address (big-endian).
func (m *testMemory) WriteWord(address uint32, value uint16) {
	addr := address & addressMask
	m.data[addr] = uint8(value >> 8)
	m.data[addr+1] = uint8(value)
}

// WriteLong writes a 32-bit long word to memory at the given address (big-endian).
func (m *testMemory) WriteLong(address uint32, value uint32) {
	addr := address & addressMask
	m.data[addr] = uint8(value >> 24)
	m.data[addr+1] = uint8(value >> 16)
	m.data[addr+2] = uint8(value >> 8)
	m.data[addr+3] = uint8(value)
}

// IRQAcknowledge acknowledges an interrupt and returns the autovector number.
func (b *testBusForSingleStep) IRQAcknowledge(level uint8) uint32 {
	return uint32(VectorAutoVector1) + uint32(level) - 1
}

// IRQLevel returns 0 (no pending interrupts).
func (b *testBusForSingleStep) IRQLevel() uint8 { return 0 }

// OnReset handles the RESET instruction (no-op for tests).
func (b *testBusForSingleStep) OnReset() {}

func getTestDataDir(t *testing.T) string {
	t.Helper()

	if dir := os.Getenv("M68000_TESTDATA"); dir != "" {
		if _, err := os.Stat(dir); err != nil {
			t.Skip("M68000_TESTDATA directory not found:", dir)
		}
		return dir
	}

	_, thisFile, _, ok := runtime.Caller(0)
	assert.True(t, ok)

	dir := filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "testdata", "m68000", "680x0")
	if _, err := os.Stat(dir); err != nil {
		t.Skip("test data not found; run 'make -C testdata m68000' to download")
	}

	return dir
}

func loadTestCases(path string) ([]testCase, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	var cases []testCase
	if err := json.NewDecoder(gz).Decode(&cases); err != nil {
		return nil, err
	}

	return cases, nil
}

func runTestFile(t *testing.T, path string) (pass, fail int) {
	t.Helper()

	cases, err := loadTestCases(path)
	assert.NoError(t, err)

	const maxFailures = 10
	reported := 0

	for i := range cases {
		tc := &cases[i]

		ok := runSingleTest(t, tc)
		if ok {
			pass++
		} else {
			fail++
			reported++
			if reported >= maxFailures {
				t.Logf("stopping after %d failures (of %d tests)", maxFailures, len(cases))
				break
			}
		}
	}

	t.Logf("%s: %d passed, %d failed out of %d",
		filepath.Base(path), pass, fail, len(cases))

	return pass, fail
}

func runSingleTest(t *testing.T, tc *testCase) bool {
	t.Helper()

	mem := newTestMemory()
	bus := &testBusForSingleStep{Memory: mem}

	// Load initial RAM.
	for _, entry := range tc.Initial.RAM {
		mem.data[entry[0]&addressMask] = uint8(entry[1])
	}

	cpu := &CPU{
		bus: bus,
		opts: Options{
			initialPC: tc.Initial.PC,
			initialSP: tc.Initial.SSP,
		},
	}

	// Set registers.
	setDataRegisters(cpu, &tc.Initial)
	setAddrRegisters(cpu, &tc.Initial)

	// Set SR directly to avoid stack pointer swap side effects.
	cpu.sr = tc.Initial.SR & MaskSystem
	cpu.SetCCR(uint8(tc.Initial.SR & MaskCCR))

	cpu.PC = tc.Initial.PC
	cpu.USP = tc.Initial.USP
	cpu.SSP = tc.Initial.SSP

	// Set active stack pointer based on supervisor mode.
	if cpu.sr&MaskSupervisor != 0 {
		cpu.sp = cpu.SSP
	} else {
		cpu.sp = cpu.USP
	}

	// Execute one step.
	err := cpu.Step()
	if err != nil {
		t.Run(tc.Name, func(t *testing.T) {
			t.Helper()
			t.Errorf("step error: %v", err)
		})
		return false
	}

	// Compare final state.
	diffs := compareFinalState(cpu, mem, &tc.Final)
	if len(diffs) == 0 {
		return true
	}

	t.Run(tc.Name, func(t *testing.T) {
		t.Helper()
		for _, d := range diffs {
			t.Error(d)
		}
	})

	return false
}

func setDataRegisters(cpu *CPU, s *testState) {
	cpu.D[0] = s.D0
	cpu.D[1] = s.D1
	cpu.D[2] = s.D2
	cpu.D[3] = s.D3
	cpu.D[4] = s.D4
	cpu.D[5] = s.D5
	cpu.D[6] = s.D6
	cpu.D[7] = s.D7
}

func setAddrRegisters(cpu *CPU, s *testState) {
	cpu.A[0] = s.A0
	cpu.A[1] = s.A1
	cpu.A[2] = s.A2
	cpu.A[3] = s.A3
	cpu.A[4] = s.A4
	cpu.A[5] = s.A5
	cpu.A[6] = s.A6
}

func getDataRegisters(cpu *CPU) [8]uint32 {
	return cpu.D
}

func getAddrRegisters(cpu *CPU) [7]uint32 {
	return [7]uint32{
		cpu.A[0], cpu.A[1], cpu.A[2], cpu.A[3],
		cpu.A[4], cpu.A[5], cpu.A[6],
	}
}

func compareFinalState(cpu *CPU, mem *testMemory, final *testState) []string {
	var diffs []string

	// Compare data registers.
	expected := [8]uint32{
		final.D0, final.D1, final.D2, final.D3,
		final.D4, final.D5, final.D6, final.D7,
	}
	actual := getDataRegisters(cpu)

	for i := range 8 {
		if expected[i] != actual[i] {
			diffs = append(diffs,
				fmt.Sprintf("D%d: expected 0x%08X, got 0x%08X", i, expected[i], actual[i]))
		}
	}

	// Compare address registers.
	expectedA := [7]uint32{
		final.A0, final.A1, final.A2, final.A3,
		final.A4, final.A5, final.A6,
	}
	actualA := getAddrRegisters(cpu)

	for i := range 7 {
		if expectedA[i] != actualA[i] {
			diffs = append(diffs,
				fmt.Sprintf("A%d: expected 0x%08X, got 0x%08X", i, expectedA[i], actualA[i]))
		}
	}

	// Compare USP.
	if final.USP != cpu.USP {
		diffs = append(diffs,
			fmt.Sprintf("USP: expected 0x%08X, got 0x%08X", final.USP, cpu.USP))
	}

	// Compare SSP.
	if final.SSP != cpu.SSP {
		diffs = append(diffs,
			fmt.Sprintf("SSP: expected 0x%08X, got 0x%08X", final.SSP, cpu.SSP))
	}

	// Compare active stack pointer.
	expectedSP := final.SSP
	if final.SR&MaskSupervisor == 0 {
		expectedSP = final.USP
	}

	if expectedSP != cpu.sp {
		diffs = append(diffs,
			fmt.Sprintf("SP(active): expected 0x%08X, got 0x%08X", expectedSP, cpu.sp))
	}

	// Compare SR.
	actualSR := cpu.GetSR()
	if final.SR != actualSR {
		diffs = append(diffs,
			fmt.Sprintf("SR: expected 0x%04X, got 0x%04X", final.SR, actualSR))
	}

	// Compare PC.
	if final.PC != cpu.PC {
		diffs = append(diffs,
			fmt.Sprintf("PC: expected 0x%08X, got 0x%08X", final.PC, cpu.PC))
	}

	// Compare RAM.
	for _, entry := range final.RAM {
		addr := entry[0] & addressMask
		expectedByte := uint8(entry[1])
		actualByte := mem.data[addr]

		if expectedByte != actualByte {
			diffs = append(diffs,
				fmt.Sprintf("RAM[0x%06X]: expected 0x%02X, got 0x%02X",
					addr, expectedByte, actualByte))
		}
	}

	return diffs
}
