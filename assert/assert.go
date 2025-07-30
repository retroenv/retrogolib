// Package assert contains test assertion helpers.
package assert

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Testing is an interface that includes the methods used from *testing.T.
type Testing interface {
	Helper()
	Error(args ...any)
	FailNow()
}

// Fail fails the test with a message and optional format arguments.
func Fail(t Testing, message string, msgAndArgs ...any) {
	t.Helper()
	if len(msgAndArgs) > 0 {
		var builder strings.Builder
		builder.WriteString(message)
		builder.WriteByte('\n')
		builder.WriteString(fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...))
		message = builder.String()
	}
	t.Error(message)
	t.FailNow()
}

// Equal asserts that two objects are equal.
//
// Example:
//
//	assert.Equal(t, 42, result)
//	assert.Equal(t, "hello", greeting, "greeting should be hello")
func Equal(t Testing, expected, actual any, msgAndArgs ...any) {
	t.Helper()
	if equal(expected, actual) {
		return
	}

	msg := fmt.Sprintf("Not equal: \nexpected: %v\nactual  : %v", expected, actual)
	Fail(t, msg, msgAndArgs...)
}

// NotEqual asserts that two objects are not equal.
//
// Example:
//
//	assert.NotEqual(t, 0, count)
//	assert.NotEqual(t, oldValue, newValue, "value should have changed")
func NotEqual(t Testing, expected, actual any, msgAndArgs ...any) {
	t.Helper()
	if !equal(expected, actual) {
		return
	}

	msg := fmt.Sprintf("Equal: \nexpected: %v\nactual  : %v", expected, actual)
	Fail(t, msg, msgAndArgs...)
}

// NoError asserts that a function returned no error.
//
// Example:
//
//	err := processData()
//	assert.NoError(t, err)
//	assert.NoError(t, err, "data processing should succeed")
func NoError(t Testing, err error, msgAndArgs ...any) {
	t.Helper()
	if err == nil {
		return
	}

	msg := fmt.Sprintf("Unexpected error:\n%+v", err)
	Fail(t, msg, msgAndArgs...)
}

// Error asserts that a function returned an error.
//
// Example:
//
//	err := divide(1, 0)
//	assert.Error(t, err)
//	assert.Error(t, err, "division by zero should fail")
func Error(t Testing, err error, msgAndArgs ...any) {
	t.Helper()
	if err != nil {
		return
	}

	msg := "Expected an error"
	Fail(t, msg, msgAndArgs...)
}

// ErrorIs asserts that a function returned an error that matches the specified error.
// Uses errors.Is for comparison, which supports error wrapping.
//
// Example:
//
//	err := processFile("missing.txt")
//	assert.ErrorIs(t, err, os.ErrNotExist)
//	assert.ErrorIs(t, err, ErrInvalidInput, "should be input validation error")
func ErrorIs(t Testing, err, expectedError error, msgAndArgs ...any) {
	t.Helper()
	if err == nil {
		msg := fmt.Sprintf("Error not returned: \nexpected: %v\nactual  : nil", expectedError)
		Fail(t, msg, msgAndArgs...)
		return
	}

	if errors.Is(err, expectedError) {
		return
	}

	msg := fmt.Sprintf("Error not equal: \nexpected: %v\nactual  : %v", expectedError, err)
	Fail(t, msg, msgAndArgs...)
}

// ErrorAs asserts that the error matches the expected error type using errors.As.
//
// Example:
//
//	var pathErr *fs.PathError
//	assert.ErrorAs(t, err, &pathErr)
func ErrorAs(t Testing, err error, target any, msgAndArgs ...any) {
	t.Helper()
	if err == nil {
		msg := fmt.Sprintf("Error not returned: \nexpected: error that matches %T\nactual  : nil", target)
		Fail(t, msg, msgAndArgs...)
		return
	}

	if errors.As(err, target) {
		return
	}

	msg := fmt.Sprintf("Error type mismatch: \nexpected: error that matches %T\nactual  : %v", target, err)
	Fail(t, msg, msgAndArgs...)
}

