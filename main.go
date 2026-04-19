package main

import (
	"fmt"
	"html/template"
	"log"
	"log/slog"
	infrastructure "meetup/_mac_infrastructure"
	"meetup/env"
	"meetup/handler"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	godotenv.Load(".env")
}

func main() {
	db, err := setupDB()
	if err != nil {
		log.Fatalln(err)
	}

	e := setupEcho()
	hub := handler.NewHub()
	handlerManager := handler.NewHandlerManager(db, e, hub)
	fmt.Println(handlerManager.SetupHandlers())

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

	// テンプレートの設定
	e.Renderer = &echo.TemplateRenderer{
		Template: template.Must(template.ParseGlob("templates/**/*.html")),
	}

	// ミドルウェア
	e.Use(middleware.Recover())
	return e
}

func setupDB() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(env.GetDSN()))
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&infrastructure.Role{},
		&infrastructure.Category{},
		&infrastructure.SupportStatus{},
		&infrastructure.User{},
		&infrastructure.Support{},
		&infrastructure.Question{},
		&infrastructure.Answer{},
		&infrastructure.Memo{},
		&infrastructure.Refer{},
		&infrastructure.Tag{},
		&infrastructure.ReferManager{},
		&infrastructure.TagManager{},
		&infrastructure.Escalation{},
	); err != nil {
		return nil, err
	}

	return db, nil
}
