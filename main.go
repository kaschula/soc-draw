package main

import (
	"fmt"
	"net/http"

	socketServer "github.com/kaschula/soc-draw/app"
)

func main() {
	appClientPath := "public/simple-chat"
	appServer := NewSimpleChatApplication()
	go appServer.Run()

	// Data
	userOne := &socketServer.User{"user:1"}
	userTwo := &socketServer.User{"user:2"}
	userThree := &socketServer.User{"user:3"}
	userFour := &socketServer.User{"user:4"}

	users := map[string]*socketServer.User{
		userOne.ID:   userOne,
		userTwo.ID:   userTwo,
		userThree.ID: userThree,
		userFour.ID:  userFour,
	}

	globalRooms := []socketServer.RoomI{
		socketServer.NewDefaultRoom("R:1", "First Room", appServer),
		socketServer.NewDefaultRoom("R:2", "Second Room", appServer),
	}

	roomService := map[*socketServer.User][]socketServer.RoomI{
		userOne:   globalRooms,
		userTwo:   globalRooms,
		userThree: globalRooms,
		userFour:  globalRooms,
	}

	app := newApp(users, roomService)
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
	roomRepo map[*socketServer.User][]socketServer.RoomI,
) *socketServer.App {
	roomService := socketServer.NewDefaultRoomService(roomRepo)
	userService := socketServer.NewInMemoryUserService(users)
	userClientService := socketServer.NewInMemoryUserClientService(nil)

	lobby := socketServer.NewRoomLobby(userService, roomService, userClientService)

	factory := &socketServer.GorillaWebsocketUpgradeFactory{}
	clientService := socketServer.NewDefaultClientService(factory)

	return socketServer.NewApp(clientService, lobby)
}
