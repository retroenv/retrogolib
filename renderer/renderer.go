package renderer

// Renderer describes a renderer that is registered with audio.Register or
// gui.Register. S is the setup function type of the package, for example
// audio.Initializer or gui.Initializer.
type Renderer[S any] struct {
	// Setup starts the renderer.
	Setup S
}

// Info describes the renderer that audio.SetupWithInfo or gui.SetupWithInfo started.
type Info struct {
	// Name is the registered name of the renderer, for example "sdl3".
	Name string
	// Skipped contains the errors of renderers with a higher priority that
	// failed before this renderer started. It is nil when none failed.
	Skipped error
}
