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
	// models := []interface{}{
	// 	&UserTest{},
	// }

	// err := db.WithContext(ctx).AutoMigrate(models...) // Crear tablas
	// if err != nil {
	// 	return err
	// }
	return nil
}

// type UserTest struct {
// 	repos_db.BaseDBModel
// 	FirstName      string
// 	FirstLastName  string
// 	SecondLastName sql.NullString
// 	Email          string `gorm:"uniqueIndex"`
// }
