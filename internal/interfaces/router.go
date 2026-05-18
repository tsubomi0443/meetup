package interfaces

import (
	"context"
	"log/slog"

	authuc "meetup/internal/usecases/auth"
	masteruc "meetup/internal/usecases/master"
	noticeuc "meetup/internal/usecases/notice"
	questionuc "meetup/internal/usecases/question"
	taguc "meetup/internal/usecases/tag"
	useruc "meetup/internal/usecases/user"

	"meetup/internal/infrastructures/config"
	"meetup/internal/interfaces/sse"

	"github.com/labstack/echo/v5"
)

// Deps は HTTP ハンドラが利用するユースケースとインフラストラクチャをまとめる。
type Deps struct {
	Hub          *sse.Hub
	NoticeEvents *noticeuc.Event
	NoticePoller *noticeuc.Poller
	Auth         *authuc.UseCase
	User         *useruc.UseCase
	Question     *questionuc.UseCase
	Tag          *taguc.UseCase
	Notice       *noticeuc.UseCase
	Master       *masteruc.UseCase
}

// Router は Echo ルート登録とハンドラ設定を担う。
type Router struct {
	e    *echo.Echo
	deps Deps
}

// NewRouter はルーターを生成する。
//
// args:
//   - e *echo.Echo: HTTP サーバ
//   - deps Deps: ハンドラ用依存関係
//
// return:
//   - *Router: 生成したルーター
func NewRouter(e *echo.Echo, deps Deps) *Router {
	return &Router{e: e, deps: deps}
}

// SetHubLogger は SSE ハブのロガーを設定する。
//
// args:
//   - logger func(ctx context.Context, level slog.Level, msg string, args ...any): 構造化ログ出力関数
func (r *Router) SetHubLogger(logger func(ctx context.Context, level slog.Level, msg string, args ...any)) {
	r.deps.Hub.SetLogger(logger)
}

// SetupHandlers は認証・ページ・SSE・API の各ハンドラを登録する。
//
// return:
//   - []echo.RouteInfo: 登録したルート情報の一覧
func (r *Router) SetupHandlers() (routeInfos []echo.RouteInfo) {
	routeInfos = append(routeInfos, r.setupAuthHandler()...)
	routeInfos = append(routeInfos, r.setPageHandler()...)
	routeInfos = append(routeInfos, r.setSSEHandler()...)
	routeInfos = append(routeInfos, r.setAPIHandler()...)
	return
}

// PollingStart は通知ポーラーのバックグラウンド処理を開始する。
//
// args:
//   - ctx context.Context: キャンセル用コンテキスト
//
// return:
//   - error: ポーラー実行中のエラー
func (r *Router) PollingStart(ctx context.Context) error {
	return r.deps.NoticePoller.Run(ctx)
}

// Logging は環境に応じて Echo ロガーへメッセージを出力する。
//
// args:
//   - ctx context.Context: 本番時のコンテキスト付きログ用
//   - level slog.Level: ログレベル
//   - msg string: メッセージ
//   - args ...any: 構造化ログの追加フィールド
func (r *Router) Logging(ctx context.Context, level slog.Level, msg string, args ...any) {
	switch level {
	case slog.LevelDebug:
		if config.IsProduct() {
			r.e.Logger.DebugContext(ctx, msg, args...)
			break
		}
		r.e.Logger.Debug(msg, args...)
	case slog.LevelInfo:
		if config.IsProduct() {
			r.e.Logger.DebugContext(ctx, msg, args...)
			break
		}
		r.e.Logger.Info(msg, args...)
	case slog.LevelWarn:
		if config.IsProduct() {
			r.e.Logger.DebugContext(ctx, msg, args...)
			break
		}
		r.e.Logger.Warn(msg, args...)
	case slog.LevelError:
		if config.IsProduct() {
			r.e.Logger.DebugContext(ctx, msg, args...)
			break
		}
		r.e.Logger.Error(msg, args...)
	}
}
