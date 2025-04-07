package database

import (
	"context"

	"github.com/Manolo-Esc/gommence/src/internal/infra/database"
	"gorm.io/gorm"
)

func createTestDatabase(ctx context.Context, db *gorm.DB) error {
	err := database.Migrate(ctx, db)
	if err != nil {
		return err
	}
	return nil
}
