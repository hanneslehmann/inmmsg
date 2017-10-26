package inmmsg

import (
	"./server"
	"testing"
	log "github.com/Sirupsen/logrus"
	"os"
)

var logger *log.Logger

func init() {
	logger=log.New()
	logger.Out=os.Stdout
	logger.Level=log.DebugLevel
}

func TestServer_AddListener(t *testing.T) {
	srv := server.New(logger)
	tc := &server.TestClient{
		Name: "Client1",
	}

	tc.AssignToServer(srv)

	e := tc.ListenTo("topic1", func(msg []byte) {})
	if e != nil {
		t.Fatalf("Expected listener to work, but got %s", e)
	}

	if len(srv.Topics["topic1"]) != 1 {
		t.Fatalf("Expected # listeners on topics 1, got %v", len(srv.Topics["topic1"]))
	}
}

func TestServer_SendMessageToTopic_WithoutListeners(t *testing.T) {
	srv := server.New(logger)

	e := srv.SendMessageToTopic("topic1", []byte("test"))

	if e.Error() != "Topic does not exist" {
		t.Fatalf("Expected Error <Topic does not exist>, got %s", e)
	}
}

func TestServer_SendMessageToTopic_WithListener(t *testing.T) {
	srv := server.New(logger)
	tc := &server.TestClient{
		Name: "Client1",
	}

	tc.AssignToServer(srv)
	tc.ListenTo("topic1", func(msg []byte) {
		t.Logf("Client <%s> received <%s>\n", tc.Name, string(msg))
	})

	e := srv.SendMessageToTopic("topic1", []byte("test"))

	if e != nil {
		t.Fatalf("Got error %s", e.Error())
	}
}

func TestServer_SendMessageToTopic_WithListeners(t *testing.T) {
	srv := server.New(logger)
	tc1 := &server.TestClient{
		Name: "Client1",
	}
	tc2 := &server.TestClient{
		Name: "Client2",
	}

	tc1.AssignToServer(srv)
	tc2.AssignToServer(srv)
	tc1.ListenTo("general", func(msg []byte) {
		t.Logf("Client <%s> received <%s>\n", tc1.Name, string(msg))
	})
	tc2.ListenTo("general", func(msg []byte) {
		t.Logf("Client <%s> received <%s>\n", tc2.Name, string(msg))
	})

	t.Log("List of listening clients: ", srv.GetListenersOnTopic("general"))

	if len(srv.Topics["general"]) != 2 {
		t.Fatalf("Expected # listeners on topics 2, got %v", len(srv.Topics["general"]))
	}

	e := srv.SendMessageToTopic("general", []byte("test"))

	if e != nil {
		t.Fatalf("Got error %s", e.Error())
	}
}


func TestServer_ClientSendMessageToTopic_WithListeners(t *testing.T) {
	srv := server.New(logger)
	tc1 := &server.TestClient{
		Name: "Client1",
	}
	tc2 := &server.TestClient{
		Name: "Client2",
	}

	tc1.AssignToServer(srv)
	tc2.AssignToServer(srv)
	tc1.ListenTo("general", func(msg []byte) {
		t.Logf("Client <%s> received <%s>\n", tc1.Name, string(msg))
	})
	tc2.ListenTo("general", func(msg []byte) {
		t.Logf("Client <%s> received <%s>\n", tc2.Name, string(msg))
	})

	t.Log("List of listening clients: ", srv.GetListenersOnTopic("general"))

	if len(srv.Topics["general"]) != 2 {
		t.Fatalf("Expected # listeners on topics 2, got %v", len(srv.Topics["general"]))
	}

	e := tc1.SendToTopic("general", []byte("test"))

	if e != nil {
		t.Fatalf("Got error %s", e.Error())
	}
}