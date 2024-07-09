package set2

func PKCS7Padding(ciphertext []byte, blockSize uint) []byte {
	if blockSize == 0 || uint(len(ciphertext))%blockSize == 0 {
		return ciphertext
	}

	padding := blockSize - uint(len(ciphertext))%blockSize
	if padding > 255 {
		return ciphertext
	}

	for range padding {
		ciphertext = append(ciphertext, byte(padding))
	}

	return ciphertext
}
