package cacheutil

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"log"
	"os"
)

// Hàm giải mã
func Decrypt(data []byte) ([]byte, error) {

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
	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
