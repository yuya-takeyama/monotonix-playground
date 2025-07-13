package database

import (
	"testing"

	"github.com/yuya-takeyama/monotonix-playground/apps/foo/pkg/logger"
)

func TestNew(t *testing.T) {
	log := logger.New("TEST")
	db := New("localhost", log)

	if db == nil {
		t.Error("Expected database instance, got nil")
	}

	if db.host != "localhost" {
		t.Errorf("Expected host localhost, got %s", db.host)
	}
}

func TestConnect(t *testing.T) {
	log := logger.New("TEST")
	db := New("localhost", log)

	err := db.Connect()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestQuery(t *testing.T) {
	log := logger.New("TEST")
	db := New("localhost", log)

	results, err := db.Query("SELECT * FROM users")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}

	if results[0] != "mock_result" {
		t.Errorf("Expected mock_result, got %s", results[0])
	}
}
