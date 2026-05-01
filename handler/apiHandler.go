package handler

import (
	"encoding/json"
	"fmt"
	"io"
	infrastructure "meetup/_mac_infrastructure"
	"meetup/crypto"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

const (
	apiVersion = "v1"
	apiPath    = "/api/" + apiVersion
)

func (hm *HandlerManager) SetAPIHandler() (routeInfos []echo.RouteInfo) {
	group := hm.e.Group(apiPath, GetJWTConfig())
	routeInfos = append(routeInfos, hm.setupNoticeHandler(group)...)
	routeInfos = append(routeInfos, hm.setupUserHandler(group)...)
	routeInfos = append(routeInfos, hm.setupQuestionHandler(group)...)
	routeInfos = append(routeInfos, hm.setupTagHandler(group)...)

	return
}

func (hm *HandlerManager) setupNoticeHandler(group *echo.Group) (routeInfos []echo.RouteInfo) {
	const uri = "/notice"
	routeInfos = append(routeInfos, group.GET(uri, hm.getNotice()))
	return
}
func (hm *HandlerManager) setupUserHandler(group *echo.Group) (routeInfos []echo.RouteInfo) {
	const uri = "/user"
	const uriWithID = uri + "/:id"
	const uriWithToken = uri + "/t"
	const api = "user"

	routeInfos = append(routeInfos, group.GET(uri, hm.getUsers()))
	routeInfos = append(routeInfos, group.GET(uriWithToken, hm.getUserFromToken()))
	routeInfos = append(routeInfos, group.POST(uri, hm.registerUser(api, hm.hub.sendCreateEvent)))
	routeInfos = append(routeInfos, group.PUT(uri, hm.updateUserByID(api, hm.hub.sendUpdateEvent)))
	routeInfos = append(routeInfos, group.DELETE(uriWithID, hm.deleteUserByID(api, hm.hub.sendDeleteEvent)))
	return
}

func (hm *HandlerManager) setupQuestionHandler(group *echo.Group) (routeInfos []echo.RouteInfo) {
	const uri = "/question"
	const uriWithID = uri + "/:id"
	const api = "question"

	routeInfos = append(routeInfos, group.POST(uri, hm.registerQuestion(api, hm.hub.sendCreateEvent)))
	routeInfos = append(routeInfos, group.GET(uri, hm.getQuestions()))
	routeInfos = append(routeInfos, group.GET(uriWithID, hm.getQuestionByID()))
	routeInfos = append(routeInfos, group.DELETE(uriWithID, hm.deleteQuestionByID(api, hm.hub.sendDeleteEvent)))
	routeInfos = append(routeInfos, group.PUT(uri, hm.updateQuestionByID(api, hm.hub.sendUpdateEvent)))
	return
}

func (hm *HandlerManager) setupTagHandler(group *echo.Group) (routeInfos []echo.RouteInfo) {
	const uri = "/tag"
	const uriWithID = uri + "/:id"
	const api = "tag"

	routeInfos = append(routeInfos, group.GET(uri, hm.getTags()))
	routeInfos = append(routeInfos, group.POST(uri, hm.registerTag(api, hm.hub.sendCreateEvent)))
	routeInfos = append(routeInfos, group.PUT(uri, hm.updateTag(api, hm.hub.sendUpdateEvent)))
	routeInfos = append(routeInfos, group.DELETE(uriWithID, hm.deleteTagByID(api, hm.hub.sendDeleteEvent)))
	return
}

func (hm *HandlerManager) registerUser(api string, sendEvent func(string, string)) echo.HandlerFunc {
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
		if strings.TrimSpace(form.Password) != "" {
			form.Password = crypto.EncryptSHA256(form.Password)
		}
		data := infrastructure.UserToEntityNoRole(form)
		if err := infrastructure.Register(c.Request().Context(), hm.db, &data); err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Errorf("Create user server error %w", err))
		}

		sendEvent(api, string(body))
		return c.JSON(http.StatusOK, nil)
	}
}

