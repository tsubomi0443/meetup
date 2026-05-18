package postgres

import (
	"meetup/internal/domains/entity"

	
	"context"
	"database/sql"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GetUserByID は指定 ID のユーザーをロール込みで取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - id int64: ユーザー ID
//
// return:
//   - entity.User: ユーザー
//   - error: DB エラー
func GetUserByID(ctx context.Context, db *gorm.DB, id int64) (model entity.User, err error) {
	model, err = gorm.G[entity.User](db).Where("id = ?", id).Select("id, name, email, memo, role_id").
		Preload("Role", commonPreloadBuilder()).
		First(ctx)
	return
}

// GetUserPasswordByEmail はメールアドレスに紐づく保存パスワードハッシュを取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - email string: メールアドレス
//
// return:
//   - string: パスワードハッシュ
//   - error: DB エラー
func GetUserPasswordByEmail(ctx context.Context, db *gorm.DB, email string) (password string, err error) {
	u, err := gorm.G[entity.User](db).Where("email = ?", email).Select("id, password").First(ctx)
	if err != nil {
		return "", err
	}
	return u.Password, nil
}

// GetUserInfo はメールとパスワードハッシュでユーザーを取得する（Preload 指定可）。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - email string: メールアドレス
//   - pass string: 保存済みパスワードハッシュ
//   - preloads ...string: GORM Preload 名
//
// return:
//   - entity.User: ユーザー
//   - error: DB エラー
func GetUserInfo(ctx context.Context, db *gorm.DB, email, pass string, preloads ...string) (model entity.User, err error) {
	chain := gorm.G[entity.User](db).Where("email = ? AND password = ?", email, pass)
	for _, preload := range preloads {
		chain = chain.Preload(preload, commonPreloadBuilder())
	}
	model, err = chain.First(ctx)
	return
}

// GetUsers は管理者（role_id=1）以外のユーザー一覧を取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//
// return:
//   - []entity.User: ユーザー一覧
//   - error: DB エラー
func GetUsers(ctx context.Context, db *gorm.DB) (models []entity.User, err error) {
	models, err = gorm.G[entity.User](db).
		Where("role_id <> ?", 1).
		Preload("Role", commonPreloadBuilder()).
		Not("role_id = 1").
		Select("id, name, email, memo, role_id").
		Order("id").
		Find(ctx)
	return
}

// GetMasterData は entity.Role / entity.Category / entity.SupportStatus などマスタテーブルの一覧を取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//
// return:
//   - []T: マスタ行の一覧
//   - error: DB エラー
func GetMasterData[T entity.Role | entity.Category | entity.SupportStatus](ctx context.Context, db *gorm.DB) (models []T, err error) {
	models, err = gorm.G[T](db).Find(ctx)
	return
}

// Register は任意型のモデルを新規登録する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - model T: 登録するモデル
//   - preloads ...string: 作成後 Preload 名（未使用の場合あり）
//
// return:
//   - error: DB エラー
func Register[T any](ctx context.Context, db *gorm.DB, model T, preloads ...string) error {
	var v = gorm.G[T](db)
	for _, preload := range preloads {
		v.Preload(preload, nil)
	}
	return v.Create(ctx, &model)
}

// Updates はモデルを一括更新する（WHERE なしの Updates）。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - model T: 更新内容
//   - preloads ...string: Preload 名（未使用の場合あり）
//
// return:
//   - int: 更新行数
//   - error: DB エラー
func Updates[T any](ctx context.Context, db *gorm.DB, model T, preloads ...string) (int, error) {
	var v = gorm.G[T](db)
	for _, preload := range preloads {
		v.Preload(preload, nil)
	}
	return v.Updates(ctx, model)
}

// UpdateByID は主キー id を WHERE に固定して単一モデルを更新する。
// gorm.G[T](db) の Updates を利用し、トランザクションは張らない。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - id int64: 主キー
//   - model T: 更新内容
//   - omit ...string: Omit に渡す関連名・列名
//
// return:
//   - int: 更新行数
//   - error: DB エラー
func UpdateByID[T any](ctx context.Context, db *gorm.DB, id int64, model T, omit ...string) (int, error) {
	return gorm.G[T](db.WithContext(ctx)).
		Omit(omit...).
		Where("id = ?", id).
		Updates(ctx, model)
}

// UpdateInTransaction は単一トランザクション内で gorm.Updates を実行する。
// omit に関連名を渡し、中間テーブル向け関連を更新対象から外す。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - model T: 更新内容
//   - omit ...string: 更新から除外する関連名
//
// return:
//   - int: 更新行数
//   - error: DB エラー
func UpdateInTransaction[T any](ctx context.Context, db *gorm.DB, model T, omit ...string) (rowsAffected int, err error) {
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		m := model
		res := tx.Omit(omit...).Model(&m).Updates(&m)
		rowsAffected = int(res.RowsAffected)
		return res.Error
	})
	return
}

