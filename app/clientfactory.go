package app

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type UpgraderFactory interface {
	NewUpgrader() Upgrade
}

type Upgrade interface {
	Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error)
}

type GorillaWebsocketUpgradeFactory struct {
}

func (f *GorillaWebsocketUpgradeFactory) NewUpgrader() Upgrade {
	return &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
}
