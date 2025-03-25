package app

import (
	"context"
	"net/http"

	"github.com/Manolo-Esc/gommence/src/internal/domain"
	"github.com/Manolo-Esc/gommence/src/internal/dtos"
	"github.com/Manolo-Esc/gommence/src/internal/ports"
	"github.com/Manolo-Esc/gommence/src/pkg/validator"
)

type UserServiceImpl struct {
	repo ports.UserRepository
	si   *ServiceInfra
}

func NewUserService(repo ports.UserRepository, serviceInfra *ServiceInfra) ports.UserService {
	return &UserServiceImpl{repo: repo, si: serviceInfra}
}

// CreateUser creates a new user in the platform. It is intended to be used only internally. REST calls shall be targeted to the auth_svc.
func (s *UserServiceImpl) CreateUser(ctx context.Context, creationData *dtos.UserCreate) (string, ports.APIError) {
	if err := validator.ValidateStruct(creationData); err != nil {
		return "", ports.NewAPIError(http.StatusBadRequest, err.Error())
	}
	if creationData.AuthMethod == domain.AuthMethPassword {
		if creationData.HashedPassword == "" {
			return "", ports.NewAPIError(http.StatusBadRequest, "Password is required")
		}
	} else {
		return "", ports.NewAPIError(http.StatusBadRequest, "Only password authentication is currently supported")
	}

	existingUser := s.repo.GetUserIdByEmail(ctx, creationData.Email)
	if existingUser != "" {
		return "", ports.NewAPIError(http.StatusBadRequest, "User already exists")
	}

	userId, err := s.repo.Create(ctx, creationData)
	if err != nil {
		return "", err
	}
	return userId, nil

}

// GetUserByEmail retrieves a domain.User by its email or nil if not found
func (s *UserServiceImpl) GetUserByEmail(ctx context.Context, email string) (*domain.User, ports.APIError) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUsers retrieves a all users
func (s *UserServiceImpl) GetUsers(ctx context.Context, byUser string) ([]*domain.User, ports.APIError) {
	users, err := s.repo.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}
