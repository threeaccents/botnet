package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/rs/zerolog"
)

type CryptoService struct {
	Key string

	Log zerolog.Logger
}

func (s CryptoService) key() []byte {
	return []byte(s.Key)
}

func (s *CryptoService) Encrypt(plaintext []byte) ([]byte, error) {
	c, err := aes.NewCipher(s.key())
	if err != nil {
		return nil, fmt.Errorf("error creating new cipher %v", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, fmt.Errorf("error creating new GCM block %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("error creating nonce %v", err)
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func (s *CryptoService) EncryptToString(plaintext []byte) (string, error) {
	ciphertext, err := s.Encrypt(plaintext)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(ciphertext)

	return encoded, nil
}

func (s *CryptoService) Decrypt(ciphertext []byte) ([]byte, error) {
	c, err := aes.NewCipher(s.key())
	if err != nil {
		return nil, fmt.Errorf("error creating new cipher %v", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, fmt.Errorf("error creating new GCM block %v", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("error ciphertext too small")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("error opening ciphertext %v", err)
	}

	return plaintext, nil
}

func (s *CryptoService) DecryptString(ciphertext string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	return s.Decrypt(decoded)
}
