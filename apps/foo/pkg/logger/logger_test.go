package logger

import (
	"testing"
)

func TestNew(t *testing.T) {
	prefix := "TEST"
	logger := New(prefix)

	if logger == nil {
		t.Error("Expected logger instance, got nil")
	}

	if logger.prefix != prefix {
		t.Errorf("Expected prefix %s, got %s", prefix, logger.prefix)
	}
}

func TestLoggerMethods(t *testing.T) {
	logger := New("TEST")

	// これらのメソッドはログ出力するだけなので、パニックしないことを確認
	logger.Info("test info message")
	logger.Error("test error message")
}
