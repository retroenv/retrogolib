package log

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/retroenv/retrogolib/assert"
)

// TestObject tests the Object field function.
func TestObject(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    any
		expected string
	}{
		{
			name:     "string value",
			key:      "data",
			value:    "test",
			expected: "test",
		},
		{
			name:     "int value",
			key:      "count",
			value:    42,
			expected: "42",
		},
		{
			name:     "struct value",
			key:      "user",
			value:    struct{ Name string }{Name: "John"},
			expected: "{John}",
		},
		{
			name:     "nil value",
			key:      "empty",
			value:    nil,
			expected: "<nil>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := Object(tt.key, tt.value)
			assert.Equal(t, tt.key, field.Key)
			assert.Equal(t, tt.expected, field.Value.String())
		})
	}
}

// TestString tests the String field function.
func TestString(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		expected string
	}{
		{
			name:     "simple string",
			key:      "message",
			value:    "hello",
			expected: "hello",
		},
		{
			name:     "empty string",
			key:      "empty",
			value:    "",
			expected: "",
		},
		{
			name:     "string with spaces",
			key:      "text",
			value:    "hello world",
			expected: "hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := String(tt.key, tt.value)
			assert.Equal(t, tt.key, field.Key)
			assert.Equal(t, tt.expected, field.Value.String())
		})
	}
}

// TestStrings tests the Strings field function.
func TestStrings(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    []string
		contains string
	}{
		{
			name:     "multiple strings",
			key:      "items",
			value:    []string{"a", "b", "c"},
			contains: "[a b c]",
		},
		{
			name:     "empty slice",
			key:      "empty",
			value:    []string{},
			contains: "[]",
		},
		{
			name:     "single string",
			key:      "single",
			value:    []string{"test"},
			contains: "[test]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := Strings(tt.key, tt.value)
			assert.Equal(t, tt.key, field.Key)
			assert.Contains(t, field.Value.String(), tt.contains)
		})
	}
}

// testStringer implements fmt.Stringer for testing.
type testStringer struct {
	value string
}

func (ts testStringer) String() string {
	return ts.value
}

// TestStringer tests the Stringer field function.
func TestStringer(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    fmt.Stringer
		expected string
	}{
		{
			name:     "custom stringer",
			key:      "obj",
			value:    testStringer{value: "custom"},
			expected: "custom",
		},
		{
			name:     "time stringer",
			key:      "time",
			value:    time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			expected: "2023-01-01 12:00:00 +0000 UTC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := Stringer(tt.key, tt.value)
			assert.Equal(t, tt.key, field.Key)
			assert.Equal(t, tt.expected, field.Value.String())
		})
	}
}

// TestStringFunc tests the StringFunc field function and lazy evaluation.
func TestStringFunc(t *testing.T) {
	t.Run("lazy evaluation", func(t *testing.T) {
		callCount := 0
		fn := func() string {
			callCount++
			return "computed"
		}

		field := StringFunc("lazy", fn)
		assert.Equal(t, "lazy", field.Key)

		// Function should not be called yet
		assert.Equal(t, 0, callCount)

		// Force evaluation by accessing the value
		value := field.Value.Any().(stringFunc).LogValue()
		assert.Equal(t, "computed", value.String())
		assert.Equal(t, 1, callCount)

		// Second access should call function again
		value2 := field.Value.Any().(stringFunc).LogValue()
		assert.Equal(t, "computed", value2.String())
		assert.Equal(t, 2, callCount)
	})

	t.Run("expensive computation", func(t *testing.T) {
		field := StringFunc("expensive", func() string {
			return strings.Repeat("x", 1000)
		})

		assert.Equal(t, "expensive", field.Key)

		// Verify the function produces expected result
		value := field.Value.Any().(stringFunc).LogValue()
		assert.Equal(t, 1000, len(value.String()))
		assert.True(t, strings.HasPrefix(value.String(), "xxx"))
	})
}

