// Package mapper はドメインエンティティと dto フォーム DTO の相互変換を提供する。
package mapper

import (
	"strconv"
	"strings"
	"time"

	"meetup/internal/domains/entity"
	"meetup/internal/usecases/dto"

	"gorm.io/gorm"
)

// =====================
// 日時ヘルパー（ISO8601 / RFC3339）
// =====================

// timePtrToISO は *time.Time を RFC3339Nano 形式の ISO 8601 文字列ポインタに変換する。
//
// args:
//   - t *time.Time: 変換元の日時（nil またはゼロ値の場合は nil を返す）
//
// return:
//   - *string: ISO 8601 文字列（変換不可時は nil）
func timePtrToISO(t *time.Time) *string {
	if t == nil || t.IsZero() {
		return nil
	}
	s := t.UTC().Format(time.RFC3339Nano)
	return &s
}

// timeToISO は time.Time を RFC3339Nano 形式の ISO 8601 文字列ポインタに変換する。
//
// args:
//   - t time.Time: 変換元の日時（ゼロ値の場合は nil を返す）
//
// return:
//   - *string: ISO 8601 文字列（変換不可時は nil）
func timeToISO(t time.Time) *string {
	if t.IsZero() {
		return nil
	}
	s := t.UTC().Format(time.RFC3339Nano)
	return &s
}

// isoToTimePtr は ISO 8601 文字列を *time.Time に変換する。RFC3339Nano および RFC3339 に対応する。
//
// args:
//   - s *string: ISO 8601 文字列（nil または空文字の場合は nil を返す）
//
// return:
//   - *time.Time: 変換後の日時（パース失敗時は nil）
func isoToTimePtr(s *string) *time.Time {
	if s == nil || *s == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339Nano, *s)
	if err != nil {
		t, err = time.Parse(time.RFC3339, *s)
	}
	if err != nil {
		return nil
	}
	return &t
}

// isoToTime は ISO 8601 文字列を time.Time に変換する。未設定・パース失敗時はゼロ値を返す。
//
// args:
//   - s *string: ISO 8601 文字列
//
// return:
//   - time.Time: 変換後の日時（未設定・失敗時はゼロ値）
func isoToTime(s *string) time.Time {
	if s == nil || *s == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339Nano, *s)
	if err != nil {
		t, _ = time.Parse(time.RFC3339, *s)
	}
	return t
}

// deletedAtToISO は gorm.DeletedAt を論理削除日時の ISO 8601 文字列に変換する。
//
// args:
//   - d gorm.DeletedAt: GORM の論理削除日時
//
// return:
//   - *string: 削除日時の ISO 8601 文字列（未削除・無効時は nil）
func deletedAtToISO(d gorm.DeletedAt) *string {
	// 論理削除されていない、または値がない場合は nil を返す。
	if !d.Valid || d.Time.IsZero() {
		return nil
	}
	s := d.Time.UTC().Format(time.RFC3339Nano)
	return &s
}

// =====================
// ロール（entity.Role）
// =====================

// RoleFromEntity は entity.Role を dto.RoleForm に変換する。Users があればネストして変換する。
//
// args:
//   - e entity.Role: 変換元エンティティ
//
// return:
//   - dto.RoleForm: ロールフォーム DTO
func RoleFromEntity(e entity.Role) dto.RoleForm {
	f := dto.RoleForm{
		ID:        e.ID,
		Name:      e.Name,
		CreatedAt: timeToISO(e.CreatedAt),
		UpdatedAt: timeToISO(e.UpdatedAt),
		DeletedAt: deletedAtToISO(e.DeletedAt),
	}
	if len(e.Users) > 0 {
		f.Users = make([]dto.UserForm, len(e.Users))
		for i := range e.Users {
			f.Users[i] = UserFromEntityNoRole(e.Users[i])
		}
	}
	return f
}

// roleFromEntityShallow は entity.Role を関連 Users なしの dto.RoleForm に変換する（循環参照回避用）。
//
// args:
//   - e entity.Role: 変換元エンティティ
//
// return:
//   - dto.RoleForm: ロールフォーム DTO（Users なし）
func roleFromEntityShallow(e entity.Role) dto.RoleForm {
	return dto.RoleForm{
		ID:        e.ID,
		Name:      e.Name,
		CreatedAt: timeToISO(e.CreatedAt),
		UpdatedAt: timeToISO(e.UpdatedAt),
		DeletedAt: deletedAtToISO(e.DeletedAt),
	}
}

// RoleToEntity は dto.RoleForm を entity.Role に変換する。
//
// args:
//   - f dto.RoleForm: 変換元フォーム DTO
//
// return:
//   - entity.Role: ロールエンティティ
func RoleToEntity(f dto.RoleForm) entity.Role {
	e := entity.Role{
		ID:   f.ID,
		Name: f.Name,
	}
	for _, uf := range f.Users {
		e.Users = append(e.Users, UserToEntity(uf))
	}
	return e
}

// =====================
// サポートステータス（entity.SupportStatus）
// =====================

// SupportStatusFromEntity は entity.SupportStatus を dto.SupportStatusForm に変換する。
//
// args:
//   - e entity.SupportStatus: 変換元エンティティ
//
// return:
//   - dto.SupportStatusForm: サポートステータスフォーム DTO
func SupportStatusFromEntity(e entity.SupportStatus) dto.SupportStatusForm {
	f := dto.SupportStatusForm{
		ID:        e.ID,
		Name:      e.Name,
		CreatedAt: timeToISO(e.CreatedAt),
		UpdatedAt: timeToISO(e.UpdatedAt),
		DeletedAt: deletedAtToISO(e.DeletedAt),
	}
	for _, s := range e.Supports {
		f.Supports = append(f.Supports, SupportFromEntity(s))
	}
	return f
}

