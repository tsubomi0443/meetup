package config

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	jwtkey   = "JWK_KEY"
	pepper   = "PEPPER"
	modeIs   = "MODE"
	host     = "PSQL_HOST"
	port     = "PSQL_PORT"
	user     = "PSQL_USER"
	password = "PSQL_PASSWORD"
	dbname   = "PSQL_DBNAME"
	sslmode  = "PSQL_SSLMODE"
)

// GetDSN は環境変数から PostgreSQL 接続文字列を組み立てる（.env.example 参照）。
//
// return:
//   - string: lib/pq 形式の DSN
func GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getenv(host, "localhost"),
		getenv(port, "5432"),
		getenv(user, "postgres"),
		os.Getenv(password),
		getenv(dbname, "postgresdb"),
		getenv(sslmode, "disable"),
	)
}

// getenv は環境変数を取得し、未設定時は既定値を返す。
//
// args:
//   - key string: 環境変数名
//   - def string: 既定値
//
// return:
//   - string: 環境変数の値または既定値
func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// GetJWTKey は JWT 署名用の秘密鍵を返す。未設定時はプロセスを終了する。
//
// return:
//   - string: JWT 署名鍵
func GetJWTKey() string {
	_jwtKey := os.Getenv(jwtkey)
	if _jwtKey == "" {
		log.Fatalln(fmt.Errorf("configure to %s", jwtkey))
	}
	return _jwtKey
}

// IsDevelop は開発モード（MODE=DEV/DEVELOP）かどうかを返す。
//
// return:
//   - bool: 開発モードなら true
func IsDevelop() bool {
	switch strings.ToUpper(os.Getenv(modeIs)) {
	case "DEV", "DEVELOP":
		return true
	default:
		return false
	}
}

// IsProduct は本番モード（MODE=PRO/PROD/PRODUCT）かどうかを返す。
//
// return:
//   - bool: 本番モードなら true
func IsProduct() bool {
	switch strings.ToUpper(os.Getenv(modeIs)) {
	case "PRO", "PROD", "PRODUCT":
		return true
	default:
		return false
	}
}

// GetPepper はパスワードハッシュ（Argon2 / SHA 補助）用のアプリケーションペッパーを返す。未設定時はプロセスを終了する。
//
// return:
//   - string: ペッパー文字列
func GetPepper() string {
	_pepper := os.Getenv(pepper)
	if _pepper == "" {
		log.Fatalln(fmt.Errorf("configure to %s", pepper))
	}
	return _pepper
}
