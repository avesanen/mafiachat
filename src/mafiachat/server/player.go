package server

import (
	"encoding/json"
	"log"
)

type player struct {
	State      string `json:"-"`
	Name       string
	Password   string `json:"-"`
	Id         string
	Connection *connection `json:"-"`
}

func newPlayer() *player {
	p := &player{}
	p.State = "new"
	p.Id = uuid()
	return p
}

func (p *player) msgParser(g *game) {
	log.Println("player msg parser starting")
	defer func() {
		log.Println("player msg parser stopping")
		g.rmPlayer(p)
	}()

	for {
		s, ok := <-p.Connection.Inbound
		if !ok {
			return
		}

		var msg message
		err := json.Unmarshal([]byte(s), &msg)
		if err != nil {
			log.Println("json can't unmarshal message:", s, err)
			continue
		}
		switch {
		case msg.MsgType == "chatMessage":
			log.Println("Chatmessage received", s)
			g.broadcast(&msg)
		case msg.MsgType == "joinGame":
			log.Println("joinGame received", s)
			g.broadcast(&msg)
		default:
			log.Println("Unknown message type '", msg.MsgType, "', ", s, msg)
		}
	}
}
