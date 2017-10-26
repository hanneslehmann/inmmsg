package server

import (
	"errors"
)

type Client interface {
	AssignToServer(srv *Server)
	SendToTopic(topic string, message []byte)
	ListenTo(topic string, callback func(message []byte)) (e error)
}

// Client implementation for testing
type testClient struct {
	Name   string
	server *Server
}

func (tc *testClient) AssignToServer(srv *Server) {
	tc.server = srv
}

func (tc *testClient) SendToTopic(topic string, message []byte) {
}

func (tc *testClient) ListenTo(topic string, cb func(message []byte)) (e error) {
	if tc.server == nil {
		return errors.New("Cannot listen without assigned server!")
	}
	tc.server.AddListener(topic, struct {
		Name     string
		callback func(msg []byte)
	}{
		Name:     tc.Name,
		callback: cb,
	})
	return nil
}
