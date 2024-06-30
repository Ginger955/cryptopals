package set1

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/bits"
)

var frequency = map[rune]float32{
	'e': 56.88,
	'E': 56.88,
	'a': 43.31,
	'A': 43.31,
	'r': 38.64,
	'R': 38.64,
	'i': 38.45,
	'I': 38.45,
	'o': 36.51,
	'O': 36.51,
	't': 35.54,
	'T': 35.54,
	'n': 33.92,
	'N': 33.92,
	's': 29.23,
	'S': 29.23,
	'l': 27.98,
	'L': 27.98,
	'c': 23.13,
	'C': 23.13,
}

func frequencySum() float32 {
	var total float32
	for _, v := range frequency {
		total += v
	}

	return total / 2
}

func hexToBase64(src []byte) ([]byte, error) {
	decodedHex, err := decodeHex(src)
	if err != nil {
		return nil, err
	}

	dst := make([]byte, base64.StdEncoding.EncodedLen(len(decodedHex)))
	base64.StdEncoding.Encode(dst, decodedHex)
	return dst, nil
}

func xor(a, b []byte) ([]byte, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("different input length")
	}

	var result = make([]byte, 0, len(a))
	for i, e := range a {
		//xor each element
		result = append(result, e^b[i])
	}

	return result, nil
}

func xorWithByteKey(input []byte, key byte) []byte {
	var result = make([]byte, 0, len(input))
	for _, a := range input {
		result = append(result, a^key)
	}

	return result
}

func xorRepeatingCycle(input, key []byte) []byte {
	var result = make([]byte, 0)

	for i, b := range input {
		result = append(result, b^key[i%len(key)])
	}

	return result
}

func computeEnglishCompatibilityScore(input []byte) float32 {
	var score float32
	runes := []rune(string(input))
	for _, r := range runes {
		score += frequency[r]
	}

	return score
}

func getMostFrequentRune(input []byte) rune {
	//this may have issues with non-roman letters
	freq := make(map[rune]uint32)
	runes := []rune(string(input))
	for _, r := range runes {
		freq[r]++
	}

	max := uint32(0)
	var most rune
	for r, count := range freq {
		if count > max {
			max = count
			most = r
		}
	}

	return most
}

func hammingDistance(a, b []byte) (int, error) {
	xored, err := xor(a, b)
	if err != nil {
		return 0, err
	}

	var distance = 0
	for _, x := range xored {
		distance += bits.OnesCount8(x)
	}

	return distance, nil
}

func decodeHex(input []byte) ([]byte, error) {
	var r = make([]byte, hex.DecodedLen(len(input)))
	_, err := hex.Decode(r, input)
	return r, err
}

func encodeHex(input []byte) []byte {
	var r = make([]byte, hex.EncodedLen(len(input)))
	hex.Encode(r, input)
	return r
}
