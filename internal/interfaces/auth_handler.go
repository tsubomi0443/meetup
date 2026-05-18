package interfaces

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"meetup/internal/infrastructures/config"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

const (
	CookieNameToken          = "access_token"
	ErrorRedirectCookieName  = "error-redirect"
)

type CustomClaims struct {
	UserID   int64  `json:"user_id"`
	Email    string `json:"email"`
	RoleID   int64  `json:"role_id"`
	RoleName string `json:"name"`
	jwt.RegisteredClaims
}

func SetErrorFlashCookie(c *echo.Context, message string) {
	c.SetCookie(&http.Cookie{
		Name:     ErrorRedirectCookieName,
		Value:    url.QueryEscape(message),
		Path:     "/",
		MaxAge:   120,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func clearAccessTokenCookie(c *echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     CookieNameToken,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   config.IsProduct(),
		SameSite: http.SameSiteLaxMode,
	})
}

func GetJWTConfig() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c *echo.Context) jwt.Claims {
			return new(CustomClaims)
		},
		SigningKey:  []byte(config.GetJWTKey()),
		TokenLookup: "cookie:" + CookieNameToken,
		ErrorHandler: func(c *echo.Context, err error) error {
			return c.Redirect(http.StatusFound, "/login")
		},
	})
}

func (r *Router) setupAuthHandler() (routeInfos []echo.RouteInfo) {
	routeInfos = append(routeInfos, r.e.POST("/login", r.loginHandler()))
	return
}

func (r *Router) loginHandler() echo.HandlerFunc {
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

		result, err := r.deps.Auth.Login(c.Request().Context(), info.E, info.P)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "ログインに失敗しました。"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		user := result.User
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
		signed, err := token.SignedString([]byte(config.GetJWTKey()))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		c.SetCookie(&http.Cookie{
			Name:     CookieNameToken,
			Value:    signed,
			Path:     "/",
			HttpOnly: true,
			Secure:   config.IsProduct(),
			SameSite: http.SameSiteLaxMode,
			Expires:  time.Now().Add(1 * time.Hour),
		})

		return c.JSON(http.StatusOK, map[string]string{
			"redirect": "/mock/5",
		})
	}
}
