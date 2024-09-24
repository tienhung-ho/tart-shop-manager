package cacheutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"
	"os"
)

// Hàm mã hóa
func Encrypt(data []byte) ([]byte, error) {

	key := os.Getenv("ENCRYPTION_KEY")
	if len(key) != 32 { // AES-256 yêu cầu khóa 32 byte
		log.Print(len(key))
		//panic("ENCRYPTION_KEY must be 32 bytes long")
	}
	encryptionKey := []byte(key)

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := aesGCM.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}
