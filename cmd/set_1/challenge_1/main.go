package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
)

func main () {
	src := []byte("49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d")
	decodedBytes := make([]byte, hex.DecodedLen(len(src)))
	decodedLen, err := hex.Decode(decodedBytes, src)

	if err != nil {
		log.Fatal(err)
	}

	encodedBytes := make([]byte, base64.StdEncoding.EncodedLen(decodedLen))
	base64.StdEncoding.Encode(encodedBytes, decodedBytes[:decodedLen])

	fmt.Printf("%s\n", encodedBytes)
}