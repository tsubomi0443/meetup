package main

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func init() {
	godotenv.Load(".env")
}

func main() {
	e := echo.New()

	// テンプレートの設定
	e.Renderer = &echo.TemplateRenderer{
		Template: template.Must(template.ParseGlob("templates/**/*.html")),
	}

	// ミドルウェア
	e.Use(middleware.Recover())

	// 静的ファイル配信
	e.Static("/static", "static")

	// ルート
	e.GET("/", func(c *echo.Context) error {
		return c.Render(http.StatusOK, "index.html", map[string]any{
			"Title": "Go + Echo + HTMX + Alpine.js",
		})
	})

	// ルート
	e.GET("/mock/:id", func(c *echo.Context) error {
		id := c.Param("id")
		return c.Render(http.StatusOK, fmt.Sprintf("mock%s.html", id), nil)
	})

	// HTMX API エンドポイント
	e.GET("/api/joke", func(c *echo.Context) error {
		// サンプルレスポンス
		return c.HTML(http.StatusOK, `<div class="alert alert-success">サンプルジョーク: なぜプログラマーはハロウィンが好きか？クリスマスが怖いから（Oct 31 == Dec 25）</div>`)
	})

	// SSE Hub
	hub := NewHub()
	go hub.Run()
	go func() {
		for t := range time.Tick(1 * time.Second) {
			hub.Broadcast <- Event{Event: "time-tick", Data: fmt.Sprintf(`<div>%s</div>`, t.Format("15:04:05"))}
		}
	}()

	e.GET("/sse", sseHandler(hub))

	// サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := e.Start(":" + port); err != nil {
		slog.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
