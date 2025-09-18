package datatypes

type TransportMethod string

const (
	WebHook   TransportMethod = "webhook"
	WebSocket TransportMethod = "websocket"
)

type RequestVersion string

const (
	V1 RequestVersion = "1"
	V2 RequestVersion = "2"
)
