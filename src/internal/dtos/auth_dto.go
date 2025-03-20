package dtos

import "github.com/Manolo-Esc/gommence/src/internal/domain"

// @Name LoginCredentials
// @Description Request to sign in the program
type LoginCredentials struct {
	Email  string `json:"email" validate:"required,email" example:"john.doe@example.com"` // Email of the user signing in
	Secret string `json:"secret" validate:"required" example:"password"`                  // Password or third party token of the user signing in
}

// @Name LoggedUser
// @Description Logged user information
type LoggedUser struct {
	AccessToken string `json:"access_token"` // Authenticaton token
}

// @Name UserSignUp
// @Description Request to register a new user in the platform
type UserSignUp struct {
	FirstName      string            `json:"first_name" validate:"required" example:"John"`                  // First name of the new user
	FirstLastName  string            `json:"last_name" validate:"required" example:"Doe"`                    // First last name of the new user
	SecondLastName string            `json:"second_last_name" example:"Smith"`                               // Second last name of the new user
	Email          string            `json:"email" validate:"required,email" example:"john.doe@example.com"` // Email of the new user
	AuthMethod     domain.AuthMethod // Authentication method of the new user
	Secret         string            `json:"secret" validate:"required" example:"password"` // Password or third party token of the user signing in
}
