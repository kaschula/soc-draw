package app

import "fmt"

// A struct that holds websocket client
// type Client struct {
// }

type IsClient interface {
	GetID() string
	Listen()
	WriteJson(message ClientResponse) error
	Subscribe(b Broadcasts)
}

func NewDefaultClient(id string, socket Socket) IsClient {
	return &DefaultClient{
		id,
		[]Broadcasts{},
		socket,
	}
}

type DefaultClient struct {
	id           string
	broadcasters []Broadcasts
	socket       Socket
}

func (c *DefaultClient) GetID() string {
	return c.id
}

func (c *DefaultClient) Listen() {
	for {
		var msg AppMessage

		if err := c.socket.ReadJSON(&msg); err != nil {
			fmt.Println("Error: reading from socket: ", err.Error())
			break
		}

		for _, broadcaster := range c.broadcasters {
			broadcaster.Broadcast(c, NewClientAppMessage(c.GetID(), msg))
		}
	}
}

func (c *DefaultClient) WriteJson(message ClientResponse) error {
	return c.socket.WriteJSON(&message)
}

func (c *DefaultClient) Subscribe(b Broadcasts) {
	c.broadcasters = append(c.broadcasters, b)
}
