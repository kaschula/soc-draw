package app_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/kaschula/socket-server/app"
	. "github.com/stretchr/testify/assert"
)

type testData struct {
	title                         string
	client                        *ClientStub2
	userRepository                app.UserRepository
	roomService                   app.RoomService
	userClientService             UserClientServiceStub
	messageToSend                 string
	expectedAppMessageType        string
	expectedPayload               string
	expectedcreateUserClientCalls int
}

func TestLobbyResolveUserRequest(t *testing.T) {
	tests := []testData{
		ALobbyCanGetAUserRequestAndSendUserResponseWithRoomDataSetUp(),
		AnErrorResponseIsSentToTheClientWhenUserCantBeResolvedSetUp(),
		AnErrorResponseIsSentToTheClientWhenLobbyDataCantBeResolvedSetUp(),
		AnErrorResponseIsSentToTheClientWhenUserClientCantBeCreatedSetUp(),
	}

	for _, test := range tests {
		t.Run(test.title, runUserRequestTest(test))
	}
}

func runUserRequestTest(test testData) func(t *testing.T) {
	return func(t *testing.T) {
		client := test.client

		lobby := app.NewLobby(test.userRepository, test.roomService, &test.userClientService)
		lobby.AddClient(client)

		go client.Listen()

		client.SendMessage(app.MessageTypeLobbyUserJoinRequest, test.messageToSend)
		// Wait for something to be written
		<-client.returnChan

		if len(client.writtenMessages) == 0 {
			t.Fatal("No Messages Written To Client")
			return
		}

		appMessage := client.writtenMessages[0]

		// The expected and actual are the wrong way round
		Equal(t, appMessage.Type, test.expectedAppMessageType, "appMessageType does not match")
		Equal(t, appMessage.Payload, test.expectedPayload, "Payload does not match")
		Equal(t, test.userClientService.createUserClientCalls, test.expectedcreateUserClientCalls, "User Client Repository Create Method ")
	}
}

func ALobbyCanGetAUserRequestAndSendUserResponseWithRoomDataSetUp() testData {
	client := &(ClientStub2{
		"client:1",
		[]app.ClientResponse{},
		make(chan app.AppMessage),
		make(chan string),
		[]app.Broadcasts{},
	})
	user := app.NewUser("UserID124")
	userRepository := &app.InMemoryUserRepository{
		map[string]*app.User{
			"UserID124": user,
		},
	}

	roomOne := app.NewDefaultRoom("r1", "Room1", nil)
	roomTwo := app.NewDefaultRoom("r2", "Room2", nil)

	roomService := &app.DefaultRoomService{
		map[*app.User][]app.RoomI{
			user: []app.RoomI{roomOne, roomTwo},
		},
	}

	userClientService := UserClientServiceStub{0, nil}

	return testData{
		title:                  "Test A Lobby Can Recieve A User Request And Send Lobby Data Response",
		client:                 client,
		userRepository:         userRepository,
		roomService:            roomService,
		userClientService:      userClientService,
		messageToSend:          "{\"user\": \"UserID124\"}",
		expectedAppMessageType: "USER_LOBBY_DATA",
		// This expectPayload may become a problem as Room struct grows
		expectedPayload:               `{"User":{"ID":"UserID124"},"Rooms":[{"ID":"r1","Name":"Room1"},{"ID":"r2","Name":"Room2"}]}`,
		expectedcreateUserClientCalls: 1,
	}
}

func AnErrorResponseIsSentToTheClientWhenUserCantBeResolvedSetUp() testData {
	client := &(ClientStub2{
		"client:1",
		[]app.ClientResponse{},
		make(chan app.AppMessage),
		make(chan string),
		[]app.Broadcasts{},
	})

	userRepository := &app.InMemoryUserRepository{
		map[string]*app.User{},
	}

	roomService := &app.DefaultRoomService{}

	userClientService := UserClientServiceStub{0, nil}

	return testData{
		title:                         "Test An Error Response Is Sent To The Client When A User Cant Be Resolved",
		client:                        client,
		userRepository:                userRepository,
		roomService:                   roomService,
		userClientService:             userClientService,
		messageToSend:                 "{\"user\": \"NoExistantID\"}",
		expectedAppMessageType:        app.ClientResponseTypes().ERROR,
		expectedPayload:               app.GetResponseErrorMessages("USER"),
		expectedcreateUserClientCalls: 0,
	}
}

