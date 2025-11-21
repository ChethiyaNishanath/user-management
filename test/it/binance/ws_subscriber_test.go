package binance

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"user-management/internal/config"
	"user-management/internal/db"
	"user-management/internal/events"
	binance "user-management/internal/integration/binance"
	"user-management/internal/ws"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebSocket_ReceiveBinancePriceUpdates(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	config.Init()
	cfg := config.GetConfig()
	dbConn := db.Connect(cfg.DBDsn)
	defer dbConn.Close()

	bus := events.NewBus()
	connMgr := ws.NewConnectionManager(ctx)
	router := ws.NewRouter()
	handler := ws.NewHandler(router, connMgr)

	router.Handle("subscribe", func(ctx context.Context, conn *websocket.Conn, payload json.RawMessage) {
		var p struct {
			Topic string `json:"topic"`
		}
		_ = json.Unmarshal(payload, &p)

		client := connMgr.GetClient(conn)
		require.NotNil(t, client)

		connMgr.Subscribe(client, p.Topic)

		_ = wsjson.Write(ctx, conn, map[string]any{
			"action":  "subscribe",
			"success": true,
			"topic":   p.Topic,
		})
	})

	bus.Subscribe("btcusdt@depth", func(e events.Event) {
		evt := e.Data.(binance.DepthStreamEvent)
		connMgr.Broadcast(
			context.Background(),
			e.Topic,
			ws.WSMessage{
				Data: evt,
			},
		)
	})

	streamer := binance.NewStreamer(ctx, bus)
	streamer.Start(ctx)

	server := httptest.NewServer(http.HandlerFunc(handler.HandleWebSocket))
	defer server.Close()

	wsURL := "ws" + server.URL[len("http"):]
	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	require.NoError(t, err)
	defer conn.Close(websocket.StatusNormalClosure, "")

	var hello map[string]any
	require.NoError(t, wsjson.Read(ctx, conn, &hello))

	_, ok := hello["client_id"].(string)
	require.True(t, ok)

	subPayload, _ := json.Marshal(map[string]any{"topic": "btcusdt@depth"})
	require.NoError(t, wsjson.Write(ctx, conn, ws.WSRequest{
		Method: "subscribe",
		Params: subPayload,
	}))

	var subResp map[string]any
	require.NoError(t, wsjson.Read(ctx, conn, &subResp))
	assert.Equal(t, "btcusdt@depth", subResp["topic"])
}
