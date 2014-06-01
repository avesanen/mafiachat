package server

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	writeWait  = 10 * time.Second
	readWait   = 10 * time.Second
	pingPeriod = (readWait * 9) / 10
)

// connection type will have outbound and inbound channels, and
// the websocket connection.
type connection struct {
	ws       *websocket.Conn
	Outbound chan []byte
	Inbound  chan []byte
}

// Return a new connection with channels inited and reader/writer
// started. Closing the websocket connection will cause reader
// and writer routines to stop and inbout/outbound channels to close
func newConnetion(ws *websocket.Conn) *connection {
	c := &connection{}
	c.ws = ws
	c.Outbound = make(chan []byte)
	c.Inbound = make(chan []byte)
	go c.reader()
	go c.writer()
	return c
}

// reader is started as a routine, it will continue to read data from
// websocket connection and sends it to the connections inbound channel
// as strings
func (c *connection) reader() {
	log.Print("connection reader gorouting starting.")
	defer func() {
		log.Print("connection reader gorouting stopping.")
		// Close in and outbound channels to message listening goroutines
		// that this connection has closed
		close(c.Inbound)
		c.ws.Close()
	}()
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(writeWait))
		return nil
	})
	c.ws.SetReadDeadline(time.Now().Add(readWait))
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		c.Inbound <- message
	}
}

// Write message as byte array to connection, with messagetype
func (c *connection) write(mt int, payload []byte) error {
	//log.Print("connection.write() called")
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// Routine to continue to write from outbound channel to websocket
// connection. Will close outbound channel when closed.
func (c *connection) writer() {
	log.Print("connection writer gorouting starting.")
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		log.Print("connection writer gorouting stopping.")
		ticker.Stop()
		close(c.Outbound)
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.Outbound:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				log.Println("[connection.writePump] !ok.")
				return
			}
			if err := c.write(websocket.TextMessage, []byte(message)); err != nil {
				log.Println("[connection.writePump] err: '", err, "'.")
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				log.Println("[connection.writePump] ticker err: '", err, "'.")
				log.Println("ping")
				return
			}
		}
	}
}
