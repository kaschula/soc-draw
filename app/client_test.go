package app_test

import (
	"errors"
	"testing"

	"github.com/kaschula/socket-server/app"
	"github.com/kaschula/socket-server/tests/stubs"
	. "github.com/stretchr/testify/assert"
)

func TestThatTheClientCanWriteAMessageToTheSocket(t *testing.T) {
	socket := stubs.NewSocketStub(nil)
	client := app.NewDefaultClient("C:1", socket)
	message := app.ClientResponse{"Type", "Payload", "R:1"}

	client.WriteJson(message)

	Equal(t, 1, len(socket.Written), "Expect Only one message to be written")

}

func TestThatTheClientCanReturnAnId(t *testing.T) {
	socket := stubs.NewSocketStub(nil)
	client := app.NewDefaultClient("C:1", socket)

	Equal(t, "C:1", client.GetID(), "Incorrect client ID")
}

func TestThatASocketErrorIsReturnedDuringWrite(t *testing.T) {
	expectedError := errors.New("Error")
	s := stubs.NewSocketStub(expectedError)
	c := app.NewDefaultClient("C:1", s)

	err := c.WriteJson(app.ClientResponse{})

	Error(t, expectedError)
	Equal(t, expectedError.Error(), err.Error(), "Incorrect client ID")

}

func TestThatAListeningClientCanRecieveAMessagePublishToItsBroadcasters(t *testing.T) {
	broadcaster := NewBroadcasterStub()
	socket := stubs.NewSocketStub(nil)
	client := app.NewDefaultClient("C:1", socket)
	message := `{"messageType":"FakeType", "payload":"FakePayload"}`

	client.Subscribe(broadcaster)
	go client.Listen()
	socket.SendMessage(message)
	// wait
	<-broadcaster.doneChan

	Equal(t, len(broadcaster.broadcasts), 1, "Should have at least 1 broadcast")
	Equal(t, broadcaster.broadcasts[0].AppMessage.Type, "FakeType", "Message should have been broadcast")
	Equal(t, broadcaster.broadcasts[0].AppMessage.Payload, "FakePayload", "Message should have been broadcast")

}

func NewBroadcasterStub() *BroadcasterStub {
	return &BroadcasterStub{[]app.ClientAppMessage{}, make(chan bool)}
}

type BroadcasterStub struct {
	broadcasts []app.ClientAppMessage
	doneChan   chan bool
}

func (b *BroadcasterStub) Broadcast(message app.ClientAppMessage) {
	b.broadcasts = append(b.broadcasts, message)
	b.doneChan <- true
}

func (b *BroadcasterStub) GetID() string {
	return ""
}

func (b *BroadcasterStub) RemoveUserClient(us app.UserClient) {
}
