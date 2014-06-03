package server

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var gameList = make(map[string]*game)

func staticHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./www/game.html")
}

func gamesList(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(gameList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
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
	p := newPlayer()
	p.Connection = c

	if _, ok := gameList[gameId]; ok {
		gameList[gameId].addPlayer(p)
	} else {
		g := newGame()
		g.Id = gameId
		gameList[g.Id] = g
		g.addPlayer(p)
	}
}

func Init() {
	log.Println("Starting MafiaChat Server")

	r := mux.NewRouter()

	r.Path("/g/{gameId:[a-f0-9]{16}}/").
		HandlerFunc(staticHandler).
		Name("static files")

	r.Path("/g/{gameId:[a-f0-9]{16}}/ws/").
		HandlerFunc(websocketHandler).
		Name("websocket")

	r.Path("/g/").
		HandlerFunc(rootHandler).
		Name("root")

	r.Path("/games.json").
		HandlerFunc(gamesList).
		Name("gameslist")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./www/")))

	// Start server
	http.Handle("/", r)
	go http.ListenAndServe(":8080", nil)
}
