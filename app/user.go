package app

import "errors"

// User
func NewUser(ID string) *User {
	return &User{ID}
}

type User struct {
	ID string `json:"ID"`
}

// This should be a service
type UserRepository interface {
	GetUser(ID string) (*User, error)
}

func NewInMemoryUserRepository(users map[string]*User) UserRepository {
	return &InMemoryUserRepository{users}
}

// Should have constructor so Users can be private
type InMemoryUserRepository struct {
	Users map[string]*User
}

func (ur *InMemoryUserRepository) GetUser(ID string) (*User, error) {
	user, ok := ur.Users[ID]
	if !ok || user == nil {
		return nil, errors.New("Unable to resolve User")
	}

	return user, nil
}
