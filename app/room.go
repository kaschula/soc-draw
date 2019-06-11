package app

import "fmt"

type RoomI interface {
	// GetID() string
	AddUserClient(UserClient)
	WriteMessage(clientMessage ClientAppMessage)
	Broadcasts
}

type BaseRoom struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`
	// Todo: Change this to Map
	clients        []UserClient
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
		[]UserClient{},
		false,
		minimumClients,
		maxClients,
		roomApp,
	}
}

func (r *BaseRoom) AddUserClient(uc UserClient) {
	r.clients = append(r.clients, uc)

	// subscribe room to client
	client := uc.GetClient() // <-------- Not Unit Untested
	client.Subscribe(r)
	fmt.Println("Room::AddUserClient() UC added")
	r.updated()
}

// retest
func (r *BaseRoom) updated() {
	fmt.Println("Room::updated()")
	// if r.isRunning() || !r.isReady() {
	// 	fmt.Println("Room::updated() Room is running or not ready: (running) (isReady)", r.isRunning(), !r.isReady())
	// 	return
	// }

	if !r.isReady() {
		fmt.Println("Room::updated() is not ready", !r.isReady())
		if r.isRunning() {
			fmt.Println("Room::updated() Room is running and not ready: room should stop")
			// r.stop()
		}

		return
	}

	fmt.Println("Room::updated() room is ready starting room")
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

/**
I think it might be required to create a RoomResponse,
This should extend ClientResponse but contain RoomId and Room.Running
The payload should be untouched and what ever is given
This will help when the room app is returning data
*/

func (r *BaseRoom) Broadcast(message ClientAppMessage) {
	if !r.isRunning() {
		fmt.Println("Room is not running writing to client")
		for _, userClient := range r.clients {
			userClient.GetClient().WriteJson(NewRoomWaitingToStart(r.GetID()))
		}
		return // Room is not running, should this now return an error
	}

	messageType := message.Type

	// 1* remove this first if
	if !isRoomType(messageType) {
		fmt.Println("Room::Broadcast(), message type note rooms: ", messageType)
		return
	}

	if isRoomRequest(messageType) {
		r.WriteMessage(message)
		return
	}

	// 1* add this, if request goes to app, if broadcast goes to clients, ignore everything else
	// if isRoomBroadcast(messageType) {
	// 	r.broadcast(message)
	// 	return
	// }

	// *1 remove this when adding extra if
	r.broadcast(message)
}

func (r *BaseRoom) broadcast(message ClientAppMessage) {
	// move NewRoomResponse into Room
	response := NewRoomResponse(message, r.GetID())

	fmt.Println("Room::broadcast()::response", response)
	for _, userClient := range r.clients {
		userClient.GetClient().WriteJson(response)
	}
}

// Todo: should this be public????
func (r *BaseRoom) WriteMessage(message ClientAppMessage) {
	fmt.Println("Room:: WriteMessage()::message", message)
	if !r.isRunning() {
		fmt.Println("Cant write to app as room is not running")
		return
	}
	fmt.Println("Room:: WriteMessage()::writing")
	r.roomApp.WriteMessage(NewRoomMessage(r, message.Payload))
}

// Untested
func isRoomType(t string) bool {
	responses := GetResponseTypes()

	return t == responses.ROOM_BROADCAST ||
		t == GetRequestTypes().ROOM_REQUEST || // add this Type to requests Types
		t == responses.ROOM_BROADCAST_INIT ||
		t == responses.ROOM_BROADCAST_MESSAGE
}

//Untested
func (r *BaseRoom) RemoveUserClient(uc UserClient) {
	// This is a terrible way of doing this, need to change to slice to map

	swop := make([]UserClient, 0)
	for _, client := range r.clients {
		if uc.GetClient().GetID() == client.GetClient().GetID() {
			continue
		}

		swop = append(swop, client)
	}

	r.clients = swop

	r.updated()
}

func isRoomRequest(messageType string) bool {
	return messageType == "ROOM_REQUEST"
}

func isRoomBroadcast(messageType string) bool {
	return messageType == "ROOM_BROADCAST"
}
