package app_test

import (
	"encoding/json"
	"testing"

	"github.com/kaschula/socket-server/app"
	. "github.com/stretchr/testify/assert"
)

func TestThatTheClientCanWriteAMessageToTheSocket(t *testing.T) {
	socket := NewJsonReaderWriterStub()
	client := app.NewDefaultClient("C:1", socket)
	message := app.ClientResponse{"Type", "Payload", "R:1"}

	client.WriteJson(message)

	Equal(t, 1, len(socket.written), "Expect Only one message to be written")

}

func TestThatTheClientCanReturnAnId(t *testing.T) {
	socket := NewJsonReaderWriterStub()
	client := app.NewDefaultClient("C:1", socket)

	Equal(t, "C:1", client.GetID(), "Incorrect client ID")
}

func TestThatAListeningClientCanRecieveAMessagePublishToItsBroadcasters(t *testing.T) {
	broadcaster := NewBroadcasterStub()
	socket := NewJsonReaderWriterStub()
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

func NewJsonReaderWriterStub() *JsonReaderWriterStub {
	return &JsonReaderWriterStub{
		[]interface{}{},
		make(chan string),
		make(chan bool),
	}
}

type JsonReaderWriterStub struct {
	written  []interface{}
	readChan chan string
	doneChan chan bool
}

func (w *JsonReaderWriterStub) WriteJSON(v interface{}) error {
	w.written = append(w.written, v)
	return nil
}

func (w *JsonReaderWriterStub) SendMessage(jsonMessage string) {
	w.readChan <- jsonMessage
}

func (w *JsonReaderWriterStub) ReadJSON(v interface{}) error {
	jsonData := <-w.readChan

	return json.Unmarshal([]byte(jsonData), v)
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
