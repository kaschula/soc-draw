package app

import (
	"errors"
)

type RoomService interface {
	GetRooms(user *User) ([]RoomI, error)
	GetRoom(user *User, roomId string) (RoomI, error)
	CanUserJoin(userClient UserClient, roomId string) bool
	AddUserClient(userClient UserClient, roomId string) error
}

func NewDefaultRoomService(userRooms map[*User][]RoomI) RoomService {
	return &InMemoryRoomRepository{userRooms}
}

type InMemoryRoomRepository struct {
	UserRooms map[*User][]RoomI
}

func (r *InMemoryRoomRepository) GetRooms(user *User) ([]RoomI, error) {
	rooms, ok := r.UserRooms[user]
	if !ok || rooms == nil {
		return nil, errors.New("Unable to resolve User's Rooms")
	}

	return rooms, nil
}

func (r *InMemoryRoomRepository) GetRoom(user *User, roomId string) (RoomI, error) {
	var room RoomI
	for _, r := range r.UserRooms[user] {
		if r.GetID() == roomId {
			room = r
		}
	}

	// How do you check if room is empty or not?

	return room, nil
}

func (r *InMemoryRoomRepository) CanUserJoin(userClient UserClient, roomId string) bool {
	user := userClient.GetUser()

	rooms, err := r.GetRooms(user)
	if err != nil {
		return false
	}

	for _, room := range rooms {
		if room.GetID() == roomId {
			return true
		}
	}

	return false
}

func (rs *InMemoryRoomRepository) AddUserClient(client UserClient, roomId string) error {
	room, err := rs.GetRoom(client.GetUser(), roomId)

	// fmt.Println(room)

	if err != nil {
		return errors.New("Could not resolve Room")
	}

	room.AddUserClient(client)

	return nil
}
