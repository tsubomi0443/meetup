package handler

import (
	"net/http"

	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
)

// TODO; 基本的なJWT認証の一部。まだ本実装には進まない
func GetJWTConfig() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte("secret"),
		ErrorHandler: func(c *echo.Context, err error) error {
			return c.Redirect(http.StatusFound, "/login")
		},
	})
}
