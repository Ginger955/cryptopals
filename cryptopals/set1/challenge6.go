package set1

import (
	"bytes"
	"cmp"
	"encoding/base64"
	"fmt"
	"os"
	"slices"
)

const (
	WINDOWS_NEWLINE = "\r\n"
	LINUX_NEWLINE   = "\n"
)

type data struct {
	keysize int
	hamming float64
}

type scoreTracker struct {
	key      byte
	decoding []byte
	score    float32
}

type keyScore struct {
	score float32
	key   []byte
}

func Challenge6() {
	filename := `/home/cristian/GolandProjects/cryptopals/cryptopals/set1/6.txt`
	content, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	const KEYSIZE = 40

	lines := bytes.Split(content, []byte(LINUX_NEWLINE))
	for index, line := range lines {
		fmt.Printf("line #%d\n", index+1)
		var decoded = make([]byte, base64.StdEncoding.DecodedLen(len(line)))
		_, err := base64.StdEncoding.Decode(decoded, line)
		if err != nil {
			panic(err)
		}

		hammings := make([]data, 0)
		for keysize := 2; keysize <= KEYSIZE && keysize*2 <= len(decoded); keysize++ {
			first := decoded[:keysize]
			second := decoded[keysize : keysize*2]

			h, err := hammingDistance(first, second)
			if err != nil {
				panic(err)
			}

			hammings = append(hammings, data{
				keysize: keysize,
				hamming: float64(h) / float64(keysize),
			})
		}

		slices.SortFunc(hammings, func(a, b data) int {
			return cmp.Compare(a.hamming, b.hamming)
		})

		//only keep the top N performing key sizes that have the lowest hamming distance
		//var keyScores = make([]keyScore, 0)
		for i, data := range getTopN(hammings, 5) {
			fmt.Printf("key guess #%d, keysize: %d\n", i+1, data.keysize)
			blocks := makeBlocks(line, data.keysize)
			transposed := transposeBlocks(blocks)
			var guessedKey = make([]byte, 0)
			for _, transposedBlock := range transposed {
				var scores = make([]scoreTracker, 0)
				for b := 0; b < 256; b++ {
					xored := xorWithByteKey(transposedBlock, byte(b))
					englishScore := computeEnglishScore(xored)
					scores = append(scores, scoreTracker{
						key:      byte(b),
						score:    englishScore,
						decoding: xored,
					})
				}

				slices.SortFunc(scores, func(a, b scoreTracker) int {
					return cmp.Compare(a.score, b.score)
				})

				//bestN := 3
				//for i := bestN - 1; i >= 0; i-- {
				//	score := scores[len(scores)-1-i]
				//	fmt.Printf("#%d key: %c, score: %f, decoding: %s\n", i+1, score.key, score.score, string(score.decoding))
				//}
				guessedKey = append(guessedKey, scores[len(scores)-1].key)
			}

			//fmt.Printf("guessed key: %s\n", string(guessedKey))
			//keyScores = append(keyScores, keyScore{
			//	score: computeEnglishScore(guessedKey),
			//	key:   guessedKey,
			//})
		}

		//slices.SortFunc(keyScores, func(a, b keyScore) int {
		//	return cmp.Compare(a.score, b.score)
		//})
		//
		//fmt.Printf("best looking key: %s\n", string(keyScores[len(keyScores)-1].key))
	}
}

func getTopN(hammings []data, N int) []data {
	var (
		previous = hammings[0]
		count    = 0
		distinct = make([]data, 0, N)
	)

	distinct = append(distinct, previous)

	for _, hamming := range hammings {
		if previous.keysize != hamming.keysize && count < N {
			distinct = append(distinct, hamming)
			count++
		}
	}

	return distinct
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

	//blocks = dim(r x c0
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
