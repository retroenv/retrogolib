package z80

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

// TestZexdoc runs the ZEXDOC Z80 instruction exerciser (documented flags only).
func TestZexdoc(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping ZEXDOC in short mode")
	}

	dir := getZexallDir(t)
	runZex(t, filepath.Join(dir, "zexdoc.com"), 67)
}

// TestZexall runs the ZEXALL Z80 instruction exerciser (all flags including undocumented).
func TestZexall(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping ZEXALL in short mode")
	}

	dir := getZexallDir(t)
	runZex(t, filepath.Join(dir, "zexall.com"), 67)
}

// zexOutput tracks output buffering and error counting for ZEX tests.
type zexOutput struct {
	buf       bytes.Buffer
	failCount int
	okCount   int
}

// flushLine processes a completed line of output, checking for errors.
func (z *zexOutput) flushLine() {
	line := strings.TrimRight(z.buf.String(), "\r\n")
	if len(line) > 0 {
		if strings.Contains(line, "ERROR") {
			z.failCount++
		}
		if strings.HasSuffix(line, "OK") {
			z.okCount++
		}
		fmt.Println(line)
	}
	z.buf.Reset()
}

// writeByte writes a byte to the buffer, flushing on newline.
func (z *zexOutput) writeByte(ch byte) {
	z.buf.WriteByte(ch)
	if ch == '\n' {
		z.flushLine()
	}
}

// flush processes any remaining buffered output.
func (z *zexOutput) flush() {
	if z.buf.Len() > 0 {
		z.flushLine()
	}
}

// getZexallDir returns the path to the ZEXALL test data directory,
// skipping the test if it is not found.
func getZexallDir(t *testing.T) string {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	assert.True(t, ok)

	dir := filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "testdata", "zexall")
	if _, err := os.Stat(dir); err != nil {
		t.Skipf("ZEXALL data not found at %s (run 'make -C testdata zexall' to download)", dir)
	}

	return dir
}

// handleBDOS processes CP/M BDOS calls intercepted at address 0x0005.
func handleBDOS(cpu *CPU, mem *BasicMemory, out *zexOutput) {
	switch cpu.C {
	case 0x02: // C_WRITE: output character in E
		out.writeByte(cpu.E)

	case 0x09: // C_WRITESTR: output $ terminated string at DE
		addr := cpu.de()
		for {
			ch := mem.Read(addr)
			if ch == '$' {
				break
			}
			out.writeByte(ch)
			addr++
		}
	}
}

// runZex loads and runs a ZEXALL/ZEXDOC .com binary under a minimal CP/M harness.
func runZex(t *testing.T, path string, expectedTests int) {
	t.Helper()

	data, err := os.ReadFile(path)
	assert.NoError(t, err, "failed to read %s", path)

	mem := NewBasicMemory()

	// CP/M loads .com files at 0x0100
	copy(mem.data[0x0100:], data)

	// CP/M memory setup
	// 0x0000: JP 0x0000 (warm boot - we detect PC=0 as exit)
	mem.data[0x0000] = 0xC3 // JP
	mem.data[0x0001] = 0x00
	mem.data[0x0002] = 0x00

	// 0x0005: RET (BDOS entry - we intercept via preExecutionHook)
	mem.data[0x0005] = 0xC9

	// 0x0006-0x0007: BDOS entry address / top of TPA
	// Programs use LD SP,(0006h) to set stack below BDOS
	mem.data[0x0006] = 0x00
	mem.data[0x0007] = 0xFE // TPA top at 0xFE00

	var out zexOutput

	// Hook BDOS calls at 0x0005
	bdosHook := func(cpu *CPU, _ uint8, _ ...any) {
		if cpu.PC == 0x0005 {
			handleBDOS(cpu, mem, &out)
		}
	}

	cpu, err := New(mem,
		WithInitialPC(0x0100),
		WithInitialSP(0xFFFE),
		WithPreExecutionHook(bdosHook),
	)
	assert.NoError(t, err, "failed to create CPU")

	// Run until PC reaches 0x0000 (warm boot = exit) or too many cycles
	maxCycles := uint64(200_000_000_000) // ~200 billion cycles
	for cpu.PC != 0x0000 && cpu.cycles < maxCycles {
		pc := cpu.PC
		opByte := mem.Read(pc)
		if err := cpu.Step(); err != nil {
			out.flush()
			context := make([]byte, 6)
			for i := range context {
				context[i] = mem.Read(pc + uint16(i))
			}
			t.Fatalf("CPU error at PC=0x%04X opcode=0x%02X bytes=%X: %v",
				pc, opByte, context, err)
		}
	}

	// Flush remaining output
	out.flush()

	assert.Less(t, cpu.cycles, maxCycles, "exceeded maximum cycle count")
	assert.Equal(t, 0, out.failCount, "%d test groups reported ERROR", out.failCount)
	assert.Equal(t, expectedTests, out.okCount, "expected %d OK results, got %d", expectedTests, out.okCount)
}
