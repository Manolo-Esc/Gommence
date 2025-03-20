package ports

import (
	"context"

	"github.com/Manolo-Esc/gommence/src/internal/domain"
	"github.com/Manolo-Esc/gommence/src/internal/dtos"
)

type UserRepository interface {
	//GetByID(ctx context.Context, id string) (*domain.User, APIError)
	Create(ctx context.Context, creationData *dtos.UserCreate) (string, APIError)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, APIError)
	GetUserIdByEmail(ctx context.Context, email string) string
}

type UserService interface {
	CreateUser(ctx context.Context, creationData *dtos.UserCreate) (string, APIError)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, APIError)
}
