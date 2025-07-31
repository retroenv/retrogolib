package config

import (
	"errors"
	"fmt"
)

// Error definitions
var (
	ErrTypeMismatch       = errors.New("type mismatch")
	ErrWriteOnly          = errors.New("configuration was not loaded from file")
	ErrInvalidStruct      = errors.New("invalid struct type")
	ErrUnsupportedType    = errors.New("unsupported field type")
	ErrRequiredField      = errors.New("required field is missing")
	ErrDuplicateSection   = errors.New("duplicate section")
	ErrDuplicateKey       = errors.New("duplicate key")
	ErrConfigTooLarge     = errors.New("configuration file too large")
	ErrTooManyLines       = errors.New("configuration file has too many lines")
	ErrSectionNameTooLong = errors.New("section name too long")
	ErrKeyNameTooLong     = errors.New("key name too long")
)

// MarshalError represents errors during struct marshaling.
type MarshalError struct {
	Field   string
	Section string
	Key     string
	Err     error
}

func (e *MarshalError) Error() string {
	return fmt.Sprintf("marshal field %s (%s.%s): %v", e.Field, e.Section, e.Key, e.Err)
}

func (e *MarshalError) Unwrap() error {
	return e.Err
}

// UnmarshalError represents errors during struct unmarshalling.
type UnmarshalError struct {
	Field   string
	Section string
	Key     string
	Err     error
}

func (e *UnmarshalError) Error() string {
	return fmt.Sprintf("unmarshal field %s (%s.%s): %v", e.Field, e.Section, e.Key, e.Err)
}

func (e *UnmarshalError) Unwrap() error {
	return e.Err
}

// ParseError represents parsing errors with location information.
type ParseError struct {
	Line int
	Pos  int
	Msg  string
	Err  error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("line %d, pos %d: %s", e.Line, e.Pos, e.Msg)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}
