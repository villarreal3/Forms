package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run generate_password_hash.go <password>")
		os.Exit(1)
	}

	password := os.Args[1]
	
	// Generar hash SHA256
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hash := hex.EncodeToString(hasher.Sum(nil))

	fmt.Printf("Password: %s\n", password)
	fmt.Printf("Hash SHA256: %s\n", hash)
}

