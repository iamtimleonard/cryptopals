package main

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
)

func keyGen () []byte {
	res := make([]byte, aes.BlockSize)

	_, err := rand.Read(res)

	if err != nil {
		log.Fatal("something really went wrong making a key")
	}

	return res
}

func ecbEncrypt(key, plaintext []byte) []byte {
	block, _ := aes.NewCipher(key)
	ciphertext := make([]byte, len(plaintext))
	for i := 0; i < len(plaintext); i += block.BlockSize() {
			block.Encrypt(ciphertext[i:], plaintext[i:])
	}
	return ciphertext
}

func cbcEncrypt(key, src []byte) []byte {
	iv := make([]byte, aes.BlockSize)
	rand.Read(iv)
	cipher, _ := aes.NewCipher(key)
	res := make([]byte, len(src))
	prev := iv

	for i := 0; i < len(src); i += aes.BlockSize {
			// Create a temporary block to avoid modifying src
			tmp := make([]byte, aes.BlockSize)
			copy(tmp, src[i:i+aes.BlockSize])
			xor(tmp, prev)
			cipher.Encrypt(res[i:], tmp)
			prev = res[i:i+aes.BlockSize]
	}
	return res
}

// XOR two byte slices
func xor(dst []byte, src []byte) {
	for i := range dst {
			dst[i] ^= src[i]
	}
}

// Generates a secure random integer between [min, max]
func secureRandInt(min, max int) int {
	var buf [8]byte
	_, err := rand.Read(buf[:])
	if err != nil {
			panic(err)
	}
	// Convert bytes to uint64 and map to [min, max]
	return min + int(binary.BigEndian.Uint64(buf[:])%uint64(max-min+1))
}

func padPKCS7(data []byte, blockSize int) []byte {
	padLen := blockSize - (len(data) % blockSize)
	padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(data, padding...)
}

func chooseMode() int {
	var buf [1]byte
	_, err := rand.Read(buf[:])
	if err != nil {
			panic(err)
	}
	return int(buf[0]) % 2 // 0=ECB, 1=CBC
}

func encryptionOracle(input []byte) []byte {
	// encrypt with unknown key
	key := keyGen()
		// prepend 5-10 bytes and append 5-10 bytes to input
	prependLen := secureRandInt(5, 10)
	prependBytes := make([]byte, prependLen)
	rand.Read(prependBytes)
	appendLen := secureRandInt(5, 10)
	appendBytes := make([]byte, appendLen)
	rand.Read(appendBytes)
	src := append(prependBytes, input...)
	src = append(src, appendBytes...)
	src = padPKCS7(src, aes.BlockSize)

	strategy := chooseMode()
	var res []byte

	if strategy == 1 {
		fmt.Println("ecb encrypt")
		res = ecbEncrypt(key, src)
	} else {
		fmt.Println("cbc encrypt")
		res = cbcEncrypt(key, src)
	}

	return res
}

func detectECB(ciphertext []byte) bool {
	blocks := make(map[string]bool)
	for i := 0; i < len(ciphertext); i += aes.BlockSize {
			block := string(ciphertext[i : i+aes.BlockSize])
			if blocks[block] {
					fmt.Printf("probably ECB encryption")
					return true // ECB detected
			}
			blocks[block] = true
	}
	fmt.Printf("probably CBC encryption")
	return false
}

func main () {
	encryptedBytes := encryptionOracle(bytes.Repeat([]byte("A"), 48))

	detectECB(encryptedBytes)
}