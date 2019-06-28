# Soc-Draw framework

This project is the early stages of a framework built with the intent of allowing easy development of multi-user real time socket applications. The `soc-draw` application intends to provide a framework that manages user's through lobbies and rooms to begin a socket based web application. The socket application will be built into 2 parts; a Frontend (client) and the Backend application.

### The General Idea, what soc-draw actually does

This framework allows users to login into a lobby for a `room-application`. A user logs into the soc-draw application. This enters them into the lobby. Once in the lobby a user will be able to see what rooms they have access to. From here they can enter a room. Once they have entered a room, the room is in a state of being ready to start. 

The rooms in the simple-chat-application included as an example in this repository require at least two users to be be in the room before the room-application is initialised and started. See the example below for information. 

The user can interact with the room application and the soc-draw frontend client will continue to update the application with state received from the server and send messages (usually some kind of user event) to the server application.

### The `Room-Application`

This is defined by a user of the framework. It involves two parts, the frontend room application and the backend room application.

The backend part of the room application is responsible for:
1. Handling messages from the front client. These messages will typically be user events that will result in a state change. 
2. Sending state to the room to be broadcast back to every room client.

The frontend should be responsible for handling and sending user events to the server and processing received application state.

The soc-draw framework provides an application `<div>`  that a room application can attach to and build what ever mark up is required for a room-application. This would be done in the room-applications initialise method.

Both the `room-application` FE and BE need to provide an update function. 


## Motivation

This project was undertaken to help develop my knowledge and understanding of the Go-Lang language and web sockets and the ideas covered by this challenging area of real-time applications.

My focus in this project has been on the Go-lang code. The javascript code I have tried to keep as minimal and clean as possible without using any frameworks or tooling. 

## Build status

The master version of this project builds using the GO compiler. It contains only in-memory versions of the services used for persisted enities such as users and rooms. This means all the room and user data is set up in the main.go as map data. These services are in a primitive state and will evolve in the future. For example Authentication at the moment is done by entering in a username that exists in the service. In the future a proper password authentication system will be used.

## Tech/framework used

This project relies on the Gorilla Websocket (https://github.com/gorilla/websocket)

The go code test uses Testify (https://github.com/stretchr/testify)

The sample chat application front end use jQuery (https://jquery.com/)


## Tests

The `soc-draw` socket server is tested using Go-langs testing package. The tests range from unit tests for individual structs, integration tests for key components like the Lobby and feature tests which attempt to cover a journeys on the backend.

To run the test use the following command from the project root:

`go test ./...`

# Examples

## Simple Chat Application

The project currently ships with an example chat application. Running this application is simple.

Simply build and run the server by running either of the following from the project root:

`go run main.go simple-chat-application.go`

or 

`go build` and `./soc-draw`


In a browser go to the `localhost:8089` and you should see a small form to enter a user name and log in

![Login screen](/readme-resources/sd_user_login.png)

To see the user names that will allow you to enter the lobby check the main.go file and look for the user structs used to create the the `users` map.

![User Data](/readme-resources/user_data.png)

Enter a valid user ID into the form input labelled `username`.

![User Logging In](/readme-resources/user-logging-in.png)

Once a you have entered the lobby you will have access to some rooms. At this point a socket connection has been created for that user. Click on a room to enter it.

![User In lobby](/readme-resources/user-in-lobby.png)

Once a user has entered a room you might have to wait until enough users have entered the room before the room-application starts. This simple chat application that this lobby works with requires at least two users in a room to begin.

![User Entered room and waiting](/readme-resources/user-entered-room-and-waiting.png)

In separate browser window. Repeat the steps with a different user. Once the second user enters the same room as the first the room-application should start and be shown to both users

![Second user logging in](/readme-resources/second-user-logged-in.png)

At this point the soc-draw client has done its work and handled users through the lobby into rooms and initiated the room application. The state received from the initial room set up event and subsequent room broadcasts is passed to the room-application update function.

![Two-users-in-room](/readme-resources/two-users-entered.png)

The chat application itself is responsible for building the message display window, handling the the text form and attaching any application events such as, pressing the enter button.

![Two-users-in-room](/readme-resources/chat-app-running.png)

#### Multiple rooms

Although a room application is running users are able to move to other rooms.

In a third browser window login as a third user and enter the second room, you will see that user 3 is waiting for a second user to enter the room.

Using user:1 window enter the second room. This will start a second conversation notice that messages sent into room 2 by user 1 or 3 do not appear in room 1 and vice-versa. User 1 can seamlessly move between rooms without missing messages by clicking on the different rooms. 

The `soc-draw` application is able to store state by room id so that a user can receive state for multiple rooms while only displaying the state for one.

## Example, building or own room application frontend

... in progress

## Domain

The domain is generally quite simple. As it stands one soc-draw application instance runs one lobby. A lobby is responsible for managing multiple users and their access to multiple rooms. A lobby is responsible for one room application and each room manages only state for that room's application instance.

Currently the creation of user and rooms is hard coded in the main.go file as all the applications only persists state in memory at this point.

... more to come

## Installation

The soc-draw application is just a prototype at this stage and has been designed for local development only.

To start the application and run quickly use go run from the root of the project

`go run main.go simple-chat-application.go`

This will start the project and allows access via `http://localhost:8089/`

To build and run the server binary in root run:

`go build`
`./soc-draw`

To install the project run:

`go install`

The `soc-draw` binary should be in your go workspace `/bin` directory

## API Reference

There are two routes defined by the server. The home route the websocket url.

Home - http://localhost:8089

Socket creation - http://localhost:8089/ws

## Other things

The go code has been organised into just one package currently. It will be extracted into modules as the project evolves.

## Next steps

The next items to be be added and fixed in no particular order are:

- Create flow diagram to describe the Applications's process in more detail
- Framework for developing room applications
- Add room notifications, these are messages sent by the server to give updates on rooms.
- Handle closing sockets and deactivating rooms
- Update FE lobby to provide better UX
- Introduce proper authentication for users
- Create CRUD actions for users, rooms
- Improve error system and error handling on the frontend
- Introduce private rooms
