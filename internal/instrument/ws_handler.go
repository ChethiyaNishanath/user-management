package instrument

import (
	"context"
	"encoding/json"
	"log/slog"
	wsutils "user-management/internal/common/wsutils"
	ws "user-management/internal/ws"

	"github.com/coder/websocket"
)

type GetInstumentPayload struct {
	Symbol string `json:"symbol"`
}

func RegisterWsRoutes(router *ws.Router, instrumentService *Service) {
	router.Handle("get_instrument", getInstrument(instrumentService))
}

func getInstrument(instrumentService *Service) ws.HandlerFunc {
	return func(ctx context.Context, conn *websocket.Conn, payload json.RawMessage) {
		var p GetInstumentPayload
		if err := json.Unmarshal(payload, &p); err != nil {
			slog.Warn("Invalid get_instrument payload", "warning", err)

			msg := wsutils.WSMessage{
				Action:  "get_instrument_response",
				Success: false,
				Error:   "Invalid payload",
			}

			ws.WriteJSON(ctx, conn, msg)
			return
		}

		u, err := instrumentService.GetInstrumentBySymbol(ctx, p.Symbol)
		if err != nil {

			msg := wsutils.WSMessage{
				Action:  "get_instrument_response",
				Success: false,
				Error:   "Instrument not found",
			}

			ws.WriteJSON(ctx, conn, msg)
			return
		}

		msg := wsutils.WSMessage{
			Action: "get_instrument_response",
			Data:   u,
		}

		ws.WriteJSON(ctx, conn, msg)
	}
}
