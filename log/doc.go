// Package log provides nil-safe, leveled structured logging built on log/slog.
//
// A Logger created with New writes human-readable records to standard output.
// Typed fields avoid reflection, and the Func field constructors defer expensive
// values until a handler accepts the record.
//
//	logger := log.New()
//	logger.Info("Starting emulation",
//		log.String("system", "NES"),
//		log.Int("rom_size", 32768),
//	)
//
// Use NewWithConfig to change the level, destination, timestamp format, source
// reporting, or handler.
//
//	cfg := log.DefaultConfig()
//	cfg.Level = log.DebugLevel
//	cfg.Output = os.Stderr
//	logger := log.NewWithConfig(cfg)
//
// Logger methods are safe for concurrent use. Methods on a nil *Logger are
// no-ops, except Fatal and FatalContext, which still terminate the process.
// Context-taking methods treat a nil context as context.Background.
// NewTestLogger routes records through testing.TB and fails the test immediately
// when an error-level record is emitted.
package log
