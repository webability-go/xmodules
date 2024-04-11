package tools

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
)

// hashSecretFunc utiliza SHA-256 para hash el secreto y generar una clave.
func hashSecret(secret string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(secret))
	return hasher.Sum(nil)
}

// encryptFunc encripta el texto con la clave derivada del secreto.
func Encrypt(plaintext string, secret string) (string, error) {
	key := hashSecret(secret)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return hex.EncodeToString(ciphertext), nil
}

// decryptFunc desencripta el texto con la clave derivada del secreto.
func Decrypt(encryptedText string, secret string) (string, error) {
	key := hashSecret(secret)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	data, err := hex.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	if len(data) > nonceSize {
		nonce, ciphertext := data[:nonceSize], data[nonceSize:]
		plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			return "", err
		}
		return string(plaintext), nil
	}
	return "", nil
}
