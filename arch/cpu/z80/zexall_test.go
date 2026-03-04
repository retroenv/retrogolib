package z80

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

// TestZexdoc runs the ZEXDOC Z80 instruction exerciser (documented flags only).
func TestZexdoc(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping ZEXDOC in short mode")
	}

	runZex(t, "testdata/zexdoc.com")
}

// TestZexall runs the ZEXALL Z80 instruction exerciser (all flags including undocumented).
func TestZexall(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping ZEXALL in short mode")
	}

	runZex(t, "testdata/zexall.com")
}

// runZex loads and runs a ZEXALL/ZEXDOC .com binary under a minimal CP/M harness.
func runZex(t *testing.T, path string) {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read %s: %v", path, err)
	}

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

	var output bytes.Buffer
	var failCount int

	// Hook BDOS calls at 0x0005
	bdosHook := func(cpu *CPU, opcode uint8, params ...any) {
		if cpu.PC != 0x0005 {
			return
		}

		switch cpu.C {
		case 0x02: // C_WRITE: output character in E
			ch := cpu.E
			output.WriteByte(ch)

			// Flush on newline
			if ch == '\n' {
				line := strings.TrimRight(output.String(), "\r\n")
				if len(line) > 0 {
					if strings.Contains(line, "ERROR") {
						failCount++
					}
					fmt.Println(line)
				}
				output.Reset()
			}

		case 0x09: // C_WRITESTR: output $ terminated string at DE
			addr := cpu.de()
			for {
				ch := mem.Read(addr)
				if ch == '$' {
					break
				}
				output.WriteByte(ch)
				if ch == '\n' {
					line := strings.TrimRight(output.String(), "\r\n")
					if len(line) > 0 {
						if strings.Contains(line, "ERROR") {
							failCount++
						}
						fmt.Println(line)
					}
					output.Reset()
				}
				addr++
			}
		}
	}

	cpu, err := New(mem,
		WithInitialPC(0x0100),
		WithInitialSP(0xFFFE),
		WithPreExecutionHook(bdosHook),
	)
	if err != nil {
		t.Fatalf("failed to create CPU: %v", err)
	}

	// Run until PC reaches 0x0000 (warm boot = exit) or too many cycles
	maxCycles := uint64(200_000_000_000) // ~200 billion cycles
	for cpu.PC != 0x0000 && cpu.cycles < maxCycles {
		pc := cpu.PC
		opByte := mem.Read(pc)
		if err := cpu.Step(); err != nil {
			if output.Len() > 0 {
				fmt.Println(output.String())
			}
			context := make([]byte, 6)
			for i := range context {
				context[i] = mem.Read(pc + uint16(i))
			}
			t.Fatalf("CPU error at PC=0x%04X opcode=0x%02X bytes=%X: %v",
				pc, opByte, context, err)
		}
	}

	// Flush remaining output
	if output.Len() > 0 {
		remaining := strings.TrimRight(output.String(), "\r\n")
		if len(remaining) > 0 {
			fmt.Println(remaining)
			if strings.Contains(remaining, "ERROR") {
				failCount++
			}
		}
	}

	if cpu.cycles >= maxCycles {
		t.Fatal("exceeded maximum cycle count")
	}

	if failCount > 0 {
		t.Fatalf("%d tests failed", failCount)
	}
}
