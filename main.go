package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
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

/*
	outline of sub request body
	POST https://api.twitch.tv/helix/eventsub/subscriptions
	Client-Id: <AppClientId>
	Authorization: Bearer <UserAccessToken>
	Content-Type: application/json
	{
		"type": "channel.chat.message",
  		"version": "1",
  		"condition": {
    		"broadcaster_user_id": "12345",
    		"user_id": "67890"
  		},
  		"transport": {
    		"method": "websocket",
    		"session_id": <session id sent in welcome message>
  		}
	}
*/

type ConnRequest struct {
	Type      string       `json:"type"`
	Version   string       `json:"version"`
	Condition ReqCondition `json:"condition"`
	Transport ReqTransport `json:"transport"`
}

type ReqCondition struct {
	BroadCasterId string `json:"broadcaster_user_id"`
	UserId        string `json:"user_id"`
}

type ReqTransport struct {
	Method    string `json:"method"`
	SessionId string `json:"session_id"`
}

type Config struct {
	TwitchKey string `env:"TWITCH_KEY"`
	UserID    string `env:"USER_ID"`
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

	recivedWelcome := false

	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(message))
			if !recivedWelcome {
				var welcomeMessage WelcomeMessageWrapper
				if err := json.Unmarshal(message, &welcomeMessage); err != nil {
					fmt.Println(err)
					continue
				}
				twitchWebsocketId = welcomeMessage.Payload.Session.Id // sets id for websocket comunication moving forward
				var websocketRequestConnectURL url.URL
				if *urlSetting == "prod" {
					websocketRequestConnectURL = url.URL{
						Scheme: "https",
						Host:   "api.twitch.tv",
						Path:   "helix/eventsub/subscriptions",
					}
				} else {
					websocketRequestConnectURL = url.URL{
						Scheme: "http",
						Host:   "localhost:2020",
						Path:   "helix/eventsub/subscriptions",
					}
				}
				reqBod := ConnRequest{
					Type:    "channel.chat.message",
					Version: "1",
					Condition: ReqCondition{
						BroadCasterId: "123456",
						UserId:        "34567890",
					},
					Transport: ReqTransport{
						Method:    "websocket",
						SessionId: twitchWebsocketId,
					},
				}

				data, _ := json.Marshal(reqBod)
				req, _ := http.NewRequest(
					http.MethodPost,
					websocketRequestConnectURL.String(),
					strings.NewReader(string(data)))
				req.Header.Add("Authorization", "abc123")
				req.Header.Add("Content-Type", "application/json")
				req.Header.Add("AppClientId", "abs1224")
				resp, _ := (&http.Client{}).Do(req)
				fmt.Println(resp)
				recivedWelcome = true
				continue
			}
			// reciving messages
			var msg Event
			if err := json.Unmarshal(message, &msg); err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(msg.Event.ChatterUserName, ": ", msg.Event.ChatMessage.Text)
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
