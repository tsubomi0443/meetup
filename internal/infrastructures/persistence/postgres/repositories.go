package postgres

import (
	"context"

	"meetup/internal/domains/entity"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (entity.User, error) {
	return GetUserByID(ctx, r.DB, id)
}

func (r *UserRepository) GetPasswordByEmail(ctx context.Context, email string) (string, error) {
	return GetUserPasswordByEmail(ctx, r.DB, email)
}

func (r *UserRepository) GetUserInfo(ctx context.Context, email, pass string, preloads ...string) (entity.User, error) {
	return GetUserInfo(ctx, r.DB, email, pass, preloads...)
}

func (r *UserRepository) GetUsers(ctx context.Context) ([]entity.User, error) {
	return GetUsers(ctx, r.DB)
}

func (r *UserRepository) Register(ctx context.Context, model *entity.User) error {
	return Register(ctx, r.DB, model)
}

func (r *UserRepository) UpdateByID(ctx context.Context, id int64, model entity.User, omit ...string) (int, error) {
	return UpdateByID(ctx, r.DB, id, model, omit...)
}

func (r *UserRepository) DeleteByID(ctx context.Context, id int64) error {
	return DeleteUserByID(ctx, r.DB, id)
}

type QuestionRepository struct {
	DB *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) *QuestionRepository {
	return &QuestionRepository{DB: db}
}

func (r *QuestionRepository) GetByID(ctx context.Context, id int64) (entity.Question, error) {
	return GetQuestion(ctx, r.DB, id)
}

func (r *QuestionRepository) GetAll(ctx context.Context) ([]entity.Question, error) {
	return GetQuestions(ctx, r.DB)
}

func (r *QuestionRepository) Register(ctx context.Context, model *entity.Question) error {
	return Register(ctx, r.DB, model)
}

func (r *QuestionRepository) UpdateInTransaction(ctx context.Context, q entity.Question) error {
	return UpdateQuestionInTransaction(ctx, r.DB, q)
}

func (r *QuestionRepository) DeleteByID(ctx context.Context, id int64) error {
	return DeleteQuestionByID(ctx, r.DB, id)
}

type TagRepository struct {
	DB *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{DB: db}
}

func (r *TagRepository) GetAll(ctx context.Context) ([]entity.Tag, error) {
	return GetTags(ctx, r.DB)
}

func (r *TagRepository) GetByID(ctx context.Context, id int64) (entity.Tag, error) {
	return GetTagByID(ctx, r.DB, id)
}

func (r *TagRepository) Register(ctx context.Context, model *entity.Tag) error {
	return Register(ctx, r.DB, model)
}

func (r *TagRepository) UpdateByID(ctx context.Context, id int64, model entity.Tag, omit ...string) (int, error) {
	return UpdateByID(ctx, r.DB, id, model, omit...)
}

func (r *TagRepository) DeleteByID(ctx context.Context, id int64) error {
	return DeleteTagByID(ctx, r.DB, id)
}

type NoticeRepository struct {
	DB *gorm.DB
}

func NewNoticeRepository(db *gorm.DB) *NoticeRepository {
	return &NoticeRepository{DB: db}
}

func (r *NoticeRepository) GetAll(ctx context.Context) ([]entity.Notice, error) {
	return GetNotice(ctx, r.DB)
}

func (r *NoticeRepository) GetByQuestionSilent(ctx context.Context, question entity.Question) (entity.Notice, error) {
	return GetNoticeByQuestionSilent(ctx, r.DB, question)
}

func (r *NoticeRepository) RegisterByQuestionID(ctx context.Context, questionID int64) error {
	return RegisterNoticeByQuestionID(ctx, r.DB, questionID)
}

func (r *NoticeRepository) DeleteByQuestionID(ctx context.Context, questionID int64) (int64, error) {
	return DeleteNoticeByQuestionID(ctx, r.DB, questionID)
}

type MasterRepository struct {
	DB *gorm.DB
}

func NewMasterRepository(db *gorm.DB) *MasterRepository {
	return &MasterRepository{DB: db}
}

func (r *MasterRepository) GetRoles(ctx context.Context) ([]entity.Role, error) {
	return GetMasterData[entity.Role](ctx, r.DB)
}

func (r *MasterRepository) GetCategories(ctx context.Context) ([]entity.Category, error) {
	return GetMasterData[entity.Category](ctx, r.DB)
}

func (r *MasterRepository) GetSupportStatuses(ctx context.Context) ([]entity.SupportStatus, error) {
	return GetMasterData[entity.SupportStatus](ctx, r.DB)
}
