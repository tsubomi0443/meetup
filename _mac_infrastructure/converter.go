package infrastructure

import (
	"strconv"
	"time"
)

// --- time helpers (ISO8601 / RFC3339) ---

func timePtrToISO(t *time.Time) *string {
	if t == nil || t.IsZero() {
		return nil
	}
	s := t.UTC().Format(time.RFC3339Nano)
	return &s
}

func timeToISO(t time.Time) *string {
	if t.IsZero() {
		return nil
	}
	s := t.UTC().Format(time.RFC3339Nano)
	return &s
}

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

// --- Role ---

func RoleFromEntity(e Role) RoleForm {
	f := RoleForm{ID: e.ID, RoleName: e.RoleName}
	if len(e.Users) > 0 {
		f.Users = make([]UserForm, len(e.Users))
		for i := range e.Users {
			f.Users[i] = UserFromEntityNoRole(e.Users[i])
		}
	}
	return f
}

func roleFromEntityShallow(e Role) RoleForm {
	return RoleForm{ID: e.ID, RoleName: e.RoleName}
}

func RoleToEntity(f RoleForm) Role {
	e := Role{ID: f.ID, RoleName: f.RoleName}
	for _, uf := range f.Users {
		e.Users = append(e.Users, UserToEntity(uf))
	}
	return e
}

// --- SupportStatus ---

func SupportStatusFromEntity(e SupportStatus) SupportStatusForm {
	f := SupportStatusForm{ID: e.ID, Title: e.Title}
	for _, s := range e.Supports {
		f.Supports = append(f.Supports, SupportFromEntity(s))
	}
	return f
}

func SupportStatusToEntity(f SupportStatusForm) SupportStatus {
	e := SupportStatus{ID: f.ID, Title: f.Title}
	for _, sf := range f.Supports {
		e.Supports = append(e.Supports, SupportToEntity(sf))
	}
	return e
}

// --- Support ---

