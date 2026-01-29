package opengl

import (
	"github.com/ebitengine/purego"
	"github.com/retroenv/retrogolib/gui"
	"github.com/retroenv/retrogolib/input"
)

const (
	GLFW_KEY_ESCAPE = 256
)

var keyMapping = map[int]input.Key{
	GLFW_KEY_ESCAPE: input.Escape,
}

func setupKeyCallback(window uintptr, backend gui.Backend) {
	cb := purego.NewCallback(func(window uintptr, key, _ int, action, _ int) {
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

		case GLFW_RELEASE:
			backend.KeyUp(controllerKey)
		}
	})
	glfwSetKeyCallback(window, cb)
}
