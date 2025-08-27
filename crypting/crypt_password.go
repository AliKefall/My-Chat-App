package crypting

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func HashPassword(password, salt string) string {
	hash := sha256.Sum256([]byte(password + salt))
	return fmt.Sprintf("%x", hash)
}

func GenerateSalt() (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(salt), nil
}
