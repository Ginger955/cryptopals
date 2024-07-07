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
	filename := `cryptopals/set1/6.txt`
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
		interm := make([]data, 0)
		for i := 0; (i+2)*keysize < len(decoded); i++ {
			b1 := decoded[i*keysize : (i+1)*keysize]
			b2 := decoded[(i+1)*keysize : (i+2)*keysize]
			dst, err := hammingDistance(b1, b2)
			if err != nil {
				panic(err)
			}

			normalized := float32(dst) / float32(keysize)
			interm = append(interm, data{
				keysize:      keysize,
				hammingScore: normalized,
			})
		}

		total := float32(0)
		for _, inter := range interm {
			total += inter.hammingScore
		}

		normalizedDistance := total / float32(len(interm))

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
