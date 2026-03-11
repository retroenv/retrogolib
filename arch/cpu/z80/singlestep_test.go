package z80

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func newTestIOHandler(ports []singleStepPort) *testIOHandler {
	h := &testIOHandler{
		reads:  make(map[uint8]uint8),
		writes: make(map[uint8]uint8),
	}
	for _, p := range ports {
		port := uint8(p.Address)
		if p.IsRead {
			h.reads[port] = p.Value
		}
	}
	return h
}

// TestSingleStep runs the SingleStepTests Z80 test suite.
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

// singleStepState represents the CPU state in a SingleStepTests JSON test case.
type singleStepState struct {
	PC   uint16   `json:"pc"`
	SP   uint16   `json:"sp"`
	A    uint8    `json:"a"`
	B    uint8    `json:"b"`
	C    uint8    `json:"c"`
	D    uint8    `json:"d"`
	E    uint8    `json:"e"`
	F    uint8    `json:"f"`
	H    uint8    `json:"h"`
	L    uint8    `json:"l"`
	I    uint8    `json:"i"`
	R    uint8    `json:"r"`
	WZ   uint16   `json:"wz"`
	IX   uint16   `json:"ix"`
	IY   uint16   `json:"iy"`
	AF_  uint16   `json:"af_"` //nolint:revive // matches JSON field name
	BC_  uint16   `json:"bc_"` //nolint:revive // matches JSON field name
	DE_  uint16   `json:"de_"` //nolint:revive // matches JSON field name
	HL_  uint16   `json:"hl_"` //nolint:revive // matches JSON field name
	IM   uint8    `json:"im"`
	IFF1 uint8    `json:"iff1"`
	IFF2 uint8    `json:"iff2"`
	Q    uint8    `json:"q"`
	RAM  [][2]int `json:"ram"`
}

// singleStepPort represents a port access in a test case.
// Format: [address, value, "r"|"w"]
type singleStepPort struct {
	Address uint16
	Value   uint8
	IsRead  bool
}

// singleStepTest represents a single test case from SingleStepTests.
type singleStepTest struct {
	Name    string          `json:"name"`
	Initial singleStepState `json:"initial"`
	Final   singleStepState `json:"final"`
	Ports   []singleStepPort
}

// testIOHandler provides port read/write values for test cases.
type testIOHandler struct {
	reads  map[uint8]uint8
	writes map[uint8]uint8
}

