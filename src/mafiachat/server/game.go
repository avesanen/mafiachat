package server

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
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
	StateTime     time.Time      `json:"-"`
}

const (
	StateTimeout = 15 * time.Minute // 10 minute timeout
)

// Return a new game
func newGame() *game {
	g := &game{}
	g.State = "lobby"
	g.StateTime = time.Now()
	g.Name = "MafiosoGame"
	g.Players = make([]*player, 0)
	if g.Id != "" {
		gameList[g.Id] = g
	}
	return g
}

func (g *game) zeroVotes() {
	for i := 0; i < len(g.Players); i++ {
		g.Players[i].Votes = 0
		g.Players[i].VotingFor = nil
		if g.Players[i].Faction == "ghost" {
			g.Players[i].Done = true
		} else {
			g.Players[i].Done = false
		}
	}
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
	for i := 0; i < len(g.Players); i++ {
		gi := getGameInfo(g, g.Players[i])
		msg, err := json.Marshal(gi)
		if err != nil {
			log.Println("Can't marshal gameinfo message to json:", err)
			return
		}
		if g.Players[i].Connection != nil {
			g.Players[i].Connection.Outbound <- msg
		}
	}
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
	g.MessageBuffer = append(g.MessageBuffer, chatMsg)
}

func (g *game) chatMessage(chatMsg *chatMessage, p *player) {
	chatMsg.Data.Player = p.Name
	chatMsg.Data.Date = time.Now().Format("15:04:05")

	if p.Faction == "ghost" {
		chatMsg.Data.Faction = "ghost"
	} else if chatMsg.Data.Faction != p.Faction {
		chatMsg.Data.Faction = "villager"
	}
	g.MessageBuffer = append(g.MessageBuffer, chatMsg)
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
	switch g.State {

	case "lobby":
		if msg.Data.Action == "startGame" && p.Admin == true {
			g.startGame()
		}

	case "day":
		if msg.Data.Action == "vote" && p.Faction != "ghost" {
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
			p.Done = true
		}
		if g.dayDone() {
			if g.checkVictory() {
				g.endGame()
			} else {
				g.startNight()
			}
		}
		break

	case "night":
		if msg.Data.Action == "vote" && p.Faction == "mafia" {
			t, err := g.getPlayerByName(msg.Data.Target)
			if err != nil {
				return
			}
			if p.VotingFor != nil {
				p.VotingFor.Votes--
				//g.serverMessage(p.Name + " changes vote to " + t.Name + ".")
			} else {
				//g.serverMessage(p.Name + " votes for " + t.Name + ".")
			}
			p.VotingFor = t
			p.VotingFor.Votes++
			p.Done = true
		}
		if msg.Data.Action == "heal" && p.Faction == "doctor" {
			t, err := g.getPlayerByName(msg.Data.Target)
			if err != nil {
				return
			}
			//g.serverMessage(p.Name + " is now protecting " + t.Name + ".")
			p.Protecting = t
			p.Done = true
		}
		if msg.Data.Action == "identify" && p.Faction == "cop" {
			if p.Done {
				return
			}
			t, err := g.getPlayerByName(msg.Data.Target)
			if err != nil {
				return
			}
			//g.serverMessage(p.Name + " identified " + t.Name + ".")
			p.IdentifiedPlayers = append(p.IdentifiedPlayers, t)
			p.Done = true
		}
		if g.nightDone() {
			if g.checkVictory() {
				g.endGame()
			} else {
				g.startDay()
			}
		}
		break
	}
	g.broadcastGameInfo()
}

func (g *game) startGame() {
	if len(g.Players) < 5 {
		g.serverMessage("Can't start game with less than 5 players.")
		return
	}
	for i := 0; i < len(g.Players); i++ {
		g.Players[i].Faction = "villager"
	}
	if len(g.Players) < 6 {
		g.Players[0].Faction = "mafia"
		g.Players[1].Faction = "cop"
		g.Players[2].Faction = "doctor"
	} else if len(g.Players) < 10 {
		g.Players[0].Faction = "mafia"
		g.Players[1].Faction = "mafia"
		g.Players[2].Faction = "cop"
		g.Players[3].Faction = "doctor"
		g.Players[4].Faction = "doctor"
	} else if len(g.Players) < 15 {
		g.Players[0].Faction = "mafia"
		g.Players[1].Faction = "mafia"
		g.Players[2].Faction = "mafia"
		g.Players[3].Faction = "cop"
		g.Players[4].Faction = "cop"
		g.Players[5].Faction = "doctor"
		g.Players[6].Faction = "doctor"
	}

	// shuffle all player factions
	for i := range g.Players {
		g.Players[i].IdentifiedPlayers = nil
		j := rand.Intn(len(g.Players))
		g.Players[i].Faction, g.Players[j].Faction = g.Players[j].Faction, g.Players[i].Faction
	}
	g.startNight()
}