// SupportStatusToEntity は dto.SupportStatusForm を entity.SupportStatus に変換する。
//
// args:
//   - f dto.SupportStatusForm: 変換元フォーム DTO
//
// return:
//   - entity.SupportStatus: サポートステータスエンティティ
func SupportStatusToEntity(f dto.SupportStatusForm) entity.SupportStatus {
	e := entity.SupportStatus{
		ID:   f.ID,
		Name: f.Name,
	}
	for _, sf := range f.Supports {
		e.Supports = append(e.Supports, SupportToEntity(sf))
	}
	return e
}

// =====================
// サポート（entity.Support）
// =====================

// SupportFromEntity は entity.Support を dto.SupportForm に変換する。
//
// args:
//   - e entity.Support: 変換元エンティティ
//
// return:
//   - dto.SupportForm: サポートフォーム DTO
func SupportFromEntity(e entity.Support) dto.SupportForm {
	f := dto.SupportForm{
		ID:              e.ID,
		UserID:          strconv.FormatInt(e.UserID, 10),
		SupportStatusID: strconv.FormatInt(e.SupportStatusID, 10),
		CreatedAt:       timeToISO(e.CreatedAt),
		UpdatedAt:       timeToISO(e.UpdatedAt),
		DeletedAt:       deletedAtToISO(e.DeletedAt),
	}
	if e.User.ID != 0 {
		u := UserFromEntity(e.User)
		f.User = &u
	}
	if e.SupportStatus.ID != 0 {
		ss := supportStatusFromEntityShallow(e.SupportStatus)
		f.SupportStatus = &ss
	}
	return f
}

// supportStatusFromEntityShallow は entity.SupportStatus を Supports なしの dto.SupportStatusForm に変換する。
//
// args:
//   - e entity.SupportStatus: 変換元エンティティ
//
// return:
//   - dto.SupportStatusForm: サポートステータスフォーム DTO（Supports なし）
func supportStatusFromEntityShallow(e entity.SupportStatus) dto.SupportStatusForm {
	return dto.SupportStatusForm{
		ID:        e.ID,
		Name:      e.Name,
		CreatedAt: timeToISO(e.CreatedAt),
		UpdatedAt: timeToISO(e.UpdatedAt),
		DeletedAt: deletedAtToISO(e.DeletedAt),
	}
}

// SupportToEntity は dto.SupportForm を entity.Support に変換する。
//
// args:
//   - f dto.SupportForm: 変換元フォーム DTO
//
// return:
//   - entity.Support: サポートエンティティ
func SupportToEntity(f dto.SupportForm) entity.Support {
	e := entity.Support{
		ID:              f.ID,
		UserID:          f.UserIDInt64(),
		SupportStatusID: f.SupportStatusIDInt64(),
	}
	if f.User != nil {
		e.User = UserToEntity(*f.User)
	}
	if f.SupportStatus != nil {
		e.SupportStatus = SupportStatusToEntity(*f.SupportStatus)
	}
	return e
}

// =====================
// ユーザー（entity.User）
// =====================

// UserFromEntity は entity.User を dto.UserForm に変換する。Role があればネストして変換する。
//
// args:
//   - e entity.User: 変換元エンティティ
//
// return:
//   - dto.UserForm: ユーザーフォーム DTO
func UserFromEntity(e entity.User) dto.UserForm {
	f := dto.UserForm{
		ID:        e.ID,
		Name:      e.Name,
		Email:     e.Email,
		Memo:      e.Memo,
		RoleID:    strconv.FormatInt(e.RoleID, 10),
		CreatedAt: timeToISO(e.CreatedAt),
		UpdatedAt: timeToISO(e.UpdatedAt),
		DeletedAt: deletedAtToISO(e.DeletedAt),
	}
	if e.Role.ID != 0 {
		r := roleFromEntityShallow(e.Role)
		f.Role = &r
	}
	return f
}

// UserFromEntityNoRole は entity.User を dto.UserForm に変換する。entity.Role.Users 埋め込み時の循環参照を避けるため Role は含めない。
//
// args:
//   - e entity.User: 変換元エンティティ
//
// return:
//   - dto.UserForm: ユーザーフォーム DTO（Role ネストなし）
func UserFromEntityNoRole(e entity.User) dto.UserForm {
	return dto.UserForm{
		ID:        e.ID,
		Name:      e.Name,
		Email:     e.Email,
		Memo:      e.Memo,
		RoleID:    strconv.FormatInt(e.RoleID, 10),
		CreatedAt: timeToISO(e.CreatedAt),
		UpdatedAt: timeToISO(e.UpdatedAt),
		DeletedAt: deletedAtToISO(e.DeletedAt),
	}
}

// UserToEntityNoRole は dto.UserForm を entity.User に変換する。Role の関連グラフは展開しない。
//
// args:
//   - f dto.UserForm: 変換元フォーム DTO
//
// return:
//   - entity.User: ユーザーエンティティ（Role は空構造体）
func UserToEntityNoRole(f dto.UserForm) entity.User {
	e := entity.User{
		ID:     f.ID,
		Name:   f.Name,
		Email:  f.Email,
		Memo:   f.Memo,
		RoleID: f.RoleIDInt64(),
	}
	if f.RoleID == "0" && f.Role != nil {
		e.RoleID = f.Role.ID
	}
	if f.Password != "" {
		e.Password = f.Password
	}
	e.Role = entity.Role{}
	return e
}

