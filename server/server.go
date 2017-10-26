package server

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"sync"
)

type Server struct {
	sync   sync.Mutex
	log    *log.Logger
	Topics map[string][]struct {
		Name     string
		callback func(msg []byte)
	}
}

func New() *Server {
	var l = log.New()
	l.Level = log.DebugLevel
	return &Server{
		log: l,
		Topics: make(map[string][]struct {
			Name     string
			callback func(msg []byte)
		}),
	}
}

func (s *Server) AddListener(topic string, fn struct {
	Name     string
	callback func(msg []byte)
}) {
	s.sync.Lock()
	tmp, ok := s.Topics[topic]
	if ok {
		tmp = append(tmp, fn)
		s.Topics[topic] = tmp
	} else {
		s.Topics[topic] = []struct {
			Name     string
			callback func(msg []byte)
		}{fn}
	}
	s.sync.Unlock()
}

func (s *Server) GetListenersOnTopic(topic string) (list []string) {
	var l []string
	s.sync.Lock()
	tmp, ok := s.Topics[topic]
	s.sync.Unlock()
	if ok {
		if len(tmp) < 1 {
			return nil
		}
		for _, c := range tmp {
			l = append(l, c.Name)
		}
	} else {
		return nil
	}
	return l
}

func (s *Server) SendMessageToTopic(topic string, msg []byte) (e error) {
	s.sync.Lock()
	tmp, ok := s.Topics[topic]
	s.sync.Unlock()
	if ok {
		if len(tmp) < 1 {
			return errors.New("No listeners on topic " + topic)
		}
		var wg sync.WaitGroup
		for _, c := range tmp {
			wg.Add(1)
			go func(m []byte, name string, callback func(msg []byte), w *sync.WaitGroup) {
				s.log.Debugf("Sending message <%s> to client <%s>", string(m), name)
				callback(m)
				w.Done()
			}(msg, c.Name, c.callback, &wg)
		}
		wg.Wait()
		return nil
	} else {
		return errors.New("Topic does not exist")
	}
}
