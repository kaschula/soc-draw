package app_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kaschula/socket-server/app"
	"github.com/kaschula/socket-server/tests/stubs"
	. "github.com/stretchr/testify/assert"
)

// TODO: Finish this integration tets
// Watchout for the Message types been thrown around, the const and nameing of types is mess
// Need to Clarify all the different message types and try to normalise
// This system may not be asyncrous, after test is passing need to update swimming lane diagram with this test flow
func TestTwoClientsCanJoinTheLobbyAndARoomAndBeginMessagingEachOther(t *testing.T) {
	// ------ Set Up ------
	clientOne, clientTwo := stubs.NewClientStub("C:1"), stubs.NewClientStub("C:2")

	// Create Two Users
	// Each Client / user must have access to the room
	userOne, userTwo := app.NewUser("1"), app.NewUser("2")
	users := map[string]*app.User{
		"1": userOne,
		"2": userTwo,
	}
	userRepository := app.NewInMemoryUserRepository(users)

	simpleApp := app.NewSimpleRoomApplication(stubs.NewPrintsStub())
	room := app.NewRoom("r1", "Chat", 2, 4, simpleApp)
	userRooms := map[*app.User][]app.RoomI{
		userOne: []app.RoomI{room},
		userTwo: []app.RoomI{room},
	}
	roomService := app.NewDefaultRoomService(userRooms)

	userClientService := app.NewInMemoryUserClientService()

	lobby := app.NewLobby(userRepository, roomService, userClientService)

	go simpleApp.Run()
	go clientOne.Listen()
	go clientTwo.Listen()

	// ------ Actions ------
	t.Log("Users request To Join Lobby")
	{
		lobby.AddClient(clientOne)
		lobby.AddClient(clientTwo)

		Equal(t, lobby, clientOne.Lobby, "Lobby should have subscribed to client one")
		Equal(t, lobby, clientTwo.Lobby, "Lobby should have subscribed to client one")

		t.Log("User One (ClientOne) request to joins")
		{

			clientOne.SendMessage(app.MessageTypeLobbyUserJoinRequest, `{"user": "1"}`)
			clientOne.WaitForReturnChan()

			Equal(t, 1, len(clientOne.WrittenMessages),
				"App Message should have been written to ClientOne after User request to join Lobby",
			)
			msg := clientOne.WrittenMessages[0]
			payload := msg.Payload
			Equal(t, app.ClientResponseTypes().USER_LOBBY_DATA, msg.Type, "Should receive Lobby Data after requesting to join")
			True(t, contains(payload, `"User":{"ID":"1"}`), "Payload should contain user data")
			True(t, contains(payload, room.GetID()), "Payload should contain Room Id")

			t.Log("User Two (ClientTwo) request to joins")
			{
				clientTwo.SendMessage(app.MessageTypeLobbyUserJoinRequest, `{"user": "2"}`)
				clientTwo.WaitForReturnChan()

				Equal(t, len(clientOne.WrittenMessages), 1,
					"ClientOne should not have received more messages",
				)
				Equal(t, len(clientTwo.WrittenMessages), 1,
					"App Message should have been written to ClientTwo after User request to join Lobby",
				)
				msg := clientTwo.WrittenMessages[0]
				payload := msg.Payload
				Equal(t, msg.Type, app.ClientResponseTypes().USER_LOBBY_DATA, "Should receive Lobby Data after requesting to join")
				True(t, contains(payload, `"User":{"ID":"2"}`), "Payload should contain user data")
				True(t, contains(payload, room.GetID()), "Payload should contain Room Id")

				t.Log("User One request to Join Room One")
				{
					clientOne.SendMessage(app.MessageTypeJoinRoom, `{"roomId": "r1"}`)
					clientOne.WaitForReturnChan()

					Equal(t, len(clientOne.WrittenMessages), 2,
						"ClientOne should have received extra message",
					)
					Equal(t, len(clientTwo.WrittenMessages), 1,
						"ClientTwo should not have received more messages",
					)

					message := clientOne.WrittenMessages[1]
					fmt.Println(message)
					Equal(t, message.Type, app.ClientResponseTypes().USER_JOINED_ROOM, "Should comfirmation, user join roomed")
					True(t, contains(message.Payload, "r1"), "Message payload should contain success")

					t.Log("User Two request to Join Room One")
					{
						clientTwo.SendMessage(app.MessageTypeJoinRoom, `{"roomId": "r1"}`)
						clientTwo.WaitForReturnChan()

						Equal(t, 3, len(clientTwo.WrittenMessages),
							"ClientTwo should have received lobby and room message",
						)
						Equal(t, 3, len(clientOne.WrittenMessages),
							"ClientOne Also Receive a message",
						)

						clientTwoMessageOne := clientOne.WrittenMessages[1]
						clientTwoMessageTwo := clientOne.WrittenMessages[2]
						clientOneMessage := clientOne.WrittenMessages[2]

						Equal(t, app.ClientResponseTypes().USER_JOINED_ROOM, clientTwoMessageOne.Type,
							"Should comfirmation, user join roomed",
						)
						True(t, contains(clientTwoMessageOne.Payload, "r1"),
							"Message payload should contain success",
						)

						Equal(t, app.ClientResponseTypes().ROOM_BROADCAST_INIT, clientOneMessage.Type,
							"ClientOne Should receive Room Init message room",
						)

						Equal(t, app.ClientResponseTypes().ROOM_BROADCAST_INIT, clientTwoMessageTwo.Type,
							"ClinetTwo Should comfirmation, a user join roomed",
						)

						True(t, contains(clientTwoMessageTwo.Payload, "Initial State"),
							"Message payload should contain welcome",
						)

						t.Log("Clients have Received room init message. Application running...")
						{
							t.Log("ClientOne sends a message to its room")
							{
								firstMessagePayload := "Hi"
								clientOne.SendMessage(app.ClientResponseTypes().ROOM_BROADCAST, firstMessagePayload)
								clientOne.WaitForReturnChan()

								Equal(t, len(clientOne.WrittenMessages), 4, "C1 Should receive its own message")
								Equal(t, len(clientTwo.WrittenMessages), 4, "C2 Should receive a message")
								True(t, contains(clientOne.WrittenMessages[3].Payload, firstMessagePayload), "Should have message Payload")
								True(t, contains(clientOne.WrittenMessages[3].Payload, firstMessagePayload), "Should have message Payload")
							}

							t.Log("ClientTwo replies")
							{
								secondMessagePayload := "Hello"
								clientTwo.SendMessage(app.ClientResponseTypes().ROOM_BROADCAST, secondMessagePayload)
								clientTwo.WaitForReturnChan()

								Equal(t, len(clientOne.WrittenMessages), 5, "C2 Should receive its own message")
								Equal(t, len(clientTwo.WrittenMessages), 5, "C1 Should receive a message")
								True(t, contains(clientOne.WrittenMessages[4].Payload, secondMessagePayload), "Should have message Payload")
								True(t, contains(clientOne.WrittenMessages[4].Payload, secondMessagePayload), "Should have message Payload")
							}
						}
					}
				}
			}
		}
	}
}

func contains(source, target string) bool {
	return strings.Contains(source, target)
}
