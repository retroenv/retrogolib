package host

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

// FromEnvironment reads RETROHOST_HEADLESS and RETROHOST_FRAMES.
// Other option fields are left unchanged. An empty value leaves its default.
func FromEnvironment(options Options) (Options, error) {
	if value := os.Getenv("RETROHOST_HEADLESS"); value != "" {
		headless, err := strconv.ParseBool(value)
		if err != nil {
			return options, fmt.Errorf("RETROHOST_HEADLESS must be a boolean: %w", err)
		}
		options.Headless = headless
	}

	if value := os.Getenv("RETROHOST_FRAMES"); value != "" {
		limit, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return options, fmt.Errorf("RETROHOST_FRAMES must be a positive frame count: %w", err)
		}

		if limit == 0 {
			return options, errors.New("RETROHOST_FRAMES must be a positive frame count")
		}
		options.FrameLimit = limit
	}

	return options, nil
}
