package server

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"time"
)

type game struct {
	Players       []*player      `json:"players"`
	State         string         `json:"state"`
	Id            string         `json:"id"`
	Name          string         `json:"name"`
	Password      string         `json:"-"`
	MessageBuffer []*chatMessage `json:"-"`
	StateTime     time.Time      `json:"-"`
	Winner        string         `json:"-"`
}

const (
	StateTimeout = 15 * time.Minute // 10 minute timeout
)

// Return a new game
func newGame() *game {
	rand.Seed(time.Now().UTC().UnixNano())
	g := &game{}
	g.State = "lobby"
	g.StateTime = time.Now()
	g.Name = "MafiosoGame"
	g.Winner = ""
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
		if g.Players[i].Dead || g.Players[i].Spectator {
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
	// If game is at lobby, remove the player completely (also freeing up the name).
	if g.State == "lobby" || p.Spectator {
		for i := range g.Players {
			if g.Players[i] == p && !p.Admin {
				g.Players = append(g.Players[:i], g.Players[i+1:]...)
				g.broadcastGameInfo()
				return
			}
		}
	}
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

	if p.Dead {
		chatMsg.Data.Faction = "ghost"
	} else if chatMsg.Data.Faction != p.Faction {
		chatMsg.Data.Faction = "villager"
	}
	g.MessageBuffer = append(g.MessageBuffer, chatMsg)
	g.checkCycle()
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

	case "debrief":
		if msg.Data.Action == "startGame" && p.Admin == true {
			g.startGame()
		}

	case "day":
		if msg.Data.Action == "vote" && !p.Dead {
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
		g.checkCycle()
		break

	case "night":
		if msg.Data.Action == "vote" && p.Faction == "mafia" {
			t, err := g.getPlayerByName(msg.Data.Target)
			if err != nil {
				return
			}
			if p.VotingFor != nil {
				p.VotingFor.Votes--
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
			if p.VotingFor != nil {
				p.VotingFor.Votes--
			}
			p.VotingFor = t
			p.VotingFor.Votes++
			p.Done = true
		}
		if msg.Data.Action == "identify" && p.Faction == "cop" {
			t, err := g.getPlayerByName(msg.Data.Target)
			if err != nil {
				return
			}
			if p.VotingFor != nil {
				p.VotingFor.Votes--
			}
			p.VotingFor = t
			p.VotingFor.Votes++
			p.Done = true
		}
		g.checkCycle()
		break
	}
	g.broadcastGameInfo()
}

func (g *game) checkCycle() {
	switch g.State {
	case "day":
		if g.dayDone() {
			if g.checkVictory() {
				g.endGame()
			} else {
				g.startNight()
			}
		}
	case "night":
		if g.nightDone() {
			if g.checkVictory() {
				g.endGame()
			} else {
				g.startDay()
			}
		}
	}
}

func (g *game) startGame() {
	g.State = "lobby"
	g.Winner = ""
	players := g.Players
	for i := range players {
		if players[i].State == "offline" {
			plr, err := g.getPlayerByName(players[i].Name)
			if err == nil {
				g.rmPlayer(plr)
			}
		}
	}

	if len(g.Players) < 5 {
		g.serverMessage("Can't start game with less than 5 players.")
		return
	}

	for i := 0; i < len(g.Players); i++ {
		g.Players[i].Faction = "villager"
		g.Players[i].Spectator = false
		g.Players[i].Dead = false
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
	} else {
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
	g.MessageBuffer = make([]*chatMessage, 0)
	g.startNight()
}

func (g *game) endGame() {
	for i := range g.Players {
		g.Players[i].Dead = false
	}
	g.State = "debrief"
	g.zeroVotes()
	g.serverMessage("Game has been reset back to lobby, admin can restart the game")
}

func (g *game) checkVictory() bool {
	if g.countFaction("mafia") == 0 {
		g.serverMessage("Last mafioso has died. Good guys win!")
		g.Winner = "villager"
		return true
	}
	// count alivePlayers
	alivePlayers := 0
	for i := 0; i < len(g.Players); i++ {
		if !g.Players[i].Dead && !g.Players[i].Spectator && g.Players[i].Faction != "mafia" {
			alivePlayers++
		}
	}
	if g.countFaction("mafia") >= alivePlayers {
		g.serverMessage("Too few good guys to finish off mafiosos. Bad guys win!")
		g.Winner = "mafia"
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

// countVotesFor takes faction as a string, and returns the leading player of the vote
// if the faction is villager, all votes are counted.
func (g *game) countVotesFor(f string) ([]*player, int) {
	players := make([]*player, 0)
	maxVotes := 0
	for i := range g.Players {
		if !g.Players[i].Dead && !g.Players[i].Spectator {
			votes := 0
			for j := range g.Players {
				if !g.Players[i].Dead && !g.Players[i].Spectator &&
					(g.Players[j].Faction == f || f == "villager") &&
					g.Players[j].VotingFor == g.Players[i] {
					votes++
				}
			}
			if votes > maxVotes {
				players = []*player{g.Players[i]}
				maxVotes = votes
			} else if votes == maxVotes {
				players = append(players, g.Players[i])
			}
		}
	}
	return players, maxVotes
}

// countFaction returns the number of players alive
func (g *game) alivePlayers() int {
	alivePlayers := 0
	for i := 0; i < len(g.Players); i++ {
		if !g.Players[i].Dead && !g.Players[i].Spectator {
			alivePlayers++
		}
	}
	return alivePlayers
}

// countFaction returns the number of players in the faction
func (g *game) countFaction(f string) int {
	factionPlayers := 0
	for i := 0; i < len(g.Players); i++ {
		if !g.Players[i].Dead && !g.Players[i].Spectator && g.Players[i].Faction == f {
			factionPlayers++
		}
	}
	return factionPlayers
}

// dayDone will check if someone has majority vote, execute that player
// and return true so the night can begin.
func (g *game) dayDone() bool {
	plr, vts := g.countVotesFor("villager")
	if vts <= g.alivePlayers()/2 && time.Since(g.StateTime) < StateTimeout {
		return false
	}

	if vts == 0 {
		g.serverMessage("Villages didn't vote, nobody dies.")
		return true
	}

	toBeKilled := plr[rand.Intn(len(plr))]
	toBeKilled.Dead = true
	g.serverMessage(toBeKilled.Name + " was lynched by an angry mob!")
	return true
}

func (g *game) startNight() {
	g.State = "night"
	g.StateTime = time.Now()
	g.zeroVotes()
	for i := 0; i < len(g.Players); i++ {
		if g.Players[i].Faction == "villager" || g.Players[i].Dead || g.Players[i].Spectator {
			g.Players[i].Done = true
		}
	}
}

func (g *game) nightDone() bool {
	mafiaQuorum := true
	doctorQuorum := true
	copQuorum := true

	var mafiaTarget *player = nil
	var doctorTarget *player = nil
	var copTarget *player = nil

	plr, vts := g.countVotesFor("mafia")
	if vts <= g.countFaction("mafia")/2 && g.countFaction("mafia") > 0 {
		mafiaQuorum = false
	}
	if len(plr) != 0 && vts > 0 {
		mafiaTarget = plr[rand.Intn(len(plr))]
	}
	log.Println("mafia votes", vts)

	plr, vts = g.countVotesFor("doctor")
	if vts <= g.countFaction("doctor")/2 && g.countFaction("doctor") > 0 {
		doctorQuorum = false
	}
	if len(plr) != 0 && vts > 0 {
		doctorTarget = plr[rand.Intn(len(plr))]
	}
	log.Println("doctor votes", vts)

	plr, vts = g.countVotesFor("cop")
	if vts <= g.countFaction("cop")/2 && g.countFaction("cop") > 0 {
		copQuorum = false
	}
	if len(plr) != 0 && vts > 0 {
		copTarget = plr[rand.Intn(len(plr))]
	}
	log.Println("cop votes", vts)

	log.Println(mafiaQuorum, doctorQuorum, copQuorum)
	log.Println(mafiaTarget, doctorTarget, copTarget)
	if (!mafiaQuorum || !doctorQuorum || !copQuorum) && time.Since(g.StateTime) < StateTimeout {
		return false
	}

	if mafiaTarget != nil && doctorTarget != mafiaTarget {
		g.serverMessage(mafiaTarget.Name + " was found dead.")
		mafiaTarget.Dead = true
	} else {
		g.serverMessage("Night ends without casualties.")
	}

	if copTarget != nil {
		for i := range g.Players {
			if g.Players[i].Faction == "cop" {
				g.Players[i].IdentifiedPlayers = append(g.Players[i].IdentifiedPlayers, copTarget)
			}
		}
	}
	return true
}

func (g *game) loginMessage(msg *loginMessage, p *player) error {
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
					p.Connection.Outbound <- []byte(`{"msgType":"loginFailed", "reason":"alreadyLoggedIn"}`)
					return errors.New("Already logged in and online, kick not supported yet.")
				}
				return nil
			} else {
				// Name already exists, but password doesn't match.
				p.Connection.Outbound <- []byte(`{"msgType":"loginFailed", "reason":"wrongPassword"}`)
				return errors.New("Wrong password.")
			}
		}
	}
	// New player
	if g.State == "lobby" || g.state == "debrief" {
		p.Name = msg.Data.Name
		p.Password = msg.Data.Password
		p.Faction = "villager"
		p.Dead = false
		p.Spectator = false
		p.State = "online"
	} else if g.State == "day" || g.State == "night" {
		p.Name = msg.Data.Name
		p.Password = msg.Data.Password
		p.Faction = "spectator"
		p.Done = true
		p.Dead = false
		p.Spectator = true
		p.State = "online"
	}
	g.addPlayer(p)
	return nil
}
