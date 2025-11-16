package ws

import "encoding/json"

type WSRequest struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}

type WSMessage struct {
	Action  string `json:"action"`
	Success bool   `json:"success,omitempty"`
	Error   string `json:"error,omitempty"`
	Topic   string `json:"topic,omitempty"`
	Data    any    `json:"data,omitempty"`
}
