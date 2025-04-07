package repos_db

import (
	"database/sql"

	"github.com/Manolo-Esc/gommence/src/internal/domain"
	"github.com/Manolo-Esc/gommence/src/internal/dtos"
)

// This will be a table in the database
type User struct {
	BaseDBModel
	FirstName      string
	FirstLastName  string
	SecondLastName sql.NullString
	Email          string `gorm:"uniqueIndex"`
	AuthMethod     domain.AuthMethod
	HashedPassword sql.NullString // Can be null depending on the AuthMethod
}

func fromDtosUserCreate(creationData *dtos.InternalUserCreate) *User {
	return &User{
		FirstName:      creationData.FirstName,
		FirstLastName:  creationData.FirstLastName,
		SecondLastName: sql.NullString{String: creationData.SecondLastName, Valid: creationData.SecondLastName != ""},
		AuthMethod:     creationData.AuthMethod,
		Email:          creationData.Email,
		HashedPassword: sql.NullString{String: creationData.HashedPassword, Valid: creationData.HashedPassword != ""},
	}
}

func (u *User) toDomainUser() *domain.User {
	return &domain.User{
		ID:             u.ID,
		FirstName:      u.FirstName,
		FirstLastName:  u.FirstLastName,
		SecondLastName: u.SecondLastName.String,
		Email:          u.Email,
		AuthMethod:     u.AuthMethod,
		HashedPassword: u.HashedPassword.String,
	}
}
