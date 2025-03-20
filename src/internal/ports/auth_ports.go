package ports

import (
	"context"

	"github.com/Manolo-Esc/gommence/src/internal/dtos"
)

type AuthService interface {
	Login(ctx context.Context, credentials dtos.LoginCredentials) (*dtos.LoggedUser, APIError)
}
