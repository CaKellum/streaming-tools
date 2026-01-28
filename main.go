package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/gorilla/websocket"
)

/******************************************************************************************
 * TODO:											                    		          *
 * 	- need to make template and index.html for screens           	                      *
 *  - twitch wesocket will be responsible for reciving chat information                   *
 *  - there will be a separate websocket that will be used to display on a hosted website *
 ******************************************************************************************
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

/*
{
	"metadata":{
		"message_id":"e6151fe6-6b5c-42dd-c0c2-8755601882c5",
		"message_type":"session_welcome",
		"message_timestamp":"2026-01-28T03:06:40.708023Z"
	},
	"payload":{
		"session":{
			"id":"92c282ea_cbf1be21",
			"status":"connected",
			"keepalive_timeout_seconds":300,
			"reconnect_url":null,
			"connected_at":"2026-01-28T03:06:40.707995Z"
		}
	}
}
*/

type WelcomeMessageWrapper struct {
	MetaData WelcomMetaData `json:"metadata"`
	Payload  WelcomePayload `json:"payload"`
}

type WelcomMetaData struct {
	MessageType string `json:"message_type"`
}

type WelcomePayload struct {
	Session WelcomSession `json:"session"`
}

type WelcomSession struct {
	Id           string `json:"id"`
	Status       string `json:"status"`
	ReconnectUrl string `json:"reconnect_url"`
}

var urlSetting = flag.String("mode", "test", "flag determining if it should be a prod or test env")
var twitchWebsocketId string

// MARK: functions
func main() {
	flag.Parse()
	var wsUrl url.URL
	if *urlSetting == "prod" {
		wsUrl = url.URL{
			Scheme:   "wss",
			Host:     "eventsub.wss.twitch.tv",
			Path:     "ws",
			RawQuery: "keepalive_timeout_seconds=300",
		}
	} else {
		wsUrl = url.URL{
			Scheme:   "ws",
			Host:     "localhost:8080",
			Path:     "ws",
			RawQuery: "keepalive_timeout_seconds=300",
		}
	}

	interupt := make(chan os.Signal, 1)
	signal.Notify(interupt, os.Interrupt)

	conn, _, err := websocket.DefaultDialer.Dial(wsUrl.String(), nil)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(message))
			var welcomeMessage WelcomeMessageWrapper
			if err := json.Unmarshal(message, &welcomeMessage); err != nil {
				fmt.Println(err)
			}
			twitchWebsocketId = welcomeMessage.Payload.Session.Id // sets id for websocket comunication moving forward
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interupt:
			fmt.Println("interupted")
			if err := conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
				fmt.Println(err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
		}
	}
}