// UnmarshalJSON handles the ports field which is an array of [addr, val, "r"|"w"].
func (t *singleStepTest) UnmarshalJSON(data []byte) error {
	type alias singleStepTest
	var raw struct {
		alias
		RawPorts []json.RawMessage `json:"ports"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("unmarshaling test: %w", err)
	}
	*t = singleStepTest(raw.alias)

	for _, rp := range raw.RawPorts {
		var arr []json.RawMessage
		if err := json.Unmarshal(rp, &arr); err != nil {
			continue
		}
		if len(arr) < 3 {
			continue
		}
		var addr float64
		var val float64
		var dir string
		if err := json.Unmarshal(arr[0], &addr); err != nil {
			continue
		}
		if err := json.Unmarshal(arr[1], &val); err != nil {
			continue
		}
		if err := json.Unmarshal(arr[2], &dir); err != nil {
			continue
		}
		t.Ports = append(t.Ports, singleStepPort{
			Address: uint16(addr),
			Value:   uint8(val),
			IsRead:  dir == "r",
		})
	}
	return nil
}

func (h *testIOHandler) ReadPort(port uint8) uint8 {
	if val, ok := h.reads[port]; ok {
		return val
	}
	return 0xFF
}

func (h *testIOHandler) WritePort(_ uint8, _ uint8) {}

// getSingleStepDir returns the path to the z80 SingleStepTests data directory,
// skipping the test if it is not found.
func getSingleStepDir(t *testing.T) string {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	assert.True(t, ok)

	dir := filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "testdata", "z80")
	if _, err := os.Stat(filepath.Join(dir, "v1")); err != nil {
		t.Skipf("SingleStepTests z80 data not found at %s (run 'make -C testdata z80' to download)", dir)
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

	// Create CPU with I/O handler if ports are used.
	var opts []Option
	if len(tc.Ports) > 0 {
		opts = append(opts, WithIOHandler(newTestIOHandler(tc.Ports)))
	}
	cpu, err := New(mem, opts...)
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
	cpu.I = s.I
	cpu.R = s.R
	cpu.MEMPTR = s.WZ
	cpu.IX = s.IX
	cpu.IY = s.IY

	// Shadow registers: af_ is A'F', bc_ is B'C', etc.
	cpu.AltA = uint8(s.AF_ >> 8)
	setAltFlags(&cpu.AltFlags, uint8(s.AF_))
	cpu.AltB = uint8(s.BC_ >> 8)
	cpu.AltC = uint8(s.BC_)
	cpu.AltD = uint8(s.DE_ >> 8)
	cpu.AltE = uint8(s.DE_)
	cpu.AltH = uint8(s.HL_ >> 8)
	cpu.AltL = uint8(s.HL_)

	// Interrupt state.
	cpu.iff1 = s.IFF1 != 0
	cpu.iff2 = s.IFF2 != 0
	cpu.im = s.IM

	// Q register for SCF/CCF X/Y flag behavior.
	cpu.q = s.Q
}

// compareReg8 compares an 8-bit register value against expected.
func compareReg8(name string, got, want uint8) error {
	if got != want {
		return fmt.Errorf("%s: got 0x%02X, want 0x%02X", name, got, want)
	}
	return nil
}

// compareReg16 compares a 16-bit register value against expected.
func compareReg16(name string, got, want uint16) error {
	if got != want {
		return fmt.Errorf("%s: got 0x%04X, want 0x%04X", name, got, want)
	}
	return nil
}

// compareSingleStepState compares the CPU state against expected final state.
// Returns an error describing the first mismatch found.
func compareSingleStepState(cpu *CPU, mem *BasicMemory, expected *singleStepState) error {
	if err := compareSingleStepRegisters(cpu, expected); err != nil {
		return err
	}
	if err := compareSingleStepShadow(cpu, expected); err != nil {
		return err
	}
	if err := compareSingleStepInterrupts(cpu, expected); err != nil {
		return err
	}
	return compareSingleStepRAM(mem, expected)
}

// compareSingleStepRegisters compares main registers, flags, and special registers.
func compareSingleStepRegisters(cpu *CPU, expected *singleStepState) error {
	checks := []struct {
		name string
		got  uint8
		want uint8
	}{
		{"A", cpu.A, expected.A}, {"B", cpu.B, expected.B},
		{"C", cpu.C, expected.C}, {"D", cpu.D, expected.D},
		{"E", cpu.E, expected.E}, {"H", cpu.H, expected.H},
		{"L", cpu.L, expected.L}, {"F", cpu.GetFlags(), expected.F},
		{"I", cpu.I, expected.I}, {"R", cpu.R, expected.R},
	}
	for _, c := range checks {
		if err := compareReg8(c.name, c.got, c.want); err != nil {
			return err
		}
	}

	checks16 := []struct {
		name string
		got  uint16
		want uint16
	}{
		{"PC", cpu.PC, expected.PC}, {"SP", cpu.SP, expected.SP},
		{"IX", cpu.IX, expected.IX}, {"IY", cpu.IY, expected.IY},
		{"MEMPTR", cpu.MEMPTR, expected.WZ},
	}
	for _, c := range checks16 {
		if err := compareReg16(c.name, c.got, c.want); err != nil {
			return err
		}
	}
	return nil
}

// compareSingleStepShadow compares shadow register pairs.
func compareSingleStepShadow(cpu *CPU, expected *singleStepState) error {
	checks := []struct {
		name string
		got  uint16
		want uint16
	}{
		{"AF'", uint16(cpu.AltA)<<8 | uint16(getAltFlagsAsUint8(cpu.AltFlags)), expected.AF_},
		{"BC'", uint16(cpu.AltB)<<8 | uint16(cpu.AltC), expected.BC_},
		{"DE'", uint16(cpu.AltD)<<8 | uint16(cpu.AltE), expected.DE_},
		{"HL'", uint16(cpu.AltH)<<8 | uint16(cpu.AltL), expected.HL_},
	}
	for _, c := range checks {
		if err := compareReg16(c.name, c.got, c.want); err != nil {
			return err
		}
	}
	return nil
}

// compareSingleStepInterrupts compares interrupt state.
func compareSingleStepInterrupts(cpu *CPU, expected *singleStepState) error {
	if err := compareReg8("IFF1", boolToUint8(cpu.iff1), expected.IFF1); err != nil {
		return err
	}
	if err := compareReg8("IFF2", boolToUint8(cpu.iff2), expected.IFF2); err != nil {
		return err
	}
	return compareReg8("IM", cpu.im, expected.IM)
}

// setAltFlags sets shadow flag register from a byte value.
func setAltFlags(flags *Flags, value uint8) {
	flags.C = value & 0x01
	flags.N = (value >> 1) & 0x01
	flags.P = (value >> 2) & 0x01
	flags.X = (value >> 3) & 0x01
	flags.H = (value >> 4) & 0x01
	flags.Y = (value >> 5) & 0x01
	flags.Z = (value >> 6) & 0x01
	flags.S = (value >> 7) & 0x01
}

// getAltFlagsAsUint8 converts shadow flag register to a byte.
func getAltFlagsAsUint8(flags Flags) uint8 {
	return flags.C | (flags.N << 1) | (flags.P << 2) | (flags.X << 3) |
		(flags.H << 4) | (flags.Y << 5) | (flags.Z << 6) | (flags.S << 7)
}

// compareSingleStepRAM compares memory contents.
func compareSingleStepRAM(mem *BasicMemory, expected *singleStepState) error {
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
