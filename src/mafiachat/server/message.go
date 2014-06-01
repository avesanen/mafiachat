package server

import (
	"encoding/json"
)

type message struct {
	MsgType string          `json:"msgType"`
	Data    json.RawMessage `json:"data"`
}

type actionMessage struct {
	MsgType string `json:"msgType"`
	Data    struct {
		Action string `json:"action"`
		Target string `json:"target"`
	} `json:"data"`
}

type gameInfo struct {
	MsgType string `json:"msgType"`
	Data    struct {
		Game *game `json:"game"`
	} `json:"data"`
}

type serverMessage struct {
	MsgType string `json:"msgType"`
	Data    struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"data"`
}

type loginMessage struct {
	MsgType string `json:"msgType"`
	Data    struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	} `json:"data"`
}
