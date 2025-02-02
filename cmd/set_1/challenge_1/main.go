package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

// hex to base64
func main () {
	// get hex string from command line
	if len(os.Args) < 2 {
		fmt.Println("Error: missing argument")
		os.Exit(1)
	}

	// convert hex string to bytes
	src := []byte(os.Args[1])
	// create byte slice for decoded hex string
	decodedBytes := make([]byte, hex.DecodedLen(len(src)))
	// decode to byte slice
	decodedLen, err := hex.Decode(decodedBytes, src)

	if err != nil {
		log.Fatal(err)
	}

	// create byte slice for hex64 encoded string
	encodedBytes := make([]byte, base64.StdEncoding.EncodedLen(decodedLen))
	// encode string to byte slice
	base64.StdEncoding.Encode(encodedBytes, decodedBytes[:decodedLen])

	// check with user-provided expected outcome
	if len(os.Args) > 2 {
		expected := os.Args[2]
		if (expected == string(encodedBytes)) {
			fmt.Print("unit test passed\n")
		} else {
			fmt.Print("unit test failed\n")
		}
	} else {
		// otherwise just print the result
		fmt.Printf("%s\n", encodedBytes)
	}
}
