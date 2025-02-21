package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"os"
)

// deriveKey converts the input string into a fixed-length key using SHA-256.
func deriveKey(password string) []byte {
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}

// EncryptFile encrypts the given file using AES encryption and provides progress updates.
func EncryptFile(filePath string, password string, updateProgress func(int, string)) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	//aa
	fileInfo, err := file.Stat()
	if err != nil {
		return "", err
	}

	var processed int64
	blockSize := 1024 * 1024 // 1MB chunk

	// Initialize progress
	updateProgress(0, "encrypting")

	// Generate AES key from password
	keyBytes := deriveKey(password)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	totalSize := fileInfo.Size()
	// Create output file
	outputFile, err := os.Create(filePath + ".enc")
	if err != nil {
		return "", err
	}
	defer outputFile.Close()

	// Write nonce to the beginning of the output file
	if _, err := outputFile.Write(nonce); err != nil {
		updateProgress(0, "error")
		return "", err
	}

	// Encrypt file in chunks
	for {
		buf := make([]byte, blockSize)
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			updateProgress(0, "error")
			return "", err
		}
		if n == 0 {
			break
		}

		// AES encryption
		ciphertext := gcm.Seal(nil, nonce, buf[:n], nil)

		// Write encrypted data to output file
		if _, err := outputFile.Write(ciphertext); err != nil {
			updateProgress(0, "error")
			return "", err
		}

		processed += int64(n)
		percent := int(float64(processed) / float64(totalSize) * 100)
		updateProgress(percent, "encrypting")

	}

	updateProgress(100, "done")
	return filePath + ".enc", nil
}

// DecryptFile decrypts the given file using AES encryption and provides progress updates.
func DecryptFile(filePath string, password string, updateProgress func(int, string)) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return "", err
	}

	var processed int64
	blockSize := 1024 * 1024 // 1MB chunk

	// Initialize progress
	updateProgress(0, "decrypting")

	// Generate AES key from password
	keyBytes := deriveKey(password)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	// Read nonce from the beginning of the input file
	nonce := make([]byte, gcm.NonceSize())
	if _, err := file.Read(nonce); err != nil {
		return "", err
	}
	totalSize := fileInfo.Size() - int64(len(nonce))
	// Create output file
	outputFile, err := os.Create(filePath + ".dec")
	if err != nil {
		return "", err
	}
	defer outputFile.Close()

	for {
		buf := make([]byte, blockSize+gcm.Overhead())
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			updateProgress(0, "error")
			return "", err
		}
		if n == 0 {
			break
		}

		// AES decryption
		plaintext, err := gcm.Open(nil, nonce, buf[:n], nil)
		if err != nil {
			updateProgress(0, "error")
			return "", err
		}

		// Write decrypted data to output file
		if _, err := outputFile.Write(plaintext); err != nil {
			updateProgress(0, "error")
			return "", err
		}

		processed += int64(n)
		percent := int(float64(processed) / float64(totalSize) * 100)
		updateProgress(percent, "decrypting")
	}

	updateProgress(100, "done")
	return filePath + ".dec", nil
}
