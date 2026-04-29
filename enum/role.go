package enum

type RoleCode int

const (
	Admin RoleCode = iota + 1
	Manager
	Staff
	Employee
)

type CategoryCode int
type NoticeTypeCode int
type SupportStatusCode int
type TagCode int
