package opengl

import (
	"errors"
	"sync"
	"unsafe"
)

var errMu sync.RWMutex
var lastErr error

var errCallback = errorCallback

func setErrorCallback() {
	glfwSetErrorCallback(uintptr(unsafe.Pointer(&errCallback)))
}

func errorCallback(err string) {
	errMu.Lock()
	lastErr = errors.New(err)
	errMu.Unlock()
}

func getLastError() error {
	errMu.RLock()
	defer errMu.RUnlock()
	return lastErr
}
