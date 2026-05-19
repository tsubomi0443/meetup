package interfaces

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"meetup/internal/usecases/dto"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

const (
	apiVersion = "v1"
	apiPath    = "/api/" + apiVersion
)

// setAPIHandler は JWT 保護下の REST API ルートグループを登録する。
//
// return:
//   - []echo.RouteInfo: 登録したルート情報
func (r *Router) setAPIHandler() (routeInfos []echo.RouteInfo) {
	group := r.e.Group(apiPath, GetJWTConfig())
	routeInfos = append(routeInfos, r.setupNoticeHandler(group)...)
	routeInfos = append(routeInfos, r.setupUserHandler(group)...)
	routeInfos = append(routeInfos, r.setupQuestionHandler(group)...)
	routeInfos = append(routeInfos, r.setupTagHandler(group)...)
	return
}

// setupNoticeHandler は通知 API ルートを登録する。
//
// args:
//   - group *echo.Group: /api/v1 グループ
//
// return:
//   - []echo.RouteInfo: 登録したルート情報
func (r *Router) setupNoticeHandler(group *echo.Group) (routeInfos []echo.RouteInfo) {
	const uri = "/notice"
	routeInfos = append(routeInfos, group.GET(uri, r.getNotice()))
	return
}

// setupUserHandler はユーザー API ルートを登録する。
//
// args:
//   - group *echo.Group: /api/v1 グループ
//
// return:
//   - []echo.RouteInfo: 登録したルート情報
func (r *Router) setupUserHandler(group *echo.Group) (routeInfos []echo.RouteInfo) {
	const uri = "/user"
	const uriWithID = uri + "/:id"
	const uriWithToken = uri + "/t"
	const api = "user"

	routeInfos = append(routeInfos, group.GET(uri, r.getUsers()))
	routeInfos = append(routeInfos, group.GET(uriWithToken, r.getUserFromToken()))
	routeInfos = append(routeInfos, group.POST(uri, r.registerUser(api, r.deps.Hub.SendCreateEvent)))
	routeInfos = append(routeInfos, group.PUT(uri, r.updateUserByID(api, r.deps.Hub.SendUpdateEvent)))
	routeInfos = append(routeInfos, group.DELETE(uriWithID, r.deleteUserByID(api, r.deps.Hub.SendDeleteEvent)))
	return
}

// setupQuestionHandler は質問 API ルートを登録する。
//
// args:
//   - group *echo.Group: /api/v1 グループ
//
// return:
//   - []echo.RouteInfo: 登録したルート情報
func (r *Router) setupQuestionHandler(group *echo.Group) (routeInfos []echo.RouteInfo) {
	const uri = "/question"
	const uriWithID = uri + "/:id"
	const api = "question"

	routeInfos = append(routeInfos, group.POST(uri, r.registerQuestion(api, r.deps.Hub.SendCreateEvent)))
	routeInfos = append(routeInfos, group.GET(uri, r.getQuestions()))
	routeInfos = append(routeInfos, group.GET(uriWithID, r.getQuestionByID()))
	routeInfos = append(routeInfos, group.DELETE(uriWithID, r.deleteQuestionByID(api, r.deps.Hub.SendDeleteEvent)))
	routeInfos = append(routeInfos, group.PUT(uri, r.updateQuestionByID(api, r.deps.Hub.SendUpdateEvent)))
	return
}

// setupTagHandler はタグ API ルートを登録する。
//
// args:
//   - group *echo.Group: /api/v1 グループ
//
// return:
//   - []echo.RouteInfo: 登録したルート情報
func (r *Router) setupTagHandler(group *echo.Group) (routeInfos []echo.RouteInfo) {
	const uri = "/tag"
	const uriWithID = uri + "/:id"
	const api = "tag"

	routeInfos = append(routeInfos, group.GET(uri, r.getTags()))
	routeInfos = append(routeInfos, group.POST(uri, r.registerTag(api, r.deps.Hub.SendCreateEvent)))
	routeInfos = append(routeInfos, group.PUT(uri, r.updateTag(api, r.deps.Hub.SendUpdateEvent)))
	routeInfos = append(routeInfos, group.DELETE(uriWithID, r.deleteTagByID(api, r.deps.Hub.SendDeleteEvent)))
	return
}

