package mapper

import (
	"time"

	"gorm.io/gorm"
)

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