// UserToEntity は dto.UserForm を entity.User に変換する。Role があればネストして変換する。
//
// args:
//   - f dto.UserForm: 変換元フォーム DTO
//
// return:
//   - entity.User: ユーザーエンティティ
func UserToEntity(f dto.UserForm) entity.User {
	e := entity.User{
		ID:     f.ID,
		Name:   f.Name,
		Email:  f.Email,
		Memo:   f.Memo,
		RoleID: f.RoleIDInt64(),
	}
	if f.RoleID == "0" && f.Role != nil {
		e.RoleID = f.Role.ID
	}
	if f.Password != "" {
		e.Password = f.Password
	}
	if f.Role != nil {
		e.Role = RoleToEntity(*f.Role)
	}
	return e
}

// UserFormsFromEntities は entity.User のスライスを dto.UserForm のスライスに一括変換する。
//
// args:
//   - users []entity.User: 変換元エンティティ一覧
//
// return:
//   - []dto.UserForm: ユーザーフォーム DTO の一覧
func UserFormsFromEntities(users []entity.User) []dto.UserForm {
	out := make([]dto.UserForm, len(users))
	for i := range users {
		out[i] = UserFromEntity(users[i])
	}
	return out
}

// =====================
// カテゴリ（entity.Category）
// =====================

// categoryFromEntityShallow は entity.Category を Tags なしの dto.CategoryForm に変換する。
//
// args:
//   - e entity.Category: 変換元エンティティ
//
// return:
//   - dto.CategoryForm: カテゴリフォーム DTO（Tags なし）
func categoryFromEntityShallow(e entity.Category) dto.CategoryForm {
	return dto.CategoryForm{
		ID:        e.ID,
		Name:      e.Name,
		CreatedAt: timeToISO(e.CreatedAt),
		UpdatedAt: timeToISO(e.UpdatedAt),
		DeletedAt: deletedAtToISO(e.DeletedAt),
	}
}

// CategoryFromEntity は entity.Category を dto.CategoryForm に変換する。
//
// args:
//   - e entity.Category: 変換元エンティティ
//
// return:
//   - dto.CategoryForm: カテゴリフォーム DTO
func CategoryFromEntity(e entity.Category) dto.CategoryForm {
	f := categoryFromEntityShallow(e)
	for _, t := range e.Tags {
		f.Tags = append(f.Tags, TagFromEntity(t))
	}
	return f
}

// CategoryToEntity は dto.CategoryForm を entity.Category に変換する。
//
// args:
//   - f dto.CategoryForm: 変換元フォーム DTO
//
// return:
//   - entity.Category: カテゴリエンティティ
func CategoryToEntity(f dto.CategoryForm) entity.Category {
	e := entity.Category{
		ID:   f.ID,
		Name: f.Name,
	}
	for _, tf := range f.Tags {
		e.Tags = append(e.Tags, TagToEntity(tf))
	}
	return e
}

// =====================
// タグ（entity.Tag）
// =====================

// tagFromEntityShallow は entity.Tag を Questions なしの dto.TagForm に変換する。Category は浅い変換のみ。
//
// args:
//   - e entity.Tag: 変換元エンティティ
//
// return:
//   - dto.TagForm: タグフォーム DTO
func tagFromEntityShallow(e entity.Tag) dto.TagForm {
	f := dto.TagForm{
		ID:         e.ID,
		Name:       e.Name,
		Usage:      e.Usage,
		CategoryID: strconv.FormatInt(e.CategoryID, 10),
		CreatedAt:  timeToISO(e.CreatedAt),
		UpdatedAt:  timeToISO(e.UpdatedAt),
		DeletedAt:  deletedAtToISO(e.DeletedAt),
	}
	if e.Category.ID != 0 {
		c := categoryFromEntityShallow(e.Category)
		f.Category = &c
	}
	return f
}

// TagFromEntities は entity.Tag のスライスを dto.TagForm のスライスに一括変換する。
//
// args:
//   - e []entity.Tag: 変換元エンティティ一覧
//
// return:
//   - []dto.TagForm: タグフォーム DTO の一覧
func TagFromEntities(e []entity.Tag) []dto.TagForm {
	forms := []dto.TagForm{}
	for _, tag := range e {
		forms = append(forms, TagFromEntity(tag))
	}
	return forms
}

// TagFromEntity は entity.Tag を dto.TagForm に変換する。TagManagers 経由で関連質問も含める。
//
// args:
//   - e entity.Tag: 変換元エンティティ
//
// return:
//   - dto.TagForm: タグフォーム DTO
func TagFromEntity(e entity.Tag) dto.TagForm {
	f := tagFromEntityShallow(e)
	for _, tm := range e.TagManagers {
		if tm.Question.ID != 0 {
			f.Questions = append(f.Questions, QuestionFromEntity(tm.Question))
		}
	}
	return f
}

// TagToEntity は dto.TagForm を entity.Tag に変換する。Category の明示的関連はセットしない。
//
// args:
//   - f dto.TagForm: 変換元フォーム DTO
//
// return:
//   - entity.Tag: タグエンティティ
func TagToEntity(f dto.TagForm) entity.Tag {
	e := entity.Tag{
		ID:         f.ID,
		Name:       f.Name,
		Usage:      f.Usage,
		CategoryID: f.CategoryIDInt64(),
	}
	// DB に余分な category が入らないよう、明示的関連はセットしない。
	// if f.Category != nil {
	// 	e.Category = CategoryToEntity(*f.Category)
	// }
	for _, qf := range f.Questions {
		tm := entity.TagManager{
			TagID:      f.ID,
			QuestionID: qf.ID,
		}
		if qf.ID != 0 {
			tm.Question = entity.Question{ID: qf.ID}
		}
		e.TagManagers = append(e.TagManagers, tm)
	}
	return e
}

