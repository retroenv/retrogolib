// Package renderer contains the types that the audio and gui packages use to
// register and select renderers.
//
// Renderer packages register a Renderer with audio.Register or gui.Register.
// audio.SetupWithInfo and gui.SetupWithInfo return an Info that names the
// renderer that started.
package renderer
