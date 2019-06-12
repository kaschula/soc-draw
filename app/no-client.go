package app

import "errors"

const NO_CLIENT_ID = "NO_CLIENT"

// An Empty implementation of IsClient
type NoClient struct {
	id string
}

func NewNoClient() IsClient {
	return &NoClient{NO_CLIENT_ID}
}

func (c *NoClient) GetID() string {
	return c.id
}
func (c *NoClient) Listen() {}
func (c *NoClient) WriteJson(message ClientResponse) error {
	return errors.New("NoClient has no socket to write to")
}
func (c *NoClient) Subscribe(b Broadcasts) {}
func (c *NoClient) SubscribeLobby(b Lobby) {}
