package app_test

import (
	"testing"

	"github.com/kaschula/soc-draw/app"
	"github.com/kaschula/soc-draw/tests/stubs"
	. "github.com/stretchr/testify/assert"
)

func TestARoomStartsAfterBeingUpdated(t *testing.T) {
	user := &app.User{"U1"}
	userClient := &UserClientStub{user, nil}
	room := stubs.NewRoomStub("R1", "RoomStub")
	roomRepository := make(map[*app.User][]app.RoomI)
	roomRepository[user] = []app.RoomI{room}

	roomService := app.NewDefaultRoomService(roomRepository)

	err := roomService.AddUserClient(userClient, "R1")

	if err != nil {
		t.Fatal("RoomService Returned Error: ", err.Error())
	}

	Equal(t, 1, room.AddUserClientCount, "RoomStub Listen Should have been Called")
}

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
