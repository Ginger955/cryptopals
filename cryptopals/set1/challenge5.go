package set1

import "fmt"

func Challenge5() {
	input := "Burning 'em, if you ain't quick and nimble\nI go crazy when I hear a cymbal"
	key := []byte("ICE")

	xored := xorRepeatingCycle([]byte(input), key)
	encoded := encodeHex(xored)
	fmt.Println(string(encoded))
}
