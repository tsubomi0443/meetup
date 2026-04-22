package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	infrastructure "meetup/_mac_infrastructure"
	"net/http"
	"strconv"
	"time"

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
		if _, err := updateInTransaction(c.Request().Context(), hm.db, model, "Category", "TagManagers"); err != nil {
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
		if err := updateQuestionInTransaction(c.Request().Context(), hm.db, updatedModel); err != nil {
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
		if _, err := updateInTransaction(c.Request().Context(), hm.db, updatedModel, "Role"); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, nil)
	}
}

func getUsers(ctx context.Context, db *gorm.DB) (models []infrastructure.User, err error) {
	models, err = gorm.G[infrastructure.User](db).
		Where("role_id <> ?", 1).
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

// updateInTransaction は単一の DB トランザクション内で gorm.Updates を実行する。
// omit に関連名を渡し、Role/Category など中間テーブル向けの関連を更新対象から外す。
// 既存の updates と違い、こちらは更新用。事前読み込みは行わない。
func updateInTransaction[T any](ctx context.Context, db *gorm.DB, model T, omit ...string) (rowsAffected int, err error) {
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		m := model
		res := tx.Omit(omit...).Model(&m).Updates(&m)
		rowsAffected = int(res.RowsAffected)
		return res.Error
	})
	return
}

// updateQuestionInTransaction は Question の1行更新と、QuestionToEntity で組み立てた
// 1対多の関連（Answer, Support, TagManagers, Memos, Notices）の同期を1トランザクションで行う。
// フォーム上の下位の質問（SubQuestions）やエスカレーション等は本関数では永続化しない。
func updateQuestionInTransaction(ctx context.Context, db *gorm.DB, q infrastructure.Question) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1) Answer: 新規作成 or 既存更新し、親の answer_id を整合させる
		if q.Answer != nil {
			a := *q.Answer
			if a.ID == 0 {
				if err := tx.Create(&a).Error; err != nil {
					return err
				}
				aid := a.ID
				q.AnswerID = &aid
				q.Answer = &a
			} else {
				if _, err := gorm.G[infrastructure.Answer](tx).
					Omit("User", "ReferManagers").
					Where("id = ?", a.ID).
					Updates(ctx, a); err != nil {
					return err
				}
				if q.AnswerID == nil {
					aid := a.ID
					q.AnswerID = &aid
				}
			}
		}
		// 2) Support: 新規作成 or 既存更新
		if q.Support != nil {
			s := *q.Support
			if s.ID == 0 {
				if err := tx.Create(&s).Error; err != nil {
					return err
				}
				sid := s.ID
				q.SupportID = &sid
				q.Support = &s
			} else {
				if err := tx.Model(&s).Omit("User", "SupportStatus").Updates(&s).Error; err != nil {
					return err
				}
				if q.SupportID == nil {
					sid := s.ID
					q.SupportID = &sid
				}
			}
		}
		// 3) 親の questions: スカラ列と FK のみ（CreatedAt は更新で変えない）
		if _, err := gorm.G[infrastructure.Question](tx).
			Omit("Answer", "Memos", "Notices", "TagManagers", "Support", "CreatedAt").
			Where("id = ?", q.ID).
			Updates(ctx, q); err != nil {
			return err
		}
		ref := infrastructure.Question{ID: q.ID}
		// 4) タグ紐づけ（tag_managers）
		for i := range q.TagManagers {
			if q.TagManagers[i].QuestionID == 0 {
				q.TagManagers[i].QuestionID = q.ID
			}
		}
		if len(q.TagManagers) == 0 {
			if err := tx.Model(&ref).Association("TagManagers").Clear(); err != nil {
				return err
			}
		} else {
			if err := tx.Model(&ref).Association("TagManagers").Replace(q.TagManagers); err != nil {
				return err
			}
		}
		// 5) メモ
		for i := range q.Memos {
			if q.Memos[i].QuestionID == 0 {
				q.Memos[i].QuestionID = q.ID
			}
		}
		if len(q.Memos) == 0 {
			if err := tx.Model(&ref).Association("Memos").Clear(); err != nil {
				return err
			}
		} else {
			if err := tx.Model(&ref).Association("Memos").Replace(q.Memos); err != nil {
				return err
			}
		}
		// 6) 通知
		for i := range q.Notices {
			if q.Notices[i].QuestionID == nil {
				qid := q.ID
				q.Notices[i].QuestionID = &qid
			}
		}
		if len(q.Notices) == 0 {
			if err := tx.Model(&ref).Association("Notices").Clear(); err != nil {
				return err
			}
		} else {
			if err := tx.Model(&ref).Association("Notices").Replace(q.Notices); err != nil {
				return err
			}
		}
		return nil
	})
}

func deleteQuestion(ctx context.Context, db *gorm.DB, id int64) error {
	if _, err := gorm.G[infrastructure.Question](db).
		Preload("Answer", nil).
		Preload("Answer.User", nil).
		Preload("Answer.User.Role", nil).
		Preload("Answer.ReferManagers", nil).
		Preload("Answer.ReferManagers.Refer", nil).
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
		Preload("Answer", nil).
		Preload("Answer.User", func(db gorm.PreloadBuilder) error {
			db.Select("id", "name", "email", "role_id")
			return nil
		}).
		Preload("Answer.User.Role", nil).
		Preload("Answer.ReferManagers", nil).
		Preload("Answer.ReferManagers.Refer", nil).
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
		Preload("Answer", nil).
		Preload("Answer.User", func(db gorm.PreloadBuilder) error {
			db.Select("id", "name", "email", "role_id")
			return nil
		}).
		Preload("Answer.User.Role", nil).
		Preload("Answer.ReferManagers", nil).
		Preload("Answer.ReferManagers.Refer", nil).
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

func getNotice(ctx context.Context, db *gorm.DB) ([]infrastructure.Question, error) {
	questions, err := gorm.G[infrastructure.Question](db).Find(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	limit := now.Add(3 * 24 * time.Hour)
	var recent []infrastructure.Question
	for _, q := range questions {
		if q.Due == nil {
			continue
		}
		d := *q.Due
		if !d.Before(now) && d.Before(limit) {
			recent = append(recent, q)
		}
	}
	return recent, nil
}
