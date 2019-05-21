package app_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/kaschula/socket-server/app"
	"github.com/kaschula/socket-server/tests/stubs"
	. "github.com/stretchr/testify/assert"
)

func TestAServiceCanCreateAClient(t *testing.T) {
	req := &http.Request{}
	upgraderStub := &upgradeStub{0}
	factory := &factoryStub{upgraderStub}

	res := &stubs.HttpWriterStub{}
	service := app.NewDefaultClientService(factory)

	client, err := service.UpgradeToSocket(res, req)

	if err != nil {
		t.Fatal("Should not receive error", err.Error())
	}

	Equal(t, len(client.GetID()), 27, "Id should 2 character long")
	True(t, strings.Contains(client.GetID(), "CLIENT:"), "Id should 2 character long")
	Equal(t, 1, upgraderStub.upgradeCall, "Upgrade")
}

type factoryStub struct {
	returnUpgrader app.Upgrade
}

func (f *factoryStub) NewUpgrader() app.Upgrade {
	return f.returnUpgrader
}

type upgradeStub struct {
	upgradeCall int
}

func (u *upgradeStub) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error) {
	u.upgradeCall++

	return &websocket.Conn{}, nil
}
