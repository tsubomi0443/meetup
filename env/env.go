package env

import (
	"fmt"
	"os"
)

// GetDSN builds a Postgres connection string from environment variables (see .env.example).
func GetDSN() string {
	host := getenv("PSQL_HOST", "localhost")
	port := getenv("PSQL_PORT", "5432")
	user := getenv("PSQL_USER", "postgres")
	password := os.Getenv("PSQL_PASSWORD")
	dbname := getenv("PSQL_DBNAME", "postgresdb")
	sslmode := getenv("PSQL_SSLMODE", "disable")
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
