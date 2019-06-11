package app_test

import (
	"testing"

	"github.com/kaschula/socket-server/tests/stubs"

	"github.com/kaschula/socket-server/app"
	. "github.com/stretchr/testify/assert"
)

func TestThatNoLobbyReturnsErrorWhenGettingClient(t *testing.T) {
	lobby := app.NewNoLobby()

	client, err := lobby.GetClient("id")

	Nil(t, client, "client should be nil")
	Error(t, err, "NoLobby cant hold A Client. Should always return error")
}

// 	ResolveUserClient(IsClient) (UserClient, error)

func TestThatNoLobbyReturnsErrorWhenGettingUserClient(t *testing.T) {
	lobby := app.NewNoLobby()
	client := stubs.NewClientStub("1")

	userClient, err := lobby.ResolveUserClient(client)

	Nil(t, userClient, "client should be nil")
	Error(t, err, "NoLobby cant hold A Client. Should always return error")
}

func TestThatNoLobbyItsID(t *testing.T) {
	lobby := app.NewNoLobby()

	Equal(t, lobby.GetID(), app.NO_LOBBY_ID, "NoLobby cant hold A Client. Should always return error")
}
