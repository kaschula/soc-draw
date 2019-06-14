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
	return &DefaultRoomService{userRooms}
}

type DefaultRoomService struct {
	UserRooms map[*User][]RoomI
}

func (s *DefaultRoomService) GetRooms(user *User) ([]RoomI, error) {
	rooms, ok := s.UserRooms[user]
	if !ok || rooms == nil {
		return nil, errors.New("Unable to resolve User's Rooms")
	}

	return rooms, nil
}

func (s *DefaultRoomService) GetRoom(user *User, roomId string) (RoomI, error) {
	var room RoomI
	for _, r := range s.UserRooms[user] {
		if r.GetID() == roomId {
			room = r
		}
	}

	// How do you check if room is empty or not?

	return room, nil
}

func (s *DefaultRoomService) CanUserJoin(userClient UserClient, roomId string) bool {
	user := userClient.GetUser()

	rooms, err := s.GetRooms(user)
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

func (s *DefaultRoomService) AddUserClient(uc UserClient, roomId string) error {
	room, err := s.GetRoom(uc.GetUser(), roomId)

	if err != nil {
		return errors.New("Could not resolve Room")
	}

	return room.AddUserClient(uc)
}