func (g *game) endGame() {
	for i := 0; i < len(g.Players); i++ {
		g.Players[i].Faction = "villager"
	}
	g.State = "lobby"
	g.zeroVotes()
	g.serverMessage("Game has been reset back to lobby, admin can restart the game")
}

func (g *game) checkVictory() bool {
	if g.countFaction("mafia") == 0 {
		g.serverMessage("Last mafioso has died. Good guys win!")
		return true
	}
	// count alivePlayers
	alivePlayers := 0
	for i := 0; i < len(g.Players); i++ {
		if g.Players[i].Faction != "ghost" && g.Players[i].Faction != "mafia" {
			alivePlayers++
		}
	}
	if g.countFaction("mafia") >= alivePlayers {
		g.serverMessage("Too few good guys to finish off mafiosos. Bad guys win!")
		return true
	}
	g.serverMessage("No victor yet, game continues!")
	return false
}

func (g *game) startDay() {
	g.State = "day"
	g.StateTime = time.Now()
	g.zeroVotes()
}

// dayDone will check if someone has majority vote, execute that player
// and return true so the night can begin.
func (g *game) dayDone() bool {
	everyoneReady := true
	for i := 0; i < len(g.Players); i++ {
		if !g.Players[i].Done {
			everyoneReady = false
		}
	}
	if !everyoneReady && time.Since(g.StateTime) < StateTimeout {
		return false
	}

	// count alivePlayers
	alivePlayers := 0
	for i := 0; i < len(g.Players); i++ {
		if g.Players[i].Faction != "ghost" {
			alivePlayers++
		}
	}

	mostVotes := make([]*player, 0)
	votesCount := 0

	for i := 0; i < len(g.Players); i++ {
		if g.Players[i].Votes > votesCount {
			mostVotes = []*player{g.Players[i]}
			votesCount = g.Players[i].Votes
		} else if g.Players[i].Votes == votesCount {
			mostVotes = append(mostVotes, g.Players[i])
		}
	}

	toBeKilled := mostVotes[rand.Intn(len(mostVotes))]
	toBeKilled.Faction = "ghost"
	g.serverMessage(toBeKilled.Name + " was lynched by an angry mob!")

	return true
}

func (g *game) startNight() {
	g.State = "night"
	g.StateTime = time.Now()
	g.zeroVotes()
	for i := 0; i < len(g.Players); i++ {
		if g.Players[i].Faction == "villager" || g.Players[i].Faction == "ghost" {
			g.Players[i].Done = true
		}
	}
}

func (g *game) nightDone() bool {
	everyoneReady := true
	for i := 0; i < len(g.Players); i++ {
		if !g.Players[i].Done {
			everyoneReady = false
		}
	}
	if !everyoneReady && time.Since(g.StateTime) < StateTimeout {
		return false
	}
	mafiosos := g.countFaction("mafia")

	for i := 0; i < len(g.Players); i++ {
		log.Println(g.Players[i].Name, "has", g.Players[i].Votes, "votes.")
		if g.Players[i].Votes == mafiosos {
			playerProtected := false
			for j := 0; j < len(g.Players); j++ {
				log.Println(g.Players[j].Protecting)
				if g.Players[j].Protecting == g.Players[i] {
					log.Println(g.Players[j].Name + " is protecting " + g.Players[j].Protecting.Name)
					log.Println(g.Players[j] == g.Players[i])
					playerProtected = true
				}
			}
			if !playerProtected {
				g.serverMessage(g.Players[i].Name + " has been found dead.")
				g.Players[i].Faction = "ghost"
				return true
			} else {
				g.serverMessage(g.Players[i].Name + " was wounded by mafia, but saved by a doctor.")
				return true
			}
		}
	}
	g.serverMessage("The dawn breaks without victims.")
	return true
}

func (g *game) countFaction(f string) int {
	c := 0
	for i := 0; i < len(g.Players); i++ {
		if g.Players[i].Faction == f {
			c++
		}
	}
	return c
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
					p.Connection.Outbound <- []byte(`{"msgType":"loginFailed", "reason":"alreadyLoggedIn"}`);
					return errors.New("Already logged in and online, kick not supported yet.")
				}
				return nil
			} else {
				// Name already exists, but password doesn't match.
				p.Connection.Outbound <- []byte(`{"msgType":"loginFailed", "reason":"wrongPassword"}`);
				return errors.New("Wrong password.")
			}
		}
	}
	// New player
	if g.State == "lobby" {
		p.Name = msg.Data.Name
		p.Password = msg.Data.Password
		p.Faction = "villager"
	} else if g.State == "day" {
		p.Name = msg.Data.Name
		p.Password = msg.Data.Password
		p.Faction = "ghost"
		p.Done = true
	} else if g.State == "night" {
		p.Name = msg.Data.Name
		p.Password = msg.Data.Password
		p.Faction = "ghost"
		p.Done = true
	}
	g.addPlayer(p)
	return nil
}
