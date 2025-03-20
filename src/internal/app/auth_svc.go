package app

import (
	"context"
	"net/http"

	"github.com/Manolo-Esc/gommence/src/internal/domain"
	"github.com/Manolo-Esc/gommence/src/internal/dtos"
	"github.com/Manolo-Esc/gommence/src/internal/infra/jwt"
	"github.com/Manolo-Esc/gommence/src/internal/ports"
	"github.com/Manolo-Esc/gommence/src/pkg/validator"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImpl struct {
	si      *ServiceInfra
	userSvc ports.UserService
}

func NewAuthService(serviceInfra *ServiceInfra, userSvc ports.UserService) ports.AuthService {
	return &AuthServiceImpl{si: serviceInfra, userSvc: userSvc}
}

func (s *AuthServiceImpl) Login(ctx context.Context, credentials dtos.LoginCredentials) (*dtos.LoggedUser, ports.APIError) {
	if err := validator.ValidateStruct(credentials); err != nil {
		return nil, ports.NewAPIError(http.StatusBadRequest, err.Error())
	}
	user, err := s.userSvc.GetUserByEmail(ctx, credentials.Email)
	if err != nil {
		return nil, err
	}
	if user.AuthMethod == domain.AuthMethPassword {
		if !CheckPassword(credentials.Secret, user.HashedPassword) {
			return nil, ports.NewAPIError(http.StatusUnauthorized, "Invalid credentials")
		}
	} else {
		return nil, ports.NewAPIError(http.StatusBadRequest, "Only password authentication is currently supported")
	}

	token, err2 := jwt.GenerarToken(user.ID)
	if err2 != nil {
		return nil, ports.NewAPIError(http.StatusInternalServerError, err2.Error())
	}
	return &dtos.LoggedUser{AccessToken: token}, nil
}

// func (s *UserServiceImpl) CreateUser(ctx context.Context, creationData *dtos.UserCreate) (*dtos.LoggedUser, ports.APIError) {
// 	userId, err := s.repo.Create(ctx, creationData)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// xxxx get token
// 	return ret, nil
// }

// HashPassword generates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPassword compares a password with a hash to check if they match
func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil // Returns true if passwords match
}

// CreateUser saves a new user with a hashed password
// func CreateUser(db *gorm.DB, username, password string) error {
// 	hashedPassword, err := HashPassword(password)
// 	if err != nil {
// 		return err
// 	}

// 	user := User{
// 		Username:       username,
// 		HashedPassword: hashedPassword,
// 	}

// 	return db.Create(&user).Error
// }

// // AuthenticateUser checks username and password against stored data
// func AuthenticateUser(db *gorm.DB, username, password string) bool {
// 	var user User
// 	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
// 		return false // User not found
// 	}

// 	return CheckPassword(password, user.HashedPassword)
// }
