package app_test

import (
	"testing"

	"github.com/kaschula/socket-server/app"
	"github.com/kaschula/socket-server/tests/stubs"
	. "github.com/stretchr/testify/assert"
)

func TestARoomStartsAfterBeingUpdated(t *testing.T) {

	user := &app.User{"U1"}
	userClient := &UserClientStub{user, nil}
	room := stubs.NewRoomStub("R1", "RoomStub")
	// roomRepository
	roomRepository := make(map[*app.User][]app.RoomI)
	roomRepository[user] = []app.RoomI{room}

	roomService := app.NewDefaultRoomService(roomRepository)

	err := roomService.AddUserClient(userClient, "R1")

	if err != nil {
		t.Fatal("RoomService Returned Error: ", err.Error())
	}

	Equal(t, 1, room.AddUserClientCount, "RoomStub Listen Should have been Called")
}

// type RoomStub struct {
// 	ID            string
// 	Name          string
// 	addUserClient int
// }

// func (r *RoomStub) AddUserClient(client app.UserClient) {
// 	r.addUserClient++
// }

// func (r *RoomStub) GetID() string {
// 	return r.ID
// }

// func (r *RoomStub) Broadcast(c app.IsClient, message app.ClientAppMessage) {}

// func (r *RoomStub) WriteMessage(message app.ClientAppMessage) {}

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