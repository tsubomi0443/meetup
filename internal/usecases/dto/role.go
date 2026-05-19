// Package dto は API・画面との入出力で用いるフォーム DTO を定義する。
// エンティティとの相互変換は mapper パッケージが担う。
package dto

// =====================
// RoleForm はロールのフォーム DTO である。
type RoleForm struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	CreatedAt *string    `json:"createdAt,omitempty"`
	UpdatedAt *string    `json:"updatedAt,omitempty"`
	DeletedAt *string    `json:"deletedAt,omitempty"`
	Users     []UserForm `json:"users,omitempty"`
}
