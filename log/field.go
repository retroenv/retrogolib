package log

import (
	"fmt"
	"log/slog"
	"time"
)

// A Field is a marshaling operation used to add a key-value pair to a logger's
// context. Most fields are lazily marshaled, so it's inexpensive to add fields
// to disabled debug-level log statements.
type Field = slog.Attr

// Object constructs a Field with the given key and value.
// It should be used for types that are not represented by a specialized Field
// function. If the passed value type does not implement a custom array or
// object marshaller, reflection will be used for the fields of the type.
// Using reflection for performance critical code paths should be avoided.
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

// stringFunc implements slog.LogValuer for lazy string evaluation.
type stringFunc struct {
	f func() string
}

// LogValue implements slog.LogValuer, ensuring the function is only called
// when the log record is actually processed.
func (sf stringFunc) LogValue() slog.Value {
	return slog.StringValue(sf.f())
}

// StringFunc constructs a Field with the given key and a function that returns a string.
// The function is evaluated lazily - only when the log level is enabled and the
// handler processes the record. This provides significant performance benefits
// for expensive string operations when logging is disabled.
func StringFunc(key string, f func() string) Field {
	return slog.Any(key, stringFunc{f: f})
}

// intFunc implements slog.LogValuer for lazy int evaluation.
type intFunc struct {
	f func() int
}

// LogValue implements slog.LogValuer, ensuring the function is only called
// when the log record is actually processed.
func (inf intFunc) LogValue() slog.Value {
	return slog.IntValue(inf.f())
}

// IntFunc constructs a Field with the given key and a function that returns an int.
// The function is evaluated lazily - only when the log level is enabled and the
// handler processes the record. This provides significant performance benefits
// for expensive int computations when logging is disabled.
func IntFunc(key string, f func() int) Field {
	return slog.Any(key, intFunc{f: f})
}

// Err constructs a Field with the given key and value.
func Err(err error) Field {
	return slog.Any("error", err)
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

// hex implements slog.LogValuer for lazy hex formatting.
type hex struct {
	val any
}

// LogValue implements slog.LogValuer, ensuring the hex formatting is only performed
// when the log record is actually processed.
func (hf hex) LogValue() slog.Value {
	return slog.StringValue(formatHex(hf.val))
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

// formatHex formats integer values as hex strings with appropriate zero-padding.
func formatHex(val any) string {
	switch v := val.(type) {
	case uint8:
		return fmt.Sprintf("0x%02X", v)
	case int8:
		return fmt.Sprintf("0x%02X", uint8(v))
	case uint16:
		return fmt.Sprintf("0x%04X", v)
	case int16:
		return fmt.Sprintf("0x%04X", uint16(v))
	case uint32:
		return fmt.Sprintf("0x%08X", v)
	case int32:
		return fmt.Sprintf("0x%08X", uint32(v))
	case uint64:
		return fmt.Sprintf("0x%016X", v)
	case int64:
		return fmt.Sprintf("0x%016X", uint64(v))
	case uint:
		return fmt.Sprintf("0x%X", v)
	case int:
		return fmt.Sprintf("0x%X", uint(v))
	default:
		return fmt.Sprintf("0x%X", val)
	}
}

// typeOf implements slog.LogValuer for lazy type name formatting.
type typeOf struct {
	val any
}

// LogValue implements slog.LogValuer, ensuring the type reflection is only performed
// when the log record is actually processed.
func (tf typeOf) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("%T", tf.val))
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
