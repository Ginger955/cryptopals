package set1

import (
	"crypto/aes"
	"encoding/base64"
	"fmt"
	"os"
)

func Challenge7() {
	const FILE = "cryptopals/set1/7.txt"
	const KEY = "YELLOW SUBMARINE"

	content, err := os.ReadFile(FILE)
	if err != nil {
		panic(err)
	}

	var decoded = make([]byte, base64.StdEncoding.DecodedLen(len(content)))
	_, err = base64.StdEncoding.Decode(decoded, content)
	if err != nil {
		panic(err)
	}

	aes128, err := aes.NewCipher([]byte(KEY))
	if err != nil {
		panic(err)
	}

	var decrypted = make([]byte, 0)

	for len(decoded) >= aes128.BlockSize() {
		var current = make([]byte, aes128.BlockSize())
		if len(decoded[aes128.BlockSize():]) == 0 {
			aes128.Decrypt(current, decoded[:aes128.BlockSize()])
		} else {
			aes128.Decrypt(current, decoded[aes128.BlockSize():])
		}
		decrypted = append(decrypted, current...)
		decoded = decoded[aes128.BlockSize():]
	}

	fmt.Println(string(decrypted))
}
