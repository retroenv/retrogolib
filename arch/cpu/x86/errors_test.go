package x86

import (
	"errors"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

var errorTestCases = []struct {
	name   string
	err    error
	errMsg string
}{
	{
		name:   "ErrNilMemory",
		err:    ErrNilMemory,
		errMsg: "memory is nil",
	},
	{
		name:   "ErrInvalidInstruction",
		err:    ErrInvalidInstruction,
		errMsg: "invalid instruction",
	},
	{
		name:   "ErrInvalidAddressingMode",
		err:    ErrInvalidAddressingMode,
		errMsg: "invalid addressing mode",
	},
	{
		name:   "ErrInvalidRegister",
		err:    ErrInvalidRegister,
		errMsg: "invalid register",
	},
	{
		name:   "ErrInvalidOpcode",
		err:    ErrInvalidOpcode,
		errMsg: "invalid opcode",
	},
	{
		name:   "ErrInvalidOperand",
		err:    ErrInvalidOperand,
		errMsg: "invalid operand",
	},
	{
		name:   "ErrDivisionByZero",
		err:    ErrDivisionByZero,
		errMsg: "division by zero",
	},
	{
		name:   "ErrInvalidSegment",
		err:    ErrInvalidSegment,
		errMsg: "invalid segment",
	},
	{
		name:   "ErrStackOverflow",
		err:    ErrStackOverflow,
		errMsg: "stack overflow",
	},
	{
		name:   "ErrStackUnderflow",
		err:    ErrStackUnderflow,
		errMsg: "stack underflow",
	},
	{
		name:   "ErrGeneralProtectionFault",
		err:    ErrGeneralProtectionFault,
		errMsg: "general protection fault",
	},
	{
		name:   "ErrInvalidInterruptVector",
		err:    ErrInvalidInterruptVector,
		errMsg: "invalid interrupt vector",
	},
}

func TestErrors(t *testing.T) {
	for _, tt := range errorTestCases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.errMsg, tt.err.Error())
		})
	}
}

func TestErrorIs(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		target   error
		expected bool
	}{
		{
			name:     "same error",
			err:      ErrNilMemory,
			target:   ErrNilMemory,
			expected: true,
		},
		{
			name:     "different error",
			err:      ErrNilMemory,
			target:   ErrInvalidInstruction,
			expected: false,
		},
		{
			name:     "wrapped error",
			err:      errors.New("wrapper: invalid instruction"),
			target:   ErrInvalidInstruction,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := errors.Is(tt.err, tt.target)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestErrorUnwrap(t *testing.T) {
	// Test that our errors are simple and don't wrap other errors
	testErrors := []error{
		ErrNilMemory,
		ErrInvalidInstruction,
		ErrInvalidAddressingMode,
		ErrInvalidRegister,
		ErrInvalidOpcode,
		ErrInvalidOperand,
		ErrDivisionByZero,
		ErrInvalidSegment,
		ErrStackOverflow,
		ErrStackUnderflow,
		ErrGeneralProtectionFault,
		ErrInvalidInterruptVector,
	}

	for _, err := range testErrors {
		unwrapped := errors.Unwrap(err)
		assert.Nil(t, unwrapped) // Our errors don't wrap other errors
	}
}
