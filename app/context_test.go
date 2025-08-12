package app_test

import (
	"context"
	"os"
	"runtime"
	"syscall"
	"testing"
	"time"

	"github.com/retroenv/retrogolib/app"
	"github.com/retroenv/retrogolib/assert"
)

func TestContext(t *testing.T) {
	ctx := app.Context()
	assert.NotNil(t, ctx, "Context should not be nil")

	// Test that context is not cancelled initially
	select {
	case <-ctx.Done():
		assert.Fail(t, "Context should not be cancelled initially")
	default:
		// Expected behavior - context is not cancelled
	}
}

func TestContext_SIGINT(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping SIGINT test on Windows - Unix signals not supported")
	}

	ctx := app.Context()

	// Send SIGINT to current process
	process, err := os.FindProcess(os.Getpid())
	assert.NoError(t, err, "Should find current process")

	err = process.Signal(syscall.SIGINT)
	assert.NoError(t, err, "Should send SIGINT signal")

	// Context should be cancelled within reasonable time
	select {
	case <-ctx.Done():
		// Expected behavior
	case <-time.After(100 * time.Millisecond):
		assert.Fail(t, "Context should be cancelled after SIGINT")
	}

	assert.Equal(t, context.Canceled, ctx.Err(), "Context should be cancelled")
}

func TestContext_SIGTERM(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping SIGTERM test on Windows - Unix signals not supported")
	}

	ctx := app.Context()

	// Send SIGTERM to current process
	process, err := os.FindProcess(os.Getpid())
	assert.NoError(t, err, "Should find current process")

	err = process.Signal(syscall.SIGTERM)
	assert.NoError(t, err, "Should send SIGTERM signal")

	// Context should be cancelled within reasonable time
	select {
	case <-ctx.Done():
		// Expected behavior
	case <-time.After(100 * time.Millisecond):
		assert.Fail(t, "Context should be cancelled after SIGTERM")
	}

	assert.Equal(t, context.Canceled, ctx.Err(), "Context should be cancelled")
}

func TestContext_MultipleContexts(t *testing.T) {
	ctx1 := app.Context()
	ctx2 := app.Context()

	assert.NotNil(t, ctx1, "First context should not be nil")
	assert.NotNil(t, ctx2, "Second context should not be nil")

	// Both contexts should be independent - they can have the same type but different cancel functions
	// We verify they work independently by checking they don't interfere with each other
	select {
	case <-ctx1.Done():
		assert.Fail(t, "First context should not be cancelled initially")
	default:
		// Expected behavior
	}

	select {
	case <-ctx2.Done():
		assert.Fail(t, "Second context should not be cancelled initially")
	default:
		// Expected behavior
	}
}

func TestContext_Cleanup(t *testing.T) {
	// Test that multiple calls to Context() work correctly
	// and don't interfere with each other
	for i := range 5 {
		ctx := app.Context()
		assert.NotNil(t, ctx, "Context should not be nil")

		// Verify context is not cancelled
		select {
		case <-ctx.Done():
			assert.Fail(t, "Context should not be cancelled initially", "Context iteration %d failed", i)
		default:
			// Expected behavior
		}
	}
}
