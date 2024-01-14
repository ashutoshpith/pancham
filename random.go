package pancham

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRandomString(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(randomBytes)[:length], nil
}
