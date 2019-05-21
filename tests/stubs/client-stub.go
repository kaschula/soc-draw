package stubs

import (
	"github.com/kaschula/socket-server/app"
)

func NewClientStub(id string) *ClientStub {
	return &ClientStub{
		id,
		[]app.ClientResponse{},
		make(chan app.AppMessage),
		make(chan bool),
		[]app.Broadcasts{},
	}
}

type ClientStub struct {
	ID              string
	WrittenMessages []app.ClientResponse
	SendChan        chan app.AppMessage
	ReturnChan      chan bool
	Broadcasters    []app.Broadcasts
}

func (c *ClientStub) GetID() string {
	return c.ID
}

func (c *ClientStub) Listen() {
	for {
		appMessage := c.ReadMessage()
		message := app.ClientAppMessage{appMessage, c.GetID()}

		c.publish(message)
	}
}

func (c *ClientStub) publish(message app.ClientAppMessage) {
	for _, broadcaster := range c.Broadcasters {
		// should this be a in a routine
		broadcaster.Broadcast(c, message)
	}

	c.ReturnChan <- true
}

func (c *ClientStub) WriteJson(message app.ClientResponse) error {
	c.WrittenMessages = append(c.WrittenMessages, message)

	return nil
}

func (c *ClientStub) SendMessage(messageType, payload string) {
	c.SendChan <- app.AppMessage{Type: messageType, Payload: payload}
}

func (c *ClientStub) ReadMessage() app.AppMessage {
	message := <-c.SendChan

	return message
}

func (c *ClientStub) Subscribe(b app.Broadcasts) {
	c.Broadcasters = append(c.Broadcasters, b)
}

func (c *ClientStub) WaitForReturnChan() bool {
	return <-c.ReturnChan
}
