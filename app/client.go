package app

import (
	"encoding/json"
	"fmt"
)

// A struct that holds websocket client
// type Client struct {
// }

type IsClient interface {
	GetID() string
	Listen()
	WriteJson(message ClientResponse) error
	Subscribe(b Broadcasts)
	SubscribeLobby(b Lobby)
}

func NewDefaultClient(id string, socket Socket) IsClient {
	return &DefaultClient{
		id,
		nil, // Replace this with NotLobby, and empty implementation of Broadcast
		[]Broadcasts{},
		socket,
	}
}

type DefaultClient struct {
	id           string
	lobby        Lobby // Use a NotLobby on construct Init
	broadcasters []Broadcasts
	socket       Socket
}

func (c *DefaultClient) GetID() string {
	return c.id
}

func (c *DefaultClient) Listen() {
	for {
		var msg AppMessage

		// c.socket

		if err := c.socket.ReadJSON(&msg); err != nil {
			fmt.Println("Error: reading from socket: ", err.Error())
			// fmt.Printf("Client %#v, lobby %#v", c, c.lobby)
			c.close()
			break
		}

		fmt.Printf(
			"Client::Listen() Client: %#v. AppMessage: %#v. Broadcasters: %#v. \n",
			c.GetID(), msg, c.broadcasters,
		)
		if IsLobbyMessage(msg.Type) {
			fmt.Println("Client::Listen() message is lobby type")
			// creating the client message is repeated below
			c.lobby.Broadcast(c, NewClientAppMessage(c.GetID(), msg))
		}

		roomId := getRoomId(msg.Payload)

		for _, broadcaster := range c.broadcasters {
			fmt.Println("Broadcast ID", broadcaster.GetID())
			if broadcaster.GetID() == roomId {
				broadcaster.Broadcast(c, NewClientAppMessage(c.GetID(), msg))
			}
		}
	}
}
func (c *DefaultClient) close() {
	c.lobby.Remove(c)

	userClient, err := c.lobby.ResolveUserClient(c)
	if err != nil {
		return
	}

	for _, broadcaster := range c.broadcasters {
		fmt.Println("Removing UserClients from broadcasters", broadcaster.GetID())
		broadcaster.RemoveUserClient(userClient)
	}
}

func (c *DefaultClient) WriteJson(message ClientResponse) error {
	err := c.socket.WriteJSON(&message)

	fmt.Println("Client::WriteJson() ID", c.GetID())
	fmt.Println("Client::WriteJson() message", message)
	fmt.Println("Client::WriteJson() err", err)
	fmt.Println("--------------------------------------")

	return err
	// return c.socket.WriteJSON(&message)
}

func (c *DefaultClient) Subscribe(b Broadcasts) {
	// use ID to create a map of instead of slice
	c.broadcasters = append(c.broadcasters, b)
}

func (c *DefaultClient) SubscribeLobby(l Lobby) {
	// maybe this could be done on construction
	c.lobby = l
}

// Untested
func getRoomId(payload string) string {
	var roomPayload struct {
		RoomId string `json:"roomId"`
	}

	// var payload struct {
	// 	UserID string `json:"user"`
	// }

	raw := []byte(payload)
	json.Unmarshal(
		raw,
		&roomPayload,
	)

	return roomPayload.RoomId
}
