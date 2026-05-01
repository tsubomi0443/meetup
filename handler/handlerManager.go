package handler

import (
	"context"
	"log/slog"

	"meetup/env"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type HandlerManager struct {
	db  *gorm.DB
	e   *echo.Echo
	hub *Hub
	ne  *NoticeEvent
}

func NewHandlerManager(db *gorm.DB, e *echo.Echo, hub *Hub, ne *NoticeEvent) *HandlerManager {
	return &HandlerManager{
		db:  db,
		e:   e,
		hub: hub,
		ne:  ne,
	}
}

func (hm *HandlerManager) SetupHandlers() (routeInfos []echo.RouteInfo) {
	routeInfos = append(routeInfos, hm.SetupAuthHandler()...)
	routeInfos = append(routeInfos, hm.SetPageHandler()...)
	routeInfos = append(routeInfos, hm.SetSSEHandler()...)
	routeInfos = append(routeInfos, hm.SetAPIHandler()...)
	return
}

func (hm *HandlerManager) Logging(ctx context.Context, level slog.Level, msg string, args ...any) {
	switch level {
	case slog.LevelDebug:
		if env.IsProduct() {
			hm.e.Logger.DebugContext(ctx, msg, args...)
			break
		}
		hm.e.Logger.Debug(msg, args...)
	case slog.LevelInfo:
		if env.IsProduct() {
			hm.e.Logger.DebugContext(ctx, msg, args...)
			break
		}
		hm.e.Logger.Info(msg, args...)
	case slog.LevelWarn:
		if env.IsProduct() {
			hm.e.Logger.DebugContext(ctx, msg, args...)
			break
		}
		hm.e.Logger.Warn(msg, args...)
	case slog.LevelError:
		if env.IsProduct() {
			hm.e.Logger.DebugContext(ctx, msg, args...)
			break
		}
		hm.e.Logger.Error(msg, args...)
	}
}
