package datatypes

import "encoding/json"

type TwitchChatCondition struct {
	BroadcasterID string `json:"broadcaster_user_id"`
	UserID        string `json:"user_id"`
}

type TwitchTransport struct {
	Method         TransportMethod `json:"method"`
	Callback       string          `json:"callback"`
	Secret         string          `json:"secret"`
	SessionID      string          `json:"sessin_id"`
	ConnectedAt    string          `json:"connected_at"`
	DisconnectedAt string          `json:"disconnected_at"`
}

type TwitchMessageRequest struct {
	RequestType string              `json:"type"`
	Version     RequestVersion      `json:"version"`
	Condition   TwitchChatCondition `json:"condition"`
	Transport   TwitchTransport     `json:"transport"`
}

type TwitchChat struct {
	Text string `json:"text"`
}

type TwitchChatMesasageEvent struct {
	ChatterName string     `json:"chatter_user_name"`
	MessageId   string     `json:"message_id"`
	Chat        TwitchChat `json:"message"`
	Color       string     `json:"color"`
}

type TwitchChatMessagePayload struct {
	Event TwitchChatMesasageEvent `json:"event"`
}

type TwitchSessionDetails struct {
	SessionID string `json:"id"`
	Status    string `json:"status"`
}
type TwitchWelcomePayload struct {
	Session TwitchSessionDetails `json:"session"`
}

type TwitchWelcomeMessage struct {
	Payload TwitchWelcomePayload `json:"payload"`
}

// MARK: Utility functions
func ToJSON(thing any) ([]byte, error) {
	return json.Marshal(thing)
}

func FromJSON(jsonData []byte, thing *any) error {
	return json.Unmarshal(jsonData, thing)
}
