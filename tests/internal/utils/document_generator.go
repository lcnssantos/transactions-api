package utils

import "crypto/rand"

func GenerateDocument(size int) string {
	const digits = "0123456789"

	b := make([]byte, size)

	_, err := rand.Read(b)

	if err != nil {
		return ""
	}

	for i := range b {
		b[i] = digits[int(b[i])%10]
	}

	return string(b)
}