// postgresTableCoalesceMaxID はホワイトリスト上のテーブルで COALESCE(MAX(id), 0) を返す。
//
// args:
//   - tx *gorm.DB: トランザクション
//   - table string: テーブル名（memos, tag_managers 等）
//
// return:
//   - int64: 最大 ID（行なしは 0）
//   - error: 未対応テーブル・SQL エラー
func postgresTableCoalesceMaxID(tx *gorm.DB, table string) (int64, error) {
	switch table {
	case "memos", "tag_managers", "related_questions", "answers", "refer_managers":
	default:
		return 0, fmt.Errorf("postgresTableCoalesceMaxID: unsupported table %q", table)
	}
	var max int64
	q := fmt.Sprintf("SELECT COALESCE(MAX(id), 0) FROM %s", table)
	if err := tx.Raw(q).Scan(&max).Error; err != nil {
		return 0, err
	}
	return max, nil
}

// assignBulkInsertZeros は一括 INSERT 前に ID==0 の行へ連番の主キーを割り当てる。
// 混在 INSERT による *_pkey 重複を防ぐ。
//
// args:
//   - tx *gorm.DB: トランザクション
//   - rows []T: 挿入行
//   - table string: テーブル名
//   - id func(*T) *int64: 行の ID フィールド参照
//
// return:
//   - error: ID 採番・DB エラー
func assignBulkInsertZeros[T any](tx *gorm.DB, rows []T, table string, id func(*T) *int64) error {
	hasZero := false
	batchMax := int64(0)
	for i := range rows {
		v := *id(&rows[i])
		if v == 0 {
			hasZero = true
		} else if v > batchMax {
			batchMax = v
		}
	}
	if !hasZero {
		return nil
	}
	dbMax, err := postgresTableCoalesceMaxID(tx, table)
	if err != nil {
		return err
	}
	next := batchMax
	if dbMax > next {
		next = dbMax
	}
	for i := range rows {
		if *id(&rows[i]) == 0 {
			next++
			*id(&rows[i]) = next
		}
	}
	return nil
}

// syncChildrenByKey は親 FK 配下の子行をキーで突き合わせ INSERT / UPDATE / DELETE する。
// want は in-place で主キーが埋め戻される。softDelete=true は論理削除、false は物理削除。
//
// args:
//   - tx *gorm.DB: トランザクション
//   - table string: 子テーブル名
//   - parentColumn string: 親 FK 列名
//   - parentID int64: 親 ID
//   - want []T: 望ましい子行の集合
//   - keyFn func(*T) K: 自然キー抽出
//   - pkFn func(*T) *int64: 主キー参照
//   - applyUpdate func(tx *gorm.DB, prev T, next *T) error: 更新時コールバック（nil 可）
//   - softDelete bool: 論理削除かどうか
//
// return:
//   - []int64: 削除した子行の ID
//   - error: DB エラー
func syncChildrenByKey[T any, K comparable](
	tx *gorm.DB,
	table string,
	parentColumn string,
	parentID int64,
	want []T,
	keyFn func(*T) K,
	pkFn func(*T) *int64,
	applyUpdate func(tx *gorm.DB, prev T, next *T) error,
	softDelete bool,
) (deletedIDs []int64, err error) {
	var existing []T
	if err := tx.Where(parentColumn+" = ?", parentID).Find(&existing).Error; err != nil {
		return nil, err
	}

	var zeroK K
	byKey := make(map[K]T, len(existing))
	for _, row := range existing {
		byKey[keyFn(&row)] = row
	}

	var insertIdx []int
	for i := range want {
		k := keyFn(&want[i])
		if k == zeroK {
			insertIdx = append(insertIdx, i)
			continue
		}
		prev, ok := byKey[k]
		if !ok {
			insertIdx = append(insertIdx, i)
			continue
		}
		*pkFn(&want[i]) = *pkFn(&prev)
		if applyUpdate != nil {
			if err := applyUpdate(tx, prev, &want[i]); err != nil {
				return nil, err
			}
		}
		delete(byKey, k)
	}

	var toDelete []int64
	for _, row := range byKey {
		toDelete = append(toDelete, *pkFn(&row))
	}
	if len(toDelete) > 0 {
		var model T
		if softDelete {
			if err := tx.Where("id IN ?", toDelete).Delete(&model).Error; err != nil {
				return nil, err
			}
		} else {
			if err := tx.Unscoped().Where("id IN ?", toDelete).Delete(&model).Error; err != nil {
				return nil, err
			}
		}
		deletedIDs = append(deletedIDs, toDelete...)
	}

	if len(insertIdx) == 0 {
		return deletedIDs, nil
	}

	toInsert := make([]T, len(insertIdx))
	for j, i := range insertIdx {
		toInsert[j] = want[i]
	}
	if err := assignBulkInsertZeros(tx, toInsert, table, pkFn); err != nil {
		return deletedIDs, err
	}
	if err := tx.Create(&toInsert).Error; err != nil {
		return deletedIDs, err
	}
	for j, i := range insertIdx {
		*pkFn(&want[i]) = *pkFn(&toInsert[j])
	}
	return deletedIDs, nil
}