// TagToEntityNoRelations は dto.TagForm を関連なしの entity.Tag に変換する。
//
// args:
//   - f dto.TagForm: 変換元フォーム DTO
//
// return:
//   - entity.Tag: タグエンティティ（Category は空構造体）
func TagToEntityNoRelations(f dto.TagForm) entity.Tag {
	e := entity.Tag{
		ID:         f.ID,
		Name:       f.Name,
		Usage:      f.Usage,
		CategoryID: f.CategoryIDInt64(),
	}
	e.Category = entity.Category{}

	return e
}

// =====================
// 参照リンク（entity.Refer）
// =====================

// referFromEntityShallow は entity.Refer を Answers なしの dto.ReferForm に変換する。
//
// args:
//   - e entity.Refer: 変換元エンティティ
//
// return:
//   - dto.ReferForm: 参照リンクフォーム DTO
func referFromEntityShallow(e entity.Refer) dto.ReferForm {
	return dto.ReferForm{
		ID:        e.ID,
		Title:     e.Title,
		URL:       e.URL,
		CreatedAt: timeToISO(e.CreatedAt),
		UpdatedAt: timeToISO(e.UpdatedAt),
		DeletedAt: deletedAtToISO(e.DeletedAt),
	}
}

// ReferFromEntity は entity.Refer を dto.ReferForm に変換する。
//
// args:
//   - e entity.Refer: 変換元エンティティ
//
// return:
//   - dto.ReferForm: 参照リンクフォーム DTO
func ReferFromEntity(e entity.Refer) dto.ReferForm {
	f := referFromEntityShallow(e)
	for _, rm := range e.ReferManagers {
		if rm.Answer.ID != 0 {
			f.Answers = append(f.Answers, AnswerFromEntity(rm.Answer))
		}
	}
	return f
}

// ReferToEntity は dto.ReferForm を entity.Refer に変換する。
//
// args:
//   - f dto.ReferForm: 変換元フォーム DTO
//
// return:
//   - entity.Refer: 参照リンクエンティティ
func ReferToEntity(f dto.ReferForm) entity.Refer {
	e := entity.Refer{
		ID:    f.ID,
		Title: f.Title,
		URL:   f.URL,
	}
	referID := f.ID
	for _, af := range f.Answers {
		rm := entity.ReferManager{
			ReferID:  referID,
			AnswerID: af.ID,
		}
		if af.ID != 0 {
			rm.Answer = entity.Answer{ID: af.ID}
		}
		e.ReferManagers = append(e.ReferManagers, rm)
	}
	return e
}

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

// =====================
// 回答（entity.Answer）
// =====================

// AnswerFromEntity は entity.Answer を dto.AnswerForm に変換する。
//
// args:
//   - e entity.Answer: 変換元エンティティ
//
// return:
//   - dto.AnswerForm: 回答フォーム DTO
func AnswerFromEntity(e entity.Answer) dto.AnswerForm {
	f := dto.AnswerForm{
		ID:        e.ID,
		UserID:    strconv.FormatInt(e.UserID, 10),
		Content:   e.Content,
		IsFinal:   e.IsFinal,
		CreatedAt: timeToISO(e.CreatedAt),
		UpdatedAt: timeToISO(e.UpdatedAt),
		DeletedAt: deletedAtToISO(e.DeletedAt),
	}
	if e.User.ID != 0 {
		u := UserFromEntity(e.User)
		f.User = &u
	}
	for _, rm := range e.ReferManagers {
		if rm.Refer.ID != 0 {
			f.Refers = append(f.Refers, referFromEntityShallow(rm.Refer))
		}
	}
	return f
}

// AnswerToEntity は dto.AnswerForm を entity.Answer に変換する。CreatedAt 未指定時は現在時刻を設定する。
//
// args:
//   - f dto.AnswerForm: 変換元フォーム DTO
//
// return:
//   - entity.Answer: 回答エンティティ
func AnswerToEntity(f dto.AnswerForm) entity.Answer {
	e := entity.Answer{
		ID:      f.ID,
		UserID:  f.UserIDInt64(),
		Content: f.Content,
		IsFinal: f.IsFinal,
	}
	if f.CreatedAt == nil || (f.CreatedAt != nil && *f.CreatedAt == "") {
		e.CreatedAt = time.Now()
	} else {
		e.CreatedAt = isoToTime(f.CreatedAt)
	}
	if f.User != nil {
		e.User = UserToEntity(*f.User)
	}
	answerID := f.ID
	for _, rf := range f.Refers {
		rm := entity.ReferManager{
			AnswerID: answerID,
			ReferID:  rf.ID,
		}
		if rf.ID != 0 {
			rm.Refer = entity.Refer{ID: rf.ID}
		} else if strings.TrimSpace(rf.Title) != "" && strings.TrimSpace(rf.URL) != "" {
			rm.Refer = ReferToEntity(rf)
		}
		e.ReferManagers = append(e.ReferManagers, rm)
	}
	return e
}

// =====================
// エスカレーション（entity.Escalation）
// =====================

