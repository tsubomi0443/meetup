package handler

import (
	"encoding/json"
	"io"
	infrastructure "meetup/_mac_infrastructure"
	"meetup/env"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
)

const (
	COOKIE_NAME_TOKEN = "access_token"
)

func GetJWTConfig() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(env.GetJWTKey()),
		TokenLookup: "cookie:" + COOKIE_NAME_TOKEN,
		ErrorHandler: func(c *echo.Context, err error) error {
			return c.Redirect(http.StatusFound, "/login")
		},
	})
}

func (hm *HandlerManager) SetupAuthHandler() (routeInfos []echo.RouteInfo) {
	routeInfos = append(routeInfos, hm.e.POST("/login", hm.loginHandler()))
	return
}

func (hm *HandlerManager) loginHandler() echo.HandlerFunc {
	return func(c *echo.Context) error {
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}
		defer c.Request().Body.Close()

		var info struct {
			E string `json:"email"`
			P string `json:"pass"`
		}

		if err := json.Unmarshal(body, &info); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		info.P = encryptSha256(info.P)
		user, err := infrastructure.GetUserInfo(c.Request().Context(), hm.db, info.E, info.P)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}

		claims := jwt.MapClaims{
			"sub":   user.ID,
			"email": user.Email,
			"exp":   time.Now().Add(1 * time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signed, err := token.SignedString([]byte(env.GetJWTKey()))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		c.SetCookie(&http.Cookie{
			Name:     COOKIE_NAME_TOKEN,
			Value:    signed,
			Path:     "/",
			HttpOnly: true,
			Secure:   env.IsProduct(),
			SameSite: http.SameSiteLaxMode,
			Expires:  time.Now().Add(1 * time.Hour),
		})

		return c.JSON(http.StatusOK, map[string]string{
			"redirect": "/mock/5",
		})
	}
}

func logoutHandler() echo.HandlerFunc {
	return func(c *echo.Context) error {
		c.SetCookie(&http.Cookie{
			Name:     COOKIE_NAME_TOKEN,
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		})
		return c.Redirect(http.StatusSeeOther, "/login")
	}
}
