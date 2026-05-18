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

// Deps bundles use cases and infrastructure used by HTTP handlers.
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

type Router struct {
	e    *echo.Echo
	deps Deps
}

func NewRouter(e *echo.Echo, deps Deps) *Router {
	return &Router{e: e, deps: deps}
}

func (r *Router) SetHubLogger(logger func(ctx context.Context, level slog.Level, msg string, args ...any)) {
	r.deps.Hub.SetLogger(logger)
}

func (r *Router) SetupHandlers() (routeInfos []echo.RouteInfo) {
	routeInfos = append(routeInfos, r.setupAuthHandler()...)
	routeInfos = append(routeInfos, r.setPageHandler()...)
	routeInfos = append(routeInfos, r.setSSEHandler()...)
	routeInfos = append(routeInfos, r.setAPIHandler()...)
	return
}

func (r *Router) PollingStart(ctx context.Context) error {
	return r.deps.NoticePoller.Run(ctx)
}

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
