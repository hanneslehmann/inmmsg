package server

import "errors"

// Client implementation for testing
type TestClient struct {
	Name   string
	server *Server
}

func (tc *TestClient) AssignToServer(srv *Server) {
	tc.server = srv
}

func (tc *TestClient) SendToTopic(topic string, message []byte) (e error){
	e=tc.server.SendMessageToTopic(topic, message)
	return
}

func (tc *TestClient) ListenTo(topic string, cb func(message []byte)) (e error) {
	if tc.server == nil {
		return errors.New("Cannot listen without assigned server!")
	}
	tc.server.AddListener(topic, CallBack{
		Name:     tc.Name,
		Callback: cb,
	})
	return nil
}