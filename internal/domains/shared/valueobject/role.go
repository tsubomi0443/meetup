// Package valueobject はドメイン横断で使う値オブジェクト（コード定数）を定義する。
package valueobject

// RoleCode はロールを表す数値コード。
type RoleCode int

const (
	// Admin は管理者ロール。
	Admin RoleCode = iota + 1
	// Manager はマネージャーロール。
	Manager
	// Staff はスタッフロール。
	Staff
	// Employee は一般社員ロール。
	Employee
)

// CategoryCode はカテゴリを表す数値コード。
type CategoryCode int

// NoticeTypeCode は通知種別を表す数値コード。
type NoticeTypeCode int

// SupportStatusCode はサポートステータスを表す数値コード。
type SupportStatusCode int

// TagCode はタグを表す数値コード。
type TagCode int
