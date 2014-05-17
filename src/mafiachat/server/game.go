package server

import (
	"encoding/json"
	"log"
)

type game struct {
	players map[string]*player
	state   string
	id      string
}

// Return a new game
func newGame() *game {
	g := &game{}
	g.state = "new"
	return g
}

// Add player to game
func (g *game) addPlayer(p *player) {
	g.players[p.id] = p
}

// Remove player from game
func (g *game) rmPlayer(p *player) {
	// Just return if g.players map does not have player.id
	if _, ok := g.players[p.id]; !ok {
		return
	}

	// remove item in g.players with key p.id
	delete(g.players, p.id)

	// if players map length is 0, the game is empty.
	if len(g.players) == 0 {
		log.Println("Game", g.id, "has no more players.")
	}
}

// Broadcast a message to players
func (g *game) broadcast(m *message) {
	s, err = json.Marshal(m)
	if err != nil {
		log.Println("Can't marshal message to json:", m)
		return
	}
	log.Println(s)
}

func (g *game) sendMessage(m *message, p *player) {
	s, err = json.Marshal(m)
	if err != nil {
		log.Println("Can't marshal message to json:", m)
		return
	}
	p.connection.outbound <- s
}

func (g *game) parseMessage(s string) {
	var m message
	err = json.Unmarshal([]byte(s), &m)
	if err != nil {
		return
	}
}
