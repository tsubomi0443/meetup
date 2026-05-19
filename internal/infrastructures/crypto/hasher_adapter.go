package crypto

import "meetup/internal/ports"

// HasherAdapter は ports.PasswordHasher を実装する。
type HasherAdapter struct{}

// NewHasherAdapter はパスワードハッシュアダプタを生成する。
//
// return:
//   - ports.PasswordHasher: ハッシュ・照合アダプタ
func NewHasherAdapter() ports.PasswordHasher {
	return HasherAdapter{}
}

// EncryptPasswordByArgon2Encode は平文パスワードを Argon2id でハッシュする。
//
// args:
//   - pass string: 平文パスワード
//
// return:
//   - string: PHC エンコード済みハッシュ
//   - error: ハッシュ化エラー
func (HasherAdapter) EncryptPasswordByArgon2Encode(pass string) (string, error) {
	return EncryptPasswordByArgon2Encode(pass)
}

// VerifyPassword は保存ハッシュと平文パスワードを照合する。
//
// args:
//   - encodedHash string: PHC 形式の保存ハッシュ
//   - password string: 平文パスワード
//
// return:
//   - bool: 一致すれば true
//   - error: 照合処理エラー
func (HasherAdapter) VerifyPassword(encodedHash, password string) (bool, error) {
	return VerifyPassword(encodedHash, password)
}
