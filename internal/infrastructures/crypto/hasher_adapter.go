package crypto

import "meetup/internal/ports"

// HasherAdapter implements ports.PasswordHasher.
type HasherAdapter struct{}

func NewHasherAdapter() ports.PasswordHasher {
	return HasherAdapter{}
}

func (HasherAdapter) EncryptPasswordByArgon2Encode(pass string) (string, error) {
	return EncryptPasswordByArgon2Encode(pass)
}

func (HasherAdapter) VerifyPassword(encodedHash, password string) (bool, error) {
	return VerifyPassword(encodedHash, password)
}
