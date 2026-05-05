package opengl

import (
	"errors"
	"image"
	"testing"
	"unsafe"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/gui"
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

func TestCleanupOpenGL(t *testing.T) {
	originalGlDeleteTextures := glDeleteTextures
	originalGlfwTerminate := glfwTerminate
	t.Cleanup(func() {
		glDeleteTextures = originalGlDeleteTextures
		glfwTerminate = originalGlfwTerminate
	})

	var calls []string
	glDeleteTextures = func(n int32, textures *uint32) {
		calls = append(calls, "texture")
		assert.Equal(t, int32(1), n)
		assert.Equal(t, uint32(3), *textures)
	}
	glfwTerminate = func() {
		calls = append(calls, "terminate")
	}

	cleanupOpenGL(3)

	assert.Equal(t, []string{"texture", "terminate"}, calls)
}

func TestCleanupOpenGLSkipsEmptyTexture(t *testing.T) {
	originalGlDeleteTextures := glDeleteTextures
	originalGlfwTerminate := glfwTerminate
	t.Cleanup(func() {
		glDeleteTextures = originalGlDeleteTextures
		glfwTerminate = originalGlfwTerminate
	})

	var calls []string
	glDeleteTextures = func(_ int32, _ *uint32) {
		calls = append(calls, "texture")
	}
	glfwTerminate = func() {
		calls = append(calls, "terminate")
	}

	cleanupOpenGL(0)

	assert.Equal(t, []string{"terminate"}, calls)
}

func TestRenderOpenGL(t *testing.T) {
	restoreOpenGLRenderFunctions(t)

	dimensions := gui.Dimensions{
		ScaleFactor: 1,
		Height:      2,
		Width:       3,
	}
	img := image.NewRGBA(image.Rect(0, 0, dimensions.Width, dimensions.Height))
	var calls []string
	installRenderSpies(t, dimensions, &calls)

	err := renderOpenGL(dimensions, img, 11, 7)
	assert.NoError(t, err)
	assert.Equal(t, expectedRenderCalls(), calls)
}

func TestRenderOpenGLRejectsInvalidImage(t *testing.T) {
	dimensions := gui.Dimensions{
		ScaleFactor: 1,
		Height:      2,
		Width:       3,
	}

	err := renderOpenGL(dimensions, nil, 11, 7)
	assert.ErrorContains(t, err, "getting image pixels")
}

func TestSetupLibrary(t *testing.T) {
	restoreLoadFunctions(t)

	openGLLibName, err := getOpenGLSystemLibrary()
	assert.NoError(t, err)
	glfwLibName, err := getGlfwSystemLibrary()
	assert.NoError(t, err)

	var calls []libraryLoad
	loadFunctions = func(name string, imports map[string]any) (uintptr, error) {
		group := "OpenGL"
		if len(calls) == 1 {
			group = "GLFW"
		}
		assert.NotEqual(t, 0, len(imports))
		calls = append(calls, libraryLoad{name: name, group: group})
		return uintptr(len(calls)), nil
	}

	err = setupLibrary()
	assert.NoError(t, err)
	assert.Equal(t, []libraryLoad{
		{name: openGLLibName, group: "OpenGL"},
		{name: glfwLibName, group: "GLFW"},
	}, calls)
}

func TestSetupLibraryWrapsOpenGLLoadError(t *testing.T) {
	restoreLoadFunctions(t)

	loadErr := errors.New("load failed")
	var calls int
	loadFunctions = func(_ string, _ map[string]any) (uintptr, error) {
		calls++
		return 0, loadErr
	}

	err := setupLibrary()
	assert.ErrorContains(t, err, "loading OpenGL functions")
	assert.ErrorIs(t, err, loadErr)
	assert.Equal(t, 1, calls)
}

func TestSetupLibraryWrapsGLFWLoadError(t *testing.T) {
	restoreLoadFunctions(t)

	loadErr := errors.New("load failed")
	var calls int
	loadFunctions = func(_ string, _ map[string]any) (uintptr, error) {
		calls++
		if calls == 2 {
			return 0, loadErr
		}
		return uintptr(calls), nil
	}

	err := setupLibrary()
	assert.ErrorContains(t, err, "loading GLFW functions")
	assert.ErrorIs(t, err, loadErr)
	assert.Equal(t, 2, calls)
}

type libraryLoad struct {
	name  string
	group string
}

func restoreLoadFunctions(t *testing.T) {
	t.Helper()

	original := loadFunctions
	t.Cleanup(func() {
		loadFunctions = original
	})
}

func restoreOpenGLRenderFunctions(t *testing.T) {
	t.Helper()

	originalGlBindTexture := glBindTexture
	originalGlTexSubImage2D := glTexSubImage2D
	originalGlMatrixMode := glMatrixMode
	originalGlLoadIdentity := glLoadIdentity
	originalGlOrtho := glOrtho
	originalGlBegin := glBegin
	originalGlTexCoord2d := glTexCoord2d
	originalGlVertex2d := glVertex2d
	originalGlEnd := glEnd
	originalGlfwSwapBuffers := glfwSwapBuffers
	originalGlfwPollEvents := glfwPollEvents

	t.Cleanup(func() {
		glBindTexture = originalGlBindTexture
		glTexSubImage2D = originalGlTexSubImage2D
		glMatrixMode = originalGlMatrixMode
		glLoadIdentity = originalGlLoadIdentity
		glOrtho = originalGlOrtho
		glBegin = originalGlBegin
		glTexCoord2d = originalGlTexCoord2d
		glVertex2d = originalGlVertex2d
		glEnd = originalGlEnd
		glfwSwapBuffers = originalGlfwSwapBuffers
		glfwPollEvents = originalGlfwPollEvents
	})
}

func installRenderSpies(t *testing.T, dimensions gui.Dimensions, calls *[]string) {
	t.Helper()

	installTextureSpies(t, dimensions, calls)
	installProjectionSpies(t, dimensions, calls)
	installQuadSpies(t, calls)
	installWindowSpies(t, calls)
}

func installTextureSpies(t *testing.T, dimensions gui.Dimensions, calls *[]string) {
	t.Helper()

	glBindTexture = func(target, texture uint32) {
		*calls = append(*calls, "glBindTexture")
		assert.Equal(t, uint32(GL_TEXTURE_2D), target)
		assert.Equal(t, uint32(7), texture)
	}
	glTexSubImage2D = func(target, level, xoffset, yoffset, width, height, format, xtype int32, pixels uintptr) {
		*calls = append(*calls, "glTexSubImage2D")
		assert.Equal(t, GL_TEXTURE_2D, target)
		assert.Equal(t, int32(0), level)
		assert.Equal(t, int32(0), xoffset)
		assert.Equal(t, int32(0), yoffset)
		assert.Equal(t, int32(dimensions.Width), width)
		assert.Equal(t, int32(dimensions.Height), height)
		assert.Equal(t, GL_RGBA, format)
		assert.Equal(t, GL_UNSIGNED_BYTE, xtype)
		assert.NotEqual(t, uintptr(0), pixels)
	}
}

func installProjectionSpies(t *testing.T, dimensions gui.Dimensions, calls *[]string) {
	t.Helper()

	glMatrixMode = func(mode int32) {
		*calls = append(*calls, "glMatrixMode")
		assert.True(t, mode == GL_PROJECTION || mode == GL_MODELVIEW)
	}
	glLoadIdentity = func() {
		*calls = append(*calls, "glLoadIdentity")
	}
	glOrtho = func(left, right, bottom, top, near, far float64) {
		*calls = append(*calls, "glOrtho")
		assert.Equal(t, 0.0, left)
		assert.Equal(t, float64(dimensions.Width), right)
		assert.Equal(t, 0.0, bottom)
		assert.Equal(t, float64(dimensions.Height), top)
		assert.Equal(t, -1.0, near)
		assert.Equal(t, 1.0, far)
	}
}

func installQuadSpies(t *testing.T, calls *[]string) {
	t.Helper()

	glBegin = func(mode int32) {
		*calls = append(*calls, "glBegin")
		assert.Equal(t, GL_QUADS, mode)
	}
	glTexCoord2d = func(_, _ float64) {
		*calls = append(*calls, "glTexCoord2d")
	}
	glVertex2d = func(_, _ float64) {
		*calls = append(*calls, "glVertex2d")
	}
	glEnd = func() {
		*calls = append(*calls, "glEnd")
	}
}

func installWindowSpies(t *testing.T, calls *[]string) {
	t.Helper()

	glfwSwapBuffers = func(window uintptr) {
		*calls = append(*calls, "glfwSwapBuffers")
		assert.Equal(t, uintptr(11), window)
	}
	glfwPollEvents = func() {
		*calls = append(*calls, "glfwPollEvents")
	}
}

func expectedRenderCalls() []string {
	return []string{
		"glBindTexture",
		"glTexSubImage2D",
		"glMatrixMode",
		"glLoadIdentity",
		"glOrtho",
		"glMatrixMode",
		"glBegin",
		"glTexCoord2d",
		"glVertex2d",
		"glTexCoord2d",
		"glVertex2d",
		"glTexCoord2d",
		"glVertex2d",
		"glTexCoord2d",
		"glVertex2d",
		"glEnd",
		"glfwSwapBuffers",
		"glfwPollEvents",
	}
}
