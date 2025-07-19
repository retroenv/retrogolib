// Package log provides fast, structured logging based on Go's slog package.
//
// This package wraps Go's standard slog library with additional convenience
// functions and configuration options specifically designed for retro console
// emulation and tooling development.
//
// # Features
//
//   - Structured logging with key-value pairs
//   - Multiple output formats (text, JSON)
//   - Configurable log levels
//   - High performance with minimal allocations
//   - Console-friendly formatting
//   - Testing utilities for log verification
//
// # Basic Usage
//
//	import "github.com/retroenv/retrogolib/log"
//
//	func main() {
//		logger := log.New(log.Config{
//			Level:  log.LevelInfo,
//			Format: log.FormatText,
//		})
//
//		logger.Info("Starting emulation",
//			log.String("system", "NES"),
//			log.Int("rom_size", 32768),
//		)
//
//		logger.Error("Failed to load ROM",
//			log.String("filename", "game.nes"),
//			log.String("error", err.Error()),
//		)
//	}
//
// # Log Levels
//
//   - Debug: Detailed diagnostic information
//   - Info: General operational messages
//   - Warn: Warning conditions that don't halt operation
//   - Error: Error conditions that may affect functionality
//
// # Output Formats
//
//   - Text: Human-readable console output
//   - JSON: Structured JSON for log aggregation systems
//
// # Performance
//
// The logging system is designed for high performance:
//   - Zero allocation for disabled log levels
//   - Efficient field handling
//   - Minimal overhead in hot paths like CPU emulation loops
//
// # Testing Support
//
// The package includes utilities for testing log output:
//   - Capture log messages in tests
//   - Verify specific log entries were written
//   - Mock logging for isolated unit tests
//
// # Thread Safety
//
// All logging operations are thread-safe and can be used concurrently
// from multiple goroutines without external synchronization.
package log
