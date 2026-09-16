package dynlib

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestMissingLibrary(t *testing.T) {
	for range 2 {
		lib, err := LoadFunctions("retrogolib-missing-library", nil)
		assert.Error(t, err)
		assert.Equal(t, uintptr(0), lib)
	}
}

func TestMissingSymbol(t *testing.T) {
	var missing func()
	err := registerFunction(0, "retrogolib_missing_symbol", &missing)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "retrogolib_missing_symbol")
}
