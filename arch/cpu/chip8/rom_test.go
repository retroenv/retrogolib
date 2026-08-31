package chip8

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

const (
	romCyclesPerTick = 20
	romMaxCycles     = 100_000
)

func TestROMConformance(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping CHIP-8 ROM tests in short mode")
	}

	testDataDir := chip8TestDataDir(t)
	tests := []struct {
		name       string
		rom        string
		screenshot string
		configure  func(*CPU)
		run        func(*testing.T, *CPU)
	}{
		{name: "CHIP-8 logo", rom: "1-chip8-logo.ch8", screenshot: "chip-8-logo.png"},
		{name: "IBM logo", rom: "2-ibm-logo.ch8", screenshot: "ibm-logo.png"},
		{name: "opcode test", rom: "3-corax+.ch8", screenshot: "corax+.png"},
		{name: "flags test", rom: "4-flags.ch8", screenshot: "flags.png"},
		{
			name:       "quirks test",
			rom:        "5-quirks.ch8",
			screenshot: "quirks.png",
			configure: func(c *CPU) {
				c.Memory[0x1ff] = 1
			},
		},
		{
			name:       "keypad test",
			rom:        "6-keypad.ch8",
			screenshot: "keypad-getkey.png",
			configure: func(c *CPU) {
				c.Memory[0x1ff] = 3
			},
			run: runKeypadROM,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			loadROM(t, c, filepath.Join(testDataDir, "bin", tt.rom))
			if tt.configure != nil {
				tt.configure(c)
			}

			if tt.run == nil {
				runROMToTerminal(t, c)
			} else {
				tt.run(t, c)
			}
			want := screenshotDisplay(t, filepath.Join(testDataDir, "pictures", tt.screenshot))
			assertDisplayMatches(t, want, c.Display)
		})
	}
}

func chip8TestDataDir(t *testing.T) string {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	assert.True(t, ok, "resolve CHIP-8 test location")
	dir := filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "testdata", "chip8")
	_, err := os.Stat(dir)
	if errors.Is(err, os.ErrNotExist) {
		t.Skipf("CHIP-8 test data not found at %s (run 'make -C testdata chip8' to download)", dir)
	}
	assert.NoError(t, err)

	return dir
}

func loadROM(t *testing.T, c *CPU, path string) {
	t.Helper()

	rom, err := os.ReadFile(path)
	assert.NoError(t, err)
	assert.True(t, len(rom) <= len(c.Memory)-initialProgramCounter, "ROM exceeds CHIP-8 program memory")
	copy(c.Memory[initialProgramCounter:], rom)
}

func runROMToTerminal(t *testing.T, c *CPU) {
	t.Helper()

	for cycle := range romMaxCycles {
		if cycle > 0 && cycle%romCyclesPerTick == 0 {
			c.UpdateTimers()
		}
		if c.PC < uint16(len(c.Memory)-1) {
			opcode := uint16(c.Memory[c.PC])<<8 | uint16(c.Memory[c.PC+1])
			if opcode&0xf000 == 0x1000 && opcode&0x0fff == c.PC {
				return
			}
		}

		assert.NoError(t, c.Step(), "cycle %d at PC 0x%03X", cycle, c.PC)
		if c.keyWait.active {
			return
		}
	}

	assert.Fail(t, "ROM did not reach its terminal loop", "cycles: %d, PC: 0x%03X", romMaxCycles, c.PC)
}

func runKeypadROM(t *testing.T, c *CPU) {
	t.Helper()

	runROMToTerminal(t, c)
	assert.True(t, c.keyWait.active, "keypad ROM did not wait for input")
	for range 3 {
		c.UpdateTimers()
	}

	c.Key[5] = true
	assert.NoError(t, c.Step())
	assert.True(t, c.keyWait.active, "FX0A resumed before key release")
	c.Key[5] = false
	assert.NoError(t, c.Step())
	assert.False(t, c.keyWait.active)

	runROMToTerminal(t, c)
}

func assertDisplayMatches(t *testing.T, want, got [displayWidth * displayHeight]byte) {
	t.Helper()

	var firstDifferences []string
	differenceCount := 0
	for index := range want {
		if want[index] == got[index] {
			continue
		}
		differenceCount++
		if len(firstDifferences) < 10 {
			firstDifferences = append(firstDifferences, fmt.Sprintf("(%d,%d): want %d, got %d", index%displayWidth, index/displayWidth, want[index], got[index]))
		}
	}

	assert.Equal(t, 0, differenceCount, "first differences: %v", firstDifferences)
}

func screenshotDisplay(t *testing.T, path string) [displayWidth * displayHeight]byte {
	t.Helper()

	data, err := os.ReadFile(path)
	assert.NoError(t, err)
	screenshot, err := png.Decode(bytes.NewReader(data))
	assert.NoError(t, err)

	bounds := screenshot.Bounds()
	assert.Equal(t, 0, bounds.Dx()%displayWidth)
	assert.Equal(t, 0, bounds.Dy()%displayHeight)
	scaleX := bounds.Dx() / displayWidth
	scaleY := bounds.Dy() / displayHeight
	assert.True(t, scaleX > 0 && scaleY > 0, "reference screenshot is smaller than the CHIP-8 display")

	minBrightness := ^uint64(0)
	var maxBrightness uint64
	for y := range displayHeight {
		for x := range displayWidth {
			brightness := pixelBrightness(screenshot, bounds.Min.X+x*scaleX+scaleX/2, bounds.Min.Y+y*scaleY+scaleY/2)
			minBrightness = min(minBrightness, brightness)
			maxBrightness = max(maxBrightness, brightness)
		}
	}
	threshold := minBrightness + (maxBrightness-minBrightness)/2

	var display [displayWidth * displayHeight]byte
	for y := range displayHeight {
		for x := range displayWidth {
			brightness := pixelBrightness(screenshot, bounds.Min.X+x*scaleX+scaleX/2, bounds.Min.Y+y*scaleY+scaleY/2)
			if brightness > threshold {
				display[y*displayWidth+x] = 1
			}
		}
	}

	return display
}

func pixelBrightness(img image.Image, x, y int) uint64 {
	r, g, b, _ := img.At(x, y).RGBA()
	return uint64(r) + uint64(g) + uint64(b)
}
