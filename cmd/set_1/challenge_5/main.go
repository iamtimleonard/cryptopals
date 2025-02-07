package main

import (
	"encoding/hex"
	"fmt"
)

// repeating key XOR
func main () {
	// set plaintext
	var src = `Burning 'em, if you ain't quick and nimble
I go crazy when I hear a cymbal`
	// set key
	var key = "ICE"
	// set key length
	var keyLen = len(key)
	// set res
	var res = make([]byte, len(src))
	// loop thru plaintext
	for idx, char := range src {
		// plaintext char XOR key (% key length) char
		res[idx] = byte(char) ^ byte(key[idx % keyLen])
	}
	// print
	expected := "0b3637272a2b2e63622c2e69692a23693a2a3c6324202d623d63343c2a26226324272765272a282b2f20430a652e2c652a3124333a653e2b2027630c692b20283165286326302e27282f"
	if hex.EncodeToString(res) != expected {
		fmt.Print("they don't match\n")
		fmt.Printf("%s\n", expected)
		fmt.Printf("%x\n", string(res))
	} else {
		fmt.Println("success")
		fmt.Printf("%x\n", string(res))
	}
}