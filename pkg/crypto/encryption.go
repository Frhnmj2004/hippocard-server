// Encryption/decryption for IPFS
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	//"encoding/base64"
	"fmt"
	"io"
	"log"
)

func Encrypt(data []byte, key []byte) ([]byte, error) {
	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Printf("Failed to create AES cipher: %v", err)
		return nil, err
	}

	// Use GCM mode (Galois/Counter Mode) for authenticated encryption
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Printf("Failed to create GCM: %v", err)
		return nil, err
	}

	// Generate a random nonce (unique per encryption)
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Printf("Failed to generate nonce: %v", err)
		return nil, err
	}

	// Encrypt the data (nonce is prepended to ciphertext)
	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	// Encode to base64 for easy storage/transmission
	//encoded := base64.StdEncoding.EncodeToString(ciphertext)
	return ciphertext, nil
}

func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	// Decode base64 string
	//ciphertext, err := base64.StdEncoding.DecodeString(encodedData)
	//if err != nil {
	//	log.Printf("Failed to decode base64 data: %v", err)
	//	return nil, err
	//}

	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Printf("Failed to create AES cipher: %v", err)
		return nil, err
	}

	// Use GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Printf("Failed to create GCM: %v", err)
		return nil, err
	}

	// Extract nonce and ciphertext
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt the data
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Printf("Failed to decrypt data: %v", err)
		return nil, err
	}

	return plaintext, nil
}

// GenerateKey generates a 32-byte AES key (for testing or initial setup)
func GenerateKey() ([]byte, error) {
	key := make([]byte, 32) // AES-256 requires a 32-byte key
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		log.Printf("Failed to generate key: %v", err)
		return nil, err
	}
	return key, nil
}
