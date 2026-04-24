package infrastructure

import (
	"time"
)

// =====================
// ROLE
// =====================
type Role struct {
	ID       int64  `gorm:"column:id;primaryKey"`
	RoleName string `gorm:"column:role_name"`
	Users    []User `gorm:"foreignKey:RoleID;references:ID"`
}

func (Role) TableName() string {
	return "roles"
}

// =====================
// USER
// =====================
type User struct {
	ID       int64  `gorm:"column:id;primaryKey"`
	Name     string `gorm:"column:name"`
	Password string `gorm:"column:password"`
	Email    string `gorm:"column:email"`
	RoleID   int64  `gorm:"column:role_id"`
	Role     Role   `gorm:"foreignKey:RoleID;references:ID"`
}

func (User) TableName() string {
	return "users"
}

// =====================
// CATEGORY
// =====================
type Category struct {
	ID           int64  `gorm:"column:id;primaryKey"`
	CategoryName string `gorm:"column:category_name"`
	Tags         []Tag  `gorm:"foreignKey:CategoryID;references:ID"`
}

func (Category) TableName() string {
	return "categories"
}

// =====================
// TAG
// =====================
type Tag struct {
	ID          int64        `gorm:"column:id;primaryKey"`
	Title       string       `gorm:"column:title"`
	Usage       int          `gorm:"column:usage"`
	CategoryID  int64        `gorm:"column:category_id"`
	Category    Category     `gorm:"foreignKey:CategoryID;references:ID"`
	TagManagers []TagManager `gorm:"foreignKey:TagID;references:ID"`
}

func (Tag) TableName() string {
	return "tags"
}

// =====================
// QUESTION
// =====================
type Question struct {
	ID               int64        `gorm:"column:id;primaryKey"`
	OriginQuestionID *int64       `gorm:"column:origin_question_id"`
	AnswerID         *int64       `gorm:"column:answer_id"`
	SupportID        *int64       `gorm:"column:support_id"`
	Title            string       `gorm:"column:title"`
	Content          string       `gorm:"column:content"`
	Deleted          bool         `gorm:"column:deleted"`
	Due              *time.Time   `gorm:"column:due"`
	CreatedAt        time.Time    `gorm:"column:created_at"`
	Answer           *Answer      `gorm:"foreignKey:AnswerID;references:ID"`
	Memos            []Memo       `gorm:"foreignKey:QuestionID;references:ID"`
	TagManagers      []TagManager `gorm:"foreignKey:QuestionID;references:ID"`
	Support          *Support     `gorm:"foreignKey:SupportID;references:ID"`
}

func (Question) TableName() string {
	return "questions"
}

// =====================
// ANSWER
// =====================
type Answer struct {
	ID            int64          `gorm:"column:id;primaryKey"`
	UserID        int64          `gorm:"column:user_id"`
	Content       string         `gorm:"column:content"`
	AnsweredAt    *time.Time     `gorm:"column:answered_at"`
	CreatedAt     time.Time      `gorm:"column:created_at"`
	User          User           `gorm:"foreignKey:UserID;references:ID"`
	ReferManagers []ReferManager `gorm:"foreignKey:AnswerID;references:ID"`
}

func (Answer) TableName() string {
	return "answers"
}

// =====================
// MEMO
// =====================
type Memo struct {
	ID         int64    `gorm:"column:id;primaryKey"`
	UserID     int64    `gorm:"column:user_id"`
	Content    string   `gorm:"column:content"`
	QuestionID int64    `gorm:"column:question_id"`
	User       User     `gorm:"foreignKey:UserID;references:ID"`
	Question   Question `gorm:"foreignKey:QuestionID;references:ID"`
}

func (Memo) TableName() string {
	return "memos"
}

// =====================
// REFER
// =====================
type Refer struct {
	ID            int64          `gorm:"column:id;primaryKey"`
	Title         string         `gorm:"column:title"`
	URL           string         `gorm:"column:url"`
	ReferManagers []ReferManager `gorm:"foreignKey:ReferID;references:ID"`
}

func (Refer) TableName() string {
	return "refers"
}

// =====================
// REFER_MANAGER
// =====================
type ReferManager struct {
	ID       int64  `gorm:"column:id;primaryKey"`
	AnswerID int64  `gorm:"column:answer_id"`
	ReferID  int64  `gorm:"column:refer_id"`
	Answer   Answer `gorm:"foreignKey:AnswerID;references:ID"`
	Refer    Refer  `gorm:"foreignKey:ReferID;references:ID"`
}

func (ReferManager) TableName() string {
	return "refer_managers"
}

// =====================
// TAG_MANAGER
// =====================
type TagManager struct {
	ID         int64    `gorm:"column:id;primaryKey"`
	TagID      int64    `gorm:"column:tag_id"`
	QuestionID int64    `gorm:"column:question_id"`
	Tag        Tag      `gorm:"foreignKey:TagID;references:ID"`
	Question   Question `gorm:"foreignKey:QuestionID;references:ID"`
}

func (TagManager) TableName() string {
	return "tag_managers"
}

// =====================
// ESCALATION
// =====================
type Escalation struct {
	ID             int64     `gorm:"column:id;primaryKey"`
	FromQuestionID int64     `gorm:"column:from_question_id"`
	ToQuestionID   int64     `gorm:"column:to_question_id"`
	EscalatedAt    time.Time `gorm:"column:escalated_at"`
	FromQuestion   Question  `gorm:"foreignKey:FromQuestionID;references:ID"`
	ToQuestion     Question  `gorm:"foreignKey:ToQuestionID;references:ID"`
}

func (Escalation) TableName() string {
	return "escalations"
}

// =====================
// SUPPORT_STATUS
// =====================
type SupportStatus struct {
	ID       int64     `gorm:"column:id;primaryKey"`
	Title    string    `gorm:"column:title"`
	Supports []Support `gorm:"foreignKey:SupportStatusID;references:ID"`
}

func (SupportStatus) TableName() string {
	return "support_statuses"
}

// =====================
// SUPPORT
// =====================
type Support struct {
	ID              int64         `gorm:"column:id;primaryKey"`
	UserID          int64         `gorm:"column:user_id"`
	SupportStatusID int64         `gorm:"column:support_status_id"`
	User            User          `gorm:"foreignKey:UserID;references:ID"`
	SupportStatus   SupportStatus `gorm:"foreignKey:SupportStatusID;references:ID"`
}

func (Support) TableName() string {
	return "supports"
}

// =====================
// NOTICE_TYPE
// =====================
type NoticeType struct {
	ID      int64    `gorm:"column:id;primaryKey"`
	Name    string   `gorm:"column:name"`
	Notices []Notice `gorm:"foreignKey:TypeID;references:ID"`
}

func (NoticeType) TableName() string {
	return "notice_types"
}

// =====================
// NOTICE
// =====================
type Notice struct {
	ID         int64      `gorm:"column:id;primaryKey"`
	TypeID     int64      `gorm:"column:type_id"`
	QuestionID *int64     `gorm:"column:question_id"`
	Content    *string    `gorm:"column:content"`
	DisplayDue *time.Time `gorm:"column:display_due"`
	NoticeType NoticeType `gorm:"foreignKey:TypeID;references:ID"`
	Question   *Question  `gorm:"foreignKey:QuestionID;references:ID"`
}

func (Notice) TableName() string {
	return "notices"
}
