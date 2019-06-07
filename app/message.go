package app

import "errors"

// These are point less as before they are broadcasted they change to ROOM_BROADCAST Response
const (
	MESSAGE_TYPE_ROOM                = "ROOM"                // delete this
	MESSAGE_TYPE_ROOM_BROADCAST      = "ROOM_BROADCAST"      // this a duplicate
	MESSAGE_TYPE_ROOM_BROADCAST_INIT = "ROOM_BROADCAST_INIT" // this a duplicate
	MESSAGE_TYPE_ROOM_WELCOME        = "ROOM_WELCOME"        // Delete this
)

func NewAppMessage(messageType, payload string) AppMessage {
	return AppMessage{messageType, payload}
}

type AppMessage struct {
	Type    string `json:"messageType"`
	Payload string `json:"payload"`
}

var MessageTypeCreated = "CREATED"
var MessageTypeLobbyUserJoinRequest = "LOBBY_USER_JOIN_REQUEST"
var MessageTypeJoinRoom = "LOBBY_ROOM_REQUEST"

func IsLobbyMessage(message string) bool {
	messages := []string{MessageTypeLobbyUserJoinRequest, MessageTypeJoinRoom}
	for _, lobbyMessage := range messages {
		if lobbyMessage == message {
			return true
		}
	}

	return false
}

// This is in the wrong file
func welcomeMessage() ClientResponse {
	return ClientResponse{Type: "CREATED", Payload: "{}"}
}

type ClientAppMessage struct {
	Client IsClient // Change this to be the Client, This means Broadcast() just takes a Message
	AppMessage
	// ClientID string // Change this to be the Client, This means Broadcast() just takes a Message
}

func (c ClientAppMessage) GetClient() (IsClient, error) {
	if c.Client == nil {
		return nil, errors.New("Client Can not be resolved")
	}

	return c.Client, nil
}

func NewClientAppMessage(c IsClient, m AppMessage) ClientAppMessage {
	return ClientAppMessage{c, m}
}
