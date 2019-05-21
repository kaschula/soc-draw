package app

import (
	"encoding/json"
	"errors"
	"fmt"
	// "github.com/kaschula/socket-server/app"
)

type Lobby interface {
	AddClient(client IsClient)
	GetClient(id string) (IsClient, error)
	Broadcasts
}

func NewLobby(userRepository UserRepository, roomRepository RoomService, userClientService UserClientService) Lobby {
	return &RoomLobby{
		// This should be a client store
		make(map[string]IsClient),
		userRepository,
		roomRepository,
		userClientService,
	}
}

type RoomLobby struct {
	clients           map[string]IsClient
	userRepository    UserRepository
	roomRepository    RoomService
	userClientService UserClientService
}

func (l *RoomLobby) AddClient(client IsClient) {
	l.clients[client.GetID()] = client

	client.Subscribe(l)
}

func (l *RoomLobby) GetClient(id string) (IsClient, error) {
	client, ok := l.clients[id]

	if !ok || client == nil {
		return nil, fmt.Errorf("Client can not be found. OK: %v, client: %v", ok, client)
	}

	return client, nil
}

func (l *RoomLobby) Broadcast(client IsClient, message ClientAppMessage) {
	if !IsLobbyMessage(message.Type) {
		// To Test, use Client repo stub to test if the Client repository was called
		// make lobby_test unit test
		fmt.Printf("Message type of %v is not Lobby Message \n", message.Type)
		return
	}

	switch message.Type {
	case MessageTypeLobbyUserJoinRequest:
		l.resolveUser(client, message.Payload)
		return
	case MessageTypeJoinRoom:
		l.joinRoom(client, message.Payload)
		return
	default:
		return
	}
}

func (l *RoomLobby) resolveUser(client IsClient, messagePayload string) {
	var payload struct {
		UserID string `json:"user"`
	}

	raw := []byte(messagePayload)
	json.Unmarshal(
		raw,
		&payload,
	)

	if payload.UserID == "" {
		client.WriteJson(newErrorResponse(ClientResponseErrorType().NO_USER))
		return
	}

	user, err := l.getUser(payload.UserID)
	if err != nil {
		client.WriteJson(newErrorResponse(ClientResponseErrorType().USER))
		return
	}

	lobbyData, err := l.getLobbyData(user)
	if err != nil {
		client.WriteJson(newErrorResponse(ClientResponseErrorType().LOBBY_DATA))
		return
	}

	if err := l.createAndStoreUserClient(user, client); err != nil {
		client.WriteJson(newErrorResponse(ClientResponseErrorType().USER_CLIENT))
		return
	}

	_ = client.WriteJson(newUserResolvedMessage(lobbyData))
}

func (l *RoomLobby) joinRoom(client IsClient, messagePayload string) {

	userClient, err := l.userClientService.Resolve(client)
	if err != nil {
		client.WriteJson(newErrorResponse("USER_CLIENT_404"))
		return
	}

	roomId, err := l.resolveRoomId(messagePayload)
	if err != nil {
		client.WriteJson(newErrorResponse(ClientResponseErrorType().PAYLOAD_ROOM_ID))
		return
	}

	if !l.roomRepository.CanUserJoin(userClient, roomId) {
		client.WriteJson(newErrorResponse(ClientResponseErrorType().USER_ROOM_AUTH))
		return
	}

	if err := l.roomRepository.AddUserClient(userClient, roomId); err != nil {
		client.WriteJson(newErrorResponse(ClientResponseErrorType().ADD_USER_TO_ROOM))
		return
	}

	userClient.WriteJson(newUserJoinRoomMessage())
}

func (l *RoomLobby) resolveRoomId(messagePayload string) (string, error) {
	var payload struct {
		RoomId string `json:"roomId"`
	}

	raw := []byte(messagePayload)
	json.Unmarshal(
		raw,
		&payload,
	)

	if payload.RoomId == "" {
		// Untested
		return "", errors.New("Can not resolve Room ID")
	}

	return payload.RoomId, nil
}

func (l *RoomLobby) getUser(userId string) (*User, error) {
	return l.userRepository.GetUser(userId)
}

func (l *RoomLobby) getLobbyData(user *User) (*LobbyData, error) {
	// use repository
	rooms, err := l.roomRepository.GetRooms(user)
	if err != nil {
		return nil, err
	}

	return &LobbyData{
			User:  *user,
			Rooms: rooms,
		},
		nil
}

func (l *RoomLobby) getClient(clientId string) (IsClient, error) {
	return l.GetClient(clientId)
}

func (l *RoomLobby) createAndStoreUserClient(user *User, client IsClient) error {
	return l.userClientService.CreateAndStoreUserClient(user, client)
}

func newUserResolvedMessage(ld *LobbyData) ClientResponse {
	b, err := json.Marshal(ld)
	if err != nil {
		return newErrorResponse(err.Error())
	}

	return ClientResponse{Type: ClientResponseTypes().USER_LOBBY_DATA, Payload: string(b)}
}

func newUserJoinRoomMessage() ClientResponse {
	return ClientResponse{Type: ClientResponseTypes().USER_JOINED_ROOM, Payload: "success"}
}

type LobbyData struct {
	User  User    `json:"User"`
	Rooms []RoomI `json:"Rooms"`
}
