package app_test

import (
	"testing"

	"github.com/kaschula/soc-draw/app"
	"github.com/kaschula/soc-draw/tests/stubs"
	. "github.com/stretchr/testify/assert"
)

func TestThatASimpleRoomAppWritesMessagesToOneRoom(t *testing.T) {
	logger := &PrintsStub{0}
	application := app.NewSimpleRoomApplication(logger)
	room := stubs.NewRoomStub("room:1", "FirstRoom")
	appMessage := app.NewAppMessage(app.GetResponseTypes().ROOM_BROADCAST, "payload")
	go application.Run()

	application.WriteMessage(app.NewRoomMessage(room, appMessage.Payload))
	<-room.BroadcastReturn

	Equal(t, logger.called, 1, "logger should have been called")
	Equal(t, room.BroadcastCalled, 1, "Room::Broadcast was not called")
	Equal(t, room.BroadcastData[0].AppMessage, appMessage, "AppMessage Client App message should match room payload")
}

func TestThatASimpleRoomAppStartsRoom(t *testing.T) {
	logger := &PrintsStub{0}
	application := app.NewSimpleRoomApplication(logger)
	room := stubs.NewRoomStub("room:1", "FirstRoom")
	expectedMessage := app.NewAppMessage(app.GetResponseTypes().ROOM_BROADCAST_INIT, `{"message":"welcome"}`)

	go application.Start(room)
	<-room.BroadcastReturn

	Equal(t, len(room.BroadcastData), 1, "Room should have a broadcast message")
	Equal(t, room.BroadcastData[0].AppMessage.Type, expectedMessage.Type, "Welcome Message Should have been broadcast")
}

type PrintsStub struct {
	called int
}

func (p *PrintsStub) Printf(format string, v ...interface{}) {
	p.called++
}
