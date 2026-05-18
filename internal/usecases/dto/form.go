// Package dto は API・画面との入出力で用いるフォーム DTO を定義する。
// エンティティとの相互変換は mapper パッケージが担う。
package dto

import "strconv"

// =====================
// ロール（ROLE）
// =====================

// RoleForm はロールのフォーム DTO である。
type RoleForm struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	CreatedAt *string    `json:"createdAt,omitempty"`
	UpdatedAt *string    `json:"updatedAt,omitempty"`
	DeletedAt *string    `json:"deletedAt,omitempty"`
	Users     []UserForm `json:"users,omitempty"`
}

// =====================
// サポートステータス（SUPPORT_STATUS）
// =====================

// SupportStatusForm はサポートステータスのフォーム DTO である。
type SupportStatusForm struct {
	ID        int64         `json:"id"`
	Name      string        `json:"name"`
	CreatedAt *string       `json:"createdAt,omitempty"`
	UpdatedAt *string       `json:"updatedAt,omitempty"`
	DeletedAt *string       `json:"deletedAt,omitempty"`
	Supports  []SupportForm `json:"supports,omitempty"`
}

// =====================
// サポート（SUPPORT）
// =====================

// SupportForm はサポート（対応）のフォーム DTO である。
type SupportForm struct {
	ID              int64              `json:"id"`
	UserID          string             `json:"userId"`
	SupportStatusID string             `json:"supportStatusId"`
	CreatedAt       *string            `json:"createdAt,omitempty"`
	UpdatedAt       *string            `json:"updatedAt,omitempty"`
	DeletedAt       *string            `json:"deletedAt,omitempty"`
	User            *UserForm          `json:"user,omitempty"`
	SupportStatus   *SupportStatusForm `json:"supportStatus,omitempty"`
	Question        *QuestionForm      `json:"question,omitempty"`
}

// UserIDInt64 は UserID 文字列を int64 に変換する。
//
// return:
//   - int64: ユーザー ID（変換失敗時は -1）
func (f SupportForm) UserIDInt64() int64 {
	if val, err := strconv.ParseInt(f.UserID, 10, 64); err == nil {
		return val
	}
	return -1
}

// SupportStatusIDInt64 は SupportStatusID 文字列を int64 に変換する。
//
// return:
//   - int64: サポートステータス ID（変換失敗時は -1）
func (f SupportForm) SupportStatusIDInt64() int64 {
	if val, err := strconv.ParseInt(f.SupportStatusID, 10, 64); err == nil {
		return val
	}
	return -1
}

// =====================
// ユーザー（USER）
// =====================

// UserForm はユーザーのフォーム DTO である。
type UserForm struct {
	ID        int64         `json:"id"`
	Name      string        `json:"name"`
	Email     string        `json:"email"`
	Memo      string        `json:"memo"`
	RoleID    string        `json:"roleId"`
	CreatedAt *string       `json:"createdAt,omitempty"`
	UpdatedAt *string       `json:"updatedAt,omitempty"`
	DeletedAt *string       `json:"deletedAt,omitempty"`
	Role      *RoleForm     `json:"role,omitempty"`
	Supports  []SupportForm `json:"supports,omitempty"`
	Answers   []AnswerForm  `json:"answers,omitempty"`
	Memos     []MemoForm    `json:"memos,omitempty"`
	Password  string        `json:"pass,omitempty"`
}

// RoleIDInt64 は RoleID 文字列を int64 に変換する。
//
// return:
//   - int64: ロール ID（変換失敗時は -1）
func (uf UserForm) RoleIDInt64() int64 {
	if val, err := strconv.ParseInt(uf.RoleID, 10, 64); err == nil {
		return val
	}
	return -1
}

// =====================
// カテゴリ（CATEGORY）
// =====================

// CategoryForm はカテゴリのフォーム DTO である。
type CategoryForm struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt *string   `json:"createdAt,omitempty"`
	UpdatedAt *string   `json:"updatedAt,omitempty"`
	DeletedAt *string   `json:"deletedAt,omitempty"`
	Tags      []TagForm `json:"tags,omitempty"`
}

// =====================
// タグ（TAG）
// =====================

// TagForm はタグのフォーム DTO である。
type TagForm struct {
	ID         int64          `json:"id"`
	Name       string         `json:"name"`
	Usage      int            `json:"usage"`
	CategoryID string         `json:"categoryId"`
	CreatedAt  *string        `json:"createdAt,omitempty"`
	UpdatedAt  *string        `json:"updatedAt,omitempty"`
	DeletedAt  *string        `json:"deletedAt,omitempty"`
	Category   *CategoryForm  `json:"category,omitempty"`
	Questions  []QuestionForm `json:"questions,omitempty"`
}

