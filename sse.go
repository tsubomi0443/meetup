package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
)

// Event はSSEで送信するイベントを表す。
type Event struct {
	ID    string
	Event string
	Data  string
}

// Client は個々のSSE接続を表す。
type Client struct {
	send chan Event
}

// Hub は全接続クライアントを管理し、メッセージのブロードキャストを担う。
type Hub struct {
	clients    map[*Client]struct{}
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Event
}

// NewHub はHubを初期化して返す。
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]struct{}),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Event),
	}
}

// Run はHubのイベントループ。goroutineとして起動する。
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.clients[client] = struct{}{}

		case client := <-h.Unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

		case event := <-h.Broadcast:
			for client := range h.clients {
				select {
				case client.send <- event:
				default:
					// 送信バッファが詰まっているクライアントは切断
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// sseHandler はSSEストリームを提供するEchoハンドラを返す。
func sseHandler(hub *Hub) echo.HandlerFunc {
	return func(c *echo.Context) error {
		client := &Client{send: make(chan Event, 8)}
		hub.Register <- client
		defer func() {
			hub.Unregister <- client
		}()

		w := c.Response()
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(http.StatusOK)

		flusher, ok := w.(http.Flusher)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "streaming not supported")
		}

		ctx := c.Request().Context()
		for {
			select {
			case <-ctx.Done():
				return nil
			case event, ok := <-client.send:
				if !ok {
					return nil
				}
				if event.ID != "" {
					fmt.Fprintf(w, "id: %s\n", event.ID)
				}
				if event.Event != "" {
					fmt.Fprintf(w, "event: %s\n", event.Event)
				}
				fmt.Fprintf(w, "data: %s\n\n", event.Data)
				flusher.Flush()
			}
		}
	}
}
