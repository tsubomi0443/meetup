package notice

import (
	"context"

	domainnotice "meetup/internal/domains/notice"
	"meetup/internal/usecases/mapper"
	"meetup/internal/usecases/dto"
)

// UseCase は通知一覧取得ユースケースを表す。
type UseCase struct {
	notices domainnotice.Repository
}

// NewUseCase は通知ユースケースを生成する。
//
// args:
//   - notices domainnotice.Repository: 通知リポジトリ
//
// return:
//   - *UseCase: 生成したユースケース
func NewUseCase(notices domainnotice.Repository) *UseCase {
	return &UseCase{notices: notices}
}

// GetAll は全通知を DTO 一覧として取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//
// return:
//   - []dto.NoticeForm: 通知フォームの一覧
//   - error: 取得エラー
func (u *UseCase) GetAll(ctx context.Context) ([]dto.NoticeForm, error) {
	models, err := u.notices.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return mapper.NoticeFromEntities(models), nil
}
