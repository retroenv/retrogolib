package log

import (
	"fmt"
	"time"

	"golang.org/x/exp/slog"
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
