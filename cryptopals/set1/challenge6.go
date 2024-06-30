package set1

import (
	"bytes"
	"cmp"
	"encoding/base64"
	"fmt"
	"os"
	"slices"
)

type data struct {
	keysize int
	hamming int
}

func Challenge6() {
	filename := `C:\Users\crist\GolandProjects\playground\cryptopals\set1\6.txt`
	content, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	const KEYSIZE = 40

	lines := bytes.Split(content, []byte("\r\n"))
	for _, line := range lines {
		var decoded = make([]byte, base64.StdEncoding.DecodedLen(len(line)))
		_, err := base64.StdEncoding.Decode(decoded, line)
		if err != nil {
			panic(err)
		}

		hammings := make([]data, 0)
		for keysize := 2; keysize <= KEYSIZE && keysize*2 < len(decoded); keysize++ {
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

		fmt.Println(hammings)

	}
}

func makeBlock(cipher []byte, keysize int) []byte {
	var result = make([]byte, 0, len(cipher)%keysize)
	for i := 0; i < len(cipher); i += keysize {
		result = append(result, cipher[keysize*i:keysize*(i+1)]...)
	}

	return result
}
