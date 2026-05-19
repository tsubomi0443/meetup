package postgres

import (
	"context"

	"meetup/internal/domains/entity"

	"gorm.io/gorm"
)

// QuestionRepository は質問の永続化を担う。
type QuestionRepository struct {
	DB *gorm.DB
}

// NewQuestionRepository は QuestionRepository を生成する。
//
// args:
//   - db *gorm.DB: データベース接続
//
// return:
//   - *QuestionRepository: リポジトリ
func NewQuestionRepository(db *gorm.DB) *QuestionRepository {
	return &QuestionRepository{DB: db}
}

// GetByID は指定 ID の質問を関連込みで取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - id int64: 質問 ID
//
// return:
//   - entity.Question: 質問
//   - error: DB エラー
func (r *QuestionRepository) GetByID(ctx context.Context, id int64) (entity.Question, error) {
	return GetQuestion(ctx, r.DB, id)
}

// GetAll は質問一覧を関連込みで取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//
// return:
//   - []entity.Question: 質問一覧
//   - error: DB エラー
func (r *QuestionRepository) GetAll(ctx context.Context) ([]entity.Question, error) {
	return GetQuestions(ctx, r.DB)
}

// Register は質問を新規登録する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - model *entity.Question: 登録する質問
//
// return:
//   - error: DB エラー
func (r *QuestionRepository) Register(ctx context.Context, model *entity.Question) error {
	return Register(ctx, r.DB, model)
}

// UpdateInTransaction は質問と子関連を1トランザクションで更新する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - q entity.Question: 更新後の質問エンティティ
//
// return:
//   - error: DB エラー
func (r *QuestionRepository) UpdateInTransaction(ctx context.Context, q entity.Question) error {
	return UpdateQuestionInTransaction(ctx, r.DB, q)
}

// DeleteByID は指定 ID の質問を削除する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - id int64: 質問 ID
//
// return:
//   - error: DB エラー
func (r *QuestionRepository) DeleteByID(ctx context.Context, id int64) error {
	return DeleteQuestionByID(ctx, r.DB, id)
}
