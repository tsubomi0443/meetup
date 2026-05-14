// genpw prints Argon2id-encoded password strings (requires PEPPER in env; loads .env if present).
package main

import (
	"fmt"
	"os"

	"meetup/crypto"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: PEPPER=... go run ./cmd/genpw <password>...")
		fmt.Fprintln(os.Stderr, "       go run ./cmd/app   # アプリ起動")
		os.Exit(1)
	}
	for _, p := range os.Args[1:] {
		h, err := crypto.EncryptPasswordByArgon2Encode(p)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(p, "->", h)
	}
}