// registerUser はユーザー登録 POST ハンドラを返す。成功時に SSE 作成イベントを送る。
//
// args:
//   - api string: SSE イベント用 API 識別子
//   - sendEvent func(string, string): SSE 配信コールバック
//
// return:
//   - echo.HandlerFunc: POST /user 用ハンドラ
func (r *Router) registerUser(api string, sendEvent func(string, string)) echo.HandlerFunc {
	return func(c *echo.Context) error {
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}
		defer c.Request().Body.Close()

		var form dto.UserForm
		if err := json.Unmarshal(body, &form); err != nil {
			return err
		}
		created, err := r.deps.User.Register(c.Request().Context(), form)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Errorf("Create user server error %w", err))
		}
		payload, err := json.Marshal(created)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		sendEvent(api, string(payload))
		return c.JSON(http.StatusOK, nil)
	}
}

// registerQuestion は質問登録 POST ハンドラを返す。成功時に SSE 作成イベントを送る。
//
// args:
//   - api string: SSE イベント用 API 識別子
//   - sendEvent func(string, string): SSE 配信コールバック
//
// return:
//   - echo.HandlerFunc: POST /question 用ハンドラ
func (r *Router) registerQuestion(api string, sendEvent func(string, string)) echo.HandlerFunc {
	return func(c *echo.Context) error {
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}
		defer c.Request().Body.Close()

		var form dto.QuestionForm
		if err := json.Unmarshal(body, &form); err != nil {
			return err
		}
		created, err := r.deps.Question.Register(c.Request().Context(), form)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Errorf("Create question server error %w", err))
		}
		payload, err := json.Marshal(created)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		sendEvent(api, string(payload))
		return c.JSON(http.StatusOK, nil)
	}
}

// registerTag はタグ登録 POST ハンドラを返す。成功時に SSE 作成イベントを送る。
//
// args:
//   - api string: SSE イベント用 API 識別子
//   - sendEvent func(string, string): SSE 配信コールバック
//
// return:
//   - echo.HandlerFunc: POST /tag 用ハンドラ
func (r *Router) registerTag(api string, sendEvent func(string, string)) echo.HandlerFunc {
	return func(c *echo.Context) error {
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}
		defer c.Request().Body.Close()

		var form dto.TagForm
		if err := json.Unmarshal(body, &form); err != nil {
			return err
		}
		loaded, err := r.deps.Tag.Register(c.Request().Context(), form)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Errorf("Create new tag error %w\n", err))
		}
		payload, err := json.Marshal(loaded)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		sendEvent(api, string(payload))
		return c.JSON(http.StatusOK, nil)
	}
}

// updateTag はタグ更新 PUT ハンドラを返す。成功時に SSE 更新イベントを送る。
//
// args:
//   - api string: SSE イベント用 API 識別子
//   - sendEvent func(string, string): SSE 配信コールバック
//
// return:
//   - echo.HandlerFunc: PUT /tag 用ハンドラ
func (r *Router) updateTag(api string, sendEvent func(string, string)) echo.HandlerFunc {
	return func(c *echo.Context) error {
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}
		defer c.Request().Body.Close()

		var form dto.TagForm
		if err := json.Unmarshal(body, &form); err != nil {
			return err
		}
		if err := r.deps.Tag.Update(c.Request().Context(), form); err != nil {
			return err
		}
		sendEvent(api, string(body))
		return c.JSON(http.StatusOK, nil)
	}
}

