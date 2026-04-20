package infrastructure

import "time"

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
		UserID:          e.UserID,
		SupportStatusID: e.SupportStatusID,
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
		UserID:          f.UserID,
		SupportStatusID: f.SupportStatusID,
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
		RoleID: e.RoleID,
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
		RoleID: e.RoleID,
	}
}

func UserToEntity(f UserForm) User {
	e := User{
		ID:     f.ID,
		Name:   f.Name,
		Email:  f.Email,
		RoleID: f.RoleID,
	}
	if f.RoleID == 0 && f.Role != nil {
		e.RoleID = f.Role.ID
	}
	if f.Password != "" {
		e.Passwordd = f.Password
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
		CategoryID: e.CategoryID,
	}
	if e.Category.ID != 0 {
		c := categoryFromEntityShallow(e.Category)
		f.Category = &c
	}
	return f
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
		CategoryID: f.CategoryID,
	}
	if f.Category != nil {
		e.Category = CategoryToEntity(*f.Category)
	}
	tagID := f.ID
	for _, qf := range f.Questions {
		tm := TagManager{
			TagID:      tagID,
			QuestionID: qf.ID,
		}
		if qf.ID != 0 {
			tm.Question = Question{ID: qf.ID}
		}
		e.TagManagers = append(e.TagManagers, tm)
	}
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
		QuestionID: e.QuestionID,
		UserID:     e.UserID,
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
		QuestionID: f.QuestionID,
		UserID:     f.UserID,
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
		UserID:     e.UserID,
		QuestionID: e.QuestionID,
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
		UserID:     f.UserID,
		QuestionID: f.QuestionID,
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
		FromQuestionID: e.FromQuestionID,
		ToQuestionID:   e.ToQuestionID,
		EscalatedAt:    timeToISO(e.EscalatedAt),
	}
}

func EscalationToEntity(f EscalationForm) Escalation {
	e := Escalation{
		ID:             f.ID,
		FromQuestionID: f.FromQuestionID,
		ToQuestionID:   f.ToQuestionID,
	}
	if f.EscalatedAt == nil || *f.EscalatedAt == "" {
		e.EscalatedAt = time.Now()
	} else {
		e.EscalatedAt = isoToTime(f.EscalatedAt)
	}
	return e
}

// --- Question ---

func QuestionFromEntity(e Question) QuestionForm {
	f := QuestionForm{
		ID:               e.ID,
		OriginQuestionID: e.OriginQuestionID,
		SupportID:        e.SupportID,
		Title:            e.Title,
		Content:          e.Content,
		Due:              timePtrToISO(e.Due),
		CreatedAt:        timeToISO(e.CreatedAt),
	}
	for _, a := range e.Answers {
		f.Answers = append(f.Answers, AnswerFromEntity(a))
	}
	for _, m := range e.Memos {
		f.Memos = append(f.Memos, MemoFromEntity(m))
	}
	for _, tm := range e.TagManagers {
		if tm.Tag.ID != 0 {
			f.Tags = append(f.Tags, tagFromEntityShallow(tm.Tag))
		}
	}
	if e.Support.ID != 0 {
		s := SupportFromEntity(e.Support)
		f.Support = &s
	}
	return f
}

func QuestionToEntity(f QuestionForm) Question {
	e := Question{
		ID:               f.ID,
		OriginQuestionID: f.OriginQuestionID,
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
		e.Answers = append(e.Answers, a)
	}
	for _, mf := range f.Memos {
		m := MemoToEntity(mf)
		if m.QuestionID == 0 {
			m.QuestionID = qid
		}
		e.Memos = append(e.Memos, m)
	}
	for _, tf := range f.Tags {
		tm := TagManager{
			QuestionID: qid,
			TagID:      tf.ID,
		}
		if tf.ID != 0 {
			tm.Tag = Tag{ID: tf.ID}
		}
		e.TagManagers = append(e.TagManagers, tm)
	}
	if f.Support != nil {
		e.Support = SupportToEntity(*f.Support)
	}
	return e
}

// --- ReferManager / TagManager (optional full graph) ---

func ReferManagerFromEntity(e ReferManager) ReferManagerForm {
	f := ReferManagerForm{
		ID:       e.ID,
		AnswerID: e.AnswerID,
		ReferID:  e.ReferID,
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
		AnswerID: f.AnswerID,
		ReferID:  f.ReferID,
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
		TagID:      e.TagID,
		QuestionID: e.QuestionID,
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
		TagID:      f.TagID,
		QuestionID: f.QuestionID,
	}
	if f.Tag != nil {
		e.Tag = TagToEntity(*f.Tag)
	}
	if f.Question != nil {
		e.Question = QuestionToEntity(*f.Question)
	}
	return e
}
