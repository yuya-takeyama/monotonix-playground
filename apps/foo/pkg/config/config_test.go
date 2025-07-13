package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	config := Load()

	if config == nil {
		t.Error("Expected config instance, got nil")
	}

	// デフォルト値のテスト
	if config.Port != "8080" {
		t.Errorf("Expected default port 8080, got %s", config.Port)
	}

	if config.LogLevel != "info" {
		t.Errorf("Expected default log level info, got %s", config.LogLevel)
	}

	if config.DBHost != "localhost" {
		t.Errorf("Expected default DB host localhost, got %s", config.DBHost)
	}
}

func TestLoadWithEnvironment(t *testing.T) {
	// 環境変数を設定
	os.Setenv("PORT", "9000")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("DB_HOST", "remote-db")

	config := Load()

	if config.Port != "9000" {
		t.Errorf("Expected port 9000, got %s", config.Port)
	}

	if config.LogLevel != "debug" {
		t.Errorf("Expected log level debug, got %s", config.LogLevel)
	}

	if config.DBHost != "remote-db" {
		t.Errorf("Expected DB host remote-db, got %s", config.DBHost)
	}

	// クリーンアップ
	os.Unsetenv("PORT")
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("DB_HOST")
}
