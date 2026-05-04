package opengl

import (
	"errors"
	"fmt"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
)

var errMu sync.RWMutex
var errLastGLFW error

func setErrorCallback() {
	cb := purego.NewCallback(func(_ purego.CDecl, code int32, description *byte) {
		errMu.Lock()
		errLastGLFW = fmt.Errorf("GLFW error %d: %s", code, cString(description))
		errMu.Unlock()
	})
	glfwSetErrorCallback(cb)
}

func resetLastError() {
	errMu.Lock()
	errLastGLFW = nil
	errMu.Unlock()
}

func getLastError() error {
	errMu.RLock()
	defer errMu.RUnlock()
	if errLastGLFW == nil {
		return errors.New("unknown GLFW error")
	}
	return errLastGLFW
}

func cString(ptr *byte) string {
	if ptr == nil {
		return ""
	}

	cstr := unsafe.Pointer(ptr)
	var length int
	for *(*byte)(unsafe.Add(cstr, length)) != 0 {
		length++
	}
	return unsafe.String((*byte)(cstr), length)
}
