package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"os"
)

func Encrypt(plaintext []byte) []byte {
	key := readKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal(err)
	}
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	return gcm.Seal(nonce, nonce, plaintext, nil)
}

func Decrypt(encryptedText []byte) ([]byte, error) {
	key := readKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(encryptedText) < nonceSize {
		return nil, fmt.Errorf("encryptedText too short")
	}
	nonce, encryptedText := encryptedText[:nonceSize], encryptedText[nonceSize:]
	return gcm.Open(nil, nonce, encryptedText, nil)
}

func readKey() []byte {
	key, err := os.ReadFile("/keys/shared.key")
	if err != nil {
		log.Fatal("Could not read key:", err)
	}
	return key
}
