package hashing

import (
	"crypto/sha256"
	"encoding/hex"
)

func EncodeStringToSHA256(value string) string {
	hash := sha256.Sum256([]byte(value))
	data := hex.EncodeToString(hash[:])
	return data
}
