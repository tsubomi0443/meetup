package mapper

import (
	"meetup/internal/domains/entity"
	"meetup/internal/usecases/dto"
)

// =====================
// カテゴリ（entity.Category）

// =====================
// categoryFromEntityShallow は entity.Category を Tags なしの dto.CategoryForm に変換する。
//
// args:
//   - e entity.Category: 変換元エンティティ
//
// return:
//   - dto.CategoryForm: カテゴリフォーム DTO（Tags なし）
func categoryFromEntityShallow(e entity.Category) dto.CategoryForm {
	return dto.CategoryForm{
		ID:        e.ID,
		Name:      e.Name,
		CreatedAt: timeToISO(e.CreatedAt),
		UpdatedAt: timeToISO(e.UpdatedAt),
		DeletedAt: deletedAtToISO(e.DeletedAt),
	}
}

// CategoryFromEntity は entity.Category を dto.CategoryForm に変換する。
//
// args:
//   - e entity.Category: 変換元エンティティ
//
// return:
//   - dto.CategoryForm: カテゴリフォーム DTO
func CategoryFromEntity(e entity.Category) dto.CategoryForm {
	f := categoryFromEntityShallow(e)
	for _, t := range e.Tags {
		f.Tags = append(f.Tags, TagFromEntity(t))
	}
	return f
}

// CategoryToEntity は dto.CategoryForm を entity.Category に変換する。
//
// args:
//   - f dto.CategoryForm: 変換元フォーム DTO
//
// return:
//   - entity.Category: カテゴリエンティティ
func CategoryToEntity(f dto.CategoryForm) entity.Category {
	e := entity.Category{
		ID:   f.ID,
		Name: f.Name,
	}
	for _, tf := range f.Tags {
		e.Tags = append(e.Tags, TagToEntity(tf))
	}
	return e
}
