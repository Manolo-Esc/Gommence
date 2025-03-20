package app

import (
	"context"
	"net/http"
	"testing"

	"github.com/Manolo-Esc/gommence/src/internal/domain"
	"github.com/Manolo-Esc/gommence/src/internal/dtos"
	"github.com/Manolo-Esc/gommence/src/internal/mocks"
	"github.com/Manolo-Esc/gommence/src/pkg/cache"
	"github.com/Manolo-Esc/gommence/src/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserCreationWithValidData01(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	perm := mocks.NewMockPermissionService(ctrl)

	repo := mocks.NewMockUserRepository(ctrl)
	repo.EXPECT().
		GetUserIdByEmail(gomock.Eq(ctx), gomock.Any()).
		Return("") // to indicate that we don't have a user with that email
	repo.EXPECT().
		Create(gomock.Eq(ctx), gomock.Any()).
		Return("JohnId", nil)

	svc := NewUserService(repo, &ServiceInfra{Permissions: perm, Logger: logger.GetNopLogger(), Cache: cache.GetNopCache()})
	creationParams := &dtos.UserCreate{
		FirstName:      "John",
		FirstLastName:  "Doe",
		SecondLastName: "",
		Email:          "john@mail.com",
		AuthMethod:     domain.AuthMethPassword,
		HashedPassword: "password",
	}
	newUser, err := svc.CreateUser(ctx, creationParams)
	assert.Nil(t, err)
	assert.Equal(t, "JohnId", newUser)
}

func TestUserCreationWithValidData02(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	perm := mocks.NewMockPermissionService(ctrl)

	repo := mocks.NewMockUserRepository(ctrl)
	repo.EXPECT().
		GetUserIdByEmail(gomock.Eq(ctx), gomock.Any()).
		Return("") // to indicate that we don't have a user with that email
	repo.EXPECT().
		Create(gomock.Eq(ctx), gomock.Any()).
		Return("JohnId", nil)

	svc := NewUserService(repo, &ServiceInfra{Permissions: perm, Logger: logger.GetNopLogger(), Cache: cache.GetNopCache()})
	creationParams := &dtos.UserCreate{
		FirstName:      "John",
		FirstLastName:  "Doe",
		SecondLastName: "Smith",
		Email:          "john@mail.com",
		AuthMethod:     domain.AuthMethPassword,
		HashedPassword: "password",
	}
	newUser, err := svc.CreateUser(ctx, creationParams)
	assert.Nil(t, err)
	assert.Equal(t, "JohnId", newUser)
}

func TestUserCreationWithDupEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	perm := mocks.NewMockPermissionService(ctrl)

	repo := mocks.NewMockUserRepository(ctrl)
	repo.EXPECT().
		GetUserIdByEmail(gomock.Eq(ctx), gomock.Any()).
		Return("JohnId")

	svc := NewUserService(repo, &ServiceInfra{Permissions: perm, Logger: logger.GetNopLogger(), Cache: cache.GetNopCache()})
	creationParams := &dtos.UserCreate{
		FirstName:      "John",
		FirstLastName:  "Doe",
		SecondLastName: "",
		Email:          "john@mail.com",
		AuthMethod:     domain.AuthMethPassword,
		HashedPassword: "password",
	}
	newUser, err := svc.CreateUser(ctx, creationParams)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Status())
	assert.Equal(t, "", newUser)
}

var invalidUsersCreate = []dtos.UserCreate{
	{FirstName: "", FirstLastName: "Doe", SecondLastName: "", Email: "j1@mail.com", AuthMethod: domain.AuthMethPassword, HashedPassword: "password"},
	{FirstName: "John", FirstLastName: "", SecondLastName: "", Email: "j1@mail.com", AuthMethod: domain.AuthMethPassword, HashedPassword: "password"},
	{FirstName: "John", FirstLastName: "Doe", SecondLastName: "", Email: "", AuthMethod: domain.AuthMethPassword, HashedPassword: "password"},
	{FirstName: "John", FirstLastName: "Doe", SecondLastName: "", Email: "john.Doe", AuthMethod: domain.AuthMethPassword, HashedPassword: "password"},
	{FirstName: "John", FirstLastName: "Doe", SecondLastName: "", Email: "j1@mail.com", AuthMethod: domain.AuthMethGoogle, HashedPassword: "password"},
	{FirstName: "John", FirstLastName: "Doe", SecondLastName: "", Email: "j1@mail.com", AuthMethod: domain.AuthMethPassword, HashedPassword: ""},
}

func TestUserCreationWithInvalidData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	perm := mocks.NewMockPermissionService(ctrl)
	repo := mocks.NewMockUserRepository(ctrl)

	svc := NewUserService(repo, &ServiceInfra{Permissions: perm, Logger: logger.GetNopLogger(), Cache: cache.GetNopCache()})
	for _, creationParams := range invalidUsersCreate {
		newUser, err := svc.CreateUser(ctx, &creationParams)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.Status())
		assert.Equal(t, "", newUser)
	}
}