// UpdateQuestionInTransaction は質問本体と1対多関連を差分同期で1トランザクション更新する。
// related_questions のみ物理削除（ユニーク制約のため）。SubQuestions 等は永続化しない。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - q entity.Question: 更新後の質問エンティティ
//
// return:
//   - error: DB エラー
func UpdateQuestionInTransaction(ctx context.Context, db *gorm.DB, q entity.Question) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1) entity.Support: 新規作成 or 既存更新、またはフォームに無ければ既存の support を detach
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
				if err := tx.Model(&s).Omit("User", "entity.SupportStatus").Updates(&s).Error; err != nil {
					return err
				}
				if q.SupportID == nil {
					sid := s.ID
					q.SupportID = &sid
					q.Support = &s
				}
			}
		} else {
			if err := DetachQuestionSupportTx(tx, q.ID); err != nil {
				return err
			}
			q.SupportID = nil
		}
		// 2) 親の questions: スカラ列と FK のみ（CreatedAt は更新で変えない）
		if _, err := gorm.G[entity.Question](tx).
			Omit("Answer", "Memos", "Notices", "TagManagers", "entity.Support", "RelatedQuestions", "SenderTalks", "TalkroomID", "CreatedAt").
			Where("id = ?", q.ID).
			Updates(ctx, q); err != nil {
			return err
		}
		// 3) タグ紐づけ（tag_managers）— 自然キー tag_id で差分同期、削除は論理削除
		var tagRows []entity.TagManager
		for _, tm := range q.TagManagers {
			if tm.TagID == 0 {
				continue
			}
			tagRows = append(tagRows, entity.TagManager{
				ID:         tm.ID,
				TagID:      tm.TagID,
				QuestionID: q.ID,
			})
		}
		if _, err := syncChildrenByKey(tx, "tag_managers", "question_id", q.ID, tagRows,
			func(t *entity.TagManager) int64 { return t.TagID },
			func(t *entity.TagManager) *int64 { return &t.ID },
			nil,
			true,
		); err != nil {
			return err
		}
		// 4) 回答（answers）— ID で差分同期、削除は論理削除。refer_managers は各回答で refer_id 自然キーで差分・論理削除。
		var answerRows []entity.Answer
		var referRowsPerAnswer [][]entity.ReferManager
		for _, a := range q.Answer {
			content := strings.TrimSpace(a.Content)
			if a.UserID == 0 || content == "" {
				continue
			}
			answerRows = append(answerRows, entity.Answer{
				ID:         a.ID,
				UserID:     a.UserID,
				Content:    content,
				IsFinal:    a.IsFinal,
				QuestionID: q.ID,
			})
			var refs []entity.ReferManager
			for _, rm := range a.ReferManagers {
				if rm.ReferID == 0 {
					continue
				}
				refs = append(refs, entity.ReferManager{ReferID: rm.ReferID})
			}
			referRowsPerAnswer = append(referRowsPerAnswer, refs)
		}
		deletedAnswerIDs, err := syncChildrenByKey(tx, "answers", "question_id", q.ID, answerRows,
			func(a *entity.Answer) int64 { return a.ID },
			func(a *entity.Answer) *int64 { return &a.ID },
			func(tx *gorm.DB, prev entity.Answer, next *entity.Answer) error {
				return tx.Model(&entity.Answer{}).Where("id = ?", prev.ID).
					Updates(map[string]any{
						"content":  next.Content,
						"is_final": next.IsFinal,
						"user_id":  next.UserID,
					}).Error
			},
			true,
		)
		if err != nil {
			return err
		}
		if len(deletedAnswerIDs) > 0 {
			if err := tx.Where("answer_id IN ?", deletedAnswerIDs).Delete(&entity.ReferManager{}).Error; err != nil {
				return err
			}
		}
		for i := range answerRows {
			refs := referRowsPerAnswer[i]
			for j := range refs {
				refs[j].AnswerID = answerRows[i].ID
			}
			if _, err := syncChildrenByKey(tx, "refer_managers", "answer_id", answerRows[i].ID, refs,
				func(r *entity.ReferManager) int64 { return r.ReferID },
				func(r *entity.ReferManager) *int64 { return &r.ID },
				nil,
				true,
			); err != nil {
				return err
			}
		}
		// 5) メモ（memos）— ID で差分同期、削除は論理削除
		var memoRows []entity.Memo
		for _, m := range q.Memos {
			content := strings.TrimSpace(m.Content)
			if m.UserID == 0 || content == "" {
				continue
			}
			memoRows = append(memoRows, entity.Memo{
				ID:         m.ID,
				UserID:     m.UserID,
				Content:    content,
				QuestionID: q.ID,
			})
		}
		if _, err := syncChildrenByKey(tx, "memos", "question_id", q.ID, memoRows,
			func(m *entity.Memo) int64 { return m.ID },
			func(m *entity.Memo) *int64 { return &m.ID },
			func(tx *gorm.DB, prev entity.Memo, next *entity.Memo) error {
				return tx.Model(&entity.Memo{}).Where("id = ?", prev.ID).
					Updates(map[string]any{"content": next.Content, "user_id": next.UserID}).Error
			},
			true,
		); err != nil {
			return err
		}
		// 6) 関連質問（related_questions）— 自然キー related_question_id で差分同期。
		// doc/db/INIT.sql の uq_related_questions UNIQUE(question_id, related_question_id) により、
		// 論理削除（deleted_at をセット）してもユニーク制約が deleted_at を区別しないため、
		// 同じ自然キーで再追加すると衝突する。そのためここだけ物理削除（Unscoped）で同期する。
		var relatedRows []entity.RelatedQuestion
		for _, rq := range q.RelatedQuestions {
			if rq.RelatedQuestionID == 0 || rq.RelatedQuestionID == q.ID {
				continue
			}
			relatedRows = append(relatedRows, entity.RelatedQuestion{
				ID:                rq.ID,
				QuestionID:        q.ID,
				RelatedQuestionID: rq.RelatedQuestionID,
			})
		}
		if _, err := syncChildrenByKey(tx, "related_questions", "question_id", q.ID, relatedRows,
			func(r *entity.RelatedQuestion) int64 { return r.RelatedQuestionID },
			func(r *entity.RelatedQuestion) *int64 { return &r.ID },
			nil,
			false,
		); err != nil {
			return err
		}
		return nil
	})
}

