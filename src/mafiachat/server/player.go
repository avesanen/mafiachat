package server

import (
	"encoding/json"
	"log"
)

type player struct {
	Id                string      `json:"id"`
	Name              string      `json:"name"`
	Password          string      `json:"-"`
	State             string      `json:"state"`
	Faction           string      `json:"faction"`
	Connection        *connection `json:"-"`
	Votes             int         `json:"votes"`
	VotingFor         *player     `json:"-"`
	Protecting        *player     `json:"-"`
	Admin             bool        `json:"admin"`
	IdentifiedPlayers []*player   `json:"-"`
	Done              bool        `json:"-"`
	Dead              bool        `json:"dead"`
	Spectator         bool        `json:"dead"`
	PlayerIntel       []struct {
		Name    string `json:"name"`
		Faction string `json:"faction"`
		Votes   int    `json:"votes"`
	} `json:"playerIntel"`
}

func newPlayer() *player {
	p := &player{}
	p.Id = uuid()
	p.Done = false
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
			log.Println("json can't unmarshal message:", string(s), err)
			continue
		}
		switch {
		case msg.MsgType == "chatMessage":
			var chatMsg chatMessage
			err := json.Unmarshal([]byte(s), &chatMsg)
			if err != nil {
				log.Println("json can't unmarshal chatMessage:", string(s), err)
			}
			g.chatMessage(&chatMsg, p)
		case msg.MsgType == "actionMessage":
			var actionMsg actionMessage
			err := json.Unmarshal([]byte(s), &actionMsg)
			if err != nil {
				log.Println("json can't unmarshal actionMessage:", string(s), err)
			}
			g.actionMessage(&actionMsg, p)
		default:
			log.Println("Unknown message type ", msg.MsgType, ":", string(s), msg)
			log.Println(p, g)
		}
	}
}
