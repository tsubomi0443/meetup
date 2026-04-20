package infrastructure

// Form types mirror static/js/model.js (camelCase JSON for the UI layer).

// RoleForm corresponds to Role.
type RoleForm struct {
	ID       int64      `json:"id"`
	RoleName string     `json:"roleName"`
	Users    []UserForm `json:"users,omitempty"`
}

// SupportStatusForm corresponds to SupportStatus.
type SupportStatusForm struct {
	ID       int64          `json:"id"`
	Title    string         `json:"title"`
	Supports []SupportForm  `json:"supports,omitempty"`
}

// SupportForm corresponds to Support.
type SupportForm struct {
	ID              int64              `json:"id"`
	UserID          int64              `json:"userId"`
	SupportStatusID int64              `json:"supportStatusId"`
	User            *UserForm          `json:"user,omitempty"`
	SupportStatus   *SupportStatusForm `json:"supportStatus,omitempty"`
	Questions       []QuestionForm     `json:"questions,omitempty"`
}

// UserForm corresponds to User.
type UserForm struct {
	ID        int64          `json:"id"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	RoleID    int64          `json:"roleId"`
	Role      *RoleForm      `json:"role,omitempty"`
	Supports  []SupportForm  `json:"supports,omitempty"`
	Answers   []AnswerForm   `json:"answers,omitempty"`
	Memos     []MemoForm     `json:"memos,omitempty"`
	Password  string         `json:"password,omitempty"`
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
	CategoryID int64          `json:"categoryId"`
	Category   *CategoryForm  `json:"category,omitempty"`
	Questions  []QuestionForm `json:"questions,omitempty"`
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
	ID         int64          `json:"id"`
	QuestionID int64          `json:"questionId"`
	UserID     int64          `json:"userId"`
	Content    string         `json:"content"`
	Question   *QuestionForm  `json:"question,omitempty"`
	User       *UserForm      `json:"user,omitempty"`
}

// AnswerForm corresponds to Answer.
type AnswerForm struct {
	ID         int64          `json:"id"`
	UserID     int64          `json:"userId"`
	QuestionID int64          `json:"questionId"`
	Content    string         `json:"content"`
	AnsweredAt *string        `json:"answeredAt,omitempty"`
	CreatedAt  *string        `json:"createdAt,omitempty"`
	User       *UserForm      `json:"user,omitempty"`
	Question   *QuestionForm  `json:"question,omitempty"`
	Refers     []ReferForm    `json:"refers,omitempty"`
}

// EscalationForm corresponds to Escalation.
type EscalationForm struct {
	ID             int64          `json:"id"`
	FromQuestionID int64          `json:"fromQuestionId"`
	ToQuestionID   int64          `json:"toQuestionId"`
	EscalatedAt    *string        `json:"escalatedAt,omitempty"`
	FromQuestion   *QuestionForm  `json:"fromQuestion,omitempty"`
	ToQuestion     *QuestionForm  `json:"toQuestion,omitempty"`
}

// QuestionForm corresponds to Question.
type QuestionForm struct {
	ID               int64            `json:"id"`
	OriginQuestionID *int64           `json:"originQuestionId,omitempty"`
	SupportID        *int64           `json:"supportId,omitempty"`
	Title            string           `json:"title"`
	Content          string           `json:"content"`
	Due              *string          `json:"due,omitempty"`
	CreatedAt        *string          `json:"createdAt,omitempty"`
	OriginQuestion   *QuestionForm    `json:"originQuestion,omitempty"`
	SubQuestions     []QuestionForm   `json:"subQuestions,omitempty"`
	Support          *SupportForm      `json:"support,omitempty"`
	Answers          []AnswerForm      `json:"answers,omitempty"`
	Memos            []MemoForm        `json:"memos,omitempty"`
	Tags             []TagForm         `json:"tags,omitempty"`
	EscalationsFrom  []EscalationForm  `json:"escalationsFrom,omitempty"`
	EscalationsTo    []EscalationForm  `json:"escalationsTo,omitempty"`
}

// ReferManagerForm corresponds to ReferManager.
type ReferManagerForm struct {
	ID       int64       `json:"id"`
	AnswerID int64       `json:"answerId"`
	ReferID  int64       `json:"referId"`
	Answer   *AnswerForm `json:"answer,omitempty"`
	Refer    *ReferForm  `json:"refer,omitempty"`
}

// TagManagerForm corresponds to TagManager.
type TagManagerForm struct {
	ID         int64          `json:"id"`
	TagID      int64          `json:"tagId"`
	QuestionID int64          `json:"questionId"`
	Tag        *TagForm       `json:"tag,omitempty"`
	Question   *QuestionForm  `json:"question,omitempty"`
}
