package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/Manolo-Esc/gommence/src/internal/adapters/repos_db"
	"github.com/Manolo-Esc/gommence/src/internal/domain"
	"github.com/Manolo-Esc/gommence/src/internal/dtos"
	opologger "github.com/Manolo-Esc/gommence/src/pkg/logger"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type databaseIntegrationSuite struct {
	suite.Suite
	db *gorm.DB
}

func (s *databaseIntegrationSuite) SetupSuite() { // SetupSuite runs once, before all tests
	dsn := "host=localhost user=postgres password=secret dbname=integration_tests port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		s.T().Fatalf("Error connecting to database: %v", err)
	}
	s.db = db

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	//err = database.Migrate(ctx, s.db)
	err = createTestDatabase(ctx, s.db)
	if err != nil {
		s.T().Fatalf("Error running migration: %v", err)
	}
}

func (s *databaseIntegrationSuite) TearDownSuite() { // TearDownSuite runs once, after all tests
	s.db.Unscoped().Where("first_name = ?", "John").Delete(&repos_db.User{})
}

func (s *databaseIntegrationSuite) SetupTest()    {} // SetupTest before individual test
func (s *databaseIntegrationSuite) TearDownTest() {} // TearDownTest after individual test

func (s *databaseIntegrationSuite) Test_CreateEntityWithPID_DupFields() {
	dbUser := repos_db.User{
		FirstName:      "John",
		FirstLastName:  "Doe",
		Email:          fmt.Sprintf("jdoe@%d.com", time.Now().Nanosecond()),
		AuthMethod:     domain.AuthMethPassword,
		HashedPassword: sql.NullString{String: "hashedPassword", Valid: true},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err := repos_db.CreateEntityWithPID(ctx, s.db, &dbUser)
	s.Nil(err)
	err = repos_db.CreateEntityWithPID(ctx, s.db, &dbUser)
	s.Equal(http.StatusConflict, err.Status()) // Duplicated fields marked as unique (email, pk auto-fixed)
	dbUser.Email = fmt.Sprintf("jdoe@%d.com", time.Now().Nanosecond()+100)
	err = repos_db.CreateEntityWithPID(ctx, s.db, &dbUser) // Duplicated fields marked as unique (pk)
	s.Nil(err)                                             // No error: Only ID was duplicated and it was auto-fixed in the function
}

func (s *databaseIntegrationSuite) Test_FindUsersByEmail() {
	dbUser := dtos.UserCreate{
		FirstName:      "John",
		FirstLastName:  "Doe",
		Email:          "John.Doe@mailer.com",
		AuthMethod:     domain.AuthMethPassword,
		HashedPassword: "hashedPassword",
	}

	repo := repos_db.NewUserRepository(&repos_db.DBReposInfra{Db: s.db, Logger: opologger.GetNopLogger()})
	johnID, err := repo.Create(context.Background(), &dbUser)
	s.Nil(err)
	s.NotEmpty(johnID)

	userID := repo.GetUserIdByEmail(context.Background(), dbUser.Email)
	s.Equal(johnID, userID)
	userID = repo.GetUserIdByEmail(context.Background(), "john.doe@mailer.com")
	s.Equal(johnID, userID)
	userID = repo.GetUserIdByEmail(context.Background(), "JOHN.DOE@mailer.com")
	s.Equal(johnID, userID)
	userID = repo.GetUserIdByEmail(context.Background(), "JOHN.DOE@Mailer.COm")
	s.Equal(johnID, userID)
	userID = repo.GetUserIdByEmail(context.Background(), "I_am_not_registered@mailer.com")
	s.Equal("", userID)
}

func TestRunSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database integration suite in short mode") // text only seen with -v
	}
	log.Println("Running database integration suite")
	suite.Run(t, new(databaseIntegrationSuite))
}
