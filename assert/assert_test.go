package assert

import (
	"errors"
	"fmt"
	"slices"
	"testing"
)

func TestEqual(t *testing.T) {
	tst := &errorCapture{}
	Equal(tst, 1, 1)
	if tst.failed {
		t.Error("Equal failed")
	}

	tst = &errorCapture{}
	Equal(tst, 1, 2)
	if !tst.failed {
		t.Error("Equal failed")
	}
}

func TestNotEqual(t *testing.T) {
	tst := &errorCapture{}
	NotEqual(tst, 1, 2)
	if tst.failed {
		t.Error("NotEqual failed")
	}

	tst = &errorCapture{}
	NotEqual(tst, 1, 1)
	if !tst.failed {
		t.Error("NotEqual failed")
	}
}

func TestNoError(t *testing.T) {
	tst := &errorCapture{}
	NoError(tst, nil)
	if tst.failed {
		t.Error("NoError failed")
	}

	tst = &errorCapture{}
	NoError(tst, errors.New("error"))
	if !tst.failed {
		t.Error("NoError failed")
	}
}

func TestError(t *testing.T) {
	tst := &errorCapture{}
	Error(tst, errors.New("error"))
	if tst.failed {
		t.Error("Error failed")
	}

	tst = &errorCapture{}
	Error(tst, nil)
	if !tst.failed {
		t.Error("Error failed")
	}
}

func TestErrorIs(t *testing.T) {
	tst := &errorCapture{}
	ErrorIs(tst, errors.New("error"), errors.New("error"))
	if !tst.failed {
		t.Error("ErrorIs failed")
	}

	tst = &errorCapture{}
	ErrorIs(tst, errors.New("error"), errors.New("other"))
	if !tst.failed {
		t.Error("ErrorIs failed")
	}

	tst = &errorCapture{}
	ErrorIs(tst, nil, errors.New("error"))
	if !tst.failed {
		t.Error("ErrorIs failed")
	}

	tst = &errorCapture{}
	err := errors.New("error")
	ErrorIs(tst, fmt.Errorf("wrapped: %w", err), err)
	if tst.failed {
		t.Error("ErrorIs failed")
	}
}

func TestTrue(t *testing.T) {
	tst := &errorCapture{}
	True(tst, true)
	if tst.failed {
		t.Error("True failed")
	}

	tst = &errorCapture{}
	True(tst, false)
	if !tst.failed {
		t.Error("True failed")
	}
}

func TestFalse(t *testing.T) {
	tst := &errorCapture{}
	False(tst, false)
	if tst.failed {
		t.Error("False failed")
	}

	tst = &errorCapture{}
	False(tst, true)
	if !tst.failed {
		t.Error("False failed")
	}
}

func TestInterfaceNilEqual(t *testing.T) {
	tst := &errorCapture{}
	Equal(tst, nil, nil)
	if tst.failed {
		t.Error("InterfaceNilEqual failed")
	}

	tst = &errorCapture{}
	Equal(tst, nil, 1)
	if !tst.failed {
		t.Error("InterfaceNilEqual failed")
	}
}

func TestLen(t *testing.T) {
	tst := &errorCapture{}
	Len(tst, []int{1, 2}, 2)
	if tst.failed {
		t.Error("Len failed")
	}

	tst = &errorCapture{}
	Len(tst, []int{}, 2)
	if !tst.failed {
		t.Error("Len failed")
	}

	// Test nil object
	tst = &errorCapture{}
	Len(tst, nil, 0)
	if !tst.failed {
		t.Error("Len should fail for nil")
	}

	// Test invalid type
	tst = &errorCapture{}
	Len(tst, 42, 1)
	if !tst.failed {
		t.Error("Len should fail for non-length type")
	}
}

func TestNotNil(t *testing.T) {
	tst := &errorCapture{}
	NotNil(tst, 1)
	if tst.failed {
		t.Error("NotNil failed")
	}

	tst = &errorCapture{}
	NotNil(tst, nil)
	if !tst.failed {
		t.Error("NotNil failed")
	}
}

func TestNil(t *testing.T) {
	tst := &errorCapture{}
	Nil(tst, nil)
	if tst.failed {
		t.Error("Nil failed")
	}

	tst = &errorCapture{}
	Nil(tst, 1)
	if !tst.failed {
		t.Error("Nil failed")
	}

	// Test nil slice
	var nilSlice []int
	tst = &errorCapture{}
	Nil(tst, nilSlice)
	if tst.failed {
		t.Error("Nil should pass for nil slice")
	}

	// Test nil function
	var nilFunc func()
	tst = &errorCapture{}
	Nil(tst, nilFunc)
	if tst.failed {
		t.Error("Nil should pass for nil function")
	}
}

func TestFail(t *testing.T) {
	tst := &errorCapture{}
	Fail(tst, "error", "msg %d", 1)
	if !tst.failed {
		t.Error("Fail failed")
	}
	if tst.errs[0].(string) != "error\nmsg 1" {
		t.Error("Fail failed")
	}
}

