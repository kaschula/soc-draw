package app

import "errors"

type UserClientService interface {
	CreateAndStoreUserClient(user *User, client IsClient) error
	Resolve(IsClient) (UserClient, error)
}

func NewInMemoryUserClientService() UserClientService {
	return &InMemoryUserClientService{make(map[IsClient]UserClient)}
}

type InMemoryUserClientService struct {
	userClients map[IsClient]UserClient
}

func (s *InMemoryUserClientService) CreateAndStoreUserClient(user *User, client IsClient) error {
	s.userClients[client] = NewUserClient(client, user)

	return nil
}

func (s *InMemoryUserClientService) Resolve(client IsClient) (UserClient, error) {
	userClient, ok := s.userClients[client]
	if !ok {
		return userClient, errors.New("UserClient could not be resolved")
	}

	return userClient, nil
}

// User Client struct
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
