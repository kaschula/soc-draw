package app_test

import (
	"strings"
	"testing"

	"github.com/kaschula/socket-server/app"
	stubs "github.com/kaschula/socket-server/tests/stubs"
	. "github.com/stretchr/testify/assert"
)

func TestARoomStartsOnceEnoughClientHaveBeenAddedAndMessagesCanBeBroadcast(t *testing.T) {
	roomApp := NewRoomAppStub()
	user := app.NewUser("1")
	clientOne := newClientStub("client:1")
	clientTwo := newClientStub("client:2")

	userClientOne := app.NewUserClient(clientOne, user)
	userClientTwo := app.NewUserClient(clientTwo, user)

	room := app.NewDefaultRoom("1", "R2", roomApp)
	room.AddUserClient(userClientOne)
	room.AddUserClient(userClientTwo)

	message := app.AppMessage{app.GetResponseTypes().ROOM_BROADCAST, "Payload"}

	room.Broadcast(app.NewClientAppMessage(nil, message))

	Equal(t, clientOne.WrittenMessages[0].Payload, "Payload", "Client One Did not recieve correct Payload")
	Equal(t, clientOne.WrittenMessages[0].Type, app.GetResponseTypes().ROOM_BROADCAST, "Client One Did not recieve correct Type")
	Equal(t, clientTwo.WrittenMessages[0].Payload, "Payload", "Client Two Did not recieve correct Payload")
	Equal(t, clientTwo.WrittenMessages[0].Type, app.GetResponseTypes().ROOM_BROADCAST, "Client Two Did not recieve correct Type")
}

func TestARoomDoesNotBroadCastToClientsIfNotEnoughInRoom(t *testing.T) {
	roomApp := NewRoomAppStub()
	user := app.NewUser("1")
	clientOne := newClientStub("client:1")
	clientTwo := newClientStub("client:2")

	userClientOne := app.NewUserClient(clientOne, user)
	userClientTwo := app.NewUserClient(clientTwo, user)

	room := app.NewDefaultRoom("1", "R2", roomApp)
	room.AddUserClient(userClientOne)

	messageOne := app.AppMessage{app.GetResponseTypes().ROOM_BROADCAST, "First Payload"}
	room.Broadcast(app.NewClientAppMessage(nil, messageOne))

	room.AddUserClient(userClientTwo)

	messageTwo := app.AppMessage{app.GetResponseTypes().ROOM_BROADCAST, "Second Payload"}
	room.Broadcast(app.NewClientAppMessage(nil, messageTwo))

	Equal(t, true, strings.Contains(clientOne.WrittenMessages[0].Payload, `"running":"false"`),
		"Client One Did not recieve correct Payload",
	)
	Equal(t, clientOne.WrittenMessages[1].Payload, "Second Payload", "Incorrect Payload")
	Equal(t, clientTwo.WrittenMessages[0].Payload, "Second Payload", "Incorrect Payload")
}

func TestThatAMessageWrittenToARoomIsPassedToTheRoomApp(t *testing.T) {
	roomApp := NewRoomAppStub()
	user := app.NewUser("1")
	clientOne := newClientStub("client:1")
	clientTwo := newClientStub("client:2")

	userClientOne := app.NewUserClient(clientOne, user)
	userClientTwo := app.NewUserClient(clientTwo, user)

	room := app.NewDefaultRoom("1", "R2", roomApp)
	message := app.NewClientAppMessage(
		nil,
		app.AppMessage{app.GetResponseTypes().ROOM_BROADCAST, "Payload"},
	)
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

func newClientStub(id string) *stubs.ClientStub {
	return &(stubs.ClientStub{
		id,
		[]app.ClientResponse{},
		make(chan app.AppMessage),
		make(chan bool),
		[]app.Broadcasts{},
		nil,
	})
}
