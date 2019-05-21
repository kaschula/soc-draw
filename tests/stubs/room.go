package stubs

import "github.com/kaschula/socket-server/app"

func NewRoomStub(id, name string) *RoomStub {
	return &RoomStub{
		id,
		name,
		0,
		0,
		make(chan bool),
		[]app.ClientAppMessage{},
	}
}

type RoomStub struct {
	ID                 string
	Name               string
	AddUserClientCount int
	BroadcastCalled    int
	BroadcastReturn    chan bool
	BroadcastData      []app.ClientAppMessage
}

func (r *RoomStub) AddUserClient(client app.UserClient) {
	r.AddUserClientCount++
}

func (r *RoomStub) GetID() string {
	return r.ID
}

func (r *RoomStub) Broadcast(c app.IsClient, message app.ClientAppMessage) {
	r.BroadcastCalled++
	r.BroadcastData = append(r.BroadcastData, message)
	r.BroadcastReturn <- true
}

func (r *RoomStub) WriteMessage(message app.ClientAppMessage) {}

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
