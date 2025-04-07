package ports

import (
	"context"

	"github.com/Manolo-Esc/gommence/src/internal/domain"
	"github.com/Manolo-Esc/gommence/src/internal/dtos"
)

type UserRepository interface {
	//GetByID(ctx context.Context, id string) (*domain.User, APIError)
	Create(ctx context.Context, creationData *dtos.InternalUserCreate) (string, APIError)
	GetUserById(ctx context.Context, idUser string) (*domain.User, APIError)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, APIError)
	GetUserIdByEmail(ctx context.Context, email string) string
	GetUsers(ctx context.Context) ([]*domain.User, APIError)
}

type UserService interface {
	CreateUser(ctx context.Context, creationData *dtos.InternalUserCreate) (string, APIError)
	GetUserById(ctx context.Context, idUser string, byUser string) (*domain.User, APIError)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, APIError)
	GetUsers(ctx context.Context, byUser string) ([]*domain.User, APIError)
}