// EscalationFromEntity は entity.Escalation を dto.EscalationForm に変換する。
//
// args:
//   - e entity.Escalation: 変換元エンティティ
//
// return:
//   - dto.EscalationForm: エスカレーションフォーム DTO
func EscalationFromEntity(e entity.Escalation) dto.EscalationForm {
	return dto.EscalationForm{
		ID:             e.ID,
		FromQuestionID: strconv.FormatInt(e.FromQuestionID, 10),
		ToQuestionID:   strconv.FormatInt(e.ToQuestionID, 10),
		EscalatedAt:    timeToISO(e.EscalatedAt),
		CreatedAt:      timeToISO(e.CreatedAt),
		UpdatedAt:      timeToISO(e.UpdatedAt),
		DeletedAt:      deletedAtToISO(e.DeletedAt),
	}
}

// EscalationToEntity は dto.EscalationForm を entity.Escalation に変換する。EscalatedAt 未指定時は現在時刻を設定する。
//
// args:
//   - f dto.EscalationForm: 変換元フォーム DTO
//
// return:
//   - entity.Escalation: エスカレーションエンティティ
func EscalationToEntity(f dto.EscalationForm) entity.Escalation {
	e := entity.Escalation{
		ID:             f.ID,
		FromQuestionID: f.FromQuestionIDInt64(),
		ToQuestionID:   f.ToQuestionIDInt64(),
	}
	if f.EscalatedAt == nil || *f.EscalatedAt == "" {
		e.EscalatedAt = time.Now()
	} else {
		e.EscalatedAt = isoToTime(f.EscalatedAt)
	}
	return e
}

// =====================
// 通知種別・通知（entity.NoticeType / entity.Notice）
// =====================

// noticeTypeFromEntityShallow は entity.NoticeType を Notices なしの dto.NoticeTypeForm に変換する。
//
// args:
//   - e entity.NoticeType: 変換元エンティティ
//
// return:
//   - dto.NoticeTypeForm: 通知種別フォーム DTO
func noticeTypeFromEntityShallow(e entity.NoticeType) dto.NoticeTypeForm {
	return dto.NoticeTypeForm{
		ID:        e.ID,
		Name:      e.Name,
		CreatedAt: timeToISO(e.CreatedAt),
		UpdatedAt: timeToISO(e.UpdatedAt),
		DeletedAt: deletedAtToISO(e.DeletedAt),
	}
}

// NoticeTypeFromEntity は entity.NoticeType を dto.NoticeTypeForm に変換する。
//
// args:
//   - e entity.NoticeType: 変換元エンティティ
//
// return:
//   - dto.NoticeTypeForm: 通知種別フォーム DTO
func NoticeTypeFromEntity(e entity.NoticeType) dto.NoticeTypeForm {
	f := noticeTypeFromEntityShallow(e)
	for _, n := range e.Notices {
		f.Notices = append(f.Notices, NoticeFromEntity(n))
	}
	return f
}

// NoticeTypeToEntity は dto.NoticeTypeForm を entity.NoticeType に変換する。
//
// args:
//   - f dto.NoticeTypeForm: 変換元フォーム DTO
//
// return:
//   - entity.NoticeType: 通知種別エンティティ
func NoticeTypeToEntity(f dto.NoticeTypeForm) entity.NoticeType {
	e := entity.NoticeType{
		ID:   f.ID,
		Name: f.Name,
	}
	for _, nf := range f.Notices {
		e.Notices = append(e.Notices, NoticeToEntity(nf))
	}
	return e
}

// NoticeFromEntity は entity.Notice を dto.NoticeForm に変換する。
//
// args:
//   - e entity.Notice: 変換元エンティティ
//
// return:
//   - dto.NoticeForm: 通知フォーム DTO
func NoticeFromEntity(e entity.Notice) dto.NoticeForm {
	var questionID *string
	if e.QuestionID != nil {
		s := strconv.FormatInt(*e.QuestionID, 10)
		questionID = &s
	}
	f := dto.NoticeForm{
		ID:         e.ID,
		TypeID:     e.TypeID,
		QuestionID: questionID,
		Content:    e.Content,
		DisplayDue: timePtrToISO(e.DisplayDue),
		CreatedAt:  timeToISO(e.CreatedAt),
		UpdatedAt:  timeToISO(e.UpdatedAt),
		DeletedAt:  deletedAtToISO(e.DeletedAt),
	}
	if e.NoticeType.ID != 0 {
		nt := noticeTypeFromEntityShallow(e.NoticeType)
		f.NoticeType = &nt
	}
	if e.Question != nil && e.Question.ID != 0 {
		// due := e.Question.Due.Format(time.DateTime)
		// const tag = TagFromEntity(e.Question.TagManagers)
		// qf := dto.QuestionForm{ID: e.Question.ID, Title: e.Question.Title, Content: e.Question.Content, Due: &due, Tags: tag}
		qf := QuestionFromEntity(*e.Question)
		f.Question = &qf
	}
	return f
}

// NoticeFromEntities は entity.Notice のスライスを dto.NoticeForm のスライスに一括変換する。
//
// args:
//   - ns []entity.Notice: 変換元エンティティ一覧
//
// return:
//   - []dto.NoticeForm: 通知フォーム DTO の一覧
func NoticeFromEntities(ns []entity.Notice) (forms []dto.NoticeForm) {
	for _, n := range ns {
		forms = append(forms, NoticeFromEntity(n))
	}
	return forms
}

