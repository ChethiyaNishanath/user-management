package ws

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	ws "user-management/internal/ws"

	"github.com/coder/websocket"
	"github.com/stretchr/testify/require"
)

func newMockWSServer(t *testing.T) *httptest.Server {
	t.Helper()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		require.NoError(t, err)

		go func() {
			ctx := context.Background()

			for {
				_, data, err := conn.Read(ctx)
				if err != nil {
					return
				}

				var msg map[string]any
				_ = json.Unmarshal(data, &msg)

				if msg["method"] == "SUBSCRIBE" {
					resp := map[string]any{
						"result": nil,
						"id":     msg["id"],
					}
					b, _ := json.Marshal(resp)
					_ = conn.Write(ctx, websocket.MessageText, b)

					time.Sleep(200 * time.Millisecond)

					update := map[string]any{
						"e": "depthUpdate",
						"E": time.Now().UnixMilli(),
						"U": 1,
						"u": 2,
						"b": [][]string{{"50000.1", "1.23"}},
						"a": [][]string{{"50010.5", "0.99"}},
					}
					ub, _ := json.Marshal(update)
					_ = conn.Write(ctx, websocket.MessageText, ub)
				}
			}
		}()
	})

	return httptest.NewServer(handler)
}

func TestWSClient(t *testing.T) {
	server := newMockWSServer(t)
	defer server.Close()

	wsURL := "ws" + server.URL[len("http"):]

	wsClient := ws.NewWSClient(wsURL)

	received := make(chan []byte, 10)

	wsClient.OnMessage = func(msgType websocket.MessageType, data []byte) {
		slog.Info("received data", "data", string(data))
		received <- data
	}

	require.NoError(t, wsClient.Connect())

	sub := map[string]any{
		"method": "SUBSCRIBE",
		"params": []string{"btcusdt@depth"},
		"id":     1,
	}

	require.NoError(t, wsClient.SendJSON(sub))

	select {
	case msg := <-received:
		require.Contains(t, string(msg), `"result":null`)
	case <-time.After(2 * time.Second):
		t.Fatal("did not receive subscription response")
	}

	select {
	case msg := <-received:
		require.Contains(t, string(msg), `"depthUpdate"`)
	case <-time.After(2 * time.Second):
		t.Fatal("did not receive depth update")
	}
}
