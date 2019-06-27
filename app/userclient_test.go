package app_test

import (
	"errors"
	"testing"

	"github.com/kaschula/soc-draw/tests/stubs"

	"github.com/kaschula/soc-draw/app"
	. "github.com/stretchr/testify/assert"
)

func TestThatUserAndClientCanBeRetrieved(t *testing.T) {
	u := app.NewUser("U1")
	c := app.NewDefaultClient("C:1", stubs.NewSocketStub(nil))

	uc := app.NewUserClient(c, u)

	Equal(t, u, uc.GetUser())
	Equal(t, c, uc.GetClient())
}

func TestThatASocketCanBeWrittenTo(t *testing.T) {
	s := stubs.NewSocketStub(nil)
	uc := app.NewUserClient(app.NewDefaultClient("C:1", s), app.NewUser("U1"))

	message := app.ClientResponse{"Type", "Payload", "RoomID"}
	uc.WriteJson(message)

	Equal(t, 1, len(s.Written), "Should have one written message")
}

func TestThatASocketErrorIsReturned(t *testing.T) {
	expectedError := errors.New("Socket Write Error")
	s := stubs.NewSocketStub(expectedError)
	uc := app.NewUserClient(app.NewDefaultClient("C:1", s), app.NewUser("U1"))

	message := app.ClientResponse{"Type", "Payload", "RoomID"}
	err := uc.WriteJson(message)

	Error(t, err, "Socket error should be returned")
	Equal(t, expectedError.Error(), err.Error(), "Error message should be same as expected error message")
}
