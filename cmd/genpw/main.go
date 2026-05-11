// genpw prints EncryptPassword hashes for given plaintexts (requires JWK_KEY).
package main

import (
	"fmt"
	"os"

	"meetup/crypto"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: JWK_KEY=... go run ./cmd/genpw <password>...")
		os.Exit(1)
	}
	for _, p := range os.Args[1:] {
		fmt.Println(p, "->", crypto.EncryptPassword(p))
	}
}
