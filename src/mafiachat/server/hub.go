package server

import (
	"log"
)

type Hub struct {
	hubId      string
	addNode    chan *Node
	removeNode chan *Node
	broadcast  chan []byte
	nodes      map[*Node]bool
	finished   chan bool
}

func (h *Hub) run() {
	log.Println("Hub goroutine starting.")
	defer func() {
		log.Println("Hub goroutine stopping.")
		hs.removeHub <- h.hubId
		close(h.addNode)
		close(h.removeNode)
		close(h.broadcast)
	}()
	for {
		select {
		case n := <-h.addNode:
			log.Println("Hub adding node.")
			h.nodes[n] = true
			go n.writer()
			go n.reader()

		case n := <-h.removeNode:
			if !h.nodes[n] {
				continue
			}
			log.Println("Hub removing node.")
			delete(h.nodes, n)
			close(n.send)
			if len(h.nodes) == 0 {
				return
			}
		case m := <-h.broadcast:
			log.Println("Broadcast message received.")
			for n := range h.nodes {
				select {
				case n.send <- m:
					log.Println("Hub broadcasting.")
				default:
					log.Println("Hub closing node.")
					delete(h.nodes, n)
					close(n.send)
					go n.ws.Close()
				}
			}
		}
	}
}

type HubServer struct {
	hubs          map[string]*Hub
	addConnection chan *Connection
	removeHub     chan string
}

var hs = HubServer{
	hubs:          make(map[string]*Hub),
	addConnection: make(chan *Connection),
	removeHub:     make(chan string),
}

func (hs *HubServer) run() {
	log.Println("Hubserver started.")
	for {
		select {
		case h := <-hs.removeHub:
			delete(hs.hubs, h)
		case c := <-hs.addConnection:
			log.Println("Adding a connection to Hub '" + c.hubId + "'.")
			if hs.hubs[c.hubId] == nil {
				log.Println("Hub '" + c.hubId + "' does not exist, creating new one.")
				hs.hubs[c.hubId] = &Hub{
					hubId:      c.hubId,
					addNode:    make(chan *Node),
					removeNode: make(chan *Node),
					broadcast:  make(chan []byte),
					nodes:      make(map[*Node]bool),
				}
				log.Println("Starting hub '" + c.hubId + "' goroutine.")
				go hs.hubs[c.hubId].run()
			}
			n := &Node{
				ws:   c.ws,
				hub:  hs.hubs[c.hubId],
				send: make(chan []byte),
			}
			hs.hubs[c.hubId].addNode <- n
		}
	}
}
