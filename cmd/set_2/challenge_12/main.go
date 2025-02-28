package main

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
)

// new version of oracle
// base64 decode then append to plaintext: "Um9sbGluJyBpbiBteSA1LjAKV2l0aCBteSByYWctdG9wIGRvd24gc28gbXkgaGFpciBjYW4gYmxvdwpUaGUgZ2lybGllcyBvbiBzdGFuZGJ5IHdhdmluZyBqdXN0IHRvIHNheSBoaQpEaWQgeW91IHN0b3A/IE5vLCBJIGp1c3QgZHJvdmUgYnkK"


func keyGen () []byte {
	res := make([]byte, aes.BlockSize)

	rand.Read(res)

	return res
}

func padPKCS7(data []byte, blockSize int) []byte {
	padLen := blockSize - (len(data) % blockSize)
	padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(data, padding...)
}

func ecbEncrypt(key, plaintext []byte) []byte {
	block, _ := aes.NewCipher(key)
	ciphertext := make([]byte, len(plaintext))
	for i := 0; i < len(plaintext); i += block.BlockSize() {
			block.Encrypt(ciphertext[i:], plaintext[i:])
	}
	return ciphertext
}

func encryptionOracle(key, input []byte) []byte {
	extraString := "Um9sbGluJyBpbiBteSA1LjAKV2l0aCBteSByYWctdG9wIGRvd24gc28gbXkgaGFpciBjYW4gYmxvdwpUaGUgZ2lybGllcyBvbiBzdGFuZGJ5IHdhdmluZyBqdXN0IHRvIHNheSBoaQpEaWQgeW91IHN0b3A/IE5vLCBJIGp1c3QgZHJvdmUgYnkK"
	
	appendedBytes := make([]byte, base64.StdEncoding.DecodedLen(len([]byte(extraString))))
	n, err := base64.StdEncoding.Decode(appendedBytes, []byte(extraString))
	if err != nil {
			log.Fatal(err)
	}
	appendedBytes = appendedBytes[:n]

	src := append(input, appendedBytes...)
	src = padPKCS7(src, aes.BlockSize)

	var res []byte

	res = ecbEncrypt(key, src)

	return res
}

func detectECB(key []byte) bool {
	blocks := make(map[string]bool)
	ciphertext := encryptionOracle(key, bytes.Repeat([]byte("A"), 48))
	for i := 0; i < len(ciphertext); i += aes.BlockSize {
			block := string(ciphertext[i : i+aes.BlockSize])
			if blocks[block] {
					fmt.Printf("probably ECB encryption\n")
					return true // ECB detected
			}
			blocks[block] = true
	}
	fmt.Printf("probably not ECB encryption\n")
	return false
}

func findBlockSize(key []byte) int {
	baseLen := len(encryptionOracle(key, []byte("")))



	for i := 1; i < 256; i++ {
		input := bytes.Repeat([]byte("A"), i)
		currentLen := len(encryptionOracle(key, input))

		if currentLen > baseLen {
			return currentLen - baseLen
		}
	}

	return 0
}

func decryptUnknownString(blockSize int, key []byte) []byte {
	var knownBytes []byte

	for {
			inputLength := blockSize - (len(knownBytes) % blockSize) - 1
			input := bytes.Repeat([]byte{'A'}, inputLength)

			// Get target ciphertext block
			ciphertext := encryptionOracle(key, input)
			targetBlockIndex := (inputLength + len(knownBytes)) / blockSize
			targetBlock := ciphertext[targetBlockIndex*blockSize : (targetBlockIndex+1)*blockSize]

			found := false
			// Brute-force the next byte
			for b := 0; b < 256; b++ {
					guess := make([]byte, 0, inputLength+len(knownBytes)+1)
					guess = append(guess, input...)
					guess = append(guess, knownBytes...)
					guess = append(guess, byte(b))

					// Get test ciphertext block
					testCiphertext := encryptionOracle(key, guess)
					testBlock := testCiphertext[targetBlockIndex*blockSize : (targetBlockIndex+1)*blockSize]

					if bytes.Equal(testBlock, targetBlock) {
							knownBytes = append(knownBytes, byte(b))
							found = true
							break
					}
			}

			// Check for padding termination (e.g., last byte is 0x01)
			if !found || (len(knownBytes) > 0 && knownBytes[len(knownBytes)-1] == 0x01) {
					break
			}
	}

	// Remove PKCS#7 padding
	return knownBytes
}


func main() {
	key := keyGen()

	// discover block size
	blockSize := findBlockSize(key)

	// detect ECB
	isEcb := detectECB(key)
	
	if !isEcb {
		log.Fatal("Probably not ECB. Exiting.")
	}

	byte := decryptUnknownString(blockSize, key)
	fmt.Printf("matched byte: %v\n", string(byte))
}