// CategoryIDInt64 は CategoryID 文字列を int64 に変換する。
//
// return:
//   - int64: カテゴリ ID（変換失敗時は -1）
func (tf TagForm) CategoryIDInt64() int64 {
	if val, err := strconv.ParseInt(tf.CategoryID, 10, 64); err == nil {
		return val
	}
	return -1
}

// =====================
// 参照リンク（REFER）
// =====================

// ReferForm は参照リンクのフォーム DTO である。
type ReferForm struct {
	ID        int64        `json:"id"`
	Title     string       `json:"title"`
	URL       string       `json:"url"`
	CreatedAt *string      `json:"createdAt,omitempty"`
	UpdatedAt *string      `json:"updatedAt,omitempty"`
	DeletedAt *string      `json:"deletedAt,omitempty"`
	Answers   []AnswerForm `json:"answers,omitempty"`
}

// =====================
// メモ（MEMO）
// =====================

// MemoForm は質問メモのフォーム DTO である。
type MemoForm struct {
	ID         int64         `json:"id"`
	QuestionID string        `json:"questionId"`
	UserID     string        `json:"userId"`
	Content    string        `json:"content"`
	CreatedAt  *string       `json:"createdAt,omitempty"`
	UpdatedAt  *string       `json:"updatedAt,omitempty"`
	DeletedAt  *string       `json:"deletedAt,omitempty"`
	Question   *QuestionForm `json:"question,omitempty"`
	User       *UserForm     `json:"user,omitempty"`
}

// QuestionIDInt64 は QuestionID 文字列を int64 に変換する。
//
// return:
//   - int64: 質問 ID（変換失敗時は -1）
func (f MemoForm) QuestionIDInt64() int64 {
	if val, err := strconv.ParseInt(f.QuestionID, 10, 64); err == nil {
		return val
	}
	return -1
}

// UserIDInt64 は UserID 文字列を int64 に変換する。
//
// return:
//   - int64: ユーザー ID（変換失敗時は -1）
func (f MemoForm) UserIDInt64() int64 {
	if val, err := strconv.ParseInt(f.UserID, 10, 64); err == nil {
		return val
	}
	return -1
}

// =====================
// 回答（ANSWER）
// =====================

// AnswerForm は回答のフォーム DTO である。
type AnswerForm struct {
	ID        int64       `json:"id"`
	UserID    string      `json:"userId"`
	Content   string      `json:"content"`
	IsFinal   bool        `json:"isFinal"`
	CreatedAt *string     `json:"createdAt,omitempty"`
	UpdatedAt *string     `json:"updatedAt,omitempty"`
	DeletedAt *string     `json:"deletedAt,omitempty"`
	User      *UserForm   `json:"user,omitempty"`
	Refers    []ReferForm `json:"refers,omitempty"`
}

// UserIDInt64 は UserID 文字列を int64 に変換する。
//
// return:
//   - int64: ユーザー ID（変換失敗時は -1）
func (f AnswerForm) UserIDInt64() int64 {
	if val, err := strconv.ParseInt(f.UserID, 10, 64); err == nil {
		return val
	}
	return -1
}

// =====================
// エスカレーション（ESCALATION）
// =====================

// EscalationForm は質問エスカレーションのフォーム DTO である。
type EscalationForm struct {
	ID             int64         `json:"id"`
	FromQuestionID string        `json:"fromQuestionId"`
	ToQuestionID   string        `json:"toQuestionId"`
	EscalatedAt    *string       `json:"escalatedAt,omitempty"`
	CreatedAt      *string       `json:"createdAt,omitempty"`
	UpdatedAt      *string       `json:"updatedAt,omitempty"`
	DeletedAt      *string       `json:"deletedAt,omitempty"`
	FromQuestion   *QuestionForm `json:"fromQuestion,omitempty"`
	ToQuestion     *QuestionForm `json:"toQuestion,omitempty"`
}

// FromQuestionIDInt64 は FromQuestionID 文字列を int64 に変換する。
//
// return:
//   - int64: エスカレーション元質問 ID（変換失敗時は -1）
func (f EscalationForm) FromQuestionIDInt64() int64 {
	if val, err := strconv.ParseInt(f.FromQuestionID, 10, 64); err == nil {
		return val
	}
	return -1
}

// ToQuestionIDInt64 は ToQuestionID 文字列を int64 に変換する。
//
// return:
//   - int64: エスカレーション先質問 ID（変換失敗時は -1）
func (f EscalationForm) ToQuestionIDInt64() int64 {
	if val, err := strconv.ParseInt(f.ToQuestionID, 10, 64); err == nil {
		return val
	}
	return -1
}

// =====================
// 通知種別（NOTICE_TYPE）
// =====================

