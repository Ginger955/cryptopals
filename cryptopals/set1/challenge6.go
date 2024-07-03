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
	hamming int
}

func Challenge6() {
	filename := `/home/cristian/GolandProjects/cryptopals/cryptopals/set1/6.txt`
	content, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	const KEYSIZE = 40

	lines := bytes.Split(content, []byte(LINUX_NEWLINE))
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
				hamming: h / keysize,
			})
		}

		slices.SortFunc(hammings, func(a, b data) int {
			return cmp.Compare(a.hamming, b.hamming)
		})
		for i, hamming := range hammings {
			//TODO: take first n keys only
			if i > 5 {
				break
			}

			blocks := makeBlocks(line, hamming.keysize)
			transposed := transposeBlocks(blocks)
			for _, transposedBlock := range transposed {
				var histogram = make(map[rune]int)
				for b := 0; b < 256; b++ {
					xored := xorWithByteKey(transposedBlock, byte(b))
					runes := []rune(string(xored))
					for _, r := range runes {
						histogram[r]++
					}
				}

				//TODO: order histogram
				for k, v := range histogram {
					fmt.Printf("%c : %d\n", k, v)
				}
				fmt.Println()
			}
		}
	}
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
