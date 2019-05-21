package main

func main() {

	// Set up home directory, simple file sever
	// set up websocket connection
	// set up routes to create rooms

	// Write app abstraction that allows for the a user to connect and be returned lobby status
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
