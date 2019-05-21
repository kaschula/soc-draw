package app_test

import (
	"testing"

	"github.com/kaschula/socket-server/app"
	stubs "github.com/kaschula/socket-server/tests/stubs"
	. "github.com/stretchr/testify/assert"
)

func TestARoomStartsOnceEnoughClientHaveBeenAddedAndMessagesCanBeBroadcast(t *testing.T) {
	roomApp := NewRoomAppStub()
	user := app.NewUser("1")
	clientOne := &(stubs.ClientStub{
		"client:1",
		[]app.ClientResponse{},
		make(chan app.AppMessage),
		make(chan bool),
		[]app.Broadcasts{},
	})
	clientTwo := &(stubs.ClientStub{
		"client:2",
		[]app.ClientResponse{},
		make(chan app.AppMessage),
		make(chan bool),
		[]app.Broadcasts{},
	})

	userClientOne := app.NewUserClient(clientOne, user)
	userClientTwo := app.NewUserClient(clientTwo, user)

	room := app.NewDefaultRoom("1", "R2", roomApp)
	room.AddUserClient(userClientOne)
	room.AddUserClient(userClientTwo)

	message := app.AppMessage{app.ClientResponseTypes().ROOM_BROADCAST, "Payload"}

	room.Broadcast(nil, app.ClientAppMessage{message, "1"})

	Equal(t, clientOne.WrittenMessages[0].Payload, "Payload", "Client One Did not recieve correct Payload")
	Equal(t, clientOne.WrittenMessages[0].Type, app.ClientResponseTypes().ROOM_BROADCAST, "Client One Did not recieve correct Type")
	Equal(t, clientTwo.WrittenMessages[0].Payload, "Payload", "Client Two Did not recieve correct Payload")
	Equal(t, clientTwo.WrittenMessages[0].Type, app.ClientResponseTypes().ROOM_BROADCAST, "Client Two Did not recieve correct Type")
}

func TestARoomDoesNotBroadCastToClientsIfNotEnoughInRoom(t *testing.T) {
	roomApp := NewRoomAppStub()
	user := app.NewUser("1")
	clientOne := &(stubs.ClientStub{
		"client:1",
		[]app.ClientResponse{},
		make(chan app.AppMessage),
		make(chan bool),
		[]app.Broadcasts{},
	})
	clientTwo := &(stubs.ClientStub{
		"client:2",
		[]app.ClientResponse{},
		make(chan app.AppMessage),
		make(chan bool),
		[]app.Broadcasts{},
	})

	userClientOne := app.NewUserClient(clientOne, user)
	userClientTwo := app.NewUserClient(clientTwo, user)

	room := app.NewDefaultRoom("1", "R2", roomApp)
	room.AddUserClient(userClientOne)

	messageOne := app.AppMessage{app.ClientResponseTypes().ROOM_BROADCAST, "First Payload"}
	room.Broadcast(nil, app.ClientAppMessage{messageOne, "1"})

	room.AddUserClient(userClientTwo)

	messageTwo := app.AppMessage{app.ClientResponseTypes().ROOM_BROADCAST, "Second Payload"}
	room.Broadcast(nil, app.ClientAppMessage{messageTwo, "1"})

	Equal(t, clientOne.WrittenMessages[0].Payload, "Second Payload", "Client One Did not recieve correct Payload")
	Equal(t, clientTwo.WrittenMessages[0].Payload, "Second Payload", "Client Two Did not recieve correct Payload")
}

func TestThatAMessageWrittenToARoomIsPassedToTheRoomApp(t *testing.T) {
	roomApp := NewRoomAppStub()
	user := app.NewUser("1")
	clientOne := &(stubs.ClientStub{
		"client:1",
		[]app.ClientResponse{},
		make(chan app.AppMessage),
		make(chan bool),
		[]app.Broadcasts{},
	})
	clientTwo := &(stubs.ClientStub{
		"client:2",
		[]app.ClientResponse{},
		make(chan app.AppMessage),
		make(chan bool),
		[]app.Broadcasts{},
	})

	userClientOne := app.NewUserClient(clientOne, user)
	userClientTwo := app.NewUserClient(clientTwo, user)

	room := app.NewDefaultRoom("1", "R2", roomApp)
	message := app.ClientAppMessage{
		app.AppMessage{app.ClientResponseTypes().ROOM_BROADCAST, "Payload"},
		"1",
	}
	// should start room app
	room.AddUserClient(userClientOne)
	room.AddUserClient(userClientTwo)

	room.WriteMessage(message)

	Equal(t, roomApp.startCalled, 1, "Start() should have been called once")
	Equal(t, roomApp.writeCalled, 1, "WriteMessage() should have been called once")
	Equal(t, len(roomApp.written), 1, "There should one message written to the app")
	Equal(t, roomApp.written[0].GetPayload(), "Payload", "Payload should match orginal message")
}

type RoomAppStub struct {
	startCalled int
	writeCalled int
	written     []*app.RoomMessage
}

func NewRoomAppStub() *RoomAppStub {
	return &RoomAppStub{0, 0, []*app.RoomMessage{}}
}

func (a *RoomAppStub) Run() {}
func (a *RoomAppStub) Start(room app.RoomI) {
	a.startCalled++
}

func (a *RoomAppStub) WriteMessage(message app.RoomMessage) {
	a.writeCalled++
	a.written = append(a.written, &message)
}
