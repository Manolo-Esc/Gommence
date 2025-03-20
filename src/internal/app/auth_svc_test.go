package app

import (
	"context"
	"net/http"
	"testing"

	"github.com/Manolo-Esc/gommence/src/internal/domain"
	"github.com/Manolo-Esc/gommence/src/internal/dtos"
	"github.com/Manolo-Esc/gommence/src/internal/infra/jwt"
	"github.com/Manolo-Esc/gommence/src/internal/mocks"
	"github.com/Manolo-Esc/gommence/src/internal/ports"
	"github.com/Manolo-Esc/gommence/src/pkg/cache"
	"github.com/Manolo-Esc/gommence/src/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_PasswordHashes(t *testing.T) {
	var passwords = []string{"passw0rd", "PA$%57ssA", "AlongPasssword1234567890$$", "Pass34With\\SLashes/"}
	for _, pass := range passwords {
		hash, err := HashPassword(pass)
		if err != nil {
			t.Errorf("Error hashing password: %v", err)
		}
		assert.True(t, CheckPassword(pass, hash))
	}
}

var invalidLoginCredentials = []dtos.LoginCredentials{
	{Email: "", Secret: "password"},                // no email
	{Email: "john@mail.com", Secret: ""},           // no secret
	{Email: "johnmail.com", Secret: "googleToken"}, // no email
	{Email: "johnmail@", Secret: "password"},       // no email
	{Email: "666555111", Secret: "password"},       // no email
}

func Test_LoginWith_InvalidData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	perm := mocks.NewMockPermissionService(ctrl)

	userSvc := mocks.NewMockUserService(ctrl)

	svc := NewAuthService(&ServiceInfra{Permissions: perm, Logger: logger.GetNopLogger(), Cache: cache.GetNopCache()}, userSvc)
	for _, loginCredentials := range invalidLoginCredentials {
		loggedUser, err := svc.Login(ctx, loginCredentials)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.Status())
		assert.Nil(t, loggedUser)
	}
}

func Test_LoginWith_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	perm := mocks.NewMockPermissionService(ctrl)

	userSvc := mocks.NewMockUserService(ctrl)
	userSvc.EXPECT().
		GetUserByEmail(gomock.Eq(ctx), gomock.Any()).
		Return(nil, ports.NewAPIError(http.StatusNotFound, "User not found")) // to indicate that we don't have a user with that email

	svc := NewAuthService(&ServiceInfra{Permissions: perm, Logger: logger.GetNopLogger(), Cache: cache.GetNopCache()}, userSvc)
	loginCredentials := dtos.LoginCredentials{Email: "j1@mail.com", Secret: "password"}
	loggedUser, err := svc.Login(ctx, loginCredentials)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusNotFound, err.Status())
	assert.Nil(t, loggedUser)
}

func Test_LoginWith_BadAuthMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	perm := mocks.NewMockPermissionService(ctrl)

	userSvc := mocks.NewMockUserService(ctrl)
	userSvc.EXPECT().
		GetUserByEmail(gomock.Eq(ctx), gomock.Any()).
		Return(&domain.User{ID: "SampleID", AuthMethod: domain.AuthMethGoogle}, nil)

	svc := NewAuthService(&ServiceInfra{Permissions: perm, Logger: logger.GetNopLogger(), Cache: cache.GetNopCache()}, userSvc)
	loginCredentials := dtos.LoginCredentials{Email: "j1@mail.com", Secret: "password"}
	loggedUser, err := svc.Login(ctx, loginCredentials)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Status())
	assert.Nil(t, loggedUser)
}

func Test_LoginWith_WrongPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	perm := mocks.NewMockPermissionService(ctrl)

	email := "j1@mail.com"
	password := "password"
	hashedPassword, err1 := HashPassword(password)
	assert.Nil(t, err1)
	userSvc := mocks.NewMockUserService(ctrl)
	userSvc.EXPECT().
		GetUserByEmail(gomock.Eq(ctx), email).
		Return(&domain.User{ID: "SampleID", AuthMethod: domain.AuthMethPassword, HashedPassword: hashedPassword}, nil)

	svc := NewAuthService(&ServiceInfra{Permissions: perm, Logger: logger.GetNopLogger(), Cache: cache.GetNopCache()}, userSvc)
	loginCredentials := dtos.LoginCredentials{Email: email, Secret: "wrong-password"}
	loggedUser, err := svc.Login(ctx, loginCredentials)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusUnauthorized, err.Status())
	assert.Nil(t, loggedUser)
}

func Test_Login_HappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	perm := mocks.NewMockPermissionService(ctrl)

	email := "j1@mail.com"
	password := "password"
	hashedPassword, err1 := HashPassword(password)
	assert.Nil(t, err1)
	userSvc := mocks.NewMockUserService(ctrl)
	userSvc.EXPECT().
		GetUserByEmail(gomock.Eq(ctx), email).
		Return(&domain.User{ID: "SampleID", AuthMethod: domain.AuthMethPassword, HashedPassword: hashedPassword}, nil)

	svc := NewAuthService(&ServiceInfra{Permissions: perm, Logger: logger.GetNopLogger(), Cache: cache.GetNopCache()}, userSvc)
	loginCredentials := dtos.LoginCredentials{Email: email, Secret: password}
	loggedUser, err := svc.Login(ctx, loginCredentials)
	assert.Nil(t, err)
	assert.NotEqual(t, "", loggedUser.AccessToken)
	claims, err2 := jwt.ValidarToken(loggedUser.AccessToken)
	assert.Nil(t, err2)
	assert.Equal(t, "SampleID", claims["user"])
}
