package tag

import (
	"context"

	domaintag "meetup/internal/domains/tag"
	"meetup/internal/usecases/mapper"
	"meetup/internal/usecases/dto"
)

// UseCase はタグ管理ユースケースを表す。
type UseCase struct {
	tags domaintag.Repository
}

// NewUseCase はタグ管理ユースケースを生成する。
//
// args:
//   - tags domaintag.Repository: タグリポジトリ
//
// return:
//   - *UseCase: 生成したユースケース
func NewUseCase(tags domaintag.Repository) *UseCase {
	return &UseCase{tags: tags}
}

// GetAll は全タグを DTO 一覧として取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//
// return:
//   - []dto.TagForm: タグフォームの一覧
//   - error: 取得エラー
func (u *UseCase) GetAll(ctx context.Context) ([]dto.TagForm, error) {
	models, err := u.tags.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return mapper.TagFromEntities(models), nil
}

// Register は新規タグを登録する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - form dto.TagForm: 登録内容
//
// return:
//   - dto.TagForm: 登録後のタグフォーム
//   - error: 登録・取得エラー
func (u *UseCase) Register(ctx context.Context, form dto.TagForm) (dto.TagForm, error) {
	model := mapper.TagToEntity(form)
	if err := u.tags.Register(ctx, &model); err != nil {
		return dto.TagForm{}, err
	}
	loaded, err := u.tags.GetByID(ctx, model.ID)
	if err != nil {
		return dto.TagForm{}, err
	}
	return mapper.TagFromEntity(loaded), nil
}

// Update は既存タグを更新する（Category・TagManagers は更新対象外）。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - form dto.TagForm: 更新内容
//
// return:
//   - error: 更新エラー
func (u *UseCase) Update(ctx context.Context, form dto.TagForm) error {
	model := mapper.TagToEntityNoRelations(form)
	_, err := u.tags.UpdateByID(ctx, model.ID, model, "Category", "TagManagers")
	return err
}

// DeleteByID は指定 ID のタグを削除する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - id int64: タグ ID
//
// return:
//   - error: 削除エラー
func (u *UseCase) DeleteByID(ctx context.Context, id int64) error {
	return u.tags.DeleteByID(ctx, id)
}
