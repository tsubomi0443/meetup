package handler

import (
	"encoding/json"
	"errors"
	"io"
	infrastructure "meetup/_mac_infrastructure"
	"meetup/crypto"
	"meetup/env"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

const (
	COOKIE_NAME_TOKEN = "access_token"
	// ERROR_REDIRECT_COOKIE_NAME は API 等から /login へリダイレクトする際の短いフラッシュ用（loginPage がサーバー側で読み取り削除する）。
	ERROR_REDIRECT_COOKIE_NAME = "error-redirect"
)

// CustomClaims は echo-jwt でパースする JWT クレーム（ログイン時に発行）。
type CustomClaims struct {
	UserID   int64  `json:"user_id"`
	Email    string `json:"email"`
	RoleID   int64  `json:"role_id"`
	RoleName string `json:"name"`
	jwt.RegisteredClaims
}

// SetErrorFlashCookie は短寿命のフラッシュ用 Cookie を付与する（/login で読み取り削除）。
func SetErrorFlashCookie(c *echo.Context, message string) {
	c.SetCookie(&http.Cookie{
		Name:     ERROR_REDIRECT_COOKIE_NAME,
		Value:    url.QueryEscape(message),
		Path:     "/",
		MaxAge:   120,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func clearAccessTokenCookie(c *echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     COOKIE_NAME_TOKEN,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   env.IsProduct(),
		SameSite: http.SameSiteLaxMode,
	})
}

func GetJWTConfig() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c *echo.Context) jwt.Claims {
			return new(CustomClaims)
		},
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

		stored, err := infrastructure.GetUserPasswordByEmail(c.Request().Context(), hm.db, info.E)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "ログインに失敗しました。"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		ok, err := crypto.VerifyPassword(stored, info.P)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "ログインに失敗しました。"})
		}

		user, err := infrastructure.GetUserInfo(c.Request().Context(), hm.db, info.E, stored, "Role")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}

		claims := CustomClaims{
			UserID:   user.ID,
			Email:    user.Email,
			RoleID:   user.RoleID,
			RoleName: user.Role.Name,
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   strconv.FormatInt(user.ID, 10),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				Issuer:    user.Name,
			},
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
