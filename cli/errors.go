package cli

import (
	"errors"
	"strings"
)

// ErrHelpRequested is returned when the user requests help via --help or -h.
var ErrHelpRequested = errors.New("help requested")

// MissingFlagsError contains details about which required flags are missing.
type MissingFlagsError struct {
	Flags []string
}

func (e *MissingFlagsError) Error() string {
	return "missing required flag(s): " + strings.Join(e.Flags, ", ")
}

// MissingArgsError contains details about which required positional arguments are missing.
type MissingArgsError struct {
	Args []string
}

func (e *MissingArgsError) Error() string {
	return "missing required argument(s): " + strings.Join(e.Args, ", ")
}
