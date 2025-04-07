package dtos

import "github.com/Manolo-Esc/gommence/src/internal/domain"

// @Name User
// @Description User data
type User struct {
	ID             string `json:"id" example:"23GfxRTs"`                // First name of the new user
	FirstName      string `json:"first_name" example:"John"`            // First name of the new user
	FirstLastName  string `json:"last_name" example:"Doe"`              // First last name of the new user
	SecondLastName string `json:"second_last_name" example:"Smith"`     // Second last name of the new user
	Email          string `json:"email" example:"john.doe@example.com"` // Email of the new user
}

func FromDomainUser(user *domain.User) *User {
	return &User{
		ID:             user.ID,
		FirstName:      user.FirstName,
		FirstLastName:  user.FirstLastName,
		SecondLastName: user.SecondLastName,
		Email:          user.Email,
	}
}

func FromDomainUsers(users []*domain.User) []*User {
	result := make([]*User, len(users))
	for i, user := range users {
		result[i] = FromDomainUser(user)
	}
	return result
}

// Intended for internal use only. External request should go through an authentication endpoint
type InternalUserCreate struct {
	FirstName      string            `json:"first_name" validate:"required" example:"John"`                  // First name of the new user
	FirstLastName  string            `json:"last_name" validate:"required" example:"Doe"`                    // First last name of the new user
	SecondLastName string            `json:"second_last_name" example:"Smith"`                               // Second last name of the new user
	Email          string            `json:"email" validate:"required,email" example:"john.doe@example.com"` // Email of the new user
	AuthMethod     domain.AuthMethod // Authentication method of the new user
	HashedPassword string            `json:"hashed_password" example:"123GfxRTs"` // Hashed user password in case of AuthMethod is AuthMethPassword
}

func fromDtosUserSignUp(creationData *UserSignUp) *InternalUserCreate {
	return &InternalUserCreate{
		FirstName:      creationData.FirstName,
		FirstLastName:  creationData.FirstLastName,
		SecondLastName: creationData.SecondLastName,
		AuthMethod:     creationData.AuthMethod,
		Email:          creationData.Email,
		HashedPassword: creationData.Secret,
	}
}
