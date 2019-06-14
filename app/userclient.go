package app

type UserClient interface {
	GetUser() *User
	WriteJson(message ClientResponse) error
	GetClient() IsClient
}

type DefaultUserClient struct {
	client IsClient
	user   *User
}

func NewUserClient(client IsClient, user *User) UserClient {
	return &DefaultUserClient{client, user}
}

func (uc *DefaultUserClient) GetUser() *User {
	return uc.user
}

func (uc *DefaultUserClient) GetClient() IsClient {
	return uc.client
}

func (uc *DefaultUserClient) WriteJson(message ClientResponse) error {
	return uc.client.WriteJson(message)
}
