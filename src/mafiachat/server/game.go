package server

import (
	"encoding/json"
	"errors"
	"log"
	"time"
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
	MessageBuffer []*chatMessage `json:"-"`
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

func (g *game) serverMessage(msg string) {
	chatMsg := &chatMessage{}
	chatMsg.MsgType = "chatMessage"
	chatMsg.Data.Message = msg
	chatMsg.Data.Faction = "server"
	chatMsg.Data.Player = "Server"
	chatMsg.Data.Date = time.Now().Format("15:04:05")
	for i := 0; i < len(g.Players); i++ {
		g.Players[i].addChatMessage(chatMsg)
	}
}

func (g *game) chatMessage(chatMsg *chatMessage, p *player) {
	chatMsg.Data.Player = p.Name
	chatMsg.Data.Date = time.Now().Format("15:04:05")

	if chatMsg.Data.Faction != p.Faction {
		chatMsg.Data.Faction = "villager"
	}

	for i := 0; i < len(g.Players); i++ {
		if chatMsg.Data.Faction == g.Players[i].Faction || chatMsg.Data.Faction == "villager" {
			g.Players[i].addChatMessage(chatMsg)
		}
	}
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

	if g.State == "gameDay" {
		// If the action is vote, it is day and player is not a ghost.
		if msg.Data.Action == "vote" && p.Faction != "ghost" {
			t, err := g.getPlayerByName(msg.Data.Target)
			if err != nil {
				//p.Connection.Outbound <- g.newError(err.Error())
				return
			}
			if p.VotingFor != nil {
				p.VotingFor.Votes--
				g.serverMessage(p.Name + " changes vote to " + t.Name + ".")
			} else {
				g.serverMessage(p.Name + " votes for " + t.Name + ".")
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
					g.serverMessage(g.Players[i].Name + " has been executed.")
					g.Players[i].Faction = "ghost"
					g.State = "gameNight"
					g.zeroVotes()
				}
			}
		}
	} else if g.State == "gameNight" {
		if msg.Data.Action == "vote" && p.Faction == "mafia" {
			t, err := g.getPlayerByName(msg.Data.Target)
			if err != nil {
				return
			}
			if p.VotingFor != nil {
				p.VotingFor.Votes--
				g.serverMessage(p.Name + " changes vote to " + t.Name + ".")
			} else {
				g.serverMessage(p.Name + " votes for " + t.Name + ".")
			}
			p.VotingFor = t
			p.VotingFor.Votes++
			mafiosos := 0
			log.Println("there are", mafiosos, "mafiosos.")
			for i := 0; i < len(g.Players); i++ {
				log.Println(g.Players[i].Name, "is in", g.Players[i].Faction, "faction.")
				if g.Players[i].Faction == "mafia" {
					mafiosos++
				}
			}
			for i := 0; i < len(g.Players); i++ {
				log.Println(g.Players[i].Name, "has", g.Players[i].Votes, "votes.")
				if g.Players[i].Votes == mafiosos {
					g.serverMessage(g.Players[i].Name + " has been found dead.")
					g.Players[i].Faction = "ghost"
					g.State = "gameDay"
					g.zeroVotes()
				}
			}
		}
	} else if msg.Data.Action == "startGame" {
		if p.Admin == true && g.State == "lobby" {
			g.startGame()
			g.serverMessage(p.Name + " has started the game. Good luck.")
		}
	}

	g.broadcastGameInfo()
}

func (g *game) startGame() {
	// TODO: Shuffle roles
	for i := 0; i < len(g.Players); i++ {
		g.Players[i].Faction = "villager"
	}
	g.Players[0].Faction = "mafia"
	g.State = "gameDay"
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
