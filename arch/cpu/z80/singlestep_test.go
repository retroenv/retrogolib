package z80

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

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
	AF_  uint16   `json:"af_"`  //nolint:revive,stylecheck // matches JSON field name
	BC_  uint16   `json:"bc_"`  //nolint:revive,stylecheck // matches JSON field name
	DE_  uint16   `json:"de_"`  //nolint:revive,stylecheck // matches JSON field name
	HL_  uint16   `json:"hl_"`  //nolint:revive,stylecheck // matches JSON field name
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

// UnmarshalJSON handles the ports field which is an array of [addr, val, "r"|"w"].
func (t *singleStepTest) UnmarshalJSON(data []byte) error {
	type alias singleStepTest
	var raw struct {
		alias
		RawPorts []json.RawMessage `json:"ports"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
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
		json.Unmarshal(arr[0], &addr)
		json.Unmarshal(arr[1], &val)
		json.Unmarshal(arr[2], &dir)
		t.Ports = append(t.Ports, singleStepPort{
			Address: uint16(addr),
			Value:   uint8(val),
			IsRead:  dir == "r",
		})
	}
	return nil
}

// testIOHandler provides port read/write values for test cases.
type testIOHandler struct {
	reads  map[uint8]uint8
	writes map[uint8]uint8
}

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

func (h *testIOHandler) ReadPort(port uint8) uint8 {
	if val, ok := h.reads[port]; ok {
		return val
	}
	return 0xFF
}

func (h *testIOHandler) WritePort(_ uint8, _ uint8) {}

const singleStepDir = "testdata/singlestep"

// TestSingleStep runs the SingleStepTests Z80 test suite.
// Each JSON file contains 1000 test cases that verify single-instruction execution
// against known-correct hardware traces.
func TestSingleStep(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping SingleStepTests in short mode")
	}

	cloneSingleStepTests(t)

	files, err := filepath.Glob(filepath.Join(singleStepDir, "v1", "*.json"))
	if err != nil {
		t.Fatalf("globbing test files: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("no SingleStepTests JSON files found")
	}

	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			t.Parallel()
			runSingleStepFile(t, file)
		})
	}
}

// cloneSingleStepTests clones the SingleStepTests z80 repo if not already present.
func cloneSingleStepTests(t *testing.T) {
	t.Helper()

	if _, err := os.Stat(filepath.Join(singleStepDir, "v1")); err == nil {
		return
	}

	t.Log("cloning SingleStepTests/z80 repository...")
	cmd := exec.Command("git", "clone", "--depth=1",
		"https://github.com/SingleStepTests/z80", singleStepDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("cloning SingleStepTests: %v", err)
	}
}

// runSingleStepFile runs all test cases from a single JSON file.
func runSingleStepFile(t *testing.T, path string) {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}

	var tests []singleStepTest
	if err := json.Unmarshal(data, &tests); err != nil {
		t.Fatalf("parsing %s: %v", path, err)
	}

	var passed, failed int

	for i := range tests {
		tc := &tests[i]
		if err := runSingleStepCase(tc); err != nil {
			t.Errorf("%s: %s", tc.Name, err)
			failed++
		} else {
			passed++
		}
	}

	if failed > 0 {
		t.Logf("%s: %d passed, %d failed", filepath.Base(path), passed, failed)
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
	cpu.setFlagsFromUint8(&cpu.AltFlags, uint8(s.AF_))
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

// compareSingleStepState compares the CPU state against expected final state.
// Returns an error describing the first mismatch found.
func compareSingleStepState(cpu *CPU, mem *BasicMemory, expected *singleStepState) error {
	// Main registers.
	if cpu.A != expected.A {
		return fmt.Errorf("A: got 0x%02X, want 0x%02X", cpu.A, expected.A)
	}
	if cpu.B != expected.B {
		return fmt.Errorf("B: got 0x%02X, want 0x%02X", cpu.B, expected.B)
	}
	if cpu.C != expected.C {
		return fmt.Errorf("C: got 0x%02X, want 0x%02X", cpu.C, expected.C)
	}
	if cpu.D != expected.D {
		return fmt.Errorf("D: got 0x%02X, want 0x%02X", cpu.D, expected.D)
	}
	if cpu.E != expected.E {
		return fmt.Errorf("E: got 0x%02X, want 0x%02X", cpu.E, expected.E)
	}
	if cpu.H != expected.H {
		return fmt.Errorf("H: got 0x%02X, want 0x%02X", cpu.H, expected.H)
	}
	if cpu.L != expected.L {
		return fmt.Errorf("L: got 0x%02X, want 0x%02X", cpu.L, expected.L)
	}

	// Flags.
	gotF := cpu.GetFlags()
	if gotF != expected.F {
		return fmt.Errorf("F: got 0x%02X, want 0x%02X", gotF, expected.F)
	}

	// Program control.
	if cpu.PC != expected.PC {
		return fmt.Errorf("PC: got 0x%04X, want 0x%04X", cpu.PC, expected.PC)
	}
	if cpu.SP != expected.SP {
		return fmt.Errorf("SP: got 0x%04X, want 0x%04X", cpu.SP, expected.SP)
	}

	// Index registers.
	if cpu.IX != expected.IX {
		return fmt.Errorf("IX: got 0x%04X, want 0x%04X", cpu.IX, expected.IX)
	}
	if cpu.IY != expected.IY {
		return fmt.Errorf("IY: got 0x%04X, want 0x%04X", cpu.IY, expected.IY)
	}

	// Special registers.
	if cpu.I != expected.I {
		return fmt.Errorf("I: got 0x%02X, want 0x%02X", cpu.I, expected.I)
	}
	if cpu.R != expected.R {
		return fmt.Errorf("R: got 0x%02X, want 0x%02X", cpu.R, expected.R)
	}
	if cpu.MEMPTR != expected.WZ {
		return fmt.Errorf("MEMPTR: got 0x%04X, want 0x%04X", cpu.MEMPTR, expected.WZ)
	}

	// Shadow registers.
	gotAF_ := uint16(cpu.AltA)<<8 | uint16(cpu.getFlagsAsUint8(cpu.AltFlags))
	if gotAF_ != expected.AF_ {
		return fmt.Errorf("AF': got 0x%04X, want 0x%04X", gotAF_, expected.AF_)
	}
	gotBC_ := uint16(cpu.AltB)<<8 | uint16(cpu.AltC)
	if gotBC_ != expected.BC_ {
		return fmt.Errorf("BC': got 0x%04X, want 0x%04X", gotBC_, expected.BC_)
	}
	gotDE_ := uint16(cpu.AltD)<<8 | uint16(cpu.AltE)
	if gotDE_ != expected.DE_ {
		return fmt.Errorf("DE': got 0x%04X, want 0x%04X", gotDE_, expected.DE_)
	}
	gotHL_ := uint16(cpu.AltH)<<8 | uint16(cpu.AltL)
	if gotHL_ != expected.HL_ {
		return fmt.Errorf("HL': got 0x%04X, want 0x%04X", gotHL_, expected.HL_)
	}

	// Interrupt state.
	gotIFF1 := boolToUint8(cpu.iff1)
	if gotIFF1 != expected.IFF1 {
		return fmt.Errorf("IFF1: got %d, want %d", gotIFF1, expected.IFF1)
	}
	gotIFF2 := boolToUint8(cpu.iff2)
	if gotIFF2 != expected.IFF2 {
		return fmt.Errorf("IFF2: got %d, want %d", gotIFF2, expected.IFF2)
	}
	if cpu.im != expected.IM {
		return fmt.Errorf("IM: got %d, want %d", cpu.im, expected.IM)
	}

	// RAM.
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