// NoticeTypeForm は通知種別のフォーム DTO である。
type NoticeTypeForm struct {
	ID        int64        `json:"id"`
	Name      string       `json:"name"`
	CreatedAt *string      `json:"createdAt,omitempty"`
	UpdatedAt *string      `json:"updatedAt,omitempty"`
	DeletedAt *string      `json:"deletedAt,omitempty"`
	Notices   []NoticeForm `json:"notices,omitempty"`
}

// =====================
// 通知（NOTICE）
// =====================

// NoticeForm は通知のフォーム DTO である。
type NoticeForm struct {
	ID         int64           `json:"id"`
	TypeID     int64           `json:"typeId"`
	QuestionID *string         `json:"questionId,omitempty"`
	Content    *string         `json:"content,omitempty"`
	DisplayDue *string         `json:"displayDue,omitempty"`
	CreatedAt  *string         `json:"createdAt,omitempty"`
	UpdatedAt  *string         `json:"updatedAt,omitempty"`
	DeletedAt  *string         `json:"deletedAt,omitempty"`
	NoticeType *NoticeTypeForm `json:"noticeType,omitempty"`
	Question   *QuestionForm   `json:"question,omitempty"`
}

// QuestionIDInt64 は QuestionID 文字列を int64 に変換する。nil または空文字の場合は -1 を返す。
//
// return:
//   - int64: 質問 ID（未設定・変換失敗時は -1）
func (f NoticeForm) QuestionIDInt64() int64 {
	if f.QuestionID == nil || *f.QuestionID == "" {
		return -1
	}
	if val, err := strconv.ParseInt(*f.QuestionID, 10, 64); err == nil {
		return val
	}
	return -1
}

// =====================
// 質問（QUESTION）
// =====================

// QuestionForm は質問のフォーム DTO である。DeletedAt は論理削除日時を表す。
type QuestionForm struct {
	ID               int64                 `json:"id"`
	OriginQuestionID *string               `json:"originQuestionId,omitempty"`
	SupportID        *int64                `json:"supportId,omitempty"`
	Title            string                `json:"title"`
	Content          string                `json:"content"`
	Due              *string               `json:"due,omitempty"`
	CreatedAt        *string               `json:"createdAt,omitempty"`
	UpdatedAt        *string               `json:"updatedAt,omitempty"`
	DeletedAt        *string               `json:"deletedAt,omitempty"`
	OriginQuestion   *QuestionForm         `json:"originQuestion,omitempty"`
	SubQuestions     []QuestionForm        `json:"subQuestions,omitempty"`
	Support          *SupportForm          `json:"support,omitempty"`
	Answers          []AnswerForm          `json:"answers,omitempty"`
	Memos            []MemoForm            `json:"memos,omitempty"`
	Tags             []TagForm             `json:"tags,omitempty"`
	EscalationsFrom  []EscalationForm      `json:"escalationsFrom,omitempty"`
	EscalationsTo    []EscalationForm      `json:"escalationsTo,omitempty"`
	RelatedQuestions []RelatedQuestionForm `json:"relatedQuestions,omitempty"`
	SenderTalks      []SenderTalkForm      `json:"senderTalks,omitempty"`
}

// OriginQuestionIDInt64 は OriginQuestionID 文字列を int64 に変換する。nil または空文字の場合は -1 を返す。
//
// return:
//   - int64: 元質問 ID（未設定・変換失敗時は -1）
func (f QuestionForm) OriginQuestionIDInt64() int64 {
	if f.OriginQuestionID == nil || *f.OriginQuestionID == "" {
		return -1
	}
	if val, err := strconv.ParseInt(*f.OriginQuestionID, 10, 64); err == nil {
		return val
	}
	return -1
}

// =====================
// 参照リンク管理（REFER_MANAGER）
// =====================

// ReferManagerForm は回答と参照リンクの紐付けフォーム DTO である。
type ReferManagerForm struct {
	ID        int64       `json:"id"`
	AnswerID  string      `json:"answerId"`
	ReferID   string      `json:"referId"`
	CreatedAt *string     `json:"createdAt,omitempty"`
	UpdatedAt *string     `json:"updatedAt,omitempty"`
	DeletedAt *string     `json:"deletedAt,omitempty"`
	Answer    *AnswerForm `json:"answer,omitempty"`
	Refer     *ReferForm  `json:"refer,omitempty"`
}

// AnswerIDInt64 は AnswerID 文字列を int64 に変換する。
//
// return:
//   - int64: 回答 ID（変換失敗時は -1）
func (f ReferManagerForm) AnswerIDInt64() int64 {
	if val, err := strconv.ParseInt(f.AnswerID, 10, 64); err == nil {
		return val
	}
	return -1
}

