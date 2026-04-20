package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	infrastructure "meetup/_mac_infrastructure"
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
	routeInfos = append(routeInfos, hm.setupUserHandler()...)
	routeInfos = append(routeInfos, hm.setupQuestionHandler()...)

	return
}

func (hm *HandlerManager) setupUserHandler() (routeInfos []echo.RouteInfo) {
	const uri = apiPath + "/user"
	const uriWithID = uri + "/:id"

	routeInfos = append(routeInfos, hm.e.POST(uri, hm.registerUser()))
	routeInfos = append(routeInfos, hm.e.PUT(uri, hm.updateUserByID()))
	routeInfos = append(routeInfos, hm.e.DELETE(uriWithID, hm.deleteUserByID()))
	return
}

func (hm *HandlerManager) setupQuestionHandler() (routeInfos []echo.RouteInfo) {
	const uri = apiPath + "/question"
	const uriWithID = uri + "/:id"

	routeInfos = append(routeInfos, hm.e.POST(uri, hm.registerQuestion()))
	routeInfos = append(routeInfos, hm.e.GET(uriWithID, hm.getQuestionByID()))
	routeInfos = append(routeInfos, hm.e.DELETE(uriWithID, hm.deleteQuestionByID()))
	routeInfos = append(routeInfos, hm.e.PUT(uri, hm.updateQuestionByID()))
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
		return register(c.Request().Context(), hm.db, &data)
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

		return register(c.Request().Context(), hm.db, &data)
	}
}

func (hm *HandlerManager) registerTag() echo.HandlerFunc {
	return func(c *echo.Context) error {
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}
		defer c.Request().Body.Close()

		model := infrastructure.Tag{}
		if err := json.Unmarshal(body, &model); err != nil {
			return err
		}
		return register(c.Request().Context(), hm.db, model)
	}
}

func (hm *HandlerManager) getUsers() echo.HandlerFunc {
	return func(c *echo.Context) error {
		users, err := getUsers(c.Request().Context(), hm.db)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, users)
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
		updatedModel := infrastructure.Question{}
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}
		defer c.Request().Body.Close()
		if err := json.Unmarshal(body, &updatedModel); err != nil {
			return err
		}

		if _, err := updates(c.Request().Context(), hm.db, updatedModel); err != nil {
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

func (hm *HandlerManager) deleteUserByID() echo.HandlerFunc {
	return func(c *echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 0, 10)
		if err != nil {
			return err
		}
		if err := deleteUser(c.Request().Context(), hm.db, id); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, nil)
	}
}

func (hm *HandlerManager) updateUserByID() echo.HandlerFunc {
	return func(c *echo.Context) error {
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}
		defer c.Request().Body.Close()

		updatedModel := infrastructure.User{}
		if err := json.Unmarshal(body, &updatedModel); err != nil {
			return err
		}

		updatedModel.Role = infrastructure.Role{}
		fmt.Println(updatedModel)
		if _, err := updates(c.Request().Context(), hm.db, updatedModel, "Role"); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, nil)
	}
}

func getUsers(ctx context.Context, db *gorm.DB) (models []infrastructure.User, err error) {
	models, err = gorm.G[infrastructure.User](db).
		Preload("Role", nil).
		Select("id, name, email, role_id").Find(ctx)
	return
}

func register[T any](ctx context.Context, db *gorm.DB, model T, preloads ...string) error {
	var v = gorm.G[T](db)
	for _, preload := range preloads {
		v.Preload(preload, nil)
	}
	return v.Create(ctx, &model)
}

func updates[T any](ctx context.Context, db *gorm.DB, model T, preloads ...string) (int, error) {
	var v = gorm.G[T](db)
	for _, preload := range preloads {
		v.Preload(preload, nil)
	}
	return v.Updates(ctx, model)
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

func deleteUser(ctx context.Context, db *gorm.DB, id int64) error {
	if _, err := gorm.G[infrastructure.User](db).
		Preload("Role", nil).
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

func getTags(ctx context.Context, db *gorm.DB) (models []infrastructure.Tag, err error) {
	models, err = gorm.G[infrastructure.Tag](db).
		Preload("Category", nil).
		Find(ctx)
	return
}
