package stubs

import (
	"errors"

	"github.com/kaschula/socket-server/app"
)

func NewUserRepositoryStub(users map[string]*app.User) app.UserRepository {
	return &UserRepositoryStub{users}
}

type UserRepositoryStub struct {
	users map[string]*app.User
}

func (r *UserRepositoryStub) GetUser(id string) (*app.User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, errors.New("User not found")
	}

	return user, nil
}

// type UserRepository interface {
// 	GetUser(ID string) (*User, error)
// }