// deleteTagByID はタグ削除 DELETE ハンドラを返す。成功時に SSE 削除イベントを送る。
//
// args:
//   - api string: SSE イベント用 API 識別子
//   - sendEvent func(string, string): SSE 配信コールバック
//
// return:
//   - echo.HandlerFunc: DELETE /tag/:id 用ハンドラ
func (r *Router) deleteTagByID(api string, sendEvent func(string, string)) echo.HandlerFunc {
	return func(c *echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return err
		}
		if err := r.deps.Tag.DeleteByID(c.Request().Context(), id); err != nil {
			return err
		}
		sendEvent(api, idStr)
		return c.JSON(http.StatusOK, nil)
	}
}

// getUsers はユーザー一覧 GET ハンドラを返す。
//
// return:
//   - echo.HandlerFunc: GET /user 用ハンドラ
func (r *Router) getUsers() echo.HandlerFunc {
	return func(c *echo.Context) error {
		users, err := r.deps.User.GetAll(c.Request().Context())
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, users)
	}
}

// getQuestions は質問一覧 GET ハンドラを返す。
//
// return:
//   - echo.HandlerFunc: GET /question 用ハンドラ
func (r *Router) getQuestions() echo.HandlerFunc {
	return func(c *echo.Context) error {
		questions, err := r.deps.Question.GetAll(c.Request().Context())
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, questions)
	}
}

// getQuestionByID は ID 指定で質問を取得する GET ハンドラを返す。
//
// return:
//   - echo.HandlerFunc: GET /question/:id 用ハンドラ
func (r *Router) getQuestionByID() echo.HandlerFunc {
	return func(c *echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 0, 10)
		if err != nil {
			fmt.Println(err)
			return err
		}
		model, err := r.deps.Question.GetByID(c.Request().Context(), id)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, model)
	}
}

// getTags はタグ一覧 GET ハンドラを返す。
//
// return:
//   - echo.HandlerFunc: GET /tag 用ハンドラ
func (r *Router) getTags() echo.HandlerFunc {
	return func(c *echo.Context) error {
		tags, err := r.deps.Tag.GetAll(c.Request().Context())
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, tags)
	}
}

// updateQuestionByID は質問更新 PUT ハンドラを返す。通知イベントと SSE 更新を送る。
//
// args:
//   - api string: SSE イベント用 API 識別子
//   - sendEvent func(string, string): SSE 配信コールバック
//
// return:
//   - echo.HandlerFunc: PUT /question 用ハンドラ
func (r *Router) updateQuestionByID(api string, sendEvent func(string, string)) echo.HandlerFunc {
	return func(c *echo.Context) error {
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}
		defer c.Request().Body.Close()
		var form dto.QuestionForm
		if err := json.Unmarshal(body, &form); err != nil {
			return err
		}
		actorID, hasActor := actorUserIDFromToken(c)
		updatedModel, loaded, err := r.deps.Question.Update(c.Request().Context(), form, actorID, hasActor)
		if err != nil {
			return err
		}
		r.deps.NoticeEvents.UpdateQuestion(updatedModel)

		payload, err := json.Marshal(loaded)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		sendEvent(api, string(payload))
		return c.JSON(http.StatusOK, nil)
	}
}

// deleteQuestionByID は質問削除 DELETE ハンドラを返す。通知イベントと SSE 削除を送る。
//
// args:
//   - api string: SSE イベント用 API 識別子
//   - sendEvent func(string, string): SSE 配信コールバック
//
// return:
//   - echo.HandlerFunc: DELETE /question/:id 用ハンドラ
func (r *Router) deleteQuestionByID(api string, sendEvent func(string, string)) echo.HandlerFunc {
	return func(c *echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 0, 10)
		if err != nil {
			fmt.Println(err)
			return err
		}
		if err := r.deps.Question.DeleteByID(c.Request().Context(), id); err != nil {
			return err
		}
		r.deps.NoticeEvents.DeleteQuestion(id)
		sendEvent(api, idStr)
		return c.JSON(http.StatusOK, "")
	}
}

