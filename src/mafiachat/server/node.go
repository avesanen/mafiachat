package server

import (
	"github.com/garyburd/go-websocket/websocket"
	"io/ioutil"
	"log"
	"time"
)

const (
	writeWait  = 10 * time.Second
	readWait   = 60 * time.Second
	pingPeriod = (readWait * 9) / 10
)

type Node struct {
	ws   *websocket.Conn
	hub  *Hub
	send chan []byte
}

func (n *Node) reader() {
	log.Print("Node reader gorouting starting.")
	defer func() {
		log.Print("Node reader gorouting stopping.")
		n.hub.removeNode <- n
		n.ws.Close()
	}()
	n.ws.SetReadDeadline(time.Now().Add(readWait))
	for {
		op, r, err := n.ws.NextReader()
		if err != nil {
			break
		}
		switch op {
		case websocket.OpPong:
			n.send <- []byte("Ping.\n")
			n.ws.SetReadDeadline(time.Now().Add(readWait))
		case websocket.OpText:
			message, err := ioutil.ReadAll(r)
			if err != nil {
				break
			}
			n.hub.broadcast <- message
		}
	}
}

func (n *Node) write(opCode int, payload []byte) error {
	log.Print("node.write() called")
	n.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return n.ws.WriteMessage(opCode, payload)
}

func (n *Node) writer() {
	log.Print("Node writer gorouting starting.")
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		log.Print("Node writer gorouting stopping.")
		ticker.Stop()
		n.ws.Close()
	}()
	for {
		select {
		case message, ok := <-n.send:
			if !ok {
				n.write(websocket.OpClose, []byte{})
				log.Println("[connection.writePump] !ok.")
				return
			}
			if err := n.write(websocket.OpText, message); err != nil {
				log.Println("[connection.writePump] err: '", err, "'.")
				return
			}
		case <-ticker.C:
			if err := n.write(websocket.OpPing, []byte{}); err != nil {
				log.Println("[connection.writePump] err: '", err, "'.")
				return
			}
		}
	}
}
