package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type IncomingMessage struct {
	Type   string `json:"type"`
	Object struct {
		Message struct {
			FromID  int    `json:"from_id"`
			Text    string `json:"text"`
			Payload string `json:"payload"`
		} `json:"message"`
	} `json:"object"`
}

type OutgoingMessage struct {
	PeerID   int    `json:"peer_id"`
	Message  string `json:"message"`
	Keyboard string `json:"keyboard"`
	RandomID int    `json:"random_id"`
}

type Keyboard struct {
	OneTime bool       `json:"one_time"`
	Buttons [][]Button `json:"buttons"`
}

type Button struct {
	Action Action `json:"action"`
	Color  string `json:"color"`
}

type Action struct {
	Type    string `json:"type"`
	Label   string `json:"label"`
	Payload string `json:"payload"`
}

const (
	groupID   = "220471365"
	serverURL = "https://api.vk.com/method/"
)

var token = ""

func main() {
	file, _ := os.Open("token")
	buf := make([]byte, 220)
	n, _ := file.Read(buf)
	token = string(buf[:n])
	file.Close()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			return
		}

		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var incomingMessage IncomingMessage
		json.Unmarshal(body, &incomingMessage)

		if incomingMessage.Type == "message_new" {
			sendMessage(incomingMessage.Object.Message.FromID)
		}
	})

	http.ListenAndServe(":8080", nil)
}

func sendMessage(peerID int) {
	fmt.Println("@")
	keyboard := Keyboard{
		OneTime: false,
		Buttons: [][]Button{

			{
				{Action: Action{Type: "text", Label: "Button 1"}, Color: "primary"},
				{Action: Action{Type: "text", Label: "Button 2"}, Color: "primary"},
				{Action: Action{Type: "text", Label: "Button 3"}, Color: "primary"},
				{Action: Action{Type: "text", Label: "Button 4"}, Color: "primary"},
			},

			{
				{Action: Action{Type: "text", Label: "Button 5"}, Color: "secondary"},
				{Action: Action{Type: "text", Label: "Button 6"}, Color: "secondary"},
			},
		},
	}

	keyboardJSON, _ := json.Marshal(keyboard)

	outgoingMessage := OutgoingMessage{
		PeerID:   peerID,
		Message:  "Приветствую!",
		Keyboard: string(keyboardJSON),
		RandomID: 0,
	}

	data := url.Values{
		"access_token": {token},
		"group_id":     {groupID},
		"v":            {"5.131"},
		"message":      {outgoingMessage.Message},
		"keyboard":     {outgoingMessage.Keyboard},
		"random_id":    {"0"},
		"peer_id":      {fmt.Sprintf("%d", peerID)},
	}

	_, _ = http.Post(
		serverURL+"messages.send",
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)
}
