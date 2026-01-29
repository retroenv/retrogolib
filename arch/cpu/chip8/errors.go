package chip8

import "errors"

// Common errors for Chip-8 emulation
var (
	ErrRegisterOutOfBounds  = errors.New("register index out of bounds")
	ErrKeyIndexOutOfBounds  = errors.New("key index out of bounds")
	ErrMemoryOutOfBounds    = errors.New("memory access out of bounds")
	ErrStackOverflow        = errors.New("stack overflow")
	ErrStackUnderflow       = errors.New("stack underflow")
	ErrDisplayOutOfBounds   = errors.New("display index out of bounds")
	ErrFontIndexOutOfBounds = errors.New("font index out of bounds")
)
