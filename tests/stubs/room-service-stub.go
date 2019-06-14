package stubs

import (
	"errors"

	"github.com/kaschula/socket-server/app"
)

type RoomServiceStub struct {
	userRooms map[*app.User][]app.RoomI
}

func NewRoomServiceStub(userRooms map[*app.User][]app.RoomI) app.RoomService {
	return &RoomServiceStub{userRooms}
}

func (s *RoomServiceStub) GetRooms(user *app.User) ([]app.RoomI, error) {
	userRooms, ok := s.userRooms[user]
	if !ok {
		return nil, errors.New("Can not resolve rooms")
	}

	return userRooms, nil
}
func (s *RoomServiceStub) GetRoom(user *app.User, roomId string) (app.RoomI, error) {
	userRooms, err := s.GetRooms(user)
	if err != nil {
		return nil, errors.New("Can not resolve rooms")
	}

	for _, room := range userRooms {
		if room.GetID() == roomId {
			return room, nil
		}
	}

	return nil, errors.New("Can not resolve room")
}

func (s *RoomServiceStub) CanUserJoin(userClient app.UserClient, roomId string) bool {
	return true
}

func (s *RoomServiceStub) AddUserClient(uc app.UserClient, roomId string) error {
	room, err := s.GetRoom(uc.GetUser(), roomId)

	if err != nil {
		return errors.New("Could not resolve Room")
	}

	room.AddUserClient(uc)

	return nil
}

// type RoomService interface {
// 	GetRooms(user *User) ([]RoomI, error)
// 	GetRoom(user *User, roomId string) (RoomI, error)
// 	CanUserJoin(userClient UserClient, roomId string) bool
// 	AddUserClient(userClient UserClient, roomId string) error
// }
