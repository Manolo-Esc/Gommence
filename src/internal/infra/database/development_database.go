package database

import (
	"context"

	"github.com/Manolo-Esc/gommence/src/internal/domain"
	"gorm.io/gorm"
)

func populateDevelopmentDatabase(ctx context.Context, db *gorm.DB, systemUsers []domain.User) error {
	return nil
}
