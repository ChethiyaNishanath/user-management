package ws

import (
	"context"
	"encoding/json"
	"log/slog"

	wsutils "user-management/internal/common/wsutils"

	"github.com/coder/websocket"
)

func (h *Handler) HandleSubscribe(ctx context.Context, conn *websocket.Conn, payload json.RawMessage) {

	var data struct {
		Topic string `json:"topic"`
	}

	if err := json.Unmarshal(payload, &data); err != nil || data.Topic == "" {

		msg := wsutils.WSMessage{
			Action:  "subscribe",
			Success: false,
			Error:   "Invalid payload or missing topic",
		}

		WriteJSON(ctx, conn, msg)
		return
	}

	client := h.connMgr.GetClient(conn)
	if client == nil {
		slog.Warn("Subscribe request from unknown client")
		return
	}

	h.connMgr.Subscribe(client, data.Topic)

	msg := wsutils.WSMessage{
		Action:  "subscribe",
		Success: true,
		Topic:   data.Topic,
	}

	WriteJSON(ctx, conn, msg)
}

func (h *Handler) HandleUnsubscribe(ctx context.Context, conn *websocket.Conn, payload json.RawMessage) {
	var data struct {
		Topic string `json:"topic"`
	}

	if err := json.Unmarshal(payload, &data); err != nil || data.Topic == "" {
		msg := wsutils.WSMessage{
			Action:  "unsubscribe",
			Success: false,
			Error:   "invalid payload or missing topic",
		}
		WriteJSON(ctx, conn, msg)
		return
	}

	client := h.connMgr.GetClient(conn)
	if client == nil {
		slog.Warn("Unsubscribe request from unknown client")
		return
	}

	h.connMgr.Unsubscribe(client, data.Topic)

	msg := wsutils.WSMessage{
		Action:  "unsubscribe",
		Success: true,
		Topic:   data.Topic,
	}

	WriteJSON(ctx, conn, msg)
}
