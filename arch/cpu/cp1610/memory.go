package cp1610

import "reflect"

// BasicMemory is the word-addressed bus used by the CP1610.
type BasicMemory interface {
	Read(address uint16) uint16
	Write(address, value uint16)
}

// Memory wraps the caller's word-addressed bus.
type Memory struct {
	BasicMemory
}

// NewMemory creates a memory wrapper.
func NewMemory(bus BasicMemory) (*Memory, error) {
	if bus == nil {
		return nil, ErrNilMemory
	}
	value := reflect.ValueOf(bus)
	if value.Kind() == reflect.Pointer && value.IsNil() {
		return nil, ErrNilMemory
	}
	return &Memory{BasicMemory: bus}, nil
}
