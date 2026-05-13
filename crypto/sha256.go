package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"meetup/env"
)

const REQUIRE_MIN_LOOP = 3

func EncryptSHA256(target string) string {
	newTarget := fmt.Sprintf("%s-%s", target, env.GetJWTKey())
	hash := sha256.Sum256([]byte(newTarget))
	return hex.EncodeToString(hash[:])
}

// EncryptPassword applies EncryptSHA256 iteratively for stored credentials.
func EncryptPassword(pass string) string {
	password := pass
	for cnt := len(pass)/2 + REQUIRE_MIN_LOOP; cnt > 0; cnt-- {
		password = EncryptSHA256(password)
	}
	return password
}
