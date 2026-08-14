package copymanga

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"regexp"
)

func decryption(cct string, contentKey string) ([]byte, error) {
	if len(contentKey) < aes.BlockSize {
		return nil, fmt.Errorf("contentKey must contain a %d-byte IV", aes.BlockSize)
	}

	block, err := aes.NewCipher([]byte(cct))
	if err != nil {
		return nil, fmt.Errorf("create AES cipher: %w", err)
	}

	ciphertext, err := hex.DecodeString(contentKey[aes.BlockSize:])
	if err != nil {
		return nil, fmt.Errorf("decode contentKey ciphertext as hex: %w", err)
	}
	if len(ciphertext) == 0 || len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("ciphertext length must be a non-zero multiple of %d", aes.BlockSize)
	}

	plaintext := make([]byte, len(ciphertext))
	cipher.NewCBCDecrypter(block, []byte(contentKey[:aes.BlockSize])).CryptBlocks(plaintext, ciphertext)

	padding := int(plaintext[len(plaintext)-1])
	if padding == 0 || padding > aes.BlockSize || padding > len(plaintext) {
		return nil, fmt.Errorf("invalid PKCS#7 padding")
	}
	for _, value := range plaintext[len(plaintext)-padding:] {
		if int(value) != padding {
			return nil, fmt.Errorf("invalid PKCS#7 padding")
		}
	}

	return plaintext[:len(plaintext)-padding], nil
}

func keyReader(body []byte) (cct string, contentKey string, err error) {
	cctMatch := regexp.MustCompile(`(?m)\bvar\s+cct\s*=\s*'([^']*)'\s*;`).FindSubmatch(body)
	if len(cctMatch) != 2 || len(cctMatch[1]) == 0 {
		return "", "", fmt.Errorf("cct was not found in chapter HTML")
	}

	contentKeyMatch := regexp.MustCompile(`(?m)\bvar\s+contentKey\s*=\s*'([^']*)'\s*;`).FindSubmatch(body)
	if len(contentKeyMatch) != 2 || len(contentKeyMatch[1]) == 0 {
		return "", "", fmt.Errorf("contentKey was not found in chapter HTML")
	}

	return string(cctMatch[1]), string(contentKeyMatch[1]), nil
}