func TestContains(t *testing.T) {
	tst := &errorCapture{}
	Contains(tst, "hello world", "world")
	if tst.failed {
		t.Error("Contains failed")
	}

	tst = &errorCapture{}
	Contains(tst, "hello world", "foo")
	if !tst.failed {
		t.Error("Contains failed")
	}
}

func TestNotContains(t *testing.T) {
	tst := &errorCapture{}
	NotContains(tst, "hello world", "foo")
	if tst.failed {
		t.Error("NotContains failed")
	}

	tst = &errorCapture{}
	NotContains(tst, "hello world", "world")
	if !tst.failed {
		t.Error("NotContains failed")
	}
}

func TestPanics(t *testing.T) {
	tst := &errorCapture{}
	Panics(tst, func() { panic("test panic") })
	if tst.failed {
		t.Error("Panics failed")
	}

	tst = &errorCapture{}
	Panics(tst, func() {})
	if !tst.failed {
		t.Error("Panics failed")
	}
}

func TestNotPanics(t *testing.T) {
	tst := &errorCapture{}
	NotPanics(tst, func() {})
	if tst.failed {
		t.Error("NotPanics failed")
	}

	tst = &errorCapture{}
	NotPanics(tst, func() { panic("test panic") })
	if !tst.failed {
		t.Error("NotPanics failed")
	}
}

func TestEmpty(t *testing.T) {
	tst := &errorCapture{}
	Empty(tst, "")
	if tst.failed {
		t.Error("Empty failed for empty string")
	}

	tst = &errorCapture{}
	Empty(tst, []int{})
	if tst.failed {
		t.Error("Empty failed for empty slice")
	}

	tst = &errorCapture{}
	Empty(tst, make(map[string]int))
	if tst.failed {
		t.Error("Empty failed for empty map")
	}

	tst = &errorCapture{}
	Empty(tst, "hello")
	if !tst.failed {
		t.Error("Empty failed for non-empty string")
	}

	tst = &errorCapture{}
	Empty(tst, []int{1, 2})
	if !tst.failed {
		t.Error("Empty failed for non-empty slice")
	}
}

func TestNotEmpty(t *testing.T) {
	tst := &errorCapture{}
	NotEmpty(tst, "hello")
	if tst.failed {
		t.Error("NotEmpty failed for non-empty string")
	}

	tst = &errorCapture{}
	NotEmpty(tst, []int{1})
	if tst.failed {
		t.Error("NotEmpty failed for non-empty slice")
	}

	tst = &errorCapture{}
	NotEmpty(tst, "")
	if !tst.failed {
		t.Error("NotEmpty failed for empty string")
	}

	tst = &errorCapture{}
	NotEmpty(tst, []int{})
	if !tst.failed {
		t.Error("NotEmpty failed for empty slice")
	}
}

func TestGreater(t *testing.T) {
	tst := &errorCapture{}
	Greater(tst, 2, 1)
	if tst.failed {
		t.Error("Greater failed for 2 > 1")
	}

	tst = &errorCapture{}
	Greater(tst, 1.5, 1.0)
	if tst.failed {
		t.Error("Greater failed for 1.5 > 1.0")
	}

	tst = &errorCapture{}
	Greater(tst, "b", "a")
	if tst.failed {
		t.Error("Greater failed for \"b\" > \"a\"")
	}

	tst = &errorCapture{}
	Greater(tst, 1, 2)
	if !tst.failed {
		t.Error("Greater failed for 1 > 2")
	}

	tst = &errorCapture{}
	Greater(tst, 1, 1)
	if !tst.failed {
		t.Error("Greater failed for 1 > 1")
	}

	// Test incompatible types
	tst = &errorCapture{}
	Greater(tst, 1, "a")
	if !tst.failed {
		t.Error("Greater should fail for incompatible types")
	}
}

func TestGreaterOrEqual(t *testing.T) {
	tst := &errorCapture{}
	GreaterOrEqual(tst, 2, 1)
	if tst.failed {
		t.Error("GreaterOrEqual failed for 2 >= 1")
	}

	tst = &errorCapture{}
	GreaterOrEqual(tst, 1, 1)
	if tst.failed {
		t.Error("GreaterOrEqual failed for 1 >= 1")
	}

	tst = &errorCapture{}
	GreaterOrEqual(tst, 1, 2)
	if !tst.failed {
		t.Error("GreaterOrEqual failed for 1 >= 2")
	}
}

func TestLess(t *testing.T) {
	tst := &errorCapture{}
	Less(tst, 1, 2)
	if tst.failed {
		t.Error("Less failed for 1 < 2")
	}

	tst = &errorCapture{}
	Less(tst, 1.0, 1.5)
	if tst.failed {
		t.Error("Less failed for 1.0 < 1.5")
	}

	tst = &errorCapture{}
	Less(tst, "a", "b")
	if tst.failed {
		t.Error("Less failed for \"a\" < \"b\"")
	}

	tst = &errorCapture{}
	Less(tst, 2, 1)
	if !tst.failed {
		t.Error("Less failed for 2 < 1")
	}

	tst = &errorCapture{}
	Less(tst, 1, 1)
	if !tst.failed {
		t.Error("Less failed for 1 < 1")
	}
}

