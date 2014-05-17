package server

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// Serve static files as requested
func staticHandler(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Query().Get(":file")
	http.ServeFile(w, r, "./www/"+file)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// Don't redirect favicon.ico requests to random urls
	if r.RequestURI == "/favicon.ico" {
		http.Error(w, "Icon under construction :(", 404)
		return
	}

	// Create random 8 byte hex string for the new game id
	rndUrlBytes := make([]byte, 8)
	n, err := rand.Read(rndUrlBytes)
	if n != len(rndUrlBytes) || err != nil {
		return
	}
	rndUrl := hex.EncodeToString(rndUrlBytes)

	// Redirect browser to the new game's random url
	log.Print("Request on index, generated url: '", rndUrl, "'.")
	http.Redirect(w, r, "/g/"+rndUrl+"/", 303)
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	// Get gameId from mux vars
	vars := mux.Vars(r)
	gameId := vars["gameId"]
	log.Print("New websocket request on '", gameId, "'.")

	// Only get requests
	if r.Method != "GET" {
		return
	}

	// Force same origin policy
	if r.Header.Get("Origin") != "http://"+r.Host {
		http.Error(w, "Origin not allowed", 403)
		return
	}

	// Try to init websocket connection
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "Server Error", 500)
		return
	}
	c := newConnetion(ws)

	// Quick callback for testing
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered in f", r)
			}
		}()
		for {
			c.outbound <- <-c.inbound
		}
	}()

	//hs.addConnection <- &Connection{ws: ws, hubId: hubId}
}

func Init() {
	log.Println("Starting MafiaChat Server")

	// Set mux routes
	r := mux.NewRouter()

	// Handler for static files under gameId url
	r.Path("/g/{gameId:[a-f0-9]{16}}/").
		HandlerFunc(staticHandler).
		Name("static files")

	// Handler for websocket connections
	r.Path("/g/{gameId:[a-f0-9]{16}}/ws/").
		HandlerFunc(websocketHandler).
		Name("websocket")

	// Root url handler
	r.Path("/").
		HandlerFunc(rootHandler).
		Name("root")

	// Start server
	http.Handle("/", r)
	go http.ListenAndServe(":8080", nil)
}
