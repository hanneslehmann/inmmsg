package client

import (
	"../server"
)

type Client interface {
	AssignToServer(srv *server.Server)
	SendToTopic(topic string, message []byte)
	ListenTo(topic string, callback func(message []byte)) (e error)
}