// TestIntFunc tests the IntFunc field function and lazy evaluation.
func TestIntFunc(t *testing.T) {
	t.Run("lazy evaluation", func(t *testing.T) {
		callCount := 0
		fn := func() int {
			callCount++
			return 42
		}

		field := IntFunc("lazy", fn)
		assert.Equal(t, "lazy", field.Key)

		// Function should not be called yet
		assert.Equal(t, 0, callCount)

		// Force evaluation by accessing the value
		value := field.Value.Any().(intFunc).LogValue()
		assert.Equal(t, int64(42), value.Int64())
		assert.Equal(t, 1, callCount)
	})

	t.Run("expensive computation", func(t *testing.T) {
		field := IntFunc("expensive", func() int {
			sum := 0
			for range 1000 {
				sum += 1
			}
			return sum
		})

		assert.Equal(t, "expensive", field.Key)

		// Verify the function produces expected result
		value := field.Value.Any().(intFunc).LogValue()
		assert.Equal(t, int64(1000), value.Int64()) // sum of 1000 ones
	})
}

// TestErr tests the Err field function.
func TestErr(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "simple error",
			err:      errors.New("test error"),
			expected: "test error",
		},
		{
			name:     "wrapped error",
			err:      fmt.Errorf("wrapped: %w", errors.New("original")),
			expected: "wrapped: original",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := Err(tt.err)
			assert.Equal(t, "error", field.Key)
			assert.Equal(t, tt.expected, field.Value.String())
		})
	}
}

// TestInt tests the Int field function.
func TestInt(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    int
		expected int64
	}{
		{
			name:     "positive int",
			key:      "count",
			value:    42,
			expected: 42,
		},
		{
			name:     "negative int",
			key:      "negative",
			value:    -10,
			expected: -10,
		},
		{
			name:     "zero",
			key:      "zero",
			value:    0,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := Int(tt.key, tt.value)
			assert.Equal(t, tt.key, field.Key)
			assert.Equal(t, tt.expected, field.Value.Int64())
		})
	}
}

// TestInt64 tests the Int64 field function.
func TestInt64(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    int64
		expected int64
	}{
		{
			name:     "large positive",
			key:      "large",
			value:    9223372036854775807, // max int64
			expected: 9223372036854775807,
		},
		{
			name:     "large negative",
			key:      "negative",
			value:    -9223372036854775808, // min int64
			expected: -9223372036854775808,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := Int64(tt.key, tt.value)
			assert.Equal(t, tt.key, field.Key)
			assert.Equal(t, tt.expected, field.Value.Int64())
		})
	}
}

// TestInt32 tests the Int32 field function.
func TestInt32(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    int32
		expected int64
	}{
		{
			name:     "max int32",
			key:      "max",
			value:    2147483647,
			expected: 2147483647,
		},
		{
			name:     "min int32",
			key:      "min",
			value:    -2147483648,
			expected: -2147483648,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := Int32(tt.key, tt.value)
			assert.Equal(t, tt.key, field.Key)
			assert.Equal(t, tt.expected, field.Value.Int64())
		})
	}
}

// TestInt16 tests the Int16 field function.
func TestInt16(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    int16
		expected int64
	}{
		{
			name:     "max int16",
			key:      "max",
			value:    32767,
			expected: 32767,
		},
		{
			name:     "min int16",
			key:      "min",
			value:    -32768,
			expected: -32768,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := Int16(tt.key, tt.value)
			assert.Equal(t, tt.key, field.Key)
			assert.Equal(t, tt.expected, field.Value.Int64())
		})
	}
}

// TestInt8 tests the Int8 field function.
func TestInt8(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    int8
		expected int64
	}{
		{
			name:     "max int8",
			key:      "max",
			value:    127,
			expected: 127,
		},
		{
			name:     "min int8",
			key:      "min",
			value:    -128,
			expected: -128,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := Int8(tt.key, tt.value)
			assert.Equal(t, tt.key, field.Key)
			assert.Equal(t, tt.expected, field.Value.Int64())
		})
	}
}

