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

// encryptArgon2ID は Argon2id で鍵導出を行う。
//
// args:
//   - target string: ハッシュ対象（ペッパー付き平文等）
//   - salt []byte: ソルト
//   - time uint32: 反復回数
//   - memory uint32: メモリ（KiB）
//   - threads uint8: 並列度
//   - keyLen uint32: 出力鍵長
//
// return:
//   - []byte: 導出したハッシュ
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

// encodePHC は Argon2id ハッシュを PHC 文字列形式にエンコードする。
//
// args:
//   - time uint32: 反復回数
//   - memory uint32: メモリ（KiB）
//   - threads uint8: 並列度
//   - hash []byte: ハッシュ本体
//   - salt []byte: ソルト
//
// return:
//   - string: PHC 形式のエンコード文字列
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

// generateSalt は暗号論的乱数でソルトを生成する。
//
// args:
//   - length int: バイト長
//
// return:
//   - []byte: ソルト
//   - error: 乱数生成エラー
func generateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	return salt, err
}

// EncryptPasswordByArgon2Encode は DB 保存用の PHC 形式 Argon2id 文字列を返す。
//
// args:
//   - pass string: 平文パスワード
//
// return:
//   - string: PHC エンコード済みハッシュ
//   - error: ソルト生成・ハッシュ化エラー
func EncryptPasswordByArgon2Encode(pass string) (string, error) {
	salt, err := generateSalt(16)
	if err != nil {
		return "", err
	}
	hash := encryptArgon2ID(pass+config.GetPepper(), salt, argonTime, argonMemory, argonThreads, argonKeyLen)
	return encodePHC(argonTime, argonMemory, argonThreads, hash, salt), nil
}

// VerifyPassword は平文パスワードと Argon2id PHC エンコード済みハッシュを照合する。
//
// args:
//   - encodedHash string: PHC 形式の保存ハッシュ
//   - password string: 照合する平文パスワード
//
// return:
//   - bool: 一致すれば true
//   - error: 形式不正・デコード・ハッシュ計算エラー
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
