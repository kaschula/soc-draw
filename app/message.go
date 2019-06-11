package app

import "errors"

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
	Client IsClient
	AppMessage
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
