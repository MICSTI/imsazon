// This package contains the user model

package user

import "errors"

// UserId uniquely identifies a user
type UserId string

type User struct {
	Id				UserId
	Name			string
	Email			string
	Username		string
	Password		string
	Role			UserRole
}

// New creates a new user
func New(id UserId, name string, email string, username string, password string, role UserRole) *User {
	return &User{
		Id:				id,
		Name:			name,
		Email:			email,
		Username:		username,
		Password:		password,
		Role:			role,
	}
}

// Repository provides access to an in-memory user store
type Repository interface {
	// adds a user to the store
	Add(user *User) error

	// attempts to find the user by id inside the store
	Find(id UserId) (*User, error)

	// returns an array of all users inside the store
	FindAll() []*User

	// checks if the login credentials match a user inside the store
	CheckLogin(username string, password string) (*User, error)
}

// ErrUnknown is used if the user cannot be found
var ErrUnknown = errors.New("Unknown user")

// UserRole describes the role of the user
type UserRole int

// valid user roles
const (
	Nobody UserRole = iota
	Standard
	Admin
	Service
)

func (s UserRole) String() string {
	switch s {
	case Nobody:
		return "Nobody"
	case Standard:
		return "Standard"
	case Admin:
		return "Admin"
	case Service:
		return "Service"
	}
	return ""
}