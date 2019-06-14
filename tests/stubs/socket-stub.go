package stubs

import "encoding/json"

func NewSocketStub(writeError error) *SocketStub {
	return &SocketStub{
		[]interface{}{},
		make(chan string),
		make(chan bool),
		writeError,
	}
}

type SocketStub struct {
	Written    []interface{}
	readChan   chan string
	doneChan   chan bool
	writeError error
}

func (w *SocketStub) WriteJSON(v interface{}) error {
	if w.writeError != nil {
		return w.writeError
	}

	w.Written = append(w.Written, v)
	return nil
}

func (w *SocketStub) SendMessage(jsonMessage string) {
	w.readChan <- jsonMessage
}

func (w *SocketStub) ReadJSON(v interface{}) error {
	jsonData := <-w.readChan

	return json.Unmarshal([]byte(jsonData), v)
}
