package database

import (
	"fmt"

	"github.com/yuya-takeyama/monotonix-playground/apps/foo/pkg/logger"
)

type Database struct {
	host   string
	logger *logger.Logger
}

func New(host string, logger *logger.Logger) *Database {
	return &Database{
		host:   host,
		logger: logger,
	}
}

func (db *Database) Connect() error {
	db.logger.Info(fmt.Sprintf("Connecting to database at %s", db.host))
	return nil
}

func (db *Database) Query(sql string) ([]string, error) {
	db.logger.Info(fmt.Sprintf("Executing query: %s", sql))
	return []string{"mock_result"}, nil
}