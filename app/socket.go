package app

type Socket interface {
	WriteJSON(v interface{}) error
	ReadJSON(v interface{}) error
}