// NoticeToEntity は dto.NoticeForm を entity.Notice に変換する。
//
// args:
//   - f dto.NoticeForm: 変換元フォーム DTO
//
// return:
//   - entity.Notice: 通知エンティティ
func NoticeToEntity(f dto.NoticeForm) entity.Notice {
	var questionID *int64
	if v := f.QuestionIDInt64(); v >= 0 {
		questionID = &v
	}
	e := entity.Notice{
		ID:         f.ID,
		TypeID:     f.TypeID,
		QuestionID: questionID,
		Content:    f.Content,
		DisplayDue: isoToTimePtr(f.DisplayDue),
	}
	if f.NoticeType != nil {
		e.NoticeType = NoticeTypeToEntity(*f.NoticeType)
	}
	if f.Question != nil {
		q := QuestionToEntity(*f.Question)
		e.Question = &q
	}
	return e
}

// QuestionFromEntities は entity.Question のスライスを dto.QuestionForm のスライスに一括変換する。
//
// args:
//   - e []entity.Question: 変換元エンティティ一覧
//
// return:
//   - []dto.QuestionForm: 質問フォーム DTO の一覧
func QuestionFromEntities(e []entity.Question) []dto.QuestionForm {
	forms := []dto.QuestionForm{}
	for _, q := range e {
		forms = append(forms, QuestionFromEntity(q))
	}
	return forms
}

// =====================
// 質問（entity.Question）
// =====================

// QuestionFromEntity は entity.Question を dto.QuestionForm に変換する。回答・メモ・タグ・関連質問などをネストする。
//
// args:
//   - e entity.Question: 変換元エンティティ
//
// return:
//   - dto.QuestionForm: 質問フォーム DTO
func QuestionFromEntity(e entity.Question) dto.QuestionForm {
	var originQuestionID *string
	if e.OriginQuestionID != nil {
		s := strconv.FormatInt(*e.OriginQuestionID, 10)
		originQuestionID = &s
	}
	f := dto.QuestionForm{
		ID:               e.ID,
		OriginQuestionID: originQuestionID,
		SupportID:        e.SupportID,
		Title:            e.Title,
		Content:          e.Content,
		Due:              timePtrToISO(e.Due),
		CreatedAt:        timeToISO(e.CreatedAt),
		UpdatedAt:        timeToISO(e.UpdatedAt),
		DeletedAt:        deletedAtToISO(e.DeletedAt),
	}
	for _, answer := range e.Answer {
		if answer.ID == 0 {
			continue
		}
		f.Answers = append(f.Answers, AnswerFromEntity(answer))
	}
	for _, m := range e.Memos {
		f.Memos = append(f.Memos, MemoFromEntity(m))
	}
	for _, tm := range e.TagManagers {
		if tm.Tag.ID != 0 {
			f.Tags = append(f.Tags, tagFromEntityShallow(tm.Tag))
		}
	}
	seenRelated := make(map[int64]struct{})
	for _, rq := range e.RelatedQuestions {
		rid := rq.RelatedQuestionID
		if rid == 0 || rid == e.ID {
			continue
		}
		if _, ok := seenRelated[rid]; ok {
			continue
		}
		seenRelated[rid] = struct{}{}
		f.RelatedQuestions = append(f.RelatedQuestions, RelatedQuestionFromEntity(rq))
	}
	if e.Support != nil && e.Support.ID != 0 {
		s := SupportFromEntity(*e.Support)
		f.Support = &s
	}
	for _, st := range e.SenderTalks {
		f.SenderTalks = append(f.SenderTalks, SenderTalkFromEntity(st))
	}
	return f
}

// QuestionToEntity は dto.QuestionForm を entity.Question に変換する。CreatedAt 未指定時は現在時刻を設定する。
//
// args:
//   - f dto.QuestionForm: 変換元フォーム DTO
//
// return:
//   - entity.Question: 質問エンティティ
func QuestionToEntity(f dto.QuestionForm) entity.Question {
	var originQuestionID *int64
	if v := f.OriginQuestionIDInt64(); v >= 0 {
		originQuestionID = &v
	}
	e := entity.Question{
		ID:               f.ID,
		OriginQuestionID: originQuestionID,
		SupportID:        f.SupportID,
		Title:            f.Title,
		Content:          f.Content,
		Due:              isoToTimePtr(f.Due),
	}
	if f.CreatedAt == nil || *f.CreatedAt == "" {
		e.CreatedAt = time.Now()
	} else {
		e.CreatedAt = isoToTime(f.CreatedAt)
	}
	qid := f.ID
	for _, af := range f.Answers {
		a := AnswerToEntity(af)
		if a.QuestionID == 0 {
			a.QuestionID = qid
		}
		e.Answer = append(e.Answer, a)
	}
	for _, mf := range f.Memos {
		m := MemoToEntity(mf)
		if m.QuestionID == 0 {
			m.QuestionID = qid
		}
		e.Memos = append(e.Memos, m)
	}
	for _, tf := range f.Tags {
		if tf.ID == 0 {
			continue
		}
		tm := entity.TagManager{
			QuestionID: qid,
			TagID:      tf.ID,
			Tag:        entity.Tag{ID: tf.ID},
		}
		e.TagManagers = append(e.TagManagers, tm)
	}
	seenRelated := make(map[int64]struct{})
	for _, rf := range f.RelatedQuestions {
		rq := RelatedQuestionToEntity(rf, qid)
		if rq.RelatedQuestionID == 0 || rq.RelatedQuestionID == qid {
			continue
		}
		if _, ok := seenRelated[rq.RelatedQuestionID]; ok {
			continue
		}
		seenRelated[rq.RelatedQuestionID] = struct{}{}
		e.RelatedQuestions = append(e.RelatedQuestions, rq)
	}
	if f.Support != nil {
		sup := SupportToEntity(*f.Support)
		e.Support = &sup
	}
	for _, sf := range f.SenderTalks {
		st := SenderTalkToEntity(sf)
		if st.QuestionID == 0 {
			st.QuestionID = qid
		}
		e.SenderTalks = append(e.SenderTalks, st)
	}
	return e
}

