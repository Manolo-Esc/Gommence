package domain

type AuthMethod int

const (
	AuthMethPassword AuthMethod = 1
	AuthMethGoogle   AuthMethod = 2
)

type User struct {
	ID             string
	FirstName      string
	FirstLastName  string
	SecondLastName string
	Email          string
	AuthMethod     AuthMethod
	HashedPassword string
}
