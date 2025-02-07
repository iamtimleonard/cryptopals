package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
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

// https://cryptopals.com/static/challenge-data/4.txt
func main () {
	// get first arg
	if len(os.Args) < 2 {
		log.Fatal("missing arguments")
	}

	url := os.Args[1]
	// fetch from cryptopals
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error getting from cryptopals: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Request failed: %d\n", resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)

	var bestScore float32
	var bestLine []byte
	var bestKey int

	// read each line
	for scanner.Scan() {
		line := scanner.Text()
		decodedBytes := make([]byte, hex.DecodedLen(len(line)))
		_, err := hex.Decode(decodedBytes, []byte(line))

		if err != nil {
			log.Fatal("couldn't parse line")
		}

		for key := 0; key < 256; key++ {
			parsed := xor(byte(key), decodedBytes)
			scored := score(parsed)
			if scored > bestScore {
				bestScore = scored
				bestLine = parsed
				bestKey = key
			}
		}
	}

	fmt.Printf("Result: %s,\n Key: %v\n", bestLine, bestKey)
}