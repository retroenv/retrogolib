package chip8

// Quirks selects behavior that differs from the original COSMAC VIP CHIP-8
// interpreter. Its zero value preserves the original interpreter semantics.
type Quirks struct {
	// DisplayWaitDisabled allows more than one sprite instruction per 60 Hz tick.
	DisplayWaitDisabled bool
	// JumpUsesVX makes BNNN add VX, selected by NNN's high nibble, instead of V0.
	JumpUsesVX bool
	// LoadStoreLeavesI keeps I unchanged after FX55 and FX65.
	LoadStoreLeavesI bool
	// LogicPreservesVF keeps VF unchanged after 8XY1, 8XY2, and 8XY3.
	LogicPreservesVF bool
	// ShiftUsesVX makes 8XY6 and 8XYE shift VX instead of VY.
	ShiftUsesVX bool
	// WrapSprites wraps overflowing sprite pixels instead of clipping them.
	WrapSprites bool
}

type options struct {
	quirks Quirks
}

// Option configures a CPU.
type Option func(*options)

// WithQuirks configures interpreter compatibility behavior.
func WithQuirks(quirks Quirks) Option {
	return func(options *options) {
		options.quirks = quirks
	}
}

func newOptions(optionList ...Option) options {
	var opts options
	for _, option := range optionList {
		if option != nil {
			option(&opts)
		}
	}

	return opts
}
