package server

import (
	"encoding/json"
)

type message struct {
	MsgType string          `json:"msgType"`
	CallbackId int       `json:"callbackId"`
	Data    json.RawMessage `json:"data"`
}

type chatMessage struct {
	MsgType string `json:"msgType"`
	Data    struct {
		Faction string `json:"faction"`
		Message string `json:"message"`
	} `json:"data"`
}

type joinMessage struct {
	MsgType string `json:"msgType"`
	Data    struct {
		Name     string `json:"name"`
		Password string `json:"password"`
		JoinAs   string `json:"joinAs"`
	} `json:"data"`
}

type voteMessage struct {
	MsgType string `json:"msgType"`
	Data    struct {
		Vote string `json:"vote"`
	} `json:"data"`
}

type gameInfo struct {
	MsgType string `json:"msgType"`
	Data    struct {
		Game *game `json:"game"`
	} `json:"data"`
}

type errorMessage struct {
	MsgType string `json:"msgType"`
	Data    struct {
		Message string `json:"vote"`
	} `json:"data"`
}
