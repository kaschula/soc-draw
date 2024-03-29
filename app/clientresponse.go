package app

func newErrorResponse(errorType string) ClientResponse {
	return ClientResponse{"ERROR", GetResponseErrorMessages(errorType), ""}
}

type ClientResponse struct {
	Type    string
	Payload string
	RoomID  string
}

type responseTypes struct {
	CREATED                string
	USER_LOBBY_DATA        string
	USER_JOINED_ROOM       string
	ROOM_BROADCAST         string
	ROOM_BROADCAST_MESSAGE string
	ROOM_BROADCAST_INIT    string
	ERROR                  string
}

func GetResponseTypes() responseTypes {
	return responseTypes{
		"CREATED",
		"USER_LOBBY_DATA",
		"USER_JOINED_ROOM",
		"ROOM_BROADCAST",
		"ROOM_BROADCAST_MESSAGE",
		"ROOM_BROADCAST_INIT",
		"ERROR",
	}
}

type ErrorTypes struct {
	USER             string
	NO_USER          string
	LOBBY_DATA       string
	USER_CLIENT      string
	PAYLOAD_ROOM_ID  string
	USER_ROOM_AUTH   string
	ADD_USER_TO_ROOM string
}

func ClientResponseErrorType() ErrorTypes {
	return ErrorTypes{
		USER:             "USER",
		NO_USER:          "NO_USER",
		LOBBY_DATA:       "LOBBY_DATA",
		USER_CLIENT:      "USER_CLIENT",
		PAYLOAD_ROOM_ID:  "PAYLOAD_ROOM_ID",
		USER_ROOM_AUTH:   "USER_ROOM_AUTH",
		ADD_USER_TO_ROOM: "ADD_USER_TO_ROOM",
	}
}

func GetResponseErrorMessages(key string) string {
	messages := map[string]string{
		"USER":             "Cannot resolve user",
		"NO_USER":          "No user ID in payload request",
		"LOBBY_DATA":       "User Lobby Data can not be resolved",
		"USER_CLIENT":      "UserClient Could not be created",
		"USER_CLIENT_404":  "UserClient Could not be found",
		"PAYLOAD_ROOM_ID":  "Payload Room ID is invalid",
		"USER_ROOM_AUTH":   "User Not Authorized to Join Room",
		"ADD_USER_TO_ROOM": "Failed to add user to room",
	}

	return messages[key]
}

func NewRoomResponse(message ClientAppMessage, roomId string) ClientResponse {
	return ClientResponse{message.Type, message.Payload, roomId}
}

func NewRoomWaitingToStart(roomId string) ClientResponse {
	return ClientResponse{
		GetResponseTypes().ROOM_BROADCAST,
		`{"running":"false", "message":"waiting for room to start"}`,
		roomId,
	}
}

func welcomeMessage() ClientResponse {
	return ClientResponse{Type: "CREATED", Payload: "{}"}
}
