package app

type requestTypes struct {
	LOBBY_USER_JOIN_REQUEST string
	ROOM_REQUEST            string
	LOBBY_ROOM_REQUEST      string
}

func GetRequestTypes() requestTypes {
	return requestTypes{
		"LOBBY_USER_JOIN_REQUEST",
		"ROOM_REQUEST",
		"LOBBY_ROOM_REQUEST",
	}
}

func NewAppMessage(messageType, payload string) AppMessage {
	return AppMessage{messageType, payload}
}

type AppMessage struct {
	Type    string `json:"messageType"`
	Payload string `json:"payload"`
}

type ClientAppMessage struct {
	client IsClient
	AppMessage
}

// Untested validation
func NewClientAppMessage(c IsClient, m AppMessage) ClientAppMessage {
	if c == nil {
		c = NewNoClient()
	}

	return ClientAppMessage{c, m}
}

func (c ClientAppMessage) GetClient() IsClient {
	return c.client
}
