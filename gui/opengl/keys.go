package opengl

import (
	"github.com/retroenv/retrogolib/gui"
	"github.com/retroenv/retrogolib/input"
)

const (
	GLFW_KEY_ESCAPE = 256
)

var keyMapping = map[int]input.Key{
	GLFW_KEY_ESCAPE: input.Escape,
}

type keyCallbackFunc = func(window uintptr, key int, _ int, action int, _ int)

var keyCallback keyCallbackFunc

func onGLFWKey(backend gui.Backend) keyCallbackFunc {
	return func(window uintptr, key, _, action, _ int) {
		if action == GLFW_PRESS && key == GLFW_KEY_ESCAPE {
			glfwSetWindowShouldClose(window, GLFW_TRUE)
			return
		}

		controllerKey, ok := keyMapping[key]
		if !ok {
			return
		}

		switch action {
		case GLFW_PRESS:
			backend.KeyDown(controllerKey)

		case GLFW_KEY_ESCAPE:
			backend.KeyUp(controllerKey)
		}
	}
}
