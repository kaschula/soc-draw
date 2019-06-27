package app_test

import (
	"strings"
	"testing"

	"github.com/kaschula/soc-draw/app"
	stubs "github.com/kaschula/soc-draw/tests/stubs"
	. "github.com/stretchr/testify/assert"
)

func TestARoomStartsOnceEnoughClientHaveBeenAddedAndMessagesCanBeBroadcast(t *testing.T) {
	roomApp := NewRoomAppStub()
	u1 := app.NewUser("1")
	u2 := app.NewUser("2")
	clientOne := newClientStub("client:1")
	clientTwo := newClientStub("client:2")

	userClientOne := app.NewUserClient(clientOne, u1)
	userClientTwo := app.NewUserClient(clientTwo, u2)

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
	u1 := app.NewUser("1")
	u2 := app.NewUser("2")
	clientOne := newClientStub("client:1")
	clientTwo := newClientStub("client:2")

	userClientOne := app.NewUserClient(clientOne, u1)
	userClientTwo := app.NewUserClient(clientTwo, u2)

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

func TestThatTheSameUserWithDifferentClientsCantJoinTwice(t *testing.T) {
	roomApp := NewRoomAppStub()
	user := app.NewUser("1")
	clientOne := newClientStub("client:1")
	clientTwo := newClientStub("client:2")

	userClientOne := app.NewUserClient(clientOne, user)
	userClientTwo := app.NewUserClient(clientTwo, user)

	room := app.NewDefaultRoom("1", "R2", roomApp)

	errOne := room.AddUserClient(userClientOne)
	errTwo := room.AddUserClient(userClientTwo)

	Nil(t, errOne, "UserClient one should be able to join")
	Error(t, errTwo, "UserClient two shouldnt be able to join because User has already joined")
}

func TestThatAUserClientCanBeRemovedFromTheRoom(t *testing.T) {
	roomApp := NewRoomAppStub()
	user := app.NewUser("1")
	clientOne := newClientStub("client:1")
	uc := app.NewUserClient(clientOne, user)

	room := app.NewDefaultRoom("1", "R2", roomApp)

	errOne := room.AddUserClient(uc)
	errTwo := room.AddUserClient(uc)
	room.RemoveUserClient(uc)
	errThree := room.AddUserClient(uc)

	Nil(t, errOne, "UserClient should be able to join with first attempt")
	Error(t, errTwo, "UserClient two shouldnt be able to join with second attemp")
	Nil(t, errThree, "UserClient should be able again after being removed")
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