// questionFormShallowFromEntity は entity.Question を関連グラフなしの dto.QuestionForm に変換する。
//
// args:
//   - e entity.Question: 変換元エンティティ
//
// return:
//   - dto.QuestionForm: 質問フォーム DTO（ネスト関連なし）
func questionFormShallowFromEntity(e entity.Question) dto.QuestionForm {
	var originQuestionID *string
	if e.OriginQuestionID != nil {
		s := strconv.FormatInt(*e.OriginQuestionID, 10)
		originQuestionID = &s
	}
	return dto.QuestionForm{
		ID:               e.ID,
		OriginQuestionID: originQuestionID,
		SupportID:        e.SupportID,
		Title:            e.Title,
		Content:          e.Content,
		Due:              timePtrToISO(e.Due),
		CreatedAt:        timeToISO(e.CreatedAt),
		UpdatedAt:        timeToISO(e.UpdatedAt),
		DeletedAt:        deletedAtToISO(e.DeletedAt),
	}
}

// RelatedQuestionFromEntity は entity.RelatedQuestion を dto.RelatedQuestionForm に変換する。
//
// args:
//   - r entity.RelatedQuestion: 変換元エンティティ
//
// return:
//   - dto.RelatedQuestionForm: 関連質問フォーム DTO
func RelatedQuestionFromEntity(r entity.RelatedQuestion) dto.RelatedQuestionForm {
	f := dto.RelatedQuestionForm{
		ID:                r.ID,
		QuestionID:        strconv.FormatInt(r.QuestionID, 10),
		RelatedQuestionID: strconv.FormatInt(r.RelatedQuestionID, 10),
		CreatedAt:         timeToISO(r.CreatedAt),
		UpdatedAt:         timeToISO(r.UpdatedAt),
		DeletedAt:         deletedAtToISO(r.DeletedAt),
	}
	if r.RelatedQuestion.ID != 0 {
		q := questionFormShallowFromEntity(r.RelatedQuestion)
		f.RelatedQuestion = &q
	}
	return f
}

// RelatedQuestionToEntity は dto.RelatedQuestionForm を entity.RelatedQuestion に変換する。
//
// args:
//   - f dto.RelatedQuestionForm: 変換元フォーム DTO
//   - parentQuestionID int64: 親質問 ID（フォームの QuestionID が無効な場合に使用）
//
// return:
//   - entity.RelatedQuestion: 関連質問エンティティ
func RelatedQuestionToEntity(f dto.RelatedQuestionForm, parentQuestionID int64) entity.RelatedQuestion {
	qid := f.QuestionIDInt64()
	if qid < 0 || qid == 0 {
		qid = parentQuestionID
	}
	rid := f.RelatedQuestionIDInt64()
	if rid < 0 {
		rid = 0
	}
	if rid == 0 && f.RelatedQuestion != nil && f.RelatedQuestion.ID != 0 {
		rid = f.RelatedQuestion.ID
	}
	e := entity.RelatedQuestion{
		ID:                f.ID,
		QuestionID:        qid,
		RelatedQuestionID: rid,
	}
	if f.RelatedQuestion != nil && f.RelatedQuestion.ID != 0 {
		e.RelatedQuestion = QuestionToEntity(*f.RelatedQuestion)
	}
	return e
}

// =====================
// 参照リンク管理・タグ管理（entity.ReferManager / entity.TagManager）
// =====================

// ReferManagerFromEntity は entity.ReferManager を dto.ReferManagerForm に変換する。
//
// args:
//   - e entity.ReferManager: 変換元エンティティ
//
// return:
//   - dto.ReferManagerForm: 参照リンク管理フォーム DTO
func ReferManagerFromEntity(e entity.ReferManager) dto.ReferManagerForm {
	f := dto.ReferManagerForm{
		ID:        e.ID,
		AnswerID:  strconv.FormatInt(e.AnswerID, 10),
		ReferID:   strconv.FormatInt(e.ReferID, 10),
		CreatedAt: timeToISO(e.CreatedAt),
		UpdatedAt: timeToISO(e.UpdatedAt),
		DeletedAt: deletedAtToISO(e.DeletedAt),
	}
	if e.Answer.ID != 0 {
		a := AnswerFromEntity(e.Answer)
		f.Answer = &a
	}
	if e.Refer.ID != 0 {
		r := ReferFromEntity(e.Refer)
		f.Refer = &r
	}
	return f
}

// ReferManagerToEntity は dto.ReferManagerForm を entity.ReferManager に変換する。
//
// args:
//   - f dto.ReferManagerForm: 変換元フォーム DTO
//
// return:
//   - entity.ReferManager: 参照リンク管理エンティティ
func ReferManagerToEntity(f dto.ReferManagerForm) entity.ReferManager {
	e := entity.ReferManager{
		ID:       f.ID,
		AnswerID: f.AnswerIDInt64(),
		ReferID:  f.ReferIDInt64(),
	}
	if f.Answer != nil {
		e.Answer = AnswerToEntity(*f.Answer)
	}
	if f.Refer != nil {
		e.Refer = ReferToEntity(*f.Refer)
	}
	return e
}

