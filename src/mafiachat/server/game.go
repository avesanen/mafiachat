package server

import (
	"encoding/json"
	"log"
)

type game struct {
	Players []*player `json:"players"`
	State   string    `json:"state"`
	Id      string    `json:"id"`
	Name	string	  `json:"name"`
}

// Return a new game
func newGame() *game {
	g := &game{}
	g.State = "new"
	g.Name = "Gamee"
	g.Players = make([]*player, 0)
	gameList[g.Id] = g
	return g
}

// Add player to game
func (g *game) addPlayer(p *player) {
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
		log.Println("Can't marshal gameinfo message to json:", msg)
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