// DetachQuestionSupportTx は質問の support_id を NULL にし、紐づく supports 行を削除する（1:1 前提）。
//
// args:
//   - tx *gorm.DB: トランザクション
//   - questionID int64: 質問 ID
//
// return:
//   - error: DB エラー
func DetachQuestionSupportTx(tx *gorm.DB, questionID int64) error {
	var current entity.Question
	if err := tx.Select("id", "support_id").
		Where("id = ?", questionID).
		Take(&current).Error; err != nil {
		return err
	}
	if err := tx.Model(&entity.Question{}).
		Where("id = ?", questionID).
		Update("support_id", nil).Error; err != nil {
		return err
	}
	if current.SupportID != nil && *current.SupportID != 0 {
		if err := tx.Unscoped().
			Where("id = ?", *current.SupportID).
			Delete(&entity.Support{}).Error; err != nil {
			return err
		}
	}
	return nil
}

// DeleteQuestionByID は指定 ID の質問と関連行を削除する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - id int64: 質問 ID
//
// return:
//   - error: DB エラー
func DeleteQuestionByID(ctx context.Context, db *gorm.DB, id int64) error {
	if err := db.WithContext(ctx).
		Where("question_id = ? OR related_question_id = ?", id, id).
		Delete(&entity.RelatedQuestion{}).Error; err != nil {
		return err
	}
	if _, err := gorm.G[entity.Question](db).
		Preload("Answer", commonPreloadBuilder()).
		Preload("Answer.User", commonPreloadBuilder()).
		Preload("Answer.User.Role", commonPreloadBuilder()).
		Preload("Answer.ReferManagers", commonPreloadBuilder()).
		Preload("Answer.ReferManagers.Refer", commonPreloadBuilder()).
		Preload("Memos", commonPreloadBuilder()).
		Preload("Memos.User", commonPreloadBuilder()).
		Preload("Memos.User.Role", commonPreloadBuilder()).
		Preload("TagManagers", commonPreloadBuilder()).
		Preload("TagManagers.Tag", commonPreloadBuilder()).
		Preload("TagManagers.Tag.Category", commonPreloadBuilder()).
		Preload("Support", commonPreloadBuilder()).
		Preload("Support.User", commonPreloadBuilder()).
		Preload("Support.User.Role", commonPreloadBuilder()).
		Preload("Support.SupportStatus", commonPreloadBuilder()).
		Where("id = ?", id).
		Limit(1).
		Delete(ctx); err != nil {
		return err
	}
	return nil
}

