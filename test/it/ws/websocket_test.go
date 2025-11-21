package ws

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-management/internal/config"
	"user-management/internal/db"
	"user-management/internal/db/sqlc"
	"user-management/internal/events"
	"user-management/internal/instrument"
	ws "user-management/internal/ws"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebSocket_SubscribeIntegration(t *testing.T) {
	ctx := context.Background()

	connMgr := ws.NewConnectionManager(ctx)
	router := ws.NewRouter()
	handler := ws.NewHandler(router, connMgr)

	router.Handle("subscribe", func(ctx context.Context, conn *websocket.Conn, payload json.RawMessage) {
		var p struct {
			Topic string `json:"topic"`
		}
		_ = json.Unmarshal(payload, &p)

		client := connMgr.GetClient(conn)
		assert.NotNil(t, client, "client should be registered")

		connMgr.Subscribe(client, p.Topic)

		_ = wsjson.Write(ctx, conn, map[string]any{
			"success": true,
			"topic":   p.Topic,
		})
	})

	server := httptest.NewServer(http.HandlerFunc(handler.HandleWebSocket))
	defer server.Close()

	wsURL := "ws" + server.URL[len("http"):]
	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	require.NoError(t, err)
	defer conn.Close(websocket.StatusNormalClosure, "")

	var hello map[string]any
	require.NoError(t, wsjson.Read(ctx, conn, &hello))

	clientId, ok := hello["client_id"].(string)
	require.True(t, ok)

	payload, _ := json.Marshal(map[string]any{"topic": "AAPL"})
	req := ws.WSRequest{Method: "subscribe", Params: payload}

	require.NoError(t, wsjson.Write(ctx, conn, req))

	var resp map[string]any
	require.NoError(t, wsjson.Read(ctx, conn, &resp))

	assert.Equal(t, "AAPL", resp["topic"])

	client := connMgr.GetClientByID(clientId)
	require.NotNil(t, client)

	assert.True(t, client.Topics["AAPL"])
}

func TestWebSocket_UnsubscribeIntegration(t *testing.T) {
	ctx := context.Background()

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

	router.Handle("unsubscribe", func(ctx context.Context, conn *websocket.Conn, payload json.RawMessage) {
		var p struct {
			Topic string `json:"topic"`
		}
		_ = json.Unmarshal(payload, &p)

		client := connMgr.GetClient(conn)
		require.NotNil(t, client)

		connMgr.Unsubscribe(client, p.Topic)

		_ = wsjson.Write(ctx, conn, map[string]any{
			"action":  "unsubscribe",
			"success": true,
			"topic":   p.Topic,
		})
	})

	server := httptest.NewServer(http.HandlerFunc(handler.HandleWebSocket))
	defer server.Close()

	wsURL := "ws" + server.URL[len("http"):]
	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	require.NoError(t, err)
	defer conn.Close(websocket.StatusNormalClosure, "")

	var hello map[string]any
	require.NoError(t, wsjson.Read(ctx, conn, &hello))

	clientId, ok := hello["client_id"].(string)
	require.True(t, ok)

	subPayload, _ := json.Marshal(map[string]any{"topic": "AAPL"})
	require.NoError(t, wsjson.Write(ctx, conn, ws.WSRequest{Method: "subscribe", Params: subPayload}))

	var subResp map[string]any
	require.NoError(t, wsjson.Read(ctx, conn, &subResp))
	assert.Equal(t, "AAPL", subResp["topic"])

	client := connMgr.GetClientByID(clientId)
	require.NotNil(t, client)
	assert.True(t, client.Topics["AAPL"])

	unsubPayload, _ := json.Marshal(map[string]any{"topic": "AAPL"})
	require.NoError(t, wsjson.Write(ctx, conn, ws.WSRequest{Method: "unsubscribe", Params: unsubPayload}))

	var unsubResp map[string]any
	require.NoError(t, wsjson.Read(ctx, conn, &unsubResp))
	assert.Equal(t, "AAPL", unsubResp["topic"])

	assert.False(t, client.Topics["AAPL"])
}

func TestWebSocket_ReceivePriceUpdate(t *testing.T) {
	ctx := context.Background()

	config := config.GetConfig()
	dbConn := db.Connect(config.DBDsn)
	defer dbConn.Close()

	queries := sqlc.New(dbConn)
	eventBus := events.NewBus()
	connMgr := ws.NewConnectionManager(ctx)
	router := ws.NewRouter()
	handler := ws.NewHandler(router, connMgr)

	repo := instrument.NewRepository(queries)
	service := instrument.NewService(repo, eventBus)

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

	eventBus.Subscribe(events.InstrumentUpdated, func(ev events.Event) {
		msg := ev.Data.(events.InstrumentUpdatedEvent)

		wsMsg := map[string]any{
			"action": ws.PriceUpdate,
			"topic":  events.InstrumentUpdated,
			"data":   msg,
		}

		connMgr.Broadcast(ctx, "AAPL", wsMsg)
	})

	server := httptest.NewServer(http.HandlerFunc(handler.HandleWebSocket))
	defer server.Close()

	wsURL := "ws" + server.URL[len("http"):]
	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	require.NoError(t, err)
	defer conn.Close(websocket.StatusNormalClosure, "")

	var hello map[string]any
	require.NoError(t, wsjson.Read(ctx, conn, &hello))

	subPayload, _ := json.Marshal(map[string]any{"topic": "AAPL"})
	require.NoError(t, wsjson.Write(ctx, conn, ws.WSRequest{
		Method: "subscribe",
		Params: subPayload,
	}))

	var subResp map[string]any
	require.NoError(t, wsjson.Read(ctx, conn, &subResp))
	assert.Equal(t, "AAPL", subResp["topic"])

	updateRequest := instrument.InstrumentUpdateRequest{
		LastPrice: 190.0105,
	}

	_, err = service.UpdateInstrument(ctx,
		"ef4837d8-49fb-43f7-827c-133706388119",
		&updateRequest,
		connMgr)
	require.NoError(t, err)

	var updateMsg map[string]any
	require.NoError(t, wsjson.Read(ctx, conn, &updateMsg))

	assert.Equal(t, ws.PriceUpdate, updateMsg["action"])
	assert.Equal(t, events.InstrumentUpdated, updateMsg["topic"])

	data := updateMsg["data"].(map[string]any)

	assert.Equal(t, "AAPL", data["symbol"])
	assert.Equal(t, "190.010500", data["price"])
	assert.NotEmpty(t, data["updated_at"])
}
