package main

import (
	"fmt"
	"net/http"

	socketServer "github.com/kaschula/socket-server/app"
)

func main() {

	appClientPath := "public/simple-chat"
	appServer := socketServer.NewSimpleRoomApplication(socketServer.NewSimpleLogger())
	// go appServer.Run()

	// Data

	userOne, userTwo := &socketServer.User{"U:1"}, &socketServer.User{"U:2"}

	users := map[string]*socketServer.User{
		userOne.ID: userOne,
		userTwo.ID: userTwo,
	}

	globalRooms := []socketServer.RoomI{
		socketServer.NewDefaultRoom("R:1", "First Room", appServer),
		socketServer.NewDefaultRoom("R:2", "Second Room", appServer),
	}

	roomRepository := map[*socketServer.User][]socketServer.RoomI{
		userOne: globalRooms,
	}

	app := newApp(users, roomRepository)
	fs := http.FileServer(http.Dir(appClientPath))
	http.Handle("/", fs)

	http.HandleFunc("/ws", app.CreateSocketHandler)

	fmt.Println("http server started on :8089")
	err := http.ListenAndServe(":8089", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}

func newApp(
	users map[string]*socketServer.User,
	roomRepository map[*socketServer.User][]socketServer.RoomI,
) *socketServer.App {
	roomService := socketServer.NewDefaultRoomService(roomRepository)
	userRepository := socketServer.NewInMemoryUserRepository(users)
	userClientService := socketServer.NewInMemoryUserClientService()

	lobby := socketServer.NewLobby(userRepository, roomService, userClientService)

	factory := &socketServer.GorillaWebsocketUpgradeFactory{}
	clientService := socketServer.NewDefaultClientService(factory)

	return socketServer.NewApp(clientService, lobby)
}

// Next steps
// Need to think about filter messages in the Client
// message is Sent from the Client, The client needs to send this onto its subscriber
// If the client is part of lots of rooms then it should filter its broadcasters for that partiocular Broadcaster
// This means the lobby and rooms will have to have some kind of broadcast Ideas

// Things that need to be look at
// The message types bit is really messy and not clear

// Need to start ClientService and factory
// Actual Client implementation will need to do the filtering discussed

// Then ready for web test

// -------- Main App tasks --------

// Key DOMAIN point. A lobby belongs to a specifc application, this application is what is started when inside a room. Lobby depends on RoomApplication

// The web application will allow users to join application (games) through a lobby.

// Once journey Is working, hook up to front end

// Create Crud apis for handling Rooms, Users, and user relationships
// Will need an Authenticate App to log Users in
// Users can be linked in a friendship relation ship.

// Refactor
// Create  A RoomService that depends on RoomRepository
// Lobby must user RoomService
// Mover Lobby Tests into an Integration test directory
// Make lobby_test for lobby unit test for testing Errors

// Continue resolving Rooms for a User
// For this prototype Rooms a Room will be created for every other User on the app. Every user can access every other for now

// Need to Make Crud Operation for Making rooms and adding contacts
// When making crud, a room will have to have an App assigned to it
// this will mean roomApp service will be required Or is the lobby assigned a room app and that creates the rooms
// Room service will have to be given an room app to create it?
