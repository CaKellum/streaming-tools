package main

import (
	"github.com/caarlos0/env/v11"
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

type Event struct {
	Event MessageEvent `json:"event"`
}

type MessageEvent struct {
	ChatterUserName string  `json:"chatter_user_name"`
	ChatMessage     Message `json:"message"`
	Color           string  `json:"color"`
}

type Message struct {
	Text string `json:"text"`
}

// MARK: Constants
const (
	twitch_sub_url = "wss://eventsub.wss.twitch.tv/ws?keepalive_timeout_seconds=300"
)

// MARK: Global VARS
var websocketID string

// MARK: functions
func main() {

}
