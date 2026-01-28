package main

import (
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

// MARK: functions
func main() {
	url := url.URL{
		Scheme:   "wss",
		Host:     "eventsub.wss.twitch.tv",
		Path:     "ws",
		RawQuery: "keepalive_timeout_seconds=300",
	}

	interupt := make(chan os.Signal, 1)
	signal.Notify(interupt, os.Interrupt)

	conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
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
