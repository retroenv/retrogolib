package cpu68000

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

// officialMapAliases maps mnemonic families of the official opcode map to the
// instruction names of this package where the two differ.
var officialMapAliases = map[string]string{
	"ANDItoCCR":   ANDIName,
	"ANDItoSR":    ANDIName,
	"EORItoCCR":   EORIName,
	"EORItoSR":    EORIName,
	"MOVE.q":      MOVEQName,
	"MOVEfromSR":  MOVEName,
	"MOVEfromUSP": MOVEName,
	"MOVEtoCCR":   MOVEName,
	"MOVEtoSR":    MOVEName,
	"MOVEtoUSP":   MOVEName,
	"ORItoCCR":    ORIName,
	"ORItoSR":     ORIName,
	"UNLINK":      UNLKName,
}

// officialMapFamilies lists the instruction names that the official map
// records under a different family: quick, immediate, and memory forms.
var officialMapFamilies = map[string][]string{
	ADDQName: {ADDName},
	SUBQName: {SUBName},
	ADDIName: {ADDName},
	SUBIName: {SUBName},
	ANDIName: {ANDName},
	ORIName:  {ORName},
	EORIName: {EORName},
	CMPIName: {CMPName},
	CMPMName: {CMPName},
	BRAName:  {BccName},
}

func TestDecoderMatchesOfficialOpcodeMap(t *testing.T) {
	// Illegal encodings previously decoded as real instructions, for example
	// JSR Dn jumped to address zero. Every word must agree with the map.
	_, thisFile, _, ok := runtime.Caller(0)
	assert.True(t, ok)
	path := filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "testdata", "cpu68000", "680x0",
		"map", "68000.official.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Skip("official opcode map not found; run 'make -C testdata cpu68000' to download")
	}

	var officialMap map[string]string
	assert.NoError(t, json.Unmarshal(data, &officialMap))
	assert.Len(t, officialMap, 0x10000)

	for word, description := range officialMap {
		opcode, err := strconv.ParseUint(word, 16, 16)
		assert.NoError(t, err)
		decoded := decodeOpcode(uint16(opcode))

		family, _, _ := strings.Cut(description, " ")
		if family == "None" {
			assert.Equal(t, insILLEGAL, decoded.Instruction, "opcode %s decodes as %s, want ILLEGAL", word, decoded.Instruction.Name)
			continue
		}

		want := officialMapFamily(family)
		got := decoded.Instruction.Name
		if got != want && !familyMatches(got, want) {
			t.Errorf("opcode %s (%s) decodes as %s, want %s", word, description, got, want)
		}
	}
}

// officialMapFamily converts a family of the official map to an instruction name.
func officialMapFamily(family string) string {
	if alias, ok := officialMapAliases[family]; ok {
		return alias
	}
	name, _, _ := strings.Cut(family, ".")
	return name
}

// familyMatches reports whether the decoded name is a form that the official
// map records under the given family.
func familyMatches(got, family string) bool {
	for _, candidate := range officialMapFamilies[got] {
		if candidate == family {
			return true
		}
	}
	return false
}
