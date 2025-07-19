// Package assert contains test assertion helpers.
package assert

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
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

// NotEqual asserts that two objects are not equal.
func NotEqual(t Testing, expected, actual any, msgAndArgs ...any) {
	t.Helper()
	if !equal(expected, actual) {
		return
	}

	msg := fmt.Sprintf("Equal: \nexpected: %v\nactual  : %v", expected, actual)
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

	msg := fmt.Sprintf("Expected value to be nil, got: %v", object)
	fail(t, msg, msgAndArgs...)
}

// Contains asserts that the string contains the substring.
func Contains(t Testing, s, substr string, msgAndArgs ...any) {
	t.Helper()
	if strings.Contains(s, substr) {
		return
	}

	msg := fmt.Sprintf("String does not contain substring:\nstring: %s\nsubstring: %s", s, substr)
	fail(t, msg, msgAndArgs...)
}

// NotContains asserts that the string does not contain the substring.
func NotContains(t Testing, s, substr string, msgAndArgs ...any) {
	t.Helper()
	if !strings.Contains(s, substr) {
		return
	}

	msg := fmt.Sprintf("String contains substring:\nstring: %s\nsubstring: %s", s, substr)
	fail(t, msg, msgAndArgs...)
}

// Panics asserts that the function panics when called.
func Panics(t Testing, fn func(), msgAndArgs ...any) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			fail(t, "Function did not panic", msgAndArgs...)
		}
	}()
	fn()
}

// NotPanics asserts that the function does not panic when called.
func NotPanics(t Testing, fn func(), msgAndArgs ...any) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("Function panicked with: %v", r)
			fail(t, msg, msgAndArgs...)
		}
	}()
	fn()
}

// Empty asserts that the object is empty.
func Empty(t Testing, object any, msgAndArgs ...any) {
	t.Helper()
	if isEmpty(object) {
		return
	}

	msg := fmt.Sprintf("Expected empty, but got: %v", object)
	fail(t, msg, msgAndArgs...)
}

// NotEmpty asserts that the object is not empty.
func NotEmpty(t Testing, object any, msgAndArgs ...any) {
	t.Helper()
	if !isEmpty(object) {
		return
	}

	msg := "Expected not empty, but got empty"
	fail(t, msg, msgAndArgs...)
}

// Greater asserts that the first value is greater than the second.
func Greater(t Testing, first, second any, msgAndArgs ...any) {
	t.Helper()
	if isGreater(first, second) {
		return
	}

	msg := fmt.Sprintf("Expected greater:\nfirst: %v\nsecond: %v", first, second)
	fail(t, msg, msgAndArgs...)
}

// GreaterOrEqual asserts that the first value is greater than or equal to the second.
func GreaterOrEqual(t Testing, first, second any, msgAndArgs ...any) {
	t.Helper()
	if isGreater(first, second) || equal(first, second) {
		return
	}

	msg := fmt.Sprintf("Expected greater or equal:\nfirst: %v\nsecond: %v", first, second)
	fail(t, msg, msgAndArgs...)
}

// Less asserts that the first value is less than the second.
func Less(t Testing, first, second any, msgAndArgs ...any) {
	t.Helper()
	if isLess(first, second) {
		return
	}

	msg := fmt.Sprintf("Expected less:\nfirst: %v\nsecond: %v", first, second)
	fail(t, msg, msgAndArgs...)
}

// LessOrEqual asserts that the first value is less than or equal to the second.
func LessOrEqual(t Testing, first, second any, msgAndArgs ...any) {
	t.Helper()
	if isLess(first, second) || equal(first, second) {
		return
	}

	msg := fmt.Sprintf("Expected less or equal:\nfirst: %v\nsecond: %v", first, second)
	fail(t, msg, msgAndArgs...)
}

// ErrorContains asserts that the error message contains the substring.
func ErrorContains(t Testing, err error, substr string, msgAndArgs ...any) {
	t.Helper()
	if err == nil {
		msg := fmt.Sprintf("Expected error containing: %s\nActual: nil", substr)
		fail(t, msg, msgAndArgs...)
		return
	}

	if strings.Contains(err.Error(), substr) {
		return
	}

	msg := fmt.Sprintf("Error does not contain substring:\nerror: %v\nsubstring: %s", err, substr)
	fail(t, msg, msgAndArgs...)
}

// Implements asserts that the object implements the interface.
func Implements(t Testing, interfaceType, object any, msgAndArgs ...any) {
	t.Helper()
	interfacePtr := reflect.TypeOf(interfaceType).Elem()

	if !reflect.TypeOf(object).Implements(interfacePtr) {
		msg := fmt.Sprintf("%T does not implement %v", object, interfacePtr)
		fail(t, msg, msgAndArgs...)
	}
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

func isEmpty(value any) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		return v.Len() == 0
	default:
		return false
	}
}

func isGreater(first, second any) bool {
	fv := reflect.ValueOf(first)
	sv := reflect.ValueOf(second)

	switch fv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fv.Int() > sv.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fv.Uint() > sv.Uint()
	case reflect.Float32, reflect.Float64:
		return fv.Float() > sv.Float()
	case reflect.String:
		return fv.String() > sv.String()
	default:
		return false
	}
}

func isLess(first, second any) bool {
	fv := reflect.ValueOf(first)
	sv := reflect.ValueOf(second)

	switch fv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fv.Int() < sv.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fv.Uint() < sv.Uint()
	case reflect.Float32, reflect.Float64:
		return fv.Float() < sv.Float()
	case reflect.String:
		return fv.String() < sv.String()
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
