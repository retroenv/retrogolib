package assert

import (
	"errors"
	"fmt"
	"testing"
)

func TestEqual(t *testing.T) {
	Equal(t, 1, 1)
}

func TestNoError(t *testing.T) {
	NoError(t, nil)
}

func TestError(t *testing.T) {
	err := errors.New("error text")
	Error(t, err, err.Error())
}

func TestErrorIs(t *testing.T) {
	errTest := errors.New("error")
	err := fmt.Errorf("error: %w", errTest)
	ErrorIs(t, err, errTest)
}

func TestTrue(t *testing.T) {
	True(t, true)
}

func TestFalse(t *testing.T) {
	False(t, false)
}

func TestInterfaceNilEqual(t *testing.T) {
	var values []int
	Equal(t, nil, values)
}
