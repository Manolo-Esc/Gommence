package dtos

import "github.com/Manolo-Esc/gommence/src/internal/domain"

// @Name UserCreate
// @Description Request to create a new user in the platform
type UserCreate struct {
	FirstName      string            `json:"first_name" validate:"required" example:"John"`                  // First name of the new user
	FirstLastName  string            `json:"last_name" validate:"required" example:"Doe"`                    // First last name of the new user
	SecondLastName string            `json:"second_last_name" example:"Smith"`                               // Second last name of the new user
	Email          string            `json:"email" validate:"required,email" example:"john.doe@example.com"` // Email of the new user
	AuthMethod     domain.AuthMethod // Authentication method of the new user
	HashedPassword string            `json:"hashed_password" example:"123GfxRTs"` // Hashed user password in case of AuthMethod is AuthMethPassword
}

func fromDtosUserSignUp(creationData *UserSignUp) *UserCreate {
	return &UserCreate{
		FirstName:      creationData.FirstName,
		FirstLastName:  creationData.FirstLastName,
		SecondLastName: creationData.SecondLastName,
		AuthMethod:     creationData.AuthMethod,
		Email:          creationData.Email,
		HashedPassword: creationData.Secret,
	}
}
