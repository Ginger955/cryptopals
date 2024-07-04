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

func Challenge6() {
	filename := `C:\Users\crist\GolandProjects\playground\cryptopals\set1\6.txt`
	content, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	const KEYSIZE = 40

	lines := bytes.Split(content, []byte(WINDOWS_NEWLINE))
	for _, line := range lines {
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
		for i, data := range getTopN(hammings, 5) {
			fmt.Printf("dataset %d\n", i)
			blocks := makeBlocks(line, data.keysize)
			transposed := transposeBlocks(blocks)
			for _, transposedBlock := range transposed {
				var byteFrequency = make(map[byte]int)
				for b := 0; b < 256; b++ {
					xored := xorWithByteKey(transposedBlock, byte(b))
					for _, x := range xored {
						byteFrequency[x]++
					}
				}

				var highest int
				var order = make(map[int]byte)
				for k, v := range byteFrequency {
					//fmt.Printf("%c : %d\n", k, v)
					if v > highest {
						highest = v
						order[v] = k
					}
				}

				fmt.Printf("most frequent character: %c with %d occurences\n", order[highest], highest)
			}
		}
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

func mostLikelyKey(d []data) {

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
