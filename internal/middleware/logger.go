package middleware

import (
	"go.uber.org/zap"
)

// Log is a global logger instance that can be used throughout the codebase.
var Log *zap.Logger

// Init initializes the zap logger. Use zap.NewProduction() in production mode.
// For now, we are using zap.NewDevelopment() for human-friendly logs.
func Init() error {
	var err error
	Log, err = zap.NewDevelopment()
	if err != nil {
		return err
	}
	return nil
}

// Sync flushes any buffered log entries.
func Sync() error {
	return Log.Sync()
}
