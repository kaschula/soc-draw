package stubs

import (
	"errors"

	"github.com/kaschula/soc-draw/app"
)

func NewUserServiceStub(users map[string]*app.User) app.UserService {
	return &UserServiceStub{users}
}

type UserServiceStub struct {
	users map[string]*app.User
}

func (r *UserServiceStub) GetUser(id string) (*app.User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, errors.New("User not found")
	}

	return user, nil
}