// True asserts that the specified value is true.
//
// Example:
//
//	assert.True(t, isValid)
//	assert.True(t, user.IsActive(), "user should be active")
func True(t Testing, value bool, msgAndArgs ...any) {
	t.Helper()
	if value {
		return
	}
	Fail(t, "Unexpected false", msgAndArgs...)
}

// False asserts that the specified value is false.
//
// Example:
//
//	assert.False(t, isEmpty)
//	assert.False(t, user.IsBlocked(), "user should not be blocked")
func False(t Testing, value bool, msgAndArgs ...any) {
	t.Helper()
	if !value {
		return
	}
	Fail(t, "Unexpected true", msgAndArgs...)
}

// Len asserts that the specified object has the expected length.
//
// Example:
//
//	assert.Len(t, items, 5)
//	assert.Len(t, "hello", 5, "string should have 5 characters")
func Len(t Testing, object any, expectedLen int, msgAndArgs ...any) {
	t.Helper()
	v := reflect.ValueOf(object)
	if !v.IsValid() {
		Fail(t, "Cannot get length of nil", msgAndArgs...)
		return
	}

	switch v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		actualLen := v.Len()
		if actualLen == expectedLen {
			return
		}
		msg := "Length not equal: \nexpected: " + strconv.Itoa(expectedLen) + "\nactual  : " + strconv.Itoa(actualLen)
		Fail(t, msg, msgAndArgs...)
	default:
		Fail(t, fmt.Sprintf("Object of type %T does not have a length", object), msgAndArgs...)
	}
}

// NotNil asserts that the specified object is not nil.
//
// Example:
//
//	assert.NotNil(t, user)
//	assert.NotNil(t, response, "response should not be nil")
func NotNil(t Testing, object any, msgAndArgs ...any) {
	t.Helper()
	if !isNil(object) {
		return
	}

	msg := "Expected value to be not nil"
	Fail(t, msg, msgAndArgs...)
}

// Nil asserts that the specified object is nil.
//
// Example:
//
//	assert.Nil(t, ptr)
//	assert.Nil(t, result, "result should be nil for invalid input")
func Nil(t Testing, object any, msgAndArgs ...any) {
	t.Helper()
	if isNil(object) {
		return
	}

	msg := fmt.Sprintf("Expected value to be nil, got: %v", object)
	Fail(t, msg, msgAndArgs...)
}

// Contains asserts that the string contains the substring.
func Contains(t Testing, s, substr string, msgAndArgs ...any) {
	t.Helper()
	if strings.Contains(s, substr) {
		return
	}

	msg := fmt.Sprintf("String does not contain substring:\nstring: %s\nsubstring: %s", s, substr)
	Fail(t, msg, msgAndArgs...)
}

// NotContains asserts that the string does not contain the substring.
func NotContains(t Testing, s, substr string, msgAndArgs ...any) {
	t.Helper()
	if !strings.Contains(s, substr) {
		return
	}

	msg := fmt.Sprintf("String contains substring:\nstring: %s\nsubstring: %s", s, substr)
	Fail(t, msg, msgAndArgs...)
}

// Panics asserts that the function panics when called.
func Panics(t Testing, fn func(), msgAndArgs ...any) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			Fail(t, "Function did not panic", msgAndArgs...)
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
			Fail(t, msg, msgAndArgs...)
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
	Fail(t, msg, msgAndArgs...)
}

// NotEmpty asserts that the object is not empty.
func NotEmpty(t Testing, object any, msgAndArgs ...any) {
	t.Helper()
	if !isEmpty(object) {
		return
	}

	msg := "Expected not empty, but got empty"
	Fail(t, msg, msgAndArgs...)
}

// Greater asserts that the first value is greater than the second.
func Greater(t Testing, first, second any, msgAndArgs ...any) {
	t.Helper()
	if isGreater(first, second) {
		return
	}

	msg := fmt.Sprintf("Expected greater:\nfirst: %v\nsecond: %v", first, second)
	Fail(t, msg, msgAndArgs...)
}

// GreaterOrEqual asserts that the first value is greater than or equal to the second.
func GreaterOrEqual(t Testing, first, second any, msgAndArgs ...any) {
	t.Helper()
	if isGreater(first, second) || equal(first, second) {
		return
	}

	msg := fmt.Sprintf("Expected greater or equal:\nfirst: %v\nsecond: %v", first, second)
	Fail(t, msg, msgAndArgs...)
}

