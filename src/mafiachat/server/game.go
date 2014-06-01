package server

import (
	"encoding/json"
	"errors"
	"log"
)

type game struct {
	Players       []*player      `json:"players"`
	State         string         `json:"state"`
	Id            string         `json:"id"`
	Name          string         `json:"name"`
	MessageBuffer []*chatMessage `json:"messageBuffer"`
}

// Return a new game
func newGame() *game {
	g := &game{}
	g.State = "lobby"
	g.Name = "MafiosoGame"
	g.Players = make([]*player, 0)
	gameList[g.Id] = g
	return g
}

// Add player to game
func (g *game) addPlayer(p *player) {
	// If the player is first, make him admin.
	if len(g.Players) == 0 {
		p.Admin = true
	}
	g.Players = append(g.Players, p)
	go p.msgParser(g)
	if p.Name != "anonymous" {
		g.broadcastGameInfo()
	}
}

// Remove player from game
func (g *game) rmPlayer(p *player) {
	for i := 0; i < len(g.Players); i++ {
		if g.Players[i] == p {
			g.Players = append(g.Players[:i], g.Players[i+1:]...)
			break
		}
	}
	g.broadcastGameInfo()
}

// Broadcast a message to players
func (g *game) broadcastGameInfo() {
	log.Println("broadcastGameInfo")
	gameInfo := &gameInfo{}
	gameInfo.MsgType = "gameInfo"
	gameInfo.Data.Game = g
	msg, err := json.Marshal(gameInfo)
	if err != nil {
		log.Println("Can't marshal gameinfo message to json:", err)
		return
	}
	g.broadcast(msg)
}

func (g *game) loginPlayer(p *player) {
	g.broadcastGameInfo()
	return
}

// Broadcast a message to players
func (g *game) broadcast(msg []byte) {
	for i := 0; i < len(g.Players); i++ {
		g.Players[i].Connection.Outbound <- msg
	}
}

func (g *game) newError(msg string) *chatMessage {
	chatMsg := &chatMessage{}
	chatMsg.MsgType = "chatMessage"
	chatMsg.Data.Message = msg
	chatMsg.Data.Faction = "serverError"
	return chatMsg
}

func (g *game) newInfo(msg string) *chatMessage {
	chatMsg := &chatMessage{}
	chatMsg.MsgType = "chatMessage"
	chatMsg.Data.Message = msg
	chatMsg.Data.Faction = "serverInfo"
	return chatMsg
}

func (g *game) chatMessage(chatMsg *chatMessage, p *player) {
	if g.State == "lobby" {
		chatMsg.Data.Faction = "lobby"
		chatMsg.Data.Player = p
		/*msg, err := json.Marshal(chatMsg)
		if err != nil {
			log.Println("Can't marshal chatMessage to json:", err)
			return
		}*/
		g.MessageBuffer = append(g.MessageBuffer, chatMsg)
		g.broadcastGameInfo()
	}
}

func (g *game) getPlayerByName(s string) (*player, error) {
	for i := 0; i < len(g.Players); i++ {
		log.Println("<" + g.Players[i].Name + "|" + s + ">")
		if g.Players[i].Name == s {
			log.Println("<" + g.Players[i].Name + "> equal to <" + s + ">")
			return g.Players[i], nil
		} else {
			log.Println("<" + g.Players[i].Name + "> not equal to <" + s + ">")
		}
	}
	return nil, errors.New("Can't find player")
}

func (g *game) actionMessage(msg *actionMessage, p *player) {
	log.Println(msg.Data.Target)
	log.Println("action message")
	if msg.Data.Action == "vote" {
		t, err := g.getPlayerByName(msg.Data.Target)
		if err != nil {
			//p.Connection.Outbound <- g.newError(err.Error())
			return
		}
		log.Println("Got t:", t)
		if p.VotingFor != nil {
			p.VotingFor.Votes--
		}
		p.VotingFor = t
		p.VotingFor.Votes++
		g.chatMessage(g.newInfo(p.Name+" votes for "+t.Name+"."), p)
	}
	g.broadcastGameInfo()
}

func (g *game) loginMessage(msg *loginMessage, p *player) {
	log.Println("login message")
	for i := 0; i < len(g.Players); i++ {
		if g.Players[i].Name == msg.Data.Name {
			if g.Players[i].Password == msg.Data.Password {
				// Login succesfull as existing user.
				g.Players[i] = p
				g.broadcastGameInfo()
				return
			} else {
				// Name already exists, but password doesn't match.
				return
			}
		}
	}
	if g.State == "lobby" {
		// New player, game not started
		p.Name = msg.Data.Name
		p.Password = msg.Data.Name
		p.State = "new"
		p.Faction = "villager"
		g.broadcastGameInfo()
	} else if g.State == "game" {
		// New player, game in progress, set as spectator
		p.Name = msg.Data.Name
		p.Password = msg.Data.Name
		p.State = "spectator"
		p.Faction = "ghost"
		g.broadcastGameInfo()
	} else if g.State == "debrief" {
		// Game already over, join the chat anyway.
		p.Name = msg.Data.Name
		p.Password = msg.Data.Name
		p.State = "spectator"
		p.Faction = "ghost"
		g.broadcastGameInfo()
	}
}
