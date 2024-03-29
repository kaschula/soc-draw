package app

import (
	"encoding/json"
	"errors"
	"fmt"
)

const ROOM_LOBBY_ID = "ROOM_LOBBY"

type Lobby interface {
	AddClient(client IsClient)
	GetClient(id string) (IsClient, error)
	Remove(IsClient)
	ResolveUserClient(IsClient) (UserClient, error)
	IsLobbyMessage(string) bool
	Broadcasts
}

func NewRoomLobby(userService UserService, roomService RoomService, userClientService UserClientService) Lobby {
	return &RoomLobby{
		make(map[string]IsClient),
		userService,
		roomService,
		userClientService,
	}
}

type RoomLobby struct {
	clients           map[string]IsClient
	userService       UserService
	roomService       RoomService
	userClientService UserClientService
}

func (l *RoomLobby) GetID() string {
	return ROOM_LOBBY_ID
}

func (l *RoomLobby) AddClient(client IsClient) {
	l.clients[client.GetID()] = client

	client.SubscribeLobby(l)
}

func (l *RoomLobby) Remove(client IsClient) {
	delete(l.clients, client.GetID())
}

func (l *RoomLobby) ResolveUserClient(c IsClient) (UserClient, error) {
	return l.userClientService.Resolve(c)
}

func (l *RoomLobby) GetClient(id string) (IsClient, error) {
	client, ok := l.clients[id]

	if !ok || client == nil {
		return nil, fmt.Errorf("Client can not be found. OK: %v, client: %v", ok, client)
	}

	return client, nil
}

func (l *RoomLobby) Broadcast(message ClientAppMessage) {
	if !l.IsLobbyMessage(message.Type) {
		fmt.Printf("Message type of %v is not Lobby Message \n", message.Type)
		return
	}

	client := message.GetClient()

	t := GetRequestTypes()

	switch message.Type {
	case t.LOBBY_USER_JOIN_REQUEST:
		l.resolveUser(client, message.Payload)
		return
	case t.LOBBY_ROOM_REQUEST:
		l.joinRoom(client, message.Payload)
		return
	default:
		// Should log unreconhgised lobby message type
		return
	}
}

func (l *RoomLobby) IsLobbyMessage(m string) bool {
	return isLobbyMessage(m)
}

func (l *RoomLobby) RemoveUserClient(uc UserClient) {
	l.userClientService.Delete(uc)
}

func (l *RoomLobby) resolveUser(client IsClient, messagePayload string) {
	fmt.Println("lobby::resolveUser()")
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

	fmt.Printf("Lobby::resolveUser():: resolved user writing to socket \n")
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

	if !l.roomService.CanUserJoin(userClient, roomId) {
		client.WriteJson(newErrorResponse(ClientResponseErrorType().USER_ROOM_AUTH))
		return
	}

	if err := l.roomService.AddUserClient(userClient, roomId); err != nil {
		client.WriteJson(newErrorResponse(ClientResponseErrorType().ADD_USER_TO_ROOM))
		return
	}

	userClient.WriteJson(newUserJoinRoomMessage(roomId))
}

func (l *RoomLobby) resolveRoomId(messagePayload string) (string, error) {
	fmt.Println("Lobby::resolveRoomId()")
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
	return l.userService.GetUser(userId)
}

func (l *RoomLobby) getLobbyData(user *User) (*LobbyData, error) {
	fmt.Println("Lobby::getLobbyData()")

	rooms, err := l.roomService.GetRooms(user)
	if err != nil {
		return nil, err
	}

	fmt.Println("Lobby::getLobbyData() return")
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

	return ClientResponse{Type: GetResponseTypes().USER_LOBBY_DATA, Payload: string(b)}
}

func newUserJoinRoomMessage(roomId string) ClientResponse {
	data := struct {
		RoomId string
	}{
		roomId,
	}
	b, err := json.Marshal(data)
	if err != nil {
		return newErrorResponse(err.Error())
	}

	return ClientResponse{Type: GetResponseTypes().USER_JOINED_ROOM, Payload: string(b), RoomID: roomId}
}

func isLobbyMessage(m string) bool {
	return m == GetRequestTypes().LOBBY_USER_JOIN_REQUEST ||
		m == GetRequestTypes().LOBBY_ROOM_REQUEST
}

type LobbyData struct {
	User  User    `json:"User"`
	Rooms []RoomI `json:"Rooms"`
}
