package server

import (
	"bytes"
	"io"
	"log"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

type rosServer struct {
	conn *websocket.Conn
	// Mutex to protect multiple threads making more than one outstanding service call.
	mu        sync.Mutex
	responses chan io.Reader
}

func newRosServer(addr string) *rosServer {
	u := url.URL{Scheme: "ws", Host: addr}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(err)
	}
	// Listen for responses, and block until someone receives on the responses
	// channel. Must only have one thread making a service call at a time.
	responses := make(chan io.Reader)
	go func() {
		for {
			_, b, err := conn.ReadMessage()
			if err != nil {
				log.Printf("[ros server] failed to read response: %s\n", err.Error())
				continue
			}
			responses <- bytes.NewBuffer(b)
		}
	}()
	return &rosServer{
		conn:      conn,
		responses: responses,
	}
}

func (s *rosServer) write(msg interface{}) {
	if err := s.conn.WriteJSON(msg); err != nil {
		panic(err)
	}
}

func (s *rosServer) advertise(topic, msgType string) {
	s.write(dict{
		"op":    "advertise",
		"topic": topic,
		"type":  msgType,
	})
}

func (s *rosServer) publish(topic string, msg interface{}) {
	s.write(dict{
		"op":    "publish",
		"topic": topic,
		"msg":   msg,
	})
}

func (s *rosServer) callService(service string, args []interface{}) io.Reader {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.write(dict{
		"op":      "call_service",
		"service": service,
		"args":    args,
	})
	return <-s.responses
}
