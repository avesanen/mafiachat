package server

import "time"

type gameInfo struct {
	Name     string        `json:"name"`
	State    string        `json:"state"`
	Messages []*chatInfo   `json:"messages"`
	Players  []*playerInfo `json:"players"`
	MyPlayer *playerInfo   `json:"myPlayer"`
	TimeLeft int           `json:"timeLeft"`
}

func (g *gameInfo) addChatMessage(c *chatInfo) {
	g.Messages = append(g.Messages, c)
}

func (g *gameInfo) addPlayer(p *playerInfo) {
	g.Players = append(g.Players, p)
}

type playerInfo struct {
	Name      string `json:"name"`
	Admin     bool   `json:"admin"`
	Faction   string `json:"faction"`
	Votes     int    `json:"votes"`
	Online    bool   `json:"online"`
	Done      bool   `json:"done"`
	Dead      bool   `json:"dead"`
	Spectator bool   `json:"spectator"`
}

type chatInfo struct {
	Date    string `json:"date"`
	Faction string `json:"faction"`
	Message string `json:"message"`
	Player  string `json:"player"`
}

func getGameInfo(g *game, p *player) *gameInfo {
	// Generate game facts
	gi := &gameInfo{}
	gi.Name = g.Name
	gi.State = g.State
	if g.State != "lobby" && g.State != "debrief" {
		gi.TimeLeft = int((StateTimeout - time.Since(g.StateTime)).Seconds())
	} else {
		gi.TimeLeft = 0
	}
	// Generate player facts
	for i := 0; i < len(g.Players); i++ {
		pi := &playerInfo{}
		pi.Name = g.Players[i].Name
		pi.Dead = g.Players[i].Dead
		pi.Spectator = g.Players[i].Spectator

		// Generate shown faction
		identifiedAs := "unknown"
		if g.State == "debrief" {
			identifiedAs = g.Players[i].Faction
		} else if g.Players[i].Name == p.Name {
			identifiedAs = p.Faction
		} else if p.Spectator {
			identifiedAs = "unknown"
		} else if pi.Dead {
			identifiedAs = "ghost"
		} else if pi.Spectator {
			identifiedAs = "spectator"
		} else if p.Dead {
			identifiedAs = g.Players[i].Faction
		} else if g.Players[i].Faction == p.Faction && p.Faction != "villager" {
			identifiedAs = g.Players[i].Faction
		} else {
			for j := 0; j < len(p.IdentifiedPlayers); j++ {
				if p.IdentifiedPlayers[j].Name == g.Players[i].Name {
					if g.Players[i].Faction == "mafia" {
						identifiedAs = "mafia"
					} else {
						identifiedAs = "villager"
					}
				}
			}
		}
		pi.Faction = identifiedAs

		// Generate "player done" fact.
		playerDone := true
		if g.State == "debrief" {
			playerDone = true
		} else if p.Spectator {
			playerDone = true
		} else if pi.Spectator {
			playerDone = true
		} else if g.State == "day" {
			playerDone = g.Players[i].Done
		} else if g.State == "night" {
			if p.Faction == g.Players[i].Faction {
				playerDone = g.Players[i].Done
			} else {
				playerDone = true
			}
		}
		pi.Done = playerDone

		// Generate shown votes fact
		if p.Spectator {
			pi.Votes = 0
		} else if g.State == "night" && p.Faction == "mafia" {
			pi.Votes = g.Players[i].Votes
		} else if g.State == "day" {
			pi.Votes = g.Players[i].Votes
		} else {
			pi.Votes = 0
		}

		// Generate myplayer fact
		pi.Admin = g.Players[i].Admin
		if pi.Name == p.Name {
			gi.MyPlayer = pi
		}

		// Online status
		pi.Online = g.Players[i].Connection != nil

		gi.addPlayer(pi)
	}

	for i := 0; i < len(g.MessageBuffer); i++ {
		visible := false
		if g.State == "debrief" {
			visible = true
		} else if g.MessageBuffer[i].Data.Faction == p.Faction {
			visible = true
		} else if g.MessageBuffer[i].Data.Faction == "server" {
			visible = true
		} else if g.MessageBuffer[i].Data.Faction == "villager" {
			visible = true
		} else if p.Dead {
			visible = true
		} else if g.State == "debrief" {
			visible = true
		}
		if visible {
			ci := &chatInfo{}
			ci.Date = g.MessageBuffer[i].Data.Date
			ci.Player = g.MessageBuffer[i].Data.Player
			ci.Message = g.MessageBuffer[i].Data.Message
			ci.Faction = g.MessageBuffer[i].Data.Faction
			gi.addChatMessage(ci)
		}
	}
	return gi
}
