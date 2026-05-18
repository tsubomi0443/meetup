package question

import (
	"context"

	"meetup/internal/domains/entity"
	domainquestion "meetup/internal/domains/question"
	"meetup/internal/usecases/mapper"
	"meetup/internal/usecases/dto"
)

// UseCase は質問（問い合わせ）管理ユースケースを表す。
type UseCase struct {
	questions domainquestion.Repository
}

// NewUseCase は質問管理ユースケースを生成する。
//
// args:
//   - questions domainquestion.Repository: 質問リポジトリ
//
// return:
//   - *UseCase: 生成したユースケース
func NewUseCase(questions domainquestion.Repository) *UseCase {
	return &UseCase{questions: questions}
}

// GetAll は全質問を DTO 一覧として取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//
// return:
//   - []dto.QuestionForm: 質問フォームの一覧
//   - error: 取得エラー
func (u *UseCase) GetAll(ctx context.Context) ([]dto.QuestionForm, error) {
	models, err := u.questions.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return mapper.QuestionFromEntities(models), nil
}

// GetByID は指定 ID の質問を取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - id int64: 質問 ID
//
// return:
//   - dto.QuestionForm: 質問フォーム
//   - error: 取得エラー
func (u *UseCase) GetByID(ctx context.Context, id int64) (dto.QuestionForm, error) {
	model, err := u.questions.GetByID(ctx, id)
	if err != nil {
		return dto.QuestionForm{}, err
	}
	return mapper.QuestionFromEntity(model), nil
}

// Register は新規質問を登録する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - form dto.QuestionForm: 登録内容
//
// return:
//   - dto.QuestionForm: 登録後の質問フォーム
//   - error: 登録・取得エラー
func (u *UseCase) Register(ctx context.Context, form dto.QuestionForm) (dto.QuestionForm, error) {
	data := mapper.QuestionToEntity(form)
	if err := u.questions.Register(ctx, &data); err != nil {
		return dto.QuestionForm{}, err
	}
	created, err := u.questions.GetByID(ctx, data.ID)
	if err != nil {
		return dto.QuestionForm{}, err
	}
	return mapper.QuestionFromEntity(created), nil
}

// Update は質問を更新する。未割当時の Support クリアや対応中の担当者自動設定を正規化してから永続化する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - form dto.QuestionForm: 更新内容
//   - actorUserID int64: 操作ユーザ ID（担当者自動設定用）
//   - hasActor bool: 操作ユーザが特定できているか
//
// return:
//   - entity.Question: 更新後のエンティティ（トランザクション反映後）
//   - dto.QuestionForm: 再読込した質問フォーム
//   - error: 更新・取得エラー
func (u *UseCase) Update(ctx context.Context, form dto.QuestionForm, actorUserID int64, hasActor bool) (entity.Question, dto.QuestionForm, error) {
	NormalizeQuestionFormClearSupportWhenUnassigned(&form)
	if hasActor {
		NormalizeQuestionFormAssignSupportUserWhenInProgress(&form, actorUserID)
	}
	updatedModel := mapper.QuestionToEntity(form)
	if err := u.questions.UpdateInTransaction(ctx, updatedModel); err != nil {
		return entity.Question{}, dto.QuestionForm{}, err
	}
	loaded, err := u.questions.GetByID(ctx, updatedModel.ID)
	if err != nil {
		return entity.Question{}, dto.QuestionForm{}, err
	}
	return updatedModel, mapper.QuestionFromEntity(loaded), nil
}

// DeleteByID は指定 ID の質問を削除する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - id int64: 質問 ID
//
// return:
//   - error: 削除エラー
func (u *UseCase) DeleteByID(ctx context.Context, id int64) error {
	return u.questions.DeleteByID(ctx, id)
}
