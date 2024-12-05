package m6502

import (
	. "github.com/retroenv/retrogolib/addressing"
)

// hasAccumulatorParam returns whether the passed or missing parameter
// indicates usage of the accumulator register.
func hasAccumulatorParam(params ...any) bool {
	if params == nil {
		return true
	}
	param := params[0]
	_, ok := param.(Accumulator)
	return ok
}
