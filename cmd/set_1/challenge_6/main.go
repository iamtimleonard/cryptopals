package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"math/bits"
	"net/http"
	"os"
	"sort"
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
	freqMap[' '] = 12
	freqMap['\''] = 1.0
	freqMap[':'] = 0.5
	freqMap[','] = 0.5
	freqMap['.'] = 0.5
	freqMap['!'] = 0.3

	var score float32
	for _, char := range text {
		r := rune(char)
		if r > unicode.MaxASCII || !unicode.IsPrint(r) {
				// Penalize heavily for non-printable/non-ASCII
				score -= 100
		} else {
			score += freqMap[rune(char)]
		}
	}
	return score
}

func hamming (a []byte, b []byte) (int, error) {
	if len(a) != len(b) {
		return 0, errors.New("Input values must be of equal length")
	}
	score := 0
	for i := 0; i < len(a); i++ {
		score += bits.OnesCount(uint(a[i] ^ b[i]))
	}
	return score, nil
}

func findKeysizes (decodedBytes []byte) []int {
	type key struct {
		size int
		score float32
	}

	var sizes []key

	for testsize := 2; testsize <= 40; testsize++ {
		// Take 4 blocks instead of 2
		blocks := decodedBytes[:testsize*4]
		totalDistance := 0
		pairs := 0

		// Compare each pair of blocks (0-1, 0-2, 0-3, 1-2, 1-3, 2-3)
		for i := 0; i < 4; i++ {
				for j := i + 1; j < 4; j++ {
						a := blocks[i*testsize : (i+1)*testsize]
						b := blocks[j*testsize : (j+1)*testsize]
						distance, err := hamming(a, b)
						if err == nil {
								totalDistance += distance
								pairs++
						}
				}
		}

		if pairs == 0 {
				continue
		}

		// Average the distances and normalize
		normalized := float32(totalDistance) / float32(pairs*testsize)
		sizes = append(sizes, key{score: normalized, size: testsize})
	}
	sort.Slice(sizes, func (i int, j int) bool {
		return sizes[i].score < sizes[j].score
	})

	bestSizes := make([]int, len(sizes))
	for i := 0; i < len(sizes); i++ {
		bestSizes[i] = sizes[i].size
	}
	return bestSizes[:1]
}

func splitIntoBlocks (keysize int, decodedBytes []byte) [][]byte {
	var blocks [][]byte
	for p := 0; p < len(decodedBytes); p += keysize {
		end := p + keysize
		if end > len(decodedBytes) - 1 {
			end = len(decodedBytes) - 1
		}
		blocks = append(blocks, decodedBytes[p:end])
	}
	return blocks
}

func transposeBlocks (blocks [][]byte) [][]byte {
	var transposed [][]byte
	for place := 0; place < len(blocks[0]); place++ {
		var newBlock []byte
		for _, block := range blocks {
			if place < len(block) {
				newBlock = append(newBlock, block[place])
			}
		}
		transposed = append(transposed, newBlock)
	}
	return transposed
}

// https://cryptopals.com/static/challenge-data/6.txt
func main () {
	type keyRes struct {
		keysize int
		hamming float32
	}
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

	// set and find keysize, lowest score
	bestSizes := findKeysizes(decodedBytes)
	fmt.Println("Top keysizes:", bestSizes)
	// for each keysize
	for _, keysize := range(bestSizes) {
		// break ciphertext into blocks of keysize
		// transpose (block of all first bytes, all second bytes, etc)
		transposed := transposeBlocks(splitIntoBlocks(keysize, decodedBytes))

		// for each block
		var keys []byte
		for _, block := range(transposed) {
			// solve single character XOR
			var bestScore float32
			var bestKey int
			for key := 0; key < 256; key++ {
				scored := score(xor(byte(key), block))
				if scored > bestScore {
					bestScore = scored
					bestKey = key
				}
			}
			keys = append(keys, byte(bestKey))
		}
		// put all single characters together
	if len(keys) != keysize {
		continue // Skip invalid keys
	}

		var res []byte
		for idx := range(decodedBytes) {
			keyIdx := idx % keysize
			res = append(res, decodedBytes[idx] ^ keys[keyIdx])
		}

		fmt.Printf("%s,\n%s\n", keys, res)
	}
}
