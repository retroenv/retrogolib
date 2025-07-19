/*
Package assert provides a set of testing assertion helpers that make test code more readable and maintainable.

# Overview

The assert package offers a comprehensive set of assertion functions similar to those found in popular
testing libraries. It provides clear error messages when assertions fail and integrates seamlessly
with Go's testing package through the Testing interface.

# Basic Usage

All assertion functions follow a similar pattern: they take a Testing interface (usually *testing.T),
the values to compare or check, and optional message formatting arguments.

	func TestExample(t *testing.T) {
		result := Calculate()
		assert.Equal(t, 42, result, "calculation should return 42")

		err := DoSomething()
		assert.NoError(t, err, "operation should succeed")
	}

# Available Assertions

Equality and Comparison:
  - Equal: Asserts two values are equal
  - NotEqual: Asserts two values are not equal
  - Greater: Asserts first value is greater than second
  - GreaterOrEqual: Asserts first value is greater than or equal to second
  - Less: Asserts first value is less than second
  - LessOrEqual: Asserts first value is less than or equal to second

Boolean Assertions:
  - True: Asserts value is true
  - False: Asserts value is false

Nil Checks:
  - Nil: Asserts value is nil
  - NotNil: Asserts value is not nil

Collection Assertions:
  - Len: Asserts collection has expected length
  - Empty: Asserts collection is empty
  - NotEmpty: Asserts collection is not empty

String Assertions:
  - Contains: Asserts string contains substring
  - NotContains: Asserts string does not contain substring

Error Handling:
  - NoError: Asserts error is nil
  - Error: Asserts error is not nil and has expected message
  - ErrorIs: Asserts error matches expected error using errors.Is
  - ErrorContains: Asserts error message contains substring

Panic Assertions:
  - Panics: Asserts function panics when called
  - NotPanics: Asserts function does not panic when called

Type Assertions:
  - Implements: Asserts object implements interface

# Custom Testing Interface

The package uses a Testing interface that matches *testing.T, allowing for easy mocking in tests:

	type Testing interface {
		Helper()
		Error(args ...any)
		FailNow()
	}

# Examples

	// Equality checks
	assert.Equal(t, expected, actual, "values should be equal")
	assert.NotEqual(t, 1, 2, "values should be different")

	// Comparisons
	assert.Greater(t, 10, 5, "10 should be greater than 5")
	assert.LessOrEqual(t, score, maxScore, "score should not exceed maximum")

	// Error handling
	assert.NoError(t, err, "operation should not fail")
	assert.ErrorIs(t, err, ErrNotFound, "should return not found error")
	assert.ErrorContains(t, err, "permission denied", "should be permission error")

	// Collections
	assert.Len(t, items, 5, "should have 5 items")
	assert.Empty(t, errors, "no errors should occur")
	assert.NotEmpty(t, results, "results should not be empty")

	// Panic handling
	assert.Panics(t, func() { divide(1, 0) }, "division by zero should panic")
	assert.NotPanics(t, func() { process(data) }, "processing should not panic")

	// String checks
	assert.Contains(t, output, "success", "output should indicate success")
	assert.NotContains(t, err.Error(), "panic", "error should not mention panic")

	// Type checks
	assert.Implements(t, (*io.Writer)(nil), &MyWriter{}, "should implement io.Writer")
*/
package assert
