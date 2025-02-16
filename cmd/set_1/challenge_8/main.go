package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main () {
	// get src and parse
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

	scanner := bufio.NewScanner(resp.Body)
	var line int
	var bestScore int
	var bestLine int

	for scanner.Scan() {
	// decode hex
		data := scanner.Text()
		decoded, err := hex.DecodeString(data)

		if err != nil {
			log.Fatalf("Could not parse source data")
		}

		blockSize := 16

		blockCounts := make(map[string]int)

		for i := 0; i <= len(decoded); i += blockSize {
			block := decoded[i:i+blockSize]
			blockCounts[string(block)]++
		}

		score := 0

		for _, count := range blockCounts {
			if count > 1 {
				score += count
			}
		}
		
		if score > bestScore {
			bestScore = score
			bestLine = line
		}

		line++
	}

	fmt.Printf("best line number: %v\n", bestLine)
}
