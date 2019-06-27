package stubs

import "github.com/kaschula/soc-draw/app"

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

func (r *RoomStub) AddUserClient(client app.UserClient) error {
	r.AddUserClientCount++
	return nil
}

func (r *RoomStub) GetID() string {
	return r.ID
}

func (r *RoomStub) Broadcast(message app.ClientAppMessage) {
	r.BroadcastCalled++
	r.BroadcastData = append(r.BroadcastData, message)
	r.BroadcastReturn <- true
}

func (r *RoomStub) WriteMessage(message app.ClientAppMessage) {}

func (r *RoomStub) RemoveUserClient(us app.UserClient) {
}
