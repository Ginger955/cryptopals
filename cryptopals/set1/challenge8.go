package set1

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
)

func Challenge8() {
	const FILE = "cryptopals/set1/8.txt"

	content, err := os.ReadFile(FILE)
	if err != nil {
		panic(err)
	}

	split := bytes.Split(content, []byte("\n"))

	for i, line := range split {
		//to detect, take blocks the size of 16 or 32 bytes and check how frequent the same sequence repeats
		countFrequency16(line, i)
	}

}

func countFrequency16(input []byte, j int) {
	const BlockSize = 16
	var freq = make(map[string]uint)
	for i := 0; i <= len(input)-BlockSize; i += BlockSize {
		block := input[i : i+BlockSize]

		//the same byte slice will encode to the same hex string, so this is a good way to keep track of repeating byte sequences
		hex := hex.EncodeToString(block)
		freq[hex]++
	}

	var max uint64
	for _, v := range freq {
		if max < uint64(v) {
			max = uint64(v)
		}
	}

	if max > 1 {
		fmt.Printf("line #%d is likely AES ECB encrypted\n", j)
	}
}
