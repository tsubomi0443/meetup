// Package entity はアプリケーションの永続化モデル（GORM エンティティ）を定義する。
package entity

import (
	"time"

	"gorm.io/gorm"
)

// =====================
// ロール
// =====================

// Role はユーザーのロール（権限）を表すエンティティ。
type Role struct {
	ID        int64          `gorm:"column:id;primaryKey" json:"id"`
	Name      string         `gorm:"column:name" json:"name"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`

	Users []User `gorm:"foreignKey:RoleID;references:ID"`
}

// TableName は GORM のテーブル名を返す。
//
// return:
//   - string: ロールテーブル名
func (Role) TableName() string {
	return "roles"
}

// =====================
// ユーザー
// =====================

// User はシステム利用者を表すエンティティ。
type User struct {
	ID        int64          `gorm:"column:id;primaryKey" json:"id"`
	Name      string         `gorm:"column:name" json:"name"`
	Password  string         `gorm:"column:password" json:"-"`
	Email     string         `gorm:"column:email" json:"email"`
	Memo      string         `gorm:"column:memo" json:"memo"`
	RoleID    int64          `gorm:"column:role_id" json:"role_id"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`

	Role Role `gorm:"foreignKey:RoleID;references:ID" json:"role"`
}

// TableName は GORM のテーブル名を返す。
//
// return:
//   - string: ユーザーテーブル名
func (User) TableName() string {
	return "users"
}

// =====================
// カテゴリ
// =====================

// Category はタグの分類を表すエンティティ。
type Category struct {
	ID        int64          `gorm:"column:id;primaryKey"`
	Name      string         `gorm:"column:name"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`

	Tags []Tag `gorm:"foreignKey:CategoryID;references:ID"`
}

// TableName は GORM のテーブル名を返す。
//
// return:
//   - string: カテゴリテーブル名
func (Category) TableName() string {
	return "categories"
}

// =====================
// タグ
// =====================

// Tag は質問に付与するタグを表すエンティティ。
type Tag struct {
	ID         int64          `gorm:"column:id;primaryKey"`
	Name       string         `gorm:"column:name"`
	Usage      int            `gorm:"column:usage"`
	CategoryID int64          `gorm:"column:category_id"`
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index"`

	Category    Category     `gorm:"foreignKey:CategoryID;references:ID"`
	TagManagers []TagManager `gorm:"foreignKey:TagID;references:ID"`
}

// TableName は GORM のテーブル名を返す。
//
// return:
//   - string: タグテーブル名
func (Tag) TableName() string {
	return "tags"
}

// =====================
// 質問
// =====================

// Question は問い合わせ・相談案件を表すエンティティ。
type Question struct {
	ID               int64          `gorm:"column:id;primaryKey"`
	OriginQuestionID *int64         `gorm:"column:origin_question_id"`
	SupportID        *int64         `gorm:"column:support_id"`
	TalkroomID       string         `gorm:"column:talkroom_id"`
	Title            string         `gorm:"column:title"`
	Content          string         `gorm:"column:content"`
	Due              *time.Time     `gorm:"column:due"`
	CreatedAt        time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt        time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at;index"`

	Answer           []Answer          `gorm:"foreignKey:QuestionID;references:ID"`
	Memos            []Memo            `gorm:"foreignKey:QuestionID;references:ID"`
	TagManagers      []TagManager      `gorm:"foreignKey:QuestionID;references:ID"`
	RelatedQuestions []RelatedQuestion `gorm:"foreignKey:QuestionID;references:ID"`
	SenderTalks      []SenderTalk      `gorm:"foreignKey:QuestionID;references:ID"`
	Support          *Support          `gorm:"foreignKey:SupportID;references:ID"`
}

// TableName は GORM のテーブル名を返す。
//
// return:
//   - string: 質問テーブル名
func (Question) TableName() string {
	return "questions"
}

// =====================
// 回答
// =====================

// Answer は質問に対する回答を表すエンティティ。
type Answer struct {
	ID         int64          `gorm:"column:id;primaryKey"`
	UserID     int64          `gorm:"column:user_id"`
	QuestionID int64          `gorm:"column:question_id"`
	Content    string         `gorm:"column:content"`
	IsFinal    bool           `gorm:"column:is_final"`
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index"`

	Question      Question       `gorm:"foreignKey:QuestionID;references:ID"`
	User          User           `gorm:"foreignKey:UserID;references:ID"`
	ReferManagers []ReferManager `gorm:"foreignKey:AnswerID;references:ID"`
}

