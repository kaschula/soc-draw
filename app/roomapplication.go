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
	room.Broadcast(NewClientAppMessage(nil, NewAppMessage("ROOM_BROADCAST_INIT", initAppState)))
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
			NewClientAppMessage(nil, NewAppMessage(ClientResponseTypes().ROOM_BROADCAST, message.GetPayload())),
		)
	}
}

// Messaging App uses new Run. For messaging
// func (app *SimpleRoomApplication) Run() {
// 	fmt.Println("RoomApplication::Run() app listen")
// 	for {
// 		message := <-app.writeChan
// 		fmt.Printf("RoomApplication::Run()::message %#v \n", message)

// 		state, err := app.processRequest(message.GetPayload())

// 		if err != nil {
// 			message.GetRoom().Broadcast(
// 				nil,
// 				NewClientAppMessage("", NewAppMessage("ROOM_BROADCAST_ERROR", err.Error())),
// 			)
// 		}

// 		fmt.Printf("RoomApplication::Run():: processed state %#v \n", state)
// 		// Should maybe send message back to Room and Room should send to its broadcasters
// 		message.GetRoom().Broadcast(
// 			nil,
// 			NewClientAppMessage("", NewAppMessage(
// 				ClientResponseTypes().ROOM_BROADCAST,
// 				// fmt.Sprintf(`{"state": %#v}`, state),
// 				state,
// 			)),
// 		)
// 	}
// }

// AppSpecific
// func (app *SimpleRoomApplication) processRequest(payload string) (string, error) {
// 	fmt.Println("RoomApplication::processRequest() payload", payload)

// 	request, err := parseRequest(payload)

// 	if err != nil {
// 		return "", err
// 	}

// 	fmt.Println("RoomApplication::processRequest() request", request)

// 	switch request.RequestType {
// 	case "STATE":
// 		return app.currentState(request)
// 	case "ROOM_EVENT":
// 		return app.update(request, payload)
// 	default:
// 		fmt.Println("RoomApplication::processRequest() request Type")
// 		return "", nil
// 	}
// }

// func (app *SimpleRoomApplication) currentState(r MessageAppProcessRequest) (string, error) {
// 	roomState, ok := app.roomState[r.RoomID]
// 	if !ok {
// 		return "", errors.New("Room state could not be found")
// 	}

// 	return roomState, nil
// }

// func (app *SimpleRoomApplication) update(request MessageAppProcessRequest, originalPayload string) (string, error) {
// 	fmt.Printf("RoomApplication::update() originalPayload: %#v , request: %#v \n", originalPayload, request)

// 	userMessage, err := parseRoomEventUserMessage(originalPayload)

// 	if err != nil {
// 		return "", err
// 	}

// 	fmt.Printf("RoomApplication::update() userMessage: %#v ", userMessage)

// 	updatedRoomStateJson, err := app.appendNewMessageToCurrentState(request.RoomID, userMessage)
// 	if err != nil {
// 		return "", err
// 	}

// 	app.roomState[request.RoomID] = updatedRoomStateJson

// 	return updatedRoomStateJson, nil
// }

// func (app *SimpleRoomApplication) appendNewMessageToCurrentState(roomId string, newMessage RoomEventUserMessage) (string, error) {
// 	fmt.Println("RoomApplication::appendNewMessageToCurrentState()")

// 	storedRoomState, ok := app.roomState[roomId]
// 	if !ok {
// 		return "", errors.New("Room state could not be found")
// 	}
// 	fmt.Println("RoomApplication::appendNewMessageToCurrentState():: roomState %#v \n", storedRoomState)

// 	roomState, err := parseStoredRoomState(storedRoomState)
// 	if err != nil {
// 		return "", errors.New("Room state could not be found")
// 	}

// 	fmt.Printf("RoomApplication::appendNewMessageToCurrentState() %#v \n", roomState)

// 	roomState.State.Messages = append(roomState.State.Messages, ApplicationStateMessage{newMessage.Sender, newMessage.Message})

// 	fmt.Printf("RoomApplication::appendNewMessageToCurrentState() appended message %#v \n", roomState.State)

// 	// marshall into json
// 	updatedStateJson, err := marshallState(roomState)
// 	if err != nil {
// 		return "", err
// 	}

// 	fmt.Printf("Updated Json marshalled: %#v \n ", string(updatedStateJson))

// 	return string(updatedStateJson), nil
// }

// Messaging App Stuff

// type MessageAppProcessRequest struct {
// 	RoomID      string `json:"roomId"`
// 	RequestType string `json:"requestType"`
// }

// type RoomEventUserMessage struct {
// 	Sender  string `json:"username"`
// 	Message string `json:"message"`
// }

// type ApplicationStoredState struct {
// 	State ApplicationState `json:"state"`
// }

// type ApplicationState struct {
// 	Messages []ApplicationStateMessage `json:"messages"`
// }

// type ApplicationStateMessage struct {
// 	Sender  string `json:"sender"`
// 	Message string `json:"message"`
// }

// makePrivateMethod
// func marshallState(state ApplicationStoredState) ([]byte, error) {
// 	return json.Marshal(state)
// }

// func parseStoredRoomState(storedState string) (ApplicationStoredState, error) {
// 	fmt.Println("RoomApplication::parseStoredRoomState() storedState", storedState)
// 	var roomState ApplicationStoredState

// 	raw := []byte(storedState)
// 	// raw := []byte("")
// 	json.Unmarshal(
// 		raw,
// 		&roomState,
// 	)

// 	fmt.Printf("RoomApplication::parseStoredRoomState() unmarshalled %#v \n", roomState)

// 	// Any Validation?

// 	return roomState, nil
// }

// func parseRequest(payload string) (MessageAppProcessRequest, error) {
// 	var request MessageAppProcessRequest

// 	raw := []byte(payload)
// 	json.Unmarshal(
// 		raw,
// 		&request,
// 	)

// 	if request.RoomID == "" {
// 		return request, errors.New("Can not resolve Room ID")
// 	}

// 	if request.RequestType == "" {
// 		return request, errors.New("Can not resolve Request type")
// 	}

// 	return request, nil
// }

// func parseRoomEventUserMessage(payload string) (RoomEventUserMessage, error) {
// 	var request RoomEventUserMessage

// 	raw := []byte(payload)
// 	json.Unmarshal(
// 		raw,
// 		&request,
// 	)

// 	if request.Sender == "" {
// 		return request, errors.New("Can not resolve Sender")
// 	}

// 	if request.Message == "" {
// 		return request, errors.New("Can not resolve Message")
// 	}

// 	return request, nil
// }

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
