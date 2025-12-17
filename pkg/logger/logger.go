package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func init() {
	// Configure zerolog
	zerolog.TimeFieldFormat = time.RFC3339

	// Pretty console output for development
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "15:04:05",
	}

	log = zerolog.New(output).
		With().
		Timestamp().
		Caller().
		Logger()
}

// Get returns the global logger instance
func Get() zerolog.Logger {
	return log
}

// WithService returns a logger with service name
func WithService(serviceName string) zerolog.Logger {
	return log.With().Str("service", serviceName).Logger()
}

// WithRequestID returns a logger with request ID
func WithRequestID(requestID string) zerolog.Logger {
	return log.With().Str("request_id", requestID).Logger()
}

// Info logs an info message
func Info(msg string) {
	log.Info().Msg(msg)
}

// Error logs an error message
func Error(err error, msg string) {
	log.Error().Err(err).Msg(msg)
}

// Debug logs a debug message
func Debug(msg string) {
	log.Debug().Msg(msg)
}

// Fatal logs a fatal message and exits
func Fatal(err error, msg string) {
	log.Fatal().Err(err).Msg(msg)
}
