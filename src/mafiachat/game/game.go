package game

type Game struct {
	Id    string `json:"id"`
	State int32
}

type Player struct {
	Id   string `json:"id"`
	Role int32
}

// Enum for game states
const (
	Lobby   = 0
	Day     = 1
	Night   = 2
	Voting  = 3
	DeBrief = 4
)

// Enum for roles
const (
	Dead     = 0
	Villager = 1
	Mafioso  = 2
	Police   = 3
	Doctor   = 4
	Crazy    = 5
)

func Init() *Game {
	return &Game{}
}
