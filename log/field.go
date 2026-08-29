package log

import (
	"fmt"
	"log/slog"
	"reflect"
	"time"
)

const upperHexDigits = "0123456789ABCDEF"

// Field is a structured logging attribute. Fields backed by slog.LogValuer are
// resolved only after a handler accepts the record.
type Field = slog.Attr

// Object constructs a Field for values without a specialized constructor.
// Prefer typed constructors on hot paths to avoid handler-side reflection.
func Object(key string, val any) Field {
	return slog.Any(key, val)
}

// String constructs a string Field.
func String(key, val string) Field {
	return slog.String(key, val)
}

// Strings constructs a string-slice Field.
func Strings(key string, val []string) Field {
	return slog.Any(key, val)
}

// Stringer constructs a lazily resolved fmt.Stringer Field.
func Stringer(key string, val fmt.Stringer) Field {
	return lazyField(key, func() slog.Value {
		return stringerValue(val)
	})
}

// StringFunc constructs a lazily evaluated string Field.
// The function runs only when a handler processes an enabled record.
func StringFunc(key string, f func() string) Field {
	return lazyField(key, func() slog.Value {
		return slog.StringValue(f())
	})
}

// IntFunc constructs a lazily evaluated int Field.
// The function runs only when a handler processes an enabled record.
func IntFunc(key string, f func() int) Field {
	return lazyField(key, func() slog.Value {
		return slog.IntValue(f())
	})
}

// Int64Func constructs a lazily evaluated int64 Field.
// The function runs only when a handler processes an enabled record.
func Int64Func(key string, f func() int64) Field {
	return lazyField(key, func() slog.Value {
		return slog.Int64Value(f())
	})
}

// Float64Func constructs a lazily evaluated float64 Field.
// The function runs only when a handler processes an enabled record.
func Float64Func(key string, f func() float64) Field {
	return lazyField(key, func() slog.Value {
		return slog.Float64Value(f())
	})
}

// BoolFunc constructs a lazily evaluated bool Field.
// The function runs only when a handler processes an enabled record.
func BoolFunc(key string, f func() bool) Field {
	return lazyField(key, func() slog.Value {
		return slog.BoolValue(f())
	})
}

// DurationFunc constructs a lazily evaluated duration Field.
// The function runs only when a handler processes an enabled record.
func DurationFunc(key string, f func() time.Duration) Field {
	return lazyField(key, func() slog.Value {
		return slog.DurationValue(f())
	})
}

// StringerFunc constructs a lazily evaluated fmt.Stringer Field.
// The function runs only when a handler processes an enabled record.
func StringerFunc(key string, f func() fmt.Stringer) Field {
	return lazyField(key, func() slog.Value {
		return stringerValue(f())
	})
}

// Err constructs an error Field with key "error".
func Err(err error) Field {
	return slog.Any("error", err)
}

// Int constructs an int Field.
func Int(key string, val int) Field {
	return slog.Int(key, val)
}

// Int64 constructs an int64 Field.
func Int64(key string, val int64) Field {
	return slog.Int64(key, val)
}

// Int32 constructs an int32 Field.
func Int32(key string, val int32) Field {
	return slog.Int64(key, int64(val))
}

// Int16 constructs an int16 Field.
func Int16(key string, val int16) Field {
	return slog.Int64(key, int64(val))
}

// Int8 constructs an int8 Field.
func Int8(key string, val int8) Field {
	return slog.Int64(key, int64(val))
}

// Uint constructs a uint Field.
func Uint(key string, val uint) Field {
	return slog.Uint64(key, uint64(val))
}

// Uint64 constructs a uint64 Field.
func Uint64(key string, val uint64) Field {
	return slog.Uint64(key, val)
}

// Uint32 constructs a uint32 Field.
func Uint32(key string, val uint32) Field {
	return slog.Uint64(key, uint64(val))
}

// Uint16 constructs a uint16 Field.
func Uint16(key string, val uint16) Field {
	return slog.Uint64(key, uint64(val))
}

// Uint8 constructs a uint8 Field.
func Uint8(key string, val uint8) Field {
	return slog.Uint64(key, uint64(val))
}

// Time constructs a time Field.
func Time(key string, val time.Time) Field {
	return slog.Time(key, val)
}

// Duration constructs a duration Field.
func Duration(key string, val time.Duration) Field {
	return slog.Duration(key, val)
}

// Bool constructs a bool Field.
func Bool(key string, val bool) Field {
	return slog.Bool(key, val)
}

// Float32 constructs a float32 Field.
func Float32(key string, val float32) Field {
	return slog.Float64(key, float64(val))
}

// Float64 constructs a float64 Field.
func Float64(key string, val float64) Field {
	return slog.Float64(key, val)
}

// Hex constructs a lazily formatted hexadecimal Field.
// Fixed-width integer types are zero-padded to their full width.
//
//	log.Hex("addr", uint16(0x1234)) // "0x1234"
//	log.Hex("byte", uint8(0xFF))    // "0xFF"
//	log.Hex("opcode", 0x4C)         // "0x4C"
func Hex(key string, val any) Field {
	return lazyField(key, func() slog.Value {
		return slog.StringValue(formatHex(val))
	})
}

// Type constructs a lazily resolved type-name Field.
//
//	log.Type("addr_type", typedInstr.Addr) // "*nes.IndirectX"
//	log.Type("value_type", value)          // "int"
//	log.Type("handler_type", handler)      // "*http.Handler"
func Type(key string, val any) Field {
	return lazyField(key, func() slog.Value {
		typ := reflect.TypeOf(val)
		if typ == nil {
			return slog.StringValue("<nil>")
		}

		return slog.StringValue(typ.String())
	})
}

type lazyValue func() slog.Value

func (f lazyValue) LogValue() slog.Value {
	return f()
}

func lazyField(key string, f lazyValue) Field {
	return slog.Any(key, f)
}

func stringerValue(val fmt.Stringer) slog.Value {
	if isNil(val) {
		return slog.StringValue("<nil>")
	}

	return slog.StringValue(val.String())
}

func isNil(val any) bool {
	if val == nil {
		return true
	}

	value := reflect.ValueOf(val)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}

func formatHex(val any) string {
	switch value := val.(type) {
	case uint8:
		return formatHexUint(uint64(value), 2)
	case int8:
		return formatHexUint(uint64(uint8(value)), 2)
	case uint16:
		return formatHexUint(uint64(value), 4)
	case int16:
		return formatHexUint(uint64(uint16(value)), 4)
	case uint32:
		return formatHexUint(uint64(value), 8)
	case int32:
		return formatHexUint(uint64(uint32(value)), 8)
	case uint64:
		return formatHexUint(value, 16)
	case int64:
		return formatHexUint(uint64(value), 16)
	case uint:
		return formatHexUint(uint64(value), 0)
	case int:
		return formatHexUint(uint64(uint(value)), 0)
	default:
		return fmt.Sprintf("0x%X", val)
	}
}

func formatHexUint(val uint64, width int) string {
	if width == 0 {
		width = minHexDigits(val)
	}

	var buf [18]byte
	start := len(buf) - width - 2
	buf[start] = '0'
	buf[start+1] = 'x'

	for i := width - 1; i >= 0; i-- {
		buf[start+2+i] = upperHexDigits[val&0xF]
		val >>= 4
	}

	return string(buf[start:])
}

func minHexDigits(val uint64) int {
	digits := 1
	for val >= 16 {
		digits++
		val >>= 4
	}

	return digits
}