// DeleteUserByID は指定 ID のユーザーを削除する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - id int64: ユーザー ID
//
// return:
//   - error: DB エラー
func DeleteUserByID(ctx context.Context, db *gorm.DB, id int64) error {
	if _, err := gorm.G[entity.User](db).
		Preload("Role", commonPreloadBuilder()).
		Where("id = ?", id).
		Limit(1).
		Delete(ctx); err != nil {
		return err
	}
	return nil
}

// GetQuestion は指定 ID の質問を関連込みで取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - id int64: 質問 ID
//
// return:
//   - entity.Question: 質問
//   - error: DB エラー
func GetQuestion(ctx context.Context, db *gorm.DB, id int64) (model entity.Question, err error) {
	model, err = gorm.G[entity.Question](db).
		Preload("Answer", commonPreloadBuilder()).
		Preload("Answer.User", userPreloadBuilder(false)).
		Preload("Answer.User.Role", commonPreloadBuilder()).
		Preload("Answer.ReferManagers", commonPreloadBuilder()).
		Preload("Answer.ReferManagers.Refer", commonPreloadBuilder()).
		Preload("Memos", commonPreloadBuilder()).
		Preload("Memos.User", userPreloadBuilder(false)).
		Preload("Memos.User.Role", commonPreloadBuilder()).
		Preload("TagManagers", commonPreloadBuilder()).
		Preload("TagManagers.Tag", commonPreloadBuilder()).
		Preload("TagManagers.Tag.Category", commonPreloadBuilder()).
		Preload("RelatedQuestions", commonPreloadBuilder()).
		Preload("RelatedQuestions.RelatedQuestion", commonPreloadBuilder()).
		Preload("SenderTalks", commonPreloadBuilder()).
		Preload("SenderTalks.Sender", commonPreloadBuilder()).
		Preload("Support", commonPreloadBuilder()).
		Preload("Support.User", userPreloadBuilder(false)).
		Preload("Support.User.Role", commonPreloadBuilder()).
		Preload("Support.SupportStatus", commonPreloadBuilder()).
		Where("id = ?", id).
		First(ctx)
	return
}

