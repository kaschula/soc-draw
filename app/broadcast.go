package app

type Broadcasts interface {
	Broadcast(message ClientAppMessage)
	GetID() string
	RemoveUserClient(UserClient)
}
