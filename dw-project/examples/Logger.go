package examples

import (
	"etl_pipeline/pkg/logger"
	"os"
)

func Run() {
	// ── 1. Default logger (colors, caller, timestamp all on) ──────────────
	logger.Debug("Initializing application config")
	logger.Info("Server listening on :8080")
	logger.Warn("Connection pool at 80% capacity")
	logger.Error("Failed to validate JWT token")

	// ── 2. Structured fields ──────────────────────────────────────────────
	logger.Info("Request handled",
		logger.F("method", "GET"),
		logger.F("path", "/api/users"),
		logger.F("latency", "12ms"),
		logger.F("status", 200),
	)

	// ── 3. Printf-style helpers ───────────────────────────────────────────
	logger.Infof("Worker pool started with %d goroutines", 8)
	logger.Warnf("Retry attempt %d/%d", 2, 5)

	// ── 4. Custom logger (warnings and above, no caller, no color) ────────
	plain := logger.New(logger.Options{
		Level:      logger.LevelWarn,
		Output:     os.Stderr,
		ShowTime:   true,
		ShowCaller: false,
		NoColor:    true,
	})
	plain.Warn("This goes to stderr without colors")
	plain.Error("Something bad happened")

	// ── 5. Child logger with pre-attached fields (e.g. per-request) ───────
	reqLog := logger.Default.WithFields(
		logger.F("requestID", "req-abc123"),
		logger.F("userID", 42),
	)
	reqLog.Info("Handler entered")
	reqLog.Warn("Slow query detected", logger.F("query", "SELECT * FROM orders"), logger.F("ms", 520))
	reqLog.Error("Database timeout")

	// ── 6. Change level at runtime ────────────────────────────────────────
	logger.Default.SetLevel(logger.LevelError)
	logger.Debug("This won't print — below the new level")
	logger.Error("This will print")
}
