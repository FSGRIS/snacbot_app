package server

import (
	"net/url"

	"github.com/gorilla/websocket"
)

type rosServer struct {
	conn *websocket.Conn
}

func newRosServer(addr string) *rosServer {
	u := url.URL{Scheme: "ws", Host: addr}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(err)
	}
	return &rosServer{
		conn: c,
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
