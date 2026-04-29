package handler

import (
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type HandlerManager struct {
	db *gorm.DB
	e *echo.Echo
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
