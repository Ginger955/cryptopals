package set1

import (
	"bytes"
	"hash/maphash"
	"os"
)

func Challenge8() {
	const FILE = "cryptopals/set1/8.txt"
	const KEY = "YELLOW SUBMARINE"

	content, err := os.ReadFile(FILE)
	if err != nil {
		panic(err)
	}

	split := bytes.Split(content, []byte("\n"))

	for _, line := range split {
		decoded, err := decodeHex(line)
		if err != nil {
			panic(err)
		}

		//fmt.Println(string(decoded))
		//to crack, take a block the size of 16 or 32 bytes and check how frequent the same sequence repeats
		countFrequency16(decoded)

	}

}

func countFrequency16(input []byte) {
	var freq = make(map[uint64]int)
	hash := maphash.Hash{}
	for i := 0; i <= len(input)-16; i += 16 {
		block := input[i : i+16]

		_, err := hash.Write(block)
		if err != nil {
			panic(err)
		}

		v := hash.Sum64()
		freq[v]++

		hash.Reset()
	}

	for k, v := range freq {

	}
}
