package opengl

import (
	"fmt"
	"sync"

	"github.com/ebitengine/purego"
)

var errMu sync.RWMutex
var lastErr error

func setErrorCallback() {
	cb := purego.NewCallback(func(code int, description string) {
		errMu.Lock()
		lastErr = fmt.Errorf("GLFW error %d: %s", code, description)
		errMu.Unlock()
	})
	glfwSetErrorCallback(cb)
}

func getLastError() error {
	errMu.RLock()
	defer errMu.RUnlock()
	return lastErr
}
