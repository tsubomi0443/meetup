package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"meetup/internal/infrastructures/config"

	"golang.org/x/crypto/argon2"
)

const (
	argonTime    = 3
	argonMemory  = 128 * 1024 // 128 MiB
	argonThreads = 4
	argonKeyLen  = 32
)

func encryptArgon2ID(target string, salt []byte, time uint32, memory uint32, threads uint8, keyLen uint32) []byte {
	return argon2.IDKey(
		[]byte(target),
		salt,
		time,
		memory,
		threads,
		keyLen,
	)
}

func encodePHC(time, memory uint32, threads uint8, hash, salt []byte) string {
	return fmt.Sprintf(
		"$argon2id$v=19$t=%d,m=%d,p=%d$%s$%s",
		time,
		memory,
		threads,
		base64.StdEncoding.EncodeToString(salt),
		base64.StdEncoding.EncodeToString(hash),
	)
}

func generateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	return salt, err
}

// EncryptPasswordByArgon2Encode returns a PHC-style Argon2id string for storage in DB.
func EncryptPasswordByArgon2Encode(pass string) (string, error) {
	salt, err := generateSalt(16)
	if err != nil {
		return "", err
	}
	hash := encryptArgon2ID(pass+config.GetPepper(), salt, argonTime, argonMemory, argonThreads, argonKeyLen)
	return encodePHC(argonTime, argonMemory, argonThreads, hash, salt), nil
}

// VerifyPassword checks plaintext password against an Argon2id PHC-encoded hash.
func VerifyPassword(encodedHash, password string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("invalid hash format")
	}

	var time uint32
	var memory uint32
	var threads uint8

	if _, err := fmt.Sscanf(parts[3], "t=%d,m=%d,p=%d", &time, &memory, &threads); err != nil {
		return false, err
	}

	salt, err := base64.StdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	expectedHash, err := base64.StdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	computedHash := encryptArgon2ID(
		password+config.GetPepper(),
		salt,
		time,
		memory,
		threads,
		uint32(len(expectedHash)),
	)

	if subtle.ConstantTimeCompare(computedHash, expectedHash) == 1 {
		return true, nil
	}
	return false, nil
}