// TagManagerFromEntity は entity.TagManager を dto.TagManagerForm に変換する。
//
// args:
//   - e entity.TagManager: 変換元エンティティ
//
// return:
//   - dto.TagManagerForm: タグ管理フォーム DTO
func TagManagerFromEntity(e entity.TagManager) dto.TagManagerForm {
	f := dto.TagManagerForm{
		ID:         e.ID,
		TagID:      strconv.FormatInt(e.TagID, 10),
		QuestionID: strconv.FormatInt(e.QuestionID, 10),
		CreatedAt:  timeToISO(e.CreatedAt),
		UpdatedAt:  timeToISO(e.UpdatedAt),
		DeletedAt:  deletedAtToISO(e.DeletedAt),
	}
	if e.Tag.ID != 0 {
		t := TagFromEntity(e.Tag)
		f.Tag = &t
	}
	if e.Question.ID != 0 {
		q := QuestionFromEntity(e.Question)
		f.Question = &q
	}
	return f
}

// TagManagerToEntity は dto.TagManagerForm を entity.TagManager に変換する。
//
// args:
//   - f dto.TagManagerForm: 変換元フォーム DTO
//
// return:
//   - entity.TagManager: タグ管理エンティティ
func TagManagerToEntity(f dto.TagManagerForm) entity.TagManager {
	e := entity.TagManager{
		ID:         f.ID,
		TagID:      f.TagIDInt64(),
		QuestionID: f.QuestionIDInt64(),
	}
	if f.Tag != nil {
		e.Tag = TagToEntity(*f.Tag)
	}
	if f.Question != nil {
		e.Question = QuestionToEntity(*f.Question)
	}
	return e
}

// =====================
// 送信者・送信者トーク（entity.Sender / entity.SenderTalk）
// =====================

// SenderFromEntity は entity.Sender を dto.SenderForm に変換する。
//
// args:
//   - e entity.Sender: 変換元エンティティ
//
// return:
//   - dto.SenderForm: 送信者フォーム DTO
func SenderFromEntity(e entity.Sender) dto.SenderForm {
	f := dto.SenderForm{
		ID:             e.ID,
		Name:           e.Name,
		DepartmentName: e.DepartmentName,
	}
	for _, st := range e.SenderTalks {
		f.SenderTalks = append(f.SenderTalks, SenderTalkFromEntityNoSender(st))
	}
	return f
}

// senderFromEntityShallow は entity.Sender を SenderTalks なしの dto.SenderForm に変換する。
//
// args:
//   - e entity.Sender: 変換元エンティティ
//
// return:
//   - dto.SenderForm: 送信者フォーム DTO（SenderTalks なし）
func senderFromEntityShallow(e entity.Sender) dto.SenderForm {
	return dto.SenderForm{
		ID:             e.ID,
		Name:           e.Name,
		DepartmentName: e.DepartmentName,
	}
}

// SenderToEntity は dto.SenderForm を entity.Sender に変換する。
//
// args:
//   - f dto.SenderForm: 変換元フォーム DTO
//
// return:
//   - entity.Sender: 送信者エンティティ
func SenderToEntity(f dto.SenderForm) entity.Sender {
	e := entity.Sender{
		ID:             f.ID,
		Name:           f.Name,
		DepartmentName: f.DepartmentName,
	}
	for _, stf := range f.SenderTalks {
		e.SenderTalks = append(e.SenderTalks, SenderTalkToEntity(stf))
	}
	return e
}

// SenderTalkFromEntity は entity.SenderTalk を dto.SenderTalkForm に変換する。Sender があればネストする。
//
// args:
//   - e entity.SenderTalk: 変換元エンティティ
//
// return:
//   - dto.SenderTalkForm: 送信者トークフォーム DTO
func SenderTalkFromEntity(e entity.SenderTalk) dto.SenderTalkForm {
	f := SenderTalkFromEntityNoSender(e)
	if e.Sender.ID != 0 {
		s := senderFromEntityShallow(e.Sender)
		f.Sender = &s
	}
	return f
}

// SenderTalkFromEntityNoSender は entity.SenderTalk を Sender なしの dto.SenderTalkForm に変換する。
//
// args:
//   - e entity.SenderTalk: 変換元エンティティ
//
// return:
//   - dto.SenderTalkForm: 送信者トークフォーム DTO（Sender なし）
func SenderTalkFromEntityNoSender(e entity.SenderTalk) dto.SenderTalkForm {
	return dto.SenderTalkForm{
		ID:         e.ID,
		Content:    e.Content,
		SenderID:   strconv.FormatInt(e.SenderID, 10),
		QuestionID: strconv.FormatInt(e.QuestionID, 10),
		CreatedAt:  timeToISO(e.CreatedAt),
		UpdatedAt:  timeToISO(e.UpdatedAt),
		DeletedAt:  deletedAtToISO(e.DeletedAt),
	}
}

// SenderTalkToEntity は dto.SenderTalkForm を entity.SenderTalk に変換する。
//
// args:
//   - f dto.SenderTalkForm: 変換元フォーム DTO
//
// return:
//   - entity.SenderTalk: 送信者トークエンティティ
func SenderTalkToEntity(f dto.SenderTalkForm) entity.SenderTalk {
	e := entity.SenderTalk{
		ID:         f.ID,
		Content:    f.Content,
		SenderID:   f.SenderIDInt64(),
		QuestionID: f.QuestionIDInt64(),
	}
	if f.Sender != nil {
		e.Sender = SenderToEntity(*f.Sender)
		if e.SenderID <= 0 {
			e.SenderID = e.Sender.ID
		}
	}
	return e
}
