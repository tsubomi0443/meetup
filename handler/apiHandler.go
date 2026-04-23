package handler

import (
	"encoding/json"
	"fmt"
	"io"
	infrastructure "meetup/_mac_infrastructure"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
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
		return infrastructure.Register(c.Request().Context(), hm.db, &data)
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
			return err
		}
		data := infrastructure.QuestionToEntity(form)
		return infrastructure.Register(c.Request().Context(), hm.db, &data)
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
		return infrastructure.Register(c.Request().Context(), hm.db, &model)
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
		if _, err := infrastructure.UpdateByID(c.Request().Context(), hm.db, model.ID, model, "Category", "TagManagers"); err != nil {
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

		if err := infrastructure.DeleteTagByID(c.Request().Context(), hm.db, id); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, nil)
	}
}

func (hm *HandlerManager) getUsers() echo.HandlerFunc {
	return func(c *echo.Context) error {
		users, err := infrastructure.GetUsers(c.Request().Context(), hm.db)
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
		model, err := infrastructure.GetQuestion(c.Request().Context(), hm.db, id)
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
		if err := infrastructure.UpdateQuestionInTransaction(c.Request().Context(), hm.db, updatedModel); err != nil {
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
		if err := infrastructure.DeleteQuestionByID(c.Request().Context(), hm.db, id); err != nil {
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
		if err := infrastructure.DeleteUserByID(c.Request().Context(), hm.db, id); err != nil {
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
		if _, err := infrastructure.UpdateByID(c.Request().Context(), hm.db, updatedModel.ID, updatedModel, "Role"); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, nil)
	}
}
