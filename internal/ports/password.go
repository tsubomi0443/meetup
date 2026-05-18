package ports

// PasswordHasher hashes and verifies passwords (Argon2id).
type PasswordHasher interface {
	EncryptPasswordByArgon2Encode(pass string) (string, error)
	VerifyPassword(encodedHash, password string) (bool, error)
}