// GetQuestions は質問一覧を関連込みで取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//
// return:
//   - []entity.Question: 質問一覧
//   - error: DB エラー
func GetQuestions(ctx context.Context, db *gorm.DB) (models []entity.Question, err error) {
	models, err = gorm.G[entity.Question](db).
		Preload("Answer", commonPreloadBuilder()).
		Preload("Answer.User", userPreloadBuilder(false)).
		Preload("Answer.User.Role", commonPreloadBuilder()).
		Preload("Answer.ReferManagers", commonPreloadBuilder()).
		Preload("Answer.ReferManagers.Refer", commonPreloadBuilder()).
		Preload("Memos", commonPreloadBuilder()).
		Preload("Memos.User", userPreloadBuilder(false)).
		Preload("Memos.User.Role", commonPreloadBuilder()).
		Preload("TagManagers", commonPreloadBuilder()).
		Preload("TagManagers.Tag", commonPreloadBuilder()).
		Preload("TagManagers.Tag.Category", commonPreloadBuilder()).
		Preload("RelatedQuestions", commonPreloadBuilder()).
		Preload("RelatedQuestions.RelatedQuestion", commonPreloadBuilder()).
		Preload("SenderTalks", commonPreloadBuilder()).
		Preload("SenderTalks.Sender", commonPreloadBuilder()).
		Preload("Support", commonPreloadBuilder()).
		Preload("Support.User", userPreloadBuilder(false)).
		Preload("Support.User.Role", commonPreloadBuilder()).
		Preload("Support.SupportStatus", commonPreloadBuilder()).
		Order("id").
		Find(ctx)
	return
}

// GetTags はタグ一覧をカテゴリ込みで取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//
// return:
//   - []entity.Tag: タグ一覧
//   - error: DB エラー
func GetTags(ctx context.Context, db *gorm.DB) (models []entity.Tag, err error) {
	models, err = gorm.G[entity.Tag](db).
		Preload("Category", commonPreloadBuilder()).
		Order("id").
		Find(ctx)
	return
}

// GetTagByID は指定 ID のタグをカテゴリ込みで取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - id int64: タグ ID
//
// return:
//   - entity.Tag: タグ
//   - error: DB エラー
func GetTagByID(ctx context.Context, db *gorm.DB, id int64) (models entity.Tag, err error) {
	models, err = gorm.G[entity.Tag](db).
		Preload("Category", commonPreloadBuilder()).
		Where("id = ?", id).
		First(ctx)
	return
}

// DeleteTagByID は指定 ID のタグを削除する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - id int64: タグ ID
//
// return:
//   - error: DB エラー
func DeleteTagByID(ctx context.Context, db *gorm.DB, id int64) error {
	if _, err := gorm.G[entity.Tag](db).Where("id = ?", id).Limit(1).Delete(ctx); err != nil {
		return err
	}
	return nil
}

// GetNoticeByQuestionIDs は質問 ID 一覧に紐づく通知を取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - questionIDs []int64: 質問 ID 一覧
//
// return:
//   - []entity.Notice: 通知一覧
//   - error: DB エラー
func GetNoticeByQuestionIDs(ctx context.Context, db *gorm.DB, questionIDs []int64) (models []entity.Notice, err error) {
	if len(questionIDs) > 0 {
		models, err = gorm.G[entity.Notice](db).Where("question_id IN ?", questionIDs).Order("id").Find(ctx)
	}
	return
}

// GetNoticeByQuestion は質問に紐づく通知を1件取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - question entity.Question: 対象質問
//
// return:
//   - entity.Notice: 通知
//   - error: DB エラー
func GetNoticeByQuestion(ctx context.Context, db *gorm.DB, question entity.Question) (models entity.Notice, err error) {
	models, err = gorm.G[entity.Notice](db).Where("question_id = ?", question.ID).First(ctx)
	return
}

// GetNotice は通知一覧を関連込みで取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//
// return:
//   - []entity.Notice: 通知一覧
//   - error: DB エラー
func GetNotice(ctx context.Context, db *gorm.DB) (models []entity.Notice, err error) {
	models, err = gorm.G[entity.Notice](db).
		Preload("NoticeType", commonPreloadBuilder()).
		Preload("Question", commonPreloadBuilder()).
		Preload("Question.Support", commonPreloadBuilder()).
		Preload("Question.Support.SupportStatus", commonPreloadBuilder()).
		Preload("Question.TagManagers", commonPreloadBuilder()).
		Preload("Question.TagManagers.Tag", commonPreloadBuilder()).
		Order("id").
		Find(ctx)
	return
}

