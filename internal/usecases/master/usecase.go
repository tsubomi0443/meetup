package master

import (
	"context"

	domainmaster "meetup/internal/domains/master"
	"meetup/internal/usecases/mapper"
	"meetup/internal/usecases/dto"
)

// UseCase はマスタデータ（ロール・カテゴリ・支援ステータス）取得ユースケースを表す。
type UseCase struct {
	master domainmaster.Repository
}

// NewUseCase はマスタデータユースケースを生成する。
//
// args:
//   - master domainmaster.Repository: マスタリポジトリ
//
// return:
//   - *UseCase: 生成したユースケース
func NewUseCase(master domainmaster.Repository) *UseCase {
	return &UseCase{master: master}
}

// ListRoles は全ロールを DTO 一覧として取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//
// return:
//   - []dto.RoleForm: ロールフォームの一覧
//   - error: 取得エラー
func (u *UseCase) ListRoles(ctx context.Context) ([]dto.RoleForm, error) {
	roles, err := u.master.GetRoles(ctx)
	if err != nil {
		return nil, err
	}
	roleForms := make([]dto.RoleForm, 0, len(roles))
	for i := range roles {
		roleForms = append(roleForms, mapper.RoleFromEntity(roles[i]))
	}
	return roleForms, nil
}

// ListCategories は全カテゴリを DTO 一覧として取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//
// return:
//   - []dto.CategoryForm: カテゴリフォームの一覧
//   - error: 取得エラー
func (u *UseCase) ListCategories(ctx context.Context) ([]dto.CategoryForm, error) {
	categories, err := u.master.GetCategories(ctx)
	if err != nil {
		return nil, err
	}
	categoryForms := make([]dto.CategoryForm, 0, len(categories))
	for i := range categories {
		categoryForms = append(categoryForms, mapper.CategoryFromEntity(categories[i]))
	}
	return categoryForms, nil
}

// ListSupportStatuses は全支援ステータスを DTO 一覧として取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//
// return:
//   - []dto.SupportStatusForm: 支援ステータスフォームの一覧
//   - error: 取得エラー
func (u *UseCase) ListSupportStatuses(ctx context.Context) ([]dto.SupportStatusForm, error) {
	statuses, err := u.master.GetSupportStatuses(ctx)
	if err != nil {
		return nil, err
	}
	supportStatusForms := make([]dto.SupportStatusForm, 0, len(statuses))
	for i := range statuses {
		supportStatusForms = append(supportStatusForms, mapper.SupportStatusFromEntity(statuses[i]))
	}
	return supportStatusForms, nil
}
