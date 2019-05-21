package app

import "fmt"

type RoomApplication interface {
	// Give a room this method when a room is ready to begin
	// This method should set up inital data required for an application and broadcast to clients
	Start(room RoomI)

	// Write an event to the application from room
	WriteMessage(roomMessage RoomMessage)
	// Start the game in go routine to listen for actions
	Run()
}

// This Is the simplest application,
// It recives a RoomMessage, logs the payload and
// then broadcasts that message back to the room
type SimpleRoomApplication struct {
	writeChan chan RoomMessage
	logger    Prints
	rooms     []RoomI
}

func NewSimpleRoomApplication(logger Prints) RoomApplication {
	return &SimpleRoomApplication{make(chan RoomMessage), logger, []RoomI{}}
}

func (app *SimpleRoomApplication) Start(room RoomI) {
	room.Broadcast(nil, NewClientAppMessage("", NewAppMessage(MESSAGE_TYPE_ROOM_WELCOME, `{"message":"welcome"}`)))
	app.rooms = append(app.rooms, room)
}

func (app *SimpleRoomApplication) WriteMessage(message RoomMessage) {
	app.writeChan <- message
}

func (app *SimpleRoomApplication) Run() {
	for {
		message := <-app.writeChan

		app.logger.Printf("Room %#v is Broadcasting %#v", message.GetRoom().GetID(), message)
		// Should maybe send message back to Room and Room should send to its broadcasters
		message.GetRoom().Broadcast(
			nil,
			NewClientAppMessage("", NewAppMessage(MESSAGE_TYPE_ROOM, message.GetPayload())),
		)
	}
}

type Prints interface {
	Printf(format string, v ...interface{})
}

type SimpleLogger struct{}

func NewSimpleLogger() *SimpleLogger {
	return &SimpleLogger{}
}

func (l *SimpleLogger) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}
