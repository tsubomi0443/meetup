package postgres

import (
	"context"

	"meetup/internal/domains/entity"

	"gorm.io/gorm"
)

// UserRepository はユーザーの永続化を担う。
type UserRepository struct {
	DB *gorm.DB
}

// NewUserRepository は UserRepository を生成する。
//
// args:
//   - db *gorm.DB: データベース接続
//
// return:
//   - *UserRepository: リポジトリ
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// GetByID は指定 ID のユーザーを取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - id int64: ユーザー ID
//
// return:
//   - entity.User: ユーザー
//   - error: DB エラー
func (r *UserRepository) GetByID(ctx context.Context, id int64) (entity.User, error) {
	return GetUserByID(ctx, r.DB, id)
}

// GetPasswordByEmail はメールアドレスに紐づく保存パスワードハッシュを取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - email string: メールアドレス
//
// return:
//   - string: パスワードハッシュ
//   - error: DB エラー
func (r *UserRepository) GetPasswordByEmail(ctx context.Context, email string) (string, error) {
	return GetUserPasswordByEmail(ctx, r.DB, email)
}

// GetUserInfo はメールとパスワードハッシュでユーザーを取得する（Preload 指定可）。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - email string: メールアドレス
//   - pass string: 保存済みパスワードハッシュ
//   - preloads ...string: GORM Preload 名
//
// return:
//   - entity.User: ユーザー
//   - error: DB エラー
func (r *UserRepository) GetUserInfo(ctx context.Context, email, pass string, preloads ...string) (entity.User, error) {
	return GetUserInfo(ctx, r.DB, email, pass, preloads...)
}

// GetUsers は管理者以外のユーザー一覧を取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//
// return:
//   - []entity.User: ユーザー一覧
//   - error: DB エラー
func (r *UserRepository) GetUsers(ctx context.Context) ([]entity.User, error) {
	return GetUsers(ctx, r.DB)
}

// Register はユーザーを新規登録する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - model *entity.User: 登録するユーザー
//
// return:
//   - error: DB エラー
func (r *UserRepository) Register(ctx context.Context, model *entity.User) error {
	return Register(ctx, r.DB, model)
}

// UpdateByID は指定 ID のユーザーを更新する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - id int64: ユーザー ID
//   - model entity.User: 更新内容
//   - omit ...string: 更新から除外する関連・列
//
// return:
//   - int: 更新行数
//   - error: DB エラー
func (r *UserRepository) UpdateByID(ctx context.Context, id int64, model entity.User, omit ...string) (int, error) {
	return UpdateByID(ctx, r.DB, id, model, omit...)
}

// DeleteByID は指定 ID のユーザーを削除する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - id int64: ユーザー ID
//
// return:
//   - error: DB エラー
func (r *UserRepository) DeleteByID(ctx context.Context, id int64) error {
	return DeleteUserByID(ctx, r.DB, id)
}

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

// TagRepository はタグの永続化を担う。
type TagRepository struct {
	DB *gorm.DB
}

// NewTagRepository は TagRepository を生成する。
//
// args:
//   - db *gorm.DB: データベース接続
//
// return:
//   - *TagRepository: リポジトリ
func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{DB: db}
}

// GetAll はタグ一覧を取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//
// return:
//   - []entity.Tag: タグ一覧
//   - error: DB エラー
func (r *TagRepository) GetAll(ctx context.Context) ([]entity.Tag, error) {
	return GetTags(ctx, r.DB)
}

// GetByID は指定 ID のタグを取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - id int64: タグ ID
//
// return:
//   - entity.Tag: タグ
//   - error: DB エラー
func (r *TagRepository) GetByID(ctx context.Context, id int64) (entity.Tag, error) {
	return GetTagByID(ctx, r.DB, id)
}

// Register はタグを新規登録する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - model *entity.Tag: 登録するタグ
//
// return:
//   - error: DB エラー
func (r *TagRepository) Register(ctx context.Context, model *entity.Tag) error {
	return Register(ctx, r.DB, model)
}

