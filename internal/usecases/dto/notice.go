package dto

import "strconv"

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
