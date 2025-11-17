package ws

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/coder/websocket"
)

type WSClient struct {
	URL    string
	Conn   *websocket.Conn
	mu     sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc

	OnMessage func(msgType websocket.MessageType, data []byte)

	Reconnect    bool
	RetryDelay   time.Duration
	PingInterval time.Duration
}

func NewWSClient(url string) *WSClient {
	ctx, cancel := context.WithCancel(context.Background())
	return &WSClient{
		URL:          url,
		ctx:          ctx,
		cancel:       cancel,
		Reconnect:    true,
		RetryDelay:   2 * time.Second,
		PingInterval: 20 * time.Second,
	}
}

func (c *WSClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var err error
	c.Conn, _, err = websocket.Dial(c.ctx, c.URL, nil)
	if err != nil {
		return err
	}

	go c.readLoop()
	go c.pingLoop()

	return nil
}

func (c *WSClient) readLoop() {
	for {
		msgType, data, err := c.Conn.Read(c.ctx)
		if err != nil {
			slog.Error("WS read error", "error", err)
			c.reconnect()
			return
		}

		if c.OnMessage != nil {
			c.OnMessage(msgType, data)
		}
	}
}

func (c *WSClient) reconnect() {
	if !c.Reconnect {
		return
	}

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		time.Sleep(c.RetryDelay)
		slog.Info("Attempting reconnect...")

		err := c.Connect()
		if err == nil {
			slog.Info("Reconnected")
			return
		}
	}
}

func (c *WSClient) pingLoop() {
	ticker := time.NewTicker(c.PingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			if c.Conn != nil {
				_ = c.Conn.Ping(c.ctx)
			}
			c.mu.Unlock()

		case <-c.ctx.Done():
			return
		}
	}
}

func (c *WSClient) SendJSON(v any) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return c.Conn.Write(c.ctx, websocket.MessageText, data)
}

func (c *WSClient) SendRaw(msgType websocket.MessageType, data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.Conn.Write(c.ctx, msgType, data)
}
