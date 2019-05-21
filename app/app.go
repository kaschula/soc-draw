package app

import (
	"net/http"
)

// NewApp notes Change the name app ^^ and below
func NewApp(clientService ClientService, lobby Lobby) *App {
	return &App{clientService, lobby}
}

type App struct {
	clientService ClientService
	lobby         Lobby
}

func (a *App) CreateSocketHandler(res http.ResponseWriter, req *http.Request) {
	client, err := a.createSocketFromRequest(res, req)

	if err != nil { // || client == nil
		res.Write([]byte("There was an error creating client: " + err.Error()))
		return
	}

	a.addClientToLobby(client)

	go client.Listen()
	client.WriteJson(welcomeMessage())
}

func (a *App) createSocketFromRequest(res http.ResponseWriter, req *http.Request) (IsClient, error) {
	return a.clientService.UpgradeToSocket(res, req)
}

func (a *App) addClientToLobby(client IsClient) {
	a.lobby.AddClient(client)
}

type Broadcasts interface {
	Broadcast(client IsClient, message ClientAppMessage)
}
