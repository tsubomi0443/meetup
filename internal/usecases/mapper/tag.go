package mapper

import (
	"strconv"

	"meetup/internal/domains/entity"
	"meetup/internal/usecases/dto"
)

// =====================
// タグ（entity.Tag）

// =====================
// tagFromEntityShallow は entity.Tag を Questions なしの dto.TagForm に変換する。Category は浅い変換のみ。
//
// args:
//   - e entity.Tag: 変換元エンティティ
//
// return:
//   - dto.TagForm: タグフォーム DTO
func tagFromEntityShallow(e entity.Tag) dto.TagForm {
	f := dto.TagForm{
		ID:         e.ID,
		Name:       e.Name,
		Usage:      e.Usage,
		CategoryID: strconv.FormatInt(e.CategoryID, 10),
		CreatedAt:  timeToISO(e.CreatedAt),
		UpdatedAt:  timeToISO(e.UpdatedAt),
		DeletedAt:  deletedAtToISO(e.DeletedAt),
	}
	if e.Category.ID != 0 {
		c := categoryFromEntityShallow(e.Category)
		f.Category = &c
	}
	return f
}

// TagFromEntities は entity.Tag のスライスを dto.TagForm のスライスに一括変換する。
//
// args:
//   - e []entity.Tag: 変換元エンティティ一覧
//
// return:
//   - []dto.TagForm: タグフォーム DTO の一覧
func TagFromEntities(e []entity.Tag) []dto.TagForm {
	forms := []dto.TagForm{}
	for _, tag := range e {
		forms = append(forms, TagFromEntity(tag))
	}
	return forms
}

// TagFromEntity は entity.Tag を dto.TagForm に変換する。TagManagers 経由で関連質問も含める。
//
// args:
//   - e entity.Tag: 変換元エンティティ
//
// return:
//   - dto.TagForm: タグフォーム DTO
func TagFromEntity(e entity.Tag) dto.TagForm {
	f := tagFromEntityShallow(e)
	for _, tm := range e.TagManagers {
		if tm.Question.ID != 0 {
			f.Questions = append(f.Questions, QuestionFromEntity(tm.Question))
		}
	}
	return f
}

// TagToEntity は dto.TagForm を entity.Tag に変換する。Category の明示的関連はセットしない。
//
// args:
//   - f dto.TagForm: 変換元フォーム DTO
//
// return:
//   - entity.Tag: タグエンティティ
func TagToEntity(f dto.TagForm) entity.Tag {
	e := entity.Tag{
		ID:         f.ID,
		Name:       f.Name,
		Usage:      f.Usage,
		CategoryID: f.CategoryIDInt64(),
	}
	// DB に余分な category が入らないよう、明示的関連はセットしない。
	// if f.Category != nil {
	// 	e.Category = CategoryToEntity(*f.Category)
	// }
	for _, qf := range f.Questions {
		tm := entity.TagManager{
			TagID:      f.ID,
			QuestionID: qf.ID,
		}
		if qf.ID != 0 {
			tm.Question = entity.Question{ID: qf.ID}
		}
		e.TagManagers = append(e.TagManagers, tm)
	}
	return e
}

// TagToEntityNoRelations は dto.TagForm を関連なしの entity.Tag に変換する。
//
// args:
//   - f dto.TagForm: 変換元フォーム DTO
//
// return:
//   - entity.Tag: タグエンティティ（Category は空構造体）
func TagToEntityNoRelations(f dto.TagForm) entity.Tag {
	e := entity.Tag{
		ID:         f.ID,
		Name:       f.Name,
		Usage:      f.Usage,
		CategoryID: f.CategoryIDInt64(),
	}
	e.Category = entity.Category{}

	return e
}

// =====================
// 参照リンク管理・タグ管理（entity.ReferManager / entity.TagManager）
