package market

import (
	"encoding/json"
	"fmt"
	"net/http"
	httputils "user-management/internal/common/httputils"
	stringutils "user-management/internal/common/stringutils"
	binance "user-management/internal/integration/binance"
)

type Handler struct {
	streamer *binance.Streamer
}

func NewHandler(streamer *binance.Streamer) *Handler {
	return &Handler{streamer: streamer}
}

func (h *Handler) FetchOrderBook(w http.ResponseWriter, r *http.Request) {
	exchange := r.URL.Query().Get("exchange")
	symbol := r.URL.Query().Get("symbol")

	w.Header().Set("Content-Type", "application/json")

	if stringutils.Equals(exchange, "binance") {

		data := h.streamer.GetOrderBook(symbol)
		if data == nil {
			httputils.WriteError(w, http.StatusServiceUnavailable, fmt.Sprintf("order book not ready symbol=%s", symbol), r)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
		return
	}

	http.Error(w, "unsupported exchange: "+exchange, http.StatusBadRequest)
}
