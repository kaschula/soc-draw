package stubs

import "github.com/kaschula/soc-draw/app"

type UserClientStub struct {
	User   *app.User
	Client app.IsClient
}

func (u *UserClientStub) GetUser() *app.User {
	return u.User
}

func (u *UserClientStub) GetClient() app.IsClient {
	return u.Client
}

func (u *UserClientStub) WriteJson(client app.ClientResponse) error {
	return nil
}
