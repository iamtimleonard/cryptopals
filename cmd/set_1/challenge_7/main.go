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

// AES-128 in ECB mode
func main () {
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

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("Error reading body: %v\n", err)
		os.Exit(1)
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(string(body))

	if err != nil {
		fmt.Printf("Error decoding: %v\n", err)
	}

	key := []byte("YELLOW SUBMARINE")

	if len(key) % aes.BlockSize != 0 {
		log.Fatal("Incorrect key size")
	}

	cipher, err := aes.NewCipher(key)

	if err != nil {
		log.Fatal("error constructing cipher")
	}

	decrypted := make([]byte, len(decodedBytes))

	
	for i := 0; i < len(decodedBytes); i += aes.BlockSize {
		cipher.Decrypt(decrypted[i:i+aes.BlockSize], decodedBytes[i:i+aes.BlockSize])
	}
	fmt.Printf("%s\n", decrypted)
}
