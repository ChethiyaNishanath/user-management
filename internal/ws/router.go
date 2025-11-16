package ws

import (
	"context"
	"encoding/json"
	"log/slog"

	wsutils "user-management/internal/common/wsutils"

	"github.com/coder/websocket"
)

type HandlerFunc func(ctx context.Context, conn *websocket.Conn, payload json.RawMessage)

type Router struct {
	Routes map[string]HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		Routes: make(map[string]HandlerFunc),
	}
}

func (r *Router) Handle(action string, handler HandlerFunc) {
	r.Routes[action] = handler
}

func (r *Router) Dispatch(ctx context.Context, conn *websocket.Conn, msg wsutils.WSRequest) {
	handler, ok := r.Routes[msg.Action]
	if !ok {
		slog.Warn("Unknown WebSocket action", "action", msg.Action)

		wsmsg := wsutils.WSMessage{
			Action:  msg.Action,
			Success: false,
			Error:   "unknown action",
		}

		WriteJSON(ctx, conn, wsmsg)
		return
	}
	handler(ctx, conn, msg.Payload)
}

func WriteJSON(ctx context.Context, conn *websocket.Conn, v any) {

	data, err := json.Marshal(v)

	if err != nil {
		slog.Error("Failed to marshal websocket response", "error", err)
		return
	}

	if err = conn.Write(ctx, websocket.MessageText, data); err != nil {
		slog.Error("Failed to write websocket message", "error", err)
		return
	}

	slog.Debug("WebSocket message sent successfully", "data", string(data))
}
