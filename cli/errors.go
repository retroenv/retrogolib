package cli

import "strings"

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
