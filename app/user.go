package app

import "errors"

func NewUser(ID string) *User {
	return &User{ID}
}

type User struct {
	ID string `json:"ID"`
}

type UserService interface {
	GetUser(ID string) (*User, error)
}

func NewInMemoryUserService(users map[string]*User) UserService {
	return &InMemoryUserService{users}
}

type InMemoryUserService struct {
	Users map[string]*User
}

func (ur *InMemoryUserService) GetUser(ID string) (*User, error) {
	user, ok := ur.Users[ID]
	if !ok || user == nil {
		return nil, errors.New("Unable to resolve User")
	}

	return user, nil
}