// deleteUserByID はユーザー削除 DELETE ハンドラを返す。成功時に SSE 削除イベントを送る。
//
// args:
//   - api string: SSE イベント用 API 識別子
//   - sendEvent func(string, string): SSE 配信コールバック
//
// return:
//   - echo.HandlerFunc: DELETE /user/:id 用ハンドラ
func (r *Router) deleteUserByID(api string, sendEvent func(string, string)) echo.HandlerFunc {
	return func(c *echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 0, 10)
		if err != nil {
			return err
		}
		if err := r.deps.User.DeleteByID(c.Request().Context(), id); err != nil {
			return err
		}
		sendEvent(api, idStr)
		return c.JSON(http.StatusOK, nil)
	}
}

// getUserFromToken は JWT からユーザ ID を解決しユーザー情報を返す GET ハンドラを返す。
//
// return:
//   - echo.HandlerFunc: GET /user/t 用ハンドラ
func (r *Router) getUserFromToken() echo.HandlerFunc {
	return func(c *echo.Context) error {
		token := c.Get("user").(*jwt.Token)
		claims, ok := token.Claims.(*CustomClaims)
		if !ok || claims == nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid claims")
		}
		id := claims.UserID
		if id <= 0 && claims.Subject != "" {
			if parsed, err := strconv.ParseInt(claims.Subject, 10, 64); err == nil {
				id = parsed
			}
		}
		if id <= 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
		}

		user, err := r.deps.User.GetByID(c.Request().Context(), id)
		if err != nil {
			clearAccessTokenCookie(c)
			SetErrorFlashCookie(c, "ログイン情報が見つかりません。\nログインからやり直してください。")
			return c.Redirect(http.StatusSeeOther, "/login")
		}
		return c.JSON(http.StatusOK, user)
	}
}

// updateUserByID はユーザー更新 PUT ハンドラを返す。成功時に SSE 更新イベントを送る。
//
// args:
//   - api string: SSE イベント用 API 識別子
//   - sendEvent func(string, string): SSE 配信コールバック
//
// return:
//   - echo.HandlerFunc: PUT /user 用ハンドラ
func (r *Router) updateUserByID(api string, sendEvent func(string, string)) echo.HandlerFunc {
	return func(c *echo.Context) error {
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}
		defer c.Request().Body.Close()

		var form dto.UserForm
		if err := json.Unmarshal(body, &form); err != nil {
			return err
		}
		if err := r.deps.User.Update(c.Request().Context(), form); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		sendEvent(api, string(body))
		return c.JSON(http.StatusOK, nil)
	}
}

// actorUserIDFromToken は JWT から操作者ユーザ ID を取り出す。
//
// args:
//   - c *echo.Context: JWT ミドルウェア適用済みコンテキスト
//
// return:
//   - int64: ユーザ ID（取得できない場合は 0）
//   - bool: 取得できた場合 true
func actorUserIDFromToken(c *echo.Context) (int64, bool) {
	raw := c.Get("user")
	if raw == nil {
		return 0, false
	}
	token, ok := raw.(*jwt.Token)
	if !ok || token == nil {
		return 0, false
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok || claims == nil {
		return 0, false
	}
	id := claims.UserID
	if id <= 0 && claims.Subject != "" {
		if parsed, err := strconv.ParseInt(claims.Subject, 10, 64); err == nil {
			id = parsed
		}
	}
	if id <= 0 {
		return 0, false
	}
	return id, true
}

// getNotice は通知一覧 GET ハンドラを返す。
//
// return:
//   - echo.HandlerFunc: GET /notice 用ハンドラ
func (r *Router) getNotice() echo.HandlerFunc {
	return func(c *echo.Context) error {
		models, err := r.deps.Notice.GetAll(c.Request().Context())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, models)
	}
}
