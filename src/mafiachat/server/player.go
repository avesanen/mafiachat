package server

import ()

type player struct {
	state      string
	name       string
	password   string
	id         string
	connection connection
}

func newPlayer() *player {
	p := &player{}
	p.state = "new"
	return p
}