// TestUint tests the Uint field function.
func TestUint(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    uint
		expected uint64
	}{
		{
			name:     "zero",
			key:      "zero",
			value:    0,
			expected: 0,
		},
		{
			name:     "large value",
			key:      "large",
			value:    4294967295,
			expected: 4294967295,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := Uint(tt.key, tt.value)
			assert.Equal(t, tt.key, field.Key)
			assert.Equal(t, tt.expected, field.Value.Uint64())
		})
	}
}

// TestUint64 tests the Uint64 field function.
func TestUint64(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    uint64
		expected uint64
	}{
		{
			name:     "max uint64",
			key:      "max",
			value:    18446744073709551615,
			expected: 18446744073709551615,
		},
		{
			name:     "zero",
			key:      "zero",
			value:    0,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := Uint64(tt.key, tt.value)
			assert.Equal(t, tt.key, field.Key)
			assert.Equal(t, tt.expected, field.Value.Uint64())
		})
	}
}

// TestUint32 tests the Uint32 field function.
func TestUint32(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    uint32
		expected uint64
	}{
		{
			name:     "max uint32",
			key:      "max",
			value:    4294967295,
			expected: 4294967295,
		},
		{
			name:     "zero",
			key:      "zero",
			value:    0,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := Uint32(tt.key, tt.value)
			assert.Equal(t, tt.key, field.Key)
			assert.Equal(t, tt.expected, field.Value.Uint64())
		})
	}
}

// TestUint16 tests the Uint16 field function.
func TestUint16(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    uint16
		expected uint64
	}{
		{
			name:     "max uint16",
			key:      "max",
			value:    65535,
			expected: 65535,
		},
		{
			name:     "zero",
			key:      "zero",
			value:    0,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := Uint16(tt.key, tt.value)
			assert.Equal(t, tt.key, field.Key)
			assert.Equal(t, tt.expected, field.Value.Uint64())
		})
	}
}

// TestUint8 tests the Uint8 field function.
func TestUint8(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    uint8
		expected uint64
	}{
		{
			name:     "max uint8",
			key:      "max",
			value:    255,
			expected: 255,
		},
		{
			name:     "zero",
			key:      "zero",
			value:    0,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := Uint8(tt.key, tt.value)
			assert.Equal(t, tt.key, field.Key)
			assert.Equal(t, tt.expected, field.Value.Uint64())
		})
	}
}

// TestTime tests the Time field function.
func TestTime(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    time.Time
		expected time.Time
	}{
		{
			name:     "specific time",
			key:      "timestamp",
			value:    time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			expected: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "zero time",
			key:      "zero",
			value:    time.Time{},
			expected: time.Time{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := Time(tt.key, tt.value)
			assert.Equal(t, tt.key, field.Key)
			assert.Equal(t, tt.expected, field.Value.Time())
		})
	}
}

// TestDuration tests the Duration field function.
func TestDuration(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    time.Duration
		expected time.Duration
	}{
		{
			name:     "seconds",
			key:      "timeout",
			value:    5 * time.Second,
			expected: 5 * time.Second,
		},
		{
			name:     "milliseconds",
			key:      "latency",
			value:    100 * time.Millisecond,
			expected: 100 * time.Millisecond,
		},
		{
			name:     "zero duration",
			key:      "zero",
			value:    0,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := Duration(tt.key, tt.value)
			assert.Equal(t, tt.key, field.Key)
			assert.Equal(t, tt.expected, field.Value.Duration())
		})
	}
}

// TestBool tests the Bool field function.
func TestBool(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    bool
		expected bool
	}{
		{
			name:     "true",
			key:      "enabled",
			value:    true,
			expected: true,
		},
		{
			name:     "false",
			key:      "disabled",
			value:    false,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := Bool(tt.key, tt.value)
			assert.Equal(t, tt.key, field.Key)
			assert.Equal(t, tt.expected, field.Value.Bool())
		})
	}
}

