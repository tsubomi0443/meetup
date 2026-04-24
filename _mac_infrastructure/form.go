package infrastructure

import "strconv"

// Form types mirror static/js/model.js (camelCase JSON for the UI layer).

// RoleForm corresponds to Role.
type RoleForm struct {
	ID       int64      `json:"id"`
	RoleName string     `json:"roleName"`
	Users    []UserForm `json:"users,omitempty"`
}

// SupportStatusForm corresponds to SupportStatus.
type SupportStatusForm struct {
	ID       int64         `json:"id"`
	Title    string        `json:"title"`
	Supports []SupportForm `json:"supports,omitempty"`
}

// SupportForm corresponds to Support.
type SupportForm struct {
	ID              int64              `json:"id"`
	UserID          string             `json:"userId"`
	SupportStatusID string             `json:"supportStatusId"`
	User            *UserForm          `json:"user,omitempty"`
	SupportStatus   *SupportStatusForm `json:"supportStatus,omitempty"`
	Question        *QuestionForm      `json:"question,omitempty"`
}

func (f SupportForm) UserIDInt64() int64 {
	if val, err := strconv.ParseInt(f.UserID, 10, 64); err == nil {
		return val
	}
	return -1
}

func (f SupportForm) SupportStatusIDInt64() int64 {
	if val, err := strconv.ParseInt(f.SupportStatusID, 10, 64); err == nil {
		return val
	}
	return -1
}

// UserForm corresponds to User.
type UserForm struct {
	ID       int64         `json:"id"`
	Name     string        `json:"name"`
	Email    string        `json:"email"`
	RoleID   string        `json:"roleId"`
	Role     *RoleForm     `json:"role,omitempty"`
	Supports []SupportForm `json:"supports,omitempty"`
	Answers  []AnswerForm  `json:"answers,omitempty"`
	Memos    []MemoForm    `json:"memos,omitempty"`
	Password string        `json:"password,omitempty"`
}

func (uf UserForm) RoleIDInt64() int64 {
	if val, err := strconv.ParseInt(uf.RoleID, 10, 64); err == nil {
		return val
	}
	return -1
}

// CategoryForm corresponds to Category.
type CategoryForm struct {
	ID           int64     `json:"id"`
	CategoryName string    `json:"categoryName"`
	Tags         []TagForm `json:"tags,omitempty"`
}

// TagForm corresponds to Tag.
type TagForm struct {
	ID         int64          `json:"id"`
	Title      string         `json:"title"`
	Usage      int            `json:"usage"`
	CategoryID string         `json:"categoryId"`
	Category   *CategoryForm  `json:"category,omitempty"`
	Questions  []QuestionForm `json:"questions,omitempty"`
}

func (tf TagForm) CategoryIDInt64() int64 {
	if val, err := strconv.ParseInt(tf.CategoryID, 10, 64); err == nil {
		return val
	}
	return -1
}

// ReferForm corresponds to Refer.
type ReferForm struct {
	ID      int64        `json:"id"`
	Title   string       `json:"title"`
	URL     string       `json:"url"`
	Answers []AnswerForm `json:"answers,omitempty"`
}

// MemoForm corresponds to Memo.
type MemoForm struct {
	ID         int64         `json:"id"`
	QuestionID string        `json:"questionId"`
	UserID     string        `json:"userId"`
	Content    string        `json:"content"`
	Question   *QuestionForm `json:"question,omitempty"`
	User       *UserForm     `json:"user,omitempty"`
}

func (f MemoForm) QuestionIDInt64() int64 {
	if val, err := strconv.ParseInt(f.QuestionID, 10, 64); err == nil {
		return val
	}
	return -1
}

func (f MemoForm) UserIDInt64() int64 {
	if val, err := strconv.ParseInt(f.UserID, 10, 64); err == nil {
		return val
	}
	return -1
}

// AnswerForm corresponds to Answer.
type AnswerForm struct {
	ID         int64       `json:"id"`
	UserID     string      `json:"userId"`
	Content    string      `json:"content"`
	AnsweredAt *string     `json:"answeredAt,omitempty"`
	CreatedAt  *string     `json:"createdAt,omitempty"`
	User       *UserForm   `json:"user,omitempty"`
	Refers     []ReferForm `json:"refers,omitempty"`
}

func (f AnswerForm) UserIDInt64() int64 {
	if val, err := strconv.ParseInt(f.UserID, 10, 64); err == nil {
		return val
	}
	return -1
}