// GetNoticeByQuestionSilent は質問に紐づく通知をサイレントログで取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - question entity.Question: 対象質問
//
// return:
//   - entity.Notice: 通知
//   - error: DB エラー
func GetNoticeByQuestionSilent(ctx context.Context, db *gorm.DB, question entity.Question) (model entity.Notice, err error) {
	model, err = gorm.G[entity.Notice](db.Session(&gorm.Session{
		Logger: db.Logger.LogMode(logger.Silent),
	})).Where("question_id = ?", question.ID).First(ctx)
	return
}

// RegisterNoticeByQuestionID は回答期日接近の通知を質問 ID に紐づけて登録する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - questionID int64: 質問 ID
//
// return:
//   - error: DB エラー
func RegisterNoticeByQuestionID(ctx context.Context, db *gorm.DB, questionID int64) error {
	var content = "質問の回答期日が近づいています。"
	notice := entity.Notice{
		TypeID:     3,
		QuestionID: &questionID,
		Content:    &content,
	}
	return gorm.G[entity.Notice](db).Create(ctx, &notice)
}

// DeleteNoticeByID は指定 ID の通知を削除する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - id int64: 通知 ID
//
// return:
//   - error: DB エラー
func DeleteNoticeByID(ctx context.Context, db *gorm.DB, id int64) error {
	if _, err := gorm.G[entity.Notice](db).Where("id = ?", id).Delete(ctx); err != nil {
		return err
	}
	return nil
}

// DeleteNoticeByQuestion は質問に紐づく通知を削除する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - question entity.Question: 対象質問
//
// return:
//   - int64: 削除した通知 ID
//   - error: DB エラー
func DeleteNoticeByQuestion(ctx context.Context, db *gorm.DB, question entity.Question) (noticeID int64, err error) {
	n, err := GetNoticeByQuestionSilent(ctx, db, question)
	if err != nil {
		return -1, err
	}
	if err := DeleteNoticeByID(ctx, db, n.ID); err != nil {
		return -1, err
	}
	return n.ID, nil
}

// DeleteNoticeByQuestionID は質問 ID に紐づく通知を検索して削除する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - questionID int64: 質問 ID
//
// return:
//   - int64: 削除した通知 ID（見つからない場合は -1）
//   - error: DB エラー
func DeleteNoticeByQuestionID(ctx context.Context, db *gorm.DB, questionID int64) (deletedID int64, err error) {
	notices, err := GetNotice(ctx, db)
	if err != nil {
		return -1, err
	}
	for _, n := range notices {
		if n.QuestionID != nil && *n.QuestionID == questionID {
			if err := DeleteNoticeByID(ctx, db, n.ID); err != nil {
				return n.ID, err
			}
			return n.ID, nil
		}
	}
	return -1, gorm.ErrRecordNotFound
}

// GetMaxByColumn は指定列の MAX 値を返す（無効・未存在時は -1）。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - columnName string: 集計列名
//
// return:
//   - int64: 最大値
func GetMaxByColumn[T any](ctx context.Context, db *gorm.DB, columnName string) int64 {
	var max sql.NullInt64
	err := db.WithContext(ctx).Model(new(T)).
		Select(fmt.Sprintf("MAX(%s)", columnName)).
		Take(&max).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return -1
	}

	if !max.Valid {
		return -1
	}
	return max.Int64

}

// commonPreloadBuilder は Preload 時に id 昇順で並べるビルダを返す。
//
// return:
//   - func(db gorm.PreloadBuilder) error: Preload コールバック
func commonPreloadBuilder() func(db gorm.PreloadBuilder) error {
	return func(db gorm.PreloadBuilder) error {
		db.Order("id")
		return nil
	}
}

// userPreloadBuilder はユーザー Preload 用の Select・Order を返す。
//
// args:
//   - includePassword bool: password 列を含めるか
//
// return:
//   - func(db gorm.PreloadBuilder) error: Preload コールバック
func userPreloadBuilder(includePassword bool) func(db gorm.PreloadBuilder) error {
	return func(db gorm.PreloadBuilder) error {
		if includePassword {
			db.Select("id, name, email, password, memo, role_id")
		} else {
			db.Select("id, name, email, memo, role_id")
		}
		db.Order("id ASC")
		return nil
	}
}