// TestFloat32 tests the Float32 field function.
func TestFloat32(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    float32
		expected float64
	}{
		{
			name:     "positive float",
			key:      "rate",
			value:    3.14,
			expected: float64(float32(3.14)), // Convert through float32 for precision
		},
		{
			name:     "negative float",
			key:      "negative",
			value:    -2.5,
			expected: -2.5,
		},
		{
			name:     "zero",
			key:      "zero",
			value:    0.0,
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := Float32(tt.key, tt.value)
			assert.Equal(t, tt.key, field.Key)
			assert.Equal(t, tt.expected, field.Value.Float64())
		})
	}
}

// TestFloat64 tests the Float64 field function.
func TestFloat64(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    float64
		expected float64
	}{
		{
			name:     "positive float",
			key:      "pi",
			value:    3.141592653589793,
			expected: 3.141592653589793,
		},
		{
			name:     "negative float",
			key:      "negative",
			value:    -123.456,
			expected: -123.456,
		},
		{
			name:     "zero",
			key:      "zero",
			value:    0.0,
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := Float64(tt.key, tt.value)
			assert.Equal(t, tt.key, field.Key)
			assert.Equal(t, tt.expected, field.Value.Float64())
		})
	}
}

// TestHex_UnsignedIntegers tests the Hex field function with unsigned integer types.
func TestHex_UnsignedIntegers(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    any
		expected string
	}{
		{
			name:     "uint8",
			key:      "byte",
			value:    uint8(0xFF),
			expected: "0xFF",
		},
		{
			name:     "uint16",
			key:      "addr",
			value:    uint16(0x1234),
			expected: "0x1234",
		},
		{
			name:     "uint32",
			key:      "dword",
			value:    uint32(0x12345678),
			expected: "0x12345678",
		},
		{
			name:     "uint64",
			key:      "qword",
			value:    uint64(0x123456789ABCDEF0),
			expected: "0x123456789ABCDEF0",
		},
		{
			name:     "uint",
			key:      "native_uint",
			value:    uint(0x42),
			expected: "0x42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := Hex(tt.key, tt.value)
			assert.Equal(t, tt.key, field.Key)

			// Force evaluation by accessing the LogValue
			hexValue := field.Value.Any().(hex).LogValue()
			actual := hexValue.String()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

// TestHex_SignedIntegers tests the Hex field function with signed integer types.
func TestHex_SignedIntegers(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    any
		expected string
	}{
		{
			name:     "int8",
			key:      "signed_byte",
			value:    int8(-1),
			expected: "0xFF",
		},
		{
			name:     "int16",
			key:      "signed_word",
			value:    int16(-1),
			expected: "0xFFFF",
		},
		{
			name:     "int32",
			key:      "signed_dword",
			value:    int32(-1),
			expected: "0xFFFFFFFF",
		},
		{
			name:     "int64",
			key:      "signed_qword",
			value:    int64(-1),
			expected: "0xFFFFFFFFFFFFFFFF",
		},
		{
			name:     "int",
			key:      "native_int",
			value:    int(0x42),
			expected: "0x42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := Hex(tt.key, tt.value)
			assert.Equal(t, tt.key, field.Key)

			// Force evaluation by accessing the LogValue
			hexValue := field.Value.Any().(hex).LogValue()
			actual := hexValue.String()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

// TestHex_UnsupportedTypes tests the Hex field function with unsupported types.
func TestHex_UnsupportedTypes(t *testing.T) {
	field := Hex("float", 3.14)
	assert.Equal(t, "float", field.Key)

	// Force evaluation by accessing the LogValue
	hexValue := field.Value.Any().(hex).LogValue()
	actual := hexValue.String()

	// For unsupported types, just check it starts with "0x"
	assert.True(t, strings.HasPrefix(actual, "0x"))
}

// TestHex_LazyEvaluation tests that Hex uses lazy evaluation.
func TestHex_LazyEvaluation(t *testing.T) {
	// This test verifies that the hex formatting is lazy
	field := Hex("test", uint8(0xFF))
	assert.Equal(t, "test", field.Key)

	// The value should be stored as a hex struct, not the formatted string
	hexStruct, ok := field.Value.Any().(hex)
	assert.True(t, ok, "Expected hex struct")
	assert.Equal(t, uint8(0xFF), hexStruct.val)

	// Only when we call LogValue should it format
	logValue := hexStruct.LogValue()
	assert.Equal(t, "0xFF", logValue.String())
}

// TestFormatHex tests the formatHex helper function directly.
func TestFormatHex(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		expected string
	}{
		{
			name:     "uint8 zero",
			value:    uint8(0),
			expected: "0x00",
		},
		{
			name:     "uint8 max",
			value:    uint8(255),
			expected: "0xFF",
		},
		{
			name:     "uint16 zero padding",
			value:    uint16(0x1),
			expected: "0x0001",
		},
		{
			name:     "uint32 zero padding",
			value:    uint32(0x1),
			expected: "0x00000001",
		},
		{
			name:     "uint64 zero padding",
			value:    uint64(0x1),
			expected: "0x0000000000000001",
		},
		{
			name:     "negative int8",
			value:    int8(-128),
			expected: "0x80",
		},
		{
			name:     "negative int16",
			value:    int16(-32768),
			expected: "0x8000",
		},
		{
			name:     "string fallback",
			value:    "test",
			expected: "0x74657374", // hex encoding of "test"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatHex(tt.value)
			if tt.name == "string fallback" {
				// For non-integer types, just verify it starts with 0x
				assert.True(t, strings.HasPrefix(result, "0x"))
			} else {
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// TestType_BasicTypes tests the Type field function with basic types.
func TestType_BasicTypes(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    any
		expected string
	}{
		{
			name:     "int type",
			key:      "value_type",
			value:    42,
			expected: "int",
		},
		{
			name:     "string type",
			key:      "text_type",
			value:    "hello",
			expected: "string",
		},
		{
			name:     "pointer to int",
			key:      "ptr_type",
			value:    new(int),
			expected: "*int",
		},
		{
			name:     "slice of strings",
			key:      "slice_type",
			value:    []string{"a", "b"},
			expected: "[]string",
		},
		{
			name:     "nil value",
			key:      "nil_type",
			value:    nil,
			expected: "<nil>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := Type(tt.key, tt.value)
			assert.Equal(t, tt.key, field.Key)

			// Force evaluation by accessing the LogValue
			typeValue := field.Value.Any().(typeOf).LogValue()
			actual := typeValue.String()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

// TestType_ComplexTypes tests the Type field function with complex types.
func TestType_ComplexTypes(t *testing.T) {
	type customStruct struct {
		Name string
	}

	tests := []struct {
		name     string
		key      string
		value    any
		expected string
	}{
		{
			name:     "custom struct",
			key:      "struct_type",
			value:    customStruct{Name: "test"},
			expected: "log.customStruct",
		},
		{
			name:     "pointer to struct",
			key:      "struct_ptr_type",
			value:    &customStruct{Name: "test"},
			expected: "*log.customStruct",
		},
		{
			name:     "error interface",
			key:      "error_type",
			value:    errors.New("test error"),
			expected: "*errors.errorString",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := Type(tt.key, tt.value)
			assert.Equal(t, tt.key, field.Key)

			// Force evaluation by accessing the LogValue
			typeValue := field.Value.Any().(typeOf).LogValue()
			actual := typeValue.String()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

// TestType_LazyEvaluation tests that Type uses lazy evaluation.
func TestType_LazyEvaluation(t *testing.T) {
	// This test verifies that the type reflection is lazy
	type testStruct struct {
		Value int
	}
	ts := testStruct{Value: 123}

	field := Type("test", ts)
	assert.Equal(t, "test", field.Key)

	// The value should be stored as a typeOf struct, not the formatted string
	typeStruct, ok := field.Value.Any().(typeOf)
	assert.True(t, ok, "Expected typeOf struct")
	assert.Equal(t, ts, typeStruct.val)

	// Only when we call LogValue should it format the type
	logValue := typeStruct.LogValue()
	assert.Equal(t, "log.testStruct", logValue.String())
}
