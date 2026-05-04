package opengl

import (
	"testing"
	"unsafe"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/input"
)

func TestCString(t *testing.T) {
	raw := []byte("GLFW failed\x00ignored")

	assert.Equal(t, "GLFW failed", cString((*byte)(unsafe.Pointer(&raw[0]))))
	assert.Equal(t, "", cString(nil))
}

func TestKeyMapping(t *testing.T) {
	assert.Equal(t, input.Escape, keyMapping[GLFW_KEY_ESCAPE])
	assert.Equal(t, input.A, keyMapping[GLFW_KEY_A])
	assert.Equal(t, input.KPEnter, keyMapping[GLFW_KEY_KP_ENTER])
	assert.Equal(t, input.Menu, keyMapping[GLFW_KEY_MENU])
}