// ReferIDInt64 は ReferID 文字列を int64 に変換する。
//
// return:
//   - int64: 参照リンク ID（変換失敗時は -1）
func (f ReferManagerForm) ReferIDInt64() int64 {
	if val, err := strconv.ParseInt(f.ReferID, 10, 64); err == nil {
		return val
	}
	return -1
}

// =====================
// タグ管理（TAG_MANAGER）
// =====================

// TagManagerForm は質問とタグの紐付けフォーム DTO である。
type TagManagerForm struct {
	ID         int64         `json:"id"`
	TagID      string        `json:"tagId"`
	QuestionID string        `json:"questionId"`
	CreatedAt  *string       `json:"createdAt,omitempty"`
	UpdatedAt  *string       `json:"updatedAt,omitempty"`
	DeletedAt  *string       `json:"deletedAt,omitempty"`
	Tag        *TagForm      `json:"tag,omitempty"`
	Question   *QuestionForm `json:"question,omitempty"`
}

// TagIDInt64 は TagID 文字列を int64 に変換する。
//
// return:
//   - int64: タグ ID（変換失敗時は -1）
func (f TagManagerForm) TagIDInt64() int64 {
	if val, err := strconv.ParseInt(f.TagID, 10, 64); err == nil {
		return val
	}
	return -1
}

// QuestionIDInt64 は QuestionID 文字列を int64 に変換する。
//
// return:
//   - int64: 質問 ID（変換失敗時は -1）
func (f TagManagerForm) QuestionIDInt64() int64 {
	if val, err := strconv.ParseInt(f.QuestionID, 10, 64); err == nil {
		return val
	}
	return -1
}

// =====================
// 関連質問（RELATED_QUESTION）
// =====================

// RelatedQuestionForm は関連質問の紐付けフォーム DTO である。
type RelatedQuestionForm struct {
	ID                int64         `json:"id"`
	QuestionID        string        `json:"questionId"`
	RelatedQuestionID string        `json:"relatedQuestionId"`
	CreatedAt         *string       `json:"createdAt,omitempty"`
	UpdatedAt         *string       `json:"updatedAt,omitempty"`
	DeletedAt         *string       `json:"deletedAt,omitempty"`
	Question          *QuestionForm `json:"question,omitempty"`
	RelatedQuestion   *QuestionForm `json:"relatedQuestion,omitempty"`
}

// QuestionIDInt64 は QuestionID 文字列を int64 に変換する。
//
// return:
//   - int64: 質問 ID（変換失敗時は -1）
func (f RelatedQuestionForm) QuestionIDInt64() int64 {
	if val, err := strconv.ParseInt(f.QuestionID, 10, 64); err == nil {
		return val
	}
	return -1
}

// RelatedQuestionIDInt64 は RelatedQuestionID 文字列を int64 に変換する。
//
// return:
//   - int64: 関連質問 ID（変換失敗時は -1）
func (f RelatedQuestionForm) RelatedQuestionIDInt64() int64 {
	if val, err := strconv.ParseInt(f.RelatedQuestionID, 10, 64); err == nil {
		return val
	}
	return -1
}

// =====================
// 送信者（SENDER）
// =====================

// SenderForm は質問送信者のフォーム DTO である。
type SenderForm struct {
	ID             int64            `json:"id"`
	Name           string           `json:"name"`
	DepartmentName string           `json:"departmentName"`
	SenderTalks    []SenderTalkForm `json:"senderTalks,omitempty"`
}

// =====================
// 送信者トーク（SENDER_TALK）
// =====================

// SenderTalkForm は送信者の発言（トーク）のフォーム DTO である。
type SenderTalkForm struct {
	ID         int64       `json:"id"`
	Content    string      `json:"content"`
	SenderID   string      `json:"senderId"`
	QuestionID string      `json:"questionId"`
	CreatedAt  *string     `json:"createdAt,omitempty"`
	UpdatedAt  *string     `json:"updatedAt,omitempty"`
	DeletedAt  *string     `json:"deletedAt,omitempty"`
	Sender     *SenderForm `json:"sender,omitempty"`
}

// SenderIDInt64 は SenderID 文字列を int64 に変換する。
//
// return:
//   - int64: 送信者 ID（変換失敗時は -1）
func (f SenderTalkForm) SenderIDInt64() int64 {
	if val, err := strconv.ParseInt(f.SenderID, 10, 64); err == nil {
		return val
	}
	return -1
}

// QuestionIDInt64 は QuestionID 文字列を int64 に変換する。
//
// return:
//   - int64: 質問 ID（変換失敗時は -1）
func (f SenderTalkForm) QuestionIDInt64() int64 {
	if val, err := strconv.ParseInt(f.QuestionID, 10, 64); err == nil {
		return val
	}
	return -1
}
