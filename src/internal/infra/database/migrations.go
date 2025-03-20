package database

import (
	"context"
	"fmt"

	"github.com/Manolo-Esc/gommence/src/internal/adapters/repos_db"
	repos "github.com/Manolo-Esc/gommence/src/internal/adapters/repos_db"
	"github.com/Manolo-Esc/gommence/src/internal/app"
	"github.com/Manolo-Esc/gommence/src/internal/domain"
	"github.com/Manolo-Esc/gommence/src/internal/dtos"
	"github.com/Manolo-Esc/gommence/src/pkg/logger"
	"gorm.io/gorm"
)

type VersionDBEntity struct {
	Major int `gorm:"column:major"`
	Minor int `gorm:"column:minor"`
	Patch int `gorm:"column:patch"`
}

func (VersionDBEntity) TableName() string {
	return "version_db"
}

var systemUsersCreate = []dtos.UserCreate{
	{FirstName: "Esmerelda", FirstLastName: "Weatherwax", Email: "granny@lancre.dw", AuthMethod: domain.AuthMethPassword},
	{FirstName: "Sam", FirstLastName: "Vimes", Email: "theduke@ankh.dw", AuthMethod: domain.AuthMethPassword},
}
var systemUsers []domain.User // The above systemUsersCreate information plus the ID

func Migrate(ctx context.Context, db *gorm.DB) error {
	if !db.Migrator().HasTable(&VersionDBEntity{}) {
		fmt.Println("version_db table not found. Initializing database")
		err := createDatabase(ctx, db)
		if err != nil {
			return err
		}
	} else {
		var version VersionDBEntity
		result := db.Limit(1).Find(&version)
		if result.Error != nil {
			return result.Error
		}
		runMigrations(ctx, db, version)
	}
	return nil
}

func createDatabase(ctx context.Context, db *gorm.DB) error {
	models := []interface{}{
		&VersionDBEntity{},
		&repos.User{},
	}
	err := db.WithContext(ctx).AutoMigrate(models...) // Create tables
	if err != nil {
		return err
	}
	err = populateDatabase(ctx, db)
	if err == nil {
		err = populateDevelopmentDatabase(ctx, db, systemUsers) // Populate with development data
	}
	return err
}

func runMigrations(ctx context.Context, db *gorm.DB, version VersionDBEntity) {
	fmt.Printf("Current database version: %d.%d.%d\n", version.Major, version.Minor, version.Patch)

	if version.Major == 1 && version.Minor == 0 && version.Patch == 0 {
		/*
			fmt.Println("Migrating database to version 1.1.0...")
			// Add here the code to migrate from 1.0.0 to 1.1.0
		*/
	}
}

func populateDatabase(ctx context.Context, db *gorm.DB) error {
	if err := db.WithContext(ctx).Create(&VersionDBEntity{Major: 1, Minor: 0, Patch: 0}); err.Error != nil {
		return err.Error
	}
	if err := createUsers(ctx, db); err != nil {
		return err
	}
	// create more entities here
	return nil
}

func createUsers(ctx context.Context, db *gorm.DB) error {
	var repo = repos.NewUserRepository(&repos_db.DBReposInfra{Db: db, Logger: logger.GetLogger()})

	for _, user := range systemUsersCreate {
		user.HashedPassword, _ = app.HashPassword("password")
		if userId, err := repo.Create(ctx, &user); err == nil {
			domainUser := domain.User{ID: userId, Email: user.Email}
			systemUsers = append(systemUsers, domainUser)
			return err
		}
	}
	return nil
}