// UpdateByID は指定 ID のタグを更新する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - id int64: タグ ID
//   - model entity.Tag: 更新内容
//   - omit ...string: 更新から除外する関連・列
//
// return:
//   - int: 更新行数
//   - error: DB エラー
func (r *TagRepository) UpdateByID(ctx context.Context, id int64, model entity.Tag, omit ...string) (int, error) {
	return UpdateByID(ctx, r.DB, id, model, omit...)
}

// DeleteByID は指定 ID のタグを削除する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - id int64: タグ ID
//
// return:
//   - error: DB エラー
func (r *TagRepository) DeleteByID(ctx context.Context, id int64) error {
	return DeleteTagByID(ctx, r.DB, id)
}

// NoticeRepository は通知の永続化を担う。
type NoticeRepository struct {
	DB *gorm.DB
}

// NewNoticeRepository は NoticeRepository を生成する。
//
// args:
//   - db *gorm.DB: データベース接続
//
// return:
//   - *NoticeRepository: リポジトリ
func NewNoticeRepository(db *gorm.DB) *NoticeRepository {
	return &NoticeRepository{DB: db}
}

// GetAll は通知一覧を関連込みで取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//
// return:
//   - []entity.Notice: 通知一覧
//   - error: DB エラー
func (r *NoticeRepository) GetAll(ctx context.Context) ([]entity.Notice, error) {
	return GetNotice(ctx, r.DB)
}

// GetByQuestionSilent は質問に紐づく通知をサイレントログで取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - question entity.Question: 対象質問
//
// return:
//   - entity.Notice: 通知
//   - error: DB エラー
func (r *NoticeRepository) GetByQuestionSilent(ctx context.Context, question entity.Question) (entity.Notice, error) {
	return GetNoticeByQuestionSilent(ctx, r.DB, question)
}

// RegisterByQuestionID は質問 ID に紐づく期限通知を登録する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - questionID int64: 質問 ID
//
// return:
//   - error: DB エラー
func (r *NoticeRepository) RegisterByQuestionID(ctx context.Context, questionID int64) error {
	return RegisterNoticeByQuestionID(ctx, r.DB, questionID)
}

// DeleteByQuestionID は質問 ID に紐づく通知を削除する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - questionID int64: 質問 ID
//
// return:
//   - int64: 削除した通知 ID（見つからない場合は -1）
//   - error: DB エラー
func (r *NoticeRepository) DeleteByQuestionID(ctx context.Context, questionID int64) (int64, error) {
	return DeleteNoticeByQuestionID(ctx, r.DB, questionID)
}

// MasterRepository はマスタデータの永続化を担う。
type MasterRepository struct {
	DB *gorm.DB
}

// NewMasterRepository は MasterRepository を生成する。
//
// args:
//   - db *gorm.DB: データベース接続
//
// return:
//   - *MasterRepository: リポジトリ
func NewMasterRepository(db *gorm.DB) *MasterRepository {
	return &MasterRepository{DB: db}
}

// GetRoles はロールマスタ一覧を取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//
// return:
//   - []entity.Role: ロール一覧
//   - error: DB エラー
func (r *MasterRepository) GetRoles(ctx context.Context) ([]entity.Role, error) {
	return GetMasterData[entity.Role](ctx, r.DB)
}

// GetCategories はカテゴリマスタ一覧を取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//
// return:
//   - []entity.Category: カテゴリ一覧
//   - error: DB エラー
func (r *MasterRepository) GetCategories(ctx context.Context) ([]entity.Category, error) {
	return GetMasterData[entity.Category](ctx, r.DB)
}

// GetSupportStatuses は支援ステータスマスタ一覧を取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//
// return:
//   - []entity.SupportStatus: 支援ステータス一覧
//   - error: DB エラー
func (r *MasterRepository) GetSupportStatuses(ctx context.Context) ([]entity.SupportStatus, error) {
	return GetMasterData[entity.SupportStatus](ctx, r.DB)
}