func SupportFromEntity(e Support) SupportForm {
	f := SupportForm{
		ID:              e.ID,
		UserID:          strconv.FormatInt(e.UserID, 10),
		SupportStatusID: strconv.FormatInt(e.SupportStatusID, 10),
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

func supportStatusFromEntityShallow(e SupportStatus) SupportStatusForm {
	return SupportStatusForm{ID: e.ID, Title: e.Title}
}

func SupportToEntity(f SupportForm) Support {
	e := Support{
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

// --- User ---

func UserFromEntity(e User) UserForm {
	f := UserForm{
		ID:     e.ID,
		Name:   e.Name,
		Email:  e.Email,
		Memo:   e.Memo,
		RoleID: strconv.FormatInt(e.RoleID, 10),
	}
	if e.Role.ID != 0 {
		r := roleFromEntityShallow(e.Role)
		f.Role = &r
	}
	return f
}

// UserFromEntityNoRole avoids Role when embedding User under Role.Users.
func UserFromEntityNoRole(e User) UserForm {
	return UserForm{
		ID:     e.ID,
		Name:   e.Name,
		Email:  e.Email,
		Memo:   e.Memo,
		RoleID: strconv.FormatInt(e.RoleID, 10),
	}
}

func UserToEntityNoRole(f UserForm) User {
	e := User{
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
	e.Role = Role{}
	return e
}

func UserToEntity(f UserForm) User {
	e := User{
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

func UserFormsFromEntities(users []User) []UserForm {
	out := make([]UserForm, len(users))
	for i := range users {
		out[i] = UserFromEntity(users[i])
	}
	return out
}

// --- Category ---

func categoryFromEntityShallow(e Category) CategoryForm {
	return CategoryForm{ID: e.ID, CategoryName: e.CategoryName}
}

func CategoryFromEntity(e Category) CategoryForm {
	f := categoryFromEntityShallow(e)
	for _, t := range e.Tags {
		f.Tags = append(f.Tags, TagFromEntity(t))
	}
	return f
}

func CategoryToEntity(f CategoryForm) Category {
	e := Category{ID: f.ID, CategoryName: f.CategoryName}
	for _, tf := range f.Tags {
		e.Tags = append(e.Tags, TagToEntity(tf))
	}
	return e
}

// --- Tag ---

func tagFromEntityShallow(e Tag) TagForm {
	f := TagForm{
		ID:         e.ID,
		Title:      e.Title,
		Usage:      e.Usage,
		CategoryID: strconv.FormatInt(e.CategoryID, 10),
	}
	if e.Category.ID != 0 {
		c := categoryFromEntityShallow(e.Category)
		f.Category = &c
	}
	return f
}

func TagFromEntities(e []Tag) []TagForm {
	forms := []TagForm{}
	for _, tag := range e {
		forms = append(forms, TagFromEntity(tag))
	}
	return forms
}

func TagFromEntity(e Tag) TagForm {
	f := tagFromEntityShallow(e)
	for _, tm := range e.TagManagers {
		if tm.Question.ID != 0 {
			f.Questions = append(f.Questions, QuestionFromEntity(tm.Question))
		}
	}
	return f
}

func TagToEntity(f TagForm) Tag {
	e := Tag{
		ID:         f.ID,
		Title:      f.Title,
		Usage:      f.Usage,
		CategoryID: f.CategoryIDInt64(),
	}
	// DB に追加の category が入らないよう、明示的関連はセットしない。
	// if f.Category != nil {
	// 	e.Category = CategoryToEntity(*f.Category)
	// }
	for _, qf := range f.Questions {
		tm := TagManager{
			TagID:      f.ID,
			QuestionID: qf.ID,
		}
		if qf.ID != 0 {
			tm.Question = Question{ID: qf.ID}
		}
		e.TagManagers = append(e.TagManagers, tm)
	}
	return e
}

func TagToEntityNoRelations(f TagForm) Tag {
	e := Tag{
		ID:         f.ID,
		Title:      f.Title,
		Usage:      f.Usage,
		CategoryID: f.CategoryIDInt64(),
	}
	e.Category = Category{}

	return e
}

// --- Refer ---

func referFromEntityShallow(e Refer) ReferForm {
	return ReferForm{ID: e.ID, Title: e.Title, URL: e.URL}
}

func ReferFromEntity(e Refer) ReferForm {
	f := referFromEntityShallow(e)
	for _, rm := range e.ReferManagers {
		if rm.Answer.ID != 0 {
			f.Answers = append(f.Answers, AnswerFromEntity(rm.Answer))
		}
	}
	return f
}

func ReferToEntity(f ReferForm) Refer {
	e := Refer{ID: f.ID, Title: f.Title, URL: f.URL}
	referID := f.ID
	for _, af := range f.Answers {
		rm := ReferManager{
			ReferID:  referID,
			AnswerID: af.ID,
		}
		if af.ID != 0 {
			rm.Answer = Answer{ID: af.ID}
		}
		e.ReferManagers = append(e.ReferManagers, rm)
	}
	return e
}

// --- Memo ---

func MemoFromEntity(e Memo) MemoForm {
	f := MemoForm{
		ID:         e.ID,
		QuestionID: strconv.FormatInt(e.QuestionID, 10),
		UserID:     strconv.FormatInt(e.UserID, 10),
		Content:    e.Content,
	}
	if e.User.ID != 0 {
		u := UserFromEntity(e.User)
		f.User = &u
	}
	return f
}

func MemoToEntity(f MemoForm) Memo {
	e := Memo{
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

// --- Answer ---

func AnswerFromEntity(e Answer) AnswerForm {
	f := AnswerForm{
		ID:         e.ID,
		UserID:     strconv.FormatInt(e.UserID, 10),
		Content:    e.Content,
		AnsweredAt: timePtrToISO(e.AnsweredAt),
		CreatedAt:  timeToISO(e.CreatedAt),
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

func AnswerToEntity(f AnswerForm) Answer {
	e := Answer{
		ID:         f.ID,
		UserID:     f.UserIDInt64(),
		Content:    f.Content,
		AnsweredAt: isoToTimePtr(f.AnsweredAt),
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
		rm := ReferManager{
			AnswerID: answerID,
			ReferID:  rf.ID,
		}
		if rf.ID != 0 {
			rm.Refer = Refer{ID: rf.ID}
		}
		e.ReferManagers = append(e.ReferManagers, rm)
	}
	return e
}

// --- Escalation ---

func EscalationFromEntity(e Escalation) EscalationForm {
	return EscalationForm{
		ID:             e.ID,
		FromQuestionID: strconv.FormatInt(e.FromQuestionID, 10),
		ToQuestionID:   strconv.FormatInt(e.ToQuestionID, 10),
		EscalatedAt:    timeToISO(e.EscalatedAt),
	}
}

func EscalationToEntity(f EscalationForm) Escalation {
	e := Escalation{
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

// --- NoticeType / Notice ---

func noticeTypeFromEntityShallow(e NoticeType) NoticeTypeForm {
	return NoticeTypeForm{ID: e.ID, Name: e.Name}
}

func NoticeTypeFromEntity(e NoticeType) NoticeTypeForm {
	f := noticeTypeFromEntityShallow(e)
	for _, n := range e.Notices {
		f.Notices = append(f.Notices, NoticeFromEntity(n))
	}
	return f
}

func NoticeTypeToEntity(f NoticeTypeForm) NoticeType {
	e := NoticeType{ID: f.ID, Name: f.Name}
	for _, nf := range f.Notices {
		e.Notices = append(e.Notices, NoticeToEntity(nf))
	}
	return e
}

func NoticeFromEntity(e Notice) NoticeForm {
	var questionID *string
	if e.QuestionID != nil {
		s := strconv.FormatInt(*e.QuestionID, 10)
		questionID = &s
	}
	f := NoticeForm{
		ID:         e.ID,
		TypeID:     e.TypeID,
		QuestionID: questionID,
		Content:    e.Content,
		DisplayDue: timePtrToISO(e.DisplayDue),
	}
	if e.NoticeType.ID != 0 {
		nt := noticeTypeFromEntityShallow(e.NoticeType)
		f.NoticeType = &nt
	}
	if e.Question != nil && e.Question.ID != 0 {
		// due := e.Question.Due.Format(time.DateTime)
		// const tag = TagFromEntity(e.Question.TagManagers)
		// qf := QuestionForm{ID: e.Question.ID, Title: e.Question.Title, Content: e.Question.Content, Due: &due, Tags: tag}
		qf := QuestionFromEntity(*e.Question)
		f.Question = &qf
	}
	return f
}

func NoticeFromEntities(ns []Notice) (forms []NoticeForm) {
	for _, n := range ns {
		forms = append(forms, NoticeFromEntity(n))
	}
	return forms
}

func NoticeToEntity(f NoticeForm) Notice {
	var questionID *int64
	if v := f.QuestionIDInt64(); v >= 0 {
		questionID = &v
	}
	e := Notice{
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

func QuestionFromEntities(e []Question) []QuestionForm {
	forms := []QuestionForm{}
	for _, q := range e {
		forms = append(forms, QuestionFromEntity(q))
	}
	return forms
}

// --- Question ---

func QuestionFromEntity(e Question) QuestionForm {
	var originQuestionID *string
	if e.OriginQuestionID != nil {
		s := strconv.FormatInt(*e.OriginQuestionID, 10)
		originQuestionID = &s
	}
	f := QuestionForm{
		ID:               e.ID,
		OriginQuestionID: originQuestionID,
		AnswerID:         e.AnswerID,
		SupportID:        e.SupportID,
		Title:            e.Title,
		Content:          e.Content,
		Deleted:          e.Deleted,
		Due:              timePtrToISO(e.Due),
		CreatedAt:        timeToISO(e.CreatedAt),
	}
	if e.Answer != nil && e.Answer.ID != 0 {
		a := AnswerFromEntity(*e.Answer)
		f.Answer = &a
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
	return f
}

func QuestionToEntity(f QuestionForm) Question {
	var originQuestionID *int64
	if v := f.OriginQuestionIDInt64(); v >= 0 {
		originQuestionID = &v
	}
	e := Question{
		ID:               f.ID,
		OriginQuestionID: originQuestionID,
		AnswerID:         f.AnswerID,
		SupportID:        f.SupportID,
		Title:            f.Title,
		Content:          f.Content,
		Deleted:          f.Deleted,
		Due:              isoToTimePtr(f.Due),
	}
	if f.CreatedAt == nil || *f.CreatedAt == "" {
		e.CreatedAt = time.Now()
	} else {
		e.CreatedAt = isoToTime(f.CreatedAt)
	}
	qid := f.ID
	if f.Answer != nil {
		a := AnswerToEntity(*f.Answer)
		e.Answer = &a
		if e.AnswerID == nil && a.ID != 0 {
			aid := a.ID
			e.AnswerID = &aid
		}
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
		tm := TagManager{
			QuestionID: qid,
			TagID:      tf.ID,
			Tag:        Tag{ID: tf.ID},
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
	return e
}

func questionFormShallowFromEntity(e Question) QuestionForm {
	var originQuestionID *string
	if e.OriginQuestionID != nil {
		s := strconv.FormatInt(*e.OriginQuestionID, 10)
		originQuestionID = &s
	}
	return QuestionForm{
		ID:               e.ID,
		OriginQuestionID: originQuestionID,
		AnswerID:         e.AnswerID,
		SupportID:        e.SupportID,
		Title:            e.Title,
		Content:          e.Content,
		Deleted:          e.Deleted,
		Due:              timePtrToISO(e.Due),
		CreatedAt:        timeToISO(e.CreatedAt),
	}
}

func RelatedQuestionFromEntity(r RelatedQuestion) RelatedQuestionForm {
	f := RelatedQuestionForm{
		ID:                r.ID,
		QuestionID:        strconv.FormatInt(r.QuestionID, 10),
		RelatedQuestionID: strconv.FormatInt(r.RelatedQuestionID, 10),
	}
	if r.RelatedQuestion.ID != 0 {
		q := questionFormShallowFromEntity(r.RelatedQuestion)
		f.RelatedQuestion = &q
	}
	return f
}

func RelatedQuestionToEntity(f RelatedQuestionForm, parentQuestionID int64) RelatedQuestion {
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
	e := RelatedQuestion{
		ID:                f.ID,
		QuestionID:        qid,
		RelatedQuestionID: rid,
	}
	if f.RelatedQuestion != nil && f.RelatedQuestion.ID != 0 {
		e.RelatedQuestion = QuestionToEntity(*f.RelatedQuestion)
	}
	return e
}

// --- ReferManager / TagManager (optional full graph) ---

func ReferManagerFromEntity(e ReferManager) ReferManagerForm {
	f := ReferManagerForm{
		ID:       e.ID,
		AnswerID: strconv.FormatInt(e.AnswerID, 10),
		ReferID:  strconv.FormatInt(e.ReferID, 10),
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

func ReferManagerToEntity(f ReferManagerForm) ReferManager {
	e := ReferManager{
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

func TagManagerFromEntity(e TagManager) TagManagerForm {
	f := TagManagerForm{
		ID:         e.ID,
		TagID:      strconv.FormatInt(e.TagID, 10),
		QuestionID: strconv.FormatInt(e.QuestionID, 10),
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

func TagManagerToEntity(f TagManagerForm) TagManager {
	e := TagManager{
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
