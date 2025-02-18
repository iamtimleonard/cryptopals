package main

import (
	"fmt"
	"log"
	"os"
)

// Implement PKCS#7 padding
func main () {
	if len(os.Args) < 2 {
		log.Fatal("Missing args")
	}
	input := os.Args[1]
	inputBytes := []byte(input)
	blockSize := 20
	
	paddingNeeded := blockSize - (len(input) % blockSize)

	if paddingNeeded == 0 {
		paddingNeeded = blockSize
	}

	padding := make([]byte, paddingNeeded)

	for i := range padding {
		padding[i] = byte(paddingNeeded)
	}

	padded := append(inputBytes, padding...)

	fmt.Printf("%v\n", padded)
}