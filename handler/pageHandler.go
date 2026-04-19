package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
)

func (hm *HandlerManager) SetPageHandler() (routeInfos []echo.RouteInfo) {
	// 静的ファイル配信
	routeInfos = append(routeInfos, hm.e.Static("/static", "static"))
	routeInfos = append(routeInfos, hm.e.GET("/login", hm.loginPage()))
	routeInfos = append(routeInfos, hm.e.GET("/", hm.homePage(), GetJWTConfig()))
	routeInfos = append(routeInfos, hm.e.GET("/mock/:id", hm.mockup()))
	return
}

func (hm *HandlerManager) loginPage() echo.HandlerFunc {
	return func(c *echo.Context) error {
		return c.Render(http.StatusOK, "login.html", nil)
	}
}

func (hm *HandlerManager) homePage() echo.HandlerFunc {
	return func(c *echo.Context) error {
		return c.Render(http.StatusOK, "index.html", map[string]any{
			"Title": "Go + Echo + HTMX + Alpine.js",
		})
	}
}

func (hm *HandlerManager) mockup() echo.HandlerFunc {
	return func(c *echo.Context) error {
		id := c.Param("id")
		return c.Render(http.StatusOK, fmt.Sprintf("mock%s.html", id), nil)
	}
}
