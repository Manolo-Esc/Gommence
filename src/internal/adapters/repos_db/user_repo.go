package repos_db

import (
	"context"
	"net/http"

	"github.com/Manolo-Esc/gommence/src/internal/domain"
	"github.com/Manolo-Esc/gommence/src/internal/dtos"
	"github.com/Manolo-Esc/gommence/src/internal/ports"
	opentelemetry "github.com/Manolo-Esc/gommence/src/pkg/open_telemetry"
)

type UserRepositoryDB struct {
	dbInfra *DBReposInfra
}

func NewUserRepository(dbInfra *DBReposInfra) ports.UserRepository {
	return &UserRepositoryDB{dbInfra: dbInfra}
}

func (r *UserRepositoryDB) Create(ctx context.Context, creationData *dtos.UserCreate) (string, ports.APIError) {
	ctx, span := opentelemetry.GetTracer().Start(ctx, "UserRepositoryDB.Create")
	defer span.End()

	dbUser := fromDtosUserCreate(creationData)
	err := CreateEntityWithPID(ctx, r.dbInfra.Db, dbUser)
	return dbUser.ID, err
}

// GetUserIdByEmail returns the user ID associated with the given email. If the email is not found, it returns an empty string.
func (r *UserRepositoryDB) GetUserIdByEmail(ctx context.Context, email string) string {
	var user User
	r.dbInfra.Db.WithContext(ctx).Where("email ILIKE ?", email).First(&user)
	return user.ID // If not found, it will return an empty string
}

// GetUserByEmail retrieves a domain.User by its email or nil if not found
func (r *UserRepositoryDB) GetUserByEmail(ctx context.Context, email string) (*domain.User, ports.APIError) {
	var user User
	r.dbInfra.Db.WithContext(ctx).Where("email ILIKE ?", email).First(&user)
	if user.ID == "" {
		return nil, ports.NewAPIError(http.StatusNotFound, "User not found")
	}
	return user.toDomainUser(), nil
}

func (r *UserRepositoryDB) GetUsers(ctx context.Context) ([]*domain.User, ports.APIError) {
	ctx, span := opentelemetry.GetTracer().Start(ctx, "UserRepositoryDB.GetUsers")
	defer span.End()

	var records []User
	result := r.dbInfra.Db.WithContext(ctx).Find(&records)

	if result.Error != nil {
		return nil, ports.NewAPIError(http.StatusInternalServerError, result.Error.Error())
	}
	users := make([]*domain.User, len(records))
	for i, _ := range records {
		users[i] = records[i].toDomainUser()
	}
	return users, nil
}
