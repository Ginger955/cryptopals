package set2

func PKCS7Padding(ciphertext []byte, blockSize uint) []byte {
	padding := blockSize - uint(len(ciphertext))%blockSize
	for range padding {
		ciphertext = append(ciphertext, byte(padding))
	}

	return ciphertext
}
