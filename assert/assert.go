// Package assert contains test assertion helpers.
package assert

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

// Equal asserts that two objects are equal.
func Equal(t *testing.T, expected, actual any, errorMessage ...string) {
	t.Helper()
	if equal(expected, actual) {
		return
	}

	msg := fmt.Sprintf("Not equal: \nexpected: %v\nactual  : %v", expected, actual)
	fail(t, msg, errorMessage...)
}

// NoError asserts that a function returned no error.
func NoError(t *testing.T, err error, errorMessage ...string) {
	t.Helper()
	if err == nil {
		return
	}

	msg := fmt.Sprintf("Unexpected error:\n%+v", err)
	fail(t, msg, errorMessage...)
}

// Error asserts that a function returned an error.
func Error(t *testing.T, err error, expectedError string, errorMessage ...string) {
	t.Helper()
	if err == nil {
		msg := fmt.Sprintf("Error message not equal: \nexpected: %v\nactual  : nil", expectedError)
		fail(t, msg, errorMessage...)
		return
	}

	actual := err.Error()
	if actual == expectedError {
		return
	}

	msg := fmt.Sprintf("Error message not equal: \nexpected: %v\nactual  : %v", expectedError, actual)
	fail(t, msg, errorMessage...)
}

// ErrorIs asserts that a function returned an error that matches the specified error.
func ErrorIs(t *testing.T, err, expectedError error, errorMessage ...string) {
	t.Helper()
	if err == nil {
		msg := fmt.Sprintf("Error not returned: \nexpected: %v\nactual  : nil", expectedError)
		fail(t, msg, errorMessage...)
		return
	}

	if errors.Is(err, expectedError) {
		return
	}

	msg := fmt.Sprintf("Error not equal: \nexpected: %v\nactual  : %v", expectedError, err)
	fail(t, msg, errorMessage...)
}

// True asserts that the specified value is true.
func True(t *testing.T, value bool, errorMessage ...string) {
	t.Helper()
	if value {
		return
	}
	fail(t, "Unexpected false", errorMessage...)
}

// False asserts that the specified value is false.
func False(t *testing.T, value bool, errorMessage ...string) {
	t.Helper()
	if !value {
		return
	}
	fail(t, "Unexpected true", errorMessage...)
}

// Len asserts that the specified object has the expected length.
func Len(t *testing.T, object any, expectedLen int, errorMessage ...string) {
	t.Helper()
	actualLen := reflect.ValueOf(object).Len()
	if actualLen == expectedLen {
		return
	}

	msg := fmt.Sprintf("Length not equal: \nexpected: %d\nactual  : %d", expectedLen, actualLen)
	fail(t, msg, errorMessage...)
}

// NotNil asserts that the specified object is not nil.
func NotNil(t *testing.T, object any, errorMessage ...string) {
	t.Helper()
	if !isNil(object) {
		return
	}

	msg := "Expected value to be not nil"
	fail(t, msg, errorMessage...)
}

// Nil asserts that the specified object is nil.
func Nil(t *testing.T, object any, errorMessage ...string) {
	t.Helper()
	if isNil(object) {
		return
	}

	msg := "Expected value to be nil"
	fail(t, msg, errorMessage...)
}

func equal(expected, actual any) bool {
	if expected == nil || actual == nil {
		return isNil(expected) == isNil(actual)
	}

	if reflect.DeepEqual(expected, actual) {
		return true
	}

	actualType := reflect.TypeOf(actual)
	if actualType == nil {
		return false
	}
	expectedValue := reflect.ValueOf(expected)
	if expectedValue.IsValid() && expectedValue.Type().ConvertibleTo(actualType) {
		return reflect.DeepEqual(expectedValue.Convert(actualType).Interface(), actual)
	}

	return false
}

func isNil(value any) bool {
	if value == nil {
		return true
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(value).IsNil()
	default:
		return false
	}
}

func fail(t *testing.T, message string, errorMessage ...string) {
	t.Helper()
	if len(errorMessage) != 0 {
		message = fmt.Sprintf("%s\n%s", message, errorMessage)
	}
	t.Error(message)
	t.FailNow()
}
