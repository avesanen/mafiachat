package main

import (
	"log"
	"mafiachat/game"
	"mafiachat/server"
	"runtime"
	"time"
)

func main() {
	g := game.Init()
	server.Init()
	log.Println(g)

	// Loop forever and log goroutine number.
	for {
		time.Sleep(time.Second * 10)
		log.Println("- Goroutines running: ", runtime.NumGoroutine())
	}
}
