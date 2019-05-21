package app

type RoomI interface {
	GetID() string
	AddUserClient(UserClient)
	WriteMessage(clientMessage ClientAppMessage)
	Broadcasts
}

type BaseRoom struct {
	ID             string `json:"ID"`
	Name           string `json:"Name"`
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
	uc.GetClient().Subscribe(r) // <-------- Not Unit Untested

	r.updated()
}

func (r *BaseRoom) updated() {
	if r.isRunning() || !r.isReady() {
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
	return r.running
}

func (r *BaseRoom) Broadcast(client IsClient, message ClientAppMessage) {
	if !r.isRunning() {
		return // Room is not running, should this now return an error
	}

	if !isRoomType(message.Type) {
		return
	}

	response := NewRoomResponse(message.Payload)

	for _, userClient := range r.clients {
		userClient.GetClient().WriteJson(response)
	}
}

func (r *BaseRoom) WriteMessage(message ClientAppMessage) {
	if !r.isRunning() {
		return
	}

	r.roomApp.WriteMessage(NewRoomMessage(r, message.Payload))
}

// Untested
func isRoomType(messageType string) bool {
	return messageType == ClientResponseTypes().ROOM_BROADCAST ||
		messageType == MESSAGE_TYPE_ROOM ||
		messageType == MESSAGE_TYPE_ROOM_WELCOME
}
