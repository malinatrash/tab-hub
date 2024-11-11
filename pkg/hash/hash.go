package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

func Password(str string) (string, error) {
	if str == "" {
		return "", errors.New("password is empty")
	}

	hash := sha256.New()
	hash.Write([]byte(str))

	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	return hashString, nil
}
