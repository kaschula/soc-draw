package app

import (
	"fmt"
)

// Need to think about the abstraction here
// It might be worth having a base room app that handles all the set
// This way someone writing an application just needs to write the ProcessApp method and Initial State method
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
	rooms     []RoomI // This probably wont be needed can be deleted
	roomState map[string]string
}

func NewSimpleRoomApplication(logger Prints) RoomApplication {
	return &SimpleRoomApplication{make(chan RoomMessage), logger, []RoomI{}, make(map[string]string)}
}

func (app *SimpleRoomApplication) Start(room RoomI) {
	initAppState := "Initial State"
	app.roomState[room.GetID()] = initAppState
	app.rooms = append(app.rooms, room)
	room.Broadcast(NewClientAppMessage(NewNoClient(), NewAppMessage("ROOM_BROADCAST_INIT", initAppState)))
}

func (app *SimpleRoomApplication) WriteMessage(message RoomMessage) {
	app.writeChan <- message
}

// Old run for Simple app
func (app *SimpleRoomApplication) Run() {
	fmt.Println("RoomApplication::Run() app listen")
	for {
		message := <-app.writeChan
		fmt.Printf("RoomApplication::Run()::message %#v \n", message)

		app.logger.Printf("Room %#v is Broadcasting %#v", message.GetRoom().GetID(), message)
		// Should maybe send message back to Room and Room should send to its broadcasters
		message.GetRoom().Broadcast(
			NewClientAppMessage(NewNoClient(), NewAppMessage(GetResponseTypes().ROOM_BROADCAST, message.GetPayload())),
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
