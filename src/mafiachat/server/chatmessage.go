package server

type chatMessage struct {
	MsgType string `json:"msgType"`
	Data    struct {
		Date    string `json:"date"`
		Faction string `json:"faction"`
		Message string `json:"message"`
		Player  string `json:"player"`
	} `json:"data"`
}

func newMessage(msg string) *chatMessage {
	return nil
}
