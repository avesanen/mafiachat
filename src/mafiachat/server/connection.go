package server

import (
	"github.com/garyburd/go-websocket/websocket"
)

type Connection struct {
	ws    *websocket.Conn
	hubId string
}
