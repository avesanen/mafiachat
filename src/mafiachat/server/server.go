package server

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/garyburd/go-websocket/websocket"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// ------------ main.go ------------
func nodeWsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hubId := vars["hub"]
	log.Print("Websocket request on '", hubId, "'.")
	log.Print("REQUEST:", r)

	if r.Method != "GET" {
		return
	}
	if r.Header.Get("Origin") != "http://"+r.Host {
		http.Error(w, "Origin not allowed", 403)
		return
	}
	ws, err := websocket.Upgrade(w, r.Header, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "Server Error", 500)
		return
	}
	hs.addConnection <- &Connection{ws: ws, hubId: hubId}
}

func hubHandler(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/favicon.ico" {
		http.Error(w, "Icon under construction :(", 404)
		return
	}
	rndUrlBytes := make([]byte, 8)
	n, err := rand.Read(rndUrlBytes)
	if n != len(rndUrlBytes) || err != nil {
		return
	}
	rndUrl := hex.EncodeToString(rndUrlBytes)
	log.Print("Request on index, generated url: '", rndUrl, "'.")
	http.Redirect(w, r, "/g/"+rndUrl+"/", 303)
}

func nodeHandler(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Query().Get(":file")
	log.Println("nodeHandler")
	log.Print("The whole request: ", r)
	http.ServeFile(w, r, "./www/"+file)
}

func Init() {
	fmt.Println("Roligo v0.0.1")
	log.Println("Starting Roligo Web Server")
	go hs.run()

	r := mux.NewRouter()

	r.Path("/g/{hub:[a-f0-9]{16}}/ws/").
		HandlerFunc(nodeWsHandler).
		Name("websocket")

	r.Path("/g/{hub:[a-f0-9]{16}}/").
		HandlerFunc(nodeHandler).
		Name("static files")

	r.PathPrefix("/js/").
		Handler(http.StripPrefix("/js/",
		http.FileServer(http.Dir("www/js/"))))

	r.PathPrefix("/img/").
		Handler(http.StripPrefix("/img/",
		http.FileServer(http.Dir("www/img/"))))

	r.PathPrefix("/css/").
		Handler(http.StripPrefix("/css/",
		http.FileServer(http.Dir("www/css/"))))

	r.Path("/").
		HandlerFunc(hubHandler).
		Name("root")

	http.Handle("/", r)

	go http.ListenAndServe(":8080", nil)
}
