package server

type gameInfo struct {
	Name     string        `json:"name"`
	State    string        `json:"state"`
	Messages []*chatInfo   `json:"messages"`
	Players  []*playerInfo `json:"players"`
	MyPlayer *playerInfo   `json:"myPlayer"`
}

func (g *gameInfo) addChatMessage(c *chatInfo) {
	g.Messages = append(g.Messages, c)
}

func (g *gameInfo) addPlayer(p *playerInfo) {
	g.Players = append(g.Players, p)
}

type playerInfo struct {
	Name    string `json:"name"`
	Admin   bool   `json:"admin"`
	Faction string `json:"faction"`
	Votes   int    `json:"votes"`
	Online  bool   `json:"online"`
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
	// Generate player facts
	for i := 0; i < len(g.Players); i++ {
		pi := &playerInfo{}
		pi.Name = g.Players[i].Name
		identified := false
		if g.Players[i].Faction == p.Faction && p.Faction != "villager" {
			identified = true
		}
		if g.Players[i].Name == p.Name {
			identified = true
		}
		for j := 0; j < len(p.IdentifiedPlayers); j++ {
			if p.IdentifiedPlayers[j].Name == g.Players[i].Name {
				identified = true
			}
		}
		if g.Players[i].Faction == "ghost" {
			identified = true
		}
		if identified {
			pi.Faction = g.Players[i].Faction
		} else {
			pi.Faction = "unknown"
		}
		if g.State == "night" && p.Faction == "mafia" {
			pi.Votes = g.Players[i].Votes
		} else if g.State == "day" {
			pi.Votes = g.Players[i].Votes
		} else {
			pi.Votes = 0
		}
		pi.Admin = g.Players[i].Admin
		if pi.Name == p.Name {
			gi.MyPlayer = pi
		}
		pi.Online = g.Players[i].Connection != nil
		gi.addPlayer(pi)
	}

	for i := 0; i < len(g.MessageBuffer); i++ {
		visible := false
		if g.MessageBuffer[i].Data.Faction == p.Faction {
			visible = true
		} else if g.MessageBuffer[i].Data.Faction == "server" {
			visible = true
		} else if g.MessageBuffer[i].Data.Faction == "villager" {
			visible = true
		} else if p.Faction == "ghost" {
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