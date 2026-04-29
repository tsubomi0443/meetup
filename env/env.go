package env

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	jwtkey   = "JWK_KEY"
	modeIs   = "MODE"
	host     = "PSQL_HOST"
	port     = "PSQL_PORT"
	user     = "PSQL_USER"
	password = "PSQL_PASSWORD"
	dbname   = "PSQL_DBNAME"
	sslmode  = "PSQL_SSLMODE"
)

// GetDSN builds a Postgres connection string from environment variables (see .env.example).
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

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func GetJWTKey() string {
	_jwtKey := os.Getenv(jwtkey)
	if _jwtKey == "" {
		log.Fatalln(fmt.Errorf("configure to %s", jwtkey))
	}
	return _jwtKey
}

func IsDevelop() bool {
	switch strings.ToUpper(os.Getenv(modeIs)) {
	case "DEV", "DEVELOP":
		return true
	default:
		return false
	}
}

func IsProduct() bool {
	switch strings.ToUpper(os.Getenv(modeIs)) {
	case "PRO", "PROD", "PRODUCT":
		return true
	default:
		return false
	}
}
