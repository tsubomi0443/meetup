package question

import (
	"strconv"

	"meetup/internal/usecases/dto"
)

// NormalizeQuestionFormClearSupportWhenUnassigned は、Support の支援ステータスが
// 未対応 (ID = 1) の場合に QuestionForm から Support / SupportID を削除する。
func NormalizeQuestionFormClearSupportWhenUnassigned(f *dto.QuestionForm) {
	if f == nil || f.Support == nil {
		return
	}
	if f.Support.SupportStatusIDInt64() != 1 {
		return
	}
	f.Support = nil
	f.SupportID = nil
}

// NormalizeQuestionFormAssignSupportUserWhenInProgress は、対応中 (SupportStatusID = 2) の PUT で
// Support の担当者 (UserID) が未設定なら、リクエスト元ユーザ (actorUserID) を担当として埋める。
func NormalizeQuestionFormAssignSupportUserWhenInProgress(f *dto.QuestionForm, actorUserID int64) {
	if f == nil || actorUserID <= 0 || f.Support == nil {
		return
	}
	if f.Support.SupportStatusIDInt64() != 2 {
		return
	}
	if f.Support.UserIDInt64() > 0 {
		return
	}
	f.Support.UserID = strconv.FormatInt(actorUserID, 10)
	f.SupportID = nil
	f.Support.ID = 0
}
