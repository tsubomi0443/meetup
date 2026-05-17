package handler

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

var masterLoginDefaults = map[string]any{
	"isAdmin":   true,
	"isManager": false,
}

func (hm *HandlerManager) SetPageHandler() (routeInfos []echo.RouteInfo) {
	// 静的ファイル配信
	routeInfos = append(routeInfos, hm.e.Static("/static", "static"))
	routeInfos = append(routeInfos, hm.e.GET("/login", hm.loginPage()))
	routeInfos = append(routeInfos, hm.e.GET("/", hm.homePage()))
	routeInfos = append(routeInfos, hm.e.GET("/mock/:id", hm.mockup(), GetJWTConfig()))
	return
}

func (hm *HandlerManager) loginPage() echo.HandlerFunc {
	return func(c *echo.Context) error {
		viewData := map[string]any{
			"isAdmin":   masterLoginDefaults["isAdmin"],
			"isManager": masterLoginDefaults["isManager"],
		}
		if ck, err := c.Request().Cookie(ERROR_REDIRECT_COOKIE_NAME); err == nil && ck.Value != "" {
			if msg, err := url.QueryUnescape(ck.Value); err == nil {
				viewData["FlashError"] = msg
			} else {
				viewData["FlashError"] = ck.Value
			}
			c.SetCookie(&http.Cookie{
				Name:     ERROR_REDIRECT_COOKIE_NAME,
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
			})
		}
		return c.Render(http.StatusOK, "login.html", viewData)
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
		viewData := map[string]any{}
		if claims, err := hm.getJWTCustomClaims(c); err == nil {
			viewData["isManager"] = claims.RoleID == 2
			viewData["isAdmin"] = claims.RoleID == 1
		} else {
			clearAccessTokenCookie(c)
			SetErrorFlashCookie(c, "ログイン情報が見つかりません。\nログインからやり直してください。")
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		return c.Render(http.StatusOK, fmt.Sprintf("mock%s.html", id), viewData)
	}
}

func (hm *HandlerManager) getJWTCustomClaims(c *echo.Context) (claims *CustomClaims, err error) {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok || user == nil {
		return nil, fmt.Errorf("ユーザトークン不正")
	}
	claims, ok = user.Claims.(*CustomClaims)
	if !ok || claims == nil {
		return nil, fmt.Errorf("Claims取得・変換失敗")
	}
	return claims, nil
}
