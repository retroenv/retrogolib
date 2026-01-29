package input_test

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/input"
)

func TestKey_Constants(t *testing.T) {
	// Test that key constants have expected values
	assert.Equal(t, input.Key(0), input.Unknown, "Unknown should be 0")
	assert.Equal(t, input.Key(1), input.Space, "Space should be 1")
	assert.True(t, input.Last > input.Menu, "Last should be greater than Menu")
}

func TestKey_AlphabeticKeys(t *testing.T) {
	// Test alphabetic key ordering
	assert.True(t, input.A < input.B, "A should be less than B")
	assert.True(t, input.B < input.C, "B should be less than C")
	assert.True(t, input.Y < input.Z, "Y should be less than Z")
}

func TestKey_NumericKeys(t *testing.T) {
	// Test numeric key ordering
	assert.True(t, input.Key0 < input.Key1, "Key0 should be less than Key1")
	assert.True(t, input.Key1 < input.Key2, "Key1 should be less than Key2")
	assert.True(t, input.Key8 < input.Key9, "Key8 should be less than Key9")
}

func TestKey_FunctionKeys(t *testing.T) {
	// Test function key ordering
	assert.True(t, input.F1 < input.F2, "F1 should be less than F2")
	assert.True(t, input.F2 < input.F3, "F2 should be less than F3")
	assert.True(t, input.F24 < input.F25, "F24 should be less than F25")
}

func TestKey_KeypadKeys(t *testing.T) {
	// Test keypad key ordering
	assert.True(t, input.KP0 < input.KP1, "KP0 should be less than KP1")
	assert.True(t, input.KP8 < input.KP9, "KP8 should be less than KP9")
	assert.True(t, input.KP9 < input.KPDecimal, "KP9 should be less than KPDecimal")
}

func TestKey_ModifierKeys(t *testing.T) {
	// Test modifier key grouping
	leftModifiers := []input.Key{
		input.LeftShift, input.LeftControl, input.LeftAlt, input.LeftSuper,
	}
	rightModifiers := []input.Key{
		input.RightShift, input.RightControl, input.RightAlt, input.RightSuper,
	}

	// Verify left modifiers are sequential
	for i := 1; i < len(leftModifiers); i++ {
		assert.True(t, leftModifiers[i-1] < leftModifiers[i],
			"Left modifier keys should be in sequential order")
	}

	// Verify right modifiers are sequential
	for i := 1; i < len(rightModifiers); i++ {
		assert.True(t, rightModifiers[i-1] < rightModifiers[i],
			"Right modifier keys should be in sequential order")
	}
}

func TestKey_ArrowKeys(t *testing.T) {
	// Test arrow key grouping
	arrowKeys := []input.Key{input.Right, input.Left, input.Down, input.Up}

	// All arrow keys should be defined
	for _, key := range arrowKeys {
		assert.True(t, key > input.Unknown, "Arrow key should be greater than Unknown")
		assert.True(t, key < input.Last, "Arrow key should be less than Last")
	}
}

func TestKey_NavigationKeys(t *testing.T) {
	// Test navigation key grouping
	navKeys := []input.Key{
		input.Home, input.End, input.PageUp, input.PageDown,
		input.Insert, input.Delete,
	}

	// All navigation keys should be defined
	for _, key := range navKeys {
		assert.True(t, key > input.Unknown, "Navigation key should be greater than Unknown")
		assert.True(t, key < input.Last, "Navigation key should be less than Last")
	}
}

func TestKey_LockKeys(t *testing.T) {
	// Test lock key grouping
	lockKeys := []input.Key{
		input.CapsLock, input.ScrollLock, input.NumLock,
	}

	// All lock keys should be defined
	for _, key := range lockKeys {
		assert.True(t, key > input.Unknown, "Lock key should be greater than Unknown")
		assert.True(t, key < input.Last, "Lock key should be less than Last")
	}
}

func TestKey_PunctuationKeys(t *testing.T) {
	// Test punctuation key grouping
	punctKeys := []input.Key{
		input.Apostrophe, input.Comma, input.Minus, input.Period, input.Slash,
		input.Semicolon, input.Equal, input.LeftBracket, input.Backslash, input.RightBracket,
	}

	// All punctuation keys should be defined
	for _, key := range punctKeys {
		assert.True(t, key > input.Unknown, "Punctuation key should be greater than Unknown")
		assert.True(t, key < input.Last, "Punctuation key should be less than Last")
	}
}

func TestKey_SpecialKeys(t *testing.T) {
	// Test special keys
	specialKeys := []input.Key{
		input.Escape, input.Enter, input.Tab, input.Backspace,
		input.PrintScreen, input.Pause, input.Menu,
	}

	// All special keys should be defined
	for _, key := range specialKeys {
		assert.True(t, key > input.Unknown, "Special key should be greater than Unknown")
		assert.True(t, key < input.Last, "Special key should be less than Last")
	}
}

func TestKey_KeypadOperators(t *testing.T) {
	// Test keypad operator keys
	kpOperators := []input.Key{
		input.KPDecimal, input.KPDivide, input.KPMultiply,
		input.KPSubtract, input.KPAdd, input.KPEnter, input.KPEqual,
	}

	// All keypad operators should be defined
	for _, key := range kpOperators {
		assert.True(t, key > input.Unknown, "Keypad operator should be greater than Unknown")
		assert.True(t, key < input.Last, "Keypad operator should be less than Last")
	}
}

func TestKey_Ranges(t *testing.T) {
	// Test that all keys are within valid range
	assert.True(t, input.Unknown >= 0, "Unknown should be non-negative")
	assert.True(t, input.Last > 0, "Last should be positive")

	// Test that all defined keys are unique
	keys := []input.Key{
		input.Space, input.A, input.Z, input.Key0, input.Key9,
		input.F1, input.F25, input.KP0, input.KP9,
		input.LeftShift, input.RightSuper, input.Menu,
	}

	for i := range keys {
		for j := i + 1; j < len(keys); j++ {
			assert.NotEqual(t, keys[i], keys[j], "Keys should be unique")
		}
	}
}

func TestKey_String(t *testing.T) {
	// Test that keys can be converted to strings (implicit test)
	// This verifies the Key type works as expected with Go's type system
	key := input.A

	// Should be able to use in maps, comparisons, etc.
	keyMap := make(map[input.Key]bool)
	keyMap[key] = true
	assert.True(t, keyMap[input.A], "Key should work as map key")

	// Should be able to compare
	assert.Equal(t, input.A, key, "Key comparison should work")
	assert.NotEqual(t, input.B, key, "Key comparison should work")
}
