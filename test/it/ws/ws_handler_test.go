package ws

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	wsutils "user-management/internal/common/wsutils"

	ws "user-management/internal/ws"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/stretchr/testify/require"
)

func waitWithTimeout(t *testing.T, wg *sync.WaitGroup) {
	t.Helper()
	ch := make(chan struct{})
	go func() {
		wg.Wait()
		close(ch)
	}()
	select {
	case <-ch:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for goroutines")
	}
}

func TestHandleSubscribe_Valid(t *testing.T) {
	ctx := context.Background()
	connMgr := ws.NewConnectionManager(ctx)
	router := ws.NewRouter()
	h := ws.NewHandler(router, connMgr)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		websocket.Accept(w, r, nil)
	}))
	defer srv.Close()

	wsConn, _, err := websocket.Dial(ctx, "ws"+srv.URL[4:], nil)
	require.NoError(t, err)

	client := &ws.Client{
		Conn:   wsConn,
		SendCh: make(chan []byte, 10),
		Topics: make(map[string]bool),
	}
	connMgr.Register(client)

	payload := []byte(`{"topic":"BTCUSD"}`)
	h.HandleSubscribe(ctx, wsConn, payload)

	if !client.Topics["BTCUSD"] {
		t.Fatalf("expected client to be subscribed to BTCUSD")
	}
}

func TestHandleSubscribe_InvalidPayload(t *testing.T) {
	ctx := context.Background()
	connMgr := ws.NewConnectionManager(ctx)
	router := ws.NewRouter()
	h := ws.NewHandler(router, connMgr)

	router.Handle("subscribe", func(ctx context.Context, conn *websocket.Conn, payload json.RawMessage) {
		h.HandleSubscribe(ctx, conn, payload)
	})

	srv := httptest.NewServer(http.HandlerFunc(h.HandleWebSocket))
	defer srv.Close()

	wsURL := "ws" + srv.URL[len("http"):]
	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	require.NoError(t, err)

	var hello map[string]any
	require.NoError(t, wsjson.Read(ctx, conn, &hello))

	clientId, ok := hello["client_id"].(string)
	require.True(t, ok)

	client := connMgr.GetClientByID(clientId)
	if client == nil {
		t.Fatal("client not registered")
	}

	payload := []byte(`{"topic":"ABC"}`)

	h.HandleSubscribe(ctx, conn, payload)

	if len(client.Topics) > 0 {
		t.Fatalf("expected client to be subscribed to BTCUSD")
	}
}

func TestHandleUnsubscribe_Valid(t *testing.T) {
	ctx := context.Background()
	connMgr := ws.NewConnectionManager(ctx)
	router := ws.NewRouter()
	h := ws.NewHandler(router, connMgr)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		websocket.Accept(w, r, nil)
	}))
	defer srv.Close()

	wsConn, _, _ := websocket.Dial(ctx, "ws"+srv.URL[4:], nil)
	client := &ws.Client{
		Conn:   wsConn,
		SendCh: make(chan []byte, 10),
		Topics: map[string]bool{"BTCUSD": true},
	}
	connMgr.Register(client)

	payload := []byte(`{"topic":"BTCUSD"}`)
	h.HandleUnsubscribe(ctx, wsConn, payload)

	if client.Topics["BTCUSD"] {
		t.Fatal("expected BTCUSD to be removed from client's topics")
	}
}

func TestHandleWebSocket_SubscribeFlow(t *testing.T) {
	ctx := context.Background()
	connMgr := ws.NewConnectionManager(ctx)
	router := ws.NewRouter()
	h := ws.NewHandler(router, connMgr)

	router.Handle("subscribe", func(ctx context.Context, conn *websocket.Conn, payload json.RawMessage) {
		h.HandleSubscribe(ctx, conn, payload)
	})

	server := httptest.NewServer(http.HandlerFunc(h.HandleWebSocket))
	defer server.Close()

	wsURL := "ws" + server.URL[len("http"):]
	conn, _, err := websocket.Dial(ctx, wsURL, nil)

	var hello map[string]any
	require.NoError(t, wsjson.Read(ctx, conn, &hello))

	require.NoError(t, err)
	defer conn.Close(websocket.StatusNormalClosure, "")

	payload, _ := json.Marshal(map[string]any{"topic": "ETHUSD"})
	req := wsutils.WSRequest{Action: "subscribe", Payload: payload}

	require.NoError(t, wsjson.Write(ctx, conn, req))

	var resp wsutils.WSMessage
	require.NoError(t, wsjson.Read(ctx, conn, &resp))

	require.True(t, resp.Success)
	require.Equal(t, "ETHUSD", resp.Topic)
}

func TestHandleWebSocket_BroadcastFromManager(t *testing.T) {
	ctx := context.Background()
	connMgr := ws.NewConnectionManager(ctx)
	router := ws.NewRouter()
	h := ws.NewHandler(router, connMgr)

	server := httptest.NewServer(http.HandlerFunc(h.HandleWebSocket))
	defer server.Close()

	wsConn, _, err := websocket.Dial(ctx, "ws"+server.URL[4:], nil)
	require.NoError(t, err)

	time.Sleep(50 * time.Millisecond)

	var hello map[string]any
	require.NoError(t, wsjson.Read(ctx, wsConn, &hello))

	clientId, ok := hello["client_id"].(string)
	require.True(t, ok)

	client := connMgr.GetClientByID(clientId)
	if client == nil {
		t.Fatal("client not registered")
	}

	connMgr.Subscribe(client, "BTCUSD")

	msg := wsutils.WSMessage{
		Action: ws.PriceUpdate,
		Topic:  "BTCUSD",
		Data:   map[string]any{"price": 123.45},
	}

	connMgr.Broadcast(ctx, "BTCUSD", msg)

	_, data, err := wsConn.Read(ctx)
	require.NoError(t, err)

	var resp wsutils.WSMessage
	json.Unmarshal(data, &resp)

	if resp.Action != ws.PriceUpdate {
		t.Fatalf("expected action %s, got %s", ws.PriceUpdate, resp.Action)
	}
	if resp.Topic != "BTCUSD" {
		t.Fatalf("expected topic 'BTCUSD', got %s", resp.Topic)
	}
	d, ok := resp.Data.(map[string]any)
	if !ok || d["price"] != 123.45 {
		t.Fatalf("unexpected data: %+v", resp.Data)
	}
}
