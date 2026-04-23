package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	infrastructure "meetup/_mac_infrastructure"

	"gorm.io/gorm"
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
		Register:   make(chan *Client, 4),
		Unregister: make(chan *Client, 4),
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
				close(client.send)
				delete(h.clients, client)
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

func (hub *Hub) RunSSE(db *gorm.DB) error {
	ctx := context.Background()
	tickSecond := time.NewTicker(1 * time.Second)
	tickHalfSecond := time.NewTicker(500 * time.Millisecond)
	local, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return fmt.Errorf("time.LoadLocationの処理に失敗しました。: %w\n", err)
	}

	go func() {
		for t := range tickSecond.C {
			hub.Broadcast <- Event{Event: "time-tick", Data: t.In(local).Format(time.DateTime)}
		}
	}()

	go func() {
		for range tickHalfSecond.C {
			models, err := infrastructure.GetQuestions(ctx, db)
			if err != nil {
				hub.Broadcast <- Event{Event: "error", Data: fmt.Sprintf(`Error: %v\n`, err)}
				continue
			}
			for _, model := range models {
				qf := infrastructure.QuestionFromEntity(model)
				if data, err := json.Marshal(qf); err != nil {
					hub.Broadcast <- Event{Event: "error", Data: fmt.Sprintf(`Error: %v\n`, err)}
					continue
				} else {
					hub.Broadcast <- Event{Event: "question", Data: string(data)}
				}
			}
		}
	}()

	go func() {
		for range tickHalfSecond.C {
			models, err := infrastructure.GetUsers(ctx, db)
			if err != nil {
				hub.Broadcast <- Event{Event: "error", Data: fmt.Sprintf(`Error: %v\n`, err)}
				continue
			}
			for _, model := range models {
				uf := infrastructure.UserFromEntity(model)
				if data, err := json.Marshal(uf); err != nil {
					hub.Broadcast <- Event{Event: "error", Data: fmt.Sprintf(`Error: %v\n`, err)}
					continue
				} else {
					hub.Broadcast <- Event{Event: "user", Data: string(data)}
				}
			}
		}
	}()

	go func() {
		for range tickHalfSecond.C {
			models, err := infrastructure.GetTags(ctx, db)
			if err != nil {
				hub.Broadcast <- Event{Event: "error", Data: fmt.Sprintf(`Error: %v\n`, err)}
				continue
			}
			for _, model := range models {
				tf := infrastructure.TagFromEntity(model)
				if data, err := json.Marshal(tf); err != nil {
					hub.Broadcast <- Event{Event: "error", Data: fmt.Sprintf(`Error: %v\n`, err)}
					continue
				} else {
					hub.Broadcast <- Event{Event: "tag", Data: string(data)}
				}
			}
		}
	}()

	go func() {
		for range tickHalfSecond.C {
			models, err := infrastructure.GetNotice(ctx, db)
			if err != nil {
				hub.Broadcast <- Event{Event: "error", Data: fmt.Sprintf(`Error: %v\n`, err)}
				continue
			}
			for _, model := range models {
				nf := infrastructure.NoticeFromEntity(model)
				if data, err := json.Marshal(nf); err != nil {
					hub.Broadcast <- Event{Event: "error", Data: fmt.Sprintf(`Error: %v\n`, err)}
					continue
				} else {
					hub.Broadcast <- Event{Event: "notice", Data: string(data)}
				}
			}
		}
	}()

	<-ctx.Done()
	return nil
}
