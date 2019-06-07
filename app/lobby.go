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
	Remove(IsClient)
	ResolveUserClient(IsClient) (UserClient, error)
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

func (l *RoomLobby) GetID() string {
	return "Lobby"
}

func (l *RoomLobby) AddClient(client IsClient) {
	l.clients[client.GetID()] = client

	client.SubscribeLobby(l)
}

// Untested
func (l *RoomLobby) Remove(client IsClient) {
	delete(l.clients, client.GetID())
}

// Untested
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
	fmt.Println("Lobby::Broadcast()")
	// This may not be needed any more as this check is don't in the client
	if !IsLobbyMessage(message.Type) {
		// To Test, use Client repo stub to test if the Client repository was called
		// make lobby_test unit test
		fmt.Printf("Message type of %v is not Lobby Message \n", message.Type)
		return
	}

	client, err := message.GetClient()
	if err != nil {
		fmt.Println("Error: No client resolved from message")
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
		// Should log unreconhgised lobby message type
		return
	}
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
	fmt.Println("Lobby::joinRoom()")
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

	fmt.Println("Lobby::joinRoom() Can user join ")
	if !l.roomRepository.CanUserJoin(userClient, roomId) {
		fmt.Println("Lobby::joinRoom() user cant join ")

		client.WriteJson(newErrorResponse(ClientResponseErrorType().USER_ROOM_AUTH))
		return
	}

	fmt.Println("Lobby::joinRoom() Adding user client to repo ")
	if err := l.roomRepository.AddUserClient(userClient, roomId); err != nil {
		fmt.Println("Lobby::joinRoom() error adding UserClient to Room Repo", err)
		client.WriteJson(newErrorResponse(ClientResponseErrorType().ADD_USER_TO_ROOM))
		return
	}

	fmt.Printf(
		"Lobby::joinRoom()::roomId: %#v, userId: %#v, clientId: %#v \n",
		roomId,
		userClient.GetUser().ID,
		userClient.GetClient().GetID(),
	)
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
	return l.userRepository.GetUser(userId)
}

func (l *RoomLobby) getLobbyData(user *User) (*LobbyData, error) {
	fmt.Println("Lobby::getLobbyData()")
	// use repository
	rooms, err := l.roomRepository.GetRooms(user)
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

	return ClientResponse{Type: ClientResponseTypes().USER_LOBBY_DATA, Payload: string(b)}
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

	return ClientResponse{Type: ClientResponseTypes().USER_JOINED_ROOM, Payload: string(b), RoomID: roomId}
}

// Untests
func (l *RoomLobby) RemoveUserClient(uc UserClient) {
	// To Implement when needed
}

type LobbyData struct {
	User  User    `json:"User"`
	Rooms []RoomI `json:"Rooms"`
}
