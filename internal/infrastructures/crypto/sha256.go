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

// EncryptSHA256 はペッパーを挟んで target を SHA256 ハッシュする（保存パスワード用ではなくログ等向け）。
//
// args:
//   - target string: ハッシュ対象文字列
//
// return:
//   - string: 16進エンコードされたハッシュ
func EncryptSHA256(target string) string {
	newTarget := fmt.Sprintf("%[2]s-%[1]s-%[2]s", target, config.GetPepper())
	hash := sha256.Sum256([]byte(newTarget))
	return hex.EncodeToString(hash[:])
}

// EncryptPasswordBySHA256 は EncryptSHA256 を反復適用する（レガシー／非 Argon 経路向け）。
//
// args:
//   - pass string: 平文パスワード
//
// return:
//   - string: ストレッチ後のハッシュ文字列
func EncryptPasswordBySHA256(pass string) string {
	password := pass
	stretchSeed := []rune{rune(pass[0]), rune(pass[len(pass)/2]), rune(pass[len(pass)-1])}
	stretch := calcTotalRune(stretchSeed...)/10 + REQUIRE_MIN_LOOP
	for cnt := stretch; cnt > 0; cnt-- {
		password = EncryptSHA256(password)
	}
	return password
}

// calcTotalRune は与えられたルーンの Unicode コードポイント合計を返す。
//
// args:
//   - runes ...rune: 対象ルーン
//
// return:
//   - int: 合計値
func calcTotalRune(runes ...rune) int {
	var cnt int
	for _, r := range runes {
		cnt += int(r)
	}
	return cnt
}
