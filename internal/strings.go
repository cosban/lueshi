package internal

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"strings"
)

const TOKEN_LENGTH = 64

// Hash takes a salt and provided string and returns their corresponding
// combined sha256 string
func Hash(salt, provided string) string {
	hasher := sha256.New()
	hasher.Write([]byte(provided))
	first := hex.EncodeToString(hasher.Sum(nil))

	hasher = sha256.New()
	hasher.Write(append([]byte(first), salt...))

	return hex.EncodeToString(hasher.Sum(nil))
}

// RandomString returns a random string with the specified length
func RandomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)[:length]
}

func URLToken() string {
	tok := RandomString(TOKEN_LENGTH)
	return strings.Replace(tok, "+", "", -1)
}
