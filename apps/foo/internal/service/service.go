package service

import (
	"fmt"

	"github.com/yuya-takeyama/monotonix-playground/apps/foo/pkg/database"
	"github.com/yuya-takeyama/monotonix-playground/apps/foo/pkg/logger"
)

type UserService struct {
	db     *database.Database
	logger *logger.Logger
}

func NewUserService(db *database.Database, logger *logger.Logger) *UserService {
	return &UserService{
		db:     db,
		logger: logger,
	}
}

func (s *UserService) GetUser(id string) (string, error) {
	s.logger.Info(fmt.Sprintf("Getting user with ID: %s", id))

	results, err := s.db.Query(fmt.Sprintf("SELECT * FROM users WHERE id = '%s'", id))
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to get user: %v", err))
		return "", err
	}

	if len(results) == 0 {
		return "", fmt.Errorf("user not found")
	}

	return fmt.Sprintf("User data: %s", results[0]), nil
}
