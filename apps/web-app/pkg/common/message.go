package common

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	// Pretty console logging
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

// GetTimestampedMessage returns a message with the current timestamp
func GetTimestampedMessage(prefix string, message string) string {
	return fmt.Sprintf("[%s] %s: %s", time.Now().Format(time.RFC3339), prefix, message)
}

// GetVersion returns the common library version
func GetVersion() string {
	return "v1.1.0"
}

// FormatMessage creates a nicely formatted message with emoji
func FormatMessage(level, message string) string {
	emoji := "ℹ️"
	switch level {
	case "info":
		emoji = "ℹ️"
	case "warn":
		emoji = "⚠️"
	case "error":
		emoji = "❌"
	case "success":
		emoji = "✅"
	}
	return fmt.Sprintf("%s %s", emoji, message)
}

// LogMessage logs a message using zerolog
func LogMessage(level, message string) {
	switch level {
	case "info":
		log.Info().Msg(message)
	case "warn":
		log.Warn().Msg(message)
	case "error":
		log.Error().Msg(message)
	case "debug":
		log.Debug().Msg(message)
	default:
		log.Info().Msg(message)
	}
}
