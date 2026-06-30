package log

import (
	"fmt"
	"log/slog"
	"reflect"
	"time"
)

const upperHexDigits = "0123456789ABCDEF"

// A Field is a marshaling operation used to add a key-value pair to a logger's
// context. Most fields are lazily marshaled, so it's inexpensive to add fields
// to disabled debug-level log statements.
type Field = slog.Attr

// Object constructs a Field for values without a specialized constructor.
// Prefer typed constructors on hot paths to avoid handler-side reflection.
func Object(key string, val any) Field {
	return slog.Any(key, val)
}

// String constructs a Field with the given key and value.
func String(key, val string) Field {
	return slog.String(key, val)
}

// Strings constructs a Field with the given key and value.
func Strings(key string, val []string) Field {
	return slog.Any(key, val)
}

// Stringer constructs a Field with the given key and value.
func Stringer(key string, val fmt.Stringer) Field {
	return slog.Any(key, val)
}

// StringFunc constructs a Field with the given key and a function that returns a string.
// The function is evaluated lazily - only when the log level is enabled and the
// handler processes the record. This provides significant performance benefits
// for expensive string operations when logging is disabled.
func StringFunc(key string, f func() string) Field {
	return slog.Any(key, stringFunc{f: f})
}

// IntFunc constructs a Field with the given key and a function that returns an int.
// The function is evaluated lazily - only when the log level is enabled and the
// handler processes the record. This provides significant performance benefits
// for expensive int computations when logging is disabled.
func IntFunc(key string, f func() int) Field {
	return slog.Any(key, intFunc{f: f})
}

// Int64Func constructs a Field with the given key and a function that returns an int64.
func Int64Func(key string, f func() int64) Field {
	return slog.Any(key, int64Func{f: f})
}

// Float64Func constructs a Field with the given key and a function that returns a float64.
func Float64Func(key string, f func() float64) Field {
	return slog.Any(key, float64Func{f: f})
}

// BoolFunc constructs a Field with the given key and a function that returns a bool.
func BoolFunc(key string, f func() bool) Field {
	return slog.Any(key, boolFunc{f: f})
}

// DurationFunc constructs a Field with the given key and a function that returns a duration.
func DurationFunc(key string, f func() time.Duration) Field {
	return slog.Any(key, durationFunc{f: f})
}

// StringerFunc constructs a Field with the given key and a function that returns a Stringer.
func StringerFunc(key string, f func() fmt.Stringer) Field {
	return slog.Any(key, stringerFunc{f: f})
}

// Err constructs an error Field with key "error".
func Err(err error) Field {
	if err == nil {
		return slog.String("error", "<nil>")
	}

	return slog.String("error", err.Error())
}

// Int constructs a Field with the given key and value.
func Int(key string, val int) Field {
	return slog.Int(key, val)
}

// Int64 constructs a Field with the given key and value.
func Int64(key string, val int64) Field {
	return slog.Int64(key, val)
}

// Int32 constructs a Field with the given key and value.
func Int32(key string, val int32) Field {
	return slog.Int64(key, int64(val))
}

// Int16 constructs a Field with the given key and value.
func Int16(key string, val int16) Field {
	return slog.Int64(key, int64(val))
}

// Int8 constructs a Field with the given key and value.
func Int8(key string, val int8) Field {
	return slog.Int64(key, int64(val))
}

// Uint constructs a Field with the given key and value.
func Uint(key string, val uint) Field {
	return slog.Uint64(key, uint64(val))
}

// Uint64 constructs a Field with the given key and value.
func Uint64(key string, val uint64) Field {
	return slog.Uint64(key, val)
}

// Uint32 constructs a Field with the given key and value.
func Uint32(key string, val uint32) Field {
	return slog.Uint64(key, uint64(val))
}

// Uint16 constructs a Field with the given key and value.
func Uint16(key string, val uint16) Field {
	return slog.Uint64(key, uint64(val))
}

// Uint8 constructs a Field with the given key and value.
func Uint8(key string, val uint8) Field {
	return slog.Uint64(key, uint64(val))
}

// Time constructs a Field with the given key and value.
func Time(key string, val time.Time) Field {
	return slog.Time(key, val)
}

