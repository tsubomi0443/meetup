package sse

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"meetup/internal/infrastructures/crypto"
)

const (
	SSETimeTick = "time-tick"
	SSENotice   = "notice"
	SSEError    = "error"
	SSEGet      = "get"
	SSECreate   = "create"
	SSEUpdate   = "update"
	SSEDelete   = "delete"
)

type Event struct {
	ID    string
	Event string
	Data  string
}

type Client struct {
	Send chan Event
}

type Hub struct {
	clients    map[*Client]struct{}
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Event
	logger     func(ctx context.Context, level slog.Level, msg string, args ...any)
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]struct{}),
		Register:   make(chan *Client, 4),
		Unregister: make(chan *Client, 4),
		Broadcast:  make(chan Event),
	}
}

func (h *Hub) SetLogger(logger func(ctx context.Context, level slog.Level, msg string, args ...any)) {
	h.logger = logger
}

func (h *Hub) Send(data, event string) {
	h.Broadcast <- Event{Event: event, Data: data}
}

func (h *Hub) Sends(data string, events ...string) {
	for _, event := range events {
		h.Send(data, event)
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.clients[client] = struct{}{}
		case client := <-h.Unregister:
			if _, ok := h.clients[client]; ok {
				close(client.Send)
				delete(h.clients, client)
			}
		case event := <-h.Broadcast:
			for client := range h.clients {
				select {
				case client.Send <- event:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) RunSSE() {
	go func() {
		tickSecond := time.NewTicker(1 * time.Second)
		local, err := time.LoadLocation("Asia/Tokyo")
		if err != nil {
			return
		}
		for t := range tickSecond.C {
			h.sendTimeTicker(t.In(local).Format(time.DateTime))
		}
	}()
}

func (h *Hub) sendTimeTicker(data string) {
	h.Send(data, SSETimeTick)
}

func (h *Hub) SendError(api, data string) {
	event := fmt.Sprintf("%s-%s", SSEError, api)
	h.logger(context.Background(), slog.LevelError, event, "data", crypto.EncryptSHA256(data))
	h.Send(data, event)
}

func (h *Hub) SendGetEvent(api, data string) {
	event := fmt.Sprintf("%s-%s", SSEGet, api)
	h.logger(context.Background(), slog.LevelInfo, event, "data", crypto.EncryptSHA256(data))
	h.Send(data, event)
}

func (h *Hub) SendCreateEvent(api, data string) {
	event := fmt.Sprintf("%s-%s", SSECreate, api)
	h.logger(context.Background(), slog.LevelInfo, event, "data", crypto.EncryptSHA256(data))
	h.Send(data, event)
}

func (h *Hub) SendUpdateEvent(api, data string) {
	event := fmt.Sprintf("%s-%s", SSEUpdate, api)
	h.logger(context.Background(), slog.LevelInfo, event, "data", crypto.EncryptSHA256(data))
	h.Send(data, event)
}

func (h *Hub) SendDeleteEvent(api, data string) {
	event := fmt.Sprintf("%s-%s", SSEDelete, api)
	h.logger(context.Background(), slog.LevelInfo, event, "data", crypto.EncryptSHA256(data))
	h.Send(data, event)
}

func (h *Hub) SendNotice(api, data string) {
	event := fmt.Sprintf("%s-%s", SSENotice, api)
	h.logger(context.Background(), slog.LevelInfo, event, "data", crypto.EncryptSHA256(data))
	h.Send(data, event)
}
