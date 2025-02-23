package main

import (
	"crypto/aes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func decryptCBC(ciphertext []byte, key []byte, iv []byte) []byte {
    block, _ := aes.NewCipher(key)
    plaintext := make([]byte, len(ciphertext))

    prevBlock := iv
    for i := 0; i < len(ciphertext); i += aes.BlockSize {
        currentBlock := ciphertext[i : i+aes.BlockSize]
        block.Decrypt(plaintext[i:], currentBlock)

        xor(plaintext[i:i+aes.BlockSize], prevBlock)
        prevBlock = currentBlock 
    }

    return unpadPKCS7(plaintext)
}

// XOR two byte slices
func xor(dst []byte, src []byte) {
    for i := range dst {
        dst[i] ^= src[i]
    }
}

// Remove PKCS#7 padding
func unpadPKCS7(data []byte) []byte {
    padLength := int(data[len(data)-1])
    return data[:len(data)-padLength]
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("missing arguments")
	}

	url := os.Args[1]

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error getting from cryptopals: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("Error reading body: %v\n", err)
		os.Exit(1)
	}

	ciphertext, _ := base64.StdEncoding.DecodeString(string(data))

	key := []byte("YELLOW SUBMARINE")
	iv := make([]byte, aes.BlockSize) // All-zero IV (for this challenge)

	plaintext := decryptCBC(ciphertext, key, iv)
	fmt.Printf("%s\n", plaintext)
}
