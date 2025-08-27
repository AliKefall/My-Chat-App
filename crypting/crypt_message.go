package crypting

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
)

func GenerateIV() (string, error) {
	iv := make([]byte, 16)
	_, err := rand.Read(iv)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(iv), nil
}

func EncryptMessage(content, key, ivHex string) (string, error) {
	iv, _ := hex.DecodeString(ivHex)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	plaintext := []byte(content)
	cfb := cipher.NewCFBEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)

	return hex.EncodeToString(ciphertext), nil

}
