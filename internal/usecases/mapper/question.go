package mapper

import (
	"strconv"
	"time"

	"meetup/internal/domains/entity"
	"meetup/internal/usecases/dto"
)

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
