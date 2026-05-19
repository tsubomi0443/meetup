package mapper

import (
	"strconv"

	"meetup/internal/domains/entity"
	"meetup/internal/usecases/dto"
)

// =====================
// メモ（entity.Memo）

// =====================
// MemoFromEntity は entity.Memo を dto.MemoForm に変換する。
//
// args:
//   - e entity.Memo: 変換元エンティティ
//
// return:
//   - dto.MemoForm: メモフォーム DTO
func MemoFromEntity(e entity.Memo) dto.MemoForm {
	f := dto.MemoForm{
		ID:         e.ID,
		QuestionID: strconv.FormatInt(e.QuestionID, 10),
		UserID:     strconv.FormatInt(e.UserID, 10),
		Content:    e.Content,
		CreatedAt:  timeToISO(e.CreatedAt),
		UpdatedAt:  timeToISO(e.UpdatedAt),
		DeletedAt:  deletedAtToISO(e.DeletedAt),
	}
	if e.User.ID != 0 {
		u := UserFromEntity(e.User)
		f.User = &u
	}
	return f
}

// MemoToEntity は dto.MemoForm を entity.Memo に変換する。
//
// args:
//   - f dto.MemoForm: 変換元フォーム DTO
//
// return:
//   - entity.Memo: メモエンティティ
func MemoToEntity(f dto.MemoForm) entity.Memo {
	e := entity.Memo{
		ID:         f.ID,
		QuestionID: f.QuestionIDInt64(),
		UserID:     f.UserIDInt64(),
		Content:    f.Content,
	}
	if f.User != nil {
		e.User = UserToEntity(*f.User)
	}
	return e
}
