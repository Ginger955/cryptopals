package set1

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"slices"
)

func Challenge4() {
	minScore := frequencySum()
	filename := `cryptopals/set1/5.txt`
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("failed to read file: %s", err.Error())
	}

	scoreToKey := make(map[float32]byte)
	keyToXOR := make(map[byte][]byte)
	scores := []float32{}

	lines := bytes.Split(content, []byte("\r\n"))
	for _, line := range lines {
		decoded, err := decodeHex(line)
		if err != nil {
			log.Fatalf("failed to decode hex: %s", err.Error())
		}

		//try all ASCII characters as potential XOR key
		for i := 0; i < 256; i++ {
			xored := xorWithByteKey(decoded, byte(i))
			score := computeEnglishScore(xored)
			if score > minScore {
				keyToXOR[byte(i)] = xored
				scores = append(scores, score)
				scoreToKey[score] = byte(i)
			}
		}
	}

	if len(scores) > 0 {
		slices.Sort(scores)

		bestN := 15
		for i := bestN - 1; i >= 0; i-- {
			score := scores[len(scores)-1-i]
			key := scoreToKey[score]
			decoding := keyToXOR[key]
			fmt.Printf("#%d key: %s, score: %f, decoding: %s\n", i+1, string(key), score, string(decoding))
		}
	}
}
