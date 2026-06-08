// Package logger provides a thin zerolog wrapper so the rest of the codebase
// depends on this package rather than importing zerolog (or its global logger)
// directly.
package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Logger is an alias for zerolog.Logger. Callers accept *Logger via constructor
// injection instead of reaching for the global logger.
type Logger = zerolog.Logger

// New returns a console logger writing to stderr. When verbose is true the
// level is set to debug, otherwise info.
func New(verbose bool) Logger {
	level := zerolog.InfoLevel
	if verbose {
		level = zerolog.DebugLevel
	}

	writer := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	return zerolog.New(writer).Level(level).With().Timestamp().Logger()
}
