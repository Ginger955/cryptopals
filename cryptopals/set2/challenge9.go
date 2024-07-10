package set2

import (
	"fmt"
)

func Challenge9() {
	input := []byte("YELLOW SUBMARINE")
	out := PKCS7Padding(input, 32)
	fmt.Println(out)
}
