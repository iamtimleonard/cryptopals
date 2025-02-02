package main

import (
	"encoding/hex"
	"fmt"
	"os"
)

// fixed XOR
func main () {
	if len(os.Args) < 3 {
		fmt.Print("missing arguments")
		os.Exit(1)
	}
	// check buffers of equal length
	if len(os.Args[1]) != len(os.Args[2]) {
		fmt.Print("Buffers must be of equal length\n")
		os.Exit(1)
	}
	// decode both buffers
	input := make([]byte, hex.DecodedLen(len(os.Args[1])))
	inputLen, err := hex.Decode(input, []byte(os.Args[1]))

	if err != nil {
		fmt.Print("Error decoding input\n")
		os.Exit(1)
	}

	cipher := make([]byte, hex.DecodedLen(len(os.Args[2])))
	cipherLen, err := hex.Decode(cipher, []byte(os.Args[2]))

	if err != nil {
		fmt.Print("Error decoding cipher\n")
		os.Exit(1)
	}

	if inputLen != cipherLen {
		fmt.Print("Decoded lengths don't match :shrug:\n")
		os.Exit(1)
	}

	// XOR
	xor := make([]byte, inputLen)

	for i := 0; i < inputLen; i++ {
		xor[i] = input[i] ^ cipher[i]
	}

	encoded := make([]byte, hex.EncodedLen(len(xor)))
	encodedLen := hex.Encode(encoded, xor)

	if len(os.Args) == 4 {
		expected := os.Args[3]
		if expected == string(encoded) {
			fmt.Print("unit test passed\n")
		} else {
			fmt.Print("unit test failed\n")
		}
	} else {
		fmt.Printf("%s\n", encoded[:encodedLen])
	}
}