package main

import (
	"log"
	"mafiachat/server"
	"runtime"
	"time"
)

func main() {
	server.Init()

	// Loop forever and log goroutine number.
	goRoutines := 0
	for {
		time.Sleep(time.Second * 1)
		if goRoutines != runtime.NumGoroutine() {
			goRoutines = runtime.NumGoroutine()
			log.Println("Goroutines running: ", goRoutines)
		}
	}
}
