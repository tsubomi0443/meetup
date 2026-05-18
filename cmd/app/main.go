package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"meetup/internal/di"
	"meetup/internal/infrastructures/config"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf(".envファイルの読み込みでエラーが発生しました。詳細: %v", err)
	}

	// アプリ上の時刻をJSTへと設定（UTC+9:00）
	time.Local = time.FixedZone("JST", 9*60*60)
}

func main() {
	db, err := setupDB()
	if err != nil {
		log.Fatalln(err)
	}

	e := setupEcho()
	app := di.NewApp(db, e)
	app.Router.SetHubLogger(app.Router.Logging)
	go app.Router.PollingStart(context.Background())
	fmt.Println(app.Router.SetupHandlers())

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

func setupEcho() *echo.Echo {
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	tmpl := template.New("")
	err := filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// .html ファイルのみを対象とする
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			if _, err := tmpl.ParseFiles(path); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalln(err.Error())
	}

	// テンプレートの設定
	e.Renderer = &echo.TemplateRenderer{
		Template: tmpl,
	}

	// ミドルウェア
	e.Use(middleware.Recover())
	return e
}

func setupDB() (*gorm.DB, error) {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.Open(config.GetDSN()), &gorm.Config{
		// GORMがレコード作成・更新時に使う時刻をJSTに固定
		NowFunc: func() time.Time {
			return time.Now().In(jst)
		},
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}