// TableName は GORM のテーブル名を返す。
//
// return:
//   - string: 回答テーブル名
func (Answer) TableName() string {
	return "answers"
}

// =====================
// メモ
// =====================

// Memo は質問に紐づく内部メモを表すエンティティ。
type Memo struct {
	ID         int64          `gorm:"column:id;primaryKey"`
	UserID     int64          `gorm:"column:user_id"`
	Content    string         `gorm:"column:content"`
	QuestionID int64          `gorm:"column:question_id"`
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index"`

	User     User      `gorm:"foreignKey:UserID;references:ID"`
	Question *Question `gorm:"foreignKey:QuestionID;references:ID"`
}

// TableName は GORM のテーブル名を返す。
//
// return:
//   - string: メモテーブル名
func (Memo) TableName() string {
	return "memos"
}

// =====================
// 参照資料
// =====================

// Refer は回答に添付する参照 URL を表すエンティティ。
type Refer struct {
	ID        int64          `gorm:"column:id;primaryKey"`
	Title     string         `gorm:"column:title"`
	URL       string         `gorm:"column:url"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`

	ReferManagers []ReferManager `gorm:"foreignKey:ReferID;references:ID"`
}

// TableName は GORM のテーブル名を返す。
//
// return:
//   - string: 参照資料テーブル名
func (Refer) TableName() string {
	return "refers"
}

// =====================
// 参照資料紐付け
// =====================

// ReferManager は回答と参照資料の中間テーブルを表すエンティティ。
type ReferManager struct {
	ID        int64          `gorm:"column:id;primaryKey"`
	AnswerID  int64          `gorm:"column:answer_id"`
	ReferID   int64          `gorm:"column:refer_id"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`

	Answer Answer `gorm:"foreignKey:AnswerID;references:ID"`
	Refer  Refer  `gorm:"foreignKey:ReferID;references:ID"`
}

// TableName は GORM のテーブル名を返す。
//
// return:
//   - string: 参照資料紐付けテーブル名
func (ReferManager) TableName() string {
	return "refer_managers"
}

// =====================
// タグ紐付け
// =====================

// TagManager は質問とタグの中間テーブルを表すエンティティ。
type TagManager struct {
	ID         int64          `gorm:"column:id;primaryKey"`
	TagID      int64          `gorm:"column:tag_id"`
	QuestionID int64          `gorm:"column:question_id"`
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index"`

	Tag      Tag      `gorm:"foreignKey:TagID;references:ID"`
	Question Question `gorm:"foreignKey:QuestionID;references:ID"`
}

// TableName は GORM のテーブル名を返す。
//
// return:
//   - string: タグ紐付けテーブル名
func (TagManager) TableName() string {
	return "tag_managers"
}

// =====================
// エスカレーション
// =====================

// Escalation は質問のエスカレーション履歴を表すエンティティ。
type Escalation struct {
	ID             int64          `gorm:"column:id;primaryKey"`
	FromQuestionID int64          `gorm:"column:from_question_id"`
	ToQuestionID   int64          `gorm:"column:to_question_id"`
	EscalatedAt    time.Time      `gorm:"column:escalated_at"`
	CreatedAt      time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index"`

	FromQuestion Question `gorm:"foreignKey:FromQuestionID;references:ID"`
	ToQuestion   Question `gorm:"foreignKey:ToQuestionID;references:ID"`
}

// TableName は GORM のテーブル名を返す。
//
// return:
//   - string: エスカレーションテーブル名
func (Escalation) TableName() string {
	return "escalations"
}

// =====================
// サポートステータス
// =====================

// SupportStatus はサポート対応の状態マスタを表すエンティティ。
type SupportStatus struct {
	ID        int64          `gorm:"column:id;primaryKey"`
	Name      string         `gorm:"column:name"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`

	Supports []Support `gorm:"foreignKey:SupportStatusID;references:ID"`
}

