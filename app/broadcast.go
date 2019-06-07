package app

type Broadcasts interface {
	// Broadcast(client IsClient, message ClientAppMessage)
	Broadcast(message ClientAppMessage)
	GetID() string
	RemoveUserClient(UserClient)
}
