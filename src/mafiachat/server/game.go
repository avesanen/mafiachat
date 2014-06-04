package server

import (
	"encoding/json"
	"errors"
	"log"
)

/* snippets:
Shuffle
	dest := make([]int, len(src))
	perm := rand.Perm(len(src))
	for i, v := range perm {
	    dest[v] = src[i]
	}
*/
type game struct {
	Players       []*player      `json:"players"`
	State         string         `json:"state"`
	Id            string         `json:"id"`
	Name          string         `json:"name"`
	Password      string         `json:"password"`
	MessageBuffer []*chatMessage `json:"messageBuffer"`
}

// Return a new game
func newGame() *game {
	g := &game{}
	g.State = "lobby"
	g.Name = "MafiosoGame"
	g.Players = make([]*player, 0)
	if g.Id != "" {
		gameList[g.Id] = g
	}
	return g
}

func (g *game) startDay(p *player) {
}

func (g *game) zeroVotes() {
	for i := 0; i < len(g.Players); i++ {
		g.Players[i].Votes = 0
		g.Players[i].VotingFor = nil
	}
}

func (g *game) countFaction(f string) int {
	c := 0
	for i := 0; i < len(g.Players); i++ {
		if g.Players[i].Faction == f {
			c++ // lol
		}
	}
	return c
}

// Add player to game
func (g *game) addPlayer(p *player) {
	// If the player is first, make him admin.
	if len(g.Players) == 0 {
		p.Admin = true
	}
	g.Players = append(g.Players, p)
	p.State = "online"
	go p.msgParser(g)
	g.broadcastGameInfo()
}

// Remove player from game
func (g *game) rmPlayer(p *player) {
	/*
		for i := 0; i < len(g.Players); i++ {
			if g.Players[i] == p {
				g.Players = append(g.Players[:i], g.Players[i+1:]...)
				break
			}
		}*/
	p.State = "offline"
	p.Connection = nil
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
		if g.Players[i].Connection != nil {
			g.Players[i].Connection.Outbound <- msg
		}
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
	}
	g.MessageBuffer = append(g.MessageBuffer, chatMsg)
	if p.Faction == "toBeExecuted" {
		p.Faction = "ghost"
		g.MessageBuffer = append(g.MessageBuffer, g.newInfo(p.Name+" has been executed."))
		g.zeroVotes()
		g.State = "gameDay"
		if g.countFaction("mafia") == 0 {
			g.MessageBuffer = append(g.MessageBuffer, g.newInfo("Villagers win!"))
			g.State = "debrief"
		} else if g.countFaction("villager") <= g.countFaction("mafia") {
			g.MessageBuffer = append(g.MessageBuffer, g.newInfo("Mafioso win!"))
		}
	}
	g.broadcastGameInfo()
}

func (g *game) getPlayerByName(s string) (*player, error) {
	for i := 0; i < len(g.Players); i++ {
		if g.Players[i].Name == s {
			return g.Players[i], nil
		}
	}
	return nil, errors.New("Can't find player")
}

func (g *game) actionMessage(msg *actionMessage, p *player) {
	log.Println(msg.Data.Target)
	log.Println("action message")

	// If the action is vote, it is day and player is not a ghost.
	if msg.Data.Action == "vote" && g.State == "gameDay" && p.Faction != "ghost" {
		if g.State == "gameDay" {
			t, err := g.getPlayerByName(msg.Data.Target)
			if err != nil {
				//p.Connection.Outbound <- g.newError(err.Error())
				return
			}
			if p.VotingFor != nil {
				p.VotingFor.Votes--
				g.chatMessage(g.newInfo(p.Name+" changes vote to "+t.Name+"."), p)
			} else {
				g.chatMessage(g.newInfo(p.Name+" votes for "+t.Name+"."), p)
			}
			p.VotingFor = t
			p.VotingFor.Votes++
			alivePlayers := 0
			for i := 0; i < len(g.Players); i++ {
				if g.Players[i].Faction != "ghost" {
					alivePlayers++
				}
			}

			for i := 0; i < len(g.Players); i++ {
				if g.Players[i].Votes > alivePlayers/2 {
					g.chatMessage(g.newInfo(g.Players[i].Name+" has majority vote. Any last words?"), p)
					g.Players[i].Faction = "toBeExecuted"
					g.State = "gameEvening"
				}
			}
		}
	}

	if msg.Data.Action == "startGame" {
		if p.Admin == true && g.State == "lobby" {
			g.State = "gameDay"
			for i := 0; i < len(g.Players); i++ {
				g.Players[i].Faction = "villager"
			}
			g.chatMessage(g.newInfo(p.Name+" has started the game. Good luck."), p)
		}
	}
	g.broadcastGameInfo()
}

func (g *game) loginMessage(msg *loginMessage, p *player) error {
	log.Println("login message")
	for i := 0; i < len(g.Players); i++ {
		if g.Players[i].Name == msg.Data.Name {
			if g.Players[i].Password == msg.Data.Password {
				if g.Players[i].State == "offline" {
					g.Players[i].Connection = p.Connection
					go g.Players[i].msgParser(g)
					g.Players[i].State = "online"
					g.broadcastGameInfo()
					return nil
				} else {
					return errors.New("Already logged in, kick not supported yet.")
				}
				return nil
			} else {
				// Name already exists, but password doesn't match.
				return errors.New("Wrong password.")
			}
		}
	}
	// New player
	if g.State == "lobby" {
		p.Name = msg.Data.Name
		p.Password = msg.Data.Password
		p.Faction = "villager"
	} else if g.State == "gameDay" {
		p.Name = msg.Data.Name
		p.Password = msg.Data.Password
		p.Faction = "ghost"
	} else if g.State == "debrief" {
		p.Name = msg.Data.Name
		p.Password = msg.Data.Password
		p.Faction = "ghost"
	}
	g.addPlayer(p)
	g.broadcastGameInfo()
	return nil
}
