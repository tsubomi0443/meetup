package interfaces

import (
	"fmt"
	"net/http"
	"net/url"

	"meetup/internal/usecases/dto"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

const (
	keyIsAdmin         = "isAdmin"
	keyIsManager       = "isManager"
	keyRoles           = "roles"
	keyCategories      = "categories"
	keySupportStatuses = "supportStatuses"
)

var masterLoginDefaults = map[string]any{
	keyIsAdmin:         true,
	keyIsManager:       false,
	keyRoles:           []dto.RoleForm{},
	keyCategories:      []dto.CategoryForm{},
	keySupportStatuses: []dto.SupportStatusForm{},
}

func (r *Router) setPageHandler() (routeInfos []echo.RouteInfo) {
	routeInfos = append(routeInfos, r.e.Static("/static", "static"))
	routeInfos = append(routeInfos, r.e.GET("/login", r.loginPage()))
	routeInfos = append(routeInfos, r.e.GET("/mock/:id", r.app(), GetJWTConfig()))
	return
}

func (r *Router) loginPage() echo.HandlerFunc {
	return func(c *echo.Context) error {
		viewData := map[string]any{
			keyIsAdmin:         masterLoginDefaults[keyIsAdmin],
			keyIsManager:       masterLoginDefaults[keyIsManager],
			keyRoles:           masterLoginDefaults[keyRoles],
			keyCategories:      masterLoginDefaults[keyCategories],
			keySupportStatuses: masterLoginDefaults[keySupportStatuses],
		}
		if ck, err := c.Request().Cookie(ErrorRedirectCookieName); err == nil && ck.Value != "" {
			if msg, err := url.QueryUnescape(ck.Value); err == nil {
				viewData["FlashError"] = msg
			} else {
				viewData["FlashError"] = ck.Value
			}
			c.SetCookie(&http.Cookie{
				Name:     ErrorRedirectCookieName,
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
			})
		}
		return c.Render(http.StatusOK, "login.html", viewData)
	}
}

func (r *Router) app() echo.HandlerFunc {
	return func(c *echo.Context) error {
		viewData := map[string]any{}

		loginErr := func(errMsg string) error {
			clearAccessTokenCookie(c)
			SetErrorFlashCookie(c, errMsg)
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		claims, err := r.getJWTCustomClaims(c)
		if err != nil {
			return loginErr("ログイン情報が見つかりません。\nログインからやり直してください。")
		}

		viewData[keyIsManager] = claims.RoleID == 2
		viewData[keyIsAdmin] = claims.RoleID == 1

		const dataLoadErrMsg = "データの読み込みに問題が発生しました。\n管理者へ問い合わせてください。"

		roleForms, err := r.deps.Master.ListRoles(c.Request().Context())
		if err != nil {
			return loginErr(dataLoadErrMsg)
		}
		viewData[keyRoles] = roleForms

		categoryForms, err := r.deps.Master.ListCategories(c.Request().Context())
		if err != nil {
			return loginErr(dataLoadErrMsg)
		}
		viewData[keyCategories] = categoryForms

		supportStatusForms, err := r.deps.Master.ListSupportStatuses(c.Request().Context())
		if err != nil {
			return loginErr(dataLoadErrMsg)
		}
		viewData[keySupportStatuses] = supportStatusForms

		return c.Render(http.StatusOK, "app.html", viewData)
	}
}

func (r *Router) getJWTCustomClaims(c *echo.Context) (claims *CustomClaims, err error) {
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
