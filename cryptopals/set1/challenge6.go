package set1

import (
	"encoding/base64"
	"fmt"
	"os"
	"sort"
)

const (
	WINDOWS_NEWLINE = "\r\n"
	LINUX_NEWLINE   = "\n"
)

type data struct {
	keysize      int
	hammingScore float32
}

type keyOutputScoreTracker struct {
	key          byte
	englishScore float32
}

func Challenge6() {
	filename := `/home/cristian/GolandProjects/cryptopals/cryptopals/set1/6.txt`
	content, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var decoded = make([]byte, base64.StdEncoding.DecodedLen(len(content)))
	n, err := base64.StdEncoding.Decode(decoded, content)
	if err != nil {
		panic(err)
	}
	decoded = decoded[:n]

	const KEYSIZE = 40
	hammings := make([]data, 0)
	for keysize := 2; keysize <= KEYSIZE; keysize++ {
		normalizedDistance := averageHammingDistance(decoded, keysize)
		hammings = append(hammings, data{
			keysize:      keysize,
			hammingScore: normalizedDistance,
		})
	}

	sort.Slice(hammings, func(i, j int) bool {
		return hammings[i].hammingScore < hammings[j].hammingScore
	})

	blocks := makeBlocks(decoded, hammings[0].keysize)
	transposed := transposeBlocks(blocks)
	var guessedKey = make([]byte, 0)
	for _, transposedBlock := range transposed {
		var scores = make([]keyOutputScoreTracker, 0)
		for b := 0; b < 256; b++ {
			xored := xorWithByteKey(transposedBlock, byte(b))
			englishScore := computeEnglishScore(xored)
			scores = append(scores, keyOutputScoreTracker{
				key:          byte(b),
				englishScore: englishScore,
			})
		}

		sort.Slice(scores, func(i, j int) bool {
			return scores[i].englishScore > scores[j].englishScore
		})

		guessedKey = append(guessedKey, scores[0].key)
	}

	decrypted := xorRepeatingCycle(decoded, guessedKey)
	fmt.Printf("guessed output: %s\n", decrypted)
	fmt.Printf("guessed key: %s\n", guessedKey)
}

func averageHammingDistance(decoded []byte, keysize int) float32 {
	var distances []float64
	for i := 0; i+2*keysize <= len(decoded); i += keysize {
		d1 := decoded[i:keysize]
		d2 := decoded[i+keysize : i+(2*keysize)]
		distance, err := hammingDistance(d1, d2)
		if err != nil {
			panic(err)
		}
		normalized := float64(distance) / float64(keysize)
		distances = append(distances, normalized)
	}
	var total float64
	for _, d := range distances {
		total += d
	}
	return float32(total / float64(len(distances)))
}

func makeBlocks(cipher []byte, keysize int) [][]byte {
	var result = make([][]byte, 0, len(cipher)/keysize)
	for i := 0; keysize*(i+1) <= len(cipher); i++ {
		result = append(result, cipher[keysize*i:keysize*(i+1)])
	}

	return result
}

func transposeBlocks(blocks [][]byte) [][]byte {
	rows := len(blocks)
	cols := len(blocks[0])

	//blocks = dim(r x c)
	//transposed = dim(c x r)

	var transposedBlocks = make([][]byte, cols)
	for i := range transposedBlocks {
		transposedBlocks[i] = make([]byte, rows)
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			transposedBlocks[j][i] = blocks[i][j]
		}
	}

	return transposedBlocks
}
