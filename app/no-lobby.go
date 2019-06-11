package app

import "errors"

const NO_LOBBY_ID = "NO_LOBBY"

type NoLobby struct {
	id string
}

func NewNoLobby() Lobby {
	return &NoLobby{NO_LOBBY_ID}
}

// func (l *NoLobby)

///type Lobby interface

func (l *NoLobby) AddClient(client IsClient) {}

func (l *NoLobby) GetClient(id string) (IsClient, error) {
	return nil, errors.New("NoLobby can not hold IsClient")
}

func (l *NoLobby) Remove(IsClient) {}

func (l *NoLobby) ResolveUserClient(IsClient) (UserClient, error) {
	return nil, errors.New("NoLobby can not hold UserClient")
}

func (l *NoLobby) Broadcast(message ClientAppMessage) {}

func (l *NoLobby) GetID() string {
	return l.id
}

func (l *NoLobby) RemoveUserClient(UserClient) {}
