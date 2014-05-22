package server

import (
	"encoding/json"
	"log"
)

type player struct {
	Id         string      `json:"id"`
	Name       string      `json:"name"`
	Password   string      `json:"-"`
	State      string      `json:"state"`
	Connection *connection `json:"-"`
}

func newPlayer() *player {
	p := &player{}
	p.State = "new"
	p.Name = "anonymous"
	p.Id = uuid()
	return p
}

func (p *player) sendError(error string) {
	errorMsg := &errorMessage{}
	errorMsg.MsgType = "error"
	errorMsg.Data.Message = error
	errorMsgJson, err := json.Marshal(errorMsg)
	if err != nil {
		log.Println("json can't marshal errormsg", error)
		return
	}
	p.Connection.Outbound <- errorMsgJson
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
			log.Println("Chatmessage received", string(s))
			g.broadcast([]byte(s))
		case msg.MsgType == "login":
			log.Println("login received", string(s))
			var login loginMessage
			err := json.Unmarshal([]byte(s), &login)
			if err != nil {
				log.Println("json can't unmarshal loginmessage", s, err)
			}
			p.Name = login.Data.Name
			p.Password = login.Data.Password
			g.loginPlayer(p)
		default:
			log.Println("Unknown message type ", msg.MsgType, ",", s, msg)
		}
	}
}