// Less asserts that the first value is less than the second.
func Less(t Testing, first, second any, msgAndArgs ...any) {
	t.Helper()
	if isLess(first, second) {
		return
	}

	msg := fmt.Sprintf("Expected less:\nfirst: %v\nsecond: %v", first, second)
	Fail(t, msg, msgAndArgs...)
}

// LessOrEqual asserts that the first value is less than or equal to the second.
func LessOrEqual(t Testing, first, second any, msgAndArgs ...any) {
	t.Helper()
	if isLess(first, second) || equal(first, second) {
		return
	}

	msg := fmt.Sprintf("Expected less or equal:\nfirst: %v\nsecond: %v", first, second)
	Fail(t, msg, msgAndArgs...)
}

// ErrorContains asserts that the error message contains the substring.
//
// Example:
//
//	err := authenticate("invalid-token")
//	assert.ErrorContains(t, err, "permission denied")
//	assert.ErrorContains(t, err, "authentication", "should be auth error")
func ErrorContains(t Testing, err error, substr string, msgAndArgs ...any) {
	t.Helper()
	if err == nil {
		msg := fmt.Sprintf("Expected error containing: %s\nActual: nil", substr)
		Fail(t, msg, msgAndArgs...)
		return
	}

	if strings.Contains(err.Error(), substr) {
		return
	}

	msg := fmt.Sprintf("Error does not contain substring:\nerror: %v\nsubstring: %s", err, substr)
	Fail(t, msg, msgAndArgs...)
}

// Implements asserts that the object implements the interface.
func Implements(t Testing, interfaceType, object any, msgAndArgs ...any) {
	t.Helper()
	interfacePtr := reflect.TypeOf(interfaceType).Elem()

	if !reflect.TypeOf(object).Implements(interfacePtr) {
		msg := fmt.Sprintf("%T does not implement %v", object, interfacePtr)
		Fail(t, msg, msgAndArgs...)
	}
}

// equal checks if two values are equal, handling type conversions and nil values efficiently.
func equal(expected, actual any) bool {
	// Handle nil cases efficiently
	if expected == nil || actual == nil {
		return isNil(expected) == isNil(actual)
	}

	// Fast path for exact equality (but only for comparable types)
	// Check if types are comparable first to avoid panic
	if reflect.TypeOf(expected).Comparable() && reflect.TypeOf(actual).Comparable() {
		if expected == actual {
			return true
		}
	}

	// Use DeepEqual for comprehensive comparison (handles slices, maps, etc.)
	if reflect.DeepEqual(expected, actual) {
		return true
	}

	// Try type conversion as fallback
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
	case reflect.Ptr, reflect.Map, reflect.Chan, reflect.Slice, reflect.Interface, reflect.Func:
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

// isGreater checks if first value is greater than second, with optimized type checking.
func isGreater(first, second any) bool {
	// Fast path for common numeric types
	switch f := first.(type) {
	case int:
		if s, ok := second.(int); ok {
			return f > s
		}
	case int64:
		if s, ok := second.(int64); ok {
			return f > s
		}
	case float64:
		if s, ok := second.(float64); ok {
			return f > s
		}
	case string:
		if s, ok := second.(string); ok {
			return f > s
		}
	}

	// Fallback to reflection for other types
	fv := reflect.ValueOf(first)
	sv := reflect.ValueOf(second)

	// Type compatibility check
	if fv.Kind() != sv.Kind() {
		return false
	}

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

// isLess checks if first value is less than second, with optimized type checking.
func isLess(first, second any) bool {
	// Fast path for common numeric types
	switch f := first.(type) {
	case int:
		if s, ok := second.(int); ok {
			return f < s
		}
	case int64:
		if s, ok := second.(int64); ok {
			return f < s
		}
	case float64:
		if s, ok := second.(float64); ok {
			return f < s
		}
	case string:
		if s, ok := second.(string); ok {
			return f < s
		}
	}

	// Fallback to reflection for other types
	fv := reflect.ValueOf(first)
	sv := reflect.ValueOf(second)

	// Type compatibility check
	if fv.Kind() != sv.Kind() {
		return false
	}

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
