package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"unicode"
)

func xor (key byte, text []byte) []byte {
	res := make([]byte, len(text))
	for idx, char := range text {
		res[idx] = char ^ key
	}
	return res
}

func score (text []byte) float32 {
	freqMap := make(map[rune]float32)
	freqMap['e'] = 12.7
	freqMap['t'] = 9.1
	freqMap['a'] = 8.2
	freqMap['o'] = 7.5
	freqMap['i'] = 7.0
	freqMap['n'] = 6.7
	freqMap['s'] = 6.3
	freqMap['h'] = 6.1
	freqMap['r'] = 6.0
	freqMap['l'] = 4.0
	freqMap['c'] = 2.8
	freqMap['d'] = 4.3
	freqMap['u'] = 2.8
	freqMap['m'] = 2.4
	freqMap['w'] = 2.4
	freqMap['f'] = 2.2
	freqMap['g'] = 2.0
	freqMap['y'] = 2.0
	freqMap['p'] = 1.9
	freqMap['b'] = 1.5
	freqMap['v'] = 1.0
	freqMap['k'] = 0.8
	freqMap['j'] = 0.2
	freqMap['x'] = 0.2
	freqMap['q'] = 0.1
	freqMap['z'] = 0.1
	freqMap[' '] = 10

	var score float32
	for _, char := range text {
		score += freqMap[unicode.ToLower(rune(char))]
	}
	return score
}

// single byte XOR cipher
func main () {
	// decode hex string
	src := []byte("1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736")
	decodedBytes := make([]byte, hex.DecodedLen(len(src)))
	decodedLen, err := hex.Decode(decodedBytes, src)

	if err != nil {
		log.Fatal(err)
	}

	var bestScore float32
	var bestKey int
	var decoded []byte

	// loop thru letters of alphabet
	for key := 0; key < 256; key++ {
			// XOR
			parsed := xor(byte(key), decodedBytes[:decodedLen])
			// score using frequency analysis
			scored := score(parsed)
			// track best score and cipher
			if scored > bestScore {
				bestScore = scored
				bestKey = key
				decoded = parsed
			}
	}
	// print key and result
	fmt.Printf("Decoded: %s\nKey: %q\n", decoded, bestKey)
}