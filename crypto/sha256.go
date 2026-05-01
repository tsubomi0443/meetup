package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"meetup/env"
)

func EncryptSHA256(target string) string {
	newTarget := fmt.Sprintf("%s-%s", target, env.GetJWTKey())
	hash := sha256.Sum256([]byte(newTarget))
	return hex.EncodeToString(hash[:])
}
