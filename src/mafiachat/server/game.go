package server

import (
	"encoding/json"
	"log"
)

type game struct {
	Players map[string]*player
	State   string
	Id      string
}

// Return a new game
func newGame() *game {
	g := &game{}
	g.State = "new"
	g.Players = make(map[string]*player)
	gameList[g.Id] = g
	return g
}

// Add player to game
func (g *game) addPlayer(p *player) {
	g.Players[p.Id] = p
	go p.msgParser(g)
}

// Remove player from game
func (g *game) rmPlayer(p *player) {
	// Just return if g.players map does not have player.id
	if _, ok := g.Players[p.Id]; !ok {
		return
	}

	// remove item in g.players with key p.id
	delete(g.Players, p.Id)

	// if players map length is 0, the game is empty.
	if len(g.Players) == 0 {
		log.Println("Game", g.Id, "has no more players.")
	}
}

// Broadcast a message to players
func (g *game) broadcast(m *message) {
	s, err := json.Marshal(m)
	if err != nil {
		log.Println("Can't marshal message to json:", m)
		return
	}

	for id, p := range g.Players {
		log.Println("sending to", id)
		p.Connection.Outbound <- string(s)
	}
}
