package app

import (
	"net/http"

	"github.com/kaschula/random"
)

type ClientService interface {
	UpgradeToSocket(res http.ResponseWriter, req *http.Request) (IsClient, error)
}

func NewDefaultClientService(factory UpgraderFactory) *DefaultClientService {
	return &DefaultClientService{random.NewRandomStringGenerator(20), factory}
}

type DefaultClientService struct {
	idGenerator *random.RandomStringGenerator
	factory     UpgraderFactory
}

func (cs *DefaultClientService) UpgradeToSocket(res http.ResponseWriter, req *http.Request) (IsClient, error) {
	upgrader := cs.factory.NewUpgrader()

	socket, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		return nil, err
	}

	hash, err := cs.idGenerator.Generate()
	if err != nil {
		return nil, err
	}

	return NewDefaultClient("CLIENT:"+hash, socket), nil
}
