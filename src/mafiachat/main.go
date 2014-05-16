package main

import (
	"log"
	"mafiachat/game"
	"mafiachat/server"
)

func main() {
	g := game.Init()
	server.Init()
	log.Println(g)
}