func (hm *HandlerManager) registerQuestion(api string, sendEvent func(string, string)) echo.HandlerFunc {
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
		if err := infrastructure.Register(c.Request().Context(), hm.db, &data); err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Errorf("Create question server error %w", err))
		}
		sendEvent(api, string(body))
		return c.JSON(http.StatusOK, nil)
	}
}

func (hm *HandlerManager) registerTag(api string, sendEvent func(string, string)) echo.HandlerFunc {
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
		if err := infrastructure.Register(c.Request().Context(), hm.db, &model); err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Errorf("Create new tag error %w\n", err))
		}

		sendEvent(api, string(body))
		return c.JSON(http.StatusOK, nil)
	}
}

func (hm *HandlerManager) updateTag(api string, sendEvent func(string, string)) echo.HandlerFunc {
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
		sendEvent(api, string(body))
		return c.JSON(http.StatusOK, nil)
	}
}

func (hm *HandlerManager) deleteTagByID(api string, sendEvent func(string, string)) echo.HandlerFunc {
	return func(c *echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return err
		}

		if err := infrastructure.DeleteTagByID(c.Request().Context(), hm.db, id); err != nil {
			return err
		}
		sendEvent(api, idStr)
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

func (hm *HandlerManager) getQuestions() echo.HandlerFunc {
	return func(c *echo.Context) error {
		questions, err := infrastructure.GetQuestions(c.Request().Context(), hm.db)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, infrastructure.QuestionFromEntities(questions))
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

func (hm *HandlerManager) getTags() echo.HandlerFunc {
	return func(c *echo.Context) error {
		tags, err := infrastructure.GetTags(c.Request().Context(), hm.db)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, infrastructure.TagFromEntities(tags))
	}
}

func (hm *HandlerManager) updateQuestionByID(api string, sendEvent func(string, string)) echo.HandlerFunc {
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
		hm.ne.UpdateQuestion(updatedModel)

		sendData := body
		for idx, memo := range form.Memos {
			if memo.ID == 0 {
				form.Memos[idx].ID = infrastructure.GetMaxByColumn[infrastructure.Memo](c.Request().Context(), hm.db, "id") + 1
			}
		}
		if sendData, err = json.Marshal(form); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		sendEvent(api, string(sendData))
		return c.JSON(http.StatusOK, nil)
	}
}

func (hm *HandlerManager) deleteQuestionByID(api string, sendEvent func(string, string)) echo.HandlerFunc {
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
		hm.ne.DeleteQuestion(id)
		sendEvent(api, idStr)
		return c.JSON(http.StatusOK, "")
	}
}

func (hm *HandlerManager) deleteUserByID(api string, sendEvent func(string, string)) echo.HandlerFunc {
	return func(c *echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 0, 10)
		if err != nil {
			return err
		}
		if err := infrastructure.DeleteUserByID(c.Request().Context(), hm.db, id); err != nil {
			return err
		}
		sendEvent(api, idStr)

		return c.JSON(http.StatusOK, nil)
	}
}

func (hm *HandlerManager) getUserFromToken() echo.HandlerFunc {
	return func(c *echo.Context) error {
		token := c.Get("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		subFloat, ok := claims["sub"].(float64)
		if !ok {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
		}
		id := int64(subFloat)

		user, err := infrastructure.GetUserByID(c.Request().Context(), hm.db, id)
		if err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, infrastructure.UserFromEntity(user))
	}
}

func (hm *HandlerManager) updateUserByID(api string, sendEvent func(string, string)) echo.HandlerFunc {
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
		if strings.TrimSpace(form.Password) != "" {
			form.Password = crypto.EncryptSHA256(form.Password)
		}
		updatedModel := infrastructure.UserToEntityNoRole(form)
		if _, err := infrastructure.UpdateByID(c.Request().Context(), hm.db, updatedModel.ID, updatedModel, "Role"); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		sendEvent(api, string(body))
		return c.JSON(http.StatusOK, nil)
	}
}

func (hm *HandlerManager) getNotice() echo.HandlerFunc {
	return func(c *echo.Context) error {
		models, err := infrastructure.GetNotice(c.Request().Context(), hm.db)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, infrastructure.NoticeFromEntities(models))
	}
}
