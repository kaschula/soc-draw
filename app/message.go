package app

// These are point less as before they are broadcasted they change to ROOM_BROADCAST Response
const (
	MESSAGE_TYPE_ROOM         = "ROOM"
	MESSAGE_TYPE_ROOM_WELCOME = "ROOM_WELCOME" // Delete this
	// MESSAGE_TYPE_ROOM_WELCOME = "ROOM_"
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
	AppMessage
	ClientID string // Change this to be the Client, This means Broadcast() just takes a Message
}

func NewClientAppMessage(clientId string, appMessage AppMessage) ClientAppMessage {
	return ClientAppMessage{appMessage, clientId}
}
