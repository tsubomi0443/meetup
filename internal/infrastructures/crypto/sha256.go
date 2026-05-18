package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"meetup/internal/infrastructures/config"
)

const (
	REQUIRE_MIN_LOOP = 4
)

// EncryptSHA256 hashes target with pepper (not for stored passwords; used for logging etc.).
func EncryptSHA256(target string) string {
	newTarget := fmt.Sprintf("%[2]s-%[1]s-%[2]s", target, config.GetPepper())
	hash := sha256.Sum256([]byte(newTarget))
	return hex.EncodeToString(hash[:])
}

// EncryptPasswordBySHA256 applies EncryptSHA256 iteratively for legacy / non-Argon paths.
func EncryptPasswordBySHA256(pass string) string {
	password := pass
	stretchSeed := []rune{rune(pass[0]), rune(pass[len(pass)/2]), rune(pass[len(pass)-1])}
	stretch := calcTotalRune(stretchSeed...)/10 + REQUIRE_MIN_LOOP
	for cnt := stretch; cnt > 0; cnt-- {
		password = EncryptSHA256(password)
	}
	return password
}

func calcTotalRune(runes ...rune) int {
	var cnt int
	for _, r := range runes {
		cnt += int(r)
	}
	return cnt
}
