package app_test

import (
	"testing"

	"github.com/kaschula/socket-server/app"
	"github.com/kaschula/socket-server/tests/stubs"
	. "github.com/stretchr/testify/assert"
)

func TestRoomLobbyCanReturnID(t *testing.T) {
	lobby := app.NewRoomLobby(
		stubs.NewUserRepositoryStub(emptyUserMap()),
		stubs.NewRoomServiceStub(emptyUserRoomsMap()),
		app.NewInMemoryUserClientService(emptyClientUserClientsMap()),
	)

	Equal(t, app.ROOM_LOBBY_ID, lobby.GetID(), "Lobby ID should be ROOM_LOBBY_ID")
}

func TestRoomLobbyCanRemoveUserClient(t *testing.T) {
	c := &stubs.ClientStub{}
	uc := &stubs.UserClientStub{nil, c}
	userClientData := map[app.IsClient]app.UserClient{c: uc}

	lobby := app.NewRoomLobby(
		stubs.NewUserRepositoryStub(emptyUserMap()),
		stubs.NewRoomServiceStub(emptyUserRoomsMap()),
		app.NewInMemoryUserClientService(userClientData),
	)

	lobby.RemoveUserClient(uc)
	Equal(t, 0, len(userClientData), "UserClient should be removed")
}

func TestLobbyCanResolveUserClient(t *testing.T) {
	c := &stubs.ClientStub{}
	uc := &stubs.UserClientStub{nil, c}
	userClientData := map[app.IsClient]app.UserClient{c: uc}

	lobby := app.NewRoomLobby(
		stubs.NewUserRepositoryStub(emptyUserMap()),
		stubs.NewRoomServiceStub(emptyUserRoomsMap()),
		app.NewInMemoryUserClientService(userClientData),
	)

	resolvedUserClient, err := lobby.ResolveUserClient(c)

	Nil(t, err, "Lobby should not return error when resolved UserClient")
	Equal(t, uc, resolvedUserClient, "Lobby should be able to resolve UserClient based on Client")
}

func TestRoomLobbyCanAddRetrieveClient(t *testing.T) {
	client := stubs.NewClientStub("1")

	lobby := app.NewRoomLobby(
		stubs.NewUserRepositoryStub(emptyUserMap()),
		stubs.NewRoomServiceStub(emptyUserRoomsMap()),
		app.NewInMemoryUserClientService(emptyClientUserClientsMap()),
	)

	lobby.AddClient(client)
	retrievedClient, firstRetrievedErr := lobby.GetClient(client.GetID())

	lobby.Remove(client)
	retrievedClientAfterDelete, secondRetrievedErr := lobby.GetClient(client.GetID())

	Nil(t, firstRetrievedErr, "Lobby should be able to retrive clieng")
	Equal(t, client, retrievedClient, "retrived Client should match set up Client")

	Error(t, secondRetrievedErr, "lobby should be unable to resolve client")
	Nil(t, retrievedClientAfterDelete, "Retrived Client should be nil")
}

func TestIsLobbyMessage(t *testing.T) {
	lobby := app.NewRoomLobby(
		stubs.NewUserRepositoryStub(emptyUserMap()),
		stubs.NewRoomServiceStub(emptyUserRoomsMap()),
		app.NewInMemoryUserClientService(emptyClientUserClientsMap()),
	)

	True(t, lobby.IsLobbyMessage(app.GetRequestTypes().LOBBY_ROOM_REQUEST), "LOBBY_ROOM_REQUEST should resolve as lobby message")
	True(t, lobby.IsLobbyMessage(app.GetRequestTypes().LOBBY_USER_JOIN_REQUEST), "LOBBY_USER_JOIN_REQUEST should resolve as lobby message")
	False(t, lobby.IsLobbyMessage("None Lobby request"), "LOBBY_USER_JOIN_REQUEST should resolve as lobby message")
}

func emptyUserMap() map[string]*app.User {
	return map[string]*app.User{}
}

func emptyUserRoomsMap() map[*app.User][]app.RoomI {
	return map[*app.User][]app.RoomI{}
}

func emptyClientUserClientsMap() map[app.IsClient]app.UserClient {
	return map[app.IsClient]app.UserClient{}
}
