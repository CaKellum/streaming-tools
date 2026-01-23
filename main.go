package main

import (
	"encoding/json"
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/gorilla/websocket"
	"strconv"
	datatypes "streaming-tools/dataTypes"
)

/*****************************************************************************************************************
 * TODO:																										 *
 *  - useable way to store creds in the application																 *
 *		> likely will use some dotenv package or roll my own													 *
 * 	- need to make template and index.html for screens															 *
 * 	- make a updatable webpage via net/http																		 *
 *		> web can do obs integration but could use the chance to use bobatea more and make tui version of client *
 *  - make it so on websocket recipt we push new screen to website												 *
 *  - <Nice To Have> Random Color selector function																 *
 *****************************************************************************************************************
 */

type Config struct {
	TwitchKey string `env:"TWITCH_KEY"`
}

func getTwitchKey() (string, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return "", err
	}
	return cfg.TwitchKey, nil
}

// MARK: Constants
const (
	twitch_sub_url = "wss://eventsub.wss.twitch.tv/ws?keepalive_timeout_seconds=300"
)

// MARK: Global VARS
var websocketID string

// MARK: functions
func main() {

	// TODO: Want to move this some where else then have a list of messages and a Permenant reference to client
	// 		 that when a new message comes in we write a new view to it rendering in UTF-8 for simplicity (EnGlIsH #1!)
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
