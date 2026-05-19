package dto

import "strconv"

// =====================
// SenderForm は質問送信者のフォーム DTO である。
type SenderForm struct {
	ID             int64            `json:"id"`
	Name           string           `json:"name"`
	DepartmentName string           `json:"departmentName"`
	SenderTalks    []SenderTalkForm `json:"senderTalks,omitempty"`
}

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