// Duration constructs a Field with the given key and value.
func Duration(key string, val time.Duration) Field {
	return slog.Duration(key, val)
}

// Bool constructs a Field with the given key and value.
func Bool(key string, val bool) Field {
	return slog.Bool(key, val)
}

// Float32 constructs a Field with the given key and value.
func Float32(key string, val float32) Field {
	return slog.Float64(key, float64(val))
}

// Float64 constructs a Field with the given key and value.
func Float64(key string, val float64) Field {
	return slog.Float64(key, val)
}

// Hex constructs a Field with the given key and formats integer values in hex format.
// The hex formatting is evaluated lazily - only when the log level is enabled and the
// handler processes the record. This provides significant performance benefits for
// expensive hex formatting when logging is disabled.
//
// Supports signed and unsigned integers of various bit widths with appropriate zero-padding.
//
// Examples:
//
//	log.Hex("addr", uint16(0x1234))  // "addr": "0x1234"
//	log.Hex("byte", uint8(0xFF))     // "byte": "0xFF"
//	log.Hex("opcode", 0x4C)          // "opcode": "0x4C"
func Hex(key string, val any) Field {
	return slog.Any(key, hex{val: val})
}

// Type constructs a Field with the given key and formats the value's type name.
// The type reflection is evaluated lazily - only when the log level is enabled and the
// handler processes the record. This provides significant performance benefits by
// avoiding reflection overhead when logging is disabled.
//
// Examples:
//
//	log.Type("addr_type", typedInstr.Addr)  // "addr_type": "*nes.IndirectX"
//	log.Type("value_type", myVar)           // "value_type": "int"
//	log.Type("handler_type", handler)       // "handler_type": "*http.Handler"
func Type(key string, val any) Field {
	return slog.Any(key, typeOf{val: val})
}

type (
	stringFunc struct {
		f func() string
	}

	intFunc struct {
		f func() int
	}

	int64Func struct {
		f func() int64
	}

	float64Func struct {
		f func() float64
	}

	boolFunc struct {
		f func() bool
	}

	durationFunc struct {
		f func() time.Duration
	}

	stringerFunc struct {
		f func() fmt.Stringer
	}

	hex struct {
		val any
	}

	typeOf struct {
		val any
	}
)

func (sf stringFunc) LogValue() slog.Value {
	return slog.StringValue(sf.f())
}

func (inf intFunc) LogValue() slog.Value {
	return slog.IntValue(inf.f())
}

func (inf int64Func) LogValue() slog.Value {
	return slog.Int64Value(inf.f())
}

func (ff float64Func) LogValue() slog.Value {
	return slog.Float64Value(ff.f())
}

func (bf boolFunc) LogValue() slog.Value {
	return slog.BoolValue(bf.f())
}

func (df durationFunc) LogValue() slog.Value {
	return slog.DurationValue(df.f())
}

func (sf stringerFunc) LogValue() slog.Value {
	val := sf.f()
	if val == nil {
		return slog.StringValue("<nil>")
	}

	return slog.StringValue(val.String())
}

func (hf hex) LogValue() slog.Value {
	return slog.StringValue(formatHex(hf.val))
}

func (tf typeOf) LogValue() slog.Value {
	typ := reflect.TypeOf(tf.val)
	if typ == nil {
		return slog.StringValue("<nil>")
	}

	return slog.StringValue(typ.String())
}

// formatHex formats integer values as hex strings with appropriate zero-padding.
func formatHex(val any) string {
	switch v := val.(type) {
	case uint8:
		return formatHexUint(uint64(v), 2)
	case int8:
		return formatHexUint(uint64(uint8(v)), 2)
	case uint16:
		return formatHexUint(uint64(v), 4)
	case int16:
		return formatHexUint(uint64(uint16(v)), 4)
	case uint32:
		return formatHexUint(uint64(v), 8)
	case int32:
		return formatHexUint(uint64(uint32(v)), 8)
	case uint64:
		return formatHexUint(v, 16)
	case int64:
		return formatHexUint(uint64(v), 16)
	case uint:
		return formatHexUint(uint64(v), 0)
	case int:
		return formatHexUint(uint64(uint(v)), 0)
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