// EscalationForm corresponds to Escalation.
type EscalationForm struct {
	ID             int64         `json:"id"`
	FromQuestionID string        `json:"fromQuestionId"`
	ToQuestionID   string        `json:"toQuestionId"`
	EscalatedAt    *string       `json:"escalatedAt,omitempty"`
	FromQuestion   *QuestionForm `json:"fromQuestion,omitempty"`
	ToQuestion     *QuestionForm `json:"toQuestion,omitempty"`
}

func (f EscalationForm) FromQuestionIDInt64() int64 {
	if val, err := strconv.ParseInt(f.FromQuestionID, 10, 64); err == nil {
		return val
	}
	return -1
}

func (f EscalationForm) ToQuestionIDInt64() int64 {
	if val, err := strconv.ParseInt(f.ToQuestionID, 10, 64); err == nil {
		return val
	}
	return -1
}

// NoticeTypeForm corresponds to NoticeType.
type NoticeTypeForm struct {
	ID      int64        `json:"id"`
	Name    string       `json:"name"`
	Notices []NoticeForm `json:"notices,omitempty"`
}

// NoticeForm corresponds to Notice.
type NoticeForm struct {
	ID         int64           `json:"id"`
	TypeID     int64           `json:"typeId"`
	QuestionID *string         `json:"questionId,omitempty"`
	Content    *string         `json:"content,omitempty"`
	DisplayDue *string         `json:"displayDue,omitempty"`
	NoticeType *NoticeTypeForm `json:"noticeType,omitempty"`
	Question   *QuestionForm   `json:"question,omitempty"`
}

func (f NoticeForm) QuestionIDInt64() int64 {
	if f.QuestionID == nil || *f.QuestionID == "" {
		return -1
	}
	if val, err := strconv.ParseInt(*f.QuestionID, 10, 64); err == nil {
		return val
	}
	return -1
}

// QuestionForm corresponds to Question.
type QuestionForm struct {
	ID               int64            `json:"id"`
	OriginQuestionID *string          `json:"originQuestionId,omitempty"`
	AnswerID         *int64           `json:"answerId,omitempty"`
	SupportID        *int64           `json:"supportId,omitempty"`
	Title            string           `json:"title"`
	Content          string           `json:"content"`
	Deleted          bool             `json:"deleted"`
	Due              *string          `json:"due,omitempty"`
	CreatedAt        *string          `json:"createdAt,omitempty"`
	OriginQuestion   *QuestionForm    `json:"originQuestion,omitempty"`
	SubQuestions     []QuestionForm   `json:"subQuestions,omitempty"`
	Support          *SupportForm     `json:"support,omitempty"`
	Answer           *AnswerForm      `json:"answer,omitempty"`
	Memos            []MemoForm       `json:"memos,omitempty"`
	Tags             []TagForm        `json:"tags,omitempty"`
	EscalationsFrom  []EscalationForm `json:"escalationsFrom,omitempty"`
	EscalationsTo    []EscalationForm `json:"escalationsTo,omitempty"`
}

func (f QuestionForm) OriginQuestionIDInt64() int64 {
	if f.OriginQuestionID == nil || *f.OriginQuestionID == "" {
		return -1
	}
	if val, err := strconv.ParseInt(*f.OriginQuestionID, 10, 64); err == nil {
		return val
	}
	return -1
}

// ReferManagerForm corresponds to ReferManager.
type ReferManagerForm struct {
	ID       int64       `json:"id"`
	AnswerID string      `json:"answerId"`
	ReferID  string      `json:"referId"`
	Answer   *AnswerForm `json:"answer,omitempty"`
	Refer    *ReferForm  `json:"refer,omitempty"`
}

func (f ReferManagerForm) AnswerIDInt64() int64 {
	if val, err := strconv.ParseInt(f.AnswerID, 10, 64); err == nil {
		return val
	}
	return -1
}

func (f ReferManagerForm) ReferIDInt64() int64 {
	if val, err := strconv.ParseInt(f.ReferID, 10, 64); err == nil {
		return val
	}
	return -1
}

// TagManagerForm corresponds to TagManager.
type TagManagerForm struct {
	ID         int64         `json:"id"`
	TagID      string        `json:"tagId"`
	QuestionID string        `json:"questionId"`
	Tag        *TagForm      `json:"tag,omitempty"`
	Question   *QuestionForm `json:"question,omitempty"`
}

func (f TagManagerForm) TagIDInt64() int64 {
	if val, err := strconv.ParseInt(f.TagID, 10, 64); err == nil {
		return val
	}
	return -1
}

func (f TagManagerForm) QuestionIDInt64() int64 {
	if val, err := strconv.ParseInt(f.QuestionID, 10, 64); err == nil {
		return val
	}
	return -1
}