// TableName は GORM のテーブル名を返す。
//
// return:
//   - string: サポートステータステーブル名
func (SupportStatus) TableName() string {
	return "support_statuses"
}

// =====================
// サポート
// =====================

// Support は質問に対する担当サポートを表すエンティティ。
type Support struct {
	ID              int64          `gorm:"column:id;primaryKey"`
	UserID          int64          `gorm:"column:user_id"`
	SupportStatusID int64          `gorm:"column:support_status_id"`
	CreatedAt       time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index"`

	User          User          `gorm:"foreignKey:UserID;references:ID"`
	SupportStatus SupportStatus `gorm:"foreignKey:SupportStatusID;references:ID"`
}

// TableName は GORM のテーブル名を返す。
//
// return:
//   - string: サポートテーブル名
func (Support) TableName() string {
	return "supports"
}

// =====================
// 通知種別
// =====================

// NoticeType は通知の種類マスタを表すエンティティ。
type NoticeType struct {
	ID        int64          `gorm:"column:id;primaryKey"`
	Name      string         `gorm:"column:name"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`

	Notices []Notice `gorm:"foreignKey:TypeID;references:ID"`
}

// TableName は GORM のテーブル名を返す。
//
// return:
//   - string: 通知種別テーブル名
func (NoticeType) TableName() string {
	return "notice_types"
}

// =====================
// 通知
// =====================

// Notice は画面上に表示する通知を表すエンティティ。
type Notice struct {
	ID         int64          `gorm:"column:id;primaryKey"`
	TypeID     int64          `gorm:"column:type_id"`
	QuestionID *int64         `gorm:"column:question_id"`
	Content    *string        `gorm:"column:content"`
	DisplayDue *time.Time     `gorm:"column:display_due"`
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index"`

	NoticeType NoticeType `gorm:"foreignKey:TypeID;references:ID"`
	Question   *Question  `gorm:"foreignKey:QuestionID;references:ID"`
}

// TableName は GORM のテーブル名を返す。
//
// return:
//   - string: 通知テーブル名
func (Notice) TableName() string {
	return "notices"
}

// =====================
// 関連質問
// =====================

// RelatedQuestion は質問同士の関連を表すエンティティ。
type RelatedQuestion struct {
	ID                int64          `gorm:"column:id;primaryKey"`
	QuestionID        int64          `gorm:"column:question_id"`
	RelatedQuestionID int64          `gorm:"column:related_question_id"`
	CreatedAt         time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt         time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt         gorm.DeletedAt `gorm:"column:deleted_at;index"`

	Question        Question `gorm:"foreignKey:QuestionID;references:ID"`
	RelatedQuestion Question `gorm:"foreignKey:RelatedQuestionID;references:ID"`
}

// TableName は GORM のテーブル名を返す。
//
// return:
//   - string: 関連質問テーブル名
func (RelatedQuestion) TableName() string {
	return "related_questions"
}

// =====================
// 送信者
// =====================

// Sender はトークルーム上の問い合わせ送信者を表すエンティティ。
type Sender struct {
	ID             int64  `gorm:"column:id;primaryKey"`
	UID            string `gorm:"column:uid;uniqueIndex"`
	Name           string `gorm:"column:name"`
	DepartmentName string `gorm:"column:department_name"`

	SenderTalks []SenderTalk `gorm:"foreignKey:SenderID;references:ID"`
}

// TableName は GORM のテーブル名を返す。
//
// return:
//   - string: 送信者テーブル名
func (Sender) TableName() string {
	return "senders"
}

// =====================
// 送信者トーク
// =====================

// SenderTalk は送信者の発言内容を表すエンティティ。
type SenderTalk struct {
	ID         int64          `gorm:"column:id;primaryKey"`
	Content    string         `gorm:"column:content"`
	SenderID   int64          `gorm:"column:sender_id"`
	QuestionID int64          `gorm:"column:question_id"`
	TalkroomID string         `gorm:"column:talkroom_id"`
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index"`

	Sender   Sender   `gorm:"foreignKey:SenderID;references:ID"`
	Question Question `gorm:"foreignKey:QuestionID;references:ID"`
}

// TableName は GORM のテーブル名を返す。
//
// return:
//   - string: 送信者トークテーブル名
func (SenderTalk) TableName() string {
	return "sender_talks"
}
