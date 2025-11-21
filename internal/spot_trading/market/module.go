package market

import binance "user-management/internal/integration/binance"

type Module struct {
	Handler *Handler
}

func NewModule(streamer *binance.Streamer) *Module {
	handler := NewHandler(streamer)

	return &Module{
		Handler: handler,
	}
}
