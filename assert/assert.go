// Package assert contains test assertion helpers.
package assert

import (
	"errors"
	"fmt"
	"reflect"
)

// Testing is an interface that includes the methods used from *testing.T.
type Testing interface {
	Helper()
	Error(args ...any)
	FailNow()
}

// Equal asserts that two objects are equal.
func Equal(t Testing, expected, actual any, msgAndArgs ...any) {
	t.Helper()
	if equal(expected, actual) {
		return
	}

	msg := fmt.Sprintf("Not equal: \nexpected: %v\nactual  : %v", expected, actual)
	fail(t, msg, msgAndArgs...)
}

// NoError asserts that a function returned no error.
func NoError(t Testing, err error, msgAndArgs ...any) {
	t.Helper()
	if err == nil {
		return
	}

	msg := fmt.Sprintf("Unexpected error:\n%+v", err)
	fail(t, msg, msgAndArgs...)
}

// Error asserts that a function returned an error.
func Error(t Testing, err error, expectedError string, msgAndArgs ...any) {
	t.Helper()
	if err == nil {
		msg := fmt.Sprintf("Error message not equal: \nexpected: %v\nactual  : nil", expectedError)
		fail(t, msg, msgAndArgs...)
		return
	}

	actual := err.Error()
	if actual == expectedError {
		return
	}

	msg := fmt.Sprintf("Error message not equal: \nexpected: %v\nactual  : %v", expectedError, actual)
	fail(t, msg, msgAndArgs...)
}

// ErrorIs asserts that a function returned an error that matches the specified error.
func ErrorIs(t Testing, err, expectedError error, msgAndArgs ...any) {
	t.Helper()
	if err == nil {
		msg := fmt.Sprintf("Error not returned: \nexpected: %v\nactual  : nil", expectedError)
		fail(t, msg, msgAndArgs...)
		return
	}

	if errors.Is(err, expectedError) {
		return
	}

	msg := fmt.Sprintf("Error not equal: \nexpected: %v\nactual  : %v", expectedError, err)
	fail(t, msg, msgAndArgs...)
}

// True asserts that the specified value is true.
func True(t Testing, value bool, msgAndArgs ...any) {
	t.Helper()
	if value {
		return
	}
	fail(t, "Unexpected false", msgAndArgs...)
}

// False asserts that the specified value is false.
func False(t Testing, value bool, msgAndArgs ...any) {
	t.Helper()
	if !value {
		return
	}
	fail(t, "Unexpected true", msgAndArgs...)
}

// Len asserts that the specified object has the expected length.
func Len(t Testing, object any, expectedLen int, msgAndArgs ...any) {
	t.Helper()
	actualLen := reflect.ValueOf(object).Len()
	if actualLen == expectedLen {
		return
	}

	msg := fmt.Sprintf("Length not equal: \nexpected: %d\nactual  : %d", expectedLen, actualLen)
	fail(t, msg, msgAndArgs...)
}

// NotNil asserts that the specified object is not nil.
func NotNil(t Testing, object any, msgAndArgs ...any) {
	t.Helper()
	if !isNil(object) {
		return
	}

	msg := "Expected value to be not nil"
	fail(t, msg, msgAndArgs...)
}

// Nil asserts that the specified object is nil.
func Nil(t Testing, object any, msgAndArgs ...any) {
	t.Helper()
	if isNil(object) {
		return
	}

	msg := "Expected value to be nil"
	fail(t, msg, msgAndArgs...)
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

func fail(t Testing, message string, msgAndArgs ...any) {
	t.Helper()
	if len(msgAndArgs) > 0 {
		message += "\n" + fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}
	t.Error(message)
	t.FailNow()
}
