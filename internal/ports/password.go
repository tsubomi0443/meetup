package ports

// PasswordHasher はパスワードのハッシュ化と照合（Argon2id）を行う。
type PasswordHasher interface {
	EncryptPasswordByArgon2Encode(pass string) (string, error)
	VerifyPassword(encodedHash, password string) (bool, error)
}
