package app_test

import (
	"errors"
	"testing"

	"github.com/kaschula/soc-draw/app"
	"github.com/kaschula/soc-draw/tests/stubs"
	. "github.com/stretchr/testify/assert"
)

type roomTest struct {
	title                  string
	client                 *stubs.ClientStub
	payload                string
	userService            app.UserService
	roomService            app.RoomService
	userClientService      app.UserClientService
	expectedAppMessageType string
	expectedPayload        string
}

func TestUserCanAccessRoomThroughLobby(t *testing.T) {
	tests := []roomTest{
		AUserRequestToJoinARoomItHasAccessTo(),
		AUserGetsAnErrorResponseWhenUserClientCanNotBeResolved(),
		AUserGetsAnErrorResponseWhenRoomIdPayloadIsInvalid(),
		AUserGetsAnErrorResponseWhenUserCanNotJoinRoom(),
		AUserGetsAnErrorResponseWhenUserTrysToJoinARoom(),
	}

	for _, test := range tests {
		t.Run(test.title, runLobbyRoomTests(test))
	}
}

func runLobbyRoomTests(test roomTest) func(t *testing.T) {
	return func(t *testing.T) {
		client := test.client

		lobby := app.NewRoomLobby(test.userService, test.roomService, test.userClientService)
		lobby.AddClient(client)
		go client.Listen()

		client.SendMessage(app.GetRequestTypes().LOBBY_ROOM_REQUEST, test.payload)

		client.WaitForReturnChan()

		if len(client.WrittenMessages) == 0 {
			t.Fatal("No Messages Written To Client")
			return
		}

		appMessage := client.WrittenMessages[0]

		Equal(t, test.expectedAppMessageType, appMessage.Type, "appMessageType does not match")
		Equal(t, test.expectedPayload, appMessage.Payload, "Payload does not match")
	}
}

func AUserRequestToJoinARoomItHasAccessTo() roomTest {
	client := &(stubs.ClientStub{
		"client:1",
		[]app.ClientResponse{},
		make(chan app.AppMessage),
		make(chan bool),
		[]app.Broadcasts{},
		nil,
	})

	user := &(app.User{"U1"})
	userService := &app.InMemoryUserService{}
	roomService := &app.DefaultRoomService{
		map[*app.User][]app.RoomI{
			user: []app.RoomI{
				app.NewDefaultRoom("r1", "Room1", nil),
			},
		},
	}
	userClientService := app.NewInMemoryUserClientService(nil)
	userClientService.CreateAndStoreUserClient(user, client)

	return roomTest{
		"A User Request To Join A Room It Has Access To Is Successful",
		client,
		`{"roomId": "r1"}`,
		userService,
		roomService,
		userClientService,
		app.GetResponseTypes().USER_JOINED_ROOM,
		`{"RoomId":"r1"}`,
	}
}

func AUserGetsAnErrorResponseWhenUserClientCanNotBeResolved() roomTest {
	client := &(stubs.ClientStub{
		"client:1",
		[]app.ClientResponse{},
		make(chan app.AppMessage),
		make(chan bool),
		[]app.Broadcasts{},
		nil,
	})

	userService := &app.InMemoryUserService{}
	roomService := &app.DefaultRoomService{
		make(map[*app.User][]app.RoomI),
	}
	userClientService := app.NewInMemoryUserClientService(nil)

	return roomTest{
		"A User Gets An Error Response When UserClient Can Not Be Resolved",
		client,
		`{"This test doesn't need a payload"}`,
		userService,
		roomService,
		userClientService,
		app.GetResponseTypes().ERROR,
		app.GetResponseErrorMessages("USER_CLIENT_404"),
	}
}

func AUserGetsAnErrorResponseWhenRoomIdPayloadIsInvalid() roomTest {
	client := &(stubs.ClientStub{
		"client:1",
		[]app.ClientResponse{},
		make(chan app.AppMessage),
		make(chan bool),
		[]app.Broadcasts{},
		nil,
	})

	user := &(app.User{"U1"})
	userService := &app.InMemoryUserService{}
	roomService := &app.DefaultRoomService{
		make(map[*app.User][]app.RoomI),
	}

	userClientService := app.NewInMemoryUserClientService(nil)
	userClientService.CreateAndStoreUserClient(user, client)

	return roomTest{
		"A User Gets An Error Response When RoomId Payload Is Invalid",
		client,
		`{"Invalid payload": "r1"}`,
		userService,
		roomService,
		userClientService,
		app.GetResponseTypes().ERROR,
		app.GetResponseErrorMessages("PAYLOAD_ROOM_ID"),
	}
}

func AUserGetsAnErrorResponseWhenUserCanNotJoinRoom() roomTest {
	client := &(stubs.ClientStub{
		"client:1",
		[]app.ClientResponse{},
		make(chan app.AppMessage),
		make(chan bool),
		[]app.Broadcasts{},
		nil,
	})

	user := &(app.User{"U1"})
	userService := &app.InMemoryUserService{}
	roomService := &app.DefaultRoomService{
		map[*app.User][]app.RoomI{
			user: []app.RoomI{
				app.NewDefaultRoom("r1", "Room1", nil),
			},
		},
	}
	userClientService := app.NewInMemoryUserClientService(nil)
	userClientService.CreateAndStoreUserClient(user, client)

	return roomTest{
		"A User Gets An Error Response When User Can Not Join Room",
		client,
		`{"roomId": "r2"}`,
		userService,
		roomService,
		userClientService,
		app.GetResponseTypes().ERROR,
		app.GetResponseErrorMessages("USER_ROOM_AUTH"),
	}
}

func AUserGetsAnErrorResponseWhenUserTrysToJoinARoom() roomTest {
	client := &(stubs.ClientStub{
		"client:1",
		[]app.ClientResponse{},
		make(chan app.AppMessage),
		make(chan bool),
		[]app.Broadcasts{},
		nil,
	})

	user := &(app.User{"U1"})
	userService := &app.InMemoryUserService{}
	inMemRoomRepo := &app.DefaultRoomService{
		map[*app.User][]app.RoomI{
			user: []app.RoomI{
				app.NewDefaultRoom("r1", "Room1", nil),
			},
		},
	}
	roomService := roomServiceStub{inMemRoomRepo}

	userClientService := app.NewInMemoryUserClientService(nil)
	userClientService.CreateAndStoreUserClient(user, client)

	return roomTest{
		"A User Gets An Error Response When User Trys To Join A Room",
		client,
		`{"roomId": "r1"}`,
		userService,
		&roomService,
		userClientService,
		app.GetResponseTypes().ERROR,
		app.GetResponseErrorMessages("ADD_USER_TO_ROOM"),
	}
}

type roomServiceStub struct {
	*app.DefaultRoomService
}

// Overide the app.DefaultRoomService AddUserClient to return error
func (r *roomServiceStub) AddUserClient(userClient app.UserClient, roomId string) error {
	return errors.New("Add User Error")
}
