package sse

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"meetup/internal/infrastructures/crypto"
)

const (
	// SSETimeTick は時刻ティック用 SSE イベント名。
	SSETimeTick = "time-tick"
	// SSENotice は通知用 SSE イベント名プレフィックス。
	SSENotice = "notice"
	// SSEError はエラー用 SSE イベント名プレフィックス。
	SSEError = "error"
	// SSEGet は取得系 SSE イベント名プレフィックス。
	SSEGet = "get"
	// SSECreate は作成系 SSE イベント名プレフィックス。
	SSECreate = "create"
	// SSEUpdate は更新系 SSE イベント名プレフィックス。
	SSEUpdate = "update"
	// SSEDelete は削除系 SSE イベント名プレフィックス。
	SSEDelete = "delete"
)

// Event は SSE で配信する1件のイベント。
type Event struct {
	ID    string
	Event string
	Data  string
}

// Client は SSE 接続1件分の送信チャネルを保持する。
type Client struct {
	Send chan Event
}

// Hub は SSE クライアントの登録・配信を管理する。
type Hub struct {
	clients    map[*Client]struct{}
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Event
	logger     func(ctx context.Context, level slog.Level, msg string, args ...any)
}

// NewHub は SSE ハブを生成する。
//
// return:
//   - *Hub: 初期化済みハブ
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]struct{}),
		Register:   make(chan *Client, 4),
		Unregister: make(chan *Client, 4),
		Broadcast:  make(chan Event),
	}
}

// SetLogger は構造化ログ出力関数を設定する。
//
// args:
//   - logger func(ctx context.Context, level slog.Level, msg string, args ...any): ログ関数
func (h *Hub) SetLogger(logger func(ctx context.Context, level slog.Level, msg string, args ...any)) {
	h.logger = logger
}

// Send は全接続クライアントへイベントをブロードキャストする。
//
// args:
//   - data string: イベント本文（data フィールド）
//   - event string: イベント種別名
func (h *Hub) Send(data, event string) {
	h.Broadcast <- Event{Event: event, Data: data}
}

// Sends は複数のイベント種別名に対して同一 data を順に配信する。
//
// args:
//   - data string: イベント本文
//   - events ...string: イベント種別名の列
func (h *Hub) Sends(data string, events ...string) {
	for _, event := range events {
		h.Send(data, event)
	}
}

// Run はクライアント登録・解除とブロードキャストのメインループを実行する（ブロッキング）。
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

// RunSSE は1秒間隔の時刻ティックを配信するゴルーチンを起動する。
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

// sendTimeTicker は時刻文字列を time-tick イベントとして配信する。
//
// args:
//   - data string: フォーマット済み日時文字列
func (h *Hub) sendTimeTicker(data string) {
	h.Send(data, SSETimeTick)
}

// SendError は API 向けエラー SSE イベントをログ出力して配信する。
//
// args:
//   - api string: API 識別子（イベント名に付与）
//   - data string: ペイロード（ログは SHA256 ハッシュ化）
func (h *Hub) SendError(api, data string) {
	event := fmt.Sprintf("%s-%s", SSEError, api)
	h.logger(context.Background(), slog.LevelError, event, "data", crypto.EncryptSHA256(data))
	h.Send(data, event)
}

// SendGetEvent は取得系 SSE イベントをログ出力して配信する。
//
// args:
//   - api string: API 識別子
//   - data string: ペイロード（ログは SHA256 ハッシュ化）
func (h *Hub) SendGetEvent(api, data string) {
	event := fmt.Sprintf("%s-%s", SSEGet, api)
	h.logger(context.Background(), slog.LevelInfo, event, "data", crypto.EncryptSHA256(data))
	h.Send(data, event)
}

// SendCreateEvent は作成系 SSE イベントをログ出力して配信する。
//
// args:
//   - api string: API 識別子
//   - data string: ペイロード（ログは SHA256 ハッシュ化）
func (h *Hub) SendCreateEvent(api, data string) {
	event := fmt.Sprintf("%s-%s", SSECreate, api)
	h.logger(context.Background(), slog.LevelInfo, event, "data", crypto.EncryptSHA256(data))
	h.Send(data, event)
}

// SendUpdateEvent は更新系 SSE イベントをログ出力して配信する。
//
// args:
//   - api string: API 識別子
//   - data string: ペイロード（ログは SHA256 ハッシュ化）
func (h *Hub) SendUpdateEvent(api, data string) {
	event := fmt.Sprintf("%s-%s", SSEUpdate, api)
	h.logger(context.Background(), slog.LevelInfo, event, "data", crypto.EncryptSHA256(data))
	h.Send(data, event)
}

// SendDeleteEvent は削除系 SSE イベントをログ出力して配信する。
//
// args:
//   - api string: API 識別子
//   - data string: ペイロード（ログは SHA256 ハッシュ化）
func (h *Hub) SendDeleteEvent(api, data string) {
	event := fmt.Sprintf("%s-%s", SSEDelete, api)
	h.logger(context.Background(), slog.LevelInfo, event, "data", crypto.EncryptSHA256(data))
	h.Send(data, event)
}

// SendNotice は通知 SSE イベントをログ出力して配信する。
//
// args:
//   - api string: API 識別子
//   - data string: ペイロード（ログは SHA256 ハッシュ化）
func (h *Hub) SendNotice(api, data string) {
	event := fmt.Sprintf("%s-%s", SSENotice, api)
	h.logger(context.Background(), slog.LevelInfo, event, "data", crypto.EncryptSHA256(data))
	h.Send(data, event)
}