func AnErrorResponseIsSentToTheClientWhenLobbyDataCantBeResolvedSetUp() testData {
	client := &(ClientStub2{
		"client:1",
		[]app.ClientResponse{},
		make(chan app.AppMessage),
		make(chan string),
		[]app.Broadcasts{},
	})

	user := app.NewUser("UserID124")
	userRepository := &app.InMemoryUserRepository{
		map[string]*app.User{
			"UserID124": user,
		},
	}

	roomService := &app.DefaultRoomService{}

	userClientService := UserClientServiceStub{0, nil}

	return testData{
		title:                         "An Error Response Is Sent To The Client When Lobby Data Cant Be Resolved",
		client:                        client,
		userRepository:                userRepository,
		roomService:                   roomService,
		userClientService:             userClientService,
		messageToSend:                 "{\"user\": \"UserID124\"}",
		expectedAppMessageType:        app.ClientResponseTypes().ERROR,
		expectedPayload:               app.GetResponseErrorMessages("LOBBY_DATA"),
		expectedcreateUserClientCalls: 0,
	}
}

func AnErrorResponseIsSentToTheClientWhenUserClientCantBeCreatedSetUp() testData {
	client := &(ClientStub2{
		"client:1",
		[]app.ClientResponse{},
		make(chan app.AppMessage),
		make(chan string),
		[]app.Broadcasts{},
	})

	user := app.NewUser("UserID124")
	userRepository := &app.InMemoryUserRepository{
		map[string]*app.User{
			"UserID124": user,
		},
	}

	roomOne := app.NewDefaultRoom("r1", "Room1", nil)
	roomTwo := app.NewDefaultRoom("r2", "Room2", nil)

	roomService := &app.DefaultRoomService{
		map[*app.User][]app.RoomI{
			user: []app.RoomI{roomOne, roomTwo},
		},
	}

	userClientService := UserClientServiceStub{0, errors.New("UserClient Error")}

	return testData{
		title:                         "An Error Response Is Sent To The Client When Lobby Data Cant Be Resolved",
		client:                        client,
		userRepository:                userRepository,
		roomService:                   roomService,
		userClientService:             userClientService,
		messageToSend:                 "{\"user\": \"UserID124\"}",
		expectedAppMessageType:        app.ClientResponseTypes().ERROR,
		expectedPayload:               app.GetResponseErrorMessages("USER_CLIENT"),
		expectedcreateUserClientCalls: 1,
	}
}

// This can be replaced by test.ClientStub
type ClientStub2 struct {
	id              string
	writtenMessages []app.ClientResponse
	sendChan        chan app.AppMessage
	returnChan      chan string
	broadcasters    []app.Broadcasts
}

func (c *ClientStub2) GetID() string {
	return c.id
}

func (c *ClientStub2) Listen() {
	for {
		appMessage := c.ReadMessage()
		message := app.ClientAppMessage{c, appMessage}

		c.broadcastToObservers(message)
	}
}

func (c *ClientStub2) broadcastToObservers(message app.ClientAppMessage) {
	for _, broadcaster := range c.broadcasters {
		broadcaster.Broadcast(message)
	}

	c.returnChan <- ""
}

func (c *ClientStub2) WriteJson(message app.ClientResponse) error {
	c.writtenMessages = append(c.writtenMessages, message)

	return nil
}

func (c *ClientStub2) SendMessage(messageType, payload string) {
	c.sendChan <- app.AppMessage{Type: messageType, Payload: payload}
}

func (c *ClientStub2) ReadMessage() app.AppMessage {
	message := <-c.sendChan

	return message
}

func (c *ClientStub2) Subscribe(b app.Broadcasts) {
	c.broadcasters = append(c.broadcasters, b)
}

func decodePayload(payload string) app.LobbyData {
	var lobbyData app.LobbyData

	raw := []byte(payload)
	json.Unmarshal(
		raw,
		&lobbyData,
	)

	return lobbyData
}

type UserClientServiceStub struct {
	createUserClientCalls  int
	createUserClientReturn error
}

func (r *UserClientServiceStub) CreateAndStoreUserClient(_ *app.User, _ app.IsClient) error {
	r.createUserClientCalls++

	return r.createUserClientReturn
}

func (r *UserClientServiceStub) Resolve(client app.IsClient) (app.UserClient, error) {
	return nil, nil
}
