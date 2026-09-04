package cpu6502

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/retroenv/retrogolib/assert"
)

const (
	releaseSingleStepRevision = "2f6980a2d95757486c7bee24355c360e40e2a224"
	releaseDormannRevision    = "7954e2dbb49c469ea286070bf46cdd71aeb29e4b"
	releaseCasesPerOpcode     = 10000
)

// TestReleaseNMOS is the fail-closed legal-opcode qualification gate. Unlike
// optional developer integration tests, enabled qualification may not skip data.
func TestReleaseNMOS(t *testing.T) {
	if os.Getenv("CPU6502_QUALIFY") != "1" {
		t.Skip("set CPU6502_QUALIFY=1 for pinned release qualification")
	}
	assert.False(t, testing.Short(), "release qualification cannot use short mode")
	root := os.Getenv("CPU6502_TESTDATA")
	verifyReleaseCorpus(t, root, releaseSingleStepRevision)
	files, cases := 0, 0
	for opcode, info := range Opcodes {
		if info.Instruction.Unofficial {
			continue
		}
		name := fmt.Sprintf("%02x", opcode)
		vectors, err := readReleaseVectors(filepath.Join(root, "6502", "v1", name+".json"), byte(opcode))
		assert.NoError(t, err)
		files++
		cases += len(vectors)
		t.Run(name, func(t *testing.T) {
			for _, vector := range vectors {
				runReleaseVector(t, vector)
			}
		})
	}
	assert.Equal(t, 151, files)
	assert.Equal(t, 151*releaseCasesPerOpcode, cases)
	if !t.Failed() {
		t.Logf("Qualified NMOS legal opcodes: %d files, %d vectors; final state and instruction cycle counts", files, cases)
	}
}

// TestReleaseDormann requires both pinned functional binaries, without skips.
func TestReleaseDormann(t *testing.T) {
	if os.Getenv("CPU6502_QUALIFY") != "1" {
		t.Skip("set CPU6502_QUALIFY=1 for pinned release qualification")
	}
	assert.False(t, testing.Short(), "release qualification cannot use short mode")
	root := os.Getenv("CPU6502_DORMANN_TESTDATA")
	verifyReleaseCorpus(t, root, releaseDormannRevision)
	for _, test := range []dormannTest{
		{name: "NMOS", binary: "bin_files/6502_functional_test.bin", variant: VariantNMOS6502, startPC: dormannStartPC, successPC: nmos6502SuccessPC, maxInstructions: dormannMaxInstructions},
		{name: "65C02-regression", binary: "bin_files/65C02_extended_opcodes_test.bin", variant: Variant65C02, startPC: dormannStartPC, successPC: c65c02SuccessPC, maxInstructions: dormannMaxInstructions},
	} {
		t.Run(test.name, func(t *testing.T) {
			info, err := os.Stat(filepath.Join(root, test.binary))
			assert.NoError(t, err)
			assert.Equal(t, int64(65536), info.Size())
			runDormannTest(t, root, test, true)
		})
	}
}

func TestReadReleaseVectorsRejectsMissingAndEmptyData(t *testing.T) {
	for _, contents := range []string{"", "[]", "[{}]", "not json"} {
		path := filepath.Join(t.TempDir(), "ea.json")
		assert.NoError(t, os.WriteFile(path, []byte(contents), 0o644))
		_, err := readReleaseVectors(path, 0xea)
		assert.Error(t, err)
	}
	_, err := readReleaseVectors(filepath.Join(t.TempDir(), "missing.json"), 0xea)
	assert.Error(t, err)
}

func TestValidateReleaseVector(t *testing.T) {
	valid := ss6502TestCase{Name: "nop", Initial: ss6502State{PC: 0x8000, RAM: [][2]uint32{{0x8000, 0xea}}}, Cycles: [][]any{{32768, 234, "read"}, {32769, 0, "read"}}}
	assert.NoError(t, validateReleaseVector(valid, 0xea))
	assert.Error(t, validateReleaseVector(valid, 0xa9))
	valid.Cycles = nil
	assert.Error(t, validateReleaseVector(valid, 0xea))
	valid.Cycles = [][]any{{}}
	valid.Initial.RAM = nil
	assert.Error(t, validateReleaseVector(valid, 0xea))
}

func verifyReleaseCorpus(t *testing.T, root, revision string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	assert.NotEmpty(t, root, "release test data path is required")
	command := exec.CommandContext(ctx, "git", "-C", root, "rev-parse", "HEAD")
	output, err := command.Output()
	assert.NoError(t, err)
	assert.Equal(t, revision, strings.TrimSpace(string(output)))
	command = exec.CommandContext(ctx, "git", "-C", root, "status", "--porcelain", "--untracked-files=all")
	output, err = command.Output()
	assert.NoError(t, err)
	assert.Empty(t, strings.TrimSpace(string(output)), "test corpus must be unchanged")
}

func readReleaseVectors(path string, opcode byte) ([]ss6502TestCase, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading required opcode vectors: %w", err)
	}
	var vectors []ss6502TestCase
	if err := json.Unmarshal(data, &vectors); err != nil {
		return nil, fmt.Errorf("decoding required opcode vectors: %w", err)
	}
	if len(vectors) != releaseCasesPerOpcode {
		return nil, fmt.Errorf("opcode %02x has %d vectors, require %d", opcode, len(vectors), releaseCasesPerOpcode)
	}
	for _, vector := range vectors {
		if err := validateReleaseVector(vector, opcode); err != nil {
			return nil, err
		}
	}
	return vectors, nil
}

func validateReleaseVector(vector ss6502TestCase, opcode byte) error {
	if len(vector.Cycles) == 0 {
		return fmt.Errorf("vector %q has no cycle trace", vector.Name)
	}
	for _, cell := range vector.Initial.RAM {
		if cell[0] == uint32(vector.Initial.PC) && cell[1] == uint32(opcode) {
			return nil
		}
	}
	return fmt.Errorf("vector %q does not initialize opcode %02x at PC", vector.Name, opcode)
}

func runReleaseVector(t *testing.T, vector ss6502TestCase) {
	t.Helper()
	mem := &ssSparseMemory{data: make(map[uint16]uint8, len(vector.Initial.RAM))}
	for _, cell := range vector.Initial.RAM {
		mem.data[uint16(cell[0])] = uint8(cell[1])
	}
	memory, err := NewMemory(mem)
	assert.NoError(t, err)
	cpu := New(memory, WithVariant(VariantNMOS6502))
	cpu.PC, cpu.SP = vector.Initial.PC, vector.Initial.S
	cpu.A, cpu.X, cpu.Y = vector.Initial.A, vector.Initial.X, vector.Initial.Y
	cpu.setFlags(vector.Initial.P)
	before := cpu.Cycles()
	assert.NoError(t, cpu.Step(), vector.Name)
	assert.Equal(t, uint64(len(vector.Cycles)), cpu.Cycles()-before, vector.Name)
	assert.True(t, verifySingleStepCase(t, vector, cpu, mem), vector.Name)
}
