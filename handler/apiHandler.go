package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"meetup/_mac_infrastructure"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

const (
	apiVersion = "v1"
	apiPath    = "/api/" + apiVersion
)

func (hm *HandlerManager) SetAPIHandler() (routeInfos []echo.RouteInfo) {
	// HTMX API エンドポイント
	routeInfos = append(routeInfos, hm.e.GET(apiPath+"/joke", func(c *echo.Context) error {
		// サンプルレスポンス
		return c.HTML(http.StatusOK, `<div class="alert alert-success">サンプルジョーク: なぜプログラマーはハロウィンが好きか？クリスマスが怖いから（Oct 31 == Dec 25）</div>`)
	}))
	routeInfos = append(routeInfos, hm.e.POST(apiPath+"/user", hm.registerUser()))
	routeInfos = append(routeInfos, hm.e.GET(apiPath+"/question/:id", hm.getQuestionByID()))
	routeInfos = append(routeInfos, hm.e.POST(apiPath+"/question", hm.registerQuestion()))
	routeInfos = append(routeInfos, hm.e.DELETE(apiPath+"/question/:id", hm.deleteQuestionByID()))
	routeInfos = append(routeInfos, hm.e.PUT(apiPath+"/question", hm.updateQuestionByID()))

	return
}

func (hm *HandlerManager) registerUser() echo.HandlerFunc {
	return func(c *echo.Context) error {
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}
		defer c.Request().Body.Close()

		data := infrastructure.User{}
		if err := json.Unmarshal(body, &data); err != nil {
			return err
		}
		fmt.Println(data)
		return registerUser(c.Request().Context(), hm.db, &data)
	}
}

func (hm *HandlerManager) registerQuestion() echo.HandlerFunc {
	return func(c *echo.Context) error {
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}
		defer c.Request().Body.Close()

		data := infrastructure.Question{}
		if err := json.Unmarshal(body, &data); err != nil {
			fmt.Println(err.Error())
			return err
		}

		return registerQuestion(c.Request().Context(), hm.db, &data)
	}
}

// 後でDomainフォルダを作成し、そちらで管理。DTO関係の関数群です。
func (hm *HandlerManager) getQuestionByID() echo.HandlerFunc {
	return func(c *echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 0, 10)
		if err != nil {
			fmt.Println(err)
			return err
		}
		model, err := getQuestion(c.Request().Context(), hm.db, id)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, model)
	}
}

func (hm *HandlerManager) updateQuestionByID() echo.HandlerFunc {
	return func(c *echo.Context) error {
		fmt.Println("IN")
		updatedModel := infrastructure.Question{}
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			fmt.Println(err)
			return err
		}
		defer c.Request().Body.Close()
		fmt.Println("READBODY")
		fmt.Println(string(body))
		if err := json.Unmarshal(body, &updatedModel); err != nil {
			fmt.Println("UNMARSHAL")
			return err
		}

		if _, err := gorm.G[infrastructure.Question](hm.db).Updates(c.Request().Context(), updatedModel); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, nil)
	}
}

func (hm *HandlerManager) deleteQuestionByID() echo.HandlerFunc {
	return func(c *echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 0, 10)
		if err != nil {
			fmt.Println(err)
			return err
		}
		if err := deleteQuestion(c.Request().Context(), hm.db, id); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, "")
	}
}

func registerUser(ctx context.Context, db *gorm.DB, user *infrastructure.User) error {
	return gorm.G[infrastructure.User](db).Create(ctx, user)
}

func getUsers(ctx context.Context, db *gorm.DB) (models []infrastructure.User, err error) {
	models, err = gorm.G[infrastructure.User](db).Find(ctx)
	return
}

func deleteQuestion(ctx context.Context, db *gorm.DB, id int64) error {
	if _, err := gorm.G[infrastructure.Question](db).
		Preload("Answers", nil).
		Preload("Answers.User", nil).
		Preload("Answers.User.Role", nil).
		Preload("Answers.ReferManagers", nil).
		Preload("Answers.ReferManagers.Refer", nil).
		Preload("Memos", nil).
		Preload("Memos.User", nil).
		Preload("Memos.User.Role", nil).
		Preload("TagManagers", nil).
		Preload("TagManagers.Tag", nil).
		Preload("TagManagers.Tag.Category", nil).
		Preload("Support", nil).
		Preload("Support.User", nil).
		Preload("Support.User.Role", nil).
		Preload("Support.SupportStatus", nil).
		Where("id = ?", id).
		Delete(ctx); err != nil {
		return err
	}
	return nil
}

func getQuestion(ctx context.Context, db *gorm.DB, id int64) (model infrastructure.Question, err error) {
	model, err = gorm.G[infrastructure.Question](db).
		Preload("Answers", nil).
		Preload("Answers.User", func(db gorm.PreloadBuilder) error {
			db.Select("id", "name", "email", "role_id")
			return nil
		}).
		Preload("Answers.User.Role", nil).
		Preload("Answers.ReferManagers", nil).
		Preload("Answers.ReferManagers.Refer", nil).
		Preload("Memos", nil).
		Preload("Memos.User", func(db gorm.PreloadBuilder) error {
			db.Select("id", "name", "email", "role_id")
			return nil
		}).
		Preload("Memos.User.Role", nil).
		Preload("TagManagers", nil).
		Preload("TagManagers.Tag", nil).
		Preload("TagManagers.Tag.Category", nil).
		Preload("Support", nil).
		Preload("Support.User", func(db gorm.PreloadBuilder) error {
			db.Select("id", "name", "email", "role_id")
			return nil
		}).
		Preload("Support.User.Role", nil).
		Preload("Support.SupportStatus", nil).
		Where("id = ?", id).
		First(ctx)
	return
}

func getQuestions(ctx context.Context, db *gorm.DB) (models []infrastructure.Question, err error) {
	models, err = gorm.G[infrastructure.Question](db).
		Preload("Answers", nil).
		Preload("Answers.User", func(db gorm.PreloadBuilder) error {
			db.Select("id", "name", "email", "role_id")
			return nil
		}).
		Preload("Answers.User.Role", nil).
		Preload("Answers.ReferManagers", nil).
		Preload("Answers.ReferManagers.Refer", nil).
		Preload("Memos", nil).
		Preload("Memos.User", func(db gorm.PreloadBuilder) error {
			db.Select("id", "name", "email", "role_id")
			return nil
		}).
		Preload("Memos.User.Role", nil).
		Preload("TagManagers", nil).
		Preload("TagManagers.Tag", nil).
		Preload("TagManagers.Tag.Category", nil).
		Preload("Support", nil).
		Preload("Support.User", func(db gorm.PreloadBuilder) error {
			db.Select("id", "name", "email", "role_id")
			return nil
		}).
		Preload("Support.User.Role", nil).
		Preload("Support.SupportStatus", nil).
		Find(ctx)
	return
}

func registerQuestion(ctx context.Context, db *gorm.DB, body *infrastructure.Question) error {
	return gorm.G[infrastructure.Question](db).Create(ctx, body)
}
