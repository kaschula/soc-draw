package app_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/kaschula/socket-server/app"
	"github.com/kaschula/socket-server/tests/stubs"
	. "github.com/stretchr/testify/assert"
)

func TestAWebRequestCanBeConvertedToASocketASocketStartsListeningAndAddedToLobby(t *testing.T) {
	req := http.Request{}
	res := stubs.HttpWriterStub{}
	testMessage := "Test Connection"
	lobby := LobbyStub{make(map[string]app.IsClient)}

	client := &ClientStub{"1", []app.ClientResponse{}, make(chan string), make(chan string)}
	clientService := clientServiceStub{client, nil}
	app := app.NewApp(&clientService, &lobby)

	app.CreateSocketHandler(&res, &req)

	resolvedClient, err := lobby.GetClient(client.GetID())

	if err != nil {
		t.Fatal("Client can not be resolved from lobby")
	}

	createdMessage := client.writtenMessages[0]

	client.SendMessage(testMessage)
	receivedMessage := <-client.returnChan

	Equal(t, client, resolvedClient, "Client not added to lobby")
	Equal(t, receivedMessage, testMessage, "Client is not listening")
	Equal(t, createdMessage.Type, "CREATED", "Client did not receive socket Created message")
}

func TestAnErrorIsWrittenToResponseIfClientCanNotBeCreated(t *testing.T) {
	req := http.Request{}
	res := stubs.HttpWriterStub{}
	lobby := LobbyStub{make(map[string]app.IsClient)}
	clientErr := errors.New("Creating Client Error")
	clientFactory := clientServiceStub{nil, errors.New("Creating Client Error")}
	app := app.NewApp(&clientFactory, &lobby)

	app.CreateSocketHandler(&res, &req)

	Equal(t, string(res.WrittenData[0]), "There was an error creating client: "+clientErr.Error(), "Error Response not written")
}

type clientServiceStub struct {
	returnClient app.IsClient
	returnError  error
}

func (c *clientServiceStub) UpgradeToSocket(res http.ResponseWriter, req *http.Request) (app.IsClient, error) {
	if c.returnClient == nil && c.returnError == nil {
		panic("Need to set a return value")
	}

	if c.returnClient == nil && c.returnError != nil {
		return nil, c.returnError
	}

	return c.returnClient, nil
}

type mockClientsStore struct {
	clients []app.IsClient
}

func (cs *mockClientsStore) GetClients() []app.IsClient {
	return cs.clients
}

func (cs *mockClientsStore) AddClient(client app.IsClient) {
	cs.clients = append(cs.clients, client)
}

type ClientStub struct {
	id              string
	writtenMessages []app.ClientResponse
	sendChan        chan string
	returnChan      chan string
}

func (c *ClientStub) GetID() string {
	return c.id
}

func (c *ClientStub) Listen() {
	for {
		c.returnChan <- c.ReadMessage()
	}
}

func (c *ClientStub) WriteJson(message app.ClientResponse) error {
	c.writtenMessages = append(c.writtenMessages, message)

	return nil
}

func (c *ClientStub) SendMessage(message string) {
	c.sendChan <- message
}

func (c *ClientStub) ReadMessage() string {
	message := <-c.sendChan

	return message
}

func (c *ClientStub) Subscribe(b app.Broadcasts) {

}

type LobbyStub struct {
	clients map[string]app.IsClient
}

func (l *LobbyStub) AddClient(c app.IsClient) {
	l.clients[c.GetID()] = c
}

func (lobby *LobbyStub) GetClient(id string) (app.IsClient, error) {
	client, ok := lobby.clients[id]
	if !ok {
		return client, errors.New("Client not found")
	}

	return client, nil
}

func (l *LobbyStub) Broadcast(client app.IsClient, message app.ClientAppMessage) {
	//
}
