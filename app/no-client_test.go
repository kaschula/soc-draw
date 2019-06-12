package app_test

import (
	"testing"

	"github.com/kaschula/socket-server/app"
	. "github.com/stretchr/testify/assert"
)

func TestThatNoClientReturnsID(t *testing.T) {
	client := app.NewNoClient()

	Equal(t, app.NO_CLIENT_ID, client.GetID(), "NoClient Should return NO_CLIENT_ID")
}

func TestThatNoClientAlwaysReturnsErrorWhenWritingJson(t *testing.T) {
	client := app.NewNoClient()
	err := client.WriteJson(app.ClientResponse{})

	Error(t, err, "WriteJson() should always return an error")
}