func TestLessOrEqual(t *testing.T) {
	tst := &errorCapture{}
	LessOrEqual(tst, 1, 2)
	if tst.failed {
		t.Error("LessOrEqual failed for 1 <= 2")
	}

	tst = &errorCapture{}
	LessOrEqual(tst, 1, 1)
	if tst.failed {
		t.Error("LessOrEqual failed for 1 <= 1")
	}

	tst = &errorCapture{}
	LessOrEqual(tst, 2, 1)
	if !tst.failed {
		t.Error("LessOrEqual failed for 2 <= 1")
	}
}

func TestLessOrEqualCrossType(t *testing.T) {
	// Test uint64 vs int (common case from benchmark tests)
	tst := &errorCapture{}
	LessOrEqual(tst, uint64(162), int(165))
	if tst.failed {
		t.Error("LessOrEqual failed for uint64(162) <= int(165)")
	}

	tst = &errorCapture{}
	LessOrEqual(tst, uint64(165), int(165))
	if tst.failed {
		t.Error("LessOrEqual failed for uint64(165) <= int(165)")
	}

	tst = &errorCapture{}
	LessOrEqual(tst, uint64(166), int(165))
	if !tst.failed {
		t.Error("LessOrEqual should fail for uint64(166) <= int(165)")
	}

	// Test int vs uint64
	tst = &errorCapture{}
	LessOrEqual(tst, int(100), uint64(200))
	if tst.failed {
		t.Error("LessOrEqual failed for int(100) <= uint64(200)")
	}

	// Test negative int vs uint64
	tst = &errorCapture{}
	LessOrEqual(tst, int(-1), uint64(0))
	if tst.failed {
		t.Error("LessOrEqual failed for int(-1) <= uint64(0)")
	}

	// Test int vs float64
	tst = &errorCapture{}
	LessOrEqual(tst, int(10), float64(10.5))
	if tst.failed {
		t.Error("LessOrEqual failed for int(10) <= float64(10.5)")
	}

	// Test uint64 vs float64
	tst = &errorCapture{}
	LessOrEqual(tst, uint64(100), float64(100.0))
	if tst.failed {
		t.Error("LessOrEqual failed for uint64(100) <= float64(100.0)")
	}
}

func TestErrorContains(t *testing.T) {
	tst := &errorCapture{}
	ErrorContains(tst, errors.New("this is an error message"), "error message")
	if tst.failed {
		t.Error("ErrorContains failed")
	}

	tst = &errorCapture{}
	ErrorContains(tst, errors.New("this is an error message"), "foo")
	if !tst.failed {
		t.Error("ErrorContains failed")
	}

	tst = &errorCapture{}
	ErrorContains(tst, nil, "error")
	if !tst.failed {
		t.Error("ErrorContains failed for nil error")
	}
}

type testInterface interface {
	Method()
}

type testImpl struct{}

func (t testImpl) Method() {}

func TestImplements(t *testing.T) {
	tst := &errorCapture{}
	Implements(tst, (*testInterface)(nil), testImpl{})
	if tst.failed {
		t.Error("Implements failed")
	}

	tst = &errorCapture{}
	Implements(tst, (*testInterface)(nil), struct{}{})
	if !tst.failed {
		t.Error("Implements failed")
	}
}

type errorCapture struct {
	errs   []any
	failed bool
}

func (e *errorCapture) Helper() {
}

func (e *errorCapture) Error(args ...any) {
	e.errs = slices.Clone(args)
}

func (e *errorCapture) FailNow() {
	e.failed = true
}

// Additional edge case tests
func TestEqual_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		expected any
		actual   any
		wantFail bool
	}{
		{"nil vs nil", nil, nil, false},
		{"nil vs zero int", nil, 0, true},
		{"zero vs zero", 0, 0, false},
		{"empty string vs empty string", "", "", false},
		{"slice comparison", []int{1, 2}, []int{1, 2}, false},
		{"different slice", []int{1, 2}, []int{2, 1}, true},
		{"type conversion", 42, int64(42), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tst := &errorCapture{}
			Equal(tst, tt.expected, tt.actual)
			if tst.failed != tt.wantFail {
				t.Errorf("Equal() failed = %v, wantFail = %v", tst.failed, tt.wantFail)
			}
		})
	}
}

func TestErrorAs_EdgeCases(t *testing.T) {
	tst := &errorCapture{}
	var pathErr *AssertTestError
	ErrorAs(tst, &AssertTestError{msg: "test"}, &pathErr)
	if tst.failed {
		t.Error("ErrorAs should pass for matching error type")
	}

	tst = &errorCapture{}
	var wrongType *AssertTestSecondError
	ErrorAs(tst, &AssertTestError{msg: "test"}, &wrongType)
	if !tst.failed {
		t.Error("ErrorAs should fail for non-matching error type")
	}
}

type AssertTestError struct {
	msg string
}

func (e *AssertTestError) Error() string {
	return e.msg
}

type AssertTestSecondError struct {
	msg string
}

func (e *AssertTestSecondError) Error() string {
	return e.msg
}
