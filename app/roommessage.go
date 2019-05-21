package app


type RoomMessage struct {
	room    RoomI
	payload string
}

func (m *RoomMessage) GetPayload() string {
	return m.payload
}

func (m *RoomMessage) GetRoom() RoomI {
	return m.room
}

func NewRoomMessage(room RoomI, payload string) RoomMessage {
	return RoomMessage{room, payload}
}
