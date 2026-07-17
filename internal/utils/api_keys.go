package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"

	"github.com/google/uuid"
)


func GenerateApiSecret() (plainTextSecret string, dbHash []byte, err error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", nil, errors.New("failed to generate secure bytes")
	}

	plainTextSecret = hex.EncodeToString(bytes)

	hash := sha256.Sum256(bytes)
	
	return plainTextSecret, hash[:], nil
}


func CompareSecret(incomingSecret string, dbHash []byte) bool {
	incomingHash := sha256.Sum256([]byte(incomingSecret))
	return subtle.ConstantTimeCompare(incomingHash[:], dbHash) == 1
}


func GenerateApiKey() (uuid.UUID, error) {
	return uuid.NewV7() 
}