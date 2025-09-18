package main

import (
	"fmt"
	"github.com/gorilla/websocket"
)

const (
	twitch_sub_url = "wss://eventsub.wss.twitch.tv/ws?keepalive_timeout_seconds=300"
)

func main() {
	websocket.DefaultDialer.Dial(twitch_sub_url, nil)

	fmt.Println("Hello World!")
}
