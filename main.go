package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	datatypes "streaming-tools/dataTypes"

	"github.com/gorilla/websocket"
)

// MARK: Constants
const (
	twitch_sub_url = "wss://eventsub.wss.twitch.tv/ws?keepalive_timeout_seconds=300"
)

// MARK: Global VARS
var websocketID string
var broadcaster_user_id string

// MARK: functions
func main() {
	conn, resp, err := websocket.DefaultDialer.Dial(twitch_sub_url, nil)
	defer conn.Close()
	if err != nil && resp.StatusCode >= 200 && resp.StatusCode < 299 {
		fmt.Println(err)
		return
	}

	lengthStr := resp.Header.Get("Content-Length")
	length, _ := strconv.Atoi(lengthStr)
	buf := make([]byte, length)
	resp.Body.Read(buf)
	var welcomeMSG datatypes.TwitchWelcomeMessage
	fmt.Println(string(buf))
	if welcomeErr := json.Unmarshal(buf, &welcomeMSG); welcomeErr != nil {
		fmt.Println(welcomeErr)
		return
	}
	if welcomeMSG.Payload.Session.Status != "connected" {
		fmt.Println("issues may arrise")
	}
	websocketID = welcomeMSG.Payload.Session.SessionID

	data, err := datatypes.ToJSON(
		datatypes.TwitchMessageRequest{
			RequestType: "channel.chat.message",
			Version:     datatypes.V1,
			Condition:   datatypes.TwitchChatCondition{},
			Transport: datatypes.TwitchTransport{
				Method:    datatypes.WebSocket,
				SessionID: websocketID,
			},
		})
	if err != nil {
		fmt.Println(err)
		return
	}
	conn.WriteMessage(0, data)

	fmt.Println("Hello World!")
}
