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
	routeInfos = append(routeInfos, hm.setupUserHandler()...)
	routeInfos = append(routeInfos, hm.setupQuestionHandler()...)
	routeInfos = append(routeInfos, hm.setupTagHandler()...)

	return
}

func (hm *HandlerManager) setupUserHandler() (routeInfos []echo.RouteInfo) {
	const uri = "/user"
	const uriWithID = uri + "/:id"
	var api = hm.e.Group(apiPath, GetJWTConfig())

	routeInfos = append(routeInfos, api.POST(uri, hm.registerUser()))
	routeInfos = append(routeInfos, api.PUT(uri, hm.updateUserByID()))
	routeInfos = append(routeInfos, api.DELETE(uriWithID, hm.deleteUserByID()))
	return
}

func (hm *HandlerManager) setupQuestionHandler() (routeInfos []echo.RouteInfo) {
	const uri = "/question"
	const uriWithID = uri + "/:id"
	var api = hm.e.Group(apiPath, GetJWTConfig())

	routeInfos = append(routeInfos, api.POST(uri, hm.registerQuestion()))
	routeInfos = append(routeInfos, api.GET(uriWithID, hm.getQuestionByID()))
	routeInfos = append(routeInfos, api.DELETE(uriWithID, hm.deleteQuestionByID()))
	routeInfos = append(routeInfos, api.PUT(uri, hm.updateQuestionByID()))
	return
}

func (hm *HandlerManager) setupTagHandler() (routeInfos []echo.RouteInfo) {
	const uri = "/tag"
	const uriWithID = uri + "/:id"
	var api = hm.e.Group(apiPath, GetJWTConfig())

	routeInfos = append(routeInfos, api.POST(uri, hm.registerTag()))
	routeInfos = append(routeInfos, api.PUT(uri, hm.updateTag()))
	routeInfos = append(routeInfos, api.DELETE(uriWithID, hm.deleteTagByID()))
	return
}

func (hm *HandlerManager) registerUser() echo.HandlerFunc {
	return func(c *echo.Context) error {
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}
		defer c.Request().Body.Close()

		var form infrastructure.UserForm
		if err := json.Unmarshal(body, &form); err != nil {
			return err
		}
		data := infrastructure.UserToEntityNoRole(form)
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

		var form infrastructure.QuestionForm
		if err := json.Unmarshal(body, &form); err != nil {
			fmt.Println(err.Error())
			return err
		}
		data := infrastructure.QuestionToEntity(form)
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

		var form infrastructure.TagForm
		if err := json.Unmarshal(body, &form); err != nil {
			return err
		}
		model := infrastructure.TagToEntity(form)
		return register(c.Request().Context(), hm.db, &model)
	}
}

func (hm *HandlerManager) updateTag() echo.HandlerFunc {
	return func(c *echo.Context) error {
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}
		defer c.Request().Body.Close()

		var form infrastructure.TagForm
		if err := json.Unmarshal(body, &form); err != nil {
			return err
		}
		model := infrastructure.TagToEntityNoRelations(form)
		fmt.Println(model)
		if _, err := updates(c.Request().Context(), hm.db, model); err != nil {
			fmt.Println(err)
			return err
		}
		return c.JSON(http.StatusOK, nil)
	}
}

func (hm *HandlerManager) deleteTagByID() echo.HandlerFunc {
	return func(c *echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return err
		}

		if err := deleteTag(c.Request().Context(), hm.db, id); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, nil)
	}
}

func (hm *HandlerManager) getUsers() echo.HandlerFunc {
	return func(c *echo.Context) error {
		users, err := getUsers(c.Request().Context(), hm.db)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, infrastructure.UserFormsFromEntities(users))
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
		return c.JSON(http.StatusOK, infrastructure.QuestionFromEntity(model))
	}
}

func (hm *HandlerManager) updateQuestionByID() echo.HandlerFunc {
	return func(c *echo.Context) error {
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}
		defer c.Request().Body.Close()
		var form infrastructure.QuestionForm
		if err := json.Unmarshal(body, &form); err != nil {
			return err
		}
		updatedModel := infrastructure.QuestionToEntity(form)
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

		var form infrastructure.UserForm
		if err := json.Unmarshal(body, &form); err != nil {
			return err
		}
		updatedModel := infrastructure.UserToEntityNoRole(form)
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

func deleteTag(ctx context.Context, db *gorm.DB, id int64) error {
	if _, err := gorm.G[infrastructure.Tag](db).Where("id = ?", id).Delete(ctx); err != nil {
		return err
	}
	return nil
}
