package stubs

import "github.com/kaschula/socket-server/app"

type UserClientStub struct {
	user   *app.User
	client app.IsClient
}

func (u *UserClientStub) GetUser() *app.User {
	return u.user
}

func (u *UserClientStub) GetClient() app.IsClient {
	return u.client
}

func (u *UserClientStub) WriteJson(client app.ClientResponse) error {
	return nil
}
