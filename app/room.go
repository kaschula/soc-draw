package app

import (
	"errors"
	"fmt"
)

type RoomI interface {
	AddUserClient(UserClient) error
	Broadcasts
}

type BaseRoom struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`
	// Todo: Change this to Map
	clients        map[string]UserClient
	running        bool
	minimumClients int
	maxClients     int
	roomApp        RoomApplication
}

// Creates room with default configuration
func NewDefaultRoom(id, name string, roomApp RoomApplication) RoomI {
	return NewRoom(id, name, 2, 5, roomApp)
}

func NewRoom(id, name string, minimumClients, maxClients int, roomApp RoomApplication) RoomI {
	return &BaseRoom{
		id,
		name,
		map[string]UserClient{},
		false,
		minimumClients,
		maxClients,
		roomApp,
	}
}

func (r *BaseRoom) AddUserClient(uc UserClient) error {
	userId := uc.GetUser().ID

	_, ok := r.clients[userId]
	if ok {
		return errors.New("User is already part of room")
	}

	r.clients[userId] = uc

	client := uc.GetClient()
	client.Subscribe(r)
	r.updated()

	return nil
}

func (r *BaseRoom) updated() {
	if !r.isReady() {
		fmt.Println("Room is not ready")
		if r.isRunning() {
			fmt.Println("Room::updated() Room is running and not ready: room should stop")
			// r.stop()
		}

		return
	}

	r.start()
}

func (r *BaseRoom) isReady() bool {
	clientsCount := len(r.clients)
	if clientsCount < r.minimumClients {
		return false
	}

	if clientsCount > r.maxClients {
		return false
	}

	return true
}

func (r *BaseRoom) start() {
	r.running = true
	r.roomApp.Start(r)
}

func (r *BaseRoom) GetID() string {
	return r.ID
}

func (r *BaseRoom) isRunning() bool {
	fmt.Printf("Room::isRunning() %#v \n", r.running)
	return r.running
}

func (r *BaseRoom) Broadcast(message ClientAppMessage) {
	fmt.Println("Room broadcast")
	if !r.isRunning() {
		fmt.Println("Room is not running writing to client")
		for _, userClient := range r.clients {
			userClient.GetClient().WriteJson(NewRoomWaitingToStart(r.GetID()))
		}

		return // Room is not running, should this now return an error
	}

	messageType := message.Type

	if isRoomRequest(messageType) {
		r.writeToRoom(message)
		return
	}

	if isRoomBroadcast(messageType) {
		r.broadcast(message)
		return
	}
}

func (r *BaseRoom) broadcast(message ClientAppMessage) {
	response := NewRoomResponse(message, r.GetID())

	fmt.Println("Room::broadcast()::response", response)
	for _, userClient := range r.clients {
		userClient.GetClient().WriteJson(response)
	}
}

func (r *BaseRoom) writeToRoom(message ClientAppMessage) {
	if !r.isRunning() {
		fmt.Println("Cant write to app as room is not running")
		return
	}

	r.roomApp.WriteMessage(NewRoomMessage(r, message.Payload))
}

func (r *BaseRoom) RemoveUserClient(uc UserClient) {
	userId := uc.GetUser().ID

	uc, ok := r.clients[userId]
	if !ok {
		return
	}
	delete(r.clients, userId)

	r.updated()
}

func isRoomRequest(messageType string) bool {
	return messageType == GetRequestTypes().ROOM_REQUEST
}

func isRoomBroadcast(messageType string) bool {
	return messageType == GetResponseTypes().ROOM_BROADCAST || messageType == GetResponseTypes().ROOM_BROADCAST_INIT
}
