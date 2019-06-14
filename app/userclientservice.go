package app

import "errors"

type UserClientService interface {
	CreateAndStoreUserClient(user *User, client IsClient) error
	Resolve(IsClient) (UserClient, error)
	Delete(UserClient)
}

func NewInMemoryUserClientService(userClients map[IsClient]UserClient) UserClientService {
	if userClients == nil {
		userClients = map[IsClient]UserClient{}
	}

	return &InMemoryUserClientService{userClients}
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

func (s *InMemoryUserClientService) Delete(uc UserClient) {
	delete(s.userClients, uc.GetClient())
}
